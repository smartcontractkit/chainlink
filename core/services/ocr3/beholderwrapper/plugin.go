package beholderwrapper

import (
	"context"
	"time"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink/v2/core/services/ocr3/beholderwrapper/metrics"
)

// MetricPrefix is the prefix for all OCR3 beholder metrics
const MetricPrefix = "platform_ocr3_reporting_plugin"

var _ ocr3types.ReportingPlugin[any] = &reportingPlugin[any]{}

type reportingPlugin[RI any] struct {
	ocr3types.ReportingPlugin[RI]
	metrics *metrics.PluginMetrics
}

func newReportingPlugin[RI any](
	origin ocr3types.ReportingPlugin[RI],
	m *metrics.PluginMetrics,
) *reportingPlugin[RI] {
	return &reportingPlugin[RI]{
		ReportingPlugin: origin,
		metrics:         m,
	}
}

func (p *reportingPlugin[RI]) Query(ctx context.Context, outctx ocr3types.OutcomeContext) (ocrtypes.Query, error) {
	return withObservedExecution(ctx, p.metrics, metrics.Query, func() (ocrtypes.Query, error) {
		return p.ReportingPlugin.Query(ctx, outctx)
	})
}

func (p *reportingPlugin[RI]) Observation(ctx context.Context, outctx ocr3types.OutcomeContext, query ocrtypes.Query) (ocrtypes.Observation, error) {
	result, err := withObservedExecution(ctx, p.metrics, metrics.Observation, func() (ocrtypes.Observation, error) {
		return p.ReportingPlugin.Observation(ctx, outctx, query)
	})
	if err == nil {
		p.metrics.TrackSize(ctx, metrics.Observation, len(result))
	}
	return result, err
}

func (p *reportingPlugin[RI]) ValidateObservation(ctx context.Context, outctx ocr3types.OutcomeContext, query ocrtypes.Query, ao ocrtypes.AttributedObservation) error {
	_, err := withObservedExecution(ctx, p.metrics, metrics.ValidateObservation, func() (any, error) {
		err := p.ReportingPlugin.ValidateObservation(ctx, outctx, query, ao)
		return nil, err
	})
	return err
}

func (p *reportingPlugin[RI]) ObservationQuorum(ctx context.Context, outctx ocr3types.OutcomeContext, query ocrtypes.Query, aos []ocrtypes.AttributedObservation) (bool, error) {
	return withObservedExecution(ctx, p.metrics, metrics.ObservationQuorum, func() (bool, error) {
		return p.ReportingPlugin.ObservationQuorum(ctx, outctx, query, aos)
	})
}

func (p *reportingPlugin[RI]) Outcome(ctx context.Context, outctx ocr3types.OutcomeContext, query ocrtypes.Query, aos []ocrtypes.AttributedObservation) (ocr3types.Outcome, error) {
	result, err := withObservedExecution(ctx, p.metrics, metrics.Outcome, func() (ocr3types.Outcome, error) {
		return p.ReportingPlugin.Outcome(ctx, outctx, query, aos)
	})
	if err == nil {
		p.metrics.TrackSize(ctx, metrics.Outcome, len(result))
	}
	return result, err
}

func (p *reportingPlugin[RI]) Reports(ctx context.Context, seqNr uint64, outcome ocr3types.Outcome) ([]ocr3types.ReportPlus[RI], error) {
	result, err := withObservedExecution(ctx, p.metrics, metrics.Reports, func() ([]ocr3types.ReportPlus[RI], error) {
		return p.ReportingPlugin.Reports(ctx, seqNr, outcome)
	})
	p.metrics.TrackReports(ctx, metrics.Reports, len(result), err == nil)
	return result, err
}

func (p *reportingPlugin[RI]) ShouldAcceptAttestedReport(ctx context.Context, seqNr uint64, reportWithInfo ocr3types.ReportWithInfo[RI]) (bool, error) {
	result, err := withObservedExecution(ctx, p.metrics, metrics.ShouldAccept, func() (bool, error) {
		return p.ReportingPlugin.ShouldAcceptAttestedReport(ctx, seqNr, reportWithInfo)
	})
	p.metrics.TrackReports(ctx, metrics.ShouldAccept, boolToInt(result), err == nil)
	return result, err
}

func (p *reportingPlugin[RI]) ShouldTransmitAcceptedReport(ctx context.Context, seqNr uint64, reportWithInfo ocr3types.ReportWithInfo[RI]) (bool, error) {
	result, err := withObservedExecution(ctx, p.metrics, metrics.ShouldTransmit, func() (bool, error) {
		return p.ReportingPlugin.ShouldTransmitAcceptedReport(ctx, seqNr, reportWithInfo)
	})
	p.metrics.TrackReports(ctx, metrics.ShouldTransmit, boolToInt(result), err == nil)
	return result, err
}

func (p *reportingPlugin[RI]) Close() error {
	p.metrics.UpdateStatus(context.Background(), false)
	return p.ReportingPlugin.Close()
}

func boolToInt(arg bool) int {
	if arg {
		return 1
	}
	return 0
}

func withObservedExecution[R any](
	ctx context.Context,
	m *metrics.PluginMetrics,
	function metrics.FunctionType,
	exec func() (R, error),
) (R, error) {
	start := time.Now()
	result, err := exec()

	success := err == nil
	m.RecordDuration(ctx, function, time.Since(start), success)
	m.UpdateStatus(ctx, true)

	return result, err
}
