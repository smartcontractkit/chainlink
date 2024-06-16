package logprovider

import (
	"sync"
)

type DequeueCoordinator interface {
	// CountEnqueuedLogsForWindow tracks how many logs are added for a particular block during the enqueue process.
	CountEnqueuedLogsForWindow(block int64, blockRate uint32, added int)
	// GetDequeueBlockWindow identifies a block window ready for processing between the given start and latest block numbers.
	// It prioritizes windows that need to have the minimum guaranteed logs dequeued before considering windows with
	// remaining logs to be dequeued, as a best effort.
	GetDequeueBlockWindow(start int64, latestBlock int64, blockRate int, minGuarantee int) (int64, int64, bool)
	// CountDequeuedLogsForWindow updates the status of a block window based on the number of logs dequeued,
	// remaining logs, and the number of upkeeps. This function tracks remaining and dequeued logs for the specified
	// block window, determines if a block window has had the minimum number of guaranteed logs dequeued, and marks a
	// window as not ready if there are not yet any logs available to dequeue from the window.
	CountDequeuedLogsForWindow(startWindow int64, logs, minGuaranteedLogs int)
	// MarkReorg handles the detection of a reorg  by resetting the state of the affected block window. It ensures that
	// upkeeps within the specified block window are marked as not having the minimum number of guaranteed logs dequeued.
	MarkReorg(block int64, blockRate uint32)
	// Clean removes any data that is older than the block window of the blockThreshold from the dequeueCoordinator
	Clean(blockThreshold int64, blockRate uint32)
}

type dequeueCoordinator struct {
	dequeuedMinimum map[int64]bool
	enqueuedLogs    map[int64]int
	dequeuedLogs    map[int64]int
	completeWindows map[int64]bool
	mu              sync.Mutex
}

func NewDequeueCoordinator() *dequeueCoordinator {
	return &dequeueCoordinator{
		dequeuedMinimum: map[int64]bool{},
		enqueuedLogs:    map[int64]int{},
		dequeuedLogs:    map[int64]int{},
		completeWindows: map[int64]bool{},
	}
}

func (c *dequeueCoordinator) CountEnqueuedLogsForWindow(block int64, blockRate uint32, added int) {
	c.mu.Lock()
	defer c.mu.Unlock()

	startWindow, _ := getBlockWindow(block, int(blockRate))
	c.enqueuedLogs[startWindow] += added
}

func (c *dequeueCoordinator) GetDequeueBlockWindow(start int64, latestBlock int64, blockRate int, minGuarantee int) (int64, int64, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// check if minimum logs have been dequeued
	for i := start; i <= latestBlock; i += int64(blockRate) {
		startWindow, end := getBlockWindow(i, blockRate)
		if latestBlock >= end {
			c.completeWindows[startWindow] = true
		} else if c.enqueuedLogs[startWindow] <= 0 { // the latest window is incomplete and has no logs to provide yet
			break
		}

		enqueuedLogs := c.enqueuedLogs[startWindow]
		dequeuedLogs := c.dequeuedLogs[startWindow]

		if enqueuedLogs > 0 && dequeuedLogs < minGuarantee {
			return startWindow, end, true
		}
	}

	// check best effort dequeue
	for i := start; i < latestBlock; i += int64(blockRate) {
		startWindow, end := getBlockWindow(i, blockRate)

		if remainingLogs, ok := c.enqueuedLogs[startWindow]; ok && remainingLogs > 0 {
			return startWindow, end, true
		}
	}

	return 0, 0, false
}

func (c *dequeueCoordinator) CountDequeuedLogsForWindow(startWindow int64, logs, minGuaranteedLogs int) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.enqueuedLogs[startWindow] -= logs
	c.dequeuedLogs[startWindow] += logs

	if c.completeWindows[startWindow] {
		if c.enqueuedLogs[startWindow] <= 0 || c.dequeuedLogs[startWindow] >= minGuaranteedLogs {
			// if the window is complete, and there are no more logs, then we have to consider this as min dequeued, even if no logs were dequeued
			c.dequeuedMinimum[startWindow] = true
		}
	} else if c.dequeuedLogs[startWindow] >= minGuaranteedLogs {
		// if the window is not complete, but we were able to dequeue min guaranteed logs from the blocks that were available
		c.dequeuedMinimum[startWindow] = true
	}
}

func (c *dequeueCoordinator) MarkReorg(block int64, blockRate uint32) {
	c.mu.Lock()
	defer c.mu.Unlock()

	startWindow, _ := getBlockWindow(block, int(blockRate))
	c.dequeuedMinimum[startWindow] = false
	c.enqueuedLogs[startWindow] = 0
	c.dequeuedLogs[startWindow] = 0
}

func (c *dequeueCoordinator) Clean(blockThreshold int64, blockRate uint32) {
	c.mu.Lock()
	defer c.mu.Unlock()

	blockThresholdStartWindow, _ := getBlockWindow(blockThreshold, int(blockRate))

	for block := range c.enqueuedLogs {
		if blockThresholdStartWindow > block {
			delete(c.enqueuedLogs, block)
			delete(c.dequeuedLogs, block)
			delete(c.dequeuedMinimum, block)
			delete(c.completeWindows, block)
		}
	}
}
