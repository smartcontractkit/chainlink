package config

import "time"

type CRE interface {
	WsURL() string
	RestURL() string
	StreamsAPIKey() string
	StreamsAPISecret() string
	WorkflowFetcher() WorkflowFetcher
	UseLocalTimeProvider() bool
	EnableDKGRecipient() bool
	Linking() CRELinking
	// DebugMode returns true if debug mode is enabled for workflow engines.
	// When enabled, additional OTel tracing and logging is performed.
	DebugMode() bool
	// CachedTriggerSubscriptionsEnabled returns true if workflow engines should
	// reuse a previously-persisted trigger subscription payload
	// (workflow_specs_v2.trigger_subscriptions) instead of executing the
	// workflow's WASM Subscribe() call on every engine start.
	CachedTriggerSubscriptionsEnabled() bool
	LocalSecretOverrides() map[string]map[string]string
	ConfidentialRelay() CREConfidentialRelay
}

// WorkflowFetcher defines configuration for fetching workflow files
type WorkflowFetcher interface {
	// URL returns the configured URL for fetching workflow files
	URL() string
}

// CREConfidentialRelay defines configuration for the confidential relay handler.
type CREConfidentialRelay interface {
	Enabled() bool
	// TrustEnclaves reports whether the relay should trust fake (non-Nitro)
	// enclaves by relaxing TEE attestation validation. INSECURE; test-only.
	TrustEnclaves() bool
	// RequireBFTQuorum selects the required signature quorum
	RequireBFTQuorum() bool
}

// CRELinking defines configuration for connecting to the CRE linking service
type CRELinking interface {
	URL() string
	TLSEnabled() bool
	// RequestTimeout bounds each organization lookup against the linking service.
	RequestTimeout() time.Duration
	// DurableCacheEnabled turns on durable caching of owner->orgID mappings (backed by Postgres).
	DurableCacheEnabled() bool
}
