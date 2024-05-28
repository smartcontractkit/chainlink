package targets

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ethereum/go-ethereum/common"

	chainselectors "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/loop"
	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-common/pkg/values"
	txmgrcommon "github.com/smartcontractkit/chainlink/v2/common/txmgr"
	evmconfig "github.com/smartcontractkit/chainlink/v2/core/chains/evm/config"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/txmgr"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/chains/legacyevm"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/forwarder"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	relayevmtypes "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/types"
)

var (
	_ capabilities.ActionCapability = &WriteTarget{}
)

type WriteTarget struct {
	relayer loop.Relayer
	config  evmconfig.ChainScopedConfig
	cr      commontypes.ContractReader
	cw      EvmChainWriter
	capabilities.CapabilityInfo
	lggr logger.Logger
}

func NewWriteTarget(ctx context.Context, relayer loop.Relayer, chain legacyevm.Chain, lggr logger.Logger) (*WriteTarget, error) {
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
	config := chain.Config()

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

	cw := EvmChainWriter{
		chain,
	}

	return &WriteTarget{
		relayer,
		config,
		cr,
		cw,
		info,
		lggr.Named("WriteTarget"),
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

type EvmChainWriter struct {
	chain legacyevm.Chain
}

func (cw *EvmChainWriter) CreateTransaction(ctx context.Context, reqConfig EvmConfig) (tx txmgr.Tx, err error) {
	txm := cw.chain.TxManager()
	config := cw.chain.Config().EVM().ChainWriter()

	// construct forwarder payload
	calldata, err := forwardABI.Pack("report", common.HexToAddress(reqConfig.Address), inputs.Report, inputs.Signatures)
	if err != nil {
		return tx, err
	}

	txMeta := &txmgr.TxMeta{
		// FwdrDestAddress could also be set for better logging but it's used for various purposes around Operator Forwarders
		// WorkflowExecutionID: &request.Metadata.WorkflowExecutionID, // TODO: remove?
	}
	req := txmgr.TxRequest{
		FromAddress:    config.FromAddress().Address(),
		ToAddress:      config.ForwarderAddress().Address(),
		EncodedPayload: calldata,
		FeeLimit:       uint64(defaultGasLimit),
		Meta:           txMeta,
		Strategy:       txmgrcommon.NewSendEveryStrategy(),
		Checker: txmgr.TransmitCheckerSpec{
			CheckerType: txmgr.TransmitCheckerTypeSimulate,
		},
	}
	return txm.CreateTransaction(ctx, req)
}

var forwardABI = evmtypes.MustGetABI(forwarder.KeystoneForwarderMetaData.ABI)

const defaultGasLimit = 200000

func (cap *WriteTarget) Execute(ctx context.Context, request capabilities.CapabilityRequest) (<-chan capabilities.CapabilityResponse, error) {
	cap.lggr.Debugw("Execute", "request", request)
	// TODO: idempotency

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

	config := cap.config.EVM().ChainWriter()

	// Check whether value was already transmitted on chain
	cap.cr.Bind(ctx, []commontypes.BoundContract{{
		Address: config.ForwarderAddress().String(),
		Name:    "forwarder",
	}})
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

	tx, err := cap.cw.CreateTransaction(ctx, reqConfig)
	if err != nil {
		return nil, err
	}
	cap.lggr.Debugw("Transaction submitted", "request", request, "transaction", tx)

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

func (cap *WriteTarget) RegisterToWorkflow(ctx context.Context, request capabilities.RegisterToWorkflowRequest) error {
	return nil
}

func (cap *WriteTarget) UnregisterFromWorkflow(ctx context.Context, request capabilities.UnregisterFromWorkflowRequest) error {
	return nil
}
