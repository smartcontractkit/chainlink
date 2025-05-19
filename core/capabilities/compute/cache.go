package compute

import (
	"fmt"
	"sync"
	"time"

	"github.com/jonboulle/clockwork"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"

	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/workflows/wasm/host"
)

var (
	moduleCacheHit = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "compute_module_cache_hit",
		Help: "hit vs non-hits of the module cache for custom compute",
	}, []string{"hit"})
	moduleCacheEviction = promauto.NewCounter(prometheus.CounterOpts{
		Name: "compute_module_cache_eviction",
		Help: "evictions from the module cache",
	})
	moduleCacheAddition = promauto.NewCounter(prometheus.CounterOpts{
		Name: "compute_module_cache_addition",
		Help: "additions to the module cache",
	})
)

type moduleCache struct {
	m  map[string]*module
	mu sync.RWMutex

	wg       sync.WaitGroup
	stopChan services.StopChan

	tickInterval   time.Duration
	timeout        time.Duration
	evictAfterSize int

	clock      clockwork.Clock
	reapTicker <-chan time.Time
	onReaper   chan struct{}
}

func newModuleCache(clock clockwork.Clock, tick, timeout time.Duration, evictAfterSize int) *moduleCache {
	return &moduleCache{
		m:              map[string]*module{},
		tickInterval:   tick,
		timeout:        timeout,
		evictAfterSize: evictAfterSize,
		clock:          clock,
		reapTicker:     clock.NewTicker(tick).Chan(),
		stopChan:       make(chan struct{}),
	}
}

func (mc *moduleCache) start() {
	mc.wg.Add(1)
	go func() {
		defer mc.wg.Done()
		mc.reapLoop()
	}()
}

func (mc *moduleCache) close() {
	close(mc.stopChan)
	mc.wg.Wait()
}

func (mc *moduleCache) reapLoop() {
	for {
		select {
		case <-mc.reapTicker:
			mc.evictOlderThan(mc.timeout)
			if mc.onReaper != nil {
				mc.onReaper <- struct{}{}
			}
		case <-mc.stopChan:
			return
		}
	}
}

func (mc *moduleCache) add(id string, mod *module) error {
	mc.mu.Lock()
	defer mc.mu.Unlock()

	if mc.m[id] != nil {
		return fmt.Errorf("module with id %q already exists in cache", id)
	}

	mod.lastFetchedAt = mc.clock.Now()
	mc.m[id] = mod
	moduleCacheAddition.Inc()
	return nil
}

func (mc *moduleCache) get(id string) (*module, bool) {
	mc.mu.Lock()
	defer mc.mu.Unlock()
	gotModule, ok := mc.m[id]
	if !ok {
		moduleCacheHit.WithLabelValues("false").Inc()
		return nil, false
	}

	moduleCacheHit.WithLabelValues("true").Inc()
	gotModule.lastFetchedAt = mc.clock.Now()
	return gotModule, true
}

func (mc *moduleCache) evictOlderThan(duration time.Duration) {
	mc.mu.Lock()
	defer mc.mu.Unlock()

	evicted := 0

	if len(mc.m) > mc.evictAfterSize {
		for id, m := range mc.m {
			if mc.clock.Now().Sub(m.lastFetchedAt) > duration {
				delete(mc.m, id)
				m.module.Close()
				evicted++
			}

			if len(mc.m) <= mc.evictAfterSize {
				break
			}
		}
	}

	moduleCacheEviction.Add(float64(evicted))
}

type module struct {
	module        host.ModuleV1
	lastFetchedAt time.Time
}
