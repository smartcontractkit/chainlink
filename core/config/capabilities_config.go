package config

import (
	"time"

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
	MaxActivationRetries() int
	WorkflowStorage() WorkflowStorage
	ModuleCache() ModuleCache
	AdditionalSources() []AdditionalWorkflowSource
}

type WorkflowStorage interface {
	ArtifactStorageHost() string
	URL() string
	TLSEnabled() bool
}

type ModuleCache interface {
	Enabled() bool
	DiskMonitorEnabled() bool
	IdleEviction() bool
	IdleTimeout() time.Duration
	MaxLoaded() int
	CacheDir() string
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
	DonID() string
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
	HTTPTrigger() HTTPTriggerCapability
	HTTPAction() HTTPActionCapability
}

// HTTPTriggerCapability provides node TOML configuration for the http-trigger
// capability (http-trigger@1.0.0-alpha). Values set here are injected into the
// capability's service config at runtime; node TOML is the authoritative source.
type HTTPTriggerCapability interface {
	// MetadataBatchSize is the number of metadata items sent in a single batch to the gateway.
	// Returns 0 when unset (the capability applies its default of 50).
	MetadataBatchSize() uint16
	// SendChannelBufferSize is the size of the channel used to trigger workflows.
	// Returns 0 when unset (the capability applies its default of 1000).
	SendChannelBufferSize() uint16
	// MaxAuthorizedKeysPerWorkflow limits keys registered per workflow.
	// Returns 0 when unset (the capability applies its default of 100).
	MaxAuthorizedKeysPerWorkflow() uint16
	// RequestCacheTTL is the time-to-live for cached request responses in seconds.
	// Returns 0 when unset (the capability applies its default of 86400).
	RequestCacheTTL() uint32
	// GatewayConnection holds capability-specific gateway connection tuning.
	GatewayConnection() HTTPTriggerGatewayConnection
}

// HTTPTriggerGatewayConnection holds gateway connection tuning for the http-trigger capability.
type HTTPTriggerGatewayConnection interface {
	// RetryConfig configures the exponential backoff retry strategy.
	RetryConfig() HTTPTriggerRetryConfig
	// MaxPushMetadataDurationMs is the max duration in ms for broadcasting metadata to the gateway.
	// Returns 0 when unset (the capability applies its default of 30000).
	MaxPushMetadataDurationMs() uint32
	// MaxPullMetadataDurationMs is the max duration in ms for responding to pull metadata from the gateway.
	// Returns 0 when unset (the capability applies its default of 30000).
	MaxPullMetadataDurationMs() uint32
}

// HTTPTriggerRetryConfig configures the exponential backoff retry strategy.
type HTTPTriggerRetryConfig interface {
	// InitialIntervalMs is the initial retry interval in milliseconds.
	// Returns 0 when unset (the capability applies its default of 100).
	InitialIntervalMs() int
	// MaxIntervalTimeMs is the maximum retry interval in milliseconds.
	// Returns 0 when unset (the capability applies its default of 30000).
	MaxIntervalTimeMs() int
	// Multiplier is the backoff multiplier applied between retries.
	// Returns 0 when unset (the capability applies its default of 2.0).
	Multiplier() float64
}

// HTTPActionCapability provides node TOML configuration for the http-action
// capability (http-actions@1.0.0-alpha). Values set here are injected into the
// capability's service config at runtime; node TOML is the authoritative source.
type HTTPActionCapability interface {
	// ProxyMode is the outbound proxy mode: "gateway" or "direct".
	// Returns "" when unset (the capability applies its default of "gateway").
	ProxyMode() string
	// GatewayConnection holds capability-specific gateway connection tuning.
	GatewayConnection() HTTPActionGatewayConnection
	// HTTPClient configures the HTTP client used in "direct" mode (no Gateway).
	// These network restrictions are potentially sensitive and are never emitted
	// into job specs; they are only read from node TOML.
	HTTPClient() HTTPActionHTTPClient
}

// HTTPActionGatewayConnection holds gateway connection tuning for the http-action capability.
type HTTPActionGatewayConnection interface {
	// InitialIntervalMs is the initial interval in milliseconds for the exponential backoff retry strategy.
	// Returns 0 when unset (the capability applies its default of 100).
	InitialIntervalMs() uint32
	// MaxElapsedTimeMs is the maximum elapsed time in milliseconds for the exponential backoff retry strategy.
	// Returns 0 when unset (the capability applies its default of 30000).
	MaxElapsedTimeMs() uint32
	// Multiplier is the multiplier for the exponential backoff retry strategy.
	// Returns 0 when unset (the capability applies its default of 2.0).
	Multiplier() float64
}

// HTTPActionHTTPClient configures the HTTP client used in "direct" mode.
type HTTPActionHTTPClient interface {
	// BlockedIPs is a list of IP addresses that are not allowed to be accessed.
	BlockedIPs() []string
	// BlockedIPsCIDR is a list of CIDR blocks that are not allowed to be accessed.
	BlockedIPsCIDR() []string
	// AllowedPorts is a list of ports that are allowed for outgoing HTTP requests.
	// Returns nil when unset (the capability applies its default of [443]).
	AllowedPorts() []int
	// AllowedSchemes is a list of URL schemes (e.g., "http", "https") that are allowed.
	// Returns nil when unset (the capability applies its default of ["https"]).
	AllowedSchemes() []string
	// AllowedIPs is a list of IP addresses that are explicitly allowed to be accessed.
	AllowedIPs() []string
	// AllowedIPsCIDR is a list of CIDR blocks that are explicitly allowed to be accessed.
	AllowedIPsCIDR() []string
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
