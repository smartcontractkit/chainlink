package client_test

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
)

func TestNullClient(t *testing.T) {
	t.Parallel()

	t.Run("chain id", func(t *testing.T) {
		lggr := logger.Test(t)
		cid := big.NewInt(123)
		nc := client.NewNullClient(cid, lggr)
		require.Equal(t, cid, nc.ConfiguredChainID())

		nc = client.NewNullClient(nil, lggr)
		require.Equal(t, big.NewInt(client.NullClientChainID), nc.ConfiguredChainID())
	})

	t.Run("CL client methods", func(t *testing.T) {
		lggr, logs := logger.TestObserved(t, zapcore.DebugLevel)
		nc := client.NewNullClient(nil, lggr)
		ctx := tests.Context(t)

		err := nc.Dial(ctx)
		require.NoError(t, err)
		require.Equal(t, 1, logs.FilterMessage("Dial").Len())

		nc.Close()
		require.Equal(t, 1, logs.FilterMessage("Close").Len())

		b, err := nc.TokenBalance(ctx, common.Address{}, common.Address{})
		require.NoError(t, err)
		require.Zero(t, b.Int64())
		require.Equal(t, 1, logs.FilterMessage("TokenBalance").Len())

		l, err := nc.LINKBalance(ctx, common.Address{}, common.Address{})
		require.NoError(t, err)
		require.True(t, l.IsZero())
		require.Equal(t, 1, logs.FilterMessage("LINKBalance").Len())

		err = nc.CallContext(ctx, nil, "")
		require.NoError(t, err)
		require.Equal(t, 1, logs.FilterMessage("CallContext").Len())

		h, err := nc.HeadByNumber(ctx, nil)
		require.NoError(t, err)
		require.Nil(t, h)
		require.Equal(t, 1, logs.FilterMessage("HeadByNumber").Len())

		chHeads := make(chan *evmtypes.Head)
		sub, err := nc.SubscribeNewHead(ctx, chHeads)
		require.NoError(t, err)
		require.Equal(t, 1, logs.FilterMessage("SubscribeNewHead").Len())
		require.Nil(t, sub.Err())
		require.Equal(t, 1, logs.FilterMessage("Err").Len())
		sub.Unsubscribe()
		require.Equal(t, 1, logs.FilterMessage("Unsubscribe").Len())

		chLogs := make(chan types.Log)
		_, err = nc.SubscribeFilterLogs(ctx, ethereum.FilterQuery{}, chLogs)
		require.NoError(t, err)
		require.Equal(t, 1, logs.FilterMessage("SubscribeFilterLogs").Len())
	})

	t.Run("Geth client methods", func(t *testing.T) {
		lggr, logs := logger.TestObserved(t, zapcore.DebugLevel)
		nc := client.NewNullClient(nil, lggr)
		ctx := tests.Context(t)

		h, err := nc.HeaderByNumber(ctx, nil)
		require.NoError(t, err)
		require.Nil(t, h)
		require.Equal(t, 1, logs.FilterMessage("HeaderByNumber").Len())

		err = nc.SendTransaction(ctx, nil)
		require.NoError(t, err)
		require.Equal(t, 1, logs.FilterMessage("SendTransaction").Len())

		c, err := nc.PendingCodeAt(ctx, common.Address{})
		require.NoError(t, err)
		require.Empty(t, c)
		require.Equal(t, 1, logs.FilterMessage("PendingCodeAt").Len())

		n, err := nc.PendingNonceAt(ctx, common.Address{})
		require.NoError(t, err)
		require.Zero(t, n)
		require.Equal(t, 1, logs.FilterMessage("PendingNonceAt").Len())

		s, err := nc.SequenceAt(ctx, common.Address{}, nil)
		require.NoError(t, err)
		require.Zero(t, s)
		require.Equal(t, 1, logs.FilterMessage("SequenceAt").Len())

		r, err := nc.TransactionReceipt(ctx, common.Hash{})
		require.NoError(t, err)
		require.Nil(t, r)
		require.Equal(t, 1, logs.FilterMessage("TransactionReceipt").Len())

		b, err := nc.BlockByNumber(ctx, nil)
		require.NoError(t, err)
		require.Nil(t, b)
		require.Equal(t, 1, logs.FilterMessage("BlockByNumber").Len())

		b, err = nc.BlockByHash(ctx, common.Hash{})
		require.NoError(t, err)
		require.Nil(t, b)
		require.Equal(t, 1, logs.FilterMessage("BlockByHash").Len())

		bal, err := nc.BalanceAt(ctx, common.Address{}, nil)
		require.NoError(t, err)
		require.Zero(t, bal.Int64())
		require.Equal(t, 1, logs.FilterMessage("BalanceAt").Len())

		log, err := nc.FilterLogs(ctx, ethereum.FilterQuery{})
		require.NoError(t, err)
		require.Nil(t, log)
		require.Equal(t, 1, logs.FilterMessage("FilterLogs").Len())

		gas, err := nc.EstimateGas(ctx, ethereum.CallMsg{})
		require.NoError(t, err)
		require.Zero(t, gas)
		require.Equal(t, 1, logs.FilterMessage("EstimateGas").Len())

		gp, err := nc.SuggestGasPrice(ctx)
		require.NoError(t, err)
		require.Zero(t, gp.Int64())
		require.Equal(t, 1, logs.FilterMessage("SuggestGasPrice").Len())

		cc, err := nc.CallContract(ctx, ethereum.CallMsg{}, nil)
		require.NoError(t, err)
		require.Nil(t, cc)
		require.Equal(t, 1, logs.FilterMessage("CallContract").Len())

		ca, err := nc.CodeAt(ctx, common.Address{}, nil)
		require.NoError(t, err)
		require.Nil(t, ca)
		require.Equal(t, 1, logs.FilterMessage("CodeAt").Len())

		err = nc.BatchCallContext(ctx, []rpc.BatchElem{})
		require.NoError(t, err)

		err = nc.BatchCallContextAll(ctx, []rpc.BatchElem{})
		require.NoError(t, err)

		tip, err := nc.SuggestGasTipCap(ctx)
		require.NoError(t, err)
		require.Nil(t, tip)

		m := nc.NodeStates()
		require.Nil(t, m)
	})
}
