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
//
// That being said, adding caching layer to the new contract is as simple as:
// * implementing ContractOrigin interface
// * registering proper events in log poller
type CachedChain[T any] struct {
	// Static configuration
	observedEvents          []common.Hash
	logPoller               logpoller.LogPoller
	address                 []common.Address
	optimisticConfirmations int64

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

	lastChangeBlock := c.readLastChangeBlock()

	// Handles first call, because cache is not eagerly populated
	if lastChangeBlock == 0 {
		return c.initializeCache(ctx)
	}

	currentBlockNumber, err := c.logPoller.LatestBlockByEventSigsAddrsWithConfs(lastChangeBlock, c.observedEvents, c.address, logpoller.Confirmations(c.optimisticConfirmations), pg.WithParentCtx(ctx))

	if err != nil {
		return empty, err
	}

	// In case of new updates, fetch fresh data from the origin
	if currentBlockNumber > lastChangeBlock {
		return c.fetchFromOrigin(ctx, currentBlockNumber)
	}
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

	c.updateCache(value, latestBlock-c.optimisticConfirmations)
	return c.copyCachedValue(), nil

}

// fetchFromOrigin fetches data from origin. This action is performed when logpoller.LogPoller says there were events
// emitted since last update.
func (c *CachedChain[T]) fetchFromOrigin(ctx context.Context, currentBlockNumber int64) (T, error) {
	var empty T
	value, err := c.origin.CallOrigin(ctx)
	if err != nil {
		return empty, err
	}
	c.updateCache(value, currentBlockNumber)

	return c.copyCachedValue(), nil
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
