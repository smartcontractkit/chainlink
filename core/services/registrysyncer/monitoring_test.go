package registrysyncer

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/metrics"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"
)

func Test_InitMonitoringResources(t *testing.T) {
	tests.BelongsToCISuite(t, "unit")
	_, err := newSyncerMetricLabeler()
	require.NoError(t, err)
}

func Test_SyncerMetricsLabeler(t *testing.T) {
	tests.BelongsToCISuite(t, "unit")
	testSyncerMetricLabeler := syncerMetricLabeler{metrics.NewLabeler(), nil, nil}
	testSyncerMetricLabeler2 := testSyncerMetricLabeler.with("foo", "baz")
	require.Equal(t, "baz", testSyncerMetricLabeler2.Labels["foo"])
}
