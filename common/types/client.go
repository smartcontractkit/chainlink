package types

import (
	"context"
	"math/big"
)

// A generic client interface for communication with the RPC node
// Every native chain must implement independently
type Client[CHAINID any, ADDR comparable, BLOCK any, TX any, HASH comparable, TXRECEIPT any, EVENT any, EVENTOPS any] interface {
	// ChainID stored for quick access
	ConfiguredChainID() CHAINID
	// ChainID RPC call
	ChainID() (CHAINID, error)

	Accounts[ADDR]
	Transactions[TX, HASH, TXRECEIPT]
	Events[EVENT, EVENTOPS]
	Blocks[BLOCK, HASH]

	CallContext(ctx context.Context, result interface{}, method string, args ...interface{}) error
}

type Accounts[ADDR comparable] interface {
	BalanceAt(ctx context.Context, accountAddress ADDR, blockNumber *big.Int) (*big.Int, error)
	TokenBalance(ctx context.Context, accountAddress ADDR, tokenAddress ADDR) (*big.Int, error)
	SequenceAt(ctx context.Context, accountAddress ADDR, blockNumber *big.Int) (uint64, error)
}

type Transactions[TX any, HASH comparable, TXRECEIPT any] interface {
	SendTransaction(ctx context.Context, tx *TX) error
	SimulateTransaction(ctx context.Context, tx *TX) error
	TransactionByHash(ctx context.Context, txHash HASH) (*TX, error)
	TransactionReceipt(ctx context.Context, txHash HASH) (*TXRECEIPT, error)
}

type Blocks[BLOCK any, HASH any] interface {
	BlockByNumber(ctx context.Context, number *big.Int) (*BLOCK, error)
	BlockByHash(ctx context.Context, hash HASH) (*BLOCK, error)
	LatestBlockHeight(context.Context) (*big.Int, error)
}

type Events[EVENT any, EVENTOPS any] interface {
	FilterEvents(ctx context.Context, query EVENTOPS) ([]EVENT, error)
}
