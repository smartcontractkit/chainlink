package config_test

import (
	"math/big"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/smartcontractkit/chainlink/core/assets"
	evmconfig "github.com/smartcontractkit/chainlink/core/chains/evm/config"
	evmmocks "github.com/smartcontractkit/chainlink/core/chains/evm/mocks"
	evmtypes "github.com/smartcontractkit/chainlink/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/configtest"
	config "github.com/smartcontractkit/chainlink/core/store/config"
	"github.com/smartcontractkit/chainlink/core/utils"
)

func TestChainScopedConfig_EvmGasPriceDefault(t *testing.T) {
	orm := new(evmmocks.ORM)
	orm.Test(t)
	chainID := big.NewInt(rand.Int63())
	cfg := config.NewGeneralConfig()
	config := evmconfig.NewChainScopedConfig(orm, cfg.CreateProductionLogger(), cfg, evmtypes.Chain{ID: *utils.NewBig(chainID)})

	t.Run("sets the gas price", func(t *testing.T) {
		assert.Equal(t, big.NewInt(20000000000), config.EvmGasPriceDefault())

		orm.On("StoreString", chainID, "EvmGasPriceDefault", "42000000000").Return(nil)
		err := config.SetEvmGasPriceDefault(big.NewInt(42000000000))
		assert.NoError(t, err)

		assert.Equal(t, big.NewInt(42000000000), config.EvmGasPriceDefault())

		orm.AssertExpectations(t)
	})
	t.Run("is not allowed to set gas price to below EvmMinGasPriceWei", func(t *testing.T) {
		assert.Equal(t, big.NewInt(1000000000), config.EvmMinGasPriceWei())

		err := config.SetEvmGasPriceDefault(big.NewInt(1))
		assert.EqualError(t, err, "cannot set default gas price to 1, it is below the minimum allowed value of 1000000000")

		assert.Equal(t, big.NewInt(42000000000), config.EvmGasPriceDefault())
	})
	t.Run("is not allowed to set gas price to above EvmMaxGasPriceWei", func(t *testing.T) {
		assert.Equal(t, big.NewInt(5000000000000), config.EvmMaxGasPriceWei())

		err := config.SetEvmGasPriceDefault(big.NewInt(999999999999999))
		assert.EqualError(t, err, "cannot set default gas price to 999999999999999, it is above the maximum allowed value of 5000000000000")

		assert.Equal(t, big.NewInt(42000000000), config.EvmGasPriceDefault())
	})
}

func TestChainScopedConfig_Profiles(t *testing.T) {
	tests := []struct {
		name                           string
		chainID                        int64
		expectedGasLimitDefault        uint64
		expectedMinimumContractPayment int64
	}{
		{"default", 0, 500000, 100000000000000},
		{"mainnet", 1, 500000, 1000000000000000000},
		{"kovan", 42, 500000, 1000000000000000000},

		{"optimism", 10, 500000, 100000000000000},
		{"optimism", 69, 500000, 100000000000000},
		{"optimism", 420, 500000, 100000000000000},

		{"bscMainnet", 56, 500000, 100000000000000},
		{"hecoMainnet", 128, 500000, 100000000000000},
		{"fantomMainnet", 250, 500000, 100000000000000},
		{"fantomTestnet", 4002, 500000, 100000000000000},
		{"polygonMatic", 800001, 500000, 100000000000000},

		{"xDai", 100, 500000, 100000000000000},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gcfg := configtest.NewTestGeneralConfig(t)
			config := evmconfig.NewChainScopedConfig(nil, gcfg.CreateProductionLogger(), gcfg, evmtypes.Chain{ID: *utils.NewBigI(tt.chainID)})

			assert.Equal(t, tt.expectedGasLimitDefault, config.EvmGasLimitDefault())
			assert.Equal(t, assets.NewLinkFromJuels(tt.expectedMinimumContractPayment).String(), config.MinimumContractPayment().String())
		})
	}
}
