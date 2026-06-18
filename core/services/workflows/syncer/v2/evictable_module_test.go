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

func TestEvictable_Execute_TryAcquireExhausted(t *testing.T) {
	prevTryAcquireAttempts := tryAcquireMaxAttempts
	tryAcquireMaxAttempts = 3
	t.Cleanup(func() { tryAcquireMaxAttempts = prevTryAcquireAttempts })

	prevExecuteAttempts := executePinMaxAttempts
	executePinMaxAttempts = 2
	t.Cleanup(func() { executePinMaxAttempts = prevExecuteAttempts })

	prevCAS := tryAcquireCompareAndSwap
	tryAcquireCompareAndSwap = func(_ *loadedModule, _, _ int64) bool { return false }
	t.Cleanup(func() { tryAcquireCompareAndSwap = prevCAS })

	cm, err := NewCacheMetrics()
	require.NoError(t, err)

	var tryAcquireExhaustedRecorded atomic.Int32
	prevTryAcquireHook := cacheTryAcquireExhaustedHook
	cacheTryAcquireExhaustedHook = func() { tryAcquireExhaustedRecorded.Add(1) }
	t.Cleanup(func() { cacheTryAcquireExhaustedHook = prevTryAcquireHook })

	var pinExhaustedRecorded atomic.Int32
	prevPinHook := cachePinExhaustedHook
	cachePinExhaustedHook = func() { pinExhaustedRecorded.Add(1) }
	t.Cleanup(func() { cachePinExhaustedHook = prevPinHook })

	inner := modulemocks.NewModuleV2(t)
	inner.EXPECT().Close()

	em, _ := newTestEvictableModule(t, inner, nil)
	em.metrics = cm
	em.started.Store(true)

	_, err = em.Execute(context.Background(), &sdkpb.ExecuteRequest{}, nil)
	require.ErrorIs(t, err, ErrExecutePinExhausted)
	assert.Equal(t, int32(2), tryAcquireExhaustedRecorded.Load(), "one tryAcquire exhaustion per execute retry attempt")
	assert.Equal(t, int32(1), pinExhaustedRecorded.Load(), "pin exhaustion should still be emitted once on terminal failure")

	em.Close()
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

// Ensure that calling evict once ends all concurrent execution attempts
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
	execErrs := make(chan error, 5)
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, err := em.Execute(context.Background(), &sdkpb.ExecuteRequest{}, nil)
			execErrs <- err
		}()
	}
	wg.Add(1)
	go func() {
		defer wg.Done()
		em.Evict()
	}()
	wg.Wait()
	close(execErrs)
	for err := range execErrs {
		require.NoError(t, err)
	}
}

func TestEvictable_EvictDoesNotWaitForExecution(t *testing.T) {
	var executing atomic.Bool
	var closeCalled atomic.Bool
	executeStarted := make(chan struct{})
	releaseExecute := make(chan struct{})

	inner := modulemocks.NewModuleV2(t)
	inner.EXPECT().Execute(mock.Anything, mock.Anything, mock.Anything).RunAndReturn(
		func(_ context.Context, _ *sdkpb.ExecuteRequest, _ host.ExecutionHelper) (*sdkpb.ExecutionResult, error) {
			executing.Store(true)
			close(executeStarted)
			<-releaseExecute
			executing.Store(false)
			return &sdkpb.ExecutionResult{}, nil
		},
	)
	inner.EXPECT().Close().Run(func() { closeCalled.Store(true) })

	em, _ := newTestEvictableModule(t, inner, nil)
	t.Cleanup(em.Close)

	execDone := make(chan error, 1)
	go func() {
		_, err := em.Execute(context.Background(), &sdkpb.ExecuteRequest{}, nil)
		execDone <- err
	}()

	<-executeStarted
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
	assert.False(t, closeCalled.Load(), "inner.Close must not fire while a pin is held")
	assert.True(t, em.IsLoaded(), "eviction is skipped while a pin is held")

	select {
	case <-execDone:
		t.Fatal("Execute returned before being released")
	default:
	}

	close(releaseExecute)
	require.NoError(t, <-execDone)
	assert.False(t, executing.Load())
	assert.False(t, closeCalled.Load(), "skipped eviction must not close the module")

	em.Close()
	assert.True(t, closeCalled.Load(), "close still releases module ownership")
}

func TestEvictable_NewExecuteUsesExistingModuleWhenEvictSkipped(t *testing.T) {
	firstExecuteStarted := make(chan struct{})
	releaseFirstExecute := make(chan struct{})
	firstExecuteDone := make(chan error, 1)
	secondExecuteDone := make(chan error, 1)
	evictReturned := make(chan struct{})

	var callCount atomic.Int32
	var closeCalled atomic.Bool
	inner := modulemocks.NewModuleV2(t)
	inner.EXPECT().Execute(mock.Anything, mock.Anything, mock.Anything).RunAndReturn(
		func(_ context.Context, _ *sdkpb.ExecuteRequest, _ host.ExecutionHelper) (*sdkpb.ExecutionResult, error) {
			switch callCount.Add(1) {
			case 1:
				close(firstExecuteStarted)
				<-releaseFirstExecute
			case 2:
				// second Execute should run against the same inner module while eviction is skipped
			default:
				t.Fatalf("unexpected execute call count: %d", callCount.Load())
			}
			return &sdkpb.ExecutionResult{}, nil
		},
	)
	inner.EXPECT().Close().Run(func() { closeCalled.Store(true) })

	em, _ := newTestEvictableModule(t, inner, nil)
	em.started.Store(true)

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

	require.NoError(t, <-secondExecuteDone)
	assert.Equal(t, int32(2), callCount.Load(), "both executes should run on the existing module")
	assert.True(t, em.IsLoaded(), "module should remain loaded when eviction is skipped")

	select {
	case <-firstExecuteDone:
		t.Fatal("first Execute returned before being released")
	default:
	}

	close(releaseFirstExecute)
	require.NoError(t, <-firstExecuteDone)
	assert.False(t, closeCalled.Load(), "skipped eviction should not close the module")

	em.Close()
	assert.True(t, closeCalled.Load(), "Close should eventually release module ownership")
}

func TestLRU_FrequentReapSkipsPinnedModuleAndEvictsAfterDrain(t *testing.T) {
	clock := clockwork.NewFakeClock()
	reapTicker := make(chan time.Time, 64)
	onReaped := make(chan struct{}, 64)

	const concurrentExecs = 5

	execStarted := make(chan struct{}, concurrentExecs)
	releaseExec := make(chan struct{})

	var activeExecs atomic.Int32
	var closeCalls atomic.Int32
	inner := modulemocks.NewModuleV2(t)
	inner.EXPECT().Execute(mock.Anything, mock.Anything, mock.Anything).RunAndReturn(
		func(_ context.Context, _ *sdkpb.ExecuteRequest, _ host.ExecutionHelper) (*sdkpb.ExecutionResult, error) {
			activeExecs.Add(1)
			execStarted <- struct{}{}
			<-releaseExec
			activeExecs.Add(-1)
			return &sdkpb.ExecutionResult{}, nil
		},
	).Times(concurrentExecs)
	inner.EXPECT().Close().Run(func() { closeCalls.Add(1) }).Once()

	var factoryCalls atomic.Int32
	factory := func(_ context.Context, _ *host.ModuleConfig, _ []byte, _ ...func(*host.ModuleConfig)) (host.ModuleV2, error) {
		factoryCalls.Add(1)
		return nil, errors.New("unexpected module reload while pinned")
	}

	em, _ := newTestEvictableModule(t, inner, factory)
	lru := NewModuleLRU(clock,
		WithIdleTimeout(time.Nanosecond),
		WithReapTicker(reapTicker),
		WithOnReaped(onReaped),
	)
	lru.Register("wf-test", em)
	lru.Start()
	defer lru.Close()

	var wg sync.WaitGroup
	execErrs := make(chan error, concurrentExecs)
	for i := 0; i < concurrentExecs; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, err := em.Execute(context.Background(), &sdkpb.ExecuteRequest{}, nil)
			execErrs <- err
		}()
	}

	for i := 0; i < concurrentExecs; i++ {
		<-execStarted
	}
	require.Equal(t, int32(concurrentExecs), activeExecs.Load(), "all executes must overlap")

	// Keep forcing eviction while work is pinned; all these attempts should be skipped.
	for i := 0; i < 25; i++ {
		em.lastUsed.Store(clock.Now().Add(-time.Hour).UnixNano())
		clock.Advance(time.Second)
		reapTicker <- clock.Now()
		<-onReaped

		assert.True(t, em.IsLoaded(), "pinned module must remain loaded")
		assert.Equal(t, int32(0), closeCalls.Load(), "Close must not fire while executes are pinned")
	}
	assert.Equal(t, int32(0), factoryCalls.Load(), "pinned eviction skips must not create a new module")

	close(releaseExec)
	wg.Wait()
	close(execErrs)
	for err := range execErrs {
		require.NoError(t, err)
	}
	require.Equal(t, int32(0), activeExecs.Load())

	// Once all pins are gone, next reap should evict immediately.
	em.lastUsed.Store(clock.Now().Add(-time.Hour).UnixNano())
	clock.Advance(time.Second)
	reapTicker <- clock.Now()
	<-onReaped

	assert.False(t, em.IsLoaded(), "module should be evicted after executions drain")
	// Evict drops the strong ref only; mod.Close runs via runtime.AddCleanup after GC
	// reclaims the weak holder (same contract as TestEvictable_GCFiresCloseAfterEvict).
	require.Eventually(t, func() bool {
		runtime.GC()
		return closeCalls.Load() == 1
	}, 5*time.Second, 50*time.Millisecond,
		"exactly one module instance should be closed after GC + AddCleanup")
	assert.Equal(t, int32(0), factoryCalls.Load(), "eviction itself should not reload a module")
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

func newTestLRU(t *testing.T, maxLoaded int) (*clockwork.FakeClock, *ModuleLRU, chan time.Time, chan struct{}) {
	t.Helper()
	clock := clockwork.NewFakeClock()
	reap := make(chan time.Time, 4)
	done := make(chan struct{}, 4)
	lru := NewModuleLRU(clock,
		WithMaxLoadedModules(maxLoaded),
		WithIdleTimeout(time.Hour),
		WithReapTicker(reap),
		WithOnReaped(done),
	)
	lru.Start()
	t.Cleanup(lru.Close)
	return clock, lru, reap, done
}

func triggerLRUReap(t *testing.T, clock *clockwork.FakeClock, reap chan time.Time, done chan struct{}) {
	t.Helper()
	reap <- clock.Now()
	<-done
}

func TestLRU_AtCapacity_noEvictionUntilOver(t *testing.T) {
	t.Parallel()
	clock, lru, reap, done := newTestLRU(t, 2)

	store, err := artifacts.NewFileModuleStore(t.TempDir(), false)
	require.NoError(t, err)

	m1 := newLRUModule(t, store, "wf-1")
	m1.lastUsed.Store(clock.Now().Add(-3 * time.Minute).UnixNano())
	m2 := newLRUModule(t, store, "wf-2")
	m2.lastUsed.Store(clock.Now().Add(-2 * time.Minute).UnixNano())

	lru.Register("wf-1", m1)
	lru.Register("wf-2", m2)
	triggerLRUReap(t, clock, reap, done)

	assert.True(t, m1.IsLoaded(), "at capacity: oldest should not be evicted yet")
	assert.True(t, m2.IsLoaded(), "at capacity: newest should remain loaded")

	m3 := newLRUModule(t, store, "wf-3")
	m3.lastUsed.Store(clock.Now().Add(-1 * time.Minute).UnixNano())
	lru.Register("wf-3", m3)
	triggerLRUReap(t, clock, reap, done)

	assert.False(t, m1.IsLoaded(), "over capacity: least recently used should be evicted")
	assert.True(t, m2.IsLoaded())
	assert.True(t, m3.IsLoaded())
}

func TestLRU_RecencyBump_changesEvictionVictim(t *testing.T) {
	t.Parallel()
	clock, lru, reap, done := newTestLRU(t, 2)

	store, err := artifacts.NewFileModuleStore(t.TempDir(), false)
	require.NoError(t, err)

	mA := newLRUModule(t, store, "wf-A")
	mA.lastUsed.Store(clock.Now().Add(-10 * time.Minute).UnixNano())
	mB := newLRUModule(t, store, "wf-B")
	mB.lastUsed.Store(clock.Now().Add(-1 * time.Minute).UnixNano())

	lru.Register("wf-A", mA)
	lru.Register("wf-B", mB)

	_, err = mA.Execute(context.Background(), &sdkpb.ExecuteRequest{}, nil)
	require.NoError(t, err)

	mC := newLRUModule(t, store, "wf-C")
	mC.lastUsed.Store(clock.Now().UnixNano())
	lru.Register("wf-C", mC)
	triggerLRUReap(t, clock, reap, done)

	assert.True(t, mA.IsLoaded(), "wf-A was touched recently and should survive")
	assert.False(t, mB.IsLoaded(), "wf-B should be evicted after wf-A recency bump")
	assert.True(t, mC.IsLoaded())
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

	lru.mu.Lock()
	defer lru.mu.Unlock()
	assert.Empty(t, lru.modules)
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

func TestLRU_MaxLoaded_zero_disablesCapEnforcement(t *testing.T) {
	t.Parallel()
	clock, lru, reap, done := newTestLRU(t, 0)

	store, err := artifacts.NewFileModuleStore(t.TempDir(), false)
	require.NoError(t, err)

	modules := make([]*EvictableModule, 3)
	for i := 0; i < 3; i++ {
		wfID := string(rune('A' + i))
		modules[i] = newLRUModule(t, store, wfID)
		modules[i].lastUsed.Store(clock.Now().Add(-time.Duration(i+1) * time.Minute).UnixNano())
		lru.Register(wfID, modules[i])
	}
	triggerLRUReap(t, clock, reap, done)

	for i, m := range modules {
		assert.True(t, m.IsLoaded(), "module %d: maxLoaded=0 must not enforce a loaded cap", i)
	}
}

func TestLRU_Register_duplicateWorkflowID_replaces(t *testing.T) {
	t.Parallel()
	_, lru, _, _ := newTestLRU(t, 10)

	store, err := artifacts.NewFileModuleStore(t.TempDir(), false)
	require.NoError(t, err)

	em1 := newLRUModule(t, store, "wf-1")
	em2 := newLRUModule(t, store, "wf-1")
	lru.Register("wf-1", em1)
	lru.Register("wf-1", em2)

	lru.mu.Lock()
	got := lru.modules["wf-1"]
	lru.mu.Unlock()
	require.Equal(t, em2, got)
	assert.True(t, lru.Contains("wf-1"))
}

func TestLRU_ConcurrentReapAndRegister(t *testing.T) {
	t.Parallel()
	clock := clockwork.NewFakeClock()
	reap := make(chan time.Time, 64)
	done := make(chan struct{}, 64)
	lru := NewModuleLRU(clock,
		WithMaxLoadedModules(5),
		WithIdleTimeout(time.Hour),
		WithReapTicker(reap),
		WithOnReaped(done),
	)
	lru.Start()
	defer lru.Close()

	store, err := artifacts.NewFileModuleStore(t.TempDir(), false)
	require.NoError(t, err)

	const workers = 20
	type entry struct {
		wfID string
		em   *EvictableModule
	}
	entries := make([]entry, workers)
	for i := 0; i < workers; i++ {
		wfID := string(rune('A' + i))
		entries[i] = entry{wfID: wfID, em: newLRUModule(t, store, wfID)}
	}

	var wg sync.WaitGroup
	for i := range entries {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			e := entries[idx]
			lru.Register(e.wfID, e.em)
			if idx%3 == 0 {
				reap <- clock.Now()
				<-done
			}
			lru.Contains(e.wfID)
		}(i)
	}
	wg.Wait()

	lru.mu.Lock()
	n := len(lru.modules)
	lru.mu.Unlock()
	assert.Positive(t, n)
}

func TestEvictable_Execute_L1_hit(t *testing.T) {
	t.Parallel()
	inner := modulemocks.NewModuleV2(t)
	inner.EXPECT().Start()
	inner.EXPECT().Execute(mock.Anything, mock.Anything, mock.Anything).Return(&sdkpb.ExecutionResult{}, nil)
	inner.EXPECT().Close()

	store, err := artifacts.NewFileModuleStore(t.TempDir(), false)
	require.NoError(t, err)
	require.NoError(t, store.StoreModule("wf-test", []byte("fake-binary"), ""))
	cs := &countingStore{SerialisedModuleStore: store}

	em := NewEvictableModule(inner, &host.ModuleConfig{}, cs, "wf-test", "", nil, nil, int64(len("fake-binary")))
	em.Start()
	t.Cleanup(em.Close)

	_, err = em.Execute(context.Background(), &sdkpb.ExecuteRequest{}, nil)
	require.NoError(t, err)
	assert.True(t, em.IsLoaded())
	assert.Equal(t, int32(0), cs.getModuleCalls.Load(), "L1 hit must not read from disk")
}

func TestEvictable_Evict_then_reloadWithoutDisk(t *testing.T) {
	t.Parallel()
	inner := modulemocks.NewModuleV2(t)
	inner.EXPECT().Execute(mock.Anything, mock.Anything, mock.Anything).Return(&sdkpb.ExecutionResult{}, nil).Once()
	inner.EXPECT().Close()

	factory := func(_ context.Context, _ *host.ModuleConfig, _ []byte, _ ...func(*host.ModuleConfig)) (host.ModuleV2, error) {
		t.Fatal("factory must not run when weak L2 resurrects after Evict")
		return nil, nil
	}

	store, err := artifacts.NewFileModuleStore(t.TempDir(), false)
	require.NoError(t, err)
	require.NoError(t, store.StoreModule("wf-test", []byte("fake-binary"), ""))
	cs := &countingStore{SerialisedModuleStore: store}

	em := NewEvictableModule(inner, &host.ModuleConfig{}, cs, "wf-test", "", factory, nil, int64(len("fake-binary")))
	em.started.Store(true)
	t.Cleanup(em.Close)

	em.Evict()
	assert.False(t, em.IsLoaded())

	_, err = em.Execute(context.Background(), &sdkpb.ExecuteRequest{}, nil)
	require.NoError(t, err)
	assert.True(t, em.IsLoaded())
	assert.Equal(t, int32(0), cs.getModuleCalls.Load(), "weak L2 reload must not touch disk")
}

func TestEvictable_emptyWorkflowID_diskMiss(t *testing.T) {
	t.Parallel()
	store, err := artifacts.NewFileModuleStore(t.TempDir(), false)
	require.NoError(t, err)

	em := NewEvictableModule(nil, &host.ModuleConfig{}, store, "", "", nil, nil, 0)
	em.started.Store(true)

	_, err = em.Execute(context.Background(), &sdkpb.ExecuteRequest{}, nil)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "no cached binary")
	assert.False(t, em.IsLoaded())
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
	require.Equal(t, []int64{3072}, observed)
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
