package monitoring

import "time"

type ChainConfig interface {
	GetRPCEndpoint() string
	GetNetworkName() string
	GetNetworkID() string
	GetChainID() string
	GetReadTimeout() time.Duration
	GetPollInterval() time.Duration
	// Useful for serializing to avro.
	ToMapping() map[string]interface{}
}
