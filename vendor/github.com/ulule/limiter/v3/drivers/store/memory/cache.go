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
	wrapper.cleaner = nil
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

// Counter is a simple counter with an expiration.
type Counter struct {
	mutex      sync.RWMutex
	value      int64
	expiration int64
}

// Value returns the counter current value.
func (counter *Counter) Value() int64 {
	counter.mutex.RLock()
	defer counter.mutex.RUnlock()
	return counter.value
}

// Expiration returns the counter expiration.
func (counter *Counter) Expiration() int64 {
	counter.mutex.RLock()
	defer counter.mutex.RUnlock()
	return counter.expiration
}

// Expired returns true if the counter has expired.
func (counter *Counter) Expired() bool {
	counter.mutex.RLock()
	defer counter.mutex.RUnlock()

	return counter.expiration == 0 || time.Now().UnixNano() > counter.expiration
}

// Load returns the value and the expiration of this counter.
// If the counter is expired, it will use the given expiration.
func (counter *Counter) Load(expiration int64) (int64, int64) {
	counter.mutex.RLock()
	defer counter.mutex.RUnlock()

	if counter.expiration == 0 || time.Now().UnixNano() > counter.expiration {
		return 0, expiration
	}

	return counter.value, counter.expiration
}

// Increment increments given value on this counter.
// If the counter is expired, it will use the given expiration.
// It returns its current value and expiration.
func (counter *Counter) Increment(value int64, expiration int64) (int64, int64) {
	counter.mutex.Lock()
	defer counter.mutex.Unlock()

	if counter.expiration == 0 || time.Now().UnixNano() > counter.expiration {
		counter.value = value
		counter.expiration = expiration
		return counter.value, counter.expiration
	}

	counter.value += value
	return counter.value, counter.expiration
}

// Cache contains a collection of counters.
type Cache struct {
	counters sync.Map
	cleaner  *cleaner
}

// NewCache returns a new cache.
func NewCache(cleanInterval time.Duration) *CacheWrapper {

	cache := &Cache{}
	wrapper := &CacheWrapper{Cache: cache}

	if cleanInterval > 0 {
		startCleaner(cache, cleanInterval)
		runtime.SetFinalizer(wrapper, stopCleaner)
	}

	return wrapper
}

// LoadOrStore returns the existing counter for the key if present.
// Otherwise, it stores and returns the given counter.
// The loaded result is true if the counter was loaded, false if stored.
func (cache *Cache) LoadOrStore(key string, counter *Counter) (*Counter, bool) {
	val, loaded := cache.counters.LoadOrStore(key, counter)
	if val == nil {
		return counter, false
	}

	actual := val.(*Counter)
	return actual, loaded
}

// Load returns the counter stored in the map for a key, or nil if no counter is present.
// The ok result indicates whether counter was found in the map.
func (cache *Cache) Load(key string) (*Counter, bool) {
	val, ok := cache.counters.Load(key)
	if val == nil || !ok {
		return nil, false
	}
	actual := val.(*Counter)
	return actual, true
}

// Store sets the counter for a key.
func (cache *Cache) Store(key string, counter *Counter) {
	cache.counters.Store(key, counter)
}

// Delete deletes the value for a key.
func (cache *Cache) Delete(key string) {
	cache.counters.Delete(key)
}

// Range calls handler sequentially for each key and value present in the cache.
// If handler returns false, range stops the iteration.
func (cache *Cache) Range(handler func(key string, counter *Counter)) {
	cache.counters.Range(func(k interface{}, v interface{}) bool {
		if v == nil {
			return true
		}

		key := k.(string)
		counter := v.(*Counter)

		handler(key, counter)

		return true
	})
}

// Increment increments given value on key.
// If key is undefined or expired, it will create it.
func (cache *Cache) Increment(key string, value int64, duration time.Duration) (int64, time.Time) {
	expiration := time.Now().Add(duration).UnixNano()

	// If counter is in cache, try to load it first.
	counter, loaded := cache.Load(key)
	if loaded {
		value, expiration = counter.Increment(value, expiration)
		return value, time.Unix(0, expiration)
	}

	// If it's not in cache, try to atomically create it.
	// We do that in two step to reduce memory allocation.
	counter, loaded = cache.LoadOrStore(key, &Counter{
		mutex:      sync.RWMutex{},
		value:      value,
		expiration: expiration,
	})
	if loaded {
		value, expiration = counter.Increment(value, expiration)
		return value, time.Unix(0, expiration)
	}

	// Otherwise, it has been created, return given value.
	return value, time.Unix(0, expiration)
}

// Get returns key's value and expiration.
func (cache *Cache) Get(key string, duration time.Duration) (int64, time.Time) {
	expiration := time.Now().Add(duration).UnixNano()

	counter, ok := cache.Load(key)
	if !ok {
		return 0, time.Unix(0, expiration)
	}

	value, expiration := counter.Load(expiration)
	return value, time.Unix(0, expiration)
}

// Clean will deleted any expired keys.
func (cache *Cache) Clean() {
	cache.Range(func(key string, counter *Counter) {
		if counter.Expired() {
			cache.Delete(key)
		}
	})
}

// Reset changes the key's value and resets the expiration.
func (cache *Cache) Reset(key string, duration time.Duration) (int64, time.Time) {
	cache.Delete(key)

	expiration := time.Now().Add(duration).UnixNano()
	return 0, time.Unix(0, expiration)
}
