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
