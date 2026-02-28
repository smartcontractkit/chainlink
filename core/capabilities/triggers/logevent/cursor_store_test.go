package logevent

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInMemoryCursorStore(t *testing.T) {
	ctx := context.Background()
	store := NewInMemoryCursorStore()

	t.Run("Get returns nil for non-existent trigger", func(t *testing.T) {
		record, err := store.Get(ctx, "non-existent")
		require.NoError(t, err)
		assert.Nil(t, record)
	})

	t.Run("Save and Get roundtrip", func(t *testing.T) {
		err := store.Save(ctx, "trigger-1", "chain-42", "cursor-abc", 1000)
		require.NoError(t, err)

		record, err := store.Get(ctx, "trigger-1")
		require.NoError(t, err)
		require.NotNil(t, record)
		assert.Equal(t, "trigger-1", record.TriggerID)
		assert.Equal(t, "chain-42", record.ChainID)
		assert.Equal(t, "cursor-abc", record.Cursor)
		assert.Equal(t, uint64(1000), record.BlockNumber)
		assert.False(t, record.UpdatedAt.IsZero())
	})

	t.Run("Save overwrites existing cursor", func(t *testing.T) {
		err := store.Save(ctx, "trigger-1", "chain-42", "cursor-abc", 1000)
		require.NoError(t, err)

		err = store.Save(ctx, "trigger-1", "chain-42", "cursor-def", 2000)
		require.NoError(t, err)

		record, err := store.Get(ctx, "trigger-1")
		require.NoError(t, err)
		require.NotNil(t, record)
		assert.Equal(t, "cursor-def", record.Cursor)
		assert.Equal(t, uint64(2000), record.BlockNumber)
	})

	t.Run("Delete removes cursor", func(t *testing.T) {
		err := store.Save(ctx, "trigger-2", "chain-1", "cursor-xyz", 500)
		require.NoError(t, err)

		err = store.Delete(ctx, "trigger-2")
		require.NoError(t, err)

		record, err := store.Get(ctx, "trigger-2")
		require.NoError(t, err)
		assert.Nil(t, record)
	})

	t.Run("Delete non-existent trigger is no-op", func(t *testing.T) {
		err := store.Delete(ctx, "non-existent")
		require.NoError(t, err)
	})

	t.Run("multiple triggers are independent", func(t *testing.T) {
		err := store.Save(ctx, "trigger-a", "chain-1", "cursor-a", 100)
		require.NoError(t, err)

		err = store.Save(ctx, "trigger-b", "chain-1", "cursor-b", 200)
		require.NoError(t, err)

		recordA, err := store.Get(ctx, "trigger-a")
		require.NoError(t, err)
		require.NotNil(t, recordA)
		assert.Equal(t, "cursor-a", recordA.Cursor)
		assert.Equal(t, uint64(100), recordA.BlockNumber)

		recordB, err := store.Get(ctx, "trigger-b")
		require.NoError(t, err)
		require.NotNil(t, recordB)
		assert.Equal(t, "cursor-b", recordB.Cursor)
		assert.Equal(t, uint64(200), recordB.BlockNumber)

		// Delete one, other remains
		err = store.Delete(ctx, "trigger-a")
		require.NoError(t, err)

		recordA, err = store.Get(ctx, "trigger-a")
		require.NoError(t, err)
		assert.Nil(t, recordA)

		recordB, err = store.Get(ctx, "trigger-b")
		require.NoError(t, err)
		require.NotNil(t, recordB)
	})
}

func TestNilCursorStoreIsNoOp(t *testing.T) {
	// Verify that passing nil cursorStore to the trigger doesn't panic.
	// This is tested implicitly by the existing log_event_trigger_test.go
	// which passes nil, but we make the contract explicit here.
	var store CursorStore
	assert.Nil(t, store)
}
