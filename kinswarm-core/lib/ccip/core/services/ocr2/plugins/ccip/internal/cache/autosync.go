package cache

import (
	"context"
	"database/sql"
	"fmt"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
)

type AutoSync[T any] interface {
	Get(ctx context.Context, syncFunc func(ctx context.Context) (T, error)) (T, error)
}

// LogpollerEventsBased IMPORTANT: Cache refresh relies on the events that are finalized.
// This introduces some delay between the event onchain occurrence and cache refreshing.
// This is intentional, because we want to prevent handling reorgs within the cache.
type LogpollerEventsBased[T any] struct {
	logPoller      logpoller.LogPoller
	observedEvents []common.Hash
	address        common.Address

	lock            *sync.RWMutex
	value           T
	lastChangeBlock int64
}

func NewLogpollerEventsBased[T any](
	lp logpoller.LogPoller,
	observedEvents []common.Hash,
	contractAddress common.Address,
) *LogpollerEventsBased[T] {
	var emptyValue T
	return &LogpollerEventsBased[T]{
		logPoller:      lp,
		observedEvents: observedEvents,
		address:        contractAddress,

		lock:            &sync.RWMutex{},
		value:           emptyValue,
		lastChangeBlock: 0,
	}
}

func (c *LogpollerEventsBased[T]) Get(ctx context.Context, syncFunc func(ctx context.Context) (T, error)) (T, error) {
	var empty T

	hasExpired, newEventBlockNum, err := c.hasExpired(ctx)
	if err != nil {
		return empty, fmt.Errorf("check cache expiration: %w", err)
	}

	if hasExpired {
		var latestValue T
		latestValue, err = syncFunc(ctx)
		if err != nil {
			return empty, fmt.Errorf("sync func: %w", err)
		}

		c.set(latestValue, newEventBlockNum)
		return latestValue, nil
	}

	cachedValue := c.get()
	if err != nil {
		return empty, fmt.Errorf("get cached value: %w", err)
	}

	c.lock.Lock()
	if newEventBlockNum > c.lastChangeBlock {
		// update the most recent block number
		// that way the scanning window is shorter in the next run
		c.lastChangeBlock = newEventBlockNum
	}
	c.lock.Unlock()

	return cachedValue, nil
}

func (c *LogpollerEventsBased[T]) hasExpired(ctx context.Context) (expired bool, blockOfLatestEvent int64, err error) {
	c.lock.RLock()
	blockOfCurrentValue := c.lastChangeBlock
	c.lock.RUnlock()

	// NOTE: latest block should be fetched before LatestBlockByEventSigsAddrsWithConfs
	// Otherwise there might be new events between LatestBlockByEventSigsAddrsWithConfs and
	// latestBlock which will be missed.
	latestBlock, err := c.logPoller.LatestBlock(ctx)
	latestFinalizedBlock := int64(0)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return false, 0, fmt.Errorf("get latest log poller block: %w", err)
	} else if err == nil {
		// Since we know that we have all the events till latestBlock.FinalizedBlockNumber
		// we want to return the block number instead of the block of the latest event
		// for reducing the scan window on the next call.
		latestFinalizedBlock = latestBlock.FinalizedBlockNumber
	}

	if blockOfCurrentValue == 0 {
		return true, latestFinalizedBlock, nil
	}

	blockOfLatestEvent, err = c.logPoller.LatestBlockByEventSigsAddrsWithConfs(
		ctx,
		blockOfCurrentValue,
		c.observedEvents,
		[]common.Address{c.address},
		evmtypes.Finalized,
	)
	if err != nil {
		return false, 0, fmt.Errorf("get latest events form lp: %w", err)
	}

	if blockOfLatestEvent > latestFinalizedBlock {
		latestFinalizedBlock = blockOfLatestEvent
	}
	return blockOfLatestEvent > blockOfCurrentValue, latestFinalizedBlock, nil
}

func (c *LogpollerEventsBased[T]) set(value T, blockNum int64) {
	c.lock.Lock()
	defer c.lock.Unlock()

	if c.lastChangeBlock > blockNum {
		return
	}

	c.value = value
	c.lastChangeBlock = blockNum
}

func (c *LogpollerEventsBased[T]) get() T {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.value
}
