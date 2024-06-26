// Code generated by mockery v2.43.2. DO NOT EDIT.

package mocks

import (
	context "context"

	api "github.com/smartcontractkit/chainlink/v2/core/services/gateway/api"

	handlers "github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers"

	mock "github.com/stretchr/testify/mock"
)

// Handler is an autogenerated mock type for the Handler type
type Handler struct {
	mock.Mock
}

// Close provides a mock function with given fields:
func (_m *Handler) Close() error {
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

// HandleNodeMessage provides a mock function with given fields: ctx, msg, nodeAddr
func (_m *Handler) HandleNodeMessage(ctx context.Context, msg *api.Message, nodeAddr string) error {
	ret := _m.Called(ctx, msg, nodeAddr)

	if len(ret) == 0 {
		panic("no return value specified for HandleNodeMessage")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *api.Message, string) error); ok {
		r0 = rf(ctx, msg, nodeAddr)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// HandleUserMessage provides a mock function with given fields: ctx, msg, callbackCh
func (_m *Handler) HandleUserMessage(ctx context.Context, msg *api.Message, callbackCh chan<- handlers.UserCallbackPayload) error {
	ret := _m.Called(ctx, msg, callbackCh)

	if len(ret) == 0 {
		panic("no return value specified for HandleUserMessage")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *api.Message, chan<- handlers.UserCallbackPayload) error); ok {
		r0 = rf(ctx, msg, callbackCh)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Start provides a mock function with given fields: _a0
func (_m *Handler) Start(_a0 context.Context) error {
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

// NewHandler creates a new instance of Handler. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewHandler(t interface {
	mock.TestingT
	Cleanup(func())
}) *Handler {
	mock := &Handler{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
