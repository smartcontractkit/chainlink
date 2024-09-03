package config

import (
	"math/big"

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
	DonId() string
	Gateways() []ConnectorGatewayConfig
	WsHandshakeTimeoutMillis() uint32
	AuthMinChallengeLen() int
	AuthTimestampToleranceSec() uint32
}

type ConnectorGatewayConfig interface {
	Id() string
	URL() string
}

type WorkflowConnectorConfig interface {
	ChainIDForNodeKey() big.Int
	GatewayConnectorConfig() GatewayConnectorConfig
}
type Capabilities interface {
	Peering() P2P
	Dispatcher() Dispatcher
	ExternalRegistry() CapabilitiesExternalRegistry
	WorkflowConnectorConfig() WorkflowConnectorConfig
}
