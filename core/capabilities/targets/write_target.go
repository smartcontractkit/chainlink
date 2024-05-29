package targets

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"

	chainselectors "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/loop"
	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-common/pkg/values"
	"github.com/smartcontractkit/chainlink/v2/core/chains/legacyevm"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/forwarder"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm"
	relayevmtypes "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/types"
)

var (
	_ capabilities.ActionCapability = &WriteTarget{}
)

type WriteTarget struct {
	relayer          loop.Relayer
	cr               commontypes.ContractReader
	cw               commontypes.ChainWriter
	forwarderAddress string
	capabilities.CapabilityInfo
	lggr logger.Logger
}

func NewEvmWriteTarget(ctx context.Context, relayer loop.Relayer, chain legacyevm.Chain, lggr logger.Logger) (*WriteTarget, error) {
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

	// EVM-specific init
	config := chain.Config().EVM().ChainWriter()

	// Initialize a reader to check whether a value was already transmitted on chain
	contractReaderConfigEncoded, err := json.Marshal(relayevmtypes.ChainReaderConfig{
		Contracts: map[string]relayevmtypes.ChainContractReader{
			"forwarder": {
				ContractABI: forwarder.KeystoneForwarderABI,
				Configs: map[string]*relayevmtypes.ChainReaderDefinition{
					"getTransmitter": {
						ChainSpecificName: "getTransmitter",
					},
				},
			},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal contract reader config %v", err)
	}
	cr, err := relayer.NewContractReader(ctx, contractReaderConfigEncoded)
	if err != nil {
		return nil, err
	}
	cr.Bind(ctx, []commontypes.BoundContract{{
		Address: config.ForwarderAddress().String(),
		Name:    "forwarder",
	}})

	logger := lggr.Named("WriteTarget")

	chainWriterConfig := relayevmtypes.ChainWriterConfig{
		Contracts: map[string]relayevmtypes.ChainWriter{
			"forwarder": {
				ContractABI: forwarder.KeystoneForwarderABI,
				Configs: map[string]*relayevmtypes.ChainWriterDefinition{
					"report": {
						ChainSpecificName: "report",
						Checker:           "simulate",
						FromAddress:       config.FromAddress().Address(),
						GasLimit:          200_000,
					},
				},
			},
		},
	}
	cw := evm.NewChainWriterService(logger, chain.Client(), chain.TxManager(), chainWriterConfig)

	return &WriteTarget{
		relayer,
		cr,
		cw,
		config.ForwarderAddress().String(),
		info,
		logger,
	}, nil
}

type EvmConfig struct {
	Address string
}

func parseConfig(rawConfig *values.Map) (config EvmConfig, err error) {
	if err := rawConfig.UnwrapTo(&config); err != nil {
		return config, err
	}
	if !common.IsHexAddress(config.Address) {
		return config, fmt.Errorf("'%v' is not a valid address", config.Address)
	}
	return config, nil
}

var inputs struct {
	Report     []byte
	Signatures [][]byte
}

func success() <-chan capabilities.CapabilityResponse {
	callback := make(chan capabilities.CapabilityResponse)
	go func() {
		callback <- capabilities.CapabilityResponse{}
		close(callback)
	}()
	return callback
}

func (cap *WriteTarget) Execute(ctx context.Context, request capabilities.CapabilityRequest) (<-chan capabilities.CapabilityResponse, error) {
	cap.lggr.Debugw("Execute", "request", request)

	reqConfig, err := parseConfig(request.Config)
	if err != nil {
		return nil, err
	}
	if err = request.Inputs.UnwrapTo(&inputs); err != nil {
		return nil, err
	}

	if inputs.Report == nil {
		// We received any empty report -- this means we should skip transmission.
		cap.lggr.Debugw("Skipping empty report", "request", request)
		return success(), nil
	}

	// TODO: validate encoded report is prefixed with workflowID and executionID that match the request meta

	// Check whether value was already transmitted on chain
	queryInputs := struct {
		Receiver            string
		WorkflowExecutionID []byte
	}{
		Receiver:            reqConfig.Address,
		WorkflowExecutionID: []byte(request.Metadata.WorkflowExecutionID),
	}
	var transmitter common.Address
	if err := cap.cr.GetLatestValue(ctx, "forwarder", "getTransmitter", queryInputs, &transmitter); err != nil {
		return nil, err
	}
	if transmitter != common.HexToAddress("0x0") {
		// report already transmitted, early return
		return success(), nil
	}

	txID, err := uuid.NewUUID() // TODO(archseer): it seems odd that CW expects us to generate an ID, rather than return one
	if err != nil {
		return nil, err
	}
	args := []any{common.HexToAddress(reqConfig.Address), inputs.Report, inputs.Signatures}
	meta := commontypes.TxMeta{WorkflowExecutionID: &request.Metadata.WorkflowExecutionID}
	value := big.NewInt(0)
	if err := cap.cw.SubmitTransaction(ctx, "forwarder", "report", args, txID, cap.forwarderAddress, &meta, *value); err != nil {
		return nil, err
	}
	cap.lggr.Debugw("Transaction submitted", "request", request, "transaction", txID)
	return success(), nil
}

func (cap *WriteTarget) RegisterToWorkflow(ctx context.Context, request capabilities.RegisterToWorkflowRequest) error {
	return nil
}

func (cap *WriteTarget) UnregisterFromWorkflow(ctx context.Context, request capabilities.UnregisterFromWorkflowRequest) error {
	return nil
}
