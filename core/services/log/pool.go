package log

import (
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/store/models"
)

type (
	logPool struct {
		allLogs []types.Log
	}
)

func newLogPool() *logPool {
	return &logPool{
		allLogs: make([]types.Log, 0),
	}
}

func (pool *logPool) addLog(log types.Log) {
	pool.allLogs = append(pool.allLogs, log)
}

func (pool *logPool) getLogsToSend(head *models.Head, highestNumConfirmations uint64) []types.Log {
	latestBlockNum := uint64(head.Number)
	logsToReturn := pool.allLogs
	logsToKeep := make([]types.Log, 0)

	// deleting old logs that will never be sent for any listener anymore
	if latestBlockNum > highestNumConfirmations && len(pool.allLogs) > 0 {
		for _, log := range pool.allLogs {
			if log.BlockNumber >= latestBlockNum-highestNumConfirmations {
				logsToKeep = append(logsToKeep, log)
			}
		}
		logger.Tracef("latestBlockNum (%v) > highestNumConfirmations (%v) so deleting older logs. Remaining: %v", latestBlockNum, highestNumConfirmations, len(logsToKeep))
		pool.allLogs = logsToKeep
	}
	return logsToReturn
}
