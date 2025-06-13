package fakes

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	commonCap "github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	evmserver "github.com/smartcontractkit/chainlink-common/pkg/capabilities/v2/chain-capabilities/evm/server"
	evmpb "github.com/smartcontractkit/chainlink-common/pkg/chains/evm"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/types/core"
	"github.com/smartcontractkit/chainlink-common/pkg/values/pb"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"google.golang.org/protobuf/types/known/emptypb"
)

type fakeEvmChain struct {
	commonCap.CapabilityInfo
	services.Service
	eng *services.Engine

	gethClient *ethclient.Client

	lggr logger.Logger
}

var evmExecInfo = capabilities.MustNewCapabilityInfo(
	"mainnet-evm@1.0.0",
	capabilities.CapabilityTypeTrigger,
	"A fake evm chain capability that can be used to execute evm chain actions.",
)

var _ services.Service = (*fakeEvmChain)(nil)
var _ evmserver.ClientCapability = (*fakeEvmChain)(nil)
var _ commonCap.ExecutableCapability = (*fakeEvmChain)(nil)

func NewFakeEvmChain(lggr logger.Logger, gethClient *ethclient.Client) *fakeEvmChain {
	fc := &fakeEvmChain{
		CapabilityInfo: evmExecInfo,
		lggr:           lggr,
		gethClient:     gethClient,
	}
	fc.Service, fc.eng = services.Config{
		Name:  "fakeEvmChain",
		Start: fc.Start,
		Close: fc.Close,
	}.NewServiceEngine(lggr)
	return fc
}

func (fc *fakeEvmChain) Initialise(ctx context.Context, config string, _ core.TelemetryService,
	_ core.KeyValueStore,
	_ core.ErrorLog,
	_ core.PipelineRunnerService,
	_ core.RelayerSet,
	_ core.OracleFactory) error {

	// TODO: do validation of config here

	err := fc.Start(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (fc *fakeEvmChain) CallContract(ctx context.Context, metadata capabilities.RequestMetadata, input *evmpb.CallContractRequest) (*evmpb.CallContractReply, error) {
	fc.eng.Infow("Fake EVM Chain CallContract Started", "input", input)

	// Prepare call msg
	// toAddress := common.Address(common.Address(input.Call.To.Address))
	toAddress := common.HexToAddress("0x779877A7B0D9E8603169DdbD7836e478b4624789")
	walletAddress := common.HexToAddress("0x437bb34CbdB6c0Eaf859FfDC2DfC424d710e4C5B")

	// balanceOf(address) selector
	methodID := []byte{0x70, 0xa0, 0x82, 0x31}

	// Pad the address to 32 bytes
	paddedAddress := common.LeftPadBytes(walletAddress.Bytes(), 32)

	// Combine method selector and padded address
	data := append(methodID, paddedAddress...)

	// Make the call
	msg := ethereum.CallMsg{
		To:   &toAddress,
		Data: data,
	}

	// Call contract
	data, err := fc.gethClient.CallContract(ctx, msg, nil)
	if err != nil {
		return nil, err
	}

	fc.eng.Infow("Fake EVM Chain CallContract Finished", "data", data)

	// Convert data to protobuf
	return &evmpb.CallContractReply{
		Data: data,
	}, nil
}

func (fc *fakeEvmChain) IsTxFinalized(ctx context.Context, metadata capabilities.RequestMetadata, input *evmpb.IsTxFinalizedRequest) (*evmpb.IsTxFinalizedReply, error) {
	fc.eng.Infow("Fake EVM Chain IsTxFinalized Started", "input", input)

	// Prepare is tx finalized request
	hash := common.Hash(input.TxHash)

	// Get transaction receipt
	receipt, err := fc.gethClient.TransactionReceipt(ctx, hash)
	if err != nil {
		return nil, err
	}

	return &evmpb.IsTxFinalizedReply{
		IsFinalized: receipt.Status == 1,
	}, nil
}

func (fc *fakeEvmChain) WriteReport(ctx context.Context, metadata capabilities.RequestMetadata, input *evmpb.WriteReportRequest) (*evmpb.WriteReportReply, error) {
	fc.eng.Infow("Fake EVM Chain WriteReport Started", "input", input)

	errMsg := ""
	return &evmpb.WriteReportReply{
		TxStatus:                        evmpb.TransactionStatus_TX_SUCCESS,
		TxHash:                          []byte{},
		ReceiverContractExecutionStatus: evmpb.ReceiverContractExecutionStatus_SUCCESS.Enum(),
		TransactionFee:                  pb.NewBigIntFromInt(big.NewInt(0)),
		ErrorMessage:                    &errMsg,
	}, nil
}

func (fc *fakeEvmChain) FilterLogs(ctx context.Context, metadata capabilities.RequestMetadata, input *evmpb.FilterLogsRequest) (*evmpb.FilterLogsReply, error) {
	fc.eng.Infow("Fake EVM Chain FilterLogs Started", "input", input)

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

	fc.eng.Infow("Fake EVM Chain FilterLogs Finished", "logs", logs)

	// Convert logs to protobuf
	logsPb := make([]*evmpb.Log, len(logs))
	for i, log := range logs {
		logsPb[i] = &evmpb.Log{
			Address: log.Address.Bytes(),
			Data:    log.Data,
			Topics:  logsPb[i].Topics,
		}
	}
	return &evmpb.FilterLogsReply{
		Logs: logsPb,
	}, nil
}

func (fc *fakeEvmChain) BalanceAt(ctx context.Context, metadata capabilities.RequestMetadata, input *evmpb.BalanceAtRequest) (*evmpb.BalanceAtReply, error) {
	fc.eng.Infow("Fake EVM Chain BalanceAt Started", "input", input)

	// Prepare balance at request
	address := common.Address(input.Account)
	blockNumber := new(big.Int).SetBytes(input.BlockNumber.AbsVal)

	// Get balance at block number
	balance, err := fc.gethClient.BalanceAt(ctx, address, blockNumber)
	if err != nil {
		return nil, err
	}

	// Convert balance to protobuf
	return &evmpb.BalanceAtReply{
		Balance: pb.NewBigIntFromInt(balance),
	}, nil
}

func (fc *fakeEvmChain) EstimateGas(ctx context.Context, metadata capabilities.RequestMetadata, input *evmpb.EstimateGasRequest) (*evmpb.EstimateGasReply, error) {
	fc.eng.Infow("Fake EVM Chain EstimateGas Started", "input", input)

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
	fc.eng.Infow("Fake EVM Chain EstimateGas Finished", "gas", gas)
	return &evmpb.EstimateGasReply{
		Gas: gas,
	}, nil
}

func (fc *fakeEvmChain) GetTransactionByHash(ctx context.Context, metadata capabilities.RequestMetadata, input *evmpb.GetTransactionByHashRequest) (*evmpb.GetTransactionByHashReply, error) {
	fc.eng.Infow("Fake EVM Chain GetTransactionByHash Started", "input", input)

	// Prepare get transaction by hash request
	hash := common.Hash(input.Hash)

	// Get transaction by hash
	transaction, pending, err := fc.gethClient.TransactionByHash(ctx, hash)
	if err != nil {
		return nil, err
	}

	fc.eng.Infow("Fake EVM Chain GetTransactionByHash Finished", "transaction", transaction, "pending", pending)

	// Convert transaction to protobuf
	transactionPb := &evmpb.Transaction{
		To:       transaction.To().Bytes(),
		Data:     transaction.Data(),
		Hash:     transaction.Hash().Bytes(),
		Value:    pb.NewBigIntFromInt(transaction.Value()),
		GasPrice: pb.NewBigIntFromInt(transaction.GasPrice()),
		Nonce:    transaction.Nonce(),
	}
	return &evmpb.GetTransactionByHashReply{
		Transaction: transactionPb,
	}, nil
}

func (fc *fakeEvmChain) GetTransactionReceipt(ctx context.Context, metadata capabilities.RequestMetadata, input *evmpb.GetTransactionReceiptRequest) (*evmpb.GetTransactionReceiptReply, error) {
	fc.eng.Infow("Fake EVM Chain GetTransactionReceipt Started", "input", input)

	// Prepare get transaction receipt request
	hash := common.Hash(input.Hash)

	// Get transaction receipt
	receipt, err := fc.gethClient.TransactionReceipt(ctx, hash)
	if err != nil {
		return nil, err
	}

	fc.eng.Infow("Fake EVM Chain GetTransactionReceipt Finished", "receipt", receipt)

	// Convert transaction receipt to protobuf
	receiptPb := &evmpb.Receipt{
		Status:            receipt.Status,
		Logs:              make([]*evmpb.Log, len(receipt.Logs)),
		GasUsed:           receipt.GasUsed,
		TxIndex:           uint64(receipt.TransactionIndex),
		BlockHash:         receipt.BlockHash.Bytes(),
		TxHash:            receipt.TxHash.Bytes(),
		EffectiveGasPrice: pb.NewBigIntFromInt(receipt.EffectiveGasPrice),
		BlockNumber:       pb.NewBigIntFromInt(receipt.BlockNumber),
		ContractAddress:   receipt.ContractAddress.Bytes(),
	}
	for i, log := range receipt.Logs {
		receiptPb.Logs[i] = &evmpb.Log{
			Address: log.Address.Bytes(),
		}
	}
	return &evmpb.GetTransactionReceiptReply{
		Receipt: receiptPb,
	}, nil
}

func (fc *fakeEvmChain) LatestAndFinalizedHead(ctx context.Context, metadata capabilities.RequestMetadata, input *emptypb.Empty) (*evmpb.LatestAndFinalizedHeadReply, error) {
	fc.eng.Infow("Fake EVM Chain latest and finalized head", "input", input)

	// Get latest and finalized head
	head, err := fc.gethClient.HeaderByNumber(ctx, nil)
	if err != nil {
		return nil, err
	}

	// Convert head to protobuf
	headPb := &evmpb.LatestAndFinalizedHeadReply{
		Latest: &evmpb.Head{
			Timestamp:   head.Time,
			BlockNumber: pb.NewBigIntFromInt(head.Number),
			Hash:        head.Hash().Bytes(),
			ParentHash:  head.ParentHash.Bytes(),
		},
	}
	return headPb, nil
}

func (fc *fakeEvmChain) QueryTrackedLogs(ctx context.Context, metadata capabilities.RequestMetadata, input *evmpb.QueryTrackedLogsRequest) (*evmpb.QueryTrackedLogsReply, error) {
	fc.eng.Infow("Fake EVM Chain QueryTrackedLogs Started", "input", input)

	// Prepare query tracked logs request
	// fc.gethClient.
	// 	filterQueryPb := input.GetFilterQuery()
	// addresses := make([]common.Address, len(filterQueryPb.Addresses))
	// for i, address := range filterQueryPb.Addresses {
	// 	addresses[i] = common.Address(address.Address)
	// }
	// filterQuery := ethereum.FilterQuery{
	// 	FromBlock: new(big.Int).SetBytes(filterQueryPb.FromBlock.AbsVal),
	// 	ToBlock:   new(big.Int).SetBytes(filterQueryPb.ToBlock.AbsVal),
	// 	Addresses: addresses,
	// }
	return nil, nil
}

func (fc *fakeEvmChain) RegisterLogTracking(ctx context.Context, metadata capabilities.RequestMetadata, input *evmpb.RegisterLogTrackingRequest) (*emptypb.Empty, error) {
	fc.eng.Infow("Fake Evm Chain registered log tracking", "input", input)
	return nil, nil
}

func (fc *fakeEvmChain) UnregisterLogTracking(ctx context.Context, metadata capabilities.RequestMetadata, input *evmpb.UnregisterLogTrackingRequest) (*emptypb.Empty, error) {
	fc.eng.Infow("Fake Evm Chain unregistered log tracking", "input", input)
	return nil, nil
}

func (fc *fakeEvmChain) Name() string {
	return fc.CapabilityInfo.ID
}

func (fc *fakeEvmChain) HealthReport() map[string]error {
	return map[string]error{fc.Name(): nil}
}

func (fc *fakeEvmChain) Start(ctx context.Context) error {
	fc.eng.Infow("Fake Evm Chain started")
	return nil
}

func (fc *fakeEvmChain) Close() error {
	fc.eng.Infow("Fake Evm Chain closed")
	return nil
}

func (fc *fakeEvmChain) RegisterToWorkflow(ctx context.Context, request commonCap.RegisterToWorkflowRequest) error {
	fc.eng.Infow("Registered to Fake Evm Chain", "workflowID", request.Metadata.WorkflowID)
	return nil
}

func (fc *fakeEvmChain) UnregisterFromWorkflow(ctx context.Context, request commonCap.UnregisterFromWorkflowRequest) error {
	fc.eng.Infow("Unregistered from Fake Evm Chain", "workflowID", request.Metadata.WorkflowID)
	return nil
}

func (fc *fakeEvmChain) Execute(ctx context.Context, request commonCap.CapabilityRequest) (commonCap.CapabilityResponse, error) {
	fc.eng.Infow("Fake Evm Chain executed", "request", request)
	return commonCap.CapabilityResponse{}, nil
}

func (fc *fakeEvmChain) Description() string {
	return "Fake Evm Chain"
}
