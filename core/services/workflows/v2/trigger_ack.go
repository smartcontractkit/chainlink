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
// handle may be nil (i.e., registration not found)
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
