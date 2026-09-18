package v2

import (
	"context"

	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/monitoring"
)

// WorkflowLimitReporter records workflow count limit rejections (gate 0:
// global and per-owner workflow count).
//
// This is deliberately its own type rather than part of TriggerCoordinator:
// gate 0 is a syncer-level concern (CAPPL-794 tracks moving its check out of
// the engine into the syncer), unrelated to trigger registration/ACK. Note
// that today Engine.init() still performs the check itself and increments
// these same counters directly via its own metrics labeler (see
// GlobalWorkflowLimit.Use in engine.go) — this type is for when that check
// moves to the syncer; until then, constructing one is a no-op in practice
// since nothing calls it.
type WorkflowLimitReporter interface {
	ReportWorkflowLimitPerOwner(ctx context.Context)
	ReportWorkflowLimitGlobal(ctx context.Context)
}

type workflowLimitReporter struct {
	metrics *monitoring.WorkflowsMetricLabeler
}

var _ WorkflowLimitReporter = (*workflowLimitReporter)(nil)

// NewWorkflowLimitReporter returns a WorkflowLimitReporter backed by the
// given metrics labeler.
func NewWorkflowLimitReporter(metrics *monitoring.WorkflowsMetricLabeler) WorkflowLimitReporter {
	return &workflowLimitReporter{metrics: metrics}
}

func (r *workflowLimitReporter) ReportWorkflowLimitPerOwner(ctx context.Context) {
	r.metrics.IncrementWorkflowLimitPerOwnerCounter(ctx)
}

func (r *workflowLimitReporter) ReportWorkflowLimitGlobal(ctx context.Context) {
	r.metrics.IncrementWorkflowLimitGlobalCounter(ctx)
}
