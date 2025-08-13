package testhelpers

import (
	"errors"
	"fmt"
	"maps"
	"math/big"
	"slices"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	chainsel "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_5_1/burn_mint_token_pool"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/onramp"
	solconfig "github.com/smartcontractkit/chainlink-ccip/chains/solana/contracts/tests/config"
	solCommon "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/v0_1_1/ccip_common"
	solRouter "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/v0_1_1/ccip_router"
	solFeeQuoter "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/v0_1_1/fee_quoter"
	solTestTokenPoolV0_1_1 "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/v0_1_1/test_token_pool"
	solccip "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/ccip"
	solcommon "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/common"
	solstate "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/state"
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

// SendRequest similar to TestSendRequest but returns an error.
func SendRequestV0_1_1(
	e cldf.Environment,
	state stateview.CCIPOnChainState,
	opts ...SendReqOpts,
) (*AnyMsgSentEvent, error) {
	cfg := &CCIPSendReqConfig{}
	for _, opt := range opts {
		opt(cfg)
	}
	family, err := chainsel.GetSelectorFamily(cfg.SourceChain)
	if err != nil {
		return nil, err
	}

	switch family {
	case chainsel.FamilyEVM:
		return SendRequestEVM(e, state, cfg)
	case chainsel.FamilySolana:
		return SendRequestSolV0_1_1(e, state, cfg)
	case chainsel.FamilyAptos:
		return SendRequestAptos(e, state, cfg)
	default:
		return nil, fmt.Errorf("send request: unsupported chain family: %v", family)
	}
}

func SendRequestSolV0_1_1(
	e cldf.Environment,
	state stateview.CCIPOnChainState,
	cfg *CCIPSendReqConfig,
) (*AnyMsgSentEvent, error) { // TODO: chain independent return value
	ctx := e.GetContext()

	s := state.SolChains[cfg.SourceChain]
	c := e.BlockChains.SolanaChains()[cfg.SourceChain]

	destinationChainSelector := cfg.DestChain
	message := cfg.Message.(solRouter.SVM2AnyMessage)
	feeToken := message.FeeToken
	client := c.Client

	// TODO: sender from cfg is EVM specific - need to revisit for Solana
	sender := c.DeployerKey

	e.Logger.Infof("Sending CCIP request from chain selector %d to chain selector %d from sender %s",
		cfg.SourceChain, cfg.DestChain, sender.PublicKey().String())

	feeTokenProgramID := solana.TokenProgramID
	feeTokenUserATA := solana.PublicKey{}
	if feeToken.IsZero() {
		// If the fee token is native SOL (i.e. message.FeeToken is the zero address), then we will
		// leave message.FeeToken as it is, but specify the WSOL mint account in the accounts list
		feeToken = solana.SolMint
	} else {
		feeTokenInfo, err := client.GetAccountInfo(ctx, feeToken)
		if err != nil {
			return nil, err
		}
		feeTokenProgramID = feeTokenInfo.Value.Owner

		_, err = GetSolanaTokenMintInfo(feeTokenInfo)
		if err != nil {
			return nil, fmt.Errorf("the provided fee token is not a valid token: (err = %w)", err)
		}

		ata, _, err := soltokens.FindAssociatedTokenAddress(feeTokenProgramID, feeToken, sender.PublicKey())
		if err != nil {
			return nil, err
		}

		feeTokenUserATA = ata
	}

	destinationChainStatePDA, err := solstate.FindDestChainStatePDA(destinationChainSelector, s.Router)
	if err != nil {
		return nil, err
	}

	noncePDA, err := solstate.FindNoncePDA(cfg.DestChain, sender.PublicKey(), s.Router)
	if err != nil {
		return nil, err
	}

	linkFqBillingConfigPDA, _, err := solstate.FindFqBillingTokenConfigPDA(s.LinkToken, s.FeeQuoter)
	if err != nil {
		return nil, err
	}

	feeTokenFqBillingConfigPDA, _, err := solstate.FindFqBillingTokenConfigPDA(feeToken, s.FeeQuoter)
	if err != nil {
		return nil, err
	}

	billingSignerPDA, _, err := solstate.FindFeeBillingSignerPDA(s.Router)
	if err != nil {
		return nil, err
	}

	feeTokenReceiverATA, _, err := soltokens.FindAssociatedTokenAddress(feeTokenProgramID, feeToken, billingSignerPDA)
	if err != nil {
		return nil, err
	}

	fqDestChainPDA, _, err := solstate.FindFqDestChainPDA(cfg.DestChain, s.FeeQuoter)
	if err != nil {
		return nil, err
	}

	rmnRemoteCursesPDA, _, err := solstate.FindRMNRemoteCursesPDA(s.RMNRemote)
	if err != nil {
		return nil, err
	}

	base := solRouter.NewCcipSendInstruction(
		destinationChainSelector,
		message,
		[]byte{}, // starting indices for accounts, calculated later
		s.RouterConfigPDA,
		destinationChainStatePDA,
		noncePDA,
		sender.PublicKey(),
		solana.SystemProgramID,
		feeTokenProgramID,
		feeToken,
		feeTokenUserATA,
		feeTokenReceiverATA,
		billingSignerPDA,
		s.FeeQuoter,
		s.FeeQuoterConfigPDA,
		fqDestChainPDA,
		feeTokenFqBillingConfigPDA,
		linkFqBillingConfigPDA,
		s.RMNRemote,
		rmnRemoteCursesPDA,
		s.RMNRemoteConfigPDA,
	)

	// When paying with a non-native token (i.e. any SPL token), the user ATA must be writable so we
	// can debit the fees. If paying with native SOL, then the ATA passed in is just a zero-address
	// placeholder, and that can't be marked as writable.
	if !feeTokenUserATA.IsZero() {
		base.GetFeeTokenUserAssociatedAccountAccount().WRITE()
	}

	addressTables := map[solana.PublicKey]solana.PublicKeySlice{}

	requiredAccounts := len(base.AccountMetaSlice)
	tokenIndexes := []byte{}

	// set config.FeeQuoterProgram and CcipRouterProgram since they point to wrong addresses
	solconfig.FeeQuoterProgram = s.FeeQuoter
	solconfig.CcipRouterProgram = s.Router

	// Append token accounts to the account metas
	for _, tokenAmount := range message.TokenAmounts {
		tokenPubKey := tokenAmount.Token

		allTokenPools := solana.PublicKeySlice{}
		allTokenPools = slices.AppendSeq(allTokenPools, maps.Values(s.LockReleaseTokenPools))
		allTokenPools = slices.AppendSeq(allTokenPools, maps.Values(s.BurnMintTokenPools))
		allTokenPools = append(allTokenPools, s.CCTPTokenPool)

		e.Logger.Infof("Found %d token pools in state - searching for matching token pool", len(allTokenPools))
		tokenPoolPubKey, err := MatchTokenToTokenPool(ctx, client, tokenPubKey, allTokenPools)
		if err != nil {
			return nil, err
		}

		e.Logger.Infof("Token '%s' was matched to token pool '%s'",
			tokenPubKey.String(),
			tokenPoolPubKey.String(),
		)

		tokenProgramID, err := InferSolanaTokenProgramID(ctx, client, tokenPubKey)
		if err != nil {
			return nil, err
		}

		tokenPool, err := soltokens.NewTokenPool(tokenProgramID, tokenPoolPubKey, tokenPubKey)
		if err != nil {
			return nil, err
		}

		// Set the token pool's lookup table address
		var tokenAdminRegistry solCommon.TokenAdminRegistry
		err = solcommon.GetAccountDataBorshInto(ctx, client, tokenPool.AdminRegistryPDA, solconfig.DefaultCommitment, &tokenAdminRegistry)
		if err != nil {
			return nil, err
		}

		tokenPool.PoolLookupTable = tokenAdminRegistry.LookupTable

		// invalid config account, maybe this billing stuff isn't right

		chainPDA, _, err := soltokens.TokenPoolChainConfigPDA(cfg.DestChain, tokenPubKey, tokenPoolPubKey)
		if err != nil {
			return nil, err
		}

		tokenPool.Chain[cfg.DestChain] = chainPDA

		billingPDA, _, err := solstate.FindFqPerChainPerTokenConfigPDA(cfg.DestChain, tokenPubKey, s.FeeQuoter)
		if err != nil {
			return nil, err
		}

		tokenPool.Billing[cfg.DestChain] = billingPDA

		userTokenAccount, _, err := soltokens.FindAssociatedTokenAddress(tokenProgramID, tokenPubKey, sender.PublicKey())
		if err != nil {
			return nil, err
		}

		tokenMetas, tokenAddressTables, err := soltokens.ParseTokenLookupTableWithChain(ctx, client, tokenPool, userTokenAccount, cfg.DestChain)
		if err != nil {
			return nil, err
		}

		tokenIndexes = append(tokenIndexes, byte(len(base.AccountMetaSlice)-requiredAccounts))
		base.AccountMetaSlice = append(base.AccountMetaSlice, tokenMetas...)
		maps.Copy(addressTables, tokenAddressTables)
	}

	base.SetTokenIndexes(tokenIndexes)

	tempIx, err := base.ValidateAndBuild()
	if err != nil {
		return nil, err
	}
	ixData, err := tempIx.Data()
	if err != nil {
		return nil, fmt.Errorf("failed to extract data payload from router ccip send instruction: %w", err)
	}
	ix := solana.NewInstruction(s.Router, tempIx.Accounts(), ixData)

	// for some reason onchain doesn't see extraAccounts

	ixs := []solana.Instruction{ix}
	result, err := solcommon.SendAndConfirmWithLookupTables(ctx, client, ixs, *sender, solconfig.DefaultCommitment, addressTables, solcommon.AddComputeUnitLimit(400_000))
	if err != nil {
		return nil, err
	}

	// check CCIP event
	ccipMessageSentEvent := solccip.EventCCIPMessageSent{}
	printEvents := true
	err = solcommon.ParseEvent(result.Meta.LogMessages, "CCIPMessageSent", &ccipMessageSentEvent, printEvents)
	if err != nil {
		return nil, err
	}

	if len(message.TokenAmounts) != len(ccipMessageSentEvent.Message.TokenAmounts) {
		return nil, errors.New("token amounts mismatch")
	}

	// TODO: fee bumping?

	transactionID := "N/A"
	if tx, err := result.Transaction.GetTransaction(); err != nil {
		e.Logger.Warnf("could not obtain transaction details (err = %s)", err.Error())
	} else if len(tx.Signatures) == 0 {
		e.Logger.Warnf("transaction has no signatures: %v", tx)
	} else {
		transactionID = tx.Signatures[0].String()
	}

	e.Logger.Infof("CCIP message (id %s) sent from chain selector %d to chain selector %d tx %s seqNum %d nonce %d sender %s testRouterEnabled %t",
		common.Bytes2Hex(ccipMessageSentEvent.Message.Header.MessageId[:]),
		cfg.SourceChain,
		cfg.DestChain,
		transactionID,
		ccipMessageSentEvent.SequenceNumber,
		ccipMessageSentEvent.Message.Header.Nonce,
		ccipMessageSentEvent.Message.Sender.String(),
		cfg.IsTestRouter,
	)

	return &AnyMsgSentEvent{
		SequenceNumber: ccipMessageSentEvent.SequenceNumber,
		RawEvent: &onramp.OnRampCCIPMessageSent{
			DestChainSelector: ccipMessageSentEvent.DestinationChainSelector,
			SequenceNumber:    ccipMessageSentEvent.SequenceNumber,
			Message: onramp.InternalEVM2AnyRampMessage{
				Header: onramp.InternalRampMessageHeader{
					SourceChainSelector: ccipMessageSentEvent.Message.Header.SourceChainSelector,
					DestChainSelector:   ccipMessageSentEvent.Message.Header.DestChainSelector,
					MessageId:           ccipMessageSentEvent.Message.Header.MessageId,
					SequenceNumber:      ccipMessageSentEvent.SequenceNumber,
					Nonce:               ccipMessageSentEvent.Message.Header.Nonce,
				},
				FeeTokenAmount: ConvertSolanaCrossChainAmountToBigInt(ccipMessageSentEvent.Message.FeeTokenAmount.LeBytes),
				FeeValueJuels:  ConvertSolanaCrossChainAmountToBigInt(ccipMessageSentEvent.Message.FeeValueJuels.LeBytes),
				ExtraArgs:      ccipMessageSentEvent.Message.ExtraArgs,
				Receiver:       ccipMessageSentEvent.Message.Receiver,
				Data:           ccipMessageSentEvent.Message.Data,

				// TODO: these fields are EVM specific - need to revisit for Solana
				FeeToken:     common.Address{}, // ccipMessageSentEvent.Message.FeeToken
				Sender:       common.Address{}, // ccipMessageSentEvent.Message.Sender
				TokenAmounts: []onramp.InternalEVM2AnyTokenTransfer{},
			},

			// TODO: EVM specific - need to revisit for Solana
			Raw: types.Log{},
		},
	}, nil
}
