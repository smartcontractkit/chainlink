package vrfcommon

import (
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// pruneInterval is the interval in blocks at which to prune old data from the delivered set.
const pruneInterval = 100

func NewLogDeduper(lookback int) *LogDeduper {
	return &LogDeduper{
		delivered: make(map[logKey]struct{}),
		lookback:  lookback,
	}
}

// LogDeduper prevents duplicate logs from being reprocessed.
type LogDeduper struct {
	// delivered is the set of logs within the lookback that have already been delivered.
	delivered map[logKey]struct{}

	// lookback defines how long state should be kept for. Logs included in blocks older than
	// lookback may or may not be redelivered.
	lookback int

	// lastPruneHeight is the blockheight at which logs were last pruned.
	lastPruneHeight uint64

	// mu synchronizes access to the delivered map.
	mu sync.Mutex
}

// logKey represents uniquely identifying information for a single log broadcast.
type logKey struct {

	// blockHash of the block the log was included in.
	blockHash common.Hash

	// blockNumber of the block the log was included in. This is necessary to prune old logs.
	blockNumber uint64

	// logIndex of the log in the block.
	logIndex uint
}

func (l *LogDeduper) ShouldDeliver(log types.Log) bool {
	l.mu.Lock()
	defer l.mu.Unlock() // unlock in the last defer, so that we hold the lock when pruning.
	defer l.Prune(log.BlockNumber)

	key := logKey{
		blockHash:   log.BlockHash,
		blockNumber: log.BlockNumber,
		logIndex:    log.Index,
	}

	if _, ok := l.delivered[key]; ok {
		return false
	}

	l.delivered[key] = struct{}{}
	return true
}

func (l *LogDeduper) Prune(logBlock uint64) {
	// Only prune every pruneInterval blocks
	if int(logBlock)-int(l.lastPruneHeight) < pruneInterval {
		return
	}

	for key := range l.delivered {
		if int(key.blockNumber) < int(logBlock)-l.lookback {
			delete(l.delivered, key)
		}
	}

	l.lastPruneHeight = logBlock
}

// Clear clears the log deduper's internal cache.
func (l *LogDeduper) Clear() {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.delivered = make(map[logKey]struct{})
}
