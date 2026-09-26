package v2

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/metrics"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/monitoring"
)

func Test_initFailureReasonForSubscribe(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		err  error
		want string
	}{
		{
			name: "nil error is a generic subscribe failure",
			err:  nil,
			want: initFailureSubscribe,
		},
		{
			// mirrors the chain: host rejects the secrets call, the guest
			// re-reports the host's error text, Subscribe wraps it as a string
			name: "guest-reported secrets rejection is classified",
			err:  fmt.Errorf("failed to execute subscribe: failed to get secrets for call 1: %s", ErrSecretsCallDuringSubscription.Error()),
			want: initFailureDisallowedSecretsCall,
		},
		{
			name: "guest-reported capability rejection is classified",
			err:  fmt.Errorf("failed to execute subscribe: %s", ErrCapabilityCallDuringSubscription.Error()),
			want: initFailureDisallowedCapabilityCall,
		},
		{
			name: "module error is a generic subscribe failure",
			err:  fmt.Errorf("failed to execute subscribe: %w", errors.New("wasm trap")),
			want: initFailureSubscribe,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, initFailureReasonForSubscribe(tt.err))
		})
	}
}

// TestEngine_WorkflowInitializationFailureCounter verifies the counter
// records failures with the reason label, defaulting an empty reason to
// "unknown".
//
//nolint:paralleltest // setupTestMeter swaps the global beholder client
func TestEngine_WorkflowInitializationFailureCounter(t *testing.T) {
	const counterName = "platform_engine_workflow_initialization_failures_total"

	t.Run("records failures with reason label", func(t *testing.T) {
		reader := setupTestMeter(t)

		em, err := monitoring.InitMonitoringResources()
		require.NoError(t, err)
		labeler := monitoring.NewWorkflowsMetricLabeler(metrics.NewLabeler(), em)

		labeler.IncrementWorkflowInitializationFailureCounter(t.Context(), initFailureDisallowedSecretsCall)
		labeler.IncrementWorkflowInitializationFailureCounter(t.Context(), initFailureTriggerRegistration)
		labeler.IncrementWorkflowInitializationFailureCounter(t.Context(), "")

		gotCount, gotReasons := collectCounterValue(t, reader, counterName)
		require.Equal(t, int64(3), gotCount)
		require.ElementsMatch(t, []string{initFailureDisallowedSecretsCall, initFailureTriggerRegistration, "unknown"}, gotReasons)
	})

	t.Run("not incremented on successful initialization", func(t *testing.T) {
		reader := setupTestMeter(t)

		em, err := monitoring.InitMonitoringResources()
		require.NoError(t, err)
		labeler := monitoring.NewWorkflowsMetricLabeler(metrics.NewLabeler(), em)

		labeler.IncrementWorkflowInitializationCounter(t.Context())

		gotCount, _ := collectCounterValue(t, reader, counterName)
		require.Zero(t, gotCount)
	})
}
