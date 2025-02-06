package prices

import (
	"context"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"

	cciptypes "github.com/smartcontractkit/chainlink-common/pkg/types/ccip"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"
	"github.com/smartcontractkit/chainlink-integrations/evm/assets"
	"github.com/smartcontractkit/chainlink-integrations/evm/gas"
	"github.com/smartcontractkit/chainlink-integrations/evm/gas/mocks"
)

func TestExecPriceEstimator_GetGasPrice(t *testing.T) {
	ctx := context.Background()

	testCases := []struct {
		name                      string
		sourceFeeEstimatorRespFee gas.EvmFee
		sourceFeeEstimatorRespErr error
		maxGasPrice               *big.Int
		expPrice                  *big.Int
		expErr                    bool
	}{
		{
			name: "gets legacy gas price",
			sourceFeeEstimatorRespFee: gas.EvmFee{
				GasPrice:   assets.NewWei(big.NewInt(10)),
				DynamicFee: gas.DynamicFee{},
			},
			sourceFeeEstimatorRespErr: nil,
			maxGasPrice:               big.NewInt(1),
			expPrice:                  big.NewInt(10),
			expErr:                    false,
		},
		{
			name: "gets dynamic gas price",
			sourceFeeEstimatorRespFee: gas.EvmFee{
				GasPrice:   nil,
				DynamicFee: gas.DynamicFee{GasFeeCap: assets.NewWei(big.NewInt(20))},
			},
			sourceFeeEstimatorRespErr: nil,
			maxGasPrice:               big.NewInt(1),
			expPrice:                  big.NewInt(20),
			expErr:                    false,
		},
		{
			name: "gets dynamic gas price over legacy gas price",
			sourceFeeEstimatorRespFee: gas.EvmFee{
				GasPrice:   assets.NewWei(big.NewInt(10)),
				DynamicFee: gas.DynamicFee{GasFeeCap: assets.NewWei(big.NewInt(20))},
			},
			sourceFeeEstimatorRespErr: nil,
			maxGasPrice:               big.NewInt(1),
			expPrice:                  big.NewInt(20),
			expErr:                    false,
		},
		{
			name: "fee estimator error",
			sourceFeeEstimatorRespFee: gas.EvmFee{
				GasPrice:   assets.NewWei(big.NewInt(10)),
				DynamicFee: gas.DynamicFee{},
			},
			sourceFeeEstimatorRespErr: errors.New("fee estimator error"),
			maxGasPrice:               big.NewInt(1),
			expPrice:                  nil,
			expErr:                    true,
		},
		{
			name: "nil gas price error",
			sourceFeeEstimatorRespFee: gas.EvmFee{
				GasPrice:   nil,
				DynamicFee: gas.DynamicFee{},
			},
			sourceFeeEstimatorRespErr: nil,
			maxGasPrice:               big.NewInt(1),
			expPrice:                  nil,
			expErr:                    true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			sourceFeeEstimator := mocks.NewEvmFeeEstimator(t)
			sourceFeeEstimator.On("GetFee", ctx, []byte(nil), uint64(0), assets.NewWei(tc.maxGasPrice), (*common.Address)(nil), (*common.Address)(nil)).Return(
				tc.sourceFeeEstimatorRespFee, uint64(0), tc.sourceFeeEstimatorRespErr)

			g := ExecGasPriceEstimator{
				estimator:   sourceFeeEstimator,
				maxGasPrice: tc.maxGasPrice,
			}

			gasPrice, err := g.GetGasPrice(ctx)
			if tc.expErr {
				assert.Nil(t, gasPrice)
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tc.expPrice, gasPrice)
		})
	}
}

func TestExecPriceEstimator_DenoteInUSD(t *testing.T) {
	val1e18 := func(val int64) *big.Int { return new(big.Int).Mul(big.NewInt(1e18), big.NewInt(val)) }

	testCases := []struct {
		name        string
		gasPrice    *big.Int
		nativePrice *big.Int
		expPrice    *big.Int
	}{
		{
			name:        "base",
			gasPrice:    big.NewInt(1e9),
			nativePrice: val1e18(2_000),
			expPrice:    big.NewInt(2_000e9),
		},
		{
			name:        "low price truncates to 0",
			gasPrice:    big.NewInt(1e9),
			nativePrice: big.NewInt(1),
			expPrice:    big.NewInt(0),
		},
		{
			name:        "high price",
			gasPrice:    val1e18(1),
			nativePrice: val1e18(2_000),
			expPrice:    val1e18(2_000),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := tests.Context(t)
			g := ExecGasPriceEstimator{}

			gasPrice, err := g.DenoteInUSD(ctx, tc.gasPrice, tc.nativePrice)
			assert.NoError(t, err)
			assert.True(t, tc.expPrice.Cmp(gasPrice) == 0)
		})
	}
}

func TestExecPriceEstimator_Median(t *testing.T) {
	val1e18 := func(val int64) *big.Int { return new(big.Int).Mul(big.NewInt(1e18), big.NewInt(val)) }

	testCases := []struct {
		name      string
		gasPrices []*big.Int
		expMedian *big.Int
	}{
		{
			name:      "base",
			gasPrices: []*big.Int{big.NewInt(1), big.NewInt(2), big.NewInt(3)},
			expMedian: big.NewInt(2),
		},
		{
			name:      "median 1",
			gasPrices: []*big.Int{big.NewInt(1)},
			expMedian: big.NewInt(1),
		},
		{
			name:      "median 2",
			gasPrices: []*big.Int{big.NewInt(1), big.NewInt(2)},
			expMedian: big.NewInt(2),
		},
		{
			name:      "large values",
			gasPrices: []*big.Int{val1e18(5), val1e18(4), val1e18(3), val1e18(2), val1e18(1)},
			expMedian: val1e18(3),
		},
		{
			name:      "zeros",
			gasPrices: []*big.Int{big.NewInt(0), big.NewInt(0), big.NewInt(0)},
			expMedian: big.NewInt(0),
		},
		{
			name:      "unsorted even number of prices",
			gasPrices: []*big.Int{big.NewInt(4), big.NewInt(2), big.NewInt(3), big.NewInt(1)},
			expMedian: big.NewInt(3),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := tests.Context(t)
			g := ExecGasPriceEstimator{}

			gasPrice, err := g.Median(ctx, tc.gasPrices)
			assert.NoError(t, err)
			assert.True(t, tc.expMedian.Cmp(gasPrice) == 0)
		})
	}
}

func TestExecPriceEstimator_Deviates(t *testing.T) {
	testCases := []struct {
		name         string
		gasPrice1    *big.Int
		gasPrice2    *big.Int
		deviationPPB int64
		expDeviates  bool
	}{
		{
			name:         "base",
			gasPrice1:    big.NewInt(100e8),
			gasPrice2:    big.NewInt(79e8),
			deviationPPB: 2e8,
			expDeviates:  true,
		},
		{
			name:         "negative difference also deviates",
			gasPrice1:    big.NewInt(100e8),
			gasPrice2:    big.NewInt(121e8),
			deviationPPB: 2e8,
			expDeviates:  true,
		},
		{
			name:         "larger difference deviates",
			gasPrice1:    big.NewInt(100e8),
			gasPrice2:    big.NewInt(70e8),
			deviationPPB: 2e8,
			expDeviates:  true,
		},
		{
			name:         "smaller difference does not deviate",
			gasPrice1:    big.NewInt(100e8),
			gasPrice2:    big.NewInt(90e8),
			deviationPPB: 2e8,
			expDeviates:  false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := tests.Context(t)
			g := ExecGasPriceEstimator{
				deviationPPB: tc.deviationPPB,
			}

			deviated, err := g.Deviates(ctx, tc.gasPrice1, tc.gasPrice2)
			assert.NoError(t, err)
			if tc.expDeviates {
				assert.True(t, deviated)
			} else {
				assert.False(t, deviated)
			}
		})
	}
}

func TestExecPriceEstimator_EstimateMsgCostUSD(t *testing.T) {
	testCases := []struct {
		name               string
		gasPrice           *big.Int
		wrappedNativePrice *big.Int
		msg                cciptypes.EVM2EVMOnRampCCIPSendRequestedWithMeta
		expUSD             *big.Int
	}{
		{
			name:               "base",
			gasPrice:           big.NewInt(1e9),  // 1 gwei
			wrappedNativePrice: big.NewInt(1e18), // $1
			msg: cciptypes.EVM2EVMOnRampCCIPSendRequestedWithMeta{
				EVM2EVMMessage: cciptypes.EVM2EVMMessage{
					GasLimit:     big.NewInt(100_000),
					Data:         []byte{},
					TokenAmounts: []cciptypes.TokenAmount{},
				},
			},
			expUSD: big.NewInt(300_000e9),
		},
		{
			name:               "base with data",
			gasPrice:           big.NewInt(1e9),  // 1 gwei
			wrappedNativePrice: big.NewInt(1e18), // $1
			msg: cciptypes.EVM2EVMOnRampCCIPSendRequestedWithMeta{
				EVM2EVMMessage: cciptypes.EVM2EVMMessage{
					GasLimit:     big.NewInt(100_000),
					Data:         make([]byte, 1_000),
					TokenAmounts: []cciptypes.TokenAmount{},
				},
			},
			expUSD: big.NewInt(316_000e9),
		},
		{
			name:               "base with data and tokens",
			gasPrice:           big.NewInt(1e9),  // 1 gwei
			wrappedNativePrice: big.NewInt(1e18), // $1
			msg: cciptypes.EVM2EVMOnRampCCIPSendRequestedWithMeta{
				EVM2EVMMessage: cciptypes.EVM2EVMMessage{
					GasLimit:     big.NewInt(100_000),
					Data:         make([]byte, 1_000),
					TokenAmounts: make([]cciptypes.TokenAmount, 5),
				},
			},
			expUSD: big.NewInt(366_000e9),
		},
		{
			name:               "empty msg",
			gasPrice:           big.NewInt(1e9),  // 1 gwei
			wrappedNativePrice: big.NewInt(1e18), // $1
			msg: cciptypes.EVM2EVMOnRampCCIPSendRequestedWithMeta{
				EVM2EVMMessage: cciptypes.EVM2EVMMessage{
					GasLimit:     big.NewInt(0),
					Data:         []byte{},
					TokenAmounts: []cciptypes.TokenAmount{},
				},
			},
			expUSD: big.NewInt(200_000e9),
		},
		{
			name:               "double native price",
			gasPrice:           big.NewInt(1e9),  // 1 gwei
			wrappedNativePrice: big.NewInt(2e18), // $1
			msg: cciptypes.EVM2EVMOnRampCCIPSendRequestedWithMeta{
				EVM2EVMMessage: cciptypes.EVM2EVMMessage{
					GasLimit:     big.NewInt(0),
					Data:         []byte{},
					TokenAmounts: []cciptypes.TokenAmount{},
				},
			},
			expUSD: big.NewInt(400_000e9),
		},
		{
			name:               "zero gas price",
			gasPrice:           big.NewInt(0),    // 1 gwei
			wrappedNativePrice: big.NewInt(1e18), // $1
			msg: cciptypes.EVM2EVMOnRampCCIPSendRequestedWithMeta{
				EVM2EVMMessage: cciptypes.EVM2EVMMessage{
					GasLimit:     big.NewInt(0),
					Data:         []byte{},
					TokenAmounts: []cciptypes.TokenAmount{},
				},
			},
			expUSD: big.NewInt(0),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := tests.Context(t)
			g := ExecGasPriceEstimator{}

			costUSD, err := g.EstimateMsgCostUSD(ctx, tc.gasPrice, tc.wrappedNativePrice, tc.msg)
			assert.NoError(t, err)
			assert.Equal(t, tc.expUSD, costUSD)
		})
	}
}
