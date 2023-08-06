package client

import (
	"context"
	"math/big"

	feetypes "github.com/smartcontractkit/chainlink/v2/common/fee/types"
	"github.com/smartcontractkit/chainlink/v2/common/types"
	"github.com/smartcontractkit/chainlink/v2/core/assets"
)

type RPCClient[
	CHAIN_ID types.ID,
	SEQ types.Sequence,
	ADDR types.Hashable,
	BLOCK_HASH types.Hashable,
	TX any,
	TX_HASH types.Hashable,
	EVENT any,
	EVENT_OPS any, // event filter query options
	TX_RECEIPT any,
	FEE feetypes.Fee,
	HEAD types.Head[BLOCK_HASH],
	SUB types.Subscription,
] interface {
	Accounts[ADDR, SEQ]
	Transactions[ADDR, TX, TX_HASH, TX_RECEIPT, SEQ, FEE]
	Events[EVENT, EVENT_OPS]
	Blocks[HEAD, BLOCK_HASH]

	BatchCallContext(ctx context.Context, b []any) error
	CallContract(
		ctx context.Context,
		attempt interface{},
		blockNumber *big.Int,
	) (rpcErr []byte, extractErr error)
	CallContext(ctx context.Context, result interface{}, method string, args ...interface{}) error
	ChainID() (CHAIN_ID, error)
	CodeAt(ctx context.Context, account ADDR, blockNumber *big.Int) ([]byte, error)
	ConfiguredChainID() CHAIN_ID

	Subscribe(ctx context.Context, channel chan<- HEAD, args ...interface{}) (SUB, error)
}

type Accounts[ADDR types.Hashable, SEQ types.Sequence] interface {
	BalanceAt(ctx context.Context, accountAddress ADDR, blockNumber *big.Int) (*big.Int, error)
	TokenBalance(ctx context.Context, accountAddress ADDR, tokenAddress ADDR) (*big.Int, error)
	SequenceAt(ctx context.Context, accountAddress ADDR, blockNumber *big.Int) (SEQ, error)
	LINKBalance(ctx context.Context, accountAddress ADDR, linkAddress ADDR) (*assets.Link, error)
	PendingSequenceAt(ctx context.Context, addr ADDR) (SEQ, error)
	EstimateGas(ctx context.Context, call any) (gas uint64, err error)
}

type Transactions[
	ADDR types.Hashable,
	TX any,
	TX_HASH types.Hashable,
	TX_RECEIPT any,
	SEQ types.Sequence,
	FEE feetypes.Fee,
] interface {
	SendTransaction(ctx context.Context, tx *TX) error
	SimulateTransaction(ctx context.Context, tx *TX) error
	TransactionByHash(ctx context.Context, txHash TX_HASH) (*TX, error)
	TransactionReceipt(ctx context.Context, txHash TX_HASH) (*TX_RECEIPT, error)
	SendEmptyTransaction(
		ctx context.Context,
		newTxAttempt func(seq SEQ, feeLimit uint32, fee FEE, fromAddress ADDR) (attempt any, err error),
		seq SEQ,
		gasLimit uint32,
		fee FEE,
		fromAddress ADDR,
	) (txhash string, err error)
	SendTransactionReturnCode(
		ctx context.Context,
		tx *TX,
	) (SendTxReturnCode, error)
}

type Blocks[HEAD types.Head[BLOCK_HASH], BLOCK_HASH types.Hashable] interface {
	BlockByNumber(ctx context.Context, number *big.Int) (HEAD, error)
	BlockByHash(ctx context.Context, hash BLOCK_HASH) (HEAD, error)
	LatestBlockHeight(context.Context) (*big.Int, error)
}

type Events[EVENT any, EVENT_OPS any] interface {
	FilterEvents(ctx context.Context, query EVENT_OPS) ([]EVENT, error)
}
