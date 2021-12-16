package headtracker

import (
	"github.com/smartcontractkit/chainlink/core/services/eth"
)

// to be used for testing only
func AddHeads(saver HeadSaver, heads []*eth.Head, historyDepth uint) {
	saver.(*headSaver).addHeads(heads, historyDepth)
}

// to be used for testing only
func Heads(saver HeadSaver) []*eth.Head {
	hs := saver.(*headSaver)
	hs.mu.RLock()
	defer hs.mu.RUnlock()
	return hs.heads
}
