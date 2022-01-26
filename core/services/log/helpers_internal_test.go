package log

import (
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/eth"
)

// NewTestBroadcaster creates a broadcaster with Pause/Resume enabled.
func NewTestBroadcaster(orm ORM, ethClient eth.Client, config Config, lggr logger.Logger, highestSavedHead *eth.Head) *broadcaster {
	b := NewBroadcaster(orm, ethClient, config, lggr, highestSavedHead)
	b.testPause, b.testResume = make(chan struct{}), make(chan struct{})
	return b
}

func (b *broadcaster) ExportedAppendLogChannel(ch1, ch2 <-chan types.Log) chan types.Log {
	return b.appendLogChannel(ch1, ch2)
}
