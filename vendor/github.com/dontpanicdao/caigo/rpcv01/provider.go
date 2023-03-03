package rpcv01

import (
	"context"
	"errors"

	ctypes "github.com/dontpanicdao/caigo/types"
	"github.com/ethereum/go-ethereum/rpc"
)

// ErrNotFound is returned by API methods if the requested item does not exist.
var (
	errNotFound = errors.New("not found")
)

// Provider provides the provider for caigo/rpc implementation.
type Provider struct {
	c callCloser
}

// NewProvider creates a *Provider from an existing `go-ethereum/rpc` *Client.
func NewProvider(c *rpc.Client) *Provider {
	return &Provider{c: c}
}

type api interface {
	BlockHashAndNumber(ctx context.Context) (*BlockHashAndNumberOutput, error)
	BlockNumber(ctx context.Context) (uint64, error)
	BlockTransactionCount(ctx context.Context, blockID BlockID) (uint64, error)
	BlockWithTxHashes(ctx context.Context, blockID BlockID) (Block, error)
	BlockWithTxs(ctx context.Context, blockID BlockID) (interface{}, error)
	Call(ctx context.Context, call ctypes.FunctionCall, block BlockID) ([]string, error)
	ChainID(ctx context.Context) (string, error)
	Class(ctx context.Context, classHash string) (*ctypes.ContractClass, error)
	ClassAt(ctx context.Context, blockID BlockID, contractAddress ctypes.Hash) (*ctypes.ContractClass, error)
	ClassHashAt(ctx context.Context, blockID BlockID, contractAddress ctypes.Hash) (*string, error)
	EstimateFee(ctx context.Context, request ctypes.FunctionInvoke, blockID BlockID) (*ctypes.FeeEstimate, error)
	Events(ctx context.Context, filter EventFilter) (*EventsOutput, error)
	Nonce(ctx context.Context, contractAddress ctypes.Hash) (*string, error)
	StateUpdate(ctx context.Context, blockID BlockID) (*StateUpdateOutput, error)
	StorageAt(ctx context.Context, contractAddress ctypes.Hash, key string, blockID BlockID) (string, error)
	Syncing(ctx context.Context) (*SyncResponse, error)
	TransactionByBlockIdAndIndex(ctx context.Context, blockID BlockID, index uint64) (Transaction, error)
	TransactionByHash(ctx context.Context, hash ctypes.Hash) (Transaction, error)
	TransactionReceipt(ctx context.Context, transactionHash ctypes.Hash) (TransactionReceipt, error)
}

var _ api = &Provider{}
