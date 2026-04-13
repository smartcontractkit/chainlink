package v2

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	caperrors "github.com/smartcontractkit/chainlink-common/pkg/capabilities/errors"
	"github.com/smartcontractkit/chainlink-common/pkg/settings"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/cresettings"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/limits"
	"github.com/smartcontractkit/chainlink-protos/cre/go/sdk"
	eventsv2 "github.com/smartcontractkit/chainlink-protos/workflows/go/v2"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

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
	require.Equal(t, caperrors.InvalidArgument, capErr.Code())
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
