package arbiter

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestState_NewState(t *testing.T) {
	state := NewState()

	require.NotNil(t, state)
	assert.Equal(t, 0, state.GetCurrentReplicaCount())
	assert.Equal(t, 0, state.GetConsensusWantShards())
}

func TestState_SetCurrentReplicas(t *testing.T) {
	state := NewState()

	replicas := map[string]ShardReplica{
		"shard-0": {Status: "READY", Message: "Running"},
		"shard-1": {Status: "READY", Message: "Running"},
		"shard-2": {Status: "INSTALLING", Message: "In progress"},
	}

	state.SetCurrentReplicas(replicas)

	assert.Equal(t, 3, state.GetCurrentReplicaCount())
}

func TestState_SetConsensusWantShards(t *testing.T) {
	state := NewState()

	state.SetConsensusWantShards(7)

	assert.Equal(t, 7, state.GetConsensusWantShards())
}

func TestState_GetRoutableShards(t *testing.T) {
	state := NewState()

	replicas := map[string]ShardReplica{
		"shard-0": {Status: StatusReady, Message: "Running"},
		"shard-1": {Status: StatusReady, Message: "Running"},
		"shard-2": {Status: "INSTALLING", Message: "In progress"},
	}

	state.SetCurrentReplicas(replicas)

	routable := state.GetRoutableShards()

	assert.Equal(t, 2, routable.ReadyCount)
	assert.Len(t, routable.ShardInfo, 3)
}

func TestState_GetRoutableShards_Empty(t *testing.T) {
	state := NewState()

	routable := state.GetRoutableShards()

	assert.Equal(t, 0, routable.ReadyCount)
	assert.Empty(t, routable.ShardInfo)
}

func TestState_Concurrency(t *testing.T) {
	state := NewState()
	var wg sync.WaitGroup

	// Run concurrent writes and reads
	for i := 0; i < 100; i++ {
		wg.Add(3)

		// Writer goroutine - SetCurrentReplicas
		go func(i int) {
			defer wg.Done()
			replicas := map[string]ShardReplica{
				"shard-0": {Status: StatusReady, Message: "Running"},
			}
			state.SetCurrentReplicas(replicas)
		}(i)

		// Writer goroutine - SetConsensusWantShards
		go func(i int) {
			defer wg.Done()
			state.SetConsensusWantShards(i)
		}(i)

		// Reader goroutine - GetRoutableShards
		go func() {
			defer wg.Done()
			_ = state.GetRoutableShards()
		}()
	}

	wg.Wait()

	// If we got here without data race, the test passes
	// The actual values don't matter, we're testing thread safety
	routable := state.GetRoutableShards()
	assert.NotNil(t, routable)
}

func TestState_SetCurrentReplicas_Empty(t *testing.T) {
	state := NewState()

	// Start with some replicas
	replicas := map[string]ShardReplica{
		"shard-0": {Status: StatusReady, Message: "Running"},
	}
	state.SetCurrentReplicas(replicas)
	assert.Equal(t, 1, state.GetCurrentReplicaCount())

	// Update with empty replicas
	state.SetCurrentReplicas(map[string]ShardReplica{})
	assert.Equal(t, 0, state.GetCurrentReplicaCount())
}

func TestState_SetCurrentReplicas_WithMetrics(t *testing.T) {
	state := NewState()

	replicas := map[string]ShardReplica{
		"shard-0": {
			Status:  StatusReady,
			Message: "Running",
			Metrics: map[string]float64{
				"cpu_usage":    0.75,
				"memory_usage": 0.60,
			},
		},
	}

	state.SetCurrentReplicas(replicas)

	assert.Equal(t, 1, state.GetCurrentReplicaCount())
	routable := state.GetRoutableShards()
	assert.Equal(t, 1, routable.ReadyCount)
}
