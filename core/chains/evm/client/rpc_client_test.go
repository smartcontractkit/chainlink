package client_test

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"net/url"
	"testing"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"

	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	commonclient "github.com/smartcontractkit/chainlink/v2/common/client"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/testutils"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
)

func TestRPCClient_SubscribeNewHead(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithTimeout(tests.Context(t), tests.WaitTimeout(t))
	defer cancel()

	chainId := big.NewInt(123456)
	lggr := logger.Test(t)

	type rpcServer struct {
		Head *evmtypes.Head
		URL  *url.URL
	}
	createRPCServer := func() *rpcServer {
		server := &rpcServer{}
		server.URL = testutils.NewWSServer(t, chainId, func(method string, params gjson.Result) (resp testutils.JSONRPCResponse) {
			if method == "eth_unsubscribe" {
				resp.Result = "true"
				return
			}
			assert.Equal(t, "eth_subscribe", method)
			if assert.True(t, params.IsArray()) && assert.Equal(t, "newHeads", params.Array()[0].String()) {
				resp.Result = `"0x00"`
				head := server.Head
				jsonHead, err := json.Marshal(head)
				if err != nil {
					panic(fmt.Errorf("failed to marshal head: %w", err))
				}
				resp.Notify = string(jsonHead)
			}
			return
		}).WSURL()

		return server
	}
	receiveNewHead := func(rpc client.RPCClient) *evmtypes.Head {
		ch := make(chan *evmtypes.Head)
		sub, err := rpc.SubscribeNewHead(tests.Context(t), ch)
		require.NoError(t, err)
		result := <-ch
		sub.Unsubscribe()
		return result
	}
	t.Run("Updates latest block info in InterceptedChainInfo", func(t *testing.T) {
		server := createRPCServer()
		rpc := client.NewRPCClient(lggr, *server.URL, nil, "rpc", 1, chainId, commonclient.Primary)
		defer rpc.Close()
		require.NoError(t, rpc.Dial(ctx))
		chainInfo := rpc.GetInterceptedChainInfo()
		require.Equal(t, int64(0), chainInfo.BlockNumber)
		require.Equal(t, int64(0), chainInfo.FinalizedBlockNumber)
		server.Head = &evmtypes.Head{
			Number:          256,
			TotalDifficulty: big.NewInt(1000),
		}
		_ = receiveNewHead(rpc)
		server.Head = &evmtypes.Head{
			Number:          128,
			TotalDifficulty: big.NewInt(1000),
		}
		_ = receiveNewHead(rpc)
		chainInfo = rpc.GetInterceptedChainInfo()
		assert.Equal(t, int64(256), chainInfo.BlockNumber)
		assert.Equal(t, int64(0), chainInfo.FinalizedBlockNumber)
	})
	t.Run("Block's chain ID matched configured", func(t *testing.T) {
		server := createRPCServer()
		rpc := client.NewRPCClient(lggr, *server.URL, nil, "rpc", 1, chainId, commonclient.Primary)
		defer rpc.Close()
		require.NoError(t, rpc.Dial(ctx))
		server.Head = &evmtypes.Head{
			Number: 256,
		}
		head := receiveNewHead(rpc)
		require.Equal(t, chainId, head.ChainID())
	})
	t.Run("Failed SubscribeNewHead returns proper error", func(t *testing.T) {
		server := testutils.NewWSServer(t, chainId, func(reqMethod string, reqParams gjson.Result) (resp testutils.JSONRPCResponse) {
			return resp
		})
		wsURL := server.WSURL()
		rpc := client.NewRPCClient(lggr, *wsURL, nil, "rpc", 1, chainId, commonclient.Primary)
		require.NoError(t, rpc.Dial(ctx))
		server.Close()
		_, err := rpc.SubscribeNewHead(ctx, make(chan *evmtypes.Head))
		require.ErrorContains(t, err, "RPCClient returned error (rpc)")
	})
	t.Run("Subscription error is properly wrapper", func(t *testing.T) {
		server := testutils.NewWSServer(t, chainId, func(method string, params gjson.Result) (resp testutils.JSONRPCResponse) {
			assert.Equal(t, "eth_subscribe", method)
			if assert.True(t, params.IsArray()) && assert.Equal(t, "newHeads", params.Array()[0].String()) {
				resp.Result = `"0x00"`
				resp.Notify = "{}"
			}
			return resp
		})
		wsURL := server.WSURL()
		rpc := client.NewRPCClient(lggr, *wsURL, nil, "rpc", 1, chainId, commonclient.Primary)
		defer rpc.Close()
		require.NoError(t, rpc.Dial(ctx))
		sub, err := rpc.SubscribeNewHead(ctx, make(chan *evmtypes.Head))
		require.NoError(t, err)
		server.MustWriteBinaryMessageSync(t, "invalid msg")
		select {
		case err = <-sub.Err():
			require.ErrorContains(t, err, "RPCClient returned error (rpc): invalid character")
		case <-ctx.Done():
			t.Errorf("Expected subscription to return an error, but test timeout instead")
		}
	})
}

func TestRPCClient_SubscribeFilterLogs(t *testing.T) {
	t.Parallel()

	chainId := big.NewInt(123456)
	lggr := logger.Test(t)
	ctx, cancel := context.WithTimeout(tests.Context(t), tests.WaitTimeout(t))
	defer cancel()
	t.Run("Failed SubscribeFilterLogs returns proper error", func(t *testing.T) {
		server := testutils.NewWSServer(t, chainId, func(reqMethod string, reqParams gjson.Result) (resp testutils.JSONRPCResponse) {
			return resp
		})
		wsURL := server.WSURL()
		rpc := client.NewRPCClient(lggr, *wsURL, nil, "rpc", 1, chainId, commonclient.Primary)
		require.NoError(t, rpc.Dial(ctx))
		server.Close()
		_, err := rpc.SubscribeFilterLogs(ctx, ethereum.FilterQuery{}, make(chan types.Log))
		require.ErrorContains(t, err, "RPCClient returned error (rpc)")
	})
	t.Run("Subscription error is properly wrapper", func(t *testing.T) {
		server := testutils.NewWSServer(t, chainId, func(method string, params gjson.Result) (resp testutils.JSONRPCResponse) {
			assert.Equal(t, "eth_subscribe", method)
			if assert.True(t, params.IsArray()) && assert.Equal(t, "logs", params.Array()[0].String()) {
				resp.Result = `"0x00"`
				resp.Notify = "{}"
			}
			return resp
		})
		wsURL := server.WSURL()
		rpc := client.NewRPCClient(lggr, *wsURL, nil, "rpc", 1, chainId, commonclient.Primary)
		defer rpc.Close()
		require.NoError(t, rpc.Dial(ctx))
		sub, err := rpc.SubscribeFilterLogs(ctx, ethereum.FilterQuery{}, make(chan types.Log))
		require.NoError(t, err)
		server.MustWriteBinaryMessageSync(t, "invalid msg")
		errorCtx, cancel := context.WithTimeout(ctx, tests.DefaultWaitTimeout)
		defer cancel()
		select {
		case err = <-sub.Err():
			require.ErrorContains(t, err, "RPCClient returned error (rpc): invalid character")
		case <-errorCtx.Done():
			t.Errorf("Expected subscription to return an error, but test timeout instead")
		}
	})
}

func TestRPCClient_LatestFinalizedBlock(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithTimeout(tests.Context(t), tests.WaitTimeout(t))
	defer cancel()

	chainId := big.NewInt(123456)
	lggr := logger.Test(t)

	type rpcServer struct {
		Head *evmtypes.Head
		URL  *url.URL
	}
	createRPCServer := func() *rpcServer {
		server := &rpcServer{}
		server.URL = testutils.NewWSServer(t, chainId, func(method string, params gjson.Result) (resp testutils.JSONRPCResponse) {
			assert.Equal(t, "eth_getBlockByNumber", method)
			if assert.True(t, params.IsArray()) && assert.Equal(t, "finalized", params.Array()[0].String()) {
				head := server.Head
				jsonHead, err := json.Marshal(head)
				if err != nil {
					panic(fmt.Errorf("failed to marshal head: %w", err))
				}
				resp.Result = string(jsonHead)
			}
			return
		}).WSURL()

		return server
	}

	server := createRPCServer()
	rpc := client.NewRPCClient(lggr, *server.URL, nil, "rpc", 1, chainId, commonclient.Primary)
	require.NoError(t, rpc.Dial(ctx))
	defer rpc.Close()
	server.Head = &evmtypes.Head{Number: 128}
	// updates chain info
	_, err := rpc.LatestFinalizedBlock(ctx)
	require.NoError(t, err)
	chainInfo := rpc.GetInterceptedChainInfo()
	assert.Equal(t, int64(0), chainInfo.BlockNumber)
	assert.Equal(t, int64(128), chainInfo.FinalizedBlockNumber)

	// lower block number does not update Highest
	server.Head = &evmtypes.Head{Number: 127}
	_, err = rpc.LatestFinalizedBlock(ctx)
	require.NoError(t, err)
	chainInfo = rpc.GetInterceptedChainInfo()
	assert.Equal(t, int64(0), chainInfo.BlockNumber)
	assert.Equal(t, int64(128), chainInfo.FinalizedBlockNumber)
}
