package transmit

import (
	"encoding/hex"
	"sync"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	ocr2keepers "github.com/smartcontractkit/ocr2keepers/pkg/v3/types"
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
		buffer: make([]cacheBlock, cap+1),
		cap:    cap,
	}
}

func (c *transmitEventCache) get(logID string) (ocr2keepers.TransmitEvent, bool) {
	c.lock.RLock()
	defer c.lock.RUnlock()

	for i := range c.buffer {
		b := c.buffer[i]
		if len(b.records) == 0 {
			continue
		}
		e, ok := b.records[logID]
		if ok {
			return e, true
		}
	}
	return ocr2keepers.TransmitEvent{}, false
}

func (c *transmitEventCache) add(logID string, e ocr2keepers.TransmitEvent) {
	c.lock.Lock()
	defer c.lock.Unlock()

	i := int64(e.TransmitBlock) % int64(c.cap)
	b := c.buffer[i]
	isBlockEmpty := len(b.records) == 0
	isNewBlock := b.block < e.TransmitBlock
	if isBlockEmpty || isNewBlock {
		b = newCacheBlock()
	}
	b.records[logID] = e
	c.buffer[i] = b
}

func logKey(log logpoller.Log) string {
	logExt := ocr2keepers.LogTriggerExtension{
		TxHash: log.TxHash,
		Index:  uint32(log.LogIndex),
	}
	logId := logExt.LogIdentifier()
	return hex.EncodeToString(logId)
}

type cacheBlock struct {
	block   ocr2keepers.BlockNumber
	records map[string]ocr2keepers.TransmitEvent
}

func newCacheBlock() cacheBlock {
	return cacheBlock{
		records: make(map[string]ocr2keepers.TransmitEvent),
	}
}
