package observation

import (
	"strconv"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"

	llotypes "github.com/smartcontractkit/chainlink-common/pkg/types/llo"
	"github.com/smartcontractkit/chainlink-data-streams/llo"
)

var (
	promCacheHitCount = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "llo",
		Subsystem: "datasource",
		Name:      "cache_hit_count",
		Help:      "Number of local observation cache hits",
	},
		[]string{"streamID"},
	)
	promCacheMissCount = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "llo",
		Subsystem: "datasource",
		Name:      "cache_miss_count",
		Help:      "Number of local observation cache misses",
	},
		[]string{"streamID", "reason"},
	)
)

// Cache of stream values.
// It maintains a cache of stream values fetched from adapters until the last
// transmission sequence number is greater or equal the sequence number at which
// the value was observed or until the maxAge is reached.
//
// The cache is cleaned up periodically to remove decommissioned stream values
// if the provided cleanupInterval is greater than 0.
type Cache struct {
	mu              sync.RWMutex
	values          map[llotypes.StreamID]item
	cleanupInterval time.Duration

	closeChan chan struct{}
}

type item struct {
	value     llo.StreamValue
	expiresAt time.Time
}

// NewCache creates a new cache.
//
// maxAge is the maximum age of a stream value to keep in the cache.
// cleanupInterval is the interval to clean up the cache.
func NewCache(cleanupInterval time.Duration) *Cache {
	c := &Cache{
		values:          make(map[llotypes.StreamID]item),
		cleanupInterval: cleanupInterval,
		closeChan:       make(chan struct{}),
	}

	if cleanupInterval > 0 {
		go func() {
			ticker := time.NewTicker(cleanupInterval)
			defer ticker.Stop()
			for {
				select {
				case <-ticker.C:
					c.cleanup()
				case <-c.closeChan:
					return
				}
			}
		}()
	}

	return c
}

// Add adds a stream value to the cache.
func (c *Cache) Add(id llotypes.StreamID, value llo.StreamValue, ttl time.Duration) {
	var expiresAt time.Time
	if ttl > 0 {
		expiresAt = time.Now().Add(ttl)
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	c.values[id] = item{value: value, expiresAt: expiresAt}
}

func (c *Cache) AddMany(values map[llotypes.StreamID]llo.StreamValue, ttl time.Duration) {
	var expiresAt time.Time
	if ttl > 0 {
		expiresAt = time.Now().Add(ttl)
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	for id, value := range values {
		c.values[id] = item{value: value, expiresAt: expiresAt}
	}
}

type cacheOutcome string

const (
	cacheOutcomeNotFound cacheOutcome = "notFound"
	cacheOutcomeMaxAge   cacheOutcome = "maxAge"
	cacheOutcomeHit      cacheOutcome = "" // empty string means cache hit
)

type metricEvent struct {
	id           llotypes.StreamID
	cacheOutcome cacheOutcome
}

//nolint:revive // GetMany mutates streamValues in-place for zero-allocation reads.
func (c *Cache) GetMany(streamValues llo.StreamValues) {
	events := make([]metricEvent, 0, len(streamValues))

	c.mu.RLock()
	now := time.Now()
	for id := range streamValues {
		itm, ok := c.values[id]
		if !ok {
			events = append(events, metricEvent{id: id, cacheOutcome: cacheOutcomeNotFound})
			streamValues[id] = nil
			continue
		}
		if now.After(itm.expiresAt) {
			events = append(events, metricEvent{id: id, cacheOutcome: cacheOutcomeMaxAge})
			streamValues[id] = nil
			continue
		}
		events = append(events, metricEvent{id: id, cacheOutcome: cacheOutcomeHit})
		streamValues[id] = itm.value
	}
	c.mu.RUnlock()

	// defer metric updates until after the read lock is released
	for _, e := range events {
		if e.cacheOutcome == cacheOutcomeHit {
			promCacheHitCount.WithLabelValues(strconv.FormatUint(uint64(e.id), 10)).Inc()
		} else {
			promCacheMissCount.WithLabelValues(strconv.FormatUint(uint64(e.id), 10), string(e.cacheOutcome)).Inc()
		}
	}
}

func (c *Cache) Get(id llotypes.StreamID) (llo.StreamValue, time.Time) {
	c.mu.RLock()
	value, expiresAt, metricEvent := c.get(id)
	c.mu.RUnlock()

	// defer metric updates until after the read lock is released
	if metricEvent.cacheOutcome != cacheOutcomeHit {
		promCacheMissCount.WithLabelValues(strconv.FormatUint(uint64(id), 10), string(metricEvent.cacheOutcome)).Inc()
	} else {
		promCacheHitCount.WithLabelValues(strconv.FormatUint(uint64(id), 10)).Inc()
	}
	return value, expiresAt
}

func (c *Cache) get(id llotypes.StreamID) (llo.StreamValue, time.Time, metricEvent) {
	item, ok := c.values[id]
	if !ok {
		return nil, time.Time{}, metricEvent{id: id, cacheOutcome: cacheOutcomeNotFound}
	}

	if time.Now().After(item.expiresAt) {
		return nil, time.Time{}, metricEvent{id: id, cacheOutcome: cacheOutcomeMaxAge}
	}

	return item.value, item.expiresAt, metricEvent{id: id, cacheOutcome: cacheOutcomeHit}
}

func (c *Cache) cleanup() {
	c.mu.Lock()
	defer c.mu.Unlock()

	for id, item := range c.values {
		if item.expiresAt.IsZero() {
			continue
		}

		if time.Now().After(item.expiresAt) {
			delete(c.values, id)
		}
	}
}

func (c *Cache) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.cleanupInterval > 0 {
		close(c.closeChan)
	}
	c.values = nil
	return nil
}
