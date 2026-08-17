package connector

import (
	"github.com/smartcontractkit/chainlink/v2/core/config"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/network"
)

type ConnectorConfig struct {
	NodeAddress               string
	DonId                     string
	Gateways                  []ConnectorGatewayConfig
	WsClientConfig            network.WebSocketClientConfig
	AuthMinChallengeLen       int
	AuthTimestampToleranceSec uint32
}

type ConnectorGatewayConfig struct {
	ID    string `toml:"Id"`
	DonID string `toml:"DonId"`
	URL   string
}

func (ConnectorConfig) From(c config.GatewayConnector) ConnectorConfig {
	r := ConnectorConfig{
		NodeAddress:               c.NodeAddress(),
		DonId:                     c.DonID(),
		WsClientConfig:            network.WebSocketClientConfig{HandshakeTimeoutMillis: c.WSHandshakeTimeoutMillis()},
		AuthMinChallengeLen:       c.AuthMinChallengeLen(),
		AuthTimestampToleranceSec: c.AuthTimestampToleranceSec(),
	}

	// AuthTimestampToleranceSec should be at least 5 seconds
	// Under network instability nodes may lose connection to the gateway
	// having a smaller tolerance range would reduce the chances of success on retries.
	// this is what happened on https://smartcontract-it.atlassian.net/browse/INCIDENT-2556
	if r.AuthTimestampToleranceSec < 5 {
		r.AuthTimestampToleranceSec = 5
	}

	if len(c.Gateways()) != 0 {
		r.Gateways = make([]ConnectorGatewayConfig, len(c.Gateways()))
		for i, gateway := range c.Gateways() {
			r.Gateways[i] = ConnectorGatewayConfig{
				ID:    gateway.ID(),
				DonID: gateway.DonID(),
				URL:   gateway.URL(),
			}
		}
	}

	return r
}
