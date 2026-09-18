package v2

import (
	"context"
	"fmt"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink/v2/core/platform"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/monitoring"
)

// AckTriggerHandle acknowledges a trigger event against the given handle,
// logging and bumping the same success/failure metrics regardless of caller.
// handle may be nil (registration not found) — the caller resolves it under
// its own lock, since Engine and TriggerCoordinator hold trigger handles in
// different shapes (a flat map vs. one nested per workflow) and shouldn't be
// forced to share the lookup itself.
//
// Shared by Engine.Ack and TriggerCoordinator.Ack.
func AckTriggerHandle(
	ctx context.Context,
	lggr logger.Logger,
	metrics *monitoring.WorkflowsMetricLabeler,
	triggerCapID, triggerRegistrationID, eventID string,
	handle *TriggerHandle,
) error {
	lggr.Infow("ACKing trigger event", "triggerRegistrationID", triggerRegistrationID, "eventID", eventID)

	tm := metrics.With(platform.KeyTriggerID, triggerCapID)

	if handle == nil {
		tm.IncrementTriggerEventAckFailureCounter(ctx)
		return fmt.Errorf("failed to find trigger %s", triggerRegistrationID)
	}
	if err := handle.AckEvent(ctx, triggerRegistrationID, eventID, handle.Method); err != nil {
		tm.IncrementTriggerEventAckFailureCounter(ctx)
		return err
	}
	tm.IncrementTriggerEventAckSuccessCounter(ctx)
	return nil
}
