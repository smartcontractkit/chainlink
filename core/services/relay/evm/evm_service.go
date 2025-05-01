package evm

import (
	"context"
	"math/big"
	"strconv"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	gethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"
	evmtypes "github.com/smartcontractkit/chainlink-common/pkg/types/chains/evm"

	"github.com/smartcontractkit/chainlink-common/pkg/types/query"
	"github.com/smartcontractkit/chainlink-common/pkg/types/query/primitives"
	"github.com/smartcontractkit/chainlink-evm/pkg/heads"
	"github.com/smartcontractkit/chainlink-evm/pkg/logpoller"
	"github.com/smartcontractkit/chainlink-evm/pkg/types"
	"github.com/smartcontractkit/chainlink-framework/chains"
)

// Direct RPC
func (r *Relayer) CallContract(ctx context.Context, msg *evmtypes.CallMsg, confidenceLevel primitives.ConfidenceLevel) ([]byte, error) {
	blockNumber, err := blockFromConfidence(ctx, r.chain.HeadTracker(), confidenceLevel)
	if err != nil {
		return nil, err
	}

	return r.chain.Client().CallContract(ctx, toEthMsg(msg), blockNumber)
}

func (r *Relayer) GetLogs(ctx context.Context, filterQuery evmtypes.FilterQuery) ([]*evmtypes.Log, error) {
	logs, err := r.Chain().Client().FilterLogs(ctx, convertEthFilter(filterQuery))
	if err != nil {
		return nil, err
	}

	ret := make([]*evmtypes.Log, 0, len(logs))

	for _, l := range logs {
		ret = append(ret, convertLog(&l))
	}

	return ret, nil
}

func (r *Relayer) BalanceAt(ctx context.Context, account string, blockNumber *big.Int) (*big.Int, error) {
	return r.chain.Client().BalanceAt(ctx, common.HexToAddress(account), blockNumber)
}

func (r *Relayer) EstimateGas(ctx context.Context, call *evmtypes.CallMsg) (uint64, error) {
	return r.chain.Client().EstimateGas(ctx, toEthMsg(call))
}

func (r *Relayer) TransactionByHash(ctx context.Context, hash string) (*evmtypes.Transaction, error) {
	tx, err := r.chain.Client().TransactionByHash(ctx, common.HexToHash(hash))
	if err != nil {
		return nil, err
	}

	return convertTransaction(tx), nil
}

func (r *Relayer) TransactionReceipt(ctx context.Context, txHash string) (*evmtypes.Receipt, error) {
	receipt, err := r.chain.Client().TransactionReceipt(ctx, common.HexToHash(txHash))
	if err != nil {
		return nil, err
	}

	return convertReceipt(receipt), nil
}

// ChainService
func (r *Relayer) GetTransactionFee(ctx context.Context, transactionID string) (*commontypes.TransactionFee, error) {
	return r.chain.TxManager().GetTransactionFee(ctx, transactionID)
}

func (r *Relayer) LatestAndFinalizedHead(ctx context.Context) (commontypes.Head, commontypes.Head, error) {
	latest, finalized, err := r.chain.HeadTracker().LatestAndFinalizedBlock(ctx)
	if err != nil {
		return commontypes.Head{}, commontypes.Head{}, err
	}

	return convertHead(latest), convertHead(finalized), nil
}

func (r *Relayer) QueryLogsFromCache(ctx context.Context, filterQuery []query.Expression,
	limitAndSort query.LimitAndSort, confidenceLevel primitives.ConfidenceLevel) ([]*evmtypes.Log, error) {
	// TODO move evm specific expressions to common
	// TODO check if required filters[event sig and address] are set
	// TODO specify query name BCFR-1328
	// BCFR-1328
	conformations := confidenceToConformations(confidenceLevel)
	filterQuery = append(filterQuery, logpoller.NewConfirmationsFilter(conformations))
	logs, err := r.chain.LogPoller().FilteredLogs(ctx, filterQuery, limitAndSort, "")

	if err != nil {
		return nil, err
	}

	return convertLPLogs(logs), nil
}

func (r *Relayer) SubscribeLogTrigger(ctx context.Context, filterQuery []query.Expression) (chan<- *evmtypes.Log, error) {
	// TODO BCFR-1332
	return nil, errors.New("unimplemented")
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

func (r *Relayer) GetTransactionStatus(ctx context.Context, transactionID string) (commontypes.TransactionStatus, error) {
	return r.chain.TxManager().GetTransactionStatus(ctx, transactionID)
}

func blockFromConfidence(ctx context.Context, ht heads.Tracker, confidence primitives.ConfidenceLevel) (*big.Int, error) {
	latest, finalized, err := ht.LatestAndFinalizedBlock(ctx)
	if err != nil {
		return nil, err
	}
	if confidence == primitives.Finalized {
		return big.NewInt(finalized.BlockNumber()), nil
	}

	return big.NewInt(latest.BlockNumber()), nil
}

func convertHead[H chains.Head[BLOCK_HASH], BLOCK_HASH chains.Hashable](h H) commontypes.Head {
	return commontypes.Head{
		Timestamp: uint64(h.GetTimestamp().Unix()),
		Hash:      h.BlockHash().Bytes(),
		Height:    strconv.FormatInt(h.BlockNumber(), 10),
	}
}

func convertReceipt(r *gethtypes.Receipt) *evmtypes.Receipt {
	return &evmtypes.Receipt{
		Status:            r.Status,
		Logs:              convertLogs(r.Logs),
		TxHash:            r.TxHash.Hex(),
		ContractAddress:   r.ContractAddress.Hex(),
		GasUsed:           r.GasUsed,
		BlockHash:         r.BlockHash.Hex(),
		BlockNumber:       r.BlockNumber,
		TransactionIndex:  uint64(r.TransactionIndex),
		EffectiveGasPrice: r.EffectiveGasPrice,
	}
}

func convertEthFilter(q evmtypes.FilterQuery) ethereum.FilterQuery {
	addresses := stringsToAddresses(q.Addresses)
	topics := stringsToHashMatrix(q.Topics)

	return ethereum.FilterQuery{
		FromBlock: q.FromBlock,
		ToBlock:   q.ToBlock,
		Addresses: addresses,
		Topics:    topics,
	}
}

var errEmptyFilterName = errors.New("filter name can't be empty")

func convertLPFilter(q evmtypes.LPFilterQuery) (logpoller.Filter, error) {
	if q.Name == "" {
		return logpoller.Filter{}, errEmptyFilterName
	}
	return logpoller.Filter{
		Name:         q.Name,
		Addresses:    stringsToAddresses(q.Addresses),
		EventSigs:    stringsToHashes(q.EventSigs),
		Topic2:       stringsToHashes(q.Topic2),
		Topic3:       stringsToHashes(q.Topic3),
		Topic4:       stringsToHashes(q.Topic4),
		Retention:    q.Retention,
		MaxLogsKept:  q.MaxLogsKept,
		LogsPerBlock: q.LogsPerBlock,
	}, nil
}

func convertTransaction(tx *gethtypes.Transaction) *evmtypes.Transaction {
	var to string
	if tx.To() != nil {
		to = tx.To().Hex()
	}

	return &evmtypes.Transaction{
		To:       to,
		Data:     tx.Data(),
		Hash:     tx.Hash().Hex(),
		Nonce:    tx.Nonce(),
		Gas:      tx.Gas(),
		GasPrice: tx.GasPrice(),
		Value:    tx.Value(),
	}
}

func stringsToHashMatrix(input []string) [][]common.Hash {
	result := make([][]common.Hash, 0, len(input))
	for _, row := range input {
		result = append(result, []common.Hash{common.HexToHash(row)})
	}
	return result
}

func stringsToAddresses(input []string) []common.Address {
	res := make([]common.Address, 0, len(input))
	for _, s := range input {
		addr := common.HexToAddress(s)
		res = append(res, addr)
	}

	return res
}

func stringsToHashes(input []string) []common.Hash {
	res := make([]common.Hash, 0, len(input))
	for _, s := range input {
		hash := common.HexToHash(s)
		res = append(res, hash)
	}

	return res
}

func toEthMsg(msg *evmtypes.CallMsg) ethereum.CallMsg {
	to := common.HexToAddress(msg.To)
	from := common.HexToAddress(msg.From)

	return ethereum.CallMsg{
		From: from,
		To:   &to,
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
	topics := make([]string, len(log.Topics))
	for i, topic := range log.Topics {
		topics[i] = topic.Hex()
	}

	var eventSig string
	if len(log.Topics) > 0 {
		eventSig = log.Topics[0].Hex()
	}

	return &evmtypes.Log{
		LogIndex:    uint32(log.Index),
		BlockHash:   log.BlockHash.Hex(),
		BlockNumber: new(big.Int).SetUint64(log.BlockNumber),
		Topics:      topics,
		EventSig:    eventSig,
		Address:     log.Address.Hex(),
		TxHash:      log.TxHash.Hex(),
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
