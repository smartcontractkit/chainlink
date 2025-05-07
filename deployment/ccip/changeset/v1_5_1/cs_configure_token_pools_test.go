package v1_5_1_test

import (
	"bytes"
	"math/big"
	"slices"
	"sort"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gagliardetto/solana-go"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_5_1/token_pool"
	solTestTokenPool "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/test_token_pool"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/shared/generated/burn_mint_erc677"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	changeset_solana "github.com/smartcontractkit/chainlink/deployment/ccip/changeset/solana"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/v1_5_1"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
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
	tokens map[uint64]*cldf.ContractDeploy[*burn_mint_erc677.BurnMintERC677],
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

func validateSolanaConfig(t *testing.T, state changeset.CCIPOnChainState, solChainUpdates map[uint64]v1_5_1.SolChainUpdate, selector uint64, solanaSelector uint64) {
	tokenPool := state.Chains[selector].BurnMintTokenPools[testhelpers.TestTokenSymbol][deployment.Version1_5_1]
	isSupported, err := tokenPool.IsSupportedChain(nil, solanaSelector)
	require.NoError(t, err)
	require.True(t, isSupported)

	remoteToken, remoteTokenPool, err := solChainUpdates[solanaSelector].GetSolanaTokenAndTokenPool(state.SolChains[solanaSelector])
	require.NoError(t, err)
	remoteTokenAddress, err := tokenPool.GetRemoteToken(nil, solanaSelector)
	require.NoError(t, err)
	require.Equal(t, remoteToken.Bytes(), remoteTokenAddress)
	remotePoolAddresses, err := tokenPool.GetRemotePools(nil, solanaSelector)
	require.NoError(t, err)
	require.Len(t, remotePoolAddresses, 1)
	require.Equal(t, remoteTokenPool.Bytes(), remotePoolAddresses[0])

	inboundRateLimiterConfig, err := tokenPool.GetCurrentInboundRateLimiterState(nil, solanaSelector)
	require.NoError(t, err)
	require.Equal(t, solChainUpdates[solanaSelector].RateLimiterConfig.Inbound.Rate.Int64(), inboundRateLimiterConfig.Rate.Int64())
	require.Equal(t, solChainUpdates[solanaSelector].RateLimiterConfig.Inbound.Capacity.Int64(), inboundRateLimiterConfig.Capacity.Int64())
	require.Equal(t, solChainUpdates[solanaSelector].RateLimiterConfig.Inbound.IsEnabled, inboundRateLimiterConfig.IsEnabled)

	outboundRateLimiterConfig, err := tokenPool.GetCurrentOutboundRateLimiterState(nil, solanaSelector)
	require.NoError(t, err)
	require.Equal(t, solChainUpdates[solanaSelector].RateLimiterConfig.Outbound.Rate.Int64(), outboundRateLimiterConfig.Rate.Int64())
	require.Equal(t, solChainUpdates[solanaSelector].RateLimiterConfig.Outbound.Capacity.Int64(), outboundRateLimiterConfig.Capacity.Int64())
	require.Equal(t, solChainUpdates[solanaSelector].RateLimiterConfig.Outbound.IsEnabled, outboundRateLimiterConfig.IsEnabled)
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
			err := test.TokenPoolConfig.Validate(e.GetContext(), e.Chains[selectorA], state, test.UseMcms, testhelpers.TestTokenSymbol)
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
		for _, mcmsConfig := range []*proposalutils.TimelockConfig{nil, {MinDelay: 0 * time.Second}} { // Run all tests with and without MCMS
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
							cldf.CreateLegacyChangeSet(v1_5_1.ConfigureTokenPoolContractsChangeset),
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
							cldf.CreateLegacyChangeSet(v1_5_1.ProposeAdminRoleChangeset),
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
							cldf.CreateLegacyChangeSet(v1_5_1.AcceptAdminRoleChangeset),
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
							cldf.CreateLegacyChangeSet(v1_5_1.SetPoolChangeset),
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
							cldf.CreateLegacyChangeSet(v1_5_1.ConfigureTokenPoolContractsChangeset),
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

func TestValidateConfigureTokenPoolContractsForSolana(t *testing.T) {
	t.Parallel()
	var err error

	deployedEnvironment, _ := testhelpers.NewMemoryEnvironment(t, func(testCfg *testhelpers.TestConfigs) {
		testCfg.Chains = 2
		testCfg.SolChains = 1
	})
	e := deployedEnvironment.Env

	evmSelectors := []uint64{e.AllChainSelectors()[0]}
	solanaSelectors := e.AllChainSelectorsSolana()

	addressBook := deployment.NewMemoryAddressBook()

	///////////////////////////
	// DEPLOY EVM TOKEN POOL //
	///////////////////////////
	for _, selector := range evmSelectors {
		token, err := cldf.DeployContract(e.Logger, e.Chains[selector], addressBook,
			func(chain deployment.Chain) cldf.ContractDeploy[*burn_mint_erc677.BurnMintERC677] {
				tokenAddress, tx, token, err := burn_mint_erc677.DeployBurnMintERC677(
					e.Chains[selector].DeployerKey,
					e.Chains[selector].Client,
					string(testhelpers.TestTokenSymbol),
					string(testhelpers.TestTokenSymbol),
					testhelpers.LocalTokenDecimals,
					big.NewInt(0).Mul(big.NewInt(1e9), big.NewInt(1e18)),
				)
				return cldf.ContractDeploy[*burn_mint_erc677.BurnMintERC677]{
					Address:  tokenAddress,
					Contract: token,
					Tv:       deployment.NewTypeAndVersion(changeset.BurnMintToken, deployment.Version1_0_0),
					Tx:       tx,
					Err:      err,
				}
			},
		)
		require.NoError(t, err)
		e = testhelpers.DeployTestTokenPools(t, e, map[uint64]v1_5_1.DeployTokenPoolInput{
			selector: {
				Type:               changeset.BurnMintTokenPool,
				TokenAddress:       token.Address,
				LocalTokenDecimals: testhelpers.LocalTokenDecimals,
			},
		}, false)
	}

	//////////////////////////////
	// DEPLOY SOLANA TOKEN POOL //
	//////////////////////////////
	for _, selector := range solanaSelectors {
		e, err = commonchangeset.Apply(t, e, nil,
			commonchangeset.Configure(
				cldf.CreateLegacyChangeSet(changeset_solana.DeploySolanaToken),
				changeset_solana.DeploySolanaTokenConfig{
					ChainSelector:    selector,
					TokenProgramName: changeset.SPL2022Tokens,
					TokenDecimals:    testhelpers.LocalTokenDecimals,
					TokenSymbol:      string(testhelpers.TestTokenSymbol),
				},
			),
		)
		require.NoError(t, err)
		state, err := changeset.LoadOnchainState(e)
		require.NoError(t, err)
		tokenAddress := state.SolChains[selector].SPL2022Tokens[0]
		e, _, err = commonchangeset.ApplyChangesetsV2(t, e, []commonchangeset.ConfiguredChangeSet{
			commonchangeset.Configure(
				cldf.CreateLegacyChangeSet(changeset_solana.AddTokenPoolAndLookupTable),
				changeset_solana.TokenPoolConfig{
					ChainSelector: selector,
					TokenPubKey:   tokenAddress,
					PoolType:      solTestTokenPool.BurnAndMint_PoolType,
				},
			),
		})
		require.NoError(t, err)
	}

	state, err := changeset.LoadOnchainState(e)
	require.NoError(t, err)

	/////////////////////////////
	// ADD SOLANA CHAIN CONFIG //
	/////////////////////////////
	for _, selector := range evmSelectors {
		solChainUpdates := make(map[uint64]v1_5_1.SolChainUpdate)
		for _, remoteSelector := range solanaSelectors {
			solChainUpdates[remoteSelector] = v1_5_1.SolChainUpdate{
				Type:              changeset.BurnMintTokenPool,
				TokenAddress:      state.SolChains[remoteSelector].SPL2022Tokens[0].String(),
				RateLimiterConfig: testhelpers.CreateSymmetricRateLimits(0, 0),
			}
		}
		e, err = commonchangeset.Apply(t, e, nil,
			commonchangeset.Configure(
				cldf.CreateLegacyChangeSet(v1_5_1.ConfigureTokenPoolContractsChangeset),
				v1_5_1.ConfigureTokenPoolContractsConfig{
					TokenSymbol: testhelpers.TestTokenSymbol,
					PoolUpdates: map[uint64]v1_5_1.TokenPoolConfig{
						selector: {
							Type:            changeset.BurnMintTokenPool,
							Version:         deployment.Version1_5_1,
							SolChainUpdates: solChainUpdates,
						},
					},
				},
			),
		)
		require.NoError(t, err)

		for _, remoteSelector := range solanaSelectors {
			validateSolanaConfig(t, state, solChainUpdates, selector, remoteSelector)
		}
	}

	////////////////////////////////
	// UPDATE SOLANA CHAIN CONFIG //
	////////////////////////////////
	for _, selector := range evmSelectors {
		solChainUpdates := make(map[uint64]v1_5_1.SolChainUpdate)
		for _, remoteSelector := range solanaSelectors {
			solChainUpdates[remoteSelector] = v1_5_1.SolChainUpdate{
				Type:              changeset.BurnMintTokenPool,
				TokenAddress:      state.SolChains[remoteSelector].SPL2022Tokens[0].String(),
				RateLimiterConfig: testhelpers.CreateSymmetricRateLimits(100, 1000),
			}
		}
		e.Chains[selector].DeployerKey.GasLimit = 1_000_000 // Hack: Increase gas limit to avoid out of gas error (could this be a cause for test flakiness?)
		e, err = commonchangeset.Apply(t, e, nil,
			commonchangeset.Configure(
				cldf.CreateLegacyChangeSet(v1_5_1.ConfigureTokenPoolContractsChangeset),
				v1_5_1.ConfigureTokenPoolContractsConfig{
					TokenSymbol: testhelpers.TestTokenSymbol,
					PoolUpdates: map[uint64]v1_5_1.TokenPoolConfig{
						selector: {
							Type:            changeset.BurnMintTokenPool,
							Version:         deployment.Version1_5_1,
							SolChainUpdates: solChainUpdates,
						},
					},
				},
			),
		)
		require.NoError(t, err)

		for _, remoteSelector := range solanaSelectors {
			validateSolanaConfig(t, state, solChainUpdates, selector, remoteSelector)
		}
	}

	///////////////////////////
	// REDEPLOY SOLANA TOKEN //
	///////////////////////////
	remoteTokenAddresses := make(map[uint64]solana.PublicKey, len(solanaSelectors))
	for _, selector := range solanaSelectors {
		tokensBefore := state.SolChains[selector].SPL2022Tokens
		e, err = commonchangeset.Apply(t, e, nil,
			commonchangeset.Configure(
				cldf.CreateLegacyChangeSet(changeset_solana.DeploySolanaToken),
				changeset_solana.DeploySolanaTokenConfig{
					ChainSelector:    selector,
					TokenProgramName: changeset.SPL2022Tokens,
					TokenDecimals:    testhelpers.LocalTokenDecimals,
					TokenSymbol:      string(testhelpers.TestTokenSymbol),
				},
			),
		)
		require.NoError(t, err)
		onchainState, err := changeset.LoadOnchainState(e)
		require.NoError(t, err)
		for _, tokenAddress := range onchainState.SolChains[selector].SPL2022Tokens {
			if slices.Contains(tokensBefore, tokenAddress) {
				continue
			}
			e, _, err = commonchangeset.ApplyChangesetsV2(t, e, []commonchangeset.ConfiguredChangeSet{
				commonchangeset.Configure(
					cldf.CreateLegacyChangeSet(changeset_solana.AddTokenPoolAndLookupTable),
					changeset_solana.TokenPoolConfig{
						ChainSelector: selector,
						TokenPubKey:   tokenAddress,
						PoolType:      solTestTokenPool.BurnAndMint_PoolType,
					},
				),
			})
			require.NoError(t, err)
			remoteTokenAddresses[selector] = tokenAddress
		}
	}

	//////////////////////////////////////
	// REMOVE & ADD SOLANA CHAIN CONFIG //
	//////////////////////////////////////
	for _, selector := range evmSelectors {
		solChainUpdates := make(map[uint64]v1_5_1.SolChainUpdate)
		for remoteSelector, remoteTokenAddress := range remoteTokenAddresses {
			solChainUpdates[remoteSelector] = v1_5_1.SolChainUpdate{
				Type:              changeset.BurnMintTokenPool,
				TokenAddress:      remoteTokenAddress.String(),
				RateLimiterConfig: testhelpers.CreateSymmetricRateLimits(0, 0),
			}
		}
		e, err = commonchangeset.Apply(t, e, nil,
			commonchangeset.Configure(
				cldf.CreateLegacyChangeSet(v1_5_1.ConfigureTokenPoolContractsChangeset),
				v1_5_1.ConfigureTokenPoolContractsConfig{
					TokenSymbol: testhelpers.TestTokenSymbol,
					PoolUpdates: map[uint64]v1_5_1.TokenPoolConfig{
						selector: {
							Type:            changeset.BurnMintTokenPool,
							Version:         deployment.Version1_5_1,
							SolChainUpdates: solChainUpdates,
						},
					},
				},
			),
		)
		require.NoError(t, err)

		for _, remoteSelector := range solanaSelectors {
			validateSolanaConfig(t, state, solChainUpdates, selector, remoteSelector)
		}
	}
}
