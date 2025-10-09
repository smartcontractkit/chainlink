package observation

import (
	"context"

	"github.com/smartcontractkit/chainlink-data-streams/llo"
	"github.com/smartcontractkit/chainlink/v2/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/v2/core/services/streams"
)

type Registry interface {
	Get(streamID streams.StreamID) (p streams.Pipeline, exists bool)
}

type Telemeter interface {
	EnqueueV3PremiumLegacy(run *pipeline.Run, trrs pipeline.TaskRunResults, streamID uint32, opts llo.DSOpts, val llo.StreamValue, err error)
	MakeObservationScopedTelemetryCh(opts llo.DSOpts, size int) (ch chan<- any)
	CaptureEATelemetry() bool
	CaptureObservationTelemetry() bool
}

type contextKey string

const ctxObservationTelemetryKey contextKey = "observation-telemetry"

func WithObservationTelemetryCh(ctx context.Context, ch chan<- any) context.Context {
	if ch == nil {
		return ctx
	}
	return context.WithValue(ctx, ctxObservationTelemetryKey, ch)
}

func GetObservationTelemetryCh(ctx context.Context) chan<- any {
	ch, ok := ctx.Value(ctxObservationTelemetryKey).(chan<- any)
	if !ok {
		return nil
	}
	return ch
}
