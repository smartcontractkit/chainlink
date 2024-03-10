package vrfcommon

import (
	"sync"

	"github.com/ethereum/go-ethereum/core/types"
)

type InflightCache interface {
	Add(lg types.Log)
	Contains(lg types.Log) bool
	Size() int
}

var _ InflightCache = (*inflightCache)(nil)

const cachePruneInterval = 1000

type inflightCache struct {
	// cache stores the logs whose fulfillments are currently in flight or already fulfilled.
	cache map[logKey]struct{}

	// lookback defines how long state should be kept for. Logs included in blocks older than
	// lookback may or may not be redelivered.
	lookback int

	// lastPruneHeight is the blockheight at which logs were last pruned.
	lastPruneHeight uint64

	// mu synchronizes access to the delivered map.
	mu sync.RWMutex
}

func NewInflightCache(lookback int) InflightCache {
	return &inflightCache{
		cache:    make(map[logKey]struct{}),
		lookback: lookback,
		mu:       sync.RWMutex{},
	}
}

func (c *inflightCache) Size() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.cache)
}

func (c *inflightCache) Add(lg types.Log) {
	c.mu.Lock()
	defer c.mu.Unlock() // unlock in the last defer, so that we hold the lock when pruning.
	defer c.prune(lg.BlockNumber)

	c.cache[logKey{
		blockHash:   lg.BlockHash,
		blockNumber: lg.BlockNumber,
		logIndex:    lg.Index,
	}] = struct{}{}
}

func (c *inflightCache) Contains(lg types.Log) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	_, ok := c.cache[logKey{
		blockHash:   lg.BlockHash,
		blockNumber: lg.BlockNumber,
		logIndex:    lg.Index,
	}]
	return ok
}

func (c *inflightCache) prune(logBlock uint64) {
	// Only prune every pruneInterval blocks
	if int(logBlock)-int(c.lastPruneHeight) < cachePruneInterval {
		return
	}

	for key := range c.cache {
		if int(key.blockNumber) < int(logBlock)-c.lookback {
			delete(c.cache, key)
		}
	}

	c.lastPruneHeight = logBlock
}
