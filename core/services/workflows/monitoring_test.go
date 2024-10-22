package workflows

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/metrics"
)

func Test_InitMonitoringResources(t *testing.T) {
	require.NoError(t, initMonitoringResources())
}

func Test_WorkflowMetricsLabeler(t *testing.T) {
	testWorkflowsMetricLabeler := workflowsMetricLabeler{metrics.NewLabeler()}
	testWorkflowsMetricLabeler2 := testWorkflowsMetricLabeler.with("foo", "baz")
	require.EqualValues(t, testWorkflowsMetricLabeler2.Labels["foo"], "baz")
}
