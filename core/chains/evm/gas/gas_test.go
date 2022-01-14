package gas_test

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/core/assets"
	"github.com/smartcontractkit/chainlink/core/chains/evm/gas"
	gasmocks "github.com/smartcontractkit/chainlink/core/chains/evm/gas/mocks"
	"github.com/smartcontractkit/chainlink/core/logger"
)

func Test_BumpLegacyGasPriceOnly(t *testing.T) {
	t.Parallel()

	for _, test := range []struct {
		name                   string
		currentGasPrice        *big.Int
		originalGasPrice       *big.Int
		bumpPercent            uint16
		bumpWei                *big.Int
		maxGasPriceWei         *big.Int
		expectedGasPrice       *big.Int
		originalLimit          uint64
		limitMultiplierPercent float32
		expectedLimit          uint64
	}{
		{
			name:                   "defaults",
			currentGasPrice:        toBigInt("2e10"), // 20 GWei
			originalGasPrice:       toBigInt("3e10"), // 30 GWei
			bumpPercent:            20,
			bumpWei:                toBigInt("5e9"),    // 0.5 GWei
			maxGasPriceWei:         toBigInt("5e11"),   // 0.5 uEther
			expectedGasPrice:       toBigInt("3.6e10"), // 36 GWei
			originalLimit:          100000,
			limitMultiplierPercent: 1.0,
			expectedLimit:          100000,
		},
		{
			name:                   "defaults with nil currentGasPrice",
			currentGasPrice:        nil,
			originalGasPrice:       toBigInt("3e10"), // 30 GWei
			bumpPercent:            20,
			bumpWei:                toBigInt("5e9"),    // 0.5 GWei
			maxGasPriceWei:         toBigInt("5e11"),   // 0.5 uEther
			expectedGasPrice:       toBigInt("3.6e10"), // 36 GWei
			originalLimit:          100000,
			limitMultiplierPercent: 1.0,
			expectedLimit:          100000,
		},
		{
			name:                   "original + percentage wins",
			currentGasPrice:        toBigInt("2e10"), // 20 GWei
			originalGasPrice:       toBigInt("3e10"), // 30 GWei
			bumpPercent:            30,
			bumpWei:                toBigInt("5e9"),    // 0.5 GWei
			maxGasPriceWei:         toBigInt("5e11"),   // 0.5 uEther
			expectedGasPrice:       toBigInt("3.9e10"), // 39 GWei
			originalLimit:          100000,
			limitMultiplierPercent: 1.1,
			expectedLimit:          110000,
		},
		{
			name:                   "original + fixed wins",
			currentGasPrice:        toBigInt("2e10"), // 20 GWei
			originalGasPrice:       toBigInt("3e10"), // 30 GWei
			bumpPercent:            20,
			bumpWei:                toBigInt("8e9"),    // 0.8 GWei
			maxGasPriceWei:         toBigInt("5e11"),   // 0.5 uEther
			expectedGasPrice:       toBigInt("3.8e10"), // 38 GWei
			originalLimit:          100000,
			limitMultiplierPercent: 0.8,
			expectedLimit:          80000,
		},
		{
			name:                   "current wins",
			currentGasPrice:        toBigInt("4e10"),
			originalGasPrice:       toBigInt("3e10"), // 30 GWei
			bumpPercent:            20,
			bumpWei:                toBigInt("9e9"),  // 0.9 GWei
			maxGasPriceWei:         toBigInt("5e11"), // 0.5 uEther
			expectedGasPrice:       toBigInt("4e10"), // 40 GWei
			originalLimit:          100000,
			limitMultiplierPercent: 1.0,
			expectedLimit:          100000,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			cfg := new(gasmocks.Config)
			cfg.Test(t)
			cfg.On("EvmGasBumpPercent").Return(test.bumpPercent)
			cfg.On("EvmGasBumpWei").Return(test.bumpWei)
			cfg.On("EvmMaxGasPriceWei").Return(test.maxGasPriceWei)
			cfg.On("EvmGasLimitMultiplier").Return(test.limitMultiplierPercent)
			actual, limit, err := gas.BumpLegacyGasPriceOnly(cfg, logger.TestLogger(t), test.currentGasPrice, test.originalGasPrice, test.originalLimit)
			require.NoError(t, err)
			if actual.Cmp(test.expectedGasPrice) != 0 {
				t.Fatalf("Expected %s but got %s", test.expectedGasPrice.String(), actual.String())
			}
			assert.Equal(t, int(test.expectedLimit), int(limit))
			cfg.AssertExpectations(t)
		})
	}
}

func Test_BumpLegacyGasPriceOnly_HitsMaxError(t *testing.T) {
	t.Parallel()
	cfg := new(gasmocks.Config)
	cfg.On("EvmGasBumpPercent").Return(uint16(50))
	cfg.On("EvmGasPriceDefault").Return(assets.GWei(20))
	cfg.On("EvmGasBumpWei").Return(assets.Wei(5000000000))
	cfg.On("EvmMaxGasPriceWei").Return(assets.GWei(40))

	originalGasPrice := toBigInt("3e10") // 30 GWei
	_, _, err := gas.BumpLegacyGasPriceOnly(cfg, logger.TestLogger(t), nil, originalGasPrice, 42)
	require.Error(t, err)
	require.Contains(t, err.Error(), "bumped gas price of 45000000000 would exceed configured max gas price of 40000000000 (original price was 30000000000)")
}

func Test_BumpLegacyGasPriceOnly_NoBumpError(t *testing.T) {
	t.Parallel()
	lggr := logger.TestLogger(t)
	cfg := new(gasmocks.Config)
	cfg.On("EvmGasBumpPercent").Return(uint16(0))
	cfg.On("EvmGasBumpWei").Return(big.NewInt(0))
	cfg.On("EvmMaxGasPriceWei").Return(assets.GWei(40))
	cfg.On("EvmGasPriceDefault").Return(assets.GWei(20))

	originalGasPrice := toBigInt("3e10") // 30 GWei
	_, _, err := gas.BumpLegacyGasPriceOnly(cfg, lggr, nil, originalGasPrice, 42)
	require.Error(t, err)
	require.Contains(t, err.Error(), "bumped gas price of 30000000000 is equal to original gas price of 30000000000. ACTION REQUIRED: This is a configuration error, you must increase either ETH_GAS_BUMP_PERCENT or ETH_GAS_BUMP_WEI")

	// Even if it's exactly the maximum
	originalGasPrice = toBigInt("4e10") // 40 GWei
	_, _, err = gas.BumpLegacyGasPriceOnly(cfg, lggr, nil, originalGasPrice, 42)
	require.Error(t, err)
	require.Contains(t, err.Error(), "bumped gas price of 40000000000 is equal to original gas price of 40000000000. ACTION REQUIRED: This is a configuration error, you must increase either ETH_GAS_BUMP_PERCENT or ETH_GAS_BUMP_WEI")
}

func Test_BumpDynamicFeeOnly(t *testing.T) {
	t.Parallel()

	for _, test := range []struct {
		name                   string
		currentTipCap          *big.Int
		originalFee            gas.DynamicFee
		tipCapDefault          *big.Int
		bumpPercent            uint16
		bumpWei                *big.Int
		maxGasPriceWei         *big.Int
		expectedFee            gas.DynamicFee
		originalLimit          uint64
		limitMultiplierPercent float32
		expectedLimit          uint64
	}{
		{
			name:                   "defaults",
			currentTipCap:          nil,
			originalFee:            gas.DynamicFee{TipCap: assets.GWei(30), FeeCap: assets.GWei(5000)},
			tipCapDefault:          assets.GWei(20),
			bumpPercent:            20,
			bumpWei:                toBigInt("5e9"), // 0.5 GWei
			maxGasPriceWei:         assets.GWei(5000),
			expectedFee:            gas.DynamicFee{TipCap: assets.GWei(36), FeeCap: assets.GWei(5000)},
			originalLimit:          100000,
			limitMultiplierPercent: 1.0,
			expectedLimit:          100000,
		},
		{
			name:                   "original + percentage wins",
			currentTipCap:          nil,
			originalFee:            gas.DynamicFee{TipCap: assets.GWei(30), FeeCap: assets.GWei(5000)},
			tipCapDefault:          assets.GWei(20),
			bumpPercent:            30,
			bumpWei:                toBigInt("5e9"),  // 0.5 GWei
			maxGasPriceWei:         toBigInt("5e11"), // 0.5 uEther
			expectedFee:            gas.DynamicFee{TipCap: assets.GWei(39), FeeCap: assets.GWei(5000)},
			originalLimit:          100000,
			limitMultiplierPercent: 1.1,
			expectedLimit:          110000,
		},
		{
			name:                   "original + fixed wins",
			currentTipCap:          nil,
			originalFee:            gas.DynamicFee{TipCap: assets.GWei(30), FeeCap: assets.GWei(5000)},
			tipCapDefault:          assets.GWei(20),
			bumpPercent:            20,
			bumpWei:                toBigInt("8e9"),  // 0.8 GWei
			maxGasPriceWei:         toBigInt("5e11"), // 0.5 uEther
			expectedFee:            gas.DynamicFee{TipCap: assets.GWei(38), FeeCap: assets.GWei(5000)},
			originalLimit:          100000,
			limitMultiplierPercent: 0.8,
			expectedLimit:          80000,
		},
		{
			name:                   "default + percentage wins",
			currentTipCap:          nil,
			originalFee:            gas.DynamicFee{TipCap: assets.GWei(30), FeeCap: assets.GWei(5000)},
			tipCapDefault:          assets.GWei(40),
			bumpPercent:            20,
			bumpWei:                toBigInt("5e9"),  // 0.5 GWei
			maxGasPriceWei:         toBigInt("5e11"), // 0.5 uEther
			expectedFee:            gas.DynamicFee{TipCap: assets.GWei(48), FeeCap: assets.GWei(5000)},
			originalLimit:          100000,
			limitMultiplierPercent: 1.0,
			expectedLimit:          100000,
		},
		{
			name:                   "default + fixed wins",
			currentTipCap:          assets.GWei(48),
			originalFee:            gas.DynamicFee{TipCap: assets.GWei(30), FeeCap: assets.GWei(5000)},
			tipCapDefault:          assets.GWei(40),
			bumpPercent:            20,
			bumpWei:                toBigInt("9e9"),  // 0.9 GWei
			maxGasPriceWei:         toBigInt("5e11"), // 0.5 uEther
			expectedFee:            gas.DynamicFee{TipCap: assets.GWei(49), FeeCap: assets.GWei(5000)},
			originalLimit:          100000,
			limitMultiplierPercent: 1.0,
			expectedLimit:          100000,
		},
		{
			name:                   "higher current tip cap wins",
			currentTipCap:          assets.GWei(50),
			originalFee:            gas.DynamicFee{TipCap: assets.GWei(30), FeeCap: assets.GWei(5000)},
			tipCapDefault:          assets.GWei(40),
			bumpPercent:            20,
			bumpWei:                toBigInt("9e9"),  // 0.9 GWei
			maxGasPriceWei:         toBigInt("5e11"), // 0.5 uEther
			expectedFee:            gas.DynamicFee{TipCap: assets.GWei(50), FeeCap: assets.GWei(5000)},
			originalLimit:          100000,
			limitMultiplierPercent: 1.0,
			expectedLimit:          100000,
		},
		{
			name:                   "max increased uses new higher max for FeeCap",
			currentTipCap:          nil,
			originalFee:            gas.DynamicFee{TipCap: assets.GWei(30), FeeCap: assets.GWei(5000)},
			tipCapDefault:          assets.GWei(20),
			bumpPercent:            20,
			bumpWei:                toBigInt("5e9"), // 0.5 GWei
			maxGasPriceWei:         assets.GWei(8000),
			expectedFee:            gas.DynamicFee{TipCap: assets.GWei(36), FeeCap: assets.GWei(8000)},
			originalLimit:          100000,
			limitMultiplierPercent: 1.0,
			expectedLimit:          100000,
		},
		{
			name:                   "max decreased uses previous higher max for FeeCap",
			currentTipCap:          nil,
			originalFee:            gas.DynamicFee{TipCap: assets.GWei(30), FeeCap: assets.GWei(5000)},
			tipCapDefault:          assets.GWei(20),
			bumpPercent:            20,
			bumpWei:                toBigInt("5e9"), // 0.5 GWei
			maxGasPriceWei:         assets.GWei(3000),
			expectedFee:            gas.DynamicFee{TipCap: assets.GWei(36), FeeCap: assets.GWei(5000)},
			originalLimit:          100000,
			limitMultiplierPercent: 1.0,
			expectedLimit:          100000,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			cfg := new(gasmocks.Config)
			cfg.Test(t)
			cfg.On("EvmGasBumpPercent").Return(test.bumpPercent)
			cfg.On("EvmGasTipCapDefault").Return(test.tipCapDefault)
			cfg.On("EvmGasBumpWei").Return(test.bumpWei)
			cfg.On("EvmMaxGasPriceWei").Return(test.maxGasPriceWei)
			cfg.On("EvmGasLimitMultiplier").Return(test.limitMultiplierPercent)
			actual, limit, err := gas.BumpDynamicFeeOnly(cfg, logger.TestLogger(t), test.currentTipCap, test.originalFee, test.originalLimit)
			require.NoError(t, err)
			if actual.TipCap.Cmp(test.expectedFee.TipCap) != 0 {
				t.Fatalf("TipCap not equal, expected %s but got %s", test.expectedFee.TipCap.String(), actual.TipCap.String())
			}
			if actual.FeeCap.Cmp(test.expectedFee.FeeCap) != 0 {
				t.Fatalf("FeeCap not equal, expected %s but got %s", test.expectedFee.FeeCap.String(), actual.FeeCap.String())
			}
			assert.Equal(t, int(test.expectedLimit), int(limit))
			cfg.AssertExpectations(t)
		})
	}
}

func Test_BumpDynamicFeeOnly_HitsMaxError(t *testing.T) {
	t.Parallel()
	cfg := new(gasmocks.Config)
	cfg.Test(t)
	cfg.On("EvmGasBumpPercent").Return(uint16(50))
	cfg.On("EvmGasTipCapDefault").Return(assets.GWei(0))
	cfg.On("EvmGasBumpWei").Return(assets.Wei(5000000000))
	cfg.On("EvmMaxGasPriceWei").Return(assets.GWei(40))

	originalFee := gas.DynamicFee{TipCap: assets.GWei(30), FeeCap: assets.GWei(100)}
	_, _, err := gas.BumpDynamicFeeOnly(cfg, logger.TestLogger(t), nil, originalFee, 42)
	require.Error(t, err)
	require.Contains(t, err.Error(), "bumped tip cap of 45000000000 would exceed configured max gas price of 40000000000 (original fee: tip cap 30000000000, fee cap 100000000000)")
}

// toBigInt is used to convert scientific notation string to a *big.Int
func toBigInt(input string) *big.Int {
	flt, _, err := big.ParseFloat(input, 10, 0, big.ToNearestEven)
	if err != nil {
		panic(fmt.Sprintf("unable to parse '%s' into a big.Float: %v", input, err))
	}
	var i = new(big.Int)
	i, _ = flt.Int(i)
	return i
}
