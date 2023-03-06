package shim

import (
	"time"

	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/internal/loghelper"
	"github.com/smartcontractkit/libocr/offchainreporting2/internal/serialization"
	"github.com/smartcontractkit/libocr/offchainreporting2/types"
)

type TelemetrySender struct {
	chTelemetry chan<- *serialization.TelemetryWrapper
	logger      commontypes.Logger
	taper       loghelper.LogarithmicTaper
}

func MakeTelemetrySender(chTelemetry chan<- *serialization.TelemetryWrapper, logger commontypes.Logger) TelemetrySender {
	return TelemetrySender{chTelemetry, logger, loghelper.LogarithmicTaper{}}
}

func (ts TelemetrySender) send(t *serialization.TelemetryWrapper) {
	select {
	case ts.chTelemetry <- t:
		ts.taper.Reset(func(oldCount uint64) {
			ts.logger.Info("TelemetrySender: stopped dropping telemetry", commontypes.LogFields{
				"droppedCount": oldCount,
			})
		})
	default:
		ts.taper.Trigger(func(newCount uint64) {
			ts.logger.Warn("TelemetrySender: dropping telemetry", commontypes.LogFields{
				"droppedCount": newCount,
			})
		})
	}
}

func (ts TelemetrySender) RoundStarted(
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
