package integration_tests

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"fmt"
	"net/http"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/onsi/gomega"
	"github.com/pelletier/go-toml/v2"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/api"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/common"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/config"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/connector"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
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
ReadTimeoutMillis = 100
RequestTimeoutMillis = 100
WriteTimeoutMillis = 100

[UserServerConfig]
Path = "/user"
Port = 0
ContentTypeHeader = "application/jsonrpc"
MaxRequestBytes = 20_000
ReadTimeoutMillis = 100
RequestTimeoutMillis = 100
WriteTimeoutMillis = 100

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

	nodeKeys := common.NewTestNodes(t, 1)[0]
	// Verify that addresses in config are case-insensitive
	nodeKeys.Address = strings.ToUpper(nodeKeys.Address)

	// Launch Gateway
	lggr := logger.TestLogger(t)
	gatewayConfig := fmt.Sprintf(gatewayConfigTemplate, nodeKeys.Address)
	gateway, err := gateway.NewGatewayFromConfig(parseGatewayConfig(t, gatewayConfig), gateway.NewHandlerFactory(nil, lggr), lggr)
	require.NoError(t, err)
	require.NoError(t, gateway.Start(testutils.Context(t)))
	userPort, nodePort := gateway.GetUserPort(), gateway.GetNodePort()
	userUrl := fmt.Sprintf("http://localhost:%d/user", userPort)
	nodeUrl := fmt.Sprintf("ws://localhost:%d/node", nodePort)

	// Launch Connector
	client := &client{privateKey: nodeKeys.PrivateKey}
	connector, err := connector.NewGatewayConnector(parseConnectorConfig(t, nodeConfigTemplate, nodeKeys.Address, nodeUrl), client, client, utils.NewRealClock(), lggr)
	require.NoError(t, err)
	client.connector = connector
	require.NoError(t, connector.Start(testutils.Context(t)))

	// Send requests until one of them reaches Connector
	gomega.NewGomegaWithT(t).Eventually(func() bool {
		msg := &api.Message{Body: api.MessageBody{MessageId: "123", Method: "test", DonId: "test_don"}}
		require.NoError(t, msg.Sign(nodeKeys.PrivateKey))
		codec := api.JsonRPCCodec{}
		rawMsg, err := codec.EncodeRequest(msg)
		require.NoError(t, err)
		req, err := http.NewRequestWithContext(testutils.Context(t), "POST", userUrl, bytes.NewBuffer(rawMsg))
		require.NoError(t, err)
		httpClient := &http.Client{}
		_, _ = httpClient.Do(req) // could initially return error if Gateway is not fully initialized yet
		return client.done.Load()
	}, testutils.WaitTimeout(t), testutils.TestInterval).Should(gomega.Equal(true))

	require.NoError(t, connector.Close())
	require.NoError(t, gateway.Close())
}
