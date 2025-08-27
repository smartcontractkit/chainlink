package config

type CRE interface {
	WsURL() string
	RestURL() string
	StreamsAPIKey() string
	StreamsAPISecret() string
	WorkflowFetcher() WorkflowFetcher
	UseLocalTimeProvider() bool
}

// WorkflowFetcher defines configuration for fetching workflow files
type WorkflowFetcher interface {
	// URL returns the configured URL for fetching workflow files
	URL() string
}
