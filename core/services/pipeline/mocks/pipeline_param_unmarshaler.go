// Code generated by mockery v2.43.2. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// PipelineParamUnmarshaler is an autogenerated mock type for the PipelineParamUnmarshaler type
type PipelineParamUnmarshaler struct {
	mock.Mock
}

// UnmarshalPipelineParam provides a mock function with given fields: val
func (_m *PipelineParamUnmarshaler) UnmarshalPipelineParam(val interface{}) error {
	ret := _m.Called(val)

	if len(ret) == 0 {
		panic("no return value specified for UnmarshalPipelineParam")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(interface{}) error); ok {
		r0 = rf(val)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewPipelineParamUnmarshaler creates a new instance of PipelineParamUnmarshaler. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewPipelineParamUnmarshaler(t interface {
	mock.TestingT
	Cleanup(func())
}) *PipelineParamUnmarshaler {
	mock := &PipelineParamUnmarshaler{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
