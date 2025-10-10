package testhelpers

import (
	"context"
	"encoding/hex"
	"errors"
	"math/big"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/block-vision/sui-go-sdk/models"
	"github.com/block-vision/sui-go-sdk/sui"
	suitx "github.com/block-vision/sui-go-sdk/transaction"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_5_1/burn_mint_token_pool"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_3/message_hasher"

	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"

	"github.com/smartcontractkit/chainlink-evm/gethwrappers/shared/generated/initial/burn_mint_erc677"

	suiBind "github.com/smartcontractkit/chainlink-sui/bindings/bind"
	sui_cs "github.com/smartcontractkit/chainlink-sui/deployment/changesets"
	sui_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops"
	ccipops "github.com/smartcontractkit/chainlink-sui/deployment/ops/ccip"
	burnminttokenpoolops "github.com/smartcontractkit/chainlink-sui/deployment/ops/ccip_burn_mint_token_pool"
	suiofframp_helper "github.com/smartcontractkit/chainlink-sui/relayer/chainwriter/ptb/offramp"

	suideps "github.com/smartcontractkit/chainlink/deployment/ccip/changeset/sui"
	ccipclient "github.com/smartcontractkit/chainlink/deployment/ccip/shared/client"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
	commoncs "github.com/smartcontractkit/chainlink/deployment/common/changeset"

	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/ccipevm"

	cldf_sui "github.com/smartcontractkit/chainlink-deployments-framework/chain/sui"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
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
	ccipPackageID := state.SuiChains[cfg.SourceChain].CCIPAddress
	onRampPackageID := state.SuiChains[cfg.SourceChain].OnRampAddress
	onRampStateObjectID := state.SuiChains[cfg.SourceChain].OnRampStateObjectId
	linkTokenPkgID := state.SuiChains[cfg.SourceChain].LinkTokenAddress
	linkTokenObjectMetadataID := state.SuiChains[cfg.SourceChain].LinkTokenCoinMetadataId
	ccipOwnerCapID := state.SuiChains[cfg.SourceChain].CCIPOwnerCapObjectId

	bigIntSourceUsdPerToken, parsed := new(big.Int).SetString("21377040000000000000000000000", 10) // 1e27 since sui is 1e9
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
		BurnMintTPPkgID := state.SuiChains[cfg.SourceChain].CCIPBurnMintTokenPool
		BurnMintTPState := state.SuiChains[cfg.SourceChain].CCIPBurnMintTokenPoolState

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

		BurnMintTPContract, err := suiBind.NewBoundContract(
			BurnMintTPPkgID,
			"burn_mint_token_pool",
			"burn_mint_token_pool",
			client,
		)
		if err != nil {
			return nil, errors.New("failed to create burn_mint_token_pool bound contract when appending PTB command: " + err.Error())
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

		createTokenTransferParamsResult, err := ccipStateHelperContract.AppendPTB(ctx, deps.SuiChain.GetCallOpts(), ptb, onRampCreateTokenTransferParamsCall)
		if err != nil {
			return nil, errors.New("failed to build PTB (get_token_param_data) using bindings: " + err.Error())
		}

		/*********  2. lock_or_burn *******/
		normalizedModuleBMTP, err := client.SuiGetNormalizedMoveModule(ctx, models.GetNormalizedMoveModuleRequest{
			Package:    BurnMintTPPkgID,
			ModuleName: "burn_mint_token_pool",
		})
		if err != nil {
			return nil, errors.New("failed to get normalized module: " + err.Error())
		}

		functionSignatureLnB, isValidLockOrBurn := normalizedModuleBMTP.ExposedFunctions["lock_or_burn"]
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
		paramValuesLockBurn := []any{
			suiBind.Object{Id: ccipObjectRefID},           // ref
			createTokenTransferParamsResult,               // token_params
			suiBind.Object{Id: msg.TokenAmounts[0].Token}, // minted token to send to EVM
			cfg.DestChain,
			suiBind.Object{Id: "0x6"},           // clock
			suiBind.Object{Id: BurnMintTPState}, // BurnMintstate
		}

		lockOrBurnParamsCall, err := BurnMintTPContract.EncodeCallArgsWithGenerics(
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

		_, err = BurnMintTPContract.AppendPTB(ctx, deps.SuiChain.GetCallOpts(), ptb, lockOrBurnParamsCall)
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
			[]byte("hello evm from sui"),
			createTokenTransferParamsResult,               // tokenParams from the original create_token_transfer_params
			suiBind.Object{Id: linkTokenObjectMetadataID}, // feeTokenMetadata
			suiBind.Object{Id: msg.FeeToken},
			[]byte{}, // extraArgs
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
		[]byte("hello evm from sui"),
		extractedAny2SuiMessageResult,                 // tokenParams
		suiBind.Object{Id: linkTokenObjectMetadataID}, // feeTokenMetadata
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

	if len(executeCCIPSend.Events) == 0 {
		return nil, errors.New("no events returned from Sui CCIPSend")
	}

	suiEvent := executeCCIPSend.Events[0].ParsedJson

	seqStr, _ := suiEvent["sequence_number"].(string)
	seq, _ := strconv.ParseUint(seqStr, 10, 64)

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

func HandleTokenAndPoolDeploymentForSUI(e cldf.Environment, suiChainSel, evmChainSel uint64) (cldf.Environment, *burn_mint_erc677.BurnMintERC677, *burn_mint_token_pool.BurnMintTokenPool, error) {
	suiChains := e.BlockChains.SuiChains()
	suiChain := suiChains[suiChainSel]

	evmChain := e.BlockChains.EVMChains()[evmChainSel]

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

	// Deploy transferrable token on EVM
	evmToken, evmPool, err := deployTransferTokenOneEnd(e.Logger, evmChain, evmDeployerKey, e.ExistingAddresses, "TOKEN")
	if err != nil {
		return cldf.Environment{}, nil, nil, errors.New("failed to deploy transfer token for evm chain " + err.Error())
	}

	err = attachTokenToTheRegistry(evmChain, state.MustGetEVMChainState(evmChain.Selector), evmDeployerKey, evmToken.Address(), evmPool.Address())
	if err != nil {
		return cldf.Environment{}, nil, nil, errors.New("failed to attach token to registry for evm " + err.Error())
	}

	// eploy & Configure BurnMint TP on SUI
	e, _, err = commoncs.ApplyChangesets(&testing.T{}, e, []commoncs.ConfiguredChangeSet{
		commoncs.Configure(sui_cs.DeployTPAndConfigure{}, sui_cs.DeployTPAndConfigureConfig{
			SuiChainSelector: suiChainSel,
			TokenPoolTypes:   []string{"bnm"},
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
	suiPoolBytes, err := hex.DecodeString(strings.TrimPrefix(state.SuiChains[suiChainSel].CCIPBurnMintTokenPool, "0x"))
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
