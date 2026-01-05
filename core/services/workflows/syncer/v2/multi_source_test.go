package v2

import (
	"context"
	"crypto/sha256"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/types"
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

	workflowID := types.WorkflowID(sha256.Sum256([]byte("workflowID")))

	// Use ContractWorkflowSource name to get real head
	source := &mockWorkflowSource{
		name: ContractWorkflowSourceName,
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

	workflowID1 := types.WorkflowID(sha256.Sum256([]byte("workflowID1")))
	workflowID2 := types.WorkflowID(sha256.Sum256([]byte("workflowID2")))

	// ContractWorkflowSource provides the real blockchain head
	source1 := &mockWorkflowSource{
		name: ContractWorkflowSourceName,
		workflows: []WorkflowMetadataView{
			{
				WorkflowID:   workflowID1,
				WorkflowName: "contract-workflow",
				Status:       WorkflowStatusActive,
			},
		},
		head: &commontypes.Head{Height: "100"},
	}

	// FileSource head is ignored (only ContractWorkflowSource head is used)
	source2 := &mockWorkflowSource{
		name: FileWorkflowSourceName,
		workflows: []WorkflowMetadataView{
			{
				WorkflowID:   workflowID2,
				WorkflowName: "file-workflow",
				Status:       WorkflowStatusActive,
			},
		},
		head: &commontypes.Head{Height: "50"}, // This is ignored
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
	// Only ContractWorkflowSource head is used
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

	workflowID := types.WorkflowID(sha256.Sum256([]byte("workflowID")))

	// ContractWorkflowSource is not ready
	source1 := &mockWorkflowSource{
		name:  ContractWorkflowSourceName,
		ready: errors.New("contract reader not initialized"),
	}

	// FileSource is ready and its head is used as fallback
	source2 := &mockWorkflowSource{
		name: FileWorkflowSourceName,
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

	// Should still succeed with the ready source, using fallback head
	workflows, head, err := aggregator.ListWorkflowMetadata(ctx, don)
	require.NoError(t, err)
	assert.Len(t, workflows, 1)
	assert.Equal(t, "ready-workflow", workflows[0].WorkflowName)
	// Since ContractWorkflowSource is not ready, we get fallback head from FileSource
	assert.NotNil(t, head)
	assert.Equal(t, "100", head.Height)
}

func TestMultiSourceWorkflowAggregator_SourceError(t *testing.T) {
	lggr := logger.TestLogger(t)

	workflowID := types.WorkflowID(sha256.Sum256([]byte("workflowID")))

	// ContractWorkflowSource fails
	source1 := &mockWorkflowSource{
		name: ContractWorkflowSourceName,
		err:  errors.New("failed to fetch"),
	}

	// Alternative source succeeds and its head is used as fallback
	source2 := &mockWorkflowSource{
		name: "GRPCSource",
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
	// and use fallback head from GRPCSource since ContractWorkflowSource failed
	workflows, head, err := aggregator.ListWorkflowMetadata(ctx, don)
	require.NoError(t, err)
	assert.Len(t, workflows, 1)
	assert.Equal(t, "good-workflow", workflows[0].WorkflowName)
	assert.NotNil(t, head)
	assert.Equal(t, "100", head.Height)
}

func TestMultiSourceWorkflowAggregator_AllSourcesFail(t *testing.T) {
	lggr := logger.TestLogger(t)

	source1 := &mockWorkflowSource{
		name:  ContractWorkflowSourceName,
		ready: errors.New("not ready"),
	}

	source2 := &mockWorkflowSource{
		name: "GRPCSource",
		err:  errors.New("failed to fetch"),
	}

	aggregator := NewMultiSourceWorkflowAggregator(lggr, source1, source2)

	ctx := context.Background()
	don := capabilities.DON{
		ID:       1,
		Families: []string{"workflow"},
	}

	// Should return empty list, not error (graceful degradation)
	// Gets synthetic head since all sources failed
	workflows, head, err := aggregator.ListWorkflowMetadata(ctx, don)
	require.NoError(t, err)
	assert.Empty(t, workflows)
	assert.NotNil(t, head)
	assert.Equal(t, []byte("synthetic-multi-source"), head.Hash)
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
	assert.Empty(t, aggregator.Sources())

	source := &mockWorkflowSource{
		name: "AddedSource",
	}

	aggregator.AddSource(source)
	assert.Len(t, aggregator.Sources(), 1)
	assert.Equal(t, "AddedSource", aggregator.Sources()[0].Name())
}

func TestMultiSourceWorkflowAggregator_HeadPriority(t *testing.T) {
	lggr := logger.TestLogger(t)

	// Alternative source comes first with valid head (but ignored)
	source1 := &mockWorkflowSource{
		name:      "GRPCSource",
		workflows: []WorkflowMetadataView{},
		head:      &commontypes.Head{Height: "300"}, // Ignored
	}

	// ContractWorkflowSource comes second but its head is used
	source2 := &mockWorkflowSource{
		name:      ContractWorkflowSourceName,
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
	// Should use ContractWorkflowSource head, not the first source
	assert.Equal(t, "200", head.Height)
}

func TestMultiSourceWorkflowAggregator_FallbackHeadForAlternativeOnly(t *testing.T) {
	lggr := logger.TestLogger(t)

	// Only alternative sources (no ContractWorkflowSource)
	source1 := &mockWorkflowSource{
		name:      "GRPCSource",
		workflows: []WorkflowMetadataView{},
		head:      &commontypes.Head{Height: "100"}, // First source, used as fallback
	}

	source2 := &mockWorkflowSource{
		name:      FileWorkflowSourceName,
		workflows: []WorkflowMetadataView{},
		head:      &commontypes.Head{Height: "50"},
	}

	aggregator := NewMultiSourceWorkflowAggregator(lggr, source1, source2)

	ctx := context.Background()
	don := capabilities.DON{
		ID:       1,
		Families: []string{"workflow"},
	}

	_, head, err := aggregator.ListWorkflowMetadata(ctx, don)
	require.NoError(t, err)
	// Should get fallback head from first source (GRPCSource) since no ContractWorkflowSource
	assert.NotNil(t, head)
	assert.Equal(t, "100", head.Height)
}
