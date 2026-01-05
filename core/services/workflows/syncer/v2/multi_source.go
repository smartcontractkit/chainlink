package v2

import (
	"context"
	"time"

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
	metrics *metrics
}

// NewMultiSourceWorkflowAggregator creates a new aggregator with the given sources.
// Sources are queried in order; the first source's head is used if multiple return heads.
func NewMultiSourceWorkflowAggregator(lggr logger.Logger, sources ...WorkflowMetadataSource) *MultiSourceWorkflowAggregator {
	return &MultiSourceWorkflowAggregator{
		lggr:    lggr.Named("MultiSourceWorkflowAggregator"),
		sources: sources,
	}
}

// NewMultiSourceWorkflowAggregatorWithMetrics creates a new aggregator with the given sources and metrics.
func NewMultiSourceWorkflowAggregatorWithMetrics(lggr logger.Logger, m *metrics, sources ...WorkflowMetadataSource) *MultiSourceWorkflowAggregator {
	return &MultiSourceWorkflowAggregator{
		lggr:    lggr.Named("MultiSourceWorkflowAggregator"),
		sources: sources,
		metrics: m,
	}
}

// ListWorkflowMetadata aggregates workflow metadata from all configured sources.
// It continues to query all sources even if some fail, logging errors for failed sources.
//
// Head handling: The contract source's head is preferred (real blockchain head). If no
// contract source is present, the first successful source's head is used. All sources
// guarantee a non-nil head (synthetic if not from blockchain).
//
// Graceful degradation: Even if all sources fail, we return an empty list and nil error
// to allow retry on the next polling cycle. Errors are logged at appropriate levels
// (WARN when all sources fail, ERROR for individual source failures).
func (m *MultiSourceWorkflowAggregator) ListWorkflowMetadata(ctx context.Context, don capabilities.DON) ([]WorkflowMetadataView, *commontypes.Head, error) {
	var allWorkflows []WorkflowMetadataView
	var primaryHead *commontypes.Head
	var sourceErrors []error
	successfulSources := 0

	for _, source := range m.sources {
		sourceName := source.Name()
		start := time.Now()

		// Check if source is ready
		if err := source.Ready(); err != nil {
			m.lggr.Debugw("Source not ready, skipping",
				"source", sourceName,
				"error", err)
			sourceErrors = append(sourceErrors, err)
			// Record metrics for not-ready source
			if m.metrics != nil {
				m.metrics.recordSourceFetch(ctx, sourceName, 0, time.Since(start), err)
			}
			continue
		}

		// Fetch workflows from this source
		workflows, head, err := source.ListWorkflowMetadata(ctx, don)
		duration := time.Since(start)

		// Record metrics for this source fetch
		if m.metrics != nil {
			m.metrics.recordSourceFetch(ctx, sourceName, len(workflows), duration, err)
		}

		if err != nil {
			m.lggr.Errorw("Failed to fetch workflows from source",
				"source", sourceName,
				"error", err,
				"durationMs", duration.Milliseconds())
			sourceErrors = append(sourceErrors, err)
			// Continue to other sources - don't fail completely if one source fails
			continue
		}

		successfulSources++
		m.lggr.Debugw("Fetched workflows from source",
			"source", sourceName,
			"count", len(workflows),
			"durationMs", duration.Milliseconds())

		allWorkflows = append(allWorkflows, workflows...)

		// Prefer contract source head (real blockchain head), fall back to any source's head.
		// All sources guarantee a non-nil head, so no synthetic fallback is needed.
		if head != nil {
			if sourceName == ContractWorkflowSourceName {
				primaryHead = head // Always prefer contract head
			} else if primaryHead == nil {
				primaryHead = head // Use first non-contract head as fallback
			}
		}
	}

	if len(m.sources) > 0 && successfulSources == 0 {
		m.lggr.Warnw("All workflow sources failed - will retry next cycle",
			"sourceCount", len(m.sources),
			"errorCount", len(sourceErrors))
	} else if len(sourceErrors) > 0 {
		m.lggr.Debugw("Some workflow sources failed",
			"successfulSources", successfulSources,
			"failedSources", len(sourceErrors),
			"totalSources", len(m.sources))
	}

	m.lggr.Debugw("Aggregated workflows from all sources",
		"totalWorkflows", len(allWorkflows),
		"sourceCount", len(m.sources),
		"successfulSources", successfulSources)

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
