package monitoring

import "io"

type FeedParser func(buf io.ReadCloser) ([]FeedConfig, error)

// FeedConfig is the interface for feed configurations extracted from the RDD.
type FeedConfig interface {
	GetName() string
	GetPath() string
	GetSymbol() string
	GetHeartbeatSec() int64
	GetContractType() string
	GetContractStatus() string
	// This functions as a feed identifier.
	GetContractAddress() string
	GetContractAddressBytes() []byte
	// Useful for mapping to kafka messages.
	ToMapping() map[string]interface{}
}
