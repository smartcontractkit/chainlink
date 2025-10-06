package solana_test

import (
	"context"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gagliardetto/solana-go"
	solRpc "github.com/gagliardetto/solana-go/rpc"
	chain_selectors "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/quarantine"
	"github.com/stretchr/testify/require"

	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"

	solTestUtil "github.com/smartcontractkit/chainlink-ccip/chains/solana/contracts/tests/testutils"
	solBaseTokenPool "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/v0_1_1/base_token_pool"
	"github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/v0_1_1/cctp_token_pool"
	solTestTokenPool "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/v0_1_1/test_token_pool"
	solCommon "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/common"
	solTokenUtil "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/tokens"

	"github.com/smartcontractkit/chainlink-evm/gethwrappers/shared/generated/initial/burn_mint_erc677"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/testcontext"

	cldf_evm "github.com/smartcontractkit/chainlink-deployments-framework/chain/evm"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_5_1/token_pool"

	ccipChangesetSolana "github.com/smartcontractkit/chainlink/deployment/ccip/changeset/solana_v0_1_1"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/v1_5_1"

	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"

	"github.com/smartcontractkit/chainlink/deployment"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
)

const (
	PartnerMetadata  = "partner_testing"
	PartnerMetadata2 = "partner_testing_2"
	PartnerMetadata3 = "partner_testing_3"
)

func TestAddTokenPoolWithoutMcms(t *testing.T) {
	t.Parallel()
	skipInCI(t)

	tenv, _ := testhelpers.NewMemoryEnvironment(t, testhelpers.WithSolChains(1), testhelpers.WithCCIPSolanaContractVersion(ccipChangesetSolana.SolanaContractV0_1_1))
	doTestTokenPool(t, tenv.Env, TokenPoolTestConfig{MCMS: false, TokenMetadata: shared.CLLMetadata})
}

func TestAddTokenPoolWithMcms(t *testing.T) {
	quarantine.Flaky(t, "DX-1797")
	t.Parallel()
	tenv, _ := testhelpers.NewMemoryEnvironment(t, testhelpers.WithSolChains(1), testhelpers.WithCCIPSolanaContractVersion(ccipChangesetSolana.SolanaContractV0_1_1))
	doTestTokenPool(t, tenv.Env, TokenPoolTestConfig{MCMS: true, TokenMetadata: shared.CLLMetadata})
}

func deployEVMTokenPool(t *testing.T, e cldf.Environment, evmChain uint64) (cldf.Environment, common.Address, error) {
	addressBook := cldf.NewMemoryAddressBook()
	evmToken, err := cldf.DeployContract(e.Logger, e.BlockChains.EVMChains()[evmChain], addressBook,
		func(chain cldf_evm.Chain) cldf.ContractDeploy[*burn_mint_erc677.BurnMintERC677] {
			tokenAddress, tx, token, err := burn_mint_erc677.DeployBurnMintERC677(
				e.BlockChains.EVMChains()[evmChain].DeployerKey,
				e.BlockChains.EVMChains()[evmChain].Client,
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
	}, true)
	return e, evmToken.Address, nil
}

type TokenPoolTestConfig struct {
	MCMS          bool
	TokenMetadata string
}

func doTestTokenPool(t *testing.T, e cldf.Environment, config TokenPoolTestConfig) {
	mcms := config.MCMS
	tokenMetadata := config.TokenMetadata

	ctx := testcontext.Get(t)
	evmChain := e.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chain_selectors.FamilyEVM))[0]
	solChain := e.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chain_selectors.FamilySolana))[0]
	deployerKey := e.BlockChains.SolanaChains()[solChain].DeployerKey.PublicKey()
	testUser, _ := solana.NewRandomPrivateKey()
	testUserPubKey := testUser.PublicKey()
	e, newTokenAddress, err := deployTokenAndMint(t, e, solChain, []string{deployerKey.String(), testUserPubKey.String()}, "TEST_TOKEN")
	require.NoError(t, err)
	e, newTokenAddress2, err := deployTokenAndMint(t, e, solChain, []string{deployerKey.String(), testUserPubKey.String()}, "TEST_TOKEN_2")
	require.NoError(t, err)
	state, err := stateview.LoadOnchainStateSolana(e)
	require.NoError(t, err)
	testUserATA, _, err := solTokenUtil.FindAssociatedTokenAddress(solana.TokenProgramID, newTokenAddress, testUserPubKey)
	require.NoError(t, err)
	deployerATA, _, err := solTokenUtil.FindAssociatedTokenAddress(
		solana.TokenProgramID,
		newTokenAddress,
		deployerKey,
	)
	rebalancer := deployerKey
	var mcmsConfig *proposalutils.TimelockConfig
	if mcms {
		timelockSignerPDA, _ := testhelpers.TransferOwnershipSolanaV0_1_1(t, &e, solChain, true,
			ccipChangesetSolana.CCIPContractsToTransfer{
				Router:    true,
				FeeQuoter: true,
				OffRamp:   true,
			})
		mcmsConfig = &proposalutils.TimelockConfig{
			MinDelay: 1 * time.Second,
		}
		rebalancer = timelockSignerPDA
	}
	require.NoError(t, err)

	rateLimitConfig := solBaseTokenPool.RateLimitConfig{
		Enabled:  true,
		Capacity: uint64(50e11),
		Rate:     uint64(167000000000),
	}
	inboundConfig := rateLimitConfig
	outboundConfig := rateLimitConfig

	type poolTestType struct {
		poolType    cldf.ContractType
		poolAddress solana.PublicKey
	}
	testCases := []poolTestType{
		{
			poolType:    shared.BurnMintTokenPool,
			poolAddress: state.SolChains[solChain].BurnMintTokenPools[tokenMetadata],
		},
		{
			poolType:    shared.LockReleaseTokenPool,
			poolAddress: state.SolChains[solChain].LockReleaseTokenPools[tokenMetadata],
		},
		{
			poolType:    shared.CCTPTokenPool,
			poolAddress: state.SolChains[solChain].CCTPTokenPool,
		},
	}

	// evm deployment
	e, _, err = deployEVMTokenPool(t, e, evmChain)
	require.NoError(t, err)

	tokenAddress := newTokenAddress
	timelockSignerPDA, err := ccipChangesetSolana.FetchTimelockSigner(e, solChain)
	require.NoError(t, err)

	// svm deployment
	for _, testCase := range testCases {
		var cctpMessengerMinter, cctpMessageTransmitter solana.PublicKey
		// only set for CCTP token pool migrations
		if testCase.poolType == shared.CCTPTokenPool {
			cctpMessengerMinter = getRandomPubKey(t)    // using mock minter
			cctpMessageTransmitter = getRandomPubKey(t) // using mock transmitter
			if config.TokenMetadata == PartnerMetadata ||
				config.TokenMetadata == PartnerMetadata2 ||
				config.TokenMetadata == PartnerMetadata3 {
				continue // CCTP token pool not supported for partner metadata
			}
		}
		e, _, err = commonchangeset.ApplyChangesets(t, e, []commonchangeset.ConfiguredChangeSet{
			commonchangeset.Configure(
				cldf.CreateLegacyChangeSet(ccipChangesetSolana.InitGlobalConfigTokenPoolProgram),
				ccipChangesetSolana.TokenPoolConfigWithMCM{
					ChainSelector: solChain,
					MCMS:          mcmsConfig,
					TokenPoolConfigs: []ccipChangesetSolana.TokenPoolConfig{
						{
							TokenPubKey: tokenAddress,
							PoolType:    testCase.poolType,
							Metadata:    tokenMetadata,
						},
					},
				},
			),
			commonchangeset.Configure(
				cldf.CreateLegacyChangeSet(ccipChangesetSolana.AddTokenPoolAndLookupTable),
				ccipChangesetSolana.AddTokenPoolAndLookupTableConfig{
					ChainSelector: solChain,
					TokenPoolConfigs: []ccipChangesetSolana.TokenPoolConfig{
						{
							TokenPubKey:              tokenAddress,
							PoolType:                 testCase.poolType,
							Metadata:                 tokenMetadata,
							CCTPTokenMessengerMinter: cctpMessengerMinter,
							CCTPMessageTransmitter:   cctpMessageTransmitter,
						},
						{
							TokenPubKey:              newTokenAddress2,
							PoolType:                 testCase.poolType,
							Metadata:                 tokenMetadata,
							CCTPTokenMessengerMinter: cctpMessengerMinter,
							CCTPMessageTransmitter:   cctpMessageTransmitter,
						},
					},
				},
			),
			commonchangeset.Configure(
				cldf.CreateLegacyChangeSet(ccipChangesetSolana.SetupTokenPoolForRemoteChain),
				ccipChangesetSolana.SetupTokenPoolForRemoteChainConfig{
					SolChainSelector: solChain,
					RemoteTokenPoolConfigs: []ccipChangesetSolana.RemoteChainTokenPoolConfig{
						{
							SolTokenPubKey: tokenAddress,
							SolPoolType:    testCase.poolType,
							Metadata:       tokenMetadata,
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
						},
					},
					MCMS: mcmsConfig,
				},
			),
			commonchangeset.Configure(
				cldf.CreateLegacyChangeSet(ccipChangesetSolana.SetRateLimitAdmin),
				ccipChangesetSolana.SetRateLimitAdminConfig{
					SolChainSelector: solChain,
					RateLimitAdminConfigs: []ccipChangesetSolana.RateLimitAdminConfig{
						{
							SolTokenPubKey:    tokenAddress.String(),
							PoolType:          testCase.poolType,
							Metadata:          tokenMetadata,
							NewRateLimitAdmin: timelockSignerPDA,
						},
						{
							SolTokenPubKey:    newTokenAddress2.String(),
							PoolType:          testCase.poolType,
							Metadata:          tokenMetadata,
							NewRateLimitAdmin: timelockSignerPDA,
						},
					},
					MCMS: mcmsConfig,
				},
			),
			commonchangeset.Configure(
				cldf.CreateLegacyChangeSet(ccipChangesetSolana.SetRateLimitAdmin),
				ccipChangesetSolana.SetRateLimitAdminConfig{
					SolChainSelector: solChain,
					RateLimitAdminConfigs: []ccipChangesetSolana.RateLimitAdminConfig{
						{
							SolTokenPubKey:    tokenAddress.String(),
							PoolType:          testCase.poolType,
							Metadata:          tokenMetadata,
							NewRateLimitAdmin: deployerKey,
						},
						{
							SolTokenPubKey:    newTokenAddress2.String(),
							PoolType:          testCase.poolType,
							Metadata:          tokenMetadata,
							NewRateLimitAdmin: deployerKey,
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
		err = e.BlockChains.SolanaChains()[solChain].GetAccountDataBorshInto(ctx, poolConfigPDA, &configAccount)
		require.NoError(t, err)
		require.Equal(t, tokenAddress, configAccount.Config.Mint)
		// test SetupTokenPoolForRemoteChain results
		remoteChainConfigPDA, _, _ := solTokenUtil.TokenPoolChainConfigPDA(evmChain, tokenAddress, testCase.poolAddress)
		base := getTokenPoolBaseChainConfig(t, e, solChain, testCase.poolType, remoteChainConfigPDA)
		require.Equal(t, testhelpers.LocalTokenDecimals, int(base.Remote.Decimals))
		e.Logger.Infof("Pool addresses: %v", base.Remote.PoolAddresses)
		require.Len(t, base.Remote.PoolAddresses, 1)
		require.Equal(t, inboundConfig.Rate, base.InboundRateLimit.Cfg.Rate)
		require.Equal(t, outboundConfig.Rate, base.OutboundRateLimit.Cfg.Rate)

		allowedAccount1, _ := solana.NewRandomPrivateKey()
		allowedAccount2, _ := solana.NewRandomPrivateKey()

		newRateLimitConfig := solBaseTokenPool.RateLimitConfig{
			Enabled:  true,
			Capacity: uint64(50e12),
			Rate:     uint64(1670000000000),
		}
		newOutboundConfig := newRateLimitConfig
		newInboundConfig := newRateLimitConfig

		if mcms {
			e.Logger.Debugf("Configuring MCMS for token pool %v", testCase.poolType)
			switch testCase.poolType {
			case shared.BurnMintTokenPool:
				_, _ = testhelpers.TransferOwnershipSolanaV0_1_1(
					t, &e, solChain, false,
					ccipChangesetSolana.CCIPContractsToTransfer{
						BurnMintTokenPools: map[string][]solana.PublicKey{
							tokenMetadata: {tokenAddress},
						},
					})
			case shared.LockReleaseTokenPool:
				_, _ = testhelpers.TransferOwnershipSolanaV0_1_1(
					t, &e, solChain, false,
					ccipChangesetSolana.CCIPContractsToTransfer{
						LockReleaseTokenPools: map[string][]solana.PublicKey{
							tokenMetadata: {tokenAddress},
						},
					})
			case shared.CCTPTokenPool:
				_, _ = testhelpers.TransferOwnershipSolanaV0_1_1(
					t, &e, solChain, false,
					ccipChangesetSolana.CCIPContractsToTransfer{
						CCTPTokenPoolMints: []solana.PublicKey{tokenAddress},
					})
			default:
				panic("unhandled default case")
			}
			e.Logger.Debugf("MCMS Configured for token pool %v with token address %v", testCase.poolType, tokenAddress)
		}

		e, _, err = commonchangeset.ApplyChangesets(t, e, []commonchangeset.ConfiguredChangeSet{
			commonchangeset.Configure(
				cldf.CreateLegacyChangeSet(ccipChangesetSolana.ConfigureTokenPoolAllowList),
				ccipChangesetSolana.ConfigureTokenPoolAllowListConfig{
					SolChainSelector: solChain,
					SolTokenPubKey:   tokenAddress.String(),
					PoolType:         testCase.poolType,
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
					PoolType:         testCase.poolType,
					Metadata:         tokenMetadata,
					Accounts:         []solana.PublicKey{allowedAccount1.PublicKey(), allowedAccount2.PublicKey()},
					MCMS:             mcmsConfig,
				},
			),
			// test update
			commonchangeset.Configure(
				cldf.CreateLegacyChangeSet(ccipChangesetSolana.SetupTokenPoolForRemoteChain),
				ccipChangesetSolana.SetupTokenPoolForRemoteChainConfig{
					SolChainSelector: solChain,
					MCMS:             mcmsConfig,
					RemoteTokenPoolConfigs: []ccipChangesetSolana.RemoteChainTokenPoolConfig{
						{
							SolTokenPubKey: tokenAddress,
							SolPoolType:    testCase.poolType,
							Metadata:       tokenMetadata,
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
						},
					},
				},
			),
		})
		require.NoError(t, err)

		base = getTokenPoolBaseChainConfig(t, e, solChain, testCase.poolType, remoteChainConfigPDA)
		require.Equal(t, newInboundConfig.Rate, base.InboundRateLimit.Cfg.Rate)
		require.Equal(t, newOutboundConfig.Rate, base.OutboundRateLimit.Cfg.Rate)

		if testCase.poolType == shared.LockReleaseTokenPool && tokenAddress == newTokenAddress {
			e, _, err = commonchangeset.ApplyChangesets(t, e, []commonchangeset.ConfiguredChangeSet{
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
						RebalancerCfg: &ccipChangesetSolana.RebalancerConfig{
							Rebalancer: rebalancer,
						},
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
			outDec, outVal, err := solTokenUtil.TokenBalance(e.GetContext(), e.BlockChains.SolanaChains()[solChain].Client, deployerATA, solRpc.CommitmentConfirmed)
			require.NoError(t, err)
			require.Equal(t, int(900), outVal)
			require.Equal(t, 9, int(outDec))

			outDec, outVal, err = solTokenUtil.TokenBalance(e.GetContext(), e.BlockChains.SolanaChains()[solChain].Client, testUserATA, solRpc.CommitmentConfirmed)
			require.NoError(t, err)
			require.Equal(t, int(1050), outVal)
			require.Equal(t, 9, int(outDec))

			err = e.BlockChains.SolanaChains()[solChain].GetAccountDataBorshInto(ctx, poolConfigPDA, &configAccount)
			require.NoError(t, err)
			outDec, outVal, err = solTokenUtil.TokenBalance(e.GetContext(), e.BlockChains.SolanaChains()[solChain].Client, configAccount.Config.PoolTokenAccount, solRpc.CommitmentConfirmed)
			require.NoError(t, err)
			require.Equal(t, int(50), outVal)
			require.Equal(t, 9, int(outDec))

			// transfer away from timelock if metadata is set and not ccipChangeset.CLLMetadata
			if mcms && tokenMetadata != "" && tokenMetadata != shared.CLLMetadata {
				require.NoError(t, err)
				e.Logger.Debugf("Transferring away from MCMS for token pool %v", testCase.poolType)
				switch testCase.poolType {
				case shared.BurnMintTokenPool:
					e, _, err = commonchangeset.ApplyChangesets(t, e, []commonchangeset.ConfiguredChangeSet{
						commonchangeset.Configure(
							cldf.CreateLegacyChangeSet(ccipChangesetSolana.TransferCCIPToMCMSWithTimelockSolana),
							ccipChangesetSolana.TransferCCIPToMCMSWithTimelockSolanaConfig{
								MCMSCfg:       proposalutils.TimelockConfig{MinDelay: 1 * time.Second},
								CurrentOwner:  timelockSignerPDA,
								ProposedOwner: deployerKey,
								ContractsByChain: map[uint64]ccipChangesetSolana.CCIPContractsToTransfer{
									solChain: {
										BurnMintTokenPools: map[string][]solana.PublicKey{
											tokenMetadata: {tokenAddress},
										},
									},
								},
							},
						),
					})
					require.NoError(t, err)
				case shared.LockReleaseTokenPool:
					e, _, err = commonchangeset.ApplyChangesets(t, e, []commonchangeset.ConfiguredChangeSet{
						commonchangeset.Configure(
							cldf.CreateLegacyChangeSet(ccipChangesetSolana.TransferCCIPToMCMSWithTimelockSolana),
							ccipChangesetSolana.TransferCCIPToMCMSWithTimelockSolanaConfig{
								MCMSCfg:       proposalutils.TimelockConfig{MinDelay: 1 * time.Second},
								CurrentOwner:  timelockSignerPDA,
								ProposedOwner: deployerKey,
								ContractsByChain: map[uint64]ccipChangesetSolana.CCIPContractsToTransfer{
									solChain: {
										LockReleaseTokenPools: map[string][]solana.PublicKey{
											tokenMetadata: {tokenAddress},
										},
									},
								},
							},
						),
					})
					require.NoError(t, err)
				case shared.CCTPTokenPool:
					e, _, err = commonchangeset.ApplyChangesets(t, e, []commonchangeset.ConfiguredChangeSet{
						commonchangeset.Configure(
							cldf.CreateLegacyChangeSet(ccipChangesetSolana.TransferCCIPToMCMSWithTimelockSolana),
							ccipChangesetSolana.TransferCCIPToMCMSWithTimelockSolanaConfig{
								MCMSCfg:       proposalutils.TimelockConfig{MinDelay: 1 * time.Second},
								CurrentOwner:  timelockSignerPDA,
								ProposedOwner: deployerKey,
								ContractsByChain: map[uint64]ccipChangesetSolana.CCIPContractsToTransfer{
									solChain: {
										CCTPTokenPoolMints: []solana.PublicKey{tokenAddress},
									},
								},
							},
						),
					})
					require.NoError(t, err)
				default:
					panic("unhandled default case")
				}
				e.Logger.Debugf("MCMS Configured for token pool %v with token address %v", testCase.poolType, tokenAddress)
				transferKeys := []solana.PublicKey{}
				switch testCase.poolType {
				case shared.BurnMintTokenPool:
					transferKeys = append(transferKeys, state.SolChains[solChain].BurnMintTokenPools[tokenMetadata])
				case shared.LockReleaseTokenPool:
					transferKeys = append(transferKeys, state.SolChains[solChain].LockReleaseTokenPools[tokenMetadata])
				case shared.CCTPTokenPool:
					transferKeys = append(transferKeys, state.SolChains[solChain].CCTPTokenPool)
				}
				e, _, err = commonchangeset.ApplyChangesets(t, e, []commonchangeset.ConfiguredChangeSet{
					// upgrade authority
					commonchangeset.Configure(
						cldf.CreateLegacyChangeSet(ccipChangesetSolana.SetUpgradeAuthorityChangeset),
						ccipChangesetSolana.SetUpgradeAuthorityConfig{
							ChainSelector:       solChain,
							NewUpgradeAuthority: timelockSignerPDA,
							TransferKeys:        transferKeys,
						},
					),
					commonchangeset.Configure(
						cldf.CreateLegacyChangeSet(ccipChangesetSolana.SetUpgradeAuthorityChangeset),
						ccipChangesetSolana.SetUpgradeAuthorityConfig{
							ChainSelector:       solChain,
							NewUpgradeAuthority: e.BlockChains.SolanaChains()[solChain].DeployerKey.PublicKey(),
							TransferKeys:        transferKeys,
							MCMS: &proposalutils.TimelockConfig{
								MinDelay: 1 * time.Second,
							},
						},
					),
				})
				require.NoError(t, err)
			}
		}

		// NOTE: the ModifyMintAuthority changeset only supports BnM token pools at the moment, so
		// we'll only create the multisig account and run the changeset if the pool type is BnM.
		if !mcms && testCase.poolType == shared.BurnMintTokenPool {
			tokenPoolSignerPDA, err := solTokenUtil.TokenPoolSignerAddress(tokenAddress, testCase.poolAddress)
			require.NoError(t, err)

			deployerPrivKey := *e.BlockChains.SolanaChains()[solChain].DeployerKey
			solanaRPCClient := e.BlockChains.SolanaChains()[solChain].Client

			multisig1 := createMultiSig(ctx, t, deployerKey, tokenPoolSignerPDA, solanaRPCClient, deployerPrivKey)
			multisig2 := createMultiSig(ctx, t, deployerKey, tokenPoolSignerPDA, solanaRPCClient, deployerPrivKey)

			e, _, err = commonchangeset.ApplyChangesets(t, e, []commonchangeset.ConfiguredChangeSet{commonchangeset.Configure(
				cldf.CreateLegacyChangeSet(ccipChangesetSolana.ExtendTokenPoolLookupTable),
				ccipChangesetSolana.ExtendTokenPoolLookupTableConfig{
					SkipValidationsForDuplicates: false,
					ChainSelector:                solChain,
					TokenPubKey:                  tokenAddress,
					PoolType:                     testCase.poolType,
					Metadata:                     tokenMetadata,
					Accounts: []solana.PublicKey{
						multisig1.PublicKey(),
						multisig2.PublicKey(),
					},
				},
			)})
			require.NoError(t, err)

			e, _, err = commonchangeset.ApplyChangesets(t, e, []commonchangeset.ConfiguredChangeSet{commonchangeset.Configure(
				cldf.CreateLegacyChangeSet(ccipChangesetSolana.ModifyMintAuthority),
				ccipChangesetSolana.NewMintTokenPoolConfig{
					NewMintAuthority: multisig1.PublicKey(),
					ChainSelector:    solChain,
					TokenPubKey:      tokenAddress,
					PoolType:         testCase.poolType,
					MCMS:             mcmsConfig,
					Metadata:         tokenMetadata,
				},
			)})
			require.NoError(t, err)

			// When changing the mint authority when the current is a multisig, we need to add the previous one as a remaining account. https://github.com/smartcontractkit/chainlink-ccip/blob/c01850356f7f653f46e6f36ac04443122c49d04b/chains/solana/contracts/programs/burnmint-token-pool/src/lib.rs#L165C13-L176C18
			e, _, err = commonchangeset.ApplyChangesets(t, e, []commonchangeset.ConfiguredChangeSet{commonchangeset.Configure(
				cldf.CreateLegacyChangeSet(ccipChangesetSolana.ModifyMintAuthority),
				ccipChangesetSolana.NewMintTokenPoolConfig{
					NewMintAuthority: multisig2.PublicKey(),
					OldMintAuthority: multisig1.PublicKey(),
					ChainSelector:    solChain,
					TokenPubKey:      tokenAddress,
					PoolType:         testCase.poolType,
					MCMS:             mcmsConfig,
					Metadata:         tokenMetadata,
				},
			)})
			require.NoError(t, err)
		}
	}
}

func createMultiSig(ctx context.Context, t *testing.T, deployerKey solana.PublicKey, tokenPoolSignerPDA solana.PublicKey, solanaRPCClient *solRpc.Client, deployerPrivateKey solana.PrivateKey) solana.PrivateKey {
	multisig, err := solana.NewRandomPrivateKey()
	require.NoError(t, err)

	ixMsig, ixErrMsig := solTokenUtil.CreateMultisig(ctx,
		deployerKey,
		solana.TokenProgramID,
		multisig.PublicKey(),
		1,
		[]solana.PublicKey{deployerKey, tokenPoolSignerPDA},
		solanaRPCClient,
		solRpc.CommitmentConfirmed,
	)
	require.NoError(t, ixErrMsig)

	res := solTestUtil.SendAndConfirm(ctx, t,
		solanaRPCClient,
		ixMsig,
		deployerPrivateKey,
		solRpc.CommitmentConfirmed,
		solCommon.AddSigners(multisig, deployerPrivateKey),
	)
	require.NotNil(t, res)
	return multisig
}

var zeroRateLimitConfig = ccipChangesetSolana.RateLimiterConfig{
	Inbound: solBaseTokenPool.RateLimitConfig{
		Enabled:  false,
		Capacity: 0,
		Rate:     0,
	},
	Outbound: solBaseTokenPool.RateLimitConfig{
		Enabled:  false,
		Capacity: 0,
		Rate:     0,
	},
}

func TestAddTokenPoolE2EWithMcms(t *testing.T) {
	quarantine.Flaky(t, "DX-1774")
	t.Parallel()
	tenv, _ := testhelpers.NewMemoryEnvironment(t, testhelpers.WithSolChains(1), testhelpers.WithCCIPSolanaContractVersion(ccipChangesetSolana.SolanaContractV0_1_1))
	solChain := tenv.Env.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chain_selectors.FamilySolana))[0]
	evmChain := tenv.Env.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chain_selectors.FamilyEVM))[0]
	deployerKey := tenv.Env.BlockChains.SolanaChains()[solChain].DeployerKey.PublicKey()
	poolType := shared.BurnMintTokenPool
	e, newTokenAddress, err := deployTokenAndMint(t, tenv.Env, solChain, []string{deployerKey.String()}, "TEST_TOKEN")
	require.NoError(t, err)
	e, newTokenAddress2, err := deployTokenAndMint(t, e, solChain, []string{deployerKey.String()}, "TEST_TOKEN_2")
	require.NoError(t, err)
	// evm deployment
	e, _, err = deployEVMTokenPool(t, e, evmChain)
	require.NoError(t, err)
	_, _ = testhelpers.TransferOwnershipSolanaV0_1_1(t, &e, solChain, true,
		ccipChangesetSolana.CCIPContractsToTransfer{
			Router:    true,
			FeeQuoter: true,
			OffRamp:   true,
		})
	mcmsConfig := &proposalutils.TimelockConfig{
		MinDelay: 1 * time.Second,
	}
	_, _, err = commonchangeset.ApplyChangesets(t, e, []commonchangeset.ConfiguredChangeSet{
		commonchangeset.Configure(
			cldf.CreateLegacyChangeSet(ccipChangesetSolana.E2ETokenPoolv2),
			ccipChangesetSolana.E2ETokenPoolConfigv2{
				ChainSelector: solChain,
				MCMS:          mcmsConfig,
				E2ETokens: []ccipChangesetSolana.E2ETokenConfig{
					{
						TokenPubKey: newTokenAddress,
						PoolType:    poolType,
						Metadata:    shared.CLLMetadata,
						SolanaToEVMRemoteConfigs: map[uint64]ccipChangesetSolana.EVMRemoteConfig{
							evmChain: {
								TokenSymbol:       testhelpers.TestTokenSymbol,
								PoolType:          shared.BurnMintTokenPool,
								PoolVersion:       shared.CurrentTokenPoolVersion,
								RateLimiterConfig: zeroRateLimitConfig,
							},
						},
						EVMToSolanaRemoteConfigs: v1_5_1.ConfigureTokenPoolContractsConfig{
							TokenSymbol: testhelpers.TestTokenSymbol,
							MCMS:        mcmsConfig,
							PoolUpdates: map[uint64]v1_5_1.TokenPoolConfig{
								evmChain: {
									Type:    shared.BurnMintTokenPool,
									Version: shared.CurrentTokenPoolVersion,
									SolChainUpdates: map[uint64]v1_5_1.SolChainUpdate{
										solChain: {
											RateLimiterConfig: v1_5_1.RateLimiterConfig{
												Inbound: token_pool.RateLimiterConfig{
													IsEnabled: false,
													Capacity:  big.NewInt(0),
													Rate:      big.NewInt(0),
												},
												Outbound: token_pool.RateLimiterConfig{
													IsEnabled: false,
													Capacity:  big.NewInt(0),
													Rate:      big.NewInt(0),
												},
											},
											TokenAddress: newTokenAddress.String(),
											Type:         shared.BurnMintTokenPool,
											Metadata:     shared.CLLMetadata,
										},
									},
								},
							},
						},
					},
					{
						TokenPubKey: newTokenAddress2,
						PoolType:    poolType,
						Metadata:    shared.CLLMetadata,
						SolanaToEVMRemoteConfigs: map[uint64]ccipChangesetSolana.EVMRemoteConfig{
							evmChain: {
								TokenSymbol:       testhelpers.TestTokenSymbol,
								PoolType:          shared.BurnMintTokenPool,
								PoolVersion:       shared.CurrentTokenPoolVersion,
								RateLimiterConfig: zeroRateLimitConfig,
							},
						},
					},
				},
			},
		),
	})
	require.NoError(t, err)
}

func TestPartnerTokenPools(t *testing.T) {
	skipInCI(t)
	tenv, _ := testhelpers.NewMemoryEnvironment(t, testhelpers.WithSolChains(1), testhelpers.WithCCIPSolanaContractVersion(ccipChangesetSolana.SolanaContractV0_1_1))
	e := tenv.Env
	solChainSelectors := e.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chain_selectors.FamilySolana))

	e, _, err := commonchangeset.ApplyChangesets(t, e, []commonchangeset.ConfiguredChangeSet{commonchangeset.Configure(
		cldf.CreateLegacyChangeSet(ccipChangesetSolana.DeployChainContractsChangeset),
		ccipChangesetSolana.DeployChainContractsConfig{
			HomeChainSelector: e.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chain_selectors.FamilyEVM))[0],
			ChainSelector:     solChainSelectors[0],
			BuildConfig: &ccipChangesetSolana.BuildSolanaConfig{
				SolanaContractVersion: ccipChangesetSolana.VersionSolanaV0_1_1,
				DestinationDir:        e.BlockChains.SolanaChains()[solChainSelectors[0]].ProgramsPath,
				LocalBuild: ccipChangesetSolana.LocalBuildConfig{
					BuildLocally: true,
				},
			},
			LockReleaseTokenPoolMetadata: []string{PartnerMetadata},
			BurnMintTokenPoolMetadata:    []string{PartnerMetadata},
		},
	)})
	require.NoError(t, err)
	err = testhelpers.ValidateSolanaState(e, solChainSelectors)
	require.NoError(t, err)
	doTestTokenPool(t, e, TokenPoolTestConfig{MCMS: false, TokenMetadata: PartnerMetadata})
	doTestPoolLookupTable(t, e, false, PartnerMetadata)
	doTestTokenPool(t, e, TokenPoolTestConfig{MCMS: true, TokenMetadata: PartnerMetadata})
}

func TestPartnerTokenPoolsWithUpgrades(t *testing.T) {
	skipInCI(t)
	tenv, _ := testhelpers.NewMemoryEnvironment(t, testhelpers.WithSolChains(1), testhelpers.WithCCIPSolanaContractVersion(ccipChangesetSolana.SolanaContractV0_1_1))
	e := tenv.Env
	var err error
	solChainSelectors := e.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chain_selectors.FamilySolana))

	allMetadata := []string{PartnerMetadata, PartnerMetadata2, PartnerMetadata3}
	for _, metadata := range allMetadata {
		e, _, err = commonchangeset.ApplyChangesets(t, e, []commonchangeset.ConfiguredChangeSet{commonchangeset.Configure(
			cldf.CreateLegacyChangeSet(ccipChangesetSolana.DeployChainContractsChangeset),
			ccipChangesetSolana.DeployChainContractsConfig{
				HomeChainSelector: e.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chain_selectors.FamilyEVM))[0],
				ChainSelector:     solChainSelectors[0],
				BuildConfig: &ccipChangesetSolana.BuildSolanaConfig{
					SolanaContractVersion: ccipChangesetSolana.VersionSolanaV0_1_1,
					DestinationDir:        e.BlockChains.SolanaChains()[solChainSelectors[0]].ProgramsPath,
					LocalBuild: ccipChangesetSolana.LocalBuildConfig{
						BuildLocally: true,
						CleanGitDir:  true,
					},
				},
				LockReleaseTokenPoolMetadata: []string{metadata},
				BurnMintTokenPoolMetadata:    []string{metadata},
			},
		)})
		require.NoError(t, err)
		err = testhelpers.ValidateSolanaState(e, solChainSelectors)
		require.NoError(t, err)
		doTestTokenPool(t, e, TokenPoolTestConfig{MCMS: false, TokenMetadata: metadata})
		doTestPoolLookupTable(t, e, false, metadata)
	}
	state, err := stateview.LoadOnchainState(e)
	require.NoError(t, err)
	// CLL + partner + partner2 + partner3
	require.Len(t, state.SolChains[solChainSelectors[0]].BurnMintTokenPools, len(allMetadata)+1)
	require.Len(t, state.SolChains[solChainSelectors[0]].LockReleaseTokenPools, len(allMetadata)+1)
	e, _, err = commonchangeset.ApplyChangesets(t, e, []commonchangeset.ConfiguredChangeSet{commonchangeset.Configure(
		cldf.CreateLegacyChangeSet(ccipChangesetSolana.DeployChainContractsChangeset),
		ccipChangesetSolana.DeployChainContractsConfig{
			HomeChainSelector: e.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chain_selectors.FamilyEVM))[0],
			ChainSelector:     solChainSelectors[0],
			BuildConfig: &ccipChangesetSolana.BuildSolanaConfig{
				SolanaContractVersion: ccipChangesetSolana.VersionSolanaV0_1_2,
				DestinationDir:        e.BlockChains.SolanaChains()[solChainSelectors[0]].ProgramsPath,
			},
			UpgradeConfig: ccipChangesetSolana.UpgradeConfig{
				NewBurnMintTokenPoolVersion:    &deployment.Version1_0_0,
				NewLockReleaseTokenPoolVersion: &deployment.Version1_0_0,
				UpgradeAuthority:               e.BlockChains.SolanaChains()[solChainSelectors[0]].DeployerKey.PublicKey(),
			},
			LockReleaseTokenPoolMetadata: allMetadata,
			BurnMintTokenPoolMetadata:    allMetadata,
		},
	)})
	require.NoError(t, err)
	err = testhelpers.ValidateSolanaState(e, solChainSelectors)
	require.NoError(t, err)

	for _, metadata := range allMetadata {
		doTestTokenPool(t, e, TokenPoolTestConfig{MCMS: false, TokenMetadata: metadata})
		doTestPoolLookupTable(t, e, false, metadata)
	}
}

func getRandomPubKey(t *testing.T) solana.PublicKey {
	t.Helper()
	privKey, err := solana.NewRandomPrivateKey()
	require.NoError(t, err)
	return privKey.PublicKey()
}

func getTokenPoolBaseChainConfig(t *testing.T, e cldf.Environment, solChain uint64, poolType cldf.ContractType, remoteChainConfigPDA solana.PublicKey) solBaseTokenPool.BaseChain {
	t.Helper()

	var base solBaseTokenPool.BaseChain
	switch poolType {
	case shared.BurnMintTokenPool, shared.LockReleaseTokenPool:
		var remoteChainConfigAccount solTestTokenPool.ChainConfig
		err := e.BlockChains.SolanaChains()[solChain].GetAccountDataBorshInto(t.Context(), remoteChainConfigPDA, &remoteChainConfigAccount)
		require.NoError(t, err)
		base = remoteChainConfigAccount.Base
	case shared.CCTPTokenPool:
		var remoteChainConfigAccount cctp_token_pool.ChainConfig
		err := e.BlockChains.SolanaChains()[solChain].GetAccountDataBorshInto(t.Context(), remoteChainConfigPDA, &remoteChainConfigAccount)
		require.NoError(t, err)
		base = remoteChainConfigAccount.Base
	default:
		require.Fail(t, "unsupported token pool type")
	}
	return base
}
