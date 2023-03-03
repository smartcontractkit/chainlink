package util

import (
	"sync"
	"time"
)

const (
	// convenience value for setting expiration to the default value
	DefaultCacheExpiration time.Duration = 0
)

func NewIntervalCacheCleaner[T any](interval time.Duration) *IntervalCacheCleaner[T] {
	return &IntervalCacheCleaner[T]{
		interval: interval,
		stop:     make(chan struct{}),
	}
}

func NewCache[T any](expiration time.Duration) *Cache[T] {
	return &Cache[T]{
		defaultExpiration: expiration,
		data:              make(map[string]CacheItem[T]),
	}
}

type CacheItem[T any] struct {
	Item    T
	Expires int64
}

type Cache[T any] struct {
	defaultExpiration time.Duration
	mu                sync.RWMutex
	data              map[string]CacheItem[T]
}

func (c *Cache[T]) Set(key string, value T, expire time.Duration) {
	var exp int64
	if expire == DefaultCacheExpiration {
		expire = c.defaultExpiration
	}

	if expire > 0 {
		exp = time.Now().Add(expire).UnixNano()
	}

	c.mu.Lock()
	c.data[key] = CacheItem[T]{
		Item:    value,
		Expires: exp,
	}
	c.mu.Unlock()
}

func (c *Cache[T]) Get(key string) (T, bool) {
	c.mu.RLock()
	value, found := c.data[key]
	if !found {
		c.mu.RUnlock()
		return getZero[T](), false
	}

	if value.Expires > 0 {
		if time.Now().UnixNano() > value.Expires {
			c.mu.RUnlock()
			return getZero[T](), false
		}
	}

	c.mu.RUnlock()
	return value.Item, true
}

func (c *Cache[T]) Keys() []string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	keys := make([]string, 0, len(c.data))
	for key := range c.data {
		keys = append(keys, key)
	}

	return keys
}

func (c *Cache[T]) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.data, key)
}

// ClearExpired loops through all keys and evaluates the value
// expire time. If an item is expired, it is removed from the
// cache. This function places a read lock on the data set and
// only obtains a write lock if needed.
func (c *Cache[T]) ClearExpired() {
	now := time.Now().UnixNano()
	c.mu.RLock()
	toclear := make([]string, 0, len(c.data))
	for k, item := range c.data {
		if item.Expires > 0 && now > item.Expires {
			toclear = append(toclear, k)
		}
	}
	c.mu.RUnlock()

	if len(toclear) > 0 {
		c.mu.Lock()
		for _, k := range toclear {
			delete(c.data, k)
		}
		c.mu.Unlock()
	}
}

func getZero[T any]() T {
	var result T
	return result
}

type IntervalCacheCleaner[T any] struct {
	interval time.Duration
	stopper  sync.Once
	stop     chan struct{}
}

func (ic *IntervalCacheCleaner[T]) Run(c *Cache[T]) {
	ticker := time.NewTicker(ic.interval)
	for {
		select {
		case <-ticker.C:
			c.ClearExpired()
		case <-ic.stop:
			ticker.Stop()
			return
		}
	}
}

func (ic *IntervalCacheCleaner[T]) Stop() {
	ic.stopper.Do(func() {
		close(ic.stop)
	})
}
