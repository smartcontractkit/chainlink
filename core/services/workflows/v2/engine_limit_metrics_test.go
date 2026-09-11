package v2

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/attribute"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"

	"github.com/smartcontractkit/chainlink-common/pkg/metrics"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/monitoring"
)

// collectCounterByAttr sums the named counter and returns the values observed for
// attrKey. Companion to collectCounterValue (engine_org_id_missing_test.go), which
// reads the "reason" attribute instead.
func collectCounterByAttr(t *testing.T, reader *sdkmetric.ManualReader, name, attrKey string) (int64, []string) {
	t.Helper()

	var rm metricdata.ResourceMetrics
	require.NoError(t, reader.Collect(context.Background(), &rm))

	var total int64
	var attrs []string
	for _, sm := range rm.ScopeMetrics {
		for _, m := range sm.Metrics {
			if m.Name != name {
				continue
			}
			data, ok := m.Data.(metricdata.Sum[int64])
			require.True(t, ok, "expected Sum[int64] for %s, got %T", name, m.Data)
			for _, dp := range data.DataPoints {
				total += dp.Value
				if v, ok := dp.Attributes.Value(attribute.Key(attrKey)); ok {
					attrs = append(attrs, v.AsString())
				}
			}
		}
	}
	return total, attrs
}

// TestEngine_LimitDegradationCounters covers the two ways a CRE settings read failure
// degrades a limiter at a call site, which are tracked separately on purpose:
//
//   - read fallback: the limit could not be read, so the static default was substituted.
//     The limit is still applied.
//   - check unenforced: the Check could not be evaluated, so the limit was skipped and
//     the operation allowed through. The limit is NOT applied, which is the more serious
//     of the two and warrants its own alert rather than a label on the other counter.
//
//nolint:paralleltest // setupTestMeter swaps the global beholder client
func TestEngine_LimitDegradationCounters(t *testing.T) {
	const (
		unenforcedCounter = "platform_engine_limit_check_unenforced_total"
		fallbackCounter   = "platform_engine_limit_read_fallback_total"
		limitKeyAttr      = "limitKey"
	)

	t.Run("check unenforced records the limit key", func(t *testing.T) {
		reader := setupTestMeter(t)
		em, err := monitoring.InitMonitoringResources()
		require.NoError(t, err)
		labeler := monitoring.NewWorkflowsMetricLabeler(metrics.NewLabeler(), em)

		labeler.IncrementLimitCheckUnenforcedCounter(t.Context(), "PerWorkflow.LogLineLimit")
		labeler.IncrementLimitCheckUnenforcedCounter(t.Context(), "PerWorkflow.LogEventLimit")

		total, keys := collectCounterByAttr(t, reader, unenforcedCounter, limitKeyAttr)
		assert.Equal(t, int64(2), total)
		assert.ElementsMatch(t, []string{"PerWorkflow.LogLineLimit", "PerWorkflow.LogEventLimit"}, keys)
	})

	t.Run("the two degradation modes are separate series", func(t *testing.T) {
		reader := setupTestMeter(t)
		em, err := monitoring.InitMonitoringResources()
		require.NoError(t, err)
		labeler := monitoring.NewWorkflowsMetricLabeler(metrics.NewLabeler(), em)

		labeler.IncrementLimitReadFallbackCounter(t.Context(), "PerWorkflow.ExecutionTimeout")

		unenforced, _ := collectCounterByAttr(t, reader, unenforcedCounter, limitKeyAttr)
		fallback, fallbackKeys := collectCounterByAttr(t, reader, fallbackCounter, limitKeyAttr)

		assert.Equal(t, int64(1), fallback)
		assert.Equal(t, []string{"PerWorkflow.ExecutionTimeout"}, fallbackKeys)
		assert.Zero(t, unenforced, "a read fallback must not be counted as an unenforced check")
	})
}
