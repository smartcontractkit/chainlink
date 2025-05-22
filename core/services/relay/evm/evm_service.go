package evm

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	gethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"

	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-common/pkg/types/chains/evm"
	evmtypes "github.com/smartcontractkit/chainlink-common/pkg/types/chains/evm"

	"github.com/smartcontractkit/chainlink-common/pkg/types/query"
	"github.com/smartcontractkit/chainlink-common/pkg/types/query/primitives"
	evmprimitives "github.com/smartcontractkit/chainlink-common/pkg/types/query/primitives/evm"
	"github.com/smartcontractkit/chainlink-evm/pkg/logpoller"
	"github.com/smartcontractkit/chainlink-evm/pkg/types"
	"github.com/smartcontractkit/chainlink-framework/chains"
)

// Direct RPC
func (r *Relayer) CallContract(ctx context.Context, msg *evmtypes.CallMsg, blockNumber *big.Int) ([]byte, error) {
	return r.chain.Client().CallContract(ctx, toEthMsg(msg), blockNumber)
}

func (r *Relayer) FilterLogs(ctx context.Context, filterQuery evmtypes.FilterQuery) ([]*evmtypes.Log, error) {
	logs, err := r.chain.Client().FilterLogs(ctx, convertEthFilter(filterQuery))
	if err != nil {
		return nil, err
	}

	ret := make([]*evmtypes.Log, 0, len(logs))

	for _, l := range logs {
		ret = append(ret, convertLog(&l))
	}

	return ret, nil
}

func (r *Relayer) BalanceAt(ctx context.Context, account evmtypes.Address, blockNumber *big.Int) (*big.Int, error) {
	return r.chain.Client().BalanceAt(ctx, account, blockNumber)
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

func (r *Relayer) LatestAndFinalizedHead(ctx context.Context) (evmtypes.Head, evmtypes.Head, error) {
	latest, finalized, err := r.chain.HeadTracker().LatestAndFinalizedBlock(ctx)
	if err != nil {
		return evmtypes.Head{}, evmtypes.Head{}, err
	}

	return convertHead(latest), convertHead(finalized), nil
}

// TODO introduce parameters validation PLEX-1437
func (r *Relayer) QueryTrackedLogs(ctx context.Context, filterQuery []query.Expression,
	limitAndSort query.LimitAndSort, confidenceLevel primitives.ConfidenceLevel) ([]*evmtypes.Log, error) {
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

	return commontypes.TransactionStatus(status), nil
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

func convertHead[H chains.Head[BLOCK_HASH], BLOCK_HASH chains.Hashable](h H) evmtypes.Head {
	return evmtypes.Head{
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
