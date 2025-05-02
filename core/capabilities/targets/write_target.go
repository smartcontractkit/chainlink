package targets

import (
	"bytes"
	"context"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/consensus/ocr3/types"
	"github.com/smartcontractkit/chainlink-common/pkg/custmsg"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-common/pkg/types/query/primitives"

	"github.com/smartcontractkit/chainlink/v2/core/platform"
)

var (
	_           capabilities.ExecutableCapability = &WriteTarget{}
	ErrTxFailed                                   = errors.New("submitted transaction failed")
)

const transactionStatusCheckInterval = 2 * time.Second

type WriteTarget struct {
	cr               ContractValueGetter
	cw               commontypes.ContractWriter
	binding          commontypes.BoundContract
	forwarderAddress string
	// The minimum amount of gas that the receiver contract must get to process the forwarder report
	receiverGasMinimum uint64
	capabilities.CapabilityInfo

	emitter custmsg.MessageEmitter
	lggr    logger.Logger

	bound bool
}

const (
	TransmissionStateNotAttempted uint8 = iota
	TransmissionStateSucceeded
	TransmissionStateInvalidReceiver
	TransmissionStateFailed
)

type TransmissionInfo struct {
	GasLimit        *big.Int
	InvalidReceiver bool
	State           uint8
	Success         bool
	TransmissionID  [32]byte
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
	cw commontypes.ContractWriter,
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
		custmsg.NewLabeler(),
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
		return r, errors.New("missing config field")
	}

	if err = rawRequest.Config.UnwrapTo(&r.Config); err != nil {
		return r, err
	}

	if !common.IsHexAddress(r.Config.Address) {
		return r, fmt.Errorf("'%v' is not a valid address", r.Config.Address)
	}

	if rawRequest.Inputs == nil {
		return r, errors.New("missing inputs field")
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
		return r, fmt.Errorf("WorkflowExecutionID in the report does not match WorkflowExecutionID in the request metadata. Report WorkflowExecutionID: %+v, request WorkflowExecutionID: %+v", hex.EncodeToString(reportMetadata.WorkflowExecutionID[:]), rawRequest.Metadata.WorkflowExecutionID)
	}

	// case-insensitive verification of the owner address (so that a check-summed address matches its non-checksummed version).
	if !strings.EqualFold(hex.EncodeToString(reportMetadata.WorkflowOwner[:]), rawRequest.Metadata.WorkflowOwner) {
		return r, fmt.Errorf("WorkflowOwner in the report does not match WorkflowOwner in the request metadata. Report WorkflowOwner: %+v, request WorkflowOwner: %+v", hex.EncodeToString(reportMetadata.WorkflowOwner[:]), rawRequest.Metadata.WorkflowOwner)
	}

	// workflowNames are padded to 10bytes
	decodedName, err := hex.DecodeString(rawRequest.Metadata.WorkflowName)
	if err != nil {
		return r, err
	}
	var workflowName [10]byte
	copy(workflowName[:], decodedName)
	if !bytes.Equal(reportMetadata.WorkflowName[:], workflowName[:]) {
		return r, fmt.Errorf("WorkflowName in the report does not match WorkflowName in the request metadata. Report WorkflowName: %+v, request WorkflowName: %+v", hex.EncodeToString(reportMetadata.WorkflowName[:]), hex.EncodeToString(workflowName[:]))
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
	transmissionInfo, err := cap.getTransmissionInfo(ctx, request, rawExecutionID)
	if err != nil {
		return capabilities.CapabilityResponse{}, err
	}

	switch {
	case transmissionInfo.State == TransmissionStateNotAttempted:
		cap.lggr.Infow("non-empty report - transmission not attempted - attempting to push to txmgr", "request", request, "reportLen", len(request.Inputs.SignedReport.Report), "reportContextLen", len(request.Inputs.SignedReport.Context), "nSignatures", len(request.Inputs.SignedReport.Signatures), "executionID", request.Metadata.WorkflowExecutionID)
	case transmissionInfo.State == TransmissionStateSucceeded:
		cap.lggr.Infow("returning without a transmission attempt - report already onchain ", "executionID", request.Metadata.WorkflowExecutionID)
		return capabilities.CapabilityResponse{}, nil
	case transmissionInfo.State == TransmissionStateInvalidReceiver:
		cap.lggr.Infow("returning without a transmission attempt - transmission already attempted, receiver was marked as invalid", "executionID", request.Metadata.WorkflowExecutionID)
		return capabilities.CapabilityResponse{}, ErrTxFailed
	case transmissionInfo.State == TransmissionStateFailed:
		receiverGasMinimum := cap.receiverGasMinimum
		if request.Config.GasLimit != nil {
			receiverGasMinimum = *request.Config.GasLimit - ForwarderContractLogicGasCost
		}
		if transmissionInfo.GasLimit.Uint64() > receiverGasMinimum {
			cap.lggr.Infow("returning without a transmission attempt - transmission already attempted and failed, sufficient gas was provided", "executionID", request.Metadata.WorkflowExecutionID, "receiverGasMinimum", receiverGasMinimum, "transmissionGasLimit", transmissionInfo.GasLimit)
			return capabilities.CapabilityResponse{}, ErrTxFailed
		}
		cap.lggr.Infow("non-empty report - retrying a failed transmission - attempting to push to txmgr", "request", request, "reportLen", len(request.Inputs.SignedReport.Report), "reportContextLen", len(request.Inputs.SignedReport.Context), "nSignatures", len(request.Inputs.SignedReport.Signatures), "executionID", request.Metadata.WorkflowExecutionID, "receiverGasMinimum", receiverGasMinimum, "transmissionGasLimit", transmissionInfo.GasLimit)

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

	tick := time.NewTicker(transactionStatusCheckInterval)
	defer tick.Stop()
	for {
		select {
		case <-ctx.Done():
			return capabilities.CapabilityResponse{}, nil
		case <-tick.C:
			txStatus, err := cap.cw.GetTransactionStatus(ctx, txID.String())
			if err != nil {
				cap.lggr.Errorw("Failed to get transaction status", "request", request, "transaction", txID, "err", err)
				continue
			}
			switch txStatus {
			case commontypes.Pending:
				cap.lggr.Debugw("Transaction pending, retrying...", "request", request, "transaction", txID)
			// TxStatus Unconfirmed actually means "Confirmed" for the transaction manager, i.e. the transaction has landed on chain but isn't finalized.
			case commontypes.Unconfirmed:
				transmissionInfo, err = cap.getTransmissionInfo(ctx, request, rawExecutionID)
				if err != nil {
					return capabilities.CapabilityResponse{}, err
				}

				// This is counterintuitive, but the tx manager is currently returning unconfirmed whenever the tx is confirmed
				// current implementation here: https://github.com/smartcontractkit/chainlink-framework/blob/main/chains/txmgr/txmgr.go#L697
				// so we need to check if we were able to write to the consumer contract to determine if the transaction was successful
				switch transmissionInfo.State {
				case TransmissionStateSucceeded:
					cap.lggr.Debugw("Transaction confirmed", "request", request, "transaction", txID)
					return capabilities.CapabilityResponse{}, nil
				case TransmissionStateFailed, TransmissionStateInvalidReceiver:
					cap.lggr.Errorw("Transaction written to the forwarder, but failed to be written to the consumer contract", "request", request, "transaction", txID, "transmissionState", transmissionInfo.State)
					msg := "transaction written to the forwarder, but failed to be written to the consumer contract, transaction ID: " + txID.String()
					err = cap.emitter.With(
						platform.KeyWorkflowID, request.Metadata.WorkflowID,
						platform.KeyWorkflowName, request.Metadata.DecodedWorkflowName,
						platform.KeyWorkflowOwner, request.Metadata.WorkflowOwner,
						platform.KeyWorkflowExecutionID, request.Metadata.WorkflowExecutionID,
					).Emit(ctx, msg)
					if err != nil {
						cap.lggr.Errorf("failed to send custom message with msg: %s, err: %v", msg, err)
					}
					return capabilities.CapabilityResponse{}, ErrTxFailed
				}
				// TransmissionStateNotAttempted is not expected here, but we'll log it just in case
				cap.lggr.Debugw("Transaction confirmed but transmission not attempted, this should never happen", "request", request, "transaction", txID)
				return capabilities.CapabilityResponse{}, errors.New("transmission not attempted")

			case commontypes.Finalized:
				cap.lggr.Debugw("Transaction finalized", "request", request, "transaction", txID)
				return capabilities.CapabilityResponse{}, nil
			case commontypes.Failed, commontypes.Fatal:
				cap.lggr.Error("Transaction failed", "request", request, "transaction", txID)

				msg := "transaction failed to be written to the forwarder, transaction ID: " + txID.String()
				err = cap.emitter.With(
					platform.KeyWorkflowID, request.Metadata.WorkflowID,
					platform.KeyWorkflowName, request.Metadata.DecodedWorkflowName,
					platform.KeyWorkflowOwner, request.Metadata.WorkflowOwner,
					platform.KeyWorkflowExecutionID, request.Metadata.WorkflowExecutionID,
				).Emit(ctx, msg)
				if err != nil {
					cap.lggr.Errorf("failed to send custom message with msg: %s, err: %v", msg, err)
				}
				return capabilities.CapabilityResponse{}, ErrTxFailed
			default:
				cap.lggr.Debugw("Unexpected transaction status", "request", request, "transaction", txID, "status", txStatus)
			}
		}
	}
}

func (cap *WriteTarget) RegisterToWorkflow(ctx context.Context, request capabilities.RegisterToWorkflowRequest) error {
	return nil
}

func (cap *WriteTarget) UnregisterFromWorkflow(ctx context.Context, request capabilities.UnregisterFromWorkflowRequest) error {
	return nil
}

func (cap *WriteTarget) getTransmissionInfo(ctx context.Context, request Request, rawExecutionID []byte) (*TransmissionInfo, error) {
	queryInputs := struct {
		Receiver            string
		WorkflowExecutionID []byte
		ReportID            []byte
	}{
		Receiver:            request.Config.Address,
		WorkflowExecutionID: rawExecutionID,
		ReportID:            request.Inputs.SignedReport.ID,
	}
	var transmissionInfo TransmissionInfo
	if err := cap.cr.GetLatestValue(ctx, cap.binding.ReadIdentifier("getTransmissionInfo"), primitives.Unconfirmed, queryInputs, &transmissionInfo); err != nil {
		return nil, fmt.Errorf("failed to getTransmissionInfo latest value: %w", err)
	}
	return &transmissionInfo, nil
}
