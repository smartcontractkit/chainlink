package v1_5_1_test

import (
	"bytes"
	"math/big"
	"sort"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/v1_5_1"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/v1_5_1/token_pool"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/shared/generated/burn_mint_erc677"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

// createSymmetricRateLimits is a utility to quickly create a rate limiter config with equal inbound and outbound values.
func createSymmetricRateLimits(rate int64, capacity int64) v1_5_1.RateLimiterConfig {
	return v1_5_1.RateLimiterConfig{
		Inbound: token_pool.RateLimiterConfig{
			IsEnabled: rate != 0 || capacity != 0,
			Rate:      big.NewInt(rate),
			Capacity:  big.NewInt(capacity),
		},
		Outbound: token_pool.RateLimiterConfig{
			IsEnabled: rate != 0 || capacity != 0,
			Rate:      big.NewInt(rate),
			Capacity:  big.NewInt(capacity),
		},
	}
}

// validateMemberOfTokenPoolPair performs checks required to validate that a token pool is fully configured for cross-chain transfer.
func validateMemberOfTokenPoolPair(
	t *testing.T,
	state changeset.CCIPOnChainState,
	tokenPool *token_pool.TokenPool,
	expectedRemotePools []common.Address,
	tokens map[uint64]*deployment.ContractDeploy[*burn_mint_erc677.BurnMintERC677],
	tokenSymbol changeset.TokenSymbol,
	chainSelector uint64,
	rate *big.Int,
	capacity *big.Int,
	expectedOwner common.Address,
) {
	// Verify that the owner is expected
	owner, err := tokenPool.Owner(nil)
	require.NoError(t, err)
	require.Equal(t, expectedOwner, owner)

	// Fetch the supported remote chains
	supportedChains, err := tokenPool.GetSupportedChains(nil)
	require.NoError(t, err)

	// Verify that the rate limits and remote addresses are correct
	for _, supportedChain := range supportedChains {
		inboundConfig, err := tokenPool.GetCurrentInboundRateLimiterState(nil, supportedChain)
		require.NoError(t, err)
		require.True(t, inboundConfig.IsEnabled)
		require.Equal(t, capacity, inboundConfig.Capacity)
		require.Equal(t, rate, inboundConfig.Rate)

		outboundConfig, err := tokenPool.GetCurrentOutboundRateLimiterState(nil, supportedChain)
		require.NoError(t, err)
		require.True(t, outboundConfig.IsEnabled)
		require.Equal(t, capacity, outboundConfig.Capacity)
		require.Equal(t, rate, outboundConfig.Rate)

		remoteTokenAddress, err := tokenPool.GetRemoteToken(nil, supportedChain)
		require.NoError(t, err)
		require.Equal(t, common.LeftPadBytes(tokens[supportedChain].Address.Bytes(), 32), remoteTokenAddress)

		remotePoolAddresses, err := tokenPool.GetRemotePools(nil, supportedChain)
		require.NoError(t, err)

		require.Equal(t, len(expectedRemotePools), len(remotePoolAddresses))
		expectedRemotePoolAddressesBytes := make([][]byte, len(expectedRemotePools))
		for i, remotePool := range expectedRemotePools {
			expectedRemotePoolAddressesBytes[i] = common.LeftPadBytes(remotePool.Bytes(), 32)
		}
		sort.Slice(expectedRemotePoolAddressesBytes, func(i, j int) bool {
			return bytes.Compare(expectedRemotePoolAddressesBytes[i], expectedRemotePoolAddressesBytes[j]) < 0
		})
		sort.Slice(remotePoolAddresses, func(i, j int) bool {
			return bytes.Compare(remotePoolAddresses[i], remotePoolAddresses[j]) < 0
		})
		for i := range expectedRemotePoolAddressesBytes {
			require.Equal(t, expectedRemotePoolAddressesBytes[i], remotePoolAddresses[i])
		}
	}
}

func TestValidateRemoteChains(t *testing.T) {
	t.Parallel()

	tests := []struct {
		IsEnabled bool
		Rate      *big.Int
		Capacity  *big.Int
		ErrStr    string
	}{
		{
			IsEnabled: false,
			Rate:      big.NewInt(1),
			Capacity:  big.NewInt(10),
			ErrStr:    "rate and capacity must be 0",
		},
		{
			IsEnabled: true,
			Rate:      big.NewInt(0),
			Capacity:  big.NewInt(10),
			ErrStr:    "rate must be greater than 0 and less than capacity",
		},
		{
			IsEnabled: true,
			Rate:      big.NewInt(11),
			Capacity:  big.NewInt(10),
			ErrStr:    "rate must be greater than 0 and less than capacity",
		},
	}

	for _, test := range tests {
		t.Run(test.ErrStr, func(t *testing.T) {
			remoteChains := v1_5_1.RateLimiterPerChain{
				1: {
					Inbound: token_pool.RateLimiterConfig{
						IsEnabled: test.IsEnabled,
						Rate:      test.Rate,
						Capacity:  test.Capacity,
					},
					Outbound: token_pool.RateLimiterConfig{
						IsEnabled: test.IsEnabled,
						Rate:      test.Rate,
						Capacity:  test.Capacity,
					},
				},
			}

			err := remoteChains.Validate()
			require.Error(t, err)
			require.Contains(t, err.Error(), test.ErrStr)
		})
	}
}

func TestValidateTokenPoolConfig(t *testing.T) {
	t.Parallel()

	e, selectorA, _, tokens, _ := testhelpers.SetupTwoChainEnvironmentWithTokens(t, logger.TestLogger(t), true)

	e = testhelpers.DeployTestTokenPools(t, e, map[uint64]v1_5_1.DeployTokenPoolInput{
		selectorA: {
			Type:               changeset.BurnMintTokenPool,
			TokenAddress:       tokens[selectorA].Address,
			LocalTokenDecimals: testhelpers.LocalTokenDecimals,
		},
	}, true)

	state, err := changeset.LoadOnchainState(e)
	require.NoError(t, err)

	tests := []struct {
		UseMcms         bool
		TokenPoolConfig v1_5_1.TokenPoolConfig
		ErrStr          string
		Msg             string
	}{
		{
			Msg:             "Pool type is invalid",
			TokenPoolConfig: v1_5_1.TokenPoolConfig{},
			ErrStr:          "is not a known token pool type",
		},
		{
			Msg: "Pool version is invalid",
			TokenPoolConfig: v1_5_1.TokenPoolConfig{
				Type: changeset.BurnMintTokenPool,
			},
			ErrStr: "is not a known token pool version",
		},
		{
			Msg: "Pool is not owned by required address",
			TokenPoolConfig: v1_5_1.TokenPoolConfig{
				Type:    changeset.BurnMintTokenPool,
				Version: deployment.Version1_5_1,
			},
			ErrStr: "failed ownership validation",
		},
	}

	for _, test := range tests {
		t.Run(test.Msg, func(t *testing.T) {
			err := test.TokenPoolConfig.Validate(e.GetContext(), e.Chains[selectorA], state.Chains[selectorA], test.UseMcms, testhelpers.TestTokenSymbol)
			require.Error(t, err)
			require.ErrorContains(t, err, test.ErrStr)
		})
	}
}

func TestValidateConfigureTokenPoolContractsConfig(t *testing.T) {
	t.Parallel()

	lggr := logger.TestLogger(t)
	e := memory.NewMemoryEnvironment(t, lggr, zapcore.InfoLevel, memory.MemoryEnvironmentConfig{
		Chains: 2,
	})

	tests := []struct {
		TokenSymbol changeset.TokenSymbol
		Input       v1_5_1.ConfigureTokenPoolContractsConfig
		ErrStr      string
		Msg         string
	}{
		{
			Msg:    "Token symbol is missing",
			Input:  v1_5_1.ConfigureTokenPoolContractsConfig{},
			ErrStr: "token symbol must be defined",
		},
		{
			Msg: "Chain selector is invalid",
			Input: v1_5_1.ConfigureTokenPoolContractsConfig{
				TokenSymbol: testhelpers.TestTokenSymbol,
				PoolUpdates: map[uint64]v1_5_1.TokenPoolConfig{
					0: v1_5_1.TokenPoolConfig{},
				},
			},
			ErrStr: "failed to validate chain selector 0",
		},
		{
			Msg: "Chain selector doesn't exist in environment",
			Input: v1_5_1.ConfigureTokenPoolContractsConfig{
				TokenSymbol: testhelpers.TestTokenSymbol,
				PoolUpdates: map[uint64]v1_5_1.TokenPoolConfig{
					5009297550715157269: v1_5_1.TokenPoolConfig{},
				},
			},
			ErrStr: "does not exist in environment",
		},
		{
			Msg: "Corresponding pool update missing",
			Input: v1_5_1.ConfigureTokenPoolContractsConfig{
				TokenSymbol: testhelpers.TestTokenSymbol,
				PoolUpdates: map[uint64]v1_5_1.TokenPoolConfig{
					e.AllChainSelectors()[0]: v1_5_1.TokenPoolConfig{
						ChainUpdates: v1_5_1.RateLimiterPerChain{
							e.AllChainSelectors()[1]: v1_5_1.RateLimiterConfig{},
						},
					},
				},
			},
			ErrStr: "is expecting a pool update to be defined for chain with selector",
		},
		/* This test condition is flakey, as we will see "missing tokenAdminRegistry" if e.AllChainSelectors()[1] is checked first
		{
			Msg: "Corresponding pool update missing a chain update",
			Input: changeset.ConfigureTokenPoolContractsConfig{
				TokenSymbol: testhelpers.TestTokenSymbol,
				PoolUpdates: map[uint64]changeset.TokenPoolConfig{
					e.AllChainSelectors()[0]: changeset.TokenPoolConfig{
						ChainUpdates: changeset.RateLimiterPerChain{
							e.AllChainSelectors()[1]: changeset.RateLimiterConfig{},
						},
					},
					e.AllChainSelectors()[1]: changeset.TokenPoolConfig{},
				},
			},
			ErrStr: "to define a chain config pointing back to it",
		},
		*/
		{
			Msg: "Token admin registry is missing",
			Input: v1_5_1.ConfigureTokenPoolContractsConfig{
				TokenSymbol: testhelpers.TestTokenSymbol,
				PoolUpdates: map[uint64]v1_5_1.TokenPoolConfig{
					e.AllChainSelectors()[0]: v1_5_1.TokenPoolConfig{
						ChainUpdates: v1_5_1.RateLimiterPerChain{
							e.AllChainSelectors()[1]: v1_5_1.RateLimiterConfig{},
						},
					},
					e.AllChainSelectors()[1]: v1_5_1.TokenPoolConfig{
						ChainUpdates: v1_5_1.RateLimiterPerChain{
							e.AllChainSelectors()[0]: v1_5_1.RateLimiterConfig{},
						},
					},
				},
			},
			ErrStr: "missing tokenAdminRegistry",
		},
	}

	for _, test := range tests {
		t.Run(test.Msg, func(t *testing.T) {
			err := test.Input.Validate(e)
			require.Contains(t, err.Error(), test.ErrStr)
		})
	}
}

func TestValidateConfigureTokenPoolContracts(t *testing.T) {
	t.Parallel()

	type regPass struct {
		SelectorA2B v1_5_1.RateLimiterConfig
		SelectorB2A v1_5_1.RateLimiterConfig
	}

	type updatePass struct {
		UpdatePoolOnA bool
		UpdatePoolOnB bool
		SelectorA2B   v1_5_1.RateLimiterConfig
		SelectorB2A   v1_5_1.RateLimiterConfig
	}

	type tokenPools struct {
		LockRelease *token_pool.TokenPool
		BurnMint    *token_pool.TokenPool
	}

	acceptLiquidity := false

	tests := []struct {
		Msg              string
		RegistrationPass *regPass
		UpdatePass       *updatePass
	}{
		{
			Msg: "Configure new pools on registry",
			RegistrationPass: &regPass{
				SelectorA2B: createSymmetricRateLimits(100, 1000),
				SelectorB2A: createSymmetricRateLimits(100, 1000),
			},
		},
		{
			Msg: "Configure new pools on registry, update their rate limits",
			RegistrationPass: &regPass{
				SelectorA2B: createSymmetricRateLimits(100, 1000),
				SelectorB2A: createSymmetricRateLimits(100, 1000),
			},
			UpdatePass: &updatePass{
				UpdatePoolOnA: false,
				UpdatePoolOnB: false,
				SelectorA2B:   createSymmetricRateLimits(200, 2000),
				SelectorB2A:   createSymmetricRateLimits(200, 2000),
			},
		},
		{
			Msg: "Configure new pools on registry, update both pools",
			RegistrationPass: &regPass{
				SelectorA2B: createSymmetricRateLimits(100, 1000),
				SelectorB2A: createSymmetricRateLimits(100, 1000),
			},
			UpdatePass: &updatePass{
				UpdatePoolOnA: true,
				UpdatePoolOnB: true,
				SelectorA2B:   createSymmetricRateLimits(100, 1000),
				SelectorB2A:   createSymmetricRateLimits(100, 1000),
			},
		},
		{
			Msg: "Configure new pools on registry, update only one pool",
			RegistrationPass: &regPass{
				SelectorA2B: createSymmetricRateLimits(100, 1000),
				SelectorB2A: createSymmetricRateLimits(100, 1000),
			},
			UpdatePass: &updatePass{
				UpdatePoolOnA: false,
				UpdatePoolOnB: true,
				SelectorA2B:   createSymmetricRateLimits(200, 2000),
				SelectorB2A:   createSymmetricRateLimits(200, 2000),
			},
		},
	}

	for _, test := range tests {
		for _, mcmsConfig := range []*changeset.MCMSConfig{nil, &changeset.MCMSConfig{MinDelay: 0 * time.Second}} { // Run all tests with and without MCMS
			t.Run(test.Msg, func(t *testing.T) {
				e, selectorA, selectorB, tokens, timelockContracts := testhelpers.SetupTwoChainEnvironmentWithTokens(t, logger.TestLogger(t), mcmsConfig != nil)

				e = testhelpers.DeployTestTokenPools(t, e, map[uint64]v1_5_1.DeployTokenPoolInput{
					selectorA: {
						Type:               changeset.BurnMintTokenPool,
						TokenAddress:       tokens[selectorA].Address,
						LocalTokenDecimals: testhelpers.LocalTokenDecimals,
					},
					selectorB: {
						Type:               changeset.BurnMintTokenPool,
						TokenAddress:       tokens[selectorB].Address,
						LocalTokenDecimals: testhelpers.LocalTokenDecimals,
					},
				}, mcmsConfig != nil)

				e = testhelpers.DeployTestTokenPools(t, e, map[uint64]v1_5_1.DeployTokenPoolInput{
					selectorA: {
						Type:               changeset.LockReleaseTokenPool,
						TokenAddress:       tokens[selectorA].Address,
						LocalTokenDecimals: testhelpers.LocalTokenDecimals,
						AcceptLiquidity:    &acceptLiquidity,
					},
					selectorB: {
						Type:               changeset.LockReleaseTokenPool,
						TokenAddress:       tokens[selectorB].Address,
						LocalTokenDecimals: testhelpers.LocalTokenDecimals,
						AcceptLiquidity:    &acceptLiquidity,
					},
				}, mcmsConfig != nil)

				state, err := changeset.LoadOnchainState(e)
				require.NoError(t, err)

				lockReleaseA, _ := token_pool.NewTokenPool(state.Chains[selectorA].LockReleaseTokenPools[testhelpers.TestTokenSymbol][deployment.Version1_5_1].Address(), e.Chains[selectorA].Client)
				burnMintA, _ := token_pool.NewTokenPool(state.Chains[selectorA].BurnMintTokenPools[testhelpers.TestTokenSymbol][deployment.Version1_5_1].Address(), e.Chains[selectorA].Client)

				lockReleaseB, _ := token_pool.NewTokenPool(state.Chains[selectorB].LockReleaseTokenPools[testhelpers.TestTokenSymbol][deployment.Version1_5_1].Address(), e.Chains[selectorB].Client)
				burnMintB, _ := token_pool.NewTokenPool(state.Chains[selectorB].BurnMintTokenPools[testhelpers.TestTokenSymbol][deployment.Version1_5_1].Address(), e.Chains[selectorB].Client)

				pools := map[uint64]tokenPools{
					selectorA: tokenPools{
						LockRelease: lockReleaseA,
						BurnMint:    burnMintA,
					},
					selectorB: tokenPools{
						LockRelease: lockReleaseB,
						BurnMint:    burnMintB,
					},
				}
				expectedOwners := make(map[uint64]common.Address, 2)
				if mcmsConfig != nil {
					expectedOwners[selectorA] = state.Chains[selectorA].Timelock.Address()
					expectedOwners[selectorB] = state.Chains[selectorB].Timelock.Address()
				} else {
					expectedOwners[selectorA] = e.Chains[selectorA].DeployerKey.From
					expectedOwners[selectorB] = e.Chains[selectorB].DeployerKey.From
				}

				if test.RegistrationPass != nil {
					// Configure & set the active pools on the registry
					e, err = commonchangeset.Apply(t, e, timelockContracts,
						commonchangeset.Configure(
							deployment.CreateLegacyChangeSet(v1_5_1.ConfigureTokenPoolContractsChangeset),
							v1_5_1.ConfigureTokenPoolContractsConfig{
								TokenSymbol: testhelpers.TestTokenSymbol,
								MCMS:        mcmsConfig,
								PoolUpdates: map[uint64]v1_5_1.TokenPoolConfig{
									selectorA: {
										Type:    changeset.LockReleaseTokenPool,
										Version: deployment.Version1_5_1,
										ChainUpdates: v1_5_1.RateLimiterPerChain{
											selectorB: test.RegistrationPass.SelectorA2B,
										},
									},
									selectorB: {
										Type:    changeset.LockReleaseTokenPool,
										Version: deployment.Version1_5_1,
										ChainUpdates: v1_5_1.RateLimiterPerChain{
											selectorA: test.RegistrationPass.SelectorB2A,
										},
									},
								},
							},
						),
						commonchangeset.Configure(
							deployment.CreateLegacyChangeSet(v1_5_1.ProposeAdminRoleChangeset),
							changeset.TokenAdminRegistryChangesetConfig{
								MCMS: mcmsConfig,
								Pools: map[uint64]map[changeset.TokenSymbol]changeset.TokenPoolInfo{
									selectorA: {
										testhelpers.TestTokenSymbol: {
											Type:    changeset.LockReleaseTokenPool,
											Version: deployment.Version1_5_1,
										},
									},
									selectorB: {
										testhelpers.TestTokenSymbol: {
											Type:    changeset.LockReleaseTokenPool,
											Version: deployment.Version1_5_1,
										},
									},
								},
							},
						),
						commonchangeset.Configure(
							deployment.CreateLegacyChangeSet(v1_5_1.AcceptAdminRoleChangeset),
							changeset.TokenAdminRegistryChangesetConfig{
								MCMS: mcmsConfig,
								Pools: map[uint64]map[changeset.TokenSymbol]changeset.TokenPoolInfo{
									selectorA: {
										testhelpers.TestTokenSymbol: {
											Type:    changeset.LockReleaseTokenPool,
											Version: deployment.Version1_5_1,
										},
									},
									selectorB: {
										testhelpers.TestTokenSymbol: {
											Type:    changeset.LockReleaseTokenPool,
											Version: deployment.Version1_5_1,
										},
									},
								},
							},
						),
						commonchangeset.Configure(
							deployment.CreateLegacyChangeSet(v1_5_1.SetPoolChangeset),
							changeset.TokenAdminRegistryChangesetConfig{
								MCMS: mcmsConfig,
								Pools: map[uint64]map[changeset.TokenSymbol]changeset.TokenPoolInfo{
									selectorA: {
										testhelpers.TestTokenSymbol: {
											Type:    changeset.LockReleaseTokenPool,
											Version: deployment.Version1_5_1,
										},
									},
									selectorB: {
										testhelpers.TestTokenSymbol: {
											Type:    changeset.LockReleaseTokenPool,
											Version: deployment.Version1_5_1,
										},
									},
								},
							},
						),
					)
					require.NoError(t, err)

					for _, selector := range e.AllChainSelectors() {
						var remoteChainSelector uint64
						var rateLimiterConfig v1_5_1.RateLimiterConfig
						switch selector {
						case selectorA:
							remoteChainSelector = selectorB
							rateLimiterConfig = test.RegistrationPass.SelectorA2B
						case selectorB:
							remoteChainSelector = selectorA
							rateLimiterConfig = test.RegistrationPass.SelectorB2A
						}
						validateMemberOfTokenPoolPair(
							t,
							state,
							pools[selector].LockRelease,
							[]common.Address{pools[remoteChainSelector].LockRelease.Address()},
							tokens,
							testhelpers.TestTokenSymbol,
							selector,
							rateLimiterConfig.Inbound.Rate, // inbound & outbound are the same in this test
							rateLimiterConfig.Inbound.Capacity,
							expectedOwners[selector],
						)
					}
				}

				if test.UpdatePass != nil {
					// Only configure, do not update registry
					aType := changeset.LockReleaseTokenPool
					if test.UpdatePass.UpdatePoolOnA {
						aType = changeset.BurnMintTokenPool
					}
					bType := changeset.LockReleaseTokenPool
					if test.UpdatePass.UpdatePoolOnB {
						bType = changeset.BurnMintTokenPool
					}
					e, err = commonchangeset.Apply(t, e, timelockContracts,
						commonchangeset.Configure(
							deployment.CreateLegacyChangeSet(v1_5_1.ConfigureTokenPoolContractsChangeset),
							v1_5_1.ConfigureTokenPoolContractsConfig{
								TokenSymbol: testhelpers.TestTokenSymbol,
								MCMS:        mcmsConfig,
								PoolUpdates: map[uint64]v1_5_1.TokenPoolConfig{
									selectorA: {
										Type:    aType,
										Version: deployment.Version1_5_1,
										ChainUpdates: v1_5_1.RateLimiterPerChain{
											selectorB: test.UpdatePass.SelectorA2B,
										},
									},
									selectorB: {
										Type:    bType,
										Version: deployment.Version1_5_1,
										ChainUpdates: v1_5_1.RateLimiterPerChain{
											selectorA: test.UpdatePass.SelectorB2A,
										},
									},
								},
							},
						),
					)
					require.NoError(t, err)

					for _, selector := range e.AllChainSelectors() {
						var updatePool bool
						var updateRemotePool bool
						var remoteChainSelector uint64
						var rateLimiterConfig v1_5_1.RateLimiterConfig
						switch selector {
						case selectorA:
							remoteChainSelector = selectorB
							rateLimiterConfig = test.UpdatePass.SelectorA2B
							updatePool = test.UpdatePass.UpdatePoolOnA
							updateRemotePool = test.UpdatePass.UpdatePoolOnB
						case selectorB:
							remoteChainSelector = selectorA
							rateLimiterConfig = test.UpdatePass.SelectorB2A
							updatePool = test.UpdatePass.UpdatePoolOnB
							updateRemotePool = test.UpdatePass.UpdatePoolOnA
						}
						remotePoolAddresses := []common.Address{pools[remoteChainSelector].LockRelease.Address()} // add registered pool by default
						if updateRemotePool {                                                                     // if remote pool address is being updated, we push the new address
							remotePoolAddresses = append(remotePoolAddresses, pools[remoteChainSelector].BurnMint.Address())
						}
						tokenPool := pools[selector].LockRelease
						if updatePool {
							tokenPool = pools[selector].BurnMint
						}
						validateMemberOfTokenPoolPair(
							t,
							state,
							tokenPool,
							remotePoolAddresses,
							tokens,
							testhelpers.TestTokenSymbol,
							selector,
							rateLimiterConfig.Inbound.Rate, // inbound & outbound are the same in this test
							rateLimiterConfig.Inbound.Capacity,
							expectedOwners[selector],
						)
					}
				}
			})
		}
	}
}
