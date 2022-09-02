package coordinator

import (
	"errors"
	"sync"

	"github.com/ethereum/go-ethereum/common"
)

type BlockCache[T any] struct {
	blockEvictionWindow int64
	cacheMu             sync.RWMutex
	cache               map[common.Hash]*cacheItem[T]
	lastItem            *cacheItem[T]
}

type cacheItem[T any] struct {
	item        T
	itemKey     common.Hash
	blockStored int64
	nextItem    *cacheItem[T]
}

// NewBlockCache constructs a new cache.
func NewBlockCache[T any](blockEvictionWindow int64) *BlockCache[T] {
	return &BlockCache[T]{
		cacheMu:             sync.RWMutex{},
		cache:               make(map[common.Hash]*cacheItem[T]),
		blockEvictionWindow: blockEvictionWindow,
		lastItem:            nil,
	}
}

// AddItem adds an item to the cache.
func (l *BlockCache[T]) AddItem(item T, itemKey common.Hash, blockStored int64) {

	// Construct new item to be stored.
	newItem := &cacheItem[T]{
		item:        item,
		itemKey:     itemKey,
		blockStored: blockStored,
		nextItem:    nil,
	}

	// Lock, and defer unlock.
	l.cacheMu.Lock()
	defer l.cacheMu.Unlock()

	// Set last item to the new item if nil, otherwise add new last item.
	if l.lastItem == nil {
		l.lastItem = newItem
	} else {
		newItem.nextItem = l.lastItem
		l.lastItem = newItem
	}

	// Assign item to key.
	l.cache[itemKey] = newItem
}

// AddItem adds an item to the cache.
func (l *BlockCache[T]) GetItem(itemKey common.Hash) (item *T, err error) {

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
func (l *BlockCache[T]) EvictExpiredItems(newestBlock int64) {

	// Lock, and defer unlock.
	l.cacheMu.Lock()
	defer l.cacheMu.Unlock()

	// Iteratively delete items starting at the last item,
	// until all items have been deleted or a non-expired item is found.
	for l.lastItem != nil && (newestBlock-l.lastItem.blockStored > l.blockEvictionWindow) {
		lastItem := l.lastItem
		l.lastItem = lastItem.nextItem
		delete(l.cache, lastItem.itemKey)
	}
}
