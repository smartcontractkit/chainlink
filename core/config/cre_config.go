package config

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
	LocalSecrets() map[string]string
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
}

// CRELinking defines configuration for connecting to the CRE linking service
type CRELinking interface {
	URL() string
	TLSEnabled() bool
}
