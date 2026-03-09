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
}

// WorkflowFetcher defines configuration for fetching workflow files
type WorkflowFetcher interface {
	// URL returns the configured URL for fetching workflow files
	URL() string
}

// CRELinking defines configuration for connecting to the CRE linking service
type CRELinking interface {
	URL() string
	TLSEnabled() bool
}
