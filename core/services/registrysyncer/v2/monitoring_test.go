package v2

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/metric/noop"

	"github.com/smartcontractkit/chainlink-common/pkg/metrics"
)

func Test_InitMonitoringResources(t *testing.T) {
	_, err := newSyncerMetricLabeler(noop.NewMeterProvider().Meter("test"))
	require.NoError(t, err)
}

func Test_SyncerMetricsLabeler(t *testing.T) {
	testSyncerMetricLabeler := syncerMetricLabeler{metrics.NewLabeler(), nil, nil}
	testSyncerMetricLabeler2 := testSyncerMetricLabeler.with("foo", "baz")
	require.Equal(t, "baz", testSyncerMetricLabeler2.Labels["foo"])
}
