package gas_test

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/smartcontractkit/chainlink/core/services/gas"
	"github.com/smartcontractkit/chainlink/core/store/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_BumpGasPriceOnly(t *testing.T) {
	t.Parallel()

	for _, test := range []struct {
		name                   string
		originalGasPrice       *big.Int
		priceDefault           *big.Int
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
			originalGasPrice:       toBigInt("3e10"), // 30 GWei
			priceDefault:           toBigInt("2e10"), // 20 GWei
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
			originalGasPrice:       toBigInt("3e10"), // 30 GWei
			priceDefault:           toBigInt("2e10"), // 20 GWei
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
			originalGasPrice:       toBigInt("3e10"), // 30 GWei
			priceDefault:           toBigInt("2e10"), // 20 GWei
			bumpPercent:            20,
			bumpWei:                toBigInt("8e9"),    // 0.8 GWei
			maxGasPriceWei:         toBigInt("5e11"),   // 0.5 uEther
			expectedGasPrice:       toBigInt("3.8e10"), // 38 GWei
			originalLimit:          100000,
			limitMultiplierPercent: 0.8,
			expectedLimit:          80000,
		},
		{
			name:                   "default + percentage wins",
			originalGasPrice:       toBigInt("3e10"), // 30 GWei
			priceDefault:           toBigInt("4e10"), // 40 GWei
			bumpPercent:            20,
			bumpWei:                toBigInt("5e9"),    // 0.5 GWei
			maxGasPriceWei:         toBigInt("5e11"),   // 0.5 uEther
			expectedGasPrice:       toBigInt("4.8e10"), // 48 GWei
			originalLimit:          100000,
			limitMultiplierPercent: 1.0,
			expectedLimit:          100000,
		},
		{
			name:                   "default + fixed wins",
			originalGasPrice:       toBigInt("3e10"), // 30 GWei
			priceDefault:           toBigInt("4e10"), // 40 GWei
			bumpPercent:            20,
			bumpWei:                toBigInt("9e9"),    // 0.9 GWei
			maxGasPriceWei:         toBigInt("5e11"),   // 0.5 uEther
			expectedGasPrice:       toBigInt("4.9e10"), // 49 GWei
			originalLimit:          100000,
			limitMultiplierPercent: 1.0,
			expectedLimit:          100000,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			cfg := config.NewConfig()
			cfg.Set("ETH_GAS_PRICE_DEFAULT", test.priceDefault)
			cfg.Set("ETH_GAS_BUMP_PERCENT", test.bumpPercent)
			cfg.Set("ETH_GAS_BUMP_WEI", test.bumpWei)
			cfg.Set("ETH_MAX_GAS_PRICE_WEI", test.maxGasPriceWei)
			cfg.Set("ETH_GAS_LIMIT_MULTIPLIER", test.limitMultiplierPercent)
			actual, limit, err := gas.BumpGasPriceOnly(cfg, test.originalGasPrice, test.originalLimit)
			require.NoError(t, err)
			if actual.Cmp(test.expectedGasPrice) != 0 {
				t.Fatalf("Expected %s but got %s", test.expectedGasPrice.String(), actual.String())
			}
			assert.Equal(t, int(test.expectedLimit), int(limit))
		})
	}
}

func Test_BumpGasPriceOnly_HitsMaxError(t *testing.T) {
	t.Parallel()
	cfg := config.NewConfig()
	cfg.Set("ETH_GAS_BUMP_PERCENT", "50")
	cfg.Set("ETH_GAS_PRICE_DEFAULT", toBigInt("2e10")) // 20 GWei
	cfg.Set("ETH_GAS_BUMP_WEI", toBigInt("5e9"))       // 0.5 GWei
	cfg.Set("ETH_MAX_GAS_PRICE_WEI", toBigInt("4e10")) // 40 Gwei

	originalGasPrice := toBigInt("3e10") // 30 GWei
	_, _, err := gas.BumpGasPriceOnly(cfg, originalGasPrice, 42)
	require.Error(t, err)
	require.Contains(t, err.Error(), "bumped gas price of 45000000000 would exceed configured max gas price of 40000000000 (original price was 30000000000)")
}

func Test_BumpGasPriceOnly_NoBumpError(t *testing.T) {
	t.Parallel()
	cfg := config.NewConfig()
	cfg.Set("ETH_GAS_BUMP_PERCENT", "0")
	cfg.Set("ETH_GAS_BUMP_WEI", "0")
	cfg.Set("ETH_MAX_GAS_PRICE_WEI", "40000000000")

	originalGasPrice := toBigInt("3e10") // 30 GWei
	_, _, err := gas.BumpGasPriceOnly(cfg, originalGasPrice, 42)
	require.Error(t, err)
	require.Contains(t, err.Error(), "bumped gas price of 30000000000 is equal to original gas price of 30000000000. ACTION REQUIRED: This is a configuration error, you must increase either ETH_GAS_BUMP_PERCENT or ETH_GAS_BUMP_WEI")

	// Even if it's exactly the maximum
	originalGasPrice = toBigInt("4e10") // 40 GWei
	_, _, err = gas.BumpGasPriceOnly(cfg, originalGasPrice, 42)
	require.Error(t, err)
	require.Contains(t, err.Error(), "bumped gas price of 40000000000 is equal to original gas price of 40000000000. ACTION REQUIRED: This is a configuration error, you must increase either ETH_GAS_BUMP_PERCENT or ETH_GAS_BUMP_WEI")
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
