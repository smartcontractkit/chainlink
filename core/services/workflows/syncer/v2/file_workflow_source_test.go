package v2

import (
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func TestFileWorkflowSource_FileNotExists(t *testing.T) {
	lggr := logger.TestLogger(t)
	_, err := NewFileWorkflowSourceWithPath(lggr, "test-file-source", "/nonexistent/path/workflows.json")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "does not exist")
}

func TestFileWorkflowSource_EmptyName(t *testing.T) {
	lggr := logger.TestLogger(t)

	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "workflows.json")
	err := os.WriteFile(tmpFile, []byte("{}"), 0600)
	require.NoError(t, err)

	_, err = NewFileWorkflowSourceWithPath(lggr, "", tmpFile)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "source name is required")
}

func TestFileWorkflowSource_ListWorkflowMetadata_EmptyFile(t *testing.T) {
	lggr := logger.TestLogger(t)

	// Create a temp file
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "workflows.json")
	err := os.WriteFile(tmpFile, []byte(""), 0600)
	require.NoError(t, err)

	source, err := NewFileWorkflowSourceWithPath(lggr, "test-file-source", tmpFile)
	require.NoError(t, err)

	ctx := t.Context()
	don := capabilities.DON{
		ID:       1,
		Families: []string{"workflow"},
	}

	workflows, head, err := source.ListWorkflowMetadata(ctx, don)
	require.NoError(t, err)
	assert.Empty(t, workflows)
	assert.NotNil(t, head)
}

func TestFileWorkflowSource_ListWorkflowMetadata_ValidFile(t *testing.T) {
	lggr := logger.TestLogger(t)

	// Create workflow ID (32 bytes)
	workflowID := make([]byte, 32)
	for i := range workflowID {
		workflowID[i] = byte(i)
	}

	// Create owner (20 bytes for Ethereum address)
	owner := make([]byte, 20)
	for i := range owner {
		owner[i] = byte(i + 100)
	}

	sourceData := FileWorkflowSourceData{
		Workflows: []FileWorkflowMetadata{
			{
				WorkflowID:   hex.EncodeToString(workflowID),
				Owner:        hex.EncodeToString(owner),
				CreatedAt:    1234567890,
				Status:       WorkflowStatusActive, // File format uses 0=active
				WorkflowName: "test-workflow",
				BinaryURL:    "file:///path/to/binary.wasm",
				ConfigURL:    "file:///path/to/config.json",
				Tag:          "v1.0.0",
				DonFamily:    "workflow",
			},
			{
				WorkflowID:   hex.EncodeToString(workflowID),
				Owner:        hex.EncodeToString(owner),
				CreatedAt:    1234567891,
				Status:       WorkflowStatusActive, // File format uses 0=active
				WorkflowName: "other-workflow",
				BinaryURL:    "file:///path/to/other.wasm",
				ConfigURL:    "file:///path/to/other.json",
				Tag:          "v2.0.0",
				DonFamily:    "other-don", // Different DON family
			},
		},
	}

	data, err := json.Marshal(sourceData)
	require.NoError(t, err)

	// Create a temp file
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "workflows.json")
	err = os.WriteFile(tmpFile, data, 0600)
	require.NoError(t, err)

	source, err := NewFileWorkflowSourceWithPath(lggr, "test-file-source", tmpFile)
	require.NoError(t, err)

	ctx := t.Context()
	don := capabilities.DON{
		ID:       1,
		Families: []string{"workflow"},
	}

	workflows, head, err := source.ListWorkflowMetadata(ctx, don)
	require.NoError(t, err)
	assert.Len(t, workflows, 1) // Only one matches the DON family
	assert.NotNil(t, head)

	// Verify the workflow metadata
	wf := workflows[0]
	assert.Equal(t, "test-workflow", wf.WorkflowName)
	assert.Equal(t, "file:///path/to/binary.wasm", wf.BinaryURL)
	assert.Equal(t, "file:///path/to/config.json", wf.ConfigURL)
	assert.Equal(t, "v1.0.0", wf.Tag)
	assert.Equal(t, "workflow", wf.DonFamily)
	assert.Equal(t, WorkflowStatusActive, wf.Status)
	assert.Equal(t, uint64(1234567890), wf.CreatedAt)
}

func TestFileWorkflowSource_ListWorkflowMetadata_MultipleDONFamilies(t *testing.T) {
	lggr := logger.TestLogger(t)

	// Create workflow ID (32 bytes)
	workflowID1 := make([]byte, 32)
	workflowID2 := make([]byte, 32)
	for i := range workflowID1 {
		workflowID1[i] = byte(i)
		workflowID2[i] = byte(i + 50)
	}

	owner := make([]byte, 20)
	for i := range owner {
		owner[i] = byte(i + 100)
	}

	sourceData := FileWorkflowSourceData{
		Workflows: []FileWorkflowMetadata{
			{
				WorkflowID:   hex.EncodeToString(workflowID1),
				Owner:        hex.EncodeToString(owner),
				Status:       WorkflowStatusActive, // File format uses 0=active
				WorkflowName: "workflow-a",
				BinaryURL:    "file:///a.wasm",
				ConfigURL:    "file:///a.json",
				DonFamily:    "family-a",
			},
			{
				WorkflowID:   hex.EncodeToString(workflowID2),
				Owner:        hex.EncodeToString(owner),
				Status:       WorkflowStatusActive, // File format uses 0=active
				WorkflowName: "workflow-b",
				BinaryURL:    "file:///b.wasm",
				ConfigURL:    "file:///b.json",
				DonFamily:    "family-b",
			},
		},
	}

	data, err := json.Marshal(sourceData)
	require.NoError(t, err)

	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "workflows.json")
	err = os.WriteFile(tmpFile, data, 0600)
	require.NoError(t, err)

	source, err := NewFileWorkflowSourceWithPath(lggr, "test-file-source", tmpFile)
	require.NoError(t, err)

	ctx := t.Context()
	don := capabilities.DON{
		ID:       1,
		Families: []string{"family-a", "family-b"},
	}

	workflows, _, err := source.ListWorkflowMetadata(ctx, don)
	require.NoError(t, err)
	assert.Len(t, workflows, 2) // Both workflows match
}

func TestFileWorkflowSource_ListWorkflowMetadata_PausedWorkflow(t *testing.T) {
	lggr := logger.TestLogger(t)

	workflowID := make([]byte, 32)
	for i := range workflowID {
		workflowID[i] = byte(i)
	}
	owner := make([]byte, 20)

	sourceData := FileWorkflowSourceData{
		Workflows: []FileWorkflowMetadata{
			{
				WorkflowID:   hex.EncodeToString(workflowID),
				Owner:        hex.EncodeToString(owner),
				Status:       WorkflowStatusPaused, // File format uses 1=paused
				WorkflowName: "paused-workflow",
				BinaryURL:    "file:///paused.wasm",
				ConfigURL:    "file:///paused.json",
				DonFamily:    "workflow",
			},
		},
	}

	data, err := json.Marshal(sourceData)
	require.NoError(t, err)

	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "workflows.json")
	err = os.WriteFile(tmpFile, data, 0600)
	require.NoError(t, err)

	source, err := NewFileWorkflowSourceWithPath(lggr, "test-file-source", tmpFile)
	require.NoError(t, err)

	ctx := t.Context()
	don := capabilities.DON{
		ID:       1,
		Families: []string{"workflow"},
	}

	workflows, _, err := source.ListWorkflowMetadata(ctx, don)
	require.NoError(t, err)
	assert.Len(t, workflows, 1)
	assert.Equal(t, WorkflowStatusPaused, workflows[0].Status)
}

func TestFileWorkflowSource_Name(t *testing.T) {
	lggr := logger.TestLogger(t)

	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "workflows.json")
	err := os.WriteFile(tmpFile, []byte("{}"), 0600)
	require.NoError(t, err)

	source, err := NewFileWorkflowSourceWithPath(lggr, "my-custom-file-source", tmpFile)
	require.NoError(t, err)
	assert.Equal(t, "my-custom-file-source", source.Name())
}

func TestFileWorkflowSource_Ready(t *testing.T) {
	lggr := logger.TestLogger(t)

	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "workflows.json")
	err := os.WriteFile(tmpFile, []byte("{}"), 0600)
	require.NoError(t, err)

	source, err := NewFileWorkflowSourceWithPath(lggr, "test-file-source", tmpFile)
	require.NoError(t, err)
	assert.NoError(t, source.Ready())

	// Delete the file and check Ready returns error
	err = os.Remove(tmpFile)
	require.NoError(t, err)
	assert.Error(t, source.Ready())
}

func TestFileWorkflowSource_InvalidJSON(t *testing.T) {
	lggr := logger.TestLogger(t)

	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "workflows.json")
	err := os.WriteFile(tmpFile, []byte("invalid json"), 0600)
	require.NoError(t, err)

	source, err := NewFileWorkflowSourceWithPath(lggr, "test-file-source", tmpFile)
	require.NoError(t, err)

	ctx := t.Context()
	don := capabilities.DON{
		ID:       1,
		Families: []string{"workflow"},
	}

	_, _, err = source.ListWorkflowMetadata(ctx, don)
	require.Error(t, err)
}

func TestFileWorkflowSource_InvalidWorkflowID(t *testing.T) {
	lggr := logger.TestLogger(t)

	owner := make([]byte, 20)

	sourceData := FileWorkflowSourceData{
		Workflows: []FileWorkflowMetadata{
			{
				WorkflowID:   "invalid-hex",
				Owner:        hex.EncodeToString(owner),
				Status:       WorkflowStatusActive, // File format uses 0=active
				WorkflowName: "invalid-workflow",
				BinaryURL:    "file:///invalid.wasm",
				ConfigURL:    "file:///invalid.json",
				DonFamily:    "workflow",
			},
		},
	}

	data, err := json.Marshal(sourceData)
	require.NoError(t, err)

	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "workflows.json")
	err = os.WriteFile(tmpFile, data, 0600)
	require.NoError(t, err)

	source, err := NewFileWorkflowSourceWithPath(lggr, "test-file-source", tmpFile)
	require.NoError(t, err)

	ctx := t.Context()
	don := capabilities.DON{
		ID:       1,
		Families: []string{"workflow"},
	}

	// Invalid workflows are skipped, not errored
	workflows, _, err := source.ListWorkflowMetadata(ctx, don)
	require.NoError(t, err)
	assert.Empty(t, workflows)
}
