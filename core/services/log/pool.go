package log

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/smartcontractkit/chainlink/core/services/eth"
	"github.com/smartcontractkit/chainlink/core/store/models"
)

type (
	logPool struct {
		ethClient eth.Client
		logs      map[common.Address]logGroup
	}
	logGroup struct {
		logsByHeight map[uint64][]types.Log
	}
)

func newLogGroup() logGroup {
	return logGroup{
		logsByHeight: make(map[uint64][]types.Log),
	}
}
func (pool logGroup) addLog(log types.Log) {
	height := log.BlockNumber
	pool.logsByHeight[height] = append(pool.logsByHeight[height], log)
}

func (pool logGroup) getByHeight(height uint64, blockHash common.Hash) []types.Log {
	toReturn := make([]types.Log, 0)
	logs := pool.logsByHeight[height]
	for _, log := range logs {
		if log.BlockHash == blockHash {
			toReturn = append(toReturn, log)
		}
	}
	return toReturn
}

func newLogPool() logPool {
	return logPool{
		logs: map[common.Address]logGroup{},
	}
}

func (pool logPool) addLog(log types.Log) {
	if _, exists := pool.logs[log.Address]; !exists {
		pool.logs[log.Address] = newLogGroup()
	}
	pool.logs[log.Address].addLog(log)
}

func (pool logPool) getLogsToSend(chains []*models.Head, confirmationDepths map[uint64]struct{}) []models.Log {
	latestBlock := chains[len(chains)-1]
	latestBlockNum := uint64(latestBlock.Number)

	logs := make([]models.Log, 0)
	for depth := range confirmationDepths {
		for _, group := range pool.logs {

			height := latestBlockNum - depth
			logs = append(logs, group.getByHeight(height, latestBlock.Hash)...)
		}
	}
	return logs
}

/*

	// Defer processing more logs from this contract if we haven't received its block yet
	if log.BlockNumber > b.latestBlock {
		logger.Infof("xyzzy no block yet")
		nextLogs[contractAddr] = logs[i:]
		break
	}

	// Skip logs that have been reorged away
	if _, exists := b.canonicalChain[log.BlockHash]; !exists {
		logger.Infof("xyzzy reorged away: %v %v", log.BlockNumber, log.BlockHash)
		continue
	}
*/

//func (b *broadcaster) broadcastPendingLogs() {
//	b.logsMu.Lock()
//	defer b.logsMu.Unlock()
//
//	currentChain := b.currentChain()
//	if currentChain == nil {
//		return
//	}
//
//	nextLogs := make(map[common.Address]map[int64][]types.Log)
//
//	for _, logsByHeight := range b.logsByHeight {
//
//		cutOff := int64(5)
//		minConfirmations := int64(1)
//		blockHeightToSend := currentChain.Number - minConfirmations
//		blockHeightToDelete := currentChain.Number - cutOff
//		logs := logsByHeight[blockHeightToSend]
//		delete(logsByHeight, blockHeightToSend)
//		delete(logsByHeight, blockHeightToDelete)
//
//		for _, log := range logs {
//			logger.Infof("log: %v", log.BlockNumber)
//			b.registrations.sendLog(log, b.orm)
//		}
//	}
//	b.logsByHeight = nextLogs
//}
