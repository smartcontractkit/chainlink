package logprovider

import (
	"encoding/hex"

	ocr2keepers "github.com/smartcontractkit/chainlink-common/pkg/types/automation"

	"github.com/smartcontractkit/chainlink-evm/pkg/logpoller"
)

// LogSorter sorts the logs primarily by block number, then by log index, and finally by tx hash.
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
	blockDiff := int(a.BlockNumber - b.BlockNumber)
	if blockDiff != 0 {
		return normalizeCompareResult(blockDiff)
	}
	logIndexDiff := int(a.LogIndex - b.LogIndex)
	if logIndexDiff != 0 {
		return normalizeCompareResult(logIndexDiff)
	}
	return a.TxHash.Big().Cmp(b.TxHash.Big())
}

// normalizeCompareResult normalizes the result of a comparison to -1, 0, 1
func normalizeCompareResult(res int) int {
	switch {
	case res < 0:
		return -1
	case res > 0:
		return 1
	default:
		return 0
	}
}

// logID returns a unique identifier for a log, which is an hex string
// of ocr2keepers.LogTriggerExtension.LogIdentifier()
func logID(l logpoller.Log) string {
	ext := ocr2keepers.LogTriggerExtension{
		Index: uint32(l.LogIndex),
	}
	copy(ext.TxHash[:], l.TxHash[:])
	copy(ext.BlockHash[:], l.BlockHash[:])
	return hex.EncodeToString(ext.LogIdentifier())
}
