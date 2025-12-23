package ring_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/workflows/ring"
)

func TestRingStoreIntegration(t *testing.T) {
	t.Run("RingStore can be created and used", func(t *testing.T) {
		store := ring.NewStore()
		require.NotNil(t, store)

		// Test shard health - required for consistent hashing to work
		store.SetShardHealth(1, true)
		health := store.GetShardHealth()
		require.True(t, health[1])

		// Test workflow routing via consistent hashing
		shardID := store.GetShardForWorkflow("workflow1")
		// Should route to shard 1 since it's the only healthy shard
		require.Equal(t, uint32(1), shardID)
	})
}

func TestRingFactoryIntegration(t *testing.T) {
	t.Run("RingFactory can be created", func(t *testing.T) {
		lggr := logger.Test(t)
		store := ring.NewStore()

		factory, err := ring.NewFactory(store, lggr, nil)
		require.NoError(t, err)
		require.NotNil(t, factory)
	})
}
