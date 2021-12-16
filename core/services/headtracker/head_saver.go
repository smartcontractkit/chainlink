package headtracker

import (
	"context"
	"sort"
	"sync"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/eth"
)

// HeadSaver maintains chains persisted in DB. All methods are thread-safe.
type HeadSaver interface {
	// Save updates the latest block number, if indeed the latest, and persists
	// this number in case of reboot.
	Save(ctx context.Context, head *eth.Head) error
	// LoadFromDB loads latest EvmHeadTrackerHistoryDepth heads, returns the latest chain.
	LoadFromDB(ctx context.Context) (*eth.Head, error)
	// LatestHeadFromDB returns the highest seen head from DB.
	LatestHeadFromDB(ctx context.Context) (*eth.Head, error)
	// LatestChain returns the block header with the highest number that has been seen, or nil.
	LatestChain() *eth.Head
	// Chain returns a head for the specified hash, or nil.
	Chain(hash common.Hash) *eth.Head
}

type headSaver struct {
	orm    ORM
	config Config
	logger logger.Logger
	heads  []*eth.Head
	mu     sync.RWMutex
}

func NewHeadSaver(lggr logger.Logger, orm ORM, config Config) HeadSaver {
	return &headSaver{
		orm:    orm,
		config: config,
		logger: lggr.Named(logger.HeadSaver),
	}
}

func (hs *headSaver) Save(ctx context.Context, head *eth.Head) error {
	if err := hs.orm.IdempotentInsertHead(ctx, head); err != nil {
		return err
	}

	historyDepth := uint(hs.config.EvmHeadTrackerHistoryDepth())
	hs.mu.Lock()
	hs.addHead(head, historyDepth)
	hs.mu.Unlock()

	return hs.orm.TrimOldHeads(ctx, historyDepth)
}

func (hs *headSaver) LoadFromDB(ctx context.Context) (chain *eth.Head, err error) {
	historyDepth := uint(hs.config.EvmHeadTrackerHistoryDepth())
	heads, err := hs.orm.LatestHeads(ctx, historyDepth)
	if err != nil {
		return nil, err
	}

	hs.mu.Lock()
	defer hs.mu.Unlock()

	hs.addHeads(heads, historyDepth)
	return hs.latestChain(), nil
}

func (hs *headSaver) LatestHeadFromDB(ctx context.Context) (head *eth.Head, err error) {
	return hs.orm.LatestHead(ctx)
}

func (hs *headSaver) LatestChain() *eth.Head {
	hs.mu.RLock()
	defer hs.mu.RUnlock()

	ch := hs.latestChain()
	if ch == nil {
		return nil
	}
	if ch.ChainLength() < hs.config.EvmFinalityDepth() {
		hs.logger.Debugw("chain shorter than EvmFinalityDepth", "chainLen", ch.ChainLength(), "evmFinalityDepth", hs.config.EvmFinalityDepth())
	}
	return ch
}

func (hs *headSaver) Chain(hash common.Hash) *eth.Head {
	hs.mu.RLock()
	defer hs.mu.RUnlock()

	return hs.headByHash(hash)
}

func (hs *headSaver) latestChain() *eth.Head {
	if len(hs.heads) == 0 {
		return nil
	}
	return hs.heads[0]
}

func (hs *headSaver) headByHash(hash common.Hash) (h *eth.Head) {
	for _, h := range hs.heads {
		if h.Hash == hash {
			return h
		}
	}
	return nil
}

func (hs *headSaver) addHead(h *eth.Head, historyDepth uint) {
	hs.addHeads([]*eth.Head{h}, historyDepth)
}

func (hs *headSaver) addHeads(newHeads []*eth.Head, historyDepth uint) {
	headsMap := make(map[common.Hash]*eth.Head, len(hs.heads)+len(newHeads))
	for _, head := range append(hs.heads, newHeads...) {
		if head.Hash == head.ParentHash {
			// shouldn't happen but it is untrusted input
			hs.logger.Errorf("ignoring head %s that points to itself", head)
			continue
		}
		// copy all head objects to avoid races when a previous head chain is used
		// elsewhere (since we mutate Parent here)
		headCopy := *head
		headCopy.Parent = nil // always build it from scratch in case it points to a head too old to be included
		// map eliminates duplicates
		headsMap[head.Hash] = &headCopy
	}
	heads := make([]*eth.Head, len(headsMap))
	// unsorted unique heads
	{
		var i int
		for _, head := range headsMap {
			heads[i] = head
			i++
		}
	}
	// sort the heads
	sort.SliceStable(heads, func(i, j int) bool {
		// sorting from the highest number to lowest
		return heads[i].Number > heads[j].Number
	})
	// cut off the oldest
	if uint(len(heads)) > historyDepth {
		heads = heads[:historyDepth]
	}
	// assign parents
	for i := 0; i < len(heads)-1; i++ {
		head := heads[i]
		parent, exists := headsMap[head.ParentHash]
		if exists {
			head.Parent = parent
		}
	}
	// set
	hs.heads = heads
}
