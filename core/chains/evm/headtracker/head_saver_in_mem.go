package headtracker

import (
	"context"
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	txmgrtypes "github.com/smartcontractkit/chainlink/v2/common/txmgr/types"
	"github.com/smartcontractkit/chainlink/v2/common/types"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

type inMemoryHeadSaver[H types.HeadTrackerHead[BLOCK_HASH, CHAIN_ID], BLOCK_HASH types.Hashable, CHAIN_ID txmgrtypes.ID] struct {
	config      Config
	logger      logger.Logger
	latestHead  H
	Heads       map[BLOCK_HASH]H
	HeadsNumber map[int64][]H
	mu          sync.RWMutex
	getNilHead  func() H
	setParent   func(H, H)
}

type EvmInMemoryHeadSaver = inMemoryHeadSaver[*evmtypes.Head, common.Hash, *big.Int]

var _ types.InMemoryHeadSaver[*evmtypes.Head, common.Hash, *big.Int] = (*EvmInMemoryHeadSaver)(nil)

func NewInMemoryHeadSaver[
	H types.HeadTrackerHead[BLOCK_HASH, CHAIN_ID],
	BLOCK_HASH types.Hashable,
	CHAIN_ID types.ID](
	config Config,
	lggr logger.Logger,
	getNilHead func() H,
	setParent func(H, H),
) *inMemoryHeadSaver[H, BLOCK_HASH, CHAIN_ID] {
	return &inMemoryHeadSaver[H, BLOCK_HASH, CHAIN_ID]{
		config:      config,
		logger:      lggr.Named("InMemoryHeadSaver"),
		Heads:       make(map[BLOCK_HASH]H),
		HeadsNumber: make(map[int64][]H),
		getNilHead:  getNilHead,
		setParent:   setParent,
	}
}

func NewEvmInMemoryHeadSaver() *EvmInMemoryHeadSaver {
	return &EvmInMemoryHeadSaver{
		Heads:       make(map[common.Hash]*evmtypes.Head),
		HeadsNumber: make(map[int64][]*evmtypes.Head),
		getNilHead:  func() *evmtypes.Head { return nil },
		setParent: func(h *evmtypes.Head, parent *evmtypes.Head) {
			h.Parent = parent
		},
	}
}

func (hs *inMemoryHeadSaver[H, BLOCK_HASH, CHAIN_ID]) Save(ctx context.Context, head H) error {
	hs.mu.Lock()
	defer hs.mu.Unlock()

	blockHash := head.BlockHash()
	blockNumber := head.BlockNumber()

	hs.Heads[blockHash] = head
	hs.HeadsNumber[blockNumber] = append(hs.HeadsNumber[blockNumber], head)

	if head.BlockNumber() > hs.latestHead.BlockNumber() {
		hs.latestHead = head
	}
	return nil
}

// No OP function for EVM
func (hs *inMemoryHeadSaver[H, BLOCK_HASH]) Load(ctx context.Context) (H, error) {

	// Pseudo Code
	// 1.Gets Heads from client
	// 2.Calls AddHeads to link the heads together and populate the Map struct

	return hs.latestHead, nil
}

func (hs *inMemoryHeadSaver[H, BLOCK_HASH]) LatestChain() H {
	hs.mu.RLock()
	defer hs.mu.RUnlock()

	return hs.latestHead
}

func (hs *inMemoryHeadSaver[H, BLOCK_HASH]) Chain(blockHash BLOCK_HASH) H {
	hs.mu.RLock()
	defer hs.mu.RUnlock()

	if head, exists := hs.Heads[blockHash]; exists {
		return head
	}

	return hs.getNilHead()
}

func (hs *inMemoryHeadSaver[H, BLOCK_HASH]) HeadByNumber(blockNumber int64) []H {
	hs.mu.RLock()
	defer hs.mu.RUnlock()

	return hs.HeadsNumber[blockNumber]
}

// Assembles the heads together and populates the Heads Map
func (hs *inMemoryHeadSaver[H, BLOCK_HASH]) AddHeads(historyDepth int64, newHeads ...H) {
	hs.mu.Lock()
	defer hs.mu.Unlock()

	// Trim heads to avoid including head that is too old
	// Triming occurs to remove outdated data before adding
	hs.trimHeads(historyDepth)

	for _, head := range newHeads {
		blockHash := head.BlockHash()
		blockNumber := head.BlockNumber()

		// Check if the head already exists
		if _, exists := hs.Heads[blockHash]; exists {
			continue
		}

		if parent, exists := hs.Heads[blockHash]; exists {
			hs.setParent(head, parent)
		} else {
			// Ignore Parent's head if is too old
			hs.setParent(head, hs.getNilHead())
		}

		hs.Heads[blockHash] = head
		hs.HeadsNumber[blockNumber] = append(hs.HeadsNumber[blockNumber], head)

		// Update the latest head if Block number is higher
		if head.BlockNumber() > hs.latestHead.BlockNumber() {
			hs.latestHead = head
		}
	}
}

// TrimOldHeads() removes old heads such that only N new heads remain
// This function can be called externally to remove old heads.
func (hs *inMemoryHeadSaver[H, BLOCK_HASH]) TrimOldHeads(number int64) {
	hs.mu.Lock()
	defer hs.mu.Unlock()

	hs.trimHeads(number)
}

// trimHeads is an internal function without locking to prevent deadlocks
// This function has no mutex locking as it is supposed to be called by functions which already have mutex locking in place.
func (hs *inMemoryHeadSaver[H, BLOCK_HASH]) trimHeads(number int64) {

	// Create a list to store block numbers to remove
	var blockNumbersToRemove []int64

	// Iterate through the map and identify the block numbers to remove
	for headNumber, headList := range hs.HeadsNumber {
		if headNumber < number { // TODO: Check if this is correct
			blockNumbersToRemove = append(blockNumbersToRemove, headNumber)
		}

		// Remove each of the head using blockHash in the Heads map
		for _, head := range headList {
			delete(hs.Heads, head.BlockHash())
		}
	}

	// Remove the corresponding heads from the Heads map according to the block hashes
	for _, headList := range hs.HeadsNumber {
		for _, head := range headList {
			delete(hs.Heads, head.BlockHash())
		}
	}
}
