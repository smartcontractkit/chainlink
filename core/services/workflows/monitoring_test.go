package workflows

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/metrics"
)

func Test_InitMonitoringResources(t *testing.T) {
	_, err := initEngineMonitoringResources()
	require.NoError(t, err)
}

func Test_WorkflowMetricsLabeler(t *testing.T) {
	testWorkflowsMetricLabeler := engineMetricLabeler{metrics.NewLabeler(), engineMetrics{}}
	testWorkflowsMetricLabeler2 := testWorkflowsMetricLabeler.with("foo", "baz")
	require.EqualValues(t, "baz", testWorkflowsMetricLabeler2.Labels["foo"])
}
