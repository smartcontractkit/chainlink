package logprovider

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDequeueCoordinator_DequeueBlockWindow(t *testing.T) {
	t.Run("an empty dequeue coordinator should tell us not to dequeue", func(t *testing.T) {
		c := NewDequeueCoordinator()

		start, end, canDequeue := c.GetDequeueBlockWindow(1, 10, 1, 10)
		assert.Equal(t, int64(0), start)
		assert.Equal(t, int64(0), end)
		assert.Equal(t, false, canDequeue)
	})

	t.Run("a populated dequeue coordinator should tell us to dequeue the first window with logs", func(t *testing.T) {
		c := NewDequeueCoordinator()

		c.CountEnqueuedLogsForWindow(3, 1, 0)
		c.CountEnqueuedLogsForWindow(4, 1, 10)
		c.CountEnqueuedLogsForWindow(5, 1, 10)

		start, end, canDequeue := c.GetDequeueBlockWindow(1, 10, 1, 10)

		// block 4 is the first block with no logs dequeued yet
		assert.Equal(t, int64(4), start)
		assert.Equal(t, int64(4), end)
		assert.Equal(t, true, canDequeue)
	})

	t.Run("a populated dequeue coordinator should tell us to dequeue the next window with logs", func(t *testing.T) {
		c := NewDequeueCoordinator()

		c.CountEnqueuedLogsForWindow(3, 1, 0)
		c.CountEnqueuedLogsForWindow(4, 1, 10)
		c.CountEnqueuedLogsForWindow(5, 1, 10)

		start, end, canDequeue := c.GetDequeueBlockWindow(1, 10, 1, 10)

		// block 4 is the first block with no logs dequeued yet
		assert.Equal(t, int64(4), start)
		assert.Equal(t, int64(4), end)
		assert.Equal(t, true, canDequeue)

		c.CountDequeuedLogsForWindow(start, 10, 10)

		start, end, canDequeue = c.GetDequeueBlockWindow(1, 10, 1, 10)

		// block 4 has been dequeued, so block 5 is the next window to dequeue
		assert.Equal(t, int64(5), start)
		assert.Equal(t, int64(5), end)
		assert.Equal(t, true, canDequeue)
	})

	t.Run("a populated dequeue coordinator with minimum dequeue met should tell us to dequeue the next window with logs as best effort", func(t *testing.T) {
		c := NewDequeueCoordinator()

		c.CountEnqueuedLogsForWindow(3, 1, 0)
		c.CountEnqueuedLogsForWindow(4, 1, 20)
		c.CountEnqueuedLogsForWindow(5, 1, 20)

		start, end, canDequeue := c.GetDequeueBlockWindow(1, 10, 1, 10)

		// block 4 is the first block with no logs dequeued yet
		assert.Equal(t, int64(4), start)
		assert.Equal(t, int64(4), end)
		assert.Equal(t, true, canDequeue)

		c.CountDequeuedLogsForWindow(start, 10, 10)

		start, end, canDequeue = c.GetDequeueBlockWindow(1, 10, 1, 10)

		// block 4 has been dequeued, so block 5 is the next window to dequeue
		assert.Equal(t, int64(5), start)
		assert.Equal(t, int64(5), end)
		assert.Equal(t, true, canDequeue)

		c.CountDequeuedLogsForWindow(start, 10, 10)

		start, end, canDequeue = c.GetDequeueBlockWindow(1, 10, 1, 10)

		// all windows have had minimum logs dequeued, so we go back to block 4 to dequeue as best effort
		assert.Equal(t, int64(4), start)
		assert.Equal(t, int64(4), end)
		assert.Equal(t, true, canDequeue)
	})

	t.Run("a fully exhausted dequeue coordinator should not tell us to dequeue", func(t *testing.T) {
		c := NewDequeueCoordinator()

		c.CountEnqueuedLogsForWindow(3, 1, 0)
		c.CountEnqueuedLogsForWindow(4, 1, 20)
		c.CountEnqueuedLogsForWindow(5, 1, 20)

		start, end, canDequeue := c.GetDequeueBlockWindow(1, 10, 1, 10)

		// block 4 is the first block with no logs dequeued yet
		assert.Equal(t, int64(4), start)
		assert.Equal(t, int64(4), end)
		assert.Equal(t, true, canDequeue)

		c.CountDequeuedLogsForWindow(start, 10, 10)

		start, end, canDequeue = c.GetDequeueBlockWindow(1, 10, 1, 10)

		// block 4 has been dequeued, so block 5 is the next window to dequeue
		assert.Equal(t, int64(5), start)
		assert.Equal(t, int64(5), end)
		assert.Equal(t, true, canDequeue)

		c.CountDequeuedLogsForWindow(start, 10, 10)

		start, end, canDequeue = c.GetDequeueBlockWindow(1, 10, 1, 10)

		// all windows have had minimum logs dequeued, so we go back to block 4 to dequeue as best effort
		assert.Equal(t, int64(4), start)
		assert.Equal(t, int64(4), end)
		assert.Equal(t, true, canDequeue)

		c.CountDequeuedLogsForWindow(start, 10, 10)

		start, end, canDequeue = c.GetDequeueBlockWindow(1, 10, 1, 10)

		// block 4 has been fully dequeued, so we dequeue block 5
		assert.Equal(t, int64(5), start)
		assert.Equal(t, int64(5), end)
		assert.Equal(t, true, canDequeue)

		c.CountDequeuedLogsForWindow(start, 10, 10)

		start, end, canDequeue = c.GetDequeueBlockWindow(1, 10, 1, 10)

		// all block windows have been fully dequeued so the coordinator tells us not to dequeue
		assert.Equal(t, int64(0), start)
		assert.Equal(t, int64(0), end)
		assert.Equal(t, false, canDequeue)
	})

	t.Run("an incomplete latest window without logs to dequeue gets passed over and best effort is executed", func(t *testing.T) {
		c := NewDequeueCoordinator()

		c.CountEnqueuedLogsForWindow(0, 4, 10)
		c.CountEnqueuedLogsForWindow(1, 4, 10)
		c.CountEnqueuedLogsForWindow(2, 4, 10)
		c.CountEnqueuedLogsForWindow(3, 4, 10)
		c.CountEnqueuedLogsForWindow(4, 4, 0)
		c.CountEnqueuedLogsForWindow(5, 4, 0)
		c.CountEnqueuedLogsForWindow(6, 4, 0)

		start, end, canDequeue := c.GetDequeueBlockWindow(1, 6, 4, 10)

		assert.Equal(t, int64(0), start)
		assert.Equal(t, int64(3), end)
		assert.Equal(t, true, canDequeue)

		c.CountDequeuedLogsForWindow(start, 10, 10)

		start, end, canDequeue = c.GetDequeueBlockWindow(1, 6, 4, 10)

		assert.Equal(t, int64(0), start)
		assert.Equal(t, int64(3), end)
		assert.Equal(t, true, canDequeue)

		// multiple dequeues in best effort now exhaust block window 0
		c.CountDequeuedLogsForWindow(start, 10, 10)
		c.CountDequeuedLogsForWindow(start, 10, 10)
		c.CountDequeuedLogsForWindow(start, 10, 10)

		start, end, canDequeue = c.GetDequeueBlockWindow(1, 6, 4, 10)

		assert.Equal(t, int64(0), start)
		assert.Equal(t, int64(0), end)
		assert.Equal(t, false, canDequeue)
	})

	t.Run("an incomplete latest window with min logs to dequeue gets dequeued", func(t *testing.T) {
		c := NewDequeueCoordinator()

		c.CountEnqueuedLogsForWindow(0, 4, 10)
		c.CountEnqueuedLogsForWindow(1, 4, 10)
		c.CountEnqueuedLogsForWindow(2, 4, 10)
		c.CountEnqueuedLogsForWindow(3, 4, 10)
		c.CountEnqueuedLogsForWindow(4, 4, 10)
		c.CountEnqueuedLogsForWindow(5, 4, 0)
		c.CountEnqueuedLogsForWindow(6, 4, 0)

		start, end, canDequeue := c.GetDequeueBlockWindow(1, 6, 4, 10)

		assert.Equal(t, int64(0), start)
		assert.Equal(t, int64(3), end)
		assert.Equal(t, true, canDequeue)

		c.CountDequeuedLogsForWindow(start, 10, 10)

		start, end, canDequeue = c.GetDequeueBlockWindow(1, 6, 4, 10)

		assert.Equal(t, int64(4), start)
		assert.Equal(t, int64(7), end)
		assert.Equal(t, true, canDequeue)

		c.CountDequeuedLogsForWindow(start, 10, 10)

		start, end, canDequeue = c.GetDequeueBlockWindow(1, 6, 4, 10)

		assert.Equal(t, int64(0), start)
		assert.Equal(t, int64(3), end)
		assert.Equal(t, true, canDequeue)
	})

	t.Run("an incomplete latest window with less than min logs to dequeue gets dequeued", func(t *testing.T) {
		c := NewDequeueCoordinator()

		c.CountEnqueuedLogsForWindow(0, 4, 10)
		c.CountEnqueuedLogsForWindow(1, 4, 10)
		c.CountEnqueuedLogsForWindow(2, 4, 10)
		c.CountEnqueuedLogsForWindow(3, 4, 10)
		c.CountEnqueuedLogsForWindow(4, 4, 5)
		c.CountEnqueuedLogsForWindow(5, 4, 0)
		c.CountEnqueuedLogsForWindow(6, 4, 0)

		start, end, canDequeue := c.GetDequeueBlockWindow(1, 6, 4, 10)

		assert.Equal(t, int64(0), start)
		assert.Equal(t, int64(3), end)
		assert.Equal(t, true, canDequeue)

		c.CountDequeuedLogsForWindow(start, 10, 10)

		start, end, canDequeue = c.GetDequeueBlockWindow(1, 6, 4, 10)

		assert.Equal(t, int64(4), start)
		assert.Equal(t, int64(7), end)
		assert.Equal(t, true, canDequeue)

		c.CountDequeuedLogsForWindow(start, 5, 10)

		start, end, canDequeue = c.GetDequeueBlockWindow(1, 6, 4, 10)

		assert.Equal(t, int64(0), start)
		assert.Equal(t, int64(3), end)
		assert.Equal(t, true, canDequeue)

		// now that the second block window is complete and has enough logs to meet min dequeue, we revert to min guaranteed dequeue
		c.CountEnqueuedLogsForWindow(7, 4, 5)

		start, end, canDequeue = c.GetDequeueBlockWindow(1, 6, 4, 10)

		assert.Equal(t, int64(4), start)
		assert.Equal(t, int64(7), end)
		assert.Equal(t, true, canDequeue)

		c.CountDequeuedLogsForWindow(start, 5, 10)

		// now we revert to best effort dequeue for the first block window
		start, end, canDequeue = c.GetDequeueBlockWindow(1, 6, 4, 10)

		assert.Equal(t, int64(0), start)
		assert.Equal(t, int64(3), end)
		assert.Equal(t, true, canDequeue)
	})

	t.Run("a reorg causes us to revert to min guaranteed log dequeue", func(t *testing.T) {
		c := NewDequeueCoordinator()

		c.CountEnqueuedLogsForWindow(3, 1, 0)
		c.CountEnqueuedLogsForWindow(4, 1, 20)
		c.CountEnqueuedLogsForWindow(5, 1, 20)
		c.CountEnqueuedLogsForWindow(6, 1, 20)
		c.CountEnqueuedLogsForWindow(7, 1, 20)

		start, end, canDequeue := c.GetDequeueBlockWindow(1, 10, 1, 10)

		// block 4 is the first block with no logs dequeued yet
		assert.Equal(t, int64(4), start)
		assert.Equal(t, int64(4), end)
		assert.Equal(t, true, canDequeue)

		c.CountDequeuedLogsForWindow(start, 10, 10)

		start, end, canDequeue = c.GetDequeueBlockWindow(1, 10, 1, 10)

		// block 4 has been dequeued, so block 5 is the next window to dequeue
		assert.Equal(t, int64(5), start)
		assert.Equal(t, int64(5), end)
		assert.Equal(t, true, canDequeue)

		c.CountDequeuedLogsForWindow(start, 10, 10)

		start, end, canDequeue = c.GetDequeueBlockWindow(1, 10, 1, 10)

		// block 5 has been dequeued, so block 6 is the next window to dequeue
		assert.Equal(t, int64(6), start)
		assert.Equal(t, int64(6), end)
		assert.Equal(t, true, canDequeue)

		c.CountDequeuedLogsForWindow(start, 10, 10)

		// reorg happens and only block 4 has been re orgd
		c.MarkReorg(4, 1)

		c.CountEnqueuedLogsForWindow(4, 1, 10)

		start, end, canDequeue = c.GetDequeueBlockWindow(1, 10, 1, 10)

		// we now have to go back to block 4 to dequeue minimum guaranteed logs
		assert.Equal(t, int64(4), start)
		assert.Equal(t, int64(4), end)
		assert.Equal(t, true, canDequeue)

		c.CountDequeuedLogsForWindow(start, 10, 10)

		start, end, canDequeue = c.GetDequeueBlockWindow(1, 10, 1, 10)

		// now that block 4 has been min dequeued, we jump forward to block 7 to continue min dequeue
		assert.Equal(t, int64(7), start)
		assert.Equal(t, int64(7), end)
		assert.Equal(t, true, canDequeue)

		c.CountDequeuedLogsForWindow(start, 10, 10)

		start, end, canDequeue = c.GetDequeueBlockWindow(1, 10, 1, 10)

		// now that all block windows have had min logs dequeued, we go back to the earliest block window with remaining logs to dequeue best effort, i.e. block window 5
		assert.Equal(t, int64(5), start)
		assert.Equal(t, int64(5), end)
		assert.Equal(t, true, canDequeue)
	})

	t.Run("cleaning deletes data from the coordinator older than the block window of block threshold", func(t *testing.T) {
		c := NewDequeueCoordinator()

		c.CountEnqueuedLogsForWindow(1, 1, 20)
		c.CountEnqueuedLogsForWindow(2, 1, 20)
		c.CountEnqueuedLogsForWindow(3, 1, 20)
		c.CountEnqueuedLogsForWindow(4, 1, 20)
		c.CountEnqueuedLogsForWindow(5, 1, 20)
		c.CountEnqueuedLogsForWindow(6, 1, 20)
		c.CountEnqueuedLogsForWindow(7, 1, 20)
		c.CountEnqueuedLogsForWindow(8, 1, 20)
		c.CountEnqueuedLogsForWindow(9, 1, 20)

		start, end, canDequeue := c.GetDequeueBlockWindow(1, 10, 1, 10)

		assert.Equal(t, int64(1), start)
		assert.Equal(t, int64(1), end)
		assert.Equal(t, true, canDequeue)

		c.CountDequeuedLogsForWindow(1, 10, 10)

		start, end, canDequeue = c.GetDequeueBlockWindow(1, 10, 1, 10)

		assert.Equal(t, int64(2), start)
		assert.Equal(t, int64(2), end)
		assert.Equal(t, true, canDequeue)

		c.CountDequeuedLogsForWindow(2, 10, 10)

		start, end, canDequeue = c.GetDequeueBlockWindow(1, 10, 1, 10)

		assert.Equal(t, int64(3), start)
		assert.Equal(t, int64(3), end)
		assert.Equal(t, true, canDequeue)

		c.CountDequeuedLogsForWindow(3, 10, 10)

		assert.Equal(t, 10, c.enqueuedLogs[1])
		assert.Equal(t, 10, c.enqueuedLogs[2])
		assert.Equal(t, 10, c.enqueuedLogs[3])

		assert.Equal(t, 10, c.dequeuedLogs[1])
		assert.Equal(t, 10, c.dequeuedLogs[2])
		assert.Equal(t, 10, c.dequeuedLogs[3])

		assert.Equal(t, true, c.dequeuedMinimum[1])
		assert.Equal(t, true, c.dequeuedMinimum[2])
		assert.Equal(t, true, c.dequeuedMinimum[3])

		assert.Equal(t, true, c.completeWindows[1])
		assert.Equal(t, true, c.completeWindows[2])
		assert.Equal(t, true, c.completeWindows[3])

		c.Clean(3, 1)

		assert.Equal(t, 0, c.enqueuedLogs[1])
		assert.Equal(t, 0, c.enqueuedLogs[2])
		assert.Equal(t, 10, c.enqueuedLogs[3])

		assert.Equal(t, 0, c.dequeuedLogs[1])
		assert.Equal(t, 0, c.dequeuedLogs[2])
		assert.Equal(t, 10, c.dequeuedLogs[3])

		assert.Equal(t, false, c.dequeuedMinimum[1])
		assert.Equal(t, false, c.dequeuedMinimum[2])
		assert.Equal(t, true, c.dequeuedMinimum[3])

		assert.Equal(t, false, c.completeWindows[1])
		assert.Equal(t, false, c.completeWindows[2])
		assert.Equal(t, true, c.completeWindows[3])
	})
}
