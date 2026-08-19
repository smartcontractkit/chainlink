package gateway

import (
	"context"
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

func (m *mockPingConn) Start(context.Context) error          { return nil }
func (m *mockPingConn) HealthReport() map[string]error       { return nil }
func (m *mockPingConn) Name() string                         { return "mockPingConn" }
func (m *mockPingConn) Ready() error                         { return nil }
func (m *mockPingConn) Reset(*websocket.Conn) <-chan error   { return nil }
func (m *mockPingConn) ReadChannel() <-chan network.ReadItem { return nil }
func (m *mockPingConn) Close() error                         { return nil }

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
		donConfig:  &config.DONConfig{DonId: "test_don"},
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
