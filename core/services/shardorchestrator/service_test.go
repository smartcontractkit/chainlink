package shardorchestrator_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	ringpb "github.com/smartcontractkit/chainlink-protos/ring/go"
	"github.com/smartcontractkit/chainlink/v2/core/services/shardorchestrator"
)

func TestServer_GetWorkflowShardMapping(t *testing.T) {
	ctx := context.Background()
	lggr := logger.Test(t)

	t.Run("returns_mappings_for_multiple_workflows", func(t *testing.T) {
		store := shardorchestrator.NewStore(lggr)
		server := shardorchestrator.NewServer(store, lggr)

		// Set up some workflow mappings
		mappings := []*shardorchestrator.WorkflowMappingState{
			{
				WorkflowID:      "wf-alpha",
				OldShardID:      0,
				NewShardID:      1,
				TransitionState: shardorchestrator.StateSteady,
			},
			{
				WorkflowID:      "wf-beta",
				OldShardID:      0,
				NewShardID:      2,
				TransitionState: shardorchestrator.StateSteady,
			},
			{
				WorkflowID:      "wf-gamma",
				OldShardID:      1,
				NewShardID:      0,
				TransitionState: shardorchestrator.StateTransitioning,
			},
		}
		err := store.BatchUpdateWorkflowMappings(ctx, mappings)
		require.NoError(t, err)

		// Request all three workflows
		req := &ringpb.GetWorkflowShardMappingRequest{
			WorkflowIds: []string{"wf-alpha", "wf-beta", "wf-gamma"},
		}

		resp, err := server.GetWorkflowShardMapping(ctx, req)
		require.NoError(t, err)
		require.NotNil(t, resp)

		// Verify simple mappings
		require.Len(t, resp.Mappings, 3)
		require.Equal(t, uint32(1), resp.Mappings["wf-alpha"])
		require.Equal(t, uint32(2), resp.Mappings["wf-beta"])
		require.Equal(t, uint32(0), resp.Mappings["wf-gamma"])

		// Verify detailed mapping states
		require.Len(t, resp.MappingStates, 3)

		// wf-alpha: steady state
		alphaState := resp.MappingStates["wf-alpha"]
		require.Equal(t, uint32(0), alphaState.OldShardId)
		require.Equal(t, uint32(1), alphaState.NewShardId)
		require.False(t, alphaState.InTransition, "steady state should not be in transition")

		// wf-gamma: transitioning state
		gammaState := resp.MappingStates["wf-gamma"]
		require.Equal(t, uint32(1), gammaState.OldShardId)
		require.Equal(t, uint32(0), gammaState.NewShardId)
		require.True(t, gammaState.InTransition, "transitioning state should be in transition")

		// Verify version
		require.Equal(t, uint64(1), resp.MappingVersion)
	})

	t.Run("rejects_empty_workflow_ids", func(t *testing.T) {
		store := shardorchestrator.NewStore(lggr)
		server := shardorchestrator.NewServer(store, lggr)

		req := &ringpb.GetWorkflowShardMappingRequest{
			WorkflowIds: []string{},
		}

		resp, err := server.GetWorkflowShardMapping(ctx, req)
		require.Error(t, err)
		require.Nil(t, resp)
		require.Contains(t, err.Error(), "required")
	})

	t.Run("handles_partial_results_for_nonexistent_workflows", func(t *testing.T) {
		store := shardorchestrator.NewStore(lggr)
		server := shardorchestrator.NewServer(store, lggr)

		// Add one workflow
		err := store.BatchUpdateWorkflowMappings(ctx, []*shardorchestrator.WorkflowMappingState{
			{WorkflowID: "exists", NewShardID: 1, TransitionState: shardorchestrator.StateSteady},
		})
		require.NoError(t, err)

		// Request one that exists and one that doesn't - batch query silently skips missing workflows
		req := &ringpb.GetWorkflowShardMappingRequest{
			WorkflowIds: []string{"exists", "does-not-exist"},
		}

		resp, err := server.GetWorkflowShardMapping(ctx, req)
		require.NoError(t, err)
		require.NotNil(t, resp)

		// Only the existing workflow is returned
		require.Len(t, resp.Mappings, 1)
		require.Equal(t, uint32(1), resp.Mappings["exists"])
		require.NotContains(t, resp.Mappings, "does-not-exist")
	})
}
