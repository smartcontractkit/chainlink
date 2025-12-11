package monitoring_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/metrics"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/monitoring"
)

func TestUnit_InitMonitoringResources(t *testing.T) {
	_, err := monitoring.InitMonitoringResources()
	require.NoError(t, err)
}

func TestUnit_WorkflowMetricsLabeler(t *testing.T) {
	testWorkflowsMetricLabeler := monitoring.NewWorkflowsMetricLabeler(metrics.NewLabeler(), &monitoring.EngineMetrics{})
	testWorkflowsMetricLabeler2 := testWorkflowsMetricLabeler.With("foo", "baz")
	require.Equal(t, "baz", testWorkflowsMetricLabeler2.Labels["foo"])
}
