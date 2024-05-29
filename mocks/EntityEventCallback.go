// Code generated by mockery v2.42.1. DO NOT EDIT.

package mocks

import (
	eebus_goapi "github.com/enbility/eebus-go/api"
	api "github.com/enbility/spine-go/api"

	mock "github.com/stretchr/testify/mock"
)

// EntityEventCallback is an autogenerated mock type for the EntityEventCallback type
type EntityEventCallback struct {
	mock.Mock
}

type EntityEventCallback_Expecter struct {
	mock *mock.Mock
}

func (_m *EntityEventCallback) EXPECT() *EntityEventCallback_Expecter {
	return &EntityEventCallback_Expecter{mock: &_m.Mock}
}

// Execute provides a mock function with given fields: ski, device, entity, event
func (_m *EntityEventCallback) Execute(ski string, device api.DeviceRemoteInterface, entity api.EntityRemoteInterface, event eebus_goapi.EventType) {
	_m.Called(ski, device, entity, event)
}

// EntityEventCallback_Execute_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Execute'
type EntityEventCallback_Execute_Call struct {
	*mock.Call
}

// Execute is a helper method to define mock.On call
//   - ski string
//   - device api.DeviceRemoteInterface
//   - entity api.EntityRemoteInterface
//   - event eebus_goapi.EventType
func (_e *EntityEventCallback_Expecter) Execute(ski interface{}, device interface{}, entity interface{}, event interface{}) *EntityEventCallback_Execute_Call {
	return &EntityEventCallback_Execute_Call{Call: _e.mock.On("Execute", ski, device, entity, event)}
}

func (_c *EntityEventCallback_Execute_Call) Run(run func(ski string, device api.DeviceRemoteInterface, entity api.EntityRemoteInterface, event eebus_goapi.EventType)) *EntityEventCallback_Execute_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].(api.DeviceRemoteInterface), args[2].(api.EntityRemoteInterface), args[3].(eebus_goapi.EventType))
	})
	return _c
}

func (_c *EntityEventCallback_Execute_Call) Return() *EntityEventCallback_Execute_Call {
	_c.Call.Return()
	return _c
}

func (_c *EntityEventCallback_Execute_Call) RunAndReturn(run func(string, api.DeviceRemoteInterface, api.EntityRemoteInterface, eebus_goapi.EventType)) *EntityEventCallback_Execute_Call {
	_c.Call.Return(run)
	return _c
}

// NewEntityEventCallback creates a new instance of EntityEventCallback. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewEntityEventCallback(t interface {
	mock.TestingT
	Cleanup(func())
}) *EntityEventCallback {
	mock := &EntityEventCallback{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}