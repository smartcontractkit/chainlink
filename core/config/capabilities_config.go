package config

import (
	"github.com/smartcontractkit/chainlink-common/pkg/types"
)

type CapabilitiesExternalRegistry interface {
	Address() string
	NetworkID() string
	ChainID() string
	RelayID() types.RelayID
}

type GatewayConnectorConfig interface {
	NodeAddress() string
	DonID() string
	Gateways() []ConnectorGatewayConfig
	WsHandshakeTimeoutMillis() uint32
	AuthMinChallengeLen() int
	AuthTimestampToleranceSec() uint32
}

type ConnectorGatewayConfig interface {
	ID() string
	URL() string
}

type WorkflowConnectorConfig interface {
	ChainIDForNodeKey() string
	GatewayConnectorConfig() GatewayConnectorConfig
}
type Capabilities interface {
	Peering() P2P
	Dispatcher() Dispatcher
	ExternalRegistry() CapabilitiesExternalRegistry
	WorkflowConnectorConfig() WorkflowConnectorConfig
}
