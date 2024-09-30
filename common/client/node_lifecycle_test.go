package client

import (
	"errors"
	"fmt"
	"math/big"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/cometbft/cometbft/libs/rand"
	prom "github.com/prometheus/client_model/go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"

	clientMocks "github.com/smartcontractkit/chainlink/v2/common/client/mocks"
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
		rpc.On("DisconnectAll").Once()
		rpc.On("SubscribeToHeads", mock.Anything).Return(nil, nil, expectedError).Once()
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
		rpc.On("GetInterceptedChainInfo").Return(ChainInfo{}, ChainInfo{}).Once()
		defer func() { assert.NoError(t, node.close()) }()

		sub := mocks.NewSubscription(t)
		errChan := make(chan error)
		close(errChan)
		sub.On("Err").Return((<-chan error)(errChan)).Once()
		sub.On("Unsubscribe").Once()
		rpc.On("SubscribeToHeads", mock.Anything).Return(nil, sub, nil).Once()
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
		opts.rpc.On("SubscribeToHeads", mock.Anything).Return(make(<-chan Head), sub, nil)
		opts.rpc.On("SetAliveLoopSub", sub).Once()
		return newDialedNode(t, opts)
	}
	t.Run("Stays alive and waits for signal", func(t *testing.T) {
		t.Parallel()
		rpc := newMockNodeClient[types.ID, Head](t)
		rpc.On("GetInterceptedChainInfo").Return(ChainInfo{}, ChainInfo{}).Once()
		lggr, observedLogs := logger.TestObserved(t, zap.DebugLevel)
		node := newSubscribedNode(t, testNodeOpts{
			config: testNodeConfig{},
			rpc:    rpc,
			lggr:   lggr,
		})
		defer func() { assert.NoError(t, node.close()) }()
		node.declareAlive()
		tests.AssertLogEventually(t, observedLogs, "Subscription liveness checking disabled")
		tests.AssertLogEventually(t, observedLogs, "Polling disabled")
		assert.Equal(t, nodeStateAlive, node.State())
	})
	t.Run("stays alive while below pollFailureThreshold and resets counter on success", func(t *testing.T) {
		t.Parallel()
		rpc := newMockNodeClient[types.ID, Head](t)
		rpc.On("GetInterceptedChainInfo").Return(ChainInfo{}, ChainInfo{})
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
		rpc.On("GetInterceptedChainInfo").Return(ChainInfo{}, ChainInfo{})
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
		poolInfo := newMockPoolChainInfoProvider(t)
		poolInfo.On("LatestChainInfo").Return(1, ChainInfo{
			BlockNumber: 20,
		}).Once()
		node.SetPoolChainInfoProvider(poolInfo)
		rpc.On("GetInterceptedChainInfo").Return(ChainInfo{BlockNumber: 20}, ChainInfo{BlockNumber: 20})
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
		const mostRecentBlock = 20
		rpc.On("GetInterceptedChainInfo").Return(ChainInfo{BlockNumber: mostRecentBlock}, ChainInfo{BlockNumber: 30})
		poolInfo := newMockPoolChainInfoProvider(t)
		poolInfo.On("LatestChainInfo").Return(10, ChainInfo{
			BlockNumber:     syncThreshold + mostRecentBlock + 1,
			TotalDifficulty: big.NewInt(10),
		}).Once()
		node.SetPoolChainInfoProvider(poolInfo)
		rpc.On("ClientVersion", mock.Anything).Return("", nil)
		// tries to redial in outOfSync
		rpc.On("Dial", mock.Anything).Return(errors.New("failed to dial")).Run(func(_ mock.Arguments) {
			assert.Equal(t, nodeStateOutOfSync, node.State())
		}).Once()
		// disconnects all on transfer to unreachable or outOfSync
		rpc.On("DisconnectAll").Maybe()
		// might be called in unreachable loop
		rpc.On("Dial", mock.Anything).Run(func(_ mock.Arguments) {
			require.Equal(t, nodeStateOutOfSync, node.State())
		}).Return(errors.New("failed to dial")).Maybe()
		node.declareAlive()
		tests.AssertLogEventually(t, observedLogs, "Dial failed: Node is unreachable")
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
		const mostRecentBlock = 20
		rpc.On("GetInterceptedChainInfo").Return(ChainInfo{BlockNumber: mostRecentBlock}, ChainInfo{BlockNumber: 30})
		poolInfo := newMockPoolChainInfoProvider(t)
		poolInfo.On("LatestChainInfo").Return(1, ChainInfo{
			BlockNumber:     syncThreshold + mostRecentBlock + 1,
			TotalDifficulty: big.NewInt(10),
		}).Once()
		node.SetPoolChainInfoProvider(poolInfo)
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
		const mostRecentBlock = 20
		rpc.On("GetInterceptedChainInfo").Return(ChainInfo{BlockNumber: mostRecentBlock}, ChainInfo{BlockNumber: 30})
		rpc.On("ClientVersion", mock.Anything).Return("", nil)
		node.declareAlive()
		tests.AssertLogCountEventually(t, observedLogs, "Version poll successful", 2)
		assert.Equal(t, nodeStateAlive, node.State())
	})
	t.Run("when no new heads received for threshold, transitions to out of sync", func(t *testing.T) {
		t.Parallel()
		rpc := newMockNodeClient[types.ID, Head](t)
		rpc.On("GetInterceptedChainInfo").Return(ChainInfo{}, ChainInfo{}).Once()
		node := newSubscribedNode(t, testNodeOpts{
			config: testNodeConfig{},
			chainConfig: clientMocks.ChainConfig{
				NoNewHeadsThresholdVal: tests.TestInterval,
			},
			rpc: rpc,
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
		rpc.On("GetInterceptedChainInfo").Return(ChainInfo{}, ChainInfo{}).Once()
		lggr, observedLogs := logger.TestObserved(t, zap.DebugLevel)
		node := newSubscribedNode(t, testNodeOpts{
			config: testNodeConfig{},
			lggr:   lggr,
			chainConfig: clientMocks.ChainConfig{
				NoNewHeadsThresholdVal: tests.TestInterval,
			},
			rpc: rpc,
		})
		defer func() { assert.NoError(t, node.close()) }()
		poolInfo := newMockPoolChainInfoProvider(t)
		poolInfo.On("LatestChainInfo").Return(1, ChainInfo{
			BlockNumber:     20,
			TotalDifficulty: big.NewInt(10),
		}).Once()
		node.SetPoolChainInfoProvider(poolInfo)
		node.declareAlive()
		tests.AssertLogEventually(t, observedLogs, fmt.Sprintf("RPC endpoint detected out of sync; %s %s", msgCannotDisable, msgDegradedState))
		assert.Equal(t, nodeStateAlive, node.State())
	})
	newSub := func(t *testing.T) *mocks.Subscription {
		sub := mocks.NewSubscription(t)
		sub.On("Err").Return((<-chan error)(nil))
		sub.On("Unsubscribe").Once()
		return sub
	}
	t.Run("rpc closed head channel", func(t *testing.T) {
		t.Parallel()
		rpc := newMockNodeClient[types.ID, Head](t)
		ch := make(chan Head)
		rpc.On("SubscribeToHeads", mock.Anything).Run(func(args mock.Arguments) {
			close(ch)
		}).Return((<-chan Head)(ch), newSub(t), nil).Once()
		rpc.On("SetAliveLoopSub", mock.Anything).Once()
		rpc.On("GetInterceptedChainInfo").Return(ChainInfo{}, ChainInfo{}).Once()
		lggr, observedLogs := logger.TestObserved(t, zap.ErrorLevel)
		node := newDialedNode(t, testNodeOpts{
			lggr:   lggr,
			config: testNodeConfig{},
			chainConfig: clientMocks.ChainConfig{
				NoNewHeadsThresholdVal: tests.TestInterval,
			},
			rpc: rpc,
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
	t.Run("If finality tag is not enabled updates finalized block metric using finality depth and latest head", func(t *testing.T) {
		t.Parallel()
		rpc := newMockNodeClient[types.ID, Head](t)
		sub := mocks.NewSubscription(t)
		sub.On("Err").Return((<-chan error)(nil))
		sub.On("Unsubscribe").Once()
		const blockNumber = 1000
		const finalityDepth = 10
		const expectedBlock = 990
		rpc.On("GetInterceptedChainInfo").Return(ChainInfo{}, ChainInfo{}).Once()
		ch := make(chan Head)
		rpc.On("SubscribeToHeads", mock.Anything).Run(func(args mock.Arguments) {
			go writeHeads(t, ch, head{BlockNumber: blockNumber - 1}, head{BlockNumber: blockNumber}, head{BlockNumber: blockNumber - 1})
		}).Return((<-chan Head)(ch), sub, nil).Once()
		rpc.On("SetAliveLoopSub", sub).Once()
		name := "node-" + rand.Str(5)
		node := newDialedNode(t, testNodeOpts{
			config:      testNodeConfig{},
			chainConfig: clientMocks.ChainConfig{FinalityDepthVal: finalityDepth},
			rpc:         rpc,
			name:        name,
			chainID:     big.NewInt(1),
		})
		defer func() { assert.NoError(t, node.close()) }()
		node.declareAlive()
		tests.AssertEventually(t, func() bool {
			metric, err := promPoolRPCNodeHighestFinalizedBlock.GetMetricWithLabelValues(big.NewInt(1).String(), name)
			require.NoError(t, err)
			var m = &prom.Metric{}
			require.NoError(t, metric.Write(m))
			return float64(expectedBlock) == m.Gauge.GetValue()
		})
	})
	t.Run("If fails to subscribe to latest finalized blocks, transitions to unreachable ", func(t *testing.T) {
		t.Parallel()
		rpc := newMockNodeClient[types.ID, Head](t)
		expectedError := errors.New("failed to subscribe to finalized heads")
		rpc.On("SubscribeToFinalizedHeads", mock.Anything).Return(nil, mocks.NewSubscription(t), expectedError).Once()
		lggr, observedLogs := logger.TestObserved(t, zap.DebugLevel)
		node := newSubscribedNode(t, testNodeOpts{
			config: testNodeConfig{
				finalizedBlockPollInterval: tests.TestInterval,
			},
			chainConfig: clientMocks.ChainConfig{
				IsFinalityTagEnabled: true,
			},
			rpc:  rpc,
			lggr: lggr,
		})
		defer func() { assert.NoError(t, node.close()) }()
		// disconnects all on transfer to unreachable or outOfSync
		rpc.On("DisconnectAll").Once()
		// might be called in unreachable loop
		rpc.On("Dial", mock.Anything).Return(errors.New("failed to dial")).Maybe()
		node.declareAlive()
		tests.AssertLogEventually(t, observedLogs, "Failed to subscribe to finalized heads")
		tests.AssertEventually(t, func() bool {
			return nodeStateUnreachable == node.State()
		})
	})
	t.Run("Logs warning if latest finalized block is not valid", func(t *testing.T) {
		t.Parallel()
		rpc := newMockNodeClient[types.ID, Head](t)
		ch := make(chan Head, 1)
		head := newMockHead(t)
		head.On("IsValid").Return(false)
		rpc.On("SubscribeToFinalizedHeads", mock.Anything).Run(func(args mock.Arguments) {
			ch <- head
		}).Return((<-chan Head)(ch), newSub(t), nil).Once()
		rpc.On("GetInterceptedChainInfo").Return(ChainInfo{}, ChainInfo{}).Once()
		rpc.On("SubscribeToHeads", mock.Anything).Return(make(<-chan Head), newSub(t), nil).Once()
		rpc.On("SetAliveLoopSub", mock.Anything).Once()
		lggr, observedLogs := logger.TestObserved(t, zap.DebugLevel)
		node := newDialedNode(t, testNodeOpts{
			config: testNodeConfig{},
			chainConfig: clientMocks.ChainConfig{
				IsFinalityTagEnabled: true,
			},
			rpc:  rpc,
			lggr: lggr,
		})
		defer func() { assert.NoError(t, node.close()) }()
		node.declareAlive()
		tests.AssertLogEventually(t, observedLogs, "Latest finalized block is not valid")
	})
	t.Run("On new finalized block updates corresponding metric", func(t *testing.T) {
		t.Parallel()
		rpc := newMockNodeClient[types.ID, Head](t)
		const expectedBlock = 1101
		const finalityDepth = 10
		ch := make(chan Head)
		rpc.On("SubscribeToFinalizedHeads", mock.Anything).Return((<-chan Head)(ch), newSub(t), nil).Once()
		rpc.On("GetInterceptedChainInfo").Return(ChainInfo{}, ChainInfo{}).Once()
		name := "node-" + rand.Str(5)
		node := newSubscribedNode(t, testNodeOpts{
			config: testNodeConfig{},
			chainConfig: clientMocks.ChainConfig{
				FinalityDepthVal:     finalityDepth,
				IsFinalityTagEnabled: true,
			},
			rpc:     rpc,
			name:    name,
			chainID: big.NewInt(1),
		})
		defer func() { assert.NoError(t, node.close()) }()
		node.declareAlive()
		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			defer wg.Done()
			writeHeads(t, ch, head{BlockNumber: expectedBlock - 1}, head{BlockNumber: expectedBlock}, head{BlockNumber: expectedBlock - 1})
		}()
		tests.AssertEventually(t, func() bool {
			metric, err := promPoolRPCNodeHighestFinalizedBlock.GetMetricWithLabelValues(big.NewInt(1).String(), name)
			require.NoError(t, err)
			var m = &prom.Metric{}
			require.NoError(t, metric.Write(m))
			return float64(expectedBlock) == m.Gauge.GetValue()
		})
	})
	t.Run("If finalized heads channel is closed, transitions to unreachable", func(t *testing.T) {
		t.Parallel()
		rpc := newMockNodeClient[types.ID, Head](t)
		rpc.On("GetInterceptedChainInfo").Return(ChainInfo{}, ChainInfo{}).Once()
		ch := make(chan Head)
		close(ch)
		rpc.On("SubscribeToFinalizedHeads", mock.Anything).Return((<-chan Head)(ch), newSub(t), nil).Once()
		lggr, observedLogs := logger.TestObserved(t, zap.DebugLevel)
		node := newSubscribedNode(t, testNodeOpts{
			chainConfig: clientMocks.ChainConfig{
				IsFinalityTagEnabled: true,
			},
			rpc:  rpc,
			lggr: lggr,
		})
		defer func() { assert.NoError(t, node.close()) }()
		// disconnects all on transfer to unreachable or outOfSync
		rpc.On("DisconnectAll").Once()
		// might be called in unreachable loop
		rpc.On("Dial", mock.Anything).Return(errors.New("failed to dial")).Maybe()
		node.declareAlive()
		tests.AssertLogEventually(t, observedLogs, "Finalized heads subscription channel unexpectedly closed")
		tests.AssertEventually(t, func() bool {
			return nodeStateUnreachable == node.State()
		})
	})
	t.Run("when no new finalized heads received for threshold, transitions to out of sync", func(t *testing.T) {
		t.Parallel()
		rpc := newMockNodeClient[types.ID, Head](t)
		rpc.On("GetInterceptedChainInfo").Return(ChainInfo{}, ChainInfo{}).Once()
		ch := make(chan Head, 1)
		ch <- head{BlockNumber: 10}.ToMockHead(t)
		rpc.On("SubscribeToFinalizedHeads", mock.Anything).Return((<-chan Head)(ch), newSub(t), nil).Once()
		lggr, observed := logger.TestObserved(t, zap.DebugLevel)
		noNewFinalizedHeadsThreshold := tests.TestInterval
		node := newSubscribedNode(t, testNodeOpts{
			config: testNodeConfig{},
			chainConfig: clientMocks.ChainConfig{
				NoNewFinalizedHeadsThresholdVal: noNewFinalizedHeadsThreshold,
				IsFinalityTagEnabled:            true,
			},
			rpc:  rpc,
			lggr: lggr,
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
		tests.AssertLogEventually(t, observed, fmt.Sprintf("RPC's finalized state is out of sync; no new finalized heads received for %s (last finalized head received was 10)", noNewFinalizedHeadsThreshold))
		tests.AssertEventually(t, func() bool {
			// right after outOfSync we'll transfer to unreachable due to returned error on Dial
			// we check that we were in out of sync state on first Dial call
			return node.State() == nodeStateUnreachable
		})
	})
	t.Run("when no new finalized heads received for threshold but we are the last live node, forcibly stays alive", func(t *testing.T) {
		t.Parallel()
		rpc := newMockNodeClient[types.ID, Head](t)
		rpc.On("GetInterceptedChainInfo").Return(ChainInfo{}, ChainInfo{}).Once()
		rpc.On("SubscribeToFinalizedHeads", mock.Anything).Return(make(<-chan Head), newSub(t), nil).Once()
		lggr, observed := logger.TestObserved(t, zap.DebugLevel)
		noNewFinalizedHeadsThreshold := tests.TestInterval
		node := newSubscribedNode(t, testNodeOpts{
			config: testNodeConfig{},
			chainConfig: clientMocks.ChainConfig{
				NoNewFinalizedHeadsThresholdVal: noNewFinalizedHeadsThreshold,
				IsFinalityTagEnabled:            true,
			},
			rpc:  rpc,
			lggr: lggr,
		})
		defer func() { assert.NoError(t, node.close()) }()
		poolInfo := newMockPoolChainInfoProvider(t)
		poolInfo.On("LatestChainInfo").Return(1, ChainInfo{
			BlockNumber:     20,
			TotalDifficulty: big.NewInt(10),
		}).Once()
		node.SetPoolChainInfoProvider(poolInfo)
		node.declareAlive()
		tests.AssertLogEventually(t, observed, fmt.Sprintf("RPC's finalized state is out of sync; %s %s", msgCannotDisable, msgDegradedState))
		assert.Equal(t, nodeStateAlive, node.State())
	})
	t.Run("If finalized subscription returns an error, transitions to unreachable", func(t *testing.T) {
		t.Parallel()
		rpc := newMockNodeClient[types.ID, Head](t)
		rpc.On("DisconnectAll").Once()
		rpc.On("GetInterceptedChainInfo").Return(ChainInfo{}, ChainInfo{}).Once()
		sub := mocks.NewSubscription(t)
		errCh := make(chan error, 1)
		errCh <- errors.New("subscription failed")
		sub.On("Err").Return((<-chan error)(errCh))
		sub.On("Unsubscribe").Once()
		rpc.On("SubscribeToFinalizedHeads", mock.Anything).Return((<-chan Head)(nil), sub, nil).Once()
		lggr, observedLogs := logger.TestObserved(t, zap.DebugLevel)
		node := newSubscribedNode(t, testNodeOpts{
			chainConfig: clientMocks.ChainConfig{
				IsFinalityTagEnabled: true,
			},
			rpc:  rpc,
			lggr: lggr,
		})
		defer func() { assert.NoError(t, node.close()) }()
		// disconnects all on transfer to unreachable or outOfSync
		// might be called in unreachable loop
		rpc.On("Dial", mock.Anything).Return(errors.New("failed to dial")).Maybe()
		node.declareAlive()
		tests.AssertLogEventually(t, observedLogs, "Finalized heads subscription was terminated")
		tests.AssertEventually(t, func() bool {
			return nodeStateUnreachable == node.State()
		})
	})
}

type head struct {
	BlockNumber     int64
	BlockDifficulty *big.Int
}

func (h head) ToMockHead(t *testing.T) *mockHead {
	m := newMockHead(t)
	m.On("BlockNumber").Return(h.BlockNumber).Maybe()
	m.On("BlockDifficulty").Return(h.BlockDifficulty).Maybe()
	m.On("IsValid").Return(true).Maybe()
	return m
}

func writeHeads(t *testing.T, ch chan<- Head, heads ...head) {
	for _, head := range heads {
		h := head.ToMockHead(t)
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
	aliveSubscription.On("Err").Return(nil).Maybe()
	aliveSubscription.On("Unsubscribe").Maybe()
	rpc.On("SubscribeToHeads", mock.Anything).Return(make(<-chan Head), aliveSubscription, nil).Maybe()
	rpc.On("SubscribeToFinalizedHeads", mock.Anything).Return(make(<-chan Head), aliveSubscription, nil).Maybe()
	rpc.On("SetAliveLoopSub", mock.Anything).Maybe()
	rpc.On("GetInterceptedChainInfo").Return(ChainInfo{}, ChainInfo{}).Maybe()
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

	t.Run("returns on closed", func(t *testing.T) {
		t.Parallel()
		node := newTestNode(t, testNodeOpts{})
		node.setState(nodeStateClosed)
		node.wg.Add(1)
		node.outOfSyncLoop(syncStatusNotInSyncWithPool)
	})
	t.Run("on old blocks stays outOfSync and returns on close", func(t *testing.T) {
		t.Parallel()
		rpc := newMockNodeClient[types.ID, Head](t)
		nodeChainID := types.RandomID()
		lggr := logger.Test(t)
		node := newAliveNode(t, testNodeOpts{
			rpc:     rpc,
			chainID: nodeChainID,
			lggr:    lggr,
		})
		defer func() { assert.NoError(t, node.close()) }()

		rpc.On("Dial", mock.Anything).Return(nil).Once()
		rpc.On("ChainID", mock.Anything).Return(nodeChainID, nil).Once()
		rpc.On("GetInterceptedChainInfo").Return(ChainInfo{BlockNumber: 0}, ChainInfo{BlockNumber: 13}).Once()

		outOfSyncSubscription := mocks.NewSubscription(t)
		outOfSyncSubscription.On("Err").Return((<-chan error)(nil))
		outOfSyncSubscription.On("Unsubscribe").Once()
		heads := []head{{BlockNumber: 7}, {BlockNumber: 11}, {BlockNumber: 13}}
		ch := make(chan Head)
		var wg sync.WaitGroup
		wg.Add(1)
		rpc.On("SubscribeToHeads", mock.Anything).Run(func(args mock.Arguments) {
			go func() {
				defer wg.Done()
				writeHeads(t, ch, heads...)
			}()
		}).Return((<-chan Head)(ch), outOfSyncSubscription, nil).Once()

		rpc.On("Dial", mock.Anything).Return(errors.New("failed to redial")).Maybe()

		node.declareOutOfSync(syncStatusNoNewHead)
		// wait until all heads are consumed
		wg.Wait()
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
		node.declareOutOfSync(syncStatusNoNewHead)
		tests.AssertEventually(t, func() bool {
			return node.State() == nodeStateUnreachable
		})
	})
	t.Run("if fail to get chainID, transitions to unreachable", func(t *testing.T) {
		t.Parallel()
		rpc := newMockNodeClient[types.ID, Head](t)
		node := newAliveNode(t, testNodeOpts{
			rpc: rpc,
		})
		defer func() { assert.NoError(t, node.close()) }()

		// for out-of-sync
		rpc.On("Dial", mock.Anything).Return(nil).Once()
		// for unreachable
		rpc.On("Dial", mock.Anything).Return(errors.New("failed to dial")).Maybe()
		expectedError := errors.New("failed to get chain ID")
		// might be called multiple times
		rpc.On("ChainID", mock.Anything).Return(types.NewIDFromInt(0), expectedError)
		node.declareOutOfSync(syncStatusNoNewHead)
		tests.AssertEventually(t, func() bool {
			return node.State() == nodeStateUnreachable
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

		// one for out-of-sync & one for invalid chainID
		rpc.On("Dial", mock.Anything).Return(nil).Twice()
		// might be called multiple times
		rpc.On("ChainID", mock.Anything).Return(rpcChainID, nil)
		node.declareOutOfSync(syncStatusNoNewHead)
		tests.AssertEventually(t, func() bool {
			return node.State() == nodeStateInvalidChainID
		})
	})
	t.Run("if syncing, transitions to syncing", func(t *testing.T) {
		t.Parallel()
		rpc := newMockNodeClient[types.ID, Head](t)
		nodeChainID := types.NewIDFromInt(10)
		node := newAliveNode(t, testNodeOpts{
			rpc:     rpc,
			chainID: nodeChainID,
			config:  testNodeConfig{nodeIsSyncingEnabled: true},
		})
		defer func() { assert.NoError(t, node.close()) }()

		rpc.On("Dial", mock.Anything).Return(nil)
		rpc.On("ChainID", mock.Anything).Return(nodeChainID, nil)
		// might be called multiple times
		rpc.On("IsSyncing", mock.Anything).Return(true, nil)
		node.declareOutOfSync(syncStatusNoNewHead)
		tests.AssertEventually(t, func() bool {
			return node.State() == nodeStateSyncing
		})
	})
	t.Run("if fails to fetch syncing status, transitions to unreachable", func(t *testing.T) {
		t.Parallel()
		rpc := newMockNodeClient[types.ID, Head](t)
		nodeChainID := types.NewIDFromInt(10)
		node := newAliveNode(t, testNodeOpts{
			rpc:     rpc,
			chainID: nodeChainID,
			config:  testNodeConfig{nodeIsSyncingEnabled: true},
		})
		defer func() { assert.NoError(t, node.close()) }()

		// one for out-of-sync
		rpc.On("Dial", mock.Anything).Return(nil).Once()
		// for unreachable
		rpc.On("Dial", mock.Anything).Return(errors.New("failed to dial")).Maybe()
		rpc.On("ChainID", mock.Anything).Return(nodeChainID, nil).Once()
		// might be called multiple times
		rpc.On("IsSyncing", mock.Anything).Return(false, errors.New("failed to check syncing"))
		node.declareOutOfSync(syncStatusNoNewHead)
		tests.AssertEventually(t, func() bool {
			return node.State() == nodeStateUnreachable
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
		rpc.On("ChainID", mock.Anything).Return(nodeChainID, nil).Once()
		expectedError := errors.New("failed to subscribe")
		rpc.On("SubscribeToHeads", mock.Anything).Return(nil, nil, expectedError).Once()
		rpc.On("Dial", mock.Anything).Return(errors.New("failed to redial")).Maybe()
		node.declareOutOfSync(syncStatusNoNewHead)
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
		rpc.On("ChainID", mock.Anything).Return(nodeChainID, nil).Once()
		rpc.On("GetInterceptedChainInfo").Return(ChainInfo{}, ChainInfo{}).Once()
		sub := mocks.NewSubscription(t)
		errChan := make(chan error, 1)
		errChan <- errors.New("subscription was terminate")
		sub.On("Err").Return((<-chan error)(errChan))
		sub.On("Unsubscribe").Once()
		rpc.On("SubscribeToHeads", mock.Anything).Return(make(<-chan Head), sub, nil).Once()
		rpc.On("Dial", mock.Anything).Return(errors.New("failed to redial")).Maybe()
		node.declareOutOfSync(syncStatusNoNewHead)
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
		rpc.On("ChainID", mock.Anything).Return(nodeChainID, nil).Once()
		rpc.On("GetInterceptedChainInfo").Return(ChainInfo{}, ChainInfo{}).Once()

		sub := mocks.NewSubscription(t)
		sub.On("Err").Return((<-chan error)(nil))
		sub.On("Unsubscribe").Once()
		ch := make(chan Head)
		rpc.On("SubscribeToHeads", mock.Anything).Run(func(args mock.Arguments) {
			close(ch)
		}).Return((<-chan Head)(ch), sub, nil).Once()
		rpc.On("Dial", mock.Anything).Return(errors.New("failed to redial")).Maybe()
		node.declareOutOfSync(syncStatusNoNewHead)
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
		rpc.On("ChainID", mock.Anything).Return(nodeChainID, nil).Once()

		outOfSyncSubscription := mocks.NewSubscription(t)
		outOfSyncSubscription.On("Err").Return((<-chan error)(nil))
		outOfSyncSubscription.On("Unsubscribe").Once()
		const highestBlock = 1000
		ch := make(chan Head)
		rpc.On("SubscribeToHeads", mock.Anything).Run(func(args mock.Arguments) {
			go writeHeads(t, ch, head{BlockNumber: highestBlock - 1}, head{BlockNumber: highestBlock}, head{BlockNumber: highestBlock + 1})
		}).Return((<-chan Head)(ch), outOfSyncSubscription, nil).Once()
		rpc.On("GetInterceptedChainInfo").Return(ChainInfo{BlockNumber: highestBlock}, ChainInfo{BlockNumber: highestBlock})
		setupRPCForAliveLoop(t, rpc)

		node.declareOutOfSync(syncStatusNoNewHead)
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
			chainConfig: clientMocks.ChainConfig{
				NoNewHeadsThresholdVal: tests.TestInterval,
			},
			rpc:     rpc,
			chainID: nodeChainID,
			lggr:    lggr,
		})
		defer func() { assert.NoError(t, node.close()) }()
		poolInfo := newMockPoolChainInfoProvider(t)
		poolInfo.On("LatestChainInfo").Return(0, ChainInfo{
			BlockNumber:     100,
			TotalDifficulty: big.NewInt(200),
		})
		node.SetPoolChainInfoProvider(poolInfo)
		rpc.On("GetInterceptedChainInfo").Return(ChainInfo{}, ChainInfo{})

		rpc.On("Dial", mock.Anything).Return(nil).Once()
		rpc.On("ChainID", mock.Anything).Return(nodeChainID, nil).Once()

		outOfSyncSubscription := mocks.NewSubscription(t)
		outOfSyncSubscription.On("Err").Return((<-chan error)(nil))
		outOfSyncSubscription.On("Unsubscribe").Once()
		rpc.On("SubscribeToHeads", mock.Anything).Return(make(<-chan Head), outOfSyncSubscription, nil).Once()
		setupRPCForAliveLoop(t, rpc)

		node.declareOutOfSync(syncStatusNoNewHead)
		tests.AssertLogEventually(t, observedLogs, "RPC endpoint is still out of sync, but there are no other available nodes. This RPC node will be forcibly moved back into the live pool in a degraded state")
		tests.AssertEventually(t, func() bool {
			return node.State() == nodeStateAlive
		})
	})
	t.Run("Stays out-of-sync if received new head, but lags behind pool", func(t *testing.T) {
		t.Parallel()
		rpc := newMockNodeClient[types.ID, Head](t)
		nodeChainID := types.RandomID()
		lggr, observedLogs := logger.TestObserved(t, zap.DebugLevel)
		node := newAliveNode(t, testNodeOpts{
			chainConfig: clientMocks.ChainConfig{
				NoNewHeadsThresholdVal: tests.TestInterval,
			},
			config: testNodeConfig{
				syncThreshold: 1,
				selectionMode: NodeSelectionModeHighestHead,
			},
			rpc:     rpc,
			chainID: nodeChainID,
			lggr:    lggr,
		})
		defer func() { assert.NoError(t, node.close()) }()
		poolInfo := newMockPoolChainInfoProvider(t)
		const highestBlock = 20
		poolInfo.On("LatestChainInfo").Return(1, ChainInfo{
			BlockNumber:     highestBlock * 2,
			TotalDifficulty: big.NewInt(200),
		})
		node.SetPoolChainInfoProvider(poolInfo)
		rpc.On("GetInterceptedChainInfo").Return(ChainInfo{}, ChainInfo{BlockNumber: highestBlock})

		rpc.On("Dial", mock.Anything).Return(nil).Once()
		rpc.On("ChainID", mock.Anything).Return(nodeChainID, nil).Once()

		outOfSyncSubscription := mocks.NewSubscription(t)
		outOfSyncSubscription.On("Err").Return((<-chan error)(nil))
		outOfSyncSubscription.On("Unsubscribe").Once()
		ch := make(chan Head)
		rpc.On("SubscribeToHeads", mock.Anything).Run(func(args mock.Arguments) {
			go writeHeads(t, ch, head{BlockNumber: highestBlock - 1}, head{BlockNumber: highestBlock}, head{BlockNumber: highestBlock + 1})
		}).Return((<-chan Head)(ch), outOfSyncSubscription, nil).Once()

		node.declareOutOfSync(syncStatusNoNewHead)
		tests.AssertLogEventually(t, observedLogs, msgReceivedBlock)
		tests.AssertLogEventually(t, observedLogs, "No new heads received for")
		tests.AssertEventually(t, func() bool {
			return node.State() == nodeStateOutOfSync
		})
	})

	// creates RPC mock with all calls necessary to create heads subscription that won't produce any events
	newRPCWithNoOpHeads := func(t *testing.T, chainID types.ID) *mockNodeClient[types.ID, Head] {
		rpc := newMockNodeClient[types.ID, Head](t)
		rpc.On("Dial", mock.Anything).Return(nil).Once()
		rpc.On("ChainID", mock.Anything).Return(chainID, nil).Once()
		sub := mocks.NewSubscription(t)
		sub.On("Err").Return((<-chan error)(nil))
		sub.On("Unsubscribe").Once()
		rpc.On("SubscribeToHeads", mock.Anything).Return(make(<-chan Head), sub, nil).Once()
		return rpc
	}

	t.Run("if fails to subscribe to finalized, becomes unreachable", func(t *testing.T) {
		t.Parallel()
		nodeChainID := types.RandomID()
		rpc := newRPCWithNoOpHeads(t, nodeChainID)
		node := newAliveNode(t, testNodeOpts{
			rpc:     rpc,
			chainID: nodeChainID,
			chainConfig: clientMocks.ChainConfig{
				IsFinalityTagEnabled: true,
			},
		})
		defer func() { assert.NoError(t, node.close()) }()

		rpc.On("SubscribeToFinalizedHeads", mock.Anything).Return((<-chan Head)(nil), nil, errors.New("failed to subscribe")).Once()
		// unreachable
		rpc.On("Dial", mock.Anything).Return(errors.New("failed to redial")).Maybe()

		node.declareOutOfSync(syncStatusNoNewHead)
		tests.AssertEventually(t, func() bool {
			return node.State() == nodeStateUnreachable
		})
	})
	t.Run("on subscription termination becomes unreachable", func(t *testing.T) {
		t.Parallel()
		nodeChainID := types.RandomID()
		rpc := newRPCWithNoOpHeads(t, nodeChainID)
		lggr, observedLogs := logger.TestObserved(t, zap.ErrorLevel)
		node := newAliveNode(t, testNodeOpts{
			rpc:     rpc,
			chainID: nodeChainID,
			lggr:    lggr,
			chainConfig: clientMocks.ChainConfig{
				IsFinalityTagEnabled: true,
			},
		})
		defer func() { assert.NoError(t, node.close()) }()

		sub := mocks.NewSubscription(t)
		errChan := make(chan error, 1)
		errChan <- errors.New("subscription was terminate")
		sub.On("Err").Return((<-chan error)(errChan))
		sub.On("Unsubscribe").Once()
		rpc.On("SubscribeToFinalizedHeads", mock.Anything).Return(make(<-chan Head), sub, nil).Once()
		rpc.On("GetInterceptedChainInfo").Return(ChainInfo{}, ChainInfo{}).Once()
		// unreachable
		rpc.On("Dial", mock.Anything).Return(errors.New("failed to redial")).Maybe()
		node.declareOutOfSync(syncStatusNoNewHead)
		tests.AssertLogEventually(t, observedLogs, "Finalized head subscription was terminated")
		tests.AssertEventually(t, func() bool {
			return node.State() == nodeStateUnreachable
		})
	})
	t.Run("becomes unreachable if head channel is closed", func(t *testing.T) {
		t.Parallel()
		nodeChainID := types.RandomID()
		rpc := newRPCWithNoOpHeads(t, nodeChainID)
		lggr, observedLogs := logger.TestObserved(t, zap.ErrorLevel)
		node := newAliveNode(t, testNodeOpts{
			rpc:     rpc,
			chainID: nodeChainID,
			lggr:    lggr,
			chainConfig: clientMocks.ChainConfig{
				IsFinalityTagEnabled: true,
			},
		})
		defer func() { assert.NoError(t, node.close()) }()

		sub := mocks.NewSubscription(t)
		sub.On("Err").Return((<-chan error)(nil))
		sub.On("Unsubscribe").Once()
		ch := make(chan Head)
		rpc.On("SubscribeToFinalizedHeads", mock.Anything).Run(func(args mock.Arguments) {
			close(ch)
		}).Return((<-chan Head)(ch), sub, nil).Once()
		rpc.On("GetInterceptedChainInfo").Return(ChainInfo{}, ChainInfo{}).Once()
		// unreachable
		rpc.On("Dial", mock.Anything).Return(errors.New("failed to redial")).Maybe()
		node.declareOutOfSync(syncStatusNoNewHead)
		tests.AssertLogEventually(t, observedLogs, "Finalized heads subscription channel unexpectedly closed")
		tests.AssertEventually(t, func() bool {
			return node.State() == nodeStateUnreachable
		})
	})
	t.Run("becomes alive on new finalized block", func(t *testing.T) {
		t.Parallel()
		nodeChainID := types.RandomID()
		rpc := newRPCWithNoOpHeads(t, nodeChainID)
		lggr := logger.Test(t)
		node := newAliveNode(t, testNodeOpts{
			rpc:     rpc,
			chainID: nodeChainID,
			lggr:    lggr,
			chainConfig: clientMocks.ChainConfig{
				IsFinalityTagEnabled:            true,
				NoNewFinalizedHeadsThresholdVal: tests.TestInterval,
			},
		})
		defer func() { assert.NoError(t, node.close()) }()

		const highestBlock = 13
		rpc.On("GetInterceptedChainInfo").Return(ChainInfo{}, ChainInfo{FinalizedBlockNumber: highestBlock}).Once()

		outOfSyncSubscription := mocks.NewSubscription(t)
		outOfSyncSubscription.On("Err").Return((<-chan error)(nil))
		outOfSyncSubscription.On("Unsubscribe").Once()
		ch := make(chan Head)
		rpc.On("SubscribeToFinalizedHeads", mock.Anything).Return((<-chan Head)(ch), outOfSyncSubscription, nil).Once()

		setupRPCForAliveLoop(t, rpc)

		node.declareOutOfSync(syncStatusNoNewFinalizedHead)
		heads := []head{{BlockNumber: highestBlock - 1}, {BlockNumber: highestBlock}}
		writeHeads(t, ch, heads...)
		assert.Equal(t, nodeStateOutOfSync, node.State())
		writeHeads(t, ch, head{BlockNumber: highestBlock + 1})
		tests.AssertEventually(t, func() bool {
			return node.State() == nodeStateAlive
		})
	})
	t.Run("adds finalized block is not increasing flag, if there is no new finalized heads for too long", func(t *testing.T) {
		t.Parallel()
		nodeChainID := types.RandomID()
		rpc := newRPCWithNoOpHeads(t, nodeChainID)
		lggr, observed := logger.TestObserved(t, zap.DebugLevel)
		const noNewFinalizedHeads = tests.TestInterval
		node := newAliveNode(t, testNodeOpts{
			rpc:     rpc,
			chainID: nodeChainID,
			lggr:    lggr,
			chainConfig: clientMocks.ChainConfig{
				IsFinalityTagEnabled:            true,
				NoNewFinalizedHeadsThresholdVal: noNewFinalizedHeads,
			},
		})
		defer func() { assert.NoError(t, node.close()) }()

		const highestBlock = 13
		rpc.On("GetInterceptedChainInfo").Return(ChainInfo{}, ChainInfo{FinalizedBlockNumber: highestBlock}).Once()

		outOfSyncSubscription := mocks.NewSubscription(t)
		outOfSyncSubscription.On("Err").Return((<-chan error)(nil))
		outOfSyncSubscription.On("Unsubscribe").Once()
		ch := make(chan Head)
		rpc.On("SubscribeToFinalizedHeads", mock.Anything).Return((<-chan Head)(ch), outOfSyncSubscription, nil).Once()

		node.declareOutOfSync(syncStatusNotInSyncWithPool)
		heads := []head{{BlockNumber: highestBlock - 1}, {BlockNumber: highestBlock}}
		writeHeads(t, ch, heads...)
		assert.Equal(t, nodeStateOutOfSync, node.State())
		tests.AssertLogEventually(t, observed, fmt.Sprintf("No new finalized heads received for %s. Node stays "+
			"out-of-sync due to sync issues: NotInSyncWithRPCPool,NoNewFinalizedHead", noNewFinalizedHeads))
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
		tests.AssertLogCountEventually(t, observedLogs, "Failed to verify chain ID for node", 2)
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
	t.Run("on syncing status check failure, keeps trying", func(t *testing.T) {
		t.Parallel()
		rpc := newMockNodeClient[types.ID, Head](t)
		nodeChainID := types.RandomID()
		lggr, observedLogs := logger.TestObserved(t, zap.DebugLevel)
		node := newAliveNode(t, testNodeOpts{
			rpc:     rpc,
			chainID: nodeChainID,
			lggr:    lggr,
			config:  testNodeConfig{nodeIsSyncingEnabled: true},
		})
		defer func() { assert.NoError(t, node.close()) }()

		rpc.On("Dial", mock.Anything).Return(nil)
		rpc.On("ChainID", mock.Anything).Run(func(_ mock.Arguments) {
			assert.Equal(t, nodeStateDialed, node.State())
		}).Return(nodeChainID, nil)
		rpc.On("IsSyncing", mock.Anything).Return(false, errors.New("failed to check syncing status"))
		node.declareUnreachable()
		tests.AssertLogCountEventually(t, observedLogs, "Unexpected error while verifying RPC node synchronization status", 2)
	})
	t.Run("on syncing, transitions to syncing state", func(t *testing.T) {
		t.Parallel()
		rpc := newMockNodeClient[types.ID, Head](t)
		nodeChainID := types.RandomID()
		node := newAliveNode(t, testNodeOpts{
			rpc:     rpc,
			chainID: nodeChainID,
			config:  testNodeConfig{nodeIsSyncingEnabled: true},
		})
		defer func() { assert.NoError(t, node.close()) }()

		rpc.On("Dial", mock.Anything).Return(nil)
		rpc.On("ChainID", mock.Anything).Return(nodeChainID, nil)
		rpc.On("IsSyncing", mock.Anything).Return(true, nil)

		setupRPCForAliveLoop(t, rpc)

		node.declareUnreachable()
		tests.AssertEventually(t, func() bool {
			return node.State() == nodeStateSyncing
		})
	})
	t.Run("on successful verification becomes alive", func(t *testing.T) {
		t.Parallel()
		rpc := newMockNodeClient[types.ID, Head](t)
		nodeChainID := types.RandomID()
		node := newAliveNode(t, testNodeOpts{
			rpc:     rpc,
			chainID: nodeChainID,
			config:  testNodeConfig{nodeIsSyncingEnabled: true},
		})
		defer func() { assert.NoError(t, node.close()) }()

		rpc.On("Dial", mock.Anything).Return(nil)
		rpc.On("ChainID", mock.Anything).Return(nodeChainID, nil)
		rpc.On("IsSyncing", mock.Anything).Return(false, nil)

		setupRPCForAliveLoop(t, rpc)

		node.declareUnreachable()
		tests.AssertEventually(t, func() bool {
			return node.State() == nodeStateAlive
		})
	})
	t.Run("on successful verification without isSyncing becomes alive", func(t *testing.T) {
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
	t.Run("on invalid dial becomes unreachable", func(t *testing.T) {
		t.Parallel()
		rpc := newMockNodeClient[types.ID, Head](t)
		nodeChainID := types.RandomID()
		node := newDialedNode(t, testNodeOpts{
			rpc:     rpc,
			chainID: nodeChainID,
		})
		defer func() { assert.NoError(t, node.close()) }()

		rpc.On("Dial", mock.Anything).Return(errors.New("failed to dial"))
		node.declareInvalidChainID()
		tests.AssertEventually(t, func() bool {
			return node.State() == nodeStateUnreachable
		})
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
		// once for chainID and maybe another one for unreachable
		rpc.On("Dial", mock.Anything).Return(nil).Once()
		rpc.On("Dial", mock.Anything).Return(errors.New("failed to dial")).Maybe()
		node.declareInvalidChainID()
		tests.AssertLogEventually(t, observedLogs, "Failed to verify chain ID for node")
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

		rpc.On("Dial", mock.Anything).Return(nil).Once()
		rpc.On("ChainID", mock.Anything).Return(rpcChainID, nil)
		node.declareInvalidChainID()
		tests.AssertLogCountEventually(t, observedLogs, "Failed to verify RPC node; remote endpoint returned the wrong chain ID", 2)
		tests.AssertEventually(t, func() bool {
			return node.State() == nodeStateInvalidChainID
		})
	})
	t.Run("on successful verification without isSyncing becomes alive", func(t *testing.T) {
		t.Parallel()
		rpc := newMockNodeClient[types.ID, Head](t)
		nodeChainID := types.NewIDFromInt(10)
		rpcChainID := types.NewIDFromInt(11)
		node := newDialedNode(t, testNodeOpts{
			rpc:     rpc,
			chainID: nodeChainID,
		})
		defer func() { assert.NoError(t, node.close()) }()

		rpc.On("Dial", mock.Anything).Return(nil).Once()
		rpc.On("ChainID", mock.Anything).Return(rpcChainID, nil).Once()
		rpc.On("ChainID", mock.Anything).Return(nodeChainID, nil).Once()

		setupRPCForAliveLoop(t, rpc)

		node.declareInvalidChainID()
		tests.AssertEventually(t, func() bool {
			return node.State() == nodeStateAlive
		})
	})
	t.Run("on successful verification becomes alive", func(t *testing.T) {
		t.Parallel()
		rpc := newMockNodeClient[types.ID, Head](t)
		nodeChainID := types.NewIDFromInt(10)
		rpcChainID := types.NewIDFromInt(11)
		node := newDialedNode(t, testNodeOpts{
			rpc:     rpc,
			chainID: nodeChainID,
			config:  testNodeConfig{nodeIsSyncingEnabled: true},
		})
		defer func() { assert.NoError(t, node.close()) }()

		rpc.On("Dial", mock.Anything).Return(nil).Once()
		rpc.On("ChainID", mock.Anything).Return(rpcChainID, nil).Once()
		rpc.On("ChainID", mock.Anything).Return(nodeChainID, nil).Once()
		rpc.On("IsSyncing", mock.Anything).Return(false, nil).Once()

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
		tests.AssertLogEventually(t, observedLogs, "Failed to verify chain ID for node")
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
	t.Run("if syncing verification fails, becomes unreachable", func(t *testing.T) {
		t.Parallel()
		rpc := newMockNodeClient[types.ID, Head](t)
		nodeChainID := types.RandomID()
		lggr, observedLogs := logger.TestObserved(t, zap.DebugLevel)
		node := newNode(t, testNodeOpts{
			rpc:     rpc,
			chainID: nodeChainID,
			lggr:    lggr,
			config:  testNodeConfig{nodeIsSyncingEnabled: true},
		})
		defer func() { assert.NoError(t, node.close()) }()

		rpc.On("Dial", mock.Anything).Return(nil).Once()
		rpc.On("ChainID", mock.Anything).Run(func(_ mock.Arguments) {
			assert.Equal(t, nodeStateDialed, node.State())
		}).Return(nodeChainID, nil).Once()
		rpc.On("IsSyncing", mock.Anything).Return(false, errors.New("failed to check syncing status"))
		// disconnects all on transfer to unreachable
		rpc.On("DisconnectAll")
		// fail to redial to stay in unreachable state
		rpc.On("Dial", mock.Anything).Return(errors.New("failed to redial"))
		err := node.Start(tests.Context(t))
		assert.NoError(t, err)
		tests.AssertLogEventually(t, observedLogs, "Unexpected error while verifying RPC node synchronization status")
		tests.AssertEventually(t, func() bool {
			return node.State() == nodeStateUnreachable
		})
	})
	t.Run("on isSyncing transitions to syncing", func(t *testing.T) {
		t.Parallel()
		rpc := newMockNodeClient[types.ID, Head](t)
		nodeChainID := types.NewIDFromInt(10)
		node := newNode(t, testNodeOpts{
			rpc:     rpc,
			chainID: nodeChainID,
			config:  testNodeConfig{nodeIsSyncingEnabled: true},
		})
		defer func() { assert.NoError(t, node.close()) }()

		rpc.On("Dial", mock.Anything).Return(nil)
		rpc.On("ChainID", mock.Anything).Return(nodeChainID, nil)
		rpc.On("IsSyncing", mock.Anything).Return(true, nil)
		// disconnects all on transfer to unreachable
		rpc.On("DisconnectAll")
		err := node.Start(tests.Context(t))
		assert.NoError(t, err)
		tests.AssertEventually(t, func() bool {
			return node.State() == nodeStateSyncing
		})
	})
	t.Run("on successful verification becomes alive", func(t *testing.T) {
		t.Parallel()
		rpc := newMockNodeClient[types.ID, Head](t)
		nodeChainID := types.RandomID()
		node := newNode(t, testNodeOpts{
			rpc:     rpc,
			chainID: nodeChainID,
			config:  testNodeConfig{nodeIsSyncingEnabled: true},
		})
		defer func() { assert.NoError(t, node.close()) }()

		rpc.On("Dial", mock.Anything).Return(nil)
		rpc.On("ChainID", mock.Anything).Return(nodeChainID, nil)
		rpc.On("IsSyncing", mock.Anything).Return(false, nil)

		setupRPCForAliveLoop(t, rpc)

		err := node.Start(tests.Context(t))
		assert.NoError(t, err)
		tests.AssertEventually(t, func() bool {
			return node.State() == nodeStateAlive
		})
	})
	t.Run("on successful verification without isSyncing becomes alive", func(t *testing.T) {
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

func TestUnit_NodeLifecycle_outOfSyncWithPool(t *testing.T) {
	t.Parallel()
	t.Run("skip if nLiveNodes is not configured", func(t *testing.T) {
		node := newTestNode(t, testNodeOpts{})
		outOfSync, liveNodes := node.isOutOfSyncWithPool(ChainInfo{})
		assert.Equal(t, false, outOfSync)
		assert.Equal(t, 0, liveNodes)
	})
	t.Run("skip if syncThreshold is not configured", func(t *testing.T) {
		node := newTestNode(t, testNodeOpts{})
		poolInfo := newMockPoolChainInfoProvider(t)
		node.SetPoolChainInfoProvider(poolInfo)
		outOfSync, liveNodes := node.isOutOfSyncWithPool(ChainInfo{})
		assert.Equal(t, false, outOfSync)
		assert.Equal(t, 0, liveNodes)
	})
	t.Run("panics on invalid selection mode", func(t *testing.T) {
		node := newTestNode(t, testNodeOpts{
			config: testNodeConfig{syncThreshold: 1},
		})
		poolInfo := newMockPoolChainInfoProvider(t)
		poolInfo.On("LatestChainInfo").Return(1, ChainInfo{}).Once()
		node.SetPoolChainInfoProvider(poolInfo)
		assert.Panics(t, func() {
			_, _ = node.isOutOfSyncWithPool(ChainInfo{})
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
			poolInfo := newMockPoolChainInfoProvider(t)
			poolInfo.On("LatestChainInfo").Return(nodesNum, ChainInfo{
				BlockNumber:     highestBlock,
				TotalDifficulty: big.NewInt(totalDifficulty),
			})
			node.SetPoolChainInfoProvider(poolInfo)
			for _, td := range []int64{totalDifficulty - syncThreshold - 1, totalDifficulty - syncThreshold, totalDifficulty, totalDifficulty + 1} {
				for _, testCase := range testCases {
					t.Run(fmt.Sprintf("%s: SelectionModeVal: %s: total difficulty: %d", testCase.name, selectionMode, td), func(t *testing.T) {
						outOfSync, liveNodes := node.isOutOfSyncWithPool(ChainInfo{BlockNumber: testCase.blockNumber, TotalDifficulty: big.NewInt(td)})
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

		poolInfo := newMockPoolChainInfoProvider(t)
		poolInfo.On("LatestChainInfo").Return(nodesNum, ChainInfo{
			BlockNumber:     highestBlock,
			TotalDifficulty: big.NewInt(totalDifficulty),
		})
		node.SetPoolChainInfoProvider(poolInfo)
		for _, hb := range []int64{highestBlock - syncThreshold - 1, highestBlock - syncThreshold, highestBlock, highestBlock + 1} {
			for _, testCase := range testCases {
				t.Run(fmt.Sprintf("%s: SelectionModeVal: %s: highest block: %d", testCase.name, NodeSelectionModeTotalDifficulty, hb), func(t *testing.T) {
					outOfSync, liveNodes := node.isOutOfSyncWithPool(ChainInfo{BlockNumber: hb, TotalDifficulty: big.NewInt(testCase.totalDifficulty)})
					assert.Equal(t, nodesNum, liveNodes)
					assert.Equal(t, testCase.outOfSync, outOfSync)
				})
			}
		}
	})
}

func TestUnit_NodeLifecycle_SyncingLoop(t *testing.T) {
	t.Parallel()
	newDialedNode := func(t *testing.T, opts testNodeOpts) testNode {
		opts.config.nodeIsSyncingEnabled = true
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
		node.syncingLoop()
	})
	t.Run("on invalid dial becomes unreachable", func(t *testing.T) {
		t.Parallel()
		rpc := newMockNodeClient[types.ID, Head](t)
		nodeChainID := types.RandomID()
		node := newDialedNode(t, testNodeOpts{
			rpc:     rpc,
			chainID: nodeChainID,
		})
		defer func() { assert.NoError(t, node.close()) }()

		rpc.On("Dial", mock.Anything).Return(errors.New("failed to dial"))
		node.declareSyncing()
		tests.AssertEventually(t, func() bool {
			return node.State() == nodeStateUnreachable
		})
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
		// once for syncing and maybe another one for unreachable
		rpc.On("Dial", mock.Anything).Return(nil).Once()
		rpc.On("Dial", mock.Anything).Return(errors.New("failed to dial")).Maybe()
		node.declareSyncing()
		tests.AssertLogEventually(t, observedLogs, "Failed to verify chain ID for node")
		tests.AssertEventually(t, func() bool {
			return node.State() == nodeStateUnreachable
		})
	})
	t.Run("on chainID mismatch transitions to invalidChainID", func(t *testing.T) {
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

		rpc.On("Dial", mock.Anything).Return(nil).Twice()
		rpc.On("ChainID", mock.Anything).Return(rpcChainID, nil)
		node.declareSyncing()
		tests.AssertLogCountEventually(t, observedLogs, "Failed to verify RPC node; remote endpoint returned the wrong chain ID", 2)
		tests.AssertEventually(t, func() bool {
			return node.State() == nodeStateInvalidChainID
		})
	})
	t.Run("on failed Syncing check - becomes unreachable", func(t *testing.T) {
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

		rpc.On("ChainID", mock.Anything).Return(nodeChainID, nil).Once()
		// first one is needed to enter internal loop
		rpc.On("IsSyncing", mock.Anything).Return(true, nil).Once()
		rpc.On("IsSyncing", mock.Anything).Return(false, errors.New("failed to check if syncing")).Once()
		rpc.On("Dial", mock.Anything).Return(nil).Once()
		rpc.On("Dial", mock.Anything).Return(errors.New("failed to dial")).Maybe()
		node.declareSyncing()
		tests.AssertLogEventually(t, observedLogs, "Unexpected error while verifying RPC node synchronization status")
		tests.AssertEventually(t, func() bool {
			return node.State() == nodeStateUnreachable
		})
	})
	t.Run("on IsSyncing - keeps trying", func(t *testing.T) {
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

		rpc.On("ChainID", mock.Anything).Return(nodeChainID, nil).Once()
		rpc.On("IsSyncing", mock.Anything).Return(true, nil)
		rpc.On("Dial", mock.Anything).Return(nil).Once()
		node.declareSyncing()
		tests.AssertLogCountEventually(t, observedLogs, "Verification failed: Node is syncing", 2)
		tests.AssertEventually(t, func() bool {
			return node.State() == nodeStateSyncing
		})
	})
	t.Run("on successful verification becomes alive", func(t *testing.T) {
		t.Parallel()
		rpc := newMockNodeClient[types.ID, Head](t)
		nodeChainID := types.RandomID()
		node := newDialedNode(t, testNodeOpts{
			rpc:     rpc,
			chainID: nodeChainID,
		})
		defer func() { assert.NoError(t, node.close()) }()

		rpc.On("Dial", mock.Anything).Return(nil).Once()
		rpc.On("ChainID", mock.Anything).Return(nodeChainID, nil).Once()
		rpc.On("IsSyncing", mock.Anything).Return(true, nil).Once()
		rpc.On("IsSyncing", mock.Anything).Return(false, nil).Once()

		setupRPCForAliveLoop(t, rpc)

		node.declareSyncing()
		tests.AssertEventually(t, func() bool {
			return node.State() == nodeStateAlive
		})
	})
}

func TestNode_State(t *testing.T) {
	t.Run("If not Alive, returns as is", func(t *testing.T) {
		for state := nodeState(0); state < nodeStateLen; state++ {
			if state == nodeStateAlive {
				continue
			}

			node := newTestNode(t, testNodeOpts{})
			node.setState(state)
			assert.Equal(t, state, node.State())
		}
	})
	t.Run("If repeatable read is not enforced, returns alive", func(t *testing.T) {
		node := newTestNode(t, testNodeOpts{})
		node.setState(nodeStateAlive)
		assert.Equal(t, nodeStateAlive, node.State())
	})
	testCases := []struct {
		Name                    string
		FinalizedBlockOffsetVal uint32
		IsFinalityTagEnabled    bool
		PoolChainInfo           ChainInfo
		NodeChainInfo           ChainInfo
		ExpectedState           nodeState
	}{
		{
			Name:                    "If finality lag does not exceeds offset, returns alive (FinalityDepth)",
			FinalizedBlockOffsetVal: 15,
			PoolChainInfo: ChainInfo{
				BlockNumber: 20,
			},
			NodeChainInfo: ChainInfo{
				BlockNumber: 5,
			},
			ExpectedState: nodeStateAlive,
		},
		{
			Name:                    "If finality lag does not exceeds offset, returns alive (FinalityTag)",
			FinalizedBlockOffsetVal: 15,
			IsFinalityTagEnabled:    true,
			PoolChainInfo: ChainInfo{
				FinalizedBlockNumber: 20,
			},
			NodeChainInfo: ChainInfo{
				FinalizedBlockNumber: 5,
			},
			ExpectedState: nodeStateAlive,
		},
		{
			Name:                    "If finality lag exceeds offset, returns nodeStateFinalizedBlockOutOfSync (FinalityDepth)",
			FinalizedBlockOffsetVal: 15,
			PoolChainInfo: ChainInfo{
				BlockNumber: 20,
			},
			NodeChainInfo: ChainInfo{
				BlockNumber: 4,
			},
			ExpectedState: nodeStateFinalizedBlockOutOfSync,
		},
		{
			Name:                    "If finality lag exceeds offset, returns nodeStateFinalizedBlockOutOfSync (FinalityTag)",
			FinalizedBlockOffsetVal: 15,
			IsFinalityTagEnabled:    true,
			PoolChainInfo: ChainInfo{
				FinalizedBlockNumber: 20,
			},
			NodeChainInfo: ChainInfo{
				FinalizedBlockNumber: 4,
			},
			ExpectedState: nodeStateFinalizedBlockOutOfSync,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			rpc := newMockNodeClient[types.ID, Head](t)
			rpc.On("GetInterceptedChainInfo").Return(tc.NodeChainInfo, tc.PoolChainInfo).Once()
			node := newTestNode(t, testNodeOpts{
				config: testNodeConfig{
					enforceRepeatableRead: true,
				},
				chainConfig: clientMocks.ChainConfig{
					FinalizedBlockOffsetVal: tc.FinalizedBlockOffsetVal,
					IsFinalityTagEnabled:    tc.IsFinalityTagEnabled,
				},
				rpc: rpc,
			})
			poolInfo := newMockPoolChainInfoProvider(t)
			poolInfo.On("HighestUserObservations").Return(tc.PoolChainInfo).Once()
			node.SetPoolChainInfoProvider(poolInfo)
			node.setState(nodeStateAlive)
			assert.Equal(t, tc.ExpectedState, node.State())
		})
	}
}
