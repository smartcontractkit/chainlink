package client

import (
	"context"
	"fmt"
	"time"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"

	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"
	"github.com/smartcontractkit/chainlink/v2/common/types"

	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
)

type sendTxMultiNode struct {
	*MultiNode[types.ID, SendTxRPCClient[any]]
}

type sendTxMultiNodeOpts struct {
	logger        logger.Logger
	selectionMode string
	leaseDuration time.Duration
	nodes         []Node[types.ID, SendTxRPCClient[any]]
	sendonlys     []SendOnlyNode[types.ID, SendTxRPCClient[any]]
	chainID       types.ID
	chainFamily   string
}

func newSendTxMultiNode(t *testing.T, opts sendTxMultiNodeOpts) sendTxMultiNode {
	if opts.logger == nil {
		opts.logger = logger.Test(t)
	}

	result := NewMultiNode[types.ID, SendTxRPCClient[any]](
		opts.logger, opts.selectionMode, opts.leaseDuration, opts.nodes, opts.sendonlys, opts.chainID, opts.chainFamily)
	return sendTxMultiNode{
		result,
	}
}

type sendTxRPC struct {
	sendTxErr error
}

var _ SendTxRPCClient[any] = (*sendTxRPC)(nil)

func newSendTxRPC(sendTxErr error) *sendTxRPC {
	return &sendTxRPC{sendTxErr: sendTxErr}
}

func (rpc *sendTxRPC) SendTransaction(ctx context.Context, tx any) error {
	return rpc.sendTxErr
}

func TestMultiNode_SendTransaction(t *testing.T) {
	t.Parallel()
	classifySendTxError := func(tx any, err error) SendTxReturnCode {
		if err != nil {
			return Fatal
		}
		return Successful
	}

	newNodeWithState := func(t *testing.T, state NodeState, returnCode SendTxReturnCode, txErr error, sendTxRun func(args mock.Arguments)) *mockNode[types.ID, SendTxRPCClient[any]] {
		rpc := newSendTxRPC(txErr)
		node := newMockNode[types.ID, SendTxRPCClient[any]](t)
		node.On("String").Return("node name").Maybe()
		node.On("RPC").Return(rpc).Maybe()
		node.On("State").Return(state).Maybe()
		node.On("Close").Return(nil).Once()
		return node
	}

	newNode := func(t *testing.T, returnCode SendTxReturnCode, txErr error, sendTxRun func(args mock.Arguments)) *mockNode[types.ID, SendTxRPCClient[any]] {
		return newNodeWithState(t, NodeStateAlive, returnCode, txErr, sendTxRun)
	}

	newStartedMultiNode := func(t *testing.T, opts sendTxMultiNodeOpts) sendTxMultiNode {
		mn := newSendTxMultiNode(t, opts)
		err := mn.StartOnce("startedTestMultiNode", func() error { return nil })
		require.NoError(t, err)
		t.Cleanup(func() {
			require.NoError(t, mn.Close())
		})
		return mn
	}

	sendTxSoftTimeout := tests.TestInterval

	t.Run("Fails if there is no nodes available", func(t *testing.T) {
		mn := newStartedMultiNode(t, sendTxMultiNodeOpts{
			selectionMode: NodeSelectionModeRoundRobin,
			chainID:       types.RandomID(),
		})

		lggr, _ := logger.TestObserved(t, zap.DebugLevel)

		txSender := NewTransactionSender[any, types.ID, SendTxRPCClient[any]](lggr, mn.chainID, mn.chainFamily, mn.MultiNode, classifySendTxError, sendTxSoftTimeout)
		_, err := txSender.SendTransaction(tests.Context(t), nil)
		assert.EqualError(t, err, "no calls were completed")
	})

	t.Run("Transaction failure happy path", func(t *testing.T) {
		chainID := types.RandomID()
		expectedError := errors.New("transaction failed")
		mainNode := newNode(t, 0, expectedError, nil)

		lggr, observedLogs := logger.TestObserved(t, zap.DebugLevel)
		mn := newStartedMultiNode(t, sendTxMultiNodeOpts{
			selectionMode: NodeSelectionModeRoundRobin,
			chainID:       chainID,
			nodes:         []Node[types.ID, SendTxRPCClient[any]]{mainNode},
			sendonlys:     []SendOnlyNode[types.ID, SendTxRPCClient[any]]{newNode(t, 0, errors.New("unexpected error"), nil)},
			logger:        lggr,
		})

		txSender := NewTransactionSender[any, types.ID, SendTxRPCClient[any]](lggr, chainID, mn.chainFamily, mn.MultiNode, classifySendTxError, sendTxSoftTimeout)

		result, sendErr := txSender.SendTransaction(tests.Context(t), nil)
		require.ErrorIs(t, sendErr, expectedError)
		require.Equal(t, Fatal, result)
		tests.AssertLogCountEventually(t, observedLogs, "Node sent transaction", 2)
		tests.AssertLogCountEventually(t, observedLogs, "RPC returned error", 2)
	})

	t.Run("Transaction success happy path", func(t *testing.T) {
		chainID := types.RandomID()
		mainNode := newNode(t, 0, nil, nil)

		lggr, observedLogs := logger.TestObserved(t, zap.DebugLevel)
		mn := newStartedMultiNode(t, sendTxMultiNodeOpts{
			selectionMode: NodeSelectionModeRoundRobin,
			chainID:       chainID,
			nodes:         []Node[types.ID, SendTxRPCClient[any]]{mainNode},
			sendonlys:     []SendOnlyNode[types.ID, SendTxRPCClient[any]]{newNode(t, 0, errors.New("unexpected error"), nil)},
			logger:        lggr,
		})

		txSender := NewTransactionSender[any, types.ID, SendTxRPCClient[any]](lggr, chainID, mn.chainFamily, mn.MultiNode, classifySendTxError, sendTxSoftTimeout)

		result, sendErr := txSender.SendTransaction(tests.Context(t), nil)
		require.NoError(t, sendErr)
		require.Equal(t, Successful, result)
		tests.AssertLogCountEventually(t, observedLogs, "Node sent transaction", 2)
		tests.AssertLogCountEventually(t, observedLogs, "RPC returned error", 1)
	})

	t.Run("Context expired before collecting sufficient results", func(t *testing.T) {
		chainID := types.RandomID()
		testContext, testCancel := context.WithCancel(tests.Context(t))
		defer testCancel()

		mainNode := newNode(t, 0, nil, func(_ mock.Arguments) {
			// block caller til end of the test
			<-testContext.Done()
		})

		lggr, _ := logger.TestObserved(t, zap.DebugLevel)
		mn := newStartedMultiNode(t, sendTxMultiNodeOpts{
			selectionMode: NodeSelectionModeRoundRobin,
			chainID:       chainID,
			nodes:         []Node[types.ID, SendTxRPCClient[any]]{mainNode},
			logger:        lggr,
		})

		txSender := NewTransactionSender[any, types.ID, SendTxRPCClient[any]](lggr, chainID, mn.chainFamily, mn.MultiNode, classifySendTxError, sendTxSoftTimeout)

		requestContext, cancel := context.WithCancel(tests.Context(t))
		cancel()
		_, sendErr := txSender.SendTransaction(requestContext, nil)
		require.EqualError(t, sendErr, "context canceled")
	})

	t.Run("Soft timeout stops results collection", func(t *testing.T) {
		chainID := types.RandomID()
		expectedError := errors.New("tmp error")
		fastNode := newNode(t, 0, expectedError, nil)

		// hold reply from the node till end of the test
		testContext, testCancel := context.WithCancel(tests.Context(t))
		defer testCancel()
		slowNode := newNode(t, 0, errors.New("transaction failed"), func(_ mock.Arguments) {
			// block caller til end of the test
			<-testContext.Done()
		})

		lggr, _ := logger.TestObserved(t, zap.DebugLevel)
		mn := newStartedMultiNode(t, sendTxMultiNodeOpts{
			selectionMode: NodeSelectionModeRoundRobin,
			chainID:       chainID,
			nodes:         []Node[types.ID, SendTxRPCClient[any]]{fastNode, slowNode},
			logger:        lggr,
		})

		txSender := NewTransactionSender[any, types.ID, SendTxRPCClient[any]](lggr, chainID, mn.chainFamily, mn.MultiNode, classifySendTxError, sendTxSoftTimeout)
		_, sendErr := txSender.SendTransaction(tests.Context(t), nil)
		require.EqualError(t, sendErr, expectedError.Error())
	})
	t.Run("Returns success without waiting for the rest of the nodes", func(t *testing.T) {
		chainID := types.RandomID()
		fastNode := newNode(t, 0, nil, nil)
		// hold reply from the node till end of the test
		testContext, testCancel := context.WithCancel(tests.Context(t))
		defer testCancel()
		slowNode := newNode(t, 0, errors.New("transaction failed"), func(_ mock.Arguments) {
			// block caller til end of the test
			<-testContext.Done()
		})
		slowSendOnly := newNode(t, 0, errors.New("send only failed"), func(_ mock.Arguments) {
			// block caller til end of the test
			<-testContext.Done()
		})
		lggr, observedLogs := logger.TestObserved(t, zap.WarnLevel)

		mn := newSendTxMultiNode(t, sendTxMultiNodeOpts{
			selectionMode: NodeSelectionModeRoundRobin,
			chainID:       chainID,
			nodes:         []Node[types.ID, SendTxRPCClient[any]]{fastNode, slowNode},
			sendonlys:     []SendOnlyNode[types.ID, SendTxRPCClient[any]]{slowSendOnly},
			logger:        lggr,
		})

		txSender := NewTransactionSender[any, types.ID, SendTxRPCClient[any]](lggr, chainID, mn.chainFamily, mn.MultiNode, classifySendTxError, sendTxSoftTimeout)
		assert.NoError(t, mn.StartOnce("startedTestMultiNode", func() error { return nil }))
		_, err := txSender.SendTransaction(tests.Context(t), nil)
		require.NoError(t, err)
		testCancel()
		require.NoError(t, mn.Close())
		tests.AssertLogEventually(t, observedLogs, "observed invariant violation on SendTransaction")
	})
	t.Run("Fails when closed", func(t *testing.T) {
		chainID := types.RandomID()
		fastNode := newNode(t, 0, nil, nil)
		// hold reply from the node till end of the test
		testContext, testCancel := context.WithCancel(tests.Context(t))
		defer testCancel()
		slowNode := newNode(t, 0, errors.New("transaction failed"), func(_ mock.Arguments) {
			// block caller til end of the test
			<-testContext.Done()
		})
		slowSendOnly := newNode(t, 0, errors.New("send only failed"), func(_ mock.Arguments) {
			// block caller til end of the test
			<-testContext.Done()
		})

		lggr, _ := logger.TestObserved(t, zap.DebugLevel)

		mn := newSendTxMultiNode(t, sendTxMultiNodeOpts{
			selectionMode: NodeSelectionModeRoundRobin,
			chainID:       chainID,
			nodes:         []Node[types.ID, SendTxRPCClient[any]]{fastNode, slowNode},
			sendonlys:     []SendOnlyNode[types.ID, SendTxRPCClient[any]]{slowSendOnly},
			logger:        lggr,
		})

		txSender := NewTransactionSender[any, types.ID, SendTxRPCClient[any]](lggr, chainID, mn.chainFamily, mn.MultiNode, classifySendTxError, sendTxSoftTimeout)
		require.NoError(t, mn.StartOnce("startedTestMultiNode", func() error { return nil }))
		require.NoError(t, mn.Close())
		_, err := txSender.SendTransaction(tests.Context(t), nil)
		require.EqualError(t, err, "aborted while broadcasting tx - MultiNode is stopped: context canceled")
	})
	t.Run("Returns error if there is no healthy primary nodes", func(t *testing.T) {
		chainID := types.RandomID()
		primary := newNodeWithState(t, NodeStateUnreachable, 0, nil, nil)
		sendOnly := newNodeWithState(t, NodeStateUnreachable, 0, nil, nil)

		lggr, _ := logger.TestObserved(t, zap.DebugLevel)

		mn := newStartedMultiNode(t, sendTxMultiNodeOpts{
			selectionMode: NodeSelectionModeRoundRobin,
			chainID:       chainID,
			nodes:         []Node[types.ID, SendTxRPCClient[any]]{primary},
			sendonlys:     []SendOnlyNode[types.ID, SendTxRPCClient[any]]{sendOnly},
			logger:        lggr,
		})

		txSender := NewTransactionSender[any, types.ID, SendTxRPCClient[any]](lggr, chainID, mn.chainFamily, mn.MultiNode, classifySendTxError, sendTxSoftTimeout)
		_, sendErr := txSender.SendTransaction(tests.Context(t), nil)
		assert.EqualError(t, sendErr, "no calls were completed")
	})

	// TODO: Get this test to pass
	t.Run("Transaction success even if one of the nodes is unhealthy", func(t *testing.T) {
		chainID := types.RandomID()
		mainNode := newNode(t, Successful, nil, nil)
		unexpectedCall := func(args mock.Arguments) {
			panic("SendTx must not be called for unhealthy node")
		}
		unhealthyNode := newNodeWithState(t, NodeStateUnreachable, 0, nil, unexpectedCall)
		unhealthySendOnlyNode := newNodeWithState(t, NodeStateUnreachable, 0, nil, unexpectedCall)

		lggr, observedLogs := logger.TestObserved(t, zap.DebugLevel)
		mn := newStartedMultiNode(t, sendTxMultiNodeOpts{
			selectionMode: NodeSelectionModeRoundRobin,
			chainID:       chainID,
			nodes:         []Node[types.ID, SendTxRPCClient[any]]{mainNode, unhealthyNode},
			sendonlys:     []SendOnlyNode[types.ID, SendTxRPCClient[any]]{unhealthySendOnlyNode},
			logger:        lggr,
		})

		txSender := NewTransactionSender[any, types.ID, SendTxRPCClient[any]](lggr, chainID, mn.chainFamily, mn.MultiNode, classifySendTxError, sendTxSoftTimeout)
		returnCode, sendErr := txSender.SendTransaction(tests.Context(t), nil)
		require.NoError(t, sendErr)
		require.Equal(t, Successful, returnCode)
		tests.AssertLogCountEventually(t, observedLogs, "Node sent transaction", 2)
		tests.AssertLogCountEventually(t, observedLogs, "RPC returned error", 1)
	})
}

func TestMultiNode_SendTransaction_aggregateTxResults(t *testing.T) {
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
		ResultsByCode       sendTxErrors
	}{
		{
			Name:                "Returns success and logs critical error on success and Fatal",
			ExpectedTxResult:    "success",
			ExpectedCriticalErr: "found contradictions in nodes replies on SendTransaction: got success and severe error",
			ResultsByCode: sendTxErrors{
				Successful: {errors.New("success")},
				Fatal:      {errors.New("fatal")},
			},
		},
		{
			Name:                "Returns TransactionAlreadyKnown and logs critical error on TransactionAlreadyKnown and Fatal",
			ExpectedTxResult:    "tx_already_known",
			ExpectedCriticalErr: "found contradictions in nodes replies on SendTransaction: got success and severe error",
			ResultsByCode: sendTxErrors{
				TransactionAlreadyKnown: {errors.New("tx_already_known")},
				Unsupported:             {errors.New("unsupported")},
			},
		},
		{
			Name:                "Prefers sever error to temporary",
			ExpectedTxResult:    "underpriced",
			ExpectedCriticalErr: "",
			ResultsByCode: sendTxErrors{
				Retryable:   {errors.New("retryable")},
				Underpriced: {errors.New("underpriced")},
			},
		},
		{
			Name:                "Returns temporary error",
			ExpectedTxResult:    "retryable",
			ExpectedCriticalErr: "",
			ResultsByCode: sendTxErrors{
				Retryable: {errors.New("retryable")},
			},
		},
		{
			Name:                "Insufficient funds is treated as  error",
			ExpectedTxResult:    "",
			ExpectedCriticalErr: "",
			ResultsByCode: sendTxErrors{
				Successful:        {nil},
				InsufficientFunds: {errors.New("insufficientFunds")},
			},
		},
		{
			Name:                "Logs critical error on empty ResultsByCode",
			ExpectedTxResult:    "expected at least one response on SendTransaction",
			ExpectedCriticalErr: "expected at least one response on SendTransaction",
			ResultsByCode:       sendTxErrors{},
		},
		{
			Name:                "Zk out of counter error",
			ExpectedTxResult:    "not enough keccak counters to continue the execution",
			ExpectedCriticalErr: "",
			ResultsByCode: sendTxErrors{
				OutOfCounters: {errors.New("not enough keccak counters to continue the execution")},
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
	//but dedicated test cases won't be beneficial
	for _, codeToIgnore := range []SendTxReturnCode{Unknown, ExceedsMaxFee, FeeOutOfValidRange} {
		delete(codesToCover, codeToIgnore)
	}
	assert.Empty(t, codesToCover, "all of the SendTxReturnCode must be covered by this test")
}
