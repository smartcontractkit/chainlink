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
	MaxConcurrency() int
	WorkflowStorage() WorkflowStorage
	AdditionalSources() []AdditionalWorkflowSource
}

type WorkflowStorage interface {
	ArtifactStorageHost() string
	URL() string
	TLSEnabled() bool
}

// AdditionalWorkflowSource represents a single additional workflow metadata source
// that can be configured to load workflows from sources other than the on-chain registry.
type AdditionalWorkflowSource interface {
	GetURL() string
	GetTLSEnabled() bool
	GetName() string
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
	Local() LocalCapabilities
}

// LocalCapabilities provides configuration for registry-based capability launching.
type LocalCapabilities interface {
	// RegistryBasedLaunchAllowlist returns regex patterns that match capability IDs to be
	// launched from the capabilities registry instead of via job specs.
	RegistryBasedLaunchAllowlist() []string
	// Capabilities returns per-capability node configuration, keyed by capability ID.
	Capabilities() map[string]CapabilityNodeConfig
	// IsAllowlisted returns true if the capability ID matches any pattern in the allowlist.
	IsAllowlisted(capabilityID string) bool
	// GetCapabilityConfig returns the node config for a specific capability, or nil if not configured.
	GetCapabilityConfig(capabilityID string) CapabilityNodeConfig
}

// CapabilityNodeConfig provides node-specific configuration for a capability.
type CapabilityNodeConfig interface {
	// BinaryPathOverride returns the override path for the capability binary, or empty if not set.
	BinaryPathOverride() string
	// Config returns capability-specific configuration as key-value pairs.
	Config() map[string]string
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
