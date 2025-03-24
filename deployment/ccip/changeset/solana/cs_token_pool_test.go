package solana_test

import (
	"testing"
	"time"

	"github.com/gagliardetto/solana-go"
	solRpc "github.com/gagliardetto/solana-go/rpc"
	"github.com/stretchr/testify/require"

	solBaseTokenPool "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/base_token_pool"
	solTestTokenPool "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/test_token_pool"
	solTokenUtil "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/tokens"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/testcontext"

	ccipChangeset "github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	ccipChangesetSolana "github.com/smartcontractkit/chainlink/deployment/ccip/changeset/solana"
	changeset_solana "github.com/smartcontractkit/chainlink/deployment/ccip/changeset/solana"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers"

	"github.com/smartcontractkit/chainlink/deployment"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
)

func TestAddTokenPool(t *testing.T) {
	t.Parallel()
	doTestTokenPool(t, false)
}

func TestAddTokenPoolMcms(t *testing.T) {
	t.Parallel()
	doTestTokenPool(t, true)
}

func doTestTokenPool(t *testing.T, mcms bool) {
	ctx := testcontext.Get(t)
	tenv, _ := testhelpers.NewMemoryEnvironment(t, testhelpers.WithSolChains(1))

	evmChain := tenv.Env.AllChainSelectors()[0]
	solChain := tenv.Env.AllChainSelectorsSolana()[0]
	e, newTokenAddress, err := deployToken(t, tenv.Env, solChain)
	require.NoError(t, err)
	state, err := ccipChangeset.LoadOnchainStateSolana(e)
	require.NoError(t, err)
	// MintTo does not support native tokens
	deployerKey := e.SolChains[solChain].DeployerKey.PublicKey()
	testUser, _ := solana.NewRandomPrivateKey()
	testUserPubKey := testUser.PublicKey()
	e, err = commonchangeset.ApplyChangesetsV2(t, e, []commonchangeset.ConfiguredChangeSet{
		commonchangeset.Configure(
			// deployer creates ATA for itself and testUser
			deployment.CreateLegacyChangeSet(ccipChangesetSolana.CreateSolanaTokenATA),
			ccipChangesetSolana.CreateSolanaTokenATAConfig{
				ChainSelector: solChain,
				TokenPubkey:   newTokenAddress,
				TokenProgram:  ccipChangeset.SPL2022Tokens,
				ATAList:       []string{deployerKey.String(), testUserPubKey.String()},
			},
		),
		commonchangeset.Configure(
			// deployer mints token to itself and testUser
			deployment.CreateLegacyChangeSet(changeset_solana.MintSolanaToken),
			ccipChangesetSolana.MintSolanaTokenConfig{
				ChainSelector: solChain,
				TokenPubkey:   newTokenAddress.String(),
				AmountToAddress: map[string]uint64{
					deployerKey.String():    uint64(1000),
					testUserPubKey.String(): uint64(1000),
				},
			},
		),
	},
	)
	require.NoError(t, err)
	testUserATA, _, err := solTokenUtil.FindAssociatedTokenAddress(solana.Token2022ProgramID, newTokenAddress, testUserPubKey)
	require.NoError(t, err)
	deployerATA, _, err := solTokenUtil.FindAssociatedTokenAddress(
		solana.Token2022ProgramID,
		newTokenAddress,
		e.SolChains[solChain].DeployerKey.PublicKey(),
	)
	var mcmsConfig *ccipChangesetSolana.MCMSConfigSolana
	if mcms {
		_, _ = testhelpers.TransferOwnershipSolana(t, &e, solChain, true,
			ccipChangesetSolana.CCIPContractsToTransfer{
				Router:    true,
				FeeQuoter: true,
				OffRamp:   true,
			})
		mcmsConfig = &ccipChangesetSolana.MCMSConfigSolana{
			MCMS: &ccipChangeset.MCMSConfig{
				MinDelay: 1 * time.Second,
			},
			RouterOwnedByTimelock:    true,
			FeeQuoterOwnedByTimelock: true,
			OffRampOwnedByTimelock:   true,
		}
	}
	require.NoError(t, err)
	remoteConfig := solBaseTokenPool.RemoteConfig{
		PoolAddresses: []solTestTokenPool.RemoteAddress{{Address: []byte{1, 2, 3}}},
		TokenAddress:  solTestTokenPool.RemoteAddress{Address: []byte{4, 5, 6}},
		Decimals:      9,
	}
	inboundConfig := solBaseTokenPool.RateLimitConfig{
		Enabled:  true,
		Capacity: uint64(1000),
		Rate:     1,
	}
	outboundConfig := solBaseTokenPool.RateLimitConfig{
		Enabled:  false,
		Capacity: 0,
		Rate:     0,
	}

	tokenMap := map[deployment.ContractType]solana.PublicKey{
		ccipChangeset.SPL2022Tokens: newTokenAddress,
		ccipChangeset.SPLTokens:     state.SolChains[solChain].WSOL,
	}

	type poolTestType struct {
		poolType    solTestTokenPool.PoolType
		poolAddress solana.PublicKey
		mcms        bool
	}
	testCases := []poolTestType{
		{
			poolType:    solTestTokenPool.BurnAndMint_PoolType,
			poolAddress: state.SolChains[solChain].BurnMintTokenPool,
		},
		{
			poolType:    solTestTokenPool.LockAndRelease_PoolType,
			poolAddress: state.SolChains[solChain].LockReleaseTokenPool,
		},
	}
	burnAndMintOwnedByTimelock := make(map[solana.PublicKey]bool)
	lockAndReleaseOwnedByTimelock := make(map[solana.PublicKey]bool)
	for _, testCase := range testCases {
		for _, tokenAddress := range tokenMap {
			e, err = commonchangeset.ApplyChangesetsV2(t, e, []commonchangeset.ConfiguredChangeSet{
				commonchangeset.Configure(
					deployment.CreateLegacyChangeSet(ccipChangesetSolana.AddTokenPool),
					ccipChangesetSolana.TokenPoolConfig{
						ChainSelector: solChain,
						TokenPubKey:   tokenAddress.String(),
						PoolType:      testCase.poolType,
					},
				),
				commonchangeset.Configure(
					deployment.CreateLegacyChangeSet(ccipChangesetSolana.SetupTokenPoolForRemoteChain),
					ccipChangesetSolana.RemoteChainTokenPoolConfig{
						SolChainSelector:    solChain,
						RemoteChainSelector: evmChain,
						SolTokenPubKey:      tokenAddress.String(),
						RemoteConfig:        remoteConfig,
						InboundRateLimit:    inboundConfig,
						OutboundRateLimit:   outboundConfig,
						PoolType:            testCase.poolType,
						MCMSSolana:          mcmsConfig,
					},
				),
			},
			)
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
			require.Equal(t, uint8(9), remoteChainConfigAccount.Base.Remote.Decimals)

			allowedAccount1, _ := solana.NewRandomPrivateKey()
			allowedAccount2, _ := solana.NewRandomPrivateKey()

			newRemoteConfig := solBaseTokenPool.RemoteConfig{
				PoolAddresses: []solTestTokenPool.RemoteAddress{{Address: []byte{7, 8, 9}}},
				TokenAddress:  solTestTokenPool.RemoteAddress{Address: []byte{10, 11, 12}},
				Decimals:      9,
			}
			newOutboundConfig := solBaseTokenPool.RateLimitConfig{
				Enabled:  true,
				Capacity: uint64(1000),
				Rate:     1,
			}
			newInboundConfig := solBaseTokenPool.RateLimitConfig{
				Enabled:  false,
				Capacity: 0,
				Rate:     0,
			}

			if mcms {
				e.Logger.Debugf("Configuring MCMS for token pool %v", testCase.poolType)
				if testCase.poolType == solTestTokenPool.BurnAndMint_PoolType {
					_, _ = testhelpers.TransferOwnershipSolana(
						t, &e, solChain, false,
						ccipChangesetSolana.CCIPContractsToTransfer{
							BurnMintTokenPools: map[solana.PublicKey]solana.PublicKey{
								poolConfigPDA: tokenAddress,
							},
						})
					burnAndMintOwnedByTimelock[tokenAddress] = true
				} else {
					_, _ = testhelpers.TransferOwnershipSolana(
						t, &e, solChain, false,
						ccipChangesetSolana.CCIPContractsToTransfer{
							LockReleaseTokenPools: map[solana.PublicKey]solana.PublicKey{
								poolConfigPDA: tokenAddress,
							},
						})
					lockAndReleaseOwnedByTimelock[tokenAddress] = true
				}
				mcmsConfig.BurnMintTokenPoolOwnedByTimelock = burnAndMintOwnedByTimelock
				mcmsConfig.LockReleaseTokenPoolOwnedByTimelock = lockAndReleaseOwnedByTimelock
				e.Logger.Debugf("MCMS Configured for token pool %v with token address %v", testCase.poolType, tokenAddress)
			}

			e, err = commonchangeset.ApplyChangesetsV2(t, e, []commonchangeset.ConfiguredChangeSet{
				commonchangeset.Configure(
					deployment.CreateLegacyChangeSet(ccipChangesetSolana.ConfigureTokenPoolAllowList),
					ccipChangesetSolana.ConfigureTokenPoolAllowListConfig{
						SolChainSelector: solChain,
						SolTokenPubKey:   tokenAddress.String(),
						PoolType:         testCase.poolType,
						Accounts:         []solana.PublicKey{allowedAccount1.PublicKey(), allowedAccount2.PublicKey()},
						Enabled:          true,
						MCMSSolana:       mcmsConfig,
					},
				),
				commonchangeset.Configure(
					deployment.CreateLegacyChangeSet(ccipChangesetSolana.RemoveFromTokenPoolAllowList),
					ccipChangesetSolana.RemoveFromAllowListConfig{
						SolChainSelector: solChain,
						SolTokenPubKey:   tokenAddress.String(),
						PoolType:         testCase.poolType,
						Accounts:         []solana.PublicKey{allowedAccount1.PublicKey(), allowedAccount2.PublicKey()},
						MCMSSolana:       mcmsConfig,
					},
				),
				commonchangeset.Configure(
					deployment.CreateLegacyChangeSet(ccipChangesetSolana.ConfigureTokenPoolAllowList),
					ccipChangesetSolana.ConfigureTokenPoolAllowListConfig{
						SolChainSelector: solChain,
						SolTokenPubKey:   tokenAddress.String(),
						PoolType:         testCase.poolType,
						Accounts:         []solana.PublicKey{allowedAccount1.PublicKey(), allowedAccount2.PublicKey()},
						Enabled:          false,
						MCMSSolana:       mcmsConfig,
					},
				),
				commonchangeset.Configure(
					deployment.CreateLegacyChangeSet(ccipChangesetSolana.RemoveFromTokenPoolAllowList),
					ccipChangesetSolana.RemoveFromAllowListConfig{
						SolChainSelector: solChain,
						SolTokenPubKey:   tokenAddress.String(),
						PoolType:         testCase.poolType,
						Accounts:         []solana.PublicKey{allowedAccount1.PublicKey(), allowedAccount2.PublicKey()},
						MCMSSolana:       mcmsConfig,
					},
				),
				commonchangeset.Configure(
					deployment.CreateLegacyChangeSet(ccipChangesetSolana.SetupTokenPoolForRemoteChain),
					ccipChangesetSolana.RemoteChainTokenPoolConfig{
						SolChainSelector:    solChain,
						RemoteChainSelector: evmChain,
						SolTokenPubKey:      tokenAddress.String(),
						RemoteConfig:        newRemoteConfig,
						InboundRateLimit:    newInboundConfig,
						OutboundRateLimit:   newOutboundConfig,
						PoolType:            testCase.poolType,
						MCMSSolana:          mcmsConfig,
						IsUpdate:            true,
					},
				),
			},
			)
			require.NoError(t, err)
			if testCase.poolType == solTestTokenPool.LockAndRelease_PoolType && tokenAddress == newTokenAddress {
				e, err = commonchangeset.ApplyChangesetsV2(t, e, []commonchangeset.ConfiguredChangeSet{
					commonchangeset.Configure(
						deployment.CreateLegacyChangeSet(ccipChangesetSolana.LockReleaseLiquidityOps),
						ccipChangesetSolana.LockReleaseLiquidityOpsConfig{
							SolChainSelector: solChain,
							SolTokenPubKey:   tokenAddress.String(),
							SetCfg: &ccipChangesetSolana.SetLiquidityConfig{
								Enabled: true,
							},
							MCMSSolana: mcmsConfig,
						},
					),
					commonchangeset.Configure(
						deployment.CreateLegacyChangeSet(ccipChangesetSolana.LockReleaseLiquidityOps),
						ccipChangesetSolana.LockReleaseLiquidityOpsConfig{
							SolChainSelector: solChain,
							SolTokenPubKey:   tokenAddress.String(),
							LiquidityCfg: &ccipChangesetSolana.LiquidityConfig{
								Amount:             100,
								RemoteTokenAccount: deployerATA,
								Type:               ccipChangesetSolana.Provide,
							},
							MCMSSolana: mcmsConfig,
						},
					),
					commonchangeset.Configure(
						deployment.CreateLegacyChangeSet(ccipChangesetSolana.LockReleaseLiquidityOps),
						ccipChangesetSolana.LockReleaseLiquidityOpsConfig{
							SolChainSelector: solChain,
							SolTokenPubKey:   tokenAddress.String(),
							LiquidityCfg: &ccipChangesetSolana.LiquidityConfig{
								Amount:             50,
								RemoteTokenAccount: testUserATA,
								Type:               ccipChangesetSolana.Withdraw,
							},
							MCMSSolana: mcmsConfig,
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
			}
		}
	}
}
