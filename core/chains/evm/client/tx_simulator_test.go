package client_test

import (
	"testing"

	"github.com/ethereum/go-ethereum"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	"github.com/smartcontractkit/chainlink/v2/common/config"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
)

func TestSimulateTx_Default(t *testing.T) {
	t.Parallel()

	fromAddress := testutils.NewAddress()
	toAddress := testutils.NewAddress()
	ctx := testutils.Context(t)

	t.Run("returns without error if simulation passes", func(t *testing.T) {
		wsURL := testutils.NewWSServer(t, &cltest.FixtureChainID, func(method string, params gjson.Result) (resp testutils.JSONRPCResponse) {
			switch method {
			case "eth_subscribe":
				resp.Result = `"0x00"`
				resp.Notify = headResult
				return
			case "eth_unsubscribe":
				resp.Result = "true"
				return
			case "eth_estimateGas":
				resp.Result = `"0x100"`
			}
			return
		}).WSURL().String()

		ethClient := mustNewChainClient(t, wsURL)
		err := ethClient.Dial(ctx)
		require.NoError(t, err)

		msg := ethereum.CallMsg{
			From: fromAddress,
			To:   &toAddress,
			Data: []byte("0x00"),
		}
		err = client.SimulateTransaction(ctx, ethClient, logger.TestSugared(t), "", msg)
		require.NoError(t, err)
	})

	t.Run("returns error if simulation returns zk out-of-counters error", func(t *testing.T) {
		wsURL := testutils.NewWSServer(t, &cltest.FixtureChainID, func(method string, params gjson.Result) (resp testutils.JSONRPCResponse) {
			switch method {
			case "eth_subscribe":
				resp.Result = `"0x00"`
				resp.Notify = headResult
				return
			case "eth_unsubscribe":
				resp.Result = "true"
				return
			case "eth_estimateGas":
				resp.Error.Code = -32000
				resp.Result = `"0x100"`
				resp.Error.Message = "not enough keccak counters to continue the execution"
			}
			return
		}).WSURL().String()

		ethClient := mustNewChainClient(t, wsURL)
		err := ethClient.Dial(ctx)
		require.NoError(t, err)

		msg := ethereum.CallMsg{
			From: fromAddress,
			To:   &toAddress,
			Data: []byte("0x00"),
		}
		err = client.SimulateTransaction(ctx, ethClient, logger.TestSugared(t), "", msg)
		require.Error(t, err, client.ErrOutOfCounters)
	})

	t.Run("returns without error if simulation returns non-OOC error", func(t *testing.T) {
		wsURL := testutils.NewWSServer(t, &cltest.FixtureChainID, func(method string, params gjson.Result) (resp testutils.JSONRPCResponse) {
			switch method {
			case "eth_subscribe":
				resp.Result = `"0x00"`
				resp.Notify = headResult
				return
			case "eth_unsubscribe":
				resp.Result = "true"
				return
			case "eth_estimateGas":
				resp.Error.Code = -32000
				resp.Result = `"0x100"`
				resp.Error.Message = "txpool is full"
			}
			return
		}).WSURL().String()

		ethClient := mustNewChainClient(t, wsURL)
		err := ethClient.Dial(ctx)
		require.NoError(t, err)

		msg := ethereum.CallMsg{
			From: fromAddress,
			To:   &toAddress,
			Data: []byte("0x00"),
		}
		err = client.SimulateTransaction(ctx, ethClient, logger.TestSugared(t), "", msg)
		require.NoError(t, err)
	})
}

func TestSimulateTx_ZkEvm(t *testing.T) {
	t.Parallel()

	fromAddress := testutils.NewAddress()
	toAddress := testutils.NewAddress()
	ctx := testutils.Context(t)

	t.Run("returns without error if simulation passes", func(t *testing.T) {
		wsURL := testutils.NewWSServer(t, &cltest.FixtureChainID, func(method string, params gjson.Result) (resp testutils.JSONRPCResponse) {
			switch method {
			case "eth_subscribe":
				resp.Result = `"0x00"`
				resp.Notify = headResult
				return
			case "eth_unsubscribe":
				resp.Result = "true"
				return
			case "zkevm_estimateCounters":
				resp.Result = `{
					"countersUsed": {
						"gasUsed": "0x5360",
						"usedKeccakHashes": "0x7",
						"usedPoseidonHashes": "0x2bb",
						"usedPoseidonPaddings": "0x4",
						"usedMemAligns": "0x0",
						"usedArithmetics": "0x263",
						"usedBinaries": "0x40c",
						"usedSteps": "0x3288",
						"usedSHA256Hashes": "0x0"
					},
					"countersLimit": {
						"maxGasUsed": "0x1c9c380",
						"maxKeccakHashes": "0x861",
						"maxPoseidonHashes": "0x3d9c5",
						"maxPoseidonPaddings": "0x21017",
						"maxMemAligns": "0x39c29",
						"maxArithmetics": "0x39c29",
						"maxBinaries": "0x73852",
						"maxSteps": "0x73846a",
						"maxSHA256Hashes": "0x63c"
					}
				}`
			}
			return
		}).WSURL().String()

		ethClient := mustNewChainClient(t, wsURL)
		err := ethClient.Dial(ctx)
		require.NoError(t, err)

		msg := ethereum.CallMsg{
			From: fromAddress,
			To:   &toAddress,
			Data: []byte("0x00"),
		}
		err = client.SimulateTransaction(ctx, ethClient, logger.TestSugared(t), config.ChainZkEvm, msg)
		require.NoError(t, err)
	})

	t.Run("returns error if simulation returns zk out-of-counters error", func(t *testing.T) {
		wsURL := testutils.NewWSServer(t, &cltest.FixtureChainID, func(method string, params gjson.Result) (resp testutils.JSONRPCResponse) {
			switch method {
			case "eth_subscribe":
				resp.Result = `"0x00"`
				resp.Notify = headResult
				return
			case "eth_unsubscribe":
				resp.Result = "true"
				return
			case "zkevm_estimateCounters":
				resp.Result = `{
					"countersUsed": {
						"gasUsed": "0x12f3bd",
						"usedKeccakHashes": "0x8d3",
						"usedPoseidonHashes": "0x222",
						"usedPoseidonPaddings": "0x16",
						"usedMemAligns": "0x1a69",
						"usedArithmetics": "0x2619",
						"usedBinaries": "0x2d738",
						"usedSteps": "0x72e223",
						"usedSHA256Hashes": "0x0"
					},
					"countersLimit": {
						"maxGasUsed": "0x1c9c380",
						"maxKeccakHashes": "0x861",
						"maxPoseidonHashes": "0x3d9c5",
						"maxPoseidonPaddings": "0x21017",
						"maxMemAligns": "0x39c29",
						"maxArithmetics": "0x39c29",
						"maxBinaries": "0x73852",
						"maxSteps": "0x73846a",
						"maxSHA256Hashes": "0x63c"
					},
					"oocError": "not enough keccak counters to continue the execution"
				}`
			}
			return
		}).WSURL().String()

		ethClient := mustNewChainClient(t, wsURL)
		err := ethClient.Dial(ctx)
		require.NoError(t, err)

		msg := ethereum.CallMsg{
			From: fromAddress,
			To:   &toAddress,
			Data: []byte("0x00"),
		}
		err = client.SimulateTransaction(ctx, ethClient, logger.TestSugared(t), config.ChainZkEvm, msg)
		require.Error(t, err, client.ErrOutOfCounters)
	})
}
