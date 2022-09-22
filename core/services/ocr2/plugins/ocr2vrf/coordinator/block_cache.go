package coordinator

import (
	"errors"
	"runtime"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
)

type blockCache[T any] struct {
	blockEvictionWindow int64
	cacheMu             sync.RWMutex
	cache               map[common.Hash]*cacheItem[T]
	cleaner             *intervalCacheCleaner[T]
}

type cacheItem[T any] struct {
	item        T
	itemKey     common.Hash
	blockStored int64
}

// NewBlockCache constructs a new cache.
func NewBlockCache[T any](blockEvictionWindow int64) *blockCache[T] {

	// Construct cache cleaner to evict old items.
	cleaner := &intervalCacheCleaner[T]{
		interval: blockEvictionWindow,
		stop:     make(chan struct{}, 1),
	}

	// Instantiate the cache for type T.
	cache := &blockCache[T]{
		cacheMu:             sync.RWMutex{},
		cache:               make(map[common.Hash]*cacheItem[T]),
		blockEvictionWindow: blockEvictionWindow,
		cleaner:             cleaner,
	}

	// Stop the cleaner upon garbage collection of the cache.
	runtime.SetFinalizer(cache, func(b *blockCache[T]) { b.cleaner.stop <- struct{}{} })

	return cache
}

// AddItem adds an item to the cache.
func (l *blockCache[T]) CacheItem(item T, itemKey common.Hash, blockStored int64) error {

	// Construct new item to be stored.
	newItem := &cacheItem[T]{
		item:        item,
		itemKey:     itemKey,
		blockStored: blockStored,
	}

	// Lock, and defer unlock.
	l.cacheMu.Lock()
	defer l.cacheMu.Unlock()

	// Assign item to key.
	l.cache[itemKey] = newItem

	return nil
}

// AddItem adds an item to the cache.
func (l *blockCache[T]) GetItem(itemKey common.Hash) (item *T, err error) {

	// Lock, and defer unlock.
	l.cacheMu.Lock()
	defer l.cacheMu.Unlock()

	// Construct new item to be stored.
	cacheItem := l.cache[itemKey]

	// Return error if item is not found, otherwise return item.
	if cacheItem == nil {
		err = errors.New("requested item not found")
	} else {
		item = &cacheItem.item
	}

	return
}

// EvictExpiredItems removes all expired items stored in the cache.
func (l *blockCache[T]) EvictExpiredItems(newestBlock int64) {

	// Lock, and defer unlock.
	l.cacheMu.Lock()
	defer l.cacheMu.Unlock()

	// Iteratively check all item ages, and delete an item if it is expired.
	for key, item := range l.cache {
		if newestBlock-item.blockStored > l.blockEvictionWindow {
			delete(l.cache, key)
		}
	}
}

// A cache cleaner that evicts items on a regular interval.
type intervalCacheCleaner[T any] struct {
	interval int64
	stop     chan struct{}
}

// Run evicts expired items every n seconds, until the "stop" channel is triggered.
func (ic *intervalCacheCleaner[T]) Run(c *blockCache[T]) {
	ticker := time.NewTicker(time.Duration(ic.interval * int64(time.Second)))
	for {
		select {
		case <-ticker.C:
			c.EvictExpiredItems(time.Now().Unix())
		case <-ic.stop:
			ticker.Stop()
			return
		}
	}
}
