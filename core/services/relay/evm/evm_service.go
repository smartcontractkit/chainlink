package evm

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	gethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-common/pkg/types/chains/evm"
	"github.com/smartcontractkit/chainlink-common/pkg/types/query"
	"github.com/smartcontractkit/chainlink-common/pkg/types/query/primitives"
	evmprimitives "github.com/smartcontractkit/chainlink-common/pkg/types/query/primitives/evm"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/retry"
	"github.com/smartcontractkit/chainlink-evm/pkg/chains/legacyevm"
	"github.com/smartcontractkit/chainlink-evm/pkg/logpoller"
	evmtxmgr "github.com/smartcontractkit/chainlink-evm/pkg/txmgr"
	"github.com/smartcontractkit/chainlink-evm/pkg/types"
	"github.com/smartcontractkit/chainlink-framework/chains/txmgr"
	txmgrtypes "github.com/smartcontractkit/chainlink-framework/chains/txmgr/types"
)

type evmService struct {
	chain  legacyevm.Chain
	logger logger.Logger
}

// Direct RPC
func (e *evmService) CallContract(ctx context.Context, request evm.CallContractRequest) (*evm.CallContractReply, error) {
	opts := types.CallContractOpts{
		ConfidenceLevel:   request.ConfidenceLevel,
		IsExternalRequest: request.IsExternal,
	}
	if request.Msg == nil {
		return nil, errors.New("request.Msg can not be nil")
	}
	result, err := e.chain.Client().CallContractWithOpts(ctx, toEthMsg(*request.Msg), request.BlockNumber, opts)
	if err != nil {
		return nil, err
	}

	return &evm.CallContractReply{Data: result}, nil
}

func (e *evmService) FilterLogs(ctx context.Context, request evm.FilterLogsRequest) (*evm.FilterLogsReply, error) {
	opts := types.FilterLogsOpts{
		ConfidenceLevel:   request.ConfidenceLevel,
		IsExternalRequest: request.IsExternal,
	}
	rawLogs, err := e.chain.Client().FilterLogsWithOpts(ctx, convertEthFilter(request.FilterQuery), opts)
	if err != nil {
		return nil, err
	}

	logs := make([]*evm.Log, 0, len(rawLogs))
	for _, l := range rawLogs {
		logs = append(logs, convertLog(&l))
	}

	return &evm.FilterLogsReply{Logs: logs}, nil
}

func (e *evmService) BalanceAt(ctx context.Context, request evm.BalanceAtRequest) (*evm.BalanceAtReply, error) {
	balance, err := e.chain.Client().BalanceAtWithOpts(ctx, request.Address, request.BlockNumber, types.BalanceAtOpts{ConfidenceLevel: request.ConfidenceLevel})
	if err != nil {
		return nil, err
	}

	return &evm.BalanceAtReply{Balance: balance}, nil
}

func (e *evmService) EstimateGas(ctx context.Context, call *evm.CallMsg) (uint64, error) {
	if call == nil {
		return 0, errors.New("call can not be nil")
	}
	return e.chain.Client().EstimateGas(ctx, toEthMsg(*call))
}

func (e *evmService) GetTransactionByHash(ctx context.Context, request evm.GetTransactionByHashRequest) (*evm.Transaction, error) {
	tx, err := e.chain.Client().TransactionByHashWithOpts(ctx, request.Hash, types.TransactionByHashOpts{IsExternalRequest: request.IsExternal})
	if err != nil {
		return nil, err
	}

	return convertTransaction(tx), nil
}

func (e *evmService) GetTransactionReceipt(ctx context.Context, request evm.GeTransactionReceiptRequest) (*evm.Receipt, error) {
	receipt, err := e.chain.Client().TransactionReceiptWithOpts(ctx, request.Hash, types.TransactionReceiptOpts{IsExternalRequest: request.IsExternal})
	if err != nil {
		return nil, err
	}

	return convertReceipt(receipt), nil
}

// ChainService
func (e *evmService) GetTransactionFee(ctx context.Context, transactionID commontypes.IdempotencyKey) (*evm.TransactionFee, error) {
	return e.chain.TxManager().GetTransactionFee(ctx, transactionID)
}

func (e *evmService) HeaderByNumber(ctx context.Context, request evm.HeaderByNumberRequest) (*evm.HeaderByNumberReply, error) {
	var err error
	var h *types.Head
	switch {
	case request.Number != nil && !request.Number.IsInt64():
		// chain-evm.RpcClient does not support block numbers larger than int64
		return nil, fmt.Errorf("block number %s is larger than int64: %w", request.Number.String(), ethereum.NotFound)
	// latest block
	case request.Number == nil || request.Number.Int64() == rpc.LatestBlockNumber.Int64():
		h, _, err = e.chain.HeadTracker().LatestAndFinalizedBlock(ctx)
		// non-special block or larger that int64
	case request.Number.Sign() >= 0:
		var header *types.Header
		header, err = e.chain.Client().HeaderByNumberWithOpts(ctx, request.Number, types.HeaderByNumberOpts{ConfidenceLevel: request.ConfidenceLevel})
		h = (*types.Head)(header)
	case request.Number.Int64() == rpc.FinalizedBlockNumber.Int64():
		_, h, err = e.chain.HeadTracker().LatestAndFinalizedBlock(ctx)
	case request.Number.Int64() == rpc.SafeBlockNumber.Int64():
		h, err = e.chain.HeadTracker().LatestSafeBlock(ctx)
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
	return &evm.HeaderByNumberReply{Header: header}, nil
}

// TODO introduce parameters validation PLEX-1437
func (e *evmService) QueryTrackedLogs(ctx context.Context, filterQuery []query.Expression,
	limitAndSort query.LimitAndSort, confidenceLevel primitives.ConfidenceLevel,
) ([]*evm.Log, error) {
	conformations := confidenceToConformations(confidenceLevel)
	filterQuery = append(filterQuery, logpoller.NewConfirmationsFilter(conformations))
	queryName := queryNameFromFilter(filterQuery)
	logs, err := e.chain.LogPoller().FilteredLogs(ctx, filterQuery, limitAndSort, queryName)
	if err != nil {
		return nil, err
	}

	return convertLPLogs(logs), nil
}

func (e *evmService) GetFiltersNames(_ context.Context) ([]string, error) {
	// TODO PLEX-1465: once code is moved away, remove this GetFiltersNames method
	filters := e.chain.LogPoller().GetFilters()
	filterNames := make([]string, 0, len(filters))
	for name := range filters {
		filterNames = append(filterNames, name)
	}
	return filterNames, nil
}

func (e *evmService) RegisterLogTracking(ctx context.Context, filter evm.LPFilterQuery) error {
	lpfilter, err := convertLPFilter(filter)
	if err != nil {
		return err
	}
	if e.chain.LogPoller().HasFilter(lpfilter.Name) {
		return nil
	}

	return e.chain.LogPoller().RegisterFilter(ctx, lpfilter)
}

func (e *evmService) UnregisterLogTracking(ctx context.Context, filterName string) error {
	if filterName == "" {
		return errEmptyFilterName
	}
	if !e.chain.LogPoller().HasFilter(filterName) {
		return nil
	}

	return e.chain.LogPoller().UnregisterFilter(ctx, filterName)
}

func (e *evmService) GetTransactionStatus(ctx context.Context, transactionID commontypes.IdempotencyKey) (commontypes.TransactionStatus, error) {
	status, err := e.chain.TxManager().GetTransactionStatus(ctx, transactionID)
	if err != nil {
		return commontypes.Unknown, err
	}

	return status, nil
}

func (e *evmService) SubmitTransaction(ctx context.Context, txRequest evm.SubmitTransactionRequest) (*evm.TransactionResult, error) {
	config := e.chain.Config()

	fromAddress := config.EVM().Workflow().FromAddress().Address()
	var gasLimit uint64
	if txRequest.GasConfig != nil && txRequest.GasConfig.GasLimit != nil {
		gasLimit = *txRequest.GasConfig.GasLimit
	}

	id, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}
	txID := id.String()
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
		Value:    *value,
	}

	_, err = e.chain.TxManager().CreateTransaction(ctx, txmReq)
	if err != nil {
		return nil, fmt.Errorf("%w; failed to create tx", err)
	}

	maximumWaitTimeForConfirmation := config.EVM().ConfirmationTimeout()

	retryContext, cancel := context.WithTimeout(ctx, maximumWaitTimeForConfirmation)
	defer cancel()
	txStatus, err := retry.Do(retryContext, e.logger, func(ctx context.Context) (evm.TransactionStatus, error) {
		txStatus, txStatusErr := e.chain.TxManager().GetTransactionStatus(ctx, txID)
		if txStatusErr != nil {
			return evm.TxFatal, txStatusErr
		}
		switch txStatus {
		case commontypes.Fatal, commontypes.Failed:
			return evm.TxFatal, nil
		case commontypes.Unconfirmed, commontypes.Finalized:
			return evm.TxSuccess, nil
		case commontypes.Pending, commontypes.Unknown:
			return evm.TxFatal, fmt.Errorf("tx still in state pending or unknown, tx status is %d for tx with ID %s", txStatus, txID)
		default:
			return evm.TxFatal, fmt.Errorf("unexpected transaction status %d for tx with ID %s", txStatus, txID)
		}
	})

	if err != nil {
		return nil, fmt.Errorf("failed getting transaction status. %w", err)
	}

	if txStatus == evm.TxFatal {
		return &evm.TransactionResult{TxStatus: txStatus, TxIdempotencyKey: txID}, nil
	}

	receipt, err := retry.Do(retryContext, e.logger, func(ctx context.Context) (*evmtxmgr.ChainReceipt, error) {
		receipt, receiptErr := e.chain.TxManager().GetTransactionReceipt(ctx, txID)
		if receiptErr != nil {
			return nil, fmt.Errorf("failed to get TX receipt for tx with ID %s: %w", txID, receiptErr)
		}
		if receipt == nil {
			return nil, fmt.Errorf("receipt was nil for TX with ID %s", txID)
		}
		return receipt, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed getting transaction receipt. %w", err)
	}

	return &evm.TransactionResult{
		TxStatus:         evm.TxSuccess,
		TxHash:           (*receipt).GetTxHash(),
		TxIdempotencyKey: txID,
	}, nil
}

func (e *evmService) CalculateTransactionFee(ctx context.Context, receipt evm.ReceiptGasInfo) (*evm.TransactionFee, error) {
	txFee := e.chain.TxManager().CalculateFee(txmgr.FeeParts{
		GasUsed:           receipt.GasUsed,
		EffectiveGasPrice: receipt.EffectiveGasPrice,
	})
	return &evm.TransactionFee{
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

func convertHead(h *types.Head) *evm.Header {
	return &evm.Header{
		Timestamp:  uint64(h.GetTimestamp().Unix()),
		Hash:       bytesToHash(h.BlockHash().Bytes()),
		Number:     big.NewInt(h.BlockNumber()),
		ParentHash: bytesToHash(h.GetParentHash().Bytes()),
	}
}

func convertReceipt(r *gethtypes.Receipt) *evm.Receipt {
	return &evm.Receipt{
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

func convertEthFilter(q evm.FilterQuery) ethereum.FilterQuery {
	return ethereum.FilterQuery{
		FromBlock: q.FromBlock,
		ToBlock:   q.ToBlock,
		Addresses: arraysToAddresses(q.Addresses),
		Topics:    arraysToHashMatrix(q.Topics),
	}
}

var errEmptyFilterName = errors.New("filter name can't be empty")

func convertLPFilter(q evm.LPFilterQuery) (logpoller.Filter, error) {
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

func convertTransaction(tx *gethtypes.Transaction) *evm.Transaction {
	var to evm.Address
	if tx.To() != nil {
		to = *tx.To()
	}

	return &evm.Transaction{
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

func toEthMsg(msg evm.CallMsg) ethereum.CallMsg {
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

func convertLogs(logs []*gethtypes.Log) []*evm.Log {
	ret := make([]*evm.Log, 0, len(logs))

	for _, l := range logs {
		ret = append(ret, convertLog(l))
	}

	return ret
}

func convertLPLogs(logs []logpoller.Log) []*evm.Log {
	ret := make([]*evm.Log, 0, len(logs))
	for _, l := range logs {
		gl := l.ToGethLog()
		ret = append(ret, convertLog(&gl))
	}

	return ret
}

func convertLog(log *gethtypes.Log) *evm.Log {
	topics := hashesToArrays(log.Topics)

	var eventSig [32]byte
	if len(log.Topics) > 0 {
		eventSig = log.Topics[0]
	}

	return &evm.Log{
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
