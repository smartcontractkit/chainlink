package v2

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
)

// fakeModulePanicStore mimics job.KVStore: Get returns a wrapped sql.ErrNoRows
// for absent keys.
type fakeModulePanicStore struct {
	mu       sync.Mutex
	m        map[string][]byte
	storeErr error
	getErr   error
}

func newFakeModulePanicStore() *fakeModulePanicStore {
	return &fakeModulePanicStore{m: map[string][]byte{}}
}

func (f *fakeModulePanicStore) Store(_ context.Context, key string, val []byte) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.storeErr != nil {
		return f.storeErr
	}
	f.m[key] = append([]byte(nil), val...)
	return nil
}

func (f *fakeModulePanicStore) Get(_ context.Context, key string) ([]byte, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.getErr != nil {
		return nil, f.getErr
	}
	v, ok := f.m[key]
	if !ok {
		return nil, fmt.Errorf("get %q: %w", key, sql.ErrNoRows)
	}
	return v, nil
}

func (f *fakeModulePanicStore) Delete(_ context.Context, key string) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	delete(f.m, key)
	return nil
}

func TestRecordModulePanic_IncrementsAndPersists(t *testing.T) {
	t.Parallel()
	s := newFakeModulePanicStore()

	for want := uint64(1); want <= 3; want++ {
		got, err := recordModulePanic(t.Context(), s, "wf")
		require.NoError(t, err)
		require.Equal(t, want, got)
	}

	count, err := modulePanicCount(t.Context(), s, "wf")
	require.NoError(t, err)
	require.Equal(t, uint64(3), count)
}

func TestModulePanicCount_MissingIsZero(t *testing.T) {
	t.Parallel()
	s := newFakeModulePanicStore()

	count, err := modulePanicCount(t.Context(), s, "never-seen")
	require.NoError(t, err)
	require.Zero(t, count)
}

func TestModulePanicCount_IsPerWorkflow(t *testing.T) {
	t.Parallel()
	s := newFakeModulePanicStore()

	_, err := recordModulePanic(t.Context(), s, "wf-a")
	require.NoError(t, err)

	count, err := modulePanicCount(t.Context(), s, "wf-b")
	require.NoError(t, err)
	require.Zero(t, count)
}

func newGuardTestEngine(t *testing.T, s ModulePanicStore) *Engine {
	return &Engine{
		cfg:  &EngineConfig{PanicStore: s, WorkflowID: "wf"},
		lggr: logger.Sugared(logger.Test(t)),
	}
}

func TestCheckModulePanicQuarantine_ThresholdBoundary(t *testing.T) {
	t.Parallel()
	s := newFakeModulePanicStore()
	e := newGuardTestEngine(t, s)

	for range MaxModulePanics - 1 {
		_, err := recordModulePanic(t.Context(), s, "wf")
		require.NoError(t, err)
		require.NoError(t, e.checkModulePanicQuarantine(t.Context()), "must not quarantine below threshold")
	}

	// The panic that reaches MaxModulePanics quarantines on the next boot.
	_, err := recordModulePanic(t.Context(), s, "wf")
	require.NoError(t, err)
	require.Error(t, e.checkModulePanicQuarantine(t.Context()))
}

func TestClearModulePanic_UnquarantinesWorkflow(t *testing.T) {
	t.Parallel()
	s := newFakeModulePanicStore()
	e := newGuardTestEngine(t, s)

	for range MaxModulePanics {
		_, err := recordModulePanic(t.Context(), s, "wf")
		require.NoError(t, err)
	}
	require.Error(t, e.checkModulePanicQuarantine(t.Context()))

	require.NoError(t, ClearModulePanic(t.Context(), s, "wf"))

	count, err := modulePanicCount(t.Context(), s, "wf")
	require.NoError(t, err)
	require.Zero(t, count)
	require.NoError(t, e.checkModulePanicQuarantine(t.Context()))
}

func TestModulePanicGuard_DisabledWhenStoreNil(t *testing.T) {
	t.Parallel()
	e := newGuardTestEngine(t, nil)

	require.Zero(t, e.recordModulePanic(t.Context()))
	require.NoError(t, e.checkModulePanicQuarantine(t.Context()))
}

func TestCheckModulePanicQuarantine_FailsOpenOnReadError(t *testing.T) {
	t.Parallel()
	s := newFakeModulePanicStore()
	s.getErr = errors.New("db unavailable")
	e := newGuardTestEngine(t, s)

	require.NoError(t, e.checkModulePanicQuarantine(t.Context()),
		"a store read error must not block execution")
}
