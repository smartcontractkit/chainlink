package integration_tests

import (
	"bytes"
	"context"
	"crypto/ecdsa"
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

	"github.com/smartcontractkit/chainlink-common/pkg/services/servicetest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/api"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/common"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/config"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/connector"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/network"
)

const gatewayConfigTemplate = `
[ConnectionManagerConfig]
AuthChallengeLen = 32
AuthGatewayId = "test_gateway"
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
	messageId1 = "123"
	messageId2 = "456"

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
	connector  connector.GatewayConnector
	done       atomic.Bool
}

func (c *client) HandleGatewayMessage(ctx context.Context, gatewayId string, msg *api.Message) {
	c.done.Store(true)
	// send back user's message without re-signing - should be ignored by the Gateway
	_ = c.connector.SendToGateway(ctx, gatewayId, msg)
	// send back a correct response
	responseMsg := &api.Message{Body: api.MessageBody{
		MessageId: msg.Body.MessageId,
		Method:    "test",
		DonId:     "test_don",
		Receiver:  msg.Body.Sender,
		Payload:   []byte(nodeResponsePayload),
	}}
	err := responseMsg.Sign(c.privateKey)
	if err != nil {
		panic(err)
	}
	_ = c.connector.SendToGateway(ctx, gatewayId, responseMsg)
}

func (c *client) Sign(data ...[]byte) ([]byte, error) {
	return common.SignData(c.privateKey, data...)
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
	lggr := logger.TestLogger(t)
	gatewayConfig := fmt.Sprintf(gatewayConfigTemplate, nodeKeys.Address)
	c, err := network.NewHTTPClient(network.HTTPClientConfig{
		DefaultTimeout:   5 * time.Second,
		MaxResponseBytes: 1000,
	}, lggr)
	require.NoError(t, err)
	gateway, err := gateway.NewGatewayFromConfig(parseGatewayConfig(t, gatewayConfig), gateway.NewHandlerFactory(nil, nil, c, lggr), lggr)
	require.NoError(t, err)
	servicetest.Run(t, gateway)
	userPort, nodePort := gateway.GetUserPort(), gateway.GetNodePort()
	userUrl := fmt.Sprintf("http://localhost:%d/user", userPort)
	nodeUrl := fmt.Sprintf("ws://localhost:%d/node", nodePort)

	// Launch Connector
	client := &client{privateKey: nodeKeys.PrivateKey}
	// client acts as a signer here
	connector, err := connector.NewGatewayConnector(parseConnectorConfig(t, nodeConfigTemplate, nodeKeys.Address, nodeUrl), client, clockwork.NewRealClock(), lggr)
	require.NoError(t, err)
	require.NoError(t, connector.AddHandler([]string{"test"}, client))
	client.connector = connector
	servicetest.Run(t, connector)

	// Send requests until one of them reaches Connector (i.e. the node)
	gomega.NewGomegaWithT(t).Eventually(func() bool {
		req := newHttpRequestObject(t, messageId1, userUrl, userKeys.PrivateKey)
		httpClient := &http.Client{}
		_, _ = httpClient.Do(req) // could initially return error if Gateway is not fully initialized yet
		return client.done.Load()
	}, testutils.WaitTimeout(t), testutils.TestInterval).Should(gomega.Equal(true))

	// Send another request and validate that response has correct content and sender
	req := newHttpRequestObject(t, messageId2, userUrl, userKeys.PrivateKey)
	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	rawResp, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	codec := api.JsonRPCCodec{}
	respMsg, err := codec.DecodeResponse(rawResp)
	require.NoError(t, err)
	require.NoError(t, respMsg.Validate())
	require.Equal(t, strings.ToLower(nodeKeys.Address), respMsg.Body.Sender)
	require.Equal(t, messageId2, respMsg.Body.MessageId)
	require.Equal(t, nodeResponsePayload, string(respMsg.Body.Payload))
}

func newHttpRequestObject(t *testing.T, messageId string, userUrl string, signerKey *ecdsa.PrivateKey) *http.Request {
	msg := &api.Message{Body: api.MessageBody{MessageId: messageId, Method: "test", DonId: "test_don"}}
	require.NoError(t, msg.Sign(signerKey))
	codec := api.JsonRPCCodec{}
	rawMsg, err := codec.EncodeRequest(msg)
	require.NoError(t, err)
	req, err := http.NewRequestWithContext(testutils.Context(t), "POST", userUrl, bytes.NewBuffer(rawMsg))
	require.NoError(t, err)
	return req
}
