package targets

import (
	"bytes"
	"context"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/consensus/ocr3/types"
	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-common/pkg/types/query/primitives"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

var (
	_ capabilities.ActionCapability = &WriteTarget{}
)

// required field of target's config in the workflow spec
const signedReportField = "signed_report"

// The gas cost of the forwarder contract logic, including state updates and event emission.
// This is a rough estimate and should be updated if the forwarder contract logic changes.
// TODO: Make this part of the on-chain capability configuration
const FORWARDER_CONTRACT_LOGIC_GAS_COST = 100_000

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

// Note: This should be a shared type that the OCR3 package validates as well
type ReportV1Metadata struct {
	Version             uint8
	WorkflowExecutionID [32]byte
	Timestamp           uint32
	DonID               uint32
	DonConfigVersion    uint32
	WorkflowCID         [32]byte
	WorkflowName        [10]byte
	WorkflowOwner       [20]byte
	ReportID            [2]byte
}

func (rm ReportV1Metadata) Encode() ([]byte, error) {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.BigEndian, rm)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (rm ReportV1Metadata) Length() int {
	bytes, err := rm.Encode()
	if err != nil {
		return 0
	}
	return len(bytes)
}

func decodeReportMetadata(reportPayload []byte) (ReportV1Metadata, error) {
	var metadata ReportV1Metadata

	if len(reportPayload) < metadata.Length() {
		return metadata, fmt.Errorf("reportPayload too short: %d bytes", len(reportPayload))
	}

	reportPayload = reportPayload[:metadata.Length()]

	buffer := bytes.NewReader(reportPayload)

	// Reading Version (1 byte)
	if err := binary.Read(buffer, binary.BigEndian, &metadata.Version); err != nil {
		return metadata, err
	}

	// Reading WorkflowExecutionID (32 bytes)
	if err := binary.Read(buffer, binary.BigEndian, &metadata.WorkflowExecutionID); err != nil {
		return metadata, err
	}

	// Reading Timestamp (4 bytes)
	if err := binary.Read(buffer, binary.BigEndian, &metadata.Timestamp); err != nil {
		return metadata, err
	}

	// Reading DonID (4 bytes)
	if err := binary.Read(buffer, binary.BigEndian, &metadata.DonID); err != nil {
		return metadata, err
	}

	// Reading DonConfigVersion (4 bytes)
	if err := binary.Read(buffer, binary.BigEndian, &metadata.DonConfigVersion); err != nil {
		return metadata, err
	}

	// Reading WorkflowCID (32 bytes)
	if err := binary.Read(buffer, binary.BigEndian, &metadata.WorkflowCID); err != nil {
		return metadata, err
	}

	// Reading WorkflowName (10 bytes)
	if err := binary.Read(buffer, binary.BigEndian, &metadata.WorkflowName); err != nil {
		return metadata, err
	}

	// Reading WorkflowOwner (20 bytes)
	if err := binary.Read(buffer, binary.BigEndian, &metadata.WorkflowOwner); err != nil {
		return metadata, err
	}

	// Reading ReportID (2 bytes)
	if err := binary.Read(buffer, binary.BigEndian, &metadata.ReportID); err != nil {
		return metadata, err
	}

	return metadata, nil
}

type Config struct {
	// Address of the contract that will get the forwarded report
	Address string
}

type Inputs struct {
	SignedReport types.SignedReport
}

type Request struct {
	Metadata capabilities.RequestMetadata
	Config   Config
	Inputs   Inputs
}

func evaluate(request capabilities.CapabilityRequest) (*Request, error) {
	if request.Config == nil {
		return nil, fmt.Errorf("missing config field")
	}

	// Validate and unwrap Config
	config := Config{}
	if err := request.Config.UnwrapTo(&config); err != nil {
		return nil, err
	}

	if !common.IsHexAddress(config.Address) {
		return nil, fmt.Errorf("'%v' is not a valid address", config.Address)
	}

	// Validate and unwrap Inputs
	if request.Inputs == nil {
		return nil, fmt.Errorf("missing inputs field")
	}

	signedReport, ok := request.Inputs.Underlying[signedReportField]
	if !ok {
		return nil, fmt.Errorf("missing required field %s", signedReportField)
	}

	report := types.SignedReport{}
	if err := signedReport.UnwrapTo(&report); err != nil {
		return nil, err
	}

	// TODO: Is this properly EVM decoding? In the error, should explain that it
	// might be due to different report encoding.
	reportMetadata, err := decodeReportMetadata(report.Report)
	if err != nil {
		return nil, err
	}

	if reportMetadata.Version != 1 {
		return nil, fmt.Errorf("unsupported report version: %d", reportMetadata.Version)
	}

	if hex.EncodeToString(reportMetadata.WorkflowExecutionID[:]) != request.Metadata.WorkflowExecutionID ||
		hex.EncodeToString(reportMetadata.WorkflowOwner[:]) != request.Metadata.WorkflowOwner ||
		hex.EncodeToString(reportMetadata.WorkflowName[:]) != request.Metadata.WorkflowName ||
		hex.EncodeToString(reportMetadata.WorkflowCID[:]) != request.Metadata.WorkflowID {
		return nil, fmt.Errorf("report metadata does not match request metadata. reportMetadata: %+v, requestMetadata: %+v", reportMetadata, request.Metadata)
	}

	return &Request{
		Metadata: request.Metadata,
		Config:   config,
		Inputs: Inputs{
			SignedReport: report,
		},
	}, nil
}

func success() <-chan capabilities.CapabilityResponse {
	callback := make(chan capabilities.CapabilityResponse)
	go func() {
		callback <- capabilities.CapabilityResponse{}
		close(callback)
	}()
	return callback
}

func (cap *WriteTarget) Execute(ctx context.Context, rawRequest capabilities.CapabilityRequest) (<-chan capabilities.CapabilityResponse, error) {
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

	cap.lggr.Debugw("Execute", "rawRequest", rawRequest)

	request, err := evaluate(rawRequest)
	if err != nil {
		return nil, err
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
		Receiver:            request.Config.Address,
		WorkflowExecutionID: rawExecutionID,
		ReportId:            request.Inputs.SignedReport.ID,
	}
	var transmissionInfo TransmissionInfo
	if err = cap.cr.GetLatestValue(ctx, "forwarder", "getTransmissionInfo", primitives.Unconfirmed, queryInputs, &transmissionInfo); err != nil {
		return nil, fmt.Errorf("failed to getTransmissionInfo latest value: %w", err)
	}

	switch {
	case transmissionInfo.State == 0: // NOT_ATTEMPTED
		cap.lggr.Infow("non-empty report - tranasmission not attempted - attempting to push to txmgr", "request", request, "reportLen", len(request.Inputs.SignedReport.Report), "reportContextLen", len(request.Inputs.SignedReport.Context), "nSignatures", len(request.Inputs.SignedReport.Signatures), "executionID", request.Metadata.WorkflowExecutionID)
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
			cap.lggr.Infow("non-empty report - retrying a failed transmission - attempting to push to txmgr", "request", request, "reportLen", len(request.Inputs.SignedReport.Report), "reportContextLen", len(request.Inputs.SignedReport.Context), "nSignatures", len(request.Inputs.SignedReport.Signatures), "executionID", request.Metadata.WorkflowExecutionID, "receiverGasMinimum", cap.receiverGasMinimum, "transmissionGasLimit", transmissionInfo.GasLimit)
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
	}{request.Config.Address, request.Inputs.SignedReport.Report, request.Inputs.SignedReport.Context, request.Inputs.SignedReport.Signatures}

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
