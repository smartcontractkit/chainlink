package client_test

import (
	"errors"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	commonclient "github.com/smartcontractkit/chainlink/v2/common/client"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/client/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
)

func newMockRpc(t *testing.T) *mocks.RPCClient {
	mockRpc := mocks.NewRPCClient(t)
	mockRpc.On("Dial", mock.Anything).Return(nil).Once()
	mockRpc.On("Close").Return(nil).Once()
	mockRpc.On("ChainID", mock.Anything).Return(testutils.FixtureChainID, nil).Once()
	// node does not always manage to fully setup aliveLoop, so we have to make calls optional to avoid flakes
	mockRpc.On("Subscribe", mock.Anything, mock.Anything, mock.Anything).Return(client.NewMockSubscription(), nil).Maybe()
	mockRpc.On("SetAliveLoopSub", mock.Anything).Return().Maybe()
	return mockRpc
}

func TestChainClient_BatchCallContext(t *testing.T) {
	t.Parallel()

	t.Run("batch requests return errors", func(t *testing.T) {
		ctx := testutils.Context(t)
		rpcError := errors.New("something went wrong")
		blockNumResp := ""
		blockNum := hexutil.EncodeBig(big.NewInt(42))
		b := []rpc.BatchElem{
			{
				Method: "eth_getBlockByNumber",
				Args:   []interface{}{blockNum, true},
				Result: &types.Block{},
			},
			{
				Method: "eth_blockNumber",
				Result: &blockNumResp,
			},
		}

		mockRpc := newMockRpc(t)
		mockRpc.On("BatchCallContext", mock.Anything, b).Run(func(args mock.Arguments) {
			reqs := args.Get(1).([]rpc.BatchElem)
			for i := 0; i < len(reqs); i++ {
				elem := &reqs[i]
				elem.Error = rpcError
			}
		}).Return(nil).Once()

		client := client.NewChainClientWithMockedRpc(t, commonclient.NodeSelectionModeRoundRobin, time.Second*0, time.Second*0, testutils.FixtureChainID, mockRpc)
		err := client.Dial(ctx)
		require.NoError(t, err)

		err = client.BatchCallContext(ctx, b)
		require.NoError(t, err)
		for _, elem := range b {
			require.ErrorIs(t, rpcError, elem.Error)
		}
	})
}

func TestChainClient_CheckTxValidity(t *testing.T) {
	t.Parallel()

	fromAddress := testutils.NewAddress()
	toAddress := testutils.NewAddress()
	ctx := testutils.Context(t)

	t.Run("returns without error if simulation passes", func(t *testing.T) {
		mockRpc := newMockRpc(t)
		mockRpc.On("BatchCallContext", mock.Anything, mock.IsType([]rpc.BatchElem{})).Run(func(args mock.Arguments) {
			reqs := args.Get(1).([]rpc.BatchElem)
			for i := 0; i < len(reqs); i++ {
				elem := &reqs[i]
				elem.Result = `"0x100"`
			}
		}).Return(nil).Once()

		ethClient := client.NewChainClientWithMockedRpc(t, commonclient.NodeSelectionModeRoundRobin, time.Second*0, time.Second*0, testutils.FixtureChainID, mockRpc)
		err := ethClient.Dial(ctx)
		require.NoError(t, err)

		sendErr := ethClient.CheckTxValidity(ctx, fromAddress, toAddress, []byte("0x00"))
		require.Empty(t, sendErr)
	})

	t.Run("returns error if simulation returns zk out-of-counters error", func(t *testing.T) {
		mockRpc := newMockRpc(t)
		mockRpc.On("BatchCallContext", mock.Anything, mock.IsType([]rpc.BatchElem{})).Run(func(args mock.Arguments) {
			reqs := args.Get(1).([]rpc.BatchElem)
			for i := 0; i < len(reqs); i++ {
				elem := &reqs[i]
				elem.Error = errors.New("not enough keccak counters to continue the execution")
			}
		}).Return(nil).Once()

		ethClient := client.NewChainClientWithMockedRpc(t, commonclient.NodeSelectionModeRoundRobin, time.Second*0, time.Second*0, testutils.FixtureChainID, mockRpc)
		err := ethClient.Dial(ctx)
		require.NoError(t, err)

		sendErr := ethClient.CheckTxValidity(ctx, fromAddress, toAddress, []byte("0x00"))
		require.Equal(t, true, sendErr.IsOutOfCounters())
	})

	t.Run("returns error if simulation returns non-OOC error but IsOutOfCounters returns false", func(t *testing.T) {
		errorMsg := "something went wrong"
		mockRpc := newMockRpc(t)
		mockRpc.On("BatchCallContext", mock.Anything, mock.IsType([]rpc.BatchElem{})).Run(func(args mock.Arguments) {
			reqs := args.Get(1).([]rpc.BatchElem)
			for i := 0; i < len(reqs); i++ {
				elem := &reqs[i]
				elem.Error = errors.New(errorMsg)
			}
		}).Return(nil).Once()

		ethClient := client.NewChainClientWithMockedRpc(t, commonclient.NodeSelectionModeRoundRobin, time.Second*0, time.Second*0, testutils.FixtureChainID, mockRpc)
		err := ethClient.Dial(ctx)
		require.NoError(t, err)

		sendErr := ethClient.CheckTxValidity(ctx, fromAddress, toAddress, []byte("0x00"))
		require.Equal(t, errorMsg, sendErr.Error())
		require.Equal(t, false, sendErr.IsOutOfCounters())
	})
}

func TestChainClient_BatchCheckTxValidity(t *testing.T) {
	t.Parallel()

	fromAddress1 := testutils.NewAddress()
	toAddress1 := testutils.NewAddress()
	fromAddress2 := testutils.NewAddress()
	toAddress2 := testutils.NewAddress()
	ctx := testutils.Context(t)

	t.Run("returns without error if simulation passes", func(t *testing.T) {
		mockRpc := newMockRpc(t)
		mockRpc.On("BatchCallContext", mock.Anything, mock.IsType([]rpc.BatchElem{})).Run(func(args mock.Arguments) {
			reqs := args.Get(1).([]rpc.BatchElem)
			for i := 0; i < len(reqs); i++ {
				elem := &reqs[i]
				elem.Result = `"0x100"`
			}
		}).Return(nil).Once()

		ethClient := client.NewChainClientWithMockedRpc(t, commonclient.NodeSelectionModeRoundRobin, time.Second*0, time.Second*0, testutils.FixtureChainID, mockRpc)
		err := ethClient.Dial(ctx)
		require.NoError(t, err)

		reqs := []client.TxSimulationRequest{
			{
				From: fromAddress1,
				To:   &toAddress1,
				Data: []byte("0x00"),
			},
			{
				From: fromAddress2,
				Data: []byte("0x01"),
			},
		}
		err = ethClient.BatchCheckTxValidity(ctx, reqs)
		require.NoError(t, err)
		require.Empty(t, reqs[0].Error)
		require.Empty(t, reqs[1].Error)
	})

	t.Run("returns zk out-of-counters error in request", func(t *testing.T) {
		oocError := "not enough keccak counters to continue the execution"
		mockRpc := newMockRpc(t)
		mockRpc.On("BatchCallContext", mock.Anything, mock.IsType([]rpc.BatchElem{})).Run(func(args mock.Arguments) {
			reqs := args.Get(1).([]rpc.BatchElem)
			// Return proper result for a request
			reqs[0].Result = `"0x100"`
			// Return error for a request
			reqs[1].Error = errors.New(oocError)
		}).Return(nil).Once()

		ethClient := client.NewChainClientWithMockedRpc(t, commonclient.NodeSelectionModeRoundRobin, time.Second*0, time.Second*0, testutils.FixtureChainID, mockRpc)
		err := ethClient.Dial(ctx)
		require.NoError(t, err)

		reqs := []client.TxSimulationRequest{
			{
				From: fromAddress1,
				To:   &toAddress1,
				Data: []byte("0x00"),
			},
			{
				From: fromAddress2,
				To:   &toAddress1,
				Data: []byte("0x01"),
			},
		}
		err = ethClient.BatchCheckTxValidity(ctx, reqs)
		require.NoError(t, err)

		// No error for first request
		require.Empty(t, reqs[0].Error)

		// Out-of-counter error for second request
		require.Equal(t, oocError, reqs[1].Error.Error())
		require.Equal(t, true, reqs[1].Error.IsOutOfCounters())
	})

	t.Run("returns other errors in request", func(t *testing.T) {
		errorMsg := "something went wrong"
		mockRpc := newMockRpc(t)
		mockRpc.On("BatchCallContext", mock.Anything, mock.IsType([]rpc.BatchElem{})).Run(func(args mock.Arguments) {
			reqs := args.Get(1).([]rpc.BatchElem)
			// Return proper result for a request
			reqs[0].Result = `"0x100"`
			// Return error for a request
			reqs[1].Error = errors.New(errorMsg)
		}).Return(nil).Once()

		ethClient := client.NewChainClientWithMockedRpc(t, commonclient.NodeSelectionModeRoundRobin, time.Second*0, time.Second*0, testutils.FixtureChainID, mockRpc)
		err := ethClient.Dial(ctx)
		require.NoError(t, err)

		reqs := []client.TxSimulationRequest{
			{
				From: fromAddress1,
				To:   &toAddress1,
				Data: []byte("0x00"),
			},
			{
				From: fromAddress2,
				To:   &toAddress1,
				Data: []byte("0x01"),
			},
		}
		err = ethClient.BatchCheckTxValidity(ctx, reqs)
		require.NoError(t, err)

		// No error for first request
		require.Empty(t, reqs[0].Error)

		// No Out-of-counter error for second request
		require.Equal(t, errorMsg, reqs[1].Error.Error())
	})

	t.Run("returns the proper errors to the associated requests", func(t *testing.T) {
		mockRpc := newMockRpc(t)
		mockRpc.On("BatchCallContext", mock.Anything, mock.IsType([]rpc.BatchElem{})).Run(func(args mock.Arguments) {
			reqs := args.Get(1).([]rpc.BatchElem)
			reqs[0].Error = errors.New("error0")
			reqs[1].Error = errors.New("error1")
			reqs[2].Error = errors.New("error2")
			reqs[3].Error = errors.New("error3")
			reqs[4].Error = errors.New("error4")
			reqs[5].Error = errors.New("error5")
			reqs[6].Error = errors.New("error6")
		}).Return(nil).Once()

		ethClient := client.NewChainClientWithMockedRpc(t, commonclient.NodeSelectionModeRoundRobin, time.Second*0, time.Second*0, testutils.FixtureChainID, mockRpc)
		err := ethClient.Dial(ctx)
		require.NoError(t, err)

		reqs := []client.TxSimulationRequest{
			{
				From: fromAddress1,
				To:   &toAddress1,
				Data: []byte("0x00"),
			},
			{
				From: fromAddress1,
				Data: []byte("0x00"),
			},
			{
				From: fromAddress1,
				To:   &toAddress2,
				Data: []byte("0x00"),
			},
			{
				From: fromAddress2,
				To:   &toAddress2,
				Data: []byte("0x00"),
			},
			{
				From: fromAddress2,
				To:   &toAddress1,
				Data: []byte("0x00"),
			},
			{
				From: fromAddress1,
				To:   &toAddress1,
				Data: []byte("0x01"),
			},
			{
				From: fromAddress1,
				Data: []byte("0x01"),
			},
		}
		err = ethClient.BatchCheckTxValidity(ctx, reqs)
		require.NoError(t, err)

		require.Equal(t, "error0", reqs[0].Error.Error())
		require.Equal(t, "error1", reqs[1].Error.Error())
		require.Equal(t, "error2", reqs[2].Error.Error())
		require.Equal(t, "error3", reqs[3].Error.Error())
		require.Equal(t, "error4", reqs[4].Error.Error())
		require.Equal(t, "error5", reqs[5].Error.Error())
		require.Equal(t, "error6", reqs[6].Error.Error())
	})
}
