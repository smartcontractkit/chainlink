package v2

import (
	"context"
	"fmt"
	"os"
	"sort"
	"sync"
	"sync/atomic"
	"time"

	"github.com/jonboulle/clockwork"
	sdkpb "github.com/smartcontractkit/chainlink-protos/cre/go/sdk"
	"github.com/smartcontractkit/chainlink-common/pkg/workflows/wasm/host"

	artifacts "github.com/smartcontractkit/chainlink/v2/core/services/workflows/artifacts/v2"
)

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
	inner      host.ModuleV2
	mu         sync.RWMutex
	lastUsed   atomic.Int64
	closed     atomic.Bool
	started    atomic.Bool
	workflowID string

	moduleConfig *host.ModuleConfig
	moduleOpts   []func(*host.ModuleConfig)
	store        artifacts.SerialisedModuleStore
	factory      ModuleFactoryFn
}

func NewEvictableModule(
	inner host.ModuleV2,
	moduleConfig *host.ModuleConfig,
	store artifacts.SerialisedModuleStore,
	workflowID string,
	factory ModuleFactoryFn,
	opts ...func(*host.ModuleConfig),
) *EvictableModule {
	if factory == nil {
		factory = defaultModuleFactory
	}
	m := &EvictableModule{
		inner:        inner,
		workflowID:   workflowID,
		moduleConfig: moduleConfig,
		moduleOpts:   opts,
		store:        store,
		factory:      factory,
	}
	m.lastUsed.Store(time.Now().UnixNano())
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
	if err := m.ensureLoaded(ctx); err != nil {
		return nil, err
	}
	m.lastUsed.Store(time.Now().UnixNano())
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.inner.Execute(ctx, request, handler)
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
		return fmt.Errorf("module is permanently closed")
	}
	if m.inner != nil {
		return nil
	}

	p, ok, err := m.store.GetModulePath(m.workflowID)
	if err != nil {
		return fmt.Errorf("failed to get module path: %w", err)
	}
	if !ok {
		return fmt.Errorf("no cached binary for workflow %s", m.workflowID)
	}

	binary, err := os.ReadFile(p)
	if err != nil {
		return fmt.Errorf("failed to read cached binary: %w", err)
	}

	mod, err := m.factory(ctx, m.moduleConfig, binary, m.moduleOpts...)
	if err != nil {
		return fmt.Errorf("failed to create module on reload: %w", err)
	}

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

// ModuleLRU tracks EvictableModule instances and periodically evicts idle ones.
type ModuleLRU struct {
	mu          sync.Mutex
	modules     map[string]*EvictableModule
	idleTimeout time.Duration
	maxLoaded   int
	clock       clockwork.Clock
	stopCh      chan struct{}
	wg          sync.WaitGroup

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
		if d > 0 {
			lru.idleTimeout = d
		}
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

	now := lru.clock.Now().UnixNano()
	threshold := lru.idleTimeout.Nanoseconds()

	for _, m := range lru.modules {
		if now-m.LastUsed() > threshold && m.IsLoaded() {
			m.Evict()
		}
	}

	if lru.maxLoaded > 0 {
		lru.enforceCapLocked()
	}
}

func (lru *ModuleLRU) enforceCapLocked() {
	type entry struct {
		id       string
		lastUsed int64
		loaded   bool
	}

	var loaded []entry
	for id, m := range lru.modules {
		if m.IsLoaded() {
			loaded = append(loaded, entry{id: id, lastUsed: m.LastUsed(), loaded: true})
		}
	}

	excess := len(loaded) - lru.maxLoaded
	if excess <= 0 {
		return
	}

	sort.Slice(loaded, func(i, j int) bool {
		return loaded[i].lastUsed < loaded[j].lastUsed
	})

	for i := 0; i < excess; i++ {
		if m, ok := lru.modules[loaded[i].id]; ok {
			m.Evict()
		}
	}
}
