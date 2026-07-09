package main

import (
	"context"
	"errors"
	"math/big"
	"time"

	commoncap "github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	caperrors "github.com/smartcontractkit/chainlink-common/pkg/capabilities/errors"
	evmcappb "github.com/smartcontractkit/chainlink-common/pkg/capabilities/v2/chain-capabilities/evm"
	evmserver "github.com/smartcontractkit/chainlink-common/pkg/capabilities/v2/chain-capabilities/evm/server"
	commonlogger "github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/types/core"
	"github.com/smartcontractkit/chainlink-protos/cre/go/values/pb"
)

// fakeEVMWriteLag is an artificial delay injected into EVM WriteReport calls to
// simulate the latency of submitting a report on-chain and awaiting inclusion.
const fakeEVMWriteLag = 60 * time.Second

// fakeEVMWrite is a minimal, RPC-free EVM chain capability. It implements just
// enough of evmserver.ClientCapability for v2 workflows whose only on-chain
// interaction is WriteReport: that method sleeps for fakeEVMWriteLag and returns
// success; every other method is unimplemented. Wrap it with
// evmserver.NewClientServer so it registers under evm:ChainSelector:<sel>@1.0.0.
type fakeEVMWrite struct {
	chainSelector uint64
	lggr          commonlogger.Logger
}

var _ evmserver.ClientCapability = (*fakeEVMWrite)(nil)

func newFakeEVMWrite(lggr commonlogger.Logger, chainSelector uint64) *fakeEVMWrite {
	return &fakeEVMWrite{chainSelector: chainSelector, lggr: commonlogger.Named(lggr, "fakeEVMWrite")}
}

func (f *fakeEVMWrite) ChainSelector() uint64 { return f.chainSelector }

func (f *fakeEVMWrite) WriteReport(ctx context.Context, metadata commoncap.RequestMetadata, input *evmcappb.WriteReportRequest) (*commoncap.ResponseAndMetadata[*evmcappb.WriteReportReply], caperrors.Error) {
	f.lggr.Infow("Fake EVM WriteReport started", "workflowID", metadata.WorkflowID, "executionID", metadata.WorkflowExecutionID)
	time.Sleep(fakeEVMWriteLag) // simulate on-chain write latency

	receiverStatus := evmcappb.ReceiverContractExecutionStatus_RECEIVER_CONTRACT_EXECUTION_STATUS_SUCCESS
	resp := &evmcappb.WriteReportReply{
		TxStatus:                        evmcappb.TxStatus_TX_STATUS_SUCCESS,
		ReceiverContractExecutionStatus: &receiverStatus,
		TxHash:                          make([]byte, 32),
		TransactionFee:                  pb.NewBigIntFromInt(big.NewInt(0)),
	}
	f.lggr.Infow("Fake EVM WriteReport finished", "executionID", metadata.WorkflowExecutionID)
	return &commoncap.ResponseAndMetadata[*evmcappb.WriteReportReply]{Response: resp}, nil
}

func unimplemented() caperrors.Error {
	return caperrors.NewPublicSystemError(errors.New("not implemented by fakeEVMWrite"), caperrors.Unknown)
}

func (f *fakeEVMWrite) CallContract(context.Context, commoncap.RequestMetadata, *evmcappb.CallContractRequest) (*commoncap.ResponseAndMetadata[*evmcappb.CallContractReply], caperrors.Error) {
	return nil, unimplemented()
}

func (f *fakeEVMWrite) FilterLogs(context.Context, commoncap.RequestMetadata, *evmcappb.FilterLogsRequest) (*commoncap.ResponseAndMetadata[*evmcappb.FilterLogsReply], caperrors.Error) {
	return nil, unimplemented()
}

func (f *fakeEVMWrite) BalanceAt(context.Context, commoncap.RequestMetadata, *evmcappb.BalanceAtRequest) (*commoncap.ResponseAndMetadata[*evmcappb.BalanceAtReply], caperrors.Error) {
	return nil, unimplemented()
}

func (f *fakeEVMWrite) EstimateGas(context.Context, commoncap.RequestMetadata, *evmcappb.EstimateGasRequest) (*commoncap.ResponseAndMetadata[*evmcappb.EstimateGasReply], caperrors.Error) {
	return nil, unimplemented()
}

func (f *fakeEVMWrite) GetTransactionByHash(context.Context, commoncap.RequestMetadata, *evmcappb.GetTransactionByHashRequest) (*commoncap.ResponseAndMetadata[*evmcappb.GetTransactionByHashReply], caperrors.Error) {
	return nil, unimplemented()
}

func (f *fakeEVMWrite) GetTransactionReceipt(context.Context, commoncap.RequestMetadata, *evmcappb.GetTransactionReceiptRequest) (*commoncap.ResponseAndMetadata[*evmcappb.GetTransactionReceiptReply], caperrors.Error) {
	return nil, unimplemented()
}

func (f *fakeEVMWrite) HeaderByNumber(context.Context, commoncap.RequestMetadata, *evmcappb.HeaderByNumberRequest) (*commoncap.ResponseAndMetadata[*evmcappb.HeaderByNumberReply], caperrors.Error) {
	return nil, unimplemented()
}

func (f *fakeEVMWrite) RegisterLogTrigger(context.Context, string, commoncap.RequestMetadata, *evmcappb.FilterLogTriggerRequest) (<-chan commoncap.TriggerAndId[*evmcappb.Log], caperrors.Error) {
	return nil, unimplemented()
}

func (f *fakeEVMWrite) UnregisterLogTrigger(context.Context, string, commoncap.RequestMetadata, *evmcappb.FilterLogTriggerRequest) caperrors.Error {
	return unimplemented()
}

func (f *fakeEVMWrite) AckEvent(context.Context, string, string, string) caperrors.Error {
	return unimplemented()
}

func (f *fakeEVMWrite) Initialise(context.Context, core.StandardCapabilitiesDependencies) error {
	return nil
}
func (f *fakeEVMWrite) Start(context.Context) error      { return nil }
func (f *fakeEVMWrite) Close() error                     { return nil }
func (f *fakeEVMWrite) Ready() error                     { return nil }
func (f *fakeEVMWrite) HealthReport() map[string]error   { return map[string]error{f.Name(): nil} }
func (f *fakeEVMWrite) Name() string                     { return "fakeEVMWrite" }
func (f *fakeEVMWrite) Description() string               { return "Fake EVM write capability (no RPC)" }
