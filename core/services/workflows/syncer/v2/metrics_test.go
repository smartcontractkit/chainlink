package v2

import (
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/noop"
)

// testMeter returns a private no-op meter for tests that only need instrument
// registration to succeed. Tests that assert recorded values build their own
// meter provider with a ManualReader.
func testMeter() metric.Meter {
	return noop.NewMeterProvider().Meter("test")
}
