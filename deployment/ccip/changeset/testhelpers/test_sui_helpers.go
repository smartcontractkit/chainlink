package testhelpers

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/block-vision/sui-go-sdk/models"
	"github.com/block-vision/sui-go-sdk/sui"
	suitx "github.com/block-vision/sui-go-sdk/transaction"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_5_1/burn_mint_token_pool"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_3/message_hasher"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/shared/generated/initial/burn_mint_erc677"
	suiBind "github.com/smartcontractkit/chainlink-sui/bindings/bind"
	module_fee_quoter "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/ccip/fee_quoter"
	sui_cs "github.com/smartcontractkit/chainlink-sui/deployment/changesets"
	sui_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops"
	ccipops "github.com/smartcontractkit/chainlink-sui/deployment/ops/ccip"
	"github.com/smartcontractkit/chainlink-sui/relayer/chainwriter/config"
	suiofframp_helper "github.com/smartcontractkit/chainlink-sui/relayer/chainwriter/ptb/offramp"
	suicodec "github.com/smartcontractkit/chainlink-sui/relayer/codec"
	"github.com/smartcontractkit/chainlink-sui/relayer/testutils"
	suideps "github.com/smartcontractkit/chainlink/deployment/ccip/changeset/sui"
	ccipclient "github.com/smartcontractkit/chainlink/deployment/ccip/shared/client"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
	commoncs "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/ccipevm"
	"github.com/stretchr/testify/require"

	cldf_sui "github.com/smartcontractkit/chainlink-deployments-framework/chain/sui"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	sui_module_bnmtp "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/ccip_token_pools/burn_mint_token_pool"
	burnminttokenpoolops "github.com/smartcontractkit/chainlink-sui/deployment/ops/ccip_burn_mint_token_pool"
)

type SuiSendRequest struct {
	Receiver      []byte
	Data          []byte
	ExtraArgs     []byte
	FeeToken      string
	FeeTokenStore string
	TokenAmounts  []SuiTokenAmount
}

type SuiTokenAmount struct {
	Token  string
	Amount uint64
}

type RampMessageHeader struct {
	MessageId           []byte `json:"message_id"`
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

func baseCCIPConfig(
	ccipPkg string,
	pubKey []byte,
	extra []config.ChainWriterPTBCommand,
) config.ChainWriterConfig {
	// common PTB command 0: create_token_params
	cmds := []config.ChainWriterPTBCommand{{
		Type:      suicodec.SuiPTBCommandMoveCall,
		PackageId: strPtr(ccipPkg),
		ModuleId:  strPtr("onramp_state_helper"),
		Function:  strPtr("create_token_transfer_params"),
		Params: []suicodec.SuiFunctionParam{
			{
				Name:     "token_receiver",
				Type:     "address",
				Required: true,
			},
		},
	}}
	// append the variant commands
	cmds = append(cmds, extra...)

	return config.ChainWriterConfig{
		Modules: map[string]*config.ChainWriterModule{
			config.PTBChainWriterModuleName: {
				Name:     config.PTBChainWriterModuleName,
				ModuleID: "0x123",
				Functions: map[string]*config.ChainWriterFunction{
					"ccip_send": {
						Name:        "ccip_send",
						PublicKey:   pubKey,
						Params:      []suicodec.SuiFunctionParam{},
						PTBCommands: cmds,
					},
				},
			},
		},
	}
}

// 2a) Simple Message → EVM
func configureChainWriterForMsg(
	ccipPkg, onRampPkg string,
	pubKey []byte,
	feeTokenPkgId string,
) config.ChainWriterConfig {
	feeTokenType := fmt.Sprintf("%s::link::LINK", feeTokenPkgId)
	extra := []config.ChainWriterPTBCommand{{
		Type:      suicodec.SuiPTBCommandMoveCall,
		PackageId: strPtr(onRampPkg),
		ModuleId:  strPtr("onramp"),
		Function:  strPtr("ccip_send"),
		Params: []suicodec.SuiFunctionParam{
			{Name: "ref", Type: "object_id", Required: true, IsMutable: testutils.BoolPointer(true)},
			{Name: "state", Type: "object_id", Required: true, IsMutable: testutils.BoolPointer(true)},
			{Name: "clock", Type: "object_id", Required: true, IsMutable: testutils.BoolPointer(false)},
			{Name: "dest_chain_selector", Type: "u64", Required: true},
			{Name: "receiver", Type: "vector<u8>", Required: true},
			{Name: "data", Type: "vector<u8>", Required: true},
			{Name: "token_params", Type: "ptb_dependency", Required: true,
				PTBDependency: &suicodec.PTBCommandDependency{CommandIndex: 0}},
			{Name: "fee_token_metadata", Type: "object_id", Required: true, IsMutable: testutils.BoolPointer(false), GenericType: strPtr(feeTokenType)},
			{Name: "fee_token", Type: "object_id", Required: true},
			{Name: "extra_args", Type: "vector<u8>", Required: true},
		},
	}}
	return baseCCIPConfig(ccipPkg, pubKey, extra)
}

// 2b) Message + BurnMintTP → EVM
func configureChainWriterForMultipleTokens(
	ccipPkg, onRampPkg string,
	pubKey []byte,
	tokenPool string,
) config.ChainWriterConfig {
	extra := []config.ChainWriterPTBCommand{
		// lock-or-burn command
		{
			Type:      suicodec.SuiPTBCommandMoveCall,
			PackageId: strPtr(tokenPool),
			ModuleId:  strPtr("burn_mint_token_pool"),
			Function:  strPtr("lock_or_burn"),
			Params: []suicodec.SuiFunctionParam{
				{Name: "ref", Type: "object_id", Required: true},
				{Name: "clock", Type: "object_id", Required: true, IsMutable: testutils.BoolPointer(false)},
				{Name: "state", Type: "object_id", Required: true},
				{Name: "c", Type: "object_id", Required: true},
				{Name: "token_params", Type: "ptb_dependency", Required: true,
					PTBDependency: &suicodec.PTBCommandDependency{CommandIndex: 0}},
			},
		},
		// the same onramp send
		{
			Type:      suicodec.SuiPTBCommandMoveCall,
			PackageId: strPtr(onRampPkg),
			ModuleId:  strPtr("onramp"),
			Function:  strPtr("ccip_send"),
			Params: []suicodec.SuiFunctionParam{
				{Name: "ref", Type: "object_id", Required: true},
				{Name: "state", Type: "object_id", Required: true},
				{Name: "clock", Type: "object_id", Required: true, IsMutable: testutils.BoolPointer(false)},
				{Name: "dest_chain_selector", Type: "u64", Required: true},
				{Name: "receiver", Type: "vector<u8>", Required: true},
				{Name: "data", Type: "vector<u8>", Required: true},
				{Name: "token_params", Type: "ptb_dependency", Required: true,
					PTBDependency: &suicodec.PTBCommandDependency{CommandIndex: 1}},
				{Name: "fee_token_metadata", Type: "object_id", Required: true, IsMutable: testutils.BoolPointer(false)},
				{Name: "fee_token", Type: "object_id", Required: true},
				{Name: "extra_args", Type: "vector<u8>", Required: true},
			},
		},
	}
	return baseCCIPConfig(ccipPkg, pubKey, extra)
}

func BuildPTBArgs(baseArgs map[string]any, coinType string, extraArgs map[string]any) config.Arguments {
	args := make(map[string]any, len(baseArgs)+len(extraArgs))
	for k, v := range baseArgs {
		args[k] = v
	}
	for k, v := range extraArgs {
		args[k] = v
	}

	argTypes := map[string]string{
		"fee_token": coinType,
	}
	if _, ok := args["c"]; ok {
		argTypes["c"] = coinType
	}

	return config.Arguments{
		Args:     args,
		ArgTypes: argTypes,
	}
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

	deps := suideps.SuiDeps{
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

	ccipObjectRefId := state.SuiChains[cfg.SourceChain].CCIPObjectRef
	ccipPackageId := state.SuiChains[cfg.SourceChain].CCIPAddress
	onRampPackageId := state.SuiChains[cfg.SourceChain].OnRampAddress
	onRampStateObjectId := state.SuiChains[cfg.SourceChain].OnRampStateObjectId
	linkTokenPkgId := state.SuiChains[cfg.SourceChain].LinkTokenAddress
	linkTokenObjectMetadataId := state.SuiChains[cfg.SourceChain].LinkTokenCoinMetadataId
	ccipOwnerCapId := state.SuiChains[cfg.SourceChain].CCIPOwnerCapObjectId

	bigIntSourceUsdPerToken, ok := new(big.Int).SetString("21377040000000000000000000000", 10) // 1e27 since sui is 1e9
	if !ok {
		return &ccipclient.AnyMsgSentEvent{}, fmt.Errorf("failed converting SourceUSDPerToken to bigInt")
	}

	bigIntGasUsdPerUnitGas, ok := new(big.Int).SetString("41946474500", 10) // optimism 4145822215
	if !ok {
		return &ccipclient.AnyMsgSentEvent{}, fmt.Errorf("failed converting GasUsdPerUnitGas to bigInt")
	}

	// getValidatedFee
	msg := cfg.Message.(SuiSendRequest)

	// Update Prices on FeeQuoter with minted LinkToken
	_, err = operations.ExecuteOperation(e.OperationsBundle, ccipops.FeeQuoterUpdatePricesWithOwnerCapOp, deps.SuiChain,
		ccipops.FeeQuoterUpdatePricesWithOwnerCapInput{
			CCIPPackageId:         ccipPackageId,
			CCIPObjectRef:         ccipObjectRefId,
			OwnerCapObjectId:      ccipOwnerCapId,
			SourceTokens:          []string{linkTokenObjectMetadataId},
			SourceUsdPerToken:     []*big.Int{bigIntSourceUsdPerToken},
			GasDestChainSelectors: []uint64{cfg.DestChain},
			GasUsdPerUnitGas:      []*big.Int{bigIntGasUsdPerUnitGas},
		})
	if err != nil {
		return &ccipclient.AnyMsgSentEvent{}, fmt.Errorf("failed to updatePrice for Sui chain %d: %w", cfg.SourceChain, err)
	}

	feeQuoter, err := module_fee_quoter.NewFeeQuoter(ccipPackageId, deps.SuiChain.Client)
	if err != nil {
		return &ccipclient.AnyMsgSentEvent{}, err
	}

	validatedFee, err := feeQuoter.DevInspect().GetValidatedFee(ctx, &suiBind.CallOpts{
		Signer:           deps.SuiChain.Signer,
		WaitForExecution: true,
	},
		suiBind.Object{Id: ccipObjectRefId},
		suiBind.Object{Id: "0x6"},
		cfg.DestChain,
		[]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0xdd, 0xbb, 0x6f, 0x35,
			0x8f, 0x29, 0x04, 0x08, 0xd7, 0x68, 0x47, 0xb4,
			0xf6, 0x02, 0xf0, 0xfd, 0x59, 0x92, 0x95, 0xfd,
		},
		[]byte("hello evm from sui"),
		[]string{},
		[]uint64{},
		linkTokenObjectMetadataId,
		[]byte{},
	)
	if err != nil {
		return &ccipclient.AnyMsgSentEvent{}, err
	}

	fmt.Println("VALIDATED FEE:", validatedFee)

	if len(msg.TokenAmounts) > 0 {
		fmt.Println("TOKEN TRANSFER DETECTED: ", msg)

		BurnMintTPPkgId := state.SuiChains[cfg.SourceChain].CCIPBurnMintTokenPool
		BurnMintTPState := state.SuiChains[cfg.SourceChain].CCIPBurnMintTokenPoolState

		fmt.Println("BNMTPPKGID: ", BurnMintTPPkgId, BurnMintTPState)

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
			ccipPackageId,
			"ccip",
			"onramp_state_helper",
			client,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to create onramp state helper bound contract when appending PTB command: %w", err)
		}

		BurnMintTPContract, err := suiBind.NewBoundContract(
			BurnMintTPPkgId,
			"burn_mint_token_pool",
			"burn_mint_token_pool",
			client,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to create burn_mint_token_pool bound contract when appending PTB command: %w", err)
		}

		onRampContract, err := suiBind.NewBoundContract(
			onRampPackageId,
			"ccip_onramp",
			"onramp",
			client,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to create ccip_onramp contract when appending PTB command: %w", err)
		}

		/*********  1. create_token_transfer_params *******/
		typeArgsList := []string{}
		typeParamsList := []string{}
		paramTypes := []string{
			"address",
		}
		paramValues := []any{
			// For SUI -> EVM BurnMint Pool token Transfer, we can use msg.Reciever as tokenReciever, this field is only used in usdc token pool
			// bc we need to check the recipient with Circle's packages from the onramp side before sending USDC. and it's not used anyway else.
			"0x0000000000000000000000000000000000000000000000000000000000000000",
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
			return nil, fmt.Errorf("failed to encode onRampCreateTokenTransferParamsCall call: %w", err)
		}

		createTokenTransferParamsResult, err := ccipStateHelperContract.AppendPTB(ctx, deps.SuiChain.GetCallOpts(), ptb, onRampCreateTokenTransferParamsCall)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB (get_token_param_data) using bindings: %w", err)
		}

		/*********  2. lock_or_burn *******/
		normalizedModuleBMTP, err := client.SuiGetNormalizedMoveModule(ctx, models.GetNormalizedMoveModuleRequest{
			Package:    BurnMintTPPkgId,
			ModuleName: "burn_mint_token_pool",
		})
		if err != nil {
			return nil, fmt.Errorf("failed to get normalized module: %w", err)
		}

		functionSignatureLnB, ok := normalizedModuleBMTP.ExposedFunctions["lock_or_burn"]
		if !ok {
			return nil, fmt.Errorf("missing function signature for receiver function not found in module (%s)", "lock_or_burn")
		}

		// Figure out the parameter types from the normalized module of the token pool
		paramTypesLockBurn, err := suiofframp_helper.DecodeParameters(e.Logger, functionSignatureLnB.(map[string]any), "parameters")
		if err != nil {
			return nil, fmt.Errorf("failed to decode parameters for token pool function: %w", err)
		}

		typeArgsListLinkTokenPkgId := []string{linkTokenPkgId + "::link::LINK"}
		typeParamsList = []string{}
		paramValuesLockBurn := []any{
			suiBind.Object{Id: ccipObjectRefId},           // ref
			createTokenTransferParamsResult,               // token_params
			suiBind.Object{Id: msg.TokenAmounts[0].Token}, // minted token to send to EVM
			cfg.DestChain,
			suiBind.Object{Id: "0x6"},           // clock
			suiBind.Object{Id: BurnMintTPState}, // BurnMintstate
		}

		lockOrBurnParamsCall, err := BurnMintTPContract.EncodeCallArgsWithGenerics(
			"lock_or_burn",
			typeArgsListLinkTokenPkgId,
			typeParamsList,
			paramTypesLockBurn,
			paramValuesLockBurn,
			nil,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to encode lockOrBurnParamsCall call: %w", err)
		}

		_, err = BurnMintTPContract.AppendPTB(ctx, deps.SuiChain.GetCallOpts(), ptb, lockOrBurnParamsCall)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB (get_token_param_data) using bindings: %w", err)
		}

		/********* 3. ccip_send *******/
		normalizedModule, err := client.SuiGetNormalizedMoveModule(ctx, models.GetNormalizedMoveModuleRequest{
			Package:    onRampPackageId,
			ModuleName: "onramp",
		})
		if err != nil {
			return nil, fmt.Errorf("failed to get normalized module: %w", err)
		}

		functionSignature, ok := normalizedModule.ExposedFunctions["ccip_send"]
		if !ok {
			return nil, fmt.Errorf("missing function signature for receiver function not found in module (%s)", "ccip_send")
		}

		// Figure out the parameter types from the normalized module of the token pool
		paramTypesCCIPSend, err := suiofframp_helper.DecodeParameters(e.Logger, functionSignature.(map[string]any), "parameters")
		if err != nil {
			return nil, fmt.Errorf("failed to decode parameters for token pool function: %w", err)
		}

		paramValuesCCIPSend := []any{
			suiBind.Object{Id: ccipObjectRefId},
			suiBind.Object{Id: onRampStateObjectId},
			suiBind.Object{Id: "0x6"},
			cfg.DestChain,
			msg.Receiver, // receiver
			[]byte("hello evm from sui"),
			createTokenTransferParamsResult,               // tokenParams from the original create_token_transfer_params
			suiBind.Object{Id: linkTokenObjectMetadataId}, // feeTokenMetadata
			suiBind.Object{Id: msg.FeeToken},
			[]byte{}, // extraArgs
		}

		encodedOnRampCCIPSendCall, err := onRampContract.EncodeCallArgsWithGenerics(
			"ccip_send",
			typeArgsListLinkTokenPkgId,
			[]string{},
			paramTypesCCIPSend,
			paramValuesCCIPSend,
			nil,
		)

		_, err = onRampContract.AppendPTB(ctx, deps.SuiChain.GetCallOpts(), ptb, encodedOnRampCCIPSendCall)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB (receiver call) using bindings: %w", err)
		}

		executeCCIPSend, err := suiBind.ExecutePTB(ctx, deps.SuiChain.GetCallOpts(), client, ptb)
		if err != nil {
			return nil, fmt.Errorf("failed to execute ccip_send with err: %w", err)
		}

		suiEvent := executeCCIPSend.Events[2].ParsedJson

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

	ccipStateHelperContract, err := suiBind.NewBoundContract(
		ccipPackageId,
		"ccip",
		"onramp_state_helper",
		client,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create onramp state helper bound contract when appending PTB command: %w", err)
	}

	// Note: these will be different for token transfers
	typeArgsList := []string{}
	typeParamsList := []string{}
	paramTypes := []string{
		"address",
	}
	paramValues := []any{
		"0x0000000000000000000000000000000000000000000000000000000000000000", // zero address as it's msg transfer
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
		return nil, fmt.Errorf("failed to encode onRampCreateTokenTransferParamsCall call: %w", err)
	}

	extractedAny2SuiMessageResult, err := ccipStateHelperContract.AppendPTB(ctx, deps.SuiChain.GetCallOpts(), ptb, onRampCreateTokenTransferParamsCall)
	if err != nil {
		return nil, fmt.Errorf("failed to build PTB (get_token_param_data) using bindings: %w", err)
	}

	onRampContract, err := suiBind.NewBoundContract(
		onRampPackageId,
		"ccip_onramp",
		"onramp",
		client,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create onramp bound contract when appending PTB command: %w", err)
	}

	// normalize module
	normalizedModule, err := client.SuiGetNormalizedMoveModule(ctx, models.GetNormalizedMoveModuleRequest{
		Package:    onRampPackageId,
		ModuleName: "onramp",
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get normalized module: %w", err)
	}

	functionSignature, ok := normalizedModule.ExposedFunctions["ccip_send"]
	if !ok {
		return nil, fmt.Errorf("missing function signature for receiver function not found in module (%s)", "ccip_send")
	}

	// Figure out the parameter types from the normalized module of the token pool
	paramTypes, err = suiofframp_helper.DecodeParameters(e.Logger, functionSignature.(map[string]any), "parameters")
	if err != nil {
		return nil, fmt.Errorf("failed to decode parameters for token pool function: %w", err)
	}

	typeArgsList = []string{linkTokenPkgId + "::link::LINK"}
	typeParamsList = []string{}
	paramValues = []any{
		suiBind.Object{Id: ccipObjectRefId},
		suiBind.Object{Id: onRampStateObjectId},
		suiBind.Object{Id: "0x6"},
		cfg.DestChain,
		msg.Receiver, // receiver (TODO: replace this with sender Address use environment.NormalizeTo32Bytes(ethereumAddress) from sui repo)
		[]byte("hello evm from sui"),
		extractedAny2SuiMessageResult,                 // tokenParams
		suiBind.Object{Id: linkTokenObjectMetadataId}, // feeTokenMetadata
		suiBind.Object{Id: msg.FeeToken},
		[]byte{}, // extraArgs
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
		return nil, fmt.Errorf("failed to encode receiver call: %w", err)
	}

	_, err = onRampContract.AppendPTB(ctx, deps.SuiChain.GetCallOpts(), ptb, encodedOnRampCCIPSendCall)
	if err != nil {
		return nil, fmt.Errorf("failed to build PTB (receiver call) using bindings: %w", err)
	}

	executeCCIPSend, err := suiBind.ExecutePTB(ctx, deps.SuiChain.GetCallOpts(), client, ptb)
	if err != nil {
		return nil, fmt.Errorf("failed to execute ccip_send with err: %w", err)
	}

	if len(executeCCIPSend.Events) == 0 {
		return nil, fmt.Errorf("no events returned from Sui CCIPSend")
	}

	suiEvent := executeCCIPSend.Events[0].ParsedJson

	seqStr, _ := suiEvent["sequence_number"].(string)
	seq, _ := strconv.ParseUint(seqStr, 10, 64)

	return &ccipclient.AnyMsgSentEvent{
		SequenceNumber: seq,
		RawEvent:       suiEvent, // just dump raw
	}, nil

	// onRamp, err := module_onramp.NewOnramp(onRampPackageId, client)
	// if err != nil {
	// 	return &ccipclient.AnyMsgSentEvent{}, err
	// }
	// ccipSendTx, err := onRamp.CcipSend(
	// 	context.Background(),
	// 	deps.SuiChain.GetCallOpts(),
	// 	[]string{linkTokenPkgId + "::link::LINK"},
	// 	suiBind.Object{Id: ccipObjectRefId},
	// 	suiBind.Object{Id: onRampStateObjectId},
	// 	suiBind.Object{Id: "0x6"},
	// 	cfg.DestChain,
	// 	[]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	// 		0x00, 0x00, 0x00, 0x00, 0xdd, 0xbb, 0x6f, 0x35,
	// 		0x8f, 0x29, 0x04, 0x08, 0xd7, 0x68, 0x47, 0xb4,
	// 		0xf6, 0x02, 0xf0, 0xfd, 0x59, 0x92, 0x95, 0xfd,
	// 	},
	// 	[]byte("hello evm from sui"),
	// 	suiBind.Object{
	// 		// call to onramp_state_helper contract
	// 		// function create_token_transfer_params
	// 		// input arg: token_receiver

	// 	},                                             // tokenParams
	// 	suiBind.Object{Id: linkTokenObjectMetadataId}, // feeTokenMetadata
	// 	suiBind.Object{Id: msg.FeeToken},
	// 	[]byte{},
	// )
	// if err != nil {
	// 	return &ccipclient.AnyMsgSentEvent{}, err
	// }

	// Setup new PTB client
	// keystoreInstance := suitestutils.NewTestKeystore(&testing.T{})
	// priv, err := cldf_sui.PrivateKey(suiChain.Signer)
	// if err != nil {
	// 	return nil, err
	// }
	// keystoreInstance.AddKey(priv)

	// relayerClient, err := client.NewPTBClient(e.Logger, suiChain.URL, nil, 30*time.Second, keystoreInstance, 5, "WaitForEffectsCert")
	// if err != nil {
	// 	return &ccipclient.AnyMsgSentEvent{}, err
	// }

	// e.Logger.Info("relayerClient", relayerClient)

	// store := txm.NewTxmStoreImpl(e.Logger)
	// conf := txm.DefaultConfigSet

	// retryManager := txm.NewDefaultRetryManager(5)
	// gasLimit := big.NewInt(30000000)
	// gasManager := txm.NewSuiGasManager(e.Logger, relayerClient, *gasLimit, 0)

	// txManager, err := txm.NewSuiTxm(e.Logger, relayerClient, keystoreInstance, conf, store, retryManager, gasManager)
	// if err != nil {
	// 	return &ccipclient.AnyMsgSentEvent{}, fmt.Errorf("Failed to create SuiTxm: %v", err)
	// }

	// var chainWriterConfig suicrcwconfig.ChainWriterConfig
	// if BurnMintTP != "" {
	// 	chainWriterConfig = configureChainWriterForMultipleTokens(ccipPackageId, onRampPackageId, publicKeyBytes, BurnMintTP)
	// } else {
	// 	chainWriterConfig = configureChainWriterForMsg(ccipPackageId, onRampPackageId, publicKeyBytes, linkTokenPkgId)
	// }

	// chainWriter, err := chainwriter.NewSuiChainWriter(e.Logger, txManager, chainWriterConfig, false)
	// if err != nil {
	// 	return &ccipclient.AnyMsgSentEvent{}, err
	// }

	// c := context.Background()
	// ctx, cancel := context.WithCancel(c)
	// defer cancel() // to ensure other calls associated with this context are released

	// err = chainWriter.Start(ctx)
	// if err != nil {
	// 	return &ccipclient.AnyMsgSentEvent{}, err
	// }

	// err = txManager.Start(ctx)
	// if err != nil {
	// 	return &ccipclient.AnyMsgSentEvent{}, err
	// }

	// txId := "ccip_send_msg_transfer"
	// err = chainWriter.SubmitTransaction(ctx,
	// 	suicrcwconfig.PTBChainWriterModuleName,
	// 	"ccip_send",
	// 	&ptbArgs,
	// 	txId,
	// 	onRampPackageId, // this is the contract address so onramp in this case
	// 	&commonTypes.TxMeta{GasLimit: big.NewInt(30000000)},
	// 	nil,
	// )
	// if err != nil {
	// 	return &ccipclient.AnyMsgSentEvent{}, err
	// }

	// // TODO: find a better way of handling waitForTransaction
	// time.Sleep(10 * time.Second)
	// status, err := chainWriter.GetTransactionStatus(ctx, txId)

	// if status != commonTypes.Finalized {
	// 	return &ccipclient.AnyMsgSentEvent{}, fmt.Errorf("tx failed to get finalized")
	// }

	// e.Logger.Infof("(Sui) CCIP message sent (tx %s) from chain selector %d to chain selector %d", txId, cfg.SourceChain, cfg.DestChain)

	// chainWriter.Close()
	// txManager.Close()

	// // Query the CCIPSend Event via chainReader
	// chainReaderConfig := crConfig.ChainReaderConfig{
	// 	IsLoopPlugin: false,
	// 	Modules: map[string]*crConfig.ChainReaderModule{
	// 		"onramp": {
	// 			Name: "onramp",
	// 			Events: map[string]*crConfig.ChainReaderEvent{
	// 				"CCIPMessageSent": {
	// 					Name:      "CCIPMessageSent",
	// 					EventType: "CCIPMessageSent",
	// 					EventSelector: client.EventSelector{
	// 						Package: onRampPackageId,
	// 						Module:  "onramp",
	// 						Event:   "CCIPMessageSent",
	// 					},
	// 				},
	// 			},
	// 		},
	// 	},
	// }

	// dbURL := os.Getenv("CL_DATABASE_URL")

	// err = sqltest.RegisterTxDB(dbURL)
	// if err != nil {
	// 	return &ccipclient.AnyMsgSentEvent{}, err
	// }

	// db, err := sqlx.Open(pg.DriverTxWrappedPostgres, uuid.New().String())
	// if err != nil {
	// 	return &ccipclient.AnyMsgSentEvent{}, err
	// }

	// db.MapperFunc(reflectx.CamelToSnakeASCII)

	// // attempt to connect
	// _, err = db.Connx(ctx)
	// if err != nil {
	// 	return &ccipclient.AnyMsgSentEvent{}, err
	// }

	// // Create the indexers
	// txnIndexer := indexer.NewTransactionsIndexer(
	// 	db,
	// 	e.Logger,
	// 	relayerClient,
	// 	10*time.Second,
	// 	10*time.Second,
	// 	// start without any configs, they will be set when ChainReader is initialized and gets a reference
	// 	// to the transaction indexer to avoid having to reading ChainReader configs here as well
	// 	map[string]*crConfig.ChainReaderEvent{},
	// )
	// evIndexer := indexer.NewEventIndexer(
	// 	db,
	// 	e.Logger,
	// 	relayerClient,
	// 	// start without any selectors, they will be added during .Bind() calls on ChainReader
	// 	[]*client.EventSelector{},
	// 	10*time.Second,
	// 	10*time.Second,
	// )
	// indexerInstance := indexer.NewIndexer(
	// 	e.Logger,
	// 	evIndexer,
	// 	txnIndexer,
	// )

	// chainReader, err := chainreader.NewChainReader(ctx, e.Logger, relayerClient, chainReaderConfig, db, indexerInstance)
	// if err != nil {
	// 	return &ccipclient.AnyMsgSentEvent{}, err
	// }

	// err = chainReader.Start(ctx)
	// if err != nil {
	// 	return &ccipclient.AnyMsgSentEvent{}, err
	// }

	// err = indexerInstance.Start(ctx)
	// if err != nil {
	// 	return &ccipclient.AnyMsgSentEvent{}, err
	// }

	// err = chainReader.Bind(context.Background(), []chain_reader_types.BoundContract{{
	// 	Name:    "onramp",
	// 	Address: onRampPackageId, // Package ID of the deployed counter contract
	// }})
	// if err != nil {
	// 	return &ccipclient.AnyMsgSentEvent{}, fmt.Errorf("failed to bind onramp contract with chainReader")
	// }

	// // TODO handle this better, maybe retrieve it from the bindings when we do binding ccip_send
	// e.Logger.Debugw("Querying for ccip_send events",
	// 	"filter", "CCIPMessageSent",
	// 	"limit", 50,
	// 	"packageId", onRampPackageId,
	// 	"contract", "onramp",
	// 	"eventType", "CCIPMessageSent")

	// var ccipSendEvent CCIPMessageSent
	// sequences, err := chainReader.QueryKey(
	// 	ctx,
	// 	chain_reader_types.BoundContract{
	// 		Name:    "onramp",
	// 		Address: onRampPackageId, // Package ID of the deployed counter contract
	// 	},
	// 	sui_query.KeyFilter{
	// 		Key: "CCIPMessageSent",
	// 	},
	// 	sui_query.LimitAndSort{
	// 		Limit: sui_query.Limit{
	// 			Count:  50,
	// 			Cursor: "",
	// 		},
	// 	},
	// 	&ccipSendEvent,
	// )
	// if err != nil {
	// 	return &ccipclient.AnyMsgSentEvent{}, fmt.Errorf("failed to query events: %w", err)
	// }

	// if len(sequences) < 1 {
	// 	return &ccipclient.AnyMsgSentEvent{}, fmt.Errorf("failed to fetch event sequence")
	// }
	// e.Logger.Debugw("Query results", "sequences", sequences)
	// rawevent := sequences[0].Data.(*CCIPMessageSent)

	// chainReader.Close()
	// indexerInstance.Close()

	// return &ccipclient.AnyMsgSentEvent{
	// 	SequenceNumber: rawevent.SequenceNumber,
	// 	RawEvent:       rawevent,
	// }, nil

}

func MakeSuiExtraArgs(gasLimit uint64, allowOOO bool, receiverObjectIds [][32]byte, tokenReceiver [32]byte) []byte {
	extraArgs, err := ccipevm.SerializeClientSUIExtraArgsV1(message_hasher.ClientSuiExtraArgsV1{
		GasLimit:                 new(big.Int).SetUint64(gasLimit),
		AllowOutOfOrderExecution: allowOOO,
		TokenReceiver:            tokenReceiver,
		ReceiverObjectIds:        receiverObjectIds,
	})
	if err != nil {
		panic(err)
	}
	return extraArgs
}

func HandleTokenAndPoolDeploymentForSUI(e cldf.Environment, suiChainSel, evmChainSel uint64) (cldf.Environment, *burn_mint_erc677.BurnMintERC677, *burn_mint_token_pool.BurnMintTokenPool, error) {
	suiChains := e.BlockChains.SuiChains()
	suiChain := suiChains[suiChainSel]

	evmChain := e.BlockChains.EVMChains()[evmChainSel]

	// Deploy Transferrable TOKEN on ETH
	// EVM
	evmDeployerKey := evmChain.DeployerKey
	state, err := stateview.LoadOnchainState(e)
	if err != nil {
		return cldf.Environment{}, nil, nil, fmt.Errorf("failed load onstate chains %w", err)
	}

	linkTokenPkgId := state.SuiChains[suiChainSel].LinkTokenAddress
	linkTokenObjectMetadataId := state.SuiChains[suiChainSel].LinkTokenCoinMetadataId
	linkTokenTreasuryCapId := state.SuiChains[suiChainSel].LinkTokenTreasuryCapId

	// Deploy transferrable token on EVM
	evmToken, evmPool, err := deployTransferTokenOneEnd(e.Logger, evmChain, evmDeployerKey, e.ExistingAddresses, "TOKEN")
	if err != nil {
		return cldf.Environment{}, nil, nil, fmt.Errorf("failed to deploy transfer token for evm chain %d: %w", evmChainSel, err)
	}

	err = attachTokenToTheRegistry(evmChain, state.MustGetEVMChainState(evmChain.Selector), evmDeployerKey, evmToken.Address(), evmPool.Address())
	if err != nil {
		return cldf.Environment{}, nil, nil, fmt.Errorf("failed to attach token to registry for evm %d: %w", evmChainSel, err)
	}

	// // // Deploy & Configure BurnMint TP on SUI
	e, _, err = commoncs.ApplyChangesets(&testing.T{}, e, []commoncs.ConfiguredChangeSet{
		commoncs.Configure(sui_cs.DeployTPAndConfigure{}, sui_cs.DeployTPAndConfigureConfig{
			SuiChainSelector: suiChainSel,
			TokenPoolTypes:   []string{"bnm"},
			BurnMintTpInput: burnminttokenpoolops.DeployAndInitBurnMintTokenPoolInput{
				CoinObjectTypeArg:    linkTokenPkgId + "::link::LINK",
				CoinMetadataObjectId: linkTokenObjectMetadataId,
				TreasuryCapObjectId:  linkTokenTreasuryCapId,

				// apply dest chain updates
				RemoteChainSelectorsToRemove: []uint64{},
				RemoteChainSelectorsToAdd:    []uint64{evmChainSel},
				RemotePoolAddressesToAdd:     [][]string{{evmPool.Address().String()}}, // this gets convert to 32byte bytes internally
				RemoteTokenAddressesToAdd: []string{
					evmToken.Address().String(), // this gets convert to 32byte bytes internally
				},

				// set chain rate limiter configs
				RemoteChainSelectors: []uint64{evmChainSel},
				OutboundIsEnableds:   []bool{false},
				OutboundCapacities:   []uint64{100000},
				OutboundRates:        []uint64{100},
				InboundIsEnableds:    []bool{false},
				InboundCapacities:    []uint64{100000},
				InboundRates:         []uint64{100},
			},
		}),
	})
	if err != nil {
		return cldf.Environment{}, nil, nil, err
	}

	// reload onChainState to get deployed TP contracts
	state, err = stateview.LoadOnchainState(e)
	if err != nil {
		return cldf.Environment{}, nil, nil, fmt.Errorf("failed load onstate chains %w", err)
	}

	fmt.Println("REMOTE TP PARAMS: ", state.SuiChains[suiChainSel].CCIPBurnMintTokenPool, state.SuiChains[suiChainSel].CCIPBurnMintTokenPoolState, state.SuiChains[suiChainSel].CCIPBurnMintTokenPoolOwnerId)

	// ensure tokenPool is added
	// (ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, remoteChainSelector uint64)
	bmtp, err := sui_module_bnmtp.NewBurnMintTokenPool(state.SuiChains[suiChainSel].CCIPBurnMintTokenPool, e.BlockChains.SuiChains()[suiChainSel].Client)
	if err != nil {
		return cldf.Environment{}, nil, nil, err
	}

	val, err := bmtp.DevInspect().GetRemotePools(context.Background(), &suiBind.CallOpts{
		Signer:           e.BlockChains.SuiChains()[suiChainSel].Signer,
		WaitForExecution: true,
	}, []string{linkTokenPkgId + "::link::LINK"}, suiBind.Object{Id: state.SuiChains[suiChainSel].CCIPBurnMintTokenPoolState}, evmChainSel)
	if err != nil {
		return cldf.Environment{}, nil, nil, err
	}

	fmt.Println("REMOTE POOLS ON SUII: ", val)

	val1, err := bmtp.DevInspect().IsRemotePool(context.Background(), &suiBind.CallOpts{
		Signer:           e.BlockChains.SuiChains()[suiChainSel].Signer,
		WaitForExecution: true,
	}, []string{linkTokenPkgId + "::link::LINK"}, suiBind.Object{Id: state.SuiChains[suiChainSel].CCIPBurnMintTokenPoolState}, evmChainSel, evmPool.Address().Bytes())
	if err != nil {
		return cldf.Environment{}, nil, nil, err
	}

	fmt.Println("IS REMOTE POOL ON SUII: ", val1, evmPool.Address().Bytes())

	suiTokenBytes, err := hex.DecodeString(strings.TrimPrefix(linkTokenObjectMetadataId, "0x"))
	if err != nil {
		return cldf.Environment{}, nil, nil, fmt.Errorf("Error while decoding suiToken")
	}
	suiPoolBytes, err := hex.DecodeString(strings.TrimPrefix(state.SuiChains[suiChainSel].CCIPBurnMintTokenPool, "0x"))
	if err != nil {
		return cldf.Environment{}, nil, nil, fmt.Errorf("Error while decoding suiPool")
	}

	err = setTokenPoolCounterPart(e.BlockChains.EVMChains()[evmChain.Selector], evmPool, evmDeployerKey, suiChain.Selector, suiTokenBytes[:], suiPoolBytes[:])
	if err != nil {
		return cldf.Environment{}, nil, nil, fmt.Errorf("failed to add token to the counterparty %d: %w", evmChainSel, err)
	}

	err = grantMintBurnPermissions(e.Logger, e.BlockChains.EVMChains()[evmChain.Selector], evmToken, evmDeployerKey, evmPool.Address())
	if err != nil {
		return cldf.Environment{}, nil, nil, fmt.Errorf("failed to grant burnMint %d: %w", evmChainSel, err)
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
	}, tests.WaitTimeout(t), 500*time.Millisecond)
}
