package testhelpers

import (
	"context"
	"errors"
	"fmt"
	"maps"
	"math/big"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	solbinary "github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/programs/token"
	"github.com/gagliardetto/solana-go/rpc"
	chainsel "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_5_1/burn_mint_token_pool"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/onramp"
	solconfig "github.com/smartcontractkit/chainlink-ccip/chains/solana/contracts/tests/config"
	"github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/base_token_pool"
	solCommon "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/ccip_common"
	solOffRamp "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/ccip_offramp"
	"github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/ccip_router"
	solRouter "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/ccip_router"
	solFeeQuoter "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/fee_quoter"
	solRmnRemote "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/rmn_remote"
	solTestReceiver "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/test_ccip_receiver"
	solTestTokenPool "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/test_token_pool"
	solccip "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/ccip"
	solcommon "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/common"
	solstate "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/state"
	soltokens "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/tokens"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/shared/generated/burn_mint_erc677"
	"github.com/smartcontractkit/chainlink/deployment"
	ccipChangeSetSolana "github.com/smartcontractkit/chainlink/deployment/ccip/changeset/solana"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
	solanastateview "github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview/solana"
	commoncs "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/common/changeset/state"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	commontypes "github.com/smartcontractkit/chainlink/deployment/common/types"
	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
	"github.com/stretchr/testify/require"
)

func addLaneSolanaChangesets(t *testing.T, e *DeployedEnv, solChainSelector, remoteChainSelector uint64, remoteFamily string) []commoncs.ConfiguredChangeSet {
	chainFamilySelector := [4]uint8{}
	if remoteFamily == chainsel.FamilyEVM {
		// bytes4(keccak256("CCIP ChainFamilySelector EVM"))
		chainFamilySelector = [4]uint8{40, 18, 213, 44}
	} else if remoteFamily == chainsel.FamilySolana {
		// bytes4(keccak256("CCIP ChainFamilySelector SVM"));
		chainFamilySelector = [4]uint8{30, 16, 189, 196}
	}
	solanaChangesets := []commoncs.ConfiguredChangeSet{
		commoncs.Configure(
			cldf.CreateLegacyChangeSet(ccipChangeSetSolana.AddRemoteChainToRouter),
			ccipChangeSetSolana.AddRemoteChainToRouterConfig{
				ChainSelector: solChainSelector,
				UpdatesByChain: map[uint64]*ccipChangeSetSolana.RouterConfig{
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
			cldf.CreateLegacyChangeSet(ccipChangeSetSolana.AddRemoteChainToFeeQuoter),
			ccipChangeSetSolana.AddRemoteChainToFeeQuoterConfig{
				ChainSelector: solChainSelector,
				UpdatesByChain: map[uint64]*ccipChangeSetSolana.FeeQuoterConfig{
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
			cldf.CreateLegacyChangeSet(ccipChangeSetSolana.AddRemoteChainToOffRamp),
			ccipChangeSetSolana.AddRemoteChainToOffRampConfig{
				ChainSelector: solChainSelector,
				UpdatesByChain: map[uint64]*ccipChangeSetSolana.OffRampConfig{
					remoteChainSelector: {
						EnabledAsSource: true,
					},
				},
			},
		),
	}
	return solanaChangesets
}

func ConvertSolanaCrossChainAmountToBigInt(amount ccip_router.CrossChainAmount) *big.Int {
	bytes := amount.LeBytes[:]
	slices.Reverse(bytes) // convert to big-endian
	return big.NewInt(0).SetBytes(bytes)
}

func SendRequestSol(
	e cldf.Environment,
	state stateview.CCIPOnChainState,
	cfg *CCIPSendReqConfig,
) (*onramp.OnRampCCIPMessageSent, error) { // TODO: chain independent return value
	ctx := e.GetContext()

	s := state.SolChains[cfg.SourceChain]
	c := e.BlockChains.SolanaChains()[cfg.SourceChain]

	destinationChainSelector := cfg.DestChain
	message := cfg.Message.(ccip_router.SVM2AnyMessage)
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

	base := ccip_router.NewCcipSendInstruction(
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

	ix, err := base.ValidateAndBuild()
	if err != nil {
		return nil, err
	}

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

	return &onramp.OnRampCCIPMessageSent{
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
			FeeTokenAmount: ConvertSolanaCrossChainAmountToBigInt(ccipMessageSentEvent.Message.FeeTokenAmount),
			FeeValueJuels:  ConvertSolanaCrossChainAmountToBigInt(ccipMessageSentEvent.Message.FeeValueJuels),
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
	}, nil
}

func InferSolanaTokenProgramID(ctx context.Context, client *rpc.Client, tokenPubKey solana.PublicKey) (solana.PublicKey, error) {
	tokenAcctInfo, err := client.GetAccountInfo(ctx, tokenPubKey)
	if errors.Is(err, rpc.ErrNotFound) {
		// NOTE: we use a fallback value of Token2022ProgramID to maintain backwards compatibility with the Solana tests
		return solana.Token2022ProgramID, nil
	}
	if err != nil {
		return solana.PublicKey{}, err
	}

	_, err = GetSolanaTokenMintInfo(tokenAcctInfo)
	if err != nil {
		return solana.PublicKey{}, fmt.Errorf("expected '%s' to be a token public key: (err = %w)", tokenPubKey, err)
	}

	return tokenAcctInfo.Value.Owner, nil
}

func GetSolanaTokenMintInfo(tokenAcctInfo *rpc.GetAccountInfoResult) (token.Mint, error) {
	var mint token.Mint

	err := solbinary.NewBinDecoder(tokenAcctInfo.Bytes()).Decode(&mint)
	if err != nil {
		return token.Mint{}, fmt.Errorf("failed to decode token mint data: (err = %w)", err)
	}

	return mint, nil
}

func MatchTokenToTokenPool(ctx context.Context, client *rpc.Client, tokenPubKey solana.PublicKey, tokenPoolPubKeys solana.PublicKeySlice) (solana.PublicKey, error) {
	for _, tokenPoolPubKey := range tokenPoolPubKeys {
		tokenPoolConfigAddress, err := soltokens.TokenPoolConfigAddress(tokenPubKey, tokenPoolPubKey)
		if err != nil {
			return solana.PublicKey{}, err
		}

		var tokenPoolConfig base_token_pool.BaseConfig
		err = solcommon.GetAccountDataBorshInto(ctx, client, tokenPoolConfigAddress, solconfig.DefaultCommitment, &tokenPoolConfig)
		if errors.Is(err, rpc.ErrNotFound) {
			continue
		}
		if err != nil {
			return solana.PublicKey{}, err
		}

		return tokenPoolPubKey, nil
	}

	tokenPoolPubKeyStrs := make([]string, len(tokenPoolPubKeys))
	for i, tokenPoolPubKey := range tokenPoolPubKeys {
		tokenPoolPubKeyStrs[i] = "'" + tokenPoolPubKey.String() + "'"
	}

	msg := "token with public key '%s' is not associated with any of the following token pools: [ %s ]"
	return solana.PublicKey{}, fmt.Errorf(msg, tokenPubKey.String(), strings.Join(tokenPoolPubKeyStrs, ", "))
}

// assuming one out of the src and dst is solana and the other is evm
func DeployTransferableTokenSolana(
	lggr logger.Logger,
	e cldf.Environment,
	evmChainSel, solChainSel uint64,
	evmDeployer *bind.TransactOpts,
	evmTokenName string,
) (*burn_mint_erc677.BurnMintERC677,
	*burn_mint_token_pool.BurnMintTokenPool, solana.PublicKey, error) {
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

	addresses := e.ExistingAddresses //nolint:staticcheck // addressbook still valid
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
	e, err = commoncs.Apply(nil, e, nil,
		commoncs.Configure(
			// this makes the deployer the mint authority by default
			cldf.CreateLegacyChangeSet(ccipChangeSetSolana.DeploySolanaToken),
			ccipChangeSetSolana.DeploySolanaTokenConfig{
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
	bnm := solTestTokenPool.BurnAndMint_PoolType

	// deploy and configure solana token pool
	e, err = commoncs.Apply(nil, e, nil,
		commoncs.Configure(
			// deploy token pool and set the burn/mint authority to the tokenPool
			cldf.CreateLegacyChangeSet(ccipChangeSetSolana.E2ETokenPool),
			ccipChangeSetSolana.E2ETokenPoolConfig{
				AddTokenPoolAndLookupTable: []ccipChangeSetSolana.TokenPoolConfig{
					{
						ChainSelector: solChainSel,
						TokenPubKey:   solTokenAddress,
						PoolType:      &bnm,
						Metadata:      shared.CLLMetadata,
					},
				},
				RegisterTokenAdminRegistry: []ccipChangeSetSolana.RegisterTokenAdminRegistryConfig{
					{
						ChainSelector:           solChainSel,
						TokenPubKey:             solTokenAddress,
						TokenAdminRegistryAdmin: solDeployerKey.String(),
						RegisterType:            ccipChangeSetSolana.ViaGetCcipAdminInstruction,
					},
				},
				AcceptAdminRoleTokenAdminRegistry: []ccipChangeSetSolana.AcceptAdminRoleTokenAdminRegistryConfig{
					{
						ChainSelector: solChainSel,
						TokenPubKey:   solTokenAddress,
					},
				},
				SetPool: []ccipChangeSetSolana.SetPoolConfig{
					{
						ChainSelector:   solChainSel,
						TokenPubKey:     solTokenAddress,
						PoolType:        &bnm,
						Metadata:        shared.CLLMetadata,
						WritableIndexes: []uint8{3, 4, 7},
					},
				},
				RemoteChainTokenPool: []ccipChangeSetSolana.RemoteChainTokenPoolConfig{
					{
						SolChainSelector: solChainSel,
						SolTokenPubKey:   solTokenAddress,
						SolPoolType:      &bnm,
						Metadata:         shared.CLLMetadata,
						EVMRemoteConfigs: map[uint64]ccipChangeSetSolana.EVMRemoteConfig{
							evmChainSel: {
								TokenSymbol: shared.TokenSymbol(evmTokenName),
								PoolType:    shared.BurnMintTokenPool,
								PoolVersion: shared.CurrentTokenPoolVersion,
								RateLimiterConfig: ccipChangeSetSolana.RateLimiterConfig{
									Inbound: solTestTokenPool.RateLimitConfig{
										Enabled:  false,
										Capacity: 0,
										Rate:     0,
									},
									Outbound: solTestTokenPool.RateLimitConfig{
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

// TODO: this should be linked to the solChain function
func SavePreloadedSolAddresses(e cldf.Environment, solChainSelector uint64) error {
	tv := cldf.NewTypeAndVersion(shared.Router, deployment.Version1_0_0)
	err := e.ExistingAddresses.Save(solChainSelector, memory.SolanaProgramIDs["ccip_router"], tv)
	if err != nil {
		return err
	}
	tv = cldf.NewTypeAndVersion(shared.Receiver, deployment.Version1_0_0)
	err = e.ExistingAddresses.Save(solChainSelector, memory.SolanaProgramIDs["test_ccip_receiver"], tv)
	if err != nil {
		return err
	}
	tv = cldf.NewTypeAndVersion(shared.FeeQuoter, deployment.Version1_0_0)
	err = e.ExistingAddresses.Save(solChainSelector, memory.SolanaProgramIDs["fee_quoter"], tv)
	if err != nil {
		return err
	}
	tv = cldf.NewTypeAndVersion(shared.OffRamp, deployment.Version1_0_0)
	err = e.ExistingAddresses.Save(solChainSelector, memory.SolanaProgramIDs["ccip_offramp"], tv)
	if err != nil {
		return err
	}
	tv = cldf.NewTypeAndVersion(shared.BurnMintTokenPool, deployment.Version1_0_0)
	err = e.ExistingAddresses.Save(solChainSelector, memory.SolanaProgramIDs["burnmint_token_pool"], tv)
	if err != nil {
		return err
	}
	tv = cldf.NewTypeAndVersion(shared.LockReleaseTokenPool, deployment.Version1_0_0)
	err = e.ExistingAddresses.Save(solChainSelector, memory.SolanaProgramIDs["lockrelease_token_pool"], tv)
	if err != nil {
		return err
	}
	tv = cldf.NewTypeAndVersion(commontypes.ManyChainMultisigProgram, deployment.Version1_0_0)
	err = e.ExistingAddresses.Save(solChainSelector, memory.SolanaProgramIDs["mcm"], tv)
	if err != nil {
		return err
	}
	tv = cldf.NewTypeAndVersion(commontypes.AccessControllerProgram, deployment.Version1_0_0)
	err = e.ExistingAddresses.Save(solChainSelector, memory.SolanaProgramIDs["access_controller"], tv)
	if err != nil {
		return err
	}
	tv = cldf.NewTypeAndVersion(commontypes.RBACTimelockProgram, deployment.Version1_0_0)
	err = e.ExistingAddresses.Save(solChainSelector, memory.SolanaProgramIDs["timelock"], tv)
	if err != nil {
		return err
	}
	tv = cldf.NewTypeAndVersion(shared.RMNRemote, deployment.Version1_0_0)
	err = e.ExistingAddresses.Save(solChainSelector, memory.SolanaProgramIDs["rmn_remote"], tv)
	if err != nil {
		return err
	}
	return nil
}

func ValidateSolanaState(e cldf.Environment, solChainSelectors []uint64) error {
	state, err := stateview.LoadOnchainStateSolana(e)
	if err != nil {
		return fmt.Errorf("failed to load Solana state: %w", err)
	}

	for _, sel := range solChainSelectors {
		// Validate chain exists in state
		chainState, exists := state.SolChains[sel]
		if !exists {
			return fmt.Errorf("chain selector %d not found in Solana state", sel)
		}

		// Validate addresses
		if chainState.Router.IsZero() {
			return fmt.Errorf("router address is zero for chain %d", sel)
		}
		if chainState.OffRamp.IsZero() {
			return fmt.Errorf("offRamp address is zero for chain %d", sel)
		}
		if chainState.FeeQuoter.IsZero() {
			return fmt.Errorf("feeQuoter address is zero for chain %d", sel)
		}
		if chainState.LinkToken.IsZero() {
			return fmt.Errorf("link token address is zero for chain %d", sel)
		}
		if chainState.RMNRemote.IsZero() {
			return fmt.Errorf("RMNRemote address is zero for chain %d", sel)
		}

		// Get router config
		var routerConfigAccount solRouter.Config
		err = e.BlockChains.SolanaChains()[sel].GetAccountDataBorshInto(context.Background(), chainState.RouterConfigPDA, &routerConfigAccount)
		if err != nil {
			return fmt.Errorf("failed to deserialize router config for chain %d: %w", sel, err)
		}

		// Get fee quoter config
		var feeQuoterConfigAccount solFeeQuoter.Config
		err = e.BlockChains.SolanaChains()[sel].GetAccountDataBorshInto(context.Background(), chainState.FeeQuoterConfigPDA, &feeQuoterConfigAccount)
		if err != nil {
			return fmt.Errorf("failed to deserialize fee quoter config for chain %d: %w", sel, err)
		}

		// Get offramp config
		var offRampConfigAccount solOffRamp.Config
		err = e.BlockChains.SolanaChains()[sel].GetAccountDataBorshInto(
			context.Background(),
			chainState.OffRampConfigPDA,
			&offRampConfigAccount,
		)
		if err != nil {
			return fmt.Errorf("failed to deserialize off-ramp config for chain %d: %w", sel, err)
		}
		if err != nil {
			return fmt.Errorf("failed to deserialize offramp config for chain %d: %w", sel, err)
		}

		// Get rmn remote config
		var rmnRemoteConfigAccount solRmnRemote.Config
		err = e.BlockChains.SolanaChains()[sel].GetAccountDataBorshInto(context.Background(), chainState.RMNRemoteConfigPDA, &rmnRemoteConfigAccount)
		if err != nil {
			return fmt.Errorf("failed to deserialize rmn remote config for chain %d: %w", sel, err)
		}

		addressLookupTable, err := solanastateview.FetchOfframpLookupTable(e.GetContext(), e.BlockChains.SolanaChains()[sel], chainState.OffRamp)
		if err != nil {
			return fmt.Errorf("failed to get offramp lookup table for chain %d: %w", sel, err)
		}

		addresses, err := solcommon.GetAddressLookupTable(
			e.GetContext(),
			e.BlockChains.SolanaChains()[sel].Client,
			addressLookupTable,
		)
		if err != nil {
			return fmt.Errorf("failed to get address lookup table for chain %d: %w", sel, err)
		}
		if len(addresses) < 22 {
			return fmt.Errorf("not enough addresses found in lookup table for chain %d: got %d, expected at least 22", sel, len(addresses))
		}
	}
	return nil
}

func DeploySolanaCcipReceiver(t *testing.T, e cldf.Environment) {
	state, err := stateview.LoadOnchainStateSolana(e)
	require.NoError(t, err)
	for solSelector, chainState := range state.SolChains {
		solTestReceiver.SetProgramID(chainState.Receiver)
		externalExecutionConfigPDA, _, _ := solana.FindProgramAddress([][]byte{[]byte("external_execution_config")}, chainState.Receiver)
		instruction, ixErr := solTestReceiver.NewInitializeInstruction(
			chainState.Router,
			solanastateview.FindReceiverTargetAccount(chainState.Receiver),
			externalExecutionConfigPDA,
			e.BlockChains.SolanaChains()[solSelector].DeployerKey.PublicKey(),
			solana.SystemProgramID,
		).ValidateAndBuild()
		require.NoError(t, ixErr)
		err = e.BlockChains.SolanaChains()[solSelector].Confirm([]solana.Instruction{instruction})
		require.NoError(t, err)
	}
}

func TransferOwnershipSolana(
	t *testing.T,
	e *cldf.Environment,
	solChain uint64,
	needTimelockDeployed bool,
	contractsToTransfer ccipChangeSetSolana.CCIPContractsToTransfer,
) (timelockSignerPDA solana.PublicKey, mcmSignerPDA solana.PublicKey) {
	var err error
	if needTimelockDeployed {
		*e, _, err = commoncs.ApplyChangesetsV2(t, *e, []commoncs.ConfiguredChangeSet{
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
	*e, _, err = commoncs.ApplyChangesetsV2(t, *e, []commoncs.ConfiguredChangeSet{
		commoncs.Configure(
			cldf.CreateLegacyChangeSet(ccipChangeSetSolana.TransferCCIPToMCMSWithTimelockSolana),
			ccipChangeSetSolana.TransferCCIPToMCMSWithTimelockSolanaConfig{
				MCMSCfg: proposalutils.TimelockConfig{MinDelay: 1 * time.Second},
				ContractsByChain: map[uint64]ccipChangeSetSolana.CCIPContractsToTransfer{
					solChain: contractsToTransfer,
				},
			},
		),
	})
	require.NoError(t, err)
	return timelockSignerPDA, mcmSignerPDA
}

func GenTestTransferOwnershipConfig(
	e DeployedEnv,
	chains []uint64,
	state stateview.CCIPOnChainState,
	withTestRouterTransfer bool,
) commoncs.TransferToMCMSWithTimelockConfig {
	var (
		timelocksPerChain = make(map[uint64]common.Address)
		contracts         = make(map[uint64][]common.Address)
	)

	// chain contracts
	for _, chain := range chains {
		timelocksPerChain[chain] = state.MustGetEVMChainState(chain).Timelock.Address()
		contracts[chain] = []common.Address{
			state.MustGetEVMChainState(chain).OnRamp.Address(),
			state.MustGetEVMChainState(chain).OffRamp.Address(),
			state.MustGetEVMChainState(chain).FeeQuoter.Address(),
			state.MustGetEVMChainState(chain).NonceManager.Address(),
			state.MustGetEVMChainState(chain).RMNRemote.Address(),
			state.MustGetEVMChainState(chain).Router.Address(),
			state.MustGetEVMChainState(chain).TokenAdminRegistry.Address(),
			state.MustGetEVMChainState(chain).RMNProxy.Address(),
		}
		if withTestRouterTransfer {
			contracts[chain] = append(contracts[chain], state.MustGetEVMChainState(chain).TestRouter.Address())
		}
	}

	// home chain
	homeChainTimelockAddress := state.MustGetEVMChainState(e.HomeChainSel).Timelock.Address()
	timelocksPerChain[e.HomeChainSel] = homeChainTimelockAddress
	contracts[e.HomeChainSel] = append(contracts[e.HomeChainSel],
		state.MustGetEVMChainState(e.HomeChainSel).CapabilityRegistry.Address(),
		state.MustGetEVMChainState(e.HomeChainSel).CCIPHome.Address(),
		state.MustGetEVMChainState(e.HomeChainSel).RMNHome.Address(),
	)

	return commoncs.TransferToMCMSWithTimelockConfig{
		ContractsByChain: contracts,
	}
}
