package logprovider

import (
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
)

// BlockWindow returns the start and end block for the given window.
func BlockWindow(block int64, blockRate int) (start int64, end int64) {
	windowSize := int64(blockRate)
	if windowSize == 0 {
		return block, block
	}
	start = block - (block % windowSize)
	end = start + windowSize - 1
	return
}

// LogSorter sorts the logs based on block number, tx hash and log index.
// returns true if b should come before a.
func LogSorter(a, b logpoller.Log) bool {
	return LogComparator(a, b) > 0
}

// LogComparator compares the logs based on block number, log index.
// tx hash is also checked in case the log index is not unique within a block.
//
// Returns:
//
//	-1 if a <  b
//	 0 if a == b
//	+1 if a >  b
func LogComparator(a, b logpoller.Log) int {
	if a.BlockNumber != b.BlockNumber {
		return int(a.BlockNumber - b.BlockNumber)
	}
	logIndexDiff := a.LogIndex - b.LogIndex
	if logIndexDiff == 0 {
		return a.TxHash.Big().Cmp(b.TxHash.Big())
	}
	return int(logIndexDiff)
}
