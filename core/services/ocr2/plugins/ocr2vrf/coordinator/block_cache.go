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
	oldestItem          *cacheItem[T]
	newestItem          *cacheItem[T]
	cleaner             *intervalCacheCleaner[T]
}

type cacheItem[T any] struct {
	item         T
	itemKey      common.Hash
	blockStored  int64
	nextItem     *cacheItem[T]
	previousItem *cacheItem[T]
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
		oldestItem:          nil,
		newestItem:          nil,
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
		item:         item,
		itemKey:      itemKey,
		blockStored:  blockStored,
		nextItem:     nil,
		previousItem: nil,
	}

	// Lock, and defer unlock.
	l.cacheMu.Lock()
	defer l.cacheMu.Unlock()

	// Items cannot be added that are older than the most recent item.
	if (l.newestItem != nil) && (blockStored < l.newestItem.blockStored) {
		return errors.New("cannot add item that is older than the most recent in the cache")
	}

	// If an item "I" associated with this key has already been cached,
	// remove it from the current linked list, and increment the
	// last item / decrement the most recent item if they point to "I".
	if item, ok := l.cache[itemKey]; ok {
		if item.previousItem != nil {
			item.previousItem.nextItem = item.nextItem
		}

		if item.nextItem != nil {
			item.nextItem.previousItem = item.previousItem
		}

		if l.oldestItem != nil && item.itemKey == l.oldestItem.itemKey {
			l.oldestItem = l.oldestItem.nextItem
		}

		if l.newestItem != nil && item.itemKey == l.newestItem.itemKey {
			l.newestItem = l.newestItem.previousItem
		}
	}

	// Set last item and most recent item to the new item if nil.
	// Otherwise, add new most recent item.
	if l.oldestItem == nil {
		l.oldestItem = newItem
		l.newestItem = newItem
	} else {
		newItem.previousItem = l.newestItem
		l.newestItem.nextItem = newItem
		l.newestItem = l.newestItem.nextItem
	}

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

	// Iteratively delete items starting at the last item,
	// until all items have been deleted or a non-expired item is found.
	for l.oldestItem != nil && (newestBlock-l.oldestItem.blockStored > l.blockEvictionWindow) {
		oldestItem := l.oldestItem
		l.oldestItem = oldestItem.nextItem
		delete(l.cache, oldestItem.itemKey)
	}

	// If we have evicted the entire list, update the newest item to be nil.
	if l.oldestItem == nil {
		l.newestItem = nil
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
