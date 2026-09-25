package triggers

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"google.golang.org/protobuf/types/known/anypb"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/validation"
	"github.com/smartcontractkit/chainlink/v2/core/platform"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/monitoring"
)

const _prefix = "trigger_reg"

// Handle is a registered trigger capability plus the registration
// payload/method needed to unregister and re-deliver to it.
type Handle struct {
	capabilities.TriggerCapability
	Payload *anypb.Any
	Method  string
}

// Ack acknowledges a trigger event against the given handle, logging and
// bumping the same success/failure metrics regardless of caller.
// handle may be nil (i.e., registration not found).
func Ack(
	ctx context.Context,
	lggr logger.Logger,
	metrics *monitoring.WorkflowsMetricLabeler,
	triggerCapID, triggerRegistrationID, eventID string,
	handle *Handle,
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

// RegistrationID constructs a trigger registration ID from a workflow ID and trigger index.
func RegistrationID(workflowID string, triggerIndex int) string {
	return fmt.Sprintf("%s_%s_%d", _prefix, workflowID, triggerIndex)
}

// ParseWorkflowID extracts the workflow ID from a registration ID produced by RegistrationID.
func ParseWorkflowID(registrationID string) (string, error) {
	rest, ok := strings.CutPrefix(registrationID, _prefix+"_")
	if !ok {
		return "", fmt.Errorf("invalid trigger registration ID %q: missing prefix %q", registrationID, _prefix)
	}

	idx := strings.LastIndex(rest, "_")
	if idx == -1 {
		return "", fmt.Errorf("invalid trigger registration ID %q: missing trigger index", registrationID)
	}

	workflowID, triggerIndexStr := rest[:idx], rest[idx+1:]
	if _, err := strconv.Atoi(triggerIndexStr); err != nil {
		return "", fmt.Errorf("invalid trigger registration ID %q: invalid trigger index %q: %w", registrationID, triggerIndexStr, err)
	}

	if err := validation.ValidateWorkflowOrExecutionID(workflowID); err != nil {
		return "", fmt.Errorf("invalid trigger registration ID %q: invalid workflow ID %q: %w", registrationID, workflowID, err)
	}

	return workflowID, nil
}
