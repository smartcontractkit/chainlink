package log

import (
	"github.com/ethereum/go-ethereum/core/types"
	evmclient "github.com/smartcontractkit/chainlink/core/chains/evm/client"
	evmtypes "github.com/smartcontractkit/chainlink/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/core/logger"
)

// NewTestBroadcaster creates a broadcaster with Pause/Resume enabled.
func NewTestBroadcaster(orm ORM, ethClient evmclient.Client, config Config, lggr logger.Logger, highestSavedHead *evmtypes.Head) *broadcaster {
	b := NewBroadcaster(orm, ethClient, config, lggr, highestSavedHead)
	b.testPause, b.testResume = make(chan struct{}), make(chan struct{})
	return b
}

func (b *broadcaster) ExportedAppendLogChannel(ch1, ch2 <-chan types.Log) chan types.Log {
	return b.appendLogChannel(ch1, ch2)
}
