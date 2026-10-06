package store

import (
	"reflect"
	"testing"
	"time"

	"github.com/jonboulle/clockwork"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/services/servicetest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func TestInMemoryStore_Add(t *testing.T) {
	t.Parallel()
	store := NewInMemoryStore(logger.TestLogger(t), clockwork.NewFakeClock())

	execution, err := store.Add(t.Context(), map[string]*WorkflowExecutionStep{
		"step-1": {Ref: "step-1"},
	}, "test-id", "w1", StatusStarted)
	require.NoError(t, err)
	assert.NotZero(t, execution.CreatedAt)
	assert.NotZero(t, execution.UpdatedAt)
	assert.Equal(t, "test-id", execution.ExecutionID)
	assert.Equal(t, "w1", execution.WorkflowID)
	assert.Equal(t, StatusStarted, execution.Status)
	assert.Len(t, execution.Steps, 1)
	assert.Equal(t, "step-1", execution.Steps["step-1"].Ref)

	// Try adding the same execution ID again
	_, err = store.Add(t.Context(), map[string]*WorkflowExecutionStep{}, "test-id", "", "")
	assert.Error(t, err)
}

func TestInMemoryStore_UpsertStep(t *testing.T) {
	t.Parallel()
	fakeClock := clockwork.NewFakeClock()
	store := NewInMemoryStore(logger.TestLogger(t), fakeClock)

	initialState, err := store.Add(t.Context(), map[string]*WorkflowExecutionStep{}, "test-id", "w1", StatusStarted)
	require.NoError(t, err)

	previousUpdatedAt := initialState.UpdatedAt
	fakeClock.Advance(1 * time.Hour)

	step := &WorkflowExecutionStep{ExecutionID: "test-id", Ref: "step-1"}
	updatedState, err := store.UpsertStep(t.Context(), step)
	require.NoError(t, err)
	assert.Equal(t, step, updatedState.Steps["step-1"])

	assert.True(t, updatedState.UpdatedAt.Equal(previousUpdatedAt.Add(1*time.Hour)) ||
		updatedState.UpdatedAt.After(previousUpdatedAt.Add(1*time.Hour)))
}

func TestInMemoryStore_Get(t *testing.T) {
	t.Parallel()
	store := NewInMemoryStore(logger.TestLogger(t), clockwork.NewFakeClock())
	_, err := store.Add(t.Context(), map[string]*WorkflowExecutionStep{}, "test-id", "w1", StatusStarted)
	require.NoError(t, err)

	retrievedState, err := store.Get(t.Context(), "test-id")
	require.NoError(t, err)
	assert.Equal(t, "test-id", retrievedState.ExecutionID)
	assert.Equal(t, "w1", retrievedState.WorkflowID)
	assert.Equal(t, StatusStarted, retrievedState.Status)
}

func TestInMemoryStore_FinishedExecution(t *testing.T) {
	t.Parallel()
	store := NewInMemoryStoreWithPruneConfiguration(logger.TestLogger(t), clockwork.NewRealClock(),
		10*time.Millisecond, 10*time.Millisecond, 1*time.Hour)
	servicetest.Run(t, store)

	_, err := store.Add(t.Context(), map[string]*WorkflowExecutionStep{
		"step-1": {Ref: "step-1"},
	}, "test-id", "w1", StatusStarted)
	require.NoError(t, err)

	updatedState, err := store.FinishExecution(t.Context(), "test-id", "completed")
	require.NoError(t, err)

	assert.Equal(t, "completed", updatedState.Status)

	// Assert eventually that the execution is no longer in the store
	require.Eventually(t, func() bool {
		_, err := store.Get(t.Context(), "test-id")
		return err != nil
	}, 10*time.Second, 10*time.Millisecond)
}

// TestInMemoryStore_RetainsCompletedExecutionsUntilRetentionCutoff verifies a completed
// execution survives sweeps until it's older than completedExecutionRetention, rather than
// being pruned on the very next sweep after finishing. This is what gives Add's dedup check a
// window to catch a duplicate execution ID arriving shortly after the original completed (e.g.
// a duplicate CRON trigger event, CRE-5715).
func TestInMemoryStore_RetainsCompletedExecutionsUntilRetentionCutoff(t *testing.T) {
	t.Parallel()

	completedExecutionRetention := 20 * time.Minute
	fakeClock := clockwork.NewFakeClock()
	store := NewInMemoryStoreWithPruneConfiguration(logger.TestLogger(t), fakeClock,
		time.Millisecond, completedExecutionRetention, time.Hour)
	servicetest.Run(t, store)

	_, err := store.Add(t.Context(), nil, "test-id", "w1", StatusStarted)
	require.NoError(t, err)
	_, err = store.FinishExecution(t.Context(), "test-id", StatusCompleted)
	require.NoError(t, err)

	// Sweeps run every millisecond on the real clock underneath the ticker, but the sweep
	// itself compares against the fake clock's time, so the entry must survive many sweeps
	// until the fake clock is advanced past the retention cutoff.
	fakeClock.Advance(completedExecutionRetention - time.Second)
	require.Never(t, func() bool {
		_, err := store.Get(t.Context(), "test-id")
		return err != nil
	}, 100*time.Millisecond, 10*time.Millisecond)

	fakeClock.Advance(2 * time.Second)
	require.Eventually(t, func() bool {
		_, err := store.Get(t.Context(), "test-id")
		return err != nil
	}, 10*time.Second, 10*time.Millisecond)
}

func TestInMemoryStore_DeleteByWorkflowID(t *testing.T) {
	t.Parallel()
	s := NewInMemoryStore(logger.TestLogger(t), clockwork.NewFakeClock())

	_, err := s.Add(t.Context(), nil, "exec-1", "wf-A", StatusStarted)
	require.NoError(t, err)
	_, err = s.Add(t.Context(), nil, "exec-2", "wf-A", StatusStarted)
	require.NoError(t, err)
	_, err = s.Add(t.Context(), nil, "exec-3", "wf-B", StatusStarted)
	require.NoError(t, err)

	require.NoError(t, s.DeleteByWorkflowID(t.Context(), "wf-A"))

	_, err = s.Get(t.Context(), "exec-1")
	require.Error(t, err)
	_, err = s.Get(t.Context(), "exec-2")
	require.Error(t, err)

	got, err := s.Get(t.Context(), "exec-3")
	require.NoError(t, err)
	assert.Equal(t, "wf-B", got.WorkflowID)
}

func TestInMemoryStore_ExpiresNonCompletedExecutions(t *testing.T) {
	if testing.Short() {
		t.Skip("too slow for testing.Short")
	}

	t.Parallel()
	expirationDuration := 50 * time.Millisecond

	// completedExecutionRetention is irrelevant here: the execution is never finished.
	store := NewInMemoryStoreWithPruneConfiguration(logger.TestLogger(t), clockwork.NewRealClock(),
		10*time.Millisecond, time.Hour, expirationDuration)

	servicetest.Run(t, store)

	_, err := store.Add(t.Context(), map[string]*WorkflowExecutionStep{
		"step-1": {Ref: "step-1"},
	}, "test-id", "w1", StatusStarted)
	require.NoError(t, err)

	// Expect the state to be removed from the store after the expiration duration
	require.Eventually(t, func() bool {
		_, err2 := store.Get(t.Context(), "test-id")
		return err2 != nil
	}, 10*time.Second, 50*time.Millisecond)

	// Now repeat the test but with a longer expiration duration and check that the state is not expired
	store = NewInMemoryStoreWithPruneConfiguration(logger.TestLogger(t), clockwork.NewRealClock(),
		10*time.Millisecond, time.Hour, 30*time.Second)

	_, err = store.Add(t.Context(), map[string]*WorkflowExecutionStep{
		"step-1": {Ref: "step-1"},
	}, "test-id", "w1", StatusStarted)
	require.NoError(t, err)

	require.Never(t, func() bool {
		_, err2 := store.Get(t.Context(), "test-id")
		return err2 != nil
	}, 300*time.Millisecond, 50*time.Millisecond)
}

// TestInMemoryStore_PruneRebuildsMapOnlyWhenPruning verifies that pruning drops completed and
// expired executions while retaining active ones, and that the backing map is rebuilt only when
// something is actually pruned. Rebuilding is what lets stranded bucket storage become
// GC-eligible (Go maps never shrink after deletes); skipping it on a no-op prune avoids a
// needless full copy of the map every interval.
func TestInMemoryStore_PruneRebuildsMapOnlyWhenPruning(t *testing.T) {
	t.Parallel()

	completedExecutionRetention := time.Minute
	fakeClock := clockwork.NewFakeClock()
	store := NewInMemoryStoreWithPruneConfiguration(logger.TestLogger(t), fakeClock,
		defaultPruneInterval, completedExecutionRetention, 24*time.Hour)

	// active stays; done is completed and should be pruned once past completedExecutionRetention.
	_, err := store.Add(t.Context(), nil, "active", "wf-1", StatusStarted)
	require.NoError(t, err)
	_, err = store.Add(t.Context(), nil, "done", "wf-1", StatusStarted)
	require.NoError(t, err)
	_, err = store.FinishExecution(t.Context(), "done", StatusCompleted)
	require.NoError(t, err)
	fakeClock.Advance(completedExecutionRetention + time.Second)

	mapBefore := reflect.ValueOf(store.idToExecution).Pointer()
	store.pruneCompletedAndExpiredExecutions()
	mapAfterPrune := reflect.ValueOf(store.idToExecution).Pointer()

	_, err = store.Get(t.Context(), "done")
	require.Error(t, err)
	got, err := store.Get(t.Context(), "active")
	require.NoError(t, err)
	assert.Equal(t, "active", got.ExecutionID)
	assert.NotEqual(t, mapBefore, mapAfterPrune, "map should be rebuilt when entries are pruned")

	// A prune that removes nothing must leave the same backing map in place.
	store.pruneCompletedAndExpiredExecutions()
	assert.Equal(t, mapAfterPrune, reflect.ValueOf(store.idToExecution).Pointer(),
		"map should not be rebuilt when nothing is pruned")
}
