package headreader

import (
	"sync"

	htrktypes "github.com/smartcontractkit/chainlink/v2/common/headtracker/types"
	"github.com/smartcontractkit/chainlink/v2/common/types"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

type HeadSaver[
	HTH htrktypes.Head[BLOCK_HASH, ID],
	ID types.ID,
	BLOCK_HASH types.Hashable,
] interface {
	SelectLatestBlock() (HTH, error)
	SelectBlockByNumber(blockNumber int64) (HTH, error)
	Store(hash HTH, logs []Log) error
	Delete(blockNumberAfterLCA int64) error
}

type headSaver[
	HTH htrktypes.Head[BLOCK_HASH, ID],
	ID types.ID,
	BLOCK_HASH types.Hashable,
] struct {
	config      htrktypes.Config
	logger      logger.Logger
	latestHead  HTH
	Heads       map[BLOCK_HASH]HTH
	HeadsNumber map[int64][]HTH
	mu          sync.RWMutex
	getNilHead  func() HTH
	getNilHash  func() BLOCK_HASH
	setParent   func(HTH, HTH)
}

func NewHeadSaver[
	HTH htrktypes.Head[BLOCK_HASH, ID],
	ID types.ID,
	BLOCK_HASH types.Hashable,
](
	config htrktypes.Config,
	lggr logger.Logger,
	getNilHead func() HTH,
	getNilHash func() BLOCK_HASH,
	setParent func(HTH, HTH),
) *headSaver[HTH, ID, BLOCK_HASH] {
	return &headSaver[HTH, ID, BLOCK_HASH]{
		config:      config,
		logger:      lggr.Named("HeadSaver"),
		Heads:       make(map[BLOCK_HASH]HTH),
		HeadsNumber: make(map[int64][]HTH),
		getNilHead:  getNilHead,
		getNilHash:  getNilHash,
		setParent:   setParent,
	}
}

func (hs *headSaver[HTH, ID, BLOCK_HASH]) Store (hash HTH, logs []Log) error {
	hs.mu.Lock()
	defer hs.mu.Unlock()

	return nil
}

func (hs *headSaver[HTH, ID, BLOCK_HASH]) AddHeads(historyDepth int64, newHeads ...HTH) {
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

		// Set the parent of the existing heads to the new heads added
		for _, existingHead := range hs.Heads {
			parentHash := existingHead.GetParentHash()
			if parentHash != hs.getNilHash() {
				if parent, exists := hs.Heads[parentHash]; exists {
					hs.setParent(existingHead, parent)
				}
			}
		}

		if !hs.latestHead.IsValid() {
			hs.latestHead = head
		} else if head.BlockNumber() > hs.latestHead.BlockNumber() {
			hs.latestHead = head
		}
	}
}

func (hs *headSaver[H,ID, BLOCK_HASH]) trimHeads(historyDepth int64) {
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