package client_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"net/url"
	"sync"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
	"go.uber.org/zap"

	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	commonclient "github.com/smartcontractkit/chainlink/v2/common/client"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/config/chaintype"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/testutils"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
)

func makeNewHeadWSMessage(head *evmtypes.Head) string {
	asJSON, err := json.Marshal(head)
	if err != nil {
		panic(fmt.Errorf("failed to marshal head: %w", err))
	}
	return fmt.Sprintf(`{"jsonrpc":"2.0","method":"eth_subscription","params":{"subscription":"0x00","result":%s}}`, string(asJSON))
}

func TestRPCClient_SubscribeNewHead(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithTimeout(tests.Context(t), tests.WaitTimeout(t))
	defer cancel()

	chainId := big.NewInt(123456)
	lggr := logger.Test(t)

	serverCallBack := func(method string, params gjson.Result) (resp testutils.JSONRPCResponse) {
		if method == "eth_unsubscribe" {
			resp.Result = "true"
			return
		}
		assert.Equal(t, "eth_subscribe", method)
		if assert.True(t, params.IsArray()) && assert.Equal(t, "newHeads", params.Array()[0].String()) {
			resp.Result = `"0x00"`
		}
		return
	}
	t.Run("Updates chain info on new blocks", func(t *testing.T) {
		server := testutils.NewWSServer(t, chainId, serverCallBack)
		wsURL := server.WSURL()

		rpc := client.NewRPCClient(lggr, *wsURL, nil, "rpc", 1, chainId, commonclient.Primary, 0, commonclient.QueryTimeout, commonclient.QueryTimeout, "")
		defer rpc.Close()
		require.NoError(t, rpc.Dial(ctx))
		// set to default values
		latest, highestUserObservations := rpc.GetInterceptedChainInfo()
		assert.Equal(t, int64(0), latest.BlockNumber)
		assert.Equal(t, int64(0), latest.FinalizedBlockNumber)
		assert.Nil(t, latest.TotalDifficulty)
		assert.Equal(t, int64(0), highestUserObservations.BlockNumber)
		assert.Equal(t, int64(0), highestUserObservations.FinalizedBlockNumber)
		assert.Nil(t, highestUserObservations.TotalDifficulty)

		ch := make(chan *evmtypes.Head)
		sub, err := rpc.SubscribeNewHead(tests.Context(t), ch)
		require.NoError(t, err)
		defer sub.Unsubscribe()
		go server.MustWriteBinaryMessageSync(t, makeNewHeadWSMessage(&evmtypes.Head{Number: 256, TotalDifficulty: big.NewInt(1000)}))
		// received 256 head
		<-ch
		go server.MustWriteBinaryMessageSync(t, makeNewHeadWSMessage(&evmtypes.Head{Number: 128, TotalDifficulty: big.NewInt(500)}))
		// received 128 head
		<-ch

		latest, highestUserObservations = rpc.GetInterceptedChainInfo()
		assert.Equal(t, int64(128), latest.BlockNumber)
		assert.Equal(t, int64(0), latest.FinalizedBlockNumber)
		assert.Equal(t, big.NewInt(500), latest.TotalDifficulty)

		assertHighestUserObservations := func(highestUserObservations commonclient.ChainInfo) {
			assert.Equal(t, int64(256), highestUserObservations.BlockNumber)
			assert.Equal(t, int64(0), highestUserObservations.FinalizedBlockNumber)
			assert.Equal(t, big.NewInt(1000), highestUserObservations.TotalDifficulty)
		}

		assertHighestUserObservations(highestUserObservations)

		// DisconnectAll resets latest
		rpc.DisconnectAll()

		latest, highestUserObservations = rpc.GetInterceptedChainInfo()
		assert.Equal(t, int64(0), latest.BlockNumber)
		assert.Equal(t, int64(0), latest.FinalizedBlockNumber)
		assert.Nil(t, latest.TotalDifficulty)

		assertHighestUserObservations(highestUserObservations)
	})
	t.Run("App layer observations are not affected by new block if health check flag is present", func(t *testing.T) {
		server := testutils.NewWSServer(t, chainId, serverCallBack)
		wsURL := server.WSURL()

		rpc := client.NewRPCClient(lggr, *wsURL, nil, "rpc", 1, chainId, commonclient.Primary, 0, commonclient.QueryTimeout, commonclient.QueryTimeout, "")
		defer rpc.Close()
		require.NoError(t, rpc.Dial(ctx))
		ch := make(chan *evmtypes.Head)
		sub, err := rpc.SubscribeNewHead(commonclient.CtxAddHealthCheckFlag(tests.Context(t)), ch)
		require.NoError(t, err)
		defer sub.Unsubscribe()
		go server.MustWriteBinaryMessageSync(t, makeNewHeadWSMessage(&evmtypes.Head{Number: 256, TotalDifficulty: big.NewInt(1000)}))
		// received 256 head
		<-ch

		latest, highestUserObservations := rpc.GetInterceptedChainInfo()
		assert.Equal(t, int64(256), latest.BlockNumber)
		assert.Equal(t, int64(0), latest.FinalizedBlockNumber)
		assert.Equal(t, big.NewInt(1000), latest.TotalDifficulty)

		assert.Equal(t, int64(0), highestUserObservations.BlockNumber)
		assert.Equal(t, int64(0), highestUserObservations.FinalizedBlockNumber)
		assert.Equal(t, (*big.Int)(nil), highestUserObservations.TotalDifficulty)
	})
	t.Run("Concurrent Unsubscribe and onNewHead calls do not lead to a deadlock", func(t *testing.T) {
		const numberOfAttempts = 1000 // need a large number to increase the odds of reproducing the issue
		server := testutils.NewWSServer(t, chainId, serverCallBack)
		wsURL := server.WSURL()

		rpc := client.NewRPCClient(lggr, *wsURL, nil, "rpc", 1, chainId, commonclient.Primary, 0, commonclient.QueryTimeout, commonclient.QueryTimeout, "")
		defer rpc.Close()
		require.NoError(t, rpc.Dial(ctx))
		var wg sync.WaitGroup
		for i := 0; i < numberOfAttempts; i++ {
			ch := make(chan *evmtypes.Head)
			sub, err := rpc.SubscribeNewHead(tests.Context(t), ch)
			require.NoError(t, err)
			wg.Add(2)
			go func() {
				server.MustWriteBinaryMessageSync(t, makeNewHeadWSMessage(&evmtypes.Head{Number: 256, TotalDifficulty: big.NewInt(1000)}))
				wg.Done()
			}()
			go func() {
				rpc.UnsubscribeAllExceptAliveLoop()
				sub.Unsubscribe()
				wg.Done()
			}()
			wg.Wait()
		}
	})
	t.Run("Block's chain ID matched configured", func(t *testing.T) {
		server := testutils.NewWSServer(t, chainId, serverCallBack)
		wsURL := server.WSURL()
		rpc := client.NewRPCClient(lggr, *wsURL, nil, "rpc", 1, chainId, commonclient.Primary, 0, commonclient.QueryTimeout, commonclient.QueryTimeout, "")
		defer rpc.Close()
		require.NoError(t, rpc.Dial(ctx))
		ch := make(chan *evmtypes.Head)
		sub, err := rpc.SubscribeNewHead(tests.Context(t), ch)
		require.NoError(t, err)
		defer sub.Unsubscribe()
		go server.MustWriteBinaryMessageSync(t, makeNewHeadWSMessage(&evmtypes.Head{Number: 256}))
		head := <-ch
		require.Equal(t, chainId, head.ChainID())
	})
	t.Run("Failed SubscribeNewHead returns and logs proper error", func(t *testing.T) {
		server := testutils.NewWSServer(t, chainId, func(reqMethod string, reqParams gjson.Result) (resp testutils.JSONRPCResponse) {
			return resp
		})
		wsURL := server.WSURL()
		observedLggr, observed := logger.TestObserved(t, zap.DebugLevel)
		rpc := client.NewRPCClient(observedLggr, *wsURL, nil, "rpc", 1, chainId, commonclient.Primary, 0, commonclient.QueryTimeout, commonclient.QueryTimeout, "")
		require.NoError(t, rpc.Dial(ctx))
		server.Close()
		_, err := rpc.SubscribeNewHead(ctx, make(chan *evmtypes.Head))
		require.ErrorContains(t, err, "RPCClient returned error (rpc)")
		tests.AssertLogEventually(t, observed, "evmclient.Client#EthSubscribe RPC call failure")
	})
	t.Run("Subscription error is properly wrapper", func(t *testing.T) {
		server := testutils.NewWSServer(t, chainId, serverCallBack)
		wsURL := server.WSURL()
		rpc := client.NewRPCClient(lggr, *wsURL, nil, "rpc", 1, chainId, commonclient.Primary, 0, commonclient.QueryTimeout, commonclient.QueryTimeout, "")
		defer rpc.Close()
		require.NoError(t, rpc.Dial(ctx))
		sub, err := rpc.SubscribeNewHead(ctx, make(chan *evmtypes.Head))
		require.NoError(t, err)
		go server.MustWriteBinaryMessageSync(t, "invalid msg")
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
	t.Run("Failed SubscribeFilterLogs logs and returns proper error", func(t *testing.T) {
		server := testutils.NewWSServer(t, chainId, func(reqMethod string, reqParams gjson.Result) (resp testutils.JSONRPCResponse) {
			return resp
		})
		wsURL := server.WSURL()
		observedLggr, observed := logger.TestObserved(t, zap.DebugLevel)
		rpc := client.NewRPCClient(observedLggr, *wsURL, nil, "rpc", 1, chainId, commonclient.Primary, 0, commonclient.QueryTimeout, commonclient.QueryTimeout, "")
		require.NoError(t, rpc.Dial(ctx))
		server.Close()
		_, err := rpc.SubscribeFilterLogs(ctx, ethereum.FilterQuery{}, make(chan types.Log))
		require.ErrorContains(t, err, "RPCClient returned error (rpc)")
		tests.AssertLogEventually(t, observed, "evmclient.Client#SubscribeFilterLogs RPC call failure")
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
		rpc := client.NewRPCClient(lggr, *wsURL, nil, "rpc", 1, chainId, commonclient.Primary, 0, commonclient.QueryTimeout, commonclient.QueryTimeout, "")
		defer rpc.Close()
		require.NoError(t, rpc.Dial(ctx))
		sub, err := rpc.SubscribeFilterLogs(ctx, ethereum.FilterQuery{}, make(chan types.Log))
		require.NoError(t, err)
		go server.MustWriteBinaryMessageSync(t, "invalid msg")
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
	rpc := client.NewRPCClient(lggr, *server.URL, nil, "rpc", 1, chainId, commonclient.Primary, 0, commonclient.QueryTimeout, commonclient.QueryTimeout, "")
	require.NoError(t, rpc.Dial(ctx))
	defer rpc.Close()
	server.Head = &evmtypes.Head{Number: 128}
	// updates chain info
	_, err := rpc.LatestFinalizedBlock(ctx)
	require.NoError(t, err)
	latest, highestUserObservations := rpc.GetInterceptedChainInfo()

	assert.Equal(t, int64(0), highestUserObservations.BlockNumber)
	assert.Equal(t, int64(128), highestUserObservations.FinalizedBlockNumber)

	assert.Equal(t, int64(0), latest.BlockNumber)
	assert.Equal(t, int64(128), latest.FinalizedBlockNumber)

	// lower block number does not update highestUserObservations
	server.Head = &evmtypes.Head{Number: 127}
	_, err = rpc.LatestFinalizedBlock(ctx)
	require.NoError(t, err)
	latest, highestUserObservations = rpc.GetInterceptedChainInfo()

	assert.Equal(t, int64(0), highestUserObservations.BlockNumber)
	assert.Equal(t, int64(128), highestUserObservations.FinalizedBlockNumber)

	assert.Equal(t, int64(0), latest.BlockNumber)
	assert.Equal(t, int64(127), latest.FinalizedBlockNumber)

	// health check flg prevents change in highestUserObservations
	server.Head = &evmtypes.Head{Number: 256}
	_, err = rpc.LatestFinalizedBlock(commonclient.CtxAddHealthCheckFlag(ctx))
	require.NoError(t, err)
	latest, highestUserObservations = rpc.GetInterceptedChainInfo()

	assert.Equal(t, int64(0), highestUserObservations.BlockNumber)
	assert.Equal(t, int64(128), highestUserObservations.FinalizedBlockNumber)

	assert.Equal(t, int64(0), latest.BlockNumber)
	assert.Equal(t, int64(256), latest.FinalizedBlockNumber)

	// DisconnectAll resets latest ChainInfo
	rpc.DisconnectAll()
	latest, highestUserObservations = rpc.GetInterceptedChainInfo()
	assert.Equal(t, int64(0), highestUserObservations.BlockNumber)
	assert.Equal(t, int64(128), highestUserObservations.FinalizedBlockNumber)

	assert.Equal(t, int64(0), latest.BlockNumber)
	assert.Equal(t, int64(0), latest.FinalizedBlockNumber)
}

func TestRpcClientLargePayloadTimeout(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name string
		Fn   func(ctx context.Context, rpc client.RPCClient) error
	}{
		{
			Name: "SendTransaction",
			Fn: func(ctx context.Context, rpc client.RPCClient) error {
				return rpc.SendTransaction(ctx, types.NewTx(&types.LegacyTx{}))
			},
		},
		{
			Name: "EstimateGas",
			Fn: func(ctx context.Context, rpc client.RPCClient) error {
				_, err := rpc.EstimateGas(ctx, ethereum.CallMsg{})
				return err
			},
		},
		{
			Name: "CallContract",
			Fn: func(ctx context.Context, rpc client.RPCClient) error {
				_, err := rpc.CallContract(ctx, ethereum.CallMsg{}, nil)
				return err
			},
		},
		{
			Name: "CallContext",
			Fn: func(ctx context.Context, rpc client.RPCClient) error {
				err := rpc.CallContext(ctx, nil, "rpc_call", nil)
				return err
			},
		},
		{
			Name: "BatchCallContext",
			Fn: func(ctx context.Context, rpc client.RPCClient) error {
				err := rpc.BatchCallContext(ctx, nil)
				return err
			},
		},
	}
	for _, testCase := range testCases {
		testCase := testCase
		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()
			// use background context to ensure that the DeadlineExceeded is caused by timeout we've set on request
			// level, instead of one that was set on test level.
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			chainId := big.NewInt(123456)
			rpcURL := testutils.NewWSServer(t, chainId, func(method string, params gjson.Result) (resp testutils.JSONRPCResponse) {
				// block until test is done
				<-ctx.Done()
				return
			}).WSURL()

			// use something unreasonably large for RPC timeout to ensure that we use largePayloadRPCTimeout
			const rpcTimeout = time.Hour
			const largePayloadRPCTimeout = tests.TestInterval
			rpc := client.NewRPCClient(logger.Test(t), *rpcURL, nil, "rpc", 1, chainId, commonclient.Primary, 0, largePayloadRPCTimeout, rpcTimeout, "")
			require.NoError(t, rpc.Dial(ctx))
			defer rpc.Close()
			err := testCase.Fn(ctx, rpc)
			assert.True(t, errors.Is(err, context.DeadlineExceeded), fmt.Sprintf("Expected DedlineExceeded error, but got: %v", err))
		})
	}
}

func TestAstarCustomFinality(t *testing.T) {
	t.Parallel()

	chainId := big.NewInt(123456)
	// create new server that returns 4 block for Astar custom finality and 8 block for finality tag.
	wsURL := testutils.NewWSServer(t, chainId, func(method string, params gjson.Result) (resp testutils.JSONRPCResponse) {
		switch method {
		case "chain_getFinalizedHead":
			resp.Result = `"0xf14c499253fd7bbcba142e5dd77dad8b5ad598c1dc414a66bacdd8dae14a6759"`
		case "chain_getHeader":
			if assert.True(t, params.IsArray()) && assert.Equal(t, "0xf14c499253fd7bbcba142e5dd77dad8b5ad598c1dc414a66bacdd8dae14a6759", params.Array()[0].String()) {
				resp.Result = `{"parentHash":"0x1311773bc6b4efc8f438ed1f094524b2a1233baf8a35396f641fcc42a378fc62","number":"0x4","stateRoot":"0x0e4920dc5516b587e1f74a0b65963134523a12cc11478bb314e52895758fbfa2","extrinsicsRoot":"0x5b02446dcab0659eb07d4a38f28f181c1b78a71b2aba207bb0ea1f0f3468e6bd","digest":{"logs":["0x066175726120ad678e0800000000","0x04525053529023158dc8e8fd0180bf26d88233a3d94eed2f4e43480395f0809f28791965e4d34e9b3905","0x0466726f6e88017441e97acf83f555e0deefef86db636bc8a37eb84747603412884e4df4d2280400","0x056175726101018a0a57edf70cc5474323114a47ee1e7f645b8beea5a1560a996416458e89f42bdf4955e24d32b5da54e1bf628aaa7ce4b8c0fa2b95c175a139d88786af12a88c"]}}`
			}
		case "eth_getBlockByNumber":
			assert.True(t, params.IsArray())
			switch params.Array()[0].String() {
			case "0x4":
				resp.Result = `{"author":"0x5accb3bf9194a5f81b2087d4bd6ac47c62775d49","baseFeePerGas":"0xb576270823","difficulty":"0x0","extraData":"0x","gasLimit":"0xe4e1c0","gasUsed":"0x0","hash":"0x7441e97acf83f555e0deefef86db636bc8a37eb84747603412884e4df4d22804","logsBloom":"0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000","miner":"0x5accb3bf9194a5f81b2087d4bd6ac47c62775d49","nonce":"0x0000000000000000","number":"0x4","parentHash":"0x6ba069c318b692bf2cc0bd7ea070a9382a20c2f52413c10554b57c2e381bf2bb","receiptsRoot":"0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421","sha3Uncles":"0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347","size":"0x201","stateRoot":"0x17c46d359b9af773312c747f1d20032c67658d9a2923799f00533b73789cf49b","timestamp":"0x66acdc22","totalDifficulty":"0x0","transactions":[],"transactionsRoot":"0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421","uncles":[]}`
			case "finalized":
				resp.Result = `{"author":"0x1687736326c9fea17e25fc5287613693c912909c","baseFeePerGas":"0x3b9aca00","difficulty":"0x0","extraData":"0x","gasLimit":"0xe4e1c0","gasUsed":"0x0","hash":"0x62f03413681948b06882e7d9f91c4949bc39ded98d36336ab03faea038ec8e3d","logsBloom":"0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000","miner":"0x1687736326c9fea17e25fc5287613693c912909c","nonce":"0x0000000000000000","number":"0x8","parentHash":"0x43f504afdc639cbb8daf5fd5328a37762164b73f9c70ed54e1928c1fca6d8f23","receiptsRoot":"0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421","sha3Uncles":"0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347","size":"0x200","stateRoot":"0x0cb938d51ad83bdf401e3f5f7f989e60df64fdea620d394af41a3e72629f7495","timestamp":"0x61bd8d1a","totalDifficulty":"0x0","transactions":[],"transactionsRoot":"0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421","uncles":[]}`
			default:
				assert.Fail(t, fmt.Sprintf("unexpected eth_getBlockByNumber param: %v", params.Array()))
			}
		default:
			assert.Fail(t, fmt.Sprintf("unexpected method: %s", method))
		}
		return
	}).WSURL()

	const expectedFinalizedBlockNumber = int64(4)
	const expectedFinalizedBlockHash = "0x7441e97acf83f555e0deefef86db636bc8a37eb84747603412884e4df4d22804"
	rpcClient := client.NewRPCClient(logger.Test(t), *wsURL, nil, "rpc", 1, chainId, commonclient.Primary, 0, commonclient.QueryTimeout, commonclient.QueryTimeout, chaintype.ChainAstar)
	defer rpcClient.Close()
	err := rpcClient.Dial(tests.Context(t))
	require.NoError(t, err)

	testCases := []struct {
		Name               string
		GetLatestFinalized func(ctx context.Context) (*evmtypes.Head, error)
	}{
		{
			Name: "Direct LatestFinalized call",
			GetLatestFinalized: func(ctx context.Context) (*evmtypes.Head, error) {
				return rpcClient.LatestFinalizedBlock(ctx)
			},
		},
		{
			Name: "BatchCallContext with Finalized tag as string",
			GetLatestFinalized: func(ctx context.Context) (*evmtypes.Head, error) {
				result := &evmtypes.Head{}
				req := rpc.BatchElem{
					Method: "eth_getBlockByNumber",
					Args:   []interface{}{rpc.FinalizedBlockNumber.String(), false},
					Result: result,
				}
				err := rpcClient.BatchCallContext(ctx, []rpc.BatchElem{
					req,
				})
				if err != nil {
					return nil, err
				}

				return result, req.Error
			},
		},
		{
			Name: "BatchCallContext with Finalized tag as BlockNumber",
			GetLatestFinalized: func(ctx context.Context) (*evmtypes.Head, error) {
				result := &evmtypes.Head{}
				req := rpc.BatchElem{
					Method: "eth_getBlockByNumber",
					Args:   []interface{}{rpc.FinalizedBlockNumber, false},
					Result: result,
				}
				err := rpcClient.BatchCallContext(ctx, []rpc.BatchElem{req})
				if err != nil {
					return nil, err
				}

				return result, req.Error
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			lf, err := testCase.GetLatestFinalized(tests.Context(t))
			require.NoError(t, err)
			require.NotNil(t, lf)
			assert.Equal(t, expectedFinalizedBlockHash, lf.Hash.String())
			assert.Equal(t, expectedFinalizedBlockNumber, lf.Number)
		})
	}
}
