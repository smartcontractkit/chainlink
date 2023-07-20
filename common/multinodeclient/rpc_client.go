package multinodeclient

import (
	"context"
	"fmt"
	"math/big"

	feetypes "github.com/smartcontractkit/chainlink/v2/common/fee/types"
	txmgrtypes "github.com/smartcontractkit/chainlink/v2/common/txmgr/types"
	"github.com/smartcontractkit/chainlink/v2/common/types"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

type SendTxReturnCode int

type RPCClient[
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
	Accounts[ADDR, SEQ]
	Transactions[TX, TXHASH, TXRECEIPT]
	Events[EVENT, EVENTOPS]
	Blocks[BLOCK, BLOCKHASH]

	BatchGetReceipts(
		ctx context.Context,
		attempts []txmgrtypes.TxAttempt[CHAINID, ADDR, TXHASH, BLOCKHASH, SEQ, FEE],
	) (txReceipt []TXRECEIPT, txErr []error, err error)
	CallContract(
		ctx context.Context,
		attempt txmgrtypes.TxAttempt[CHAINID, ADDR, TXHASH, BLOCKHASH, SEQ, FEE],
		blockNumber *big.Int,
	) (rpcErr fmt.Stringer, extractErr error)
	CallContext(ctx context.Context, result interface{}, method string, args ...interface{}) error
	ChainID() (CHAINID, error)
	ConfiguredChainID() CHAINID
	PendingSequenceAt(ctx context.Context, addr ADDR) (SEQ, error)
	SendEmptyTransaction(
		ctx context.Context,
		newTxAttempt func(seq SEQ, feeLimit uint32, fee FEE, fromAddress ADDR) (attempt txmgrtypes.TxAttempt[CHAINID, ADDR, TXHASH, BLOCKHASH, SEQ, FEE], err error),
		seq SEQ,
		gasLimit uint32,
		fee FEE,
		fromAddress ADDR,
	) (txhash string, err error)
	SendTransactionReturnCode(
		ctx context.Context,
		tx txmgrtypes.Tx[CHAINID, ADDR, TXHASH, BLOCKHASH, SEQ, FEE],
		attempt txmgrtypes.TxAttempt[CHAINID, ADDR, TXHASH, BLOCKHASH, SEQ, FEE],
		lggr logger.Logger,
	) (SendTxReturnCode, error)
	Subscribe(ctx context.Context, channel chan<- types.Head[BLOCKHASH], args ...interface{}) (types.Subscription, error)
}

type Accounts[ADDR types.Hashable, SEQ types.Sequence] interface {
	BalanceAt(ctx context.Context, accountAddress ADDR, blockNumber *big.Int) (*big.Int, error)
	TokenBalance(ctx context.Context, accountAddress ADDR, tokenAddress ADDR) (*big.Int, error)
	SequenceAt(ctx context.Context, accountAddress ADDR, blockNumber *big.Int) (SEQ, error)
}

type Transactions[TX any, TXHASH types.Hashable, TXRECEIPT any] interface {
	SendTransaction(ctx context.Context, tx *TX) error
	SimulateTransaction(ctx context.Context, tx *TX) error
	TransactionByHash(ctx context.Context, txHash TXHASH) (*TX, error)
	TransactionReceipt(ctx context.Context, txHash TXHASH) (*TXRECEIPT, error)
}

type Blocks[BLOCK any, BLOCKHASH types.Hashable] interface {
	BlockByNumber(ctx context.Context, number *big.Int) (*BLOCK, error)
	BlockByHash(ctx context.Context, hash BLOCKHASH) (*BLOCK, error)
	LatestBlockHeight(context.Context) (*big.Int, error)
}

type Events[EVENT any, EVENTOPS any] interface {
	FilterEvents(ctx context.Context, query EVENTOPS) ([]EVENT, error)
}
