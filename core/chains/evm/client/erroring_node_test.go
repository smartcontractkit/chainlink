package client

import (
	"testing"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
)

func TestErroringNode(t *testing.T) {
	t.Parallel()

	ctx := testutils.Context(t)
	n := &erroringNode{
		"boo",
	}

	require.Nil(t, n.ChainID())
	err := n.Start(ctx)
	require.Equal(t, n.errMsg, err.Error())

	defer func() { assert.NoError(t, n.Close()) }()

	err = n.Verify(ctx, nil)
	require.Equal(t, n.errMsg, err.Error())

	err = n.CallContext(ctx, nil, "")
	require.Equal(t, n.errMsg, err.Error())

	err = n.BatchCallContext(ctx, nil)
	require.Equal(t, n.errMsg, err.Error())

	err = n.SendTransaction(ctx, nil)
	require.Equal(t, n.errMsg, err.Error())

	_, err = n.PendingCodeAt(ctx, common.Address{})
	require.Equal(t, n.errMsg, err.Error())

	_, err = n.PendingNonceAt(ctx, common.Address{})
	require.Equal(t, n.errMsg, err.Error())

	_, err = n.NonceAt(ctx, common.Address{}, nil)
	require.Equal(t, n.errMsg, err.Error())

	_, err = n.TransactionReceipt(ctx, common.Hash{})
	require.Equal(t, n.errMsg, err.Error())

	_, err = n.BlockByNumber(ctx, nil)
	require.Equal(t, n.errMsg, err.Error())

	_, err = n.BlockByHash(ctx, common.Hash{})
	require.Equal(t, n.errMsg, err.Error())

	_, err = n.BalanceAt(ctx, common.Address{}, nil)
	require.Equal(t, n.errMsg, err.Error())

	_, err = n.FilterLogs(ctx, ethereum.FilterQuery{})
	require.Equal(t, n.errMsg, err.Error())

	_, err = n.SubscribeFilterLogs(ctx, ethereum.FilterQuery{}, nil)
	require.Equal(t, n.errMsg, err.Error())

	_, err = n.EstimateGas(ctx, ethereum.CallMsg{})
	require.Equal(t, n.errMsg, err.Error())

	_, err = n.SuggestGasPrice(ctx)
	require.Equal(t, n.errMsg, err.Error())

	_, err = n.CallContract(ctx, ethereum.CallMsg{}, nil)
	require.Equal(t, n.errMsg, err.Error())

	_, err = n.CodeAt(ctx, common.Address{}, nil)
	require.Equal(t, n.errMsg, err.Error())

	_, err = n.HeaderByNumber(ctx, nil)
	require.Equal(t, n.errMsg, err.Error())

	_, err = n.SuggestGasTipCap(ctx)
	require.Equal(t, n.errMsg, err.Error())

	_, err = n.EthSubscribe(ctx, nil)
	require.Equal(t, n.errMsg, err.Error())

	require.Equal(t, "<erroring node>", n.String())
	require.Equal(t, NodeStateUnreachable, n.State())

	state, num, _ := n.StateAndLatest()
	require.Equal(t, NodeStateUnreachable, state)
	require.Equal(t, int64(-1), num)

	n.DeclareInSync()
	n.DeclareOutOfSync()
	n.DeclareUnreachable()

	require.Zero(t, n.Name())
	require.Nil(t, n.NodeStates())
}
