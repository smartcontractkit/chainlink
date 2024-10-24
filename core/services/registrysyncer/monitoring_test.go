package registrysyncer

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/metrics"
)

func Test_InitMonitoringResources(t *testing.T) {
	require.NoError(t, initMonitoringResources())
}

func Test_SyncerMetricsLabeler(t *testing.T) {
	testSyncerMetricLabeler := syncerMetricLabeler{metrics.NewLabeler()}
	testSyncerMetricLabeler2 := testSyncerMetricLabeler.with("foo", "baz")
	require.EqualValues(t, testSyncerMetricLabeler2.Labels["foo"], "baz")
}
