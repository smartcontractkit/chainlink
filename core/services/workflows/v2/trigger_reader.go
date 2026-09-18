package v2

import (
	"context"
	"errors"

	"github.com/jonboulle/clockwork"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink/v2/core/platform"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/monitoring"
)

// RunTriggerReader consumes triggerEventCh until it closes or ctx is done,
// converting each received capabilities.TriggerResponse into a
// RoutedTriggerEvent and handing it to deliver. It is the loop body only —
// the caller owns spawning the goroutine (via its own services.Engine), since
// goroutine lifecycle differs between callers.
//
// Shared by Engine's own per-subscription reader and TriggerCoordinator's:
// the only thing that differs between them is what deliver resolves to — the
// engine calling its own Put directly, vs. the coordinator resolving an
// EventSink from the registry on every call. A deliver error is always
// logged and the loop continues; it never exits early on a delivery failure,
// only on ctx.Done() or the channel closing. That includes the coordinator's
// "engine not found" / "engine doesn't accept trigger events" cases — deliver
// returns an error for those like any other delivery failure, rather than
// the reader special-casing an early exit for them.
func RunTriggerReader(
	ctx context.Context,
	lggr logger.Logger,
	metrics *monitoring.WorkflowsMetricLabeler,
	clock clockwork.Clock,
	workflowID, triggerCapID string,
	triggerIndex int,
	triggerEventCh <-chan capabilities.TriggerResponse,
	deliver func(ctx context.Context, event RoutedTriggerEvent) error,
) {
	for {
		select {
		case <-ctx.Done():
			return
		case event, isOpen := <-triggerEventCh:
			if !isOpen {
				return
			}
			eventID := event.Event.ID
			metrics.With(platform.KeyTriggerID, triggerCapID).IncrementTriggerEventReceivedCounter(ctx)
			lggr.Debugw("Processing trigger event", "triggerID", triggerCapID, "eventID", eventID)
			if event.Err != nil {
				lggr.Errorw("Received a trigger event with error, dropping", "triggerID", triggerCapID, "err", event.Err)
				tm := metrics.With(platform.KeyTriggerID, triggerCapID)
				tm.IncrementWorkflowTriggerEventErrorCounter(ctx)
				tm.IncrementTriggerEventDroppedTotal(ctx, monitoring.TriggerDropReasonTriggerResponseError)
				continue
			}

			routed := RoutedTriggerEvent{
				WorkflowID:   workflowID,
				TriggerCapID: triggerCapID,
				TriggerIndex: triggerIndex,
				ObservedAt:   clock.Now(),
				Event:        event,
			}

			if err := deliver(ctx, routed); err != nil {
				// Draining is expected during workflow deletion, so it logs at info rather than error level.
				if errors.Is(err, ErrEngineDraining) {
					lggr.Infow("Dropping trigger event: engine draining", "triggerID", triggerCapID, "eventID", eventID)
				} else {
					lggr.Errorw("Failed to put routed trigger event", "triggerID", triggerCapID, "eventID", eventID, "err", err)
				}
			}
		}
	}
}
