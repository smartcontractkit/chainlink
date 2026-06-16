package v2

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_activationRetriesExhausted(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name                 string
		retryCount           int
		maxActivationRetries int
		want                 bool
	}{
		{name: "disabled", retryCount: 100, maxActivationRetries: 0, want: false},
		{name: "below limit", retryCount: 2, maxActivationRetries: 3, want: false},
		{name: "at limit", retryCount: 3, maxActivationRetries: 3, want: true},
		{name: "above limit", retryCount: 5, maxActivationRetries: 3, want: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, activationRetriesExhausted(tt.retryCount, tt.maxActivationRetries))
		})
	}
}

func Test_droppedActivations(t *testing.T) {
	t.Parallel()

	dropped := newDroppedActivations()
	const source = "ContractWorkflowSource"
	const workflowID = "abc"
	const signature = "WorkflowActivated-abc-1"

	require.False(t, dropped.isDropped(source, workflowID, signature))
	dropped.drop(source, workflowID, signature)
	require.True(t, dropped.isDropped(source, workflowID, signature))
	require.False(t, dropped.isDropped(source, workflowID, "other-signature"))

	dropped.clear(source, workflowID)
	require.False(t, dropped.isDropped(source, workflowID, signature))
}
