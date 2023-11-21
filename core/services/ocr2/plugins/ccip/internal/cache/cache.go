package cache

import (
	"context"
	"sync"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
)

// AutoSync cache provides only a Get method, the expiration and syncing is a black-box for the caller.
//
//go:generate mockery --quiet --name AutoSync --output . --filename cache_mock.go --inpackage --case=underscore
type AutoSync[T any] interface {
	Get(ctx context.Context) (T, error)
}

// CachedChain represents caching on-chain calls based on the events read from logpoller.LogPoller.
// Instead of directly going to on-chain to fetch data, we start with checking logpoller.LogPoller events (database request).
// If we discover that change occurred since last update, we perform RPC to the chain using ContractOrigin.CallOrigin function.
// Purpose of this struct is handle common logic in a single place, you only need to override methods from ContractOrigin
// and Get function (behaving as orchestrator) will take care of the rest.
// IMPORTANT: Cache refresh relies on the events that are finalized. This introduces some delay between the event onchain occurrence
// and cache refreshing. This is intentional, because we want to prevent handling reorgs within the cache.
//
// That being said, adding caching layer to the new contract is as simple as:
// * implementing ContractOrigin interface
// * registering proper events in log poller
type CachedChain[T any] struct {
	// Static configuration
	observedEvents []common.Hash
	logPoller      logpoller.LogPoller
	address        []common.Address

	// Cache
	lock            *sync.RWMutex
	value           T
	lastChangeBlock int64
	origin          ContractOrigin[T]
}

type ContractOrigin[T any] interface {
	// Copy must return copy of the cached data to limit locking section to the minimum
	Copy(T) T
	// CallOrigin fetches data that is next stored within cache. Usually, should perform RPC to the source (e.g. chain)
	CallOrigin(ctx context.Context) (T, error)
}

// Get is an entry point to the caching. Main function that decides whether cache content is fresh and should be returned
// to the caller, or whether we need to update it's content from on-chain data.
// This decision is made based on the events emitted by Smart Contracts
func (c *CachedChain[T]) Get(ctx context.Context) (T, error) {
	var empty T

	cachedLastChangeBlock := c.readLastChangeBlock()
	// Handles first call, because cache is not eagerly populated
	if cachedLastChangeBlock == 0 {
		return c.initializeCache(ctx)
	}

	// Ordering matters here, we need to do operations in the following order:
	// * get LatestBlock
	// * get LatestBlockByEventSigsAddrsWithConfs
	// * fetch data from Origin
	// It's because LogPoller keep progressing in the background, and we want to prevent missing data.
	// If we do it in the opposite order, we might store in cache block that after logs that
	// were not scanned by LatestBlockByEventSigsAddrsWithConfs. And therefore ignore them and not update the cache.
	// (this will ignore logs produced between LatestBlockByEventSigsAddrsWithConfs and LatestBlock calls).
	// Calling LatestBlock first gives us guarantee that we never miss anything.
	latestBlock, err := c.logPoller.LatestBlock(pg.WithParentCtx(ctx))
	if err != nil {
		// Intentionally ignore the error here
		latestBlock = logpoller.LogPollerBlock{}
	}

	if err1 := c.maybeRefreshCache(ctx, cachedLastChangeBlock); err1 != nil {
		return empty, err1
	}

	// This is performance improvement that will prevent for large db scans, by updating the lower bound of the search query
	c.maybeCacheLatestFinalizedBlock(cachedLastChangeBlock, latestBlock.FinalizedBlockNumber)
	return c.copyCachedValue(), nil
}

// initializeCache performs first call to origin when is not populated yet.
// It's done eagerly, so cache it's populated for the first time when data is needed, not at struct initialization
func (c *CachedChain[T]) initializeCache(ctx context.Context) (T, error) {
	var empty T

	// To prevent missing data when blocks are produced after calling the origin,
	// we first get the latest block and then call the origin.
	latestBlock, err := c.logPoller.LatestBlock(pg.WithParentCtx(ctx))
	if err != nil {
		return empty, err
	}

	// Init
	value, err := c.origin.CallOrigin(ctx)
	if err != nil {
		return empty, err
	}

	c.updateCache(value, latestBlock.FinalizedBlockNumber)
	return c.copyCachedValue(), nil

}

// maybeRefreshCache checks whether cache is fresh or needs to be updated.
// We fetch the last changed block from the log poller and compare that with the last change block stored within cache.
// If the last changed block is greater than the one stored within cache, we need to update the cache by fetching data from the origin.
func (c *CachedChain[T]) maybeRefreshCache(ctx context.Context, cachedLastChangeBlock int64) error {
	chainLastChangeBlock, err := c.logPoller.LatestBlockByEventSigsAddrsWithConfs(
		cachedLastChangeBlock,
		c.observedEvents,
		c.address,
		logpoller.Finalized,
		pg.WithParentCtx(ctx),
	)
	if err != nil {
		return err
	}

	// In case of new updates, fetch fresh data from the origin
	if chainLastChangeBlock > cachedLastChangeBlock {
		// Return error when cache cannot be fetched, don't return stale values
		value, err1 := c.origin.CallOrigin(ctx)
		if err1 != nil {
			return err1
		}
		c.updateCache(value, chainLastChangeBlock)
	}
	return nil
}

// updateCache performs updating two critical variables for cache to work properly:
// * value that is stored within cache
// * lastChangeBlock representing last seen event from logpoller.LogPoller
func (c *CachedChain[T]) updateCache(newValue T, currentBlockNumber int64) {
	c.lock.Lock()
	defer c.lock.Unlock()

	// Double-lock checking. No need to update if other goroutine was faster
	if currentBlockNumber <= c.lastChangeBlock {
		return
	}

	c.value = newValue
	c.lastChangeBlock = currentBlockNumber
}

func (c *CachedChain[T]) readLastChangeBlock() int64 {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.lastChangeBlock
}

func (c *CachedChain[T]) copyCachedValue() T {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.origin.Copy(c.value)
}

func (c *CachedChain[T]) maybeCacheLatestFinalizedBlock(cachedLastBlock int64, latestFinalizedBlock int64) {
	// Check if applicable to prevent unnecessary locking
	if cachedLastBlock >= latestFinalizedBlock {
		return
	}

	c.lock.Lock()
	defer c.lock.Unlock()
	// Double-lock checking. No need to update if other goroutine was faster
	if latestFinalizedBlock <= c.lastChangeBlock {
		return
	}
	c.lastChangeBlock = latestFinalizedBlock
}
