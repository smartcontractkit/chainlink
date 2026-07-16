package gateway_test

import (
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/jonboulle/clockwork"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	jsonrpc "github.com/smartcontractkit/chainlink-common/pkg/jsonrpc2"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/limits"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway"
	gc "github.com/smartcontractkit/chainlink/v2/core/services/gateway/common"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/config"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/monitoring"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/network"
)

const defaultConfig = `
[nodeServerConfig]
Path = "/node"

[[dons]]
DonId = "my_don_1"
HandlerName = "dummy"

[[dons.members]]
Name = "example_node"
Address = "0x68902D681C28119F9B2531473A417088BF008E59"

[[dons]]
DonId = "my_don_2"
HandlerName = "dummy"

[[dons.members]]
Name = "example_node"
Address = "0x68902d681c28119f9b2531473a417088bf008e59"
`

func TestConnectionManager_NewConnectionManager_ValidConfig(t *testing.T) {
	t.Parallel()

	tomlConfig := parseTOMLConfig(t, defaultConfig)

	_ = newConnectionManager(t, tomlConfig, clockwork.NewFakeClock())
}

func TestConnectionManager_NewConnectionManager_InvalidConfig(t *testing.T) {
	t.Parallel()

	invalidCases := map[string]string{
		"duplicate DON ID": `
[[dons]]
DonId = "my_don"
[[dons]]
DonId = "my_don"
`,
		"duplicate node address": `
[[dons]]
DonId = "my_don"
[[dons.members]]
Name = "node_1"
Address = "0x68902d681c28119f9b2531473a417088bf008e59"
[[dons.members]]
Name = "node_2"
Address = "0x68902d681c28119f9b2531473a417088bf008e59"
`,
		"duplicate node address with different casing": `
[[dons]]
DonId = "my_don"
[[dons.members]]
Name = "node_1"
Address = "0x68902d681c28119f9b2531473a417088bf008e59"
[[dons.members]]
Name = "node_2"
Address = "0x68902D681c28119f9b2531473a417088bf008E59"
`,
	}

	for name, config := range invalidCases {
		t.Run(name, func(t *testing.T) {
			fullConfig := `
[nodeServerConfig]
Path = "/node"` + config
			lggr := logger.Test(t)
			gMetrics, err := monitoring.NewGatewayMetrics()
			require.NoError(t, err)
			_, err = gateway.NewConnectionManager(parseTOMLConfig(t, fullConfig), clockwork.NewFakeClock(), gMetrics, lggr, limits.Factory{Logger: lggr})
			require.Error(t, err)
		})
	}
}

func newTestConfig(t *testing.T, nNodes int) (*config.GatewayConfig, []gc.TestNode) {
	nodes := gc.NewTestNodes(t, nNodes)

	var config strings.Builder
	config.WriteString(`
[nodeServerConfig]
Path = "/node"
[connectionManagerConfig]
AuthGatewayId = "my_gateway_no_3"
AuthTimestampToleranceSec = 5
AuthChallengeLen = 100
[[dons]]
DonId = "my_don_1"
HandlerName = "dummy"
`)

	for i := range nNodes {
		config.WriteString(`[[dons.members]]` + "\n")
		config.WriteString(fmt.Sprintf(`Name = "node_%d"`, i) + "\n")
		config.WriteString(fmt.Sprintf(`Address = "%s"`, nodes[i].Address) + "\n")
	}

	return parseTOMLConfig(t, config.String()), nodes
}

func signAndPackAuthHeader(t *testing.T, authHeaderElems *network.AuthHeaderElems, signerKey *ecdsa.PrivateKey) []byte {
	packedElems := network.PackAuthHeader(authHeaderElems)
	signature, err := gc.SignData(signerKey, packedElems)
	require.NoError(t, err)
	return append(packedElems, signature...)
}

func TestConnectionManager_StartHandshake(t *testing.T) {
	t.Parallel()

	config, nodes := newTestConfig(t, 4)
	unrelatedNode := gc.NewTestNodes(t, 1)[0]
	clock := clockwork.NewFakeClock()
	mgr := newConnectionManager(t, config, clock)

	authHeaderElems := network.AuthHeaderElems{
		Timestamp: uint32(clock.Now().Unix()),
		DonId:     "my_don_1",
		GatewayId: "my_gateway_no_3",
	}

	// valid
	_, _, err := mgr.StartHandshake(signAndPackAuthHeader(t, &authHeaderElems, nodes[0].PrivateKey))
	require.NoError(t, err)

	// header too short
	_, _, err = mgr.StartHandshake([]byte("ab"))
	require.ErrorIs(t, err, network.ErrAuthHeaderParse)

	// invalid DON ID
	badAuthHeaderElems := authHeaderElems
	badAuthHeaderElems.DonId = "my_don_2"
	_, _, err = mgr.StartHandshake(signAndPackAuthHeader(t, &badAuthHeaderElems, nodes[0].PrivateKey))
	require.ErrorIs(t, err, network.ErrAuthInvalidDonId)

	// invalid Gateway URL
	badAuthHeaderElems = authHeaderElems
	badAuthHeaderElems.GatewayId = "www.example.com"
	_, _, err = mgr.StartHandshake(signAndPackAuthHeader(t, &badAuthHeaderElems, nodes[0].PrivateKey))
	require.ErrorIs(t, err, network.ErrAuthInvalidGateway)

	// invalid Signer Address
	badAuthHeaderElems = authHeaderElems
	_, _, err = mgr.StartHandshake(signAndPackAuthHeader(t, &badAuthHeaderElems, unrelatedNode.PrivateKey))
	require.ErrorIs(t, err, network.ErrAuthInvalidNode)

	// invalid signature
	badAuthHeaderElems = authHeaderElems
	rawHeader := signAndPackAuthHeader(t, &badAuthHeaderElems, nodes[0].PrivateKey)
	copy(rawHeader[len(rawHeader)-65:], make([]byte, 65))
	_, _, err = mgr.StartHandshake(rawHeader)
	require.ErrorIs(t, err, network.ErrAuthHeaderParse)

	// invalid timestamp
	badAuthHeaderElems = authHeaderElems
	badAuthHeaderElems.Timestamp -= 10
	_, _, err = mgr.StartHandshake(signAndPackAuthHeader(t, &badAuthHeaderElems, nodes[0].PrivateKey))
	require.ErrorIs(t, err, network.ErrAuthInvalidTimestamp)
}

func TestConnectionManager_FinalizeHandshake(t *testing.T) {
	t.Parallel()

	config, nodes := newTestConfig(t, 4)
	clock := clockwork.NewFakeClock()
	mgr := newConnectionManager(t, config, clock)

	authHeaderElems := network.AuthHeaderElems{
		Timestamp: uint32(clock.Now().Unix()),
		DonId:     "my_don_1",
		GatewayId: "my_gateway_no_3",
	}

	// correct
	attemptId, challenge, err := mgr.StartHandshake(signAndPackAuthHeader(t, &authHeaderElems, nodes[0].PrivateKey))
	require.NoError(t, err)
	response, err := gc.SignData(nodes[0].PrivateKey, challenge)
	require.NoError(t, err)
	require.NoError(t, mgr.FinalizeHandshake(attemptId, response, nil))

	// invalid attempt
	err = mgr.FinalizeHandshake("fake_attempt", response, nil)
	require.ErrorIs(t, err, network.ErrChallengeAttemptNotFound)

	// invalid signature
	attemptId, challenge, err = mgr.StartHandshake(signAndPackAuthHeader(t, &authHeaderElems, nodes[0].PrivateKey))
	require.NoError(t, err)
	response, err = gc.SignData(nodes[1].PrivateKey, challenge)
	require.NoError(t, err)
	err = mgr.FinalizeHandshake(attemptId, response, nil)
	require.ErrorIs(t, err, network.ErrChallengeInvalidSignature)
}

func TestConnectionManager_SendToNode_Failures(t *testing.T) {
	t.Parallel()

	config, nodes := newTestConfig(t, 2)
	clock := clockwork.NewFakeClock()
	mgr := newConnectionManager(t, config, clock)

	donMgr := mgr.DONConnectionManager("my_don_1")
	err := donMgr.SendToNode(testutils.Context(t), nodes[0].Address, nil)
	require.Error(t, err)

	message := &jsonrpc.Request[json.RawMessage]{}
	err = donMgr.SendToNode(testutils.Context(t), "some_other_node", message)
	require.Error(t, err)
}

func TestConnectionManager_CleanStartClose(t *testing.T) {
	t.Parallel()

	config, _ := newTestConfig(t, 2)
	config.ConnectionManagerConfig.HeartbeatIntervalSec = 1
	clock := clockwork.NewFakeClock()
	mgr := newConnectionManager(t, config, clock)

	err := mgr.Start(testutils.Context(t))
	require.NoError(t, err)

	err = mgr.Close()
	require.NoError(t, err)
}

func TestConnectionManager_ShardedDONs_CreatesPerShardManagers(t *testing.T) {
	t.Parallel()

	tomlConfig := `
[nodeServerConfig]
Path = "/node"

[[shardedDONs]]
DonName = "myDON"
F = 1

[[shardedDONs.Shards]]
[[shardedDONs.Shards.Nodes]]
Name = "s0_n0"
Address = "0x0001020304050607080900010203040506070809"
[[shardedDONs.Shards.Nodes]]
Name = "s0_n1"
Address = "0x0002020304050607080900010203040506070809"
[[shardedDONs.Shards.Nodes]]
Name = "s0_n2"
Address = "0x0003020304050607080900010203040506070809"
[[shardedDONs.Shards.Nodes]]
Name = "s0_n3"
Address = "0x0004020304050607080900010203040506070809"

[[shardedDONs.Shards]]
[[shardedDONs.Shards.Nodes]]
Name = "s1_n0"
Address = "0x0005020304050607080900010203040506070809"
[[shardedDONs.Shards.Nodes]]
Name = "s1_n1"
Address = "0x0006020304050607080900010203040506070809"
[[shardedDONs.Shards.Nodes]]
Name = "s1_n2"
Address = "0x0007020304050607080900010203040506070809"
[[shardedDONs.Shards.Nodes]]
Name = "s1_n3"
Address = "0x0008020304050607080900010203040506070809"
`

	cfg := parseTOMLConfig(t, tomlConfig)
	mgr := newConnectionManager(t, cfg, clockwork.NewFakeClock())

	require.NotNil(t, mgr.DONConnectionManager(config.ShardDONID("myDON", 0)), "shard 0 connection manager should exist")
	require.NotNil(t, mgr.DONConnectionManager(config.ShardDONID("myDON", 1)), "shard 1 connection manager should exist")
	require.Nil(t, mgr.DONConnectionManager("myDON_2"), "shard 2 should not exist")
}

func TestConnectionManager_ShardedDONs_MultipleDONs(t *testing.T) {
	t.Parallel()

	tomlConfig := `
[nodeServerConfig]
Path = "/node"

[[shardedDONs]]
DonName = "donA"
F = 0

[[shardedDONs.Shards]]
[[shardedDONs.Shards.Nodes]]
Name = "a_n0"
Address = "0x0001020304050607080900010203040506070809"

[[shardedDONs]]
DonName = "donB"
F = 0

[[shardedDONs.Shards]]
[[shardedDONs.Shards.Nodes]]
Name = "b_n0"
Address = "0x0002020304050607080900010203040506070809"
`

	cfg := parseTOMLConfig(t, tomlConfig)
	mgr := newConnectionManager(t, cfg, clockwork.NewFakeClock())

	require.NotNil(t, mgr.DONConnectionManager(config.ShardDONID("donA", 0)))
	require.NotNil(t, mgr.DONConnectionManager(config.ShardDONID("donB", 0)))
}

func TestConnectionManager_ShardedDONs_DuplicateNodeAddress(t *testing.T) {
	t.Parallel()

	tomlConfig := `
[nodeServerConfig]
Path = "/node"

[[shardedDONs]]
DonName = "myDON"
F = 0

[[shardedDONs.Shards]]
[[shardedDONs.Shards.Nodes]]
Name = "n0"
Address = "0x0001020304050607080900010203040506070809"
[[shardedDONs.Shards.Nodes]]
Name = "n1"
Address = "0x0001020304050607080900010203040506070809"
`

	cfg := parseTOMLConfig(t, tomlConfig)
	lggr := logger.Test(t)
	gMetrics, err := monitoring.NewGatewayMetrics()
	require.NoError(t, err)
	_, err = gateway.NewConnectionManager(cfg, clockwork.NewFakeClock(), gMetrics, lggr, limits.Factory{Logger: lggr})
	require.Error(t, err)
	require.Contains(t, err.Error(), "duplicate node address")
}

func TestConnectionManager_ShardedDONs_SendToNode(t *testing.T) {
	t.Parallel()

	tomlConfig := `
[nodeServerConfig]
Path = "/node"

[[shardedDONs]]
DonName = "myDON"
F = 0

[[shardedDONs.Shards]]
[[shardedDONs.Shards.Nodes]]
Name = "n0"
Address = "0x0001020304050607080900010203040506070809"
`

	cfg := parseTOMLConfig(t, tomlConfig)
	mgr := newConnectionManager(t, cfg, clockwork.NewFakeClock())

	donMgr := mgr.DONConnectionManager(config.ShardDONID("myDON", 0))
	require.NotNil(t, donMgr)

	err := donMgr.SendToNode(testutils.Context(t), "0x0001020304050607080900010203040506070809", nil)
	require.Error(t, err, "nil request should fail")

	message := &jsonrpc.Request[json.RawMessage]{}
	err = donMgr.SendToNode(testutils.Context(t), "0xdeadbeef", message)
	require.Error(t, err, "unknown node should fail")
}

func TestConnectionManager_ShardedDONs_StartClose(t *testing.T) {
	t.Parallel()

	tomlConfig := `
[nodeServerConfig]
Path = "/node"
[connectionManagerConfig]
HeartbeatIntervalSec = 1

[[shardedDONs]]
DonName = "myDON"
F = 0

[[shardedDONs.Shards]]
[[shardedDONs.Shards.Nodes]]
Name = "n0"
Address = "0x0001020304050607080900010203040506070809"
`

	cfg := parseTOMLConfig(t, tomlConfig)
	mgr := newConnectionManager(t, cfg, clockwork.NewFakeClock())

	err := mgr.Start(testutils.Context(t))
	require.NoError(t, err)

	err = mgr.Close()
	require.NoError(t, err)
}

// newWebSocketPair creates an httptest server with a WebSocket upgrader and dials it.
// Returns the server-side conn (to pass to FinalizeHandshake) and the client-side conn
// (simulating the node). The httptest server is closed on test cleanup.
func newWebSocketPair(t *testing.T) (serverConn, clientConn *websocket.Conn) {
	t.Helper()
	upgrader := websocket.Upgrader{}
	connCh := make(chan *websocket.Conn, 1)
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		connCh <- c
	}))
	t.Cleanup(s.Close)

	clientConn, _, err := websocket.DefaultDialer.Dial("ws"+strings.TrimPrefix(s.URL, "http"), nil)
	require.NoError(t, err)
	t.Cleanup(func() { clientConn.Close() })

	select {
	case serverConn = <-connCh:
	case <-time.After(5 * time.Second):
		t.Fatal("timed out waiting for WebSocket upgrade")
	}
	t.Cleanup(func() { serverConn.Close() })
	return serverConn, clientConn
}

// doHandshake performs StartHandshake + FinalizeHandshake with the given conn.
func doHandshake(t *testing.T, mgr gateway.ConnectionManager, clock clockwork.Clock, node gc.TestNode, conn *websocket.Conn) {
	t.Helper()
	authHeaderElems := network.AuthHeaderElems{
		Timestamp: uint32(clock.Now().Unix()), //nolint:gosec // test clock is always small positive
		DonId:     "my_don_1",
		GatewayId: "my_gateway_no_3",
	}
	attemptID, challenge, err := mgr.StartHandshake(signAndPackAuthHeader(t, &authHeaderElems, node.PrivateKey))
	require.NoError(t, err)
	response, err := gc.SignData(node.PrivateKey, challenge)
	require.NoError(t, err)
	require.NoError(t, mgr.FinalizeHandshake(attemptID, response, conn))
}

func TestConnectionManager_ReadDeadline_ClosesIdleConnection(t *testing.T) {
	t.Parallel()

	cfg, nodes := newTestConfig(t, 1)
	cfg.ConnectionManagerConfig.HeartbeatIntervalSec = 1
	cfg.ConnectionManagerConfig.PongTimeoutSec = 3
	clock := clockwork.NewRealClock()
	mgr := newConnectionManager(t, cfg, clock)

	err := mgr.Start(testutils.Context(t))
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, mgr.Close()) })

	serverConn, _ := newWebSocketPair(t)
	// Client does NOT read from its connection, so pings from the gateway
	// are never processed and pongs are never sent back.

	doHandshake(t, mgr, clock, nodes[0], serverConn)

	// After PongTimeoutSec (3s) the read deadline fires, readPump closes the conn.
	// SendToNode should eventually fail because the underlying conn is closed.
	donMgr := mgr.DONConnectionManager("my_don_1")
	require.NotNil(t, donMgr)

	// Verify connection works initially.
	msg := &jsonrpc.Request[json.RawMessage]{ID: "pre-check"}
	require.NoError(t, donMgr.SendToNode(testutils.Context(t), nodes[0].Address, msg))

	assert.Eventually(t, func() bool {
		msg := &jsonrpc.Request[json.RawMessage]{ID: "test"}
		return donMgr.SendToNode(testutils.Context(t), nodes[0].Address, msg) != nil
	}, 6*time.Second, 200*time.Millisecond, "connection should be closed by read deadline")
}

func TestConnectionManager_ReadDeadline_ConnectionAliveWithPongs(t *testing.T) {
	t.Parallel()

	cfg, nodes := newTestConfig(t, 1)
	cfg.ConnectionManagerConfig.HeartbeatIntervalSec = 1
	cfg.ConnectionManagerConfig.PongTimeoutSec = 3
	clock := clockwork.NewRealClock()
	mgr := newConnectionManager(t, cfg, clock)

	err := mgr.Start(testutils.Context(t))
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, mgr.Close()) })

	serverConn, clientConn := newWebSocketPair(t)

	// Client reads in a loop - gorilla/websocket's default ping handler
	// automatically responds with pongs, which resets the server's read deadline.
	go func() {
		for {
			_, _, err := clientConn.ReadMessage()
			if err != nil {
				return // connection closed
			}
		}
	}()

	doHandshake(t, mgr, clock, nodes[0], serverConn)

	// Assert the connection never dies during a period longer than PongTimeoutSec.
	// With HeartbeatIntervalSec=1, pongs arrive every ~1s resetting the 3s deadline.
	donMgr := mgr.DONConnectionManager("my_don_1")
	require.NotNil(t, donMgr)

	assert.Never(t, func() bool {
		msg := &jsonrpc.Request[json.RawMessage]{ID: "alive-check"}
		return donMgr.SendToNode(testutils.Context(t), nodes[0].Address, msg) != nil
	}, 5*time.Second, 500*time.Millisecond, "connection should stay alive when pongs are received")
}

func TestConnectionManager_ReadDeadline_DisabledWhenZero(t *testing.T) {
	t.Parallel()

	cfg, nodes := newTestConfig(t, 1)
	cfg.ConnectionManagerConfig.HeartbeatIntervalSec = 1
	// PongTimeoutSec = 0 (default) - deadline enforcement disabled
	clock := clockwork.NewRealClock()
	mgr := newConnectionManager(t, cfg, clock)

	err := mgr.Start(testutils.Context(t))
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, mgr.Close()) })

	serverConn, _ := newWebSocketPair(t)
	// Client does NOT read - no pongs sent back.

	doHandshake(t, mgr, clock, nodes[0], serverConn)

	// Even without pongs, the connection should remain alive because
	// PongTimeoutSec=0 means no read deadline is set.
	donMgr := mgr.DONConnectionManager("my_don_1")
	require.NotNil(t, donMgr)

	assert.Never(t, func() bool {
		msg := &jsonrpc.Request[json.RawMessage]{ID: "test"}
		return donMgr.SendToNode(testutils.Context(t), nodes[0].Address, msg) != nil
	}, 4*time.Second, 500*time.Millisecond, "connection should stay alive when deadline enforcement is disabled")
}

func newConnectionManager(t *testing.T, gwConfig *config.GatewayConfig, clock clockwork.Clock) gateway.ConnectionManager {
	lggr := logger.Test(t)
	gMetrics, err := monitoring.NewGatewayMetrics()
	require.NoError(t, err)
	mgr, err := gateway.NewConnectionManager(gwConfig, clock, gMetrics, lggr, limits.Factory{Logger: lggr})
	require.NoError(t, err)
	return mgr
}
