// Code generated by mockery v2.43.2. DO NOT EDIT.

package mocks

import (
	big "math/big"

	assets "github.com/smartcontractkit/chainlink/v2/core/chains/evm/assets"

	context "context"

	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"

	gas "github.com/smartcontractkit/chainlink/v2/core/chains/evm/gas"

	mock "github.com/stretchr/testify/mock"

	rollups "github.com/smartcontractkit/chainlink/v2/core/chains/evm/gas/rollups"

	types "github.com/smartcontractkit/chainlink/v2/common/fee/types"
)

// EvmFeeEstimator is an autogenerated mock type for the EvmFeeEstimator type
type EvmFeeEstimator struct {
	mock.Mock
}

// BumpFee provides a mock function with given fields: ctx, originalFee, feeLimit, maxFeePrice, attempts
func (_m *EvmFeeEstimator) BumpFee(ctx context.Context, originalFee gas.EvmFee, feeLimit uint64, maxFeePrice *assets.Wei, attempts []gas.EvmPriorAttempt) (gas.EvmFee, uint64, error) {
	ret := _m.Called(ctx, originalFee, feeLimit, maxFeePrice, attempts)

	if len(ret) == 0 {
		panic("no return value specified for BumpFee")
	}

	var r0 gas.EvmFee
	var r1 uint64
	var r2 error
	if rf, ok := ret.Get(0).(func(context.Context, gas.EvmFee, uint64, *assets.Wei, []gas.EvmPriorAttempt) (gas.EvmFee, uint64, error)); ok {
		return rf(ctx, originalFee, feeLimit, maxFeePrice, attempts)
	}
	if rf, ok := ret.Get(0).(func(context.Context, gas.EvmFee, uint64, *assets.Wei, []gas.EvmPriorAttempt) gas.EvmFee); ok {
		r0 = rf(ctx, originalFee, feeLimit, maxFeePrice, attempts)
	} else {
		r0 = ret.Get(0).(gas.EvmFee)
	}

	if rf, ok := ret.Get(1).(func(context.Context, gas.EvmFee, uint64, *assets.Wei, []gas.EvmPriorAttempt) uint64); ok {
		r1 = rf(ctx, originalFee, feeLimit, maxFeePrice, attempts)
	} else {
		r1 = ret.Get(1).(uint64)
	}

	if rf, ok := ret.Get(2).(func(context.Context, gas.EvmFee, uint64, *assets.Wei, []gas.EvmPriorAttempt) error); ok {
		r2 = rf(ctx, originalFee, feeLimit, maxFeePrice, attempts)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// Close provides a mock function with given fields:
func (_m *EvmFeeEstimator) Close() error {
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

// GetFee provides a mock function with given fields: ctx, calldata, feeLimit, maxFeePrice, opts
func (_m *EvmFeeEstimator) GetFee(ctx context.Context, calldata []byte, feeLimit uint64, maxFeePrice *assets.Wei, opts ...types.Opt) (gas.EvmFee, uint64, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, calldata, feeLimit, maxFeePrice)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for GetFee")
	}

	var r0 gas.EvmFee
	var r1 uint64
	var r2 error
	if rf, ok := ret.Get(0).(func(context.Context, []byte, uint64, *assets.Wei, ...types.Opt) (gas.EvmFee, uint64, error)); ok {
		return rf(ctx, calldata, feeLimit, maxFeePrice, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, []byte, uint64, *assets.Wei, ...types.Opt) gas.EvmFee); ok {
		r0 = rf(ctx, calldata, feeLimit, maxFeePrice, opts...)
	} else {
		r0 = ret.Get(0).(gas.EvmFee)
	}

	if rf, ok := ret.Get(1).(func(context.Context, []byte, uint64, *assets.Wei, ...types.Opt) uint64); ok {
		r1 = rf(ctx, calldata, feeLimit, maxFeePrice, opts...)
	} else {
		r1 = ret.Get(1).(uint64)
	}

	if rf, ok := ret.Get(2).(func(context.Context, []byte, uint64, *assets.Wei, ...types.Opt) error); ok {
		r2 = rf(ctx, calldata, feeLimit, maxFeePrice, opts...)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// GetMaxCost provides a mock function with given fields: ctx, amount, calldata, feeLimit, maxFeePrice, opts
func (_m *EvmFeeEstimator) GetMaxCost(ctx context.Context, amount assets.Eth, calldata []byte, feeLimit uint64, maxFeePrice *assets.Wei, opts ...types.Opt) (*big.Int, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, amount, calldata, feeLimit, maxFeePrice)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for GetMaxCost")
	}

	var r0 *big.Int
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, assets.Eth, []byte, uint64, *assets.Wei, ...types.Opt) (*big.Int, error)); ok {
		return rf(ctx, amount, calldata, feeLimit, maxFeePrice, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, assets.Eth, []byte, uint64, *assets.Wei, ...types.Opt) *big.Int); ok {
		r0 = rf(ctx, amount, calldata, feeLimit, maxFeePrice, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*big.Int)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, assets.Eth, []byte, uint64, *assets.Wei, ...types.Opt) error); ok {
		r1 = rf(ctx, amount, calldata, feeLimit, maxFeePrice, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// HealthReport provides a mock function with given fields:
func (_m *EvmFeeEstimator) HealthReport() map[string]error {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for HealthReport")
	}

	var r0 map[string]error
	if rf, ok := ret.Get(0).(func() map[string]error); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[string]error)
		}
	}

	return r0
}

// L1Oracle provides a mock function with given fields:
func (_m *EvmFeeEstimator) L1Oracle() rollups.L1Oracle {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for L1Oracle")
	}

	var r0 rollups.L1Oracle
	if rf, ok := ret.Get(0).(func() rollups.L1Oracle); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(rollups.L1Oracle)
		}
	}

	return r0
}

// Name provides a mock function with given fields:
func (_m *EvmFeeEstimator) Name() string {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Name")
	}

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// OnNewLongestChain provides a mock function with given fields: ctx, head
func (_m *EvmFeeEstimator) OnNewLongestChain(ctx context.Context, head *evmtypes.Head) {
	_m.Called(ctx, head)
}

// Ready provides a mock function with given fields:
func (_m *EvmFeeEstimator) Ready() error {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Ready")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Start provides a mock function with given fields: _a0
func (_m *EvmFeeEstimator) Start(_a0 context.Context) error {
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

// NewEvmFeeEstimator creates a new instance of EvmFeeEstimator. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewEvmFeeEstimator(t interface {
	mock.TestingT
	Cleanup(func())
}) *EvmFeeEstimator {
	mock := &EvmFeeEstimator{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
