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

type GatewayConnector interface {
	ChainIDForNodeKey() string
	NodeAddress() string
	DonID() string
	Gateways() []ConnectorGateway
	WSHandshakeTimeoutMillis() uint32
	AuthMinChallengeLen() int
	AuthTimestampToleranceSec() uint32
}

type ConnectorGateway interface {
	ID() string
	URL() string
}

type Capabilities interface {
	Peering() P2P
	Dispatcher() Dispatcher
	ExternalRegistry() CapabilitiesExternalRegistry
	GatewayConnector() GatewayConnector
}
