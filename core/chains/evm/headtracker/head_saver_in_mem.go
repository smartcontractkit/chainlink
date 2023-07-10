package headtracker

import (
	"context"
	"errors"
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
	getNilHash  func() BLOCK_HASH
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
	getNilHash func() BLOCK_HASH,
	setParent func(H, H),
) *inMemoryHeadSaver[H, BLOCK_HASH, CHAIN_ID] {
	return &inMemoryHeadSaver[H, BLOCK_HASH, CHAIN_ID]{
		config:      config,
		logger:      lggr.Named("InMemoryHeadSaver"),
		Heads:       make(map[BLOCK_HASH]H),
		HeadsNumber: make(map[int64][]H),
		getNilHead:  getNilHead,
		getNilHash:  getNilHash,
		setParent:   setParent,
	}
}

func NewEvmInMemoryHeadSaver(config Config, lggr logger.Logger) *EvmInMemoryHeadSaver {
	evmConfig := NewWrappedConfig(config)
	return NewInMemoryHeadSaver[*evmtypes.Head, common.Hash, *big.Int](
		evmConfig,
		lggr,
		func() *evmtypes.Head { return nil },
		func() common.Hash { return common.Hash{} },
		func(head, parent *evmtypes.Head) { head.Parent = parent },
	)
}

func (hs *inMemoryHeadSaver[H, BLOCK_HASH, CHAIN_ID]) Save(ctx context.Context, head H) error {
	if !head.IsValid() {
		return errors.New("invalid head passed to Save method of InMemoryHeadSaver")
	}

	historyDepth := int64(hs.config.HeadTrackerHistoryDepth())
	hs.AddHeads(historyDepth, head)

	return nil
}

// No OP function for EVM
func (hs *inMemoryHeadSaver[H, BLOCK_HASH, CHAIN_ID]) Load(ctx context.Context) (H, error) {

	return hs.LatestChain(), nil
}

func (hs *inMemoryHeadSaver[H, BLOCK_HASH, CHAIN_ID]) LatestChain() H {
	head := hs.getLatestHead()

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

	hs.trimHeads(historyDepth)

	for _, head := range newHeads {
		blockHash := head.BlockHash()
		blockNumber := head.BlockNumber()
		parentHash := head.GetParentHash()

		if _, exists := hs.Heads[blockHash]; exists {
			continue
		}

		if parentHash != hs.getNilHash() {
			if parent, exists := hs.Heads[parentHash]; exists {
				hs.setParent(head, parent)
			} else {
				// If parent's head is too old, we should set it to nil
				hs.setParent(head, hs.getNilHead())
			}
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

func (hs *inMemoryHeadSaver[H, BLOCK_HASH, CHAIN_ID]) TrimOldHeads(historyDepth int64) {
	hs.mu.Lock()
	defer hs.mu.Unlock()

	hs.trimHeads(historyDepth)
}

// trimHeads() is should only be called by functions with mutex locking.
// trimHeads() is an internal function without locking to prevent deadlocks
func (hs *inMemoryHeadSaver[H, BLOCK_HASH, CHAIN_ID]) trimHeads(historyDepth int64) {
	for headNumber, headNumberList := range hs.HeadsNumber {
		// Checks if the block lies within the historyDepth
		if hs.latestHead.BlockNumber()-headNumber >= historyDepth {
			for _, head := range headNumberList {
				delete(hs.Heads, head.BlockHash())
			}

			delete(hs.HeadsNumber, headNumber)
		}
	}
}

func (hs *inMemoryHeadSaver[H, BLOCK_HASH, CHAIN_ID]) getLatestHead() H {
	hs.mu.RLock()
	defer hs.mu.RUnlock()

	return hs.latestHead
}
