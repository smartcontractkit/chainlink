package v2

import (
	"context"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

// MultiSourceWorkflowAggregator aggregates workflow metadata from multiple WorkflowMetadataSource
// implementations. This allows the workflow registry syncer to reconcile workflows from various
// sources (e.g., on-chain contracts, file-based sources, APIs) in a unified manner.
type MultiSourceWorkflowAggregator struct {
	lggr    logger.Logger
	sources []WorkflowMetadataSource
}

// NewMultiSourceWorkflowAggregator creates a new aggregator with the given sources.
// Sources are queried in order; the first source's head is used if multiple return heads.
func NewMultiSourceWorkflowAggregator(lggr logger.Logger, sources ...WorkflowMetadataSource) *MultiSourceWorkflowAggregator {
	return &MultiSourceWorkflowAggregator{
		lggr:    lggr.Named("MultiSourceWorkflowAggregator"),
		sources: sources,
	}
}

// ListWorkflowMetadata aggregates workflow metadata from all configured sources.
// It continues to query all sources even if some fail, logging errors for failed sources.
// The returned head is from the first source that returns a non-nil head (typically the contract source).
//
// NOTE: For the MVP, we assume workflowID collisions between sources are handled externally
// (e.g., by having separate workflow registry contracts with non-overlapping ID spaces).
// If a collision occurs, workflows from later sources will be appended (both will be present).
func (m *MultiSourceWorkflowAggregator) ListWorkflowMetadata(ctx context.Context, don capabilities.DON) ([]WorkflowMetadataView, *commontypes.Head, error) {
	var allWorkflows []WorkflowMetadataView
	var primaryHead *commontypes.Head

	for _, source := range m.sources {
		sourceName := source.Name()

		// Check if source is ready
		if err := source.Ready(); err != nil {
			m.lggr.Debugw("Source not ready, skipping",
				"source", sourceName,
				"error", err)
			continue
		}

		// Fetch workflows from this source
		workflows, head, err := source.ListWorkflowMetadata(ctx, don)
		if err != nil {
			m.lggr.Errorw("Failed to fetch workflows from source",
				"source", sourceName,
				"error", err)
			// Continue to other sources - don't fail completely if one source fails
			continue
		}

		m.lggr.Debugw("Fetched workflows from source",
			"source", sourceName,
			"count", len(workflows))

		allWorkflows = append(allWorkflows, workflows...)

		// Use the first source's head as the primary head (typically contract source)
		// This is because the contract source provides actual blockchain head data,
		// while file sources provide synthetic heads.
		if primaryHead == nil && head != nil {
			primaryHead = head
		}
	}

	// If no head was obtained from any source, create a default one
	if primaryHead == nil {
		primaryHead = &commontypes.Head{Height: "0"}
	}

	m.lggr.Debugw("Aggregated workflows from all sources",
		"totalWorkflows", len(allWorkflows),
		"sourceCount", len(m.sources))

	return allWorkflows, primaryHead, nil
}

// AddSource adds a new workflow metadata source to the aggregator.
// Sources added later will be queried after existing sources.
func (m *MultiSourceWorkflowAggregator) AddSource(source WorkflowMetadataSource) {
	m.sources = append(m.sources, source)
	m.lggr.Debugw("Added workflow metadata source",
		"source", source.Name(),
		"totalSources", len(m.sources))
}

// Sources returns the list of configured sources (for debugging/testing).
func (m *MultiSourceWorkflowAggregator) Sources() []WorkflowMetadataSource {
	return m.sources
}



