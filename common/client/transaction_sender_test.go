package client

import (
	"context"
	"fmt"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"
	"github.com/smartcontractkit/chainlink/v2/common/types"
)

type sendTxMultiNode struct {
	*MultiNode[types.ID, SendTxRPCClient[any]]
}

type sendTxRPC struct {
	sendTxRun func(args mock.Arguments)
	sendTxErr error
}

var _ SendTxRPCClient[any] = (*sendTxRPC)(nil)

func newSendTxRPC(sendTxErr error, sendTxRun func(args mock.Arguments)) *sendTxRPC {
	return &sendTxRPC{sendTxErr: sendTxErr, sendTxRun: sendTxRun}
}

func (rpc *sendTxRPC) SendTransaction(ctx context.Context, _ any) error {
	if rpc.sendTxRun != nil {
		rpc.sendTxRun(mock.Arguments{ctx})
	}
	return rpc.sendTxErr
}

func newTestTransactionSender(t *testing.T, chainID types.ID, lggr logger.Logger,
	nodes []Node[types.ID, SendTxRPCClient[any]],
	sendOnlyNodes []SendOnlyNode[types.ID, SendTxRPCClient[any]],
) (*sendTxMultiNode, *TransactionSender[any, types.ID, SendTxRPCClient[any]]) {
	mn := sendTxMultiNode{NewMultiNode[types.ID, SendTxRPCClient[any]](
		lggr, NodeSelectionModeRoundRobin, 0, nodes, sendOnlyNodes, chainID, "chainFamily", 0)}
	err := mn.StartOnce("startedTestMultiNode", func() error { return nil })
	require.NoError(t, err)

	txSender := NewTransactionSender[any, types.ID, SendTxRPCClient[any]](lggr, chainID, mn.chainFamily, mn.MultiNode, classifySendTxError, tests.TestInterval)
	err = txSender.Start(tests.Context(t))
	require.NoError(t, err)

	t.Cleanup(func() {
		err := mn.Close()
		if err != nil {
			// Allow MultiNode to be closed early for testing
			require.EqualError(t, err, "MultiNode has already been stopped: already stopped")
		}
		err = txSender.Close()
		if err != nil {
			// Allow TransactionSender to be closed early for testing
			require.EqualError(t, err, "TransactionSender has already been stopped: already stopped")
		}
	})
	return &mn, txSender
}

func classifySendTxError(_ any, err error) SendTxReturnCode {
	if err != nil {
		return Fatal
	}
	return Successful
}

func TestTransactionSender_SendTransaction(t *testing.T) {
	t.Parallel()

	newNodeWithState := func(t *testing.T, state nodeState, txErr error, sendTxRun func(args mock.Arguments)) *mockNode[types.ID, SendTxRPCClient[any]] {
		rpc := newSendTxRPC(txErr, sendTxRun)
		node := newMockNode[types.ID, SendTxRPCClient[any]](t)
		node.On("String").Return("node name").Maybe()
		node.On("RPC").Return(rpc).Maybe()
		node.On("State").Return(state).Maybe()
		node.On("Close").Return(nil).Once()
		return node
	}

	newNode := func(t *testing.T, txErr error, sendTxRun func(args mock.Arguments)) *mockNode[types.ID, SendTxRPCClient[any]] {
		return newNodeWithState(t, nodeStateAlive, txErr, sendTxRun)
	}

	t.Run("Fails if there is no nodes available", func(t *testing.T) {
		lggr, _ := logger.TestObserved(t, zap.DebugLevel)
		_, txSender := newTestTransactionSender(t, types.RandomID(), lggr, nil, nil)
		_, err := txSender.SendTransaction(tests.Context(t), nil)
		assert.EqualError(t, err, ErroringNodeError.Error())
	})

	t.Run("Transaction failure happy path", func(t *testing.T) {
		expectedError := errors.New("transaction failed")
		mainNode := newNode(t, expectedError, nil)
		lggr, observedLogs := logger.TestObserved(t, zap.DebugLevel)

		_, txSender := newTestTransactionSender(t, types.RandomID(), lggr,
			[]Node[types.ID, SendTxRPCClient[any]]{mainNode},
			[]SendOnlyNode[types.ID, SendTxRPCClient[any]]{newNode(t, errors.New("unexpected error"), nil)})

		result, sendErr := txSender.SendTransaction(tests.Context(t), nil)
		require.ErrorIs(t, sendErr, expectedError)
		require.Equal(t, Fatal, result)
		tests.AssertLogCountEventually(t, observedLogs, "Node sent transaction", 2)
		tests.AssertLogCountEventually(t, observedLogs, "RPC returned error", 2)
	})

	t.Run("Transaction success happy path", func(t *testing.T) {
		mainNode := newNode(t, nil, nil)

		lggr, observedLogs := logger.TestObserved(t, zap.DebugLevel)
		_, txSender := newTestTransactionSender(t, types.RandomID(), lggr,
			[]Node[types.ID, SendTxRPCClient[any]]{mainNode},
			[]SendOnlyNode[types.ID, SendTxRPCClient[any]]{newNode(t, errors.New("unexpected error"), nil)})

		result, sendErr := txSender.SendTransaction(tests.Context(t), nil)
		require.NoError(t, sendErr)
		require.Equal(t, Successful, result)
		tests.AssertLogCountEventually(t, observedLogs, "Node sent transaction", 2)
		tests.AssertLogCountEventually(t, observedLogs, "RPC returned error", 1)
	})

	t.Run("Context expired before collecting sufficient results", func(t *testing.T) {
		testContext, testCancel := context.WithCancel(tests.Context(t))
		defer testCancel()

		mainNode := newNode(t, nil, func(_ mock.Arguments) {
			// block caller til end of the test
			<-testContext.Done()
		})

		lggr, _ := logger.TestObserved(t, zap.DebugLevel)

		_, txSender := newTestTransactionSender(t, types.RandomID(), lggr,
			[]Node[types.ID, SendTxRPCClient[any]]{mainNode}, nil)

		requestContext, cancel := context.WithCancel(tests.Context(t))
		cancel()
		_, sendErr := txSender.SendTransaction(requestContext, nil)
		require.EqualError(t, sendErr, "context canceled")
	})

	t.Run("Soft timeout stops results collection", func(t *testing.T) {
		chainID := types.RandomID()
		expectedError := errors.New("transaction failed")
		fastNode := newNode(t, expectedError, nil)

		// hold reply from the node till end of the test
		testContext, testCancel := context.WithCancel(tests.Context(t))
		defer testCancel()
		slowNode := newNode(t, errors.New("transaction failed"), func(_ mock.Arguments) {
			// block caller til end of the test
			<-testContext.Done()
		})

		lggr, _ := logger.TestObserved(t, zap.DebugLevel)

		_, txSender := newTestTransactionSender(t, chainID, lggr, []Node[types.ID, SendTxRPCClient[any]]{fastNode, slowNode}, nil)
		_, sendErr := txSender.SendTransaction(tests.Context(t), nil)
		require.EqualError(t, sendErr, expectedError.Error())
	})
	t.Run("Returns success without waiting for the rest of the nodes", func(t *testing.T) {
		chainID := types.RandomID()
		fastNode := newNode(t, nil, nil)
		// hold reply from the node till end of the test
		testContext, testCancel := context.WithCancel(tests.Context(t))
		defer testCancel()
		slowNode := newNode(t, errors.New("transaction failed"), func(_ mock.Arguments) {
			// block caller til end of the test
			<-testContext.Done()
		})
		slowSendOnly := newNode(t, errors.New("send only failed"), func(_ mock.Arguments) {
			// block caller til end of the test
			<-testContext.Done()
		})
		lggr, _ := logger.TestObserved(t, zap.WarnLevel)
		mn, txSender := newTestTransactionSender(t, chainID, lggr,
			[]Node[types.ID, SendTxRPCClient[any]]{fastNode, slowNode},
			[]SendOnlyNode[types.ID, SendTxRPCClient[any]]{slowSendOnly})

		rtnCode, err := txSender.SendTransaction(tests.Context(t), nil)
		require.NoError(t, err)
		require.Equal(t, Successful, rtnCode)
		require.NoError(t, mn.Close())
	})
	t.Run("Fails when multinode is closed", func(t *testing.T) {
		chainID := types.RandomID()
		fastNode := newNode(t, nil, nil)
		// hold reply from the node till end of the test
		testContext, testCancel := context.WithCancel(tests.Context(t))
		defer testCancel()
		slowNode := newNode(t, errors.New("transaction failed"), func(_ mock.Arguments) {
			// block caller til end of the test
			<-testContext.Done()
		})
		slowSendOnly := newNode(t, errors.New("send only failed"), func(_ mock.Arguments) {
			// block caller til end of the test
			<-testContext.Done()
		})

		lggr, _ := logger.TestObserved(t, zap.DebugLevel)

		mn, txSender := newTestTransactionSender(t, chainID, lggr,
			[]Node[types.ID, SendTxRPCClient[any]]{fastNode, slowNode},
			[]SendOnlyNode[types.ID, SendTxRPCClient[any]]{slowSendOnly})

		require.NoError(t, mn.Close())
		_, err := txSender.SendTransaction(tests.Context(t), nil)
		require.EqualError(t, err, "MultiNode is stopped")
	})
	t.Run("Fails when closed", func(t *testing.T) {
		chainID := types.RandomID()
		fastNode := newNode(t, nil, nil)
		// hold reply from the node till end of the test
		testContext, testCancel := context.WithCancel(tests.Context(t))
		defer testCancel()
		slowNode := newNode(t, errors.New("transaction failed"), func(_ mock.Arguments) {
			// block caller til end of the test
			<-testContext.Done()
		})
		slowSendOnly := newNode(t, errors.New("send only failed"), func(_ mock.Arguments) {
			// block caller til end of the test
			<-testContext.Done()
		})

		lggr, _ := logger.TestObserved(t, zap.DebugLevel)

		_, txSender := newTestTransactionSender(t, chainID, lggr,
			[]Node[types.ID, SendTxRPCClient[any]]{fastNode, slowNode},
			[]SendOnlyNode[types.ID, SendTxRPCClient[any]]{slowSendOnly})

		require.NoError(t, txSender.Close())
		_, err := txSender.SendTransaction(tests.Context(t), nil)
		require.EqualError(t, err, "TransactionSender not started")
	})
	t.Run("Returns error if there is no healthy primary nodes", func(t *testing.T) {
		chainID := types.RandomID()
		primary := newNodeWithState(t, nodeStateUnreachable, nil, nil)
		sendOnly := newNodeWithState(t, nodeStateUnreachable, nil, nil)

		lggr, _ := logger.TestObserved(t, zap.DebugLevel)

		_, txSender := newTestTransactionSender(t, chainID, lggr,
			[]Node[types.ID, SendTxRPCClient[any]]{primary},
			[]SendOnlyNode[types.ID, SendTxRPCClient[any]]{sendOnly})

		_, sendErr := txSender.SendTransaction(tests.Context(t), nil)
		assert.EqualError(t, sendErr, ErroringNodeError.Error())
	})

	t.Run("Transaction success even if one of the nodes is unhealthy", func(t *testing.T) {
		chainID := types.RandomID()
		mainNode := newNode(t, nil, nil)
		unexpectedCall := func(args mock.Arguments) {
			panic("SendTx must not be called for unhealthy node")
		}
		unhealthyNode := newNodeWithState(t, nodeStateUnreachable, nil, unexpectedCall)
		unhealthySendOnlyNode := newNodeWithState(t, nodeStateUnreachable, nil, unexpectedCall)

		lggr, _ := logger.TestObserved(t, zap.DebugLevel)

		_, txSender := newTestTransactionSender(t, chainID, lggr,
			[]Node[types.ID, SendTxRPCClient[any]]{mainNode, unhealthyNode},
			[]SendOnlyNode[types.ID, SendTxRPCClient[any]]{unhealthySendOnlyNode})

		returnCode, sendErr := txSender.SendTransaction(tests.Context(t), nil)
		require.NoError(t, sendErr)
		require.Equal(t, Successful, returnCode)
	})
}

func TestTransactionSender_SendTransaction_aggregateTxResults(t *testing.T) {
	t.Parallel()
	// ensure failure on new SendTxReturnCode
	codesToCover := map[SendTxReturnCode]struct{}{}
	for code := Successful; code < sendTxReturnCodeLen; code++ {
		codesToCover[code] = struct{}{}
	}

	testCases := []struct {
		Name                string
		ExpectedTxResult    string
		ExpectedCriticalErr string
		ResultsByCode       sendTxResults
	}{
		{
			Name:                "Returns success and logs critical error on success and Fatal",
			ExpectedTxResult:    "success",
			ExpectedCriticalErr: "found contradictions in nodes replies on SendTransaction: got success and severe error",
			ResultsByCode: sendTxResults{
				Successful: {errors.New("success")},
				Fatal:      {errors.New("fatal")},
			},
		},
		{
			Name:                "Returns TransactionAlreadyKnown and logs critical error on TransactionAlreadyKnown and Fatal",
			ExpectedTxResult:    "tx_already_known",
			ExpectedCriticalErr: "found contradictions in nodes replies on SendTransaction: got success and severe error",
			ResultsByCode: sendTxResults{
				TransactionAlreadyKnown: {errors.New("tx_already_known")},
				Unsupported:             {errors.New("unsupported")},
			},
		},
		{
			Name:                "Prefers sever error to temporary",
			ExpectedTxResult:    "underpriced",
			ExpectedCriticalErr: "",
			ResultsByCode: sendTxResults{
				Retryable:   {errors.New("retryable")},
				Underpriced: {errors.New("underpriced")},
			},
		},
		{
			Name:                "Returns temporary error",
			ExpectedTxResult:    "retryable",
			ExpectedCriticalErr: "",
			ResultsByCode: sendTxResults{
				Retryable: {errors.New("retryable")},
			},
		},
		{
			Name:                "Insufficient funds is treated as  error",
			ExpectedTxResult:    "",
			ExpectedCriticalErr: "",
			ResultsByCode: sendTxResults{
				Successful:        {nil},
				InsufficientFunds: {errors.New("insufficientFunds")},
			},
		},
		{
			Name:                "Logs critical error on empty ResultsByCode",
			ExpectedTxResult:    "expected at least one response on SendTransaction",
			ExpectedCriticalErr: "expected at least one response on SendTransaction",
			ResultsByCode:       sendTxResults{},
		},
		{
			Name:                "Zk terminally stuck",
			ExpectedTxResult:    "not enough keccak counters to continue the execution",
			ExpectedCriticalErr: "",
			ResultsByCode: sendTxResults{
				TerminallyStuck: {errors.New("not enough keccak counters to continue the execution")},
			},
		},
	}

	for _, testCase := range testCases {
		for code := range testCase.ResultsByCode {
			delete(codesToCover, code)
		}

		t.Run(testCase.Name, func(t *testing.T) {
			_, txResult, err := aggregateTxResults(testCase.ResultsByCode)
			if testCase.ExpectedTxResult == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, txResult, testCase.ExpectedTxResult)
			}

			logger.Sugared(logger.Test(t)).Info("Map: " + fmt.Sprint(testCase.ResultsByCode))
			logger.Sugared(logger.Test(t)).Criticalw("observed invariant violation on SendTransaction", "resultsByCode", testCase.ResultsByCode, "err", err)

			if testCase.ExpectedCriticalErr == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, testCase.ExpectedCriticalErr)
			}
		})
	}

	// explicitly signal that following codes are properly handled in aggregateTxResults,
	// but dedicated test cases won't be beneficial
	for _, codeToIgnore := range []SendTxReturnCode{Unknown, ExceedsMaxFee, FeeOutOfValidRange} {
		delete(codesToCover, codeToIgnore)
	}
	assert.Empty(t, codesToCover, "all of the SendTxReturnCode must be covered by this test")
}
