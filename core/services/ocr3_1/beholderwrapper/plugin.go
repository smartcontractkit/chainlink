package beholderwrapper

import (
	"context"
	"time"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3_1types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
)

var _ ocr3_1types.ReportingPlugin[any] = &reportingPlugin[any]{}

type reportingPlugin[RI any] struct {
	ocr3_1types.ReportingPlugin[RI]
	metrics *pluginMetrics
}

func newReportingPlugin[RI any](
	origin ocr3_1types.ReportingPlugin[RI],
	metrics *pluginMetrics,
) *reportingPlugin[RI] {
	return &reportingPlugin[RI]{
		ReportingPlugin: origin,
		metrics:         metrics,
	}
}

func (p *reportingPlugin[RI]) Query(ctx context.Context, seqNr uint64, keyValueReader ocr3_1types.KeyValueStateReader, blobBroadcastFetcher ocr3_1types.BlobBroadcastFetcher) (ocrtypes.Query, error) {
	return withObservedExecution(ctx, p.metrics, query, func() (ocrtypes.Query, error) {
		return p.ReportingPlugin.Query(ctx, seqNr, keyValueReader, blobBroadcastFetcher)
	})
}

func (p *reportingPlugin[RI]) Observation(ctx context.Context, seqNr uint64, aq ocrtypes.AttributedQuery, keyValueReader ocr3_1types.KeyValueStateReader, blobBroadcastFetcher ocr3_1types.BlobBroadcastFetcher) (ocrtypes.Observation, error) {
	result, err := withObservedExecution(ctx, p.metrics, observation, func() (ocrtypes.Observation, error) {
		return p.ReportingPlugin.Observation(ctx, seqNr, aq, keyValueReader, blobBroadcastFetcher)
	})
	if err == nil {
		p.metrics.trackSize(ctx, observation, len(result))
	}
	return result, err
}

func (p *reportingPlugin[RI]) ValidateObservation(ctx context.Context, seqNr uint64, aq ocrtypes.AttributedQuery, ao ocrtypes.AttributedObservation, keyValueReader ocr3_1types.KeyValueStateReader, blobFetcher ocr3_1types.BlobFetcher) error {
	_, err := withObservedExecution(ctx, p.metrics, validateObservation, func() (any, error) {
		err := p.ReportingPlugin.ValidateObservation(ctx, seqNr, aq, ao, keyValueReader, blobFetcher)
		return nil, err
	})
	return err
}

func (p *reportingPlugin[RI]) ObservationQuorum(ctx context.Context, seqNr uint64, aq ocrtypes.AttributedQuery, aos []ocrtypes.AttributedObservation, keyValueReader ocr3_1types.KeyValueStateReader, blobFetcher ocr3_1types.BlobFetcher) (bool, error) {
	return withObservedExecution(ctx, p.metrics, observationQuorum, func() (bool, error) {
		return p.ReportingPlugin.ObservationQuorum(ctx, seqNr, aq, aos, keyValueReader, blobFetcher)
	})
}

func (p *reportingPlugin[RI]) StateTransition(ctx context.Context, seqNr uint64, aq ocrtypes.AttributedQuery, aos []ocrtypes.AttributedObservation, keyValueReadWriter ocr3_1types.KeyValueStateReadWriter, blobFetcher ocr3_1types.BlobFetcher) (ocr3_1types.ReportsPlusPrecursor, error) {
	result, err := withObservedExecution(ctx, p.metrics, stateTransition, func() (ocr3_1types.ReportsPlusPrecursor, error) {
		return p.ReportingPlugin.StateTransition(ctx, seqNr, aq, aos, keyValueReadWriter, blobFetcher)
	})
	if err == nil {
		p.metrics.trackSize(ctx, stateTransition, len(result))
	}
	return result, err
}

func (p *reportingPlugin[RI]) Committed(ctx context.Context, seqNr uint64, keyValueReader ocr3_1types.KeyValueStateReader) error {
	_, err := withObservedExecution(ctx, p.metrics, committed, func() (any, error) {
		err := p.ReportingPlugin.Committed(ctx, seqNr, keyValueReader)
		return nil, err
	})
	return err
}

func (p *reportingPlugin[RI]) Reports(ctx context.Context, seqNr uint64, reportsPlusPrecursor ocr3_1types.ReportsPlusPrecursor) ([]ocr3types.ReportPlus[RI], error) {
	result, err := withObservedExecution(ctx, p.metrics, reports, func() ([]ocr3types.ReportPlus[RI], error) {
		return p.ReportingPlugin.Reports(ctx, seqNr, reportsPlusPrecursor)
	})
	p.metrics.trackReports(ctx, reports, len(result))
	return result, err
}

func (p *reportingPlugin[RI]) ShouldAcceptAttestedReport(ctx context.Context, seqNr uint64, reportWithInfo ocr3types.ReportWithInfo[RI]) (bool, error) {
	result, err := withObservedExecution(ctx, p.metrics, shouldAccept, func() (bool, error) {
		return p.ReportingPlugin.ShouldAcceptAttestedReport(ctx, seqNr, reportWithInfo)
	})
	p.metrics.trackReports(ctx, shouldAccept, boolToInt(result))
	return result, err
}

func (p *reportingPlugin[RI]) ShouldTransmitAcceptedReport(ctx context.Context, seqNr uint64, reportWithInfo ocr3types.ReportWithInfo[RI]) (bool, error) {
	result, err := withObservedExecution(ctx, p.metrics, shouldTransmit, func() (bool, error) {
		return p.ReportingPlugin.ShouldTransmitAcceptedReport(ctx, seqNr, reportWithInfo)
	})
	p.metrics.trackReports(ctx, shouldTransmit, boolToInt(result))
	return result, err
}

func (p *reportingPlugin[RI]) Close() error {
	p.metrics.updateStatus(context.Background(), false)
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
	metrics *pluginMetrics,
	function functionType,
	exec func() (R, error),
) (R, error) {
	start := time.Now()
	result, err := exec()

	success := err == nil
	metrics.recordDuration(ctx, function, time.Since(start), success)
	metrics.updateStatus(ctx, true)

	return result, err
}
