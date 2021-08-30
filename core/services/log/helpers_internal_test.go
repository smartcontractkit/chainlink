package log

import (
	"github.com/ethereum/go-ethereum/core/types"
)

func (b *broadcaster) ExportedAppendLogChannel(ch1, ch2 <-chan types.Log) chan types.Log {
	return b.appendLogChannel(ch1, ch2)
}

func (b *broadcaster) ExportedGetPoolCount() func() int {
	return func() int {
		count := 0
		for _, logs := range b.logPool.logsByBlockHash {
			count += len(logs)
		}
		return count
	}
}
