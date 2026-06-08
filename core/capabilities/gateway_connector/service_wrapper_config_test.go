package gatewayconnector

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/config"
)

type mockConnectorGateway struct {
	id    string
	donID string
	url   string
}

func (m mockConnectorGateway) ID() string    { return m.id }
func (m mockConnectorGateway) DonID() string { return m.donID }
func (m mockConnectorGateway) URL() string   { return m.url }

type mockGatewayConnector struct {
	nodeAddress               string
	donID                     string
	gateways                  []config.ConnectorGateway
	wsHandshakeTimeoutMillis  uint32
	authMinChallengeLen       int
	authTimestampToleranceSec uint32
}

func (m mockGatewayConnector) ChainIDForNodeKey() string           { return "1" }
func (m mockGatewayConnector) NodeAddress() string                 { return m.nodeAddress }
func (m mockGatewayConnector) DonID() string                       { return m.donID }
func (m mockGatewayConnector) Gateways() []config.ConnectorGateway { return m.gateways }
func (m mockGatewayConnector) WSHandshakeTimeoutMillis() uint32    { return m.wsHandshakeTimeoutMillis }
func (m mockGatewayConnector) AuthMinChallengeLen() int            { return m.authMinChallengeLen }
func (m mockGatewayConnector) AuthTimestampToleranceSec() uint32   { return m.authTimestampToleranceSec }

func TestTranslateConfigs(t *testing.T) {
	t.Parallel()

	translated := translateConfigs(mockGatewayConnector{
		nodeAddress:               "0x68902d681c28119f9b2531473a417088bf008e59",
		donID:                     "example_don",
		wsHandshakeTimeoutMillis:  100,
		authMinChallengeLen:       10,
		authTimestampToleranceSec: 5,
		gateways: []config.ConnectorGateway{
			mockConnectorGateway{id: "example_gateway", donID: "example_gateway_don", url: "wss://localhost:8081/node"},
			mockConnectorGateway{id: "another_gateway", donID: "another_gateway_don", url: "wss://example.com:8090/node"},
		},
	})

	assert.Equal(t, "0x68902d681c28119f9b2531473a417088bf008e59", translated.NodeAddress)
	assert.Equal(t, "example_don", translated.DonId)
	assert.Equal(t, uint32(100), translated.WsClientConfig.HandshakeTimeoutMillis)
	assert.Equal(t, 10, translated.AuthMinChallengeLen)
	assert.Equal(t, uint32(5), translated.AuthTimestampToleranceSec)
	require.Len(t, translated.Gateways, 2)
	assert.Equal(t, "example_gateway", translated.Gateways[0].ID)
	assert.Equal(t, "example_gateway_don", translated.Gateways[0].DonID)
	assert.Equal(t, "wss://localhost:8081/node", translated.Gateways[0].URL)
	assert.Equal(t, "another_gateway", translated.Gateways[1].ID)
	assert.Equal(t, "another_gateway_don", translated.Gateways[1].DonID)
	assert.Equal(t, "wss://example.com:8090/node", translated.Gateways[1].URL)
}
