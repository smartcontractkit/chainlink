package shardorchestrator_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	ringpb "github.com/smartcontractkit/chainlink-protos/ring/go"
	"github.com/smartcontractkit/chainlink/v2/core/services/ring"
	"github.com/smartcontractkit/chainlink/v2/core/services/shardorchestrator"
)

func TestServer_GetWorkflowShardMapping(t *testing.T) {
	t.Run("returns_mappings_for_multiple_workflows", func(t *testing.T) {
		ctx := t.Context()
		lggr := logger.Test(t)
		ringStore := ring.NewStore()
		server := shardorchestrator.NewServer(ringStore, lggr)

		ringStore.SetShardForWorkflow("wf-alpha", 1)
		ringStore.SetShardForWorkflow("wf-beta", 2)
		ringStore.SetShardForWorkflow("wf-gamma", 0)

		req := &ringpb.GetWorkflowShardMappingRequest{
			WorkflowIds: []string{"wf-alpha", "wf-beta", "wf-gamma"},
		}

		resp, err := server.GetWorkflowShardMapping(ctx, req)
		require.NoError(t, err)
		require.NotNil(t, resp)

		require.Len(t, resp.Mappings, 3)
		require.Equal(t, uint32(1), resp.Mappings["wf-alpha"])
		require.Equal(t, uint32(2), resp.Mappings["wf-beta"])
		require.Equal(t, uint32(0), resp.Mappings["wf-gamma"])
	})

	t.Run("rejects_empty_workflow_ids", func(t *testing.T) {
		ctx := t.Context()
		lggr := logger.Test(t)
		ringStore := ring.NewStore()
		server := shardorchestrator.NewServer(ringStore, lggr)

		req := &ringpb.GetWorkflowShardMappingRequest{
			WorkflowIds: []string{},
		}

		resp, err := server.GetWorkflowShardMapping(ctx, req)
		require.Error(t, err)
		require.Nil(t, resp)
		require.Contains(t, err.Error(), "required")
	})

	t.Run("handles_partial_results_for_nonexistent_workflows", func(t *testing.T) {
		ctx := t.Context()
		lggr := logger.Test(t)
		ringStore := ring.NewStore()
		server := shardorchestrator.NewServer(ringStore, lggr)

		ringStore.SetShardForWorkflow("exists", 1)

		req := &ringpb.GetWorkflowShardMappingRequest{
			WorkflowIds: []string{"exists", "does-not-exist"},
		}

		resp, err := server.GetWorkflowShardMapping(ctx, req)
		require.NoError(t, err)
		require.NotNil(t, resp)

		require.Len(t, resp.Mappings, 1)
		require.Equal(t, uint32(1), resp.Mappings["exists"])
		require.NotContains(t, resp.Mappings, "does-not-exist")
	})
}
