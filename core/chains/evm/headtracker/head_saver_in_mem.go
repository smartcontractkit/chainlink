package headtracker

import (
	"context"
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	htrktypes "github.com/smartcontractkit/chainlink/v2/common/headtracker/types"
	"github.com/smartcontractkit/chainlink/v2/common/types"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

type inMemoryHeadSaver[H types.HeadTrackerHead[BLOCK_HASH, CHAIN_ID], BLOCK_HASH types.Hashable, CHAIN_ID types.ID] struct {
	config      htrktypes.Config
	logger      logger.Logger
	latestHead  H
	Heads       map[BLOCK_HASH]H
	HeadsNumber map[int64][]H
	mu          sync.RWMutex
	getNilHead  func() H
	setParent   func(H, H)
}

type EvmInMemoryHeadSaver = inMemoryHeadSaver[*evmtypes.Head, common.Hash, *big.Int]

var _ types.HeadSaver[*evmtypes.Head, common.Hash] = (*EvmInMemoryHeadSaver)(nil)

func NewInMemoryHeadSaver[
	H types.HeadTrackerHead[BLOCK_HASH, CHAIN_ID],
	BLOCK_HASH types.Hashable,
	CHAIN_ID types.ID](
	config htrktypes.Config,
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

func NewEvmInMemoryHeadSaver(config Config, lggr logger.Logger) *EvmInMemoryHeadSaver {
	evmConfig := NewWrappedConfig(config)
	return NewInMemoryHeadSaver[*evmtypes.Head, common.Hash, *big.Int](
		evmConfig,
		lggr,
		func() *evmtypes.Head { return nil },
		func(head, parent *evmtypes.Head) { head.Parent = parent },
	)
}

func (hs *inMemoryHeadSaver[H, BLOCK_HASH, CHAIN_ID]) Save(ctx context.Context, head H) error {
	historyDepth := int64(hs.config.HeadTrackerHistoryDepth())
	hs.AddHeads(int64(historyDepth), head) // TODO: Why int64?

	return nil
}

// No OP function for EVM
func (hs *inMemoryHeadSaver[H, BLOCK_HASH, CHAIN_ID]) Load(ctx context.Context) (H, error) {

	// Pseudo Code
	// 1.Gets Heads from client
	// 2.Calls AddHeads to link the heads together and populate the Map struct

	return hs.latestHead, nil
}

func (hs *inMemoryHeadSaver[H, BLOCK_HASH, CHAIN_ID]) LatestChain() H {
	head := hs.LatestHead()

	if !head.IsValid() {
		return hs.getNilHead()
	}

	if head.ChainLength() < hs.config.FinalityDepth() {
		hs.logger.Debugw("chain shorter than EvmFinalityDepth", "chainLen", head.ChainLength(), "evmFinalityDepth", hs.config.FinalityDepth())
	}
	return head
}

func (hs *inMemoryHeadSaver[H, BLOCK_HASH, CHAIN_ID]) Chain(blockHash BLOCK_HASH) H {
	hs.mu.RLock()
	defer hs.mu.RUnlock()

	if head, exists := hs.Heads[blockHash]; exists {
		return head
	}

	return hs.getNilHead()
}

func (hs *inMemoryHeadSaver[H, BLOCK_HASH, CHAIN_ID]) HeadByNumber(blockNumber int64) []H {
	hs.mu.RLock()
	defer hs.mu.RUnlock()

	return hs.HeadsNumber[blockNumber]
}

// Assembles the heads together and populates the Heads Map
func (hs *inMemoryHeadSaver[H, BLOCK_HASH, CHAIN_ID]) AddHeads(historyDepth int64, newHeads ...H) {
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

		if !hs.latestHead.IsValid() {
			hs.latestHead = head
		} else if head.BlockNumber() > hs.latestHead.BlockNumber() {
			hs.latestHead = head
		}
	}
}

// TrimOldHeads() removes old heads such that only N new heads remain
// This function can be called externally to remove old heads.
func (hs *inMemoryHeadSaver[H, BLOCK_HASH, CHAIN_ID]) TrimOldHeads(historyDepth int64) {
	hs.mu.Lock()
	defer hs.mu.Unlock()

	hs.trimHeads(historyDepth)
}

func (hs *inMemoryHeadSaver[H, BLOCK_HASH, CHAIN_ID]) trimHeads(historyDepth int64) {
	for headNumber, headNumberList := range hs.HeadsNumber {
		if headNumber < historyDepth {
			for _, head := range headNumberList {
				delete(hs.Heads, head.BlockHash())
			}

			delete(hs.HeadsNumber, headNumber) // TODO: Check if this is safe, and good practice
		}
	}
}

func (hs *inMemoryHeadSaver[H, BLOCK_HASH, CHAIN_ID]) LatestHead() H {
	hs.mu.RLock()
	defer hs.mu.RUnlock()

	return hs.latestHead
}
