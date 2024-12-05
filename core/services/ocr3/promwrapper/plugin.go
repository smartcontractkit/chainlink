package promwrapper

import (
	"context"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
)

var _ ocr3types.ReportingPlugin[any] = &ReportingPlugin[any]{}

type ReportingPlugin[RI any] struct {
	ocr3types.ReportingPlugin[RI]
	chainID string
	plugin  string

	// Prometheus components for tracking metrics
	reportsGenerated *prometheus.CounterVec
	durations        *prometheus.HistogramVec
}

func NewReportingPlugin[RI any](
	origin ocr3types.ReportingPlugin[RI],
	chainID string,
	plugin string,
	reportsGenerated *prometheus.CounterVec,
	durations *prometheus.HistogramVec,
) *ReportingPlugin[RI] {
	return &ReportingPlugin[RI]{
		ReportingPlugin:  origin,
		chainID:          chainID,
		plugin:           plugin,
		reportsGenerated: reportsGenerated,
		durations:        durations,
	}
}

func (p *ReportingPlugin[RI]) Query(ctx context.Context, outctx ocr3types.OutcomeContext) (ocrtypes.Query, error) {
	return withObservedExecution(p, query, func() (ocrtypes.Query, error) {
		return p.ReportingPlugin.Query(ctx, outctx)
	})
}

func (p *ReportingPlugin[RI]) Observation(ctx context.Context, outctx ocr3types.OutcomeContext, query ocrtypes.Query) (ocrtypes.Observation, error) {
	return withObservedExecution(p, observation, func() (ocrtypes.Observation, error) {
		return p.ReportingPlugin.Observation(ctx, outctx, query)
	})
}

func (p *ReportingPlugin[RI]) ValidateObservation(ctx context.Context, outctx ocr3types.OutcomeContext, query ocrtypes.Query, ao ocrtypes.AttributedObservation) error {
	_, err := withObservedExecution(p, validateObservation, func() (any, error) {
		err := p.ReportingPlugin.ValidateObservation(ctx, outctx, query, ao)
		return nil, err
	})
	return err
}

func (p *ReportingPlugin[RI]) Outcome(ctx context.Context, outctx ocr3types.OutcomeContext, query ocrtypes.Query, aos []ocrtypes.AttributedObservation) (ocr3types.Outcome, error) {
	return withObservedExecution(p, outcome, func() (ocr3types.Outcome, error) {
		return p.ReportingPlugin.Outcome(ctx, outctx, query, aos)
	})
}

func (p *ReportingPlugin[RI]) Reports(ctx context.Context, seqNr uint64, outcome ocr3types.Outcome) ([]ocr3types.ReportPlus[RI], error) {
	result, err := withObservedExecution(p, reports, func() ([]ocr3types.ReportPlus[RI], error) {
		return p.ReportingPlugin.Reports(ctx, seqNr, outcome)
	})
	p.trackReports(reports, len(result))
	return result, err
}

func (p *ReportingPlugin[RI]) ShouldAcceptAttestedReport(ctx context.Context, seqNr uint64, reportWithInfo ocr3types.ReportWithInfo[RI]) (bool, error) {
	result, err := withObservedExecution(p, shouldAccept, func() (bool, error) {
		return p.ReportingPlugin.ShouldAcceptAttestedReport(ctx, seqNr, reportWithInfo)
	})
	p.trackReports(shouldAccept, boolToInt(result))
	return result, err
}

func (p *ReportingPlugin[RI]) ShouldTransmitAcceptedReport(ctx context.Context, seqNr uint64, reportWithInfo ocr3types.ReportWithInfo[RI]) (bool, error) {
	result, err := withObservedExecution(p, shouldTransmit, func() (bool, error) {
		return p.ReportingPlugin.ShouldTransmitAcceptedReport(ctx, seqNr, reportWithInfo)
	})
	p.trackReports(shouldTransmit, boolToInt(result))
	return result, err
}

func (p *ReportingPlugin[RI]) Close() error {
	return p.ReportingPlugin.Close()
}

func (p *ReportingPlugin[RI]) trackReports(
	function functionType,
	count int,
) {
	p.reportsGenerated.
		WithLabelValues(p.chainID, p.plugin, string(function)).
		Add(float64(count))
}

func boolToInt(arg bool) int {
	if arg {
		return 1
	}
	return 0
}

func withObservedExecution[RI, R any](
	p *ReportingPlugin[RI],
	function functionType,
	exec func() (R, error),
) (R, error) {
	start := time.Now()
	result, err := exec()

	success := err == nil

	p.durations.
		WithLabelValues(p.chainID, p.plugin, string(function), strconv.FormatBool(success)).
		Observe(float64(time.Since(start)))

	return result, err
}
