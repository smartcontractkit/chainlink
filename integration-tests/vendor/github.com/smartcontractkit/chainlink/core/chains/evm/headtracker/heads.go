package headtracker

import (
	"sort"
	"sync"

	"github.com/ethereum/go-ethereum/common"

	evmtypes "github.com/smartcontractkit/chainlink/core/chains/evm/types"
)

// Heads is a collection of heads. All methods are thread-safe.
type Heads interface {
	// LatestHead returns the block header with the highest number that has been seen, or nil.
	LatestHead() *evmtypes.Head
	// HeadByHash returns a head for the specified hash, or nil.
	HeadByHash(hash common.Hash) *evmtypes.Head
	// AddHeads adds newHeads to the collection, eliminates duplicates,
	// sorts by head number, fixes parents and cuts off old heads (historyDepth).
	AddHeads(historyDepth uint, newHeads ...*evmtypes.Head)
	// Count returns number of heads in the collection.
	Count() int
}

type heads struct {
	heads []*evmtypes.Head
	mu    sync.RWMutex
}

func NewHeads() Heads {
	return &heads{}
}

func (h *heads) LatestHead() *evmtypes.Head {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if len(h.heads) == 0 {
		return nil
	}
	return h.heads[0]
}

func (h *heads) HeadByHash(hash common.Hash) *evmtypes.Head {
	h.mu.RLock()
	defer h.mu.RUnlock()

	for _, head := range h.heads {
		if head.Hash == hash {
			return head
		}
	}
	return nil
}

func (h *heads) Count() int {
	h.mu.RLock()
	defer h.mu.RUnlock()

	return len(h.heads)
}

func (h *heads) AddHeads(historyDepth uint, newHeads ...*evmtypes.Head) {
	h.mu.Lock()
	defer h.mu.Unlock()

	headsMap := make(map[common.Hash]*evmtypes.Head, len(h.heads)+len(newHeads))
	for _, head := range append(h.heads, newHeads...) {
		if head.Hash == head.ParentHash {
			// shouldn't happen but it is untrusted input
			continue
		}
		// copy all head objects to avoid races when a previous head chain is used
		// elsewhere (since we mutate Parent here)
		headCopy := *head
		headCopy.Parent = nil // always build it from scratch in case it points to a head too old to be included
		// map eliminates duplicates
		headsMap[head.Hash] = &headCopy
	}

	heads := make([]*evmtypes.Head, len(headsMap))
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
	h.heads = heads
}
