package headtracker

import (
	"context"
	"sort"
	"sync"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/eth"
)

type HeadSaver struct {
	orm    *ORM
	config Config
	heads  []*eth.Head
	logger logger.Logger
	mu     sync.RWMutex
}

func NewHeadSaver(lggr logger.Logger, orm *ORM, config Config) *HeadSaver {
	return &HeadSaver{
		orm:    orm,
		config: config,
		logger: lggr.Named("HeadSaver"),
	}
}

// Save updates the latest block number, if indeed the latest, and persists
// this number in case of reboot. Thread safe.
func (ht *HeadSaver) Save(ctx context.Context, h eth.Head) error {
	err := ht.orm.IdempotentInsertHead(ctx, h)
	if err != nil {
		return err
	}

	historyDepth := ht.config.EvmHeadTrackerHistoryDepth()
	ht.mu.Lock()
	ht.addHead(&h, int(historyDepth))
	ht.mu.Unlock()

	return ht.orm.TrimOldHeads(ctx, uint(historyDepth))
}

func (ht *HeadSaver) LoadFromDB(ctx context.Context) (chain *eth.Head, err error) {
	historyDepth := int(ht.config.EvmHeadTrackerHistoryDepth())
	heads, err := ht.orm.LatestHeads(ctx, historyDepth)
	if err != nil {
		return nil, err
	}
	ht.mu.Lock()
	defer ht.mu.Unlock()

	ht.addHeads(heads, historyDepth)
	return ht.latestChain(), nil
}

// LatestChain returns the block header with the highest number that has been seen, or nil
func (ht *HeadSaver) LatestChain() *eth.Head {
	ht.mu.RLock()
	defer ht.mu.RUnlock()
	ch := ht.latestChain()
	if ch.ChainLength() < ht.config.EvmFinalityDepth() {
		ht.logger.Debugw("chain shorter than EvmFinalityDepth", "chainLen", ch.ChainLength(), "evmFinalityDepth", ht.config.EvmFinalityDepth())
	}
	return ch
}

func (ht *HeadSaver) latestChain() *eth.Head {
	if len(ht.heads) == 0 {
		return nil
	}
	return ht.heads[0]
}

func (ht *HeadSaver) Chain(hash common.Hash) *eth.Head {
	ht.mu.RLock()
	defer ht.mu.RUnlock()

	h := ht.headByHash(hash)
	if h == nil {
		return nil
	}
	return h
}

// note: not thread-safe
func (ht *HeadSaver) headByHash(hash common.Hash) (h *eth.Head) {
	for _, h := range ht.heads {
		if h.Hash == hash {
			return h
		}
	}
	return nil
}

// note: not thread-safe
func (ht *HeadSaver) addHead(h *eth.Head, historyDepth int) {
	ht.addHeads([]*eth.Head{h}, historyDepth)
}

// note: not thread-safe
func (ht *HeadSaver) addHeads(hs []*eth.Head, historyDepth int) {
	hMap := make(map[common.Hash]*eth.Head, len(ht.heads)+len(hs))
	for _, head := range append(ht.heads, hs...) {
		if head.Hash == head.ParentHash {
			// shouldn't happen but its untrusted input
			ht.logger.Errorf("ignoring head %s that points to itself", head)
			continue
		}
		// copy all head objects to avoid races when a previous head chain is used
		// elsewhere (since we mutate Parent here)
		cphead := *head
		cphead.Parent = nil // always build it from scratch in case it points to a head too old to be included
		// map eliminates duplicates
		hMap[head.Hash] = &cphead
	}
	heads := make([]*eth.Head, len(hMap))
	// unsorted unique heads
	{
		var i int
		for _, head := range hMap {
			heads[i] = head
			i++
		}
	}
	// sort the heads
	sort.SliceStable(heads, func(i, j int) bool {
		// sorting from highest number to lowest
		return heads[i].Number > heads[j].Number
	})
	// cut off the oldest
	if len(heads) > historyDepth {
		heads = heads[:historyDepth]
	}
	// assign parents
	for i := 0; i < len(heads)-1; i++ {
		head := heads[i]
		parent, exists := hMap[head.ParentHash]
		if exists {
			head.Parent = parent
		}
	}
	// set
	ht.heads = heads
}
