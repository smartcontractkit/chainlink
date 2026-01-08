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
	assert.Equal(t, 1, state.GetDesiredReplicaCount())
	assert.Equal(t, 1, state.GetApprovedReplicaCount())
	assert.Equal(t, 0, state.GetCurrentReplicaCount())

	spec := state.GetScalingSpec()
	assert.Equal(t, "Initial state", spec.LastScalingReason)
}

func TestState_Update(t *testing.T) {
	state := NewState()

	replicas := map[string]ShardReplica{
		"shard-0": {Status: "READY", Message: "Running"},
		"shard-1": {Status: "READY", Message: "Running"},
		"shard-2": {Status: "INSTALLING", Message: "In progress"},
	}

	state.Update(replicas, 5, "high CPU utilization")

	assert.Equal(t, 3, state.GetCurrentReplicaCount())
	assert.Equal(t, 5, state.GetDesiredReplicaCount())

	spec := state.GetScalingSpec()
	assert.Equal(t, "high CPU utilization", spec.LastScalingReason)
	assert.Equal(t, 3, spec.CurrentReplicaCount)
	assert.Equal(t, 5, spec.DesiredReplicaCount)
}

func TestState_SetApprovedCount(t *testing.T) {
	state := NewState()

	state.SetApprovedCount(7)

	assert.Equal(t, 7, state.GetApprovedReplicaCount())

	spec := state.GetScalingSpec()
	assert.Equal(t, 7, spec.ApprovedReplicaCount)
}

func TestState_GetScalingSpec(t *testing.T) {
	state := NewState()

	replicas := map[string]ShardReplica{
		"shard-0": {Status: "READY", Message: "Running"},
		"shard-1": {Status: "READY", Message: "Running"},
	}

	state.Update(replicas, 10, "scale-up request")
	state.SetApprovedCount(8)

	spec := state.GetScalingSpec()

	assert.Equal(t, 2, spec.CurrentReplicaCount)
	assert.Equal(t, 10, spec.DesiredReplicaCount)
	assert.Equal(t, 8, spec.ApprovedReplicaCount)
	assert.Equal(t, "scale-up request", spec.LastScalingReason)
}

func TestState_Concurrency(t *testing.T) {
	state := NewState()
	var wg sync.WaitGroup

	// Run concurrent writes and reads
	for i := 0; i < 100; i++ {
		wg.Add(3)

		// Writer goroutine - Update
		go func(i int) {
			defer wg.Done()
			replicas := map[string]ShardReplica{
				"shard-0": {Status: "READY", Message: "Running"},
			}
			state.Update(replicas, i, "concurrent update")
		}(i)

		// Writer goroutine - SetApprovedCount
		go func(i int) {
			defer wg.Done()
			state.SetApprovedCount(i)
		}(i)

		// Reader goroutine - GetScalingSpec
		go func() {
			defer wg.Done()
			_ = state.GetScalingSpec()
		}()
	}

	wg.Wait()

	// If we got here without data race, the test passes
	// The actual values don't matter, we're testing thread safety
	spec := state.GetScalingSpec()
	assert.NotNil(t, spec)
}

func TestState_UpdateWithEmptyReplicas(t *testing.T) {
	state := NewState()

	// Start with some replicas
	replicas := map[string]ShardReplica{
		"shard-0": {Status: "READY", Message: "Running"},
	}
	state.Update(replicas, 3, "initial")
	assert.Equal(t, 1, state.GetCurrentReplicaCount())

	// Update with empty replicas
	state.Update(map[string]ShardReplica{}, 1, "scale-down")
	assert.Equal(t, 0, state.GetCurrentReplicaCount())
	assert.Equal(t, 1, state.GetDesiredReplicaCount())
}

func TestState_UpdateWithMetrics(t *testing.T) {
	state := NewState()

	replicas := map[string]ShardReplica{
		"shard-0": {
			Status:  "READY",
			Message: "Running",
			Metrics: map[string]float64{
				"cpu_usage":    0.75,
				"memory_usage": 0.60,
			},
		},
	}

	state.Update(replicas, 2, "metrics present")

	assert.Equal(t, 1, state.GetCurrentReplicaCount())
}
