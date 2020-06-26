package bulletprooftxmanager_test

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/smartcontractkit/chainlink/core/services/bulletprooftxmanager"
	"github.com/smartcontractkit/chainlink/core/store/orm"
)

func TestBulletproofTxManager_BumpGas(t *testing.T) {
	t.Parallel()

	for _, test := range []struct {
		name             string
		originalGasPrice *big.Int
		priceDefault     *big.Int
		bumpPercent      uint16
		bumpWei          *big.Int
		maxGasPriceWei   *big.Int
		expected         *big.Int
	}{
		{
			name:             "defaults",
			originalGasPrice: toBigInt("3e10"), // 30 GWei
			priceDefault:     toBigInt("2e10"), // 20 GWei
			bumpPercent:      20,
			bumpWei:          toBigInt("5e9"),    // 0.5 GWei
			maxGasPriceWei:   toBigInt("5e11"),   // 0.5 uEther
			expected:         toBigInt("3.6e10"), // 36 GWei
		},
		{
			name:             "original + percentage wins",
			originalGasPrice: toBigInt("3e10"), // 30 GWei
			priceDefault:     toBigInt("2e10"), // 20 GWei
			bumpPercent:      30,
			bumpWei:          toBigInt("5e9"),    // 0.5 GWei
			maxGasPriceWei:   toBigInt("5e11"),   // 0.5 uEther
			expected:         toBigInt("3.9e10"), // 39 GWei
		},
		{
			name:             "original + fixed wins",
			originalGasPrice: toBigInt("3e10"), // 30 GWei
			priceDefault:     toBigInt("2e10"), // 20 GWei
			bumpPercent:      20,
			bumpWei:          toBigInt("8e9"),    // 0.8 GWei
			maxGasPriceWei:   toBigInt("5e11"),   // 0.5 uEther
			expected:         toBigInt("3.8e10"), // 38 GWei
		},
		{
			name:             "default + percentage wins",
			originalGasPrice: toBigInt("3e10"), // 30 GWei
			priceDefault:     toBigInt("4e10"), // 40 GWei
			bumpPercent:      20,
			bumpWei:          toBigInt("5e9"),    // 0.5 GWei
			maxGasPriceWei:   toBigInt("5e11"),   // 0.5 uEther
			expected:         toBigInt("4.8e10"), // 48 GWei
		},
		{
			name:             "default + fixed wins",
			originalGasPrice: toBigInt("3e10"), // 30 GWei
			priceDefault:     toBigInt("4e10"), // 40 GWei
			bumpPercent:      20,
			bumpWei:          toBigInt("9e9"),    // 0.9 GWei
			maxGasPriceWei:   toBigInt("5e11"),   // 0.5 uEther
			expected:         toBigInt("4.9e10"), // 49 GWei
		},
		{
			name:             "max wins",
			originalGasPrice: toBigInt("3e10"), // 30 GWei
			priceDefault:     toBigInt("2e10"), // 20 GWei
			bumpPercent:      50,
			bumpWei:          toBigInt("5e9"),  // 0.5 GWei
			maxGasPriceWei:   toBigInt("4e10"), // 40 GWei
			expected:         toBigInt("4e10"), // 40 GWei
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			config := orm.NewConfig()
			config.Set("ETH_GAS_PRICE_DEFAULT", test.priceDefault)
			config.Set("ETH_GAS_BUMP_PERCENT", test.bumpPercent)
			config.Set("ETH_GAS_BUMP_WEI", test.bumpWei)
			config.Set("ETH_MAX_GAS_PRICE_WEI", test.maxGasPriceWei)
			actual := bulletprooftxmanager.BumpGas(config, test.originalGasPrice)
			if actual.Cmp(test.expected) != 0 {
				t.Fatalf("Expected %s but got %s", test.expected.String(), actual.String())
			}
		})
	}
}

// Helpers

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
