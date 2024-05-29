package evcc

import (
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/util"
	"github.com/stretchr/testify/assert"
)

func (s *EVCCSuite) Test_Events() {
	payload := spineapi.EventPayload{
		Entity: s.mockRemoteEntity,
	}
	s.sut.HandleEvent(payload)

	payload.Entity = s.evEntity
	s.sut.HandleEvent(payload)

	payload.EventType = spineapi.EventTypeDeviceChange
	payload.ChangeType = spineapi.ElementChangeRemove
	s.sut.HandleEvent(payload)

	payload.EventType = spineapi.EventTypeEntityChange
	payload.ChangeType = spineapi.ElementChangeAdd
	s.sut.HandleEvent(payload)

	payload.EventType = spineapi.EventTypeEntityChange
	payload.ChangeType = spineapi.ElementChangeRemove
	s.sut.HandleEvent(payload)

	payload.EventType = spineapi.EventTypeDataChange
	payload.ChangeType = spineapi.ElementChangeAdd
	s.sut.HandleEvent(payload)

	payload.EventType = spineapi.EventTypeDataChange
	payload.ChangeType = spineapi.ElementChangeUpdate
	payload.Data = util.Ptr(model.DeviceConfigurationKeyValueDescriptionListDataType{})
	s.sut.HandleEvent(payload)

	payload.Data = util.Ptr(model.DeviceConfigurationKeyValueListDataType{})
	s.sut.HandleEvent(payload)

	var value model.DeviceDiagnosisOperatingStateType
	payload.Data = &value
	s.sut.HandleEvent(payload)

	payload.Data = util.Ptr(model.DeviceClassificationManufacturerDataType{})
	s.sut.HandleEvent(payload)

	payload.Data = util.Ptr(model.ElectricalConnectionParameterDescriptionListDataType{})
	s.sut.HandleEvent(payload)

	payload.Data = util.Ptr(model.ElectricalConnectionPermittedValueSetListDataType{})
	s.sut.HandleEvent(payload)

	payload.Data = util.Ptr(model.IdentificationListDataType{})
	s.sut.HandleEvent(payload)
}

func (s *EVCCSuite) Test_Failures() {
	payload := spineapi.EventPayload{
		Entity: s.mockRemoteEntity,
	}
	s.sut.evConnected(payload)

	s.sut.evConfigurationDescriptionDataUpdate(s.mockRemoteEntity)

	s.sut.evElectricalParamerDescriptionUpdate(s.mockRemoteEntity)
}

func (s *EVCCSuite) Test_evConfigurationDataUpdate() {
	payload := spineapi.EventPayload{
		Ski:    remoteSki,
		Device: s.remoteDevice,
		Entity: s.mockRemoteEntity,
	}
	s.sut.evConfigurationDataUpdate(payload)
	assert.False(s.T(), s.eventCalled)

	payload.Entity = s.evEntity
	s.sut.evConfigurationDataUpdate(payload)
	assert.False(s.T(), s.eventCalled)

	descData := &model.DeviceConfigurationKeyValueDescriptionListDataType{
		DeviceConfigurationKeyValueDescriptionData: []model.DeviceConfigurationKeyValueDescriptionDataType{
			{
				KeyId:   util.Ptr(model.DeviceConfigurationKeyIdType(1)),
				KeyName: util.Ptr(model.DeviceConfigurationKeyNameTypeCommunicationsStandard),
			},
			{
				KeyId:   util.Ptr(model.DeviceConfigurationKeyIdType(2)),
				KeyName: util.Ptr(model.DeviceConfigurationKeyNameTypeAsymmetricChargingSupported),
			},
		},
	}

	rFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.evEntity, model.FeatureTypeTypeDeviceConfiguration, model.RoleTypeServer)
	fErr := rFeature.UpdateData(model.FunctionTypeDeviceConfigurationKeyValueDescriptionListData, descData, nil, nil)
	assert.Nil(s.T(), fErr)

	s.sut.evConfigurationDataUpdate(payload)
	assert.False(s.T(), s.eventCalled)

	data := &model.DeviceConfigurationKeyValueListDataType{
		DeviceConfigurationKeyValueData: []model.DeviceConfigurationKeyValueDataType{},
	}

	payload.Data = data

	s.sut.evConfigurationDataUpdate(payload)
	assert.False(s.T(), s.eventCalled)

	data = &model.DeviceConfigurationKeyValueListDataType{
		DeviceConfigurationKeyValueData: []model.DeviceConfigurationKeyValueDataType{
			{
				KeyId: util.Ptr(model.DeviceConfigurationKeyIdType(0)),
				Value: util.Ptr(model.DeviceConfigurationKeyValueValueType{}),
			},
			{
				KeyId: util.Ptr(model.DeviceConfigurationKeyIdType(1)),
				Value: util.Ptr(model.DeviceConfigurationKeyValueValueType{
					String: util.Ptr(model.DeviceConfigurationKeyValueStringTypeISO151182ED2),
				}),
			},
			{
				KeyId: util.Ptr(model.DeviceConfigurationKeyIdType(2)),
				Value: util.Ptr(model.DeviceConfigurationKeyValueValueType{
					Boolean: util.Ptr(false),
				}),
			},
		},
	}

	payload.Data = data

	s.sut.evConfigurationDataUpdate(payload)
	assert.True(s.T(), s.eventCalled)
}

func (s *EVCCSuite) Test_evOperatingStateDataUpdate() {
	payload := spineapi.EventPayload{
		Ski:    remoteSki,
		Device: s.remoteDevice,
		Entity: s.mockRemoteEntity,
	}
	s.sut.evOperatingStateDataUpdate(payload)
	assert.False(s.T(), s.eventCalled)

	payload.Entity = s.evEntity
	s.sut.evOperatingStateDataUpdate(payload)
	assert.False(s.T(), s.eventCalled)

	data := &model.DeviceDiagnosisStateDataType{
		OperatingState: util.Ptr(model.DeviceDiagnosisOperatingStateTypeNormalOperation),
	}

	rFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.evEntity, model.FeatureTypeTypeDeviceDiagnosis, model.RoleTypeServer)
	fErr := rFeature.UpdateData(model.FunctionTypeDeviceDiagnosisStateData, data, nil, nil)
	assert.Nil(s.T(), fErr)

	s.sut.evOperatingStateDataUpdate(payload)
	assert.True(s.T(), s.eventCalled)
}

func (s *EVCCSuite) Test_evIdentificationDataUpdate() {
	payload := spineapi.EventPayload{
		Ski:    remoteSki,
		Device: s.remoteDevice,
		Entity: s.mockRemoteEntity,
	}
	s.sut.evIdentificationDataUpdate(payload)
	assert.False(s.T(), s.eventCalled)

	payload.Entity = s.evEntity
	s.sut.evIdentificationDataUpdate(payload)
	assert.False(s.T(), s.eventCalled)

	data := &model.IdentificationListDataType{
		IdentificationData: []model.IdentificationDataType{
			{
				IdentificationId:   util.Ptr(model.IdentificationIdType(0)),
				IdentificationType: util.Ptr(model.IdentificationTypeTypeEui48),
			},
			{
				IdentificationId:    util.Ptr(model.IdentificationIdType(1)),
				IdentificationType:  util.Ptr(model.IdentificationTypeTypeEui48),
				IdentificationValue: util.Ptr(model.IdentificationValueType("test")),
			},
		},
	}

	payload.Data = data
	s.sut.evIdentificationDataUpdate(payload)
	assert.True(s.T(), s.eventCalled)
}

func (s *EVCCSuite) Test_evManufacturerDataUpdate() {
	payload := spineapi.EventPayload{
		Ski:    remoteSki,
		Device: s.remoteDevice,
		Entity: s.mockRemoteEntity,
	}
	s.sut.evManufacturerDataUpdate(payload)
	assert.False(s.T(), s.eventCalled)

	payload.Entity = s.evEntity
	s.sut.evManufacturerDataUpdate(payload)
	assert.False(s.T(), s.eventCalled)

	data := &model.DeviceClassificationManufacturerDataType{
		BrandName: util.Ptr(model.DeviceClassificationStringType("test")),
	}

	rFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.evEntity, model.FeatureTypeTypeDeviceClassification, model.RoleTypeServer)
	fErr := rFeature.UpdateData(model.FunctionTypeDeviceClassificationManufacturerData, data, nil, nil)
	assert.Nil(s.T(), fErr)

	s.sut.evManufacturerDataUpdate(payload)
	assert.True(s.T(), s.eventCalled)
}

func (s *EVCCSuite) Test_evElectricalPermittedValuesUpdate() {
	payload := spineapi.EventPayload{
		Ski:    remoteSki,
		Device: s.remoteDevice,
		Entity: s.mockRemoteEntity,
	}
	s.sut.evElectricalPermittedValuesUpdate(payload)
	assert.False(s.T(), s.eventCalled)

	payload.Entity = s.evEntity
	s.sut.evElectricalPermittedValuesUpdate(payload)
	assert.False(s.T(), s.eventCalled)

	paramData := &model.ElectricalConnectionParameterDescriptionListDataType{
		ElectricalConnectionParameterDescriptionData: []model.ElectricalConnectionParameterDescriptionDataType{
			{
				ElectricalConnectionId: util.Ptr(model.ElectricalConnectionIdType(0)),
				ParameterId:            util.Ptr(model.ElectricalConnectionParameterIdType(0)),
				ScopeType:              util.Ptr(model.ScopeTypeTypeACPowerTotal),
			},
		},
	}

	rFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.evEntity, model.FeatureTypeTypeElectricalConnection, model.RoleTypeServer)
	fErr := rFeature.UpdateData(model.FunctionTypeElectricalConnectionParameterDescriptionListData, paramData, nil, nil)
	assert.Nil(s.T(), fErr)

	s.sut.evElectricalPermittedValuesUpdate(payload)
	assert.False(s.T(), s.eventCalled)

	permData := &model.ElectricalConnectionPermittedValueSetListDataType{
		ElectricalConnectionPermittedValueSetData: []model.ElectricalConnectionPermittedValueSetDataType{
			{
				ElectricalConnectionId: util.Ptr(model.ElectricalConnectionIdType(0)),
				ParameterId:            util.Ptr(model.ElectricalConnectionParameterIdType(0)),
				PermittedValueSet: []model.ScaledNumberSetType{
					{
						Value: []model.ScaledNumberType{*model.NewScaledNumberType(0)},
						Range: []model.ScaledNumberRangeType{
							{
								Min: model.NewScaledNumberType(6),
								Max: model.NewScaledNumberType(16),
							},
						},
					},
				},
			},
		},
	}

	payload.Data = permData
	s.sut.evElectricalPermittedValuesUpdate(payload)
	assert.True(s.T(), s.eventCalled)
}
