package telem

import (
	"crypto/sha256"
	"fmt"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink-data-streams/llo"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/proto"

	"github.com/smartcontractkit/chainlink/v2/core/services/synchronization"
)

func TestFingerprint(t *testing.T) {
	ot := time.Now()
	bytes32 := sha256.Sum256([]byte(ot.String()))
	configDigest := bytes32[:]
	donID := uint32(2)
	streamID := uint32(123)
	channelID := uint32(234)

	tests := []struct {
		name        string
		msg         proto.Message
		typ         synchronization.TelemetryType
		fingerprint string
		ts          int32
		err         error
	}{
		{
			name: "successful observation",
			msg: &LLOObservationTelemetry{
				DonId:                donID,
				StreamId:             streamID,
				ConfigDigest:         configDigest,
				ObservationTimestamp: ot.UnixNano(),
			},
			typ:         synchronization.LLOObservation,
			fingerprint: fmt.Sprintf("%d-%d-%x", donID, streamID, configDigest),
			ts:          int32(ot.Unix()),
			err:         nil,
		},
		{
			name: "successful outcome",
			msg: &llo.LLOOutcomeTelemetry{
				DonId:                           donID,
				ConfigDigest:                    configDigest,
				ObservationTimestampNanoseconds: uint64(ot.UnixNano()),
			},
			typ:         synchronization.LLOOutcome,
			fingerprint: fmt.Sprintf("%d-%x", donID, configDigest),
			ts:          int32(ot.Unix()),
			err:         nil,
		},
		{
			name: "successful report",
			msg: &llo.LLOReportTelemetry{
				DonId:                           donID,
				ChannelId:                       channelID,
				ConfigDigest:                    configDigest,
				ObservationTimestampNanoseconds: uint64(ot.UnixNano()),
			},
			typ:         synchronization.LLOReport,
			fingerprint: fmt.Sprintf("%d-%d-%x", donID, channelID, configDigest),
			ts:          int32(ot.Unix()),
			err:         nil,
		},
		{
			name: "successful bridge",
			msg: &LLOBridgeTelemetry{
				DonId:                donID,
				StreamId:             &streamID,
				SpecId:               345,
				BridgeAdapterName:    "bridge-adapter",
				ConfigDigest:         configDigest,
				ObservationTimestamp: ot.UnixNano(),
			},
			typ:         synchronization.PipelineBridge,
			fingerprint: fmt.Sprintf("%d-%d-%d-%s-%x", donID, streamID, 345, "bridge-adapter", configDigest),
			ts:          int32(ot.Unix()),
			err:         nil,
		},
		{
			name: "unsupported telemetry type",
			typ:  synchronization.HeadReport,
			err:  errUnsupportedTelemetryType,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			fp, ts, err := fingerprint(test.typ, test.msg)
			if test.err != nil {
				assert.EqualError(t, err, test.err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.fingerprint, fp)
				assert.Equal(t, test.ts, ts)
			}
		})
	}
}
