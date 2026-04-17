package v2

import (
	"context"
	"errors"
	"fmt"
	"os"
	"sort"
	"sync"
	"sync/atomic"
	"time"
	"weak"

	"github.com/jonboulle/clockwork"

	"github.com/smartcontractkit/chainlink-common/pkg/workflows/wasm/host"
	sdkpb "github.com/smartcontractkit/chainlink-protos/cre/go/sdk"
	artifacts "github.com/smartcontractkit/chainlink/v2/core/services/workflows/artifacts/v2"
)

const defaultExecutePinMaxAttempts = 1024

// executePinMaxAttempts bounds how many times Execute retries after ensureLoaded races with Evict.
// It is mutable only so tests can lower it; production code leaves it at defaultExecutePinMaxAttempts.
var executePinMaxAttempts = defaultExecutePinMaxAttempts

// evictAfterEnsureLoadedHook is set by tests to force eviction after a successful ensureLoaded,
// exercising the pin retry loop and exhaustion path. It must be nil in production.
var evictAfterEnsureLoadedHook func(*EvictableModule)

// ErrExecutePinExhausted is returned when Execute could not observe a non-nil inner module after
// ensureLoaded, even after bounded retries (eviction repeatedly won the race before RLock).
var ErrExecutePinExhausted = errors.New("evictable module: failed to pin inner module after repeated eviction races")

// ErrEngineVersionMismatch is returned by ensureLoaded when the cached binary was persisted by
// a different engine version than the current one. The stale entry is deleted before returning.
var ErrEngineVersionMismatch = errors.New("evictable module: cached binary engine version mismatch")

// ModuleFactoryFn creates a host.ModuleV2 from config and binary.
// Defaults to host.NewModule in production; tests inject mocks via this.
type ModuleFactoryFn func(ctx context.Context, modCfg *host.ModuleConfig, binary []byte, opts ...func(*host.ModuleConfig)) (host.ModuleV2, error)

func defaultModuleFactory(ctx context.Context, modCfg *host.ModuleConfig, binary []byte, opts ...func(*host.ModuleConfig)) (host.ModuleV2, error) {
	return host.NewModule(ctx, modCfg, binary, opts...)
}

// EvictableModule wraps a host.ModuleV2 with idle-eviction and on-demand reload.
// Trigger registrations and event channels are owned by the engine, not by this module,
// so evicting the inner module only frees WASM memory without losing trigger connectivity.
type EvictableModule struct {
	inner         host.ModuleV2
	mu            sync.RWMutex
	lastUsed      atomic.Int64
	binarySize    atomic.Int64
	closed        atomic.Bool
	started       atomic.Bool
	workflowID    string
	engineVersion string

	moduleConfig *host.ModuleConfig
	moduleOpts   []func(*host.ModuleConfig)
	store        artifacts.SerialisedModuleStore
	factory      ModuleFactoryFn
	metrics      *CacheMetrics

	weakBinary weak.Pointer[[]byte] // L2: raw WASM binary survives eviction until GC pressure
}

func NewEvictableModule(
	inner host.ModuleV2,
	moduleConfig *host.ModuleConfig,
	store artifacts.SerialisedModuleStore,
	workflowID string,
	engineVersion string,
	factory ModuleFactoryFn,
	cm *CacheMetrics,
	initialBinaryLen int64,
	opts ...func(*host.ModuleConfig),
) *EvictableModule {
	if factory == nil {
		factory = defaultModuleFactory
	}
	m := &EvictableModule{
		inner:         inner,
		workflowID:    workflowID,
		engineVersion: engineVersion,
		moduleConfig:  moduleConfig,
		moduleOpts:    opts,
		store:         store,
		factory:       factory,
		metrics:       cm,
	}
	m.lastUsed.Store(time.Now().UnixNano())
	// Set from the bytes used to build inner (and written by StoreModule) so eviction
	// before any Execute/reload still contributes to memorySaved; ensureLoaded refreshes.
	if initialBinaryLen > 0 {
		m.binarySize.Store(initialBinaryLen)
	}
	return m
}

func (m *EvictableModule) Start() {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if m.inner != nil {
		m.inner.Start()
	}
	m.started.Store(true)
}

func (m *EvictableModule) Close() {
	m.closed.Store(true)
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.inner != nil {
		m.inner.Close()
		m.inner = nil
	}
}

func (m *EvictableModule) IsLegacyDAG() bool {
	return false
}

func (m *EvictableModule) Execute(ctx context.Context, request *sdkpb.ExecuteRequest, handler host.ExecutionHelper) (*sdkpb.ExecutionResult, error) {
	// the fundamental issue is that Go's RWMutex can't upgrade RLock to WLock
	// atomically. ensureLoaded and RLock acquisition cannot be atomic because
	// ensureLoaded may need the write lock internally to reload. This opens a
	// narrow window where Evict (called by the reaper under the write lock)
	// can nil-out m.inner between ensureLoaded returning and us grabbing the
	// RLock. The loop re-checks m.inner under the RLock and retries (bounded by
	// executePinMaxAttempts) if it was evicted in that gap; each iteration also
	// checks ctx so cancellation is not starved. Once pinned, RLock is held for
	// the duration of inner.Execute.
	var pinned host.ModuleV2
	for attempt := 0; attempt < executePinMaxAttempts; attempt++ {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		if err := m.ensureLoaded(ctx); err != nil {
			return nil, err
		}
		if h := evictAfterEnsureLoadedHook; h != nil {
			h(m)
		}
		m.lastUsed.Store(time.Now().UnixNano())
		m.mu.RLock()
		if m.inner != nil {
			pinned = m.inner
			break
		}
		m.mu.RUnlock()
	}
	if pinned == nil {
		return nil, fmt.Errorf("%w (workflow_id=%s attempts=%d)", ErrExecutePinExhausted, m.workflowID, executePinMaxAttempts)
	}
	defer m.mu.RUnlock()
	return pinned.Execute(ctx, request, handler)
}

func (m *EvictableModule) ensureLoaded(ctx context.Context) error {
	m.mu.RLock()
	if m.inner != nil {
		m.mu.RUnlock()
		return nil
	}
	m.mu.RUnlock()

	m.mu.Lock()
	defer m.mu.Unlock()

	if m.closed.Load() {
		return errors.New("module is permanently closed")
	}
	if m.inner != nil {
		return nil
	}

	var binary []byte
	var reloadSource string
	if bp := m.weakBinary.Value(); bp != nil {
		binary = *bp
		reloadSource = "weak_ref"
	} else {
		p, cachedVersion, ok, err := m.store.GetModule(m.workflowID)
		if err != nil {
			return fmt.Errorf("failed to get module path: %w", err)
		}
		if !ok {
			return fmt.Errorf("no cached binary for workflow %s", m.workflowID)
		}
		if cachedVersion != m.engineVersion {
			m.metrics.recordVersionMismatch(ctx)
			lg := m.moduleConfig.Logger
			if lg != nil {
				lg.Warnw("rejecting cached module binary: engine version mismatch",
					"workflowID", m.workflowID,
					"cachedEngineVersion", cachedVersion,
					"currentEngineVersion", m.engineVersion)
			}
			if delErr := m.store.DeleteModule(m.workflowID); delErr != nil && lg != nil {
				lg.Warnw("failed to delete stale cached module", "workflowID", m.workflowID, "err", delErr)
			}
			return fmt.Errorf("%w (workflow_id=%s cached=%q current=%q)", ErrEngineVersionMismatch, m.workflowID, cachedVersion, m.engineVersion)
		}

		binary, err = os.ReadFile(p)
		if err != nil {
			return fmt.Errorf("failed to read cached binary: %w", err)
		}
		reloadSource = "disk"
	}

	m.binarySize.Store(int64(len(binary)))

	mod, err := m.factory(ctx, m.moduleConfig, binary, m.moduleOpts...)
	if err != nil {
		return fmt.Errorf("failed to create module on reload: %w", err)
	}

	m.metrics.recordReload(ctx, reloadSource)

	binaryHeap := new([]byte)
	*binaryHeap = binary
	m.weakBinary = weak.Make(binaryHeap)

	if m.started.Load() {
		mod.Start()
	}
	m.inner = mod
	return nil
}

// Evict closes the inner module and frees its memory.
// A subsequent Execute call will reload from disk.
func (m *EvictableModule) Evict() {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.inner != nil && !m.closed.Load() {
		m.inner.Close()
		m.inner = nil
	}
}

// IsLoaded reports whether the inner module is currently in memory.
func (m *EvictableModule) IsLoaded() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.inner != nil
}

// LastUsed returns the last time Execute was called (unix nanoseconds).
func (m *EvictableModule) LastUsed() int64 {
	return m.lastUsed.Load()
}

// BinarySize returns the size of the WASM binary in bytes. It is populated from the
// disk cache when available at construction time and updated whenever the binary is loaded in ensureLoaded.
func (m *EvictableModule) BinarySize() int64 {
	return m.binarySize.Load()
}

// ModuleLRU tracks EvictableModule instances and periodically evicts idle ones.
type ModuleLRU struct {
	mu          sync.Mutex
	modules     map[string]*EvictableModule
	idleTimeout time.Duration
	maxLoaded   int
	clock       clockwork.Clock
	stopCh      chan struct{}
	wg          sync.WaitGroup
	metrics     *CacheMetrics

	// reapTicker drives the eviction scan. Injectable for deterministic tests.
	reapTicker <-chan time.Time
	// onReaped is signalled after each reap cycle completes (test hook only).
	onReaped chan struct{}
}

var (
	defaultIdleTimeout  = 10 * time.Minute
	defaultScanInterval = 30 * time.Second
)

func NewModuleLRU(clock clockwork.Clock, opts ...func(*ModuleLRU)) *ModuleLRU {
	lru := &ModuleLRU{
		modules:     make(map[string]*EvictableModule),
		idleTimeout: defaultIdleTimeout,
		clock:       clock,
		stopCh:      make(chan struct{}),
		reapTicker:  clock.NewTicker(defaultScanInterval).Chan(),
	}
	for _, o := range opts {
		o(lru)
	}
	return lru
}

func WithIdleTimeout(d time.Duration) func(*ModuleLRU) {
	return func(lru *ModuleLRU) {
		lru.idleTimeout = d
	}
}

func WithMaxLoadedModules(n int) func(*ModuleLRU) {
	return func(lru *ModuleLRU) {
		lru.maxLoaded = n
	}
}

func WithReapTicker(ch <-chan time.Time) func(*ModuleLRU) {
	return func(lru *ModuleLRU) {
		lru.reapTicker = ch
	}
}

func WithOnReaped(ch chan struct{}) func(*ModuleLRU) {
	return func(lru *ModuleLRU) {
		lru.onReaped = ch
	}
}

func WithCacheMetrics(cm *CacheMetrics) func(*ModuleLRU) {
	return func(lru *ModuleLRU) {
		lru.metrics = cm
	}
}

func (lru *ModuleLRU) Start() {
	lru.wg.Add(1)
	go func() {
		defer lru.wg.Done()
		lru.reapLoop()
	}()
}

func (lru *ModuleLRU) Close() {
	close(lru.stopCh)
	lru.wg.Wait()
}

func (lru *ModuleLRU) Register(workflowID string, m *EvictableModule) {
	lru.mu.Lock()
	defer lru.mu.Unlock()
	lru.modules[workflowID] = m
}

func (lru *ModuleLRU) Deregister(workflowID string) {
	lru.mu.Lock()
	defer lru.mu.Unlock()
	delete(lru.modules, workflowID)
}

func (lru *ModuleLRU) Contains(workflowID string) bool {
	lru.mu.Lock()
	defer lru.mu.Unlock()
	_, ok := lru.modules[workflowID]
	return ok
}

func (lru *ModuleLRU) reapLoop() {
	for {
		select {
		case <-lru.reapTicker:
			lru.reap()
			if lru.onReaped != nil {
				lru.onReaped <- struct{}{}
			}
		case <-lru.stopCh:
			return
		}
	}
}

func (lru *ModuleLRU) reap() {
	lru.mu.Lock()
	defer lru.mu.Unlock()

	evicted := 0

	if lru.idleTimeout > 0 {
		now := lru.clock.Now().UnixNano()
		threshold := lru.idleTimeout.Nanoseconds()
		for _, m := range lru.modules {
			if now-m.LastUsed() > threshold && m.IsLoaded() {
				m.Evict()
				evicted++
			}
		}
	}

	if lru.maxLoaded > 0 {
		evicted += lru.enforceCapLocked()
	}

	if evicted > 0 {
		lru.metrics.recordEviction(context.Background(), evicted)
	}

	loaded := 0
	var savedBytes int64
	for _, m := range lru.modules {
		if m.IsLoaded() {
			loaded++
		} else {
			savedBytes += m.binarySize.Load()
		}
	}
	lru.metrics.recordLoaded(context.Background(), loaded)
	lru.metrics.recordMemorySaved(context.Background(), savedBytes)
}

func (lru *ModuleLRU) enforceCapLocked() int {
	type entry struct {
		id       string
		lastUsed int64
	}

	var loaded []entry
	for id, m := range lru.modules {
		if m.IsLoaded() {
			loaded = append(loaded, entry{id: id, lastUsed: m.LastUsed()})
		}
	}

	excess := len(loaded) - lru.maxLoaded
	if excess <= 0 {
		return 0
	}

	sort.Slice(loaded, func(i, j int) bool {
		return loaded[i].lastUsed < loaded[j].lastUsed
	})

	evicted := 0
	for i := 0; i < excess; i++ {
		if m, ok := lru.modules[loaded[i].id]; ok {
			m.Evict()
			evicted++
		}
	}
	return evicted
}
