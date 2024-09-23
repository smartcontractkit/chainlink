package mercurytransmitter

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

var _ asyncDeleter = &mockAsyncDeleter{}

type mockAsyncDeleter struct {
	donID  uint32
	hashes [][32]byte
}

func (m *mockAsyncDeleter) AsyncDelete(hash [32]byte) {
	m.hashes = append(m.hashes, hash)
}
func (m *mockAsyncDeleter) DonID() uint32 {
	return m.donID
}

func Test_Queue(t *testing.T) {
	t.Parallel()
	lggr, observedLogs := logger.TestLoggerObserved(t, zapcore.ErrorLevel)
	testTransmissions := makeSampleTransmissions()
	deleter := &mockAsyncDeleter{}
	transmitQueue := NewTransmitQueue(lggr, sURL, 7, deleter)
	transmitQueue.Init([]*Transmission{})

	t.Run("successfully add transmissions to transmit queue", func(t *testing.T) {
		for _, tt := range testTransmissions {
			ok := transmitQueue.Push(tt)
			require.True(t, ok)
		}
		report := transmitQueue.HealthReport()
		assert.Nil(t, report[transmitQueue.Name()])
	})

	t.Run("transmit queue is more than 50% full", func(t *testing.T) {
		transmitQueue.Push(testTransmissions[2])
		report := transmitQueue.HealthReport()
		assert.Equal(t, report[transmitQueue.Name()].Error(), "transmit priority queue is greater than 50% full (4/7)")
	})

	t.Run("transmit queue pops the highest priority transmission", func(t *testing.T) {
		tr := transmitQueue.BlockingPop()
		assert.Equal(t, testTransmissions[2], tr)
	})

	t.Run("transmit queue is full and evicts the oldest transmission", func(t *testing.T) {
		// add 5 more transmissions to overflow the queue by 1
		for i := 0; i < 5; i++ {
			transmitQueue.Push(testTransmissions[1])
		}

		// expecting testTransmissions[0] to get evicted and not present in the queue anymore
		testutils.WaitForLogMessage(t, observedLogs, "Transmit queue is full; dropping oldest transmission (reached max length of 7)")
		var transmissions []*Transmission
		for i := 0; i < 7; i++ {
			tr := transmitQueue.BlockingPop()
			transmissions = append(transmissions, tr)
		}

		assert.NotContains(t, transmissions, testTransmissions[0])
		require.Len(t, deleter.hashes, 1)
		assert.Equal(t, testTransmissions[0].Hash(), deleter.hashes[0])
	})

	t.Run("transmit queue blocks when empty and resumes when tranmission available", func(t *testing.T) {
		assert.True(t, transmitQueue.IsEmpty())

		var wg sync.WaitGroup
		wg.Add(2)
		go func() {
			defer wg.Done()
			tr := transmitQueue.BlockingPop()
			assert.Equal(t, tr, testTransmissions[0])
		}()
		go func() {
			defer wg.Done()
			transmitQueue.Push(testTransmissions[0])
		}()
		wg.Wait()
	})

	t.Run("initializes transmissions", func(t *testing.T) {
		expected := makeSampleTransmission(1)
		transmissions := []*Transmission{
			expected,
		}
		transmitQueue := NewTransmitQueue(lggr, sURL, 7, deleter)
		transmitQueue.Init(transmissions)

		transmission := transmitQueue.BlockingPop()
		assert.Equal(t, expected, transmission)
		assert.True(t, transmitQueue.IsEmpty())
	})
}
