package types

import (
	"context"
	"math/big"

	"github.com/smartcontractkit/chainlink/v2/common/types"
)

// A generic client interface for communication with the RPC node
// Every native chain must implement independently
type Client[
	CHAINID ID,
	SEQ Sequence, // nonce
	ADDR types.Hashable,
	BLOCK any,
	BLOCKHASH types.Hashable,
	TX any,
	TXHASH types.Hashable,
	TXRECEIPT any,
	EVENT any,
	EVENTOPS any, // event filter query options
] interface {
	// ChainID stored for quick access
	ConfiguredChainID() CHAINID
	// ChainID RPC call
	ChainID() (CHAINID, error)

	Accounts[ADDR, SEQ]
	Transactions[TX, TXHASH, TXRECEIPT]
	Events[EVENT, EVENTOPS]
	Blocks[BLOCK, BLOCKHASH]

	CallContext(ctx context.Context, result interface{}, method string, args ...interface{}) error
}

type Accounts[ADDR types.Hashable, SEQ Sequence] interface {
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
