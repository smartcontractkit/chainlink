package mercury

import (
	"sync"
	"testing"

	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury/wsrpc/pb"
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
	testTransmissions := createTestTransmissions(t)
	deleter := mocks.NewAsyncDeleter(t)
	transmitQueue := NewTransmitQueue(logger.Test(t), sURL, "foo feed ID", 7, deleter)
	transmitQueue.Init([]*Transmission{})

	t.Run("successfully add transmissions to transmit queue", func(t *testing.T) {
		for _, tt := range testTransmissions {
			ok := transmitQueue.Push(tt.tr, tt.ctx)
			require.True(t, ok)
		}
		report := transmitQueue.HealthReport()
		assert.Nil(t, report[transmitQueue.Name()])
	})

	t.Run("transmit queue is more than 50% full", func(t *testing.T) {
		transmitQueue.Push(testTransmissions[2].tr, testTransmissions[2].ctx)
		report := transmitQueue.HealthReport()
		assert.Equal(t, report[transmitQueue.Name()].Error(), "transmit priority queue is greater than 50% full (4/7)")
	})

	t.Run("transmit queue pops the highest priority transmission", func(t *testing.T) {
		tr := transmitQueue.BlockingPop()
		assert.Equal(t, testTransmissions[2].tr, tr.Req)
	})

	t.Run("transmit queue is full and evicts the oldest transmission", func(t *testing.T) {
		deleter.On("AsyncDelete", testTransmissions[0].tr).Once()

		// add 5 more transmissions to overflow the queue by 1
		for range 5 {
			transmitQueue.Push(testTransmissions[1].tr, testTransmissions[1].ctx)
		}

		// expecting testTransmissions[0] to get evicted, processed by deleter and not present in the queue anymore
		for range 7 {
			tr := transmitQueue.BlockingPop()
			assert.NotEqual(t, tr.Req, testTransmissions[0].tr)
		}
	})

	t.Run("transmit queue blocks when empty and resumes when transmission available", func(t *testing.T) {
		assert.True(t, transmitQueue.IsEmpty())

		var wg sync.WaitGroup
		wg.Add(2)
		go func() {
			defer wg.Done()
			tr := transmitQueue.BlockingPop()
			assert.Equal(t, tr.Req, testTransmissions[0].tr)
		}()
		go func() {
			defer wg.Done()
			transmitQueue.Push(testTransmissions[0].tr, testTransmissions[0].ctx)
		}()
		wg.Wait()
	})

	t.Run("initializes transmissions", func(t *testing.T) {
		transmissions := []*Transmission{
			{
				Req: &pb.TransmitRequest{
					Payload: []byte("new1"),
				},
				ReportCtx: ocrtypes.ReportContext{
					ReportTimestamp: ocrtypes.ReportTimestamp{
						Epoch:        1,
						Round:        1,
						ConfigDigest: ocrtypes.ConfigDigest{},
					},
				},
			},
		}
		transmitQueue := NewTransmitQueue(logger.Test(t), sURL, "foo feed ID", 7, deleter)
		transmitQueue.Init(transmissions)

		transmission := transmitQueue.BlockingPop()
		assert.Equal(t, transmission.Req.Payload, []byte("new1"))
		assert.True(t, transmitQueue.IsEmpty())
	})
}
