package client_test

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"

	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	commonclient "github.com/smartcontractkit/chainlink/v2/common/client"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
)

func TestRPCClient_SubscribeNewHead(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithTimeout(testutils.Context(t), testutils.WaitTimeout(t))
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
		require.NoError(t, rpc.Dial(ctx))
		maxBlockNumber, maxFinalizedBlockNumber := rpc.GetInterceptedChainInfo()
		require.Equal(t, int64(0), maxBlockNumber)
		require.Equal(t, int64(0), maxFinalizedBlockNumber)
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
		maxBlockNumber, maxFinalizedBlockNumber = rpc.GetInterceptedChainInfo()
		assert.Equal(t, int64(256), maxBlockNumber)
		assert.Equal(t, int64(0), maxFinalizedBlockNumber)
	})
	t.Run("Block's chain ID matched configured", func(t *testing.T) {
		server := createRPCServer()
		rpc := client.NewRPCClient(lggr, *server.URL, nil, "rpc", 1, chainId, commonclient.Primary)
		require.NoError(t, rpc.Dial(ctx))
		server.Head = &evmtypes.Head{
			Number: 256,
		}
		head := receiveNewHead(rpc)
		require.Equal(t, chainId, head.ChainID())
	})
}

func TestRPCClient_LatestFinalizedBlock(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithTimeout(testutils.Context(t), testutils.WaitTimeout(t))
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
	server.Head = &evmtypes.Head{Number: 128}
	// updates chain info
	_, err := rpc.LatestFinalizedBlock(ctx)
	require.NoError(t, err)
	maxBlockNumber, maxFinalizedBlockNumber := rpc.GetInterceptedChainInfo()
	assert.Equal(t, int64(0), maxBlockNumber)
	assert.Equal(t, int64(128), maxFinalizedBlockNumber)

	// lower block number does not update Highest
	server.Head = &evmtypes.Head{Number: 127}
	_, err = rpc.LatestFinalizedBlock(ctx)
	require.NoError(t, err)
	maxBlockNumber, maxFinalizedBlockNumber = rpc.GetInterceptedChainInfo()
	assert.Equal(t, int64(0), maxBlockNumber)
	assert.Equal(t, int64(128), maxFinalizedBlockNumber)

}
