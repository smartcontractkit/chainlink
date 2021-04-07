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

func (pool *logPool) getLogsToSend(head *models.Head, highestNumConfirmations uint64, finalityDepth uint64) []types.Log {
	latestBlockNum := uint64(head.Number)
	logsToReturn := pool.allLogs
	logsToKeep := make([]types.Log, 0)

	keptLogsDepth := finalityDepth
	if highestNumConfirmations > keptLogsDepth {
		keptLogsDepth = highestNumConfirmations
	}
	// deleting old logs that will never be sent for any listener anymore
	if latestBlockNum > keptLogsDepth && len(pool.allLogs) > 0 {
		for _, log := range pool.allLogs {
			if log.BlockNumber >= latestBlockNum-keptLogsDepth {
				logsToKeep = append(logsToKeep, log)
			}
		}
		logger.Tracew("LogBroadcaster: latestBlockNum > highestNumConfirmations so deleting older logs", "latestBlockNum", latestBlockNum, "highestNumConfirmations", highestNumConfirmations, "remainingLogsCount", len(logsToKeep))
		pool.allLogs = logsToKeep
	}
	return logsToReturn
}
