package targets

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/consensus/ocr3/types"
	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-common/pkg/types/query/primitives"
	"github.com/smartcontractkit/chainlink-common/pkg/values"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

var (
	_ capabilities.ActionCapability = &WriteTarget{}
)

// required field of target's config in the workflow spec
const signedReportField = "signed_report"

type WriteTarget struct {
	cr               commontypes.ContractReader
	cw               commontypes.ChainWriter
	forwarderAddress string
	// The minimum amount of gas that the receiver contract must get to process the forwarder report
	receiverGasMinimum uint64
	capabilities.CapabilityInfo

	lggr logger.Logger

	bound bool
}

type TransmissionInfo struct {
	GasLimit        *big.Int
	InvalidReceiver bool
	State           uint8
	Success         bool
	TransmissionId  [32]byte
	Transmitter     common.Address
}

// The gas cost of the forwarder contract logic, including state updates and event emission.
// This is a rough estimate and should be updated if the forwarder contract logic changes.
// TODO: Make this part of the on-chain capability configuration
const FORWARDER_CONTRACT_LOGIC_GAS_COST = 100_000

func NewWriteTarget(lggr logger.Logger, id string, cr commontypes.ContractReader, cw commontypes.ChainWriter, forwarderAddress string, txGasLimit uint64) *WriteTarget {
	info := capabilities.MustNewCapabilityInfo(
		id,
		capabilities.CapabilityTypeTarget,
		"Write target.",
	)

	logger := lggr.Named("WriteTarget")

	return &WriteTarget{
		cr,
		cw,
		forwarderAddress,
		txGasLimit - FORWARDER_CONTRACT_LOGIC_GAS_COST,
		info,
		logger,
		false,
	}
}

type EvmConfig struct {
	Address string
}

func parseConfig(rawConfig *values.Map) (config EvmConfig, err error) {
	if rawConfig == nil {
		return config, fmt.Errorf("missing config field")
	}

	if err := rawConfig.UnwrapTo(&config); err != nil {
		return config, err
	}
	if !common.IsHexAddress(config.Address) {
		return config, fmt.Errorf("'%v' is not a valid address", config.Address)
	}
	return config, nil
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
	// Bind to the contract address on the write path.
	// Bind() requires a connection to the node's RPCs and
	// cannot be run during initialization.
	if !cap.bound {
		cap.lggr.Debugw("Binding to forwarder address")
		err := cap.cr.Bind(ctx, []commontypes.BoundContract{{
			Address: cap.forwarderAddress,
			Name:    "forwarder",
		}})
		if err != nil {
			return nil, err
		}
		cap.bound = true
	}

	cap.lggr.Debugw("Execute", "request", request)

	reqConfig, err := parseConfig(request.Config)
	if err != nil {
		return nil, err
	}

	if request.Inputs == nil {
		return nil, fmt.Errorf("missing inputs field")
	}

	signedReport, ok := request.Inputs.Underlying[signedReportField]
	if !ok {
		return nil, fmt.Errorf("missing required field %s", signedReportField)
	}

	inputs := types.SignedReport{}
	if err = signedReport.UnwrapTo(&inputs); err != nil {
		return nil, err
	}

	if len(inputs.Report) == 0 {
		// We received any empty report -- this means we should skip transmission.
		cap.lggr.Debugw("Skipping empty report", "request", request)
		return success(), nil
	}
	// TODO: validate encoded report is prefixed with workflowID and executionID that match the request meta

	rawExecutionID, err := hex.DecodeString(request.Metadata.WorkflowExecutionID)
	if err != nil {
		return nil, err
	}
	// Check whether value was already transmitted on chain
	queryInputs := struct {
		Receiver            string
		WorkflowExecutionID []byte
		ReportId            []byte
	}{
		Receiver:            reqConfig.Address,
		WorkflowExecutionID: rawExecutionID,
		ReportId:            inputs.ID,
	}
	var transmissionInfo TransmissionInfo
	if err = cap.cr.GetLatestValue(ctx, "forwarder", "getTransmissionInfo", primitives.Unconfirmed, queryInputs, &transmissionInfo); err != nil {
		return nil, fmt.Errorf("failed to getTransmissionInfo latest value: %w", err)
	}

	switch {
	case transmissionInfo.State == 0: // NOT_ATTEMPTED
		cap.lggr.Infow("non-empty report - tranasmission not attempted - attempting to push to txmgr", "request", request, "reportLen", len(inputs.Report), "reportContextLen", len(inputs.Context), "nSignatures", len(inputs.Signatures), "executionID", request.Metadata.WorkflowExecutionID)
	case transmissionInfo.State == 1: // SUCCEEDED
		cap.lggr.Infow("returning without a tranmission attempt - report already onchain ", "executionID", request.Metadata.WorkflowExecutionID)
		return success(), nil
	case transmissionInfo.State == 2: // INVALID_RECEIVER
		cap.lggr.Infow("returning without a tranmission attempt - transmission already attempted, receiver was marked as invalid", "executionID", request.Metadata.WorkflowExecutionID)
		return success(), nil
	case transmissionInfo.State == 3: // FAILED
		if transmissionInfo.GasLimit.Uint64() > cap.receiverGasMinimum {
			cap.lggr.Infow("returning without a tranmission attempt - transmission already attempted and failed, sufficient gas was provided", "executionID", request.Metadata.WorkflowExecutionID, "receiverGasMinimum", cap.receiverGasMinimum, "transmissionGasLimit", transmissionInfo.GasLimit)
			return success(), nil
		} else {
			cap.lggr.Infow("non-empty report - retrying a failed transmission - attempting to push to txmgr", "request", request, "reportLen", len(inputs.Report), "reportContextLen", len(inputs.Context), "nSignatures", len(inputs.Signatures), "executionID", request.Metadata.WorkflowExecutionID, "receiverGasMinimum", cap.receiverGasMinimum, "transmissionGasLimit", transmissionInfo.GasLimit)
		}
	default:
		return nil, fmt.Errorf("unexpected transmission state: %v", transmissionInfo.State)
	}

	txID, err := uuid.NewUUID() // NOTE: CW expects us to generate an ID, rather than return one
	if err != nil {
		return nil, err
	}

	// Note: The codec that ChainWriter uses to encode the parameters for the contract ABI cannot handle
	// `nil` values, including for slices. Until the bug is fixed we need to ensure that there are no
	// `nil` values passed in the request.
	req := struct {
		Receiver      string
		RawReport     []byte
		ReportContext []byte
		Signatures    [][]byte
	}{reqConfig.Address, inputs.Report, inputs.Context, inputs.Signatures}

	if req.RawReport == nil {
		req.RawReport = make([]byte, 0)
	}

	if req.ReportContext == nil {
		req.ReportContext = make([]byte, 0)
	}

	if req.Signatures == nil {
		req.Signatures = make([][]byte, 0)
	}
	cap.lggr.Debugw("Transaction raw report", "report", hex.EncodeToString(req.RawReport))

	meta := commontypes.TxMeta{WorkflowExecutionID: &request.Metadata.WorkflowExecutionID}
	value := big.NewInt(0)
	if err := cap.cw.SubmitTransaction(ctx, "forwarder", "report", req, txID.String(), cap.forwarderAddress, &meta, value); err != nil {
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
