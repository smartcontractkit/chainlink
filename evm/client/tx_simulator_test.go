package client_test

import (
	"testing"

	"github.com/ethereum/go-ethereum"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"

	"github.com/smartcontractkit/chainlink/v2/evm/client"
	"github.com/smartcontractkit/chainlink/v2/evm/testutils"
)

func TestSimulateTx_Default(t *testing.T) {
	t.Parallel()

	fromAddress := testutils.NewAddress()
	toAddress := testutils.NewAddress()

	t.Run("returns without error if simulation passes", func(t *testing.T) {
		ctx := tests.Context(t)
		wsURL := testutils.NewWSServer(t, testutils.FixtureChainID, func(method string, params gjson.Result) (resp testutils.JSONRPCResponse) {
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

		msg := ethereum.CallMsg{
			From: fromAddress,
			To:   &toAddress,
			Data: []byte("0x00"),
		}
		sendErr := client.SimulateTransaction(ctx, ethClient, logger.TestSugared(t), "", msg)
		require.Empty(t, sendErr)
	})

	t.Run("returns error if simulation returns zk out-of-counters error", func(t *testing.T) {
		ctx := tests.Context(t)
		wsURL := testutils.NewWSServer(t, testutils.FixtureChainID, func(method string, params gjson.Result) (resp testutils.JSONRPCResponse) {
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

		msg := ethereum.CallMsg{
			From: fromAddress,
			To:   &toAddress,
			Data: []byte("0x00"),
		}
		sendErr := client.SimulateTransaction(ctx, ethClient, logger.TestSugared(t), "", msg)
		require.Equal(t, true, sendErr.IsTerminallyStuckConfigError(nil))
	})

	t.Run("returns without error if simulation returns non-OOC error", func(t *testing.T) {
		ctx := tests.Context(t)
		wsURL := testutils.NewWSServer(t, testutils.FixtureChainID, func(method string, params gjson.Result) (resp testutils.JSONRPCResponse) {
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
				resp.Error.Message = "something went wrong"
			}
			return
		}).WSURL().String()

		ethClient := mustNewChainClient(t, wsURL)

		msg := ethereum.CallMsg{
			From: fromAddress,
			To:   &toAddress,
			Data: []byte("0x00"),
		}
		sendErr := client.SimulateTransaction(ctx, ethClient, logger.TestSugared(t), "", msg)
		require.Equal(t, false, sendErr.IsTerminallyStuckConfigError(nil))
	})
}
