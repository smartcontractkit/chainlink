package targets

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"

	chainselectors "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-common/pkg/values"
	txmgrcommon "github.com/smartcontractkit/chainlink/v2/common/txmgr"
	abiutil "github.com/smartcontractkit/chainlink/v2/core/chains/evm/abi"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/txmgr"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/chains/legacyevm"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/forwarder"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/pipeline"
)

var forwardABI = evmtypes.MustGetABI(forwarder.KeystoneForwarderMetaData.ABI)

func InitializeWrite(registry commontypes.CapabilitiesRegistry, legacyEVMChains legacyevm.LegacyChainContainer, lggr logger.Logger) error {
	for _, chain := range legacyEVMChains.Slice() {
		capability := NewEvmWrite(chain, lggr)
		if err := registry.Add(context.TODO(), capability); err != nil {
			return err
		}
	}
	return nil
}

var (
	_ capabilities.ActionCapability = &EvmWrite{}
)

const defaultGasLimit = 200000

type EvmWrite struct {
	chain legacyevm.Chain
	capabilities.CapabilityInfo
	lggr logger.Logger
}

func NewEvmWrite(chain legacyevm.Chain, lggr logger.Logger) *EvmWrite {
	// generate ID based on chain selector
	name := fmt.Sprintf("write_%v", chain.ID())
	chainName, err := chainselectors.NameFromChainId(chain.ID().Uint64())
	if err == nil {
		name = fmt.Sprintf("write_%v", chainName)
	}

	info := capabilities.MustNewCapabilityInfo(
		name,
		capabilities.CapabilityTypeTarget,
		"Write target.",
		"v1.0.0",
	)

	return &EvmWrite{
		chain,
		info,
		lggr.Named("EvmWrite"),
	}
}

type EvmConfig struct {
	ChainID uint
	Address string
	Params  []any
	ABI     string
}

// TODO: enforce required key presence

func parseConfig(rawConfig *values.Map) (EvmConfig, error) {
	var config EvmConfig
	configAny, err := rawConfig.Unwrap()
	if err != nil {
		return config, err
	}
	err = mapstructure.Decode(configAny, &config)
	return config, err
}

func evaluateParams(params []any, inputs map[string]any) ([]any, error) {
	vars := pipeline.NewVarsFrom(inputs)
	var args []any
	for _, param := range params {
		switch v := param.(type) {
		case string:
			val, err := pipeline.VarExpr(v, vars)()
			if err == nil {
				args = append(args, val)
			} else if errors.Is(errors.Cause(err), pipeline.ErrParameterEmpty) {
				args = append(args, param)
			} else {
				return args, err
			}
		default:
			args = append(args, param)
		}
	}

	return args, nil
}

func encodePayload(args []any, rawSelector string) ([]byte, error) {
	// TODO: do spec parsing as part of parseConfig()

	// Based on https://github.com/ethereum/go-ethereum/blob/f1c27c286ea2d0e110a507e5749e92d0a6144f08/signer/fourbyte/abi.go#L77-L102

	// NOTE: without having full ABI it's actually impossible to support function overloading
	selector, err := abiutil.ParseSignature(rawSelector)
	if err != nil {
		return nil, err
	}

	abidata, err := json.Marshal([]abi.SelectorMarshaling{selector})
	if err != nil {
		return nil, err
	}

	spec, err := abi.JSON(strings.NewReader(string(abidata)))
	if err != nil {
		return nil, err
	}

	return spec.Pack(selector.Name, args...)

	// NOTE: could avoid JSON encoding/decoding the selector
	// var args abi.Arguments
	// for _, arg := range selector.Inputs {
	// 	ty, err := abi.NewType(arg.Type, arg.InternalType, arg.Components)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	args = append(args, abi.Argument{Name: arg.Name, Type: ty})
	// }
	// // we only care about the name + inputs so we can compute the method ID
	// method := abi.NewMethod(selector.Name, selector.Name, abi.Function, "nonpayable", false, false, args, nil)
	//
	// https://github.com/ethereum/go-ethereum/blob/f1c27c286ea2d0e110a507e5749e92d0a6144f08/accounts/abi/abi.go#L77-L82
	// arguments, err := method.Inputs.Pack(args...)
	// if err != nil {
	// 	return nil, err
	// }
	// // Pack up the method ID too if not a constructor and return
	// return append(method.ID, arguments...), nil
}

func (cap *EvmWrite) Execute(ctx context.Context, callback chan<- capabilities.CapabilityResponse, request capabilities.CapabilityRequest) error {
	cap.lggr.Debugw("Execute", "request", request)
	// TODO: idempotency

	// TODO: extract into ChainWriter?
	txm := cap.chain.TxManager()

	config := cap.chain.Config().EVM().ChainWriter()

	reqConfig, err := parseConfig(request.Config)
	if err != nil {
		return err
	}

	inputsAny, err := request.Inputs.Unwrap()
	if err != nil {
		return err
	}
	inputs := inputsAny.(map[string]any)

	// evaluate any variables in reqConfig.Params
	args, err := evaluateParams(reqConfig.Params, inputs)
	if err != nil {
		return err
	}

	data, err := encodePayload(args, reqConfig.ABI)
	if err != nil {
		return err
	}

	// TODO: validate encoded report is prefixed with workflowID and executionID that match the request meta

	// No signature validation in the MVP demo
	signatures := [][]byte{}

	// construct forwarding payload
	calldata, err := forwardABI.Pack("report", common.HexToAddress(reqConfig.Address), data, signatures)
	if err != nil {
		return err
	}

	txMeta := &txmgr.TxMeta{
		// FwdrDestAddress could also be set for better logging but it's used for various purposes around Operator Forwarders
		WorkflowExecutionID: &request.Metadata.WorkflowExecutionID,
	}
	strategy := txmgrcommon.NewSendEveryStrategy()

	checker := txmgr.TransmitCheckerSpec{
		CheckerType: txmgr.TransmitCheckerTypeSimulate,
	}
	req := txmgr.TxRequest{
		FromAddress:    config.FromAddress().Address(),
		ToAddress:      config.ForwarderAddress().Address(),
		EncodedPayload: calldata,
		FeeLimit:       uint64(defaultGasLimit),
		Meta:           txMeta,
		Strategy:       strategy,
		Checker:        checker,
		// SignalCallback:   true, TODO: add code that checks if a workflow id is present, if so, route callback to chainwriter rather than pipeline
	}
	tx, err := txm.CreateTransaction(ctx, req)
	if err != nil {
		return err
	}
	fmt.Printf("Transaction submitted %v", tx.ID)
	go func() {
		// TODO: cast tx.Error to Err (or Value to Value?)
		callback <- capabilities.CapabilityResponse{
			Value: nil,
			Err:   nil,
		}
		close(callback)
	}()
	return nil
}

func (cap *EvmWrite) RegisterToWorkflow(ctx context.Context, request capabilities.RegisterToWorkflowRequest) error {
	return nil
}

func (cap *EvmWrite) UnregisterFromWorkflow(ctx context.Context, request capabilities.UnregisterFromWorkflowRequest) error {
	return nil
}
