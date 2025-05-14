package solana_test

import (
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gagliardetto/solana-go"
	solRpc "github.com/gagliardetto/solana-go/rpc"
	"github.com/stretchr/testify/require"

	solBaseTokenPool "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/base_token_pool"
	solTestTokenPool "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/test_token_pool"
	solTokenUtil "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/tokens"

	"github.com/smartcontractkit/chainlink-evm/gethwrappers/shared/generated/burn_mint_erc677"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/testcontext"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	ccipChangesetSolana "github.com/smartcontractkit/chainlink/deployment/ccip/changeset/solana"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/v1_5_1"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"

	"github.com/smartcontractkit/chainlink/deployment"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
)

func TestAddTokenPoolWithoutMcms(t *testing.T) {
	t.Parallel()
	tenv, _ := testhelpers.NewMemoryEnvironment(t, testhelpers.WithSolChains(1))
	doTestTokenPool(t, tenv.Env, false, shared.CLLMetadata)
}

func TestAddTokenPoolWithMcms(t *testing.T) {
	t.Parallel()
	tenv, _ := testhelpers.NewMemoryEnvironment(t, testhelpers.WithSolChains(1))
	doTestTokenPool(t, tenv.Env, true, shared.CLLMetadata)
}

func deployEVMTokenPool(t *testing.T, e cldf.Environment, evmChain uint64) (cldf.Environment, common.Address, error) {
	addressBook := cldf.NewMemoryAddressBook()
	evmToken, err := cldf.DeployContract(e.Logger, e.Chains[evmChain], addressBook,
		func(chain cldf.Chain) cldf.ContractDeploy[*burn_mint_erc677.BurnMintERC677] {
			tokenAddress, tx, token, err := burn_mint_erc677.DeployBurnMintERC677(
				e.Chains[evmChain].DeployerKey,
				e.Chains[evmChain].Client,
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
		evmChain: {
			Type:               shared.BurnMintTokenPool,
			TokenAddress:       evmToken.Address,
			LocalTokenDecimals: testhelpers.LocalTokenDecimals,
		},
	}, false)
	return e, evmToken.Address, nil
}

func doTestTokenPool(t *testing.T, e cldf.Environment, mcms bool, tokenMetadata string) {
	ctx := testcontext.Get(t)
	evmChain := e.AllChainSelectors()[0]
	solChain := e.AllChainSelectorsSolana()[0]
	deployerKey := e.SolChains[solChain].DeployerKey.PublicKey()
	testUser, _ := solana.NewRandomPrivateKey()
	testUserPubKey := testUser.PublicKey()
	e, newTokenAddress, err := deployTokenAndMint(t, e, solChain, []string{deployerKey.String(), testUserPubKey.String()})
	require.NoError(t, err)
	state, err := stateview.LoadOnchainStateSolana(e)
	require.NoError(t, err)
	testUserATA, _, err := solTokenUtil.FindAssociatedTokenAddress(solana.TokenProgramID, newTokenAddress, testUserPubKey)
	require.NoError(t, err)
	deployerATA, _, err := solTokenUtil.FindAssociatedTokenAddress(
		solana.TokenProgramID,
		newTokenAddress,
		e.SolChains[solChain].DeployerKey.PublicKey(),
	)
	var mcmsConfig *proposalutils.TimelockConfig
	if mcms {
		_, _ = testhelpers.TransferOwnershipSolana(t, &e, solChain, true,
			ccipChangesetSolana.CCIPContractsToTransfer{
				Router:    true,
				FeeQuoter: true,
				OffRamp:   true,
			})
		mcmsConfig = &proposalutils.TimelockConfig{
			MinDelay: 1 * time.Second,
		}
	}
	require.NoError(t, err)

	rateLimitConfig := solBaseTokenPool.RateLimitConfig{
		Enabled:  false,
		Capacity: 0,
		Rate:     0,
	}
	inboundConfig := rateLimitConfig
	outboundConfig := rateLimitConfig

	type poolTestType struct {
		poolType    solTestTokenPool.PoolType
		poolAddress solana.PublicKey
	}
	testCases := []poolTestType{
		{
			poolType:    solTestTokenPool.BurnAndMint_PoolType,
			poolAddress: state.SolChains[solChain].BurnMintTokenPools[tokenMetadata],
		},
		{
			poolType:    solTestTokenPool.LockAndRelease_PoolType,
			poolAddress: state.SolChains[solChain].LockReleaseTokenPools[tokenMetadata],
		},
	}

	// evm deployment
	e, _, err = deployEVMTokenPool(t, e, evmChain)
	require.NoError(t, err)

	tokenAddress := newTokenAddress

	for _, testCase := range testCases {
		typePtr := &testCase.poolType
		// for _, tokenAddress := range tokenMap {
		e, _, err = commonchangeset.ApplyChangesetsV2(t, e, []commonchangeset.ConfiguredChangeSet{
			commonchangeset.Configure(
				cldf.CreateLegacyChangeSet(ccipChangesetSolana.AddTokenPoolAndLookupTable),
				ccipChangesetSolana.TokenPoolConfig{
					ChainSelector: solChain,
					TokenPubKey:   tokenAddress,
					PoolType:      typePtr,
					Metadata:      tokenMetadata,
				},
			),
			commonchangeset.Configure(
				cldf.CreateLegacyChangeSet(ccipChangesetSolana.SetupTokenPoolForRemoteChain),
				ccipChangesetSolana.RemoteChainTokenPoolConfig{
					SolChainSelector: solChain,
					SolTokenPubKey:   tokenAddress,
					SolPoolType:      typePtr,
					Metadata:         tokenMetadata,
					EVMRemoteConfigs: map[uint64]ccipChangesetSolana.EVMRemoteConfig{
						evmChain: {
							TokenSymbol: testhelpers.TestTokenSymbol,
							PoolType:    shared.BurnMintTokenPool, // EVM test tokens are always burn and mint
							PoolVersion: shared.CurrentTokenPoolVersion,
							RateLimiterConfig: ccipChangesetSolana.RateLimiterConfig{
								Inbound:  rateLimitConfig,
								Outbound: rateLimitConfig,
							},
						},
					},
					MCMS: mcmsConfig,
				},
			),
		})
		require.NoError(t, err)
		// test AddTokenPool results
		configAccount := solTestTokenPool.State{}
		poolConfigPDA, _ := solTokenUtil.TokenPoolConfigAddress(tokenAddress, testCase.poolAddress)
		err = e.SolChains[solChain].GetAccountDataBorshInto(ctx, poolConfigPDA, &configAccount)
		require.NoError(t, err)
		require.Equal(t, tokenAddress, configAccount.Config.Mint)
		// test SetupTokenPoolForRemoteChain results
		remoteChainConfigPDA, _, _ := solTokenUtil.TokenPoolChainConfigPDA(evmChain, tokenAddress, testCase.poolAddress)
		var remoteChainConfigAccount solTestTokenPool.ChainConfig
		err = e.SolChains[solChain].GetAccountDataBorshInto(ctx, remoteChainConfigPDA, &remoteChainConfigAccount)
		require.NoError(t, err)
		require.Equal(t, testhelpers.LocalTokenDecimals, int(remoteChainConfigAccount.Base.Remote.Decimals))
		e.Logger.Infof("Pool addresses: %v", remoteChainConfigAccount.Base.Remote.PoolAddresses)
		require.Len(t, remoteChainConfigAccount.Base.Remote.PoolAddresses, 1)
		require.Equal(t, inboundConfig.Enabled, remoteChainConfigAccount.Base.InboundRateLimit.Cfg.Enabled)
		require.Equal(t, outboundConfig.Enabled, remoteChainConfigAccount.Base.OutboundRateLimit.Cfg.Enabled)

		allowedAccount1, _ := solana.NewRandomPrivateKey()
		allowedAccount2, _ := solana.NewRandomPrivateKey()

		newRateLimitConfig := solBaseTokenPool.RateLimitConfig{
			Enabled:  true,
			Capacity: uint64(1000),
			Rate:     1,
		}
		newOutboundConfig := newRateLimitConfig
		newInboundConfig := newRateLimitConfig

		if mcms {
			e.Logger.Debugf("Configuring MCMS for token pool %v", testCase.poolType)
			if testCase.poolType == solTestTokenPool.BurnAndMint_PoolType {
				_, _ = testhelpers.TransferOwnershipSolana(
					t, &e, solChain, false,
					ccipChangesetSolana.CCIPContractsToTransfer{
						BurnMintTokenPoolMetadata: tokenMetadata,
						BurnMintTokenPools: map[solana.PublicKey]solana.PublicKey{
							poolConfigPDA: tokenAddress,
						},
					})
			} else if testCase.poolType == solTestTokenPool.LockAndRelease_PoolType {
				_, _ = testhelpers.TransferOwnershipSolana(
					t, &e, solChain, false,
					ccipChangesetSolana.CCIPContractsToTransfer{
						LockReleaseTokenPoolMetadata: tokenMetadata,
						LockReleaseTokenPools: map[solana.PublicKey]solana.PublicKey{
							poolConfigPDA: tokenAddress,
						},
					})
			}
			e.Logger.Debugf("MCMS Configured for token pool %v with token address %v", testCase.poolType, tokenAddress)
		}

		e, _, err = commonchangeset.ApplyChangesetsV2(t, e, []commonchangeset.ConfiguredChangeSet{
			commonchangeset.Configure(
				cldf.CreateLegacyChangeSet(ccipChangesetSolana.ConfigureTokenPoolAllowList),
				ccipChangesetSolana.ConfigureTokenPoolAllowListConfig{
					SolChainSelector: solChain,
					SolTokenPubKey:   tokenAddress.String(),
					PoolType:         typePtr,
					Metadata:         tokenMetadata,
					Accounts:         []solana.PublicKey{allowedAccount1.PublicKey(), allowedAccount2.PublicKey()},
					Enabled:          true,
					MCMS:             mcmsConfig,
				},
			),
			commonchangeset.Configure(
				cldf.CreateLegacyChangeSet(ccipChangesetSolana.RemoveFromTokenPoolAllowList),
				ccipChangesetSolana.RemoveFromAllowListConfig{
					SolChainSelector: solChain,
					SolTokenPubKey:   tokenAddress.String(),
					PoolType:         typePtr,
					Metadata:         tokenMetadata,
					Accounts:         []solana.PublicKey{allowedAccount1.PublicKey(), allowedAccount2.PublicKey()},
					MCMS:             mcmsConfig,
				},
			),
			// test update
			commonchangeset.Configure(
				cldf.CreateLegacyChangeSet(ccipChangesetSolana.SetupTokenPoolForRemoteChain),
				ccipChangesetSolana.RemoteChainTokenPoolConfig{
					SolChainSelector: solChain,
					SolTokenPubKey:   tokenAddress,
					SolPoolType:      typePtr,
					Metadata:         tokenMetadata,
					EVMRemoteConfigs: map[uint64]ccipChangesetSolana.EVMRemoteConfig{
						evmChain: {
							TokenSymbol: testhelpers.TestTokenSymbol,
							PoolType:    shared.BurnMintTokenPool, // EVM test tokens are always burn and mint
							PoolVersion: shared.CurrentTokenPoolVersion,
							RateLimiterConfig: ccipChangesetSolana.RateLimiterConfig{
								Inbound:  newInboundConfig,
								Outbound: newOutboundConfig,
							},
						},
					},
					MCMS: mcmsConfig,
				},
			),
		})
		require.NoError(t, err)

		err = e.SolChains[solChain].GetAccountDataBorshInto(ctx, remoteChainConfigPDA, &remoteChainConfigAccount)
		require.NoError(t, err)
		require.Equal(t, newInboundConfig.Enabled, remoteChainConfigAccount.Base.InboundRateLimit.Cfg.Enabled)
		require.Equal(t, newOutboundConfig.Enabled, remoteChainConfigAccount.Base.OutboundRateLimit.Cfg.Enabled)

		if testCase.poolType == solTestTokenPool.LockAndRelease_PoolType && tokenAddress == newTokenAddress {
			e, _, err = commonchangeset.ApplyChangesetsV2(t, e, []commonchangeset.ConfiguredChangeSet{
				commonchangeset.Configure(
					cldf.CreateLegacyChangeSet(ccipChangesetSolana.LockReleaseLiquidityOps),
					ccipChangesetSolana.LockReleaseLiquidityOpsConfig{
						SolChainSelector: solChain,
						SolTokenPubKey:   tokenAddress.String(),
						Metadata:         tokenMetadata,
						SetCfg: &ccipChangesetSolana.SetLiquidityConfig{
							Enabled: true,
						},
						MCMS: mcmsConfig,
					},
				),
				commonchangeset.Configure(
					cldf.CreateLegacyChangeSet(ccipChangesetSolana.LockReleaseLiquidityOps),
					ccipChangesetSolana.LockReleaseLiquidityOpsConfig{
						SolChainSelector: solChain,
						SolTokenPubKey:   tokenAddress.String(),
						Metadata:         tokenMetadata,
						LiquidityCfg: &ccipChangesetSolana.LiquidityConfig{
							Amount:             100,
							RemoteTokenAccount: deployerATA,
							Type:               ccipChangesetSolana.Provide,
						},
						MCMS: mcmsConfig,
					},
				),
				commonchangeset.Configure(
					cldf.CreateLegacyChangeSet(ccipChangesetSolana.LockReleaseLiquidityOps),
					ccipChangesetSolana.LockReleaseLiquidityOpsConfig{
						SolChainSelector: solChain,
						SolTokenPubKey:   tokenAddress.String(),
						Metadata:         tokenMetadata,
						LiquidityCfg: &ccipChangesetSolana.LiquidityConfig{
							Amount:             50,
							RemoteTokenAccount: testUserATA,
							Type:               ccipChangesetSolana.Withdraw,
						},
						MCMS: mcmsConfig,
					},
				),
			},
			)
			require.NoError(t, err)
			outDec, outVal, err := solTokenUtil.TokenBalance(e.GetContext(), e.SolChains[solChain].Client, deployerATA, solRpc.CommitmentConfirmed)
			require.NoError(t, err)
			require.Equal(t, int(900), outVal)
			require.Equal(t, 9, int(outDec))

			outDec, outVal, err = solTokenUtil.TokenBalance(e.GetContext(), e.SolChains[solChain].Client, testUserATA, solRpc.CommitmentConfirmed)
			require.NoError(t, err)
			require.Equal(t, int(1050), outVal)
			require.Equal(t, 9, int(outDec))

			err = e.SolChains[solChain].GetAccountDataBorshInto(ctx, poolConfigPDA, &configAccount)
			require.NoError(t, err)
			outDec, outVal, err = solTokenUtil.TokenBalance(e.GetContext(), e.SolChains[solChain].Client, configAccount.Config.PoolTokenAccount, solRpc.CommitmentConfirmed)
			require.NoError(t, err)
			require.Equal(t, int(50), outVal)
			require.Equal(t, 9, int(outDec))

			// transfer away from timelock if metadata is set and not ccipChangeset.CLLMetadata
			if mcms && tokenMetadata != "" && tokenMetadata != shared.CLLMetadata {
				timelockSignerPDA, err := ccipChangesetSolana.FetchTimelockSigner(e, solChain)
				require.NoError(t, err)
				e.Logger.Debugf("Transferring away from MCMS for token pool %v", testCase.poolType)
				if testCase.poolType == solTestTokenPool.BurnAndMint_PoolType {
					e, _, err = commonchangeset.ApplyChangesetsV2(t, e, []commonchangeset.ConfiguredChangeSet{
						commonchangeset.Configure(
							cldf.CreateLegacyChangeSet(ccipChangesetSolana.TransferCCIPToMCMSWithTimelockSolana),
							ccipChangesetSolana.TransferCCIPToMCMSWithTimelockSolanaConfig{
								MCMSCfg:       proposalutils.TimelockConfig{MinDelay: 1 * time.Second},
								CurrentOwner:  timelockSignerPDA,
								ProposedOwner: deployerKey,
								ContractsByChain: map[uint64]ccipChangesetSolana.CCIPContractsToTransfer{
									solChain: ccipChangesetSolana.CCIPContractsToTransfer{
										BurnMintTokenPoolMetadata: tokenMetadata,
										BurnMintTokenPools: map[solana.PublicKey]solana.PublicKey{
											poolConfigPDA: tokenAddress,
										},
									},
								},
							},
						),
					})
					require.NoError(t, err)
				} else if testCase.poolType == solTestTokenPool.LockAndRelease_PoolType {
					e, _, err = commonchangeset.ApplyChangesetsV2(t, e, []commonchangeset.ConfiguredChangeSet{
						commonchangeset.Configure(
							cldf.CreateLegacyChangeSet(ccipChangesetSolana.TransferCCIPToMCMSWithTimelockSolana),
							ccipChangesetSolana.TransferCCIPToMCMSWithTimelockSolanaConfig{
								MCMSCfg:       proposalutils.TimelockConfig{MinDelay: 1 * time.Second},
								CurrentOwner:  timelockSignerPDA,
								ProposedOwner: deployerKey,
								ContractsByChain: map[uint64]ccipChangesetSolana.CCIPContractsToTransfer{
									solChain: ccipChangesetSolana.CCIPContractsToTransfer{
										LockReleaseTokenPoolMetadata: tokenMetadata,
										LockReleaseTokenPools: map[solana.PublicKey]solana.PublicKey{
											poolConfigPDA: tokenAddress,
										},
									},
								},
							},
						),
					})
					require.NoError(t, err)
				}
				e.Logger.Debugf("MCMS Configured for token pool %v with token address %v", testCase.poolType, tokenAddress)
				e, _, err = commonchangeset.ApplyChangesetsV2(t, e, []commonchangeset.ConfiguredChangeSet{
					// upgrade authority
					commonchangeset.Configure(
						cldf.CreateLegacyChangeSet(ccipChangesetSolana.SetUpgradeAuthorityChangeset),
						ccipChangesetSolana.SetUpgradeAuthorityConfig{
							ChainSelector:       solChain,
							NewUpgradeAuthority: timelockSignerPDA,
							TransferKeys: []solana.PublicKey{
								state.SolChains[solChain].BurnMintTokenPools[tokenMetadata],
								state.SolChains[solChain].LockReleaseTokenPools[tokenMetadata],
							},
						},
					),
					commonchangeset.Configure(
						cldf.CreateLegacyChangeSet(ccipChangesetSolana.SetUpgradeAuthorityChangeset),
						ccipChangesetSolana.SetUpgradeAuthorityConfig{
							ChainSelector:       solChain,
							NewUpgradeAuthority: e.SolChains[solChain].DeployerKey.PublicKey(),
							TransferKeys: []solana.PublicKey{
								state.SolChains[solChain].BurnMintTokenPools[tokenMetadata],
								state.SolChains[solChain].LockReleaseTokenPools[tokenMetadata],
							},
							MCMS: &proposalutils.TimelockConfig{
								MinDelay: 1 * time.Second,
							},
						},
					),
				})
				require.NoError(t, err)
			}
		}
	}
}

func TestPartnerTokenPools(t *testing.T) {
	skipInCI(t)
	tenv, _ := testhelpers.NewMemoryEnvironment(t, testhelpers.WithSolChains(1))
	e := tenv.Env
	solChainSelectors := e.AllChainSelectorsSolana()
	metadata := "partner_testing"
	e, _, err := commonchangeset.ApplyChangesetsV2(t, e, []commonchangeset.ConfiguredChangeSet{commonchangeset.Configure(
		cldf.CreateLegacyChangeSet(ccipChangesetSolana.DeployChainContractsChangeset),
		ccipChangesetSolana.DeployChainContractsConfig{
			HomeChainSelector: e.AllChainSelectors()[0],
			ChainSelector:     solChainSelectors[0],
			BuildConfig: &ccipChangesetSolana.BuildSolanaConfig{
				GitCommitSha:   OldSha,
				DestinationDir: e.SolChains[solChainSelectors[0]].ProgramsPath,
				LocalBuild: ccipChangesetSolana.LocalBuildConfig{
					BuildLocally: true,
				},
			},
			LockReleaseTokenPoolMetadata: metadata,
			BurnMintTokenPoolMetadata:    metadata,
		},
	)})
	require.NoError(t, err)
	testhelpers.ValidateSolanaState(t, e, solChainSelectors)
	doTestTokenPool(t, e, false, metadata)
	doTestPoolLookupTable(t, e, false, metadata)
	doTestTokenPool(t, e, true, metadata)
}
