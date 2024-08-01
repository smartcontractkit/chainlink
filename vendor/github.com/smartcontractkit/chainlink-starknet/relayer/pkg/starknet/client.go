package starknet

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/NethermindEth/juno/core/felt"
	starknetrpc "github.com/NethermindEth/starknet.go/rpc"
	ethrpc "github.com/ethereum/go-ethereum/rpc"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
)

//go:generate mockery --name Reader --output ./mocks/

type Reader interface {
	CallContract(context.Context, CallOps) ([]*felt.Felt, error)
	LatestBlockHeight(context.Context) (uint64, error)

	// provider interface
	BlockWithTxHashes(ctx context.Context, blockID starknetrpc.BlockID) (*starknetrpc.Block, error)
	Call(context.Context, starknetrpc.FunctionCall, starknetrpc.BlockID) ([]*felt.Felt, error)
	Events(ctx context.Context, input starknetrpc.EventsInput) (*starknetrpc.EventChunk, error)
	TransactionByHash(context.Context, *felt.Felt) (starknetrpc.Transaction, error)
	TransactionReceipt(context.Context, *felt.Felt) (starknetrpc.TransactionReceipt, error)
	AccountNonce(context.Context, *felt.Felt) (*felt.Felt, error)
}

type Writer interface {
}

type ReaderWriter interface {
	Reader
	Writer
}

var _ ReaderWriter = (*Client)(nil)

// var _ starknettypes.Provider = (*Client)(nil)

type Client struct {
	Provider       starknetrpc.RpcProvider
	EthClient      *ethrpc.Client
	lggr           logger.Logger
	defaultTimeout time.Duration
}

// pass nil or 0 to timeout to not use built in default timeout
func NewClient(chainID string, baseURL string, apiKey string, lggr logger.Logger, timeout *time.Duration) (*Client, error) {
	// TODO: chainID now unused

	options := []ethrpc.ClientOption{}
	if strings.TrimSpace(apiKey) != "" {
		options = append(options, ethrpc.WithHeader("x-apikey", apiKey))
	}

	provider, err := starknetrpc.NewProvider(baseURL, options...)
	if err != nil {
		return nil, err
	}

	c, err := ethrpc.DialContext(context.Background(), baseURL)
	if err != nil {
		return nil, err
	}

	client := &Client{
		Provider:  provider,
		EthClient: c,
		lggr:      lggr,
	}

	// make copy to preserve value
	// defensive in case the timeout reference is ever garbage collected or removed
	if timeout == nil {
		client.defaultTimeout = 0
	} else {
		client.defaultTimeout = *timeout
	}

	return client, nil
}

// -- Custom Wrapped Func --

func (c *Client) CallContract(ctx context.Context, ops CallOps) (data []*felt.Felt, err error) {
	tx := starknetrpc.FunctionCall{
		ContractAddress:    ops.ContractAddress,
		EntryPointSelector: ops.Selector,
		Calldata:           ops.Calldata,
	}

	res, err := c.Call(ctx, tx, starknetrpc.WithBlockTag("latest"))
	if err != nil {
		return nil, fmt.Errorf("error in client.CallContract: %w", err)
	}

	return res, nil
}

func (c *Client) LatestBlockHeight(ctx context.Context) (uint64, error) {
	if c.defaultTimeout != 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, c.defaultTimeout)
		defer cancel()
	}

	blockNum, err := c.Provider.BlockNumber(ctx)
	if err != nil {
		return 0, fmt.Errorf("error in client.LatestBlockHeight: %w", err)
	}

	return blockNum, nil
}

// -- caigo.Provider interface --

func (c *Client) BlockWithTxHashes(ctx context.Context, blockID starknetrpc.BlockID) (*starknetrpc.Block, error) {
	if c.defaultTimeout != 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, c.defaultTimeout)
		defer cancel()
	}

	out, err := c.Provider.BlockWithTxHashes(ctx, blockID)
	if err != nil {
		return out.(*starknetrpc.Block), fmt.Errorf("error in client.BlockWithTxHashes: %w", err)
	}
	return out.(*starknetrpc.Block), nil
}

func (c *Client) Call(ctx context.Context, calls starknetrpc.FunctionCall, blockHashOrTag starknetrpc.BlockID) ([]*felt.Felt, error) {
	if c.defaultTimeout != 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, c.defaultTimeout)
		defer cancel()
	}

	out, err := c.Provider.Call(ctx, calls, blockHashOrTag)
	if err != nil {
		return out, fmt.Errorf("error in client.Call: %w", err)
	}
	if out == nil {
		return out, NilResultError("client.Call")
	}
	return out, nil
}

func (c *Client) TransactionByHash(ctx context.Context, hash *felt.Felt) (starknetrpc.Transaction, error) {
	if c.defaultTimeout != 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, c.defaultTimeout)
		defer cancel()
	}

	out, err := c.Provider.TransactionByHash(ctx, hash)
	if err != nil {
		return out, fmt.Errorf("error in client.TransactionByHash: %w", err)
	}
	if out == nil {
		return out, NilResultError("client.TransactionByHash")
	}
	return out, nil
}

func (c *Client) TransactionReceipt(ctx context.Context, hash *felt.Felt) (starknetrpc.TransactionReceipt, error) {
	if c.defaultTimeout != 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, c.defaultTimeout)
		defer cancel()
	}

	out, err := c.Provider.TransactionReceipt(ctx, hash)
	if err != nil {
		return out, fmt.Errorf("error in client.TransactionReceipt: %w", err)
	}
	if out == nil {
		return out, NilResultError("client.TransactionReceipt")
	}
	return out, nil
}

func (c *Client) Events(ctx context.Context, input starknetrpc.EventsInput) (*starknetrpc.EventChunk, error) {
	if c.defaultTimeout != 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, c.defaultTimeout)
		defer cancel()
	}

	out, err := c.Provider.Events(ctx, input)
	if err != nil {
		return out, fmt.Errorf("error in client.Events: %w", err)
	}
	if out == nil {
		return out, NilResultError("client.Events")
	}
	return out, nil
}

func (c *Client) AccountNonce(ctx context.Context, accountAddress *felt.Felt) (*felt.Felt, error) {
	if c.defaultTimeout != 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, c.defaultTimeout)
		defer cancel()
	}

	return c.Provider.Nonce(ctx, starknetrpc.BlockID{Tag: "pending"}, accountAddress)
}
