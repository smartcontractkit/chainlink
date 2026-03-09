package triggerqueue

import (
	"context"
	"errors"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3_1types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

// errNotImplemented is returned by all plugin methods in this draft.
var errNotImplemented = errors.New("triggerqueue plugin: draft implementation, not yet implemented")

var _ ocr3_1types.ReportingPlugin[[]byte] = (*ReportingPlugin)(nil)

// ReportingPlugin implements OCR 3.1 ReportingPlugin for the trigger queue.
// Draft: all methods return errors.
type ReportingPlugin struct {
	lggr logger.Logger
}

// NewReportingPlugin creates a new ReportingPlugin. Draft: returns plugin that errors on all calls.
func NewReportingPlugin(lggr logger.Logger) *ReportingPlugin {
	return &ReportingPlugin{lggr: lggr.Named("TriggerQueuePlugin")}
}

func (p *ReportingPlugin) Query(ctx context.Context, seqNr uint64, keyValueReader ocr3_1types.KeyValueStateReader, blobBroadcastFetcher ocr3_1types.BlobBroadcastFetcher) (types.Query, error) {
	return nil, errNotImplemented
}

func (p *ReportingPlugin) Observation(ctx context.Context, seqNr uint64, aq types.AttributedQuery, keyValueReader ocr3_1types.KeyValueStateReader, blobBroadcastFetcher ocr3_1types.BlobBroadcastFetcher) (types.Observation, error) {
	return nil, errNotImplemented
}

func (p *ReportingPlugin) ValidateObservation(ctx context.Context, seqNr uint64, aq types.AttributedQuery, ao types.AttributedObservation, keyValueReader ocr3_1types.KeyValueStateReader, blobFetcher ocr3_1types.BlobFetcher) error {
	return errNotImplemented
}

func (p *ReportingPlugin) ObservationQuorum(ctx context.Context, seqNr uint64, aq types.AttributedQuery, aos []types.AttributedObservation, keyValueReader ocr3_1types.KeyValueStateReader, blobFetcher ocr3_1types.BlobFetcher) (bool, error) {
	return false, errNotImplemented
}

func (p *ReportingPlugin) StateTransition(ctx context.Context, seqNr uint64, aq types.AttributedQuery, aos []types.AttributedObservation, keyValueReadWriter ocr3_1types.KeyValueStateReadWriter, blobFetcher ocr3_1types.BlobFetcher) (ocr3_1types.ReportsPlusPrecursor, error) {
	return ocr3_1types.ReportsPlusPrecursor{}, errNotImplemented
}

func (p *ReportingPlugin) Reports(ctx context.Context, seqNr uint64, reportsPlusPrecursor ocr3_1types.ReportsPlusPrecursor) ([]ocr3types.ReportPlus[[]byte], error) {
	return nil, errNotImplemented
}

func (p *ReportingPlugin) Committed(ctx context.Context, seqNr uint64, keyValueReader ocr3_1types.KeyValueStateReader) error {
	return errNotImplemented
}

func (p *ReportingPlugin) ShouldAcceptAttestedReport(ctx context.Context, seqNr uint64, reportWithInfo ocr3types.ReportWithInfo[[]byte]) (bool, error) {
	return false, errNotImplemented
}

func (p *ReportingPlugin) ShouldTransmitAcceptedReport(ctx context.Context, seqNr uint64, reportWithInfo ocr3types.ReportWithInfo[[]byte]) (bool, error) {
	return false, errNotImplemented
}

func (p *ReportingPlugin) Close() error {
	return nil
}
