package connector

import (
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
	Id  string
	URL string
}
