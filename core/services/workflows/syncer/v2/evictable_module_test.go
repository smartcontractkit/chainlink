package v2

import (
	"context"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
	"time"
	"weak"

	"github.com/jonboulle/clockwork"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/workflows/wasm/host"
	modulemocks "github.com/smartcontractkit/chainlink-common/pkg/workflows/wasm/host/mocks"
	sdkpb "github.com/smartcontractkit/chainlink-protos/cre/go/sdk"
	artifacts "github.com/smartcontractkit/chainlink/v2/core/services/workflows/artifacts/v2"
)

// countingStore wraps a SerialisedModuleStore and counts GetModulePath calls
// to verify whether disk I/O occurred during module reload.
type countingStore struct {
	artifacts.SerialisedModuleStore
	getModulePathCalls atomic.Int32
}

func (s *countingStore) GetModulePath(wfID string) (string, bool, error) {
	s.getModulePathCalls.Add(1)
	return s.SerialisedModuleStore.GetModulePath(wfID)
}

func newTestEvictableModule(t *testing.T, inner host.ModuleV2, factory ModuleFactoryFn) (*EvictableModule, artifacts.SerialisedModuleStore) {
	t.Helper()
	store, err := artifacts.NewFileModuleStore(t.TempDir())
	require.NoError(t, err)
	require.NoError(t, store.StoreModule("wf-test", "bin-1", []byte("fake-binary")))
	em := NewEvictableModule(inner, &host.ModuleConfig{}, store, "wf-test", factory, nil, int64(len("fake-binary")))
	return em, store
}

func TestEvictable_DelegatesToInner(t *testing.T) {
	inner := modulemocks.NewModuleV2(t)
	inner.EXPECT().Start()
	inner.EXPECT().Execute(mock.Anything, mock.Anything, mock.Anything).Return(&sdkpb.ExecutionResult{}, nil)

	em, _ := newTestEvictableModule(t, inner, nil)
	em.Start()

	result, err := em.Execute(context.Background(), &sdkpb.ExecuteRequest{}, nil)
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.False(t, em.IsLegacyDAG())
}

func TestEvictable_LastUsedUpdated(t *testing.T) {
	inner := modulemocks.NewModuleV2(t)
	inner.EXPECT().Execute(mock.Anything, mock.Anything, mock.Anything).Return(&sdkpb.ExecutionResult{}, nil).Times(2)

	em, _ := newTestEvictableModule(t, inner, nil)

	before := em.LastUsed()
	time.Sleep(time.Millisecond)

	_, err := em.Execute(context.Background(), &sdkpb.ExecuteRequest{}, nil)
	require.NoError(t, err)
	after1 := em.LastUsed()
	assert.Greater(t, after1, before)

	time.Sleep(time.Millisecond)
	_, err = em.Execute(context.Background(), &sdkpb.ExecuteRequest{}, nil)
	require.NoError(t, err)
	assert.Greater(t, em.LastUsed(), after1)
}

func TestEvictable_EvictFreesModule(t *testing.T) {
	inner := modulemocks.NewModuleV2(t)
	inner.EXPECT().Close()

	em, _ := newTestEvictableModule(t, inner, nil)
	assert.True(t, em.IsLoaded())

	em.Evict()
	assert.False(t, em.IsLoaded())
}

func TestEvictable_ReloadFromDisk(t *testing.T) {
	inner := modulemocks.NewModuleV2(t)
	inner.EXPECT().Close()

	var reloadedBinary []byte
	reloaded := modulemocks.NewModuleV2(t)
	reloaded.EXPECT().Start()
	reloaded.EXPECT().Execute(mock.Anything, mock.Anything, mock.Anything).Return(&sdkpb.ExecutionResult{}, nil)

	factory := func(_ context.Context, _ *host.ModuleConfig, binary []byte, _ ...func(*host.ModuleConfig)) (host.ModuleV2, error) {
		reloadedBinary = make([]byte, len(binary))
		copy(reloadedBinary, binary)
		return reloaded, nil
	}

	em, _ := newTestEvictableModule(t, inner, factory)
	em.started.Store(true)
	em.Evict()

	_, err := em.Execute(context.Background(), &sdkpb.ExecuteRequest{}, nil)
	require.NoError(t, err)
	assert.True(t, em.IsLoaded())
	assert.Equal(t, []byte("fake-binary"), reloadedBinary)
}

func TestEvictable_ReloadCallsStart(t *testing.T) {
	inner := modulemocks.NewModuleV2(t)
	inner.EXPECT().Start()
	inner.EXPECT().Close()

	reloaded := modulemocks.NewModuleV2(t)
	reloaded.EXPECT().Start()
	reloaded.EXPECT().Execute(mock.Anything, mock.Anything, mock.Anything).Return(&sdkpb.ExecutionResult{}, nil)

	factory := func(_ context.Context, _ *host.ModuleConfig, _ []byte, _ ...func(*host.ModuleConfig)) (host.ModuleV2, error) {
		return reloaded, nil
	}

	em, _ := newTestEvictableModule(t, inner, factory)
	em.Start()
	em.Evict()

	_, err := em.Execute(context.Background(), &sdkpb.ExecuteRequest{}, nil)
	require.NoError(t, err)
}

func TestEvictable_ClosePreventsReload(t *testing.T) {
	inner := modulemocks.NewModuleV2(t)
	inner.EXPECT().Close()

	em, _ := newTestEvictableModule(t, inner, nil)
	em.Close()

	_, err := em.Execute(context.Background(), &sdkpb.ExecuteRequest{}, nil)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "permanently closed")
	assert.False(t, em.IsLoaded())
}

func TestEvictable_ConcurrentExecuteDuringEvict(t *testing.T) {
	inner := modulemocks.NewModuleV2(t)
	inner.EXPECT().Close()
	inner.EXPECT().Execute(mock.Anything, mock.Anything, mock.Anything).
		Return(&sdkpb.ExecutionResult{}, nil).Maybe()

	factory := func(_ context.Context, _ *host.ModuleConfig, _ []byte, _ ...func(*host.ModuleConfig)) (host.ModuleV2, error) {
		m := modulemocks.NewModuleV2(t)
		m.EXPECT().Start().Maybe()
		m.EXPECT().Close().Maybe()
		m.EXPECT().Execute(mock.Anything, mock.Anything, mock.Anything).
			Return(&sdkpb.ExecutionResult{}, nil).Maybe()
		return m, nil
	}

	em, _ := newTestEvictableModule(t, inner, factory)
	em.started.Store(true)

	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, _ = em.Execute(context.Background(), &sdkpb.ExecuteRequest{}, nil)
		}()
	}
	wg.Add(1)
	go func() {
		defer wg.Done()
		em.Evict()
	}()
	wg.Wait()
}

func TestEvictable_EvictWaitsForExecution(t *testing.T) {
	var executing atomic.Bool
	inner := modulemocks.NewModuleV2(t)
	inner.EXPECT().Execute(mock.Anything, mock.Anything, mock.Anything).RunAndReturn(
		func(_ context.Context, _ *sdkpb.ExecuteRequest, _ host.ExecutionHelper) (*sdkpb.ExecutionResult, error) {
			executing.Store(true)
			time.Sleep(50 * time.Millisecond)
			executing.Store(false)
			return &sdkpb.ExecutionResult{}, nil
		},
	)
	inner.EXPECT().Close()

	em, _ := newTestEvictableModule(t, inner, nil)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		_, _ = em.Execute(context.Background(), &sdkpb.ExecuteRequest{}, nil)
	}()

	// Give Execute goroutine time to acquire RLock
	time.Sleep(10 * time.Millisecond)
	assert.True(t, executing.Load(), "Execute should still be running")

	em.Evict()
	// By the time Evict returns (it needed WLock), Execute should have finished
	assert.False(t, executing.Load())
	wg.Wait()
}

func TestEvictable_MultipleEvictReloadCycles(t *testing.T) {
	var createCount atomic.Int32

	inner := modulemocks.NewModuleV2(t)
	inner.EXPECT().Close()

	factory := func(_ context.Context, _ *host.ModuleConfig, _ []byte, _ ...func(*host.ModuleConfig)) (host.ModuleV2, error) {
		createCount.Add(1)
		m := modulemocks.NewModuleV2(t)
		m.EXPECT().Start()
		m.EXPECT().Close()
		m.EXPECT().Execute(mock.Anything, mock.Anything, mock.Anything).
			Return(&sdkpb.ExecutionResult{}, nil)
		return m, nil
	}

	em, _ := newTestEvictableModule(t, inner, factory)
	em.started.Store(true)

	for i := 0; i < 3; i++ {
		em.Evict()
		assert.False(t, em.IsLoaded())

		_, err := em.Execute(context.Background(), &sdkpb.ExecuteRequest{}, nil)
		require.NoError(t, err)
		assert.True(t, em.IsLoaded())
	}

	assert.Equal(t, int32(3), createCount.Load())

	// Final cleanup
	em.Close()
}

func TestEvictable_ReloadFailure(t *testing.T) {
	inner := modulemocks.NewModuleV2(t)
	inner.EXPECT().Close()

	em, store := newTestEvictableModule(t, inner, nil)
	em.Evict()

	// Corrupt the cache by deleting the binary
	require.NoError(t, store.DeleteModule("wf-test"))

	_, err := em.Execute(context.Background(), &sdkpb.ExecuteRequest{}, nil)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "no cached binary")
	assert.False(t, em.IsLoaded())
}

// --- ModuleLRU tests ---

func newLRUModule(t *testing.T, store artifacts.SerialisedModuleStore, wfID string) *EvictableModule {
	t.Helper()
	inner := modulemocks.NewModuleV2(t)
	inner.EXPECT().Close().Maybe()
	require.NoError(t, store.StoreModule(wfID, "bin", []byte("binary")))
	factory := func(_ context.Context, _ *host.ModuleConfig, _ []byte, _ ...func(*host.ModuleConfig)) (host.ModuleV2, error) {
		m := modulemocks.NewModuleV2(t)
		m.EXPECT().Start().Maybe()
		m.EXPECT().Close().Maybe()
		m.EXPECT().Execute(mock.Anything, mock.Anything, mock.Anything).Return(&sdkpb.ExecutionResult{}, nil).Maybe()
		return m, nil
	}
	em := NewEvictableModule(inner, &host.ModuleConfig{}, store, wfID, factory, nil, int64(len("binary")))
	em.started.Store(true)
	return em
}

func TestLRU_EvictsIdleModule(t *testing.T) {
	clock := clockwork.NewFakeClock()
	reapTicker := make(chan time.Time, 1)
	onReaped := make(chan struct{}, 1)

	store, err := artifacts.NewFileModuleStore(t.TempDir())
	require.NoError(t, err)

	em := newLRUModule(t, store, "wf-idle")
	em.lastUsed.Store(clock.Now().UnixNano())

	lru := NewModuleLRU(clock, WithIdleTimeout(5*time.Minute), WithReapTicker(reapTicker), WithOnReaped(onReaped))
	lru.Register("wf-idle", em)
	lru.Start()
	defer lru.Close()

	assert.True(t, em.IsLoaded())

	clock.Advance(6 * time.Minute)
	reapTicker <- clock.Now()
	<-onReaped

	assert.False(t, em.IsLoaded())
}

func TestLRU_ActiveModuleNotEvicted(t *testing.T) {
	clock := clockwork.NewFakeClock()
	reapTicker := make(chan time.Time, 1)
	onReaped := make(chan struct{}, 1)

	store, err := artifacts.NewFileModuleStore(t.TempDir())
	require.NoError(t, err)

	em := newLRUModule(t, store, "wf-active")
	em.lastUsed.Store(clock.Now().UnixNano())

	lru := NewModuleLRU(clock, WithIdleTimeout(5*time.Minute), WithReapTicker(reapTicker), WithOnReaped(onReaped))
	lru.Register("wf-active", em)
	lru.Start()
	defer lru.Close()

	clock.Advance(3 * time.Minute)
	// Simulate activity: update lastUsed
	em.lastUsed.Store(clock.Now().UnixNano())

	clock.Advance(3 * time.Minute)
	reapTicker <- clock.Now()
	<-onReaped

	assert.True(t, em.IsLoaded(), "active module should not be evicted")
}

func TestLRU_MaxLoadedCap(t *testing.T) {
	clock := clockwork.NewFakeClock()
	reapTicker := make(chan time.Time, 1)
	onReaped := make(chan struct{}, 1)

	store, err := artifacts.NewFileModuleStore(t.TempDir())
	require.NoError(t, err)

	m1 := newLRUModule(t, store, "wf-1")
	m1.lastUsed.Store(clock.Now().Add(-3 * time.Minute).UnixNano())

	m2 := newLRUModule(t, store, "wf-2")
	m2.lastUsed.Store(clock.Now().Add(-2 * time.Minute).UnixNano())

	m3 := newLRUModule(t, store, "wf-3")
	m3.lastUsed.Store(clock.Now().Add(-1 * time.Minute).UnixNano())

	lru := NewModuleLRU(clock,
		WithIdleTimeout(1*time.Hour),
		WithMaxLoadedModules(2),
		WithReapTicker(reapTicker),
		WithOnReaped(onReaped),
	)
	lru.Register("wf-1", m1)
	lru.Register("wf-2", m2)
	lru.Register("wf-3", m3)
	lru.Start()
	defer lru.Close()

	reapTicker <- clock.Now()
	<-onReaped

	assert.False(t, m1.IsLoaded(), "oldest module should be evicted")
	assert.True(t, m2.IsLoaded())
	assert.True(t, m3.IsLoaded())
}

func TestLRU_DeregisterStopsTracking(t *testing.T) {
	clock := clockwork.NewFakeClock()
	reapTicker := make(chan time.Time, 1)
	onReaped := make(chan struct{}, 1)

	store, err := artifacts.NewFileModuleStore(t.TempDir())
	require.NoError(t, err)

	em := newLRUModule(t, store, "wf-dereg")
	em.lastUsed.Store(clock.Now().UnixNano())

	lru := NewModuleLRU(clock, WithIdleTimeout(5*time.Minute), WithReapTicker(reapTicker), WithOnReaped(onReaped))
	lru.Register("wf-dereg", em)
	lru.Deregister("wf-dereg")
	lru.Start()
	defer lru.Close()

	clock.Advance(10 * time.Minute)
	reapTicker <- clock.Now()
	<-onReaped

	assert.True(t, em.IsLoaded(), "deregistered module should not be evicted by LRU")
}

func TestLRU_ConcurrentRegisterDeregister(t *testing.T) {
	clock := clockwork.NewFakeClock()
	lru := NewModuleLRU(clock)

	store, err := artifacts.NewFileModuleStore(t.TempDir())
	require.NoError(t, err)

	type entry struct {
		wfID string
		em   *EvictableModule
	}
	entries := make([]entry, 20)
	for i := 0; i < 20; i++ {
		wfID := string(rune('A' + i))
		entries[i] = entry{wfID: wfID, em: newLRUModule(t, store, wfID)}
	}

	var wg sync.WaitGroup
	for i := range entries {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			lru.Register(entries[idx].wfID, entries[idx].em)
			lru.Deregister(entries[idx].wfID)
		}(i)
	}
	wg.Wait()
}

func TestLRU_StartStop(t *testing.T) {
	clock := clockwork.NewFakeClock()
	reapTicker := make(chan time.Time, 1)
	onReaped := make(chan struct{}, 1)

	store, err := artifacts.NewFileModuleStore(t.TempDir())
	require.NoError(t, err)

	em := newLRUModule(t, store, "wf-1")
	em.lastUsed.Store(clock.Now().UnixNano())

	lru := NewModuleLRU(clock, WithIdleTimeout(1*time.Second), WithReapTicker(reapTicker), WithOnReaped(onReaped))
	lru.Register("wf-1", em)
	lru.Start()
	lru.Close()

	// After Close, sending to reapTicker should not cause eviction (loop exited)
	assert.True(t, em.IsLoaded())
}

func TestLRU_EmptyScan(t *testing.T) {
	clock := clockwork.NewFakeClock()
	reapTicker := make(chan time.Time, 1)
	onReaped := make(chan struct{}, 1)

	lru := NewModuleLRU(clock, WithReapTicker(reapTicker), WithOnReaped(onReaped))
	lru.Start()
	defer lru.Close()

	reapTicker <- clock.Now()
	<-onReaped
}

func TestLRU_EvictionOrder(t *testing.T) {
	clock := clockwork.NewFakeClock()
	reapTicker := make(chan time.Time, 1)
	onReaped := make(chan struct{}, 1)

	store, err := artifacts.NewFileModuleStore(t.TempDir())
	require.NoError(t, err)

	modules := make([]*EvictableModule, 5)
	for i := 0; i < 5; i++ {
		wfID := string(rune('A' + i))
		modules[i] = newLRUModule(t, store, wfID)
		modules[i].lastUsed.Store(clock.Now().Add(time.Duration(i) * time.Minute).UnixNano())
	}

	lru := NewModuleLRU(clock,
		WithIdleTimeout(1*time.Hour),
		WithMaxLoadedModules(2),
		WithReapTicker(reapTicker),
		WithOnReaped(onReaped),
	)
	for i, m := range modules {
		lru.Register(string(rune('A'+i)), m)
	}
	lru.Start()
	defer lru.Close()

	reapTicker <- clock.Now()
	<-onReaped

	for i, m := range modules {
		if i < 3 {
			assert.False(t, m.IsLoaded(), "module %d should be evicted", i)
		} else {
			assert.True(t, m.IsLoaded(), "module %d should survive", i)
		}
	}
}

// --- Weak reference (L2 cache) tests ---

func TestEvictable_WeakRefHitAfterEvict(t *testing.T) {
	inner := modulemocks.NewModuleV2(t)
	inner.EXPECT().Close()

	var receivedBinary []byte
	reloaded := modulemocks.NewModuleV2(t)
	reloaded.EXPECT().Start()
	reloaded.EXPECT().Execute(mock.Anything, mock.Anything, mock.Anything).Return(&sdkpb.ExecutionResult{}, nil)

	factory := func(_ context.Context, _ *host.ModuleConfig, binary []byte, _ ...func(*host.ModuleConfig)) (host.ModuleV2, error) {
		receivedBinary = make([]byte, len(binary))
		copy(receivedBinary, binary)
		return reloaded, nil
	}

	realStore, err := artifacts.NewFileModuleStore(t.TempDir())
	require.NoError(t, err)
	require.NoError(t, realStore.StoreModule("wf-test", "bin-1", []byte("disk-binary")))
	cs := &countingStore{SerialisedModuleStore: realStore}

	em := NewEvictableModule(inner, &host.ModuleConfig{}, cs, "wf-test", factory, nil, int64(len("disk-binary")))
	em.started.Store(true)

	// Seed the weak reference with a known binary. Hold a strong reference
	// to binaryHeap so the GC cannot collect it during the test.
	binaryData := []byte("weak-ref-binary")
	binaryHeap := new([]byte)
	*binaryHeap = binaryData
	em.weakBinary = weak.Make(binaryHeap)

	em.Evict()

	_, err = em.Execute(context.Background(), &sdkpb.ExecuteRequest{}, nil)
	require.NoError(t, err)

	assert.Equal(t, int32(0), cs.getModulePathCalls.Load(), "disk should not be accessed when weak ref is alive")
	assert.Equal(t, []byte("weak-ref-binary"), receivedBinary)
	runtime.KeepAlive(binaryHeap)
}

func TestEvictable_WeakRefMissFallsToDisk(t *testing.T) {
	inner := modulemocks.NewModuleV2(t)
	inner.EXPECT().Close()

	var receivedBinary []byte
	reloaded := modulemocks.NewModuleV2(t)
	reloaded.EXPECT().Start()
	reloaded.EXPECT().Execute(mock.Anything, mock.Anything, mock.Anything).Return(&sdkpb.ExecutionResult{}, nil)

	factory := func(_ context.Context, _ *host.ModuleConfig, binary []byte, _ ...func(*host.ModuleConfig)) (host.ModuleV2, error) {
		receivedBinary = make([]byte, len(binary))
		copy(receivedBinary, binary)
		return reloaded, nil
	}

	realStore, err := artifacts.NewFileModuleStore(t.TempDir())
	require.NoError(t, err)
	require.NoError(t, realStore.StoreModule("wf-test", "bin-1", []byte("disk-binary")))
	cs := &countingStore{SerialisedModuleStore: realStore}

	// weakBinary is zero-valued (never populated), so L2 is a guaranteed miss.
	em := NewEvictableModule(inner, &host.ModuleConfig{}, cs, "wf-test", factory, nil, int64(len("disk-binary")))
	em.started.Store(true)
	em.Evict()

	_, err = em.Execute(context.Background(), &sdkpb.ExecuteRequest{}, nil)
	require.NoError(t, err)

	assert.Equal(t, int32(1), cs.getModulePathCalls.Load(), "disk should be accessed when weak ref is nil")
	assert.Equal(t, []byte("disk-binary"), receivedBinary)
}

func TestEvictable_WeakRefUpdatedOnReload(t *testing.T) {
	inner := modulemocks.NewModuleV2(t)
	inner.EXPECT().Close()

	reloaded := modulemocks.NewModuleV2(t)
	reloaded.EXPECT().Start().Maybe()
	reloaded.EXPECT().Close().Maybe()
	reloaded.EXPECT().Execute(mock.Anything, mock.Anything, mock.Anything).
		Return(&sdkpb.ExecutionResult{}, nil).Maybe()

	factory := func(_ context.Context, _ *host.ModuleConfig, _ []byte, _ ...func(*host.ModuleConfig)) (host.ModuleV2, error) {
		m := modulemocks.NewModuleV2(t)
		m.EXPECT().Start().Maybe()
		m.EXPECT().Close().Maybe()
		m.EXPECT().Execute(mock.Anything, mock.Anything, mock.Anything).
			Return(&sdkpb.ExecutionResult{}, nil).Maybe()
		return m, nil
	}

	realStore, err := artifacts.NewFileModuleStore(t.TempDir())
	require.NoError(t, err)
	require.NoError(t, realStore.StoreModule("wf-test", "bin-1", []byte("disk-binary")))
	cs := &countingStore{SerialisedModuleStore: realStore}

	em := NewEvictableModule(inner, &host.ModuleConfig{}, cs, "wf-test", factory, nil, int64(len("disk-binary")))
	em.started.Store(true)

	// First cycle: evict + reload from disk (populates weakBinary)
	em.Evict()
	_, err = em.Execute(context.Background(), &sdkpb.ExecuteRequest{}, nil)
	require.NoError(t, err)
	assert.Equal(t, int32(1), cs.getModulePathCalls.Load())

	// Verify weakBinary was populated
	assert.NotNil(t, em.weakBinary.Value(), "weak ref should be set after disk reload")

	// Hold a strong reference to the binary via the weak pointer so GC
	// cannot collect it between the evict and the second reload.
	bp := em.weakBinary.Value()
	require.NotNil(t, bp)

	// Second cycle: evict + reload — should use weak ref (no new disk read)
	em.Evict()
	_, err = em.Execute(context.Background(), &sdkpb.ExecuteRequest{}, nil)
	require.NoError(t, err)
	assert.Equal(t, int32(1), cs.getModulePathCalls.Load(), "second reload should use weak ref, not disk")
	runtime.KeepAlive(bp)
}

// --- Metrics integration tests ---

func TestEvictable_ReloadSourceMetric(t *testing.T) {
	cm, err := NewCacheMetrics()
	require.NoError(t, err)

	inner := modulemocks.NewModuleV2(t)
	inner.EXPECT().Close()

	reloaded := modulemocks.NewModuleV2(t)
	reloaded.EXPECT().Start()
	reloaded.EXPECT().Execute(mock.Anything, mock.Anything, mock.Anything).Return(&sdkpb.ExecutionResult{}, nil)

	factory := func(_ context.Context, _ *host.ModuleConfig, _ []byte, _ ...func(*host.ModuleConfig)) (host.ModuleV2, error) {
		return reloaded, nil
	}

	store, err := artifacts.NewFileModuleStore(t.TempDir())
	require.NoError(t, err)
	require.NoError(t, store.StoreModule("wf-test", "bin-1", []byte("binary")))

	em := NewEvictableModule(inner, &host.ModuleConfig{}, store, "wf-test", factory, cm, int64(len("binary")))
	em.started.Store(true)
	em.Evict()

	// Reload from disk — metrics.recordReload should not panic with non-nil cacheMetrics
	_, err = em.Execute(context.Background(), &sdkpb.ExecuteRequest{}, nil)
	require.NoError(t, err)
}

func TestLRU_EvictionMetric(t *testing.T) {
	cm, err := NewCacheMetrics()
	require.NoError(t, err)

	clock := clockwork.NewFakeClock()
	reapTicker := make(chan time.Time, 1)
	onReaped := make(chan struct{}, 1)

	store, storeErr := artifacts.NewFileModuleStore(t.TempDir())
	require.NoError(t, storeErr)

	em := newLRUModule(t, store, "wf-metric")
	em.lastUsed.Store(clock.Now().UnixNano())

	lru := NewModuleLRU(clock,
		WithIdleTimeout(5*time.Minute),
		WithReapTicker(reapTicker),
		WithOnReaped(onReaped),
		WithCacheMetrics(cm),
	)
	lru.Register("wf-metric", em)
	lru.Start()
	defer lru.Close()

	clock.Advance(6 * time.Minute)
	reapTicker <- clock.Now()
	<-onReaped

	// Eviction happened, and metrics.recordEviction + recordLoaded + recordMemorySaved
	// should not panic with non-nil cacheMetrics.
	assert.False(t, em.IsLoaded())
}

func TestEvictable_BinarySizeTracked(t *testing.T) {
	inner := modulemocks.NewModuleV2(t)
	inner.EXPECT().Close()

	reloaded := modulemocks.NewModuleV2(t)
	reloaded.EXPECT().Start()
	reloaded.EXPECT().Execute(mock.Anything, mock.Anything, mock.Anything).Return(&sdkpb.ExecutionResult{}, nil)

	factory := func(_ context.Context, _ *host.ModuleConfig, _ []byte, _ ...func(*host.ModuleConfig)) (host.ModuleV2, error) {
		return reloaded, nil
	}

	binaryData := make([]byte, 4096)
	store, err := artifacts.NewFileModuleStore(t.TempDir())
	require.NoError(t, err)
	require.NoError(t, store.StoreModule("wf-test", "bin-1", binaryData))

	em := NewEvictableModule(inner, &host.ModuleConfig{}, store, "wf-test", factory, nil, int64(len(binaryData)))
	em.started.Store(true)
	assert.Equal(t, int64(4096), em.BinarySize(), "binary size should match on-disk cache before first reload")

	em.Evict()
	assert.Equal(t, int64(4096), em.BinarySize(), "binary size should remain after eviction before any reload (memorySaved metric)")

	_, err = em.Execute(context.Background(), &sdkpb.ExecuteRequest{}, nil)
	require.NoError(t, err)

	assert.Equal(t, int64(4096), em.BinarySize(), "binary size should match stored binary after reload")
}

func TestLRU_MemorySavedMetric(t *testing.T) {
	cm, err := NewCacheMetrics()
	require.NoError(t, err)

	clock := clockwork.NewFakeClock()
	reapTicker := make(chan time.Time, 1)
	onReaped := make(chan struct{}, 1)

	store, storeErr := artifacts.NewFileModuleStore(t.TempDir())
	require.NoError(t, storeErr)

	m1 := newLRUModule(t, store, "wf-a")
	m1.lastUsed.Store(clock.Now().UnixNano())
	m1.binarySize.Store(1024)

	m2 := newLRUModule(t, store, "wf-b")
	m2.lastUsed.Store(clock.Now().UnixNano())
	m2.binarySize.Store(2048)

	lru := NewModuleLRU(clock,
		WithIdleTimeout(5*time.Minute),
		WithReapTicker(reapTicker),
		WithOnReaped(onReaped),
		WithCacheMetrics(cm),
	)
	lru.Register("wf-a", m1)
	lru.Register("wf-b", m2)
	lru.Start()
	defer lru.Close()

	clock.Advance(6 * time.Minute)
	reapTicker <- clock.Now()
	<-onReaped

	// Both modules evicted; their combined binary size (3072 bytes) is
	// the memory saved. recordMemorySaved should not panic.
	assert.False(t, m1.IsLoaded())
	assert.False(t, m2.IsLoaded())
}
