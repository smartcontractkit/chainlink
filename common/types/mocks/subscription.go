// Code generated by mockery v2.43.2. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// Subscription is an autogenerated mock type for the Subscription type
type Subscription struct {
	mock.Mock
}

// Err provides a mock function with given fields:
func (_m *Subscription) Err() <-chan error {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Err")
	}

	var r0 <-chan error
	if rf, ok := ret.Get(0).(func() <-chan error); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan error)
		}
	}

	return r0
}

// Unsubscribe provides a mock function with given fields:
func (_m *Subscription) Unsubscribe() {
	_m.Called()
}

// NewSubscription creates a new instance of Subscription. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewSubscription(t interface {
	mock.TestingT
	Cleanup(func())
}) *Subscription {
	mock := &Subscription{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
