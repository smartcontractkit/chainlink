package beholderwrapper

import (
	"context"
	"time"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3_1types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink/v2/core/services/ocr3/beholderwrapper/metrics"
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

func (p *reportingPlugin[RI]) wrapReader(ctx context.Context, r ocr3_1types.KeyValueStateReader) ocr3_1types.KeyValueStateReader {
	if r == nil {
		return nil
	}
	return &instrumentedKVStateReader{inner: r, ctx: ctx, metrics: p.metrics}
}

func (p *reportingPlugin[RI]) wrapReadWriter(ctx context.Context, rw ocr3_1types.KeyValueStateReadWriter) ocr3_1types.KeyValueStateReadWriter {
	if rw == nil {
		return nil
	}
	return &instrumentedKVStateReadWriter{
		instrumentedKVStateReader: instrumentedKVStateReader{inner: rw, ctx: ctx, metrics: p.metrics},
		writer:                    rw,
	}
}

func (p *reportingPlugin[RI]) wrapBroadcastFetcher(bbf ocr3_1types.BlobBroadcastFetcher) ocr3_1types.BlobBroadcastFetcher {
	if bbf == nil {
		return nil
	}
	return &instrumentedBlobBroadcastFetcher{
		inner:                   bbf,
		metrics:                 p.metrics,
		instrumentedBlobFetcher: instrumentedBlobFetcher{inner: bbf, metrics: p.metrics},
	}
}

func (p *reportingPlugin[RI]) wrapFetcher(bf ocr3_1types.BlobFetcher) ocr3_1types.BlobFetcher {
	if bf == nil {
		return nil
	}
	return &instrumentedBlobFetcher{inner: bf, metrics: p.metrics}
}

func (p *reportingPlugin[RI]) Query(ctx context.Context, seqNr uint64, keyValueReader ocr3_1types.KeyValueStateReader, blobBroadcastFetcher ocr3_1types.BlobBroadcastFetcher) (ocrtypes.Query, error) {
	return withObservedExecution(ctx, p.metrics, metrics.Query, func() (ocrtypes.Query, error) {
		return p.ReportingPlugin.Query(ctx, seqNr, p.wrapReader(ctx, keyValueReader), p.wrapBroadcastFetcher(blobBroadcastFetcher))
	})
}

func (p *reportingPlugin[RI]) Observation(ctx context.Context, seqNr uint64, aq ocrtypes.AttributedQuery, keyValueReader ocr3_1types.KeyValueStateReader, blobBroadcastFetcher ocr3_1types.BlobBroadcastFetcher) (ocrtypes.Observation, error) {
	result, err := withObservedExecution(ctx, p.metrics, metrics.Observation, func() (ocrtypes.Observation, error) {
		return p.ReportingPlugin.Observation(ctx, seqNr, aq, p.wrapReader(ctx, keyValueReader), p.wrapBroadcastFetcher(blobBroadcastFetcher))
	})
	if err == nil {
		p.metrics.TrackSize(ctx, metrics.Observation, len(result))
	}
	return result, err
}

func (p *reportingPlugin[RI]) ValidateObservation(ctx context.Context, seqNr uint64, aq ocrtypes.AttributedQuery, ao ocrtypes.AttributedObservation, keyValueReader ocr3_1types.KeyValueStateReader, blobFetcher ocr3_1types.BlobFetcher) error {
	_, err := withObservedExecution(ctx, p.metrics, metrics.ValidateObservation, func() (any, error) {
		err := p.ReportingPlugin.ValidateObservation(ctx, seqNr, aq, ao, p.wrapReader(ctx, keyValueReader), p.wrapFetcher(blobFetcher))
		return nil, err
	})
	return err
}

func (p *reportingPlugin[RI]) ObservationQuorum(ctx context.Context, seqNr uint64, aq ocrtypes.AttributedQuery, aos []ocrtypes.AttributedObservation, keyValueReader ocr3_1types.KeyValueStateReader, blobFetcher ocr3_1types.BlobFetcher) (bool, error) {
	return withObservedExecution(ctx, p.metrics, metrics.ObservationQuorum, func() (bool, error) {
		return p.ReportingPlugin.ObservationQuorum(ctx, seqNr, aq, aos, p.wrapReader(ctx, keyValueReader), p.wrapFetcher(blobFetcher))
	})
}

func (p *reportingPlugin[RI]) StateTransition(ctx context.Context, seqNr uint64, aq ocrtypes.AttributedQuery, aos []ocrtypes.AttributedObservation, keyValueReadWriter ocr3_1types.KeyValueStateReadWriter, blobFetcher ocr3_1types.BlobFetcher) (ocr3_1types.ReportsPlusPrecursor, error) {
	result, err := withObservedExecution(ctx, p.metrics, metrics.StateTransition, func() (ocr3_1types.ReportsPlusPrecursor, error) {
		return p.ReportingPlugin.StateTransition(ctx, seqNr, aq, aos, p.wrapReadWriter(ctx, keyValueReadWriter), p.wrapFetcher(blobFetcher))
	})
	if err == nil {
		p.metrics.TrackSize(ctx, metrics.StateTransition, len(result))
	}
	return result, err
}

func (p *reportingPlugin[RI]) Committed(ctx context.Context, seqNr uint64, keyValueReader ocr3_1types.KeyValueStateReader) error {
	_, err := withObservedExecution(ctx, p.metrics, metrics.Committed, func() (any, error) {
		err := p.ReportingPlugin.Committed(ctx, seqNr, p.wrapReader(ctx, keyValueReader))
		return nil, err
	})
	return err
}

func (p *reportingPlugin[RI]) Reports(ctx context.Context, seqNr uint64, reportsPlusPrecursor ocr3_1types.ReportsPlusPrecursor) ([]ocr3types.ReportPlus[RI], error) {
	result, err := withObservedExecution(ctx, p.metrics, metrics.Reports, func() ([]ocr3types.ReportPlus[RI], error) {
		return p.ReportingPlugin.Reports(ctx, seqNr, reportsPlusPrecursor)
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
	m *pluginMetrics,
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
