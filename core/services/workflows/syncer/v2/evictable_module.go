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

const (
	defaultTryAcquireMaxAttempts = 1024
	defaultExecutePinMaxAttempts = 1024
)

// tryAcquireMaxAttempts bounds CAS retries while pinning a moduleEntry.
// It is mutable only so tests can lower it; production code leaves it at defaultTryAcquireMaxAttempts.
var tryAcquireMaxAttempts = defaultTryAcquireMaxAttempts

// executePinMaxAttempts bounds how many times Execute retries after ensureLoaded races with Evict.
// It is mutable only so tests can lower it; production code leaves it at defaultExecutePinMaxAttempts.
var executePinMaxAttempts = defaultExecutePinMaxAttempts

// tryAcquireCompareAndSwap is injectable for tests to deterministically force CAS contention.
// It must be left as nil in production (defaulting to moduleEntry.refCount.CompareAndSwap).
var tryAcquireCompareAndSwap func(e *moduleEntry, old, next int64) bool

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

// moduleEntry wraps a host.ModuleV2 with a refcount so Evict can drop ownership
// without waiting for in-flight Execute calls. The owning ref held via
// EvictableModule.current counts as 1; each successful tryAcquire adds 1.
// The release that drives the count to zero closes the inner module exactly once.
type moduleEntry struct {
	mod      host.ModuleV2
	refCount atomic.Int64
}

func newModuleEntry(mod host.ModuleV2) *moduleEntry {
	e := &moduleEntry{mod: mod}
	e.refCount.Store(1)
	return e
}

// tryAcquire pins the entry by incrementing refCount only if it is non-zero.
// A plain Add(1) would race with a release that already drove the count to zero
// and called Close, leading to a use-after-close. CAS makes the increment
// conditional on the entry still being live.
func (e *moduleEntry) tryAcquire() (acquired bool, exhausted bool) {
	cas := tryAcquireCompareAndSwap
	if cas == nil {
		cas = func(e *moduleEntry, old, next int64) bool {
			return e.refCount.CompareAndSwap(old, next)
		}
	}
	for attempt := 0; attempt < tryAcquireMaxAttempts; attempt++ {
		n := e.refCount.Load()
		if n == 0 {
			return false, false
		}
		if cas(e, n, n+1) {
			return true, false
		}
	}
	return false, true
}

func (e *moduleEntry) release() {
	if e.refCount.Add(-1) == 0 {
		e.mod.Close()
	}
}

// EvictableModule wraps a host.ModuleV2 with idle-eviction and on-demand reload.
// Trigger registrations and event channels are owned by the engine, not by this module,
// so evicting the inner module only frees WASM memory without losing trigger connectivity.
type EvictableModule struct {
	current       atomic.Pointer[moduleEntry]
	mu            sync.Mutex // serializes ensureLoaded reloads; never held during inner.Execute
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
		workflowID:    workflowID,
		engineVersion: engineVersion,
		moduleConfig:  moduleConfig,
		moduleOpts:    opts,
		store:         store,
		factory:       factory,
		metrics:       cm,
	}
	if inner != nil {
		m.current.Store(newModuleEntry(inner))
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
	if e := m.current.Load(); e != nil {
		acquired, _ := e.tryAcquire()
		if acquired {
			defer e.release()
			e.mod.Start()
		}
	}
	m.started.Store(true)
}

func (m *EvictableModule) Close() {
	m.closed.Store(true)
	if e := m.current.Swap(nil); e != nil {
		e.release()
	}
}

func (m *EvictableModule) IsLegacyDAG() bool {
	return false
}

func (m *EvictableModule) Execute(ctx context.Context, request *sdkpb.ExecuteRequest, handler host.ExecutionHelper) (*sdkpb.ExecutionResult, error) {
	// Each loaded module is held behind a refcounted moduleEntry. Pinning is a
	// CAS-conditional refcount increment: it succeeds only if the entry is still
	// live (count > 0). Evict drops the owning ref atomically, so in-flight pins
	// keep the entry alive until they release; the last release calls Close.
	// ensureLoaded and pin are still not atomic, so we keep a bounded retry loop
	// for the case where Evict fires between ensureLoaded returning and pin: each
	// iteration also checks ctx so cancellation is not starved.
	var pinned *moduleEntry
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
		if e := m.current.Load(); e != nil {
			acquired, exhausted := e.tryAcquire()
			if acquired {
				pinned = e
				break
			}
			if exhausted {
				m.metrics.recordTryAcquireExhausted(ctx)
			}
		}
	}
	if pinned == nil {
		m.metrics.recordPinExhausted(ctx)
		return nil, fmt.Errorf("%w (workflow_id=%s attempts=%d)", ErrExecutePinExhausted, m.workflowID, executePinMaxAttempts)
	}
	defer pinned.release()
	return pinned.mod.Execute(ctx, request, handler)
}

func (m *EvictableModule) ensureLoaded(ctx context.Context) error {
	if m.current.Load() != nil {
		return nil
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	if m.closed.Load() {
		return errors.New("module is permanently closed")
	}
	if m.current.Load() != nil {
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
	m.current.Store(newModuleEntry(mod))
	return nil
}

// Evict drops the owning reference to the inner module only when no Execute call
// is currently pinned on that entry.
//
// Why this check exists:
//   - If we clear m.current while another Execute still holds a pin, the entry
//     remains alive (refcount > 0) but becomes unreachable through m.current.
//   - A concurrent/new Execute that observes m.current == nil will run
//     ensureLoaded and instantiate another module from weak-ref/disk.
//   - That creates transient duplicate module instances for one workflow:
//     the old one still serving in-flight work and a new one for subsequent
//     work. This is safe, but unnecessarily increases memory churn and defeats
//     the eviction intent under contention.
//
// By refusing eviction while refcount > 1 (owner + at least one pin), we make
// eviction eventually consistent: reap/cap may skip a busy module in this cycle
// and retry later, but we avoid duplicate live instances caused by evicting a
// still-pinned entry.
//
// Evict remains non-blocking and single-pass: it performs one CAS attempt and
// returns. If that CAS loses a race with a concurrent load/reload/evict, the
// caller can retry on the next reap/cap cycle.
func (m *EvictableModule) Evict() {
	if m.closed.Load() {
		return
	}
	e := m.current.Load()
	if e == nil {
		return
	}
	if e.refCount.Load() > 1 {
		return
	}
	if m.current.CompareAndSwap(e, nil) {
		e.release()
	}
}

// IsLoaded reports whether the inner module is currently in memory.
func (m *EvictableModule) IsLoaded() bool {
	return m.current.Load() != nil
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

	// reapMemorySavedHook receives savedBytes immediately before recordMemorySaved; tests observe real reap metric inputs.
	reapMemorySavedHook func(int64)
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
	if reapMemorySavedHook != nil {
		reapMemorySavedHook(savedBytes)
	}
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
