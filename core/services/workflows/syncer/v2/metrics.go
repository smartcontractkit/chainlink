package v2

import (
	"context"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"

	"github.com/smartcontractkit/chainlink-common/pkg/beholder"
)

type metrics struct {
	handleDuration metric.Int64Histogram
}

func (m *metrics) recordHandleDuration(ctx context.Context, d time.Duration, event string, success bool) {
	// Beholder doesn't support non-string attributes
	successStr := "false"
	if success {
		successStr = "true"
	}
	m.handleDuration.Record(ctx, d.Milliseconds(), metric.WithAttributes(
		attribute.String("success", successStr),
		attribute.String("eventType", event),
	))
}

func newMetrics() (*metrics, error) {
	h, err := beholder.GetMeter().Int64Histogram("platform_workflow_registry_syncer_handler_duration_ms")
	if err != nil {
		return nil, err
	}

	return &metrics{handleDuration: h}, nil
}
