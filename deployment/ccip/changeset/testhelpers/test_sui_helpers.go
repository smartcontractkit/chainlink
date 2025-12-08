package testhelpers

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/block-vision/sui-go-sdk/models"
	"github.com/block-vision/sui-go-sdk/sui"
	suitx "github.com/block-vision/sui-go-sdk/transaction"
	"github.com/stretchr/testify/require"

	chainsel "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_5_1/burn_mint_token_pool"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_3/message_hasher"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"
	cldf_sui "github.com/smartcontractkit/chainlink-deployments-framework/chain/sui"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/shared/generated/initial/burn_mint_erc677"
	suiBind "github.com/smartcontractkit/chainlink-sui/bindings/bind"
	sui_deployment "github.com/smartcontractkit/chainlink-sui/deployment"
	sui_cs "github.com/smartcontractkit/chainlink-sui/deployment/changesets"
	sui_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops"
	ccipops "github.com/smartcontractkit/chainlink-sui/deployment/ops/ccip"
	burnminttokenpoolops "github.com/smartcontractkit/chainlink-sui/deployment/ops/ccip_burn_mint_token_pool"
	lockreleasetokenpoolops "github.com/smartcontractkit/chainlink-sui/deployment/ops/ccip_lock_release_token_pool"
	managedtokenpoolops "github.com/smartcontractkit/chainlink-sui/deployment/ops/ccip_managed_token_pool"
	managedtokenops "github.com/smartcontractkit/chainlink-sui/deployment/ops/managed_token"
	suiofframp_helper "github.com/smartcontractkit/chainlink-sui/relayer/chainwriter/ptb/offramp"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	suideps "github.com/smartcontractkit/chainlink/deployment/ccip/changeset/sui"
	ccipclient "github.com/smartcontractkit/chainlink/deployment/ccip/shared/client"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
	commoncs "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/ccipevm"
)

const TokenSymbolLINK = "LINK"

type SuiSendRequest struct {
	Receiver         []byte
	Data             []byte
	ExtraArgs        []byte
	FeeToken         string
	FeeTokenStore    string
	TokenAmounts     []SuiTokenAmount
	TokenReceiverATA []byte
}

type SuiTokenAmount struct {
	TokenPoolType sui_deployment.TokenPoolType
	Token         string
	Amount        uint64
}

type RampMessageHeader struct {
	MessageID           []byte `json:"message_id"`
	SourceChainSelector string `json:"source_chain_selector"`
	DestChainSelector   string `json:"dest_chain_selector"`
	SequenceNumber      string `json:"sequence_number"`
	Nonce               string `json:"nonce"`
}

type Sui2AnyRampMessage struct {
	Header         RampMessageHeader `json:"header"`
	Sender         string            `json:"sender"`
	Data           []byte            `json:"data"`
	Receiver       []byte            `json:"receiver"`
	ExtraArgs      []byte            `json:"extra_args"`
	FeeToken       string            `json:"fee_token"`
	FeeTokenAmount string            `json:"fee_token_amount"`
	FeeValueJuels  string            `json:"fee_value_juels"`
}

type CCIPMessageSent struct {
	DestChainSelector string             `json:"dest_chain_selector"`
	SequenceNumber    string             `json:"sequence_number"`
	Message           Sui2AnyRampMessage `json:"message"`
}

type TokenPoolRateLimiterConfig struct {
	RemoteChainSelector uint64
	OutboundIsEnabled   bool
	OutboundCapacity    uint64
	OutboundRate        uint64
	InboundIsEnabled    bool
	InboundCapacity     uint64
	InboundRate         uint64
}

func SendSuiCCIPRequest(e cldf.Environment, cfg *ccipclient.CCIPSendReqConfig) (*ccipclient.AnyMsgSentEvent, error) {
	ctx := e.GetContext()
	state, err := stateview.LoadOnchainState(e)
	if err != nil {
		return &ccipclient.AnyMsgSentEvent{}, err
	}

	if err != nil {
		return &ccipclient.AnyMsgSentEvent{}, err
	}

	suiChains := e.BlockChains.SuiChains()
	suiChain := suiChains[cfg.SourceChain]

	deps := suideps.Deps{
		SuiChain: sui_ops.OpTxDeps{
			Client: suiChain.Client,
			Signer: suiChain.Signer,
			GetCallOpts: func() *suiBind.CallOpts {
				b := uint64(400_000_000)
				return &suiBind.CallOpts{
					Signer:           suiChain.Signer,
					WaitForExecution: true,
					GasBudget:        &b,
				}
			},
		},
	}

	ccipObjectRefID := state.SuiChains[cfg.SourceChain].CCIPObjectRef
	ccipPackageID := state.SuiChains[cfg.SourceChain].CCIPMockV2PackageId
	if ccipPackageID == "" {
		fmt.Println("ccip v2 not set, using ccip v1")
		ccipPackageID = state.SuiChains[cfg.SourceChain].CCIPAddress
	}
	onRampPackageID := state.SuiChains[cfg.SourceChain].OnRampMockV2PackageId
	if onRampPackageID == "" {
		fmt.Println("onRamp v2 not set, using onramp v1")
		onRampPackageID = state.SuiChains[cfg.SourceChain].OnRampAddress
	}
	onRampStateObjectID := state.SuiChains[cfg.SourceChain].OnRampStateObjectId
	linkTokenPkgID := state.SuiChains[cfg.SourceChain].LinkTokenAddress
	linkTokenObjectMetadataID := state.SuiChains[cfg.SourceChain].LinkTokenCoinMetadataId
	ccipOwnerCapID := state.SuiChains[cfg.SourceChain].CCIPOwnerCapObjectId

	bigIntSourceUsdPerToken, parsed := new(big.Int).SetString("15377040000000000000000000000", 10) // 1e27 since sui is 1e9
	if !parsed {
		return &ccipclient.AnyMsgSentEvent{}, errors.New("failed converting SourceUSDPerToken to bigInt")
	}

	bigIntGasUsdPerUnitGas, ok := new(big.Int).SetString("41946474500", 10) // optimism sep 4145822215
	if !ok {
		return &ccipclient.AnyMsgSentEvent{}, errors.New("failed converting GasUsdPerUnitGas to bigInt")
	}

	// getValidatedFee
	msg := cfg.Message.(SuiSendRequest)

	// Update Prices on FeeQuoter with minted LinkToken
	_, err = operations.ExecuteOperation(e.OperationsBundle, ccipops.FeeQuoterUpdatePricesWithOwnerCapOp, deps.SuiChain,
		ccipops.FeeQuoterUpdatePricesWithOwnerCapInput{
			CCIPPackageId:         ccipPackageID,
			CCIPObjectRef:         ccipObjectRefID,
			OwnerCapObjectId:      ccipOwnerCapID,
			SourceTokens:          []string{linkTokenObjectMetadataID},
			SourceUsdPerToken:     []*big.Int{bigIntSourceUsdPerToken},
			GasDestChainSelectors: []uint64{cfg.DestChain},
			GasUsdPerUnitGas:      []*big.Int{bigIntGasUsdPerUnitGas},
		})
	if err != nil {
		return &ccipclient.AnyMsgSentEvent{}, errors.New("failed to updatePrice for Sui chain " + err.Error())
	}

	// TODO: might be needed for validation
	// feeQuoter, err := module_fee_quoter.NewFeeQuoter(ccipPackageID, deps.SuiChain.Client)
	// if err != nil {
	// 	return &ccipclient.AnyMsgSentEvent{}, err
	// }

	// validatedFee, err := feeQuoter.DevInspect().GetValidatedFee(ctx, &suiBind.CallOpts{
	// 	Signer:           deps.SuiChain.Signer,
	// 	WaitForExecution: true,
	// },
	// 	suiBind.Object{Id: ccipObjectRefID},
	// 	suiBind.Object{Id: "0x6"},
	// 	cfg.DestChain,
	// 	[]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	// 		0x00, 0x00, 0x00, 0x00, 0xdd, 0xbb, 0x6f, 0x35,
	// 		0x8f, 0x29, 0x04, 0x08, 0xd7, 0x68, 0x47, 0xb4,
	// 		0xf6, 0x02, 0xf0, 0xfd, 0x59, 0x92, 0x95, 0xfd,
	// 	},
	// 	[]byte("hello evm from sui"),
	// 	[]string{},
	// 	[]uint64{},
	// 	linkTokenObjectMetadataID,
	// 	[]byte{},
	// )
	// if err != nil {
	// 	return &ccipclient.AnyMsgSentEvent{}, err
	// }

	// fmt.Println("VALIDATED FEE:", validatedFee)

	if len(msg.TokenAmounts) > 0 {
		var tokenPoolState sui_deployment.CCIPPoolState
		var tokenPoolPkgName string
		var tokenPoolModuleName string

		switch msg.TokenAmounts[0].TokenPoolType {
		case sui_deployment.TokenPoolTypeBurnMint:
			bnmTokenPool, exists := state.SuiChains[cfg.SourceChain].BnMTokenPools[TokenSymbolLINK]
			if !exists {
				return nil, fmt.Errorf("no BurnMintTokenPool found for token: %s", TokenSymbolLINK)
			}
			tokenPoolState = bnmTokenPool
			tokenPoolPkgName = "burn_mint_token_pool"
			tokenPoolModuleName = "burn_mint_token_pool"
		case sui_deployment.TokenPoolTypeManaged:
			managedTokenPool, exists := state.SuiChains[cfg.SourceChain].ManagedTokenPools[TokenSymbolLINK]
			if !exists {
				return nil, fmt.Errorf("no ManagedTokenPool found for token: %s", TokenSymbolLINK)
			}
			tokenPoolState = managedTokenPool
			tokenPoolPkgName = "managed_token_pool"
			tokenPoolModuleName = "managed_token_pool"
		case sui_deployment.TokenPoolTypeLockRelease:
			lnrTokenPool, exists := state.SuiChains[cfg.SourceChain].LnRTokenPools[TokenSymbolLINK]
			if !exists {
				return nil, fmt.Errorf("no LockReleaseTokenPool found for token: %s", TokenSymbolLINK)
			}
			tokenPoolState = lnrTokenPool
			tokenPoolPkgName = "lock_release_token_pool"
			tokenPoolModuleName = "lock_release_token_pool"
		default:
			return nil, fmt.Errorf("unsupported token pool type: %s", msg.TokenAmounts[0].TokenPoolType)
		}

		tokenPoolPkgID := tokenPoolState.PackageID
		tokenPoolStateObjectID := tokenPoolState.StateObjectId

		// 3 ptb calls
		// 1. create_token_transfer_params
		// 2. lock_or_burn
		// 3. ccip_send

		// 1. create_token_transfer_params
		client := sui.NewSuiClient(suiChain.URL)
		ptb := suitx.NewTransaction()
		ptb.SetSuiClient(client.(*sui.Client))

		// Bind contracts
		ccipStateHelperContract, err := suiBind.NewBoundContract(
			ccipPackageID,
			"ccip",
			"onramp_state_helper",
			client,
		)
		if err != nil {
			return nil, errors.New("failed to create onramp state helper bound contract when appending PTB command: " + err.Error())
		}

		tokenPoolContract, err := suiBind.NewBoundContract(
			tokenPoolPkgID,
			tokenPoolPkgName,
			tokenPoolModuleName,
			client,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to create %s bound contract when appending PTB command: %w", tokenPoolPkgName, err)
		}

		onRampContract, err := suiBind.NewBoundContract(
			onRampPackageID,
			"ccip_onramp",
			"onramp",
			client,
		)
		if err != nil {
			return nil, errors.New("failed to create ccip_onramp contract when appending PTB command: " + err.Error())
		}

		/*********  1. create_token_transfer_params *******/
		typeArgsList := []string{}
		typeParamsList := []string{}
		paramTypes := []string{
			"vector<u8>",
		}

		var paramValues []any
		destFamily, err := chainsel.GetSelectorFamily(cfg.DestChain)
		if err != nil {
			return nil, errors.New("failed to get selector family for destination chain: " + err.Error())
		}
		switch destFamily {
		case chainsel.FamilyEVM, chainsel.FamilyAptos:
			paramValues = []any{
				msg.Receiver,
			}
		case chainsel.FamilySui, chainsel.FamilySolana:
			paramValues = []any{
				msg.TokenReceiverATA,
			}
		default:
			return nil, errors.New("unsupported destination chain family: " + destFamily)
		}

		onRampCreateTokenTransferParamsCall, err := ccipStateHelperContract.EncodeCallArgsWithGenerics(
			"create_token_transfer_params",
			typeArgsList,
			typeParamsList,
			paramTypes,
			paramValues,
			nil,
		)
		if err != nil {
			return nil, errors.New("failed to encode onRampCreateTokenTransferParamsCall call: " + err.Error())
		}

		createTokenTransferParamsResult, err := ccipStateHelperContract.AppendPTB(ctx, deps.SuiChain.GetCallOpts(), ptb, onRampCreateTokenTransferParamsCall)
		if err != nil {
			return nil, errors.New("failed to build PTB (get_token_param_data) using bindings: " + err.Error())
		}

		/*********  2. lock_or_burn *******/
		normalizedModuleTP, err := client.SuiGetNormalizedMoveModule(ctx, models.GetNormalizedMoveModuleRequest{
			Package:    tokenPoolPkgID,
			ModuleName: tokenPoolModuleName,
		})
		if err != nil {
			return nil, errors.New("failed to get normalized module: " + err.Error())
		}

		functionSignatureLnB, isValidLockOrBurn := normalizedModuleTP.ExposedFunctions["lock_or_burn"]
		if !isValidLockOrBurn {
			return nil, errors.New("missing function signature for receiver function not found in module lock_or_burn")
		}

		// Figure out the parameter types from the normalized module of the token pool
		paramTypesLockBurn, err := suiofframp_helper.DecodeParameters(e.Logger, functionSignatureLnB.(map[string]any), "parameters")
		if err != nil {
			return nil, errors.New("failed to decode parameters for token pool function: " + err.Error())
		}

		typeArgsListLinkTokenPkgID := []string{linkTokenPkgID + "::link::LINK"}
		typeParamsList = []string{}

		var paramValuesLockBurn []any
		switch msg.TokenAmounts[0].TokenPoolType {
		case sui_deployment.TokenPoolTypeBurnMint:
			paramValuesLockBurn = []any{
				suiBind.Object{Id: ccipObjectRefID},           // ref
				createTokenTransferParamsResult,               // token_params
				suiBind.Object{Id: msg.TokenAmounts[0].Token}, // minted token to send to EVM
				cfg.DestChain,
				suiBind.Object{Id: "0x6"},                  // clock
				suiBind.Object{Id: tokenPoolStateObjectID}, // BM TP state object id
			}
		case sui_deployment.TokenPoolTypeManaged:
			paramValuesLockBurn = []any{
				suiBind.Object{Id: ccipObjectRefID},           // ref
				createTokenTransferParamsResult,               // token_params
				suiBind.Object{Id: msg.TokenAmounts[0].Token}, // minted token to send to EVM
				cfg.DestChain,
				suiBind.Object{Id: "0x6"},   // clock
				suiBind.Object{Id: "0x403"}, // deny list
				suiBind.Object{Id: state.SuiChains[cfg.SourceChain].ManagedTokens[TokenSymbolLINK].StateObjectId}, // Managed token state object id
				suiBind.Object{Id: tokenPoolStateObjectID},                                                        // Managed TP state object id
			}
		case sui_deployment.TokenPoolTypeLockRelease:
			paramValuesLockBurn = []any{
				suiBind.Object{Id: ccipObjectRefID},           // ref
				createTokenTransferParamsResult,               // token_params
				suiBind.Object{Id: msg.TokenAmounts[0].Token}, // locked token to send to EVM
				cfg.DestChain,
				suiBind.Object{Id: "0x6"},                  // clock
				suiBind.Object{Id: tokenPoolStateObjectID}, // LnR TP state object id
			}
		default:
			return nil, fmt.Errorf("unsupported token pool type: %s", msg.TokenAmounts[0].TokenPoolType)
		}

		lockOrBurnParamsCall, err := tokenPoolContract.EncodeCallArgsWithGenerics(
			"lock_or_burn",
			typeArgsListLinkTokenPkgID,
			typeParamsList,
			paramTypesLockBurn,
			paramValuesLockBurn,
			nil,
		)
		if err != nil {
			return nil, errors.New("failed to encode lockOrBurnParamsCall call: " + err.Error())
		}

		_, err = tokenPoolContract.AppendPTB(ctx, deps.SuiChain.GetCallOpts(), ptb, lockOrBurnParamsCall)
		if err != nil {
			return nil, errors.New("failed to build PTB (get_token_param_data) using bindings: " + err.Error())
		}

		/********* 3. ccip_send *******/
		normalizedModule, err := client.SuiGetNormalizedMoveModule(ctx, models.GetNormalizedMoveModuleRequest{
			Package:    onRampPackageID,
			ModuleName: "onramp",
		})
		if err != nil {
			return nil, errors.New("failed to get normalized module: " + err.Error())
		}

		functionSignature, parsedccipSend := normalizedModule.ExposedFunctions["ccip_send"]
		if !parsedccipSend {
			return nil, errors.New("missing function signature for receiver function not found in module ccip_send")
		}

		// Figure out the parameter types from the normalized module of the token pool
		paramTypesCCIPSend, err := suiofframp_helper.DecodeParameters(e.Logger, functionSignature.(map[string]any), "parameters")
		if err != nil {
			return nil, errors.New("failed to decode parameters for token pool function: " + err.Error())
		}

		paramValuesCCIPSend := []any{
			suiBind.Object{Id: ccipObjectRefID},
			suiBind.Object{Id: onRampStateObjectID},
			suiBind.Object{Id: "0x6"},
			cfg.DestChain,
			msg.Receiver, // receiver
			msg.Data,
			createTokenTransferParamsResult,               // tokenParams from the original create_token_transfer_params
			suiBind.Object{Id: linkTokenObjectMetadataID}, // feeTokenMetadata
			suiBind.Object{Id: msg.FeeToken},
			msg.ExtraArgs, // extraArgs
		}

		encodedOnRampCCIPSendCall, err := onRampContract.EncodeCallArgsWithGenerics(
			"ccip_send",
			typeArgsListLinkTokenPkgID,
			[]string{},
			paramTypesCCIPSend,
			paramValuesCCIPSend,
			nil,
		)
		if err != nil {
			return nil, errors.New("failed to encode calls for ccip_send: " + err.Error())
		}

		_, err = onRampContract.AppendPTB(ctx, deps.SuiChain.GetCallOpts(), ptb, encodedOnRampCCIPSendCall)
		if err != nil {
			return nil, errors.New("failed to build PTB (receiver call) using bindings: " + err.Error())
		}

		executeCCIPSend, err := suiBind.ExecutePTB(ctx, deps.SuiChain.GetCallOpts(), client, ptb)
		if err != nil {
			return nil, errors.New("failed to execute ccip_send with err: " + err.Error())
		}

		var suiEventResp models.SuiEventResponse
		var found bool
		for _, event := range executeCCIPSend.Events {
			// find the CCIPMessageSent event emitted by the onramp package
			if event.PackageId == onRampPackageID && strings.HasSuffix(event.Type, "CCIPMessageSent") {
				suiEventResp = event
				found = true
				break
			}
		}
		if !found {
			return nil, errors.New("no CCIPMessageSent event found")
		}
		suiEvent := suiEventResp.ParsedJson

		seqStr, _ := suiEvent["sequence_number"].(string)
		seq, _ := strconv.ParseUint(seqStr, 10, 64)

		return &ccipclient.AnyMsgSentEvent{
			SequenceNumber: seq,
			RawEvent:       suiEvent, // just dump raw
		}, nil
	}

	// TODO: SUI CCIPSend using bindings
	client := sui.NewSuiClient(suiChain.URL)
	ptb := suitx.NewTransaction()
	ptb.SetSuiClient(client.(*sui.Client))

	// ptb1
	ccipStateHelperContract, err := suiBind.NewBoundContract(
		ccipPackageID,
		"ccip",
		"onramp_state_helper",
		client,
	)
	if err != nil {
		return nil, errors.New("failed to create onramp state helper bound contract when appending PTB command: " + err.Error())
	}

	// Note: these will be different for token transfers
	typeArgsList := []string{}
	typeParamsList := []string{}
	paramTypes := []string{
		"vector<u8>",
	}

	// For SUI -> EVM BurnMint Pool token Transfer, we can use msg.Receiver as tokenReceiver, this field is only used in usdc token pool
	// bc we need to check the recipient with Circle's packages from the onramp side before sending USDC. and it's not used anyway else.
	decodedTokenReceiver, _ := hex.DecodeString("0000000000000000000000000000000000000000000000000000000000000000")
	var tokenReceiver [32]byte
	copy(tokenReceiver[:], decodedTokenReceiver)

	paramValues := []any{
		decodedTokenReceiver,
	}

	onRampCreateTokenTransferParamsCall, err := ccipStateHelperContract.EncodeCallArgsWithGenerics(
		"create_token_transfer_params",
		typeArgsList,
		typeParamsList,
		paramTypes,
		paramValues,
		nil,
	)
	if err != nil {
		return nil, errors.New("failed to encode onRampCreateTokenTransferParamsCall call: " + err.Error())
	}

	extractedAny2SuiMessageResult, err := ccipStateHelperContract.AppendPTB(ctx, deps.SuiChain.GetCallOpts(), ptb, onRampCreateTokenTransferParamsCall)
	if err != nil {
		return nil, errors.New("failed to build PTB (get_token_param_data) using bindings: " + err.Error())
	}

	// ptb2
	onRampContract, err := suiBind.NewBoundContract(
		onRampPackageID,
		"ccip_onramp",
		"onramp",
		client,
	)
	if err != nil {
		return nil, errors.New("failed to create onramp bound contract when appending PTB command: " + err.Error())
	}

	// normalize module
	normalizedModule, err := client.SuiGetNormalizedMoveModule(ctx, models.GetNormalizedMoveModuleRequest{
		Package:    onRampPackageID,
		ModuleName: "onramp",
	})
	if err != nil {
		return nil, errors.New("failed to get normalized module: " + err.Error())
	}

	functionSignature, ok := normalizedModule.ExposedFunctions["ccip_send"]
	if !ok {
		return nil, errors.New("missing function signature for receiver function not found in module ccip_send")
	}

	// Figure out the parameter types from the normalized module of the token pool
	paramTypes, err = suiofframp_helper.DecodeParameters(e.Logger, functionSignature.(map[string]any), "parameters")
	if err != nil {
		return nil, errors.New("failed to decode parameters for token pool function: " + err.Error())
	}

	typeArgsList = []string{linkTokenPkgID + "::link::LINK"}
	typeParamsList = []string{}
	paramValues = []any{
		suiBind.Object{Id: ccipObjectRefID},
		suiBind.Object{Id: onRampStateObjectID},
		suiBind.Object{Id: "0x6"},
		cfg.DestChain,
		msg.Receiver, // receiver (TODO: replace this with sender Address use environment.NormalizeTo32Bytes(ethereumAddress) from sui repo)
		msg.Data,
		extractedAny2SuiMessageResult,                 // tokenParams
		suiBind.Object{Id: linkTokenObjectMetadataID}, // feeTokenMetadata
		suiBind.Object{Id: msg.FeeToken},
		msg.ExtraArgs, // extraArgs
	}

	encodedOnRampCCIPSendCall, err := onRampContract.EncodeCallArgsWithGenerics(
		"ccip_send",
		typeArgsList,
		typeParamsList,
		paramTypes,
		paramValues,
		nil,
	)
	if err != nil {
		return nil, errors.New("failed to encode receiver call: " + err.Error())
	}

	_, err = onRampContract.AppendPTB(ctx, deps.SuiChain.GetCallOpts(), ptb, encodedOnRampCCIPSendCall)
	if err != nil {
		return nil, errors.New("failed to build PTB (receiver call) using bindings: " + err.Error())
	}

	executeCCIPSend, err := suiBind.ExecutePTB(ctx, deps.SuiChain.GetCallOpts(), client, ptb)
	if err != nil {
		return nil, errors.New("failed to execute ccip_send with err: " + err.Error())
	}

	var suiEventResp models.SuiEventResponse
	var found bool
	for _, event := range executeCCIPSend.Events {
		// find the CCIPMessageSent event emitted by the onramp package
		if event.PackageId == onRampPackageID && strings.HasSuffix(event.Type, "CCIPMessageSent") {
			suiEventResp = event
			found = true
			break
		}
	}
	if !found {
		return nil, errors.New("no CCIPMessageSent event found")
	}
	suiEvent := suiEventResp.ParsedJson

	seqStr, ok := suiEvent["sequence_number"].(string)
	if !ok {
		return nil, fmt.Errorf("failed to extract sequence_number from Sui event: %+v", suiEvent)
	}
	seq, err := strconv.ParseUint(seqStr, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse sequence number '%s': %w", seqStr, err)
	}

	return &ccipclient.AnyMsgSentEvent{
		SequenceNumber: seq,
		RawEvent:       suiEvent, // just dump raw
	}, nil
}

func MakeSuiExtraArgs(gasLimit uint64, allowOOO bool, receiverObjectIDs [][32]byte, tokenReceiver [32]byte) []byte {
	extraArgs, err := ccipevm.SerializeClientSUIExtraArgsV1(message_hasher.ClientSuiExtraArgsV1{
		GasLimit:                 new(big.Int).SetUint64(gasLimit),
		AllowOutOfOrderExecution: allowOOO,
		TokenReceiver:            tokenReceiver,
		ReceiverObjectIds:        receiverObjectIDs,
	})
	if err != nil {
		panic(err)
	}
	return extraArgs
}

// HandleTokenAndBurnMintTokenPoolDeploymentForSUI deploys a transferrable token and a burn mint token pool on the EVM chain.
// It also deploys a burn mint token pool on the SUI chain and configures it to work with the transferrable token on the EVM chain.
func HandleTokenAndBurnMintTokenPoolDeploymentForSUI(e cldf.Environment, suiChainSel, evmChainSel uint64, rateLimiterConfigs []TokenPoolRateLimiterConfig) (cldf.Environment, *burn_mint_erc677.BurnMintERC677, *burn_mint_token_pool.BurnMintTokenPool, error) {
	suiChains := e.BlockChains.SuiChains()
	suiChain := suiChains[suiChainSel]

	evmChain := e.BlockChains.EVMChains()[evmChainSel]

	evmDeployerKey := evmChain.DeployerKey
	state, err := stateview.LoadOnchainState(e)
	if err != nil {
		return cldf.Environment{}, nil, nil, errors.New("failed load onstate chains " + err.Error())
	}

	linkTokenPkgID := state.SuiChains[suiChainSel].LinkTokenAddress
	linkTokenObjectMetadataID := state.SuiChains[suiChainSel].LinkTokenCoinMetadataId
	linkTokenTreasuryCapID := state.SuiChains[suiChainSel].LinkTokenTreasuryCapId

	// Deploy transferrable token on EVM
	evmToken, evmPool, err := deployTransferTokenOneEnd(e.Logger, evmChain, evmDeployerKey, e.ExistingAddresses, "TOKEN")
	if err != nil {
		return cldf.Environment{}, nil, nil, errors.New("failed to deploy transfer token for evm chain " + err.Error())
	}

	err = attachTokenToTheRegistry(evmChain, state.MustGetEVMChainState(evmChain.Selector), evmDeployerKey, evmToken.Address(), evmPool.Address())
	if err != nil {
		return cldf.Environment{}, nil, nil, errors.New("failed to attach token to registry for evm " + err.Error())
	}

	// Deploy & Configure BurnMint TP on SUI
	e, _, err = commoncs.ApplyChangesets(&testing.T{}, e, []commoncs.ConfiguredChangeSet{
		commoncs.Configure(sui_cs.DeployTPAndConfigure{}, sui_cs.DeployTPAndConfigureConfig{
			SuiChainSelector: suiChainSel,
			TokenPoolTypes:   []sui_deployment.TokenPoolType{sui_deployment.TokenPoolTypeBurnMint},
			BurnMintTpInput: burnminttokenpoolops.DeployAndInitBurnMintTokenPoolInput{
				CoinObjectTypeArg:    linkTokenPkgID + "::link::LINK",
				CoinMetadataObjectId: linkTokenObjectMetadataID,
				TreasuryCapObjectId:  linkTokenTreasuryCapID,

				// apply dest chain updates
				RemoteChainSelectorsToRemove: []uint64{},
				RemoteChainSelectorsToAdd:    []uint64{evmChainSel},
				RemotePoolAddressesToAdd:     [][]string{{evmPool.Address().String()}}, // this gets convert to 32byte bytes internally
				RemoteTokenAddressesToAdd: []string{
					evmToken.Address().String(), // this gets convert to 32byte bytes internally
				},

				// set chain rate limiter configs
				RemoteChainSelectors: extractFields(rateLimiterConfigs, func(c TokenPoolRateLimiterConfig) uint64 { return c.RemoteChainSelector }),
				OutboundIsEnableds:   extractFields(rateLimiterConfigs, func(c TokenPoolRateLimiterConfig) bool { return c.OutboundIsEnabled }),
				OutboundCapacities:   extractFields(rateLimiterConfigs, func(c TokenPoolRateLimiterConfig) uint64 { return c.OutboundCapacity }),
				OutboundRates:        extractFields(rateLimiterConfigs, func(c TokenPoolRateLimiterConfig) uint64 { return c.OutboundRate }),
				InboundIsEnableds:    extractFields(rateLimiterConfigs, func(c TokenPoolRateLimiterConfig) bool { return c.InboundIsEnabled }),
				InboundCapacities:    extractFields(rateLimiterConfigs, func(c TokenPoolRateLimiterConfig) uint64 { return c.InboundCapacity }),
				InboundRates:         extractFields(rateLimiterConfigs, func(c TokenPoolRateLimiterConfig) uint64 { return c.InboundRate }),
			},
		}),
	})
	if err != nil {
		return cldf.Environment{}, nil, nil, err
	}

	// reload onChainState to get deployed TP contracts
	state, err = stateview.LoadOnchainState(e)
	if err != nil {
		return cldf.Environment{}, nil, nil, errors.New("failed load onstate chains " + err.Error())
	}

	// TODO: might be needed for validation
	// ensure tokenPool is added
	// (ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, remoteChainSelector uint64)
	// bmtp, err := sui_module_bnmtp.NewBurnMintTokenPool(state.SuiChains[suiChainSel].CCIPBurnMintTokenPool, e.BlockChains.SuiChains()[suiChainSel].Client)
	// if err != nil {
	// 	return cldf.Environment{}, nil, nil, err
	// }

	// val, err := bmtp.DevInspect().GetRemotePools(context.Background(), &suiBind.CallOpts{
	// 	Signer:           e.BlockChains.SuiChains()[suiChainSel].Signer,
	// 	WaitForExecution: true,
	// }, []string{linkTokenPkgID + "::link::LINK"}, suiBind.Object{Id: state.SuiChains[suiChainSel].CCIPBurnMintTokenPoolState}, evmChainSel)
	// if err != nil {
	// 	return cldf.Environment{}, nil, nil, err
	// }

	// val1, err := bmtp.DevInspect().IsRemotePool(context.Background(), &suiBind.CallOpts{
	// 	Signer:           e.BlockChains.SuiChains()[suiChainSel].Signer,
	// 	WaitForExecution: true,
	// }, []string{linkTokenPkgID + "::link::LINK"}, suiBind.Object{Id: state.SuiChains[suiChainSel].CCIPBurnMintTokenPoolState}, evmChainSel, evmPool.Address().Bytes())
	// if err != nil {
	// 	return cldf.Environment{}, nil, nil, err
	// }

	suiTokenBytes, err := hex.DecodeString(strings.TrimPrefix(linkTokenObjectMetadataID, "0x"))
	if err != nil {
		return cldf.Environment{}, nil, nil, errors.New("error while decoding suiToken")
	}

	bnmTokenPool, ok := state.SuiChains[suiChainSel].BnMTokenPools[TokenSymbolLINK]
	if !ok {
		return cldf.Environment{}, nil, nil, fmt.Errorf("no BurnMintTokenPool found for token: %s", TokenSymbolLINK)
	}

	suiPoolBytes, err := hex.DecodeString(strings.TrimPrefix(bnmTokenPool.PackageID, "0x"))
	if err != nil {
		return cldf.Environment{}, nil, nil, errors.New("error while decoding suiPool")
	}

	err = setTokenPoolCounterPart(e.BlockChains.EVMChains()[evmChain.Selector], evmPool, evmDeployerKey, suiChain.Selector, suiTokenBytes, suiPoolBytes)
	if err != nil {
		return cldf.Environment{}, nil, nil, errors.New("failed to add token to the counterparty " + err.Error())
	}

	err = grantMintBurnPermissions(e.Logger, e.BlockChains.EVMChains()[evmChain.Selector], evmToken, evmDeployerKey, evmPool.Address())
	if err != nil {
		return cldf.Environment{}, nil, nil, errors.New("failed to grant burnMint " + err.Error())
	}

	return e, evmToken, evmPool, nil
}

// HandleTokenAndManagedTokenPoolDeploymentForSUI deploys a transferrable token and a burn mint token pool on the EVM chain.
// It also deploys a managed token pool on the SUI chain and configures it to work with the transferrable token on the EVM chain.
func HandleTokenAndManagedTokenPoolDeploymentForSUI(e cldf.Environment, suiChainSel, evmChainSel uint64, rateLimiterConfigs []TokenPoolRateLimiterConfig) (cldf.Environment, *burn_mint_erc677.BurnMintERC677, *burn_mint_token_pool.BurnMintTokenPool, error) {
	evmChain := e.BlockChains.EVMChains()[evmChainSel]
	suiChain := e.BlockChains.SuiChains()[suiChainSel]
	deployerAddr, err := suiChain.Signer.GetAddress()
	if err != nil {
		return cldf.Environment{}, nil, nil, errors.New("failed to get deployer address " + err.Error())
	}

	// Deploy Transferrable TOKEN on ETH
	// EVM
	evmDeployerKey := evmChain.DeployerKey
	state, err := stateview.LoadOnchainState(e)
	if err != nil {
		return cldf.Environment{}, nil, nil, errors.New("failed load onstate chains " + err.Error())
	}

	linkTokenPkgID := state.SuiChains[suiChainSel].LinkTokenAddress
	linkTokenObjectMetadataID := state.SuiChains[suiChainSel].LinkTokenCoinMetadataId
	linkTokenTreasuryCapID := state.SuiChains[suiChainSel].LinkTokenTreasuryCapId

	// Deploy & Configure Managed Token on SUI
	e, _, err = commoncs.ApplyChangesets(&testing.T{}, e, []commoncs.ConfiguredChangeSet{
		commoncs.Configure(sui_cs.DeployManagedToken{}, sui_cs.DeployManagedTokenConfig{
			DeployAndInitManagedTokenInput: managedtokenops.DeployAndInitManagedTokenInput{
				CoinObjectTypeArg:   linkTokenPkgID + "::link::LINK",
				TreasuryCapObjectId: linkTokenTreasuryCapID,
				MinterAddress:       deployerAddr,
				Allowance:           0,
				IsUnlimited:         true,
			},
			ChainSelector: suiChainSel,
		}),
	})

	if err != nil {
		return cldf.Environment{}, nil, nil, err
	}

	// Deploy transferrable token on EVM
	evmToken, evmPool, err := deployTransferTokenOneEnd(e.Logger, evmChain, evmDeployerKey, e.ExistingAddresses, "TOKEN")
	if err != nil {
		return cldf.Environment{}, nil, nil, errors.New("failed to deploy transfer token for evm chain " + err.Error())
	}

	err = attachTokenToTheRegistry(evmChain, state.MustGetEVMChainState(evmChain.Selector), evmDeployerKey, evmToken.Address(), evmPool.Address())
	if err != nil {
		return cldf.Environment{}, nil, nil, errors.New("failed to attach token to registry for evm " + err.Error())
	}

	// reload onChainState to get deployed Managed Token contracts
	state, err = stateview.LoadOnchainState(e)
	if err != nil {
		return cldf.Environment{}, nil, nil, errors.New("failed load onstate chains " + err.Error())
	}

	e, _, err = commoncs.ApplyChangesets(&testing.T{}, e, []commoncs.ConfiguredChangeSet{
		commoncs.Configure(sui_cs.DeployTPAndConfigure{}, sui_cs.DeployTPAndConfigureConfig{
			SuiChainSelector: suiChainSel,
			TokenPoolTypes:   []sui_deployment.TokenPoolType{sui_deployment.TokenPoolTypeManaged},
			ManagedTPInput: managedtokenpoolops.DeployAndInitManagedTokenPoolInput{
				CoinObjectTypeArg:         linkTokenPkgID + "::link::LINK",
				CoinMetadataObjectId:      linkTokenObjectMetadataID,
				MintCapObjectId:           state.SuiChains[suiChainSel].ManagedTokens[TokenSymbolLINK].MinterCapObjectIds[0],
				ManagedTokenStateObjectId: state.SuiChains[suiChainSel].ManagedTokens[TokenSymbolLINK].StateObjectId,
				ManagedTokenOwnerCapId:    state.SuiChains[suiChainSel].ManagedTokens[TokenSymbolLINK].OwnerCapObjectId,
				// apply dest chain updates
				RemoteChainSelectorsToRemove: []uint64{},
				RemoteChainSelectorsToAdd:    []uint64{evmChainSel},
				RemotePoolAddressesToAdd:     [][]string{{evmPool.Address().String()}}, // this gets convert to 32byte bytes internally
				RemoteTokenAddressesToAdd: []string{
					evmToken.Address().String(), // this gets convert to 32byte bytes internally
				},

				// set chain rate limiter configs
				RemoteChainSelectors: extractFields(rateLimiterConfigs, func(c TokenPoolRateLimiterConfig) uint64 { return c.RemoteChainSelector }),
				OutboundIsEnableds:   extractFields(rateLimiterConfigs, func(c TokenPoolRateLimiterConfig) bool { return c.OutboundIsEnabled }),
				OutboundCapacities:   extractFields(rateLimiterConfigs, func(c TokenPoolRateLimiterConfig) uint64 { return c.OutboundCapacity }),
				OutboundRates:        extractFields(rateLimiterConfigs, func(c TokenPoolRateLimiterConfig) uint64 { return c.OutboundRate }),
				InboundIsEnableds:    extractFields(rateLimiterConfigs, func(c TokenPoolRateLimiterConfig) bool { return c.InboundIsEnabled }),
				InboundCapacities:    extractFields(rateLimiterConfigs, func(c TokenPoolRateLimiterConfig) uint64 { return c.InboundCapacity }),
				InboundRates:         extractFields(rateLimiterConfigs, func(c TokenPoolRateLimiterConfig) uint64 { return c.InboundRate }),
			},
		}),
	})
	if err != nil {
		return cldf.Environment{}, nil, nil, err
	}

	// reload onChainState to get deployed managed token pool contracts
	state, err = stateview.LoadOnchainState(e)
	if err != nil {
		return cldf.Environment{}, nil, nil, errors.New("failed load onstate chains " + err.Error())
	}

	suiTokenBytes, err := hex.DecodeString(strings.TrimPrefix(linkTokenObjectMetadataID, "0x"))
	if err != nil {
		return cldf.Environment{}, nil, nil, errors.New("error while decoding suiToken")
	}

	managedTokenPool, ok := state.SuiChains[suiChainSel].ManagedTokenPools[TokenSymbolLINK]
	if !ok {
		return cldf.Environment{}, nil, nil, fmt.Errorf("no ManagedTokenPool found for token: %s", TokenSymbolLINK)
	}

	suiPoolBytes, err := hex.DecodeString(strings.TrimPrefix(managedTokenPool.PackageID, "0x"))
	if err != nil {
		return cldf.Environment{}, nil, nil, errors.New("error while decoding suiPool")
	}

	err = setTokenPoolCounterPart(e.BlockChains.EVMChains()[evmChain.Selector], evmPool, evmDeployerKey, suiChain.Selector, suiTokenBytes, suiPoolBytes)
	if err != nil {
		return cldf.Environment{}, nil, nil, errors.New("failed to add token to the counterparty " + err.Error())
	}

	err = grantMintBurnPermissions(e.Logger, e.BlockChains.EVMChains()[evmChain.Selector], evmToken, evmDeployerKey, evmPool.Address())
	if err != nil {
		return cldf.Environment{}, nil, nil, errors.New("failed to grant burnMint " + err.Error())
	}

	return e, evmToken, evmPool, nil
}

func HandleTokenAndLockReleaseTokenPoolDeploymentForSUI(e cldf.Environment, suiChainSel, evmChainSel uint64, rateLimiterConfigs []TokenPoolRateLimiterConfig) (cldf.Environment, *burn_mint_erc677.BurnMintERC677, *burn_mint_token_pool.BurnMintTokenPool, error) {
	suiChains := e.BlockChains.SuiChains()
	suiChain := suiChains[suiChainSel]

	evmChain := e.BlockChains.EVMChains()[evmChainSel]

	evmDeployerKey := evmChain.DeployerKey
	deployerSuiAddr, err := suiChain.Signer.GetAddress()
	if err != nil {
		return cldf.Environment{}, nil, nil, errors.New("failed to get deployer address " + err.Error())
	}
	state, err := stateview.LoadOnchainState(e)
	if err != nil {
		return cldf.Environment{}, nil, nil, errors.New("failed load onstate chains " + err.Error())
	}

	linkTokenPkgID := state.SuiChains[suiChainSel].LinkTokenAddress
	linkTokenObjectMetadataID := state.SuiChains[suiChainSel].LinkTokenCoinMetadataId
	linkTokenTreasuryCapID := state.SuiChains[suiChainSel].LinkTokenTreasuryCapId

	// Deploy transferrable token on EVM
	evmToken, evmPool, err := deployTransferTokenOneEnd(e.Logger, evmChain, evmDeployerKey, e.ExistingAddresses, "TOKEN")
	if err != nil {
		return cldf.Environment{}, nil, nil, errors.New("failed to deploy transfer token for evm chain " + err.Error())
	}

	err = attachTokenToTheRegistry(evmChain, state.MustGetEVMChainState(evmChain.Selector), evmDeployerKey, evmToken.Address(), evmPool.Address())
	if err != nil {
		return cldf.Environment{}, nil, nil, errors.New("failed to attach token to registry for evm " + err.Error())
	}

	// Deploy & Configure LockRelease TP on SUI
	e, _, err = commoncs.ApplyChangesets(&testing.T{}, e, []commoncs.ConfiguredChangeSet{
		commoncs.Configure(sui_cs.DeployTPAndConfigure{}, sui_cs.DeployTPAndConfigureConfig{
			SuiChainSelector: suiChainSel,
			TokenPoolTypes:   []sui_deployment.TokenPoolType{sui_deployment.TokenPoolTypeLockRelease},
			LockReleaseTPInput: lockreleasetokenpoolops.DeployAndInitLockReleaseTokenPoolInput{
				CoinObjectTypeArg:    linkTokenPkgID + "::link::LINK",
				CoinMetadataObjectId: linkTokenObjectMetadataID,
				TreasuryCapObjectId:  linkTokenTreasuryCapID,
				Rebalancer:           deployerSuiAddr,

				// apply dest chain updates
				RemoteChainSelectorsToRemove: []uint64{},
				RemoteChainSelectorsToAdd:    []uint64{evmChainSel},
				RemotePoolAddressesToAdd:     [][]string{{evmPool.Address().String()}}, // this gets convert to 32byte bytes internally
				RemoteTokenAddressesToAdd: []string{
					evmToken.Address().String(), // this gets convert to 32byte bytes internally
				},

				// set chain rate limiter configs
				RemoteChainSelectors: extractFields(rateLimiterConfigs, func(c TokenPoolRateLimiterConfig) uint64 { return c.RemoteChainSelector }),
				OutboundIsEnableds:   extractFields(rateLimiterConfigs, func(c TokenPoolRateLimiterConfig) bool { return c.OutboundIsEnabled }),
				OutboundCapacities:   extractFields(rateLimiterConfigs, func(c TokenPoolRateLimiterConfig) uint64 { return c.OutboundCapacity }),
				OutboundRates:        extractFields(rateLimiterConfigs, func(c TokenPoolRateLimiterConfig) uint64 { return c.OutboundRate }),
				InboundIsEnableds:    extractFields(rateLimiterConfigs, func(c TokenPoolRateLimiterConfig) bool { return c.InboundIsEnabled }),
				InboundCapacities:    extractFields(rateLimiterConfigs, func(c TokenPoolRateLimiterConfig) uint64 { return c.InboundCapacity }),
				InboundRates:         extractFields(rateLimiterConfigs, func(c TokenPoolRateLimiterConfig) uint64 { return c.InboundRate }),
			},
		}),
	})
	if err != nil {
		return cldf.Environment{}, nil, nil, err
	}

	// reload onChainState to get deployed TP contracts
	state, err = stateview.LoadOnchainState(e)
	if err != nil {
		return cldf.Environment{}, nil, nil, errors.New("failed load onstate chains " + err.Error())
	}

	suiTokenBytes, err := hex.DecodeString(strings.TrimPrefix(linkTokenObjectMetadataID, "0x"))
	if err != nil {
		return cldf.Environment{}, nil, nil, errors.New("error while decoding suiToken")
	}

	lnrTokenPool, ok := state.SuiChains[suiChainSel].LnRTokenPools[TokenSymbolLINK]
	if !ok {
		return cldf.Environment{}, nil, nil, fmt.Errorf("no LockReleaseTokenPool found for token: %s", TokenSymbolLINK)
	}

	suiPoolBytes, err := hex.DecodeString(strings.TrimPrefix(lnrTokenPool.PackageID, "0x"))
	if err != nil {
		return cldf.Environment{}, nil, nil, errors.New("error while decoding suiPool")
	}

	err = setTokenPoolCounterPart(e.BlockChains.EVMChains()[evmChain.Selector], evmPool, evmDeployerKey, suiChain.Selector, suiTokenBytes, suiPoolBytes)
	if err != nil {
		return cldf.Environment{}, nil, nil, errors.New("failed to add token to the counterparty " + err.Error())
	}

	err = grantMintBurnPermissions(e.Logger, e.BlockChains.EVMChains()[evmChain.Selector], evmToken, evmDeployerKey, evmPool.Address())
	if err != nil {
		return cldf.Environment{}, nil, nil, errors.New("failed to grant burn mint " + err.Error())
	}

	return e, evmToken, evmPool, nil
}

func WaitForTokenBalanceSui(
	ctx context.Context,
	t *testing.T,
	fungibleAsset string,
	account string,
	chain cldf_sui.Chain,
	expected *big.Int,
) {
	require.Eventually(t, func() bool {
		balanceReq := models.SuiXGetBalanceRequest{
			Owner:    account,
			CoinType: fungibleAsset + "::link::LINK", // Sui Link token Type
		}

		response, err := chain.Client.SuiXGetBalance(ctx, balanceReq)
		require.NoError(t, err)

		balance, ok := new(big.Int).SetString(response.TotalBalance, 10)
		require.True(t, ok)

		return balance.Cmp(expected) == 0
	}, tests.WaitTimeout(t), 2000*time.Millisecond)
}

func UpgradeContractDirect(
	ctx context.Context,
	callOpts *suiBind.CallOpts, // must include Signer, GasBudget, WaitForExecution
	client sui.ISuiAPI,
	packageToUpgrade string,
	upgradeCapID string,
	modules [][]byte,
	dependencies []models.SuiAddress,
	policy byte,
	digest []byte,
) (*models.SuiTransactionBlockResponse, error) {
	lggr, _ := logger.New()

	ptb := suitx.NewTransaction()
	ptb.SetSuiClient(client.(*sui.Client))

	packageContract, err := suiBind.NewBoundContract(
		"0x2",     // Framework package
		"sui",     // Package name
		"package", // Module name
		client,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to suibind package module: %w", err)
	}

	// Encode authorize_upgrade call
	typeArgsList := []string{}
	typeParamsList := []string{}

	normalizedModulePackage, err := client.SuiGetNormalizedMoveModule(ctx, models.GetNormalizedMoveModuleRequest{
		Package:    "0x2",
		ModuleName: "package",
	})
	if err != nil {
		return nil, errors.New("failed to get normalized module: " + err.Error())
	}

	functionSignatureAuthorizeUpgrade, isValidaAU := normalizedModulePackage.ExposedFunctions["authorize_upgrade"]
	if !isValidaAU {
		return nil, errors.New("missing function signature for receiver function not found in module authorize_upgrade")
	}

	// Figure out the parameter types from the normalized module of the token pool
	paramTypesAuthorizeUpgrade, err := suiofframp_helper.DecodeParameters(lggr, functionSignatureAuthorizeUpgrade.(map[string]any), "parameters")
	if err != nil {
		return nil, errors.New("failed to decode parameters for commit upgrade function: " + err.Error())
	}

	paramValues := []any{
		suiBind.Object{Id: upgradeCapID},
		policy,
		digest,
	}

	authCall, err := packageContract.EncodeCallArgsWithGenerics(
		"authorize_upgrade",
		typeArgsList,
		typeParamsList,
		paramTypesAuthorizeUpgrade,
		paramValues,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to encode authorize_upgrade call: %w", err)
	}

	authResult, err := packageContract.AppendPTB(ctx, callOpts, ptb, authCall)
	if err != nil {
		return nil, fmt.Errorf("failed to append authorize_upgrade to PTB: %w", err)
	}

	// Append the Upgrade command (consumes UpgradeTicket)
	upgradeReceiptArg := ptb.Upgrade(
		modules,                             // Raw bytes (from Call)
		dependencies,                        // Dependencies as addresses (from Call)
		models.SuiAddress(packageToUpgrade), // Package being upgraded (from Call)
		*authResult,                         // UpgradeTicket from authorize step
	)

	// commit the ticket
	typeArgsListCommit := []string{}
	typeParamsListCommit := []string{}

	functionSignatureCommitUpgrade, isValidaCU := normalizedModulePackage.ExposedFunctions["commit_upgrade"]
	if !isValidaCU {
		return nil, errors.New("missing function signature for receiver function not found in module commit_upgrade")
	}

	// Figure out the parameter types from the normalized module of the token pool
	paramTypesCommitUpgrade, err := suiofframp_helper.DecodeParameters(lggr, functionSignatureCommitUpgrade.(map[string]any), "parameters")
	if err != nil {
		return nil, errors.New("failed to decode parameters for commit upgrade function: " + err.Error())
	}

	paramValuesCommit := []any{
		suiBind.Object{Id: upgradeCapID},
		upgradeReceiptArg,
	}

	commitEncoded, err := packageContract.EncodeCallArgsWithGenerics(
		"commit_upgrade",
		typeArgsListCommit,
		typeParamsListCommit,
		paramTypesCommitUpgrade,
		paramValuesCommit,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to encode commit_upgrade call: %w", err)
	}

	_, err = packageContract.AppendPTB(ctx, callOpts, ptb, commitEncoded)
	if err != nil {
		return nil, fmt.Errorf("failed to append commit_upgrade to PTB: %w", err)
	}

	// ️ Execute PTB
	resp, err := suiBind.ExecutePTB(ctx, callOpts, client, ptb)
	if err != nil {
		return nil, fmt.Errorf("failed executing upgrade PTB: %w", err)
	}

	return resp, nil
}

// Helper functions to extract rate limiter config fields from array using concise Go patterns
func extractFields[T any](configs []TokenPoolRateLimiterConfig, selector func(TokenPoolRateLimiterConfig) T) []T {
	result := make([]T, len(configs))
	for i, config := range configs {
		result[i] = selector(config)
	}
	return result
}
