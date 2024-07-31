package client

import (
	"context"
	"math/big"

	"github.com/smartcontractkit/chainlink-common/pkg/assets"

	feetypes "github.com/smartcontractkit/chainlink/v2/common/fee/types"
	"github.com/smartcontractkit/chainlink/v2/common/types"
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
	BATCH_ELEM any,
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
		BATCH_ELEM,
	]
}

// Head is the interface required by the NodeClient
type Head interface {
	BlockNumber() int64
	BlockDifficulty() *big.Int
	IsValid() bool
}

// NodeClient includes all the necessary RPC methods required by a node.
type NodeClient[
	CHAIN_ID types.ID,
	HEAD Head,
] interface {
	connection[CHAIN_ID, HEAD]

	DialHTTP() error
	// DisconnectAll - cancels all inflight requests, terminates all subscriptions and resets latest ChainInfo.
	DisconnectAll()
	Close()
	ClientVersion(context.Context) (string, error)
	SubscribersCount() int32
	SetAliveLoopSub(types.Subscription)
	UnsubscribeAllExceptAliveLoop()
	IsSyncing(ctx context.Context) (bool, error)
	SubscribeToFinalizedHeads(_ context.Context) (<-chan HEAD, types.Subscription, error)
	// GetInterceptedChainInfo - returns latest and highest observed by application layer ChainInfo.
	// latest ChainInfo is the most recent value received within a NodeClient's current lifecycle between Dial and DisconnectAll.
	// highestUserObservations ChainInfo is the highest ChainInfo observed excluding health checks calls.
	// Its values must not be reset.
	// The results of corresponding calls, to get the most recent head and the latest finalized head, must be
	// intercepted and reflected in ChainInfo before being returned to a caller. Otherwise, MultiNode is not able to
	// provide repeatable read guarantee.
	// DisconnectAll must reset latest ChainInfo to default value.
	// Ensure implementation does not have a race condition when values are reset before request completion and as
	// a result latest ChainInfo contains information from the previous cycle.
	GetInterceptedChainInfo() (latest, highestUserObservations ChainInfo)
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
	BATCH_ELEM any,
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
	LatestFinalizedBlock(ctx context.Context) (HEAD, error)

	// Events
	FilterEvents(ctx context.Context, query EVENT_OPS) ([]EVENT, error)

	// Misc
	BatchCallContext(ctx context.Context, b []BATCH_ELEM) error
	CallContract(
		ctx context.Context,
		msg interface{},
		blockNumber *big.Int,
	) (rpcErr []byte, extractErr error)
	PendingCallContract(
		ctx context.Context,
		msg interface{},
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
	SubscribeToHeads(ctx context.Context) (ch <-chan HEAD, sub types.Subscription, err error)
	// TODO: remove as part of merge with BCI-2875
	SubscribeNewHead(ctx context.Context, channel chan<- HEAD) (s types.Subscription, err error)
}

// PoolChainInfoProvider - provides aggregation of nodes pool ChainInfo
type PoolChainInfoProvider interface {
	// LatestChainInfo - returns number of live nodes available in the pool, so we can prevent the last alive node in a pool from being
	// moved to out-of-sync state. It is better to have one out-of-sync node than no nodes at all.
	// Returns highest latest ChainInfo within the alive nodes. E.g. most recent block number and highest block number
	// observed by Node A are 10 and 15; Node B - 12 and 14. This method will return 12.
	LatestChainInfo() (int, ChainInfo)
	// HighestUserObservations - returns highest ChainInfo ever observed by any user of MultiNode.
	HighestUserObservations() ChainInfo
}

// ChainInfo - defines RPC's or MultiNode's view on the chain
type ChainInfo struct {
	BlockNumber          int64
	FinalizedBlockNumber int64
	TotalDifficulty      *big.Int
}

func MaxTotalDifficulty(a, b *big.Int) *big.Int {
	if a == nil {
		if b == nil {
			return nil
		}

		return big.NewInt(0).Set(b)
	}

	if b == nil || a.Cmp(b) >= 0 {
		return big.NewInt(0).Set(a)
	}

	return big.NewInt(0).Set(b)
}
