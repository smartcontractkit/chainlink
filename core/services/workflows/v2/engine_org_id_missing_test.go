package v2

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"

	"github.com/smartcontractkit/chainlink-common/pkg/beholder"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/metrics"
	"github.com/smartcontractkit/chainlink-common/pkg/services/orgresolver"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/monitoring"
)

// testOrgResolver is a test implementation of orgresolver.OrgResolver.
type testOrgResolver struct {
	orgID string
	err   error
}

func (m *testOrgResolver) Get(_ context.Context, _ string) (string, error) {
	return m.orgID, m.err
}

func (m *testOrgResolver) Start(_ context.Context) error { return nil }
func (m *testOrgResolver) Close() error                  { return nil }
func (m *testOrgResolver) HealthReport() map[string]error {
	return map[string]error{m.Name(): nil}
}
func (m *testOrgResolver) Name() string { return "TestOrgResolver" }
func (m *testOrgResolver) Ready() error { return nil }

var _ orgresolver.OrgResolver = (*testOrgResolver)(nil)

// setupTestMeter creates a beholder client backed by a manual meter reader
// so tests can collect and inspect metric increments. Returns the reader and
// registers a cleanup that restores the previous global beholder client.
//
// The client is built from a fully-initialized noop client so that the
// Emitter, Tracer, and Logger fields are always non-nil — this prevents nil
// pointer panics if a parallel test calls beholder.GetEmitter() while our
// client is briefly the global one.
func setupTestMeter(t *testing.T) *sdkmetric.ManualReader {
	t.Helper()

	reader := sdkmetric.NewManualReader()
	mp := sdkmetric.NewMeterProvider(sdkmetric.WithReader(reader))
	t.Cleanup(func() { _ = mp.Shutdown(context.Background()) })

	prevClient := beholder.GetClient()
	t.Cleanup(func() { beholder.SetClient(prevClient) })

	// Start from a noop client so Emitter/Tracer/Logger are valid, then swap
	// only the meter for our manual reader.
	client := beholder.NoopClientConfig{Lggr: logger.Test(t)}.New()
	client.Meter = mp.Meter("beholder")
	client.MeterProvider = mp
	beholder.SetClient(client)

	return reader
}

// collectCounterValue reads metrics from the manual reader and returns the
// total sum of the named counter across all data points, plus the set of
// reason attribute values seen.
func collectCounterValue(t *testing.T, reader *sdkmetric.ManualReader, name string) (int64, []string) {
	t.Helper()

	var rm metricdata.ResourceMetrics
	require.NoError(t, reader.Collect(context.Background(), &rm))

	var total int64
	var reasons []string
	for _, sm := range rm.ScopeMetrics {
		for _, m := range sm.Metrics {
			if m.Name != name {
				continue
			}
			data, ok := m.Data.(metricdata.Sum[int64])
			require.True(t, ok, "expected Sum[int64] for %s, got %T", name, m.Data)
			for _, dp := range data.DataPoints {
				total += dp.Value
				if r, ok := dp.Attributes.Value("reason"); ok {
					reasons = append(reasons, r.AsString())
				}
			}
		}
	}
	return total, reasons
}

// TestEngine_OrgIDMissingReason verifies that the engine correctly records
// why org ID is missing during the resolution logic in start().
func TestEngine_OrgIDMissingReason(t *testing.T) {
	t.Parallel()

	resolverErr := errors.New("linking service unavailable")

	tests := []struct {
		name        string
		orgResolver orgresolver.OrgResolver
		wantReason  string
		wantOrgID   string
		wantErr     error
	}{
		{
			name:        "resolver_nil",
			orgResolver: nil,
			wantReason:  "resolver_nil",
			wantOrgID:   "",
			wantErr:     nil,
		},
		{
			name:        "resolver_error",
			orgResolver: &testOrgResolver{orgID: "", err: resolverErr},
			wantReason:  "resolver_error",
			wantOrgID:   "",
			wantErr:     resolverErr,
		},
		{
			name:        "empty_response",
			orgResolver: &testOrgResolver{orgID: "", err: nil},
			wantReason:  "empty_response",
			wantOrgID:   "",
			wantErr:     nil,
		},
		{
			name:        "valid_org_id",
			orgResolver: &testOrgResolver{orgID: "org-123", err: nil},
			wantReason:  "",
			wantOrgID:   "org-123",
			wantErr:     nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			resolved := resolveOrgID(context.Background(), tt.orgResolver, "owner-1", logger.Sugared(logger.Test(t)))
			require.Equal(t, tt.wantOrgID, resolved.ID)
			require.Equal(t, tt.wantReason, resolved.Reason)
			if tt.wantErr != nil {
				require.EqualError(t, resolved.Err, tt.wantErr.Error())
			} else {
				require.NoError(t, resolved.Err)
			}
		})
	}
}

// TestEngine_OrgIDMissingCounter verifies that the org_id_missing counter is
// incremented with the correct reason label when orgID is empty, and not
// incremented when orgID is set.
//
//nolint:paralleltest // subtests mutate the global beholder client via setupTestMeter
func TestEngine_OrgIDMissingCounter(t *testing.T) {
	const counterName = "platform_engine_org_id_missing_total"

	tests := []struct {
		name       string
		orgID      string
		reason     string
		wantCount  int64
		wantReason string
	}{
		{
			name:       "increments with resolver_nil reason",
			orgID:      "",
			reason:     "resolver_nil",
			wantCount:  1,
			wantReason: "resolver_nil",
		},
		{
			name:       "increments with resolver_error reason",
			orgID:      "",
			reason:     "resolver_error",
			wantCount:  1,
			wantReason: "resolver_error",
		},
		{
			name:       "increments with empty_response reason",
			orgID:      "",
			reason:     "empty_response",
			wantCount:  1,
			wantReason: "empty_response",
		},
		{
			name:      "does not increment when orgID is set",
			orgID:     "org-123",
			reason:    "",
			wantCount: 0,
		},
	}

	for _, tt := range tests {
		//nolint:paralleltest // subtests mutate the global beholder client via setupTestMeter
		t.Run(tt.name, func(t *testing.T) {
			reader := setupTestMeter(t)

			em, err := monitoring.InitMonitoringResources()
			require.NoError(t, err)

			labeler := monitoring.NewWorkflowsMetricLabeler(metrics.NewLabeler(), em)
			engine := &Engine{
				cfg:                &EngineConfig{WorkflowID: "wf-1", WorkflowOwner: "owner-1"},
				orgID:              tt.orgID,
				orgIDMissingReason: tt.reason,
				metrics:            labeler,
			}

			// Simulate the counter increment from startExecution()
			if engine.orgID == "" {
				engine.metrics.IncrementOrgIDMissingCounter(context.Background(), engine.orgIDMissingReason)
			}

			gotCount, gotReasons := collectCounterValue(t, reader, counterName)
			require.Equal(t, tt.wantCount, gotCount)
			if tt.wantCount > 0 {
				require.Contains(t, gotReasons, tt.wantReason)
			}
		})
	}
}
