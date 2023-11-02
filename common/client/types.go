package client

import (
	"context"
	"math/big"

	feetypes "github.com/smartcontractkit/chainlink/v2/common/fee/types"
	"github.com/smartcontractkit/chainlink/v2/common/types"
	"github.com/smartcontractkit/chainlink/v2/core/assets"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
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
	clientAPI[
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

// Head is the interface required by the NodeClient
type Head interface {
	BlockNumber() int64
	BlockDifficulty() *utils.Big
}

// NodeClient includes all the necessary RPC methods required by a node.
type NodeClient[
	CHAIN_ID types.ID,
	HEAD Head,
] interface {
	connection[CHAIN_ID, HEAD]

	DialHTTP() error
	DisconnectAll()
	Close()
	ClientVersion(context.Context) (string, error)
	SubscribersCount() int32
	SetAliveLoopSub(types.Subscription)
	UnsubscribeAllExceptAliveLoop()
}

// clientAPI includes all the direct RPC methods required by the generalized common client to implement its own.
type clientAPI[
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
	connection[CHAIN_ID, HEAD]

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

	// Misc
	BatchCallContext(ctx context.Context, b []any) error
	CallContract(
		ctx context.Context,
		msg interface{},
		blockNumber *big.Int,
	) (rpcErr []byte, extractErr error)
	CallContext(ctx context.Context, result interface{}, method string, args ...interface{}) error
	CodeAt(ctx context.Context, account ADDR, blockNumber *big.Int) ([]byte, error)
}

type connection[
	CHAIN_ID types.ID,
	HEAD Head,
] interface {
	ChainID(ctx context.Context) (CHAIN_ID, error)
	Dial(ctx context.Context) error
	Subscribe(ctx context.Context, channel chan<- HEAD, args ...interface{}) (types.Subscription, error)
}
