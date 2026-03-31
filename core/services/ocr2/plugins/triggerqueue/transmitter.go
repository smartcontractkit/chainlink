package triggerqueue

import (
	"context"
	"time"

	"github.com/smartcontractkit/libocr/offchainreporting2/types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	v2 "github.com/smartcontractkit/chainlink/v2/core/services/workflows/v2"
)

var _ ocr3types.ContractTransmitter[[]byte] = (*Transmitter)(nil)

// Transmitter receives OCR consensus reports and delivers decoded trigger events via callback.
// No on-chain transmit; reports are consumed locally and routed to the ConsensusEventReceiver.
type Transmitter struct {
	receiver v2.ObserverFunc[v2.EnqueuedTriggerEvent]
	lggr     logger.Logger
}

// NewTransmitter creates a transmitter that decodes reports and calls receiver.OnConsensusEvent.
func NewTransmitter(receiver v2.ObserverFunc[v2.EnqueuedTriggerEvent], lggr logger.Logger) *Transmitter {
	return &Transmitter{receiver: receiver, lggr: lggr.Named("TriggerQueueTransmitter")}
}

// decodedTriggerEvent is the result of decoding a report
type decodedTriggerEvent struct {
	workflowID   string // hex-encoded, for dispatcher routing
	triggerCapID string
	triggerIndex int
	timestamp    time.Time
	event        capabilities.TriggerResponse
}

// decodeReport extracts the consensus slice of events from the OCR report.
//
// TODO: make interface dependency to mock in tests
var decodeReport = func(rwi ocr3types.ReportWithInfo[[]byte]) ([]decodedTriggerEvent, error) {
	// TODO: implement real decode.
	// - Parse rwi.Report (proto: ordered event IDs)
	// - For each event ID, fetch payload from KV via "Event::"+eventID
	// - Unmarshal to TriggerResponse, build decodedTriggerEvent
	_ = rwi
	return nil, nil
}

// Transmit decodes the report after consensus is reached and
// enqueues each event into the internal queue. Engine Wait() will return the head
// of these events.
func (t *Transmitter) Transmit(ctx context.Context, cd types.ConfigDigest, seqNr uint64, rwi ocr3types.ReportWithInfo[[]byte], sigs []types.AttributedOnchainSignature) error {
	events, err := decodeReport(rwi)
	if err != nil {
		t.lggr.Errorw("Failed to decode trigger queue report", "err", err, "seqNr", seqNr)
		return err
	}
	for _, ev := range events {
		enqueued := v2.NewEnqueuedTriggerEvent(ev.workflowID, ev.triggerCapID, ev.triggerIndex, ev.timestamp, ev.event)
		t.receiver(ctx, enqueued)
		t.lggr.Debugw("Delivered consensus event", "triggerCapID", ev.triggerCapID, "triggerIndex", ev.triggerIndex)
	}
	return nil
}

func (t *Transmitter) FromAccount(_ context.Context) (types.Account, error) {
	return types.Account(""), nil
}
