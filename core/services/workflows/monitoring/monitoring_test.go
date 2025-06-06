package monitoring_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/metrics"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/monitoring"
)

func Test_InitMonitoringResources(t *testing.T) {
	_, err := monitoring.InitMonitoringResources()
	require.NoError(t, err)
}

func Test_WorkflowMetricsLabeler(t *testing.T) {
	testWorkflowsMetricLabeler := monitoring.NewWorkflowsMetricLabeler(metrics.NewLabeler(), &monitoring.EngineMetrics{})
	testWorkflowsMetricLabeler2 := testWorkflowsMetricLabeler.With("foo", "baz")
	require.EqualValues(t, "baz", testWorkflowsMetricLabeler2.Labels["foo"])
}
