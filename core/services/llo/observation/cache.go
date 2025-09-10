package observation

import (
	"runtime"
	"strconv"
	"sync"
	"sync/atomic"
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
	mu sync.RWMutex

	values          map[llotypes.StreamID]item
	maxAge          time.Duration
	cleanupInterval time.Duration

	lastTransmissionSeqNr atomic.Uint64
	closeChan             chan struct{}
}

type item struct {
	value     llo.StreamValue
	seqNr     uint64
	expiresAt time.Time
}

// NewCache creates a new cache.
//
// maxAge is the maximum age of a stream value to keep in the cache.
// cleanupInterval is the interval to clean up the cache.
func NewCache(maxAge time.Duration, cleanupInterval time.Duration) *Cache {
	c := &Cache{
		values:          make(map[llotypes.StreamID]item),
		maxAge:          maxAge,
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

		runtime.AddCleanup(c, func(ch chan struct{}) {
			close(ch)
		}, c.closeChan)
	}

	return c
}

// SetLastTransmissionSeqNr sets the last transmission sequence number.
func (c *Cache) SetLastTransmissionSeqNr(seqNr uint64) {
	if c == nil {
		return
	}

	c.lastTransmissionSeqNr.Store(seqNr)
}

// Add adds a stream value to the cache.
func (c *Cache) Add(id llotypes.StreamID, value llo.StreamValue, seqNr uint64) {
	if c == nil {
		return
	}

	c.mu.Lock()
	defer c.mu.Unlock()
	c.values[id] = item{value: value, seqNr: seqNr, expiresAt: time.Now().Add(c.maxAge)}
}

func (c *Cache) Get(id llotypes.StreamID) llo.StreamValue {
	if c == nil {
		return nil
	}

	c.mu.RLock()
	defer c.mu.RUnlock()

	label := strconv.FormatUint(uint64(id), 10)
	item, ok := c.values[id]
	if !ok {
		promCacheMissCount.WithLabelValues(label, "notFound").Inc()
		return nil
	}

	if item.seqNr <= c.lastTransmissionSeqNr.Load() {
		promCacheMissCount.WithLabelValues(label, "seqNr").Inc()
		return nil
	}

	if time.Now().After(item.expiresAt) {
		promCacheMissCount.WithLabelValues(label, "maxAge").Inc()
		return nil
	}

	promCacheHitCount.WithLabelValues(label).Inc()
	return item.value
}

func (c *Cache) cleanup() {
	c.mu.Lock()
	defer c.mu.Unlock()

	lastTransmissionSeqNr := c.lastTransmissionSeqNr.Load()
	now := time.Now()
	for id, item := range c.values {
		if item.seqNr <= lastTransmissionSeqNr || now.After(item.expiresAt) {
			delete(c.values, id)
		}
	}
}
