package promwrapper

import (
	"context"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3_1types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
)

var _ ocr3_1types.ReportingPlugin[any] = &reportingPlugin[any]{}

// reportingPlugin wraps an ocr3_1types.ReportingPlugin[RI] with prometheus
// instrumentation. Each method is timed; size-bearing methods (Observation,
// StateTransition precursor output, Reports) also emit byte counters.
type reportingPlugin[RI any] struct {
	ocr3_1types.ReportingPlugin[RI]
	chainFamily string
	chainID     string
	plugin      string

	configDigest     string
	reportsGenerated *prometheus.CounterVec
	durations        *prometheus.HistogramVec
	sizes            *prometheus.CounterVec
	status           *prometheus.GaugeVec
}

func newReportingPlugin[RI any](
	origin ocr3_1types.ReportingPlugin[RI],
	chainFamily string,
	chainID string,
	plugin string,
	configDigest string,
	reportsGenerated *prometheus.CounterVec,
	durations *prometheus.HistogramVec,
	sizes *prometheus.CounterVec,
	status *prometheus.GaugeVec,
) *reportingPlugin[RI] {
	return &reportingPlugin[RI]{
		ReportingPlugin:  origin,
		chainFamily:      chainFamily,
		chainID:          chainID,
		plugin:           plugin,
		configDigest:     configDigest,
		reportsGenerated: reportsGenerated,
		durations:        durations,
		sizes:            sizes,
		status:           status,
	}
}

func (p *reportingPlugin[RI]) Query(
	ctx context.Context,
	seqNr uint64,
	kv ocr3_1types.KeyValueReader,
	bbf ocr3_1types.BlobBroadcastFetcher,
) (ocrtypes.Query, error) {
	return withObservedExecution(p, query, func() (ocrtypes.Query, error) {
		return p.ReportingPlugin.Query(ctx, seqNr, kv, bbf)
	})
}

func (p *reportingPlugin[RI]) Observation(
	ctx context.Context,
	seqNr uint64,
	aq ocrtypes.AttributedQuery,
	kv ocr3_1types.KeyValueReader,
	bbf ocr3_1types.BlobBroadcastFetcher,
) (ocrtypes.Observation, error) {
	result, err := withObservedExecution(p, observation, func() (ocrtypes.Observation, error) {
		return p.ReportingPlugin.Observation(ctx, seqNr, aq, kv, bbf)
	})
	p.trackSize(observation, len(result), err)
	return result, err
}

func (p *reportingPlugin[RI]) ValidateObservation(
	ctx context.Context,
	seqNr uint64,
	aq ocrtypes.AttributedQuery,
	ao ocrtypes.AttributedObservation,
	kv ocr3_1types.KeyValueReader,
	bf ocr3_1types.BlobFetcher,
) error {
	_, err := withObservedExecution(p, validateObservation, func() (any, error) {
		err := p.ReportingPlugin.ValidateObservation(ctx, seqNr, aq, ao, kv, bf)
		return nil, err
	})
	return err
}

func (p *reportingPlugin[RI]) ObservationQuorum(
	ctx context.Context,
	seqNr uint64,
	aq ocrtypes.AttributedQuery,
	aos []ocrtypes.AttributedObservation,
	kv ocr3_1types.KeyValueReader,
	bf ocr3_1types.BlobFetcher,
) (bool, error) {
	return withObservedExecution(p, observationQuorum, func() (bool, error) {
		return p.ReportingPlugin.ObservationQuorum(ctx, seqNr, aq, aos, kv, bf)
	})
}

func (p *reportingPlugin[RI]) StateTransition(
	ctx context.Context,
	seqNr uint64,
	aq ocrtypes.AttributedQuery,
	aos []ocrtypes.AttributedObservation,
	kvrw ocr3_1types.KeyValueReadWriter,
	bf ocr3_1types.BlobFetcher,
) (ocr3_1types.ReportsPlusPrecursor, error) {
	result, err := withObservedExecution(p, stateTransition, func() (ocr3_1types.ReportsPlusPrecursor, error) {
		return p.ReportingPlugin.StateTransition(ctx, seqNr, aq, aos, kvrw, bf)
	})
	p.trackSize(stateTransition, len(result), err)
	return result, err
}

func (p *reportingPlugin[RI]) Committed(
	ctx context.Context,
	seqNr uint64,
	kv ocr3_1types.KeyValueReader,
) error {
	_, err := withObservedExecution(p, committed, func() (any, error) {
		err := p.ReportingPlugin.Committed(ctx, seqNr, kv)
		return nil, err
	})
	return err
}

func (p *reportingPlugin[RI]) Reports(
	ctx context.Context,
	seqNr uint64,
	precursor ocr3_1types.ReportsPlusPrecursor,
) ([]ocr3types.ReportPlus[RI], error) {
	result, err := withObservedExecution(p, reports, func() ([]ocr3types.ReportPlus[RI], error) {
		return p.ReportingPlugin.Reports(ctx, seqNr, precursor)
	})
	p.trackReports(reports, len(result))
	return result, err
}

func (p *reportingPlugin[RI]) ShouldAcceptAttestedReport(
	ctx context.Context,
	seqNr uint64,
	reportWithInfo ocr3types.ReportWithInfo[RI],
) (bool, error) {
	result, err := withObservedExecution(p, shouldAccept, func() (bool, error) {
		return p.ReportingPlugin.ShouldAcceptAttestedReport(ctx, seqNr, reportWithInfo)
	})
	p.trackReports(shouldAccept, boolToInt(result))
	return result, err
}

func (p *reportingPlugin[RI]) ShouldTransmitAcceptedReport(
	ctx context.Context,
	seqNr uint64,
	reportWithInfo ocr3types.ReportWithInfo[RI],
) (bool, error) {
	result, err := withObservedExecution(p, shouldTransmit, func() (bool, error) {
		return p.ReportingPlugin.ShouldTransmitAcceptedReport(ctx, seqNr, reportWithInfo)
	})
	p.trackReports(shouldTransmit, boolToInt(result))
	return result, err
}

func (p *reportingPlugin[RI]) Close() error {
	p.updateStatus(false)
	return p.ReportingPlugin.Close()
}

// ---- shared helpers ----

func (p *reportingPlugin[RI]) trackReports(function functionType, count int) {
	p.reportsGenerated.
		WithLabelValues(p.chainFamily, p.chainID, p.plugin, string(function)).
		Add(float64(count))
}

func (p *reportingPlugin[RI]) updateStatus(status bool) {
	p.status.
		WithLabelValues(p.chainFamily, p.chainID, p.plugin, p.configDigest).
		Set(float64(boolToInt(status)))
}

func (p *reportingPlugin[RI]) trackSize(function functionType, size int, err error) {
	if err != nil {
		return
	}
	p.sizes.
		WithLabelValues(p.chainFamily, p.chainID, p.plugin, string(function)).
		Add(float64(size))
}

func boolToInt(arg bool) int {
	if arg {
		return 1
	}
	return 0
}

func withObservedExecution[RI, R any](
	p *reportingPlugin[RI],
	function functionType,
	exec func() (R, error),
) (R, error) {
	start := time.Now()
	result, err := exec()

	success := err == nil

	p.durations.
		WithLabelValues(p.chainFamily, p.chainID, p.plugin, string(function), strconv.FormatBool(success)).
		Observe(float64(time.Since(start)))

	p.updateStatus(true)

	return result, err
}
