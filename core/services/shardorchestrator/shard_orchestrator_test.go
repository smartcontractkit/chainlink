package shardorchestrator_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	ringpb "github.com/smartcontractkit/chainlink-protos/ring/go"
	"github.com/smartcontractkit/chainlink/v2/core/services/shardorchestrator"
)

// setupShardOrchestrator creates a test ShardOrchestrator and returns the store, client, and cleanup function
func setupShardOrchestrator(t *testing.T) (*shardorchestrator.Store, ringpb.ShardOrchestratorServiceClient, func()) {
	lggr := logger.Test(t)
	store := shardorchestrator.NewStore(lggr)

	ctx := context.Background()
	orchestrator := shardorchestrator.New(0, store, lggr)

	err := orchestrator.Start(ctx)
	require.NoError(t, err)

	// Wait for the server to start and get the listener
	var addr string
	require.Eventually(t, func() bool {
		addr = orchestrator.GetAddress()
		return addr != ""
	}, 5*time.Second, 10*time.Millisecond, "orchestrator should have started and be listening")

	// Connect to the gRPC server
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)

	client := ringpb.NewShardOrchestratorServiceClient(conn)

	cleanup := func() {
		conn.Close()
		orchestrator.Close()
	}

	return store, client, cleanup
}

func TestShardOrchestrator_GetWorkflowShardMapping(t *testing.T) {
	t.Run("successfully retrieves workflow mappings", func(t *testing.T) {
		store, client, cleanup := setupShardOrchestrator(t)
		defer cleanup()

		ctx := context.Background()

		// Add some test mappings to the store
		err := store.UpdateWorkflowMapping(ctx, "workflow1", 0, 1, shardorchestrator.StateSteady)
		require.NoError(t, err)
		err = store.UpdateWorkflowMapping(ctx, "workflow2", 0, 2, shardorchestrator.StateSteady)
		require.NoError(t, err)

		// Call the gRPC endpoint
		resp, err := client.GetWorkflowShardMapping(ctx, &ringpb.GetWorkflowShardMappingRequest{
			WorkflowIds: []string{"workflow1", "workflow2"},
		})

		require.NoError(t, err)
		require.NotNil(t, resp)
		assert.Len(t, resp.Mappings, 2)
		assert.Equal(t, uint32(1), resp.Mappings["workflow1"])
		assert.Equal(t, uint32(2), resp.Mappings["workflow2"])
		assert.NotEmpty(t, resp.MappingVersion)
	})

	t.Run("returns error for empty workflow IDs", func(t *testing.T) {
		_, client, cleanup := setupShardOrchestrator(t)
		defer cleanup()

		ctx := context.Background()

		// Call with empty workflow IDs
		resp, err := client.GetWorkflowShardMapping(ctx, &ringpb.GetWorkflowShardMappingRequest{
			WorkflowIds: []string{},
		})

		require.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "workflow_ids is required")
	})
}

func TestShardOrchestrator_ReportWorkflowTriggerRegistration(t *testing.T) {
	t.Run("successfully reports workflow registration", func(t *testing.T) {
		_, client, cleanup := setupShardOrchestrator(t)
		defer cleanup()

		ctx := context.Background()

		// Report workflows registered on a shard
		resp, err := client.ReportWorkflowTriggerRegistration(ctx, &ringpb.ReportWorkflowTriggerRegistrationRequest{
			SourceShardId: 2,
			RegisteredWorkflows: map[string]uint32{
				"workflow1": 1,
				"workflow2": 1,
			},
			TotalActiveWorkflows: 2,
		})

		require.NoError(t, err)
		require.NotNil(t, resp)
		assert.True(t, resp.Success)
	})

	t.Run("handles empty workflow list", func(t *testing.T) {
		_, client, cleanup := setupShardOrchestrator(t)
		defer cleanup()

		ctx := context.Background()

		// Report with no workflows
		resp, err := client.ReportWorkflowTriggerRegistration(ctx, &ringpb.ReportWorkflowTriggerRegistrationRequest{
			SourceShardId:        3,
			RegisteredWorkflows:  map[string]uint32{},
			TotalActiveWorkflows: 0,
		})

		require.NoError(t, err)
		require.NotNil(t, resp)
		assert.True(t, resp.Success)
	})

	t.Run("handles multiple shards reporting", func(t *testing.T) {
		_, client, cleanup := setupShardOrchestrator(t)
		defer cleanup()

		ctx := context.Background()

		// Multiple shards reporting different workflows
		resp1, err := client.ReportWorkflowTriggerRegistration(ctx, &ringpb.ReportWorkflowTriggerRegistrationRequest{
			SourceShardId: 1,
			RegisteredWorkflows: map[string]uint32{
				"workflow1": 1,
			},
			TotalActiveWorkflows: 1,
		})
		require.NoError(t, err)
		assert.True(t, resp1.Success)

		resp2, err := client.ReportWorkflowTriggerRegistration(ctx, &ringpb.ReportWorkflowTriggerRegistrationRequest{
			SourceShardId: 2,
			RegisteredWorkflows: map[string]uint32{
				"workflow2": 1,
			},
			TotalActiveWorkflows: 1,
		})
		require.NoError(t, err)
		assert.True(t, resp2.Success)
	})
}

func TestShardOrchestrator_Integration(t *testing.T) {
	t.Run("end-to-end workflow registration and retrieval", func(t *testing.T) {
		store, client, cleanup := setupShardOrchestrator(t)
		defer cleanup()

		ctx := context.Background()

		// Step 1: Add workflows to the store
		err := store.UpdateWorkflowMapping(ctx, "workflow-a", 0, 1, shardorchestrator.StateSteady)
		require.NoError(t, err)
		err = store.UpdateWorkflowMapping(ctx, "workflow-b", 0, 2, shardorchestrator.StateSteady)
		require.NoError(t, err)

		// Step 2: Shard 1 reports it has registered workflow-a
		reportResp, err := client.ReportWorkflowTriggerRegistration(ctx, &ringpb.ReportWorkflowTriggerRegistrationRequest{
			SourceShardId: 1,
			RegisteredWorkflows: map[string]uint32{
				"workflow-a": 1,
			},
			TotalActiveWorkflows: 1,
		})
		require.NoError(t, err)
		assert.True(t, reportResp.Success)

		// Step 3: Another shard queries for the mapping
		mappingResp, err := client.GetWorkflowShardMapping(ctx, &ringpb.GetWorkflowShardMappingRequest{
			WorkflowIds: []string{"workflow-a", "workflow-b"},
		})
		require.NoError(t, err)
		require.NotNil(t, mappingResp)
		assert.Equal(t, uint32(1), mappingResp.Mappings["workflow-a"])
		assert.Equal(t, uint32(2), mappingResp.Mappings["workflow-b"])

		// Verify mapping states are included
		assert.NotNil(t, mappingResp.MappingStates)
		assert.Contains(t, mappingResp.MappingStates, "workflow-a")
		assert.Contains(t, mappingResp.MappingStates, "workflow-b")
	})
}
