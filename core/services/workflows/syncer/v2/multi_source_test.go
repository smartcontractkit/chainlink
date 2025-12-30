package v2

import (
	"context"
	"errors"
	"testing"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockWorkflowSource is a mock implementation of WorkflowMetadataSource for testing
type mockWorkflowSource struct {
	name      string
	workflows []WorkflowMetadataView
	head      *commontypes.Head
	err       error
	ready     error
}

func (m *mockWorkflowSource) ListWorkflowMetadata(_ context.Context, _ capabilities.DON) ([]WorkflowMetadataView, *commontypes.Head, error) {
	if m.err != nil {
		return nil, nil, m.err
	}
	return m.workflows, m.head, nil
}

func (m *mockWorkflowSource) Name() string {
	return m.name
}

func (m *mockWorkflowSource) Ready() error {
	return m.ready
}

func TestMultiSourceWorkflowAggregator_SingleSource(t *testing.T) {
	lggr := logger.TestLogger(t)

	workflowID := types.WorkflowID{}
	for i := range workflowID {
		workflowID[i] = byte(i)
	}

	source := &mockWorkflowSource{
		name: "MockSource",
		workflows: []WorkflowMetadataView{
			{
				WorkflowID:   workflowID,
				WorkflowName: "test-workflow",
				Status:       WorkflowStatusActive,
			},
		},
		head: &commontypes.Head{Height: "100"},
	}

	aggregator := NewMultiSourceWorkflowAggregator(lggr, source)

	ctx := context.Background()
	don := capabilities.DON{
		ID:       1,
		Families: []string{"workflow"},
	}

	workflows, head, err := aggregator.ListWorkflowMetadata(ctx, don)
	require.NoError(t, err)
	assert.Len(t, workflows, 1)
	assert.Equal(t, "test-workflow", workflows[0].WorkflowName)
	assert.Equal(t, "100", head.Height)
}

func TestMultiSourceWorkflowAggregator_MultipleSources(t *testing.T) {
	lggr := logger.TestLogger(t)

	workflowID1 := types.WorkflowID{}
	workflowID2 := types.WorkflowID{}
	for i := range workflowID1 {
		workflowID1[i] = byte(i)
		workflowID2[i] = byte(i + 50)
	}

	source1 := &mockWorkflowSource{
		name: "ContractSource",
		workflows: []WorkflowMetadataView{
			{
				WorkflowID:   workflowID1,
				WorkflowName: "contract-workflow",
				Status:       WorkflowStatusActive,
			},
		},
		head: &commontypes.Head{Height: "100"},
	}

	source2 := &mockWorkflowSource{
		name: "FileSource",
		workflows: []WorkflowMetadataView{
			{
				WorkflowID:   workflowID2,
				WorkflowName: "file-workflow",
				Status:       WorkflowStatusActive,
			},
		},
		head: &commontypes.Head{Height: "50"}, // Lower height, should be ignored
	}

	aggregator := NewMultiSourceWorkflowAggregator(lggr, source1, source2)

	ctx := context.Background()
	don := capabilities.DON{
		ID:       1,
		Families: []string{"workflow"},
	}

	workflows, head, err := aggregator.ListWorkflowMetadata(ctx, don)
	require.NoError(t, err)
	assert.Len(t, workflows, 2)
	// First source's head is used
	assert.Equal(t, "100", head.Height)

	// Check both workflows are present
	names := make(map[string]bool)
	for _, wf := range workflows {
		names[wf.WorkflowName] = true
	}
	assert.True(t, names["contract-workflow"])
	assert.True(t, names["file-workflow"])
}

func TestMultiSourceWorkflowAggregator_SourceNotReady(t *testing.T) {
	lggr := logger.TestLogger(t)

	workflowID := types.WorkflowID{}
	for i := range workflowID {
		workflowID[i] = byte(i)
	}

	source1 := &mockWorkflowSource{
		name:  "NotReadySource",
		ready: errors.New("contract reader not initialized"),
	}

	source2 := &mockWorkflowSource{
		name: "ReadySource",
		workflows: []WorkflowMetadataView{
			{
				WorkflowID:   workflowID,
				WorkflowName: "ready-workflow",
				Status:       WorkflowStatusActive,
			},
		},
		head: &commontypes.Head{Height: "100"},
	}

	aggregator := NewMultiSourceWorkflowAggregator(lggr, source1, source2)

	ctx := context.Background()
	don := capabilities.DON{
		ID:       1,
		Families: []string{"workflow"},
	}

	// Should still succeed with the ready source
	workflows, head, err := aggregator.ListWorkflowMetadata(ctx, don)
	require.NoError(t, err)
	assert.Len(t, workflows, 1)
	assert.Equal(t, "ready-workflow", workflows[0].WorkflowName)
	assert.Equal(t, "100", head.Height)
}

func TestMultiSourceWorkflowAggregator_SourceError(t *testing.T) {
	lggr := logger.TestLogger(t)

	workflowID := types.WorkflowID{}
	for i := range workflowID {
		workflowID[i] = byte(i)
	}

	source1 := &mockWorkflowSource{
		name: "ErrorSource",
		err:  errors.New("failed to fetch"),
	}

	source2 := &mockWorkflowSource{
		name: "GoodSource",
		workflows: []WorkflowMetadataView{
			{
				WorkflowID:   workflowID,
				WorkflowName: "good-workflow",
				Status:       WorkflowStatusActive,
			},
		},
		head: &commontypes.Head{Height: "100"},
	}

	aggregator := NewMultiSourceWorkflowAggregator(lggr, source1, source2)

	ctx := context.Background()
	don := capabilities.DON{
		ID:       1,
		Families: []string{"workflow"},
	}

	// Should still succeed with the good source (errors are logged, not propagated)
	workflows, head, err := aggregator.ListWorkflowMetadata(ctx, don)
	require.NoError(t, err)
	assert.Len(t, workflows, 1)
	assert.Equal(t, "good-workflow", workflows[0].WorkflowName)
	assert.Equal(t, "100", head.Height)
}

func TestMultiSourceWorkflowAggregator_AllSourcesFail(t *testing.T) {
	lggr := logger.TestLogger(t)

	source1 := &mockWorkflowSource{
		name:  "NotReadySource",
		ready: errors.New("not ready"),
	}

	source2 := &mockWorkflowSource{
		name: "ErrorSource",
		err:  errors.New("failed to fetch"),
	}

	aggregator := NewMultiSourceWorkflowAggregator(lggr, source1, source2)

	ctx := context.Background()
	don := capabilities.DON{
		ID:       1,
		Families: []string{"workflow"},
	}

	// Should return empty list, not error (graceful degradation)
	workflows, head, err := aggregator.ListWorkflowMetadata(ctx, don)
	require.NoError(t, err)
	assert.Empty(t, workflows)
	assert.NotNil(t, head)
	assert.Equal(t, "0", head.Height)
}

func TestMultiSourceWorkflowAggregator_NoSources(t *testing.T) {
	lggr := logger.TestLogger(t)

	aggregator := NewMultiSourceWorkflowAggregator(lggr)

	ctx := context.Background()
	don := capabilities.DON{
		ID:       1,
		Families: []string{"workflow"},
	}

	workflows, head, err := aggregator.ListWorkflowMetadata(ctx, don)
	require.NoError(t, err)
	assert.Empty(t, workflows)
	assert.NotNil(t, head)
}

func TestMultiSourceWorkflowAggregator_AddSource(t *testing.T) {
	lggr := logger.TestLogger(t)

	aggregator := NewMultiSourceWorkflowAggregator(lggr)
	assert.Len(t, aggregator.Sources(), 0)

	source := &mockWorkflowSource{
		name: "AddedSource",
	}

	aggregator.AddSource(source)
	assert.Len(t, aggregator.Sources(), 1)
	assert.Equal(t, "AddedSource", aggregator.Sources()[0].Name())
}

func TestMultiSourceWorkflowAggregator_HeadPriority(t *testing.T) {
	lggr := logger.TestLogger(t)

	// First source has nil head
	source1 := &mockWorkflowSource{
		name:      "NilHeadSource",
		workflows: []WorkflowMetadataView{},
		head:      nil,
	}

	// Second source has valid head
	source2 := &mockWorkflowSource{
		name:      "ValidHeadSource",
		workflows: []WorkflowMetadataView{},
		head:      &commontypes.Head{Height: "200"},
	}

	aggregator := NewMultiSourceWorkflowAggregator(lggr, source1, source2)

	ctx := context.Background()
	don := capabilities.DON{
		ID:       1,
		Families: []string{"workflow"},
	}

	_, head, err := aggregator.ListWorkflowMetadata(ctx, don)
	require.NoError(t, err)
	// Should use the first non-nil head
	assert.Equal(t, "200", head.Height)
}



