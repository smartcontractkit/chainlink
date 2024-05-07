package logprovider

import (
	"encoding/hex"

	"github.com/ethereum/go-ethereum/common"

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

// latestBlockNumber returns the latest block number from the given logs
func (b *logBuffer) latestBlockNumber(logs ...logpoller.Log) (int64, map[int64]bool) {
	var latest int64
	var latestBlockHash common.Hash

	reorgBlocks := map[int64]bool{}

	for _, l := range logs {
		b.blockHashes[l.BlockNumber] = l.BlockHash.String()
		if l.BlockNumber > latest {
			latest = l.BlockNumber
			latestBlockHash = l.BlockHash
		}
	}

	subscriberLatest := b.latestBlockHash.Load()

	// if we see a reorg, update the stored hashes for the reorg blocks, and collect the reorg block numbers
	// so that we can later evict logs for those block numbers
	if subscriberLatest.String() != latestBlockHash.String() {
		for _, block := range b.history {
			number := int64(block.Number)
			historyBlockHash := common.Hash(block.Hash).String()

			if hash, ok := b.blockHashes[number]; ok && hash != historyBlockHash {
				b.blockHashes[number] = historyBlockHash
				reorgBlocks[number] = true
			}
		}
	}

	return latest, reorgBlocks
}
