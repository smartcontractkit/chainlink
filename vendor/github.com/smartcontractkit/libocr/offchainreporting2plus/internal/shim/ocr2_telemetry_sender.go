package shim

import (
	"time"

	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/internal/loghelper"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/internal/ocr2/serialization"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
)

type OCR2TelemetrySender struct {
	chTelemetry chan<- *serialization.TelemetryWrapper
	logger      commontypes.Logger
	taper       loghelper.LogarithmicTaper
}

func MakeOCR2TelemetrySender(chTelemetry chan<- *serialization.TelemetryWrapper, logger commontypes.Logger) OCR2TelemetrySender {
	return OCR2TelemetrySender{chTelemetry, logger, loghelper.LogarithmicTaper{}}
}

func (ts OCR2TelemetrySender) send(t *serialization.TelemetryWrapper) {
	select {
	case ts.chTelemetry <- t:
		ts.taper.Reset(func(oldCount uint64) {
			ts.logger.Info("OCR2TelemetrySender: stopped dropping telemetry", commontypes.LogFields{
				"droppedCount": oldCount,
			})
		})
	default:
		ts.taper.Trigger(func(newCount uint64) {
			ts.logger.Warn("OCR2TelemetrySender: dropping telemetry", commontypes.LogFields{
				"droppedCount": newCount,
			})
		})
	}
}

func (ts OCR2TelemetrySender) RoundStarted(
	configDigest types.ConfigDigest,
	epoch uint32,
	round uint8,
	leader commontypes.OracleID,
) {
	ts.send(&serialization.TelemetryWrapper{
		Wrapped: &serialization.TelemetryWrapper_RoundStarted{&serialization.TelemetryRoundStarted{
			ConfigDigest: configDigest[:],
			Epoch:        uint64(epoch),
			Round:        uint64(round),
			Leader:       uint64(leader),
			Time:         uint64(time.Now().UnixNano()),
		}},
		UnixTimeNanoseconds: time.Now().UnixNano(),
	})
}
