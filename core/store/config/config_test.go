package config_test

import (
	"math/big"
	"os"
	"testing"

	"github.com/smartcontractkit/chainlink/core/assets"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/store/config"
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

func TestConfig_Profiles(t *testing.T) {
	config := config.NewConfig()

	tests := []struct {
		name                           string
		chainID                        string
		expectedGasLimitDefault        uint64
		expectedMinimumContractPayment int64
	}{
		{"default", "", 500000, 1000000000000000000},
		{"mainnet", "1", 500000, 1000000000000000000},
		{"kovan", "42", 500000, 1000000000000000000},

		{"optimism", "10", 500000, 100000000000000},
		{"optimism", "69", 500000, 100000000000000},
		{"optimism", "420", 500000, 100000000000000},

		{"bscMainnet", "56", 500000, 100000000000000},
		{"hecoMainnet", "128", 500000, 100000000000000},
		{"fantomMainnet", "250", 500000, 100000000000000},
		{"fantomTestnet", "4002", 500000, 100000000000000},
		{"polygonMatic", "800001", 500000, 100000000000000},

		{"xDai", "100", 500000, 100000000000000},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config.Set("ETH_CHAIN_ID", tt.chainID)
			assert.Equal(t, tt.expectedGasLimitDefault, config.EthGasLimitDefault())
			assert.Equal(t, assets.NewLink(tt.expectedMinimumContractPayment), config.MinimumContractPayment())
		})
	}
}

func TestConfig_MinimumContractPayment(t *testing.T) {
	originalJuels := os.Getenv("MINIMUM_CONTRACT_PAYMENT_LINK_JUELS")
	originalLink := os.Getenv("MINIMUM_CONTRACT_PAYMENT")
	defer func() {
		os.Setenv("MINIMUM_CONTRACT_PAYMENT_LINK_JUELS", originalJuels)
		os.Setenv("MINIMUM_CONTRACT_PAYMENT", originalLink)
	}()

	cfg := config.NewConfig()
	assert.Equal(t, assets.NewLink(1000000000000000000), cfg.MinimumContractPayment())

	os.Setenv("MINIMUM_CONTRACT_PAYMENT_LINK_JUELS", "5987")
	cfg = config.NewConfig()
	assert.Equal(t, assets.NewLink(5987), cfg.MinimumContractPayment())

	os.Setenv("MINIMUM_CONTRACT_PAYMENT", "4937")
	cfg = config.NewConfig()
	assert.Equal(t, assets.NewLink(5987), cfg.MinimumContractPayment())

	os.Setenv("MINIMUM_CONTRACT_PAYMENT_LINK_JUELS", "")
	cfg = config.NewConfig()
	assert.Equal(t, assets.NewLink(4937), cfg.MinimumContractPayment())
}
