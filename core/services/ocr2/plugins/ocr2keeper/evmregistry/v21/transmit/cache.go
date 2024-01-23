package transmit

import (
	"sync"

	ocr2keepers "github.com/smartcontractkit/chainlink-common/pkg/types/automation"
)

// transmitEventCache holds a ring buffer of the last visited blocks (transmit block),
// and their corresponding logs (by log id).
// Using a ring buffer allows us to keep a cache of the last N blocks,
// without having to iterate over the entire buffer to clean it up.
type transmitEventCache struct {
	lock   sync.RWMutex
	buffer []cacheBlock

	cap int64
}

func newTransmitEventCache(cap int64) transmitEventCache {
	return transmitEventCache{
		buffer: make([]cacheBlock, cap),
		cap:    cap,
	}
}

func (c *transmitEventCache) get(block ocr2keepers.BlockNumber, logID string) (ocr2keepers.TransmitEvent, bool) {
	c.lock.RLock()
	defer c.lock.RUnlock()

	i := int64(block) % c.cap
	b := c.buffer[i]
	if b.block != block {
		return ocr2keepers.TransmitEvent{}, false
	}
	if len(b.records) == 0 {
		return ocr2keepers.TransmitEvent{}, false
	}
	e, ok := b.records[logID]

	return e, ok
}

func (c *transmitEventCache) add(logID string, e ocr2keepers.TransmitEvent) {
	c.lock.Lock()
	defer c.lock.Unlock()

	i := int64(e.TransmitBlock) % c.cap
	b := c.buffer[i]
	isBlockEmpty := len(b.records) == 0
	isNewBlock := b.block < e.TransmitBlock
	if isBlockEmpty || isNewBlock {
		b = newCacheBlock(e.TransmitBlock)
	} else if b.block > e.TransmitBlock {
		// old log
		return
	}
	b.records[logID] = e
	c.buffer[i] = b
}

type cacheBlock struct {
	block   ocr2keepers.BlockNumber
	records map[string]ocr2keepers.TransmitEvent
}

func newCacheBlock(block ocr2keepers.BlockNumber) cacheBlock {
	return cacheBlock{
		block:   block,
		records: make(map[string]ocr2keepers.TransmitEvent),
	}
}
