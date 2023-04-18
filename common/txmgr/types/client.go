package types

import (
	"context"
	"math/big"
)

type Client[CHAINID any, ADDR comparable, BLOCK any, HEADER any, TX any, HASH comparable, TXRECEIPT any, EVENT any, EVENTOPS any] interface {
	ChainID() (CHAINID, error)

	// account
	BalanceAt(ctx context.Context, accountAddress ADDR, blockNumber *big.Int) (*big.Int, error)
	TokenBalance(ctx context.Context, accountAddress ADDR, tokenAddress ADDR) (*big.Int, error)
	NonceAt(ctx context.Context, accountAddress ADDR, blockNumber *big.Int) (uint64, error)

	// tx
	SendTransaction(ctx context.Context, tx *TX) error
	SimulateTransaction(ctx context.Context, tx *TX) error
	TransactionByHash(ctx context.Context, txHash HASH) (*TX, error)
	TransactionReceipt(ctx context.Context, txHash HASH) (*TXRECEIPT, error)

	// events
	FilterEvents(ctx context.Context, query EVENTOPS) ([]EVENT, error)

	// block
	BlockByNumber(ctx context.Context, number *big.Int) (*BLOCK, error)
	BlockByHash(ctx context.Context, hash HASH) (*BLOCK, error)
	LatestBlockHeight(context.Context) (*big.Int, error)
	HeaderByNumber(context.Context, *big.Int) (*HEADER, error)
	HeaderByHash(context.Context, HASH) (*HEADER, error)

	CallContext(ctx context.Context, result interface{}, method string, args ...interface{}) error
}