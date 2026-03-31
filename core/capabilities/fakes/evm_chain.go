package fakes

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"

	commonCap "github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	caperrors "github.com/smartcontractkit/chainlink-common/pkg/capabilities/errors"
	evmcappb "github.com/smartcontractkit/chainlink-common/pkg/capabilities/v2/chain-capabilities/evm"
	evmserver "github.com/smartcontractkit/chainlink-common/pkg/capabilities/v2/chain-capabilities/evm/server"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/types/core"
	"github.com/smartcontractkit/chainlink-protos/cre/go/values/pb"
)

type FakeEVMChain struct {
	commonCap.CapabilityInfo
	services.Service
	eng *services.Engine

	gethClient                   *ethclient.Client
	privateKey                   *ecdsa.PrivateKey
	mockKeystoneForwarder        *MockKeystoneForwarder
	mockKeystoneForwarderAddress common.Address
	chainSelector                uint64

	// if true, WriteReport will simulate the call and not broadcast
	dryRunWrites bool

	lggr logger.Logger

	// log trigger callback channels and their registered filters
	callbackCh        map[string]chan commonCap.TriggerAndId[*evmcappb.Log]
	logTriggerFilters map[string]*evmcappb.FilterLogTriggerRequest
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
	dryRunWrites bool,
) *FakeEVMChain {
	mockKeystoneForwarder, err := NewMockKeystoneForwarder(mockKeystoneForwarderAddress, gethClient)
	if err != nil {
		lggr.Errorw("Failed to create mock keystone forwarder", "error", err)
		return nil
	}

	fc := &FakeEVMChain{
		CapabilityInfo:               evmExecInfo,
		lggr:                         lggr,
		gethClient:                   gethClient,
		privateKey:                   privateKey,
		mockKeystoneForwarder:        mockKeystoneForwarder,
		mockKeystoneForwarderAddress: mockKeystoneForwarderAddress,
		chainSelector:                chainSelector,
		callbackCh:                   make(map[string]chan commonCap.TriggerAndId[*evmcappb.Log]),
		logTriggerFilters:            make(map[string]*evmcappb.FilterLogTriggerRequest),
		dryRunWrites:                 dryRunWrites,
	}
	fc.Service, fc.eng = services.Config{
		Name:  "FakeEVMChain",
		Start: fc.Start,
		Close: fc.Close,
	}.NewServiceEngine(lggr)
	return fc
}

func (fc *FakeEVMChain) Initialise(ctx context.Context, dependencies core.StandardCapabilitiesDependencies) error {
	// TODO: do validation of config here

	err := fc.Start(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (fc *FakeEVMChain) CallContract(ctx context.Context, metadata commonCap.RequestMetadata, input *evmcappb.CallContractRequest) (*commonCap.ResponseAndMetadata[*evmcappb.CallContractReply], caperrors.Error) {
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
		return nil, caperrors.NewPublicSystemError(err, caperrors.Unknown)
	}

	fc.eng.Debugw("EVM Chain CallContract Data Output", "data", new(big.Int).SetBytes(data).String())
	fc.eng.Infow("EVM Chain CallContract Finished")

	// Convert data to protobuf
	response := &evmcappb.CallContractReply{
		Data: data,
	}
	responseAndMetadata := commonCap.ResponseAndMetadata[*evmcappb.CallContractReply]{
		Response:         response,
		ResponseMetadata: commonCap.ResponseMetadata{},
	}
	return &responseAndMetadata, nil
}

func (fc *FakeEVMChain) WriteReport(
	ctx context.Context,
	metadata commonCap.RequestMetadata,
	input *evmcappb.WriteReportRequest,
) (*commonCap.ResponseAndMetadata[*evmcappb.WriteReportReply], caperrors.Error) {
	fc.eng.Infow("EVM Chain WriteReport Started")
	fc.eng.Debugw("EVM Chain WriteReport Input", "input", input)

	// Create authenticated transactor
	chainID, err := fc.gethClient.ChainID(ctx)
	if err != nil {
		return nil, caperrors.NewPublicSystemError(err, caperrors.Unknown)
	}

	auth, err := bind.NewKeyedTransactorWithChainID(fc.privateKey, chainID)
	if err != nil {
		return nil, caperrors.NewPublicSystemError(err, caperrors.Unknown)
	}

	// Set gas limit if provided
	if gc := input.GasConfig; gc != nil {
		auth.GasLimit = gc.GasLimit
	}

	signatures := make([][]byte, len(input.Report.Sigs))
	for i, sig := range input.Report.Sigs {
		signatures[i] = sig.Signature
	}

	// If dryRunWrites is enabled, simulate the transaction without broadcasting it
	if fc.dryRunWrites {
		resp, dryRunErr := fc.dryRunWriteReport(ctx, auth.From, input, signatures)
		if dryRunErr != nil {
			return nil, caperrors.NewPublicSystemError(dryRunErr, caperrors.Unknown)
		}

		return resp, nil
	}

	reportTx, err := fc.mockKeystoneForwarder.Report(
		auth,
		common.Address(input.Receiver),
		input.Report.RawReport,
		input.Report.ReportContext,
		signatures,
	)
	if err != nil {
		return nil, caperrors.NewPublicSystemError(err, caperrors.Unknown)
	}

	// TODO: should we wait for the transaction to be mined?
	receipt, err := bind.WaitMined(ctx, fc.gethClient, reportTx)
	if err != nil {
		return nil, caperrors.NewPublicSystemError(err, caperrors.Unknown)
	}

	fc.eng.Debugw("EVM Chain WriteReport Receipt", "status", receipt.Status, "gasUsed", receipt.GasUsed, "txHash", receipt.TxHash.Hex())
	txHash := receipt.TxHash.Bytes()

	// Calculate transaction fee (gas used * effective gas price)
	transactionFee := new(big.Int).Mul(new(big.Int).SetUint64(receipt.GasUsed), receipt.EffectiveGasPrice)

	if receipt.Status == types.ReceiptStatusSuccessful {
		fc.eng.Infow("EVM Chain WriteReport Successful", "txHash", receipt.TxHash.Hex(), "gasUsed", receipt.GasUsed, "fee", transactionFee.String())

		receiverStatus := evmcappb.ReceiverContractExecutionStatus_RECEIVER_CONTRACT_EXECUTION_STATUS_SUCCESS
		response := &evmcappb.WriteReportReply{
			TxStatus:                        evmcappb.TxStatus_TX_STATUS_SUCCESS,
			ReceiverContractExecutionStatus: &receiverStatus,
			TxHash:                          txHash,
			TransactionFee:                  pb.NewBigIntFromInt(transactionFee),
		}
		responseAndMetadata := commonCap.ResponseAndMetadata[*evmcappb.WriteReportReply]{
			Response:         response,
			ResponseMetadata: commonCap.ResponseMetadata{},
		}
		return &responseAndMetadata, nil
	}

	fc.eng.Infow("EVM Chain WriteReport Failed", "txHash", receipt.TxHash.Hex(), "gasUsed", receipt.GasUsed, "fee", transactionFee.String())
	receiverStatus := evmcappb.ReceiverContractExecutionStatus_RECEIVER_CONTRACT_EXECUTION_STATUS_REVERTED
	errorMsg := "Transaction reverted"
	response := &evmcappb.WriteReportReply{
		TxStatus:                        evmcappb.TxStatus_TX_STATUS_REVERTED,
		ReceiverContractExecutionStatus: &receiverStatus,
		TxHash:                          txHash,
		TransactionFee:                  pb.NewBigIntFromInt(transactionFee),
		ErrorMessage:                    &errorMsg,
	}
	responseAndMetadata := commonCap.ResponseAndMetadata[*evmcappb.WriteReportReply]{
		Response:         response,
		ResponseMetadata: commonCap.ResponseMetadata{},
	}
	return &responseAndMetadata, nil
}

func (fc *FakeEVMChain) RegisterLogTrigger(ctx context.Context, triggerID string, metadata commonCap.RequestMetadata, input *evmcappb.FilterLogTriggerRequest) (<-chan commonCap.TriggerAndId[*evmcappb.Log], caperrors.Error) {
	fc.callbackCh[triggerID] = make(chan commonCap.TriggerAndId[*evmcappb.Log])
	fc.logTriggerFilters[triggerID] = input
	return fc.callbackCh[triggerID], nil
}

func (fc *FakeEVMChain) UnregisterLogTrigger(ctx context.Context, triggerID string, metadata commonCap.RequestMetadata, input *evmcappb.FilterLogTriggerRequest) caperrors.Error {
	delete(fc.logTriggerFilters, triggerID)
	delete(fc.callbackCh, triggerID)
	return nil
}

func (fc *FakeEVMChain) AckEvent(ctx context.Context, triggerID string, eventID string, method string) caperrors.Error {
	return nil
}

func (fc *FakeEVMChain) ManualTrigger(ctx context.Context, triggerID string, log *evmcappb.Log) error {
	fc.eng.Debugf("ManualTrigger: %s", log.String())

	if filter, ok := fc.logTriggerFilters[triggerID]; ok && filter != nil {
		if err := fakeEVMLogMatchesFilter(log, filter); err != nil {
			return fmt.Errorf("log does not match registered filter for trigger %s: %w", triggerID, err)
		}
	}

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

// fakeEVMLogMatchesFilter checks whether log satisfies the FilterLogTriggerRequest
// registered for a trigger. It mirrors production EVM log filter semantics:
//   - Address: if Addresses is non-empty, log.Address must equal one of them (OR).
//   - Topics: fixed 4-slot array; slot 0 = event signature, slots 1-3 = indexed args.
//     Within each slot the match is OR (any value matches); across slots it is AND
//     (all non-empty slots must match). An empty Values slice in a slot is a wildcard.
//
// As a developer aid, this fake rejects filters that omit topic0 (the event
// signature hash). Leaving topic0 empty is a common mistake — especially when
// using the raw API or TypeScript bindings — and silently causes the trigger to
// fire for every event emitted by the contract, not just the intended one.
func fakeEVMLogMatchesFilter(log *evmcappb.Log, filter *evmcappb.FilterLogTriggerRequest) error {
	topics := filter.GetTopics()
	if len(topics) == 0 || len(topics[0].GetValues()) == 0 {
		return errors.New("filter is missing topic0 (event signature hash): " +
			"omitting topic0 would match every event emitted by the contract; " +
			"set Topics[0] to the keccak256 hash of the event signature")
	}

	if len(filter.GetAddresses()) > 0 {
		addrMatched := false
		for _, addr := range filter.GetAddresses() {
			if bytes.Equal(log.GetAddress(), addr) {
				addrMatched = true
				break
			}
		}
		if !addrMatched {
			return fmt.Errorf("log address %s does not match any of the addresses in the filter", log.GetAddress())
		}
	}

	logTopics := log.GetTopics()
	for i, topicValues := range filter.GetTopics() {
		if len(topicValues.GetValues()) == 0 {
			continue // wildcard slot (only valid for slots 1-3, slot 0 is checked above)
		}
		if i >= len(logTopics) {
			return fmt.Errorf("log topics length %d does not match the filter topics length %d", len(logTopics), len(filter.GetTopics()))
		}
		slotMatched := false
		for _, v := range topicValues.GetValues() {
			if bytes.Equal(logTopics[i], v) {
				slotMatched = true
				break
			}
		}
		if !slotMatched {
			return fmt.Errorf("log topic %d does not match any of the values in the filter", i)
		}
	}

	return nil
}

func (fc *FakeEVMChain) createManualTriggerEvent(log *evmcappb.Log) commonCap.TriggerAndId[*evmcappb.Log] {
	return commonCap.TriggerAndId[*evmcappb.Log]{
		Trigger: log,
		Id:      "manual-evm-chain-trigger-id",
	}
}

func (fc *FakeEVMChain) FilterLogs(ctx context.Context, metadata commonCap.RequestMetadata, input *evmcappb.FilterLogsRequest) (*commonCap.ResponseAndMetadata[*evmcappb.FilterLogsReply], caperrors.Error) {
	fc.eng.Infow("EVM Chain FilterLogs Started", "input", input)

	if input == nil {
		return nil, caperrors.NewPublicSystemError(errors.New("FilterLogsRequest is nil"), caperrors.Unknown)
	}

	// Prepare filter query
	filterQueryPb := input.GetFilterQuery()
	if filterQueryPb == nil {
		return nil, caperrors.NewPublicSystemError(errors.New("FilterQuery is nil"), caperrors.Unknown)
	}
	addresses := make([]common.Address, len(filterQueryPb.Addresses))
	for i, address := range filterQueryPb.Addresses {
		addresses[i] = common.Address(address)
	}

	// Convert and validate FromBlock/ToBlock using pb.NewIntFromBigInt to preserve sign
	var fromBlock, toBlock *big.Int
	if filterQueryPb.FromBlock != nil {
		fromBlock = pb.NewIntFromBigInt(filterQueryPb.FromBlock)
	}
	if filterQueryPb.ToBlock != nil {
		toBlock = pb.NewIntFromBigInt(filterQueryPb.ToBlock)
	}

	// Validate block numbers to match real capability behavior
	if err := validateBlockNumber(fromBlock, "fromBlock"); err != nil {
		return nil, caperrors.NewPublicSystemError(err, caperrors.Unknown)
	}
	if err := validateBlockNumber(toBlock, "toBlock"); err != nil {
		return nil, caperrors.NewPublicSystemError(err, caperrors.Unknown)
	}
	if fromBlock != nil && toBlock != nil && fromBlock.Sign() > 0 && toBlock.Sign() > 0 {
		if new(big.Int).Sub(toBlock, fromBlock).Sign() < 0 {
			return nil, caperrors.NewPublicSystemError(fmt.Errorf("toBlock %s is less than fromBlock %s", toBlock.String(), fromBlock.String()), caperrors.Unknown)
		}
	}

	// Convert topics from protobuf []*Topics to geth [][]common.Hash
	topics := make([][]common.Hash, 0, len(filterQueryPb.Topics))
	for _, topicSet := range filterQueryPb.Topics {
		if topicSet == nil {
			topics = append(topics, nil)
			continue
		}
		hashes := make([]common.Hash, len(topicSet.Topic))
		for j, t := range topicSet.Topic {
			hashes[j] = common.BytesToHash(t)
		}
		topics = append(topics, hashes)
	}

	filterQuery := ethereum.FilterQuery{
		FromBlock: fromBlock,
		ToBlock:   toBlock,
		Addresses: addresses,
		Topics:    topics,
	}

	// Filter logs
	logs, err := fc.gethClient.FilterLogs(ctx, filterQuery)
	if err != nil {
		return nil, caperrors.NewPublicSystemError(err, caperrors.Unknown)
	}

	fc.eng.Infow("EVM Chain FilterLogs Finished", "logs", logs)

	// Convert logs to protobuf
	logsPb := make([]*evmcappb.Log, len(logs))
	for i, log := range logs {
		logTopics := make([][]byte, len(log.Topics))
		for j, t := range log.Topics {
			logTopics[j] = t.Bytes()
		}
		var eventSig []byte
		if len(log.Topics) > 0 {
			eventSig = log.Topics[0].Bytes()
		}
		logsPb[i] = &evmcappb.Log{
			Address:     log.Address.Bytes(),
			Data:        log.Data,
			Topics:      logTopics,
			EventSig:    eventSig,
			BlockNumber: pb.NewBigIntFromInt(new(big.Int).SetUint64(log.BlockNumber)),
			BlockHash:   log.BlockHash.Bytes(),
			TxHash:      log.TxHash.Bytes(),
			Index:       uint32(log.Index), //nolint:gosec // log index will never exceed uint32
			Removed:     log.Removed,
		}
	}
	response := &evmcappb.FilterLogsReply{
		Logs: logsPb,
	}
	responseAndMetadata := commonCap.ResponseAndMetadata[*evmcappb.FilterLogsReply]{
		Response:         response,
		ResponseMetadata: commonCap.ResponseMetadata{},
	}
	return &responseAndMetadata, nil
}

func (fc *FakeEVMChain) BalanceAt(ctx context.Context, metadata commonCap.RequestMetadata, input *evmcappb.BalanceAtRequest) (*commonCap.ResponseAndMetadata[*evmcappb.BalanceAtReply], caperrors.Error) {
	fc.eng.Infow("EVM Chain BalanceAt Started", "input", input)

	if input == nil {
		return nil, caperrors.NewPublicSystemError(errors.New("BalanceAtRequest is nil"), caperrors.Unknown)
	}

	// Prepare balance at request
	address := common.Address(input.Account)

	// Convert proto big-int to *big.Int; nil ⇒ latest (handled by geth toBlockNumArg)
	blockArg := pb.NewIntFromBigInt(input.BlockNumber)

	// Get balance at block number
	balance, err := fc.gethClient.BalanceAt(ctx, address, blockArg)
	if err != nil {
		return nil, caperrors.NewPublicSystemError(err, caperrors.Unknown)
	}

	// Convert balance to protobuf
	response := &evmcappb.BalanceAtReply{
		Balance: pb.NewBigIntFromInt(balance),
	}
	responseAndMetadata := commonCap.ResponseAndMetadata[*evmcappb.BalanceAtReply]{
		Response:         response,
		ResponseMetadata: commonCap.ResponseMetadata{},
	}
	return &responseAndMetadata, nil
}

func (fc *FakeEVMChain) EstimateGas(ctx context.Context, metadata commonCap.RequestMetadata, input *evmcappb.EstimateGasRequest) (*commonCap.ResponseAndMetadata[*evmcappb.EstimateGasReply], caperrors.Error) {
	fc.eng.Infow("EVM Chain EstimateGas Started", "input", input)

	if input == nil {
		return nil, caperrors.NewPublicSystemError(errors.New("EstimateGasRequest is nil"), caperrors.Unknown)
	}

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
		return nil, caperrors.NewPublicSystemError(err, caperrors.Unknown)
	}

	// Convert gas to protobuf
	fc.eng.Infow("EVM Chain EstimateGas Finished", "gas", gas)
	response := &evmcappb.EstimateGasReply{
		Gas: gas,
	}
	responseAndMetadata := commonCap.ResponseAndMetadata[*evmcappb.EstimateGasReply]{
		Response:         response,
		ResponseMetadata: commonCap.ResponseMetadata{},
	}
	return &responseAndMetadata, nil
}

func (fc *FakeEVMChain) GetTransactionByHash(ctx context.Context, metadata commonCap.RequestMetadata, input *evmcappb.GetTransactionByHashRequest) (*commonCap.ResponseAndMetadata[*evmcappb.GetTransactionByHashReply], caperrors.Error) {
	fc.eng.Infow("EVM Chain GetTransactionByHash Started", "input", input)

	if input == nil {
		return nil, caperrors.NewPublicSystemError(errors.New("GetTransactionByHashRequest is nil"), caperrors.Unknown)
	}

	// Prepare get transaction by hash request
	hash := common.Hash(input.Hash)

	// Get transaction by hash
	transaction, pending, err := fc.gethClient.TransactionByHash(ctx, hash)
	if err != nil {
		return nil, caperrors.NewPublicSystemError(err, caperrors.Unknown)
	}

	fc.eng.Infow("EVM Chain GetTransactionByHash Finished", "transaction", transaction, "pending", pending)

	// Handle nil To() for contract creation transactions
	var toBytes []byte
	if transaction.To() != nil {
		toBytes = transaction.To().Bytes()
	}

	// Convert transaction to protobuf
	transactionPb := &evmcappb.Transaction{
		To:       toBytes,
		Data:     transaction.Data(),
		Hash:     transaction.Hash().Bytes(),
		Value:    pb.NewBigIntFromInt(transaction.Value()),
		GasPrice: pb.NewBigIntFromInt(transaction.GasPrice()),
		Nonce:    transaction.Nonce(),
	}
	response := &evmcappb.GetTransactionByHashReply{
		Transaction: transactionPb,
	}
	responseAndMetadata := commonCap.ResponseAndMetadata[*evmcappb.GetTransactionByHashReply]{
		Response:         response,
		ResponseMetadata: commonCap.ResponseMetadata{},
	}
	return &responseAndMetadata, nil
}

func (fc *FakeEVMChain) GetTransactionReceipt(ctx context.Context, metadata commonCap.RequestMetadata, input *evmcappb.GetTransactionReceiptRequest) (*commonCap.ResponseAndMetadata[*evmcappb.GetTransactionReceiptReply], caperrors.Error) {
	fc.eng.Infow("EVM Chain GetTransactionReceipt Started", "input", input)

	if input == nil {
		return nil, caperrors.NewPublicSystemError(errors.New("GetTransactionReceiptRequest is nil"), caperrors.Unknown)
	}

	// Prepare get transaction receipt request
	hash := common.Hash(input.Hash)

	// Get transaction receipt
	receipt, err := fc.gethClient.TransactionReceipt(ctx, hash)
	if err != nil {
		return nil, caperrors.NewPublicSystemError(err, caperrors.Unknown)
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
		topics := make([][]byte, len(log.Topics))
		for j, t := range log.Topics {
			topics[j] = t.Bytes()
		}
		var eventSig []byte
		if len(log.Topics) > 0 {
			eventSig = log.Topics[0].Bytes()
		}
		receiptPb.Logs[i] = &evmcappb.Log{
			Address:     log.Address.Bytes(),
			Data:        log.Data,
			Topics:      topics,
			EventSig:    eventSig,
			BlockNumber: pb.NewBigIntFromInt(new(big.Int).SetUint64(log.BlockNumber)),
			BlockHash:   log.BlockHash.Bytes(),
			TxHash:      log.TxHash.Bytes(),
			Index:       uint32(log.Index), //nolint:gosec // log index will never exceed uint32
			Removed:     log.Removed,
		}
	}
	response := &evmcappb.GetTransactionReceiptReply{
		Receipt: receiptPb,
	}
	responseAndMetadata := commonCap.ResponseAndMetadata[*evmcappb.GetTransactionReceiptReply]{
		Response:         response,
		ResponseMetadata: commonCap.ResponseMetadata{},
	}
	return &responseAndMetadata, nil
}

func (fc *FakeEVMChain) HeaderByNumber(ctx context.Context, metadata commonCap.RequestMetadata, input *evmcappb.HeaderByNumberRequest) (*commonCap.ResponseAndMetadata[*evmcappb.HeaderByNumberReply], caperrors.Error) {
	fc.eng.Infow("EVM Chain HeaderByNumber Started", "input", input)

	var (
		header *types.Header
		err    error
	)

	// Convert the request block number preserving sign.
	var reqNum *big.Int
	if input != nil {
		reqNum = pb.NewIntFromBigInt(input.BlockNumber)
	}

	// Enforce int64 constraint
	if reqNum != nil && !reqNum.IsInt64() {
		return nil, caperrors.NewPublicSystemError(fmt.Errorf("block number %s is larger than int64: %w", reqNum.String(), ethereum.NotFound), caperrors.Unknown)
	}

	switch {
	// latest block (nil or explicit "latest"): nil or -2
	case reqNum == nil || reqNum.Int64() == rpc.LatestBlockNumber.Int64():
		header, err = fc.gethClient.HeaderByNumber(ctx, nil)

	// non-special, non-negative block number (including 0): >=0
	case reqNum.Sign() >= 0:
		header, err = fc.gethClient.HeaderByNumber(ctx, reqNum)

	// finalized tag: -3
	case reqNum.Int64() == rpc.FinalizedBlockNumber.Int64():
		header, err = fc.gethClient.HeaderByNumber(ctx, big.NewInt(rpc.FinalizedBlockNumber.Int64()))

	// safe tag: -4
	case reqNum.Int64() == rpc.SafeBlockNumber.Int64():
		header, err = fc.gethClient.HeaderByNumber(ctx, big.NewInt(rpc.SafeBlockNumber.Int64()))

	// any other negative is unexpected
	default:
		return nil, caperrors.NewPublicSystemError(fmt.Errorf("unexpected block number %s: %w", reqNum.String(), ethereum.NotFound), caperrors.Unknown)
	}

	if err != nil {
		return nil, caperrors.NewPublicSystemError(err, caperrors.Unknown)
	}
	if header == nil {
		return nil, caperrors.NewPublicSystemError(ethereum.NotFound, caperrors.Unknown)
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
	responseAndMetadata := commonCap.ResponseAndMetadata[*evmcappb.HeaderByNumberReply]{
		Response:         headerPb,
		ResponseMetadata: commonCap.ResponseMetadata{},
	}
	return &responseAndMetadata, nil
}

func (fc *FakeEVMChain) Name() string {
	return fc.ID
}

// validateBlockNumber checks that a block number is valid, matching the real capability's normalizeBlockNumber.
// nil is valid (latest), positive is valid, special tags (-1 latest, -2 finalized, -3 safe) are valid and
// passed through to geth. Zero and other negatives are rejected.
func validateBlockNumber(blockNumber *big.Int, name string) error {
	if blockNumber == nil {
		return nil
	}
	if !blockNumber.IsInt64() {
		return fmt.Errorf("block number %s is not an int64", blockNumber)
	}
	bn := blockNumber.Int64()
	if bn > 0 {
		return nil
	}
	switch rpc.BlockNumber(bn) {
	case rpc.LatestBlockNumber, rpc.SafeBlockNumber, rpc.FinalizedBlockNumber:
		return nil
	default:
		return fmt.Errorf("block number %d is not supported", bn)
	}
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

// dryRunWriteReport simulates the report transaction using eth_call without broadcasting.
func (fc *FakeEVMChain) dryRunWriteReport(
	ctx context.Context,
	from common.Address,
	input *evmcappb.WriteReportRequest,
	signatures [][]byte,
) (*commonCap.ResponseAndMetadata[*evmcappb.WriteReportReply], error) {
	fc.eng.Infow("EVM Chain WriteReport Dry-Run Enabled")
	contractABI, err := abi.JSON(strings.NewReader(MockKeystoneForwarderABI))
	if err != nil {
		return nil, err
	}
	calldata, err := contractABI.Pack(
		"report",
		common.Address(input.Receiver),
		input.Report.RawReport,
		input.Report.ReportContext,
		signatures,
	)
	if err != nil {
		return nil, err
	}

	msg := ethereum.CallMsg{
		From: from,
		To:   &fc.mockKeystoneForwarderAddress,
		Data: calldata,
	}
	_, err = fc.gethClient.CallContract(ctx, msg, nil)
	if err != nil {
		fc.eng.Infow("EVM Chain WriteReport Dry-Run Reverted", "error", err)
		receiverStatus := evmcappb.ReceiverContractExecutionStatus_RECEIVER_CONTRACT_EXECUTION_STATUS_REVERTED
		errMsg := err.Error()
		response := &evmcappb.WriteReportReply{
			TxStatus:                        evmcappb.TxStatus_TX_STATUS_REVERTED,
			ReceiverContractExecutionStatus: &receiverStatus,
			TransactionFee:                  pb.NewBigIntFromInt(big.NewInt(0)),
			ErrorMessage:                    &errMsg,
		}
		responseAndMetadata := commonCap.ResponseAndMetadata[*evmcappb.WriteReportReply]{
			Response:         response,
			ResponseMetadata: commonCap.ResponseMetadata{},
		}
		return &responseAndMetadata, nil
	}

	fc.eng.Infow("EVM Chain WriteReport Dry-Run Successful")
	receiverStatus := evmcappb.ReceiverContractExecutionStatus_RECEIVER_CONTRACT_EXECUTION_STATUS_SUCCESS
	response := &evmcappb.WriteReportReply{
		TxStatus:                        evmcappb.TxStatus_TX_STATUS_SUCCESS,
		ReceiverContractExecutionStatus: &receiverStatus,
		TransactionFee:                  pb.NewBigIntFromInt(big.NewInt(0)),
	}
	responseAndMetadata := commonCap.ResponseAndMetadata[*evmcappb.WriteReportReply]{
		Response:         response,
		ResponseMetadata: commonCap.ResponseMetadata{},
	}
	return &responseAndMetadata, nil
}
