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

	jsonrpc "github.com/smartcontractkit/chainlink-common/pkg/jsonrpc2"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/limits"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services/servicetest"
	"github.com/smartcontractkit/chainlink-common/pkg/types/core"

	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/api"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/common"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/config"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/connector"
	hc "github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/common"
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

[[Dons]]
DonId = "test_don"
HandlerName = "dummy"

[[Dons.Members]]
Address = "%s"
Name = "test_node_1"
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

func parseConnectorConfig(t *testing.T, tomlConfig string, nodeAddress string, nodeURL string) *connector.ConnectorConfig {
	nodeConfig := fmt.Sprintf(tomlConfig, nodeAddress, nodeURL)
	var cfg connector.ConnectorConfig
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
		ID:      msg.Body.MessageId,
		Result:  &rawPayload,
		Method:  req.Method,
	}
	// send back user's message without re-signing - should be ignored by the Gateway
	_ = c.connector.SendToGateway(ctx, gatewayID, resp)
	// send back a correct response
	responseMsg := &api.Message{Body: api.MessageBody{
		MessageId: msg.Body.MessageId,
		Method:    "test",
		DonId:     "test_don",
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

	// Launch Gateway
	lggr := logger.Test(t)
	gatewayConfig := fmt.Sprintf(gatewayConfigTemplate, nodeKeys.Address)
	c, err := network.NewHTTPClient(network.HTTPClientConfig{
		DefaultTimeout:   5 * time.Second,
		MaxResponseBytes: 1000,
	}, lggr)
	require.NoError(t, err)
	gateway, err := gateway.NewGatewayFromConfig(parseGatewayConfig(t, gatewayConfig), gateway.NewHandlerFactory(nil, nil, c, nil, nil, lggr, limits.Factory{Logger: lggr}), lggr)
	require.NoError(t, err)
	servicetest.Run(t, gateway)
	userPort, nodePort := gateway.GetUserPort(), gateway.GetNodePort()
	userURL := fmt.Sprintf("http://localhost:%d/user", userPort)
	nodeURL := fmt.Sprintf("ws://localhost:%d/node", nodePort)

	// Launch Connector
	client := &client{privateKey: nodeKeys.PrivateKey}
	// client acts as a signer here
	connector, err := connector.NewGatewayConnector(parseConnectorConfig(t, nodeConfigTemplate, nodeKeys.Address, nodeURL), client, clockwork.NewRealClock(), lggr)
	require.NoError(t, err)
	require.NoError(t, connector.AddHandler(t.Context(), []string{"test"}, client))
	client.connector = connector
	servicetest.Run(t, connector)

	// Send requests until one of them reaches Connector (i.e. the node)
	gomega.NewGomegaWithT(t).Eventually(func() bool {
		req := newLegacyHTTPRequestObject(t, messageID1, userURL, userKeys.PrivateKey)
		httpClient := &http.Client{}
		_, _ = httpClient.Do(req) // could initially return error if Gateway is not fully initialized yet
		return client.done.Load()
	}, testutils.WaitTimeout(t), testutils.TestInterval).Should(gomega.Equal(true))

	// Send another request and validate that response has correct content and sender
	req := newLegacyHTTPRequestObject(t, messageID2, userURL, userKeys.PrivateKey)
	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	rawResp, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	codec := api.JsonRPCCodec{}
	respMsg, err := codec.DecodeLegacyResponse(rawResp)
	require.NoError(t, err)
	require.NoError(t, respMsg.Validate())
	require.Equal(t, strings.ToLower(nodeKeys.Address), respMsg.Body.Sender)
	require.Equal(t, messageID2, respMsg.Body.MessageId)
	require.JSONEq(t, nodeResponsePayload, string(respMsg.Body.Payload))
}

func newLegacyHTTPRequestObject(t *testing.T, messageID string, userURL string, signerKey *ecdsa.PrivateKey) *http.Request {
	msg := &api.Message{Body: api.MessageBody{MessageId: messageID, Method: "test", DonId: "test_don"}}
	require.NoError(t, msg.Sign(signerKey))
	codec := api.JsonRPCCodec{}
	rawMsg, err := codec.EncodeLegacyRequest(msg)
	require.NoError(t, err)
	req, err := http.NewRequestWithContext(testutils.Context(t), "POST", userURL, bytes.NewBuffer(rawMsg))
	require.NoError(t, err)
	return req
}
