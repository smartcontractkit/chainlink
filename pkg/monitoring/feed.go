package monitoring

import "io"

type FeedParser func(buf io.ReadCloser) ([]FeedConfig, error)

// FeedConfig is the interface for feed configurations extracted from the RDD.
type FeedConfig interface {
	// This functions as a feed identifier.
	GetID() string
	GetName() string
	GetPath() string
	GetSymbol() string
	GetHeartbeatSec() int64
	GetContractType() string
	GetContractStatus() string
	GetContractAddress() string
	GetContractAddressBytes() []byte
	GetMultiply() uint64
	// Useful for mapping to kafka messages.
	ToMapping() map[string]interface{}
}
