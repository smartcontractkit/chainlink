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

// StreamValueCache is used by dataSource to decouple the read/write paths for stream values.
type StreamValueCache interface {
	Get(id llotypes.StreamID) (llo.StreamValue, time.Time)
	UpdateStreamValues(streamValues llo.StreamValues)
	AddMany(values map[llotypes.StreamID]llo.StreamValue, ttl time.Duration)
	Close() error
}

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
	metricsCh       chan []metricEvent

	wg        sync.WaitGroup
	closeOnce sync.Once
	closeChan chan struct{}
}

type item struct {
	value     llo.StreamValue
	expiresAt time.Time
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

// NewCache creates a new cache.
//
// maxAge is the maximum age of a stream value to keep in the cache.
// cleanupInterval is the interval to clean up the cache.
func NewCache(cleanupInterval time.Duration) *Cache {
	c := &Cache{
		values:          make(map[llotypes.StreamID]item),
		cleanupInterval: cleanupInterval,
		metricsCh:       make(chan []metricEvent, 64),
		closeChan:       make(chan struct{}),
	}

	c.wg.Add(1)
	go c.updateMetrics()

	if cleanupInterval > 0 {
		c.wg.Add(1)
		go func() {
			defer c.wg.Done()
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

// UpdateStreamValues mutates streamValues in-place for zero-allocation reads.
// Emits cache hit/miss metrics async.
func (c *Cache) UpdateStreamValues(streamValues llo.StreamValues) {
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

	c.sendMetrics(events)
}

func (c *Cache) Get(id llotypes.StreamID) (llo.StreamValue, time.Time) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	item, ok := c.values[id]
	if !ok {
		return nil, time.Time{}
	}

	if time.Now().After(item.expiresAt) {
		return nil, time.Time{}
	}

	return item.value, item.expiresAt
}

// sendMetrics enqueues metric events for async processing, dropping if the
// channel is full to avoid blocking the caller.
func (c *Cache) sendMetrics(events []metricEvent) {
	select {
	case c.metricsCh <- events:
	default:
	}
}

func (c *Cache) updateMetrics() {
	defer c.wg.Done()
	for {
		select {
		case events := <-c.metricsCh:
			for _, e := range events {
				idStr := strconv.FormatUint(uint64(e.id), 10)
				if e.cacheOutcome == cacheOutcomeHit {
					promCacheHitCount.WithLabelValues(idStr).Inc()
				} else {
					promCacheMissCount.WithLabelValues(idStr, string(e.cacheOutcome)).Inc()
				}
			}
		case <-c.closeChan:
			return
		}
	}
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
	c.closeOnce.Do(func() {
		close(c.closeChan)
	})
	c.wg.Wait()
	c.mu.Lock()
	defer c.mu.Unlock()
	c.values = nil
	return nil
}
