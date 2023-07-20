package multinodeclient

import (
	"context"
	"fmt"
	"math/big"
	"time"

	feetypes "github.com/smartcontractkit/chainlink/v2/common/fee/types"
	txmgrtypes "github.com/smartcontractkit/chainlink/v2/common/txmgr/types"

	"github.com/smartcontractkit/chainlink/v2/core/config"
	"github.com/smartcontractkit/chainlink/v2/core/logger"

	"github.com/smartcontractkit/chainlink/v2/common/types"
)

type MultiNodeClient[
	CHAINID types.ID,
	SEQ types.Sequence,
	ADDR types.Hashable,
	BLOCK any,
	BLOCKHASH types.Hashable,
	TX any,
	TXHASH types.Hashable,
	EVENT any,
	EVENTOPS any, // event filter query options
	TXRECEIPT txmgrtypes.ChainReceipt[TXHASH, BLOCKHASH],
	FEE feetypes.Fee,
] interface {
	RPCClient[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE]
}

type multiNodeClient[
	CHAINID types.ID,
	SEQ types.Sequence,
	ADDR types.Hashable,
	BLOCK any,
	BLOCKHASH types.Hashable,
	TX any,
	TXHASH types.Hashable,
	EVENT any,
	EVENTOPS any,
	TXRECEIPT txmgrtypes.ChainReceipt[TXHASH, BLOCKHASH],
	FEE feetypes.Fee,
] struct {
	logger logger.Logger
	pool   Pool[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE]
}

func NewMultiNodeClient[
	CHAINID types.ID,
	SEQ types.Sequence,
	ADDR types.Hashable,
	BLOCK any,
	BLOCKHASH types.Hashable,
	TX any,
	TXHASH types.Hashable,
	EVENT any,
	EVENTOPS any,
	TXRECEIPT txmgrtypes.ChainReceipt[TXHASH, BLOCKHASH],
	FEE feetypes.Fee,
](logger logger.Logger, selectionMode string, noNewHeadsTreshold time.Duration, nodes []Node[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE], sendonlys []Node[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE], chainID *big.Int, chainType config.ChainType,
) MultiNodeClient[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE] {
	pool := NewPool[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE](
		logger, selectionMode, noNewHeadsTreshold, nodes, sendonlys, chainID, chainType)
	return &multiNodeClient[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE]{
		pool: pool,
	}
}

func (client *multiNodeClient[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE]) BalanceAt(ctx context.Context, account ADDR, blockNumber *big.Int) (*big.Int, error) {
	return client.pool.BalanceAt(ctx, account, blockNumber)
}

func (client *multiNodeClient[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE]) BatchGetReceipts(
	ctx context.Context,
	attempts []txmgrtypes.TxAttempt[CHAINID, ADDR, TXHASH, BLOCKHASH, SEQ, FEE],
) (txReceipt []TXRECEIPT, txErr []error, err error) {
	return client.pool.BatchGetReceipts(ctx, attempts)
}

func (client *multiNodeClient[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE]) BlockByHash(ctx context.Context, hash BLOCKHASH) (*BLOCK, error) {
	return client.pool.BlockByHash(ctx, hash)
}

func (client *multiNodeClient[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE]) BlockByNumber(ctx context.Context, number *big.Int) (*BLOCK, error) {
	return client.pool.BlockByNumber(ctx, number)
}

func (client *multiNodeClient[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE]) CallContext(ctx context.Context, result interface{}, method string, args ...interface{}) error {
	return client.pool.CallContext(ctx, result, method, args...)
}

func (client *multiNodeClient[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE]) CallContract(
	ctx context.Context,
	attempt txmgrtypes.TxAttempt[CHAINID, ADDR, TXHASH, BLOCKHASH, SEQ, FEE],
	blockNumber *big.Int,
) (rpcErr fmt.Stringer, extractErr error) {
	return client.pool.CallContract(ctx, attempt, blockNumber)
}

func (client *multiNodeClient[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE]) ChainID() (CHAINID, error) {
	return client.pool.ChainID()
}

func (client *multiNodeClient[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE]) ConfiguredChainID() CHAINID {
	return client.pool.ConfiguredChainID()
}

func (client *multiNodeClient[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE]) FilterEvents(ctx context.Context, query EVENTOPS) ([]EVENT, error) {
	return client.pool.FilterEvents(ctx, query)
}

func (client *multiNodeClient[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE]) LatestBlockHeight(ctx context.Context) (*big.Int, error) {
	return client.pool.LatestBlockHeight(ctx)
}

func (client *multiNodeClient[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE]) PendingSequenceAt(ctx context.Context, addr ADDR) (SEQ, error) {
	return client.pool.PendingSequenceAt(ctx, addr)
}

func (client *multiNodeClient[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE]) SendEmptyTransaction(
	ctx context.Context,
	newTxAttempt func(seq SEQ, feeLimit uint32, fee FEE, fromAddress ADDR) (attempt txmgrtypes.TxAttempt[CHAINID, ADDR, TXHASH, BLOCKHASH, SEQ, FEE], err error),
	seq SEQ,
	gasLimit uint32,
	fee FEE,
	fromAddress ADDR,
) (txhash string, err error) {
	return client.pool.SendEmptyTransaction(ctx, newTxAttempt, seq, gasLimit, fee, fromAddress)
}

func (client *multiNodeClient[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE]) SendTransaction(ctx context.Context, tx *TX) error {
	return client.pool.SendTransaction(ctx, tx)
}

func (client *multiNodeClient[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE]) SendTransactionReturnCode(
	ctx context.Context,
	tx txmgrtypes.Tx[CHAINID, ADDR, TXHASH, BLOCKHASH, SEQ, FEE],
	attempt txmgrtypes.TxAttempt[CHAINID, ADDR, TXHASH, BLOCKHASH, SEQ, FEE],
	lggr logger.Logger,
) (SendTxReturnCode, error) {
	return client.pool.SendTransactionReturnCode(ctx, tx, attempt, lggr)
}

func (client *multiNodeClient[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE]) SequenceAt(ctx context.Context, account ADDR, blockNumber *big.Int) (SEQ, error) {
	return client.pool.SequenceAt(ctx, account, blockNumber)
}

func (client *multiNodeClient[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE]) SimulateTransaction(ctx context.Context, tx *TX) error {
	return client.pool.SimulateTransaction(ctx, tx)
}

func (client *multiNodeClient[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE]) 	Subscribe(ctx context.Context, channel chan<- types.Head[BLOCKHASH], args ...interface{}) (types.Subscription, error) {
	return client.pool.Subscribe(ctx, channel, args)
}


func (client *multiNodeClient[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE]) TokenBalance(ctx context.Context, account ADDR, tokenAddr ADDR) (*big.Int, error) {
	return client.pool.TokenBalance(ctx, account, tokenAddr)
}

func (client *multiNodeClient[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE]) TransactionByHash(ctx context.Context, txHash TXHASH) (*TX, error) {
	return client.pool.TransactionByHash(ctx, txHash)
}

func (client *multiNodeClient[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE]) TransactionReceipt(ctx context.Context, txHash TXHASH) (*TXRECEIPT, error) {
	return client.pool.TransactionReceipt(ctx, txHash)
}
