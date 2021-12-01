package log

import (
	"math"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	heaps "github.com/theodesp/go-heaps"
	pairingHeap "github.com/theodesp/go-heaps/pairing"
)

type logPool struct {
	// A mapping of block numbers to a set of block hashes for all
	// the logs in the pool.
	hashesByBlockNumbers map[uint64]map[common.Hash]struct{}
	// A mapping of blockhashes to logs
	logsByBlockHash map[common.Hash]map[uint]types.Log
	// This min-heap maintains block numbers of logs in the pool.
	// it helps us easily determine the minimum log block number
	// in the pool (while the set of log block numbers is dynamically changing).
	heap *pairingHeap.PairHeap
}

func newLogPool() *logPool {
	return &logPool{
		hashesByBlockNumbers: make(map[uint64]map[common.Hash]struct{}),
		logsByBlockHash:      make(map[common.Hash]map[uint]types.Log),
		heap:                 pairingHeap.New(),
	}
}

// addLog adds log to the pool and returns true if its block number is a new minimum.
func (pool *logPool) addLog(log types.Log) bool {
	_, exists := pool.hashesByBlockNumbers[log.BlockNumber]
	if !exists {
		pool.hashesByBlockNumbers[log.BlockNumber] = make(map[common.Hash]struct{})
	}
	pool.hashesByBlockNumbers[log.BlockNumber][log.BlockHash] = struct{}{}
	if _, exists := pool.logsByBlockHash[log.BlockHash]; !exists {
		pool.logsByBlockHash[log.BlockHash] = make(map[uint]types.Log)
	}
	pool.logsByBlockHash[log.BlockHash][log.Index] = log
	min := pool.heap.FindMin()
	pool.heap.Insert(Uint64(log.BlockNumber))
	// first or new min
	return min == nil || log.BlockNumber < uint64(min.(Uint64))
}

func (pool *logPool) getAndDeleteAll() ([]logsOnBlock, int64, int64) {
	logsToReturn := make([]logsOnBlock, 0)
	lowest := int64(math.MaxInt64)
	highest := int64(0)

	for {
		item := pool.heap.DeleteMin()
		if item == nil {
			break
		}

		blockNum := uint64(item.(Uint64))
		hashes, exists := pool.hashesByBlockNumbers[blockNum]
		if exists {
			if int64(blockNum) < lowest {
				lowest = int64(blockNum)
			}
			if int64(blockNum) > highest {
				highest = int64(blockNum)
			}
			for hash := range hashes {
				logsToReturn = append(logsToReturn, newLogsOnBlock(blockNum, pool.logsByBlockHash[hash]))
				delete(pool.hashesByBlockNumbers[blockNum], hash)
				delete(pool.logsByBlockHash, hash)
			}
		}

		delete(pool.hashesByBlockNumbers, blockNum)
	}
	return logsToReturn, lowest, highest
}

func (pool *logPool) getLogsToSend(latestBlockNum int64) ([]logsOnBlock, int64) {
	logsToReturn := make([]logsOnBlock, 0)

	// gathering logs to return - from min block number kept, to latestBlockNum
	minBlockNumToSendItem := pool.heap.FindMin()
	if minBlockNumToSendItem == nil {
		return logsToReturn, 0
	}
	minBlockNumToSend := int64(minBlockNumToSendItem.(Uint64))

	for num := minBlockNumToSend; num <= latestBlockNum; num++ {
		for hash := range pool.hashesByBlockNumbers[uint64(num)] {
			logsToReturn = append(logsToReturn, newLogsOnBlock(uint64(num), pool.logsByBlockHash[hash]))
		}
	}
	return logsToReturn, minBlockNumToSend
}

// deleteOlderLogs - deleting all logs for block numbers under 'keptDepth'
func (pool *logPool) deleteOlderLogs(keptDepth int64) *int64 {
	min := pool.heap.FindMin
	for item := min(); item != nil; item = min() {
		blockNum := uint64(item.(Uint64))
		if i := int64(blockNum); i >= keptDepth {
			return &i
		}
		pool.heap.DeleteMin()

		for hash := range pool.hashesByBlockNumbers[blockNum] {
			delete(pool.logsByBlockHash, hash)
		}
		delete(pool.hashesByBlockNumbers, blockNum)
	}
	return nil
}

func (pool *logPool) removeLog(log types.Log) {
	// deleting all logs for this log's block hash
	delete(pool.logsByBlockHash, log.BlockHash)
	delete(pool.hashesByBlockNumbers[log.BlockNumber], log.BlockHash)
	if len(pool.hashesByBlockNumbers[log.BlockNumber]) == 0 {
		delete(pool.hashesByBlockNumbers, log.BlockNumber)
	}
}

type Uint64 uint64

func (a Uint64) Compare(b heaps.Item) int {
	a1 := a
	a2 := b.(Uint64)
	switch {
	case a1 > a2:
		return 1
	case a1 < a2:
		return -1
	default:
		return 0
	}
}

type logsOnBlock struct {
	BlockNumber uint64
	Logs        []types.Log
}

func newLogsOnBlock(num uint64, logsMap map[uint]types.Log) logsOnBlock {
	logs := make([]types.Log, 0, len(logsMap))
	for _, l := range logsMap {
		logs = append(logs, l)
	}
	return logsOnBlock{num, logs}
}
