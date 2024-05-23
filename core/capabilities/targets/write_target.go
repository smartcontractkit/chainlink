package targets

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/mitchellh/mapstructure"

	chainselectors "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/types/core"
	"github.com/smartcontractkit/chainlink-common/pkg/values"
	txmgrcommon "github.com/smartcontractkit/chainlink/v2/common/txmgr"
	"github.com/smartcontractkit/chainlink/v2/common/txmgr/types"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/config"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/txmgr"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/chains/legacyevm"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

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

	txMeta := make(map[string]interface{}) // TODO: Consider just using the TxMeta struct here and pass it to submitTransaction
	txMeta["WorkflowExecutionID"] = request.Metadata.WorkflowExecutionID

	err = cap.submitSignedTransaction(ctx, cap.chain.Config().EVM().ChainWriter(), inputs.Report, inputs.Signatures, reqConfig.Address, txMeta)
	if err != nil {
		return nil, err
	}
	// TODO: Do we want to log something here about whether the transaction was submitted successfully?
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

func (cap *EvmWrite) submitSignedTransaction(ctx context.Context, chainWriterConfig config.ChainWriter, report []byte, signatures [][]byte, toAddress string, txMeta map[string]interface{}) error {
	// construct forwarder payload
	forwardABI := evmtypes.MustGetABI(chainWriterConfig.ABI())
	calldata, err := forwardABI.Pack("report", common.HexToAddress(toAddress), report, signatures)
	if err != nil {
		return err
	}

	txMetaStruct, err := mapToStruct(txMeta)
	if err != nil {
		return err
	}

	// TODO: Turn this into config
	strategy := txmgrcommon.NewSendEveryStrategy()

	// TODO: validate the config's checker string to ensure it's a valid checker
	var checker txmgr.TransmitCheckerSpec
	if chainWriterConfig.Checker() != "" {
		checker.CheckerType = types.TransmitCheckerType(chainWriterConfig.Checker())
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

func mapToStruct(m map[string]interface{}) (*txmgr.TxMeta, error) {
	// Marshal the map to JSON
	data, err := json.Marshal(m)
	if err != nil {
		return &txmgr.TxMeta{}, err
	}

	// Unmarshal the JSON to the struct
	var result txmgr.TxMeta
	if err := json.Unmarshal(data, &result); err != nil {
		return &txmgr.TxMeta{}, err
	}

	return &result, nil
}

func (cap *EvmWrite) RegisterToWorkflow(ctx context.Context, request capabilities.RegisterToWorkflowRequest) error {
	return nil
}

func (cap *EvmWrite) UnregisterFromWorkflow(ctx context.Context, request capabilities.UnregisterFromWorkflowRequest) error {
	return nil
}
