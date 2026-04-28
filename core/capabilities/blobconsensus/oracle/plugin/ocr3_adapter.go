package plugin

import (
	"context"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3_1types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
)

// ocr3ReportingPluginAdapter wraps an OCR 3.1 reporting plugin so it can be used
// with OCR3 managed oracle and chainlink-common's ocr3types.ReportingPluginFactory
// contract, using an in-memory KeyValueState seeded from OutcomeContext.PreviousOutcome.
type ocr3ReportingPluginAdapter struct {
	inner *reportingPlugin
}

// AsOCR3ReportingPlugin wraps p for use with OCR3 orchestration APIs.
func AsOCR3ReportingPlugin(p *reportingPlugin) ocr3types.ReportingPlugin[[]byte] {
	if p == nil {
		return nil
	}
	return &ocr3ReportingPluginAdapter{inner: p}
}

func (a *ocr3ReportingPluginAdapter) Query(ctx context.Context, outctx ocr3types.OutcomeContext) (types.Query, error) {
	kv := newMemoryKVFromPreviousOutcome(outctx.PreviousOutcome)
	return a.inner.Query(ctx, outctx.SeqNr, kv, defaultNoopBlobBroadcastFetcher)
}

func (a *ocr3ReportingPluginAdapter) Observation(ctx context.Context, outctx ocr3types.OutcomeContext, query types.Query) (types.Observation, error) {
	kv := newMemoryKVFromPreviousOutcome(outctx.PreviousOutcome)
	aq := types.AttributedQuery{Query: query, Proposer: 0}
	return a.inner.Observation(ctx, outctx.SeqNr, aq, kv, defaultNoopBlobBroadcastFetcher)
}

func (a *ocr3ReportingPluginAdapter) ValidateObservation(ctx context.Context, outctx ocr3types.OutcomeContext, query types.Query, ao types.AttributedObservation) error {
	kv := newMemoryKVFromPreviousOutcome(outctx.PreviousOutcome)
	aq := types.AttributedQuery{Query: query, Proposer: 0}
	return a.inner.ValidateObservation(ctx, outctx.SeqNr, aq, ao, kv, defaultNoopBlobBroadcastFetcher)
}

func (a *ocr3ReportingPluginAdapter) ObservationQuorum(ctx context.Context, outctx ocr3types.OutcomeContext, query types.Query, aos []types.AttributedObservation) (bool, error) {
	kv := newMemoryKVFromPreviousOutcome(outctx.PreviousOutcome)
	aq := types.AttributedQuery{Query: query, Proposer: 0}
	return a.inner.ObservationQuorum(ctx, outctx.SeqNr, aq, aos, kv, defaultNoopBlobBroadcastFetcher)
}

func (a *ocr3ReportingPluginAdapter) Outcome(ctx context.Context, outctx ocr3types.OutcomeContext, query types.Query, attributedObservations []types.AttributedObservation) (ocr3types.Outcome, error) {
	kv := newMemoryKVFromPreviousOutcome(outctx.PreviousOutcome)
	aq := types.AttributedQuery{Query: query, Proposer: 0}
	precursor, err := a.inner.StateTransition(ctx, outctx.SeqNr, aq, attributedObservations, kv, defaultNoopBlobBroadcastFetcher)
	if err != nil {
		return nil, err
	}
	if err := a.inner.Committed(ctx, outctx.SeqNr, kv); err != nil {
		return nil, err
	}
	return ocr3types.Outcome(precursor), nil
}

func (a *ocr3ReportingPluginAdapter) Reports(ctx context.Context, seqNr uint64, outcome ocr3types.Outcome) ([]ocr3types.ReportPlus[[]byte], error) {
	return a.inner.Reports(ctx, seqNr, ocr3_1types.ReportsPlusPrecursor(outcome))
}

func (a *ocr3ReportingPluginAdapter) ShouldAcceptAttestedReport(ctx context.Context, seqNr uint64, rwi ocr3types.ReportWithInfo[[]byte]) (bool, error) {
	return a.inner.ShouldAcceptAttestedReport(ctx, seqNr, rwi)
}

func (a *ocr3ReportingPluginAdapter) ShouldTransmitAcceptedReport(ctx context.Context, seqNr uint64, rwi ocr3types.ReportWithInfo[[]byte]) (bool, error) {
	return a.inner.ShouldTransmitAcceptedReport(ctx, seqNr, rwi)
}

func (a *ocr3ReportingPluginAdapter) Close() error {
	return a.inner.Close()
}
