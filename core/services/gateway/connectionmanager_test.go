package gateway_test

import (
	"crypto/ecdsa"
	"fmt"
	"testing"

	"github.com/jonboulle/clockwork"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/api"
	gc "github.com/smartcontractkit/chainlink/v2/core/services/gateway/common"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/config"
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

	_, err := gateway.NewConnectionManager(tomlConfig, clockwork.NewFakeClock(), logger.TestLogger(t))
	require.NoError(t, err)
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
		config := config
		t.Run(name, func(t *testing.T) {
			fullConfig := `
[nodeServerConfig]
Path = "/node"` + config
			_, err := gateway.NewConnectionManager(parseTOMLConfig(t, fullConfig), clockwork.NewFakeClock(), logger.TestLogger(t))
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

	for i := 0; i < nNodes; i++ {
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
	mgr, err := gateway.NewConnectionManager(config, clock, logger.TestLogger(t))
	require.NoError(t, err)

	authHeaderElems := network.AuthHeaderElems{
		Timestamp: uint32(clock.Now().Unix()),
		DonId:     "my_don_1",
		GatewayId: "my_gateway_no_3",
	}

	// valid
	_, _, err = mgr.StartHandshake(signAndPackAuthHeader(t, &authHeaderElems, nodes[0].PrivateKey))
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
	mgr, err := gateway.NewConnectionManager(config, clock, logger.TestLogger(t))
	require.NoError(t, err)

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
	mgr, err := gateway.NewConnectionManager(config, clock, logger.TestLogger(t))
	require.NoError(t, err)

	donMgr := mgr.DONConnectionManager("my_don_1")
	err = donMgr.SendToNode(testutils.Context(t), nodes[0].Address, nil)
	require.Error(t, err)

	message := &api.Message{}
	err = donMgr.SendToNode(testutils.Context(t), "some_other_node", message)
	require.Error(t, err)
}

func TestConnectionManager_CleanStartClose(t *testing.T) {
	t.Parallel()

	config, _ := newTestConfig(t, 2)
	config.ConnectionManagerConfig.HeartbeatIntervalSec = 1
	clock := clockwork.NewFakeClock()
	mgr, err := gateway.NewConnectionManager(config, clock, logger.TestLogger(t))
	require.NoError(t, err)

	err = mgr.Start(testutils.Context(t))
	require.NoError(t, err)

	err = mgr.Close()
	require.NoError(t, err)
}
