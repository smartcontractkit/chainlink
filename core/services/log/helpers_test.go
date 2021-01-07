package log

import (
	"github.com/ethereum/go-ethereum/core/types"
)

func (lb *broadcaster) ExportedAppendLogChannel(ch1, ch2 <-chan types.Log) chan types.Log {
	return lb.appendLogChannel(ch1, ch2)
}
