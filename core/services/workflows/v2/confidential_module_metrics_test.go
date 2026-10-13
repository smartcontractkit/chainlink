package v2

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/attribute"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"

	caperrors "github.com/smartcontractkit/chainlink-common/pkg/capabilities/errors"
	"github.com/smartcontractkit/chainlink/v2/core/platform"
)

func collectLastFailures(t *testing.T, reader *sdkmetric.ManualReader) map[string]int64 {
	t.Helper()
	var rm metricdata.ResourceMetrics
	require.NoError(t, reader.Collect(context.Background(), &rm))
	values := make(map[string]int64)
	for _, sm := range rm.ScopeMetrics {
		for _, m := range sm.Metrics {
			if m.Name != enclaveLastFailureTimestamp {
				continue
			}
			require.Equal(t, "s", m.Unit)
			gauge, ok := m.Data.(metricdata.Gauge[int64])
			require.True(t, ok)
			for _, dp := range gauge.DataPoints {
				errorType, found := dp.Attributes.Value(attribute.Key(errorTypeAttribute))
				require.True(t, found)
				require.Equal(t, attribute.NewSet(
					attribute.String(platform.KeyWorkflowID, "wf-123"),
					attribute.String(platform.KeyWorkflowOwner, "owner-abc"),
					attribute.String(platform.KeyWorkflowName, "my-workflow"),
					attribute.String(errorTypeAttribute, errorType.AsString()),
				), dp.Attributes)
				require.NotContains(t, values, errorType.AsString())
				values[errorType.AsString()] = dp.Value
			}
		}
	}
	return values
}

func newTestConfidentialMetrics(t *testing.T) (*confidentialModuleMetrics, *sdkmetric.ManualReader) {
	t.Helper()
	reader := sdkmetric.NewManualReader()
	provider := sdkmetric.NewMeterProvider(sdkmetric.WithReader(reader))
	t.Cleanup(func() { require.NoError(t, provider.Shutdown(context.Background())) })
	m, err := newConfidentialModuleMetrics(provider.Meter("test"), "wf-123", "owner-abc", "my-workflow")
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, m.close()) })
	return m, reader
}

func TestConfidentialMetricsRetention(t *testing.T) {
	t.Parallel()
	m, reader := newTestConfidentialMetrics(t)
	assert.Empty(t, collectLastFailures(t, reader), "construction must not register a callback")
	require.NoError(t, m.start())
	require.NoError(t, m.start())
	assert.Equal(t, map[string]int64{"system": 0, "user": 0}, collectLastFailures(t, reader))

	m.recordFailure(t.Context(), errors.New("transport"), time.Unix(200, 0))
	m.recordFailure(t.Context(), errors.New("earlier concurrent completion"), time.Unix(100, 0))
	m.recordFailure(t.Context(), caperrors.NewPublicUserError(errors.New("budget"), caperrors.DeadlineExceeded), time.Unix(300, 0))
	for range 3 {
		assert.Equal(t, map[string]int64{"system": 200, "user": 300}, collectLastFailures(t, reader), "collection must not refresh event time")
	}
	m.recordFailure(t.Context(), errors.New("new failure"), time.Unix(400, 0))
	assert.Equal(t, map[string]int64{"system": 400, "user": 300}, collectLastFailures(t, reader))

	require.NoError(t, m.close())
	require.NoError(t, m.close())
	require.NoError(t, m.start())
	m.recordFailure(t.Context(), errors.New("late completion"), time.Unix(500, 0))
	assert.Empty(t, collectLastFailures(t, reader))
	assert.Empty(t, m.failures)

	replacement, err := newConfidentialModuleMetrics(m.meter, "wf-123", "owner-abc", "my-workflow")
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, replacement.close()) })
	require.NoError(t, replacement.start())
	assert.Equal(t, map[string]int64{"system": 0, "user": 0}, collectLastFailures(t, reader))
}

func TestConfidentialMetricsFirstFailureBeforeCollection(t *testing.T) {
	t.Parallel()
	m, reader := newTestConfidentialMetrics(t)
	require.NoError(t, m.start())
	m.recordFailure(t.Context(), errors.New("first attempt"), time.Unix(100, 0))
	assert.Equal(t, map[string]int64{"system": 100, "user": 0}, collectLastFailures(t, reader))
}

func TestConfidentialMetricsConcurrentLifecycle(t *testing.T) {
	t.Parallel()
	m, reader := newTestConfidentialMetrics(t)
	require.NoError(t, m.start())
	var wg sync.WaitGroup
	for i := range 100 {
		wg.Go(func() {
			m.recordFailure(t.Context(), errors.New("failure"), time.Unix(int64(i+1), 0))
			var rm metricdata.ResourceMetrics
			assert.NoError(t, reader.Collect(t.Context(), &rm))
			assert.NoError(t, m.start())
		})
	}
	wg.Go(func() { assert.NoError(t, m.close()) })
	wg.Wait()
	assert.Empty(t, collectLastFailures(t, reader))
}
