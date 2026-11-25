package standardcapabilities

import (
	"go.opentelemetry.io/otel/metric"

	"github.com/smartcontractkit/chainlink-common/pkg/beholder"
)

func newStartFailuresCounter() (metric.Int64Counter, error) {
	// The global meter is used because the delegate that creates this service gets no
	// Beholder client or meter to pass on.
	meter := beholder.GetMeter() //nolint:staticcheck // SA1019 see the comment above
	return meter.Int64Counter("platform_standard_capability_start_failures_total",
		metric.WithDescription("Failed start attempts of a standard capability. Every failed attempt is retried."))
}
