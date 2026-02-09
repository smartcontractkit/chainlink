package ring_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	ringpb "github.com/smartcontractkit/chainlink-protos/ring/go"
	"github.com/smartcontractkit/chainlink/v2/core/services/ring"
	"github.com/smartcontractkit/chainlink/v2/core/services/shardorchestrator"
)

// mockArbiterScalerClient implements ringpb.ArbiterScalerClient for testing.
// It allows configuring the return values for Status and ConsensusWantShards.
type mockArbiterScalerClient struct {
	wantShards   uint32
	shardStatus  map[uint32]*ringpb.ShardStatus
	statusErr    error
	consensusErr error
}

func newMockArbiterScalerClient() *mockArbiterScalerClient {
	return &mockArbiterScalerClient{
		wantShards:  1,
		shardStatus: map[uint32]*ringpb.ShardStatus{0: {IsHealthy: true}},
	}
}

func (m *mockArbiterScalerClient) Status(_ context.Context, _ *emptypb.Empty, _ ...grpc.CallOption) (*ringpb.ReplicaStatus, error) {
	if m.statusErr != nil {
		return nil, m.statusErr
	}
	return &ringpb.ReplicaStatus{
		WantShards: m.wantShards,
		Status:     m.shardStatus,
	}, nil
}

func (m *mockArbiterScalerClient) ConsensusWantShards(_ context.Context, _ *ringpb.ConsensusWantShardsRequest, _ ...grpc.CallOption) (*emptypb.Empty, error) {
	if m.consensusErr != nil {
		return nil, m.consensusErr
	}
	return &emptypb.Empty{}, nil
}

func TestRingStoreIntegration(t *testing.T) {
	t.Run("RingStore can be created and used", func(t *testing.T) {
		store := ring.NewStore()
		require.NotNil(t, store)

		// Set shard health - required for consistent hashing to work
		store.SetShardHealth(1, true)
		health := store.GetShardHealth()
		require.True(t, health[1])

		// Set steady state routing (required for GetShardForWorkflow without OCR)
		store.SetRoutingState(&ringpb.RoutingState{
			Id:    1,
			State: &ringpb.RoutingState_RoutableShards{RoutableShards: 1},
		})

		// Test workflow routing via consistent hashing
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		shardID, err := store.GetShardForWorkflow(ctx, "workflow1")
		require.NoError(t, err)
		// Should route to shard 1 since it's the only healthy shard
		require.Equal(t, uint32(1), shardID)
	})

	t.Run("RingStore routes to multiple shards", func(t *testing.T) {
		store := ring.NewStore()

		// Set up multiple healthy shards
		store.SetShardHealth(0, true)
		store.SetShardHealth(1, true)
		store.SetShardHealth(2, true)

		// Set steady state with 3 shards
		store.SetRoutingState(&ringpb.RoutingState{
			Id:    1,
			State: &ringpb.RoutingState_RoutableShards{RoutableShards: 3},
		})

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Test that different workflows can route to different shards
		shard1, err := store.GetShardForWorkflow(ctx, "workflow-a")
		require.NoError(t, err)

		shard2, err := store.GetShardForWorkflow(ctx, "workflow-b")
		require.NoError(t, err)

		// Both should be valid shard IDs (0, 1, or 2)
		require.LessOrEqual(t, shard1, uint32(2), "shard1 should be <= 2")
		require.LessOrEqual(t, shard2, uint32(2), "shard2 should be <= 2")
	})

	t.Run("RingStore caches workflow allocations", func(t *testing.T) {
		store := ring.NewStore()

		store.SetShardHealth(0, true)
		store.SetRoutingState(&ringpb.RoutingState{
			Id:    1,
			State: &ringpb.RoutingState_RoutableShards{RoutableShards: 1},
		})

		// Manually set a workflow allocation
		store.SetShardForWorkflow("cached-workflow", 0)

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Should return the cached value
		shardID, err := store.GetShardForWorkflow(ctx, "cached-workflow")
		require.NoError(t, err)
		require.Equal(t, uint32(0), shardID)
	})
}

func TestRingFactoryIntegration(t *testing.T) {
	t.Run("RingFactory can be created", func(t *testing.T) {
		lggr := logger.Test(t)
		store := ring.NewStore()
		shardOrchestratorStore := shardorchestrator.NewStore(lggr)
		mockArbiter := newMockArbiterScalerClient()

		factory, err := ring.NewFactory(store, shardOrchestratorStore, mockArbiter, lggr, nil)
		require.NoError(t, err)
		require.NotNil(t, factory)
	})
}
