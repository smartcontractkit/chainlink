package client

import (
	"context"
	"math/big"

	feetypes "github.com/smartcontractkit/chainlink/v2/common/fee/types"
	"github.com/smartcontractkit/chainlink/v2/common/types"
	"github.com/smartcontractkit/chainlink/v2/core/assets"
)

// RPC includes all the necessary methods for a multi-node client to interact directly with any RPC endpoint.
type RPC[
	CHAIN_ID types.ID,
	SEQ types.Sequence,
	ADDR types.Hashable,
	BLOCK_HASH types.Hashable,
	TX any,
	TX_HASH types.Hashable,
	EVENT any,
	EVENT_OPS any,
	TX_RECEIPT types.Receipt[TX_HASH, BLOCK_HASH],
	FEE feetypes.Fee,
	HEAD types.Head[BLOCK_HASH],

] interface {
	NodeClient[
		CHAIN_ID,
		HEAD,
	]
	ClientAPI[
		CHAIN_ID,
		SEQ,
		ADDR,
		BLOCK_HASH,
		TX,
		TX_HASH,
		EVENT,
		EVENT_OPS,
		TX_RECEIPT,
		FEE,
		HEAD,
	]
}

type ClientAPI[
	CHAIN_ID types.ID,
	SEQ types.Sequence,
	ADDR types.Hashable,
	BLOCK_HASH types.Hashable,
	TX any,
	TX_HASH types.Hashable,
	EVENT any,
	EVENT_OPS any, // event filter query options
	TX_RECEIPT types.Receipt[TX_HASH, BLOCK_HASH],
	FEE feetypes.Fee,
	HEAD types.Head[BLOCK_HASH],
] interface {
	// Underlying RPC methods required by generalized client.
	// Account
	BalanceAt(ctx context.Context, accountAddress ADDR, blockNumber *big.Int) (*big.Int, error)
	TokenBalance(ctx context.Context, accountAddress ADDR, tokenAddress ADDR) (*big.Int, error)
	SequenceAt(ctx context.Context, accountAddress ADDR, blockNumber *big.Int) (SEQ, error)
	LINKBalance(ctx context.Context, accountAddress ADDR, linkAddress ADDR) (*assets.Link, error)
	PendingSequenceAt(ctx context.Context, addr ADDR) (SEQ, error)
	EstimateGas(ctx context.Context, call any) (gas uint64, err error)

	// Transactions
	SendTransaction(ctx context.Context, tx TX) error
	SimulateTransaction(ctx context.Context, tx TX) error
	TransactionByHash(ctx context.Context, txHash TX_HASH) (TX, error)
	TransactionReceipt(ctx context.Context, txHash TX_HASH) (TX_RECEIPT, error)
	SendEmptyTransaction(
		ctx context.Context,
		newTxAttempt func(seq SEQ, feeLimit uint32, fee FEE, fromAddress ADDR) (attempt any, err error),
		seq SEQ,
		gasLimit uint32,
		fee FEE,
		fromAddress ADDR,
	) (txhash string, err error)

	// Blocks
	BlockByNumber(ctx context.Context, number *big.Int) (HEAD, error)
	BlockByHash(ctx context.Context, hash BLOCK_HASH) (HEAD, error)
	LatestBlockHeight(context.Context) (*big.Int, error)

	// Events
	FilterEvents(ctx context.Context, query EVENT_OPS) ([]EVENT, error)

	// Connection
	Close()
	Dial(ctx context.Context) error

	// Misc
	BatchCallContext(ctx context.Context, b []any) error
	CallContract(
		ctx context.Context,
		msg interface{},
		blockNumber *big.Int,
	) (rpcErr []byte, extractErr error)
	CallContext(ctx context.Context, result interface{}, method string, args ...interface{}) error
	ChainID(ctx context.Context) (CHAIN_ID, error)
	CodeAt(ctx context.Context, account ADDR, blockNumber *big.Int) ([]byte, error)
	Subscribe(ctx context.Context, channel chan<- HEAD, args ...interface{}) (types.Subscription, error)
}
