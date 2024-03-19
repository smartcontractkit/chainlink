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
	end = block + (windowSize - (block % windowSize) - 1)
	return
}

// LogSorter sorts the logs based on block number, tx hash and log index.
// returns true if b should come before a.
func LogSorter(a, b logpoller.Log) bool {
	return LogComparator(a, b) > 0
}

// LogComparator compares the logs based on block number, tx hash and log index.
//
// Returns:
//
//	-1 if a <  b
//	 0 if a == b
//	+1 if a >  b
func LogComparator(a, b logpoller.Log) int {
	if b.BlockNumber != a.BlockNumber {
		return int(a.BlockNumber - b.BlockNumber)
	}
	if txDiff := a.TxHash.Big().Cmp(b.TxHash.Big()); txDiff != 0 {
		return txDiff
	}
	return int(a.LogIndex - b.LogIndex)
}
