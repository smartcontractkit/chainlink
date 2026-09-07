package v2

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/types"
)

func Test_classifyActivationError(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		err  error
		want ActivationRetryPolicy
	}{
		{
			name: "artifact fetch is retryable",
			err:  &types.ArtifactFetchError{ArtifactType: "binary", URL: "http://example", Err: errors.New("503")},
			want: ActivationRetryable,
		},
		{
			name: "wrapped artifact fetch is retryable",
			err:  fmt.Errorf("create spec: %w", &types.ArtifactFetchError{ArtifactType: "config", URL: "http://example", Err: errors.New("timeout")}),
			want: ActivationRetryable,
		},
		{
			name: "global limit is non-retryable",
			err:  types.ErrGlobalWorkflowCountLimitReached,
			want: ActivationNonRetryable,
		},
		{
			name: "per owner limit is non-retryable",
			err:  types.ErrPerOwnerWorkflowCountLimitReached,
			want: ActivationNonRetryable,
		},
		{
			name: "explicit non-retryable wrapper",
			err:  nonRetryable(errors.New("workflowID mismatch")),
			want: ActivationNonRetryable,
		},
		{
			name: "workflow id mismatch message is non-retryable",
			err:  fmt.Errorf("engine initialization failed: %w", errors.New("workflowID mismatch: abc != def")),
			want: ActivationNonRetryable,
		},
		{
			name: "failure parsing cron schedule",
			err: fmt.Errorf("failed to register trigger %s: %w", "trigger_0",
				errors.New("[3]InvalidArgument: failed to initialize job: gocron: CronJob: crontab parse failure\nprovided bad location moon: unknown time zone moon")),
			want: ActivationNonRetryable,
		},
		{
			name: "cron schedule faster than the allowed minimum",
			err: fmt.Errorf("failed to register trigger %s: %w", "trigger_0",
				errors.New("[3]InvalidArgument: maximum fastest cron schedule is 30s")),
			want: ActivationNonRetryable,
		},
		{
			name: "unknown error is retryable",
			err:  errors.New("unexpected engine failure"),
			want: ActivationRetryable,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, classifyActivationError(tt.err))
		})
	}
}

func Test_activationRetryPolicyForEvent(t *testing.T) {
	t.Parallel()

	t.Run("non-activation events are retryable", func(t *testing.T) {
		t.Parallel()
		policy := activationRetryPolicyForEvent(WorkflowDeleted, &types.ArtifactFetchError{})
		require.Equal(t, ActivationRetryable, policy)
	})

	t.Run("activation artifact fetch is retryable", func(t *testing.T) {
		t.Parallel()
		policy := activationRetryPolicyForEvent(WorkflowActivated, &types.ArtifactFetchError{})
		require.Equal(t, ActivationRetryable, policy)
	})
}
