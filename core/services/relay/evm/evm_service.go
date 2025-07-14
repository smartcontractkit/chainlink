package evm

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	gethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink-common/pkg/types/chains/evm"

	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"
	evmtypes "github.com/smartcontractkit/chainlink-common/pkg/types/chains/evm"

	"github.com/smartcontractkit/chainlink-common/pkg/types/query"
	"github.com/smartcontractkit/chainlink-common/pkg/types/query/primitives"
	evmprimitives "github.com/smartcontractkit/chainlink-common/pkg/types/query/primitives/evm"
	"github.com/smartcontractkit/chainlink-evm/pkg/logpoller"
	evmtxmgr "github.com/smartcontractkit/chainlink-evm/pkg/txmgr"
	"github.com/smartcontractkit/chainlink-evm/pkg/types"
	"github.com/smartcontractkit/chainlink-framework/chains/txmgr"
	txmgrtypes "github.com/smartcontractkit/chainlink-framework/chains/txmgr/types"
)

// Direct RPC
func (r *Relayer) CallContract(ctx context.Context, request evmtypes.CallContractRequest) (*evmtypes.CallContractReply, error) {
	result, err := r.chain.Client().CallContractWithOpts(ctx, toEthMsg(request.Msg), request.BlockNumber, types.CallContractOpts{ConfidenceLevel: request.ConfidenceLevel})
	if err != nil {
		return nil, err
	}

	return &evmtypes.CallContractReply{Data: result}, nil
}

func (r *Relayer) FilterLogs(ctx context.Context, request evmtypes.FilterLogsRequest) (*evmtypes.FilterLogsReply, error) {
	rawLogs, err := r.chain.Client().FilterLogsWithOpts(ctx, convertEthFilter(request.FilterQuery), types.FilterLogsOpts{ConfidenceLevel: request.ConfidenceLevel})
	if err != nil {
		return nil, err
	}

	logs := make([]*evmtypes.Log, 0, len(rawLogs))
	for _, l := range rawLogs {
		logs = append(logs, convertLog(&l))
	}

	return &evmtypes.FilterLogsReply{Logs: logs}, nil
}

func (r *Relayer) BalanceAt(ctx context.Context, request evmtypes.BalanceAtRequest) (*evmtypes.BalanceAtReply, error) {
	balance, err := r.chain.Client().BalanceAtWithOpts(ctx, request.Address, request.BlockNumber, types.BalanceAtOpts{ConfidenceLevel: request.ConfidenceLevel})
	if err != nil {
		return nil, err
	}

	return &evmtypes.BalanceAtReply{Balance: balance}, nil
}

func (r *Relayer) EstimateGas(ctx context.Context, call *evmtypes.CallMsg) (uint64, error) {
	return r.chain.Client().EstimateGas(ctx, toEthMsg(call))
}

func (r *Relayer) GetTransactionByHash(ctx context.Context, hash evmtypes.Hash) (*evmtypes.Transaction, error) {
	tx, err := r.chain.Client().TransactionByHash(ctx, hash)
	if err != nil {
		return nil, err
	}

	return convertTransaction(tx), nil
}

func (r *Relayer) GetTransactionReceipt(ctx context.Context, txHash evmtypes.Hash) (*evmtypes.Receipt, error) {
	receipt, err := r.chain.Client().TransactionReceipt(ctx, txHash)
	if err != nil {
		return nil, err
	}

	return convertReceipt(receipt), nil
}

// ChainService
func (r *Relayer) GetTransactionFee(ctx context.Context, transactionID commontypes.IdempotencyKey) (*evmtypes.TransactionFee, error) {
	return r.chain.TxManager().GetTransactionFee(ctx, transactionID)
}

func (r *Relayer) HeaderByNumber(ctx context.Context, request evmtypes.HeaderByNumberRequest) (*evmtypes.HeaderByNumberReply, error) {
	var err error
	var h *types.Head
	switch {
	// latest block
	case request.Number == nil || request.Number.Int64() == rpc.LatestBlockNumber.Int64():
		h, _, err = r.chain.HeadTracker().LatestAndFinalizedBlock(ctx)
		// non-special block or larger that int64
	case request.Number.Sign() >= 0 || request.Number.IsInt64():
		var header *types.Header
		header, err = r.chain.Client().HeaderByNumberWithOpts(ctx, request.Number, types.HeaderByNumberOpts{ConfidenceLevel: request.ConfidenceLevel})
		h = (*types.Head)(header)
	case request.Number.Int64() == rpc.FinalizedBlockNumber.Int64():
		_, h, err = r.chain.HeadTracker().LatestAndFinalizedBlock(ctx)
	case request.Number.Int64() == rpc.SafeBlockNumber.Int64():
		h, err = r.chain.HeadTracker().LatestSafeBlock(ctx)
	default:
		return nil, fmt.Errorf("unexpected block number %s: %w", request.Number.String(), ethereum.NotFound)
	}

	if err != nil {
		return nil, err
	}

	if h == nil {
		return nil, ethereum.NotFound
	}

	header := convertHead(h)
	return &evmtypes.HeaderByNumberReply{Header: header}, nil
}

// TODO introduce parameters validation PLEX-1437
func (r *Relayer) QueryTrackedLogs(ctx context.Context, filterQuery []query.Expression,
	limitAndSort query.LimitAndSort, confidenceLevel primitives.ConfidenceLevel,
) ([]*evmtypes.Log, error) {
	conformations := confidenceToConformations(confidenceLevel)
	filterQuery = append(filterQuery, logpoller.NewConfirmationsFilter(conformations))
	queryName := queryNameFromFilter(filterQuery)
	logs, err := r.chain.LogPoller().FilteredLogs(ctx, filterQuery, limitAndSort, queryName)
	if err != nil {
		return nil, err
	}

	return convertLPLogs(logs), nil
}

func (r *Relayer) RegisterLogTracking(ctx context.Context, filter evmtypes.LPFilterQuery) error {
	lpfilter, err := convertLPFilter(filter)
	if err != nil {
		return err
	}
	if r.chain.LogPoller().HasFilter(lpfilter.Name) {
		return nil
	}

	return r.chain.LogPoller().RegisterFilter(ctx, lpfilter)
}

func (r *Relayer) UnregisterLogTracking(ctx context.Context, filterName string) error {
	if filterName == "" {
		return errEmptyFilterName
	}
	if !r.chain.LogPoller().HasFilter(filterName) {
		return nil
	}

	return r.chain.LogPoller().UnregisterFilter(ctx, filterName)
}

func (r *Relayer) GetTransactionStatus(ctx context.Context, transactionID commontypes.IdempotencyKey) (commontypes.TransactionStatus, error) {
	status, err := r.chain.TxManager().GetTransactionStatus(ctx, transactionID)
	if err != nil {
		return commontypes.Unknown, err
	}

	return status, nil
}

func (r *Relayer) SubmitTransaction(ctx context.Context, txRequest evmtypes.SubmitTransactionRequest) (*evmtypes.TransactionResult, error) {
	config := r.chain.Config()

	fromAddress := config.EVM().Workflow().FromAddress().Address()
	var gasLimit uint64
	if txRequest.GasConfig != nil && txRequest.GasConfig.GasLimit != nil {
		gasLimit = *txRequest.GasConfig.GasLimit
	}

	uuid, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}
	txID := uuid.String()

	// PLEX-1524 - review which transmitter checker we should use
	var checker evmtxmgr.TransmitCheckerSpec
	checker.CheckerType = evmtxmgr.TransmitCheckerTypeSimulate
	value := big.NewInt(0)

	// PLEX-1524 - Define how we should properly get the workflow execution ID into the meta without making the API CRE specific.
	var txMeta *txmgrtypes.TxMeta[common.Address, common.Hash]
	txmReq := evmtxmgr.TxRequest{
		FromAddress:    fromAddress,
		ToAddress:      txRequest.To,
		EncodedPayload: txRequest.Data,
		FeeLimit:       gasLimit,
		Meta:           txMeta,
		IdempotencyKey: &txID,
		// PLEX-1524 - Review strategy to be used.
		Strategy: txmgr.NewSendEveryStrategy(),
		Checker:  checker,
		Value:    *value,
	}

	_, err = r.chain.TxManager().CreateTransaction(ctx, txmReq)
	if err != nil {
		return nil, fmt.Errorf("%w; failed to create tx", err)
	}

	maximumWaitTimeForConfirmation := config.EVM().ConfirmationTimeout()
	start := time.Now()
StatusCheckingLoop:
	for {
		txStatus, txStatusErr := r.chain.TxManager().GetTransactionStatus(ctx, txID)
		if txStatusErr != nil {
			return nil, txStatusErr
		}
		switch txStatus {
		case commontypes.Fatal, commontypes.Failed:
			return &evmtypes.TransactionResult{
				TxStatus: evm.TxFatal,
				TxHash:   evmtypes.Hash{},
			}, nil

		case commontypes.Unconfirmed, commontypes.Finalized:
			break StatusCheckingLoop
		case commontypes.Pending, commontypes.Unknown:
		default:
			return nil, fmt.Errorf("unexpected transaction status %d for tx with ID %s", txStatus, txID)
		}
		if time.Since(start) > maximumWaitTimeForConfirmation {
			return nil, errors.Errorf("Wait time for Tx %s to get confirmed was greater than maximum wait time %d", txID, maximumWaitTimeForConfirmation)
		}
		// PLEX-1524 - Use ticker instead of time.Sleep and make the time configurable
		time.Sleep(100 * time.Millisecond)
	}

	receipt, err := r.chain.TxManager().GetTransactionReceipt(ctx, txID)
	if err != nil {
		return nil, fmt.Errorf("failed to get TX receipt for tx with ID %s: %w", txID, err)
	}
	if receipt == nil {
		return nil, fmt.Errorf("receipt was nil for TX with ID %s: %w", txID, err)
	}

	return &evmtypes.TransactionResult{
		TxStatus: evm.TxSuccess,
		TxHash:   (*receipt).GetTxHash(),
	}, nil
}

func (r *Relayer) CalculateTransactionFee(ctx context.Context, receipt evm.ReceiptGasInfo) (*evm.TransactionFee, error) {
	txFee := r.chain.TxManager().CalculateFee(txmgr.FeeParts{
		GasUsed:           receipt.GasUsed,
		EffectiveGasPrice: receipt.EffectiveGasPrice,
	})
	return &evmtypes.TransactionFee{
		TransactionFee: txFee,
	}, nil
}

func (r *Relayer) GetForwarderForEOA(ctx context.Context, eoa, ocr2AggregatorID evm.Address, pluginType string) (forwarder evm.Address, err error) {
	if pluginType == string(commontypes.Median) {
		return r.chain.TxManager().GetForwarderForEOAOCR2Feeds(ctx, eoa, ocr2AggregatorID)
	}
	return r.chain.TxManager().GetForwarderForEOA(ctx, eoa)
}

func queryNameFromFilter(filterQuery []query.Expression) string {
	var address string
	var eventSig string

	for _, expr := range filterQuery {
		if expr.IsPrimitive() {
			switch primitive := expr.Primitive.(type) {
			case *evmprimitives.Address:
				address = common.Address(primitive.Address).Hex()
			case *evmprimitives.EventSig:
				eventSig = common.Hash(primitive.EventSig).Hex()
			}
		}
	}

	return address + "-" + eventSig
}

func convertHead(h *types.Head) *evmtypes.Header {
	return &evmtypes.Header{
		Timestamp:  uint64(h.GetTimestamp().Unix()),
		Hash:       bytesToHash(h.BlockHash().Bytes()),
		Number:     big.NewInt(h.BlockNumber()),
		ParentHash: bytesToHash(h.GetParentHash().Bytes()),
	}
}

func convertReceipt(r *gethtypes.Receipt) *evmtypes.Receipt {
	return &evmtypes.Receipt{
		Status:            r.Status,
		Logs:              convertLogs(r.Logs),
		TxHash:            r.TxHash,
		ContractAddress:   r.ContractAddress,
		GasUsed:           r.GasUsed,
		BlockHash:         r.BlockHash,
		BlockNumber:       r.BlockNumber,
		TransactionIndex:  uint64(r.TransactionIndex),
		EffectiveGasPrice: r.EffectiveGasPrice,
	}
}

func convertEthFilter(q evmtypes.FilterQuery) ethereum.FilterQuery {
	return ethereum.FilterQuery{
		FromBlock: q.FromBlock,
		ToBlock:   q.ToBlock,
		Addresses: arraysToAddresses(q.Addresses),
		Topics:    arraysToHashMatrix(q.Topics),
	}
}

var errEmptyFilterName = errors.New("filter name can't be empty")

func convertLPFilter(q evmtypes.LPFilterQuery) (logpoller.Filter, error) {
	if q.Name == "" {
		return logpoller.Filter{}, errEmptyFilterName
	}
	return logpoller.Filter{
		Name:         q.Name,
		Addresses:    arraysToAddresses(q.Addresses),
		EventSigs:    arraysToHashes(q.EventSigs),
		Topic2:       arraysToHashes(q.Topic2),
		Topic3:       arraysToHashes(q.Topic3),
		Topic4:       arraysToHashes(q.Topic4),
		Retention:    q.Retention,
		MaxLogsKept:  q.MaxLogsKept,
		LogsPerBlock: q.LogsPerBlock,
	}, nil
}

func convertTransaction(tx *gethtypes.Transaction) *evmtypes.Transaction {
	var to evm.Address
	if tx.To() != nil {
		to = *tx.To()
	}

	return &evmtypes.Transaction{
		To:       to,
		Data:     tx.Data(),
		Hash:     tx.Hash(),
		Nonce:    tx.Nonce(),
		Gas:      tx.Gas(),
		GasPrice: tx.GasPrice(),
		Value:    tx.Value(),
	}
}

func arraysToHashMatrix(input [][][32]byte) [][]common.Hash {
	result := make([][]common.Hash, 0, len(input))
	for _, row := range input {
		result = append(result, arraysToHashes(row))
	}
	return result
}

func arraysToAddresses(input [][20]byte) []common.Address {
	res := make([]common.Address, 0, len(input))
	for _, s := range input {
		res = append(res, s)
	}

	return res
}

func arraysToHashes(input [][32]byte) []common.Hash {
	res := make([]common.Hash, 0, len(input))
	for _, s := range input {
		res = append(res, s)
	}

	return res
}

func hashesToArrays(input []common.Hash) [][32]byte {
	res := make([][32]byte, 0, len(input))
	for _, s := range input {
		res = append(res, s)
	}

	return res
}

var empty common.Address

func toEthMsg(msg *evmtypes.CallMsg) ethereum.CallMsg {
	var to *common.Address

	if empty.Cmp(msg.To) != 0 {
		to = new(common.Address)
		*to = msg.To
	}

	return ethereum.CallMsg{
		From: msg.From,
		To:   to,
		Data: msg.Data,
	}
}

func convertLogs(logs []*gethtypes.Log) []*evmtypes.Log {
	ret := make([]*evmtypes.Log, 0, len(logs))

	for _, l := range logs {
		ret = append(ret, convertLog(l))
	}

	return ret
}

func convertLPLogs(logs []logpoller.Log) []*evmtypes.Log {
	ret := make([]*evmtypes.Log, 0, len(logs))
	for _, l := range logs {
		gl := l.ToGethLog()
		ret = append(ret, convertLog(&gl))
	}

	return ret
}

func convertLog(log *gethtypes.Log) *evmtypes.Log {
	topics := hashesToArrays(log.Topics)

	var eventSig [32]byte
	if len(log.Topics) > 0 {
		eventSig = log.Topics[0]
	}

	return &evmtypes.Log{
		LogIndex:    uint32(log.Index),
		BlockHash:   log.BlockHash,
		BlockNumber: new(big.Int).SetUint64(log.BlockNumber),
		Topics:      topics,
		EventSig:    eventSig,
		Address:     log.Address,
		TxHash:      log.TxHash,
		Data:        log.Data,
		Removed:     log.Removed,
	}
}

func confidenceToConformations(conf primitives.ConfidenceLevel) types.Confirmations {
	if conf == primitives.Finalized {
		return types.Finalized
	}

	return types.Unconfirmed
}

func bytesToHash(b []byte) (h evm.Hash) {
	copy(h[:], b)
	return
}
