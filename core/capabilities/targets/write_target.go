package targets

import (
	"context"
	"fmt"

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
	"github.com/smartcontractkit/chainlink/v2/core/services/pipeline"
)

var forwardABI = evmtypes.MustGetABI(forwarder.KeystoneForwarderMetaData.ABI)

func InitializeWrite(registry commontypes.CapabilitiesRegistry, legacyEVMChains legacyevm.LegacyChainContainer) error {
	for _, chain := range legacyEVMChains.Slice() {
		capability := NewEvmWrite(chain)
		if err := registry.Add(context.TODO(), capability); err != nil {
			return err
		}
	}
	return nil
}

var (
	_ capabilities.ActionCapability = &EvmWrite{}
)

type EvmWrite struct {
	chain legacyevm.Chain
	capabilities.CapabilityInfo
}

func NewEvmWrite(chain legacyevm.Chain) *EvmWrite {
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

	// NOTE: without having full ABI it's actually impossible to support function overloading
	selector, err := abiutil.ParseSelector(rawSelector)
	if err != nil {
		return nil, err
	}

	return selector.Pack(args...)
}

func (cap *EvmWrite) Execute(ctx context.Context, callback chan<- capabilities.CapabilityResponse, request capabilities.CapabilityRequest) error {
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

	// unlimited gas in the MVP demo
	gasLimit := 0
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
		FeeLimit:       uint32(gasLimit),
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
