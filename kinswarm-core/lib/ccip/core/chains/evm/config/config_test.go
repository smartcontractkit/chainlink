package config_test

import (
	"math/big"
	"math/rand"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/assets"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/config/toml"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	ubig "github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils/big"
)

func TestChainScopedConfig(t *testing.T) {
	t.Parallel()
	cfg := testutils.NewTestChainScopedConfig(t, func(c *toml.EVMConfig) {
		c.GasEstimator.PriceMax = assets.NewWeiI(100000000000000)
	})

	cfg2 := testutils.NewTestChainScopedConfig(t, func(c *toml.EVMConfig) {
		c.GasEstimator.PriceMax = assets.NewWeiI(100000000000000)
		c.GasEstimator.PriceDefault = assets.NewWeiI(42000000000)
	})

	t.Run("EVM().GasEstimator().PriceDefault()", func(t *testing.T) {
		assert.Equal(t, assets.NewWeiI(20000000000), cfg.EVM().GasEstimator().PriceDefault())

		assert.Equal(t, assets.NewWeiI(42000000000), cfg2.EVM().GasEstimator().PriceDefault())
	})

	t.Run("EvmGasBumpTxDepthDefault", func(t *testing.T) {
		t.Run("uses MaxInFlightTransactions when not set", func(t *testing.T) {
			assert.Equal(t, cfg.EVM().Transactions().MaxInFlight(), cfg.EVM().GasEstimator().BumpTxDepth())
		})

		t.Run("uses customer configured value when set", func(t *testing.T) {
			var bumpTxDepth uint32 = 10
			cfg2 := testutils.NewTestChainScopedConfig(t, func(c *toml.EVMConfig) {
				c.GasEstimator.BumpTxDepth = &bumpTxDepth
			})
			assert.NotEqual(t, cfg2.EVM().Transactions().MaxInFlight(), cfg2.EVM().GasEstimator().BumpTxDepth())
			assert.Equal(t, bumpTxDepth, cfg2.EVM().GasEstimator().BumpTxDepth())
		})
	})

	t.Run("PriceMaxKey", func(t *testing.T) {
		addr := testutils.NewAddress()
		randomOtherAddr := testutils.NewAddress()
		cfg2 := testutils.NewTestChainScopedConfig(t, func(c *toml.EVMConfig) {
			c.KeySpecific = toml.KeySpecificConfig{
				{Key: ptr(types.EIP55AddressFromAddress(randomOtherAddr)),
					GasEstimator: toml.KeySpecificGasEstimator{
						PriceMax: assets.GWei(850),
					},
				},
			}
			c.GasEstimator.PriceMax = assets.NewWeiI(100000000000000)
			c.GasEstimator.PriceDefault = assets.NewWeiI(42000000000)
		})

		t.Run("uses chain-specific default value when nothing is set", func(t *testing.T) {
			assert.Equal(t, assets.NewWeiI(100000000000000), cfg2.EVM().GasEstimator().PriceMaxKey(addr))
		})

		t.Run("uses chain-specific override value when that is set", func(t *testing.T) {
			priceMax := assets.NewWeiI(rand.Int63())
			cfg3 := testutils.NewTestChainScopedConfig(t, func(c *toml.EVMConfig) {
				c.GasEstimator.PriceMax = priceMax
			})
			assert.Equal(t, priceMax.String(), cfg3.EVM().GasEstimator().PriceMaxKey(addr).String())
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
					cfg3 := testutils.NewTestChainScopedConfig(t, func(c *toml.EVMConfig) {
						c.KeySpecific = toml.KeySpecificConfig{
							{Key: ptr(types.EIP55AddressFromAddress(addr)),
								GasEstimator: toml.KeySpecificGasEstimator{
									PriceMax: tt.val,
								},
							},
						}
					})

					assert.Equal(t, tt.val.String(), cfg3.EVM().GasEstimator().PriceMaxKey(addr).String())
				})
			}
		})
		t.Run("uses key-specific override value when set and lower than chain specific config", func(t *testing.T) {
			keySpecificPrice := assets.GWei(900)
			chainSpecificPrice := assets.GWei(1200)
			cfg3 := testutils.NewTestChainScopedConfig(t, func(c *toml.EVMConfig) {
				c.GasEstimator.PriceMax = chainSpecificPrice
				c.KeySpecific = toml.KeySpecificConfig{
					{Key: ptr(types.EIP55AddressFromAddress(addr)),
						GasEstimator: toml.KeySpecificGasEstimator{
							PriceMax: keySpecificPrice,
						},
					},
				}
			})

			assert.Equal(t, keySpecificPrice.String(), cfg3.EVM().GasEstimator().PriceMaxKey(addr).String())
		})
		t.Run("uses chain-specific value when higher than key-specific value", func(t *testing.T) {
			keySpecificPrice := assets.GWei(1400)
			chainSpecificPrice := assets.GWei(1200)
			cfg3 := testutils.NewTestChainScopedConfig(t, func(c *toml.EVMConfig) {
				c.GasEstimator.PriceMax = chainSpecificPrice
				c.KeySpecific = toml.KeySpecificConfig{
					{Key: ptr(types.EIP55AddressFromAddress(addr)),
						GasEstimator: toml.KeySpecificGasEstimator{
							PriceMax: keySpecificPrice,
						},
					},
				}
			})

			assert.Equal(t, chainSpecificPrice.String(), cfg3.EVM().GasEstimator().PriceMaxKey(addr).String())
		})
		t.Run("uses key-specific override value when set and lower than global config", func(t *testing.T) {
			keySpecificPrice := assets.GWei(900)
			cfg3 := testutils.NewTestChainScopedConfig(t, func(c *toml.EVMConfig) {
				c.KeySpecific = toml.KeySpecificConfig{
					{Key: ptr(types.EIP55AddressFromAddress(addr)),
						GasEstimator: toml.KeySpecificGasEstimator{
							PriceMax: keySpecificPrice,
						},
					},
				}
			})

			assert.Equal(t, keySpecificPrice.String(), cfg3.EVM().GasEstimator().PriceMaxKey(addr).String())
		})
		t.Run("uses global value when higher than key-specific value", func(t *testing.T) {
			keySpecificPrice := assets.GWei(1400)
			chainSpecificPrice := assets.GWei(1200)
			cfg3 := testutils.NewTestChainScopedConfig(t, func(c *toml.EVMConfig) {
				c.GasEstimator.PriceMax = chainSpecificPrice
				c.KeySpecific = toml.KeySpecificConfig{
					{Key: ptr(types.EIP55AddressFromAddress(addr)),
						GasEstimator: toml.KeySpecificGasEstimator{
							PriceMax: keySpecificPrice,
						},
					},
				}
			})

			assert.Equal(t, chainSpecificPrice.String(), cfg3.EVM().GasEstimator().PriceMaxKey(addr).String())
		})
		t.Run("uses global value when there is no key-specific price", func(t *testing.T) {
			priceMax := assets.NewWeiI(rand.Int63())
			unsetAddr := testutils.NewAddress()
			cfg3 := testutils.NewTestChainScopedConfig(t, func(c *toml.EVMConfig) {
				c.GasEstimator.PriceMax = priceMax
			})

			assert.Equal(t, priceMax.String(), cfg3.EVM().GasEstimator().PriceMaxKey(unsetAddr).String())
		})
	})

	t.Run("LinkContractAddress", func(t *testing.T) {
		t.Run("uses chain-specific default value when nothing is set", func(t *testing.T) {
			assert.Equal(t, "", cfg.EVM().LinkContractAddress())
		})

		t.Run("uses chain-specific override value when that is set", func(t *testing.T) {
			addr := testutils.NewAddress()

			cfg3 := testutils.NewTestChainScopedConfig(t, func(c *toml.EVMConfig) {
				c.LinkContractAddress = ptr(types.EIP55AddressFromAddress(addr))
			})

			assert.Equal(t, addr.String(), cfg3.EVM().LinkContractAddress())
		})
	})

	t.Run("OperatorFactoryAddress", func(t *testing.T) {
		t.Run("uses chain-specific default value when nothing is set", func(t *testing.T) {
			assert.Equal(t, "", cfg.EVM().OperatorFactoryAddress())
		})

		t.Run("uses chain-specific override value when that is set", func(t *testing.T) {
			val := testutils.NewAddress()

			cfg3 := testutils.NewTestChainScopedConfig(t, func(c *toml.EVMConfig) {
				c.OperatorFactoryAddress = ptr(types.EIP55AddressFromAddress(val))
			})

			assert.Equal(t, val.String(), cfg3.EVM().OperatorFactoryAddress())
		})
	})

	t.Run("LogBroadcasterEnabled", func(t *testing.T) {
		t.Run("turn on LogBroadcasterEnabled by default", func(t *testing.T) {
			assert.Equal(t, true, cfg.EVM().LogBroadcasterEnabled())
		})

		t.Run("verify LogBroadcasterEnabled is set correctly", func(t *testing.T) {
			val := false
			cfg3 := testutils.NewTestChainScopedConfig(t, func(c *toml.EVMConfig) {
				c.LogBroadcasterEnabled = ptr(val)
			})

			assert.Equal(t, false, cfg3.EVM().LogBroadcasterEnabled())
		})

		t.Run("use Noop logBroadcaster when LogBroadcaster is disabled", func(t *testing.T) {

		})
	})
}

func TestChainScopedConfig_BlockHistory(t *testing.T) {
	t.Parallel()
	cfg := testutils.NewTestChainScopedConfig(t, nil)

	bh := cfg.EVM().GasEstimator().BlockHistory()
	assert.Equal(t, uint32(25), bh.BatchSize())
	assert.Equal(t, uint16(8), bh.BlockHistorySize())
	assert.Equal(t, uint16(60), bh.TransactionPercentile())
	assert.Equal(t, uint16(90), bh.CheckInclusionPercentile())
	assert.Equal(t, uint16(12), bh.CheckInclusionBlocks())
	assert.Equal(t, uint16(1), bh.BlockDelay())
	assert.Equal(t, uint16(4), bh.EIP1559FeeCapBufferBlocks())
}
func TestChainScopedConfig_FeeHistory(t *testing.T) {
	t.Parallel()
	cfg := testutils.NewTestChainScopedConfig(t, nil)

	u := cfg.EVM().GasEstimator().FeeHistory()
	assert.Equal(t, 10*time.Second, u.CacheTimeout())
}

func TestChainScopedConfig_GasEstimator(t *testing.T) {
	t.Parallel()
	cfg := testutils.NewTestChainScopedConfig(t, func(c *toml.EVMConfig) {
		c.GasEstimator.PriceMax = assets.GWei(500)
	})

	ge := cfg.EVM().GasEstimator()
	assert.Equal(t, "BlockHistory", ge.Mode())
	assert.Equal(t, assets.GWei(20), ge.PriceDefault())
	assert.Equal(t, assets.GWei(500), ge.PriceMax())
	assert.Equal(t, assets.GWei(1), ge.PriceMin())
	assert.Equal(t, uint64(8_000_000), ge.LimitDefault())
	assert.Equal(t, uint64(8_000_000), ge.LimitMax())
	assert.Equal(t, float32(1), ge.LimitMultiplier())
	assert.Equal(t, uint64(21000), ge.LimitTransfer())
	assert.Equal(t, assets.GWei(5), ge.BumpMin())
	assert.Equal(t, uint16(20), ge.BumpPercent())
	assert.Equal(t, uint64(3), ge.BumpThreshold())
	assert.False(t, ge.EIP1559DynamicFees())
	assert.Equal(t, assets.GWei(100), ge.FeeCapDefault())
	assert.Equal(t, assets.NewWeiI(1), ge.TipCapDefault())
	assert.Equal(t, assets.NewWeiI(1), ge.TipCapMin())
	assert.Equal(t, false, ge.EstimateLimit())
}

func TestChainScopedConfig_BSCDefaults(t *testing.T) {
	cfg := testutils.NewTestChainScopedConfig(t, func(c *toml.EVMConfig) {
		c.ChainID = (*ubig.Big)(big.NewInt(56))
	})

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
		{"default", 0, 8_000_000, "0.00001"},
		{"mainnet", 1, 8_000_000, "0.1"},
		{"kovan", 42, 8_000_000, "0.1"},

		{"optimism", 10, 8_000_000, "0.00001"},
		{"optimism", 69, 8_000_000, "0.00001"},
		{"optimism", 420, 8_000_000, "0.00001"},

		{"bscMainnet", 56, 8_000_000, "0.00001"},
		{"hecoMainnet", 128, 8_000_000, "0.00001"},
		{"fantomMainnet", 250, 8_000_000, "0.00001"},
		{"fantomTestnet", 4002, 8_000_000, "0.00001"},
		{"polygonMatic", 800001, 8_000_000, "0.00001"},
		{"harmonyMainnet", 1666600000, 8_000_000, "0.00001"},
		{"harmonyTestnet", 1666700000, 8_000_000, "0.00001"},

		{"gnosisMainnet", 100, 8_000_000, "0.00001"},
	}
	for _, test := range tests {
		tt := test

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			config := testutils.NewTestChainScopedConfig(t, func(c *toml.EVMConfig) {
				c.ChainID = ubig.NewI(tt.chainID)
			})

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
	cfg := testutils.NewTestChainScopedConfig(t, nil)

	ht := cfg.EVM().HeadTracker()
	assert.Equal(t, uint32(100), ht.HistoryDepth())
	assert.Equal(t, uint32(3), ht.MaxBufferSize())
	assert.Equal(t, time.Second, ht.SamplingInterval())
	assert.Equal(t, true, ht.FinalityTagBypass())
	assert.Equal(t, uint32(10000), ht.MaxAllowedFinalityDepth())
}

func TestNodePoolConfig(t *testing.T) {
	cfg := testutils.NewTestChainScopedConfig(t, nil)

	require.Equal(t, "HighestHead", cfg.EVM().NodePool().SelectionMode())
	require.Equal(t, uint32(5), cfg.EVM().NodePool().SyncThreshold())
	require.Equal(t, time.Duration(10000000000), cfg.EVM().NodePool().PollInterval())
	require.Equal(t, uint32(5), cfg.EVM().NodePool().PollFailureThreshold())
	require.Equal(t, false, cfg.EVM().NodePool().NodeIsSyncingEnabled())
	require.Equal(t, false, cfg.EVM().NodePool().EnforceRepeatableRead())
	require.Equal(t, time.Duration(10000000000), cfg.EVM().NodePool().DeathDeclarationDelay())
}

func TestClientErrorsConfig(t *testing.T) {
	t.Parallel()

	t.Run("EVM().NodePool().Errors()", func(t *testing.T) {
		cfg := testutils.NewTestChainScopedConfig(t, func(c *toml.EVMConfig) {
			id := ubig.New(big.NewInt(rand.Int63()))
			c.ChainID = id
			c.NodePool = toml.NodePool{
				Errors: toml.ClientErrors{
					NonceTooLow:                       ptr("client error nonce too low"),
					NonceTooHigh:                      ptr("client error nonce too high"),
					ReplacementTransactionUnderpriced: ptr("client error replacement underpriced"),
					LimitReached:                      ptr("client error limit reached"),
					TransactionAlreadyInMempool:       ptr("client error transaction already in mempool"),
					TerminallyUnderpriced:             ptr("client error terminally underpriced"),
					InsufficientEth:                   ptr("client error insufficient eth"),
					TxFeeExceedsCap:                   ptr("client error tx fee exceeds cap"),
					L2FeeTooLow:                       ptr("client error l2 fee too low"),
					L2FeeTooHigh:                      ptr("client error l2 fee too high"),
					L2Full:                            ptr("client error l2 full"),
					TransactionAlreadyMined:           ptr("client error transaction already mined"),
					Fatal:                             ptr("client error fatal"),
					ServiceUnavailable:                ptr("client error service unavailable"),
					TooManyResults:                    ptr("client error too many results"),
				},
			}
		})

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
		assert.Equal(t, "client error too many results", errors.TooManyResults())
	})
}

func ptr[T any](t T) *T { return &t }
