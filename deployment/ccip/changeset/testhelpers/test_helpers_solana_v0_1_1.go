package testhelpers

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"regexp"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/gagliardetto/solana-go"
	addresslookuptable "github.com/gagliardetto/solana-go/programs/address-lookup-table"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/stretchr/testify/require"

	chainsel "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_5_1/burn_mint_token_pool"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/onramp"
	solconfig "github.com/smartcontractkit/chainlink-ccip/chains/solana/contracts/tests/config"
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
	ccipclient "github.com/smartcontractkit/chainlink/deployment/ccip/shared/client"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
	solanastateview "github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview/solana"
	commoncs "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/common/changeset/state"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	commontypes "github.com/smartcontractkit/chainlink/deployment/common/types"
	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
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
						TokenPoolConfigs: []ccipChangeSetSolanaV0_1_1.TokenPoolConfig{
							{
								TokenPubKey: solTokenAddress,
								PoolType:    bnm,
								Metadata:    shared.CLLMetadata,
							},
						},
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
	opts ...ccipclient.SendReqOpts,
) (*ccipclient.AnyMsgSentEvent, error) {
	cfg := &ccipclient.CCIPSendReqConfig{}
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
	cfg *ccipclient.CCIPSendReqConfig,
) (*ccipclient.AnyMsgSentEvent, error) { // TODO: chain independent return value
	ctx := e.GetContext()

	s := state.SolChains[cfg.SourceChain]
	c := e.BlockChains.SolanaChains()[cfg.SourceChain]

	destinationChainSelector := cfg.DestChain
	message := cfg.Message.(solRouter.SVM2AnyMessage)
	client := c.Client

	// TODO: sender from cfg is EVM specific - need to revisit for Solana
	sender := c.DeployerKey

	e.Logger.Infof("Sending CCIP request from chain selector %d to chain selector %d from sender %s",
		cfg.SourceChain, cfg.DestChain, sender.PublicKey().String())

	accounts, lutAddresses, tokenIndexes, err := deriveCCIPSendAccounts(e, sender.PublicKey(), message, destinationChainSelector, client, s.Router)
	if err != nil {
		return nil, fmt.Errorf("failed to derive accounts from on-chain: %w", err)
	}

	builder := solRouter.NewCcipSendInstructionBuilder()
	builder.SetDestChainSelector(destinationChainSelector)
	builder.SetMessage(message)
	builder.SetTokenIndexes(tokenIndexes)
	err = builder.SetAccounts(accounts)
	if err != nil {
		return nil, fmt.Errorf("failed to set accounts in instruction builder: %w", err)
	}

	tempIx, err := builder.ValidateAndBuild()
	if err != nil {
		return nil, fmt.Errorf("failed to build ccip send message: %w", err)
	}
	ixData, err := tempIx.Data()
	if err != nil {
		return nil, fmt.Errorf("failed to extract payload data from instruction: %w", err)
	}

	ix := solana.NewInstruction(s.Router, tempIx.Accounts(), ixData)
	ixs := []solana.Instruction{ix}

	addressTables, err := fetchLookupTables(e.GetContext(), client, lutAddresses)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch lookup table addresses: %w", err)
	}

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

	return &ccipclient.AnyMsgSentEvent{
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

func deriveCCIPSendAccounts(
	e cldf.Environment,
	authority solana.PublicKey,
	message solRouter.SVM2AnyMessage,
	destChainSelector uint64,
	client *rpc.Client,
	router solana.PublicKey,
) (accounts []*solana.AccountMeta, lookUpTables []solana.PublicKey, tokenIndices []uint8, err error) {
	blockhash, err := client.GetLatestBlockhash(e.GetContext(), rpc.CommitmentConfirmed)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("error fetching latest blockhash: %w", err)
	}
	derivedAccounts := []*solana.AccountMeta{}
	askWith := []*solana.AccountMeta{}
	stage := "Start"
	tokenIndex := byte(0)
	routerConfigPDA, _, err := solstate.FindConfigPDA(router)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to calculate the router config PDA: %w", err)
	}
	var re = regexp.MustCompile(`^TokenTransferStaticAccounts/\d+/0$`)
	for {
		params := solRouter.DeriveAccountsCcipSendParams{
			DestChainSelector: destChainSelector,
			CcipSendCaller:    authority,
			Message:           message,
		}

		deriveRaw := solRouter.NewDeriveAccountsCcipSendInstruction(
			params,
			stage,
			routerConfigPDA,
		)
		deriveRaw.AccountMetaSlice = append(deriveRaw.AccountMetaSlice, askWith...)
		derive, err := deriveRaw.ValidateAndBuild()
		if err != nil {
			return nil, nil, nil, fmt.Errorf("failed to create derive account instruction: %w", err)
		}
		tx, err := solana.NewTransaction([]solana.Instruction{derive}, blockhash.Value.Blockhash, solana.TransactionPayer(authority))
		if err != nil {
			return nil, nil, nil, fmt.Errorf("failed to build derive ccip_send accounts transaction: %w", err)
		}
		tx.Signatures = append(tx.Signatures, solana.Signature{}) // Append empty signature since tx fails without any sigs even if SigVerify is false
		res, err := client.SimulateTransactionWithOpts(e.GetContext(), tx, &rpc.SimulateTransactionOpts{SigVerify: false, ReplaceRecentBlockhash: true})
		if err != nil {
			return nil, nil, nil, fmt.Errorf("failed to simulate derive ccip_send accounts transaction at stage %s: %w", stage, err)
		}
		if res.Value.Err != nil {
			return nil, nil, nil, fmt.Errorf("failed to simulate derive ccip_send accounts transaction at stage %s. Err: %v, Logs: %v", stage, res.Value.Err, res.Value.Logs)
		}
		derivation, err := solcommon.ExtractAnchorTypedReturnValue[solRouter.DeriveAccountsResponse](e.GetContext(), res.Value.Logs, router.String())
		if err != nil {
			e.Logger.Errorf("Error deriving accounts for stage %s: %v\n", stage, err)
			for _, line := range res.Value.Logs {
				e.Logger.Error(line)
			}
			return nil, nil, nil, fmt.Errorf("failed to exract accounts from simulated transaction log: %w", err)
		}
		e.Logger.Infof("Derive stage: %s. Len: %d\n", derivation.CurrentStage, len(derivation.AccountsToSave))

		isStartOfToken := re.MatchString(derivation.CurrentStage)
		if isStartOfToken {
			tokenIndices = append(tokenIndices, tokenIndex-byte(cap(solRouter.NewCcipSendInstructionBuilder().AccountMetaSlice)))
		}
		tokenIndex += byte(len(derivation.AccountsToSave))

		for _, meta := range derivation.AccountsToSave {
			derivedAccounts = append(derivedAccounts, &solana.AccountMeta{
				PublicKey:  meta.Pubkey,
				IsWritable: meta.IsWritable,
				IsSigner:   meta.IsSigner,
			})
		}
		askWith = []*solana.AccountMeta{}
		for _, meta := range derivation.AskAgainWith {
			askWith = append(askWith, &solana.AccountMeta{
				PublicKey:  meta.Pubkey,
				IsWritable: meta.IsWritable,
				IsSigner:   meta.IsSigner,
			})
		}
		lookUpTables = append(lookUpTables, derivation.LookUpTablesToSave...)

		if len(derivation.NextStage) == 0 {
			return derivedAccounts, lookUpTables, tokenIndices, nil
		}
		stage = derivation.NextStage
	}
}

func fetchLookupTables(ctx context.Context, client *rpc.Client, lookupTablesAddrs []solana.PublicKey) (map[solana.PublicKey]solana.PublicKeySlice, error) {
	lookupTableMap := make(map[solana.PublicKey]solana.PublicKeySlice)
	for _, addr := range lookupTablesAddrs {
		lookupTableContents, err := getLookupTableAddresses(ctx, client, addr)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch lookup table contents for address %s: %w", addr.String(), err)
		}
		lookupTableMap[addr] = lookupTableContents
	}
	return lookupTableMap, nil
}

func getLookupTableAddresses(ctx context.Context, client *rpc.Client, tableAddress solana.PublicKey) (solana.PublicKeySlice, error) {
	// Fetch the account info for the static table
	accountInfo, err := client.GetAccountInfoWithOpts(ctx, tableAddress, &rpc.GetAccountInfoOpts{
		Encoding:   "base64",
		Commitment: rpc.CommitmentFinalized,
	})

	if err != nil || accountInfo == nil || accountInfo.Value == nil {
		return nil, fmt.Errorf("error fetching account info for table: %s, error: %w", tableAddress.String(), err)
	}
	alt, err := addresslookuptable.DecodeAddressLookupTableState(accountInfo.GetBinary())
	if err != nil {
		return nil, fmt.Errorf("error decoding address lookup table state: %w", err)
	}
	return alt.Addresses, nil
}
