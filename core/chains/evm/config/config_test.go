package config_test

import (
	"fmt"
	"math/big"
	"math/rand"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/core/assets"
	evmconfig "github.com/smartcontractkit/chainlink/core/chains/evm/config"
	v2 "github.com/smartcontractkit/chainlink/core/chains/evm/config/v2"
	"github.com/smartcontractkit/chainlink/core/config"
	"github.com/smartcontractkit/chainlink/core/internal/testutils"
	configtest "github.com/smartcontractkit/chainlink/core/internal/testutils/configtest/v2"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/evmtest"
	"github.com/smartcontractkit/chainlink/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/ethkey"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"
)

func TestChainScopedConfig(t *testing.T) {
	t.Parallel()
	gcfg := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		id := utils.NewBig(big.NewInt(rand.Int63()))
		c.EVM[0] = &v2.EVMConfig{
			ChainID: id,
			Chain: v2.Defaults(id, &v2.Chain{
				GasEstimator: v2.GasEstimator{PriceMax: assets.NewWeiI(100000000000000)},
			}),
		}
	})
	cfg := evmtest.NewChainScopedConfig(t, gcfg)

	overrides := func(c *chainlink.Config, s *chainlink.Secrets) {
		id := utils.NewBig(big.NewInt(rand.Int63()))
		c.EVM[0] = &v2.EVMConfig{
			ChainID: id,
			Chain: v2.Defaults(id, &v2.Chain{
				GasEstimator: v2.GasEstimator{
					PriceMax:     assets.NewWeiI(100000000000000),
					PriceDefault: assets.NewWeiI(42000000000),
				},
			}),
		}
	}
	t.Run("EvmGasPriceDefault", func(t *testing.T) {
		assert.Equal(t, assets.NewWeiI(20000000000), cfg.EvmGasPriceDefault())

		gcfg2 := configtest.NewGeneralConfig(t, overrides)
		cfg2 := evmtest.NewChainScopedConfig(t, gcfg2)
		assert.Equal(t, assets.NewWeiI(42000000000), cfg2.EvmGasPriceDefault())
	})

	t.Run("KeySpecificMaxGasPriceWei", func(t *testing.T) {
		addr := testutils.NewAddress()
		randomOtherAddr := testutils.NewAddress()
		gcfg2 := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
			overrides(c, s)
			c.EVM[0].KeySpecific = v2.KeySpecificConfig{
				{Key: ptr(ethkey.EIP55AddressFromAddress(randomOtherAddr)),
					GasEstimator: v2.KeySpecificGasEstimator{
						PriceMax: assets.GWei(850),
					},
				},
			}
		})
		cfg2 := evmtest.NewChainScopedConfig(t, gcfg2)

		t.Run("uses chain-specific default value when nothing is set", func(t *testing.T) {
			assert.Equal(t, assets.NewWeiI(100000000000000), cfg2.KeySpecificMaxGasPriceWei(addr))
		})

		t.Run("uses chain-specific override value when that is set", func(t *testing.T) {
			val := assets.NewWeiI(rand.Int63())
			gcfg3 := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
				c.EVM[0].GasEstimator.PriceMax = val
			})
			cfg3 := evmtest.NewChainScopedConfig(t, gcfg3)

			assert.Equal(t, val.String(), cfg3.KeySpecificMaxGasPriceWei(addr).String())
		})
		t.Run("uses key-specific override value when set", func(t *testing.T) {
			val := assets.GWei(250)
			gcfg3 := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
				c.EVM[0].KeySpecific = v2.KeySpecificConfig{
					{Key: ptr(ethkey.EIP55AddressFromAddress(addr)),
						GasEstimator: v2.KeySpecificGasEstimator{
							PriceMax: val,
						},
					},
				}
			})
			cfg3 := evmtest.NewChainScopedConfig(t, gcfg3)

			assert.Equal(t, val.String(), cfg3.KeySpecificMaxGasPriceWei(addr).String())
		})
		t.Run("uses key-specific override value when set and lower than chain specific config", func(t *testing.T) {
			keySpecificPrice := assets.GWei(900)
			chainSpecificPrice := assets.GWei(1200)
			gcfg3 := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
				c.EVM[0].GasEstimator.PriceMax = chainSpecificPrice
				c.EVM[0].KeySpecific = v2.KeySpecificConfig{
					{Key: ptr(ethkey.EIP55AddressFromAddress(addr)),
						GasEstimator: v2.KeySpecificGasEstimator{
							PriceMax: keySpecificPrice,
						},
					},
				}
			})
			cfg3 := evmtest.NewChainScopedConfig(t, gcfg3)

			assert.Equal(t, keySpecificPrice.String(), cfg3.KeySpecificMaxGasPriceWei(addr).String())
		})
		t.Run("uses chain-specific value when higher than key-specific value", func(t *testing.T) {
			keySpecificPrice := assets.GWei(1400)
			chainSpecificPrice := assets.GWei(1200)
			gcfg3 := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
				c.EVM[0].GasEstimator.PriceMax = chainSpecificPrice
				c.EVM[0].KeySpecific = v2.KeySpecificConfig{
					{Key: ptr(ethkey.EIP55AddressFromAddress(addr)),
						GasEstimator: v2.KeySpecificGasEstimator{
							PriceMax: keySpecificPrice,
						},
					},
				}
			})
			cfg3 := evmtest.NewChainScopedConfig(t, gcfg3)

			assert.Equal(t, chainSpecificPrice.String(), cfg3.KeySpecificMaxGasPriceWei(addr).String())
		})
		t.Run("uses key-specific override value when set and lower than global config", func(t *testing.T) {
			keySpecificPrice := assets.GWei(900)
			gcfg3 := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
				c.EVM[0].KeySpecific = v2.KeySpecificConfig{
					{Key: ptr(ethkey.EIP55AddressFromAddress(addr)),
						GasEstimator: v2.KeySpecificGasEstimator{
							PriceMax: keySpecificPrice,
						},
					},
				}
			})
			cfg3 := evmtest.NewChainScopedConfig(t, gcfg3)

			assert.Equal(t, keySpecificPrice.String(), cfg3.KeySpecificMaxGasPriceWei(addr).String())
		})
		t.Run("uses global value when higher than key-specific value", func(t *testing.T) {
			keySpecificPrice := assets.GWei(1400)
			chainSpecificPrice := assets.GWei(1200)
			gcfg3 := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
				c.EVM[0].GasEstimator.PriceMax = chainSpecificPrice
				c.EVM[0].KeySpecific = v2.KeySpecificConfig{
					{Key: ptr(ethkey.EIP55AddressFromAddress(addr)),
						GasEstimator: v2.KeySpecificGasEstimator{
							PriceMax: keySpecificPrice,
						},
					},
				}
			})
			cfg3 := evmtest.NewChainScopedConfig(t, gcfg3)

			assert.Equal(t, chainSpecificPrice.String(), cfg3.KeySpecificMaxGasPriceWei(addr).String())
		})
		t.Run("uses global value when there is no key-specific price", func(t *testing.T) {
			val := assets.NewWeiI(rand.Int63())
			unsetAddr := testutils.NewAddress()
			gcfg3 := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
				c.EVM[0].GasEstimator.PriceMax = val
			})
			cfg3 := evmtest.NewChainScopedConfig(t, gcfg3)

			assert.Equal(t, val.String(), cfg3.KeySpecificMaxGasPriceWei(unsetAddr).String())
		})
	})

	t.Run("LinkContractAddress", func(t *testing.T) {
		t.Run("uses chain-specific default value when nothing is set", func(t *testing.T) {
			assert.Equal(t, "", cfg.LinkContractAddress())
		})

		t.Run("uses chain-specific override value when that is set", func(t *testing.T) {
			val := testutils.NewAddress()

			gcfg3 := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
				c.EVM[0].LinkContractAddress = ptr(ethkey.EIP55AddressFromAddress(val))
			})
			cfg3 := evmtest.NewChainScopedConfig(t, gcfg3)

			assert.Equal(t, val.String(), cfg3.LinkContractAddress())
		})
	})

	t.Run("OperatorFactoryAddress", func(t *testing.T) {
		t.Run("uses chain-specific default value when nothing is set", func(t *testing.T) {
			assert.Equal(t, "", cfg.OperatorFactoryAddress())
		})

		t.Run("uses chain-specific override value when that is set", func(t *testing.T) {
			val := testutils.NewAddress()

			gcfg3 := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
				c.EVM[0].OperatorFactoryAddress = ptr(ethkey.EIP55AddressFromAddress(val))
			})
			cfg3 := evmtest.NewChainScopedConfig(t, gcfg3)

			assert.Equal(t, val.String(), cfg3.OperatorFactoryAddress())
		})
	})
}

func TestChainScopedConfig_BSCDefaults(t *testing.T) {
	chainID := big.NewInt(56)
	gcfg := configtest.NewGeneralConfig(t, func(c *chainlink.Config, secrets *chainlink.Secrets) {
		id := utils.NewBig(chainID)
		cfg := v2.Defaults(id)
		c.EVM[0] = &v2.EVMConfig{
			ChainID: id,
			Enabled: ptr(true),
			Chain:   cfg,
		}
	})
	cfg := evmtest.NewChainScopedConfig(t, gcfg)

	timeout := cfg.OCRDatabaseTimeout()
	require.Equal(t, 2*time.Second, timeout)
	timeout = cfg.OCRContractTransmitterTransmitTimeout()
	require.Equal(t, 2*time.Second, timeout)
	timeout = cfg.OCRObservationGracePeriod()
	require.Equal(t, 500*time.Millisecond, timeout)
}

func TestChainScopedConfig_Profiles(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name                           string
		chainID                        int64
		expectedGasLimitDefault        uint32
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
	for _, test := range tests {
		tt := test

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			gcfg := configtest.NewGeneralConfig(t, func(c *chainlink.Config, secrets *chainlink.Secrets) {
				id := utils.NewBigI(tt.chainID)
				cfg := v2.Defaults(id)
				c.EVM[0] = &v2.EVMConfig{
					ChainID: id,
					Enabled: ptr(true),
					Chain:   cfg,
				}
			})
			config := evmtest.NewChainScopedConfig(t, gcfg)

			assert.Equal(t, tt.expectedGasLimitDefault, config.EvmGasLimitDefault())
			assert.Nil(t, config.EvmGasLimitOCRJobType())
			assert.Nil(t, config.EvmGasLimitDRJobType())
			assert.Nil(t, config.EvmGasLimitVRFJobType())
			assert.Nil(t, config.EvmGasLimitFMJobType())
			assert.Nil(t, config.EvmGasLimitKeeperJobType())
			assert.Equal(t, tt.expectedMinimumContractPayment, strings.TrimRight(config.MinimumContractPayment().Link(), "0"))
		})
	}
}

func Test_chainScopedConfig_Validate(t *testing.T) {
	configWithChains := func(t *testing.T, id int64, chains ...*v2.Chain) config.GeneralConfig {
		return configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
			chainID := utils.NewBigI(id)
			c.EVM[0] = &v2.EVMConfig{ChainID: chainID, Enabled: ptr(true), Chain: v2.Defaults(chainID, chains...),
				Nodes: v2.EVMNodes{{
					Name:    ptr("fake"),
					WSURL:   models.MustParseURL("wss://foo.test/ws"),
					HTTPURL: models.MustParseURL("http://foo.test"),
				}}}
		})
	}

	// Validate built-in
	for id := range evmconfig.ChainSpecificConfigDefaultSets() {
		id := id
		t.Run(fmt.Sprintf("chainID-%d", id), func(t *testing.T) {
			cfg := configWithChains(t, id)
			assert.NoError(t, cfg.Validate())
		})
	}

	// Invalid Cases:

	t.Run("arbitrum-estimator", func(t *testing.T) {
		t.Run("custom", func(t *testing.T) {
			cfg := configWithChains(t, 0, &v2.Chain{
				ChainType: ptr(string(config.ChainArbitrum)),
				GasEstimator: v2.GasEstimator{
					Mode: ptr("BlockHistory"),
				},
			})
			assert.NoError(t, cfg.Validate())
		})
		t.Run("mainnet", func(t *testing.T) {
			cfg := configWithChains(t, 42161, &v2.Chain{
				GasEstimator: v2.GasEstimator{
					Mode: ptr("BlockHistory"),
					BlockHistory: v2.BlockHistoryEstimator{
						BlockHistorySize: ptr[uint16](1),
					},
				},
			})
			assert.NoError(t, cfg.Validate())
		})
		t.Run("testnet", func(t *testing.T) {
			cfg := configWithChains(t, 421611, &v2.Chain{
				GasEstimator: v2.GasEstimator{
					Mode: ptr("L2Suggested"),
				},
			})
			assert.NoError(t, cfg.Validate())
		})
	})

	t.Run("optimism-estimator", func(t *testing.T) {
		t.Run("custom", func(t *testing.T) {
			cfg := configWithChains(t, 0, &v2.Chain{
				ChainType: ptr(string(config.ChainOptimism)),
				GasEstimator: v2.GasEstimator{
					Mode: ptr("BlockHistory"),
				},
			})
			assert.Error(t, cfg.Validate())
		})
		t.Run("mainnet", func(t *testing.T) {
			cfg := configWithChains(t, 10, &v2.Chain{
				GasEstimator: v2.GasEstimator{
					Mode: ptr("FixedPrice"),
				},
			})
			assert.Error(t, cfg.Validate())
		})
		t.Run("testnet", func(t *testing.T) {
			cfg := configWithChains(t, 69, &v2.Chain{
				GasEstimator: v2.GasEstimator{
					Mode: ptr("BlockHistory"),
				},
			})
			assert.Error(t, cfg.Validate())
		})
	})
}

func ptr[T any](t T) *T { return &t }
