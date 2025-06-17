package config

import (
	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

type CapabilitiesExternalRegistry interface {
	Address() string
	NetworkID() string
	ChainID() string
	RelayID() types.RelayID
}

type EngineExecutionRateLimit interface {
	GlobalRPS() float64
	GlobalBurst() int
	PerSenderRPS() float64
	PerSenderBurst() int
}

type CapabilitiesWorkflowRegistry interface {
	Address() string
	NetworkID() string
	ChainID() string
	MaxEncryptedSecretsSize() utils.FileSize
	MaxBinarySize() utils.FileSize
	MaxConfigSize() utils.FileSize
	RelayID() types.RelayID
	SyncStrategy() string
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
	RateLimit() EngineExecutionRateLimit
	Peering() P2P
	Dispatcher() Dispatcher
	ExternalRegistry() CapabilitiesExternalRegistry
	WorkflowRegistry() CapabilitiesWorkflowRegistry
	GatewayConnector() GatewayConnector
}
