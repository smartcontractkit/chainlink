// Code generated by mockery v2.43.2. DO NOT EDIT.

package mocks

import (
	url "net/url"

	mock "github.com/stretchr/testify/mock"
)

// TelemetryIngressEndpoint is an autogenerated mock type for the TelemetryIngressEndpoint type
type TelemetryIngressEndpoint struct {
	mock.Mock
}

// ChainID provides a mock function with given fields:
func (_m *TelemetryIngressEndpoint) ChainID() string {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for ChainID")
	}

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// Network provides a mock function with given fields:
func (_m *TelemetryIngressEndpoint) Network() string {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Network")
	}

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// ServerPubKey provides a mock function with given fields:
func (_m *TelemetryIngressEndpoint) ServerPubKey() string {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for ServerPubKey")
	}

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// URL provides a mock function with given fields:
func (_m *TelemetryIngressEndpoint) URL() *url.URL {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for URL")
	}

	var r0 *url.URL
	if rf, ok := ret.Get(0).(func() *url.URL); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*url.URL)
		}
	}

	return r0
}

// NewTelemetryIngressEndpoint creates a new instance of TelemetryIngressEndpoint. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewTelemetryIngressEndpoint(t interface {
	mock.TestingT
	Cleanup(func())
}) *TelemetryIngressEndpoint {
	mock := &TelemetryIngressEndpoint{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
