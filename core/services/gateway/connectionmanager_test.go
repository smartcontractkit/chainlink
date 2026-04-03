package gateway_test

import (
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/jonboulle/clockwork"
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

	config := `
[nodeServerConfig]
Path = "/node"
[connectionManagerConfig]
AuthGatewayId = "my_gateway_no_3"
AuthTimestampToleranceSec = 5
AuthChallengeLen = 100
[[dons]]
DonId = "my_don_1"
HandlerName = "dummy"
`

	for i := range nNodes {
		config += `[[dons.members]]` + "\n"
		config += fmt.Sprintf(`Name = "node_%d"`, i) + "\n"
		config += fmt.Sprintf(`Address = "%s"`, nodes[i].Address) + "\n"
	}

	return parseTOMLConfig(t, config), nodes
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

func newConnectionManager(t *testing.T, gwConfig *config.GatewayConfig, clock clockwork.Clock) gateway.ConnectionManager {
	lggr := logger.Test(t)
	gMetrics, err := monitoring.NewGatewayMetrics()
	require.NoError(t, err)
	mgr, err := gateway.NewConnectionManager(gwConfig, clock, gMetrics, lggr, limits.Factory{Logger: lggr})
	require.NoError(t, err)
	return mgr
}
