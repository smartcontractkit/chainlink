package shard

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	shardorchpb "github.com/smartcontractkit/chainlink-common/pkg/workflows/shardorchestrator/pb"
)

type mockShardOrchestratorClient struct {
	shardorchpb.ShardOrchestratorServiceClient
	callCount int
	mappings  map[string]uint32
}

func (m *mockShardOrchestratorClient) GetWorkflowShardMapping(ctx context.Context, req *shardorchpb.GetWorkflowShardMappingRequest, opts ...grpc.CallOption) (*shardorchpb.GetWorkflowShardMappingResponse, error) {
	m.callCount++
	return &shardorchpb.GetWorkflowShardMappingResponse{
		Mappings: m.mappings,
	}, nil
}

func TestFilter_ShouldExecute_CacheHit(t *testing.T) {
	lggr := logger.Test(t)
	mockClient := &mockShardOrchestratorClient{
		mappings: map[string]uint32{
			"workflow-A": 0,
			"workflow-B": 1,
		},
	}

	filter := &Filter{
		orchestratorClient: mockClient,
		myShardIndex:       0,
		cacheTTL:           5 * time.Second,
		logger:             lggr,
		cache:              make(map[string]*cacheEntry),
	}

	ctx := context.Background()

	// First call - cache miss, should call orchestrator
	shouldExecute := filter.ShouldExecute(ctx, "workflow-A")
	assert.True(t, shouldExecute, "workflow-A should execute on shard 0")
	assert.Equal(t, 1, mockClient.callCount, "Should have called orchestrator once")

	// Second call - cache hit, should not call orchestrator again
	shouldExecute = filter.ShouldExecute(ctx, "workflow-A")
	assert.True(t, shouldExecute, "workflow-A should still execute on shard 0")
	assert.Equal(t, 1, mockClient.callCount, "Should not have called orchestrator again (cache hit)")

	// Third call with different workflow - cache hit from first response
	shouldExecute = filter.ShouldExecute(ctx, "workflow-B")
	assert.False(t, shouldExecute, "workflow-B should not execute on shard 0 (assigned to shard 1)")
	assert.Equal(t, 1, mockClient.callCount, "Should not have called orchestrator (cache hit)")
}

func TestFilter_ShouldExecute_CacheExpiry(t *testing.T) {
	lggr := logger.Test(t)
	mockClient := &mockShardOrchestratorClient{
		mappings: map[string]uint32{
			"workflow-A": 0,
		},
	}

	filter := &Filter{
		orchestratorClient: mockClient,
		myShardIndex:       0,
		cacheTTL:           100 * time.Millisecond, // Short TTL for test
		logger:             lggr,
		cache:              make(map[string]*cacheEntry),
	}

	ctx := context.Background()

	// First call
	shouldExecute := filter.ShouldExecute(ctx, "workflow-A")
	assert.True(t, shouldExecute)
	assert.Equal(t, 1, mockClient.callCount)

	// Wait for cache to expire
	time.Sleep(150 * time.Millisecond)

	// Second call after expiry - should call orchestrator again
	shouldExecute = filter.ShouldExecute(ctx, "workflow-A")
	assert.True(t, shouldExecute)
	assert.Equal(t, 2, mockClient.callCount, "Should have called orchestrator again after cache expiry")
}

func TestFilter_ShouldExecute_ShardIndexMatching(t *testing.T) {
	lggr := logger.Test(t)
	mockClient := &mockShardOrchestratorClient{
		mappings: map[string]uint32{
			"workflow-A": 0,
			"workflow-B": 1,
			"workflow-C": 2,
		},
	}

	// Test with shard index 0
	filter0 := &Filter{
		orchestratorClient: mockClient,
		myShardIndex:       0,
		cacheTTL:           5 * time.Second,
		logger:             lggr,
		cache:              make(map[string]*cacheEntry),
	}

	ctx := context.Background()
	assert.True(t, filter0.ShouldExecute(ctx, "workflow-A"), "Shard 0 should execute workflow-A")
	assert.False(t, filter0.ShouldExecute(ctx, "workflow-B"), "Shard 0 should not execute workflow-B")
	assert.False(t, filter0.ShouldExecute(ctx, "workflow-C"), "Shard 0 should not execute workflow-C")

	// Test with shard index 1
	filter1 := &Filter{
		orchestratorClient: mockClient,
		myShardIndex:       1,
		cacheTTL:           5 * time.Second,
		logger:             lggr,
		cache:              make(map[string]*cacheEntry),
	}

	assert.False(t, filter1.ShouldExecute(ctx, "workflow-A"), "Shard 1 should not execute workflow-A")
	assert.True(t, filter1.ShouldExecute(ctx, "workflow-B"), "Shard 1 should execute workflow-B")
	assert.False(t, filter1.ShouldExecute(ctx, "workflow-C"), "Shard 1 should not execute workflow-C")
}

func TestFilter_MyShardIndex(t *testing.T) {
	lggr := logger.Test(t)
	filter := &Filter{
		myShardIndex: 42,
		logger:       lggr,
		cache:        make(map[string]*cacheEntry),
	}

	assert.Equal(t, uint32(42), filter.MyShardIndex())
}

func TestFilter_ClearCache(t *testing.T) {
	lggr := logger.Test(t)
	mockClient := &mockShardOrchestratorClient{
		mappings: map[string]uint32{
			"workflow-A": 0,
		},
	}

	filter := &Filter{
		orchestratorClient: mockClient,
		myShardIndex:       0,
		cacheTTL:           5 * time.Second,
		logger:             lggr,
		cache:              make(map[string]*cacheEntry),
	}

	ctx := context.Background()

	// Populate cache
	filter.ShouldExecute(ctx, "workflow-A")
	assert.Equal(t, 1, mockClient.callCount)

	// Clear cache
	filter.ClearCache()

	// Next call should fetch from orchestrator again
	filter.ShouldExecute(ctx, "workflow-A")
	assert.Equal(t, 2, mockClient.callCount, "Should have called orchestrator after cache clear")
}

func TestNewFilter_Validation(t *testing.T) {
	lggr := logger.Test(t)

	// Missing orchestrator address
	_, err := NewFilter(Config{
		ShardIndex: 0,
		Logger:     lggr,
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "orchestrator address is required")

	// Missing logger
	_, err = NewFilter(Config{
		ShardIndex:          0,
		OrchestratorAddress: "localhost:60051",
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "logger is required")
}
