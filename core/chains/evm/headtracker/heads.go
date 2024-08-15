package headtracker

import (
	"container/heap"
	"fmt"
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
	AddHeads(newHeads ...*evmtypes.Head) error
	// Count returns number of heads in the collection.
	Count() int
	// MarkFinalized - finds `finalized` in the LatestHead and marks it and all direct ancestors as finalized.
	// Trims old blocks whose height is smaller than minBlockToKeep
	MarkFinalized(finalized common.Hash, minBlockToKeep int64) bool
}

type heads struct {
	highest       *evmtypes.Head
	headsAsc      *headsHeap
	headsByHash   map[common.Hash]*evmtypes.Head
	headsByParent map[common.Hash]map[common.Hash]*evmtypes.Head
	mu            sync.RWMutex
}

func NewHeads() Heads {
	return &heads{
		headsAsc:      &headsHeap{},
		headsByHash:   make(map[common.Hash]*evmtypes.Head),
		headsByParent: map[common.Hash]map[common.Hash]*evmtypes.Head{},
	}
}

func (h *heads) LatestHead() *evmtypes.Head {
	h.mu.RLock()
	defer h.mu.RUnlock()

	return h.highest
}

func (h *heads) HeadByHash(hash common.Hash) *evmtypes.Head {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if h.headsByHash == nil {
		return nil
	}

	return h.headsByHash[hash]
}

func (h *heads) Count() int {
	h.mu.RLock()
	defer h.mu.RUnlock()

	return h.headsAsc.Len()
}

// MarkFinalized - marks block with hash equal to finalized and all it's direct ancestors as finalized.
// Trims old blocks whose height is smaller than minBlockToKeep
func (h *heads) MarkFinalized(finalized common.Hash, minBlockToKeep int64) bool {
	h.mu.Lock()
	defer h.mu.Unlock()

	if len(h.headsByHash) == 0 {
		return false
	}

	finalizedHead, ok := h.headsByHash[finalized]
	if !ok {
		return false
	}

	markFinalized(finalizedHead)

	// remove all blocks that are older than minBlockToKeep
	for h.headsAsc.Len() > 0 && h.headsAsc.Peek().Number < minBlockToKeep {
		oldBlock := heap.Pop(h.headsAsc).(*evmtypes.Head)
		delete(h.headsByHash, oldBlock.Hash)
		// clear .Parent in oldBlock's children
		for _, oldBlockChildren := range h.headsByParent[oldBlock.Hash] {
			oldBlockChildren.Parent.Store(nil)
		}
		// headsByParent are expected to be of the same height, so we can remove them all at once
		delete(h.headsByParent, oldBlock.ParentHash)
	}

	if h.highest.Number < minBlockToKeep {
		h.highest = nil
	}

	return true
}

func markFinalized(head *evmtypes.Head) {
	// we can assume that if a head was previously marked as finalized all its ancestors were marked as finalized
	for head != nil && !head.IsFinalized.Load() {
		head.IsFinalized.Store(true)
		head = head.Parent.Load()
	}
}

func (h *heads) ensureNoCycles(newHead *evmtypes.Head) error {
	if newHead.ParentHash == newHead.Hash {
		return fmt.Errorf("cycle detected: newHeads reference itself newHead(%s)", newHead.String())
	}
	if parent, ok := h.headsByHash[newHead.ParentHash]; ok {
		if parent.Number >= newHead.Number {
			return fmt.Errorf("potential cycle detected while adding newHead as child: %w", newPotentialCycleError(parent, newHead))
		}
	}

	for _, child := range h.headsByParent[newHead.Hash] {
		if newHead.Number >= child.Number {
			return fmt.Errorf("potential cycle detected while adding newHead as parent: %w", newPotentialCycleError(newHead, child))
		}
	}

	return nil
}

func (h *heads) AddHeads(newHeads ...*evmtypes.Head) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	for _, newHead := range newHeads {
		// skip blocks that were previously added
		if _, ok := h.headsByHash[newHead.Hash]; ok {
			continue
		}

		if err := h.ensureNoCycles(newHead); err != nil {
			return err
		}

		// heads now owns the newHead - reset values that are populated by heads
		newHead.IsFinalized.Store(false)
		newHead.Parent.Store(nil)

		// prefer newer head to set as highest
		if h.highest == nil || h.highest.Number <= newHead.Number {
			h.highest = newHead
		}

		heap.Push(h.headsAsc, newHead)
		h.headsByHash[newHead.Hash] = newHead
		siblings, ok := h.headsByParent[newHead.ParentHash]
		if !ok {
			siblings = make(map[common.Hash]*evmtypes.Head)
			h.headsByParent[newHead.ParentHash] = siblings
		}
		siblings[newHead.Hash] = newHead
		// populate reference to parent
		if parent, ok := h.headsByHash[newHead.ParentHash]; ok {
			newHead.Parent.Store(parent)
		}
		for _, child := range h.headsByParent[newHead.Hash] {
			// ensure all children have reference to newHead
			child.Parent.Store(newHead)
			if child.IsFinalized.Load() {
				// mark newHead as finalized if any of its children is finalized
				markFinalized(newHead)
			}
		}
	}

	return nil
}

func newPotentialCycleError(parent, child *evmtypes.Head) error {
	return fmt.Errorf("expected head number to strictly decrease in 'child -> parent' relation: "+
		"child(%s), parent(%s)", child.String(), parent.String())
}
