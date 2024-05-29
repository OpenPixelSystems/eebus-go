package lpc

import (
	"errors"
	"time"

	"github.com/enbility/eebus-go/api"
	"github.com/enbility/eebus-go/features/server"
	ucapi "github.com/enbility/eebus-go/usecases/api"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/util"
)

// Scenario 1

// return the current loadcontrol limit data
//
// return values:
//   - limit: load limit data
//
// possible errors:
//   - ErrDataNotAvailable if no such limit is (yet) available
//   - and others
func (e *CsLPC) ConsumptionLimit() (limit ucapi.LoadLimit, resultErr error) {
	limit = ucapi.LoadLimit{
		Value:        0.0,
		IsChangeable: false,
		IsActive:     false,
		Duration:     0,
	}
	resultErr = api.ErrDataNotAvailable

	lc, limidId, err := e.loadControlServerAndLimitId()
	if err != nil {
		return limit, err
	}

	value, err := lc.GetLimitDataForId(limidId)
	if err != nil || value == nil || value.LimitId == nil || value.Value == nil {
		return
	}

	limit.Value = value.Value.GetValue()
	limit.IsChangeable = (value.IsLimitChangeable != nil && *value.IsLimitChangeable)
	limit.IsActive = (value.IsLimitActive != nil && *value.IsLimitActive)
	if value.TimePeriod != nil && value.TimePeriod.EndTime != nil {
		if duration, err := value.TimePeriod.EndTime.GetTimeDuration(); err == nil {
			limit.Duration = duration
		}
	}

	return limit, nil
}

// set the current loadcontrol limit data
func (e *CsLPC) SetConsumptionLimit(limit ucapi.LoadLimit) (resultErr error) {
	loadControlf, limidId, err := e.loadControlServerAndLimitId()
	if err != nil {
		return err
	}

	limitData := model.LoadControlLimitDataType{
		LimitId:           util.Ptr(limidId),
		IsLimitChangeable: util.Ptr(limit.IsChangeable),
		IsLimitActive:     util.Ptr(limit.IsActive),
		Value:             model.NewScaledNumberType(limit.Value),
	}
	if limit.Duration > 0 {
		limitData.TimePeriod = &model.TimePeriodType{
			EndTime: model.NewAbsoluteOrRelativeTimeTypeFromDuration(limit.Duration),
		}
	}

	deleteTimePeriod := &model.LoadControlLimitDataElementsType{
		TimePeriod: util.Ptr(model.TimePeriodElementsType{}),
	}

	return loadControlf.UpdateLimitDataForId(limitData, deleteTimePeriod, limidId)
}

// return the currently pending incoming consumption write limits
func (e *CsLPC) PendingConsumptionLimits() map[model.MsgCounterType]ucapi.LoadLimit {
	result := make(map[model.MsgCounterType]ucapi.LoadLimit)

	_, limitId, err := e.loadControlServerAndLimitId()
	if err != nil {
		return result
	}

	e.pendingMux.Lock()
	defer e.pendingMux.Unlock()

	for key, msg := range e.pendingLimits {
		data := msg.Cmd.LoadControlLimitListData

		// elements are only added to the map if all required fields exist
		// therefor not check for these are needed here

		// find the item which contains the limit for this usecase
		for _, item := range data.LoadControlLimitData {
			if item.LimitId == nil ||
				limitId != *item.LimitId {
				continue
			}

			limit := ucapi.LoadLimit{}

			if item.TimePeriod != nil {
				if duration, err := item.TimePeriod.GetDuration(); err == nil {
					limit.Duration = duration
				}
			}

			if item.IsLimitActive != nil {
				limit.IsActive = *item.IsLimitActive
			}

			if item.Value != nil {
				limit.Value = item.Value.GetValue()
			}

			result[key] = limit
		}
	}

	return result
}

// accept or deny an incoming consumption write limit
//
// use PendingConsumptionLimits to get the list of currently pending requests
func (e *CsLPC) ApproveOrDenyConsumptionLimit(msgCounter model.MsgCounterType, approve bool, reason string) {
	e.pendingMux.Lock()
	defer e.pendingMux.Unlock()

	msg, ok := e.pendingLimits[msgCounter]
	if !ok {
		return
	}

	f := e.LocalEntity.FeatureOfTypeAndRole(model.FeatureTypeTypeLoadControl, model.RoleTypeServer)

	result := model.ErrorType{
		ErrorNumber: model.ErrorNumberType(0),
	}
	if !approve {
		result.ErrorNumber = model.ErrorNumberType(7)
		result.Description = util.Ptr(model.DescriptionType(reason))
	}
	f.ApproveOrDenyWrite(msg, result)
}

// Scenario 2

// return Failsafe limit for the consumed active (real) power of the
// Controllable System. This limit becomes activated in "init" state or "failsafe state".
func (e *CsLPC) FailsafeConsumptionActivePowerLimit() (limit float64, isChangeable bool, resultErr error) {
	limit = 0
	isChangeable = false
	resultErr = api.ErrDataNotAvailable

	dc, err := server.NewDeviceConfiguration(e.LocalEntity)
	if err != nil {
		return
	}

	filter := model.DeviceConfigurationKeyValueDescriptionDataType{
		KeyName: util.Ptr(model.DeviceConfigurationKeyNameTypeFailsafeConsumptionActivePowerLimit),
	}
	keyData, err := dc.GetKeyValueDataForFilter(filter)
	if err != nil || keyData == nil || keyData.KeyId == nil || keyData.Value == nil || keyData.Value.ScaledNumber == nil {
		return
	}

	limit = keyData.Value.ScaledNumber.GetValue()
	isChangeable = (keyData.IsValueChangeable != nil && *keyData.IsValueChangeable)
	resultErr = nil
	return
}

// set Failsafe limit for the consumed active (real) power of the
// Controllable System. This limit becomes activated in "init" state or "failsafe state".
func (e *CsLPC) SetFailsafeConsumptionActivePowerLimit(value float64, changeable bool) error {
	keyName := model.DeviceConfigurationKeyNameTypeFailsafeConsumptionActivePowerLimit
	keyValue := model.DeviceConfigurationKeyValueValueType{
		ScaledNumber: model.NewScaledNumberType(value),
	}

	dc, err := server.NewDeviceConfiguration(e.LocalEntity)
	if err != nil {
		return err
	}

	data := model.DeviceConfigurationKeyValueDataType{
		Value:             &keyValue,
		IsValueChangeable: &changeable,
	}
	filter := model.DeviceConfigurationKeyValueDescriptionDataType{
		KeyName: util.Ptr(keyName),
	}
	return dc.UpdateKeyValueDataForFilter(data, nil, filter)
}

// return minimum time the Controllable System remains in "failsafe state" unless conditions
// specified in this Use Case permit leaving the "failsafe state"
func (e *CsLPC) FailsafeDurationMinimum() (duration time.Duration, isChangeable bool, resultErr error) {
	duration = 0
	isChangeable = false
	resultErr = api.ErrDataNotAvailable

	dc, err := server.NewDeviceConfiguration(e.LocalEntity)
	if err != nil {
		return
	}

	filter := model.DeviceConfigurationKeyValueDescriptionDataType{
		KeyName: util.Ptr(model.DeviceConfigurationKeyNameTypeFailsafeDurationMinimum),
	}
	keyData, err := dc.GetKeyValueDataForFilter(filter)
	if err != nil || keyData == nil || keyData.KeyId == nil || keyData.Value == nil || keyData.Value.Duration == nil {
		return
	}

	durationValue, err := keyData.Value.Duration.GetTimeDuration()
	if err != nil {
		return
	}

	duration = durationValue
	isChangeable = (keyData.IsValueChangeable != nil && *keyData.IsValueChangeable)
	resultErr = nil
	return
}

// set minimum time the Controllable System remains in "failsafe state" unless conditions
// specified in this Use Case permit leaving the "failsafe state"
//
// parameters:
//   - duration: has to be >= 2h and <= 24h
//   - changeable: boolean if the client service can change this value
func (e *CsLPC) SetFailsafeDurationMinimum(duration time.Duration, changeable bool) error {
	if duration < time.Duration(time.Hour*2) || duration > time.Duration(time.Hour*24) {
		return errors.New("duration outside of allowed range")
	}
	keyName := model.DeviceConfigurationKeyNameTypeFailsafeDurationMinimum
	keyValue := model.DeviceConfigurationKeyValueValueType{
		Duration: model.NewDurationType(duration),
	}

	dc, err := server.NewDeviceConfiguration(e.LocalEntity)
	if err != nil {
		return err
	}

	data := model.DeviceConfigurationKeyValueDataType{
		Value:             &keyValue,
		IsValueChangeable: &changeable,
	}
	filter := model.DeviceConfigurationKeyValueDescriptionDataType{
		KeyName: util.Ptr(keyName),
	}
	return dc.UpdateKeyValueDataForFilter(data, nil, filter)
}

// Scenario 3

func (e *CsLPC) IsHeartbeatWithinDuration() bool {
	if e.heartbeatDiag == nil {
		return false
	}

	return e.heartbeatDiag.IsHeartbeatWithinDuration(2 * time.Minute)
}

// Scenario 4

// return nominal maximum active (real) power the Controllable System is
// allowed to consume due to the customer's contract.
func (e *CsLPC) ContractualConsumptionNominalMax() (value float64, resultErr error) {
	value = 0
	resultErr = api.ErrDataNotAvailable

	ec, err := server.NewElectricalConnection(e.LocalEntity)
	if err != nil {
		resultErr = err
		return
	}

	filter := model.ElectricalConnectionCharacteristicDataType{
		ElectricalConnectionId: util.Ptr(model.ElectricalConnectionIdType(0)),
		ParameterId:            util.Ptr(model.ElectricalConnectionParameterIdType(0)),
		CharacteristicContext:  util.Ptr(model.ElectricalConnectionCharacteristicContextTypeEntity),
		CharacteristicType:     util.Ptr(model.ElectricalConnectionCharacteristicTypeTypeContractualConsumptionNominalMax),
	}
	charData, err := ec.GetCharacteristicsForFilter(filter)
	if err != nil || len(charData) == 0 ||
		charData[0].CharacteristicId == nil ||
		charData[0].Value == nil {
		return
	}

	return charData[0].Value.GetValue(), nil
}

// set nominal maximum active (real) power the Controllable System is
// allowed to consume due to the customer's contract.
func (e *CsLPC) SetContractualConsumptionNominalMax(value float64) error {
	ec, err := server.NewElectricalConnection(e.LocalEntity)
	if err != nil {
		return err
	}

	electricalConnectionid := util.Ptr(model.ElectricalConnectionIdType(0))
	parameterId := util.Ptr(model.ElectricalConnectionParameterIdType(0))
	charList, err := ec.GetCharacteristicsForFilter(model.ElectricalConnectionCharacteristicDataType{
		ElectricalConnectionId: electricalConnectionid,
		ParameterId:            parameterId,
		CharacteristicContext:  util.Ptr(model.ElectricalConnectionCharacteristicContextTypeEntity),
		CharacteristicType:     util.Ptr(model.ElectricalConnectionCharacteristicTypeTypeContractualConsumptionNominalMax),
	})
	if err != nil || len(charList) == 0 {
		return api.ErrDataNotAvailable
	}

	data := model.ElectricalConnectionCharacteristicDataType{
		ElectricalConnectionId: electricalConnectionid,
		ParameterId:            parameterId,
		CharacteristicId:       charList[0].CharacteristicId,
		Value:                  model.NewScaledNumberType(value),
	}
	return ec.UpdateCharacteristic(data, nil)
}
