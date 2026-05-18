package v2

import (
	"context"
	"errors"
	"fmt"
	"os"
	"runtime"
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
var tryAcquireCompareAndSwap func(e *loadedModule, old, next int64) bool

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

// onceCloseModule wraps host.ModuleV2 with an idempotent Close. Both
// runtime.AddCleanup (fired when GC reaps the holder) and explicit
// EvictableModule.Close target this method, and they may run in either order;
// sync.Once ensures host.module.Close (which is not idempotent and would panic
// on double-close) runs at most once across both paths.
//
// onceCloseModule is separate from loadedModule because runtime.AddCleanup
// forbids the cleanup arg from transitively referencing the cleanup target.
type onceCloseModule struct {
	mod       host.ModuleV2
	closeOnce sync.Once
}

func (m *onceCloseModule) Start() { m.mod.Start() }

func (m *onceCloseModule) Execute(ctx context.Context, request *sdkpb.ExecuteRequest, handler host.ExecutionHelper) (*sdkpb.ExecutionResult, error) {
	return m.mod.Execute(ctx, request, handler)
}

func (m *onceCloseModule) Close() {
	m.closeOnce.Do(m.mod.Close)
}

// loadedModule is the weak-pointable holder around a compiled host.ModuleV2.
// EvictableModule holds it strongly while loaded (L1) and weakly after Evict
// (L2), so a reload before GC pressure can resurrect the same compiled module
// and skip both the disk read and the wasm compilation cost.
//
// refCount enables non-blocking Evict: 1 represents the owning reference held
// via EvictableModule.current, +1 per active Execute pin. Evict drops the
// owning ref without waiting; in-flight pins keep the holder reachable until
// they release. The actual mod.Close runs via runtime.AddCleanup once the
// holder becomes unreachable (no strong refs, no pins) and GC reclaims it,
// or eagerly via EvictableModule.Close on shutdown.
type loadedModule struct {
	mod      *onceCloseModule
	cleanup  runtime.Cleanup
	refCount atomic.Int64
}

func newLoadedModule(mod host.ModuleV2) *loadedModule {
	wrapped := &onceCloseModule{mod: mod}
	h := &loadedModule{mod: wrapped}
	h.cleanup = runtime.AddCleanup(h, (*onceCloseModule).Close, wrapped)
	h.refCount.Store(1)
	return h
}

// tryAcquire pins the entry by incrementing refCount only if it is non-zero.
// A plain Add(1) would race with a release that already drove the count to zero
// and called Close, leading to a use-after-close. CAS makes the increment
// conditional on the entry still being live.
func (e *loadedModule) tryAcquire() (acquired bool, exhausted bool) {
	cas := tryAcquireCompareAndSwap
	if cas == nil {
		cas = func(e *loadedModule, old, next int64) bool {
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

// release drops one ref. It deliberately does NOT call mod.Close on zero so the
// weak L2 cache can resurrect the holder before GC reaps it. Final close runs
// via the runtime.AddCleanup callback (or eagerly via EvictableModule.Close).
func (e *loadedModule) release() { e.refCount.Add(-1) }

// EvictableModule wraps a host.ModuleV2 with idle-eviction and on-demand reload.
// Trigger registrations and event channels are owned by the engine, not by this module,
// so evicting the inner module only frees WASM memory without losing trigger connectivity.
type EvictableModule struct {
	current       atomic.Pointer[loadedModule] // L1: strong, refcounted; cleared by Evict and Close
	weakInner     weak.Pointer[loadedModule]   // L2: survives eviction until GC reclaims the holder
	mu            sync.Mutex                   // guards weakInner and serializes ensureLoaded reloads; never held during inner.Execute
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
		holder := newLoadedModule(inner)
		m.current.Store(holder)
		m.weakInner = weak.Make(holder)
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
	if m.closed.Swap(true) {
		return
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	// weakInner is the live holder (current is a strong subset of it). Drop
	// the owning refcount and close via the holder; sync.Once makes the
	// eventual cleanup-driven Close (after GC reaps the holder) a no-op.
	h := m.weakInner.Value()
	if strong := m.current.Swap(nil); strong != nil {
		strong.release()
	}
	if h != nil {
		h.mod.Close()
	}
}

func (m *EvictableModule) IsLegacyDAG() bool {
	return false
}

func (m *EvictableModule) Execute(ctx context.Context, request *sdkpb.ExecuteRequest, handler host.ExecutionHelper) (*sdkpb.ExecutionResult, error) {
	// Each loaded module is held behind a refcounted loadedModule. Pinning is
	// a CAS-conditional refcount increment: it succeeds only if the holder is
	// still live (count > 0). Evict drops the owning ref atomically without
	// blocking, so in-flight pins keep the holder reachable until they release;
	// final close runs via runtime.AddCleanup once GC reaps the holder, or
	// eagerly via EvictableModule.Close. ensureLoaded and pin are not atomic,
	// so we keep a bounded retry loop for the case where Evict fires between
	// ensureLoaded returning and pin; each iteration also checks ctx so
	// cancellation is not starved.
	var pinned *loadedModule
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

	// L2: try to resurrect the still-live compiled module weakly held since
	// Evict. weak.Pointer.Value returns nil once GC has reclaimed the holder
	// (and cleanup is therefore queued or done), so a non-nil result means
	// the holder is reachable and its mod is still open. We deliberately do
	// NOT call h.cleanup.Stop here: the cleanup is the only path that closes
	// the module when this resurrected holder is later evicted and reaped,
	// and re-promoting it via current strongly references it again, so GC
	// won't fire until the next eviction makes it unreachable once more.
	if h := m.weakInner.Value(); h != nil {
		h.refCount.Add(1) // re-establish owning ref (weakly-held holder had refCount 0)
		m.metrics.recordReload(ctx, "weak_ref")
		m.current.Store(h)
		return nil
	}

	// L3: read binary from disk and re-instantiate via the factory.
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

	binary, err := os.ReadFile(p)
	if err != nil {
		return fmt.Errorf("failed to read cached binary: %w", err)
	}

	m.binarySize.Store(int64(len(binary)))

	mod, err := m.factory(ctx, m.moduleConfig, binary, m.moduleOpts...)
	if err != nil {
		return fmt.Errorf("failed to create module on reload: %w", err)
	}

	m.metrics.recordReload(ctx, "disk")

	if m.started.Load() {
		mod.Start()
	}
	holder := newLoadedModule(mod)
	m.current.Store(holder)
	m.weakInner = weak.Make(holder)
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
	// onReaped is signaled after each reap cycle completes (test hook only).
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
