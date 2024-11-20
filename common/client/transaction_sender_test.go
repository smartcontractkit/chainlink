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
	"github.com/smartcontractkit/chainlink-common/pkg/services/servicetest"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"
	"github.com/smartcontractkit/chainlink/v2/common/types"
)

type TestSendTxRPCClient SendTxRPCClient[any, *sendTxResult]

type sendTxMultiNode struct {
	*MultiNode[types.ID, TestSendTxRPCClient]
}

type sendTxRPC struct {
	sendTxRun func(args mock.Arguments)
	sendTxErr error
}

type sendTxResult struct {
	err  error
	code SendTxReturnCode
}

var _ SendTxResult = (*sendTxResult)(nil)

func NewSendTxResult(err error) *sendTxResult {
	result := &sendTxResult{
		err: err,
	}
	return result
}

func (r *sendTxResult) Error() error {
	return r.err
}

func (r *sendTxResult) Code() SendTxReturnCode {
	return r.code
}

var _ TestSendTxRPCClient = (*sendTxRPC)(nil)

func newSendTxRPC(sendTxErr error, sendTxRun func(args mock.Arguments)) *sendTxRPC {
	return &sendTxRPC{sendTxErr: sendTxErr, sendTxRun: sendTxRun}
}

func (rpc *sendTxRPC) SendTransaction(ctx context.Context, _ any) *sendTxResult {
	if rpc.sendTxRun != nil {
		rpc.sendTxRun(mock.Arguments{ctx})
	}
	return &sendTxResult{err: rpc.sendTxErr, code: classifySendTxError(nil, rpc.sendTxErr)}
}

// newTestTransactionSender returns a sendTxMultiNode and TransactionSender.
// Only the TransactionSender is run via Start/Close.
func newTestTransactionSender(t *testing.T, chainID types.ID, lggr logger.Logger,
	nodes []Node[types.ID, TestSendTxRPCClient],
	sendOnlyNodes []SendOnlyNode[types.ID, TestSendTxRPCClient],
) (*sendTxMultiNode, *TransactionSender[any, *sendTxResult, types.ID, TestSendTxRPCClient]) {
	mn := sendTxMultiNode{NewMultiNode[types.ID, TestSendTxRPCClient](
		lggr, NodeSelectionModeRoundRobin, 0, nodes, sendOnlyNodes, chainID, "chainFamily", 0)}

	txSender := NewTransactionSender[any, *sendTxResult, types.ID, TestSendTxRPCClient](lggr, chainID, mn.chainFamily, mn.MultiNode, NewSendTxResult, tests.TestInterval)
	servicetest.Run(t, txSender)
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

	newNodeWithState := func(t *testing.T, state nodeState, txErr error, sendTxRun func(args mock.Arguments)) *mockNode[types.ID, TestSendTxRPCClient] {
		rpc := newSendTxRPC(txErr, sendTxRun)
		node := newMockNode[types.ID, TestSendTxRPCClient](t)
		node.On("String").Return("node name").Maybe()
		node.On("RPC").Return(rpc).Maybe()
		node.On("State").Return(state).Maybe()
		node.On("Start", mock.Anything).Return(nil).Maybe()
		node.On("Close").Return(nil).Maybe()
		node.On("SetPoolChainInfoProvider", mock.Anything).Return(nil).Maybe()
		return node
	}

	newNode := func(t *testing.T, txErr error, sendTxRun func(args mock.Arguments)) *mockNode[types.ID, TestSendTxRPCClient] {
		return newNodeWithState(t, nodeStateAlive, txErr, sendTxRun)
	}

	t.Run("Fails if there is no nodes available", func(t *testing.T) {
		lggr := logger.Test(t)
		_, txSender := newTestTransactionSender(t, types.RandomID(), lggr, nil, nil)
		result := txSender.SendTransaction(tests.Context(t), nil)
		assert.EqualError(t, result.Error(), ErroringNodeError.Error())
	})

	t.Run("Transaction failure happy path", func(t *testing.T) {
		expectedError := errors.New("transaction failed")
		mainNode := newNode(t, expectedError, nil)
		lggr, observedLogs := logger.TestObserved(t, zap.DebugLevel)

		_, txSender := newTestTransactionSender(t, types.RandomID(), lggr,
			[]Node[types.ID, TestSendTxRPCClient]{mainNode},
			[]SendOnlyNode[types.ID, TestSendTxRPCClient]{newNode(t, errors.New("unexpected error"), nil)})

		result := txSender.SendTransaction(tests.Context(t), nil)
		require.ErrorIs(t, result.Error(), expectedError)
		require.Equal(t, Fatal, result.Code())
		tests.AssertLogCountEventually(t, observedLogs, "Node sent transaction", 2)
		tests.AssertLogCountEventually(t, observedLogs, "RPC returned error", 2)
	})

	t.Run("Transaction success happy path", func(t *testing.T) {
		mainNode := newNode(t, nil, nil)

		lggr, observedLogs := logger.TestObserved(t, zap.DebugLevel)
		_, txSender := newTestTransactionSender(t, types.RandomID(), lggr,
			[]Node[types.ID, TestSendTxRPCClient]{mainNode},
			[]SendOnlyNode[types.ID, TestSendTxRPCClient]{newNode(t, errors.New("unexpected error"), nil)})

		result := txSender.SendTransaction(tests.Context(t), nil)
		require.NoError(t, result.Error())
		require.Equal(t, Successful, result.Code())
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

		lggr := logger.Test(t)

		_, txSender := newTestTransactionSender(t, types.RandomID(), lggr,
			[]Node[types.ID, TestSendTxRPCClient]{mainNode}, nil)

		requestContext, cancel := context.WithCancel(tests.Context(t))
		cancel()
		result := txSender.SendTransaction(requestContext, nil)
		require.EqualError(t, result.Error(), "context canceled")
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

		lggr := logger.Test(t)

		_, txSender := newTestTransactionSender(t, chainID, lggr, []Node[types.ID, TestSendTxRPCClient]{fastNode, slowNode}, nil)
		result := txSender.SendTransaction(tests.Context(t), nil)
		require.EqualError(t, result.Error(), expectedError.Error())
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
		_, txSender := newTestTransactionSender(t, chainID, lggr,
			[]Node[types.ID, TestSendTxRPCClient]{fastNode, slowNode},
			[]SendOnlyNode[types.ID, TestSendTxRPCClient]{slowSendOnly})

		result := txSender.SendTransaction(tests.Context(t), nil)
		require.NoError(t, result.Error())
		require.Equal(t, Successful, result.Code())
	})
	t.Run("Fails when multinode is closed", func(t *testing.T) {
		chainID := types.RandomID()
		fastNode := newNode(t, nil, nil)
		fastNode.On("ConfiguredChainID").Return(chainID).Maybe()
		// hold reply from the node till end of the test
		testContext, testCancel := context.WithCancel(tests.Context(t))
		defer testCancel()
		slowNode := newNode(t, errors.New("transaction failed"), func(_ mock.Arguments) {
			// block caller til end of the test
			<-testContext.Done()
		})
		slowNode.On("ConfiguredChainID").Return(chainID).Maybe()
		slowSendOnly := newNode(t, errors.New("send only failed"), func(_ mock.Arguments) {
			// block caller til end of the test
			<-testContext.Done()
		})
		slowSendOnly.On("ConfiguredChainID").Return(chainID).Maybe()

		lggr, _ := logger.TestObserved(t, zap.DebugLevel)

		mn, txSender := newTestTransactionSender(t, chainID, lggr,
			[]Node[types.ID, TestSendTxRPCClient]{fastNode, slowNode},
			[]SendOnlyNode[types.ID, TestSendTxRPCClient]{slowSendOnly})

		require.NoError(t, mn.Start(tests.Context(t)))
		require.NoError(t, mn.Close())
		result := txSender.SendTransaction(tests.Context(t), nil)
		require.EqualError(t, result.Error(), "service is stopped")
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

		var txSender *TransactionSender[any, *sendTxResult, types.ID, TestSendTxRPCClient]

		t.Cleanup(func() { // after txSender.Close()
			result := txSender.SendTransaction(tests.Context(t), nil)
			assert.EqualError(t, result.err, "TransactionSender not started")
		})

		_, txSender = newTestTransactionSender(t, chainID, logger.Test(t),
			[]Node[types.ID, TestSendTxRPCClient]{fastNode, slowNode},
			[]SendOnlyNode[types.ID, TestSendTxRPCClient]{slowSendOnly})

	})
	t.Run("Returns error if there is no healthy primary nodes", func(t *testing.T) {
		chainID := types.RandomID()
		primary := newNodeWithState(t, nodeStateUnreachable, nil, nil)
		sendOnly := newNodeWithState(t, nodeStateUnreachable, nil, nil)

		lggr := logger.Test(t)

		_, txSender := newTestTransactionSender(t, chainID, lggr,
			[]Node[types.ID, TestSendTxRPCClient]{primary},
			[]SendOnlyNode[types.ID, TestSendTxRPCClient]{sendOnly})

		result := txSender.SendTransaction(tests.Context(t), nil)
		assert.EqualError(t, result.Error(), ErroringNodeError.Error())
	})

	t.Run("Transaction success even if one of the nodes is unhealthy", func(t *testing.T) {
		chainID := types.RandomID()
		mainNode := newNode(t, nil, nil)
		unexpectedCall := func(args mock.Arguments) {
			panic("SendTx must not be called for unhealthy node")
		}
		unhealthyNode := newNodeWithState(t, nodeStateUnreachable, nil, unexpectedCall)
		unhealthySendOnlyNode := newNodeWithState(t, nodeStateUnreachable, nil, unexpectedCall)

		lggr := logger.Test(t)

		_, txSender := newTestTransactionSender(t, chainID, lggr,
			[]Node[types.ID, TestSendTxRPCClient]{mainNode, unhealthyNode},
			[]SendOnlyNode[types.ID, TestSendTxRPCClient]{unhealthySendOnlyNode})

		result := txSender.SendTransaction(tests.Context(t), nil)
		require.NoError(t, result.Error())
		require.Equal(t, Successful, result.Code())
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
		ResultsByCode       sendTxResults[*sendTxResult]
	}{
		{
			Name:                "Returns success and logs critical error on success and Fatal",
			ExpectedTxResult:    "success",
			ExpectedCriticalErr: "found contradictions in nodes replies on SendTransaction: got success and severe error",
			ResultsByCode: sendTxResults[*sendTxResult]{
				Successful: {NewSendTxResult(errors.New("success"))},
				Fatal:      {NewSendTxResult(errors.New("fatal"))},
			},
		},
		{
			Name:                "Returns TransactionAlreadyKnown and logs critical error on TransactionAlreadyKnown and Fatal",
			ExpectedTxResult:    "tx_already_known",
			ExpectedCriticalErr: "found contradictions in nodes replies on SendTransaction: got success and severe error",
			ResultsByCode: sendTxResults[*sendTxResult]{
				TransactionAlreadyKnown: {NewSendTxResult(errors.New("tx_already_known"))},
				Unsupported:             {NewSendTxResult(errors.New("unsupported"))},
			},
		},
		{
			Name:                "Prefers sever error to temporary",
			ExpectedTxResult:    "underpriced",
			ExpectedCriticalErr: "",
			ResultsByCode: sendTxResults[*sendTxResult]{
				Retryable:   {NewSendTxResult(errors.New("retryable"))},
				Underpriced: {NewSendTxResult(errors.New("underpriced"))},
			},
		},
		{
			Name:                "Returns temporary error",
			ExpectedTxResult:    "retryable",
			ExpectedCriticalErr: "",
			ResultsByCode: sendTxResults[*sendTxResult]{
				Retryable: {NewSendTxResult(errors.New("retryable"))},
			},
		},
		{
			Name:                "Insufficient funds is treated as  error",
			ExpectedTxResult:    "insufficientFunds",
			ExpectedCriticalErr: "",
			ResultsByCode: sendTxResults[*sendTxResult]{
				InsufficientFunds: {NewSendTxResult(errors.New("insufficientFunds"))},
			},
		},
		{
			Name:                "Logs critical error on empty ResultsByCode",
			ExpectedCriticalErr: "expected at least one response on SendTransaction",
			ResultsByCode:       sendTxResults[*sendTxResult]{},
		},
		{
			Name:                "Zk terminally stuck",
			ExpectedTxResult:    "not enough keccak counters to continue the execution",
			ExpectedCriticalErr: "",
			ResultsByCode: sendTxResults[*sendTxResult]{
				TerminallyStuck: {NewSendTxResult(errors.New("not enough keccak counters to continue the execution"))},
			},
		},
	}

	for _, testCase := range testCases {
		for code := range testCase.ResultsByCode {
			delete(codesToCover, code)
		}

		t.Run(testCase.Name, func(t *testing.T) {
			txResult, err := aggregateTxResults(testCase.ResultsByCode)
			if testCase.ExpectedTxResult != "" {
				require.EqualError(t, txResult.Error(), testCase.ExpectedTxResult)
			}

			logger.Sugared(logger.Test(t)).Info("Map: " + fmt.Sprint(testCase.ResultsByCode))
			logger.Sugared(logger.Test(t)).Criticalw("observed invariant violation on SendTransaction", "resultsByCode", testCase.ResultsByCode, "err", err)

			if testCase.ExpectedCriticalErr == "" {
				require.NoError(t, err)
			} else {
				require.EqualError(t, err, testCase.ExpectedCriticalErr)
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
