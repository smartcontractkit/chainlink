package client

import (
	"errors"
	"fmt"
	big "math/big"
	"sync/atomic"
	"testing"

	"github.com/cometbft/cometbft/libs/rand"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	bigmath "github.com/smartcontractkit/chainlink-common/pkg/utils/big_math"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"

	"github.com/smartcontractkit/chainlink/v2/common/types"
	"github.com/smartcontractkit/chainlink/v2/common/types/mocks"
)

func TestUnit_NodeLifecycle_aliveLoop(t *testing.T) {
	t.Parallel()

	newDialedNode := func(t *testing.T, opts testNodeOpts) testNode {
		node := newTestNode(t, opts)
		opts.rpc.On("Close").Return(nil).Once()

		node.setState(nodeStateDialed)
		return node
	}

	t.Run("returns on closed", func(t *testing.T) {
		node := newTestNode(t, testNodeOpts{})
		node.setState(nodeStateClosed)
		node.wg.Add(1)
		node.aliveLoop()

	})
	t.Run("if initial subscribe fails, transitions to unreachable", func(t *testing.T) {
		t.Parallel()
		rpc := newMockNodeClient[types.ID, Head](t)
		node := newDialedNode(t, testNodeOpts{
			rpc: rpc,
		})
		defer func() { assert.NoError(t, node.close()) }()

		expectedError := errors.New("failed to subscribe to rpc")
		rpc.On("Subscribe", mock.Anything, mock.Anything, rpcSubscriptionMethodNewHeads).Return(nil, expectedError).Once()
		rpc.On("DisconnectAll").Once()
		// might be called in unreachable loop
		rpc.On("Dial", mock.Anything).Return(errors.New("failed to dial")).Maybe()
		node.declareAlive()
		tests.AssertEventually(t, func() bool {
			return node.State() == nodeStateUnreachable
		})

	})
	t.Run("if remote RPC connection is closed transitions to unreachable", func(t *testing.T) {
		t.Parallel()
		rpc := newMockNodeClient[types.ID, Head](t)

		lggr, observedLogs := logger.TestObserved(t, zap.WarnLevel)
		node := newDialedNode(t, testNodeOpts{
			rpc:  rpc,
			lggr: lggr,
		})
		defer func() { assert.NoError(t, node.close()) }()

		sub := mocks.NewSubscription(t)
		errChan := make(chan error)
		close(errChan)
		sub.On("Err").Return((<-chan error)(errChan)).Once()
		sub.On("Unsubscribe").Once()
		rpc.On("Subscribe", mock.Anything, mock.Anything, rpcSubscriptionMethodNewHeads).Return(sub, nil).Once()
		rpc.On("SetAliveLoopSub", sub).Once()
		// disconnects all on transfer to unreachable
		rpc.On("DisconnectAll").Once()
		// might be called in unreachable loop
		rpc.On("Dial", mock.Anything).Return(errors.New("failed to dial")).Maybe()
		node.declareAlive()
		tests.AssertLogEventually(t, observedLogs, "Subscription was terminated")
		assert.Equal(t, nodeStateUnreachable, node.State())
	})

	newSubscribedNode := func(t *testing.T, opts testNodeOpts) testNode {
		sub := mocks.NewSubscription(t)
		sub.On("Err").Return((<-chan error)(nil))
		sub.On("Unsubscribe").Once()
		opts.rpc.On("Subscribe", mock.Anything, mock.Anything, rpcSubscriptionMethodNewHeads).Return(sub, nil).Once()
		opts.rpc.On("SetAliveLoopSub", sub).Once()
		return newDialedNode(t, opts)
	}
	t.Run("Stays alive and waits for signal", func(t *testing.T) {
		t.Parallel()
		rpc := newMockNodeClient[types.ID, Head](t)
		lggr, observedLogs := logger.TestObserved(t, zap.DebugLevel)
		node := newSubscribedNode(t, testNodeOpts{
			config: testNodeConfig{},
			rpc:    rpc,
			lggr:   lggr,
		})
		defer func() { assert.NoError(t, node.close()) }()
		node.declareAlive()
		tests.AssertLogEventually(t, observedLogs, "Head liveness checking disabled")
		tests.AssertLogEventually(t, observedLogs, "Polling disabled")
		assert.Equal(t, nodeStateAlive, node.State())
	})
	t.Run("stays alive while below pollFailureThreshold and resets counter on success", func(t *testing.T) {
		t.Parallel()
		rpc := newMockNodeClient[types.ID, Head](t)
		lggr, observedLogs := logger.TestObserved(t, zap.DebugLevel)
		const pollFailureThreshold = 3
		node := newSubscribedNode(t, testNodeOpts{
			config: testNodeConfig{
				pollFailureThreshold: pollFailureThreshold,
				pollInterval:         tests.TestInterval,
			},
			rpc:  rpc,
			lggr: lggr,
		})
		defer func() { assert.NoError(t, node.close()) }()

		pollError := errors.New("failed to get ClientVersion")
		// 1. Return error several times, but below threshold
		rpc.On("ClientVersion", mock.Anything).Return("", pollError).Run(func(_ mock.Arguments) {
			// stays healthy while below threshold
			assert.Equal(t, nodeStateAlive, node.State())
		}).Times(pollFailureThreshold - 1)
		// 2. Successful call that is expected to reset counter
		rpc.On("ClientVersion", mock.Anything).Return("client_version", nil).Once()
		// 3. Return error. If we have not reset the timer, we'll transition to nonAliveState
		rpc.On("ClientVersion", mock.Anything).Return("", pollError).Once()
		// 4. Once during the call, check if node is alive
		var ensuredAlive atomic.Bool
		rpc.On("ClientVersion", mock.Anything).Return("client_version", nil).Run(func(_ mock.Arguments) {
			if ensuredAlive.Load() {
				return
			}
			ensuredAlive.Store(true)
			assert.Equal(t, nodeStateAlive, node.State())
		}).Once()
		// redundant call to stay in alive state
		rpc.On("ClientVersion", mock.Anything).Return("client_version", nil)
		node.declareAlive()
		tests.AssertLogCountEventually(t, observedLogs, fmt.Sprintf("Poll failure, RPC endpoint %s failed to respond properly", node.String()), pollFailureThreshold)
		tests.AssertLogCountEventually(t, observedLogs, "Version poll successful", 2)
		assert.True(t, ensuredAlive.Load(), "expected to ensure that node was alive")

	})
	t.Run("with threshold poll failures, transitions to unreachable", func(t *testing.T) {
		t.Parallel()
		rpc := newMockNodeClient[types.ID, Head](t)
		lggr, observedLogs := logger.TestObserved(t, zap.DebugLevel)
		const pollFailureThreshold = 3
		node := newSubscribedNode(t, testNodeOpts{
			config: testNodeConfig{
				pollFailureThreshold: pollFailureThreshold,
				pollInterval:         tests.TestInterval,
			},
			rpc:  rpc,
			lggr: lggr,
		})
		defer func() { assert.NoError(t, node.close()) }()
		pollError := errors.New("failed to get ClientVersion")
		rpc.On("ClientVersion", mock.Anything).Return("", pollError)
		// disconnects all on transfer to unreachable
		rpc.On("DisconnectAll").Once()
		// might be called in unreachable loop
		rpc.On("Dial", mock.Anything).Return(errors.New("failed to dial")).Maybe()
		node.declareAlive()
		tests.AssertLogCountEventually(t, observedLogs, fmt.Sprintf("Poll failure, RPC endpoint %s failed to respond properly", node.String()), pollFailureThreshold)
		tests.AssertEventually(t, func() bool {
			return nodeStateUnreachable == node.State()
		})
	})
	t.Run("with threshold poll failures, but we are the last node alive, forcibly keeps it alive", func(t *testing.T) {
		t.Parallel()
		rpc := newMockNodeClient[types.ID, Head](t)
		lggr, observedLogs := logger.TestObserved(t, zap.DebugLevel)
		const pollFailureThreshold = 3
		node := newSubscribedNode(t, testNodeOpts{
			config: testNodeConfig{
				pollFailureThreshold: pollFailureThreshold,
				pollInterval:         tests.TestInterval,
			},
			rpc:  rpc,
			lggr: lggr,
		})
		defer func() { assert.NoError(t, node.close()) }()
		node.nLiveNodes = func() (count int, blockNumber int64, totalDifficulty *big.Int) {
			return 1, 20, big.NewInt(10)
		}
		pollError := errors.New("failed to get ClientVersion")
		rpc.On("ClientVersion", mock.Anything).Return("", pollError)
		node.declareAlive()
		tests.AssertLogEventually(t, observedLogs, fmt.Sprintf("RPC endpoint failed to respond to %d consecutive polls", pollFailureThreshold))
		assert.Equal(t, nodeStateAlive, node.State())
	})
	t.Run("when behind more than SyncThreshold, transitions to out of sync", func(t *testing.T) {
		t.Parallel()
		rpc := newMockNodeClient[types.ID, Head](t)
		lggr, observedLogs := logger.TestObserved(t, zap.DebugLevel)
		const syncThreshold = 10
		node := newSubscribedNode(t, testNodeOpts{
			config: testNodeConfig{
				pollInterval:  tests.TestInterval,
				syncThreshold: syncThreshold,
				selectionMode: NodeSelectionModeRoundRobin,
			},
			rpc:  rpc,
			lggr: lggr,
		})
		defer func() { assert.NoError(t, node.close()) }()
		node.stateLatestBlockNumber = 20
		node.nLiveNodes = func() (count int, blockNumber int64, totalDifficulty *big.Int) {
			return 10, syncThreshold + node.stateLatestBlockNumber + 1, big.NewInt(10)
		}
		rpc.On("ClientVersion", mock.Anything).Return("", nil)
		// tries to redial in outOfSync
		rpc.On("Dial", mock.Anything).Return(errors.New("failed to dial")).Run(func(_ mock.Arguments) {
			assert.Equal(t, nodeStateOutOfSync, node.State())
		}).Once()
		// disconnects all on transfer to unreachable or outOfSync
		rpc.On("DisconnectAll").Maybe()
		// might be called in unreachable loop
		rpc.On("Dial", mock.Anything).Return(errors.New("failed to dial")).Maybe()
		node.declareAlive()
		tests.AssertLogEventually(t, observedLogs, "Failed to dial out-of-sync RPC node")
	})
	t.Run("when behind more than SyncThreshold but we are the last live node, forcibly stays alive", func(t *testing.T) {
		t.Parallel()
		rpc := newMockNodeClient[types.ID, Head](t)
		lggr, observedLogs := logger.TestObserved(t, zap.DebugLevel)
		const syncThreshold = 10
		node := newSubscribedNode(t, testNodeOpts{
			config: testNodeConfig{
				pollInterval:  tests.TestInterval,
				syncThreshold: syncThreshold,
				selectionMode: NodeSelectionModeRoundRobin,
			},
			rpc:  rpc,
			lggr: lggr,
		})
		defer func() { assert.NoError(t, node.close()) }()
		node.stateLatestBlockNumber = 20
		node.nLiveNodes = func() (count int, blockNumber int64, totalDifficulty *big.Int) {
			return 1, syncThreshold + node.stateLatestBlockNumber + 1, big.NewInt(10)
		}
		rpc.On("ClientVersion", mock.Anything).Return("", nil)
		node.declareAlive()
		tests.AssertLogEventually(t, observedLogs, fmt.Sprintf("RPC endpoint has fallen behind; %s %s", msgCannotDisable, msgDegradedState))
	})
	t.Run("when behind but SyncThreshold=0, stay alive", func(t *testing.T) {
		t.Parallel()
		rpc := newMockNodeClient[types.ID, Head](t)
		lggr, observedLogs := logger.TestObserved(t, zap.DebugLevel)
		node := newSubscribedNode(t, testNodeOpts{
			config: testNodeConfig{
				pollInterval:  tests.TestInterval,
				syncThreshold: 0,
				selectionMode: NodeSelectionModeRoundRobin,
			},
			rpc:  rpc,
			lggr: lggr,
		})
		defer func() { assert.NoError(t, node.close()) }()
		node.stateLatestBlockNumber = 20
		node.nLiveNodes = func() (count int, blockNumber int64, totalDifficulty *big.Int) {
			return 1, node.stateLatestBlockNumber + 100, big.NewInt(10)
		}
		rpc.On("ClientVersion", mock.Anything).Return("", nil)
		node.declareAlive()
		tests.AssertLogCountEventually(t, observedLogs, "Version poll successful", 2)
		assert.Equal(t, nodeStateAlive, node.State())
	})

	t.Run("when no new heads received for threshold, transitions to out of sync", func(t *testing.T) {
		t.Parallel()
		rpc := newMockNodeClient[types.ID, Head](t)
		node := newSubscribedNode(t, testNodeOpts{
			config:              testNodeConfig{},
			noNewHeadsThreshold: tests.TestInterval,
			rpc:                 rpc,
		})
		defer func() { assert.NoError(t, node.close()) }()
		// tries to redial in outOfSync
		rpc.On("Dial", mock.Anything).Return(errors.New("failed to dial")).Run(func(_ mock.Arguments) {
			assert.Equal(t, nodeStateOutOfSync, node.State())
		}).Once()
		// disconnects all on transfer to unreachable or outOfSync
		rpc.On("DisconnectAll").Maybe()
		// might be called in unreachable loop
		rpc.On("Dial", mock.Anything).Return(errors.New("failed to dial")).Maybe()
		node.declareAlive()
		tests.AssertEventually(t, func() bool {
			// right after outOfSync we'll transfer to unreachable due to returned error on Dial
			// we check that we were in out of sync state on first Dial call
			return node.State() == nodeStateUnreachable
		})
	})
	t.Run("when no new heads received for threshold but we are the last live node, forcibly stays alive", func(t *testing.T) {
		t.Parallel()
		rpc := newMockNodeClient[types.ID, Head](t)
		lggr, observedLogs := logger.TestObserved(t, zap.DebugLevel)
		node := newSubscribedNode(t, testNodeOpts{
			config:              testNodeConfig{},
			lggr:                lggr,
			noNewHeadsThreshold: tests.TestInterval,
			rpc:                 rpc,
		})
		defer func() { assert.NoError(t, node.close()) }()
		node.nLiveNodes = func() (count int, blockNumber int64, totalDifficulty *big.Int) {
			return 1, 20, big.NewInt(10)
		}
		node.declareAlive()
		tests.AssertLogEventually(t, observedLogs, fmt.Sprintf("RPC endpoint detected out of sync; %s %s", msgCannotDisable, msgDegradedState))
		assert.Equal(t, nodeStateAlive, node.State())
	})

	t.Run("rpc closed head channel", func(t *testing.T) {
		t.Parallel()
		rpc := newMockNodeClient[types.ID, Head](t)
		sub := mocks.NewSubscription(t)
		sub.On("Err").Return((<-chan error)(nil))
		sub.On("Unsubscribe").Once()
		rpc.On("Subscribe", mock.Anything, mock.Anything, rpcSubscriptionMethodNewHeads).Run(func(args mock.Arguments) {
			ch := args.Get(1).(chan<- Head)
			close(ch)
		}).Return(sub, nil).Once()
		rpc.On("SetAliveLoopSub", sub).Once()
		lggr, observedLogs := logger.TestObserved(t, zap.ErrorLevel)
		node := newDialedNode(t, testNodeOpts{
			lggr:                lggr,
			config:              testNodeConfig{},
			noNewHeadsThreshold: tests.TestInterval,
			rpc:                 rpc,
		})
		defer func() { assert.NoError(t, node.close()) }()
		// disconnects all on transfer to unreachable or outOfSync
		rpc.On("DisconnectAll").Once()
		// might be called in unreachable loop
		rpc.On("Dial", mock.Anything).Return(errors.New("failed to dial")).Maybe()
		node.declareAlive()
		tests.AssertLogEventually(t, observedLogs, "Subscription channel unexpectedly closed")
		assert.Equal(t, nodeStateUnreachable, node.State())

	})
	t.Run("updates block number and difficulty on new head", func(t *testing.T) {
		t.Parallel()
		rpc := newMockNodeClient[types.ID, Head](t)
		sub := mocks.NewSubscription(t)
		sub.On("Err").Return((<-chan error)(nil))
		sub.On("Unsubscribe").Once()
		expectedBlockNumber := rand.Int64()
		expectedDiff := big.NewInt(rand.Int64())
		rpc.On("Subscribe", mock.Anything, mock.Anything, rpcSubscriptionMethodNewHeads).Run(func(args mock.Arguments) {
			ch := args.Get(1).(chan<- Head)
			go writeHeads(t, ch, head{BlockNumber: expectedBlockNumber, BlockDifficulty: expectedDiff})
		}).Return(sub, nil).Once()
		rpc.On("SetAliveLoopSub", sub).Once()
		node := newDialedNode(t, testNodeOpts{
			config: testNodeConfig{},
			rpc:    rpc,
		})
		defer func() { assert.NoError(t, node.close()) }()
		node.declareAlive()
		tests.AssertEventually(t, func() bool {
			state, block, diff := node.StateAndLatest()
			return state == nodeStateAlive && block == expectedBlockNumber == bigmath.Equal(diff, expectedDiff)
		})
	})
}

type head struct {
	BlockNumber     int64
	BlockDifficulty *big.Int
}

func writeHeads(t *testing.T, ch chan<- Head, heads ...head) {
	for _, head := range heads {
		h := newMockHead(t)
		h.On("BlockNumber").Return(head.BlockNumber)
		h.On("BlockDifficulty").Return(head.BlockDifficulty)
		select {
		case ch <- h:
		case <-tests.Context(t).Done():
			return
		}
	}
}

func setupRPCForAliveLoop(t *testing.T, rpc *mockNodeClient[types.ID, Head]) {
	rpc.On("Dial", mock.Anything).Return(nil).Maybe()
	aliveSubscription := mocks.NewSubscription(t)
	aliveSubscription.On("Err").Return((<-chan error)(nil)).Maybe()
	aliveSubscription.On("Unsubscribe").Maybe()
	rpc.On("Subscribe", mock.Anything, mock.Anything, rpcSubscriptionMethodNewHeads).Return(aliveSubscription, nil).Maybe()
	rpc.On("SetAliveLoopSub", mock.Anything).Maybe()
}

func TestUnit_NodeLifecycle_outOfSyncLoop(t *testing.T) {
	t.Parallel()

	newAliveNode := func(t *testing.T, opts testNodeOpts) testNode {
		node := newTestNode(t, opts)
		opts.rpc.On("Close").Return(nil).Once()
		// disconnects all on transfer to unreachable or outOfSync
		opts.rpc.On("DisconnectAll")
		node.setState(nodeStateAlive)
		return node
	}

	stubIsOutOfSync := func(num int64, td *big.Int) bool {
		return false
	}

	t.Run("returns on closed", func(t *testing.T) {
		t.Parallel()
		node := newTestNode(t, testNodeOpts{})
		node.setState(nodeStateClosed)
		node.wg.Add(1)
		node.outOfSyncLoop(stubIsOutOfSync)
	})
	t.Run("on old blocks stays outOfSync and returns on close", func(t *testing.T) {
		t.Parallel()
		rpc := newMockNodeClient[types.ID, Head](t)
		nodeChainID := types.RandomID()
		lggr, observedLogs := logger.TestObserved(t, zap.DebugLevel)
		node := newAliveNode(t, testNodeOpts{
			rpc:     rpc,
			chainID: nodeChainID,
			lggr:    lggr,
		})
		defer func() { assert.NoError(t, node.close()) }()

		rpc.On("Dial", mock.Anything).Return(nil).Once()
		// might be called multiple times
		rpc.On("ChainID", mock.Anything).Return(nodeChainID, nil).Once()

		outOfSyncSubscription := mocks.NewSubscription(t)
		outOfSyncSubscription.On("Err").Return((<-chan error)(nil))
		outOfSyncSubscription.On("Unsubscribe").Once()
		heads := []head{{BlockNumber: 7}, {BlockNumber: 11}, {BlockNumber: 13}}
		rpc.On("Subscribe", mock.Anything, mock.Anything, rpcSubscriptionMethodNewHeads).Run(func(args mock.Arguments) {
			ch := args.Get(1).(chan<- Head)
			go writeHeads(t, ch, heads...)
		}).Return(outOfSyncSubscription, nil).Once()
		rpc.On("Dial", mock.Anything).Return(errors.New("failed to redial")).Maybe()

		node.declareOutOfSync(func(num int64, td *big.Int) bool {
			return true
		})
		tests.AssertLogCountEventually(t, observedLogs, msgReceivedBlock, len(heads))
		assert.Equal(t, nodeStateOutOfSync, node.State())
	})
	t.Run("if initial dial fails, transitions to unreachable", func(t *testing.T) {
		t.Parallel()
		rpc := newMockNodeClient[types.ID, Head](t)
		node := newAliveNode(t, testNodeOpts{
			rpc: rpc,
		})
		defer func() { assert.NoError(t, node.close()) }()

		expectedError := errors.New("failed to dial rpc")
		// might be called again in unreachable loop, so no need to set once
		rpc.On("Dial", mock.Anything).Return(expectedError)
		node.declareOutOfSync(stubIsOutOfSync)
		tests.AssertEventually(t, func() bool {
			return node.State() == nodeStateUnreachable
		})
	})
	t.Run("if fail to get chainID, transitions to invalidChainID", func(t *testing.T) {
		t.Parallel()
		rpc := newMockNodeClient[types.ID, Head](t)
		node := newAliveNode(t, testNodeOpts{
			rpc: rpc,
		})
		defer func() { assert.NoError(t, node.close()) }()

		rpc.On("Dial", mock.Anything).Return(nil).Once()
		expectedError := errors.New("failed to get chain ID")
		// might be called multiple times
		rpc.On("ChainID", mock.Anything).Return(types.NewIDFromInt(0), expectedError)
		node.declareOutOfSync(stubIsOutOfSync)
		tests.AssertEventually(t, func() bool {
			return node.State() == nodeStateInvalidChainID
		})
	})
	t.Run("if chainID does not match, transitions to invalidChainID", func(t *testing.T) {
		t.Parallel()
		rpc := newMockNodeClient[types.ID, Head](t)
		nodeChainID := types.NewIDFromInt(10)
		rpcChainID := types.NewIDFromInt(11)
		node := newAliveNode(t, testNodeOpts{
			rpc:     rpc,
			chainID: nodeChainID,
		})
		defer func() { assert.NoError(t, node.close()) }()

		rpc.On("Dial", mock.Anything).Return(nil).Once()
		// might be called multiple times
		rpc.On("ChainID", mock.Anything).Return(rpcChainID, nil)
		node.declareOutOfSync(stubIsOutOfSync)
		tests.AssertEventually(t, func() bool {
			return node.State() == nodeStateInvalidChainID
		})
	})
	t.Run("if fails to subscribe, becomes unreachable", func(t *testing.T) {
		t.Parallel()
		rpc := newMockNodeClient[types.ID, Head](t)
		nodeChainID := types.RandomID()
		node := newAliveNode(t, testNodeOpts{
			rpc:     rpc,
			chainID: nodeChainID,
		})
		defer func() { assert.NoError(t, node.close()) }()

		rpc.On("Dial", mock.Anything).Return(nil).Once()
		// might be called multiple times
		rpc.On("ChainID", mock.Anything).Return(nodeChainID, nil).Once()
		expectedError := errors.New("failed to subscribe")
		rpc.On("Subscribe", mock.Anything, mock.Anything, rpcSubscriptionMethodNewHeads).Return(nil, expectedError)
		rpc.On("Dial", mock.Anything).Return(errors.New("failed to redial")).Maybe()
		node.declareOutOfSync(stubIsOutOfSync)
		tests.AssertEventually(t, func() bool {
			return node.State() == nodeStateUnreachable
		})
	})
	t.Run("on subscription termination becomes unreachable", func(t *testing.T) {
		t.Parallel()
		rpc := newMockNodeClient[types.ID, Head](t)
		nodeChainID := types.RandomID()
		lggr, observedLogs := logger.TestObserved(t, zap.ErrorLevel)
		node := newAliveNode(t, testNodeOpts{
			rpc:     rpc,
			chainID: nodeChainID,
			lggr:    lggr,
		})
		defer func() { assert.NoError(t, node.close()) }()

		rpc.On("Dial", mock.Anything).Return(nil).Once()
		// might be called multiple times
		rpc.On("ChainID", mock.Anything).Return(nodeChainID, nil).Once()

		sub := mocks.NewSubscription(t)
		errChan := make(chan error, 1)
		errChan <- errors.New("subscription was terminate")
		sub.On("Err").Return((<-chan error)(errChan))
		sub.On("Unsubscribe").Once()
		rpc.On("Subscribe", mock.Anything, mock.Anything, rpcSubscriptionMethodNewHeads).Return(sub, nil).Once()
		rpc.On("Dial", mock.Anything).Return(errors.New("failed to redial")).Maybe()
		node.declareOutOfSync(stubIsOutOfSync)
		tests.AssertLogEventually(t, observedLogs, "Subscription was terminated")
		tests.AssertEventually(t, func() bool {
			return node.State() == nodeStateUnreachable
		})
	})
	t.Run("becomes unreachable if head channel is closed", func(t *testing.T) {
		t.Parallel()
		rpc := newMockNodeClient[types.ID, Head](t)
		nodeChainID := types.RandomID()
		lggr, observedLogs := logger.TestObserved(t, zap.ErrorLevel)
		node := newAliveNode(t, testNodeOpts{
			rpc:     rpc,
			chainID: nodeChainID,
			lggr:    lggr,
		})
		defer func() { assert.NoError(t, node.close()) }()

		rpc.On("Dial", mock.Anything).Return(nil).Once()
		// might be called multiple times
		rpc.On("ChainID", mock.Anything).Return(nodeChainID, nil).Once()

		sub := mocks.NewSubscription(t)
		sub.On("Err").Return((<-chan error)(nil))
		sub.On("Unsubscribe").Once()
		rpc.On("Subscribe", mock.Anything, mock.Anything, rpcSubscriptionMethodNewHeads).Run(func(args mock.Arguments) {
			ch := args.Get(1).(chan<- Head)
			close(ch)
		}).Return(sub, nil).Once()
		rpc.On("Dial", mock.Anything).Return(errors.New("failed to redial")).Maybe()
		node.declareOutOfSync(stubIsOutOfSync)
		tests.AssertLogEventually(t, observedLogs, "Subscription channel unexpectedly closed")
		tests.AssertEventually(t, func() bool {
			return node.State() == nodeStateUnreachable
		})
	})

	t.Run("becomes alive if it receives a newer head", func(t *testing.T) {
		t.Parallel()
		rpc := newMockNodeClient[types.ID, Head](t)
		nodeChainID := types.RandomID()
		lggr, observedLogs := logger.TestObserved(t, zap.DebugLevel)
		node := newAliveNode(t, testNodeOpts{
			rpc:     rpc,
			chainID: nodeChainID,
			lggr:    lggr,
		})
		defer func() { assert.NoError(t, node.close()) }()

		rpc.On("Dial", mock.Anything).Return(nil).Once()
		// might be called multiple times
		rpc.On("ChainID", mock.Anything).Return(nodeChainID, nil).Once()

		outOfSyncSubscription := mocks.NewSubscription(t)
		outOfSyncSubscription.On("Err").Return((<-chan error)(nil))
		outOfSyncSubscription.On("Unsubscribe").Once()
		const highestBlock = 1000
		rpc.On("Subscribe", mock.Anything, mock.Anything, rpcSubscriptionMethodNewHeads).Run(func(args mock.Arguments) {
			ch := args.Get(1).(chan<- Head)
			go writeHeads(t, ch, head{BlockNumber: highestBlock - 1}, head{BlockNumber: highestBlock})
		}).Return(outOfSyncSubscription, nil).Once()

		setupRPCForAliveLoop(t, rpc)

		node.declareOutOfSync(func(num int64, td *big.Int) bool {
			return num < highestBlock
		})
		tests.AssertLogEventually(t, observedLogs, msgReceivedBlock)
		tests.AssertLogEventually(t, observedLogs, msgInSync)
		tests.AssertEventually(t, func() bool {
			return node.State() == nodeStateAlive
		})
	})
	t.Run("becomes alive if there is no other nodes", func(t *testing.T) {
		t.Parallel()
		rpc := newMockNodeClient[types.ID, Head](t)
		nodeChainID := types.RandomID()
		lggr, observedLogs := logger.TestObserved(t, zap.DebugLevel)
		node := newAliveNode(t, testNodeOpts{
			noNewHeadsThreshold: tests.TestInterval,
			rpc:                 rpc,
			chainID:             nodeChainID,
			lggr:                lggr,
		})
		defer func() { assert.NoError(t, node.close()) }()
		node.nLiveNodes = func() (count int, blockNumber int64, totalDifficulty *big.Int) {
			return 0, 100, big.NewInt(200)
		}

		rpc.On("Dial", mock.Anything).Return(nil).Once()
		// might be called multiple times
		rpc.On("ChainID", mock.Anything).Return(nodeChainID, nil).Once()

		outOfSyncSubscription := mocks.NewSubscription(t)
		outOfSyncSubscription.On("Err").Return((<-chan error)(nil))
		outOfSyncSubscription.On("Unsubscribe").Once()
		rpc.On("Subscribe", mock.Anything, mock.Anything, rpcSubscriptionMethodNewHeads).Return(outOfSyncSubscription, nil).Once()

		setupRPCForAliveLoop(t, rpc)

		node.declareOutOfSync(stubIsOutOfSync)
		tests.AssertLogEventually(t, observedLogs, "RPC endpoint is still out of sync, but there are no other available nodes. This RPC node will be forcibly moved back into the live pool in a degraded state")
		tests.AssertEventually(t, func() bool {
			return node.State() == nodeStateAlive
		})
	})
}

func TestUnit_NodeLifecycle_unreachableLoop(t *testing.T) {
	t.Parallel()

	newAliveNode := func(t *testing.T, opts testNodeOpts) testNode {
		node := newTestNode(t, opts)
		opts.rpc.On("Close").Return(nil).Once()
		// disconnects all on transfer to unreachable
		opts.rpc.On("DisconnectAll")

		node.setState(nodeStateAlive)
		return node
	}
	t.Run("returns on closed", func(t *testing.T) {
		t.Parallel()
		node := newTestNode(t, testNodeOpts{})
		node.setState(nodeStateClosed)
		node.wg.Add(1)
		node.unreachableLoop()

	})
	t.Run("on failed redial, keeps trying", func(t *testing.T) {
		t.Parallel()
		rpc := newMockNodeClient[types.ID, Head](t)
		nodeChainID := types.RandomID()
		lggr, observedLogs := logger.TestObserved(t, zap.DebugLevel)
		node := newAliveNode(t, testNodeOpts{
			rpc:     rpc,
			chainID: nodeChainID,
			lggr:    lggr,
		})
		defer func() { assert.NoError(t, node.close()) }()

		rpc.On("Dial", mock.Anything).Return(errors.New("failed to dial"))
		node.declareUnreachable()
		tests.AssertLogCountEventually(t, observedLogs, "Failed to redial RPC node; still unreachable", 2)
	})
	t.Run("on failed chainID verification, keep trying", func(t *testing.T) {
		t.Parallel()
		rpc := newMockNodeClient[types.ID, Head](t)
		nodeChainID := types.RandomID()
		lggr, observedLogs := logger.TestObserved(t, zap.DebugLevel)
		node := newAliveNode(t, testNodeOpts{
			rpc:     rpc,
			chainID: nodeChainID,
			lggr:    lggr,
		})
		defer func() { assert.NoError(t, node.close()) }()

		rpc.On("Dial", mock.Anything).Return(nil)
		rpc.On("ChainID", mock.Anything).Run(func(_ mock.Arguments) {
			assert.Equal(t, nodeStateDialed, node.State())
		}).Return(nodeChainID, errors.New("failed to get chain id"))
		node.declareUnreachable()
		tests.AssertLogCountEventually(t, observedLogs, "Failed to redial RPC node; verify failed", 2)
	})
	t.Run("on chain ID mismatch transitions to invalidChainID", func(t *testing.T) {
		t.Parallel()
		rpc := newMockNodeClient[types.ID, Head](t)
		nodeChainID := types.NewIDFromInt(10)
		rpcChainID := types.NewIDFromInt(11)
		node := newAliveNode(t, testNodeOpts{
			rpc:     rpc,
			chainID: nodeChainID,
		})
		defer func() { assert.NoError(t, node.close()) }()

		rpc.On("Dial", mock.Anything).Return(nil)
		rpc.On("ChainID", mock.Anything).Return(rpcChainID, nil)
		node.declareUnreachable()
		tests.AssertEventually(t, func() bool {
			return node.State() == nodeStateInvalidChainID
		})
	})
	t.Run("on valid chain ID becomes alive", func(t *testing.T) {
		t.Parallel()
		rpc := newMockNodeClient[types.ID, Head](t)
		nodeChainID := types.RandomID()
		node := newAliveNode(t, testNodeOpts{
			rpc:     rpc,
			chainID: nodeChainID,
		})
		defer func() { assert.NoError(t, node.close()) }()

		rpc.On("Dial", mock.Anything).Return(nil)
		rpc.On("ChainID", mock.Anything).Return(nodeChainID, nil)

		setupRPCForAliveLoop(t, rpc)

		node.declareUnreachable()
		tests.AssertEventually(t, func() bool {
			return node.State() == nodeStateAlive
		})
	})
}

func TestUnit_NodeLifecycle_invalidChainIDLoop(t *testing.T) {
	t.Parallel()
	newDialedNode := func(t *testing.T, opts testNodeOpts) testNode {
		node := newTestNode(t, opts)
		opts.rpc.On("Close").Return(nil).Once()
		opts.rpc.On("DisconnectAll")

		node.setState(nodeStateDialed)
		return node
	}
	t.Run("returns on closed", func(t *testing.T) {
		t.Parallel()
		node := newTestNode(t, testNodeOpts{})
		node.setState(nodeStateClosed)
		node.wg.Add(1)
		node.invalidChainIDLoop()

	})
	t.Run("on failed chainID call becomes unreachable", func(t *testing.T) {
		t.Parallel()
		rpc := newMockNodeClient[types.ID, Head](t)
		nodeChainID := types.RandomID()
		lggr, observedLogs := logger.TestObserved(t, zap.DebugLevel)
		node := newDialedNode(t, testNodeOpts{
			rpc:     rpc,
			chainID: nodeChainID,
			lggr:    lggr,
		})
		defer func() { assert.NoError(t, node.close()) }()

		rpc.On("ChainID", mock.Anything).Return(nodeChainID, errors.New("failed to get chain id"))
		// for unreachable loop
		rpc.On("Dial", mock.Anything).Return(errors.New("failed to dial")).Maybe()
		node.declareInvalidChainID()
		tests.AssertLogEventually(t, observedLogs, "Unexpected error while verifying RPC node chain ID")
		tests.AssertEventually(t, func() bool {
			return node.State() == nodeStateUnreachable
		})
	})
	t.Run("on chainID mismatch keeps trying", func(t *testing.T) {
		t.Parallel()
		rpc := newMockNodeClient[types.ID, Head](t)
		nodeChainID := types.NewIDFromInt(10)
		rpcChainID := types.NewIDFromInt(11)
		lggr, observedLogs := logger.TestObserved(t, zap.DebugLevel)
		node := newDialedNode(t, testNodeOpts{
			rpc:     rpc,
			chainID: nodeChainID,
			lggr:    lggr,
		})
		defer func() { assert.NoError(t, node.close()) }()

		rpc.On("ChainID", mock.Anything).Return(rpcChainID, nil)
		node.declareInvalidChainID()
		tests.AssertLogCountEventually(t, observedLogs, "Failed to verify RPC node; remote endpoint returned the wrong chain ID", 2)
		tests.AssertEventually(t, func() bool {
			return node.State() == nodeStateInvalidChainID
		})
	})
	t.Run("on valid chainID becomes alive", func(t *testing.T) {
		t.Parallel()
		rpc := newMockNodeClient[types.ID, Head](t)
		nodeChainID := types.RandomID()
		node := newDialedNode(t, testNodeOpts{
			rpc:     rpc,
			chainID: nodeChainID,
		})
		defer func() { assert.NoError(t, node.close()) }()

		rpc.On("ChainID", mock.Anything).Return(nodeChainID, nil)

		setupRPCForAliveLoop(t, rpc)

		node.declareInvalidChainID()
		tests.AssertEventually(t, func() bool {
			return node.State() == nodeStateAlive
		})
	})
}

func TestUnit_NodeLifecycle_start(t *testing.T) {
	t.Parallel()

	newNode := func(t *testing.T, opts testNodeOpts) testNode {
		node := newTestNode(t, opts)
		opts.rpc.On("Close").Return(nil).Once()

		return node
	}
	t.Run("if fails on initial dial, becomes unreachable", func(t *testing.T) {
		t.Parallel()
		rpc := newMockNodeClient[types.ID, Head](t)
		nodeChainID := types.RandomID()
		lggr, observedLogs := logger.TestObserved(t, zap.DebugLevel)
		node := newNode(t, testNodeOpts{
			rpc:     rpc,
			chainID: nodeChainID,
			lggr:    lggr,
		})
		defer func() { assert.NoError(t, node.close()) }()

		rpc.On("Dial", mock.Anything).Return(errors.New("failed to dial"))
		// disconnects all on transfer to unreachable
		rpc.On("DisconnectAll")
		err := node.Start(tests.Context(t))
		assert.NoError(t, err)
		tests.AssertLogEventually(t, observedLogs, "Dial failed: Node is unreachable")
		tests.AssertEventually(t, func() bool {
			return node.State() == nodeStateUnreachable
		})
	})
	t.Run("if chainID verification fails, becomes unreachable", func(t *testing.T) {
		t.Parallel()
		rpc := newMockNodeClient[types.ID, Head](t)
		nodeChainID := types.RandomID()
		lggr, observedLogs := logger.TestObserved(t, zap.DebugLevel)
		node := newNode(t, testNodeOpts{
			rpc:     rpc,
			chainID: nodeChainID,
			lggr:    lggr,
		})
		defer func() { assert.NoError(t, node.close()) }()

		rpc.On("Dial", mock.Anything).Return(nil)
		rpc.On("ChainID", mock.Anything).Run(func(_ mock.Arguments) {
			assert.Equal(t, nodeStateDialed, node.State())
		}).Return(nodeChainID, errors.New("failed to get chain id"))
		// disconnects all on transfer to unreachable
		rpc.On("DisconnectAll")
		err := node.Start(tests.Context(t))
		assert.NoError(t, err)
		tests.AssertLogEventually(t, observedLogs, "Verify failed")
		tests.AssertEventually(t, func() bool {
			return node.State() == nodeStateUnreachable
		})
	})
	t.Run("on chain ID mismatch transitions to invalidChainID", func(t *testing.T) {
		t.Parallel()
		rpc := newMockNodeClient[types.ID, Head](t)
		nodeChainID := types.NewIDFromInt(10)
		rpcChainID := types.NewIDFromInt(11)
		node := newNode(t, testNodeOpts{
			rpc:     rpc,
			chainID: nodeChainID,
		})
		defer func() { assert.NoError(t, node.close()) }()

		rpc.On("Dial", mock.Anything).Return(nil)
		rpc.On("ChainID", mock.Anything).Return(rpcChainID, nil)
		// disconnects all on transfer to unreachable
		rpc.On("DisconnectAll")
		err := node.Start(tests.Context(t))
		assert.NoError(t, err)
		tests.AssertEventually(t, func() bool {
			return node.State() == nodeStateInvalidChainID
		})
	})
	t.Run("on valid chain ID becomes alive", func(t *testing.T) {
		t.Parallel()
		rpc := newMockNodeClient[types.ID, Head](t)
		nodeChainID := types.RandomID()
		node := newNode(t, testNodeOpts{
			rpc:     rpc,
			chainID: nodeChainID,
		})
		defer func() { assert.NoError(t, node.close()) }()

		rpc.On("Dial", mock.Anything).Return(nil)
		rpc.On("ChainID", mock.Anything).Return(nodeChainID, nil)

		setupRPCForAliveLoop(t, rpc)

		err := node.Start(tests.Context(t))
		assert.NoError(t, err)
		tests.AssertEventually(t, func() bool {
			return node.State() == nodeStateAlive
		})
	})
}

func TestUnit_NodeLifecycle_syncStatus(t *testing.T) {
	t.Parallel()
	t.Run("skip if nLiveNodes is not configured", func(t *testing.T) {
		node := newTestNode(t, testNodeOpts{})
		outOfSync, liveNodes := node.syncStatus(0, nil)
		assert.Equal(t, false, outOfSync)
		assert.Equal(t, 0, liveNodes)
	})
	t.Run("skip if syncThreshold is not configured", func(t *testing.T) {
		node := newTestNode(t, testNodeOpts{})
		node.nLiveNodes = func() (count int, blockNumber int64, totalDifficulty *big.Int) {
			return
		}
		outOfSync, liveNodes := node.syncStatus(0, nil)
		assert.Equal(t, false, outOfSync)
		assert.Equal(t, 0, liveNodes)
	})
	t.Run("panics on invalid selection mode", func(t *testing.T) {
		node := newTestNode(t, testNodeOpts{
			config: testNodeConfig{syncThreshold: 1},
		})
		node.nLiveNodes = func() (count int, blockNumber int64, totalDifficulty *big.Int) {
			return
		}
		assert.Panics(t, func() {
			_, _ = node.syncStatus(0, nil)
		})
	})
	t.Run("block height selection mode", func(t *testing.T) {
		const syncThreshold = 10
		const highestBlock = 1000
		const nodesNum = 20
		const totalDifficulty = 3000
		testCases := []struct {
			name        string
			blockNumber int64
			outOfSync   bool
		}{
			{
				name:        "below threshold",
				blockNumber: highestBlock - syncThreshold - 1,
				outOfSync:   true,
			},
			{
				name:        "equal to threshold",
				blockNumber: highestBlock - syncThreshold,
				outOfSync:   false,
			},
			{
				name:        "equal to highest block",
				blockNumber: highestBlock,
				outOfSync:   false,
			},
			{
				name:        "higher than highest block",
				blockNumber: highestBlock,
				outOfSync:   false,
			},
		}

		for _, selectionMode := range []string{NodeSelectionModeHighestHead, NodeSelectionModeRoundRobin, NodeSelectionModePriorityLevel} {
			node := newTestNode(t, testNodeOpts{
				config: testNodeConfig{
					syncThreshold: syncThreshold,
					selectionMode: selectionMode,
				},
			})
			node.nLiveNodes = func() (int, int64, *big.Int) {
				return nodesNum, highestBlock, big.NewInt(totalDifficulty)
			}
			for _, td := range []int64{totalDifficulty - syncThreshold - 1, totalDifficulty - syncThreshold, totalDifficulty, totalDifficulty + 1} {
				for _, testCase := range testCases {
					t.Run(fmt.Sprintf("%s: selectionMode: %s: total difficulty: %d", testCase.name, selectionMode, td), func(t *testing.T) {
						outOfSync, liveNodes := node.syncStatus(testCase.blockNumber, big.NewInt(td))
						assert.Equal(t, nodesNum, liveNodes)
						assert.Equal(t, testCase.outOfSync, outOfSync)
					})
				}
			}
		}

	})
	t.Run("total difficulty selection mode", func(t *testing.T) {
		const syncThreshold = 10
		const highestBlock = 1000
		const nodesNum = 20
		const totalDifficulty = 3000
		testCases := []struct {
			name            string
			totalDifficulty int64
			outOfSync       bool
		}{
			{
				name:            "below threshold",
				totalDifficulty: totalDifficulty - syncThreshold - 1,
				outOfSync:       true,
			},
			{
				name:            "equal to threshold",
				totalDifficulty: totalDifficulty - syncThreshold,
				outOfSync:       false,
			},
			{
				name:            "equal to highest block",
				totalDifficulty: totalDifficulty,
				outOfSync:       false,
			},
			{
				name:            "higher than highest block",
				totalDifficulty: totalDifficulty,
				outOfSync:       false,
			},
		}

		node := newTestNode(t, testNodeOpts{
			config: testNodeConfig{
				syncThreshold: syncThreshold,
				selectionMode: NodeSelectionModeTotalDifficulty,
			},
		})
		node.nLiveNodes = func() (int, int64, *big.Int) {
			return nodesNum, highestBlock, big.NewInt(totalDifficulty)
		}
		for _, hb := range []int64{highestBlock - syncThreshold - 1, highestBlock - syncThreshold, highestBlock, highestBlock + 1} {
			for _, testCase := range testCases {
				t.Run(fmt.Sprintf("%s: selectionMode: %s: highest block: %d", testCase.name, NodeSelectionModeTotalDifficulty, hb), func(t *testing.T) {
					outOfSync, liveNodes := node.syncStatus(hb, big.NewInt(testCase.totalDifficulty))
					assert.Equal(t, nodesNum, liveNodes)
					assert.Equal(t, testCase.outOfSync, outOfSync)
				})
			}
		}

	})
}
