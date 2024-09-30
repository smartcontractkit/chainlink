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
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-common/pkg/types/query/primitives"
)

var (
	_ capabilities.TargetCapability = &WriteTarget{}
)

type WriteTarget struct {
	cr               ContractValueGetter
	cw               commontypes.ChainWriter
	binding          commontypes.BoundContract
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
const ForwarderContractLogicGasCost = 100_000

type ContractValueGetter interface {
	Bind(context.Context, []commontypes.BoundContract) error
	GetLatestValue(context.Context, string, primitives.ConfidenceLevel, any, any) error
}

func NewWriteTarget(
	lggr logger.Logger,
	id string,
	cr ContractValueGetter,
	cw commontypes.ChainWriter,
	forwarderAddress string,
	txGasLimit uint64,
) *WriteTarget {
	info := capabilities.MustNewCapabilityInfo(
		id,
		capabilities.CapabilityTypeTarget,
		"Write target.",
	)

	return &WriteTarget{
		cr,
		cw,
		commontypes.BoundContract{
			Address: forwarderAddress,
			Name:    "forwarder",
		},
		forwarderAddress,
		txGasLimit - ForwarderContractLogicGasCost,
		info,
		logger.Named(lggr, "WriteTarget"),
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

func decodeReportMetadata(data []byte) (metadata ReportV1Metadata, err error) {
	if len(data) < metadata.Length() {
		return metadata, fmt.Errorf("data too short: %d bytes", len(data))
	}
	return metadata, binary.Read(bytes.NewReader(data[:metadata.Length()]), binary.BigEndian, &metadata)
}

type Config struct {
	// Address of the contract that will get the forwarded report
	Address string
	// Optional gas limit that overrides the default limit sent to the chain writer
	GasLimit *uint64
}

type Inputs struct {
	SignedReport types.SignedReport
}

type Request struct {
	Metadata capabilities.RequestMetadata
	Config   Config
	Inputs   Inputs
}

func evaluate(rawRequest capabilities.CapabilityRequest) (r Request, err error) {
	r.Metadata = rawRequest.Metadata

	if rawRequest.Config == nil {
		return r, fmt.Errorf("missing config field")
	}

	if err = rawRequest.Config.UnwrapTo(&r.Config); err != nil {
		return r, err
	}

	if !common.IsHexAddress(r.Config.Address) {
		return r, fmt.Errorf("'%v' is not a valid address", r.Config.Address)
	}

	if rawRequest.Inputs == nil {
		return r, fmt.Errorf("missing inputs field")
	}

	// required field of target's config in the workflow spec
	const signedReportField = "signed_report"
	signedReport, ok := rawRequest.Inputs.Underlying[signedReportField]
	if !ok {
		return r, fmt.Errorf("missing required field %s", signedReportField)
	}

	if err = signedReport.UnwrapTo(&r.Inputs.SignedReport); err != nil {
		return r, err
	}

	reportMetadata, err := decodeReportMetadata(r.Inputs.SignedReport.Report)
	if err != nil {
		return r, err
	}

	if reportMetadata.Version != 1 {
		return r, fmt.Errorf("unsupported report version: %d", reportMetadata.Version)
	}

	if hex.EncodeToString(reportMetadata.WorkflowExecutionID[:]) != rawRequest.Metadata.WorkflowExecutionID {
		return r, fmt.Errorf("WorkflowExecutionID in the report does not match WorkflowExecutionID in the request metadata. Report WorkflowExecutionID: %+v, request WorkflowExecutionID: %+v", reportMetadata.WorkflowExecutionID, rawRequest.Metadata.WorkflowExecutionID)
	}

	if hex.EncodeToString(reportMetadata.WorkflowOwner[:]) != rawRequest.Metadata.WorkflowOwner {
		return r, fmt.Errorf("WorkflowOwner in the report does not match WorkflowOwner in the request metadata. Report WorkflowOwner: %+v, request WorkflowOwner: %+v", reportMetadata.WorkflowOwner, rawRequest.Metadata.WorkflowOwner)
	}

	if hex.EncodeToString(reportMetadata.WorkflowName[:]) != rawRequest.Metadata.WorkflowName {
		return r, fmt.Errorf("WorkflowName in the report does not match WorkflowName in the request metadata. Report WorkflowName: %+v, request WorkflowName: %+v", reportMetadata.WorkflowName, rawRequest.Metadata.WorkflowName)
	}

	if hex.EncodeToString(reportMetadata.WorkflowCID[:]) != rawRequest.Metadata.WorkflowID {
		return r, fmt.Errorf("WorkflowID in the report does not match WorkflowID in the request metadata. Report WorkflowID: %+v, request WorkflowID: %+v", reportMetadata.WorkflowCID, rawRequest.Metadata.WorkflowID)
	}

	if !bytes.Equal(reportMetadata.ReportID[:], r.Inputs.SignedReport.ID) {
		return r, fmt.Errorf("ReportID in the report does not match ReportID in the inputs. reportMetadata.ReportID: %x, Inputs.SignedReport.ID: %x", reportMetadata.ReportID, r.Inputs.SignedReport.ID)
	}

	return r, nil
}

func (cap *WriteTarget) Execute(ctx context.Context, rawRequest capabilities.CapabilityRequest) (capabilities.CapabilityResponse, error) {
	// Bind to the contract address on the write path.
	// Bind() requires a connection to the node's RPCs and
	// cannot be run during initialization.
	if !cap.bound {
		cap.lggr.Debugw("Binding to forwarder address")
		err := cap.cr.Bind(ctx, []commontypes.BoundContract{cap.binding})
		if err != nil {
			return capabilities.CapabilityResponse{}, err
		}
		cap.bound = true
	}

	cap.lggr.Debugw("Execute", "rawRequest", rawRequest)

	request, err := evaluate(rawRequest)
	if err != nil {
		return capabilities.CapabilityResponse{}, err
	}

	rawExecutionID, err := hex.DecodeString(request.Metadata.WorkflowExecutionID)
	if err != nil {
		return capabilities.CapabilityResponse{}, err
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
	if err = cap.cr.GetLatestValue(ctx, cap.binding.ReadIdentifier("getTransmissionInfo"), primitives.Unconfirmed, queryInputs, &transmissionInfo); err != nil {
		return capabilities.CapabilityResponse{}, fmt.Errorf("failed to getTransmissionInfo latest value: %w", err)
	}

	switch {
	case transmissionInfo.State == 0: // NOT_ATTEMPTED
		cap.lggr.Infow("non-empty report - transmission not attempted - attempting to push to txmgr", "request", request, "reportLen", len(request.Inputs.SignedReport.Report), "reportContextLen", len(request.Inputs.SignedReport.Context), "nSignatures", len(request.Inputs.SignedReport.Signatures), "executionID", request.Metadata.WorkflowExecutionID)
	case transmissionInfo.State == 1: // SUCCEEDED
		cap.lggr.Infow("returning without a transmission attempt - report already onchain ", "executionID", request.Metadata.WorkflowExecutionID)
		return capabilities.CapabilityResponse{}, nil
	case transmissionInfo.State == 2: // INVALID_RECEIVER
		cap.lggr.Infow("returning without a transmission attempt - transmission already attempted, receiver was marked as invalid", "executionID", request.Metadata.WorkflowExecutionID)
		return capabilities.CapabilityResponse{}, nil
	case transmissionInfo.State == 3: // FAILED
		receiverGasMinimum := cap.receiverGasMinimum
		if request.Config.GasLimit != nil {
			receiverGasMinimum = *request.Config.GasLimit - ForwarderContractLogicGasCost
		}
		if transmissionInfo.GasLimit.Uint64() > receiverGasMinimum {
			cap.lggr.Infow("returning without a transmission attempt - transmission already attempted and failed, sufficient gas was provided", "executionID", request.Metadata.WorkflowExecutionID, "receiverGasMinimum", receiverGasMinimum, "transmissionGasLimit", transmissionInfo.GasLimit)
			return capabilities.CapabilityResponse{}, nil
		} else {
			cap.lggr.Infow("non-empty report - retrying a failed transmission - attempting to push to txmgr", "request", request, "reportLen", len(request.Inputs.SignedReport.Report), "reportContextLen", len(request.Inputs.SignedReport.Context), "nSignatures", len(request.Inputs.SignedReport.Signatures), "executionID", request.Metadata.WorkflowExecutionID, "receiverGasMinimum", receiverGasMinimum, "transmissionGasLimit", transmissionInfo.GasLimit)
		}
	default:
		return capabilities.CapabilityResponse{}, fmt.Errorf("unexpected transmission state: %v", transmissionInfo.State)
	}

	txID, err := uuid.NewUUID() // NOTE: CW expects us to generate an ID, rather than return one
	if err != nil {
		return capabilities.CapabilityResponse{}, err
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
	if request.Config.GasLimit != nil {
		meta.GasLimit = new(big.Int).SetUint64(*request.Config.GasLimit)
	}

	value := big.NewInt(0)
	if err := cap.cw.SubmitTransaction(ctx, "forwarder", "report", req, txID.String(), cap.forwarderAddress, &meta, value); err != nil {
		if !commontypes.ErrSettingTransactionGasLimitNotSupported.Is(err) {
			return capabilities.CapabilityResponse{}, fmt.Errorf("failed to submit transaction: %w", err)
		}
		meta.GasLimit = nil
		if err := cap.cw.SubmitTransaction(ctx, "forwarder", "report", req, txID.String(), cap.forwarderAddress, &meta, value); err != nil {
			return capabilities.CapabilityResponse{}, fmt.Errorf("failed to submit transaction: %w", err)
		}
	}

	cap.lggr.Debugw("Transaction submitted", "request", request, "transaction", txID)
	return capabilities.CapabilityResponse{}, nil
}

func (cap *WriteTarget) RegisterToWorkflow(ctx context.Context, request capabilities.RegisterToWorkflowRequest) error {
	return nil
}

func (cap *WriteTarget) UnregisterFromWorkflow(ctx context.Context, request capabilities.UnregisterFromWorkflowRequest) error {
	return nil
}
