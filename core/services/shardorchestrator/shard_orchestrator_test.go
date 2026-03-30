package shardorchestrator_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	ringpb "github.com/smartcontractkit/chainlink-protos/ring/go"
	"github.com/smartcontractkit/chainlink/v2/core/services/ring"
	"github.com/smartcontractkit/chainlink/v2/core/services/shardorchestrator"
)

func setupShardOrchestrator(t *testing.T) (*ring.Store, ringpb.ShardOrchestratorServiceClient, func()) {
	lggr := logger.Test(t)
	ringStore := ring.NewStore()

	ctx := t.Context()
	orchestrator := shardorchestrator.New(0, ringStore, lggr)

	err := orchestrator.Start(ctx)
	require.NoError(t, err)

	var addr string
	require.Eventually(t, func() bool {
		addr = orchestrator.GetAddress()
		return addr != ""
	}, 5*time.Second, 10*time.Millisecond, "orchestrator should have started and be listening")

	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)

	client := ringpb.NewShardOrchestratorServiceClient(conn)

	cleanup := func() {
		conn.Close()
		orchestrator.Close()
	}

	return ringStore, client, cleanup
}

func TestShardOrchestrator_GetWorkflowShardMapping(t *testing.T) {
	t.Run("successfully retrieves workflow mappings", func(t *testing.T) {
		ringStore, client, cleanup := setupShardOrchestrator(t)
		defer cleanup()

		ctx := t.Context()

		ringStore.SetShardForWorkflow("workflow1", 1)
		ringStore.SetShardForWorkflow("workflow2", 2)

		resp, err := client.GetWorkflowShardMapping(ctx, &ringpb.GetWorkflowShardMappingRequest{
			WorkflowIds: []string{"workflow1", "workflow2"},
		})

		require.NoError(t, err)
		require.NotNil(t, resp)
		assert.Len(t, resp.Mappings, 2)
		assert.Equal(t, uint32(1), resp.Mappings["workflow1"])
		assert.Equal(t, uint32(2), resp.Mappings["workflow2"])
		assert.NotEmpty(t, resp.MappingVersion)

		ringStore.SetRoutingState(&ringpb.RoutingState{
			Id:    42,
			State: &ringpb.RoutingState_RoutableShards{RoutableShards: 2},
		})
		resp2, err := client.GetWorkflowShardMapping(ctx, &ringpb.GetWorkflowShardMappingRequest{
			WorkflowIds: []string{"workflow1"},
		})
		require.NoError(t, err)
		assert.Equal(t, uint64(42), resp2.RoutingStateId)
		assert.True(t, resp2.RoutingSteady)

		ringStore.SetRoutingState(&ringpb.RoutingState{
			Id:    43,
			State: &ringpb.RoutingState_Transition{Transition: &ringpb.Transition{WantShards: 3}},
		})
		resp3, err := client.GetWorkflowShardMapping(ctx, &ringpb.GetWorkflowShardMappingRequest{
			WorkflowIds: []string{"workflow1"},
		})
		require.NoError(t, err)
		assert.Equal(t, uint64(43), resp3.RoutingStateId)
		assert.False(t, resp3.RoutingSteady)
	})

	t.Run("returns error for empty workflow IDs", func(t *testing.T) {
		_, client, cleanup := setupShardOrchestrator(t)
		defer cleanup()

		ctx := t.Context()

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

		ctx := t.Context()

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

		ctx := t.Context()

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

		ctx := t.Context()

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
		ringStore, client, cleanup := setupShardOrchestrator(t)
		defer cleanup()

		ctx := t.Context()

		ringStore.SetShardForWorkflow("workflow-a", 1)
		ringStore.SetShardForWorkflow("workflow-b", 2)

		reportResp, err := client.ReportWorkflowTriggerRegistration(ctx, &ringpb.ReportWorkflowTriggerRegistrationRequest{
			SourceShardId: 1,
			RegisteredWorkflows: map[string]uint32{
				"workflow-a": 1,
			},
			TotalActiveWorkflows: 1,
		})
		require.NoError(t, err)
		assert.True(t, reportResp.Success)

		mappingResp, err := client.GetWorkflowShardMapping(ctx, &ringpb.GetWorkflowShardMappingRequest{
			WorkflowIds: []string{"workflow-a", "workflow-b"},
		})
		require.NoError(t, err)
		require.NotNil(t, mappingResp)
		assert.Equal(t, uint32(1), mappingResp.Mappings["workflow-a"])
		assert.Equal(t, uint32(2), mappingResp.Mappings["workflow-b"])

		assert.NotNil(t, mappingResp.MappingStates)
		assert.Contains(t, mappingResp.MappingStates, "workflow-a")
		assert.Contains(t, mappingResp.MappingStates, "workflow-b")
	})
}
