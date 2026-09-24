package gateway

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/config"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/monitoring"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/network"
)

// mockPingConn is a minimal WSConnectionWrapper for keepalive testing.
// It counts ping writes and optionally blocks on Write to simulate a
// half-open TCP connection.
type mockPingConn struct {
	pingCount atomic.Int64
	unblock   chan struct{} // if non-nil, Write blocks until closed
}

func (m *mockPingConn) Start(context.Context) error             { return nil }
func (m *mockPingConn) HealthReport() map[string]error          { return nil }
func (m *mockPingConn) Name() string                            { return "mockPingConn" }
func (m *mockPingConn) Ready() error                            { return nil }
func (m *mockPingConn) Reset(network.WSConnection) <-chan error { return nil }
func (m *mockPingConn) ReadChannel() <-chan network.ReadItem    { return nil }
func (m *mockPingConn) IsConnected() bool                       { return true }
func (m *mockPingConn) Close() error                            { return nil }

func (m *mockPingConn) Write(ctx context.Context, msgType int, _ []byte) error {
	if msgType == websocket.PingMessage {
		m.pingCount.Add(1)
	}
	if m.unblock != nil {
		select {
		case <-m.unblock:
			return nil
		case <-ctx.Done():
			return ctx.Err()
		}
	}
	return nil
}

// TestKeepAliveLoop_StuckNodeDoesBlocksAll
// this reproduces a bug in which one stall node blocks pings to all nodes.
// This leads to a wrong representation of the communication between nodes and gateway, making it harder to diagnose an incident.
//
// Before the keepaliveLoop used to ping nodes sequentially with no per-node timeout. If one
// node's Write blocked (half-open TCP), the loop stalls and no other node
// receives pings.
func TestKeepAliveLoop_StuckNodeBlocksAll(t *testing.T) {
	if testing.Short() {
		t.Skip("too slow for testing.Short")
	}

	t.Parallel()

	lggr := logger.Test(t)
	gMetrics, err := monitoring.NewGatewayMetrics()
	require.NoError(t, err)

	unblock := make(chan struct{})
	t.Cleanup(func() { close(unblock) })

	stuckConn := &mockPingConn{unblock: unblock}
	healthyConns := make([]*mockPingConn, 3)
	mockNodes := make(map[string]*nodeState, 4)
	mockNodes["0xstuck"] = &nodeState{name: "stuck", conn: stuckConn}
	for i := range healthyConns {
		healthyConns[i] = &mockPingConn{}
		mockNodes["0xhealthy"+string(rune('0'+i))] = &nodeState{
			name: "healthy_" + string(rune('0'+i)),
			conn: healthyConns[i],
		}
	}

	donMgr := &donConnectionManager{
		donConfig:  &config.DONConfig{DonID: "test_don"},
		nodes:      mockNodes,
		handlers:   nil,
		closeWait:  sync.WaitGroup{},
		shutdownCh: make(services.StopChan),
		gMetrics:   gMetrics,
		lggr:       lggr,
	}

	// Start the keepalive loop directly with a 1-second interval.
	const heartbeatSec = 1
	donMgr.closeWait.Add(len(donMgr.nodes))
	for nodeAddress, nodeState := range donMgr.nodes {
		go donMgr.nodeKeepalive(nodeAddress, nodeState, heartbeatSec)
	}

	// Let at least 2 ticks fire. With the bug, the loop is stuck on the
	// first node and never reaches the healthy nodes. With the fix,
	// healthy nodes get pinged each tick.
	time.Sleep(3 * time.Second)

	// --- Check ping counts BEFORE stopping the loop ---
	//
	// BEFORE the fix (sequential, no per-node timeout):
	//   The loop calls Write on "0xstuck" first. It blocks on unblock.
	//   The loop never reaches the healthy nodes.
	//   → all healthyConns have pingCount == 0 → test FAILS.
	//
	// AFTER the fix (per-node timeout or concurrent sends):
	//   Each node is pinged independently. The stuck node blocks/times out
	//   but healthy nodes still receive pings.
	//   → all healthyConns have pingCount >= 1 → test PASSES.
	for i, hc := range healthyConns {
		pings := hc.pingCount.Load()
		t.Logf("healthy_%d: %d pings", i, pings)
		require.Positive(t, pings, "healthy node %d received 0 pings — keepaliveLoop is stuck on the blocked node", i)
	}

	stuckPings := stuckConn.pingCount.Load()
	t.Logf("stuck node: %d ping attempts", stuckPings)

	// Stop the loop (after assertions so the stuck Write doesn't unblock early).
	close(donMgr.shutdownCh)
	donMgr.closeWait.Wait()
}

func TestPingProbeTracker_Matching(t *testing.T) {
	t.Parallel()

	tracker := &pingProbeTracker{}

	// No probe registered: no observation.
	_, ok := tracker.consume("unknown")
	require.False(t, ok)

	// Wrong token (e.g. a pong for a superseded probe): no observation, and the
	// pending probe stays registered.
	start := time.Now()
	tracker.register("token-a", start)
	_, ok = tracker.consume("token-b")
	require.False(t, ok)

	// Matching token: exactly one observation.
	got, ok := tracker.consume("token-a")
	require.True(t, ok)
	require.Equal(t, start, got)

	// Duplicate pong: no second observation.
	_, ok = tracker.consume("token-a")
	require.False(t, ok)

	// A superseded probe produces no observation: only the latest is retained.
	tracker.register("token-old", time.Now())
	newStart := time.Now()
	tracker.register("token-new", newStart)
	_, ok = tracker.consume("token-old")
	require.False(t, ok, "superseded probe must not record")
	got, ok = tracker.consume("token-new")
	require.True(t, ok)
	require.Equal(t, newStart, got)
}

func TestPingProbeTracker_StaleConnectionCannotMatchReplacement(t *testing.T) {
	t.Parallel()

	// nodeState carries the current connection's tracker; FinalizeHandshake
	// replaces it on every handshake.
	ns := &nodeState{name: "node"}
	stale := &pingProbeTracker{}
	ns.pingTracker.Store(stale)

	// Connection replaced: a fresh tracker is installed and receives new probes.
	fresh := &pingProbeTracker{}
	ns.pingTracker.Store(fresh)
	start := time.Now()
	fresh.register("token-replacement", start)

	// A pong arriving on the old connection is matched against the stale tracker
	// via its pong-handler closure: it can never record against the replacement.
	_, ok := stale.consume("token-replacement")
	require.False(t, ok, "stale connection's tracker must not match the replacement's probe")

	// The replacement connection's tracker still matches its own probe, once.
	got, ok := fresh.consume("token-replacement")
	require.True(t, ok)
	require.Equal(t, start, got)
}

func TestPingProbeTracker_Concurrent(t *testing.T) {
	t.Parallel()

	tracker := &pingProbeTracker{}
	var wg sync.WaitGroup
	for i := range 8 {
		wg.Add(2)
		go func() {
			defer wg.Done()
			for j := range 100 {
				tracker.register(fmt.Sprintf("token-%d-%d", i, j), time.Now())
			}
		}()
		go func() {
			defer wg.Done()
			for j := range 100 {
				tracker.consume(fmt.Sprintf("token-%d", j))
			}
		}()
	}
	wg.Wait()
}
