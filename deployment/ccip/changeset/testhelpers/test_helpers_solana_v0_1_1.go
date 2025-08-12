package testhelpers

import (
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	chainsel "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_5_1/burn_mint_token_pool"
	solRouter "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/v0_1_1/ccip_router"
	solFeeQuoter "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/v0_1_1/fee_quoter"
	solTestTokenPoolV0_1_1 "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/v0_1_1/test_token_pool"
	soltokens "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/tokens"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/shared/generated/initial/burn_mint_erc677"

	"github.com/smartcontractkit/chainlink/deployment"
	ccipChangeSetSolanaV0_1_1 "github.com/smartcontractkit/chainlink/deployment/ccip/changeset/solana_v0_1_1"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
	solanastateview "github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview/solana"

	"github.com/stretchr/testify/require"

	commoncs "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/common/changeset/state"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	commontypes "github.com/smartcontractkit/chainlink/deployment/common/types"

	"github.com/smartcontractkit/chainlink/deployment/environment/memory"

	"github.com/gagliardetto/solana-go"
)

func TransferOwnershipSolanaV0_1_1(
	t *testing.T,
	e *cldf.Environment,
	solChain uint64,
	needTimelockDeployed bool,
	contractsToTransfer ccipChangeSetSolanaV0_1_1.CCIPContractsToTransfer,
) (timelockSignerPDA solana.PublicKey, mcmSignerPDA solana.PublicKey) {
	var err error
	if needTimelockDeployed {
		*e, _, err = commoncs.ApplyChangesets(t, *e, []commoncs.ConfiguredChangeSet{
			commoncs.Configure(
				cldf.CreateLegacyChangeSet(commoncs.DeployMCMSWithTimelockV2),
				map[uint64]commontypes.MCMSWithTimelockConfigV2{
					solChain: {
						Canceller:        proposalutils.SingleGroupMCMSV2(t),
						Proposer:         proposalutils.SingleGroupMCMSV2(t),
						Bypasser:         proposalutils.SingleGroupMCMSV2(t),
						TimelockMinDelay: big.NewInt(0),
					},
				},
			),
		})
		require.NoError(t, err)
	}

	addresses, err := e.ExistingAddresses.AddressesForChain(solChain)
	require.NoError(t, err)
	mcmState, err := state.MaybeLoadMCMSWithTimelockChainStateSolana(e.BlockChains.SolanaChains()[solChain], addresses)
	require.NoError(t, err)

	// Fund signer PDAs for timelock and mcm
	// If we don't fund, execute() calls will fail with "no funds" errors.
	timelockSignerPDA = state.GetTimelockSignerPDA(mcmState.TimelockProgram, mcmState.TimelockSeed)
	mcmSignerPDA = state.GetMCMSignerPDA(mcmState.McmProgram, mcmState.ProposerMcmSeed)
	err = memory.FundSolanaAccounts(e.GetContext(), []solana.PublicKey{timelockSignerPDA, mcmSignerPDA},
		100, e.BlockChains.SolanaChains()[solChain].Client)
	require.NoError(t, err)
	t.Logf("funded timelock signer PDA: %s", timelockSignerPDA.String())
	t.Logf("funded mcm signer PDA: %s", mcmSignerPDA.String())
	// Apply transfer ownership changeset
	*e, _, err = commoncs.ApplyChangesets(t, *e, []commoncs.ConfiguredChangeSet{
		commoncs.Configure(
			cldf.CreateLegacyChangeSet(ccipChangeSetSolanaV0_1_1.TransferCCIPToMCMSWithTimelockSolana),
			ccipChangeSetSolanaV0_1_1.TransferCCIPToMCMSWithTimelockSolanaConfig{
				MCMSCfg: proposalutils.TimelockConfig{MinDelay: 1 * time.Second},
				ContractsByChain: map[uint64]ccipChangeSetSolanaV0_1_1.CCIPContractsToTransfer{
					solChain: contractsToTransfer,
				},
			},
		),
	})
	require.NoError(t, err)
	return timelockSignerPDA, mcmSignerPDA
}

// assuming one out of the src and dst is solana and the other is evm
func DeployTransferableTokenSolanaV0_1_1(
	lggr logger.Logger,
	e cldf.Environment,
	evmChainSel, solChainSel uint64,
	evmDeployer *bind.TransactOpts,
	evmTokenName string,
) (*burn_mint_erc677.BurnMintERC677, *burn_mint_token_pool.BurnMintTokenPool, solana.PublicKey, error) {
	selectorFamily, err := chainsel.GetSelectorFamily(evmChainSel)
	if err != nil {
		return nil, nil, solana.PublicKey{}, err
	}
	if selectorFamily != chainsel.FamilyEVM {
		return nil, nil, solana.PublicKey{}, fmt.Errorf("evmChainSel %d is not an evm chain", evmChainSel)
	}
	selectorFamily, err = chainsel.GetSelectorFamily(solChainSel)
	if err != nil {
		return nil, nil, solana.PublicKey{}, err
	}
	if selectorFamily != chainsel.FamilySolana {
		return nil, nil, solana.PublicKey{}, fmt.Errorf("solChainSel %d is not a solana chain", solChainSel)
	}
	state, err := stateview.LoadOnchainState(e)
	if err != nil {
		return nil, nil, solana.PublicKey{}, err
	}

	addresses := e.ExistingAddresses
	// deploy evm token and pool
	evmToken, evmPool, err := deployTransferTokenOneEnd(lggr, e.BlockChains.EVMChains()[evmChainSel], evmDeployer, addresses, evmTokenName)
	if err != nil {
		return nil, nil, solana.PublicKey{}, err
	}
	// attach token and pool to the registry
	if err := attachTokenToTheRegistry(e.BlockChains.EVMChains()[evmChainSel], state.MustGetEVMChainState(evmChainSel), evmDeployer, evmToken.Address(), evmPool.Address()); err != nil {
		return nil, nil, solana.PublicKey{}, err
	}
	solDeployerKey := e.BlockChains.SolanaChains()[solChainSel].DeployerKey.PublicKey()

	// deploy solana token
	solTokenName := evmTokenName
	e, err = commoncs.Apply(nil, e,
		commoncs.Configure(
			// this makes the deployer the mint authority by default
			cldf.CreateLegacyChangeSet(ccipChangeSetSolanaV0_1_1.DeploySolanaToken),
			ccipChangeSetSolanaV0_1_1.DeploySolanaTokenConfig{
				ChainSelector:    solChainSel,
				TokenProgramName: shared.SPL2022Tokens,
				TokenDecimals:    9,
				TokenSymbol:      solTokenName,
				ATAList:          []string{solDeployerKey.String()},
				MintAmountToAddress: map[string]uint64{
					solDeployerKey.String(): uint64(1000e9),
				},
			},
		),
	)
	if err != nil {
		return nil, nil, solana.PublicKey{}, err
	}
	// find solana token address
	solAddresses, err := e.ExistingAddresses.AddressesForChain(solChainSel)
	if err != nil {
		return nil, nil, solana.PublicKey{}, err
	}
	solTokenAddress := solanastateview.FindSolanaAddress(
		cldf.TypeAndVersion{
			Type:    shared.SPL2022Tokens,
			Version: deployment.Version1_0_0,
			Labels:  cldf.NewLabelSet(solTokenName),
		},
		solAddresses,
	)
	bnm := shared.BurnMintTokenPool

	// deploy and configure solana token pool
	e, err = commoncs.Apply(nil, e,
		commoncs.Configure(
			// deploy token pool and set the burn/mint authority to the tokenPool
			cldf.CreateLegacyChangeSet(ccipChangeSetSolanaV0_1_1.E2ETokenPool),
			ccipChangeSetSolanaV0_1_1.E2ETokenPoolConfig{
				InitializeGlobalTokenPoolConfig: []ccipChangeSetSolanaV0_1_1.TokenPoolConfigWithMCM{
					{
						ChainSelector: solChainSel,
						TokenPubKey:   solTokenAddress,
						PoolType:      bnm,
						Metadata:      shared.CLLMetadata,
					},
				},
				AddTokenPoolAndLookupTable: []ccipChangeSetSolanaV0_1_1.AddTokenPoolAndLookupTableConfig{
					{
						ChainSelector: solChainSel,
						TokenPoolConfigs: []ccipChangeSetSolanaV0_1_1.TokenPoolConfig{
							{
								TokenPubKey: solTokenAddress,
								PoolType:    bnm,
								Metadata:    shared.CLLMetadata,
							},
						},
					},
				},
				RegisterTokenAdminRegistry: []ccipChangeSetSolanaV0_1_1.RegisterTokenAdminRegistryConfig{
					{
						ChainSelector: solChainSel,
						RegisterTokenConfigs: []ccipChangeSetSolanaV0_1_1.RegisterTokenConfig{
							{
								TokenPubKey:             solTokenAddress,
								TokenAdminRegistryAdmin: solDeployerKey,
								RegisterType:            ccipChangeSetSolanaV0_1_1.ViaGetCcipAdminInstruction,
							},
						},
					},
				},
				AcceptAdminRoleTokenAdminRegistry: []ccipChangeSetSolanaV0_1_1.AcceptAdminRoleTokenAdminRegistryConfig{
					{
						ChainSelector: solChainSel,
						AcceptAdminRoleTokenConfigs: []ccipChangeSetSolanaV0_1_1.AcceptAdminRoleTokenConfig{
							{
								TokenPubKey: solTokenAddress,
							},
						},
					},
				},
				SetPool: []ccipChangeSetSolanaV0_1_1.SetPoolConfig{
					{
						ChainSelector: solChainSel,
						SetPoolTokenConfigs: []ccipChangeSetSolanaV0_1_1.SetPoolTokenConfig{
							{
								TokenPubKey:     solTokenAddress,
								PoolType:        bnm,
								Metadata:        shared.CLLMetadata,
								WritableIndexes: []uint8{3, 4, 7},
							},
						},
					},
				},
				RemoteChainTokenPool: []ccipChangeSetSolanaV0_1_1.SetupTokenPoolForRemoteChainConfig{
					{
						SolChainSelector: solChainSel,
						RemoteTokenPoolConfigs: []ccipChangeSetSolanaV0_1_1.RemoteChainTokenPoolConfig{
							{
								SolTokenPubKey: solTokenAddress,
								SolPoolType:    bnm,
								Metadata:       shared.CLLMetadata,
								EVMRemoteConfigs: map[uint64]ccipChangeSetSolanaV0_1_1.EVMRemoteConfig{
									evmChainSel: {
										TokenSymbol: shared.TokenSymbol(evmTokenName),
										PoolType:    shared.BurnMintTokenPool,
										PoolVersion: shared.CurrentTokenPoolVersion,
										RateLimiterConfig: ccipChangeSetSolanaV0_1_1.RateLimiterConfig{
											Inbound: solTestTokenPoolV0_1_1.RateLimitConfig{
												Enabled:  false,
												Capacity: 0,
												Rate:     0,
											},
											Outbound: solTestTokenPoolV0_1_1.RateLimitConfig{
												Enabled:  false,
												Capacity: 0,
												Rate:     0,
											},
										},
									},
								},
							},
						},
					},
				},
			},
		),
	)
	if err != nil {
		return nil, nil, solana.PublicKey{}, err
	}

	// configure evm
	poolConfigPDA, err := soltokens.TokenPoolConfigAddress(solTokenAddress, state.SolChains[solChainSel].BurnMintTokenPools[shared.CLLMetadata])
	if err != nil {
		return nil, nil, solana.PublicKey{}, err
	}
	err = setTokenPoolCounterPart(e.BlockChains.EVMChains()[evmChainSel], evmPool, evmDeployer, solChainSel, solTokenAddress.Bytes(), poolConfigPDA.Bytes())
	if err != nil {
		return nil, nil, solana.PublicKey{}, err
	}

	err = grantMintBurnPermissions(lggr, e.BlockChains.EVMChains()[evmChainSel], evmToken, evmDeployer, evmPool.Address())
	if err != nil {
		return nil, nil, solana.PublicKey{}, err
	}

	return evmToken, evmPool, solTokenAddress, nil
}

func AddLaneSolanaChangesetsV0_1_1(e *DeployedEnv, solChainSelector, remoteChainSelector uint64, remoteFamily string) []commoncs.ConfiguredChangeSet {
	var chainFamilySelector [4]uint8
	switch remoteFamily {
	case chainsel.FamilyEVM:
		// bytes4(keccak256("CCIP ChainFamilySelector EVM"))
		chainFamilySelector = [4]uint8{40, 18, 213, 44}
	case chainsel.FamilySolana:
		// bytes4(keccak256("CCIP ChainFamilySelector SVM"));
		chainFamilySelector = [4]uint8{30, 16, 189, 196}
	case chainsel.FamilyAptos:
		// bytes4(keccak256("CCIP ChainFamilySelector APTOS"));
		chainFamilySelector = [4]uint8{0xac, 0x77, 0xff, 0xec}
	default:
		panic("unsupported remote family")
	}
	solanaChangesets := []commoncs.ConfiguredChangeSet{
		commoncs.Configure(
			cldf.CreateLegacyChangeSet(ccipChangeSetSolanaV0_1_1.AddRemoteChainToRouter),
			ccipChangeSetSolanaV0_1_1.AddRemoteChainToRouterConfig{
				ChainSelector: solChainSelector,
				UpdatesByChain: map[uint64]*ccipChangeSetSolanaV0_1_1.RouterConfig{
					remoteChainSelector: {
						RouterDestinationConfig: solRouter.DestChainConfig{
							AllowListEnabled: true,
							AllowedSenders:   []solana.PublicKey{e.Env.BlockChains.SolanaChains()[solChainSelector].DeployerKey.PublicKey()},
						},
					},
				},
			},
		),
		commoncs.Configure(
			cldf.CreateLegacyChangeSet(ccipChangeSetSolanaV0_1_1.AddRemoteChainToFeeQuoter),
			ccipChangeSetSolanaV0_1_1.AddRemoteChainToFeeQuoterConfig{
				ChainSelector: solChainSelector,
				UpdatesByChain: map[uint64]*ccipChangeSetSolanaV0_1_1.FeeQuoterConfig{
					remoteChainSelector: {
						FeeQuoterDestinationConfig: solFeeQuoter.DestChainConfig{
							IsEnabled:                   true,
							DefaultTxGasLimit:           200000,
							MaxPerMsgGasLimit:           3000000,
							MaxDataBytes:                30000,
							MaxNumberOfTokensPerMsg:     5,
							DefaultTokenDestGasOverhead: 90000,
							DestGasOverhead:             90000,
							ChainFamilySelector:         chainFamilySelector,
						},
					},
				},
			},
		),
		commoncs.Configure(
			cldf.CreateLegacyChangeSet(ccipChangeSetSolanaV0_1_1.AddRemoteChainToOffRamp),
			ccipChangeSetSolanaV0_1_1.AddRemoteChainToOffRampConfig{
				ChainSelector: solChainSelector,
				UpdatesByChain: map[uint64]*ccipChangeSetSolanaV0_1_1.OffRampConfig{
					remoteChainSelector: {
						EnabledAsSource: true,
					},
				},
			},
		),
	}
	return solanaChangesets
}
