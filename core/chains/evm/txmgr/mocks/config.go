// Code generated by mockery v2.43.2. DO NOT EDIT.

package mocks

import (
	chaintype "github.com/smartcontractkit/chainlink/v2/core/chains/evm/config/chaintype"
	mock "github.com/stretchr/testify/mock"
)

// Config is an autogenerated mock type for the ChainConfig type
type Config struct {
	mock.Mock
}

// ChainType provides a mock function with given fields:
func (_m *Config) ChainType() chaintype.ChainType {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for ChainType")
	}

	var r0 chaintype.ChainType
	if rf, ok := ret.Get(0).(func() chaintype.ChainType); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(chaintype.ChainType)
	}

	return r0
}

// FinalityDepth provides a mock function with given fields:
func (_m *Config) FinalityDepth() uint32 {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for FinalityDepth")
	}

	var r0 uint32
	if rf, ok := ret.Get(0).(func() uint32); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(uint32)
	}

	return r0
}

// FinalityTagEnabled provides a mock function with given fields:
func (_m *Config) FinalityTagEnabled() bool {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for FinalityTagEnabled")
	}

	var r0 bool
	if rf, ok := ret.Get(0).(func() bool); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// NonceAutoSync provides a mock function with given fields:
func (_m *Config) NonceAutoSync() bool {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for NonceAutoSync")
	}

	var r0 bool
	if rf, ok := ret.Get(0).(func() bool); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// RPCDefaultBatchSize provides a mock function with given fields:
func (_m *Config) RPCDefaultBatchSize() uint32 {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for RPCDefaultBatchSize")
	}

	var r0 uint32
	if rf, ok := ret.Get(0).(func() uint32); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(uint32)
	}

	return r0
}

// NewConfig creates a new instance of Config. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewConfig(t interface {
	mock.TestingT
	Cleanup(func())
}) *Config {
	mock := &Config{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
