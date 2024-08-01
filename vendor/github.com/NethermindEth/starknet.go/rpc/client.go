package rpc

import (
	"context"
	"encoding/json"

	ethrpc "github.com/ethereum/go-ethereum/rpc"
)

type callCloser interface {
	CallContext(ctx context.Context, result interface{}, method string, args ...interface{}) error
	Close()
}

// do is a function that performs a remote procedure call (RPC) using the provided callCloser.
//
// Parameters:
// - ctx: represents the current execution context
// - call: the callCloser object
// - method: the string representing the RPC method to be called
// - data: the interface{} to store the result of the RPC call
// - args: variadic and can be used to pass additional arguments to the RPC method
// Returns:
// - error: an error if any occurred during the function call
func do(ctx context.Context, call callCloser, method string, data interface{}, args ...interface{}) error {
	var raw json.RawMessage
	err := call.CallContext(ctx, &raw, method, args...)
	if err != nil {
		return err
	}
	if len(raw) == 0 {
		return errNotFound
	}
	if err := json.Unmarshal(raw, &data); err != nil {
		return err
	}
	return nil
}

// NewClient creates a new ethrpc.Client instance.
//
// Parameters:
// - url: the URL of the RPC endpoint
// Returns:
// - *ethrpc.Client: a new ethrpc.Client
// - error: an error if any occurred
func NewClient(url string) (*ethrpc.Client, error) {
	return ethrpc.DialContext(context.Background(), url)
}
