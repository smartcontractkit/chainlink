package ccipdata

import (
	"time"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
)

const (
	V1_0_0 = "1.0.0"
	V1_1_0 = "1.1.0"
	V1_2_0 = "1.2.0"
	V1_3_0 = "1.3.0-dev"
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

type Closer interface {
	Close(qopts ...pg.QOpt) error
}

func LogsConfirmations(finalized bool) logpoller.Confirmations {
	if finalized {
		return logpoller.Finalized
	}
	return logpoller.Confirmations(0)
}
