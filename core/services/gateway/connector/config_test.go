package connector_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/config/toml"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/connector"
)

func TestConnectorConfig_From(t *testing.T) {
	t.Parallel()

	cfg, err := chainlink.GeneralConfigOpts{
		Config: chainlink.Config{
			Core: toml.Core{
				Capabilities: toml.Capabilities{
					GatewayConnector: toml.GatewayConnector{
						ChainIDForNodeKey:         new("1"),
						NodeAddress:               new("0x68902d681c28119f9b2531473a417088bf008e59"),
						DonID:                     new("example_don"),
						WSHandshakeTimeoutMillis:  new(uint32(100)),
						AuthMinChallengeLen:       new(10),
						AuthTimestampToleranceSec: new(uint32(5)),
						Gateways: []toml.ConnectorGateway{
							{
								ID:    new("example_gateway"),
								DonID: new("example_gateway_don"),
								URL:   new("wss://localhost:8081/node"),
							},
							{
								ID:    new("another_gateway"),
								DonID: new("another_gateway_don"),
								URL:   new("wss://example.com:8090/node"),
							},
						},
					},
				},
			},
		},
	}.New()
	require.NoError(t, err)

	translated := connector.ConnectorConfig{}.From(cfg.Capabilities().GatewayConnector())

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
