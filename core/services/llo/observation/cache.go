package observation

import (
	"runtime"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	ocr2types "github.com/smartcontractkit/libocr/offchainreporting2/types"

	llotypes "github.com/smartcontractkit/chainlink-common/pkg/types/llo"
	"github.com/smartcontractkit/chainlink-data-streams/llo"
)

var (
	registryMu = sync.Mutex{}
	registry   = map[ocr2types.ConfigDigest]*Cache{}

	promCacheHitCount = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "llo",
		Subsystem: "datasource",
		Name:      "cache_hit_count",
		Help:      "Number of local observation cache hits",
	},
		[]string{"configDigest", "streamID"},
	)
	promCacheMissCount = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "llo",
		Subsystem: "datasource",
		Name:      "cache_miss_count",
		Help:      "Number of local observation cache misses",
	},
		[]string{"configDigest", "streamID", "reason"},
	)
)

func GetCache(configDigest ocr2types.ConfigDigest) *Cache {
	registryMu.Lock()
	defer registryMu.Unlock()

	cache, ok := registry[configDigest]
	if !ok {
		cache = NewCache(configDigest, 500*time.Millisecond, time.Minute)
		registry[configDigest] = cache
	}

	return cache
}

// Cache of stream values.
// It maintains a cache of stream values fetched from adapters until the last
// transmission sequence number is greater or equal the sequence number at which
// the value was observed or until the maxAge is reached.
//
// The cache is cleaned up periodically to remove decommissioned stream values
// if the provided cleanupInterval is greater than 0.
type Cache struct {
	mu sync.RWMutex

	configDigestStr string
	values          map[llotypes.StreamID]item
	maxAge          time.Duration
	cleanupInterval time.Duration

	lastTransmissionSeqNr atomic.Uint64
	closeChan             chan struct{}
}

type item struct {
	value     llo.StreamValue
	seqNr     uint64
	createdAt time.Time
}

// NewCache creates a new cache.
//
// maxAge is the maximum age of a stream value to keep in the cache.
// cleanupInterval is the interval to clean up the cache.
func NewCache(configDigest ocr2types.ConfigDigest, maxAge time.Duration, cleanupInterval time.Duration) *Cache {
	c := &Cache{
		configDigestStr: configDigest.Hex(),
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
	c.lastTransmissionSeqNr.Store(seqNr)
}

// Add adds a stream value to the cache.
func (c *Cache) Add(id llotypes.StreamID, value llo.StreamValue, seqNr uint64) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.values[id] = item{value: value, seqNr: seqNr, createdAt: time.Now()}
}

func (c *Cache) Get(id llotypes.StreamID) (llo.StreamValue, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	label := strconv.FormatUint(uint64(id), 10)
	item, ok := c.values[id]
	if !ok {
		promCacheMissCount.WithLabelValues(c.configDigestStr, label, "notFound").Inc()
		return nil, false
	}

	if item.seqNr <= c.lastTransmissionSeqNr.Load() {
		promCacheMissCount.WithLabelValues(c.configDigestStr, label, "seqNr").Inc()
		return nil, false
	}

	if time.Since(item.createdAt) >= c.maxAge {
		promCacheMissCount.WithLabelValues(c.configDigestStr, label, "maxAge").Inc()
		return nil, false
	}

	promCacheHitCount.WithLabelValues(c.configDigestStr, label).Inc()
	return item.value, true
}

func (c *Cache) cleanup() {
	c.mu.Lock()
	defer c.mu.Unlock()

	lastTransmissionSeqNr := c.lastTransmissionSeqNr.Load()
	for id, item := range c.values {
		if item.seqNr <= lastTransmissionSeqNr || time.Since(item.createdAt) >= c.maxAge {
			delete(c.values, id)
		}
	}
}
