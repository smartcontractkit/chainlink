// Package globalconfigtest provides test helpers for observing offchain capabilities registry
// metrics.
package globalconfigtest

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/attribute"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"

	"github.com/smartcontractkit/chainlink/v2/core/capabilities/globalconfig"
)

// NewMetrics returns Metrics backed by an in-memory reader.
func NewMetrics(t testing.TB) (*globalconfig.Metrics, *sdkmetric.ManualReader) {
	t.Helper()
	reader := sdkmetric.NewManualReader()
	m, err := globalconfig.NewMetrics(sdkmetric.NewMeterProvider(sdkmetric.WithReader(reader)).Meter("test"))
	require.NoError(t, err)
	return m, reader
}

// Series identifies the (domain, env) series of a metric.
func Series(domain, env string) attribute.Distinct {
	set := attribute.NewSet(attribute.String("domain", domain), attribute.String("env", env))
	return set.Equivalent()
}

// Collect returns metric name -> series -> current value for int64 gauges and counters.
func Collect(t testing.TB, reader *sdkmetric.ManualReader) map[string]map[attribute.Distinct]int64 {
	t.Helper()
	var rm metricdata.ResourceMetrics
	require.NoError(t, reader.Collect(t.Context(), &rm))
	out := map[string]map[attribute.Distinct]int64{}
	for _, sm := range rm.ScopeMetrics {
		for _, md := range sm.Metrics {
			out[md.Name] = map[attribute.Distinct]int64{}
			switch data := md.Data.(type) {
			case metricdata.Gauge[int64]:
				for _, dp := range data.DataPoints {
					out[md.Name][dp.Attributes.Equivalent()] = dp.Value
				}
			case metricdata.Sum[int64]:
				for _, dp := range data.DataPoints {
					out[md.Name][dp.Attributes.Equivalent()] = dp.Value
				}
			}
		}
	}
	return out
}
