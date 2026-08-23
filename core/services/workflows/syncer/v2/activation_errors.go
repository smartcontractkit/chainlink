package v2

import (
	"errors"
	"math"
	"strings"

	eventsv2 "github.com/smartcontractkit/chainlink-protos/workflows/go/v2"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/types"
)

// ActivationRetryPolicy controls how failed workflow activations are retried.
type ActivationRetryPolicy int

const (
	// ActivationRetryable retries with exponential backoff up to MaxActivationRetries
	// (~10h at default 100). Applies to all errors except permanent user/config failures.
	ActivationRetryable ActivationRetryPolicy = iota
	// ActivationNonRetryable drops immediately without further activation retries
	// (permanent user/config errors).
	ActivationNonRetryable
)

func activationRetryCountAsInt32(retryCount int) int32 {
	if retryCount < 0 {
		return 0
	}
	if retryCount > math.MaxInt32 {
		return math.MaxInt32
	}
	return int32(retryCount)
}

func activationAbandonReasonMetricLabel(reason eventsv2.ActivationAbandonReason) string {
	switch reason {
	case eventsv2.ActivationAbandonReason_ACTIVATION_ABANDON_REASON_NON_RETRYABLE:
		return "non_retryable"
	case eventsv2.ActivationAbandonReason_ACTIVATION_ABANDON_REASON_RETRY_LIMIT_EXCEEDED:
		return "retry_limit_exceeded"
	default:
		return "unknown"
	}
}

type activationPolicyError struct {
	err    error
	policy ActivationRetryPolicy
}

func (e *activationPolicyError) Error() string { return e.err.Error() }

func (e *activationPolicyError) Unwrap() error { return e.err }

func nonRetryable(err error) error {
	if err == nil {
		return nil
	}
	return &activationPolicyError{err: err, policy: ActivationNonRetryable}
}

func classifyActivationError(err error) ActivationRetryPolicy {
	if policyErr, ok := errors.AsType[*activationPolicyError](err); ok {
		return policyErr.policy
	}

	if errors.Is(err, types.ErrGlobalWorkflowCountLimitReached) ||
		errors.Is(err, types.ErrPerOwnerWorkflowCountLimitReached) {
		return ActivationNonRetryable
	}

	if isPermanentEngineInitError(err) {
		return ActivationNonRetryable
	}

	return ActivationRetryable
}

func activationRetryPolicyForEvent(eventName WorkflowRegistryEventName, handleErr error) ActivationRetryPolicy {
	if eventName != WorkflowActivated {
		return ActivationRetryable
	}
	return classifyActivationError(handleErr)
}

// isPermanentEngineInitError identifies engine init failures caused by user/config
// mistakes. Prefer sentinel errors at the source; extend this set as errors are catalogued.
func isPermanentEngineInitError(err error) bool {
	for err != nil {
		msg := strings.ToLower(err.Error())
		switch {
		case strings.Contains(msg, "workflowid mismatch"),
			strings.Contains(msg, "invalid workflow id"),
			strings.Contains(msg, "invalid workflow name"),
			strings.Contains(msg, "failed to decode workflow spec binary"),
			strings.Contains(msg, "failed to decode owner"),
			strings.Contains(msg, "invalid cron schedule"),
			strings.Contains(msg, "cron schedule must specify"),
			strings.Contains(msg, "interval exceeded"):
			return true
		}
		err = errors.Unwrap(err)
	}
	return false
}
