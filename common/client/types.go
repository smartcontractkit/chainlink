package client

import (
	"context"
	"math/big"

	"github.com/smartcontractkit/chainlink-common/pkg/assets"

	feetypes "github.com/smartcontractkit/chainlink/v2/common/fee/types"
	"github.com/smartcontractkit/chainlink/v2/common/types"
)

// RPCClient includes all the necessary generalized RPC methods along with any additional chain-specific methods.
//
//go:generate mockery --quiet --name RPCClient --output ./mocks --case=underscore
type RPCClient[
	CHAIN_ID types.ID,
	HEAD Head,
] interface {
	// ChainID - fetches ChainID from the RPC to verify that it matches config
	ChainID(ctx context.Context) (CHAIN_ID, error)
	// Dial - prepares the RPC for usage. Can be called on fresh or closed RPC
	Dial(ctx context.Context) error
	// SubscribeToHeads - returns channel and subscription for new heads.
	SubscribeToHeads(ctx context.Context) (<-chan HEAD, types.Subscription, error)
	// SubscribeToFinalizedHeads - returns channel and subscription for finalized heads.
	SubscribeToFinalizedHeads(ctx context.Context) (<-chan HEAD, types.Subscription, error)
	// Ping - returns error if RPC is not reachable
	Ping(context.Context) error
	// IsSyncing - returns true if the RPC is in Syncing state and can not process calls
	IsSyncing(ctx context.Context) (bool, error)
	// UnsubscribeAllExcept - close all subscriptions except `subs`
	UnsubscribeAllExcept(subs ...types.Subscription)
	// Close - closes all subscriptions and aborts all RPC calls
	Close()
}

// RPC includes all the necessary methods for a multi-node client to interact directly with any RPC endpoint.
//
//go:generate mockery --quiet --name RPC --structname mockRPC --inpackage --filename "mock_rpc_test.go" --case=underscore
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
//
//go:generate mockery --quiet --name Head --structname mockHead --filename "mock_head_test.go" --inpackage --case=underscore
type Head interface {
	BlockNumber() int64
	BlockDifficulty() *big.Int
	IsValid() bool
}

// NodeClient includes all the necessary RPC methods required by a node.
//
//go:generate mockery --quiet --name NodeClient --structname mockNodeClient --filename "mock_node_client_test.go" --inpackage --case=underscore
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
	IsSyncing(ctx context.Context) (bool, error)
	LatestFinalizedBlock(ctx context.Context) (HEAD, error)
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
	Subscribe(ctx context.Context, channel chan<- HEAD, args ...interface{}) (types.Subscription, error)
}
