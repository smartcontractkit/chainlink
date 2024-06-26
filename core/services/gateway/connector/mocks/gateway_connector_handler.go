// Code generated by mockery v2.43.2. DO NOT EDIT.

package mocks

import (
	api "github.com/smartcontractkit/chainlink/v2/core/services/gateway/api"

	context "context"

	mock "github.com/stretchr/testify/mock"
)

// GatewayConnectorHandler is an autogenerated mock type for the GatewayConnectorHandler type
type GatewayConnectorHandler struct {
	mock.Mock
}

// Close provides a mock function with given fields:
func (_m *GatewayConnectorHandler) Close() error {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Close")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// HandleGatewayMessage provides a mock function with given fields: ctx, gatewayId, msg
func (_m *GatewayConnectorHandler) HandleGatewayMessage(ctx context.Context, gatewayId string, msg *api.Message) {
	_m.Called(ctx, gatewayId, msg)
}

// Start provides a mock function with given fields: _a0
func (_m *GatewayConnectorHandler) Start(_a0 context.Context) error {
	ret := _m.Called(_a0)

	if len(ret) == 0 {
		panic("no return value specified for Start")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewGatewayConnectorHandler creates a new instance of GatewayConnectorHandler. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewGatewayConnectorHandler(t interface {
	mock.TestingT
	Cleanup(func())
}) *GatewayConnectorHandler {
	mock := &GatewayConnectorHandler{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
