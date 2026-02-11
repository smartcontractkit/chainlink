// Package contracts handles deployment, management, and interactions of smart contracts on various chains
package contracts

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type WETHToken interface {
	Address() string
	Approve(to string, amount *big.Int) error
	Transfer(to string, amount *big.Int) error
	BalanceOf(ctx context.Context, addr string) (*big.Int, error)
	Name(context.Context) (string, error)
	Decimals() uint
}

type MockLINKETHFeed interface {
	Address() string
	LatestRoundData() (*big.Int, error)
	LatestRoundDataUpdatedAt() (*big.Int, error)
}

type MockETHUSDFeed interface {
	Address() string
	LatestRoundData() (*big.Int, error)
	LatestRoundDataUpdatedAt() (*big.Int, error)
	Decimals() uint
}

type MockGasFeed interface {
	Address() string
}

type LinkToken interface {
	Address() string
	Approve(to string, amount *big.Int) error
	Transfer(to string, amount *big.Int) error
	BalanceOf(ctx context.Context, addr string) (*big.Int, error)
	TransferAndCall(to string, amount *big.Int, data []byte) (*types.Transaction, error)
	TransferAndCallFromKey(to string, amount *big.Int, data []byte, keyNum int) (*types.Transaction, error)
	Name(context.Context) (string, error)
	Decimals() uint
}

type LogEmitter interface {
	Address() common.Address
	EmitLogInts(ints []int) (*types.Transaction, error)
	EmitLogIntsIndexed(ints []int) (*types.Transaction, error)
	EmitLogIntMultiIndexed(ints int, ints2 int, count int) (*types.Transaction, error)
	EmitLogStrings(strings []string) (*types.Transaction, error)
	EmitLogIntsFromKey(ints []int, keyNum int) (*types.Transaction, error)
	EmitLogIntsIndexedFromKey(ints []int, keyNum int) (*types.Transaction, error)
	EmitLogIntMultiIndexedFromKey(ints int, ints2 int, count int, keyNum int) (*types.Transaction, error)
	EmitLogStringsFromKey(strings []string, keyNum int) (*types.Transaction, error)
	EmitLogInt(payload int) (*types.Transaction, error)
	EmitLogIntIndexed(payload int) (*types.Transaction, error)
	EmitLogString(strings string) (*types.Transaction, error)
}
