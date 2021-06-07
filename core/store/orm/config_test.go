package orm_test

import (
	"math/big"
	"testing"

	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/stretchr/testify/assert"
)

func TestConfig_SetEthGasPriceDefault(t *testing.T) {
	store, cleanup := cltest.NewStore(t)
	t.Cleanup(cleanup)
	config := store.Config

	config.Set("ETH_MAX_GAS_PRICE_WEI", 1500000000000)

	t.Run("sets the gas price", func(t *testing.T) {
		assert.Equal(t, big.NewInt(20000000000), config.EthGasPriceDefault())

		err := config.SetEthGasPriceDefault(big.NewInt(42000000000))
		assert.NoError(t, err)

		assert.Equal(t, big.NewInt(42000000000), config.EthGasPriceDefault())
	})
	t.Run("is not allowed to set gas price to below EthMinGasPriceWei", func(t *testing.T) {
		assert.Equal(t, big.NewInt(1000000000), config.EthMinGasPriceWei())

		err := config.SetEthGasPriceDefault(big.NewInt(1))
		assert.EqualError(t, err, "cannot set default gas price to 1, it is below the minimum allowed value of 1000000000")

		assert.Equal(t, big.NewInt(42000000000), config.EthGasPriceDefault())
	})
	t.Run("is not allowed to set gas price to above EthMaxGasPriceWei", func(t *testing.T) {
		assert.Equal(t, big.NewInt(1500000000000), config.EthMaxGasPriceWei())

		err := config.SetEthGasPriceDefault(big.NewInt(999999999999999))
		assert.EqualError(t, err, "cannot set default gas price to 999999999999999, it is above the maximum allowed value of 1500000000000")

		assert.Equal(t, big.NewInt(42000000000), config.EthGasPriceDefault())
	})
}

func TestConfig_EthGasLimitDefault_Overrides(t *testing.T) {
	store, cleanup := cltest.NewStore(t)
	t.Cleanup(cleanup)
	config := store.Config

	t.Run("returns the default", func(t *testing.T) {
		assert.Equal(t, uint64(500000), config.EthGasLimitDefault())
	})
	t.Run("for fantom returns ", func(t *testing.T) {
		config.Set("ETH_CHAIN_ID", 250)
		assert.Equal(t, uint64(500000), config.EthGasLimitDefault())

		config.Set("ETH_CHAIN_ID", 4002)
		assert.Equal(t, uint64(500000), config.EthGasLimitDefault())
	})
	t.Run("allows an override", func(t *testing.T) {
		config.Set("ETH_GAS_LIMIT_DEFAULT", 9)
		assert.Equal(t, uint64(9), config.EthGasLimitDefault())
	})
	t.Run("allows an override on fantom", func(t *testing.T) {
		config.Set("ETH_CHAIN_ID", 250)
		config.Set("ETH_GAS_LIMIT_DEFAULT", 9)
		assert.Equal(t, uint64(9), config.EthGasLimitDefault())

		config.Set("ETH_CHAIN_ID", 4002)
		assert.Equal(t, uint64(9), config.EthGasLimitDefault())
	})
}

func TestConfig_EthGasLimitDefault_AllNetworks(t *testing.T) {
	store, cleanup := cltest.NewStore(t)
	t.Cleanup(cleanup)
	config := store.Config

	tests := []struct {
		name          string
		chainID       string
		expectedValue uint64
	}{
		{"default", "", 500000},
		{"mainnet", "1", 500000},
		{"kovan", "42", 500000},

		{"optimism", "10", 500000},
		{"optimism", "69", 500000},
		{"optimism", "420", 500000},

		{"bscMainnet", "56", 500000},
		{"hecoMainnet", "128", 500000},
		{"fantomMainnet", "250", 500000},
		{"fantomTestnet", "4002", 500000},
		{"polygonMatic", "800001", 500000},

		{"xDai", "100", 500000},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config.Set("ETH_CHAIN_ID", tt.chainID)
			assert.Equal(t, tt.expectedValue, config.EthGasLimitDefault())
		})
	}
}
