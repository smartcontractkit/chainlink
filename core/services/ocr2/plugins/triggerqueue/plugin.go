package triggerqueue

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3_1types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	v2 "github.com/smartcontractkit/chainlink/v2/core/services/workflows/v2"
)

// errNotImplemented is returned by all plugin methods in this draft.
var errNotImplemented = errors.New("triggerqueue plugin: draft implementation, not yet implemented")

var _ ocr3_1types.ReportingPlugin[[]byte] = (*ReportingPlugin)(nil)

// ReportingPlugin implements OCR 3.1 ReportingPlugin for the trigger queue.
type ReportingPlugin struct {
	lggr   logger.Logger
	buffer *v2.ObservationBuffer[v2.EnqueuedTriggerEvent]
}

// NewReportingPlugin creates a new ReportingPlugin.
func NewReportingPlugin(lggr logger.Logger, buffer *v2.ObservationBuffer[v2.EnqueuedTriggerEvent]) *ReportingPlugin {
	return &ReportingPlugin{lggr: lggr.Named("TriggerQueuePlugin"), buffer: buffer}
}

func (p *ReportingPlugin) Query(ctx context.Context, seqNr uint64, keyValueReader ocr3_1types.KeyValueStateReader, blobBroadcastFetcher ocr3_1types.BlobBroadcastFetcher) (types.Query, error) {
	return nil, errNotImplemented
}

// Observation reads from the buffer (filled by OCRQueue.Put) and produces an observation.
// Draft: returns minimal observation (event IDs and Lamport as JSON); full impl would BroadcastBlob for payloads.
func (p *ReportingPlugin) Observation(ctx context.Context, seqNr uint64, aq types.AttributedQuery, keyValueReader ocr3_1types.KeyValueStateReader, blobBroadcastFetcher ocr3_1types.BlobBroadcastFetcher) (types.Observation, error) {
	events := p.buffer.TakeForObservation()
	if len(events) == 0 {
		return []byte("[]"), nil
	}
	obs := make([]struct {
		ID      string `json:"id"`
		Lamport uint64 `json:"lamport"`
	}, len(events))
	for i, be := range events {
		obs[i] = struct {
			ID      string `json:"id"`
			Lamport uint64 `json:"lamport"`
		}{ID: be.ID(), Lamport: be.Lamport()}
	}
	return json.Marshal(obs)
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
