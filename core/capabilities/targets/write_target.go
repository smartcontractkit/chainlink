package targets

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/mitchellh/mapstructure"

	chainselectors "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/types/core"
	"github.com/smartcontractkit/chainlink-common/pkg/values"
	txmgrcommon "github.com/smartcontractkit/chainlink/v2/common/txmgr"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/config"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/txmgr"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/chains/legacyevm"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/forwarder"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

var forwardABI = evmtypes.MustGetABI(forwarder.KeystoneForwarderMetaData.ABI)

func InitializeWrite(registry core.CapabilitiesRegistry, legacyEVMChains legacyevm.LegacyChainContainer, lggr logger.Logger) error {
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
		nil,
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

func (cap *EvmWrite) Execute(ctx context.Context, request capabilities.CapabilityRequest) (<-chan capabilities.CapabilityResponse, error) {
	cap.lggr.Debugw("Execute", "request", request)
	// TODO: idempotency

	reqConfig, err := parseConfig(request.Config)
	if err != nil {
		return nil, err
	}

	var inputs struct {
		Report     []byte
		Signatures [][]byte
	}
	if err = request.Inputs.UnwrapTo(&inputs); err != nil {
		return nil, err
	}

	if inputs.Report == nil {
		// We received any empty report -- this means we should skip transmission.
		cap.lggr.Debugw("Skipping empty report", "request", request)
		callback := make(chan capabilities.CapabilityResponse)
		go func() {
			// TODO: cast tx.Error to Err (or Value to Value?)
			callback <- capabilities.CapabilityResponse{
				Value: nil,
				Err:   nil,
			}
			close(callback)
		}()
		return callback, nil
	}

	// TODO: validate encoded report is prefixed with workflowID and executionID that match the request meta

	txMeta := make(map[string]string)
	txMeta["WorkflowExecutionID"] = request.Metadata.WorkflowExecutionID

	err = cap.submitTransaction(ctx, cap.chain.Config().EVM().ChainWriter(), inputs.Report, inputs.Signatures, reqConfig.Address, txMeta)
	if err != nil {
		return nil, err
	}
	//cap.lggr.Debugw("Transaction submitted", "request", request, "transaction", tx)

	callback := make(chan capabilities.CapabilityResponse)
	go func() {
		// TODO: cast tx.Error to Err (or Value to Value?)
		callback <- capabilities.CapabilityResponse{
			Value: nil,
			Err:   nil,
		}
		close(callback)
	}()
	return callback, nil
}

func (cap *EvmWrite) submitTransaction(ctx context.Context, chainWriterConfig config.ChainWriter, report []byte, signatures [][]byte, toAddress string, txMeta map[string]string) error {
	// construct forwarder payload
	// TODO: we can't assume we have the ABI locally here, we need to get it from the chain
	calldata, err := forwardABI.Pack("report", common.HexToAddress(toAddress), report, signatures)
	if err != nil {
		return err
	}

	workflowExecutionID := txMeta["WorkflowExecutionID"]

	txMetaStruct := &txmgr.TxMeta{
		WorkflowExecutionID: &workflowExecutionID,
	}
	// TODO: Turn this into config
	strategy := txmgrcommon.NewSendEveryStrategy()

	// TODO: Turn this into config
	checker := txmgr.TransmitCheckerSpec{
		CheckerType: txmgr.TransmitCheckerTypeSimulate,
	}
	req := txmgr.TxRequest{
		FromAddress:    chainWriterConfig.FromAddress().Address(),
		ToAddress:      chainWriterConfig.ForwarderAddress().Address(),
		EncodedPayload: calldata,
		FeeLimit:       chainWriterConfig.GasLimit(),
		Meta:           txMetaStruct,
		Strategy:       strategy,
		Checker:        checker,
	}
	tx, err := cap.chain.TxManager().CreateTransaction(ctx, req)
	if err != nil {
		return err
	}
	cap.lggr.Debugw("Transaction submitted", "transaction", tx)
	return nil
}

func (cap *EvmWrite) RegisterToWorkflow(ctx context.Context, request capabilities.RegisterToWorkflowRequest) error {
	return nil
}

func (cap *EvmWrite) UnregisterFromWorkflow(ctx context.Context, request capabilities.UnregisterFromWorkflowRequest) error {
	return nil
}
