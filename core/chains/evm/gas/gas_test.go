package gas_test

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/assets"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/gas"
	gasmocks "github.com/smartcontractkit/chainlink/v2/core/chains/evm/gas/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func Test_BumpLegacyGasPriceOnly(t *testing.T) {
	t.Parallel()

	for _, test := range []struct {
		name                   string
		currentGasPrice        *assets.Wei
		originalGasPrice       *assets.Wei
		bumpPercent            uint16
		bumpWei                *assets.Wei
		maxGasPriceWei         *assets.Wei
		expectedGasPrice       *assets.Wei
		originalLimit          uint32
		limitMultiplierPercent float32
		expectedLimit          uint64
	}{
		{
			name:                   "defaults",
			currentGasPrice:        toWei("2e10"), // 20 GWei
			originalGasPrice:       toWei("3e10"), // 30 GWei
			bumpPercent:            20,
			bumpWei:                toWei("5e9"),    // 0.5 GWei
			maxGasPriceWei:         toWei("5e11"),   // 0.5 uEther
			expectedGasPrice:       toWei("3.6e10"), // 36 GWei
			originalLimit:          100000,
			limitMultiplierPercent: 1.0,
			expectedLimit:          100000,
		},
		{
			name:                   "defaults with nil currentGasPrice",
			currentGasPrice:        nil,
			originalGasPrice:       toWei("3e10"), // 30 GWei
			bumpPercent:            20,
			bumpWei:                toWei("5e9"),    // 0.5 GWei
			maxGasPriceWei:         toWei("5e11"),   // 0.5 uEther
			expectedGasPrice:       toWei("3.6e10"), // 36 GWei
			originalLimit:          100000,
			limitMultiplierPercent: 1.0,
			expectedLimit:          100000,
		},
		{
			name:                   "original + percentage wins",
			currentGasPrice:        toWei("2e10"), // 20 GWei
			originalGasPrice:       toWei("3e10"), // 30 GWei
			bumpPercent:            30,
			bumpWei:                toWei("5e9"),    // 0.5 GWei
			maxGasPriceWei:         toWei("5e11"),   // 0.5 uEther
			expectedGasPrice:       toWei("3.9e10"), // 39 GWei
			originalLimit:          100000,
			limitMultiplierPercent: 1.1,
			expectedLimit:          110000,
		},
		{
			name:                   "original + fixed wins",
			currentGasPrice:        toWei("2e10"), // 20 GWei
			originalGasPrice:       toWei("3e10"), // 30 GWei
			bumpPercent:            20,
			bumpWei:                toWei("8e9"),    // 0.8 GWei
			maxGasPriceWei:         toWei("5e11"),   // 0.5 uEther
			expectedGasPrice:       toWei("3.8e10"), // 38 GWei
			originalLimit:          100000,
			limitMultiplierPercent: 0.8,
			expectedLimit:          80000,
		},
		{
			name:                   "current wins",
			currentGasPrice:        toWei("4e10"),
			originalGasPrice:       toWei("3e10"), // 30 GWei
			bumpPercent:            20,
			bumpWei:                toWei("9e9"),  // 0.9 GWei
			maxGasPriceWei:         toWei("5e11"), // 0.5 uEther
			expectedGasPrice:       toWei("4e10"), // 40 GWei
			originalLimit:          100000,
			limitMultiplierPercent: 1.0,
			expectedLimit:          100000,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			cfg := gasmocks.NewConfig(t)
			cfg.On("EvmGasBumpPercent").Return(test.bumpPercent)
			cfg.On("EvmGasBumpWei").Return(test.bumpWei)
			cfg.On("EvmMaxGasPriceWei").Return(test.maxGasPriceWei)
			cfg.On("EvmGasLimitMultiplier").Return(test.limitMultiplierPercent)
			actual, limit, err := gas.BumpLegacyGasPriceOnly(cfg, logger.TestLogger(t), test.currentGasPrice, test.originalGasPrice, test.originalLimit, test.maxGasPriceWei)
			require.NoError(t, err)
			if actual.Cmp(test.expectedGasPrice) != 0 {
				t.Fatalf("Expected %s but got %s", test.expectedGasPrice.String(), actual.String())
			}
			assert.Equal(t, int(test.expectedLimit), int(limit))
		})
	}
}

func Test_BumpLegacyGasPriceOnly_HitsMaxError(t *testing.T) {
	t.Parallel()

	cfg := gasmocks.NewConfig(t)
	maxGasPriceWei := assets.GWei(40)
	cfg.On("EvmGasBumpPercent").Return(uint16(50))
	cfg.On("EvmGasBumpWei").Return(assets.NewWeiI(5000000000))
	cfg.On("EvmMaxGasPriceWei").Return(maxGasPriceWei)

	originalGasPrice := toWei("3e10") // 30 GWei
	_, _, err := gas.BumpLegacyGasPriceOnly(cfg, logger.TestLogger(t), nil, originalGasPrice, 42, maxGasPriceWei)
	require.Error(t, err)
	require.Contains(t, err.Error(), "bumped gas price of 45 gwei would exceed configured max gas price of 40 gwei (original price was 30 gwei)")
}

func Test_BumpLegacyGasPriceOnly_NoBumpError(t *testing.T) {
	t.Parallel()

	maxGasPriceWei := assets.GWei(40)
	lggr := logger.TestLogger(t)
	cfg := gasmocks.NewConfig(t)
	cfg.On("EvmGasBumpPercent").Return(uint16(0))
	cfg.On("EvmGasBumpWei").Return(assets.NewWeiI(0))
	cfg.On("EvmMaxGasPriceWei").Return(maxGasPriceWei)

	originalGasPrice := toWei("3e10") // 30 GWei
	_, _, err := gas.BumpLegacyGasPriceOnly(cfg, lggr, nil, originalGasPrice, 42, maxGasPriceWei)
	require.Error(t, err)
	require.Contains(t, err.Error(), "bumped gas price of 30 gwei is equal to original gas price of 30 gwei. ACTION REQUIRED: This is a configuration error, you must increase either EVM.GasEstimator.BumpPercent or EVM.GasEstimator.BumpMin")

	// Even if it's exactly the maximum
	originalGasPrice = toWei("4e10") // 40 GWei
	_, _, err = gas.BumpLegacyGasPriceOnly(cfg, lggr, nil, originalGasPrice, 42, maxGasPriceWei)
	require.Error(t, err)
	require.Contains(t, err.Error(), "bumped gas price of 40 gwei is equal to original gas price of 40 gwei. ACTION REQUIRED: This is a configuration error, you must increase either EVM.GasEstimator.BumpPercent or EVM.GasEstimator.BumpMin")
}

func Test_BumpDynamicFeeOnly(t *testing.T) {
	t.Parallel()

	for _, test := range []struct {
		name                   string
		currentTipCap          *assets.Wei
		currentBaseFee         *assets.Wei
		originalFee            gas.DynamicFee
		tipCapDefault          *assets.Wei
		bumpPercent            uint16
		bumpWei                *assets.Wei
		maxGasPriceWei         *assets.Wei
		expectedFee            gas.DynamicFee
		originalLimit          uint32
		limitMultiplierPercent float32
		expectedLimit          uint64
	}{
		{
			name:                   "defaults",
			currentTipCap:          nil,
			currentBaseFee:         nil,
			originalFee:            gas.DynamicFee{TipCap: assets.GWei(30), FeeCap: assets.GWei(4000)},
			tipCapDefault:          assets.GWei(20),
			bumpPercent:            20,
			bumpWei:                toWei("5e9"), // 0.5 GWei
			maxGasPriceWei:         assets.GWei(5000),
			expectedFee:            gas.DynamicFee{TipCap: assets.GWei(36), FeeCap: assets.GWei(4800)},
			originalLimit:          100000,
			limitMultiplierPercent: 1.0,
			expectedLimit:          100000,
		},
		{
			name:                   "original + percentage wins",
			currentTipCap:          nil,
			currentBaseFee:         nil,
			originalFee:            gas.DynamicFee{TipCap: assets.GWei(30), FeeCap: assets.GWei(100)},
			tipCapDefault:          assets.GWei(20),
			bumpPercent:            30,
			bumpWei:                toWei("5e9"),  // 0.5 GWei
			maxGasPriceWei:         toWei("5e11"), // 500GWei
			expectedFee:            gas.DynamicFee{TipCap: assets.GWei(39), FeeCap: assets.GWei(130)},
			originalLimit:          100000,
			limitMultiplierPercent: 1.1,
			expectedLimit:          110000,
		},
		{
			name:                   "original + fixed wins",
			currentTipCap:          nil,
			currentBaseFee:         nil,
			originalFee:            gas.DynamicFee{TipCap: assets.GWei(30), FeeCap: assets.GWei(400)},
			tipCapDefault:          assets.GWei(20),
			bumpPercent:            20,
			bumpWei:                toWei("8e9"),  // 0.8 GWei
			maxGasPriceWei:         toWei("5e11"), // 500GWei
			expectedFee:            gas.DynamicFee{TipCap: assets.GWei(38), FeeCap: assets.GWei(480)},
			originalLimit:          100000,
			limitMultiplierPercent: 0.8,
			expectedLimit:          80000,
		},
		{
			name:                   "default + percentage wins",
			currentTipCap:          nil,
			currentBaseFee:         nil,
			originalFee:            gas.DynamicFee{TipCap: assets.GWei(30), FeeCap: assets.GWei(400)},
			tipCapDefault:          assets.GWei(40),
			bumpPercent:            20,
			bumpWei:                toWei("5e9"),  // 0.5 GWei
			maxGasPriceWei:         toWei("5e11"), // 500GWei
			expectedFee:            gas.DynamicFee{TipCap: assets.GWei(48), FeeCap: assets.GWei(480)},
			originalLimit:          100000,
			limitMultiplierPercent: 1.0,
			expectedLimit:          100000,
		},
		{
			name:                   "default + fixed wins",
			currentTipCap:          assets.GWei(48),
			currentBaseFee:         nil,
			originalFee:            gas.DynamicFee{TipCap: assets.GWei(30), FeeCap: assets.GWei(400)},
			tipCapDefault:          assets.GWei(40),
			bumpPercent:            20,
			bumpWei:                toWei("9e9"),  // 0.9 GWei
			maxGasPriceWei:         toWei("5e11"), // 500GWei
			expectedFee:            gas.DynamicFee{TipCap: assets.GWei(49), FeeCap: assets.GWei(480)},
			originalLimit:          100000,
			limitMultiplierPercent: 1.0,
			expectedLimit:          100000,
		},
		{
			name:                   "higher current tip cap wins",
			currentTipCap:          assets.GWei(50),
			currentBaseFee:         nil,
			originalFee:            gas.DynamicFee{TipCap: assets.GWei(30), FeeCap: assets.GWei(400)},
			tipCapDefault:          assets.GWei(40),
			bumpPercent:            20,
			bumpWei:                toWei("9e9"),  // 0.9 GWei
			maxGasPriceWei:         toWei("5e11"), // 500GWei
			expectedFee:            gas.DynamicFee{TipCap: assets.GWei(50), FeeCap: assets.GWei(480)},
			originalLimit:          100000,
			limitMultiplierPercent: 1.0,
			expectedLimit:          100000,
		},
		{
			name:                   "if bumped tip cap would exceed bumped fee cap, adds fixed value to expectedFee",
			currentTipCap:          nil,
			currentBaseFee:         nil,
			originalFee:            gas.DynamicFee{TipCap: assets.GWei(10), FeeCap: assets.GWei(20)},
			tipCapDefault:          assets.GWei(5),
			bumpPercent:            5,
			bumpWei:                assets.GWei(50),
			maxGasPriceWei:         toWei("5e11"), // 500GWei
			expectedFee:            gas.DynamicFee{TipCap: assets.GWei(60), FeeCap: assets.GWei(70)},
			originalLimit:          100000,
			limitMultiplierPercent: 1.0,
			expectedLimit:          100000,
		},
		{
			name:                   "ignores current base fee and uses previous fee cap if calculated fee cap would be lower",
			currentTipCap:          assets.GWei(20),
			currentBaseFee:         assets.GWei(100),
			originalFee:            gas.DynamicFee{TipCap: assets.GWei(30), FeeCap: assets.GWei(400)},
			tipCapDefault:          assets.GWei(20),
			bumpPercent:            20,
			bumpWei:                toWei("5e9"), // 0.5 GWei
			maxGasPriceWei:         assets.GWei(5000),
			expectedFee:            gas.DynamicFee{TipCap: assets.GWei(36), FeeCap: assets.GWei(480)},
			originalLimit:          100000,
			limitMultiplierPercent: 1.0,
			expectedLimit:          100000,
		},
		{
			name:           "uses current base fee to calculate fee cap if that would be higher than the existing one",
			currentTipCap:  assets.GWei(20),
			currentBaseFee: assets.GWei(1000),
			originalFee:    gas.DynamicFee{TipCap: assets.GWei(30), FeeCap: assets.GWei(400)},
			tipCapDefault:  assets.GWei(20),
			bumpPercent:    20,
			bumpWei:        toWei("5e9"), // 0.5 GWei
			maxGasPriceWei: assets.GWei(5000),
			// base fee * 4 blocks * 1.125 % plus new tip cap to give max
			// 1000 * (1.125 ^ 4) + 36 ~= 1637
			expectedFee:            gas.DynamicFee{TipCap: assets.GWei(36), FeeCap: assets.NewWeiI(1637806640625)},
			originalLimit:          100000,
			limitMultiplierPercent: 1.0,
			expectedLimit:          100000,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			cfg := gasmocks.NewConfig(t)
			cfg.On("EvmGasBumpPercent").Return(test.bumpPercent)
			cfg.On("EvmGasTipCapDefault").Return(test.tipCapDefault)
			cfg.On("EvmGasBumpWei").Return(test.bumpWei)
			cfg.On("EvmMaxGasPriceWei").Return(test.maxGasPriceWei)
			cfg.On("EvmGasLimitMultiplier").Return(test.limitMultiplierPercent)
			if test.currentBaseFee != nil {
				cfg.On("BlockHistoryEstimatorEIP1559FeeCapBufferBlocks").Return(uint16(4))
			}
			actual, limit, err := gas.BumpDynamicFeeOnly(cfg, logger.TestLogger(t), test.currentTipCap, test.currentBaseFee, test.originalFee, test.originalLimit, test.maxGasPriceWei)
			require.NoError(t, err)
			if actual.TipCap.Cmp(test.expectedFee.TipCap) != 0 {
				t.Fatalf("TipCap not equal, expected %s but got %s", test.expectedFee.TipCap.String(), actual.TipCap.String())
			}
			if actual.FeeCap.Cmp(test.expectedFee.FeeCap) != 0 {
				t.Fatalf("FeeCap not equal, expected %s but got %s", test.expectedFee.FeeCap.String(), actual.FeeCap.String())
			}
			assert.Equal(t, int(test.expectedLimit), int(limit))
		})
	}
}

func Test_BumpDynamicFeeOnly_HitsMaxError(t *testing.T) {
	t.Parallel()

	maxGasPriceWei := assets.GWei(40)
	cfg := gasmocks.NewConfig(t)
	cfg.On("EvmGasBumpPercent").Return(uint16(50))
	cfg.On("EvmGasTipCapDefault").Return(assets.GWei(0))
	cfg.On("EvmGasBumpWei").Return(assets.NewWeiI(5000000000))
	cfg.On("EvmMaxGasPriceWei").Return(maxGasPriceWei)

	t.Run("tip cap hits max", func(t *testing.T) {
		originalFee := gas.DynamicFee{TipCap: assets.GWei(30), FeeCap: assets.GWei(100)}
		_, _, err := gas.BumpDynamicFeeOnly(cfg, logger.TestLogger(t), nil, nil, originalFee, 42, maxGasPriceWei)
		require.Error(t, err)
		require.Contains(t, err.Error(), "bumped tip cap of 45 gwei would exceed configured max gas price of 40 gwei (original fee: tip cap 30 gwei, fee cap 100 gwei)")
	})

	t.Run("fee cap hits max", func(t *testing.T) {
		originalFee := gas.DynamicFee{TipCap: assets.GWei(10), FeeCap: assets.GWei(100)}
		_, _, err := gas.BumpDynamicFeeOnly(cfg, logger.TestLogger(t), nil, nil, originalFee, 42, maxGasPriceWei)
		require.Error(t, err)
		require.Contains(t, err.Error(), "bumped fee cap of 150 gwei would exceed configured max gas price of 40 gwei (original fee: tip cap 10 gwei, fee cap 100 gwei)")
	})
}

// toWei is used to convert scientific notation string to a *assets.Wei
func toWei(input string) *assets.Wei {
	flt, _, err := big.ParseFloat(input, 10, 0, big.ToNearestEven)
	if err != nil {
		panic(fmt.Sprintf("unable to parse '%s' into a big.Float: %v", input, err))
	}
	var i = new(big.Int)
	i, _ = flt.Int(i)
	return assets.NewWei(i)
}
