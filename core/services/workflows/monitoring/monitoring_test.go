package monitoring_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"

	"github.com/smartcontractkit/chainlink-common/pkg/metrics"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/monitoring"
)

func Test_InitMonitoringResources(t *testing.T) {
	t.Parallel()
	em, err := monitoring.InitMonitoringResources()
	require.NoError(t, err)
	require.NotNil(t, em)
}

func Test_InitMonitoringResourcesWithMeter(t *testing.T) {
	t.Parallel()

	reader := sdkmetric.NewManualReader()
	mp := sdkmetric.NewMeterProvider(sdkmetric.WithReader(reader))
	t.Cleanup(func() { require.NoError(t, mp.Shutdown(context.Background())) })

	meter := mp.Meter("test-workflow-monitoring")
	em, err := monitoring.InitMonitoringResourcesWithMeter(meter)
	require.NoError(t, err)
	require.NotNil(t, em)

	labeler := monitoring.NewWorkflowsMetricLabeler(metrics.NewLabeler(), em).With("workflow_id", "wf-123")

	// Verify coordinator and execution instruments record to the private meter
	labeler.UpdateTriggerCoordinatorWorkflowsGauge(t.Context(), 5)
	labeler.UpdateTriggerCoordinatorReadersGauge(t.Context(), 3)
	labeler.IncrementTriggerCoordinatorUnregisterFailuresCounter(t.Context(), 2)
	labeler.UpdateActiveExecutionsGauge(t.Context(), 7)

	var rm metricdata.ResourceMetrics
	require.NoError(t, reader.Collect(t.Context(), &rm))

	gauges := make(map[string]int64)
	counters := make(map[string]int64)
	for _, sm := range rm.ScopeMetrics {
		for _, m := range sm.Metrics {
			switch data := m.Data.(type) {
			case metricdata.Gauge[int64]:
				for _, dp := range data.DataPoints {
					gauges[m.Name] = dp.Value
				}
			case metricdata.Sum[int64]:
				for _, dp := range data.DataPoints {
					counters[m.Name] += dp.Value
				}
			}
		}
	}

	require.Equal(t, int64(5), gauges["platform_trigger_coordinator_workflows"])
	require.Equal(t, int64(3), gauges["platform_trigger_coordinator_readers"])
	require.Equal(t, int64(2), counters["platform_trigger_coordinator_unregister_failures_total"])
	require.Equal(t, int64(7), gauges["platform_engine_active_executions"])
}

func Test_WorkflowMetricsLabeler(t *testing.T) {
	t.Parallel()
	em, err := monitoring.InitMonitoringResources()
	require.NoError(t, err)
	testWorkflowsMetricLabeler := monitoring.NewWorkflowsMetricLabeler(metrics.NewLabeler(), em)
	testWorkflowsMetricLabeler2 := testWorkflowsMetricLabeler.With("foo", "baz")
	require.Equal(t, "baz", testWorkflowsMetricLabeler2.Labels["foo"])
}
