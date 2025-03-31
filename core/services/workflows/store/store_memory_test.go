package store

import (
	"context"
	"testing"
	"time"

	"github.com/jonboulle/clockwork"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/services/servicetest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func TestInMemoryStore_Add(t *testing.T) {
	store := NewInMemoryStore(logger.TestLogger(t), clockwork.NewFakeClock())
	state := &WorkflowExecution{ExecutionID: "test-id"}

	execution, err := store.Add(context.Background(), state)
	require.NoError(t, err)
	assert.NotZero(t, execution.CreatedAt)
	assert.NotZero(t, execution.UpdatedAt)

	// Try adding the same execution ID again
	_, err = store.Add(context.Background(), state)
	assert.Error(t, err)
}

func TestInMemoryStore_UpsertStep(t *testing.T) {
	fakeClock := clockwork.NewFakeClock()
	store := NewInMemoryStore(logger.TestLogger(t), fakeClock)
	state := &WorkflowExecution{ExecutionID: "test-id", Steps: make(map[string]*WorkflowExecutionStep)}
	_, err := store.Add(context.Background(), state)
	require.NoError(t, err)

	previousUpdatedAt := state.UpdatedAt
	fakeClock.Advance(1 * time.Hour)

	step := &WorkflowExecutionStep{ExecutionID: "test-id", Ref: "step-1"}
	updatedState, err := store.UpsertStep(context.Background(), step)
	require.NoError(t, err)
	assert.Equal(t, step, updatedState.Steps["step-1"])

	assert.True(t, updatedState.UpdatedAt.Equal(previousUpdatedAt.Add(1*time.Hour)) ||
		updatedState.UpdatedAt.After(previousUpdatedAt.Add(1*time.Hour)))
}

func TestInMemoryStore_Get(t *testing.T) {
	store := NewInMemoryStore(logger.TestLogger(t), clockwork.NewFakeClock())
	state := &WorkflowExecution{ExecutionID: "test-id"}
	_, err := store.Add(context.Background(), state)
	require.NoError(t, err)

	retrievedState, err := store.Get(context.Background(), "test-id")
	require.NoError(t, err)
	assert.Equal(t, state, &retrievedState)
}

func TestInMemoryStore_FinishedExecution(t *testing.T) {
	store := NewInMemoryStoreWithPruneConfiguration(logger.TestLogger(t), clockwork.NewRealClock(),
		10*time.Millisecond, 1*time.Hour)
	servicetest.Run(t, store)

	state := &WorkflowExecution{ExecutionID: "test-id", Status: "initial"}
	_, err := store.Add(context.Background(), state)
	require.NoError(t, err)

	updatedState, err := store.FinishExecution(context.Background(), "test-id", "completed")
	require.NoError(t, err)

	assert.Equal(t, "completed", updatedState.Status)

	// Assert eventually that the execution is no longer in the store
	require.Eventually(t, func() bool {
		_, err := store.Get(context.Background(), "test-id")
		return err != nil
	}, 10*time.Second, 10*time.Millisecond)
}

func TestInMemoryStore_ExpiresNonCompletedExecutions(t *testing.T) {
	expirationDuration := 50 * time.Millisecond

	store := NewInMemoryStoreWithPruneConfiguration(logger.TestLogger(t), clockwork.NewRealClock(),
		10*time.Millisecond, expirationDuration)

	servicetest.Run(t, store)

	state := &WorkflowExecution{ExecutionID: "test-id"}
	_, err := store.Add(context.Background(), state)
	require.NoError(t, err)

	// Expect the state to be removed from the store after the expiration duration
	require.Eventually(t, func() bool {
		_, err2 := store.Get(context.Background(), "test-id")
		return err2 != nil
	}, 10*time.Second, 50*time.Millisecond)

	// Now repeat the test but with a longer expiration duration and check that the state is not expired
	store = NewInMemoryStoreWithPruneConfiguration(logger.TestLogger(t), clockwork.NewRealClock(),
		10*time.Millisecond, 30*time.Second)

	state = &WorkflowExecution{ExecutionID: "test-id"}
	_, err = store.Add(context.Background(), state)
	require.NoError(t, err)

	require.Never(t, func() bool {
		_, err2 := store.Get(context.Background(), "test-id")
		return err2 != nil
	}, 300*time.Millisecond, 50*time.Millisecond)
}
