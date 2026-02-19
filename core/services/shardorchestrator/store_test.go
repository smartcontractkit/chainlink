package shardorchestrator_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/shardorchestrator"
)

func TestStore_BatchUpdateAndQuery(t *testing.T) {
	ctx := context.Background()
	lggr := logger.Test(t)
	store := shardorchestrator.NewStore(lggr)

	// Create and insert multiple workflow mappings
	mappings := []*shardorchestrator.WorkflowMappingState{
		{
			WorkflowID:      "workflow-1",
			OldShardID:      0,
			NewShardID:      1,
			TransitionState: shardorchestrator.StateSteady,
		},
		{
			WorkflowID:      "workflow-2",
			OldShardID:      0,
			NewShardID:      2,
			TransitionState: shardorchestrator.StateSteady,
		},
		{
			WorkflowID:      "workflow-3",
			OldShardID:      0,
			NewShardID:      1,
			TransitionState: shardorchestrator.StateSteady,
		},
	}

	err := store.BatchUpdateWorkflowMappings(ctx, mappings)
	require.NoError(t, err)

	// Query individual workflow
	mapping1, err := store.GetWorkflowMapping(ctx, "workflow-1")
	require.NoError(t, err)
	assert.Equal(t, uint32(1), mapping1.NewShardID)
	assert.Equal(t, shardorchestrator.StateSteady, mapping1.TransitionState)

	// Query all workflows
	allMappings, err := store.GetAllWorkflowMappings(ctx)
	require.NoError(t, err)
	assert.Len(t, allMappings, 3)

	// Query batch
	batchMappings, version, err := store.GetWorkflowMappingsBatch(ctx, []string{"workflow-1", "workflow-2"})
	require.NoError(t, err)
	assert.Len(t, batchMappings, 2)
	assert.Equal(t, uint64(1), version) // First update
}

func TestStore_WorkflowTransition(t *testing.T) {
	ctx := context.Background()
	lggr := logger.Test(t)
	store := shardorchestrator.NewStore(lggr)

	// Initial assignment
	err := store.UpdateWorkflowMapping(ctx, "workflow-123", 0, 1, shardorchestrator.StateSteady)
	require.NoError(t, err)

	mapping, err := store.GetWorkflowMapping(ctx, "workflow-123")
	require.NoError(t, err)
	assert.Equal(t, uint32(1), mapping.NewShardID)
	assert.Equal(t, shardorchestrator.StateSteady, mapping.TransitionState)

	// Move to different shard (transitioning)
	err = store.UpdateWorkflowMapping(ctx, "workflow-123", 1, 3, shardorchestrator.StateTransitioning)
	require.NoError(t, err)

	mapping, err = store.GetWorkflowMapping(ctx, "workflow-123")
	require.NoError(t, err)
	assert.Equal(t, uint32(1), mapping.OldShardID)
	assert.Equal(t, uint32(3), mapping.NewShardID)
	assert.Equal(t, shardorchestrator.StateTransitioning, mapping.TransitionState)

	// Complete transition
	err = store.UpdateWorkflowMapping(ctx, "workflow-123", 1, 3, shardorchestrator.StateSteady)
	require.NoError(t, err)

	mapping, err = store.GetWorkflowMapping(ctx, "workflow-123")
	require.NoError(t, err)
	assert.Equal(t, uint32(3), mapping.NewShardID)
	assert.Equal(t, shardorchestrator.StateSteady, mapping.TransitionState)
}

func TestStore_VersionTracking(t *testing.T) {
	ctx := context.Background()
	lggr := logger.Test(t)
	store := shardorchestrator.NewStore(lggr)

	// Initial version should be 0
	assert.Equal(t, uint64(0), store.GetMappingVersion())

	// First update increments version
	err := store.UpdateWorkflowMapping(ctx, "wf-1", 0, 1, shardorchestrator.StateSteady)
	require.NoError(t, err)
	assert.Equal(t, uint64(1), store.GetMappingVersion())

	// Batch update increments version
	err = store.BatchUpdateWorkflowMappings(ctx, []*shardorchestrator.WorkflowMappingState{
		{WorkflowID: "wf-2", NewShardID: 2, TransitionState: shardorchestrator.StateSteady},
	})
	require.NoError(t, err)
	assert.Equal(t, uint64(2), store.GetMappingVersion())

	// Version is included in batch query response
	_, version, err := store.GetWorkflowMappingsBatch(ctx, []string{"wf-1", "wf-2"})
	require.NoError(t, err)
	assert.Equal(t, uint64(2), version)
}

func TestStore_ShardRegistrations(t *testing.T) {
	ctx := context.Background()
	lggr := logger.Test(t)
	store := shardorchestrator.NewStore(lggr)

	// Shard 1 reports its workflows
	err := store.ReportShardRegistration(ctx, 1, []string{"workflow-1", "workflow-3"})
	require.NoError(t, err)

	// Shard 2 reports its workflows
	err = store.ReportShardRegistration(ctx, 2, []string{"workflow-2"})
	require.NoError(t, err)

	// Query shard registrations
	shard1Workflows, err := store.GetShardRegistrations(ctx, 1)
	require.NoError(t, err)
	assert.Len(t, shard1Workflows, 2)
	assert.Contains(t, shard1Workflows, "workflow-1")
	assert.Contains(t, shard1Workflows, "workflow-3")

	shard2Workflows, err := store.GetShardRegistrations(ctx, 2)
	require.NoError(t, err)
	assert.Len(t, shard2Workflows, 1)
	assert.Contains(t, shard2Workflows, "workflow-2")

	// Query non-existent shard returns empty
	shard3Workflows, err := store.GetShardRegistrations(ctx, 3)
	require.NoError(t, err)
	assert.Empty(t, shard3Workflows)

	// Re-reporting replaces previous registrations
	err = store.ReportShardRegistration(ctx, 1, []string{"workflow-1"})
	require.NoError(t, err)

	shard1Workflows, err = store.GetShardRegistrations(ctx, 1)
	require.NoError(t, err)
	assert.Len(t, shard1Workflows, 1)
	assert.Contains(t, shard1Workflows, "workflow-1")
	assert.NotContains(t, shard1Workflows, "workflow-3") // Removed
}

func TestStore_NotFoundError(t *testing.T) {
	ctx := context.Background()
	lggr := logger.Test(t)
	store := shardorchestrator.NewStore(lggr)

	// Query non-existent workflow
	_, err := store.GetWorkflowMapping(ctx, "non-existent")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestStore_BatchQueryPartialResults(t *testing.T) {
	ctx := context.Background()
	lggr := logger.Test(t)
	store := shardorchestrator.NewStore(lggr)

	// Insert only some workflows
	err := store.UpdateWorkflowMapping(ctx, "exists-1", 0, 1, shardorchestrator.StateSteady)
	require.NoError(t, err)
	err = store.UpdateWorkflowMapping(ctx, "exists-2", 0, 2, shardorchestrator.StateSteady)
	require.NoError(t, err)

	// Query mix of existing and non-existing workflows
	results, _, err := store.GetWorkflowMappingsBatch(ctx, []string{
		"exists-1",
		"non-existent",
		"exists-2",
	})
	require.NoError(t, err)

	// Should only return existing ones
	assert.Len(t, results, 2)
	assert.Contains(t, results, "exists-1")
	assert.Contains(t, results, "exists-2")
	assert.NotContains(t, results, "non-existent")
}
