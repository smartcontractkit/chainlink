package monitoring_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/metrics"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/monitoring"
)

func Test_InitMonitoringResources(t *testing.T) {
	em, err := monitoring.InitMonitoringResources()
	require.NoError(t, err)
	require.NotNil(t, em)
}

func Test_WorkflowMetricsLabeler(t *testing.T) {
	em, err := monitoring.InitMonitoringResources()
	require.NoError(t, err)
	testWorkflowsMetricLabeler := monitoring.NewWorkflowsMetricLabeler(metrics.NewLabeler(), em)
	testWorkflowsMetricLabeler2 := testWorkflowsMetricLabeler.With("foo", "baz")
	require.Equal(t, "baz", testWorkflowsMetricLabeler2.Labels["foo"])
}
