package v2

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/actions/vault"
	caperrors "github.com/smartcontractkit/chainlink-common/pkg/capabilities/errors"
	commonlogger "github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/settings"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/cresettings"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/limits"
	core "github.com/smartcontractkit/chainlink-common/pkg/types/core"
	"github.com/smartcontractkit/chainlink-protos/cre/go/sdk"
	eventsv2 "github.com/smartcontractkit/chainlink-protos/workflows/go/v2"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func TestExecutionHelper_SystemCapabilityBlocked(t *testing.T) {
	t.Parallel()

	resolvedInfo := capabilities.CapabilityInfo{ID: confidentialWorkflowsCapabilityID}
	reg := stubRegistry{cap: stubExecutableCapability{CapabilityInfo: resolvedInfo}}

	engine := &Engine{cfg: &EngineConfig{
		Lggr:        logger.TestLogger(t),
		CapRegistry: reg,
	}}
	engine.setLogger(commonlogger.Sugared(commonlogger.Test(t)))
	exec := &ExecutionHelper{Engine: engine}

	req := &sdk.CapabilityRequest{
		Id:         confidentialWorkflowsCapabilityID,
		Method:     "Execute",
		CallbackId: 1,
	}

	_, err := exec.callCapability(t.Context(), req)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "system-only")
}

// stubExecutableCapability satisfies capabilities.ExecutableCapability. Its
// Info reports the canonical registered ID, simulating a registry that
// resolved a requested version to a registered system-only capability.
type stubExecutableCapability struct {
	capabilities.CapabilityInfo
}

func (stubExecutableCapability) Execute(context.Context, capabilities.CapabilityRequest) (capabilities.CapabilityResponse, error) {
	return capabilities.CapabilityResponse{}, nil
}
func (stubExecutableCapability) RegisterToWorkflow(context.Context, capabilities.RegisterToWorkflowRequest) error {
	return nil
}
func (stubExecutableCapability) UnregisterFromWorkflow(context.Context, capabilities.UnregisterFromWorkflowRequest) error {
	return nil
}

// stubRegistry embeds the unimplemented base and only overrides GetExecutable.
type stubRegistry struct {
	core.UnimplementedCapabilitiesRegistry
	cap capabilities.ExecutableCapability
}

func (s stubRegistry) GetExecutable(context.Context, string) (capabilities.ExecutableCapability, error) {
	return s.cap, nil
}

// TestExecutionHelper_SystemCapabilityResolvedBypass ensures the deny list is
// re-checked against the resolved canonical capability ID. GetExecutable
// performs SemVer-compatible resolution, so a request for a lower prerelease
// (e.g. "confidential-workflows@1.0.0-0") that resolves to the registered
// "confidential-workflows@1.0.0-alpha" must still be rejected.
func TestExecutionHelper_SystemCapabilityResolvedBypass(t *testing.T) {
	t.Parallel()

	resolvedInfo := capabilities.CapabilityInfo{ID: confidentialWorkflowsCapabilityID}
	reg := stubRegistry{cap: stubExecutableCapability{CapabilityInfo: resolvedInfo}}

	engine := &Engine{cfg: &EngineConfig{
		Lggr:        logger.TestLogger(t),
		CapRegistry: reg,
	}}
	engine.setLogger(commonlogger.Sugared(commonlogger.Test(t)))
	exec := &ExecutionHelper{Engine: engine}

	// Request a version strictly below the registered system-only one; the raw
	// ID is not in the deny map, so the fast-path check does not catch it.
	req := &sdk.CapabilityRequest{
		Id:         "confidential-workflows@1.0.0-0",
		Method:     "Execute",
		CallbackId: 1,
	}

	_, err := exec.callCapability(t.Context(), req)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "system-only")
}

// TestExecutionHelper_VaultAndDonTimeSystemCapabilitiesBlocked ensures the
// vault and dontime capabilities cannot be invoked through the raw
// CallCapability path, regardless of the requested version resolving to the
// registered system-only one.
func TestExecutionHelper_VaultAndDonTimeSystemCapabilitiesBlocked(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		registeredID string
		requestID    string
	}{
		{"vault exact ID", vault.CapabilityID, vault.CapabilityID},
		{"vault version-resolved", vault.CapabilityID, "vault@1.0.0-0"},
		{"vault unversioned", vault.CapabilityID, "vault"},
		{"dontime exact ID", "dontime@1.0.0", "dontime@1.0.0"},
		{"dontime version-resolved", "dontime@1.0.0", "dontime@1.0.0-0"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			resolvedInfo := capabilities.CapabilityInfo{ID: tt.registeredID}
			reg := stubRegistry{cap: stubExecutableCapability{CapabilityInfo: resolvedInfo}}

			engine := &Engine{cfg: &EngineConfig{
				Lggr:        logger.TestLogger(t),
				CapRegistry: reg,
			}}
			engine.setLogger(commonlogger.Sugared(commonlogger.Test(t)))
			exec := &ExecutionHelper{Engine: engine}

			req := &sdk.CapabilityRequest{
				Id:         tt.requestID,
				Method:     "Execute",
				CallbackId: 1,
			}

			_, err := exec.callCapability(t.Context(), req)
			require.Error(t, err)
			assert.Contains(t, err.Error(), "system-only")
		})
	}
}

func TestIsSystemCapability(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		capID string
		want  bool
	}{
		{"confidential-workflows registered ID", confidentialWorkflowsCapabilityID, true},
		{"vault registered ID", vault.CapabilityID, true},
		{"dontime registered ID", dontimeCapabilityID, true},
		{"confidential-workflows other version", "confidential-workflows@1.0.0", false},
		{"vault other version", "vault@2.0.0", false},
		{"vault unversioned", "vault", false},
		{"vault with labels", "vault:ChainSelector:123@1.0.0", false},
		{"dontime other version", "dontime@2.0.0", false},
		{"confidential-http is not system-only", "confidential-http@1.0.0", false},
		{"evm with labels is not system-only", "evm:ChainSelector:1@1.0.0", false},
		{"empty ID", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.want, isSystemCapability(tt.capID))
		})
	}
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

// TestExecutionHelper_DefaultCallLimitForUnknownCapabilities ensures raw
// CallCapability calls for (capability, method) pairs without a dedicated
// limiter entry are bounded by the default call limit instead of unbounded.
func TestExecutionHelper_DefaultCallLimitForUnknownCapabilities(t *testing.T) {
	t.Parallel()

	lggr := logger.TestLogger(t)
	lf := limits.Factory{Logger: lggr}

	limiters, err := NewLimiters(lf, nil)
	require.NoError(t, err)
	t.Cleanup(func() { _ = limiters.Close() })

	exec := &ExecutionHelper{}
	exec.initLimiters(limiters)

	require.NoError(t, exec.defaultCallLimiter.Check(t.Context(), defaultCapabilityCallLimit))
	err = exec.defaultCallLimiter.Check(t.Context(), defaultCapabilityCallLimit+1)
	require.ErrorContains(t, err, "limited: cannot use 6, limit is 5")

	// Prime the counter to the default limit, then call an unknown capability
	// through the raw path; the next call must be rejected.
	exec.callCounts = make(map[limits.Limiter[int]]int)
	exec.callCounts[exec.defaultCallLimiter] = defaultCapabilityCallLimit

	req := &sdk.CapabilityRequest{
		Id:         "future-capability@1.0.0",
		Method:     "DoThing",
		CallbackId: 1,
	}

	_, err = exec.CallCapability(t.Context(), req)
	require.Error(t, err, "expected CallCapability to fail when the default call limit is exceeded for an unknown capability")
	var capErr caperrors.Error
	require.ErrorAs(t, err, &capErr, "expected default call limit exceedance to be classified as capability user error")
	require.Equal(t, caperrors.OriginUser, capErr.Origin())
	require.Equal(t, caperrors.LimitExceeded, capErr.Code())
	require.ErrorContains(t, err, "capability call limit exceeded for future-capability.DoThing")
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
