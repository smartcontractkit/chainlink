package mercury

import (
	"testing"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury/wsrpc/pb"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type TestTransmissionWithReport struct {
	tr  *pb.TransmitRequest
	ctx ocrtypes.ReportContext
}

func createTestTransmissions(t *testing.T) []TestTransmissionWithReport {
	t.Helper()
	return []TestTransmissionWithReport{
		{
			tr: &pb.TransmitRequest{
				Payload: []byte("test1"),
			},
			ctx: ocrtypes.ReportContext{
				ReportTimestamp: ocrtypes.ReportTimestamp{
					Epoch:        1,
					Round:        1,
					ConfigDigest: ocrtypes.ConfigDigest{},
				},
			},
		},
		{
			tr: &pb.TransmitRequest{
				Payload: []byte("test2"),
			},
			ctx: ocrtypes.ReportContext{
				ReportTimestamp: ocrtypes.ReportTimestamp{
					Epoch:        2,
					Round:        2,
					ConfigDigest: ocrtypes.ConfigDigest{},
				},
			},
		},
		{
			tr: &pb.TransmitRequest{
				Payload: []byte("test3"),
			},
			ctx: ocrtypes.ReportContext{
				ReportTimestamp: ocrtypes.ReportTimestamp{
					Epoch:        3,
					Round:        3,
					ConfigDigest: ocrtypes.ConfigDigest{},
				},
			},
		},
	}
}

func Test_Queue(t *testing.T) {
	t.Parallel()
	lggr := logger.TestLogger(t)
	testTransmissions := createTestTransmissions(t)
	transmitQueue := NewTransmitQueue(lggr, 7)

	t.Run("successfully add transmissions to transmit queue", func(t *testing.T) {
		for _, tt := range testTransmissions {
			ok := transmitQueue.Push(tt.tr, tt.ctx)
			require.True(t, ok)
		}
		report := transmitQueue.HealthReport()
		assert.Nil(t, report[transmitQueue.Name()])
	})

	t.Run("transmit queue is more than 50% full", func(t *testing.T) {
		transmitQueue.Push(testTransmissions[0].tr, testTransmissions[0].ctx)
		report := transmitQueue.HealthReport()
		assert.Equal(t, report[transmitQueue.Name()].Error(), "transmit priority queue is greater than 50% full (4/7)")
	})

}
