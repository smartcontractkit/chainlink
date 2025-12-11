package registrysyncer

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/metrics"
)

func TestUnit_InitMonitoringResources(t *testing.T) {
	_, err := newSyncerMetricLabeler()
	require.NoError(t, err)
}

func TestUnit_SyncerMetricsLabeler(t *testing.T) {
	testSyncerMetricLabeler := syncerMetricLabeler{metrics.NewLabeler(), nil, nil}
	testSyncerMetricLabeler2 := testSyncerMetricLabeler.with("foo", "baz")
	require.Equal(t, "baz", testSyncerMetricLabeler2.Labels["foo"])
}
