package v2

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/types"
)

const (
	// FileWorkflowSourceName is the name used for logging and identification.
	FileWorkflowSourceName = "FileWorkflowSource"
)

// FileWorkflowMetadata represents a single workflow entry in the JSON file.
// This mirrors the WorkflowMetadataView structure but uses JSON-friendly types.
type FileWorkflowMetadata struct {
	// WorkflowID is the hex-encoded workflow ID (without 0x prefix)
	WorkflowID string `json:"workflow_id"`
	// Owner is the hex-encoded owner address (without 0x prefix)
	Owner string `json:"owner"`
	// CreatedAt is the Unix timestamp when the workflow was created
	CreatedAt uint64 `json:"created_at"`
	// Status is the workflow status (0=active, 1=paused)
	Status uint8 `json:"status"`
	// WorkflowName is the human-readable name of the workflow
	WorkflowName string `json:"workflow_name"`
	// BinaryURL is the URL to fetch the workflow binary (same format as contract)
	BinaryURL string `json:"binary_url"`
	// ConfigURL is the URL to fetch the workflow config (same format as contract)
	ConfigURL string `json:"config_url"`
	// Tag is the workflow tag/version
	Tag string `json:"tag"`
	// Attributes is optional JSON-encoded attributes
	Attributes string `json:"attributes,omitempty"`
	// DonFamily is the DON family this workflow belongs to
	DonFamily string `json:"don_family"`
}

// FileWorkflowSourceData is the root structure of the JSON file.
type FileWorkflowSourceData struct {
	// Workflows is the list of workflow metadata entries
	Workflows []FileWorkflowMetadata `json:"workflows"`
}

// FileWorkflowSource implements WorkflowMetadataSource by reading from a JSON file.
type FileWorkflowSource struct {
	lggr     logger.Logger
	filePath string
	name     string
	mu       sync.RWMutex
}

// NewFileWorkflowSourceWithPath creates a new file-based workflow source with a custom path.
// The name parameter is required and must be unique across all workflow sources.
// Returns an error if name is empty or if the file does not exist.
func NewFileWorkflowSourceWithPath(lggr logger.Logger, name string, path string) (*FileWorkflowSource, error) {
	if name == "" {
		return nil, errors.New("source name is required")
	}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, errors.New("workflow metadata file does not exist: " + path)
	}
	return &FileWorkflowSource{
		lggr:     logger.Named(lggr, name),
		filePath: path,
		name:     name,
	}, nil
}

// ListWorkflowMetadata reads the JSON file and returns workflow metadata filtered by DON families.
func (f *FileWorkflowSource) ListWorkflowMetadata(ctx context.Context, don capabilities.DON) ([]WorkflowMetadataView, *commontypes.Head, error) {
	f.tryInitialize(ctx)

	f.mu.RLock()
	defer f.mu.RUnlock()

	filePath := f.filePath

	// Read file contents
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, nil, err
	}

	// Handle empty file
	if len(data) == 0 {
		f.lggr.Debugw("Workflow metadata file is empty, returning empty list", "path", filePath)
		return []WorkflowMetadataView{}, f.syntheticHead(), nil
	}

	// Parse JSON
	var sourceData FileWorkflowSourceData
	if err := json.Unmarshal(data, &sourceData); err != nil {
		return nil, nil, err
	}

	// Build a set of DON families for efficient lookup
	donFamilySet := make(map[string]bool)
	for _, family := range don.Families {
		donFamilySet[family] = true
	}

	// Filter and convert workflows
	workflows := make([]WorkflowMetadataView, 0, len(sourceData.Workflows))
	for _, wf := range sourceData.Workflows {
		// Filter by DON family
		if !donFamilySet[wf.DonFamily] {
			continue
		}

		// Convert to WorkflowMetadataView
		view, err := f.toWorkflowMetadataView(wf)
		if err != nil {
			f.lggr.Warnw("Failed to parse workflow metadata, skipping",
				"source", f.name,
				"workflowName", wf.WorkflowName,
				"error", err)
			continue
		}

		workflows = append(workflows, view)
	}

	f.lggr.Debugw("Loaded workflows from file",
		"path", filePath,
		"totalInFile", len(sourceData.Workflows),
		"matchingDON", len(workflows),
		"donFamilies", don.Families)

	return workflows, f.syntheticHead(), nil
}

func (f *FileWorkflowSource) Name() string {
	return f.name
}

// Ready returns nil if the file exists, or an error if it doesn't.
func (f *FileWorkflowSource) Ready() error {
	if _, err := os.Stat(f.filePath); os.IsNotExist(err) {
		return errors.New("workflow metadata file does not exist: " + f.filePath)
	}
	return nil
}

// tryInitialize always returns true (file validated in constructor).
func (f *FileWorkflowSource) tryInitialize(_ context.Context) bool {
	return true
}

// toWorkflowMetadataView converts a FileWorkflowMetadata to a WorkflowMetadataView.
func (f *FileWorkflowSource) toWorkflowMetadataView(wf FileWorkflowMetadata) (WorkflowMetadataView, error) {
	// Parse workflow ID from hex string
	workflowIDBytes, err := hex.DecodeString(wf.WorkflowID)
	if err != nil {
		return WorkflowMetadataView{}, errors.New("invalid workflow_id hex: " + err.Error())
	}
	if len(workflowIDBytes) != 32 {
		return WorkflowMetadataView{}, errors.New("workflow_id must be 32 bytes")
	}
	var workflowID types.WorkflowID
	copy(workflowID[:], workflowIDBytes)

	// Parse owner from hex string
	ownerBytes, err := hex.DecodeString(wf.Owner)
	if err != nil {
		return WorkflowMetadataView{}, errors.New("invalid owner hex: " + err.Error())
	}

	// Parse attributes if present
	var attributes []byte
	if wf.Attributes != "" {
		attributes = []byte(wf.Attributes)
	}

	return WorkflowMetadataView{
		WorkflowID:   workflowID,
		Owner:        ownerBytes,
		CreatedAt:    wf.CreatedAt,
		Status:       FileStatusToInternal(wf.Status),
		WorkflowName: wf.WorkflowName,
		BinaryURL:    wf.BinaryURL,
		ConfigURL:    wf.ConfigURL,
		Tag:          wf.Tag,
		Attributes:   attributes,
		DonFamily:    wf.DonFamily,
		Source:       f.SourceIdentifier(),
	}, nil
}

// SourceIdentifier returns the source identifier used in WorkflowMetadataView.Source.
// Format: file:{source_name}:v1
func (f *FileWorkflowSource) SourceIdentifier() string {
	return fmt.Sprintf("file:%s:v1", f.name)
}

// syntheticHead creates a synthetic head for the file source.
// Since file sources don't have blockchain blocks, we use the current timestamp.
func (f *FileWorkflowSource) syntheticHead() *commontypes.Head {
	now := time.Now().Unix()
	var timestamp uint64
	if now >= 0 { // satisfies overflow check on linter
		timestamp = uint64(now)
	}
	return &commontypes.Head{
		Height:    strconv.FormatInt(now, 10),
		Hash:      []byte("file-source"),
		Timestamp: timestamp,
	}
}
