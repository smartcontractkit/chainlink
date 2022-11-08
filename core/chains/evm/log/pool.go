package log

import (
	"math"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	heaps "github.com/theodesp/go-heaps"
	pairingHeap "github.com/theodesp/go-heaps/pairing"

	"github.com/smartcontractkit/chainlink/core/logger"
)

//go:generate mockery --quiet --name iLogPool --output ./ --inpackage --testonly

// The Log Pool interface.
type iLogPool interface {

	// AddLog adds log to the pool and returns true if its block number is a new minimum.
	addLog(log types.Log) bool

	// GetAndDeleteAll purges the pool completely, returns all logs, and also the minimum and
	// maximum block numbers retrieved.
	getAndDeleteAll() ([]logsOnBlock, int64, int64)

	// GetLogsToSend returns all logs upto the block number specified in latestBlockNum.
	// Also returns the minimum block number in the result.
	// In case the pool is empty, returns empty results, and min block number=0
	getLogsToSend(latestBlockNum int64) ([]logsOnBlock, int64)

	// DeleteOlderLogs deletes all logs in blocks that are less than specific block number keptDepth.
	// Also returns the remaining minimum block number in pool after these deletions.
	// Returns nil if this ends up emptying the pool.
	deleteOlderLogs(keptDepth int64) *int64

	// RemoveBlock removes all logs for the block identified by provided Block hash and number.
	removeBlock(hash common.Hash, number uint64)

	// TestOnly_getNumLogsForBlock FOR TESTING USE ONLY.
	// Returns all logs for the provided block hash.
	testOnly_getNumLogsForBlock(bh common.Hash) int
}

type logPool struct {
	// A mapping of block numbers to a set of block hashes for all
	// the logs in the pool.
	hashesByBlockNumbers map[uint64]map[common.Hash]struct{}

	// A mapping of block hashes, to tx index within block, to log index, to logs
	logsByBlockHash map[common.Hash]map[uint]map[uint]types.Log

	// This min-heap maintains block numbers of logs in the pool.
	// it helps us easily determine the minimum log block number
	// in the pool (while the set of log block numbers is dynamically changing).
	heap   *pairingHeap.PairHeap
	logger logger.Logger
}

func newLogPool(lggr logger.Logger) *logPool {
	return &logPool{
		hashesByBlockNumbers: make(map[uint64]map[common.Hash]struct{}),
		logsByBlockHash:      make(map[common.Hash]map[uint]map[uint]types.Log),
		heap:                 pairingHeap.New(),
		logger:               lggr.Named("LogPool"),
	}
}

func (pool *logPool) addLog(log types.Log) bool {
	_, exists := pool.hashesByBlockNumbers[log.BlockNumber]
	if !exists {
		pool.hashesByBlockNumbers[log.BlockNumber] = make(map[common.Hash]struct{})
	}
	pool.hashesByBlockNumbers[log.BlockNumber][log.BlockHash] = struct{}{}
	if _, exists := pool.logsByBlockHash[log.BlockHash]; !exists {
		pool.logsByBlockHash[log.BlockHash] = make(map[uint]map[uint]types.Log)
	}
	if _, exists := pool.logsByBlockHash[log.BlockHash][log.TxIndex]; !exists {
		pool.logsByBlockHash[log.BlockHash][log.TxIndex] = make(map[uint]types.Log)
	}
	pool.logsByBlockHash[log.BlockHash][log.TxIndex][log.Index] = log
	min := pool.heap.FindMin()
	pool.heap.Insert(Uint64(log.BlockNumber))
	pool.logger.Debugw("Inserted block to log pool", "blockNumber", log.BlockNumber, "blockHash", log.BlockHash, "index", log.Index, "prevMinBlockNumber", min)
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

func (pool *logPool) removeBlock(hash common.Hash, number uint64) {
	// deleting all logs for this log's block hash
	delete(pool.logsByBlockHash, hash)
	delete(pool.hashesByBlockNumbers[number], hash)
	if len(pool.hashesByBlockNumbers[number]) == 0 {
		delete(pool.hashesByBlockNumbers, number)
	}
}

func (pool *logPool) testOnly_getNumLogsForBlock(bh common.Hash) int {
	var numLogs int
	for _, txLogs := range pool.logsByBlockHash[bh] {
		numLogs += len(txLogs)
	}
	return numLogs
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

func newLogsOnBlock(num uint64, logsMap map[uint]map[uint]types.Log) logsOnBlock {
	logs := make([]types.Log, 0, len(logsMap))
	for _, txLogs := range logsMap {
		for _, l := range txLogs {
			logs = append(logs, l)
		}
	}
	return logsOnBlock{num, logs}
}
