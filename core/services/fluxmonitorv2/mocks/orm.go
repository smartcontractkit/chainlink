// Code generated by mockery v2.43.2. DO NOT EDIT.

package mocks

import (
	context "context"

	common "github.com/ethereum/go-ethereum/common"

	fluxmonitorv2 "github.com/smartcontractkit/chainlink/v2/core/services/fluxmonitorv2"

	mock "github.com/stretchr/testify/mock"

	sqlutil "github.com/smartcontractkit/chainlink-common/pkg/sqlutil"
)

// ORM is an autogenerated mock type for the ORM type
type ORM struct {
	mock.Mock
}

// CountFluxMonitorRoundStats provides a mock function with given fields: ctx
func (_m *ORM) CountFluxMonitorRoundStats(ctx context.Context) (int, error) {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for CountFluxMonitorRoundStats")
	}

	var r0 int
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) (int, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) int); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Get(0).(int)
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CreateEthTransaction provides a mock function with given fields: ctx, fromAddress, toAddress, payload, gasLimit, idempotencyKey
func (_m *ORM) CreateEthTransaction(ctx context.Context, fromAddress common.Address, toAddress common.Address, payload []byte, gasLimit uint64, idempotencyKey *string) error {
	ret := _m.Called(ctx, fromAddress, toAddress, payload, gasLimit, idempotencyKey)

	if len(ret) == 0 {
		panic("no return value specified for CreateEthTransaction")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, common.Address, common.Address, []byte, uint64, *string) error); ok {
		r0 = rf(ctx, fromAddress, toAddress, payload, gasLimit, idempotencyKey)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteFluxMonitorRoundsBackThrough provides a mock function with given fields: ctx, aggregator, roundID
func (_m *ORM) DeleteFluxMonitorRoundsBackThrough(ctx context.Context, aggregator common.Address, roundID uint32) error {
	ret := _m.Called(ctx, aggregator, roundID)

	if len(ret) == 0 {
		panic("no return value specified for DeleteFluxMonitorRoundsBackThrough")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, common.Address, uint32) error); ok {
		r0 = rf(ctx, aggregator, roundID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// FindOrCreateFluxMonitorRoundStats provides a mock function with given fields: ctx, aggregator, roundID, newRoundLogs
func (_m *ORM) FindOrCreateFluxMonitorRoundStats(ctx context.Context, aggregator common.Address, roundID uint32, newRoundLogs uint) (fluxmonitorv2.FluxMonitorRoundStatsV2, error) {
	ret := _m.Called(ctx, aggregator, roundID, newRoundLogs)

	if len(ret) == 0 {
		panic("no return value specified for FindOrCreateFluxMonitorRoundStats")
	}

	var r0 fluxmonitorv2.FluxMonitorRoundStatsV2
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, common.Address, uint32, uint) (fluxmonitorv2.FluxMonitorRoundStatsV2, error)); ok {
		return rf(ctx, aggregator, roundID, newRoundLogs)
	}
	if rf, ok := ret.Get(0).(func(context.Context, common.Address, uint32, uint) fluxmonitorv2.FluxMonitorRoundStatsV2); ok {
		r0 = rf(ctx, aggregator, roundID, newRoundLogs)
	} else {
		r0 = ret.Get(0).(fluxmonitorv2.FluxMonitorRoundStatsV2)
	}

	if rf, ok := ret.Get(1).(func(context.Context, common.Address, uint32, uint) error); ok {
		r1 = rf(ctx, aggregator, roundID, newRoundLogs)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MostRecentFluxMonitorRoundID provides a mock function with given fields: ctx, aggregator
func (_m *ORM) MostRecentFluxMonitorRoundID(ctx context.Context, aggregator common.Address) (uint32, error) {
	ret := _m.Called(ctx, aggregator)

	if len(ret) == 0 {
		panic("no return value specified for MostRecentFluxMonitorRoundID")
	}

	var r0 uint32
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, common.Address) (uint32, error)); ok {
		return rf(ctx, aggregator)
	}
	if rf, ok := ret.Get(0).(func(context.Context, common.Address) uint32); ok {
		r0 = rf(ctx, aggregator)
	} else {
		r0 = ret.Get(0).(uint32)
	}

	if rf, ok := ret.Get(1).(func(context.Context, common.Address) error); ok {
		r1 = rf(ctx, aggregator)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateFluxMonitorRoundStats provides a mock function with given fields: ctx, aggregator, roundID, runID, newRoundLogsAddition
func (_m *ORM) UpdateFluxMonitorRoundStats(ctx context.Context, aggregator common.Address, roundID uint32, runID int64, newRoundLogsAddition uint) error {
	ret := _m.Called(ctx, aggregator, roundID, runID, newRoundLogsAddition)

	if len(ret) == 0 {
		panic("no return value specified for UpdateFluxMonitorRoundStats")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, common.Address, uint32, int64, uint) error); ok {
		r0 = rf(ctx, aggregator, roundID, runID, newRoundLogsAddition)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// WithDataSource provides a mock function with given fields: _a0
func (_m *ORM) WithDataSource(_a0 sqlutil.DataSource) fluxmonitorv2.ORM {
	ret := _m.Called(_a0)

	if len(ret) == 0 {
		panic("no return value specified for WithDataSource")
	}

	var r0 fluxmonitorv2.ORM
	if rf, ok := ret.Get(0).(func(sqlutil.DataSource) fluxmonitorv2.ORM); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(fluxmonitorv2.ORM)
		}
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
