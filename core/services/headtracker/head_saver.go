package headtracker

import (
	"context"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink/core/services/postgres"
	strpkg "github.com/smartcontractkit/chainlink/core/store"
	"github.com/smartcontractkit/chainlink/core/store/models"
)

type HeadSaver struct {
	highestSeenHead *models.Head
	store           *strpkg.Store
	headMutex       sync.RWMutex
}

func NewHeadSaver(store *strpkg.Store) *HeadSaver {
	return &HeadSaver{
		store: store,
	}
}

// Save updates the latest block number, if indeed the latest, and persists
// this number in case of reboot. Thread safe.
func (ht *HeadSaver) Save(ctx context.Context, h models.Head) error {
	ht.headMutex.Lock()
	if h.GreaterThan(ht.highestSeenHead) {
		ht.highestSeenHead = &h
	}
	ht.headMutex.Unlock()

	err := ht.store.IdempotentInsertHead(ctx, h)
	if ctx.Err() != nil {
		return nil
	} else if err != nil {
		return err
	}
	return ht.store.TrimOldHeads(ctx, ht.store.Config.EthHeadTrackerHistoryDepth())
}

// HighestSeenHead returns the block header with the highest number that has been seen, or nil
func (ht *HeadSaver) HighestSeenHead() *models.Head {
	ht.headMutex.RLock()
	defer ht.headMutex.RUnlock()

	if ht.highestSeenHead == nil {
		return nil
	}
	h := *ht.highestSeenHead
	return &h
}

func (ht *HeadSaver) IdempotentInsertHead(ctx context.Context, head models.Head) error {
	return ht.store.IdempotentInsertHead(ctx, head)
}

func (ht *HeadSaver) SetHighestSeenHeadFromDB() (*models.Head, error) {
	ht.headMutex.RLock()
	defer ht.headMutex.RUnlock()

	head, err := ht.HighestSeenHeadFromDB()
	if err != nil {
		return nil, err
	}
	ht.highestSeenHead = head
	return head, nil
}

func (ht *HeadSaver) HighestSeenHeadFromDB() (*models.Head, error) {
	ctxQuery, _ := postgres.DefaultQueryCtx()
	return ht.store.LastHead(ctxQuery)
}

func (ht *HeadSaver) Chain(ctx context.Context, hash common.Hash, depth uint) (models.Head, error) {
	return ht.store.Chain(ctx, hash, depth)
}
