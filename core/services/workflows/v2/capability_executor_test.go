package v2

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	caperrors "github.com/smartcontractkit/chainlink-common/pkg/capabilities/errors"
	"github.com/smartcontractkit/chainlink-common/pkg/contexts"
	"github.com/smartcontractkit/chainlink-common/pkg/settings"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/cresettings"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/limits"
	"github.com/smartcontractkit/chainlink-protos/cre/go/sdk"
	eventsv2 "github.com/smartcontractkit/chainlink-protos/workflows/go/v2"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func TestExecutionHelper_SystemCapabilityBlocked(t *testing.T) {
	t.Parallel()

	exec := &ExecutionHelper{}

	req := &sdk.CapabilityRequest{
		Id:         confidentialWorkflowsCapabilityID,
		Method:     "Execute",
		CallbackId: 1,
	}

	_, err := exec.CallCapability(t.Context(), req)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "system-only")
}

func TestExecutionHelper_ConfidentialHTTPPerWorkflowLimit(t *testing.T) {
	t.Parallel()

	lggr := logger.TestLogger(t)
	lf := limits.Factory{Logger: lggr}

	// Configure per-workflow confidential-http call limit to 1
	cfgFn := func(w *cresettings.Workflows) {
		w.ConfidentialHTTP.CallLimit = settings.Int(1)
	}

	limiters, err := NewLimiters(lf, cfgFn)
	require.NoError(t, err)
	t.Cleanup(func() { _ = limiters.Close() })

	// Build ExecutionHelper and initialize its call limiters from EngineLimiters
	exec := &ExecutionHelper{}
	exec.initLimiters(limiters)

	// Grab the configured limiter instance for confidential-http SendRequest
	capCallValue := capCall{name: "confidential-http", method: "SendRequest"}
	limiter, ok := exec.callLimiters[capCallValue]
	require.True(t, ok, "expected confidential-http limiter to be configured")

	// Prime the internal callCounts to simulate one prior call so the next call will exceed the configured limit (1)
	exec.callCounts = make(map[limits.Limiter[int]]int)
	exec.callCounts[limiter] = 1

	// Prepare a request that will parse to capName == "confidential-http" and method == "SendRequest"
	req := &sdk.CapabilityRequest{
		Id:         "confidential-http",
		Method:     "SendRequest",
		CallbackId: 1,
	}

	// Call and expect an error from the bound limiter (limit exceeded)
	_, err = exec.CallCapability(t.Context(), req)
	require.Error(t, err, "expected CallCapability to fail when per-workflow confidential-http call limit is exceeded")
	var capErr caperrors.Error
	require.ErrorAs(t, err, &capErr, "expected per-workflow call limit exceedance to be classified as capability user error")
	require.Equal(t, caperrors.OriginUser, capErr.Origin())
	require.Equal(t, caperrors.LimitExceeded, capErr.Code())
}

func TestExecutionHelper_VaultSecretsGetPerWorkflowLimit(t *testing.T) {
	t.Parallel()

	lggr := logger.TestLogger(t)
	lf := limits.Factory{Logger: lggr}

	cfgFn := func(w *cresettings.Workflows) {
		w.Secrets.CallLimit = settings.Int(1)
	}

	limiters, err := NewLimiters(lf, cfgFn)
	require.NoError(t, err)
	t.Cleanup(func() { _ = limiters.Close() })

	exec := &ExecutionHelper{}
	exec.initLimiters(limiters)
	exec.sharedSecretsCounter = &secretsCallCounter{called: 1}

	req := &sdk.CapabilityRequest{
		Id:         "vault",
		Method:     "vault.secrets.get",
		CallbackId: 1,
	}

	_, err = exec.CallCapability(t.Context(), req)
	require.Error(t, err, "expected CallCapability to fail when per-workflow secrets call limit is exceeded")
	var capErr caperrors.Error
	require.ErrorAs(t, err, &capErr, "expected per-workflow call limit exceedance to be classified as capability user error")
	require.Equal(t, caperrors.OriginUser, capErr.Origin())
	require.Equal(t, caperrors.LimitExceeded, capErr.Code())
}

func TestExecutionHelper_VaultSecretsGetSharedCallCounter(t *testing.T) {
	t.Parallel()

	lggr := logger.TestLogger(t)
	lf := limits.Factory{Logger: lggr}

	cfgFn := func(w *cresettings.Workflows) {
		w.Secrets.CallLimit = settings.Int(3)
		w.SecretsConcurrencyLimit = settings.Int(100)
	}

	limiters, err := NewLimiters(lf, cfgFn)
	require.NoError(t, err)
	t.Cleanup(func() { _ = limiters.Close() })

	exec := &ExecutionHelper{Engine: &Engine{}}
	exec.initLimiters(limiters)
	exec.capCallsSemaphore = limits.GlobalResourcePoolLimiter(0)
	exec.sharedSecretsCounter = &secretsCallCounter{}

	// Simulate 2 prior calls via the GetSecrets path (secretsFetcher increments
	// the same shared counter).
	exec.sharedSecretsCounter.called = 2
	require.Equal(t, 2, exec.sharedSecretsCounter.called)

	req := &sdk.CapabilityRequest{
		Id:         "vault",
		Method:     "vault.secrets.get",
		CallbackId: 1,
	}

	// 3rd call: limit check should pass (counter 2->3, limit 3), then fail at
	// the zero-capacity semaphore. The counter must be incremented to 3.
	callCtx, cancel := context.WithTimeout(t.Context(), 50*time.Millisecond)
	defer cancel()
	_, err = exec.CallCapability(callCtx, req)
	require.Error(t, err, "expected call to fail at zero-capacity semaphore")
	require.ErrorIs(t, err, context.DeadlineExceeded)
	require.Equal(t, 3, exec.sharedSecretsCounter.called, "shared counter should be incremented to 3 after a non-limit failure")

	// 4th call: should be rejected at the limit check (counter 3 > limit 3).
	_, err = exec.CallCapability(t.Context(), req)
	require.Error(t, err, "expected CallCapability to fail when shared counter exceeds limit")
	var capErr caperrors.Error
	require.ErrorAs(t, err, &capErr)
	require.Equal(t, caperrors.LimitExceeded, capErr.Code())
	require.Equal(t, 3, exec.sharedSecretsCounter.called, "shared counter must not increment on a failed limit check")
}

func TestExecutionHelper_VaultSecretsGetConcurrencyLimit(t *testing.T) {
	t.Parallel()

	lggr := logger.TestLogger(t)
	lf := limits.Factory{Logger: lggr}

	cfgFn := func(w *cresettings.Workflows) {
		w.Secrets.CallLimit = settings.Int(100)
		w.SecretsConcurrencyLimit = settings.Int(1)
	}

	limiters, err := NewLimiters(lf, cfgFn)
	require.NoError(t, err)
	t.Cleanup(func() { _ = limiters.Close() })

	exec := &ExecutionHelper{}
	exec.initLimiters(limiters)

	capCallValue := capCall{name: "vault", method: "vault.secrets.get"}
	concurrencyLimiter, ok := exec.concurrencyLimiters[capCallValue]
	require.True(t, ok, "expected vault.secrets.get concurrency limiter to be configured")

	ctx := contexts.WithCRE(t.Context(), contexts.CRE{
		Owner:    "1111111111111111111111111111111111111111",
		Workflow: "22222222222222222222222222222222222222222222222222222222222222222",
	})

	free, err := concurrencyLimiter.Wait(ctx, 1)
	require.NoError(t, err)
	defer free()

	req := &sdk.CapabilityRequest{
		Id:         "vault",
		Method:     "vault.secrets.get",
		CallbackId: 1,
	}

	callCtx, cancel := context.WithTimeout(ctx, 50*time.Millisecond)
	defer cancel()

	_, err = exec.CallCapability(callCtx, req)
	require.Error(t, err, "expected CallCapability to fail when secrets concurrency limit is exhausted")
	assert.ErrorIs(t, err, context.DeadlineExceeded)
}

func TestUserMetricTypeSuffix(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		metricType eventsv2.UserMetricType
		wantSuffix string
		wantErr    bool
	}{
		{
			name:       "counter",
			metricType: eventsv2.UserMetricType_USER_METRIC_TYPE_COUNTER,
			wantSuffix: "_counter",
		},
		{
			name:       "gauge",
			metricType: eventsv2.UserMetricType_USER_METRIC_TYPE_GAUGE,
			wantSuffix: "_gauge",
		},
		{
			name:       "unspecified",
			metricType: eventsv2.UserMetricType_USER_METRIC_TYPE_UNSPECIFIED,
			wantErr:    true,
		},
		{
			name:       "unknown numeric value",
			metricType: eventsv2.UserMetricType(999),
			wantErr:    true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			suffix, err := userMetricTypeSuffix(tc.metricType)
			if tc.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), "unsupported user metric type")
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.wantSuffix, suffix)
			}
		})
	}
}

func TestUserMetricNameFormatting(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		metricName string
		metricType eventsv2.UserMetricType
		wantName   string
	}{
		{
			name:       "counter metric gets prefix and suffix",
			metricName: "price",
			metricType: eventsv2.UserMetricType_USER_METRIC_TYPE_COUNTER,
			wantName:   "user_workflow_price_counter",
		},
		{
			name:       "gauge metric gets prefix and suffix",
			metricName: "temperature",
			metricType: eventsv2.UserMetricType_USER_METRIC_TYPE_GAUGE,
			wantName:   "user_workflow_temperature_gauge",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			suffix, err := userMetricTypeSuffix(tc.metricType)
			require.NoError(t, err)
			got := userMetricPrefix + tc.metricName + suffix
			assert.Equal(t, tc.wantName, got)
		})
	}
}

func TestUserMetricUnsupportedTypeRejected(t *testing.T) {
	t.Parallel()

	_, err := userMetricTypeSuffix(eventsv2.UserMetricType_USER_METRIC_TYPE_UNSPECIFIED)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported user metric type")
}
