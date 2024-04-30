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

	configurl "github.com/smartcontractkit/chainlink-common/pkg/config"

	commonconfig "github.com/smartcontractkit/chainlink/v2/common/config"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/assets"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/config/toml"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	ubig "github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils/big"
	"github.com/smartcontractkit/chainlink/v2/core/config"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/configtest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/evmtest"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/store/models"
)

func TestChainScopedConfig(t *testing.T) {
	t.Parallel()
	gcfg := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		id := ubig.New(big.NewInt(rand.Int63()))
		c.EVM[0] = &toml.EVMConfig{
			ChainID: id,
			Chain: toml.Defaults(id, &toml.Chain{
				GasEstimator: toml.GasEstimator{PriceMax: assets.NewWeiI(100000000000000)},
			}),
		}
	})
	cfg := evmtest.NewChainScopedConfig(t, gcfg)

	overrides := func(c *chainlink.Config, s *chainlink.Secrets) {
		id := ubig.New(big.NewInt(rand.Int63()))
		c.EVM[0] = &toml.EVMConfig{
			ChainID: id,
			Chain: toml.Defaults(id, &toml.Chain{
				GasEstimator: toml.GasEstimator{
					PriceMax:     assets.NewWeiI(100000000000000),
					PriceDefault: assets.NewWeiI(42000000000),
				},
			}),
		}
	}
	t.Run("EVM().GasEstimator().PriceDefault()", func(t *testing.T) {
		assert.Equal(t, assets.NewWeiI(20000000000), cfg.EVM().GasEstimator().PriceDefault())

		gcfg2 := configtest.NewGeneralConfig(t, overrides)
		cfg2 := evmtest.NewChainScopedConfig(t, gcfg2)
		assert.Equal(t, assets.NewWeiI(42000000000), cfg2.EVM().GasEstimator().PriceDefault())
	})

	t.Run("EvmGasBumpTxDepthDefault", func(t *testing.T) {
		t.Run("uses MaxInFlightTransactions when not set", func(t *testing.T) {
			assert.Equal(t, cfg.EVM().Transactions().MaxInFlight(), cfg.EVM().GasEstimator().BumpTxDepth())
		})

		t.Run("uses customer configured value when set", func(t *testing.T) {
			var override uint32 = 10
			gasBumpOverrides := func(c *chainlink.Config, s *chainlink.Secrets) {
				id := ubig.New(big.NewInt(rand.Int63()))
				c.EVM[0] = &toml.EVMConfig{
					ChainID: id,
					Chain: toml.Defaults(id, &toml.Chain{
						GasEstimator: toml.GasEstimator{
							BumpTxDepth: ptr(override),
						},
					}),
				}
			}
			gcfg2 := configtest.NewGeneralConfig(t, gasBumpOverrides)
			cfg2 := evmtest.NewChainScopedConfig(t, gcfg2)
			assert.NotEqual(t, cfg2.EVM().Transactions().MaxInFlight(), cfg2.EVM().GasEstimator().BumpTxDepth())
			assert.Equal(t, override, cfg2.EVM().GasEstimator().BumpTxDepth())
		})
	})

	t.Run("PriceMaxKey", func(t *testing.T) {
		addr := testutils.NewAddress()
		randomOtherAddr := testutils.NewAddress()
		gcfg2 := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
			overrides(c, s)
			c.EVM[0].KeySpecific = toml.KeySpecificConfig{
				{Key: ptr(types.EIP55AddressFromAddress(randomOtherAddr)),
					GasEstimator: toml.KeySpecificGasEstimator{
						PriceMax: assets.GWei(850),
					},
				},
			}
		})
		cfg2 := evmtest.NewChainScopedConfig(t, gcfg2)

		t.Run("uses chain-specific default value when nothing is set", func(t *testing.T) {
			assert.Equal(t, assets.NewWeiI(100000000000000), cfg2.EVM().GasEstimator().PriceMaxKey(addr))
		})

		t.Run("uses chain-specific override value when that is set", func(t *testing.T) {
			val := assets.NewWeiI(rand.Int63())
			gcfg3 := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
				c.EVM[0].GasEstimator.PriceMax = val
			})
			cfg3 := evmtest.NewChainScopedConfig(t, gcfg3)

			assert.Equal(t, val.String(), cfg3.EVM().GasEstimator().PriceMaxKey(addr).String())
		})
		t.Run("uses key-specific override value when set", func(t *testing.T) {
			tests := []struct {
				name string
				val  *assets.Wei
			}{
				{"Test with 250 GWei", assets.GWei(250)},
				{"Test with 0 GWei", assets.GWei(0)},
			}

			for _, tt := range tests {
				t.Run(tt.name, func(t *testing.T) {
					gcfg3 := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
						c.EVM[0].KeySpecific = toml.KeySpecificConfig{
							{Key: ptr(types.EIP55AddressFromAddress(addr)),
								GasEstimator: toml.KeySpecificGasEstimator{
									PriceMax: tt.val,
								},
							},
						}
					})
					cfg3 := evmtest.NewChainScopedConfig(t, gcfg3)

					assert.Equal(t, tt.val.String(), cfg3.EVM().GasEstimator().PriceMaxKey(addr).String())
				})
			}
		})
		t.Run("uses key-specific override value when set and lower than chain specific config", func(t *testing.T) {
			keySpecificPrice := assets.GWei(900)
			chainSpecificPrice := assets.GWei(1200)
			gcfg3 := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
				c.EVM[0].GasEstimator.PriceMax = chainSpecificPrice
				c.EVM[0].KeySpecific = toml.KeySpecificConfig{
					{Key: ptr(types.EIP55AddressFromAddress(addr)),
						GasEstimator: toml.KeySpecificGasEstimator{
							PriceMax: keySpecificPrice,
						},
					},
				}
			})
			cfg3 := evmtest.NewChainScopedConfig(t, gcfg3)

			assert.Equal(t, keySpecificPrice.String(), cfg3.EVM().GasEstimator().PriceMaxKey(addr).String())
		})
		t.Run("uses chain-specific value when higher than key-specific value", func(t *testing.T) {
			keySpecificPrice := assets.GWei(1400)
			chainSpecificPrice := assets.GWei(1200)
			gcfg3 := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
				c.EVM[0].GasEstimator.PriceMax = chainSpecificPrice
				c.EVM[0].KeySpecific = toml.KeySpecificConfig{
					{Key: ptr(types.EIP55AddressFromAddress(addr)),
						GasEstimator: toml.KeySpecificGasEstimator{
							PriceMax: keySpecificPrice,
						},
					},
				}
			})
			cfg3 := evmtest.NewChainScopedConfig(t, gcfg3)

			assert.Equal(t, chainSpecificPrice.String(), cfg3.EVM().GasEstimator().PriceMaxKey(addr).String())
		})
		t.Run("uses key-specific override value when set and lower than global config", func(t *testing.T) {
			keySpecificPrice := assets.GWei(900)
			gcfg3 := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
				c.EVM[0].KeySpecific = toml.KeySpecificConfig{
					{Key: ptr(types.EIP55AddressFromAddress(addr)),
						GasEstimator: toml.KeySpecificGasEstimator{
							PriceMax: keySpecificPrice,
						},
					},
				}
			})
			cfg3 := evmtest.NewChainScopedConfig(t, gcfg3)

			assert.Equal(t, keySpecificPrice.String(), cfg3.EVM().GasEstimator().PriceMaxKey(addr).String())
		})
		t.Run("uses global value when higher than key-specific value", func(t *testing.T) {
			keySpecificPrice := assets.GWei(1400)
			chainSpecificPrice := assets.GWei(1200)
			gcfg3 := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
				c.EVM[0].GasEstimator.PriceMax = chainSpecificPrice
				c.EVM[0].KeySpecific = toml.KeySpecificConfig{
					{Key: ptr(types.EIP55AddressFromAddress(addr)),
						GasEstimator: toml.KeySpecificGasEstimator{
							PriceMax: keySpecificPrice,
						},
					},
				}
			})
			cfg3 := evmtest.NewChainScopedConfig(t, gcfg3)

			assert.Equal(t, chainSpecificPrice.String(), cfg3.EVM().GasEstimator().PriceMaxKey(addr).String())
		})
		t.Run("uses global value when there is no key-specific price", func(t *testing.T) {
			val := assets.NewWeiI(rand.Int63())
			unsetAddr := testutils.NewAddress()
			gcfg3 := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
				c.EVM[0].GasEstimator.PriceMax = val
			})
			cfg3 := evmtest.NewChainScopedConfig(t, gcfg3)

			assert.Equal(t, val.String(), cfg3.EVM().GasEstimator().PriceMaxKey(unsetAddr).String())
		})
	})

	t.Run("LinkContractAddress", func(t *testing.T) {
		t.Run("uses chain-specific default value when nothing is set", func(t *testing.T) {
			assert.Equal(t, "", cfg.EVM().LinkContractAddress())
		})

		t.Run("uses chain-specific override value when that is set", func(t *testing.T) {
			val := testutils.NewAddress()

			gcfg3 := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
				c.EVM[0].LinkContractAddress = ptr(types.EIP55AddressFromAddress(val))
			})
			cfg3 := evmtest.NewChainScopedConfig(t, gcfg3)

			assert.Equal(t, val.String(), cfg3.EVM().LinkContractAddress())
		})
	})

	t.Run("OperatorFactoryAddress", func(t *testing.T) {
		t.Run("uses chain-specific default value when nothing is set", func(t *testing.T) {
			assert.Equal(t, "", cfg.EVM().OperatorFactoryAddress())
		})

		t.Run("uses chain-specific override value when that is set", func(t *testing.T) {
			val := testutils.NewAddress()

			gcfg3 := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
				c.EVM[0].OperatorFactoryAddress = ptr(types.EIP55AddressFromAddress(val))
			})
			cfg3 := evmtest.NewChainScopedConfig(t, gcfg3)

			assert.Equal(t, val.String(), cfg3.EVM().OperatorFactoryAddress())
		})
	})
}

func TestChainScopedConfig_BlockHistory(t *testing.T) {
	t.Parallel()
	gcfg := configtest.NewTestGeneralConfig(t)
	cfg := evmtest.NewChainScopedConfig(t, gcfg)

	bh := cfg.EVM().GasEstimator().BlockHistory()
	assert.Equal(t, uint32(25), bh.BatchSize())
	assert.Equal(t, uint16(8), bh.BlockHistorySize())
	assert.Equal(t, uint16(60), bh.TransactionPercentile())
	assert.Equal(t, uint16(90), bh.CheckInclusionPercentile())
	assert.Equal(t, uint16(12), bh.CheckInclusionBlocks())
	assert.Equal(t, uint16(1), bh.BlockDelay())
	assert.Equal(t, uint16(4), bh.EIP1559FeeCapBufferBlocks())
}

func TestChainScopedConfig_GasEstimator(t *testing.T) {
	t.Parallel()
	gcfg := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		c.EVM[0].GasEstimator.PriceMax = assets.GWei(500)
	})
	cfg := evmtest.NewChainScopedConfig(t, gcfg)

	ge := cfg.EVM().GasEstimator()
	assert.Equal(t, "BlockHistory", ge.Mode())
	assert.Equal(t, assets.GWei(20), ge.PriceDefault())
	assert.Equal(t, assets.GWei(500), ge.PriceMax())
	assert.Equal(t, assets.GWei(1), ge.PriceMin())
	assert.Equal(t, uint64(500000), ge.LimitDefault())
	assert.Equal(t, uint64(500000), ge.LimitMax())
	assert.Equal(t, float32(1), ge.LimitMultiplier())
	assert.Equal(t, uint64(21000), ge.LimitTransfer())
	assert.Equal(t, assets.GWei(5), ge.BumpMin())
	assert.Equal(t, uint16(20), ge.BumpPercent())
	assert.Equal(t, uint64(3), ge.BumpThreshold())
	assert.False(t, ge.EIP1559DynamicFees())
	assert.Equal(t, assets.GWei(100), ge.FeeCapDefault())
	assert.Equal(t, assets.NewWeiI(1), ge.TipCapDefault())
	assert.Equal(t, assets.NewWeiI(1), ge.TipCapMin())
}

func TestChainScopedConfig_BSCDefaults(t *testing.T) {
	chainID := big.NewInt(56)
	gcfg := configtest.NewGeneralConfig(t, func(c *chainlink.Config, secrets *chainlink.Secrets) {
		id := ubig.New(chainID)
		cfg := toml.Defaults(id)
		c.EVM[0] = &toml.EVMConfig{
			ChainID: id,
			Enabled: ptr(true),
			Chain:   cfg,
		}
	})
	cfg := evmtest.NewChainScopedConfig(t, gcfg)

	timeout := cfg.EVM().OCR().DatabaseTimeout()
	require.Equal(t, 2*time.Second, timeout)
	timeout = cfg.EVM().OCR().ContractTransmitterTransmitTimeout()
	require.Equal(t, 2*time.Second, timeout)
	timeout = cfg.EVM().OCR().ObservationGracePeriod()
	require.Equal(t, 500*time.Millisecond, timeout)
}

func TestChainScopedConfig_Profiles(t *testing.T) {
	t.Parallel()

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

		{"gnosisMainnet", 100, 500000, "0.00001"},
	}
	for _, test := range tests {
		tt := test

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			gcfg := configtest.NewGeneralConfig(t, func(c *chainlink.Config, secrets *chainlink.Secrets) {
				id := ubig.NewI(tt.chainID)
				cfg := toml.Defaults(id)
				c.EVM[0] = &toml.EVMConfig{
					ChainID: id,
					Enabled: ptr(true),
					Chain:   cfg,
				}
			})
			config := evmtest.NewChainScopedConfig(t, gcfg)

			assert.Equal(t, tt.expectedGasLimitDefault, config.EVM().GasEstimator().LimitDefault())
			assert.Nil(t, config.EVM().GasEstimator().LimitJobType().OCR())
			assert.Nil(t, config.EVM().GasEstimator().LimitJobType().DR())
			assert.Nil(t, config.EVM().GasEstimator().LimitJobType().VRF())
			assert.Nil(t, config.EVM().GasEstimator().LimitJobType().FM())
			assert.Nil(t, config.EVM().GasEstimator().LimitJobType().Keeper())
			assert.Equal(t, tt.expectedMinimumContractPayment, strings.TrimRight(config.EVM().MinContractPayment().Link(), "0"))
		})
	}
}

func TestChainScopedConfig_HeadTracker(t *testing.T) {
	t.Parallel()
	gcfg := configtest.NewTestGeneralConfig(t)
	cfg := evmtest.NewChainScopedConfig(t, gcfg)

	ht := cfg.EVM().HeadTracker()
	assert.Equal(t, uint32(100), ht.HistoryDepth())
	assert.Equal(t, uint32(3), ht.MaxBufferSize())
	assert.Equal(t, time.Second, ht.SamplingInterval())
}

func Test_chainScopedConfig_Validate(t *testing.T) {
	configWithChains := func(t *testing.T, id int64, chains ...*toml.Chain) config.AppConfig {
		return configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
			s.Database.URL = models.MustSecretURL("postgresql://doesnotexist:justtopassvalidationtests@localhost:5432/chainlink_na_test")
			chainID := ubig.NewI(id)
			c.EVM[0] = &toml.EVMConfig{ChainID: chainID, Enabled: ptr(true), Chain: toml.Defaults(chainID, chains...),
				Nodes: toml.EVMNodes{{
					Name:    ptr("fake"),
					WSURL:   configurl.MustParseURL("wss://foo.test/ws"),
					HTTPURL: configurl.MustParseURL("http://foo.test"),
				}}}
		})
	}

	// Validate built-in
	for _, id := range toml.DefaultIDs {
		id := id
		t.Run(fmt.Sprintf("chainID-%s", id), func(t *testing.T) {
			cfg := configWithChains(t, id.Int64())
			assert.NoError(t, cfg.Validate())
		})
	}

	// Invalid Cases:

	t.Run("arbitrum-estimator", func(t *testing.T) {
		t.Run("custom", func(t *testing.T) {
			cfg := configWithChains(t, 0, &toml.Chain{
				ChainType: ptr(string(commonconfig.ChainArbitrum)),
				GasEstimator: toml.GasEstimator{
					Mode: ptr("BlockHistory"),
				},
			})
			assert.NoError(t, cfg.Validate())
		})
		t.Run("mainnet", func(t *testing.T) {
			cfg := configWithChains(t, 42161, &toml.Chain{
				GasEstimator: toml.GasEstimator{
					Mode: ptr("BlockHistory"),
					BlockHistory: toml.BlockHistoryEstimator{
						BlockHistorySize: ptr[uint16](1),
					},
				},
			})
			assert.NoError(t, cfg.Validate())
		})
		t.Run("testnet", func(t *testing.T) {
			cfg := configWithChains(t, 421611, &toml.Chain{
				GasEstimator: toml.GasEstimator{
					Mode: ptr("SuggestedPrice"),
				},
			})
			assert.NoError(t, cfg.Validate())
		})
	})

	t.Run("optimism-estimator", func(t *testing.T) {
		t.Run("custom", func(t *testing.T) {
			cfg := configWithChains(t, 0, &toml.Chain{
				ChainType: ptr(string(commonconfig.ChainOptimismBedrock)),
				GasEstimator: toml.GasEstimator{
					Mode: ptr("BlockHistory"),
				},
			})
			assert.NoError(t, cfg.Validate())
		})
		t.Run("mainnet", func(t *testing.T) {
			cfg := configWithChains(t, 10, &toml.Chain{
				GasEstimator: toml.GasEstimator{
					Mode: ptr("FixedPrice"),
				},
			})
			assert.NoError(t, cfg.Validate())
		})
		t.Run("testnet", func(t *testing.T) {
			cfg := configWithChains(t, 69, &toml.Chain{
				GasEstimator: toml.GasEstimator{
					Mode: ptr("FixedPrice"),
				},
			})
			assert.NoError(t, cfg.Validate())
		})
	})
}

func TestNodePoolConfig(t *testing.T) {
	gcfg := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		id := ubig.New(big.NewInt(rand.Int63()))
		c.EVM[0] = &toml.EVMConfig{
			ChainID: id,
			Chain:   toml.Defaults(id, &toml.Chain{}),
		}
	})
	cfg := evmtest.NewChainScopedConfig(t, gcfg)

	require.Equal(t, "HighestHead", cfg.EVM().NodePool().SelectionMode())
	require.Equal(t, uint32(5), cfg.EVM().NodePool().SyncThreshold())
	require.Equal(t, time.Duration(10000000000), cfg.EVM().NodePool().PollInterval())
	require.Equal(t, uint32(5), cfg.EVM().NodePool().PollFailureThreshold())
	require.Equal(t, false, cfg.EVM().NodePool().NodeIsSyncingEnabled())
}

func TestClientErrorsConfig(t *testing.T) {
	t.Parallel()

	t.Run("EVM().NodePool().Errors()", func(t *testing.T) {
		clientErrorsOverrides := func(c *chainlink.Config, s *chainlink.Secrets) {
			id := ubig.New(big.NewInt(rand.Int63()))
			c.EVM[0] = &toml.EVMConfig{
				ChainID: id,
				Chain: toml.Defaults(id, &toml.Chain{
					NodePool: toml.NodePool{
						Errors: toml.ClientErrors{
							NonceTooLow:                       ptr[string]("client error nonce too low"),
							NonceTooHigh:                      ptr[string]("client error nonce too high"),
							ReplacementTransactionUnderpriced: ptr[string]("client error replacement underpriced"),
							LimitReached:                      ptr[string]("client error limit reached"),
							TransactionAlreadyInMempool:       ptr[string]("client error transaction already in mempool"),
							TerminallyUnderpriced:             ptr[string]("client error terminally underpriced"),
							InsufficientEth:                   ptr[string]("client error insufficient eth"),
							TxFeeExceedsCap:                   ptr[string]("client error tx fee exceeds cap"),
							L2FeeTooLow:                       ptr[string]("client error l2 fee too low"),
							L2FeeTooHigh:                      ptr[string]("client error l2 fee too high"),
							L2Full:                            ptr[string]("client error l2 full"),
							TransactionAlreadyMined:           ptr[string]("client error transaction already mined"),
							Fatal:                             ptr[string]("client error fatal"),
							ServiceUnavailable:                ptr[string]("client error service unavailable"),
						},
					},
				}),
			}
		}

		gcfg := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
			clientErrorsOverrides(c, s)
		})
		cfg := evmtest.NewChainScopedConfig(t, gcfg)

		errors := cfg.EVM().NodePool().Errors()
		assert.Equal(t, "client error nonce too low", errors.NonceTooLow())
		assert.Equal(t, "client error nonce too high", errors.NonceTooHigh())
		assert.Equal(t, "client error replacement underpriced", errors.ReplacementTransactionUnderpriced())
		assert.Equal(t, "client error limit reached", errors.LimitReached())
		assert.Equal(t, "client error transaction already in mempool", errors.TransactionAlreadyInMempool())
		assert.Equal(t, "client error terminally underpriced", errors.TerminallyUnderpriced())
		assert.Equal(t, "client error insufficient eth", errors.InsufficientEth())
		assert.Equal(t, "client error tx fee exceeds cap", errors.TxFeeExceedsCap())
		assert.Equal(t, "client error l2 fee too low", errors.L2FeeTooLow())
		assert.Equal(t, "client error l2 fee too high", errors.L2FeeTooHigh())
		assert.Equal(t, "client error l2 full", errors.L2Full())
		assert.Equal(t, "client error transaction already mined", errors.TransactionAlreadyMined())
		assert.Equal(t, "client error fatal", errors.Fatal())
		assert.Equal(t, "client error service unavailable", errors.ServiceUnavailable())
	})
}

func ptr[T any](t T) *T { return &t }
