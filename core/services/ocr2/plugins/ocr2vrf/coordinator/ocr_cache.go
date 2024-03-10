package coordinator

import (
	"runtime"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
)

// ocrCache is a caching strucuture that allows items to be stored and then evicted
// based on an eviction window. In this package, it is being used to track in-flight
// items the coordinator includes in OCR reports, such that the items can be checked against
// the cache to avoid double-transmissions.
type ocrCache[T any] struct {
	evictionWindow time.Duration
	cacheMu        sync.Mutex
	cache          map[common.Hash]*cacheItem[T]
	cleaner        *intervalCacheCleaner[T]
}

type cacheItem[T any] struct {
	item       T
	itemKey    common.Hash
	timeStored time.Time
}

// NewBlockCache constructs a new cache.
func NewBlockCache[T any](evictionWindow time.Duration) *ocrCache[T] {

	// Construct cache cleaner to evict old items.
	cleaner := &intervalCacheCleaner[T]{
		interval: evictionWindow,
		stop:     make(chan struct{}, 1),
	}

	// Instantiate the cache for type T.
	cache := &ocrCache[T]{
		cacheMu:        sync.Mutex{},
		cache:          make(map[common.Hash]*cacheItem[T]),
		evictionWindow: evictionWindow,
		cleaner:        cleaner,
	}

	// Stop the cleaner upon garbage collection of the cache.
	runtime.SetFinalizer(cache, func(b *ocrCache[T]) { b.cleaner.stop <- struct{}{} })

	return cache
}

// AddItem adds an item to the cache.
func (l *ocrCache[T]) CacheItem(item T, itemKey common.Hash, timeStored time.Time) {

	// Construct new item to be stored.
	newItem := &cacheItem[T]{
		item:       item,
		itemKey:    itemKey,
		timeStored: timeStored,
	}

	// Lock, and defer unlock.
	l.cacheMu.Lock()
	defer l.cacheMu.Unlock()

	// Assign item to key.
	l.cache[itemKey] = newItem
}

func (l *ocrCache[T]) SetEvictonWindow(newWindow time.Duration) {
	l.evictionWindow = newWindow
}

// AddItem adds an item to the cache.
func (l *ocrCache[T]) GetItem(itemKey common.Hash) (item *T) {

	// Lock, and defer unlock.
	l.cacheMu.Lock()
	defer l.cacheMu.Unlock()

	// Construct new item to be stored.
	cacheItem := l.cache[itemKey]

	// Return nil if the item is not found, otherwise return item.
	if cacheItem == nil {
		return
	}

	return &cacheItem.item
}

// EvictExpiredItems removes all expired items stored in the cache.
func (l *ocrCache[T]) EvictExpiredItems(currentTime time.Time) {

	// Lock, and defer unlock.
	l.cacheMu.Lock()
	defer l.cacheMu.Unlock()

	// Iteratively check all item ages, and delete an item if it is expired.
	for key, item := range l.cache {
		diff := currentTime.Sub(item.timeStored)
		if diff > l.evictionWindow {
			delete(l.cache, key)
		}
	}
}

// A cache cleaner that evicts items on a regular interval.
type intervalCacheCleaner[T any] struct {
	interval time.Duration
	stop     chan struct{}
}

// Run evicts expired items every n seconds, until the "stop" channel is triggered.
func (ic *intervalCacheCleaner[T]) Run(c *ocrCache[T]) {
	ticker := time.NewTicker(ic.interval)
	for {
		select {
		case <-ticker.C:
			c.EvictExpiredItems(time.Now().UTC())
		case <-ic.stop:
			ticker.Stop()
			return
		}
	}
}
