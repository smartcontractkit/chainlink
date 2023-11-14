package ccipdata

import (
	"context"
	"time"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
)

type Event[T any] struct {
	Data T
	Meta
}

type Meta struct {
	BlockTimestamp time.Time
	BlockNumber    int64
	TxHash         common.Hash
	LogIndex       uint
}

const (
	V1_0_0 = "1.0.0"
	V1_1_0 = "1.1.0"
	V1_2_0 = "1.2.0"
)

type Closer interface {
	Close(qopts ...pg.QOpt) error
}

// Client can be used to fetch CCIP related parsed on-chain data.
//
//go:generate mockery --quiet --name Reader --filename reader_mock.go --case=underscore
type Reader interface {
	// LatestBlock returns the latest known/parsed block of the underlying implementation.
	LatestBlock(ctx context.Context) (logpoller.LogPollerBlock, error)
}
