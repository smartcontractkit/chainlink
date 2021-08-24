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
	logsByBlockHash map[common.Hash][]types.Log
	// This min-heap maintains block numbers of logs in the pool.
	// it helps us easily determine the minimum log block number
	// in the pool (while the set of log block numbers is dynamically changing).
	heap *pairingHeap.PairHeap
}

func newLogPool() *logPool {
	return &logPool{
		hashesByBlockNumbers: make(map[uint64]map[common.Hash]struct{}),
		logsByBlockHash:      make(map[common.Hash][]types.Log),
		heap:                 pairingHeap.New(),
	}
}

func (pool *logPool) addLog(log types.Log) {
	_, exists := pool.hashesByBlockNumbers[log.BlockNumber]
	if !exists {
		pool.hashesByBlockNumbers[log.BlockNumber] = make(map[common.Hash]struct{})
	}
	pool.hashesByBlockNumbers[log.BlockNumber][log.BlockHash] = struct{}{}
	pool.logsByBlockHash[log.BlockHash] = append(pool.logsByBlockHash[log.BlockHash], log)
	pool.heap.Insert(Uint64(log.BlockNumber))
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
				logsToReturn = append(logsToReturn, logsOnBlock{blockNum, pool.logsByBlockHash[hash]})
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
			logsToReturn = append(logsToReturn, logsOnBlock{uint64(num), pool.logsByBlockHash[hash]})
		}
	}
	return logsToReturn, minBlockNumToSend
}

// deleteOlderLogs - deleting all logs for block numbers under 'keptDepth'
func (pool *logPool) deleteOlderLogs(keptDepth uint64) {
	for {
		item := pool.heap.FindMin()
		if item == nil {
			break
		}

		blockNum := uint64(item.(Uint64))
		if blockNum >= keptDepth {
			break
		}
		pool.heap.DeleteMin()

		for hash := range pool.hashesByBlockNumbers[blockNum] {
			delete(pool.logsByBlockHash, hash)
		}
		delete(pool.hashesByBlockNumbers, blockNum)
	}
}

func (pool *logPool) removeLog(log types.Log) {
	// deleting all logs for this log's block hash
	delete(pool.logsByBlockHash, log.BlockHash)
	delete(pool.hashesByBlockNumbers[log.BlockNumber], log.BlockHash)
	if len(pool.hashesByBlockNumbers[log.BlockNumber]) == 0 {
		delete(pool.hashesByBlockNumbers, log.BlockNumber)
	}
}

type Uint64 int

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
