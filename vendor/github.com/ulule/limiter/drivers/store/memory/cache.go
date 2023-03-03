package memory

import (
	"runtime"
	"sync"
	"time"
)

// Forked from https://github.com/patrickmn/go-cache

// CacheWrapper is used to ensure that the underlying cleaner goroutine used to clean expired keys will not prevent
// Cache from being garbage collected.
type CacheWrapper struct {
	*Cache
}

// A cleaner will periodically delete expired keys from cache.
type cleaner struct {
	interval time.Duration
	stop     chan bool
}

// Run will periodically delete expired keys from given cache until GC notify that it should stop.
func (cleaner *cleaner) Run(cache *Cache) {
	ticker := time.NewTicker(cleaner.interval)
	for {
		select {
		case <-ticker.C:
			cache.Clean()
		case <-cleaner.stop:
			ticker.Stop()
			return
		}
	}
}

// stopCleaner is a callback from GC used to stop cleaner goroutine.
func stopCleaner(wrapper *CacheWrapper) {
	wrapper.cleaner.stop <- true
}

// startCleaner will start a cleaner goroutine for given cache.
func startCleaner(cache *Cache, interval time.Duration) {
	cleaner := &cleaner{
		interval: interval,
		stop:     make(chan bool),
	}

	cache.cleaner = cleaner
	go cleaner.Run(cache)
}

// Counter is a simple counter with an optional expiration.
type Counter struct {
	Value      int64
	Expiration int64
}

// Expired returns true if the counter has expired.
func (counter Counter) Expired() bool {
	if counter.Expiration == 0 {
		return false
	}
	return time.Now().UnixNano() > counter.Expiration
}

// Cache contains a collection of counters.
type Cache struct {
	mutex    sync.RWMutex
	counters map[string]Counter
	cleaner  *cleaner
}

// NewCache returns a new cache.
func NewCache(cleanInterval time.Duration) *CacheWrapper {

	cache := &Cache{
		counters: map[string]Counter{},
	}

	wrapper := &CacheWrapper{Cache: cache}

	if cleanInterval > 0 {
		startCleaner(cache, cleanInterval)
		runtime.SetFinalizer(wrapper, stopCleaner)
	}

	return wrapper
}

// Increment increments given value on key.
// If key is undefined or expired, it will create it.
func (cache *Cache) Increment(key string, value int64, duration time.Duration) (int64, time.Time) {
	cache.mutex.Lock()

	counter, ok := cache.counters[key]
	if !ok || counter.Expired() {
		expiration := time.Now().Add(duration).UnixNano()
		counter = Counter{
			Value:      value,
			Expiration: expiration,
		}

		cache.counters[key] = counter
		cache.mutex.Unlock()

		return value, time.Unix(0, expiration)
	}

	value = counter.Value + value
	counter.Value = value
	expiration := counter.Expiration

	cache.counters[key] = counter
	cache.mutex.Unlock()

	return value, time.Unix(0, expiration)
}

// Get returns key's value and expiration.
func (cache *Cache) Get(key string, duration time.Duration) (int64, time.Time) {
	cache.mutex.RLock()

	counter, ok := cache.counters[key]
	if !ok || counter.Expired() {
		expiration := time.Now().Add(duration).UnixNano()
		cache.mutex.RUnlock()
		return 0, time.Unix(0, expiration)
	}

	value := counter.Value
	expiration := counter.Expiration
	cache.mutex.RUnlock()

	return value, time.Unix(0, expiration)
}

// Clean will deleted any expired keys.
func (cache *Cache) Clean() {
	now := time.Now().UnixNano()

	cache.mutex.Lock()
	for key, counter := range cache.counters {
		if now > counter.Expiration {
			delete(cache.counters, key)
		}
	}
	cache.mutex.Unlock()
}
