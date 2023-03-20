package client

import (
	"fmt"
	"math/big"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
	"go.uber.org/zap"

	evmtypes "github.com/smartcontractkit/chainlink/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/utils"
)

func standardHandler(method string, _ gjson.Result) (resp testutils.JSONRPCResponse) {
	if method == "eth_subscribe" {
		resp.Result = `"0x00"`
		resp.Notify = HeadResult
		return
	}
	return
}

func newTestNode(t *testing.T, cfg NodeConfig) *node {
	return newTestNodeWithCallback(t, cfg, standardHandler)
}

func newTestNodeWithCallback(t *testing.T, cfg NodeConfig, callback testutils.JSONRPCHandler) *node {
	s := testutils.NewWSServer(t, testutils.FixtureChainID, callback)
	iN := NewNode(cfg, logger.TestLogger(t), *s.WSURL(), nil, "test node", 42, testutils.FixtureChainID)
	n := iN.(*node)
	return n
}

// dial sets up the node and puts it into the live state, bypassing the
// normal Start() method which would fire off unwanted goroutines
func dial(t *testing.T, n *node) {
	ctx := testutils.Context(t)
	require.NoError(t, n.dial(ctx))
	n.setState(NodeStateAlive)
	start(t, n)
}

func start(t *testing.T, n *node) {
	// must start to allow closing
	err := n.StartOnce("test node", func() error { return nil })
	assert.NoError(t, err)
}

func makeHeadResult(n int) string {
	return fmt.Sprintf(
		`{"difficulty":"0xf3a00","extraData":"0xd883010503846765746887676f312e372e318664617277696e","gasLimit":"0xffc001","gasUsed":"0x0","hash":"0x41800b5c3f1717687d85fc9018faac0a6e90b39deaa0b99e7fe4fe796ddeb26a","logsBloom":"0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000","miner":"0xd1aeb42885a43b72b518182ef893125814811048","mixHash":"0x0f98b15f1a4901a7e9204f3c500a7bd527b3fb2c3340e12176a44b83e414a69e","nonce":"0x0ece08ea8c49dfd9","number":"%s","parentHash":"0x41941023680923e0fe4d74a34bdac8141f2540e3ae90623718e47d66d1ca4a2d","receiptsRoot":"0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421","sha3Uncles":"0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347","size":"0x218","stateRoot":"0xc7b01007a10da045eacb90385887dd0c38fcb5db7393006bdde24b93873c334b","timestamp":"0x58318da2","totalDifficulty":"0x1f3a00","transactions":[],"transactionsRoot":"0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421","uncles":[]}`,
		testutils.IntToHex(n),
	)
}

func makeNewHeadWSMessage(n int) string {
	return fmt.Sprintf(`{"jsonrpc":"2.0","method":"eth_subscription","params":{"subscription":"0x00","result":%s}}`, makeHeadResult(n))
}

func TestUnit_NodeLifecycle_aliveLoop(t *testing.T) {
	t.Parallel()

	t.Run("with no poll and sync timeouts, exits on close", func(t *testing.T) {
		pollAndSyncTimeoutsDisabledCfg := TestNodeConfig{}
		n := newTestNode(t, pollAndSyncTimeoutsDisabledCfg)
		dial(t, n)

		ch := make(chan struct{})
		n.wg.Add(1)
		go func() {
			defer close(ch)
			n.aliveLoop()
		}()
		assert.NoError(t, n.Close())
		testutils.WaitWithTimeout(t, ch, "expected aliveLoop to exit")
	})

	t.Run("with no poll failures past threshold, stays alive", func(t *testing.T) {
		threshold := 5
		cfg := TestNodeConfig{PollFailureThreshold: uint32(threshold), PollInterval: testutils.TestInterval}
		var calls atomic.Int32
		n := newTestNodeWithCallback(t, cfg, func(method string, params gjson.Result) (resp testutils.JSONRPCResponse) {
			switch method {
			case "eth_subscribe":
				resp.Result = `"0x00"`
				resp.Notify = makeHeadResult(0)
				return
			case "eth_unsubscribe":
				resp.Result = "true"
				return
			case "web3_clientVersion":
				defer calls.Add(1)
				// It starts working right before it hits threshold
				if int(calls.Load())+1 >= threshold {
					resp.Result = `"test client version"`
					return
				}
				resp.Result = "this will error"
				return
			default:
				t.Errorf("unexpected RPC method: %s", method)
			}
			return
		})
		dial(t, n)
		defer func() { assert.NoError(t, n.Close()) }()

		n.wg.Add(1)
		go n.aliveLoop()

		testutils.AssertEventually(t, func() bool {
			// Need to wait for one complete cycle before checking state so add
			// 1 to threshold
			return int(calls.Load()) > threshold+1
		})

		assert.Equal(t, NodeStateAlive, n.State())
	})

	t.Run("with threshold poll failures, transitions to unreachable", func(t *testing.T) {
		syncTimeoutsDisabledCfg := TestNodeConfig{PollFailureThreshold: 3, PollInterval: testutils.TestInterval}
		n := newTestNode(t, syncTimeoutsDisabledCfg)
		dial(t, n)
		defer func() { assert.NoError(t, n.Close()) }()

		n.wg.Add(1)
		go n.aliveLoop()

		testutils.AssertEventually(t, func() bool {
			return n.State() == NodeStateUnreachable
		})
	})

	t.Run("with threshold poll failures, but we are the last node alive, forcibly keeps it alive", func(t *testing.T) {
		threshold := 3
		cfg := TestNodeConfig{PollFailureThreshold: uint32(threshold), PollInterval: testutils.TestInterval}
		var calls atomic.Int32
		n := newTestNodeWithCallback(t, cfg, func(method string, params gjson.Result) (resp testutils.JSONRPCResponse) {
			switch method {
			case "eth_subscribe":
				resp.Result = `"0x00"`
				resp.Notify = HeadResult
				return
			case "eth_unsubscribe":
				resp.Result = "true"
				return
			case "web3_clientVersion":
				defer calls.Add(1)
				resp.Error.Message = "this will error"
				return
			default:
				t.Errorf("unexpected RPC method: %s", method)
			}
			return
		})
		n.nLiveNodes = func() (int, int64, *utils.Big) { return 1, 0, nil }
		dial(t, n)
		defer func() { assert.NoError(t, n.Close()) }()

		n.wg.Add(1)
		go n.aliveLoop()

		testutils.AssertEventually(t, func() bool {
			// Need to wait for one complete cycle before checking state so add
			// 1 to threshold
			return int(calls.Load()) > threshold+1
		})

		assert.Equal(t, NodeStateAlive, n.State())
	})

	t.Run("if initial subscribe fails, transitions to unreachable", func(t *testing.T) {
		pollDisabledCfg := TestNodeConfig{NoNewHeadsThreshold: testutils.TestInterval}
		n := newTestNodeWithCallback(t, pollDisabledCfg, func(string, gjson.Result) (resp testutils.JSONRPCResponse) { return })
		dial(t, n)
		defer func() { assert.NoError(t, n.Close()) }()

		_, err := n.EthSubscribe(testutils.Context(t), make(chan *evmtypes.Head))
		assert.Error(t, err)

		n.wg.Add(1)
		n.aliveLoop()

		assert.Equal(t, NodeStateUnreachable, n.State())
		// sc-39341: ensure failed EthSubscribe didn't register a (*rpc.ClientSubscription)(nil) which would lead to a panic on Unsubscribe
		assert.Len(t, n.subs, 0)
	})

	t.Run("if remote RPC connection is closed transitions to unreachable", func(t *testing.T) {
		// NoNewHeadsThreshold needs to be positive but must be very large so
		// we don't time out waiting for a new head before we have a chance to
		// handle the server disconnect
		cfg := TestNodeConfig{NoNewHeadsThreshold: testutils.WaitTimeout(t), PollInterval: 1 * time.Second}
		chSubbed := make(chan struct{}, 1)
		chPolled := make(chan struct{})
		s := testutils.NewWSServer(t, testutils.FixtureChainID,
			func(method string, params gjson.Result) (resp testutils.JSONRPCResponse) {
				switch method {
				case "eth_subscribe":
					chSubbed <- struct{}{}
					resp.Result = `"0x00"`
					resp.Notify = makeHeadResult(0)
					return
				case "web3_clientVersion":
					select {
					case chPolled <- struct{}{}:
					default:
					}
					resp.Result = `"test client version 2"`
					return
				default:
					t.Errorf("unexpected RPC method: %s", method)
				}
				return
			})

		iN := NewNode(cfg, logger.TestLogger(t), *s.WSURL(), nil, "test node", 42, testutils.FixtureChainID)
		n := iN.(*node)

		dial(t, n)
		defer func() { assert.NoError(t, n.Close()) }()

		n.wg.Add(1)
		go n.aliveLoop()

		testutils.WaitWithTimeout(t, chSubbed, "timed out waiting for initial subscription")
		testutils.WaitWithTimeout(t, chPolled, "timed out waiting for initial poll")

		assert.Equal(t, NodeStateAlive, n.State())

		// Simulate remote websocket disconnect
		// This causes sub.Err() to close
		s.Close()

		testutils.AssertEventually(t, func() bool {
			return n.State() == NodeStateUnreachable
		})
	})

	t.Run("when no new heads received for threshold, transitions to out of sync", func(t *testing.T) {
		cfg := TestNodeConfig{NoNewHeadsThreshold: 1 * time.Second}
		chSubbed := make(chan struct{}, 2)
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
				case "web3_clientVersion":
					resp.Result = `"test client version 2"`
					return
				default:
					t.Errorf("unexpected RPC method: %s", method)
				}
				return
			})

		iN := NewNode(cfg, logger.TestLogger(t), *s.WSURL(), nil, "test node", 42, testutils.FixtureChainID)
		n := iN.(*node)

		dial(t, n)
		defer func() { assert.NoError(t, n.Close()) }()

		n.wg.Add(1)
		go n.aliveLoop()

		testutils.WaitWithTimeout(t, chSubbed, "timed out waiting for initial subscription for InSync")

		testutils.AssertEventually(t, func() bool {
			return n.State() == NodeStateOutOfSync
		})

		// Otherwise, there may be data race on dial() vs Close() (accessing ws.rpc)
		testutils.WaitWithTimeout(t, chSubbed, "timed out waiting for initial subscription for OutOfSync")
	})

	t.Run("when no new heads received for threshold but we are the last live node, forcibly stays alive", func(t *testing.T) {
		lggr, observedLogs := logger.TestLoggerObserved(t, zap.ErrorLevel)
		pollDisabledCfg := TestNodeConfig{NoNewHeadsThreshold: testutils.TestInterval}
		s := testutils.NewWSServer(t, testutils.FixtureChainID,
			func(method string, params gjson.Result) (resp testutils.JSONRPCResponse) {
				switch method {
				case "eth_subscribe":
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

		iN := NewNode(pollDisabledCfg, lggr, *s.WSURL(), nil, "test node", 42, testutils.FixtureChainID)
		n := iN.(*node)
		n.nLiveNodes = func() (int, int64, *utils.Big) { return 1, 0, nil }
		dial(t, n)
		defer func() { assert.NoError(t, n.Close()) }()

		n.wg.Add(1)
		go n.aliveLoop()

		// to avoid timing-dependent tests, simply wait for the log message instead
		// wait for the log twice to be sure we have fully completed the code path and gone around the loop
		testutils.WaitForLogMessageCount(t, observedLogs, msgCannotDisable, 2)

		assert.Equal(t, NodeStateAlive, n.State())
	})

	t.Run("when behind more than SyncThreshold, transitions to out of sync", func(t *testing.T) {
		cfg := TestNodeConfig{SyncThreshold: 10, PollFailureThreshold: 2, PollInterval: 100 * time.Millisecond, SelectionMode: NodeSelectionMode_HighestHead}
		chSubbed := make(chan struct{}, 2)
		var highestHead atomic.Int64
		const stall = 10
		s := testutils.NewWSServer(t, testutils.FixtureChainID,
			func(method string, params gjson.Result) (resp testutils.JSONRPCResponse) {
				switch method {
				case "eth_subscribe":
					chSubbed <- struct{}{}
					resp.Result = `"0x00"`
					resp.Notify = makeHeadResult(int(highestHead.Load()))
					return
				case "eth_unsubscribe":
					resp.Result = "true"
					return
				case "web3_clientVersion":
					resp.Result = `"test client version 2"`
					// always tick each poll, but only signal back up to stall
					if n := highestHead.Add(1); n <= stall {
						resp.Notify = makeHeadResult(int(n))
					}
					return
				default:
					t.Errorf("unexpected RPC method: %s", method)
				}
				return
			})

		iN := NewNode(cfg, logger.TestLogger(t), *s.WSURL(), nil, "test node", 42, testutils.FixtureChainID)
		n := iN.(*node)
		n.nLiveNodes = func() (count int, blockNumber int64, totalDifficulty *utils.Big) {
			return 2, highestHead.Load(), nil
		}

		dial(t, n)
		defer func() { assert.NoError(t, n.Close()) }()

		n.wg.Add(1)
		go n.aliveLoop()

		testutils.WaitWithTimeout(t, chSubbed, "timed out waiting for initial subscription for InSync")

		// ensure alive up to stall
		testutils.AssertEventually(t, func() bool {
			state, num, _ := n.StateAndLatest()
			if num < stall {
				require.Equal(t, NodeStateAlive, state)
			}
			return num == stall
		})

		testutils.AssertEventually(t, func() bool {
			state, num, _ := n.StateAndLatest()
			return state == NodeStateOutOfSync && num == stall
		})
		assert.GreaterOrEqual(t, highestHead.Load(), int64(stall+cfg.SyncThreshold))

		// Otherwise, there may be data race on dial() vs Close() (accessing ws.rpc)
		testutils.WaitWithTimeout(t, chSubbed, "timed out waiting for initial subscription for OutOfSync")
	})

	t.Run("when behind but SyncThreshold=0, stay alive", func(t *testing.T) {
		cfg := TestNodeConfig{SyncThreshold: 0, PollFailureThreshold: 2, PollInterval: 100 * time.Millisecond, SelectionMode: NodeSelectionMode_HighestHead}
		chSubbed := make(chan struct{}, 1)
		var highestHead atomic.Int64
		const stall = 10
		s := testutils.NewWSServer(t, testutils.FixtureChainID,
			func(method string, params gjson.Result) (resp testutils.JSONRPCResponse) {
				switch method {
				case "eth_subscribe":
					chSubbed <- struct{}{}
					resp.Result = `"0x00"`
					resp.Notify = makeHeadResult(int(highestHead.Load()))
					return
				case "eth_unsubscribe":
					resp.Result = "true"
					return
				case "web3_clientVersion":
					resp.Result = `"test client version 2"`
					// always tick each poll, but only signal back up to stall
					if n := highestHead.Add(1); n <= stall {
						resp.Notify = makeHeadResult(int(n))
					}
					return
				default:
					t.Errorf("unexpected RPC method: %s", method)
				}
				return
			})

		iN := NewNode(cfg, logger.TestLogger(t), *s.WSURL(), nil, "test node", 42, testutils.FixtureChainID)
		n := iN.(*node)
		n.nLiveNodes = func() (count int, blockNumber int64, totalDifficulty *utils.Big) {
			return 2, highestHead.Load(), nil
		}

		dial(t, n)
		defer func() { assert.NoError(t, n.Close()) }()

		n.wg.Add(1)
		go n.aliveLoop()

		testutils.WaitWithTimeout(t, chSubbed, "timed out waiting for initial subscription for InSync")

		// ensure alive up to stall
		testutils.AssertEventually(t, func() bool {
			state, num, _ := n.StateAndLatest()
			require.Equal(t, NodeStateAlive, state)
			return num == stall
		})

		assert.Equal(t, NodeStateAlive, n.state)
		assert.GreaterOrEqual(t, highestHead.Load(), int64(stall+cfg.SyncThreshold))
	})

	t.Run("when behind more than SyncThreshold but we are the last live node, forcibly stays alive", func(t *testing.T) {
		lggr, observedLogs := logger.TestLoggerObserved(t, zap.ErrorLevel)
		cfg := TestNodeConfig{SyncThreshold: 5, PollFailureThreshold: 2, PollInterval: 100 * time.Millisecond, SelectionMode: NodeSelectionMode_HighestHead}
		chSubbed := make(chan struct{}, 1)
		var highestHead atomic.Int64
		const stall = 10
		s := testutils.NewWSServer(t, testutils.FixtureChainID,
			func(method string, params gjson.Result) (resp testutils.JSONRPCResponse) {
				switch method {
				case "eth_subscribe":
					chSubbed <- struct{}{}
					resp.Result = `"0x00"`
					n := highestHead.Load()
					if n > stall {
						n = stall
					}
					resp.Notify = makeHeadResult(int(n))
					return
				case "eth_unsubscribe":
					resp.Result = "true"
					return
				case "web3_clientVersion":
					resp.Result = `"test client version 2"`
					// always tick each poll, but only signal back up to stall
					if n := highestHead.Add(1); n <= stall {
						resp.Notify = makeHeadResult(int(n))
					}
					return
				default:
					t.Errorf("unexpected RPC method: %s", method)
				}
				return
			})

		iN := NewNode(cfg, lggr, *s.WSURL(), nil, "test node", 42, testutils.FixtureChainID)
		n := iN.(*node)
		n.nLiveNodes = func() (count int, blockNumber int64, totalDifficulty *utils.Big) {
			return 1, highestHead.Load(), nil
		}

		dial(t, n)
		defer func() { assert.NoError(t, n.Close()) }()

		n.wg.Add(1)
		go n.aliveLoop()

		testutils.WaitWithTimeout(t, chSubbed, "timed out waiting for initial subscription for InSync")

		// ensure alive up to stall
		testutils.AssertEventually(t, func() bool {
			state, num, _ := n.StateAndLatest()
			require.Equal(t, NodeStateAlive, state)
			return num == stall
		})

		assert.Equal(t, NodeStateAlive, n.state)
		testutils.AssertEventually(t, func() bool {
			return highestHead.Load() >= int64(stall+cfg.SyncThreshold)
		})

		testutils.WaitForLogMessageCount(t, observedLogs, msgCannotDisable, 1)

		state, num, _ := n.StateAndLatest()
		assert.Equal(t, NodeStateAlive, state)
		assert.Equal(t, int64(stall), num)

	})
}

func TestUnit_NodeLifecycle_outOfSyncLoop(t *testing.T) {
	t.Parallel()

	t.Run("exits on close", func(t *testing.T) {
		cfg := TestNodeConfig{}
		n := newTestNode(t, cfg)
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
		cfg := TestNodeConfig{}
		n := newTestNodeWithCallback(t, cfg, func(string, gjson.Result) (resp testutils.JSONRPCResponse) { return })
		dial(t, n)
		n.setState(NodeStateOutOfSync)
		defer func() { assert.NoError(t, n.Close()) }()

		n.wg.Add(1)

		n.outOfSyncLoop(func(num int64, td *utils.Big) bool { return num == 0 })
		assert.Equal(t, NodeStateUnreachable, n.State())
	})

	t.Run("transitions to unreachable if remote RPC subscription channel closed", func(t *testing.T) {
		cfg := TestNodeConfig{}
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

		iN := NewNode(cfg, logger.TestLogger(t), *s.WSURL(), nil, "test node", 42, testutils.FixtureChainID)
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
		cfg := TestNodeConfig{}
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

		iN := NewNode(cfg, lggr, *s.WSURL(), nil, "test node", 0, testutils.FixtureChainID)
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
		cfg := TestNodeConfig{SyncThreshold: 5, SelectionMode: NodeSelectionMode_HighestHead}
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

		iN := NewNode(cfg, lggr, *s.WSURL(), nil, "test node", 0, testutils.FixtureChainID)
		n := iN.(*node)
		n.nLiveNodes = func() (count int, blockNumber int64, totalDifficulty *utils.Big) {
			return 2, stall + int64(cfg.SyncThreshold), nil
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
		cfg := TestNodeConfig{NoNewHeadsThreshold: testutils.TestInterval}
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

		iN := NewNode(cfg, logger.TestLogger(t), *s.WSURL(), nil, "test node", 42, testutils.FixtureChainID)
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
		cfg := TestNodeConfig{}
		n := newTestNode(t, cfg)
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
		cfg := TestNodeConfig{}
		n := newTestNode(t, cfg)
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
		cfg := TestNodeConfig{}
		s := testutils.NewWSServer(t, testutils.FixtureChainID, standardHandler)
		lggr, observedLogs := logger.TestLoggerObserved(t, zap.ErrorLevel)
		iN := NewNode(cfg, lggr, *s.WSURL(), nil, "test node", 0, big.NewInt(42))
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
		cfg := TestNodeConfig{}
		lggr, observedLogs := logger.TestLoggerObserved(t, zap.DebugLevel)
		iN := NewNode(cfg, lggr, *testutils.MustParseURL(t, "ws://test.invalid"), nil, "test node", 0, big.NewInt(42))
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
		cfg := TestNodeConfig{}
		n := newTestNode(t, cfg)
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
		cfg := TestNodeConfig{}
		n := newTestNode(t, cfg)
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
		cfg := TestNodeConfig{}
		s := testutils.NewWSServer(t, testutils.FixtureChainID, standardHandler)
		lggr, observedLogs := logger.TestLoggerObserved(t, zap.ErrorLevel)
		iN := NewNode(cfg, lggr, *s.WSURL(), nil, "test node", 0, big.NewInt(42))
		n := iN.(*node)
		defer func() { assert.NoError(t, n.Close()) }()
		dial(t, n)
		n.setState(NodeStateUnreachable)
		n.wg.Add(1)

		go n.unreachableLoop()

		testutils.WaitForLogMessageCount(t, observedLogs, "Failed to redial RPC node; remote endpoint returned the wrong chain ID", 3)

		assert.Equal(t, NodeStateInvalidChainID, n.State())
	})
}
