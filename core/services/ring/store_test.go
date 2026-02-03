package ring

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	ringpb "github.com/smartcontractkit/chainlink-protos/ring/go"
)

func TestStore_DeterministicHashing(t *testing.T) {
	store := NewStore()

	// Set up healthy shards
	store.SetAllShardHealth(map[uint32]bool{
		0: true,
		1: true,
		2: true,
	})
	// Simulate OCR having moved to steady state
	store.SetRoutingState(&ringpb.RoutingState{
		State: &ringpb.RoutingState_RoutableShards{RoutableShards: 3},
	})

	ctx := context.Background()

	// Test determinism: same workflow always gets same shard
	shard1, err := store.GetShardForWorkflow(ctx, "workflow-123")
	require.NoError(t, err)
	shard2, err := store.GetShardForWorkflow(ctx, "workflow-123")
	require.NoError(t, err)
	shard3, err := store.GetShardForWorkflow(ctx, "workflow-123")
	require.NoError(t, err)

	require.Equal(t, shard1, shard2, "Same workflow should get same shard (call 2)")
	require.Equal(t, shard2, shard3, "Same workflow should get same shard (call 3)")
	require.LessOrEqual(t, shard1, uint32(2), "Shard should be in healthy set")
}

func TestStore_ConsistentRingConsistency(t *testing.T) {
	store1 := NewStore()
	store2 := NewStore()
	store3 := NewStore()

	// All stores with same healthy shards
	healthyShards := map[uint32]bool{0: true, 1: true, 2: true}
	steadyState := &ringpb.RoutingState{
		State: &ringpb.RoutingState_RoutableShards{RoutableShards: 3},
	}
	store1.SetAllShardHealth(healthyShards)
	store1.SetRoutingState(steadyState)
	store2.SetAllShardHealth(healthyShards)
	store2.SetRoutingState(steadyState)
	store3.SetAllShardHealth(healthyShards)
	store3.SetRoutingState(steadyState)

	ctx := context.Background()

	// All compute same assignments
	workflows := []string{"workflow-A", "workflow-B", "workflow-C", "workflow-D"}
	for _, wf := range workflows {
		s1, err := store1.GetShardForWorkflow(ctx, wf)
		require.NoError(t, err)
		s2, err := store2.GetShardForWorkflow(ctx, wf)
		require.NoError(t, err)
		s3, err := store3.GetShardForWorkflow(ctx, wf)
		require.NoError(t, err)

		require.Equal(t, s1, s2, "All nodes should agree on %s assignment", wf)
		require.Equal(t, s2, s3, "All nodes should agree on %s assignment", wf)
	}
}

func TestStore_Rebalancing(t *testing.T) {
	store := NewStore()
	ctx := context.Background()

	// Start with 3 healthy shards
	store.SetAllShardHealth(map[uint32]bool{0: true, 1: true, 2: true})
	store.SetRoutingState(&ringpb.RoutingState{
		State: &ringpb.RoutingState_RoutableShards{RoutableShards: 3},
	})
	assignments1 := make(map[string]uint32)
	for i := 1; i <= 10; i++ {
		wfID := "workflow-" + string(rune(i))
		shard, err := store.GetShardForWorkflow(ctx, wfID)
		require.NoError(t, err)
		assignments1[wfID] = shard
	}

	// Shard 1 fails
	store.SetShardHealth(1, false)
	assignments2 := make(map[string]uint32)
	for i := 1; i <= 10; i++ {
		wfID := "workflow-" + string(rune(i))
		shard, err := store.GetShardForWorkflow(ctx, wfID)
		require.NoError(t, err)
		assignments2[wfID] = shard
	}

	// Check that rebalancing occurred (some workflows moved)
	healthyShards := store.GetHealthyShards()
	require.Len(t, healthyShards, 2, "Should have 2 healthy shards")
	require.NotContains(t, healthyShards, uint32(1), "Shard 1 should not be healthy")

	// Verify that workflows on healthy shards did not move
	for wfID, originalShard := range assignments1 {
		if originalShard == 0 || originalShard == 2 {
			require.Equal(t, originalShard, assignments2[wfID],
				"Workflow %s on healthy shard %d should not have moved", wfID, originalShard)
		}
	}
}

func TestStore_GetHealthyShards(t *testing.T) {
	store := NewStore()

	store.SetAllShardHealth(map[uint32]bool{
		3: true,
		1: true,
		2: true,
	})

	healthyShards := store.GetHealthyShards()
	require.Len(t, healthyShards, 3)
	// Should be sorted
	require.Equal(t, []uint32{1, 2, 3}, healthyShards)
}

func TestStore_DistributionAcrossShards(t *testing.T) {
	store := NewStore()
	ctx := context.Background()

	store.SetAllShardHealth(map[uint32]bool{
		0: true,
		1: true,
		2: true,
	})
	store.SetRoutingState(&ringpb.RoutingState{
		State: &ringpb.RoutingState_RoutableShards{RoutableShards: 3},
	})

	// Generate many workflows and check distribution
	totalWorkflows := 100
	distribution := make(map[uint32]int)
	for i := 0; i < totalWorkflows; i++ {
		wfID := "workflow-" + string(rune(i))
		shard, err := store.GetShardForWorkflow(ctx, wfID)
		require.NoError(t, err)
		distribution[shard]++
	}

	require.Equal(t, totalWorkflows, sum(distribution), "Should have 100 workflows")

	// Each shard should have roughly 33% of workflows (±5%)
	for shard, count := range distribution {
		pct := float64(count) / 100.0 * 100
		require.GreaterOrEqual(t, pct, 28.0, "Shard %d has too few workflows: %d%%", shard, int(pct))
		require.LessOrEqual(t, pct, 38.0, "Shard %d has too many workflows: %d%%", shard, int(pct))
	}
}

func sum(distribution map[uint32]int) int {
	total := 0
	for _, count := range distribution {
		total += count
	}
	return total
}

func TestStore_GetShardForWorkflow_CacheHit(t *testing.T) {
	store := NewStore()
	ctx := context.Background()

	// Set up steady state
	store.SetAllShardHealth(map[uint32]bool{0: true, 1: true, 2: true})
	store.SetRoutingState(&ringpb.RoutingState{
		State: &ringpb.RoutingState_RoutableShards{RoutableShards: 3},
	})

	// Pre-populate cache with a specific shard assignment
	store.SetShardForWorkflow("cached-workflow", 2)

	// Should return cached value, not recompute
	shard, err := store.GetShardForWorkflow(ctx, "cached-workflow")
	require.NoError(t, err)
	require.Equal(t, uint32(2), shard)
}

func TestStore_GetShardForWorkflow_ContextCancelledDuringSend(t *testing.T) {
	store := NewStore()

	// Put store in transition state
	store.SetAllShardHealth(map[uint32]bool{0: true})
	store.SetRoutingState(&ringpb.RoutingState{
		State: &ringpb.RoutingState_Transition{
			Transition: &ringpb.Transition{WantShards: 2},
		},
	})

	// Fill up the allocRequests channel
	for i := 0; i < AllocationRequestChannelCapacity; i++ {
		store.allocRequests <- AllocationRequest{WorkflowID: "filler"}
	}

	// Context that's already cancelled
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	// Should fail: channel is full and context is cancelled
	_, err := store.GetShardForWorkflow(ctx, "workflow-123")
	require.ErrorIs(t, err, context.Canceled)
}

func TestStore_PendingAllocsDuringTransition(t *testing.T) {
	store := NewStore()
	store.SetAllShardHealth(map[uint32]bool{0: true, 1: true})

	// Put store in transition state
	store.SetRoutingState(&ringpb.RoutingState{
		State: &ringpb.RoutingState_Transition{
			Transition: &ringpb.Transition{WantShards: 3},
		},
	})

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	// Start a goroutine that requests allocation (will block)
	resultCh := make(chan uint32)
	go func() {
		shard, _ := store.GetShardForWorkflow(ctx, "workflow-X")
		resultCh <- shard
	}()

	// Give goroutine time to enqueue request
	time.Sleep(10 * time.Millisecond)

	// Verify request is pending
	pending := store.GetPendingAllocations()
	require.Contains(t, pending, "workflow-X")

	// Fulfill the allocation (simulates transmitter receiving OCR outcome)
	store.SetShardForWorkflow("workflow-X", 2)

	// Blocked goroutine should now receive result
	select {
	case shard := <-resultCh:
		require.Equal(t, uint32(2), shard)
	case <-time.After(50 * time.Millisecond):
		t.Fatal("allocation was not fulfilled")
	}
}

func TestStore_AccessorMethods(t *testing.T) {
	store := NewStore()

	store.SetAllShardHealth(map[uint32]bool{0: true, 1: true, 2: false})
	store.SetRoutingState(&ringpb.RoutingState{
		State: &ringpb.RoutingState_RoutableShards{RoutableShards: 2},
	})
	store.SetShardForWorkflow("wf-1", 0)
	store.SetShardForWorkflow("wf-2", 1)

	t.Run("GetRoutingState", func(t *testing.T) {
		state := store.GetRoutingState()
		require.NotNil(t, state)
		require.Equal(t, uint32(2), state.GetRoutableShards())
	})

	t.Run("IsInTransition_steady_state", func(t *testing.T) {
		require.False(t, store.IsInTransition())
	})

	t.Run("GetShardHealth", func(t *testing.T) {
		health := store.GetShardHealth()
		require.Len(t, health, 3)
		require.True(t, health[0])
		require.True(t, health[1])
		require.False(t, health[2])
	})

	t.Run("GetAllRoutingState", func(t *testing.T) {
		routes := store.GetAllRoutingState()
		require.Len(t, routes, 2)
		require.Equal(t, uint32(0), routes["wf-1"])
		require.Equal(t, uint32(1), routes["wf-2"])
	})

	t.Run("GetHealthyShardCount", func(t *testing.T) {
		require.Equal(t, 2, store.GetHealthyShardCount())
	})

	t.Run("DeleteWorkflow", func(t *testing.T) {
		store.DeleteWorkflow("wf-1")
		routes := store.GetAllRoutingState()
		require.Len(t, routes, 1)
		require.NotContains(t, routes, "wf-1")
	})

	t.Run("IsInTransition_transition_state", func(t *testing.T) {
		store.SetRoutingState(&ringpb.RoutingState{
			State: &ringpb.RoutingState_Transition{Transition: &ringpb.Transition{WantShards: 3}},
		})
		require.True(t, store.IsInTransition())
	})

	t.Run("IsInTransition_nil_state", func(t *testing.T) {
		store.SetRoutingState(nil)
		require.True(t, store.IsInTransition())
	})
}
