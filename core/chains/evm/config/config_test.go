package config_test

import (
	"fmt"
	"math/big"
	"math/rand"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/guregu/null.v4"

	"github.com/smartcontractkit/chainlink/core/chains"
	evmconfig "github.com/smartcontractkit/chainlink/core/chains/evm/config"
	evmmocks "github.com/smartcontractkit/chainlink/core/chains/evm/mocks"
	evmtypes "github.com/smartcontractkit/chainlink/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/configtest"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/utils"
)

func TestChainScopedConfig(t *testing.T) {
	orm := new(evmmocks.ORM)
	orm.Test(t)
	chainID := big.NewInt(rand.Int63())
	gcfg := configtest.NewTestGeneralConfig(t)
	lggr := logger.TestLogger(t).With("evmChainID", chainID.String())
	cfg := evmconfig.NewChainScopedConfig(chainID, evmtypes.ChainCfg{
		KeySpecific: make(map[string]evmtypes.ChainCfg),
	}, orm, lggr, gcfg)

	t.Run("EvmGasPriceDefault", func(t *testing.T) {
		t.Run("sets the gas price", func(t *testing.T) {
			assert.Equal(t, big.NewInt(20000000000), cfg.EvmGasPriceDefault())

			orm.On("StoreString", chainID, "EvmGasPriceDefault", "42000000000").Return(nil)
			err := cfg.SetEvmGasPriceDefault(big.NewInt(42000000000))
			assert.NoError(t, err)

			assert.Equal(t, big.NewInt(42000000000), cfg.EvmGasPriceDefault())

			orm.AssertExpectations(t)
		})
		t.Run("is not allowed to set gas price to below EvmMinGasPriceWei", func(t *testing.T) {
			assert.Equal(t, big.NewInt(1000000000), cfg.EvmMinGasPriceWei())

			err := cfg.SetEvmGasPriceDefault(big.NewInt(1))
			assert.EqualError(t, err, "cannot set default gas price to 1, it is below the minimum allowed value of 1000000000")

			assert.Equal(t, big.NewInt(42000000000), cfg.EvmGasPriceDefault())
		})
		t.Run("is not allowed to set gas price to above EvmMaxGasPriceWei", func(t *testing.T) {
			assert.Equal(t, big.NewInt(5000000000000), cfg.EvmMaxGasPriceWei())

			err := cfg.SetEvmGasPriceDefault(big.NewInt(999999999999999))
			assert.EqualError(t, err, "cannot set default gas price to 999999999999999, it is above the maximum allowed value of 5000000000000")

			assert.Equal(t, big.NewInt(42000000000), cfg.EvmGasPriceDefault())
		})
	})

	t.Run("KeySpecificMaxGasPriceWei", func(t *testing.T) {
		addr := cltest.NewAddress()
		randomOtherAddr := cltest.NewAddress()
		randomOtherKeySpecific := evmtypes.ChainCfg{EvmMaxGasPriceWei: utils.NewBigI(rand.Int63())}
		evmconfig.PersistedCfgPtr(cfg).KeySpecific[randomOtherAddr.Hex()] = randomOtherKeySpecific

		t.Run("uses chain-specific default value when nothing is set", func(t *testing.T) {
			assert.Equal(t, big.NewInt(5000000000000), cfg.KeySpecificMaxGasPriceWei(addr))
		})

		t.Run("uses chain-specific override value when that is set", func(t *testing.T) {
			val := utils.NewBigI(rand.Int63())
			evmconfig.PersistedCfgPtr(cfg).EvmMaxGasPriceWei = val

			assert.Equal(t, val.String(), cfg.KeySpecificMaxGasPriceWei(addr).String())
		})
		t.Run("uses key-specific override value when that is set", func(t *testing.T) {
			val := utils.NewBigI(rand.Int63())
			keySpecific := evmtypes.ChainCfg{EvmMaxGasPriceWei: val}
			evmconfig.PersistedCfgPtr(cfg).KeySpecific[addr.Hex()] = keySpecific

			assert.Equal(t, val.String(), cfg.KeySpecificMaxGasPriceWei(addr).String())
		})
		t.Run("uses global value when that is set", func(t *testing.T) {
			val := big.NewInt(rand.Int63())
			gcfg.Overrides.GlobalEvmMaxGasPriceWei = val

			assert.Equal(t, val.String(), cfg.KeySpecificMaxGasPriceWei(addr).String())
		})
	})
}

func TestChainScopedConfig_Profiles(t *testing.T) {
	tests := []struct {
		name                           string
		chainID                        int64
		expectedGasLimitDefault        uint64
		expectedMinimumContractPayment string
	}{
		{"default", 0, 500000, "0.00001"},
		{"mainnet", 1, 500000, "0.1"},
		{"kovan", 42, 500000, "0.1"},

		{"optimism", 10, 500000, "0.00001"},
		{"optimism", 69, 500000, "0.00001"},
		{"optimism", 420, 500000, "0.00001"},

		{"bscMainnet", 56, 500000, "0.00001"},
		{"hecoMainnet", 128, 500000, "0.00001"},
		{"fantomMainnet", 250, 500000, "0.00001"},
		{"fantomTestnet", 4002, 500000, "0.00001"},
		{"polygonMatic", 800001, 500000, "0.00001"},
		{"harmonyMainnet", 1666600000, 500000, "0.00001"},
		{"harmonyTestnet", 1666700000, 500000, "0.00001"},

		{"xDai", 100, 500000, "0.00001"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gcfg := configtest.NewTestGeneralConfig(t)
			lggr := logger.TestLogger(t)
			config := evmconfig.NewChainScopedConfig(big.NewInt(tt.chainID), evmtypes.ChainCfg{}, nil, lggr, gcfg)

			assert.Equal(t, tt.expectedGasLimitDefault, config.EvmGasLimitDefault())
			assert.Equal(t, tt.expectedMinimumContractPayment, strings.TrimRight(config.MinimumContractPayment().Link(), "0"))
		})
	}
}

func Test_chainScopedConfig_Validate(t *testing.T) {
	// Validate built-in
	for id := range evmconfig.ChainSpecificConfigDefaultSets() {
		id := id
		t.Run(fmt.Sprintf("chainID-%d", id), func(t *testing.T) {
			gcfg := cltest.NewTestGeneralConfig(t)
			lggr := logger.TestLogger(t)
			cfg := evmconfig.NewChainScopedConfig(big.NewInt(id), evmtypes.ChainCfg{}, nil, lggr, gcfg)
			assert.NoError(t, cfg.Validate())
		})
	}

	// Invalid Cases:

	t.Run("arbitrum-estimator", func(t *testing.T) {
		t.Run("custom", func(t *testing.T) {
			gcfg := cltest.NewTestGeneralConfig(t)
			lggr := logger.TestLogger(t)
			cfg := evmconfig.NewChainScopedConfig(big.NewInt(0), evmtypes.ChainCfg{
				ChainType:        null.StringFrom(string(chains.Arbitrum)),
				GasEstimatorMode: null.StringFrom("BlockHistory"),
			}, nil, lggr, gcfg)
			assert.Error(t, cfg.Validate())
		})
		t.Run("mainnet", func(t *testing.T) {
			gcfg := cltest.NewTestGeneralConfig(t)
			lggr := logger.TestLogger(t)
			cfg := evmconfig.NewChainScopedConfig(big.NewInt(42161), evmtypes.ChainCfg{
				GasEstimatorMode: null.StringFrom("BlockHistory"),
			}, nil, lggr, gcfg)
			assert.Error(t, cfg.Validate())
		})
		t.Run("testnet", func(t *testing.T) {
			gcfg := cltest.NewTestGeneralConfig(t)
			lggr := logger.TestLogger(t)
			cfg := evmconfig.NewChainScopedConfig(big.NewInt(421611), evmtypes.ChainCfg{
				GasEstimatorMode: null.StringFrom("Optimism"),
			}, nil, lggr, gcfg)
			assert.Error(t, cfg.Validate())
		})
	})

	t.Run("optimism-estimator", func(t *testing.T) {
		t.Run("custom", func(t *testing.T) {
			gcfg := cltest.NewTestGeneralConfig(t)
			lggr := logger.TestLogger(t)
			cfg := evmconfig.NewChainScopedConfig(big.NewInt(0), evmtypes.ChainCfg{
				ChainType:        null.StringFrom(string(chains.Optimism)),
				GasEstimatorMode: null.StringFrom("BlockHistory"),
			}, nil, lggr, gcfg)
			assert.Error(t, cfg.Validate())
		})
		t.Run("mainnet", func(t *testing.T) {
			gcfg := cltest.NewTestGeneralConfig(t)
			lggr := logger.TestLogger(t)
			cfg := evmconfig.NewChainScopedConfig(big.NewInt(10), evmtypes.ChainCfg{
				GasEstimatorMode: null.StringFrom("FixedPrice"),
			}, nil, lggr, gcfg)
			assert.Error(t, cfg.Validate())
		})
		t.Run("testnet", func(t *testing.T) {
			gcfg := cltest.NewTestGeneralConfig(t)
			lggr := logger.TestLogger(t)
			cfg := evmconfig.NewChainScopedConfig(big.NewInt(69), evmtypes.ChainCfg{
				GasEstimatorMode: null.StringFrom("BlockHistory"),
			}, nil, lggr, gcfg)
			assert.Error(t, cfg.Validate())
		})
	})
}
