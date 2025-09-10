package config

import (
	ocrcommontypes "github.com/smartcontractkit/libocr/commontypes"

	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

type CapabilitiesExternalRegistry interface {
	Address() string
	NetworkID() string
	ChainID() string
	ContractVersion() string
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
	ContractVersion() string
	MaxEncryptedSecretsSize() utils.FileSize
	MaxBinarySize() utils.FileSize
	MaxConfigSize() utils.FileSize
	RelayID() types.RelayID
	SyncStrategy() string
	WorkflowStorage() WorkflowStorage
}

type WorkflowStorage interface {
	ArtifactStorageHost() string
	URL() string
	TLSEnabled() bool
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
	SharedPeering() SharedPeering
	Dispatcher() Dispatcher
	ExternalRegistry() CapabilitiesExternalRegistry
	WorkflowRegistry() CapabilitiesWorkflowRegistry
	GatewayConnector() GatewayConnector
}

type SharedPeering interface {
	Enabled() bool
	Bootstrappers() (locators []ocrcommontypes.BootstrapperLocator)
	StreamConfig() StreamConfig
}

type StreamConfig interface {
	IncomingMessageBufferSize() int
	OutgoingMessageBufferSize() int
	MaxMessageLenBytes() int
	MessageRateLimiterRate() float64
	MessageRateLimiterCapacity() uint32
	BytesRateLimiterRate() float64
	BytesRateLimiterCapacity() uint32
}
