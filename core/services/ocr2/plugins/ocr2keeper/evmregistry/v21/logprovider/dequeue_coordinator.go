package logprovider

import (
	"sync"
)

type DequeueCoordinator interface {
	// GetDequeueBlockWindow identifies a block window ready for processing between the given start and latest block numbers.
	// It prioritizes windows that need to have the minimum guaranteed logs dequeued before considering windows with
	// remaining logs to be dequeued, as a best effort.
	GetDequeueBlockWindow(start int64, latestBlock int64, blockRate int, minGuarantee int) (int64, int64, map[string]int, bool)
	// CountDequeuedLogsForWindow updates the status of a block window based on the number of logs dequeued,
	// remaining logs, and the number of upkeeps. This function tracks remaining and dequeued logs for the specified
	// block window, determines if a block window has had the minimum number of guaranteed logs dequeued, and marks a
	// window as not ready if there are not yet any logs available to dequeue from the window.
	CountDequeuedLogsForWindow(startWindow int64, logs []BufferedLog, minGuaranteedLogs int)
	// MarkReorg handles the detection of a reorg  by resetting the state of the affected block window. It ensures that
	// upkeeps within the specified block window are marked as not having the minimum number of guaranteed logs dequeued.
	MarkReorg(block int64, blockRate uint32)
	// Clean removes any data that is older than the block window of the blockThreshold from the dequeueCoordinator
	Clean(blockThreshold int64, blockRate uint32)
	// Sync updates the state of the coordinator based on the queue states
	Sync(map[string]*upkeepLogQueue, uint32)
}

type dequeueCoordinator struct {
	dequeuedMinimum                map[int64]bool
	dequeuedMinimumUpkeepsByWindow map[int64]map[string]bool
	enqueuedLogs                   map[int64]int
	enqueuedUpkeepLogs             map[int64]map[string]int
	dequeuedLogs                   map[int64]int
	dequeuedUpkeepLogs             map[int64]map[string]int
	completeWindows                map[int64]bool
	mu                             sync.Mutex
}

func NewDequeueCoordinator() *dequeueCoordinator {
	return &dequeueCoordinator{
		dequeuedMinimum:                map[int64]bool{},
		dequeuedMinimumUpkeepsByWindow: map[int64]map[string]bool{},
		enqueuedLogs:                   map[int64]int{},
		enqueuedUpkeepLogs:             map[int64]map[string]int{},
		dequeuedLogs:                   map[int64]int{},
		dequeuedUpkeepLogs:             map[int64]map[string]int{},
		completeWindows:                map[int64]bool{},
	}
}

func (c *dequeueCoordinator) Sync(queues map[string]*upkeepLogQueue, blockRate uint32) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.enqueuedLogs = map[int64]int{}
	c.enqueuedUpkeepLogs = map[int64]map[string]int{}

	for uid, queue := range queues {
		for blockNumber, logs := range queue.logs {
			startWindow, _ := getBlockWindow(blockNumber, int(blockRate))

			c.enqueuedLogs[startWindow] += len(logs)

			if _, ok := c.enqueuedUpkeepLogs[startWindow]; ok {
				c.enqueuedUpkeepLogs[startWindow][uid] += len(logs)
			} else {
				c.enqueuedUpkeepLogs[startWindow] = map[string]int{
					uid: len(logs),
				}
			}
		}
	}
}

func (c *dequeueCoordinator) GetDequeueBlockWindow(start int64, latestBlock int64, blockRate int, logLimitLow int) (int64, int64, map[string]int, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	upkeepIDs := map[string]int{}

	startBlockWindow, _ := getBlockWindow(start, blockRate)

	// check if minimum logs have been dequeued
	for blockWindow := startBlockWindow; blockWindow <= latestBlock; blockWindow += int64(blockRate) {

		if c.dequeuedMinimum[blockWindow] || len(c.enqueuedUpkeepLogs[blockWindow]) == 0 {
			continue
		}

		// find the upkeep IDs in this window that need min dequeue
		for upkeepID, remainingLogCount := range c.enqueuedUpkeepLogs[blockWindow] {
			if remainingLogCount == 0 {
				continue
			}

			var dequeuedLogs int
			if windowDequeues, ok := c.dequeuedUpkeepLogs[blockWindow]; ok {
				dequeuedLogs = windowDequeues[upkeepID]
			}

			if dequeuedLogs >= logLimitLow {
				continue
			}

			upkeepIDs[upkeepID] = logLimitLow - dequeuedLogs
		}

		startWindow, end := getBlockWindow(blockWindow, blockRate)
		if latestBlock >= end {
			c.completeWindows[startWindow] = true
		} else if c.enqueuedLogs[startWindow] <= 0 { // the latest window is incomplete and has no logs to provide yet
			break
		}

		if len(upkeepIDs) > 0 {
			return startWindow, end, upkeepIDs, true
		}
	}

	// check best effort dequeue
	for i := start; i < latestBlock; i += int64(blockRate) {
		startWindow, end := getBlockWindow(i, blockRate)

		if remainingLogs, ok := c.enqueuedLogs[startWindow]; ok && remainingLogs > 0 {
			return startWindow, end, nil, true
		}
	}

	return 0, 0, nil, false
}

func (c *dequeueCoordinator) CountDequeuedLogsForWindow(startWindow int64, logs []BufferedLog, logLimitLow int) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.enqueuedLogs[startWindow] -= len(logs)
	c.dequeuedLogs[startWindow] += len(logs)

	if dequeuedCounts, ok := c.dequeuedUpkeepLogs[startWindow]; ok {
		for _, log := range logs {
			if count, ok := dequeuedCounts[log.ID.String()]; ok {
				dequeuedCounts[log.ID.String()] = count + 1
			} else {
				dequeuedCounts[log.ID.String()] = 1
			}
		}
	} else {
		newDequeuedCounts := map[string]int{}

		for _, log := range logs {
			if count, ok := newDequeuedCounts[log.ID.String()]; ok {
				newDequeuedCounts[log.ID.String()] = count + 1
			} else {
				newDequeuedCounts[log.ID.String()] = 1
			}
		}

		c.dequeuedUpkeepLogs[startWindow] = newDequeuedCounts
	}

	for _, dequeueCounts := range c.dequeuedUpkeepLogs {
		for uid, count := range dequeueCounts {
			if count >= logLimitLow {
				if completeUpkeeps, ok := c.dequeuedMinimumUpkeepsByWindow[startWindow]; ok {
					completeUpkeeps[uid] = true
					c.dequeuedMinimumUpkeepsByWindow[startWindow] = completeUpkeeps
				} else {
					c.dequeuedMinimumUpkeepsByWindow[startWindow] = map[string]bool{
						uid: true,
					}
				}
			}
		}
	}

	if c.completeWindows[startWindow] {
		// if all upkeeps in this window have had min dequeue met then the whole window has had min dequeue met
		if len(c.dequeuedMinimumUpkeepsByWindow[startWindow]) == len(c.enqueuedUpkeepLogs[startWindow]) {
			c.dequeuedMinimum[startWindow] = true
		} else if c.enqueuedLogs[startWindow] <= 0 {
			// if the window is complete, and there are no more logs, then we have to consider this as min dequeued, even if no logs were dequeued
			c.dequeuedMinimum[startWindow] = true
		}
	} else {
		// if the window is not complete, but we were able to dequeue min guaranteed logs from the blocks that were available
		if len(c.dequeuedMinimumUpkeepsByWindow[startWindow]) == len(c.enqueuedUpkeepLogs[startWindow]) {
			c.dequeuedMinimum[startWindow] = true
		}
	}
}

func (c *dequeueCoordinator) MarkReorg(block int64, blockRate uint32) {
	c.mu.Lock()
	defer c.mu.Unlock()

	startWindow, _ := getBlockWindow(block, int(blockRate))
	c.dequeuedMinimum[startWindow] = false
	c.enqueuedLogs[startWindow] = 0
	c.dequeuedLogs[startWindow] = 0
	c.dequeuedUpkeepLogs[startWindow] = map[string]int{}
	c.enqueuedUpkeepLogs[startWindow] = map[string]int{}
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
			delete(c.dequeuedUpkeepLogs, block)
			delete(c.enqueuedUpkeepLogs, block)
		}
	}
}
