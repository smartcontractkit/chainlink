package v2

import (
	"context"
	"fmt"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/limits"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/types"
	v2 "github.com/smartcontractkit/chainlink/v2/core/services/workflows/v2"
)

// ConsensusEventDispatcher routes trigger events from a shared queue (or OCR callback) to the
// correct workflow engine. Uses EngineRegistry to look up engines by workflowID.
type ConsensusEventDispatcher struct {
	lggr           logger.Logger
	registry       *EngineRegistry
	queue          limits.QueueLimiter[v2.EnqueuedTriggerEvent]
	engineLimiters *v2.EngineLimiters
}

// NewConsensusEventDispatcher creates a dispatcher that routes events to engines.
// queue is used for the standard Wait loop; may be nil when using OCR callback mode only.
func NewConsensusEventDispatcher(
	lggr logger.Logger,
	registry *EngineRegistry,
	queue limits.QueueLimiter[v2.EnqueuedTriggerEvent],
	engineLimiters *v2.EngineLimiters,
) (*ConsensusEventDispatcher, error) {
	if registry == nil {
		return nil, fmt.Errorf("engine registry is required")
	}
	d := &ConsensusEventDispatcher{
		lggr:           lggr,
		registry:       registry,
		queue:          queue,
		engineLimiters: engineLimiters,
	}

	return d, nil
}

// OnConsensusEvent implements ConsensusEventReceiver. Routes the event to the correct engine.
// Called by the OCR Transmitter when using callback mode.
func (d *ConsensusEventDispatcher) OnConsensusEvent(ctx context.Context, event v2.EnqueuedTriggerEvent) error {
	return d.routeEvent(ctx, event)
}

func (d *ConsensusEventDispatcher) routeEvent(ctx context.Context, event v2.EnqueuedTriggerEvent) error {
	wid, err := types.WorkflowIDFromHex(event.WorkflowID())
	if err != nil {
		d.lggr.Warnw("Invalid workflowID in consensus event", "workflowID", event.WorkflowID(), "err", err)
		return err
	}
	entry, ok := d.registry.Get(wid)
	if !ok {
		d.lggr.Debugw("No engine for workflow, skipping consensus event", "workflowID", event.WorkflowID())
		return nil
	}
	receiver, ok := entry.Service.(v2.ConsensusEventReceiver)
	if !ok {
		d.lggr.Warnw("Engine does not implement ConsensusEventReceiver", "workflowID", event.WorkflowID())
		return nil
	}

	return receiver.OnConsensusEvent(ctx, event)
}
