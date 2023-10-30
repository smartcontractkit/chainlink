package client

import (
	"fmt"
	"sync/atomic"
	"testing"

	"github.com/cometbft/cometbft/libs/rand"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"

	"github.com/smartcontractkit/chainlink/v2/common/types"
	"github.com/smartcontractkit/chainlink/v2/common/types/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

func TestUnit_NodeLifecycle_aliveLoop(t *testing.T) {
	t.Parallel()

	newDialedNode := func(t *testing.T, opts testNodeOpts) testNode {
		node := newTestNode(t, opts)
		opts.rpc.On("Close").Return(nil).Once()

		t.Cleanup(func() {
			assert.NoError(t, node.close())
		})
		node.setState(nodeStateDialed)
		return node
	}

	t.Run("if initial subscribe fails, transitions to unreachable", func(t *testing.T) {
		t.Parallel()
		rpc := newMockNodeClient[types.ID, Head](t)
		node := newDialedNode(t, testNodeOpts{
			rpc: rpc,
		})

		expectedError := errors.New("failed to subscribe to rpc")
		rpc.On("Subscribe", mock.Anything, mock.Anything, rpcSubscriptionMethodNewHeads).Return(nil, expectedError).Once()
		rpc.On("DisconnectAll").Once()
		// might be called in unreachable loop
		rpc.On("Dial", mock.Anything).Return(errors.New("failed to dial")).Maybe()
		node.declareAlive()
		testutils.AssertEventually(t, func() bool {
			return node.State() == nodeStateUnreachable
		})
	})
	t.Run("if remote RPC connection is closed transitions to unreachable", func(t *testing.T) {
		t.Parallel()
		rpc := newMockNodeClient[types.ID, Head](t)

		lggr, observedLogs := logger.TestLoggerObserved(t, zap.WarnLevel)
		node := newDialedNode(t, testNodeOpts{
			rpc:  rpc,
			lggr: lggr,
		})

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
		testutils.WaitForLogMessage(t, observedLogs, "Subscription was terminated")
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
	t.Run("stays alive while below pollFailureThreshold and resets counter on success", func(t *testing.T) {
		t.Parallel()
		rpc := newMockNodeClient[types.ID, Head](t)
		lggr, observedLogs := logger.TestLoggerObserved(t, zap.DebugLevel)
		const pollFailureThreshold = 3
		node := newSubscribedNode(t, testNodeOpts{
			config: testNodeConfig{
				pollFailureThreshold: pollFailureThreshold,
				pollInterval:         testutils.TestInterval,
			},
			rpc:  rpc,
			lggr: lggr,
		})

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
		testutils.WaitForLogMessageCount(t, observedLogs, fmt.Sprintf("Poll failure, RPC endpoint %s failed to respond properly", node.String()), pollFailureThreshold)
		testutils.WaitForLogMessageCount(t, observedLogs, "Version poll successful", 2)
		assert.True(t, ensuredAlive.Load(), "expected to ensure that node was alive")

	})
	t.Run("becomes unreachable when exceeds pollFailureThreshold", func(t *testing.T) {
		t.Parallel()
		rpc := newMockNodeClient[types.ID, Head](t)
		lggr, observedLogs := logger.TestLoggerObserved(t, zap.DebugLevel)
		const pollFailureThreshold = 3
		node := newSubscribedNode(t, testNodeOpts{
			config: testNodeConfig{
				pollFailureThreshold: pollFailureThreshold,
				pollInterval:         testutils.TestInterval,
			},
			rpc:  rpc,
			lggr: lggr,
		})
		pollError := errors.New("failed to get ClientVersion")
		rpc.On("ClientVersion", mock.Anything).Return("", pollError)
		// disconnects all on transfer to unreachable
		rpc.On("DisconnectAll").Once()
		// might be called in unreachable loop
		rpc.On("Dial", mock.Anything).Return(errors.New("failed to dial")).Maybe()
		node.declareAlive()
		testutils.WaitForLogMessageCount(t, observedLogs, fmt.Sprintf("Poll failure, RPC endpoint %s failed to respond properly", node.String()), pollFailureThreshold)
		assert.Equal(t, nodeStateUnreachable, node.State())
	})
	t.Run("stays alive even, when exceeds pollFailureThreshold because it's last node", func(t *testing.T) {
		t.Parallel()
		rpc := newMockNodeClient[types.ID, Head](t)
		lggr, observedLogs := logger.TestLoggerObserved(t, zap.DebugLevel)
		const pollFailureThreshold = 3
		node := newSubscribedNode(t, testNodeOpts{
			config: testNodeConfig{
				pollFailureThreshold: pollFailureThreshold,
				pollInterval:         testutils.TestInterval,
			},
			rpc:  rpc,
			lggr: lggr,
		})
		node.nLiveNodes = func() (count int, blockNumber int64, totalDifficulty *utils.Big) {
			return 1, 20, utils.NewBigI(10)
		}
		pollError := errors.New("failed to get ClientVersion")
		rpc.On("ClientVersion", mock.Anything).Return("", pollError)
		node.declareAlive()
		testutils.WaitForLogMessage(t, observedLogs, fmt.Sprintf("RPC endpoint failed to respond to %d consecutive polls", pollFailureThreshold))
		assert.Equal(t, nodeStateAlive, node.State())
	})
	t.Run("outOfSync when falls behind", func(t *testing.T) {
		rpc := newMockNodeClient[types.ID, Head](t)
		lggr, observedLogs := logger.TestLoggerObserved(t, zap.DebugLevel)
		const syncThreshold = 10
		node := newSubscribedNode(t, testNodeOpts{
			config: testNodeConfig{
				pollInterval:  testutils.TestInterval,
				syncThreshold: syncThreshold,
				selectionMode: NodeSelectionMode_RoundRobin,
			},
			rpc:  rpc,
			lggr: lggr,
		})
		node.stateLatestBlockNumber = 20
		node.nLiveNodes = func() (count int, blockNumber int64, totalDifficulty *utils.Big) {
			return 10, syncThreshold + node.stateLatestBlockNumber + 1, utils.NewBigI(10)
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
		testutils.WaitForLogMessage(t, observedLogs, "RPC endpoint has fallen behind")
	})
	t.Run("stays alive even when falls behind", func(t *testing.T) {
		rpc := newMockNodeClient[types.ID, Head](t)
		lggr, observedLogs := logger.TestLoggerObserved(t, zap.DebugLevel)
		const syncThreshold = 10
		node := newSubscribedNode(t, testNodeOpts{
			config: testNodeConfig{
				pollInterval:  testutils.TestInterval,
				syncThreshold: syncThreshold,
				selectionMode: NodeSelectionMode_RoundRobin,
			},
			rpc:  rpc,
			lggr: lggr,
		})
		node.stateLatestBlockNumber = 20
		node.nLiveNodes = func() (count int, blockNumber int64, totalDifficulty *utils.Big) {
			return 1, syncThreshold + node.stateLatestBlockNumber + 1, utils.NewBigI(10)
		}
		rpc.On("ClientVersion", mock.Anything).Return("", nil)
		node.declareAlive()
		testutils.WaitForLogMessage(t, observedLogs, fmt.Sprintf("RPC endpoint has fallen behind; %s %s", msgCannotDisable, msgDegradedState))
	})
	t.Run("when no new heads received for threshold, transitions to out of sync", func(t *testing.T) {
		rpc := newMockNodeClient[types.ID, Head](t)
		node := newSubscribedNode(t, testNodeOpts{
			config:              testNodeConfig{},
			noNewHeadsThreshold: testutils.TestInterval,
			rpc:                 rpc,
		})
		// tries to redial in outOfSync
		rpc.On("Dial", mock.Anything).Return(errors.New("failed to dial")).Run(func(_ mock.Arguments) {
			assert.Equal(t, nodeStateOutOfSync, node.State())
		}).Once()
		// disconnects all on transfer to unreachable or outOfSync
		rpc.On("DisconnectAll").Maybe()
		// might be called in unreachable loop
		rpc.On("Dial", mock.Anything).Return(errors.New("failed to dial")).Maybe()
		node.declareAlive()
		testutils.AssertEventually(t, func() bool {
			// right after outOfSync we'll transfer to unreachable due to returned error on Dial
			return node.State() == nodeStateOutOfSync || node.State() == nodeStateUnreachable
		})
	})
	t.Run("when no new heads received for threshold and no nodes stay alive", func(t *testing.T) {
		rpc := newMockNodeClient[types.ID, Head](t)
		lggr, observedLogs := logger.TestLoggerObserved(t, zap.DebugLevel)
		node := newSubscribedNode(t, testNodeOpts{
			config:              testNodeConfig{},
			lggr:                lggr,
			noNewHeadsThreshold: testutils.TestInterval,
			rpc:                 rpc,
		})
		node.nLiveNodes = func() (count int, blockNumber int64, totalDifficulty *utils.Big) {
			return 1, 20, utils.NewBigI(10)
		}
		node.declareAlive()
		testutils.WaitForLogMessage(t, observedLogs, fmt.Sprintf("RPC endpoint detected out of sync; %s %s", msgCannotDisable, msgDegradedState))
		assert.Equal(t, nodeStateAlive, node.State())
	})

	t.Run("rpc closed head channel", func(t *testing.T) {
		rpc := newMockNodeClient[types.ID, Head](t)
		sub := mocks.NewSubscription(t)
		lggr, observedLogs := logger.TestLoggerObserved(t, zap.ErrorLevel)
		sub.On("Err").Return((<-chan error)(nil))
		sub.On("Unsubscribe").Once()
		rpc.On("Subscribe", mock.Anything, mock.Anything, rpcSubscriptionMethodNewHeads).Run(func(args mock.Arguments) {
			ch := args.Get(1).(chan<- Head)
			close(ch)
		}).Return(sub, nil).Once()
		rpc.On("SetAliveLoopSub", sub).Once()
		node := newDialedNode(t, testNodeOpts{
			lggr:                lggr,
			config:              testNodeConfig{},
			noNewHeadsThreshold: testutils.TestInterval,
			rpc:                 rpc,
		})
		// disconnects all on transfer to unreachable or outOfSync
		rpc.On("DisconnectAll").Once()
		// might be called in unreachable loop
		rpc.On("Dial", mock.Anything).Return(errors.New("failed to dial")).Maybe()
		node.declareAlive()
		testutils.WaitForLogMessage(t, observedLogs, "Subscription channel unexpectedly closed")
		assert.Equal(t, nodeStateUnreachable, node.State())

	})
	t.Run("updates block number and difficulty on new head", func(t *testing.T) {
		rpc := newMockNodeClient[types.ID, Head](t)
		sub := mocks.NewSubscription(t)
		sub.On("Err").Return((<-chan error)(nil))
		sub.On("Unsubscribe").Once()
		expectedBlockNumber := rand.Int64()
		expectedDiff := utils.NewBigI(rand.Int64())
		rpc.On("Subscribe", mock.Anything, mock.Anything, rpcSubscriptionMethodNewHeads).Run(func(args mock.Arguments) {
			ch := args.Get(1).(chan<- Head)
			go func() {
				head := newMockHead(t)
				head.On("BlockNumber").Return(expectedBlockNumber)
				head.On("BlockDifficulty").Return(expectedDiff)
				ch <- head
			}()
		}).Return(sub, nil).Once()
		rpc.On("SetAliveLoopSub", sub).Once()
		node := newDialedNode(t, testNodeOpts{
			config: testNodeConfig{},
			rpc:    rpc,
		})
		node.declareAlive()
		testutils.AssertEventually(t, func() bool {
			state, block, diff := node.StateAndLatest()
			return state == nodeStateAlive && block == expectedBlockNumber == diff.Equal(expectedDiff)
		})

	})
}

/*
func TestUnit_NodeLifecycle_outOfSyncLoop(t *testing.T) {
	t.Parallel()

	t.Run("exits on close", func(t *testing.T) {
		cfg := TestNodePoolConfig{}
		n := newTestNode(t, cfg, time.Second*0)
		dial(t, n)
		n.setState(NodeStateOutOfSync)

		ch := make(chan struct{})

		n.wg.Add(1)
		go func() {
			defer close(ch)
			n.aliveLoop()
		}()
		assert.NoError(t, n.Close())
		testutils.WaitWithTimeout(t, ch, "expected outOfSyncLoop to exit")
	})

	t.Run("if initial subscribe fails, transitions to unreachable", func(t *testing.T) {
		cfg := TestNodePoolConfig{}
		n := newTestNodeWithCallback(t, cfg, time.Second*0, func(string, gjson.Result) (resp testutils.JSONRPCResponse) { return })
		dial(t, n)
		n.setState(NodeStateOutOfSync)
		defer func() { assert.NoError(t, n.Close()) }()

		n.wg.Add(1)

		n.outOfSyncLoop(func(num int64, td *utils.Big) bool { return num == 0 })
		assert.Equal(t, NodeStateUnreachable, n.State())
	})

	t.Run("transitions to unreachable if remote RPC subscription channel closed", func(t *testing.T) {
		cfg := TestNodePoolConfig{}
		chSubbed := make(chan struct{}, 1)
		s := testutils.NewWSServer(t, testutils.FixtureChainID,
			func(method string, params gjson.Result) (resp testutils.JSONRPCResponse) {
				switch method {
				case "eth_subscribe":
					chSubbed <- struct{}{}
					resp.Result = `"0x00"`
					resp.Notify = makeHeadResult(0)
					return
				default:
					t.Errorf("unexpected RPC method: %s", method)
				}
				return
			})

		iN := NewNode(cfg, time.Duration(time.Second), logger.TestLogger(t), *s.WSURL(), nil, "test node", 42, testutils.FixtureChainID, 1)
		n := iN.(*node)

		dial(t, n)
		n.setState(NodeStateOutOfSync)
		defer func() { assert.NoError(t, n.Close()) }()

		n.wg.Add(1)
		go n.outOfSyncLoop(func(num int64, td *utils.Big) bool { return num == 0 })

		testutils.WaitWithTimeout(t, chSubbed, "timed out waiting for initial subscription")

		assert.Equal(t, NodeStateOutOfSync, n.State())

		// Simulate remote websocket disconnect
		// This causes sub.Err() to close
		s.Close()

		testutils.AssertEventually(t, func() bool {
			return n.State() == NodeStateUnreachable
		})
	})

	t.Run("transitions to alive if it receives a newer head", func(t *testing.T) {
		// NoNewHeadsThreshold needs to be positive but must be very large so
		// we don't time out waiting for a new head before we have a chance to
		// handle the server disconnect
		lggr, observedLogs := logger.TestLoggerObserved(t, zap.DebugLevel)
		cfg := TestNodePoolConfig{}
		chSubbed := make(chan struct{}, 1)
		s := testutils.NewWSServer(t, testutils.FixtureChainID,
			func(method string, params gjson.Result) (resp testutils.JSONRPCResponse) {
				switch method {
				case "eth_subscribe":
					chSubbed <- struct{}{}
					resp.Result = `"0x00"`
					resp.Notify = makeNewHeadWSMessage(42)
					return
				case "eth_unsubscribe":
					resp.Result = "true"
					return
				default:
					t.Errorf("unexpected RPC method: %s", method)
				}
				return
			})

		iN := NewNode(cfg, time.Second*0, lggr, *s.WSURL(), nil, "test node", 0, testutils.FixtureChainID, 1)
		n := iN.(*node)

		start(t, n)
		n.setState(NodeStateOutOfSync)
		defer func() { assert.NoError(t, n.Close()) }()

		n.wg.Add(1)
		go n.outOfSyncLoop(func(num int64, td *utils.Big) bool { return num < 43 })

		testutils.WaitWithTimeout(t, chSubbed, "timed out waiting for initial subscription")

		assert.Equal(t, NodeStateOutOfSync, n.State())

		// heads less than latest seen head are ignored; they do not make the node live
		for i := 0; i < 43; i++ {
			msg := makeNewHeadWSMessage(i)
			s.MustWriteBinaryMessageSync(t, msg)
			testutils.WaitForLogMessageCount(t, observedLogs, msgReceivedBlock, i+1)
			assert.Equal(t, NodeStateOutOfSync, n.State())
		}

		msg := makeNewHeadWSMessage(43)
		s.MustWriteBinaryMessageSync(t, msg)

		testutils.AssertEventually(t, func() bool {
			s, n, td := n.StateAndLatest()
			return s == NodeStateAlive && n != -1 && td != nil
		})

		testutils.WaitForLogMessage(t, observedLogs, msgInSync)
	})

	t.Run("transitions to alive if back in-sync", func(t *testing.T) {
		lggr, observedLogs := logger.TestLoggerObserved(t, zap.DebugLevel)
		cfg := TestNodePoolConfig{NodeSyncThreshold: 5, NodeSelectionMode: NodeSelectionMode_HighestHead}
		chSubbed := make(chan struct{}, 1)
		const stall = 42
		s := testutils.NewWSServer(t, testutils.FixtureChainID,
			func(method string, params gjson.Result) (resp testutils.JSONRPCResponse) {
				switch method {
				case "eth_subscribe":
					chSubbed <- struct{}{}
					resp.Result = `"0x00"`
					resp.Notify = makeNewHeadWSMessage(stall)
					return
				case "eth_unsubscribe":
					resp.Result = "true"
					return
				default:
					t.Errorf("unexpected RPC method: %s", method)
				}
				return
			})

		iN := NewNode(cfg, time.Second*0, lggr, *s.WSURL(), nil, "test node", 0, testutils.FixtureChainID, 1)
		n := iN.(*node)
		n.nLiveNodes = func() (count int, blockNumber int64, totalDifficulty *utils.Big) {
			return 2, stall + int64(cfg.SyncThreshold()), nil
		}

		start(t, n)
		n.setState(NodeStateOutOfSync)
		defer func() { assert.NoError(t, n.Close()) }()

		n.wg.Add(1)
		go n.outOfSyncLoop(n.isOutOfSync)

		testutils.WaitWithTimeout(t, chSubbed, "timed out waiting for initial subscription")

		assert.Equal(t, NodeStateOutOfSync, n.State())

		// heads less than stall (latest seen head - SyncThreshold) are ignored; they do not make the node live
		for i := 0; i < stall; i++ {
			msg := makeNewHeadWSMessage(i)
			s.MustWriteBinaryMessageSync(t, msg)
			testutils.WaitForLogMessageCount(t, observedLogs, msgReceivedBlock, i+1)
			assert.Equal(t, NodeStateOutOfSync, n.State())
		}

		msg := makeNewHeadWSMessage(stall)
		s.MustWriteBinaryMessageSync(t, msg)

		testutils.AssertEventually(t, func() bool {
			s, n, td := n.StateAndLatest()
			return s == NodeStateAlive && n != -1 && td != nil
		})

		testutils.WaitForLogMessage(t, observedLogs, msgInSync)
	})

	t.Run("if no live nodes are available, forcibly marks this one alive again", func(t *testing.T) {
		cfg := TestNodePoolConfig{}
		chSubbed := make(chan struct{}, 1)
		s := testutils.NewWSServer(t, testutils.FixtureChainID,
			func(method string, params gjson.Result) (resp testutils.JSONRPCResponse) {
				switch method {
				case "eth_subscribe":
					chSubbed <- struct{}{}
					resp.Result = `"0x00"`
					resp.Notify = makeHeadResult(0)
					return
				case "eth_unsubscribe":
					resp.Result = "true"
					return
				default:
					t.Errorf("unexpected RPC method: %s", method)
				}
				return
			})

		iN := NewNode(cfg, testutils.TestInterval, logger.TestLogger(t), *s.WSURL(), nil, "test node", 42, testutils.FixtureChainID, 1)
		n := iN.(*node)
		n.nLiveNodes = func() (int, int64, *utils.Big) { return 0, 0, nil }

		dial(t, n)
		n.setState(NodeStateOutOfSync)
		defer func() { assert.NoError(t, n.Close()) }()

		n.wg.Add(1)
		go n.outOfSyncLoop(func(num int64, td *utils.Big) bool { return num == 0 })

		testutils.WaitWithTimeout(t, chSubbed, "timed out waiting for initial subscription")

		testutils.AssertEventually(t, func() bool {
			return n.State() == NodeStateAlive
		})
	})
}

func TestUnit_NodeLifecycle_unreachableLoop(t *testing.T) {
	t.Parallel()

	t.Run("exits on close", func(t *testing.T) {
		cfg := TestNodePoolConfig{}
		n := newTestNode(t, cfg, time.Second*0)
		start(t, n)
		n.setState(NodeStateUnreachable)

		ch := make(chan struct{})
		n.wg.Add(1)
		go func() {
			n.unreachableLoop()
			close(ch)
		}()
		assert.NoError(t, n.Close())
		testutils.WaitWithTimeout(t, ch, "expected unreachableLoop to exit")
	})

	t.Run("on successful redial and verify, transitions to alive", func(t *testing.T) {
		cfg := TestNodePoolConfig{}
		n := newTestNode(t, cfg, time.Second*0)
		start(t, n)
		defer func() { assert.NoError(t, n.Close()) }()
		n.setState(NodeStateUnreachable)
		n.wg.Add(1)

		go n.unreachableLoop()

		testutils.AssertEventually(t, func() bool {
			return n.State() == NodeStateAlive
		})
	})

	t.Run("on successful redial but failed verify, transitions to invalid chain ID", func(t *testing.T) {
		cfg := TestNodePoolConfig{}
		s := testutils.NewWSServer(t, testutils.FixtureChainID, standardHandler)
		lggr, observedLogs := logger.TestLoggerObserved(t, zap.ErrorLevel)
		iN := NewNode(cfg, time.Second*0, lggr, *s.WSURL(), nil, "test node", 0, big.NewInt(42), 1)
		n := iN.(*node)
		defer func() { assert.NoError(t, n.Close()) }()
		start(t, n)
		n.setState(NodeStateUnreachable)
		n.wg.Add(1)

		go n.unreachableLoop()

		testutils.WaitForLogMessage(t, observedLogs, "Failed to redial RPC node; remote endpoint returned the wrong chain ID")

		testutils.AssertEventually(t, func() bool {
			return n.State() == NodeStateInvalidChainID
		})
	})

	t.Run("on failed redial, keeps trying to redial", func(t *testing.T) {
		cfg := TestNodePoolConfig{}
		lggr, observedLogs := logger.TestLoggerObserved(t, zap.DebugLevel)
		iN := NewNode(cfg, time.Second*0, lggr, *testutils.MustParseURL(t, "ws://test.invalid"), nil, "test node", 0, big.NewInt(42), 1)
		n := iN.(*node)
		defer func() { assert.NoError(t, n.Close()) }()
		start(t, n)
		n.setState(NodeStateUnreachable)
		n.wg.Add(1)

		go n.unreachableLoop()

		testutils.WaitForLogMessageCount(t, observedLogs, "Failed to redial RPC node", 3)

		assert.Equal(t, NodeStateUnreachable, n.State())
	})
}
func TestUnit_NodeLifecycle_invalidChainIDLoop(t *testing.T) {
	t.Parallel()

	t.Run("exits on close", func(t *testing.T) {
		cfg := TestNodePoolConfig{}
		n := newTestNode(t, cfg, time.Second*0)
		start(t, n)
		n.setState(NodeStateInvalidChainID)

		ch := make(chan struct{})
		n.wg.Add(1)
		go func() {
			n.invalidChainIDLoop()
			close(ch)
		}()
		assert.NoError(t, n.Close())
		testutils.WaitWithTimeout(t, ch, "expected invalidChainIDLoop to exit")
	})

	t.Run("on successful verify, transitions to alive", func(t *testing.T) {
		cfg := TestNodePoolConfig{}
		n := newTestNode(t, cfg, time.Second*0)
		dial(t, n)
		defer func() { assert.NoError(t, n.Close()) }()
		n.setState(NodeStateInvalidChainID)
		n.wg.Add(1)

		go n.invalidChainIDLoop()

		testutils.AssertEventually(t, func() bool {
			return n.State() == NodeStateAlive
		})
	})

	t.Run("on failed verify, keeps checking", func(t *testing.T) {
		cfg := TestNodePoolConfig{}
		s := testutils.NewWSServer(t, testutils.FixtureChainID, standardHandler)
		lggr, observedLogs := logger.TestLoggerObserved(t, zap.ErrorLevel)
		iN := NewNode(cfg, time.Second*0, lggr, *s.WSURL(), nil, "test node", 0, big.NewInt(42), 1)
		n := iN.(*node)
		defer func() { assert.NoError(t, n.Close()) }()
		dial(t, n)
		n.setState(NodeStateUnreachable)
		n.wg.Add(1)

		go n.unreachableLoop()

		testutils.WaitForLogMessageCount(t, observedLogs, "Failed to redial RPC node; remote endpoint returned the wrong chain ID", 3)

		assert.Equal(t, NodeStateInvalidChainID, n.State())
	})
}*/
