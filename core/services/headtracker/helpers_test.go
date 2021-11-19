package headtracker

import (
	"context"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink/core/services/eth"
)

func AddHeads(hs *HeadSaver, heads []*eth.Head, historyDepth int) {
	hs.addHeads(heads, historyDepth)
}

func LoadFromDB(ht *HeadTracker) (*eth.Head, error) {
	return ht.headSaver.LoadFromDB(context.Background())
}

func Heads(hs *HeadSaver) []*eth.Head {
	hs.mu.RLock()
	defer hs.mu.RUnlock()
	return hs.heads
}

func (ht *HeadTracker) Chain(hash common.Hash) *eth.Head {
	return ht.headSaver.Chain(hash)
}
