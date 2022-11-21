package monitoring

import "time"

// ChainConfig contains chain-specific configuration.
// It is an interface so that implementations can add extra fields as long as
// they provide data from these methods which are required by the framework.
type ChainConfig interface {
	GetRPCEndpoint() string
	GetNetworkName() string
	GetNetworkID() string
	GetChainID() string
	GetReadTimeout() time.Duration
	GetPollInterval() time.Duration
	// Useful for serializing to avro.
	// Check the latest version of the transmission schema to see what the exact return format should be.
	ToMapping() map[string]interface{}
}
