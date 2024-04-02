package logprovider

import (
	"encoding/hex"

	ocr2keepers "github.com/smartcontractkit/chainlink-common/pkg/types/automation"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
)

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

// latestBlockNumber returns the latest block number from the given logs
func latestBlockNumber(logs ...logpoller.Log) int64 {
	var latest int64
	for _, l := range logs {
		if l.BlockNumber > latest {
			latest = l.BlockNumber
		}
	}
	return latest
}
