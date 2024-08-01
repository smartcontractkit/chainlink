package rpc

import (
	"context"
	"io"
	"net/http"

	"github.com/gagliardetto/solana-go/rpc/jsonrpc"
	"go.uber.org/ratelimit"
)

var _ JSONRPCClient = &clientWithRateLimiting{}

type clientWithRateLimiting struct {
	rpcClient   jsonrpc.RPCClient
	rateLimiter ratelimit.Limiter
}

// NewWithRateLimit creates a new rate-limitted Solana RPC client.
func NewWithRateLimit(
	rpcEndpoint string,
	rps int, // requests per second
) JSONRPCClient {
	opts := &jsonrpc.RPCClientOpts{
		HTTPClient: newHTTP(),
	}

	rpcClient := jsonrpc.NewClientWithOpts(rpcEndpoint, opts)

	return &clientWithRateLimiting{
		rpcClient:   rpcClient,
		rateLimiter: ratelimit.New(rps),
	}
}

func (wr *clientWithRateLimiting) CallForInto(ctx context.Context, out interface{}, method string, params []interface{}) error {
	wr.rateLimiter.Take()
	return wr.rpcClient.CallForInto(ctx, &out, method, params)
}

func (wr *clientWithRateLimiting) CallWithCallback(
	ctx context.Context,
	method string,
	params []interface{},
	callback func(*http.Request, *http.Response) error,
) error {
	wr.rateLimiter.Take()
	return wr.rpcClient.CallWithCallback(ctx, method, params, callback)
}

func (wr *clientWithRateLimiting) CallBatch(
	ctx context.Context,
	requests jsonrpc.RPCRequests,
) (jsonrpc.RPCResponses, error) {
	wr.rateLimiter.Take()
	return wr.rpcClient.CallBatch(ctx, requests)
}

// Close closes clientWithRateLimiting.
func (cl *clientWithRateLimiting) Close() error {
	if c, ok := cl.rpcClient.(io.Closer); ok {
		return c.Close()
	}
	return nil
}
