package v2

import (
	"context"
	"errors"
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

// forceEvictForTest drops the strong reference and synchronously closes any
// weakly-held holder, simulating the GC reclaiming the compiled module.
// Production code uses Evict (strong-only drop); tests use this to
// deterministically exercise the disk-reload path without waiting for GC.
func (m *EvictableModule) forceEvictForTest() {
	m.mu.Lock()
	defer m.mu.Unlock()
	if h := m.current.Swap(nil); h != nil {
		h.release()
		h.mod.Close()
	}
	if h := m.weakInner.Value(); h != nil {
		h.mod.Close()
	}
	m.weakInner = weak.Pointer[loadedModule]{}
}

// fakeModule is a minimal host.ModuleV2 used by tests that need to observe
// Close calls without going through the mockery expectation lifecycle.
type fakeModule struct {
	closeCalls atomic.Int32
}

func (f *fakeModule) Start()            {}
func (f *fakeModule) IsLegacyDAG() bool { return false }
func (f *fakeModule) Close()            { f.closeCalls.Add(1) }
func (f *fakeModule) Execute(_ context.Context, _ *sdkpb.ExecuteRequest, _ host.ExecutionHelper) (*sdkpb.ExecutionResult, error) {
	return &sdkpb.ExecutionResult{}, nil
}

// countingStore wraps a SerialisedModuleStore and counts GetModule calls
// to verify whether disk I/O occurred during module reload.
type countingStore struct {
	artifacts.SerialisedModuleStore
	getModuleCalls atomic.Int32
}

func (s *countingStore) GetModule(wfID string) (string, string, bool, error) {
	s.getModuleCalls.Add(1)
	return s.SerialisedModuleStore.GetModule(wfID)
}

func newTestEvictableModule(t *testing.T, inner host.ModuleV2, factory ModuleFactoryFn) (*EvictableModule, artifacts.SerialisedModuleStore) {
	t.Helper()
	store, err := artifacts.NewFileModuleStore(t.TempDir(), false)
	require.NoError(t, err)
	require.NoError(t, store.StoreModule("wf-test", []byte("fake-binary"), ""))
	em := NewEvictableModule(inner, &host.ModuleConfig{}, store, "wf-test", "", factory, nil, int64(len("fake-binary")))
	return em, store
}

func TestEvictable_Execute_ContextCanceled(t *testing.T) {
	inner := modulemocks.NewModuleV2(t)
	inner.EXPECT().Start()
	inner.EXPECT().Close()

	em, _ := newTestEvictableModule(t, inner, nil)
	em.Start()
	t.Cleanup(em.Close)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := em.Execute(ctx, &sdkpb.ExecuteRequest{}, nil)
	require.ErrorIs(t, err, context.Canceled)
}

func TestEvictable_Execute_PinRetriesExhausted(t *testing.T) {
	prevAttempts := executePinMaxAttempts
	executePinMaxAttempts = 3
	t.Cleanup(func() { executePinMaxAttempts = prevAttempts })

	prevHook := evictAfterEnsureLoadedHook
	// The hook drops the strong reference each iteration; weak resurrection
	// then re-promotes the same holder until the retry budget is exhausted.
	// The factory is never invoked because L2 keeps hitting.
	evictAfterEnsureLoadedHook = func(em *EvictableModule) { em.Evict() }
	t.Cleanup(func() { evictAfterEnsureLoadedHook = prevHook })

	cm, err := NewCacheMetrics()
	require.NoError(t, err)

	var pinExhaustedRecorded atomic.Int32
	prevMetricHook := cachePinExhaustedHook
	cachePinExhaustedHook = func() { pinExhaustedRecorded.Add(1) }
	t.Cleanup(func() { cachePinExhaustedHook = prevMetricHook })

	inner := modulemocks.NewModuleV2(t)
	inner.EXPECT().Close()

	factory := func(_ context.Context, _ *host.ModuleConfig, _ []byte, _ ...func(*host.ModuleConfig)) (host.ModuleV2, error) {
		t.Fatalf("factory must not be called when weak resurrection keeps hitting")
		return nil, nil
	}

	em, _ := newTestEvictableModule(t, inner, factory)
	em.metrics = cm
	em.started.Store(true)
	t.Cleanup(em.Close)

	_, err = em.Execute(context.Background(), &sdkpb.ExecuteRequest{}, nil)
	require.ErrorIs(t, err, ErrExecutePinExhausted)
	assert.Equal(t, int32(1), pinExhaustedRecorded.Load())
}

func TestEvictable_DelegatesToInner(t *testing.T) {
	inner := modulemocks.NewModuleV2(t)
	inner.EXPECT().Start()
	inner.EXPECT().Execute(mock.Anything, mock.Anything, mock.Anything).Return(&sdkpb.ExecutionResult{}, nil)
	inner.EXPECT().Close()

	em, _ := newTestEvictableModule(t, inner, nil)
	em.Start()
	t.Cleanup(em.Close)

	result, err := em.Execute(context.Background(), &sdkpb.ExecuteRequest{}, nil)
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.False(t, em.IsLegacyDAG())
}

func TestEvictable_LastUsedUpdated(t *testing.T) {
	inner := modulemocks.NewModuleV2(t)
	inner.EXPECT().Execute(mock.Anything, mock.Anything, mock.Anything).Return(&sdkpb.ExecutionResult{}, nil).Times(2)
	inner.EXPECT().Close()

	em, _ := newTestEvictableModule(t, inner, nil)
	t.Cleanup(em.Close)

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
	// Close comes from forceEvictForTest below (simulating GC cleanup);
	// the production Evict call does not close.
	inner.EXPECT().Close()

	em, _ := newTestEvictableModule(t, inner, nil)
	assert.True(t, em.IsLoaded())

	em.Evict()
	assert.False(t, em.IsLoaded())

	// Simulate GC pressure so the mock sees Close deterministically before
	// the test harness asserts expectations.
	em.forceEvictForTest()
}

func TestEvictable_ReloadFromDisk(t *testing.T) {
	inner := modulemocks.NewModuleV2(t)
	inner.EXPECT().Close()

	var reloadedBinary []byte
	reloaded := modulemocks.NewModuleV2(t)
	reloaded.EXPECT().Start()
	reloaded.EXPECT().Execute(mock.Anything, mock.Anything, mock.Anything).Return(&sdkpb.ExecutionResult{}, nil)
	reloaded.EXPECT().Close()

	factory := func(_ context.Context, _ *host.ModuleConfig, binary []byte, _ ...func(*host.ModuleConfig)) (host.ModuleV2, error) {
		reloadedBinary = make([]byte, len(binary))
		copy(reloadedBinary, binary)
		return reloaded, nil
	}

	em, _ := newTestEvictableModule(t, inner, factory)
	em.started.Store(true)
	// Force a full evict (including L2) so the subsequent Execute must go
	// all the way to disk and invoke the factory.
	em.forceEvictForTest()
	t.Cleanup(em.Close)

	_, err := em.Execute(context.Background(), &sdkpb.ExecuteRequest{}, nil)
	require.NoError(t, err)
	assert.True(t, em.IsLoaded())
	assert.Equal(t, []byte("fake-binary"), reloadedBinary)
}

func TestEvictable_ReloadFromDisk_RejectsEngineVersionMismatch(t *testing.T) {
	inner := modulemocks.NewModuleV2(t)
	inner.EXPECT().Close()

	factoryCalls := 0
	factory := func(_ context.Context, _ *host.ModuleConfig, _ []byte, _ ...func(*host.ModuleConfig)) (host.ModuleV2, error) {
		factoryCalls++
		return nil, errors.New("factory should not be called on version mismatch")
	}

	store, err := artifacts.NewFileModuleStore(t.TempDir(), false)
	require.NoError(t, err)
	require.NoError(t, store.StoreModule("wf-test", []byte("v1-binary"), "v1"))

	cm, err := NewCacheMetrics()
	require.NoError(t, err)

	em := NewEvictableModule(inner, &host.ModuleConfig{}, store, "wf-test", "v2", factory, cm, int64(len("v1-binary")))
	em.started.Store(true)
	em.forceEvictForTest()

	_, err = em.Execute(context.Background(), &sdkpb.ExecuteRequest{}, nil)
	require.ErrorIs(t, err, ErrEngineVersionMismatch)
	assert.Equal(t, 0, factoryCalls, "factory must not be invoked on version mismatch")

	_, _, ok, err := store.GetModule("wf-test")
	require.NoError(t, err)
	assert.False(t, ok, "stale cached binary must be deleted after mismatch")
}

func TestEvictable_ReloadFromDisk_AcceptsMatchingEngineVersion(t *testing.T) {
	inner := modulemocks.NewModuleV2(t)
	inner.EXPECT().Close()

	reloaded := modulemocks.NewModuleV2(t)
	reloaded.EXPECT().Start()
	reloaded.EXPECT().Execute(mock.Anything, mock.Anything, mock.Anything).Return(&sdkpb.ExecutionResult{}, nil)
	reloaded.EXPECT().Close()

	factory := func(_ context.Context, _ *host.ModuleConfig, _ []byte, _ ...func(*host.ModuleConfig)) (host.ModuleV2, error) {
		return reloaded, nil
	}

	store, err := artifacts.NewFileModuleStore(t.TempDir(), false)
	require.NoError(t, err)
	require.NoError(t, store.StoreModule("wf-test", []byte("v2-binary"), "v2"))

	em := NewEvictableModule(inner, &host.ModuleConfig{}, store, "wf-test", "v2", factory, nil, int64(len("v2-binary")))
	em.started.Store(true)
	em.forceEvictForTest()
	t.Cleanup(em.Close)

	_, err = em.Execute(context.Background(), &sdkpb.ExecuteRequest{}, nil)
	require.NoError(t, err)
}

func TestEvictable_ReloadCallsStart(t *testing.T) {
	inner := modulemocks.NewModuleV2(t)
	inner.EXPECT().Start()
	inner.EXPECT().Close()

	reloaded := modulemocks.NewModuleV2(t)
	reloaded.EXPECT().Start()
	reloaded.EXPECT().Execute(mock.Anything, mock.Anything, mock.Anything).Return(&sdkpb.ExecutionResult{}, nil)
	reloaded.EXPECT().Close()

	factory := func(_ context.Context, _ *host.ModuleConfig, _ []byte, _ ...func(*host.ModuleConfig)) (host.ModuleV2, error) {
		return reloaded, nil
	}

	em, _ := newTestEvictableModule(t, inner, factory)
	em.Start()
	em.forceEvictForTest()
	t.Cleanup(em.Close)

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
	t.Cleanup(em.Close)

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

func TestEvictable_EvictDoesNotWaitForExecution(t *testing.T) {
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
	t.Cleanup(em.Close)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		_, _ = em.Execute(context.Background(), &sdkpb.ExecuteRequest{}, nil)
	}()

	// Give Execute goroutine time to pin the inner module.
	time.Sleep(10 * time.Millisecond)
	require.True(t, executing.Load(), "Execute should still be running")

	evictReturned := make(chan struct{})
	go func() {
		em.Evict()
		close(evictReturned)
	}()

	select {
	case <-evictReturned:
	case <-time.After(20 * time.Millisecond):
		t.Fatal("Evict must not block on in-flight Execute")
	}

	assert.True(t, executing.Load(), "Execute is still running after Evict returned")
	assert.False(t, em.IsLoaded(), "current entry must be cleared synchronously by Evict")

	wg.Wait()
	assert.False(t, executing.Load())
}

func TestEvictable_NewExecuteProceedsWhileEvictPendingOnLongExecution(t *testing.T) {
	firstExecuteStarted := make(chan struct{})
	releaseFirstExecute := make(chan struct{})
	firstExecuteDone := make(chan error, 1)
	secondExecuteStarted := make(chan struct{})
	secondExecuteDone := make(chan error, 1)
	evictReturned := make(chan struct{})

	var callCount atomic.Int32
	inner := modulemocks.NewModuleV2(t)
	inner.EXPECT().Execute(mock.Anything, mock.Anything, mock.Anything).RunAndReturn(
		func(_ context.Context, _ *sdkpb.ExecuteRequest, _ host.ExecutionHelper) (*sdkpb.ExecutionResult, error) {
			switch callCount.Add(1) {
			case 1:
				close(firstExecuteStarted)
				<-releaseFirstExecute
			case 2:
				close(secondExecuteStarted)
			default:
				t.Fatalf("unexpected execute call count: %d", callCount.Load())
			}
			return &sdkpb.ExecutionResult{}, nil
		},
	)
	inner.EXPECT().Close()

	em, _ := newTestEvictableModule(t, inner, nil)
	t.Cleanup(em.Close)

	go func() {
		_, err := em.Execute(context.Background(), &sdkpb.ExecuteRequest{}, nil)
		firstExecuteDone <- err
	}()

	<-firstExecuteStarted

	evictStart := time.Now()
	go func() {
		em.Evict()
		close(evictReturned)
	}()

	select {
	case <-evictReturned:
	case <-time.After(50 * time.Millisecond):
		t.Fatal("Evict must return without waiting for the long execution")
	}
	require.Less(t, time.Since(evictStart), 50*time.Millisecond)

	go func() {
		_, err := em.Execute(context.Background(), &sdkpb.ExecuteRequest{}, nil)
		secondExecuteDone <- err
	}()

	select {
	case <-secondExecuteStarted:
	case <-time.After(time.Second):
		t.Fatal("second Execute should resurrect via L2 and run while first Execute is still pinned")
	}

	require.NoError(t, <-secondExecuteDone)

	select {
	case <-firstExecuteDone:
		t.Fatal("first Execute returned before being released")
	default:
	}

	close(releaseFirstExecute)
	require.NoError(t, <-firstExecuteDone)
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

	// Each iteration force-evicts (including L2) so the factory is guaranteed
	// to run. Without the force, weak resurrection would skip the factory after
	// the first cycle.
	for i := 0; i < 3; i++ {
		em.forceEvictForTest()
		assert.False(t, em.IsLoaded())

		_, err := em.Execute(context.Background(), &sdkpb.ExecuteRequest{}, nil)
		require.NoError(t, err)
		assert.True(t, em.IsLoaded())
	}

	assert.Equal(t, int32(3), createCount.Load())

	em.Close()
}

func TestEvictable_ReloadFailure(t *testing.T) {
	inner := modulemocks.NewModuleV2(t)
	inner.EXPECT().Close()

	em, store := newTestEvictableModule(t, inner, nil)
	em.forceEvictForTest()

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
	inner.EXPECT().Execute(mock.Anything, mock.Anything, mock.Anything).Return(&sdkpb.ExecutionResult{}, nil).Maybe()
	inner.EXPECT().Close().Maybe()
	require.NoError(t, store.StoreModule(wfID, []byte("binary"), ""))
	factory := func(_ context.Context, _ *host.ModuleConfig, _ []byte, _ ...func(*host.ModuleConfig)) (host.ModuleV2, error) {
		m := modulemocks.NewModuleV2(t)
		m.EXPECT().Start().Maybe()
		m.EXPECT().Close().Maybe()
		m.EXPECT().Execute(mock.Anything, mock.Anything, mock.Anything).Return(&sdkpb.ExecutionResult{}, nil).Maybe()
		return m, nil
	}
	em := NewEvictableModule(inner, &host.ModuleConfig{}, store, wfID, "", factory, nil, int64(len("binary")))
	em.started.Store(true)
	t.Cleanup(em.forceEvictForTest)
	return em
}

func TestLRU_EvictsIdleModule(t *testing.T) {
	clock := clockwork.NewFakeClock()
	reapTicker := make(chan time.Time, 1)
	onReaped := make(chan struct{}, 1)

	store, err := artifacts.NewFileModuleStore(t.TempDir(), false)
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

	store, err := artifacts.NewFileModuleStore(t.TempDir(), false)
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

	store, err := artifacts.NewFileModuleStore(t.TempDir(), false)
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

	store, err := artifacts.NewFileModuleStore(t.TempDir(), false)
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

	store, err := artifacts.NewFileModuleStore(t.TempDir(), false)
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

	store, err := artifacts.NewFileModuleStore(t.TempDir(), false)
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

	store, err := artifacts.NewFileModuleStore(t.TempDir(), false)
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

// TestEvictable_WeakRefHitAfterEvict verifies that Evict drops only the strong
// reference and a subsequent Execute resurrects the still-live compiled module
// via the weak L2, skipping both disk I/O and the factory.
func TestEvictable_WeakRefHitAfterEvict(t *testing.T) {
	inner := modulemocks.NewModuleV2(t)
	inner.EXPECT().Execute(mock.Anything, mock.Anything, mock.Anything).Return(&sdkpb.ExecutionResult{}, nil)
	inner.EXPECT().Close()

	factory := func(_ context.Context, _ *host.ModuleConfig, _ []byte, _ ...func(*host.ModuleConfig)) (host.ModuleV2, error) {
		t.Fatalf("factory must not be called when weak resurrection succeeds")
		return nil, nil
	}

	realStore, err := artifacts.NewFileModuleStore(t.TempDir(), false)
	require.NoError(t, err)
	require.NoError(t, realStore.StoreModule("wf-test", []byte("disk-binary"), ""))
	cs := &countingStore{SerialisedModuleStore: realStore}

	em := NewEvictableModule(inner, &host.ModuleConfig{}, cs, "wf-test", "", factory, nil, int64(len("disk-binary")))
	em.started.Store(true)
	t.Cleanup(em.Close)

	em.Evict()

	_, err = em.Execute(context.Background(), &sdkpb.ExecuteRequest{}, nil)
	require.NoError(t, err)

	assert.Equal(t, int32(0), cs.getModuleCalls.Load(), "disk should not be accessed when weak module is alive")
}

// TestEvictable_WeakRefMissFallsToDisk verifies that when the weak L2 is
// unreachable (GC has reclaimed the holder, simulated via forceEvictForTest),
// ensureLoaded falls through to disk and invokes the factory.
func TestEvictable_WeakRefMissFallsToDisk(t *testing.T) {
	inner := modulemocks.NewModuleV2(t)
	inner.EXPECT().Close()

	var receivedBinary []byte
	reloaded := modulemocks.NewModuleV2(t)
	reloaded.EXPECT().Start()
	reloaded.EXPECT().Execute(mock.Anything, mock.Anything, mock.Anything).Return(&sdkpb.ExecutionResult{}, nil)
	reloaded.EXPECT().Close()

	factory := func(_ context.Context, _ *host.ModuleConfig, binary []byte, _ ...func(*host.ModuleConfig)) (host.ModuleV2, error) {
		receivedBinary = make([]byte, len(binary))
		copy(receivedBinary, binary)
		return reloaded, nil
	}

	realStore, err := artifacts.NewFileModuleStore(t.TempDir(), false)
	require.NoError(t, err)
	require.NoError(t, realStore.StoreModule("wf-test", []byte("disk-binary"), ""))
	cs := &countingStore{SerialisedModuleStore: realStore}

	em := NewEvictableModule(inner, &host.ModuleConfig{}, cs, "wf-test", "", factory, nil, int64(len("disk-binary")))
	em.started.Store(true)
	em.forceEvictForTest()
	t.Cleanup(em.Close)

	_, err = em.Execute(context.Background(), &sdkpb.ExecuteRequest{}, nil)
	require.NoError(t, err)

	assert.Equal(t, int32(1), cs.getModuleCalls.Load(), "disk should be accessed when weak module is dead")
	assert.Equal(t, []byte("disk-binary"), receivedBinary)
}

// TestEvictable_WeakRefPopulatedAfterReload verifies that a disk reload
// populates weakInner, so a second evict+execute cycle hits the weak L2.
func TestEvictable_WeakRefPopulatedAfterReload(t *testing.T) {
	inner := modulemocks.NewModuleV2(t)
	inner.EXPECT().Close()

	var factoryCalls atomic.Int32
	reloaded := modulemocks.NewModuleV2(t)
	reloaded.EXPECT().Start()
	reloaded.EXPECT().Execute(mock.Anything, mock.Anything, mock.Anything).
		Return(&sdkpb.ExecutionResult{}, nil).Times(2)
	reloaded.EXPECT().Close()

	factory := func(_ context.Context, _ *host.ModuleConfig, _ []byte, _ ...func(*host.ModuleConfig)) (host.ModuleV2, error) {
		factoryCalls.Add(1)
		return reloaded, nil
	}

	realStore, err := artifacts.NewFileModuleStore(t.TempDir(), false)
	require.NoError(t, err)
	require.NoError(t, realStore.StoreModule("wf-test", []byte("disk-binary"), ""))
	cs := &countingStore{SerialisedModuleStore: realStore}

	em := NewEvictableModule(inner, &host.ModuleConfig{}, cs, "wf-test", "", factory, nil, int64(len("disk-binary")))
	em.started.Store(true)
	// Drop initial inner entirely so the first reload must go to disk.
	em.forceEvictForTest()
	t.Cleanup(em.Close)

	// First cycle: disk reload populates weakInner.
	_, err = em.Execute(context.Background(), &sdkpb.ExecuteRequest{}, nil)
	require.NoError(t, err)
	assert.Equal(t, int32(1), factoryCalls.Load())
	assert.Equal(t, int32(1), cs.getModuleCalls.Load())
	require.NotNil(t, em.weakInner.Value(), "weak L2 must be populated after disk reload")
	// Hold a strong reference to the weak holder so GC cannot reclaim it
	// between Evict and the second Execute.
	holder := em.weakInner.Value()
	require.NotNil(t, holder)

	// Second cycle: strong-drop only; weak holder stays alive, so L2 hits.
	em.Evict()
	_, err = em.Execute(context.Background(), &sdkpb.ExecuteRequest{}, nil)
	require.NoError(t, err)
	assert.Equal(t, int32(1), factoryCalls.Load(), "second reload must use weak L2, not factory")
	assert.Equal(t, int32(1), cs.getModuleCalls.Load(), "second reload must not touch disk")
	runtime.KeepAlive(holder)
}

// TestEvictable_WeakRefClearedOnForceEvict proves that forceEvictForTest (the
// GC-pressure simulation) genuinely clears the weak pointer — a sanity check
// for the other weak-ref tests.
func TestEvictable_WeakRefClearedOnForceEvict(t *testing.T) {
	inner := modulemocks.NewModuleV2(t)
	inner.EXPECT().Close()

	em, _ := newTestEvictableModule(t, inner, nil)
	require.NotNil(t, em.weakInner.Value(), "weak L2 must be populated at construction")

	em.forceEvictForTest()
	assert.Nil(t, em.weakInner.Value(), "weak L2 must be cleared after forceEvictForTest")

	// runtime.KeepAlive prevents the compiler from reordering the inner
	// reference above the forceEvict call, which would defeat the check.
	runtime.KeepAlive(inner)
}

// TestEvictable_GCFiresCloseAfterEvict proves the production close path:
// after Evict drops the strong reference, the weakly-held loadedModule becomes
// GC-eligible and runtime.AddCleanup must eventually invoke mod.Close. This
// is the only path that reclaims wasm runtime resources in production
// (forceEvictForTest exists solely as a deterministic test hook).
func TestEvictable_GCFiresCloseAfterEvict(t *testing.T) {
	fake := &fakeModule{}

	em, _ := newTestEvictableModule(t, fake, nil)
	require.True(t, em.IsLoaded())
	require.Equal(t, int32(0), fake.closeCalls.Load(), "Close must not fire on construction")

	em.Evict()
	require.False(t, em.IsLoaded(), "Evict clears the strong reference synchronously")
	require.Equal(t, int32(0), fake.closeCalls.Load(),
		"Evict must not call Close; close is deferred to GC + runtime.AddCleanup")

	// runtime.AddCleanup fires asynchronously after the holder becomes
	// unreachable. Force GC repeatedly until the cleanup runs (or time out).
	require.Eventually(t, func() bool {
		runtime.GC()
		return fake.closeCalls.Load() == 1
	}, 5*time.Second, 50*time.Millisecond,
		"runtime.AddCleanup must close the wrapped module after GC reclaims the holder")

	assert.Nil(t, em.weakInner.Value(),
		"weak L2 must report nil once the holder has been GC-reclaimed")
	assert.Equal(t, int32(1), fake.closeCalls.Load(), "Close must fire exactly once")
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
	reloaded.EXPECT().Close()

	factory := func(_ context.Context, _ *host.ModuleConfig, _ []byte, _ ...func(*host.ModuleConfig)) (host.ModuleV2, error) {
		return reloaded, nil
	}

	store, err := artifacts.NewFileModuleStore(t.TempDir(), false)
	require.NoError(t, err)
	require.NoError(t, store.StoreModule("wf-test", []byte("binary"), ""))

	em := NewEvictableModule(inner, &host.ModuleConfig{}, store, "wf-test", "", factory, cm, int64(len("binary")))
	em.started.Store(true)
	em.forceEvictForTest()
	t.Cleanup(em.Close)

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

	store, storeErr := artifacts.NewFileModuleStore(t.TempDir(), false)
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
	reloaded.EXPECT().Close()

	factory := func(_ context.Context, _ *host.ModuleConfig, _ []byte, _ ...func(*host.ModuleConfig)) (host.ModuleV2, error) {
		return reloaded, nil
	}

	binaryData := make([]byte, 4096)
	store, err := artifacts.NewFileModuleStore(t.TempDir(), false)
	require.NoError(t, err)
	require.NoError(t, store.StoreModule("wf-test", binaryData, ""))

	em := NewEvictableModule(inner, &host.ModuleConfig{}, store, "wf-test", "", factory, nil, int64(len(binaryData)))
	em.started.Store(true)
	assert.Equal(t, int64(4096), em.BinarySize(), "binary size should match on-disk cache before first reload")

	em.forceEvictForTest()
	assert.Equal(t, int64(4096), em.BinarySize(), "binary size should remain after eviction before any reload (memorySaved metric)")

	_, err = em.Execute(context.Background(), &sdkpb.ExecuteRequest{}, nil)
	require.NoError(t, err)

	assert.Equal(t, int64(4096), em.BinarySize(), "binary size should match stored binary after reload")
	t.Cleanup(em.Close)
}

func TestLRU_MemorySavedMetric(t *testing.T) {
	cm, err := NewCacheMetrics()
	require.NoError(t, err)

	clock := clockwork.NewFakeClock()
	reapTicker := make(chan time.Time, 1)
	onReaped := make(chan struct{}, 1)

	store, storeErr := artifacts.NewFileModuleStore(t.TempDir(), false)
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

func TestLRU_ReapMemorySavedBytesNotCumulative(t *testing.T) {
	prevHook := reapMemorySavedHook
	var observed []int64
	reapMemorySavedHook = func(b int64) { observed = append(observed, b) }
	t.Cleanup(func() { reapMemorySavedHook = prevHook })

	cm, err := NewCacheMetrics()
	require.NoError(t, err)

	clock := clockwork.NewFakeClock()
	reapTicker := make(chan time.Time, 1)
	onReaped := make(chan struct{}, 1)

	store, storeErr := artifacts.NewFileModuleStore(t.TempDir(), false)
	require.NoError(t, storeErr)

	em := newLRUModule(t, store, "wf-repeat")
	em.lastUsed.Store(clock.Now().UnixNano())
	want := em.BinarySize()

	lru := NewModuleLRU(clock,
		WithIdleTimeout(5*time.Minute),
		WithReapTicker(reapTicker),
		WithOnReaped(onReaped),
		WithCacheMetrics(cm),
	)
	lru.Register("wf-repeat", em)
	lru.Start()
	defer lru.Close()

	clock.Advance(6 * time.Minute)
	reapTicker <- clock.Now()
	<-onReaped
	require.False(t, em.IsLoaded())

	reapTicker <- clock.Now()
	<-onReaped
	require.Equal(t, []int64{want, want}, observed, "idle snapshot must not grow on repeated reap")

	_, err = em.Execute(context.Background(), &sdkpb.ExecuteRequest{}, nil)
	require.NoError(t, err)
	em.lastUsed.Store(clock.Now().UnixNano())

	reapTicker <- clock.Now()
	<-onReaped
	require.True(t, em.IsLoaded())

	em.Evict()
	reapTicker <- clock.Now()
	<-onReaped
	require.Equal(t, []int64{want, want, 0, want}, observed)
}
