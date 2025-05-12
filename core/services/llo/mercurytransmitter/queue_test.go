package mercurytransmitter

import (
	"sync"
	"testing"

	heap "github.com/esote/minmaxheap"
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
	const maxSize = 7

	lggr, observedLogs := logger.TestLoggerObserved(t, zapcore.ErrorLevel)

	t.Run("cannot init with more transmissions than capacity", func(t *testing.T) {
		transmissions := makeSampleTransmissions(maxSize+1, sURL)
		tq := NewTransmitQueue(lggr, sURL, maxSize, &mockAsyncDeleter{})
		err := tq.Init(transmissions)
		require.Error(t, err)
	})

	t.Run("happy cases", func(t *testing.T) {
		testTransmissions := makeSampleTransmissions(3, sURL)
		deleter := &mockAsyncDeleter{}
		tq := NewTransmitQueue(lggr, sURL, maxSize, deleter)

		require.NoError(t, tq.Init([]*Transmission{}))

		t.Run("successfully add transmissions to transmit queue", func(t *testing.T) {
			for _, tt := range testTransmissions {
				ok := tq.Push(tt)
				require.True(t, ok)
			}
			report := tq.HealthReport()
			require.NoError(t, report[tq.Name()])
		})

		t.Run("transmit queue is more than 50% full", func(t *testing.T) {
			tq.Push(testTransmissions[2])
			report := tq.HealthReport()
			assert.Equal(t, "transmit priority queue is greater than 50% full (4/7)", report[tq.Name()].Error())
		})

		t.Run("transmit queue pops the highest priority transmission", func(t *testing.T) {
			tr := tq.BlockingPop()
			assert.Equal(t, testTransmissions[2], tr)
		})

		t.Run("transmit queue is full and evicts the oldest transmission", func(t *testing.T) {
			// add 5 more transmissions to overflow the queue by 1
			for i := 0; i < 5; i++ {
				tq.Push(testTransmissions[1])
			}

			// expecting testTransmissions[0] to get evicted and not present in the queue anymore
			testutils.WaitForLogMessage(t, observedLogs, "Transmit queue is full; dropping oldest transmission (reached max length of 7)")
			var transmissions []*Transmission
			for i := 0; i < 7; i++ {
				tr := tq.BlockingPop()
				transmissions = append(transmissions, tr)
			}

			assert.NotContains(t, transmissions, testTransmissions[0])
			require.Len(t, deleter.hashes, 1)
			assert.Equal(t, testTransmissions[0].Hash(), deleter.hashes[0])
		})

		t.Run("transmit queue blocks when empty and resumes when transmission available", func(t *testing.T) {
			assert.True(t, tq.IsEmpty())

			var wg sync.WaitGroup
			wg.Add(2)
			go func() {
				defer wg.Done()
				tr := tq.BlockingPop()
				assert.Equal(t, tr, testTransmissions[0])
			}()
			go func() {
				defer wg.Done()
				tq.Push(testTransmissions[0])
			}()
			wg.Wait()
		})

		t.Run("initializes transmissions", func(t *testing.T) {
			expected := makeSampleTransmission(1, sURL, []byte{1})
			transmissions := []*Transmission{
				expected,
			}
			tq := NewTransmitQueue(lggr, sURL, 7, deleter)
			require.NoError(t, tq.Init(transmissions))

			transmission := tq.BlockingPop()
			assert.Equal(t, expected, transmission)
			assert.True(t, tq.IsEmpty())
		})
	})

	t.Run("if the queue was overfilled it evicts entries until reaching maxSize", func(t *testing.T) {
		testTransmissions := makeSampleTransmissions(maxSize*3, sURL)
		deleter := &mockAsyncDeleter{}
		tq := NewTransmitQueue(lggr, sURL, maxSize, deleter)

		// add 3 over capacity to queue
		{
			// need to copy to avoid sorting original slice
			init := make([]*Transmission, maxSize+3)
			copy(init, testTransmissions)
			pq := priorityQueue(init)
			heap.Init(&pq)               // ensure the heap is ordered
			tq.(*transmitQueue).pq = &pq // directly assign to bypass Init check
		}

		tq.Push(testTransmissions[maxSize+3]) // push one more to trigger eviction
		require.Equal(t, maxSize, tq.(*transmitQueue).Len())
		require.Len(t, deleter.hashes, 4) // evicted overfill entries (3 oversize plus 1 more to make room)

		// oldest entries removed
		assert.Equal(t, testTransmissions[0].Hash(), deleter.hashes[0])
		assert.Equal(t, testTransmissions[1].Hash(), deleter.hashes[1])
		assert.Equal(t, testTransmissions[2].Hash(), deleter.hashes[2])
		assert.Equal(t, testTransmissions[3].Hash(), deleter.hashes[3])

		queueEntriesSorted := []*Transmission{}
		for {
			transmission := tq.(*transmitQueue).pop()
			if transmission == nil {
				break
			}
			queueEntriesSorted = append(queueEntriesSorted, transmission)
		}
		assert.ElementsMatch(t, testTransmissions[4:4+maxSize], queueEntriesSorted)
	})
}
