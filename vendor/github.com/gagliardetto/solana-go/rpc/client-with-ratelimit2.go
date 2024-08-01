package rpc

import (
	"context"
	"io"
	"net/http"

	"github.com/gagliardetto/solana-go/rpc/jsonrpc"
	"golang.org/x/time/rate"
)

var _ JSONRPCClient = &clientWithLimiter{}

type clientWithLimiter struct {
	rpcClient jsonrpc.RPCClient
	limiter   *rate.Limiter
}

// NewWithLimiter creates a new rate-limitted Solana RPC client.
// Example: NewWithLimiter(URL, rate.Every(time.Second), 1)
func NewWithLimiter(
	rpcEndpoint string,
	every rate.Limit, // time frame
	b int, // number of requests per time frame
) JSONRPCClient {
	opts := &jsonrpc.RPCClientOpts{
		HTTPClient: newHTTP(),
	}

	rpcClient := jsonrpc.NewClientWithOpts(rpcEndpoint, opts)
	rater := rate.NewLimiter(every, b)

	return &clientWithLimiter{
		rpcClient: rpcClient,
		limiter:   rater,
	}
}

func (wr *clientWithLimiter) CallForInto(ctx context.Context, out interface{}, method string, params []interface{}) error {
	err := wr.limiter.Wait(ctx)
	if err != nil {
		return err
	}
	return wr.rpcClient.CallForInto(ctx, &out, method, params)
}

func (wr *clientWithLimiter) CallWithCallback(
	ctx context.Context,
	method string,
	params []interface{},
	callback func(*http.Request, *http.Response) error,
) error {
	err := wr.limiter.Wait(ctx)
	if err != nil {
		return err
	}
	return wr.rpcClient.CallWithCallback(ctx, method, params, callback)
}

func (wr *clientWithLimiter) CallBatch(
	ctx context.Context,
	requests jsonrpc.RPCRequests,
) (jsonrpc.RPCResponses, error) {
	err := wr.limiter.Wait(ctx)
	if err != nil {
		return nil, err
	}
	return wr.rpcClient.CallBatch(ctx, requests)
}

// Close closes clientWithLimiter.
func (cl *clientWithLimiter) Close() error {
	if c, ok := cl.rpcClient.(io.Closer); ok {
		return c.Close()
	}
	return nil
}
