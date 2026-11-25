package integration_tests

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/jonboulle/clockwork"
	"github.com/onsi/gomega"
	"github.com/pelletier/go-toml/v2"
	"github.com/stretchr/testify/require"

	jsonrpc "github.com/smartcontractkit/chainlink-common/pkg/jsonrpc2"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services/servicetest"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/limits"
	"github.com/smartcontractkit/chainlink-common/pkg/types/core"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/api"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/common"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/config"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/connector"
	hc "github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/common"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/internal/netfault"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/network"
)

const gatewayConfigTemplate = `
[ConnectionManagerConfig]
AuthChallengeLen = 32
AuthGatewayID = "test_gateway"
AuthTimestampToleranceSec = 30

[NodeServerConfig]
Path = "/node"
Port = 0
HandshakeTimeoutMillis = 2_000
MaxRequestBytes = 20_000
ReadTimeoutMillis = 1000
RequestTimeoutMillis = 1000
WriteTimeoutMillis = 1000

[UserServerConfig]
Path = "/user"
Port = 0
ContentTypeHeader = "application/jsonrpc"
MaxRequestBytes = 20_000
ReadTimeoutMillis = 1000
RequestTimeoutMillis = 1000
WriteTimeoutMillis = 1000
CORSEnabled = true
CORSAllowedOrigins = ["https://remix.ethereum.org"]

[[shardedDONs]]
DonName = "test_don"
F = 0

[[shardedDONs.Shards]]
[[shardedDONs.Shards.Nodes]]
Address = "%s"
Name = "test_node_1"

[[services]]
ServiceName = "test"
DONs = ["test_don"]

[[services.Handlers]]
Name = "dummy"
`

const nodeConfigTemplate = `
DonID = "test_don"
AuthMinChallengeLen = 32
AuthTimestampToleranceSec = 30
NodeAddress = "%s"

[WsClientConfig]
HandshakeTimeoutMillis = 2_000

[[Gateways]]
Id = "test_gateway"
URL = "%s"
`

const (
	messageID1 = "123"
	messageID2 = "456"

	nodeResponsePayload = `{"response":"correct response"}`
)

func parseGatewayConfig(t *testing.T, tomlConfig string) *config.GatewayConfig {
	var cfg config.GatewayConfig
	err := toml.Unmarshal([]byte(tomlConfig), &cfg)
	require.NoError(t, err)
	return &cfg
}

func parseConnectorConfig(t *testing.T, tomlConfig, nodeAddress, nodeURL string) *connector.Config {
	nodeConfig := fmt.Sprintf(tomlConfig, nodeAddress, nodeURL)
	var cfg connector.Config
	require.NoError(t, toml.Unmarshal([]byte(nodeConfig), &cfg))
	return &cfg
}

type client struct {
	privateKey *ecdsa.PrivateKey
	connector  core.GatewayConnector
	done       atomic.Bool
}

func (c *client) HandleGatewayMessage(ctx context.Context, gatewayID string, req *jsonrpc.Request[json.RawMessage]) error {
	msg, err := hc.ValidatedMessageFromReq(req)
	if err != nil {
		panic(err)
	}
	c.done.Store(true)
	payload, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	rawPayload := json.RawMessage(payload)
	resp := &jsonrpc.Response[json.RawMessage]{
		Version: "2.0",
		ID:      msg.Body.MessageID,
		Result:  &rawPayload,
		Method:  req.Method,
	}
	// send back user's message without re-signing - should be ignored by the Gateway
	_ = c.connector.SendToGateway(ctx, gatewayID, resp)
	// send back a correct response
	responseMsg := &api.Message{Body: api.MessageBody{
		MessageID: msg.Body.MessageID,
		Method:    "test",
		DonID:     "test_don",
		Receiver:  msg.Body.Sender,
		Payload:   []byte(nodeResponsePayload),
	}}
	err = responseMsg.Sign(c.privateKey)
	if err != nil {
		panic(err)
	}
	resp, err = hc.ValidatedResponseFromMessage(responseMsg) // ensure the message is valid
	if err != nil {
		panic(err)
	}
	return c.connector.SendToGateway(ctx, gatewayID, resp)
}

func (c *client) Sign(ctx context.Context, data ...[]byte) ([]byte, error) {
	return common.SignData(c.privateKey, data...)
}

func (c *client) ID(ctx context.Context) (string, error) {
	return "test_client", nil
}

func (*client) Start(ctx context.Context) error {
	return nil
}

func (*client) Close() error {
	return nil
}

func TestIntegration_Gateway_NoFullNodes_BasicConnectionAndMessage(t *testing.T) {
	t.Parallel()

	testWallets := common.NewTestNodes(t, 2)
	nodeKeys := testWallets[0]
	userKeys := testWallets[1]
	// Verify that addresses in config are case-insensitive
	nodeKeys.Address = strings.ToUpper(nodeKeys.Address)

	gw := startGateway(t, nodeKeys)
	userPort, nodePort := gw.GetUserPort(), gw.GetNodePort()
	userURL := fmt.Sprintf("http://localhost:%d/user", userPort)
	require.Equal(t, http.StatusServiceUnavailable, getHTTPStatus(t, fmt.Sprintf("http://localhost:%d/health", userPort)))
	require.Equal(t, http.StatusOK, getHTTPStatus(t, fmt.Sprintf("http://localhost:%d/health", nodePort)))

	node := startNode(t, nodeKeys, fmt.Sprintf("ws://localhost:%d/node", nodePort), nil)
	require.Eventually(t, func() bool {
		return getHTTPStatus(t, fmt.Sprintf("http://localhost:%d/health", userPort)) == http.StatusOK
	}, testutils.WaitTimeout(t), testutils.TestInterval)

	// Send requests until one of them reaches Connector (i.e. the node)
	gomega.NewGomegaWithT(t).Eventually(func() bool {
		_, _ = sendUserRequest(t, userURL, userKeys, messageID1) // could initially fail if Gateway is not fully initialized yet
		return node.done.Load()
	}, testutils.WaitTimeout(t), testutils.TestInterval).Should(gomega.BeTrue())

	// Send another request and validate that response has correct content and sender
	respMsg, err := sendUserRequest(t, userURL, userKeys, messageID2)
	require.NoError(t, err)
	requireNodeResponse(t, respMsg, nodeKeys, messageID2)
}

// The node must notice a gateway that went silent without closing its connection,
// and serve requests again once it has redialed.
func TestIntegration_Gateway_NodeRecoversFromVanishedGatewayConnection(t *testing.T) {
	t.Parallel()

	testWallets := common.NewTestNodes(t, 2)
	nodeKeys, userKeys := testWallets[0], testWallets[1]

	gw := startGateway(t, nodeKeys)
	userURL := fmt.Sprintf("http://localhost:%d/user", gw.GetUserPort())
	proxy := netfault.NewProxy(t, fmt.Sprintf("localhost:%d", gw.GetNodePort()))
	startNode(t, nodeKeys, fmt.Sprintf("ws://%s/node", proxy.Addr()), func(cfg *connector.Config) {
		cfg.PingInterval, cfg.PongTimeout, cfg.WriteTimeout = 100*time.Millisecond, 500*time.Millisecond, 500*time.Millisecond
	})
	requireServed := func(within time.Duration) {
		t.Helper()
		var respMsg *api.Message
		require.Eventually(t, func() bool {
			var err error
			respMsg, err = sendUserRequest(t, userURL, userKeys, messageID2)
			return err == nil
		}, within, testutils.TestInterval)
		requireNodeResponse(t, respMsg, nodeKeys, messageID2)
	}
	requireServed(testutils.WaitTimeout(t))

	proxy.Blackhole()
	proxy.Restore()

	requireServed(5 * time.Second)
}

func startGateway(t *testing.T, nodeKeys common.TestNode) gateway.Gateway {
	t.Helper()
	lggr := logger.Test(t)
	httpClient, err := network.NewHTTPClient(network.HTTPClientConfig{
		DefaultTimeout:   5 * time.Second,
		MaxResponseBytes: 1000,
	}, lggr)
	require.NoError(t, err)
	lf := limits.Factory{Logger: lggr}
	gw, err := gateway.NewGatewayFromConfig(
		parseGatewayConfig(t, fmt.Sprintf(gatewayConfigTemplate, nodeKeys.Address)),
		gateway.NewHandlerFactory(nil, nil, httpClient, nil, nil, lggr, lf, nil, nil), lggr, lf)
	require.NoError(t, err)
	servicetest.Run(t, gw)
	return gw
}

// startNode runs a connector that dials nodeURL. tweak may adjust its config.
func startNode(t *testing.T, nodeKeys common.TestNode, nodeURL string, tweak func(*connector.Config)) *client {
	t.Helper()
	lggr := logger.Test(t)
	cfg := parseConnectorConfig(t, nodeConfigTemplate, nodeKeys.Address, nodeURL)
	if tweak != nil {
		tweak(cfg)
	}
	// client acts as a signer here
	node := &client{privateKey: nodeKeys.PrivateKey}
	conn, err := connector.NewGatewayConnector(cfg, node, clockwork.NewRealClock(), lggr, "")
	require.NoError(t, err)
	require.NoError(t, conn.AddHandler(t.Context(), []string{"test"}, node))
	node.connector = conn
	servicetest.Run(t, conn)
	return node
}

// sendUserRequest returns the node's response, or an error if the gateway could not get one.
func sendUserRequest(t *testing.T, userURL string, userKeys common.TestNode, messageID string) (*api.Message, error) {
	t.Helper()
	resp, err := (&http.Client{}).Do(newJSONRPCHTTPRequestObject(t, messageID, userURL, userKeys.PrivateKey))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status %d", resp.StatusCode)
	}
	rawResp, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	codec := api.JSONRPCCodec{}
	return codec.DecodeLegacyResponse(rawResp)
}

func requireNodeResponse(t *testing.T, respMsg *api.Message, nodeKeys common.TestNode, messageID string) {
	t.Helper()
	require.NoError(t, respMsg.Validate())
	require.Equal(t, strings.ToLower(nodeKeys.Address), respMsg.Body.Sender)
	require.Equal(t, messageID, respMsg.Body.MessageID)
	require.JSONEq(t, nodeResponsePayload, string(respMsg.Body.Payload))
}

func getHTTPStatus(t *testing.T, url string) int {
	t.Helper()
	req, err := http.NewRequestWithContext(t.Context(), http.MethodGet, url, nil)
	require.NoError(t, err)
	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	require.NoError(t, resp.Body.Close())
	return resp.StatusCode
}

func newJSONRPCHTTPRequestObject(t *testing.T, messageID, userURL string, signerKey *ecdsa.PrivateKey) *http.Request {
	msg := &api.Message{Body: api.MessageBody{MessageID: messageID, Method: "test"}}
	require.NoError(t, msg.Sign(signerKey))
	msgBytes, err := json.Marshal(msg)
	require.NoError(t, err)
	rawParams := json.RawMessage(msgBytes)
	request := jsonrpc.Request[json.RawMessage]{
		Version: jsonrpc.JsonRpcVersion,
		ID:      messageID,
		Method:  "test.test",
		Params:  &rawParams,
	}
	rawMsg, err := json.Marshal(&request)
	require.NoError(t, err)
	req, err := http.NewRequestWithContext(t.Context(), http.MethodPost, userURL, bytes.NewBuffer(rawMsg))
	require.NoError(t, err)
	return req
}
