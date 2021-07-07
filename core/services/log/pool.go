package log

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/store/models"
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

	logger.Tracew(fmt.Sprintf("LogBroadcaster: ADDED LOG NUM %v %v", log.BlockNumber, log.BlockHash))
}

func (pool *logPool) getLogsToSend(head models.Head, highestNumConfirmations uint64, finalityDepth uint64) []types.Log {
	latestBlockNum := uint64(head.Number)
	logsToReturn := make([]types.Log, 0)

	keptLogsDepth := finalityDepth
	if highestNumConfirmations > keptLogsDepth {
		keptLogsDepth = highestNumConfirmations
	}

	keptDepth := int64(latestBlockNum) - int64(keptLogsDepth)
	if keptDepth < 0 {
		keptDepth = 0
	}
	logger.Tracew(fmt.Sprintf("LogBroadcaster: keptDepth %v latestBlockNum: %v", keptDepth, latestBlockNum),
		"latestBlockNum", latestBlockNum,
		"highestNumConfirmations", highestNumConfirmations,
		"finalityDepth", finalityDepth,
	)

	minBlockNumToSendItem := pool.heap.FindMin()
	if minBlockNumToSendItem == nil {
		return logsToReturn
	}
	minBlockNumToSend := int64(minBlockNumToSendItem.(Uint64))

	for i := minBlockNumToSend; i <= int64(latestBlockNum); i++ {
		num := i

		logger.Tracew(fmt.Sprintf("LogBroadcaster: BLOCK NUM %v - num hashes: %v", num, len(pool.hashesByBlockNumbers[uint64(num)])),
			"latestBlockNum", latestBlockNum,
			"highestNumConfirmations", highestNumConfirmations,
			"finalityDepth", finalityDepth,
		)
		for _, hash := range pool.hashesByBlockNumbers[uint64(num)] {
			logger.Tracew(fmt.Sprintf("LogBroadcaster: ADDING %v %v", num, hash),
				"latestBlockNum", latestBlockNum,
				"highestNumConfirmations", highestNumConfirmations,
				"finalityDepth", finalityDepth,
			)
			logsToReturn = append(logsToReturn, pool.logsByBlockHash[hash]...)
		}
	}

	for {
		item := pool.heap.FindMin()
		if item == nil {
			break
		}

		blockNum := uint64(item.(Uint64))
		if blockNum >= uint64(keptDepth) {
			break
		}
		pool.heap.DeleteMin()

		for _, hash := range pool.hashesByBlockNumbers[blockNum] {
			logger.Tracew(fmt.Sprintf("LogBroadcaster: Will delete %v logs for block hash %v", len(pool.logsByBlockHash[hash]), hash),
				"latestBlockNum", latestBlockNum,
				"highestNumConfirmations", highestNumConfirmations,
				"finalityDepth", finalityDepth,
			)
			delete(pool.logsByBlockHash, hash)
		}
		delete(pool.hashesByBlockNumbers, blockNum)
	}

	if len(logsToReturn) > 0 {
		logger.Tracew(fmt.Sprintf("LogBroadcaster: Will return %v logs", len(logsToReturn)),
			"latestBlockNum", latestBlockNum,
			"highestNumConfirmations", highestNumConfirmations,
			"finalityDepth", finalityDepth,
		)
	}
	return logsToReturn
}

func (pool *logPool) removeLog(log types.Log) {
	// deleting all logs for this log's block hash
	delete(pool.logsByBlockHash, log.BlockHash)
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
