package registrysyncer

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/metrics"
)

func Test_InitMonitoringResources(t *testing.T) {
	_, err := newSyncerMetricLabeler()
	require.NoError(t, err)
}

func Test_SyncerMetricsLabeler(t *testing.T) {
	testSyncerMetricLabeler := syncerMetricLabeler{metrics.NewLabeler(), nil, nil}
	testSyncerMetricLabeler2 := testSyncerMetricLabeler.with("foo", "baz")
	require.EqualValues(t, "baz", testSyncerMetricLabeler2.Labels["foo"])
}
