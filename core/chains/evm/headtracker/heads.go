package headtracker

import (
	"sort"
	"sync"

	"github.com/ethereum/go-ethereum/common"

	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
)

// Heads is a collection of heads. All methods are thread-safe.
type Heads interface {
	// LatestHead returns the block header with the highest number that has been seen, or nil.
	LatestHead() *evmtypes.Head
	// HeadByHash returns a head for the specified hash, or nil.
	HeadByHash(hash common.Hash) *evmtypes.Head
	// AddHeads adds newHeads to the collection, eliminates duplicates,
	// sorts by head number, fixes parents and cuts off old heads (historyDepth).
	AddHeads(newHeads ...*evmtypes.Head)
	// Count returns number of heads in the collection.
	Count() int
	// MarkFinalized - finds `finalized` in the LatestHead and marks it and all direct ancestors as finalized.
	// Trims old blocks whose height is smaller than minBlockToKeep
	MarkFinalized(finalized common.Hash, minBlockToKeep int64) bool
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

// MarkFinalized - marks block with has equal to finalized and all it's direct ancestors as finalized.
// Trims old blocks whose height is smaller than minBlockToKeep
func (h *heads) MarkFinalized(finalized common.Hash, minBlockToKeep int64) bool {
	h.mu.Lock()
	defer h.mu.Unlock()

	if len(h.heads) == 0 {
		return false
	}

	// deep copy to avoid race on head.Parent
	h.heads = deepCopy(h.heads, minBlockToKeep)

	head := h.heads[0]
	foundFinalized := false
	for head != nil {
		if head.Hash == finalized {
			foundFinalized = true
		}

		// we might see finalized to move back in chain due to request to lagging RPC,
		// we should not override the flag in such cases
		head.IsFinalized = head.IsFinalized || foundFinalized
		head = head.Parent
	}

	return foundFinalized
}

func deepCopy(oldHeads []*evmtypes.Head, minBlockToKeep int64) []*evmtypes.Head {
	headsMap := make(map[common.Hash]*evmtypes.Head, len(oldHeads))
	for _, head := range oldHeads {
		if head.Hash == head.ParentHash {
			// shouldn't happen but it is untrusted input
			continue
		}
		if head.BlockNumber() < minBlockToKeep {
			// trim redundant blocks
			continue
		}
		// copy all head objects to avoid races when a previous head chain is used
		// elsewhere (since we mutate Parent here)
		headCopy := *head
		headCopy.Parent = nil // always build it from scratch in case it points to a head too old to be included
		// map eliminates duplicates
		// prefer head that was already in heads as it might have been marked as finalized on previous run
		if _, ok := headsMap[head.Hash]; !ok {
			headsMap[head.Hash] = &headCopy
		}
	}

	heads := make([]*evmtypes.Head, 0, len(headsMap))
	// unsorted unique heads
	{
		for _, head := range headsMap {
			heads = append(heads, head)
		}
	}

	// sort the heads
	sort.SliceStable(heads, func(i, j int) bool {
		// sorting from the highest number to lowest
		return heads[i].Number > heads[j].Number
	})

	// assign parents
	for i := 0; i < len(heads); i++ {
		head := heads[i]
		parent, exists := headsMap[head.ParentHash]
		if exists {
			head.Parent = parent
		}
	}

	return heads
}

func (h *heads) AddHeads(newHeads ...*evmtypes.Head) {
	h.mu.Lock()
	defer h.mu.Unlock()

	// deep copy to avoid race on head.Parent
	h.heads = deepCopy(append(h.heads, newHeads...), 0)
}
