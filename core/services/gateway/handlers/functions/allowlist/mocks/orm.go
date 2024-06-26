// Code generated by mockery v2.43.2. DO NOT EDIT.

package mocks

import (
	context "context"

	common "github.com/ethereum/go-ethereum/common"

	mock "github.com/stretchr/testify/mock"
)

// ORM is an autogenerated mock type for the ORM type
type ORM struct {
	mock.Mock
}

// CreateAllowedSenders provides a mock function with given fields: ctx, allowedSenders
func (_m *ORM) CreateAllowedSenders(ctx context.Context, allowedSenders []common.Address) error {
	ret := _m.Called(ctx, allowedSenders)

	if len(ret) == 0 {
		panic("no return value specified for CreateAllowedSenders")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, []common.Address) error); ok {
		r0 = rf(ctx, allowedSenders)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteAllowedSenders provides a mock function with given fields: ctx, blockedSenders
func (_m *ORM) DeleteAllowedSenders(ctx context.Context, blockedSenders []common.Address) error {
	ret := _m.Called(ctx, blockedSenders)

	if len(ret) == 0 {
		panic("no return value specified for DeleteAllowedSenders")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, []common.Address) error); ok {
		r0 = rf(ctx, blockedSenders)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetAllowedSenders provides a mock function with given fields: ctx, offset, limit
func (_m *ORM) GetAllowedSenders(ctx context.Context, offset uint, limit uint) ([]common.Address, error) {
	ret := _m.Called(ctx, offset, limit)

	if len(ret) == 0 {
		panic("no return value specified for GetAllowedSenders")
	}

	var r0 []common.Address
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uint, uint) ([]common.Address, error)); ok {
		return rf(ctx, offset, limit)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uint, uint) []common.Address); ok {
		r0 = rf(ctx, offset, limit)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]common.Address)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uint, uint) error); ok {
		r1 = rf(ctx, offset, limit)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PurgeAllowedSenders provides a mock function with given fields: ctx
func (_m *ORM) PurgeAllowedSenders(ctx context.Context) error {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for PurgeAllowedSenders")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context) error); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewORM creates a new instance of ORM. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewORM(t interface {
	mock.TestingT
	Cleanup(func())
}) *ORM {
	mock := &ORM{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
