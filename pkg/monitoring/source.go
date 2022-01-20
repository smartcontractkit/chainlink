package monitoring

import (
	"context"
	"math/big"
	"time"

	"github.com/smartcontractkit/libocr/offchainreporting2/types"
)

type Source interface {
	Fetch(context.Context) (interface{}, error)
}

type SourceFactory interface {
	NewSource(chainConfig ChainConfig, feedConfig FeedConfig) (Source, error)
}

type Envelope struct {
	// latest transmission details
	ConfigDigest    types.ConfigDigest
	Epoch           uint32
	Round           uint8
	LatestAnswer    *big.Int
	LatestTimestamp time.Time

	// latest contract config
	ContractConfig types.ContractConfig

	// extra
	BlockNumber uint64
	Transmitter types.Account
}
