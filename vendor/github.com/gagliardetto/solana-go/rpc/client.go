// Copyright 2021 github.com/gagliardetto
// This file has been modified by github.com/gagliardetto
//
// Copyright 2020 dfuse Platform Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package rpc

import (
	"context"
	"errors"
	"io"
	"net"
	"net/http"
	"time"

	"github.com/gagliardetto/solana-go/rpc/jsonrpc"
	"github.com/klauspost/compress/gzhttp"
)

var ErrNotFound = errors.New("not found")
var ErrNotConfirmed = errors.New("not confirmed")

type Client struct {
	rpcURL    string
	rpcClient JSONRPCClient
}

type JSONRPCClient interface {
	CallForInto(ctx context.Context, out interface{}, method string, params []interface{}) error
	CallWithCallback(ctx context.Context, method string, params []interface{}, callback func(*http.Request, *http.Response) error) error
	CallBatch(ctx context.Context, requests jsonrpc.RPCRequests) (jsonrpc.RPCResponses, error)
}

// New creates a new Solana JSON RPC client.
// Client is safe for concurrent use by multiple goroutines.
func New(rpcEndpoint string) *Client {
	opts := &jsonrpc.RPCClientOpts{
		HTTPClient: newHTTP(),
	}

	rpcClient := jsonrpc.NewClientWithOpts(rpcEndpoint, opts)
	return NewWithCustomRPCClient(rpcClient)
}

// New creates a new Solana JSON RPC client with the provided custom headers.
// The provided headers will be added to each RPC request sent via this RPC client.
func NewWithHeaders(rpcEndpoint string, headers map[string]string) *Client {
	opts := &jsonrpc.RPCClientOpts{
		HTTPClient:    newHTTP(),
		CustomHeaders: headers,
	}
	rpcClient := jsonrpc.NewClientWithOpts(rpcEndpoint, opts)
	return NewWithCustomRPCClient(rpcClient)
}

// Close closes the client.
func (cl *Client) Close() error {
	if cl.rpcClient == nil {
		return nil
	}
	if c, ok := cl.rpcClient.(io.Closer); ok {
		return c.Close()
	}
	return nil
}

// NewWithCustomRPCClient creates a new Solana RPC client
// with the provided RPC client.
func NewWithCustomRPCClient(rpcClient JSONRPCClient) *Client {
	return &Client{
		rpcClient: rpcClient,
	}
}

var (
	defaultMaxIdleConnsPerHost = 9
	defaultTimeout             = 5 * time.Minute
	defaultKeepAlive           = 180 * time.Second
)

func newHTTPTransport() *http.Transport {
	return &http.Transport{
		IdleConnTimeout:     defaultTimeout,
		MaxConnsPerHost:     defaultMaxIdleConnsPerHost,
		MaxIdleConnsPerHost: defaultMaxIdleConnsPerHost,
		Proxy:               http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   defaultTimeout,
			KeepAlive: defaultKeepAlive,
			DualStack: true,
		}).DialContext,
		ForceAttemptHTTP2: true,
		// MaxIdleConns:          100,
		TLSHandshakeTimeout: 10 * time.Second,
		// ExpectContinueTimeout: 1 * time.Second,
	}
}

// newHTTP returns a new Client from the provided config.
// Client is safe for concurrent use by multiple goroutines.
func newHTTP() *http.Client {
	tr := newHTTPTransport()

	return &http.Client{
		Timeout:   defaultTimeout,
		Transport: gzhttp.Transport(tr),
	}
}

// RPCCallForInto allows to access the raw RPC client and send custom requests.
func (cl *Client) RPCCallForInto(ctx context.Context, out interface{}, method string, params []interface{}) error {
	return cl.rpcClient.CallForInto(ctx, out, method, params)
}

func (cl *Client) RPCCallWithCallback(
	ctx context.Context,
	method string,
	params []interface{},
	callback func(*http.Request, *http.Response) error,
) error {
	return cl.rpcClient.CallWithCallback(ctx, method, params, callback)
}

func (cl *Client) RPCCallBatch(
	ctx context.Context,
	requests jsonrpc.RPCRequests,
) (jsonrpc.RPCResponses, error) {
	return cl.rpcClient.CallBatch(ctx, requests)
}
