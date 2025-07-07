package evm

import (
	"context"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"sync/atomic"

	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-common/pkg/types/query/primitives"

	evmtxmgr "github.com/smartcontractkit/chainlink-evm/pkg/txmgr"
	"github.com/smartcontractkit/chainlink-framework/capabilities/writetarget"
)

var (
	_           writetarget.TargetStrategy = &evmTargetStrategy{}
	ErrTxFailed                            = errors.New("submitted transaction failed")
)

type evmTargetStrategy struct {
	cr        commontypes.ContractReader
	cw        commontypes.ContractWriter
	txm       evmtxmgr.TxManager
	lggr      logger.Logger
	forwarder string

	receiverGasMinimum uint64
	binding            commontypes.BoundContract
	bound              *atomic.Bool
}

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

const (
	TransmissionStateNotAttempted uint8 = iota
	TransmissionStateSucceeded
	TransmissionStateInvalidReceiver
	TransmissionStateFailed
)

type Config struct {
	// Address of the contract that will get the forwarded report
	Address string
	// Optional gas limit that overrides the default limit sent to the chain writer
	GasLimit *uint64
}

func NewEVMTargetStrategy(cr commontypes.ContractReader, cw commontypes.ContractWriter, txm evmtxmgr.TxManager, forwarder string, gasLimitDefault uint64, lggr logger.Logger) (*evmTargetStrategy, error) {
	if gasLimitDefault < ForwarderContractLogicGasCost {
		return nil, fmt.Errorf("default gas limit '%d' is lower than forwarder estimate '%d'", gasLimitDefault, ForwarderContractLogicGasCost)
	}
	bound := atomic.Bool{}
	return &evmTargetStrategy{
		cr:                 cr,
		cw:                 cw,
		txm:                txm,
		lggr:               lggr,
		forwarder:          forwarder,
		receiverGasMinimum: gasLimitDefault - ForwarderContractLogicGasCost,
		binding: commontypes.BoundContract{
			Address: forwarder,
			Name:    "forwarder",
		},
		bound: &bound,
	}, nil
}

func (t *evmTargetStrategy) QueryTransmissionState(ctx context.Context, reportID uint16, request capabilities.CapabilityRequest) (*writetarget.TransmissionState, error) {
	rawExecutionID, err := hex.DecodeString(request.Metadata.WorkflowExecutionID)
	if err != nil {
		return nil, fmt.Errorf("failed to decode workflowExecutionID: %w", err)
	}

	b := make([]byte, 2)
	binary.BigEndian.PutUint16(b, reportID)

	if !t.bound.Load() {
		t.lggr.Debugw("Binding to forwarder address")
		err = t.cr.Bind(ctx, []commontypes.BoundContract{t.binding})
		if err != nil {
			return nil, err
		}
		t.bound.Store(true)
	}

	r, err := getEVMRequest(request)
	if err != nil {
		return nil, fmt.Errorf("failed to getEVMRequest: %w", err)
	}
	receiver := r.Config.Address

	queryInputs := struct {
		Receiver            string
		WorkflowExecutionID []byte
		ReportID            []byte
	}{
		Receiver:            receiver,
		WorkflowExecutionID: rawExecutionID,
		ReportID:            b,
	}

	var transmissionInfo TransmissionInfo
	if err := t.cr.GetLatestValue(ctx, t.binding.ReadIdentifier("getTransmissionInfo"), primitives.Unconfirmed, queryInputs, &transmissionInfo); err != nil {
		return nil, fmt.Errorf("failed to getTransmissionInfo latest value: %w", err)
	}

	// TODO: Want to confirm these states are correct - invalid receiver and failed with sufficient gas are fatal.
	switch transmissionInfo.State {
	case TransmissionStateNotAttempted:
		t.lggr.Infow("non-empty report - transmission not attempted", "request", request, "reportLen", len(r.Inputs.SignedReport.Report), "reportContextLen", len(r.Inputs.SignedReport.Context), "nSignatures", len(r.Inputs.SignedReport.Signatures), "executionID", request.Metadata.WorkflowExecutionID)
		return &writetarget.TransmissionState{
			Status:      writetarget.TransmissionStateNotAttempted,
			Transmitter: transmissionInfo.Transmitter.String(),
			Err:         nil,
		}, nil
	case TransmissionStateSucceeded:
		t.lggr.Infow("returning without a transmission attempt - report already onchain ", "executionID", request.Metadata.WorkflowExecutionID)
		return &writetarget.TransmissionState{
			Status:      writetarget.TransmissionStateSucceeded,
			Transmitter: transmissionInfo.Transmitter.String(),
			Err:         nil,
		}, nil
	case TransmissionStateInvalidReceiver:
		t.lggr.Infow("returning without a transmission attempt - transmission already attempted, receiver was marked as invalid", "executionID", request.Metadata.WorkflowExecutionID)
		return &writetarget.TransmissionState{
			Status:      writetarget.TransmissionStateFatal,
			Transmitter: transmissionInfo.Transmitter.String(),
			Err:         ErrTxFailed,
		}, nil
	case TransmissionStateFailed:
		receiverGasMinimum := t.receiverGasMinimum
		if r.Config.GasLimit != nil {
			receiverGasMinimum = *r.Config.GasLimit - ForwarderContractLogicGasCost
		}
		if transmissionInfo.GasLimit.Uint64() > receiverGasMinimum {
			t.lggr.Infow("returning without a transmission attempt - transmission already attempted and failed, sufficient gas was provided", "executionID", request.Metadata.WorkflowExecutionID, "receiverGasMinimum", receiverGasMinimum, "transmissionGasLimit", transmissionInfo.GasLimit)
			return &writetarget.TransmissionState{
				Status:      writetarget.TransmissionStateFatal,
				Transmitter: transmissionInfo.Transmitter.String(),
				Err:         ErrTxFailed,
			}, nil
		}
		t.lggr.Infow("non-empty report - transmission should be retried", "request", request, "reportLen", len(r.Inputs.SignedReport.Report), "reportContextLen", len(r.Inputs.SignedReport.Context), "nSignatures", len(r.Inputs.SignedReport.Signatures), "executionID", request.Metadata.WorkflowExecutionID, "receiverGasMinimum", receiverGasMinimum, "transmissionGasLimit", transmissionInfo.GasLimit)
		return &writetarget.TransmissionState{
			Status:      writetarget.TransmissionStateFailed,
			Transmitter: transmissionInfo.Transmitter.String(),
			Err:         ErrTxFailed,
		}, nil
	}

	return &writetarget.TransmissionState{
		Status:      writetarget.TransmissionStateFatal,
		Transmitter: transmissionInfo.Transmitter.String(),
		Err:         ErrTxFailed,
	}, fmt.Errorf("unexpected transmission state: %v", transmissionInfo.State)
}

func (t *evmTargetStrategy) GetEstimateFee(ctx context.Context, report []byte, reportContext []byte, signatures [][]byte,
	request capabilities.CapabilityRequest) (commontypes.EstimateFee, error) {
	r, err := getEVMRequest(request)
	if err != nil {
		return commontypes.EstimateFee{}, err
	}

	req := getRawReport(r)
	t.lggr.Debugw("Transaction raw report", "report", hex.EncodeToString(req.RawReport))

	meta := commontypes.TxMeta{WorkflowExecutionID: &request.Metadata.WorkflowExecutionID}
	if r.Config.GasLimit != nil {
		meta.GasLimit = new(big.Int).SetUint64(*r.Config.GasLimit)
	}

	value := big.NewInt(0)
	return t.cw.GetEstimateFee(ctx, contractName, method, req, t.forwarder, &meta, value)
}

func (t *evmTargetStrategy) GetTransactionFee(ctx context.Context, transactionID string) (decimal.Decimal, error) {
	fee, err := t.txm.GetTransactionFee(ctx, transactionID)
	if err != nil {
		return decimal.Decimal{}, err
	}
	return decimal.New(fee.TransactionFee.Int64(), -18), nil
}

var (
	contractName = "forwarder"
	method       = "report"
)

// TransmitReport constructs the tx to transmit the report, and defines
// any specific handling for sending the report via ChainWriter.
func (t *evmTargetStrategy) TransmitReport(ctx context.Context, _ []byte, _ []byte, _ [][]byte, request capabilities.CapabilityRequest) (string, error) {
	txID, err := uuid.NewUUID() // NOTE: CW expects us to generate an ID, rather than return one
	if err != nil {
		return "", err
	}

	r, err := getEVMRequest(request)
	if err != nil {
		return txID.String(), fmt.Errorf("failed to getEVMRequest: %w", err)
	}
	req := getRawReport(r)
	t.lggr.Debugw("Transaction raw report", "report", hex.EncodeToString(req.RawReport))

	meta := commontypes.TxMeta{WorkflowExecutionID: &request.Metadata.WorkflowExecutionID}
	if r.Config.GasLimit != nil {
		meta.GasLimit = new(big.Int).SetUint64(*r.Config.GasLimit)
	}

	value := big.NewInt(0)
	if err := t.cw.SubmitTransaction(ctx, contractName, method, req, txID.String(), t.forwarder, &meta, value); err != nil {
		if !commontypes.ErrSettingTransactionGasLimitNotSupported.Is(err) {
			return txID.String(), fmt.Errorf("failed to submit transaction: %w", err)
		}
		meta.GasLimit = nil
		if err := t.cw.SubmitTransaction(ctx, contractName, method, req, txID.String(), t.forwarder, &meta, value); err != nil {
			return txID.String(), fmt.Errorf("failed to submit transaction: %w", err)
		}
	}
	return txID.String(), nil
}

type rawReport struct {
	Receiver      string
	RawReport     []byte
	ReportContext []byte
	Signatures    [][]byte
}

func getRawReport(r TargetRequest) rawReport {
	// Note: The codec that ChainWriter uses to encode the parameters for the contract ABI cannot handle
	// `nil` values, including for slices. Until the bug is fixed we need to ensure that there are no
	// `nil` values passed in the request.
	req := rawReport{r.Config.Address, r.Inputs.SignedReport.Report, r.Inputs.SignedReport.Context, r.Inputs.SignedReport.Signatures}

	if req.RawReport == nil {
		req.RawReport = make([]byte, 0)
	}

	if req.ReportContext == nil {
		req.ReportContext = make([]byte, 0)
	}

	if req.Signatures == nil {
		req.Signatures = make([][]byte, 0)
	}

	return req
}

func (t *evmTargetStrategy) GetTransactionStatus(ctx context.Context, transactionID string) (commontypes.TransactionStatus, error) {
	return t.cw.GetTransactionStatus(ctx, transactionID)
}
