package client_test

import (
	"errors"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/rpc"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	commonclient "github.com/smartcontractkit/chainlink/v2/common/client"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
)

func TestSimulateTx(t *testing.T) {
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

		msg := client.TxSimulationRequest{
			From: fromAddress,
			To:   &toAddress,
			Data: []byte("0x00"),
		}
		sendErr := client.SimulateTransaction(ctx, ethClient, "", msg)
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

		msg := client.TxSimulationRequest{
			From: fromAddress,
			To:   &toAddress,
			Data: []byte("0x00"),
		}
		sendErr := client.SimulateTransaction(ctx, ethClient, "", msg)
		require.Equal(t, true, sendErr.IsOutOfCounters())
	})

	t.Run("returns without error if simulation returns non-OOC error", func(t *testing.T) {
		mockRpc := newMockRpc(t)
		mockRpc.On("BatchCallContext", mock.Anything, mock.IsType([]rpc.BatchElem{})).Run(func(args mock.Arguments) {
			reqs := args.Get(1).([]rpc.BatchElem)
			for i := 0; i < len(reqs); i++ {
				elem := &reqs[i]
				elem.Error = errors.New("something went wrong")
			}
		}).Return(nil).Once()

		ethClient := client.NewChainClientWithMockedRpc(t, commonclient.NodeSelectionModeRoundRobin, time.Second*0, time.Second*0, testutils.FixtureChainID, mockRpc)
		err := ethClient.Dial(ctx)
		require.NoError(t, err)

		msg := client.TxSimulationRequest{
			From: fromAddress,
			To:   &toAddress,
			Data: []byte("0x00"),
		}
		sendErr := client.SimulateTransaction(ctx, ethClient, "", msg)
		require.Equal(t, false, sendErr.IsOutOfCounters())
	})
}

func TestSimulateTx_Batch(t *testing.T) {
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
		err = client.BatchSimulateTransaction(ctx, ethClient, "", reqs)
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
		err = client.BatchSimulateTransaction(ctx, ethClient, "", reqs)
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
		err = client.BatchSimulateTransaction(ctx, ethClient, "", reqs)
		require.NoError(t, err)

		// No error for first request
		require.Empty(t, reqs[0].Error)

		// No Out-of-counter error for second request
		require.Equal(t, errorMsg, reqs[1].Error.Error())
	})

	t.Run("returns the proper error to the associated request", func(t *testing.T) {
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
		err = client.BatchSimulateTransaction(ctx, ethClient, "", reqs)
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
