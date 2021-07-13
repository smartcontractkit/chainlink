package log

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	heaps "github.com/theodesp/go-heaps"
	pairingHeap "github.com/theodesp/go-heaps/pairing"
)

type (
	logPool struct {
		hashesByBlockNumbers map[uint64][]common.Hash
		logsByBlockHash      map[common.Hash][]types.Log
		heap                 *pairingHeap.PairHeap
	}
)

func newLogPool() *logPool {
	return &logPool{
		hashesByBlockNumbers: make(map[uint64][]common.Hash),
		logsByBlockHash:      make(map[common.Hash][]types.Log),
		heap:                 pairingHeap.New(),
	}
}

func (pool *logPool) addLog(log types.Log) {
	pool.hashesByBlockNumbers[log.BlockNumber] = append(pool.hashesByBlockNumbers[log.BlockNumber], log.BlockHash)
	pool.logsByBlockHash[log.BlockHash] = append(pool.logsByBlockHash[log.BlockHash], log)
	pool.heap.Insert(Uint64(log.BlockNumber))
}

func (pool *logPool) getLogsToSend(latestBlockNum int64) ([]types.Log, int64) {
	logsToReturn := make([]types.Log, 0)

	// gathering logs to return - from min block number kept, to latestBlockNum
	minBlockNumToSendItem := pool.heap.FindMin()
	if minBlockNumToSendItem == nil {
		return logsToReturn, 0
	}
	minBlockNumToSend := int64(minBlockNumToSendItem.(Uint64))

	for num := minBlockNumToSend; num <= latestBlockNum; num++ {

		for _, hash := range pool.hashesByBlockNumbers[uint64(num)] {
			logsToReturn = append(logsToReturn, pool.logsByBlockHash[hash]...)
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

		for _, hash := range pool.hashesByBlockNumbers[blockNum] {
			delete(pool.logsByBlockHash, hash)
		}
		delete(pool.hashesByBlockNumbers, blockNum)
	}
}

func (pool *logPool) removeLog(log types.Log) {
	// deleting all logs for this log's block hash
	delete(pool.logsByBlockHash, log.BlockHash)

	for i, hash := range pool.hashesByBlockNumbers[log.BlockNumber] {
		num := i
		if hash == log.BlockHash {
			pool.hashesByBlockNumbers[log.BlockNumber] =
				append(pool.hashesByBlockNumbers[log.BlockNumber][:num], pool.hashesByBlockNumbers[log.BlockNumber][num+1:]...)
			break
		}
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
