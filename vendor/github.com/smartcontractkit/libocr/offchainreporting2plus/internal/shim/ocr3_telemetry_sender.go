package shim

import (
	"time"

	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/internal/loghelper"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/internal/ocr3/serialization"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
)

type OCR3TelemetrySender struct {
	chTelemetry chan<- *serialization.TelemetryWrapper
	logger      commontypes.Logger
	taper       loghelper.LogarithmicTaper
}

func MakeOCR3TelemetrySender(chTelemetry chan<- *serialization.TelemetryWrapper, logger commontypes.Logger) OCR3TelemetrySender {
	return OCR3TelemetrySender{chTelemetry, logger, loghelper.LogarithmicTaper{}}
}

func (ts OCR3TelemetrySender) send(t *serialization.TelemetryWrapper) {
	select {
	case ts.chTelemetry <- t:
		ts.taper.Reset(func(oldCount uint64) {
			ts.logger.Info("OCR3TelemetrySender: stopped dropping telemetry", commontypes.LogFields{
				"droppedCount": oldCount,
			})
		})
	default:
		ts.taper.Trigger(func(newCount uint64) {
			ts.logger.Warn("OCR3TelemetrySender: dropping telemetry", commontypes.LogFields{
				"droppedCount": newCount,
			})
		})
	}
}

func (ts OCR3TelemetrySender) RoundStarted(
	configDigest types.ConfigDigest,
	epoch uint64,
	seqNr uint64,
	round uint64,
	leader commontypes.OracleID,
) {
	ts.send(&serialization.TelemetryWrapper{
		Wrapped: &serialization.TelemetryWrapper_RoundStarted{&serialization.TelemetryRoundStarted{
			ConfigDigest: configDigest[:],
			Epoch:        epoch,
			Round:        round,
			Leader:       uint64(leader),
			Time:         uint64(time.Now().UnixNano()),
			SeqNr:        seqNr,
		}},
		UnixTimeNanoseconds: time.Now().UnixNano(),
	})
}
