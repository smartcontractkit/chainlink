package cache

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/jpillora/backoff"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	mercuryutils "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury/utils"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury/wsrpc/pb"
)

var (
	promFetchFailedCount = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "mercury_cache_fetch_failure_count",
		Help: "Number of times we tried to call LatestReport from the mercury server, but some kind of error occurred",
	},
		[]string{"serverURL", "feedID"},
	)
	promCacheHitCount = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "mercury_cache_hit_count",
		Help: "Running count of cache hits",
	},
		[]string{"serverURL", "feedID"},
	)
	promCacheWaitCount = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "mercury_cache_wait_count",
		Help: "Running count of times that we had to wait for a fetch to complete before reading from cache",
	},
		[]string{"serverURL", "feedID"},
	)
	promCacheMissCount = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "mercury_cache_miss_count",
		Help: "Running count of cache misses",
	},
		[]string{"serverURL", "feedID"},
	)
)

type Fetcher interface {
	LatestReport(ctx context.Context, req *pb.LatestReportRequest) (resp *pb.LatestReportResponse, err error)
}

type Client interface {
	Fetcher
	ServerURL() string
	RawClient() pb.MercuryClient
}

// Cache is scoped to one particular mercury server
// Use CacheSet to hold lookups for multiple servers
type Cache interface {
	Fetcher
	services.Service
}

type Config struct {
	// LatestReportTTL controls how "stale" we will allow a price to be e.g. if
	// set to 1s, a new price will always be fetched if the last result was
	// from more than 1 second ago.
	//
	// Another way of looking at it is such: the cache will _never_ return a
	// price that was queried from before now-LatestReportTTL.
	//
	// Setting to zero disables caching entirely.
	LatestReportTTL time.Duration
	// MaxStaleAge is that maximum amount of time that a value can be stale
	// before it is deleted from the cache (a form of garbage collection).
	//
	// This should generally be set to something much larger than
	// LatestReportTTL. Setting to zero disables garbage collection.
	MaxStaleAge time.Duration
	// LatestReportDeadline controls how long to wait for a response before
	// retrying. Setting this to zero will wait indefinitely.
	LatestReportDeadline time.Duration
}

func NewCache(lggr logger.Logger, client Client, cfg Config) Cache {
	return newMemCache(lggr, client, cfg)
}

type cacheVal struct {
	sync.RWMutex

	fetching bool
	fetchCh  chan (struct{})

	val *pb.LatestReportResponse
	err error

	expiresAt time.Time
}

func (v *cacheVal) read() (*pb.LatestReportResponse, error) {
	v.RLock()
	defer v.RUnlock()
	return v.val, v.err
}

// caller expected to hold lock
func (v *cacheVal) initiateFetch() <-chan struct{} {
	if v.fetching {
		panic("cannot initiateFetch on cache val that is already fetching")
	}
	v.fetching = true
	v.fetchCh = make(chan struct{})
	return v.fetchCh
}

func (v *cacheVal) setError(err error) {
	v.Lock()
	defer v.Unlock()
	v.err = err
}

func (v *cacheVal) completeFetch(val *pb.LatestReportResponse, err error, expiresAt time.Time) {
	v.Lock()
	defer v.Unlock()
	if !v.fetching {
		panic("can only completeFetch on cache val that is fetching")
	}
	v.val = val
	v.err = err
	if err == nil {
		v.expiresAt = expiresAt
	}
	close(v.fetchCh)
	v.fetchCh = nil
	v.fetching = false
}

func (v *cacheVal) abandonFetch(err error) {
	v.completeFetch(nil, err, time.Now())
}

func (v *cacheVal) waitForResult(ctx context.Context, chResult <-chan struct{}, chStop <-chan struct{}) (*pb.LatestReportResponse, error) {
	select {
	case <-ctx.Done():
		_, err := v.read()
		return nil, errors.Join(err, ctx.Err())
	case <-chStop:
		return nil, errors.New("stopped")
	case <-chResult:
		return v.read()
	}
}

// memCache stores values in memory
// it will never return a stale value older than latestPriceTTL, instead
// waiting for a successful fetch or caller context cancels, whichever comes
// first
type memCache struct {
	services.StateMachine
	lggr logger.SugaredLogger

	client Client

	cfg Config

	cache sync.Map

	wg     sync.WaitGroup
	chStop services.StopChan
}

func newMemCache(lggr logger.Logger, client Client, cfg Config) *memCache {
	return &memCache{
		services.StateMachine{},
		logger.Sugared(lggr).Named("MemCache").Named(client.ServerURL()),
		client,
		cfg,
		sync.Map{},
		sync.WaitGroup{},
		make(chan (struct{})),
	}
}

// LatestReport
// NOTE: This will actually block on all types of errors, even non-timeouts.
// Context should be set carefully and timed to be the maximum time we are
// willing to wait for a result, the background thread will keep re-querying
// until it gets one even on networking errors etc.
func (m *memCache) LatestReport(ctx context.Context, req *pb.LatestReportRequest) (resp *pb.LatestReportResponse, err error) {
	if req == nil {
		return nil, errors.New("req must not be nil")
	}
	feedIDHex := mercuryutils.BytesToFeedID(req.FeedId).String()
	if m.cfg.LatestReportTTL <= 0 {
		return m.client.RawClient().LatestReport(ctx, req)
	}
	vi, loaded := m.cache.LoadOrStore(feedIDHex, &cacheVal{
		sync.RWMutex{},
		false,
		nil,
		nil,
		nil,
		time.Now(), // first result is always "expired" and requires fetch
	})
	v := vi.(*cacheVal)

	m.lggr.Tracew("LatestReport", "feedID", feedIDHex, "loaded", loaded)

	// HOT PATH
	v.RLock()
	if time.Now().Before(v.expiresAt) {
		// CACHE HIT
		promCacheHitCount.WithLabelValues(m.client.ServerURL(), feedIDHex).Inc()
		m.lggr.Tracew("LatestReport CACHE HIT (hot path)", "feedID", feedIDHex)

		defer v.RUnlock()
		return v.val, nil
	} else if v.fetching {
		// CACHE WAIT
		promCacheWaitCount.WithLabelValues(m.client.ServerURL(), feedIDHex).Inc()
		m.lggr.Tracew("LatestReport CACHE WAIT (hot path)", "feedID", feedIDHex)
		// if someone else is fetching then wait for the fetch to complete
		ch := v.fetchCh
		v.RUnlock()
		return v.waitForResult(ctx, ch, m.chStop)
	}
	// CACHE MISS
	promCacheMissCount.WithLabelValues(m.client.ServerURL(), feedIDHex).Inc()
	// fallthrough to cold path and fetch
	v.RUnlock()

	// COLD PATH
	v.Lock()
	if time.Now().Before(v.expiresAt) {
		// CACHE HIT
		promCacheHitCount.WithLabelValues(m.client.ServerURL(), feedIDHex).Inc()
		m.lggr.Tracew("LatestReport CACHE HIT (cold path)", "feedID", feedIDHex)
		defer v.Unlock()
		return v.val, nil
	} else if v.fetching {
		// CACHE WAIT
		promCacheWaitCount.WithLabelValues(m.client.ServerURL(), feedIDHex).Inc()
		m.lggr.Tracew("LatestReport CACHE WAIT (cold path)", "feedID", feedIDHex)
		// if someone else is fetching then wait for the fetch to complete
		ch := v.fetchCh
		v.Unlock()
		return v.waitForResult(ctx, ch, m.chStop)
	}
	// CACHE MISS
	promCacheMissCount.WithLabelValues(m.client.ServerURL(), feedIDHex).Inc()
	m.lggr.Tracew("LatestReport CACHE MISS (cold path)", "feedID", feedIDHex)
	// initiate the fetch and wait for result
	ch := v.initiateFetch()
	v.Unlock()

	ok := m.IfStarted(func() {
		m.wg.Add(1)
		go m.fetch(req, v)
	})
	if !ok {
		err := fmt.Errorf("memCache must be started, but is: %v", m.State())
		v.abandonFetch(err)
		return nil, err
	}
	return v.waitForResult(ctx, ch, m.chStop)
}

const minBackoffRetryInterval = 50 * time.Millisecond

// newBackoff creates a backoff for retrying
func (m *memCache) newBackoff() backoff.Backoff {
	min := minBackoffRetryInterval
	max := m.cfg.LatestReportTTL / 2
	if min > max {
		// avoid setting a min that is greater than max
		min = max
	}
	return backoff.Backoff{
		Min:    min,
		Max:    max,
		Factor: 2,
		Jitter: true,
	}
}

// fetch continually tries to call FetchLatestReport and write the result to v
// it writes errors as they come up
func (m *memCache) fetch(req *pb.LatestReportRequest, v *cacheVal) {
	defer m.wg.Done()
	b := m.newBackoff()
	memcacheCtx, cancel := m.chStop.NewCtx()
	defer cancel()
	var t time.Time
	var val *pb.LatestReportResponse
	var err error
	defer func() {
		v.completeFetch(val, err, t.Add(m.cfg.LatestReportTTL))
	}()

	for {
		t = time.Now()

		ctx := memcacheCtx
		cancel := func() {}
		if m.cfg.LatestReportDeadline > 0 {
			ctx, cancel = context.WithTimeoutCause(memcacheCtx, m.cfg.LatestReportDeadline, errors.New("latest report fetch deadline exceeded"))
		}

		// NOTE: must drop down to RawClient here otherwise we enter an
		// infinite loop of calling a client that calls back to this same cache
		// and on and on
		val, err = m.client.RawClient().LatestReport(ctx, req)
		cancel()
		v.setError(err)
		if memcacheCtx.Err() != nil {
			// stopped
			return
		} else if err != nil {
			m.lggr.Warnw("FetchLatestReport failed", "err", err)
			promFetchFailedCount.WithLabelValues(m.client.ServerURL(), mercuryutils.BytesToFeedID(req.FeedId).String()).Inc()
			select {
			case <-m.chStop:
				return
			case <-time.After(b.Duration()):
				continue
			}
		}
		return
	}
}

func (m *memCache) Start(context.Context) error {
	return m.StartOnce(m.Name(), func() error {
		m.lggr.Debugw("MemCache starting", "config", m.cfg, "serverURL", m.client.ServerURL())
		m.wg.Add(1)
		go m.runloop()
		return nil
	})
}

func (m *memCache) runloop() {
	defer m.wg.Done()

	if m.cfg.MaxStaleAge == 0 {
		return
	}
	t := services.NewTicker(m.cfg.MaxStaleAge)
	defer t.Stop()

	for {
		select {
		case <-t.C:
			m.cleanup()
			t.Reset()
		case <-m.chStop:
			return
		}
	}
}

// remove anything that has been stale for longer than maxStaleAge so that
// cache doesn't grow forever and cause memory leaks
//
// NOTE: This should be concurrent-safe with LatestReport. The only time they
// can race is if the cache item has expired past the stale age between
// creation of the cache item and start of fetch. This is unlikely, and even if
// it does occur, the worst case is that we discard a cache item early and
// double fetch, which isn't bad at all.
func (m *memCache) cleanup() {
	m.cache.Range(func(k, vi any) bool {
		v := vi.(*cacheVal)
		v.RLock()
		defer v.RUnlock()
		if v.fetching {
			// skip cleanup if fetching
			return true
		}
		if time.Now().After(v.expiresAt.Add(m.cfg.MaxStaleAge)) {
			// garbage collection
			m.cache.Delete(k)
		}
		return true
	})
}

func (m *memCache) Close() error {
	return m.StopOnce(m.Name(), func() error {
		close(m.chStop)
		m.wg.Wait()
		return nil
	})
}
func (m *memCache) HealthReport() map[string]error {
	return map[string]error{
		m.Name(): m.Ready(),
	}
}
func (m *memCache) Name() string { return m.lggr.Name() }
