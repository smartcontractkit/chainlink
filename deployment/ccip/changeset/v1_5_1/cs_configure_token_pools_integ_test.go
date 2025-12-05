//go:build integration && solana

package v1_5_1_test

import (
	"math/big"
	"slices"
	"testing"

	"github.com/gagliardetto/solana-go"
	chain_selectors "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/quarantine"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-evm/gethwrappers/shared/generated/initial/burn_mint_erc677"

	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"
	cldf_evm "github.com/smartcontractkit/chainlink-deployments-framework/chain/evm"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment"
	changeset_solana "github.com/smartcontractkit/chainlink/deployment/ccip/changeset/solana_v0_1_1"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/v1_5_1"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"

	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
)

func validateSolanaConfig(t *testing.T, state stateview.CCIPOnChainState, solChainUpdates map[uint64]v1_5_1.SolChainUpdate, selector uint64, solanaSelector uint64) {
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

func TestValidateConfigureTokenPoolContractsForSolana(t *testing.T) {
	quarantine.Flaky(t, "DX-1726")
	t.Parallel()
	var err error

	deployedEnvironment, _ := testhelpers.NewMemoryEnvironment(t, func(testCfg *testhelpers.TestConfigs) {
		testCfg.Chains = 2
		testCfg.SolChains = 1
	})
	e := deployedEnvironment.Env

	evmSelectors := []uint64{e.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chain_selectors.FamilyEVM))[0]}
	solanaSelectors := e.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chain_selectors.FamilySolana))

	addressBook := cldf.NewMemoryAddressBook()

	///////////////////////////
	// DEPLOY EVM TOKEN POOL //
	///////////////////////////
	for _, selector := range evmSelectors {
		token, err := cldf.DeployContract(e.Logger, e.BlockChains.EVMChains()[selector], addressBook,
			func(chain cldf_evm.Chain) cldf.ContractDeploy[*burn_mint_erc677.BurnMintERC677] {
				tokenAddress, tx, token, err := burn_mint_erc677.DeployBurnMintERC677(
					e.BlockChains.EVMChains()[selector].DeployerKey,
					e.BlockChains.EVMChains()[selector].Client,
					string(testhelpers.TestTokenSymbol),
					string(testhelpers.TestTokenSymbol),
					testhelpers.LocalTokenDecimals,
					big.NewInt(0).Mul(big.NewInt(1e9), big.NewInt(1e18)),
				)
				return cldf.ContractDeploy[*burn_mint_erc677.BurnMintERC677]{
					Address:  tokenAddress,
					Contract: token,
					Tv:       cldf.NewTypeAndVersion(shared.BurnMintToken, deployment.Version1_0_0),
					Tx:       tx,
					Err:      err,
				}
			},
		)
		require.NoError(t, err)
		e = testhelpers.DeployTestTokenPools(t, e, map[uint64]v1_5_1.DeployTokenPoolInput{
			selector: {
				Type:               shared.BurnMintTokenPool,
				TokenAddress:       token.Address,
				LocalTokenDecimals: testhelpers.LocalTokenDecimals,
			},
		}, false)
	}

	//////////////////////////////
	// DEPLOY SOLANA TOKEN POOL //
	//////////////////////////////
	for _, selector := range solanaSelectors {
		e, err = commonchangeset.Apply(t, e,
			commonchangeset.Configure(
				cldf.CreateLegacyChangeSet(changeset_solana.DeploySolanaToken),
				changeset_solana.DeploySolanaTokenConfig{
					ChainSelector:    selector,
					TokenProgramName: shared.SPL2022Tokens,
					TokenDecimals:    testhelpers.LocalTokenDecimals,
					TokenSymbol:      string(testhelpers.TestTokenSymbol),
				},
			),
		)
		require.NoError(t, err)
		state, err := stateview.LoadOnchainState(e)
		require.NoError(t, err)
		tokenAddress := state.SolChains[selector].SPL2022Tokens[0]
		e, _, err = commonchangeset.ApplyChangesets(t, e, []commonchangeset.ConfiguredChangeSet{
			commonchangeset.Configure(
				cldf.CreateLegacyChangeSet(changeset_solana.InitGlobalConfigTokenPoolProgram),
				changeset_solana.TokenPoolConfigWithMCM{
					ChainSelector: selector,
					TokenPoolConfigs: []changeset_solana.TokenPoolConfig{
						{
							TokenPubKey: tokenAddress,
							PoolType:    shared.BurnMintTokenPool,
							Metadata:    shared.CLLMetadata,
						},
					},
				},
			),
			commonchangeset.Configure(
				cldf.CreateLegacyChangeSet(changeset_solana.AddTokenPoolAndLookupTable),
				changeset_solana.AddTokenPoolAndLookupTableConfig{
					ChainSelector: selector,
					TokenPoolConfigs: []changeset_solana.TokenPoolConfig{
						{
							TokenPubKey: tokenAddress,
							PoolType:    shared.BurnMintTokenPool,
							Metadata:    shared.CLLMetadata,
						},
					},
				},
			),
		})
		require.NoError(t, err)
	}

	state, err := stateview.LoadOnchainState(e)
	require.NoError(t, err)

	/////////////////////////////
	// ADD SOLANA CHAIN CONFIG //
	/////////////////////////////
	for _, selector := range evmSelectors {
		solChainUpdates := make(map[uint64]v1_5_1.SolChainUpdate)
		for _, remoteSelector := range solanaSelectors {
			solChainUpdates[remoteSelector] = v1_5_1.SolChainUpdate{
				Type:              shared.BurnMintTokenPool,
				TokenAddress:      state.SolChains[remoteSelector].SPL2022Tokens[0].String(),
				RateLimiterConfig: testhelpers.CreateSymmetricRateLimits(0, 0),
				Metadata:          shared.CLLMetadata,
			}
		}
		e, err = commonchangeset.Apply(t, e,
			commonchangeset.Configure(
				cldf.CreateLegacyChangeSet(v1_5_1.ConfigureTokenPoolContractsChangeset),
				v1_5_1.ConfigureTokenPoolContractsConfig{
					TokenSymbol: testhelpers.TestTokenSymbol,
					PoolUpdates: map[uint64]v1_5_1.TokenPoolConfig{
						selector: {
							Type:            shared.BurnMintTokenPool,
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
				Type:              shared.BurnMintTokenPool,
				TokenAddress:      state.SolChains[remoteSelector].SPL2022Tokens[0].String(),
				RateLimiterConfig: testhelpers.CreateSymmetricRateLimits(100, 1000),
				Metadata:          shared.CLLMetadata,
			}
		}
		e.BlockChains.EVMChains()[selector].DeployerKey.GasLimit = 1_000_000 // Hack: Increase gas limit to avoid out of gas error (could this be a cause for test flakiness?)
		e, err = commonchangeset.Apply(t, e,
			commonchangeset.Configure(
				cldf.CreateLegacyChangeSet(v1_5_1.ConfigureTokenPoolContractsChangeset),
				v1_5_1.ConfigureTokenPoolContractsConfig{
					TokenSymbol: testhelpers.TestTokenSymbol,
					PoolUpdates: map[uint64]v1_5_1.TokenPoolConfig{
						selector: {
							Type:            shared.BurnMintTokenPool,
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
		e, err = commonchangeset.Apply(t, e,
			commonchangeset.Configure(
				cldf.CreateLegacyChangeSet(changeset_solana.DeploySolanaToken),
				changeset_solana.DeploySolanaTokenConfig{
					ChainSelector:    selector,
					TokenProgramName: shared.SPL2022Tokens,
					TokenDecimals:    testhelpers.LocalTokenDecimals,
					TokenSymbol:      string(testhelpers.TestTokenSymbol),
				},
			),
		)
		require.NoError(t, err)
		onchainState, err := stateview.LoadOnchainState(e)
		require.NoError(t, err)
		for _, tokenAddress := range onchainState.SolChains[selector].SPL2022Tokens {
			if slices.Contains(tokensBefore, tokenAddress) {
				continue
			}
			e, _, err = commonchangeset.ApplyChangesets(t, e, []commonchangeset.ConfiguredChangeSet{
				commonchangeset.Configure(
					cldf.CreateLegacyChangeSet(changeset_solana.AddTokenPoolAndLookupTable),
					changeset_solana.AddTokenPoolAndLookupTableConfig{
						ChainSelector: selector,
						TokenPoolConfigs: []changeset_solana.TokenPoolConfig{
							{
								TokenPubKey: tokenAddress,
								PoolType:    shared.BurnMintTokenPool,
								Metadata:    shared.CLLMetadata,
							},
						},
					},
				),
			})
			require.NoError(t, err)
			remoteTokenAddresses[selector] = tokenAddress
		}
	}

	////////////////////////////////////////////////////////////
	// REMOVE & ADD SOLANA CHAIN CONFIG (due to token change) //
	////////////////////////////////////////////////////////////
	for _, selector := range evmSelectors {
		solChainUpdates := make(map[uint64]v1_5_1.SolChainUpdate)
		for remoteSelector, remoteTokenAddress := range remoteTokenAddresses {
			solChainUpdates[remoteSelector] = v1_5_1.SolChainUpdate{
				Type:              shared.BurnMintTokenPool,
				TokenAddress:      remoteTokenAddress.String(),
				RateLimiterConfig: testhelpers.CreateSymmetricRateLimits(0, 0),
				Metadata:          shared.CLLMetadata,
			}
		}
		e, err = commonchangeset.Apply(t, e,
			commonchangeset.Configure(
				cldf.CreateLegacyChangeSet(v1_5_1.ConfigureTokenPoolContractsChangeset),
				v1_5_1.ConfigureTokenPoolContractsConfig{
					TokenSymbol: testhelpers.TestTokenSymbol,
					PoolUpdates: map[uint64]v1_5_1.TokenPoolConfig{
						selector: {
							Type:            shared.BurnMintTokenPool,
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

	//////////////////////////////////
	// DEPLOY NEW SOLANA TOKEN POOL //
	//////////////////////////////////
	require.NoError(t, err)
	for _, selector := range solanaSelectors {
		for _, tokenAddress := range remoteTokenAddresses {
			e, _, err = commonchangeset.ApplyChangesets(t, e, []commonchangeset.ConfiguredChangeSet{
				commonchangeset.Configure(
					cldf.CreateLegacyChangeSet(changeset_solana.InitGlobalConfigTokenPoolProgram),
					changeset_solana.TokenPoolConfigWithMCM{
						ChainSelector: selector,
						TokenPoolConfigs: []changeset_solana.TokenPoolConfig{
							{
								TokenPubKey: tokenAddress,
								PoolType:    shared.LockReleaseTokenPool,
								Metadata:    shared.CLLMetadata,
							},
						},
					},
				),
				commonchangeset.Configure(
					cldf.CreateLegacyChangeSet(changeset_solana.AddTokenPoolAndLookupTable),
					changeset_solana.AddTokenPoolAndLookupTableConfig{
						ChainSelector: selector,
						TokenPoolConfigs: []changeset_solana.TokenPoolConfig{
							{
								TokenPubKey: tokenAddress,
								PoolType:    shared.LockReleaseTokenPool,
								Metadata:    shared.CLLMetadata,
							},
						},
					},
				),
			})
			require.NoError(t, err)
		}
	}

	/////////////////////////////////////////////////////////////////
	// REMOVE & ADD SOLANA CHAIN CONFIG (due to token pool change) //
	/////////////////////////////////////////////////////////////////
	for _, selector := range evmSelectors {
		solChainUpdates := make(map[uint64]v1_5_1.SolChainUpdate)
		for remoteSelector, remoteTokenAddress := range remoteTokenAddresses {
			solChainUpdates[remoteSelector] = v1_5_1.SolChainUpdate{
				Type:              shared.LockReleaseTokenPool,
				TokenAddress:      remoteTokenAddress.String(),
				RateLimiterConfig: testhelpers.CreateSymmetricRateLimits(0, 0),
				Metadata:          shared.CLLMetadata,
			}
		}
		e, err = commonchangeset.Apply(t, e,
			commonchangeset.Configure(
				cldf.CreateLegacyChangeSet(v1_5_1.ConfigureTokenPoolContractsChangeset),
				v1_5_1.ConfigureTokenPoolContractsConfig{
					TokenSymbol: testhelpers.TestTokenSymbol,
					PoolUpdates: map[uint64]v1_5_1.TokenPoolConfig{
						selector: {
							Type:            shared.BurnMintTokenPool,
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
