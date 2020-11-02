package eth

import (
	"github.com/ethereum/go-ethereum/core/types"
)

func (lb *logBroadcaster) ExportedAppendLogChannel(ch1, ch2 <-chan types.Log) chan types.Log {
	return lb.appendLogChannel(ch1, ch2)
}
