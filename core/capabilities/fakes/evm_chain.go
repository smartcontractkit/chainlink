package fakes

import (
	"context"
	"crypto/ecdsa"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"google.golang.org/protobuf/types/known/emptypb"

	commonCap "github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	evmcappb "github.com/smartcontractkit/chainlink-common/pkg/capabilities/v2/chain-capabilities/evm"
	evmserver "github.com/smartcontractkit/chainlink-common/pkg/capabilities/v2/chain-capabilities/evm/server"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/types/core"
	"github.com/smartcontractkit/chainlink-common/pkg/values/pb"
)

type FakeEVMChain struct {
	commonCap.CapabilityInfo
	services.Service
	eng *services.Engine

	gethClient            *ethclient.Client
	privateKey            *ecdsa.PrivateKey
	mockKeystoneForwarder *MockKeystoneForwarder
	chainSelector         uint64

	lggr logger.Logger

	// log trigger callback channel
	callbackCh map[string]chan commonCap.TriggerAndId[*evmcappb.Log]
}

var evmExecInfo = commonCap.MustNewCapabilityInfo(
	"mainnet-evm@1.0.0",
	commonCap.CapabilityTypeTrigger,
	"A fake evm chain capability that can be used to execute evm chain actions.",
)

var _ services.Service = (*FakeEVMChain)(nil)
var _ evmserver.ClientCapability = (*FakeEVMChain)(nil)
var _ commonCap.ExecutableCapability = (*FakeEVMChain)(nil)

func NewFakeEvmChain(
	lggr logger.Logger,
	gethClient *ethclient.Client,
	privateKey *ecdsa.PrivateKey,
	mockKeystoneForwarderAddress common.Address,
	chainSelector uint64,
) *FakeEVMChain {
	mockKeystoneForwarder, err := NewMockKeystoneForwarder(mockKeystoneForwarderAddress, gethClient)
	if err != nil {
		lggr.Errorw("Failed to create mock keystone forwarder", "error", err)
		return nil
	}

	fc := &FakeEVMChain{
		CapabilityInfo:        evmExecInfo,
		lggr:                  lggr,
		gethClient:            gethClient,
		privateKey:            privateKey,
		mockKeystoneForwarder: mockKeystoneForwarder,
		chainSelector:         chainSelector,
		callbackCh:            make(map[string]chan commonCap.TriggerAndId[*evmcappb.Log]),
	}
	fc.Service, fc.eng = services.Config{
		Name:  "FakeEVMChain",
		Start: fc.Start,
		Close: fc.Close,
	}.NewServiceEngine(lggr)
	return fc
}

func (fc *FakeEVMChain) Initialise(ctx context.Context, config string, _ core.TelemetryService,
	_ core.KeyValueStore,
	_ core.ErrorLog,
	_ core.PipelineRunnerService,
	_ core.RelayerSet,
	_ core.OracleFactory,
	_ core.GatewayConnector,
	_ core.Keystore) error {
	// TODO: do validation of config here

	err := fc.Start(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (fc *FakeEVMChain) CallContract(ctx context.Context, metadata commonCap.RequestMetadata, input *evmcappb.CallContractRequest) (*evmcappb.CallContractReply, error) {
	fc.eng.Infow("EVM Chain CallContract Started")
	fc.eng.Debugw("EVM Chain CallContract Input", "input", input)

	toAddress := common.Address(input.Call.To)
	data := input.Call.Data

	// Make the call
	msg := ethereum.CallMsg{
		To:   &toAddress,
		Data: data,
	}

	// Call contract
	blockNumber := pb.NewIntFromBigInt(input.BlockNumber)
	data, err := fc.gethClient.CallContract(ctx, msg, blockNumber)
	if err != nil {
		return nil, err
	}

	fc.eng.Debugw("EVM Chain CallContract Data Output", "data", new(big.Int).SetBytes(data).String())
	fc.eng.Infow("EVM Chain CallContract Finished")

	// Convert data to protobuf
	return &evmcappb.CallContractReply{
		Data: data,
	}, nil
}

func (fc *FakeEVMChain) WriteReport(ctx context.Context, metadata commonCap.RequestMetadata, input *evmcappb.WriteReportRequest) (*evmcappb.WriteReportReply, error) {
	fc.eng.Infow("EVM Chain WriteReport Started")
	fc.eng.Debugw("EVM Chain WriteReport Input", "input", input)

	// Create authenticated transactor
	chainID, err := fc.gethClient.ChainID(ctx)
	if err != nil {
		return nil, err
	}

	auth, err := bind.NewKeyedTransactorWithChainID(fc.privateKey, chainID)
	if err != nil {
		return nil, err
	}

	signatures := make([][]byte, len(input.Report.Sigs))
	for i, sig := range input.Report.Sigs {
		signatures[i] = sig.Signature
	}

	reportTx, err := fc.mockKeystoneForwarder.Report(
		auth,
		common.Address(input.Receiver),
		input.Report.RawReport,
		input.Report.ReportContext,
		signatures,
	)
	if err != nil {
		return nil, err
	}

	// TODO: should we wait for the transaction to be mined?
	receipt, err := bind.WaitMined(ctx, fc.gethClient, reportTx)
	if err != nil {
		return nil, err
	}

	fc.eng.Debugw("EVM Chain WriteReport Receipt", "status", receipt.Status, "gasUsed", receipt.GasUsed, "txHash", receipt.TxHash.Hex())
	txHash := receipt.TxHash.Bytes()

	// Calculate transaction fee (gas used * effective gas price)
	transactionFee := new(big.Int).Mul(new(big.Int).SetUint64(receipt.GasUsed), receipt.EffectiveGasPrice)

	if receipt.Status == types.ReceiptStatusSuccessful {
		fc.eng.Infow("EVM Chain WriteReport Successful", "txHash", receipt.TxHash.Hex(), "gasUsed", receipt.GasUsed, "fee", transactionFee.String())

		receiverStatus := evmcappb.ReceiverContractExecutionStatus_RECEIVER_CONTRACT_EXECUTION_STATUS_SUCCESS
		return &evmcappb.WriteReportReply{
			TxStatus:                        evmcappb.TxStatus_TX_STATUS_SUCCESS,
			ReceiverContractExecutionStatus: &receiverStatus,
			TxHash:                          txHash,
			TransactionFee:                  pb.NewBigIntFromInt(transactionFee),
		}, nil
	}

	fc.eng.Infow("EVM Chain WriteReport Failed", "txHash", receipt.TxHash.Hex(), "gasUsed", receipt.GasUsed, "fee", transactionFee.String())
	receiverStatus := evmcappb.ReceiverContractExecutionStatus_RECEIVER_CONTRACT_EXECUTION_STATUS_REVERTED
	errorMsg := "Transaction reverted"
	return &evmcappb.WriteReportReply{
		TxStatus:                        evmcappb.TxStatus_TX_STATUS_REVERTED,
		ReceiverContractExecutionStatus: &receiverStatus,
		TxHash:                          txHash,
		TransactionFee:                  pb.NewBigIntFromInt(transactionFee),
		ErrorMessage:                    &errorMsg,
	}, nil
}

func (fc *FakeEVMChain) RegisterLogTrigger(ctx context.Context, triggerID string, metadata commonCap.RequestMetadata, input *evmcappb.FilterLogTriggerRequest) (<-chan commonCap.TriggerAndId[*evmcappb.Log], error) {
	fc.callbackCh[triggerID] = make(chan commonCap.TriggerAndId[*evmcappb.Log])
	return fc.callbackCh[triggerID], nil
}

func (fc *FakeEVMChain) UnregisterLogTrigger(ctx context.Context, triggerID string, metadata commonCap.RequestMetadata, input *evmcappb.FilterLogTriggerRequest) error {
	return nil
}

func (fc *FakeEVMChain) ManualTrigger(ctx context.Context, triggerID string, log *evmcappb.Log) error {
	fc.eng.Debugf("ManualTrigger: %s", log.String())

	go func() {
		select {
		case fc.callbackCh[triggerID] <- fc.createManualTriggerEvent(log):
			// Successfully sent trigger response
		case <-ctx.Done():
			// Context cancelled, cleanup goroutine
			fc.eng.Debug("ManualTrigger goroutine cancelled due to context cancellation")
		}
	}()

	return nil
}

func (fc *FakeEVMChain) createManualTriggerEvent(log *evmcappb.Log) commonCap.TriggerAndId[*evmcappb.Log] {
	return commonCap.TriggerAndId[*evmcappb.Log]{
		Trigger: log,
		Id:      "manual-evm-chain-trigger-id",
	}
}

func (fc *FakeEVMChain) FilterLogs(ctx context.Context, metadata commonCap.RequestMetadata, input *evmcappb.FilterLogsRequest) (*evmcappb.FilterLogsReply, error) {
	fc.eng.Infow("EVM Chain FilterLogs Started", "input", input)

	// Prepare filter query
	filterQueryPb := input.GetFilterQuery()
	addresses := make([]common.Address, len(filterQueryPb.Addresses))
	for i, address := range filterQueryPb.Addresses {
		addresses[i] = common.Address(address)
	}
	filterQuery := ethereum.FilterQuery{
		FromBlock: new(big.Int).SetBytes(filterQueryPb.FromBlock.AbsVal),
		ToBlock:   new(big.Int).SetBytes(filterQueryPb.ToBlock.AbsVal),
		Addresses: addresses,
	}

	// Filter logs
	logs, err := fc.gethClient.FilterLogs(ctx, filterQuery)
	if err != nil {
		return nil, err
	}

	fc.eng.Infow("EVM Chain FilterLogs Finished", "logs", logs)

	// Convert logs to protobuf
	logsPb := make([]*evmcappb.Log, len(logs))
	for i, log := range logs {
		logsPb[i] = &evmcappb.Log{
			Address: log.Address.Bytes(),
			Data:    log.Data,
			Topics:  logsPb[i].Topics,
		}
	}
	return &evmcappb.FilterLogsReply{
		Logs: logsPb,
	}, nil
}

func (fc *FakeEVMChain) BalanceAt(ctx context.Context, metadata commonCap.RequestMetadata, input *evmcappb.BalanceAtRequest) (*evmcappb.BalanceAtReply, error) {
	fc.eng.Infow("EVM Chain BalanceAt Started", "input", input)

	// Prepare balance at request
	address := common.Address(input.Account)
	blockNumber := new(big.Int).SetBytes(input.BlockNumber.AbsVal)

	// Get balance at block number
	balance, err := fc.gethClient.BalanceAt(ctx, address, blockNumber)
	if err != nil {
		return nil, err
	}

	// Convert balance to protobuf
	return &evmcappb.BalanceAtReply{
		Balance: pb.NewBigIntFromInt(balance),
	}, nil
}

func (fc *FakeEVMChain) EstimateGas(ctx context.Context, metadata commonCap.RequestMetadata, input *evmcappb.EstimateGasRequest) (*evmcappb.EstimateGasReply, error) {
	fc.eng.Infow("EVM Chain EstimateGas Started", "input", input)

	// Prepare estimate gas request
	toAddress := common.Address(input.Msg.To)
	msg := ethereum.CallMsg{
		From: common.Address(input.Msg.From),
		To:   &toAddress,
		Data: input.Msg.Data,
	}

	// Estimate gas
	gas, err := fc.gethClient.EstimateGas(ctx, msg)
	if err != nil {
		return nil, err
	}

	// Convert gas to protobuf
	fc.eng.Infow("EVM Chain EstimateGas Finished", "gas", gas)
	return &evmcappb.EstimateGasReply{
		Gas: gas,
	}, nil
}

func (fc *FakeEVMChain) GetTransactionByHash(ctx context.Context, metadata commonCap.RequestMetadata, input *evmcappb.GetTransactionByHashRequest) (*evmcappb.GetTransactionByHashReply, error) {
	fc.eng.Infow("EVM Chain GetTransactionByHash Started", "input", input)

	// Prepare get transaction by hash request
	hash := common.Hash(input.Hash)

	// Get transaction by hash
	transaction, pending, err := fc.gethClient.TransactionByHash(ctx, hash)
	if err != nil {
		return nil, err
	}

	fc.eng.Infow("EVM Chain GetTransactionByHash Finished", "transaction", transaction, "pending", pending)

	// Convert transaction to protobuf
	transactionPb := &evmcappb.Transaction{
		To:       transaction.To().Bytes(),
		Data:     transaction.Data(),
		Hash:     transaction.Hash().Bytes(),
		Value:    pb.NewBigIntFromInt(transaction.Value()),
		GasPrice: pb.NewBigIntFromInt(transaction.GasPrice()),
		Nonce:    transaction.Nonce(),
	}
	return &evmcappb.GetTransactionByHashReply{
		Transaction: transactionPb,
	}, nil
}

func (fc *FakeEVMChain) GetTransactionReceipt(ctx context.Context, metadata commonCap.RequestMetadata, input *evmcappb.GetTransactionReceiptRequest) (*evmcappb.GetTransactionReceiptReply, error) {
	fc.eng.Infow("EVM Chain GetTransactionReceipt Started", "input", input)

	// Prepare get transaction receipt request
	hash := common.Hash(input.Hash)

	// Get transaction receipt
	receipt, err := fc.gethClient.TransactionReceipt(ctx, hash)
	if err != nil {
		return nil, err
	}

	fc.eng.Infow("EVM Chain GetTransactionReceipt Finished", "receipt", receipt)

	// Convert transaction receipt to protobuf
	receiptPb := &evmcappb.Receipt{
		Status:            receipt.Status,
		Logs:              make([]*evmcappb.Log, len(receipt.Logs)),
		GasUsed:           receipt.GasUsed,
		TxIndex:           uint64(receipt.TransactionIndex),
		BlockHash:         receipt.BlockHash.Bytes(),
		TxHash:            receipt.TxHash.Bytes(),
		EffectiveGasPrice: pb.NewBigIntFromInt(receipt.EffectiveGasPrice),
		BlockNumber:       pb.NewBigIntFromInt(receipt.BlockNumber),
		ContractAddress:   receipt.ContractAddress.Bytes(),
	}
	for i, log := range receipt.Logs {
		receiptPb.Logs[i] = &evmcappb.Log{
			Address: log.Address.Bytes(),
		}
	}
	return &evmcappb.GetTransactionReceiptReply{
		Receipt: receiptPb,
	}, nil
}

func (fc *FakeEVMChain) HeaderByNumber(ctx context.Context, metadata commonCap.RequestMetadata, input *evmcappb.HeaderByNumberRequest) (*evmcappb.HeaderByNumberReply, error) {
	fc.eng.Infow("EVM Chain HeaderByNumber Started", "input", input)

	// Prepare header by number request
	blockNumber := new(big.Int).SetBytes(input.BlockNumber.AbsVal)

	// Get header by number
	header, err := fc.gethClient.HeaderByNumber(ctx, blockNumber)
	if err != nil {
		return nil, err
	}

	// Convert header to protobuf
	headerPb := &evmcappb.HeaderByNumberReply{
		Header: &evmcappb.Header{
			Timestamp:   header.Time,
			BlockNumber: pb.NewBigIntFromInt(header.Number),
			Hash:        header.Hash().Bytes(),
			ParentHash:  header.ParentHash.Bytes(),
		},
	}

	fc.eng.Infow("EVM Chain HeaderByNumber Finished", "header", headerPb)
	return headerPb, nil
}

func (fc *FakeEVMChain) RegisterLogTracking(ctx context.Context, metadata commonCap.RequestMetadata, input *evmcappb.RegisterLogTrackingRequest) (*emptypb.Empty, error) {
	fc.eng.Infow("EVM Chain registered log tracking", "input", input)
	return nil, nil
}

func (fc *FakeEVMChain) UnregisterLogTracking(ctx context.Context, metadata commonCap.RequestMetadata, input *evmcappb.UnregisterLogTrackingRequest) (*emptypb.Empty, error) {
	fc.eng.Infow("EVM Chain unregistered log tracking", "input", input)
	return nil, nil
}

func (fc *FakeEVMChain) Name() string {
	return fc.ID
}

func (fc *FakeEVMChain) HealthReport() map[string]error {
	return map[string]error{fc.Name(): nil}
}

func (fc *FakeEVMChain) Start(ctx context.Context) error {
	fc.eng.Debugw("EVM Chain started")
	return nil
}

func (fc *FakeEVMChain) Close() error {
	fc.eng.Debugw("EVM Chain closed")
	return nil
}

func (fc *FakeEVMChain) RegisterToWorkflow(ctx context.Context, request commonCap.RegisterToWorkflowRequest) error {
	fc.eng.Infow("Registered to EVM Chain", "workflowID", request.Metadata.WorkflowID)
	return nil
}

func (fc *FakeEVMChain) UnregisterFromWorkflow(ctx context.Context, request commonCap.UnregisterFromWorkflowRequest) error {
	fc.eng.Infow("Unregistered from EVM Chain", "workflowID", request.Metadata.WorkflowID)
	return nil
}

func (fc *FakeEVMChain) Execute(ctx context.Context, request commonCap.CapabilityRequest) (commonCap.CapabilityResponse, error) {
	fc.eng.Infow("EVM Chain executed", "request", request)
	return commonCap.CapabilityResponse{}, nil
}

func (fc *FakeEVMChain) Description() string {
	return "EVM Chain"
}

func (fc *FakeEVMChain) ChainSelector() uint64 {
	return fc.chainSelector
}
