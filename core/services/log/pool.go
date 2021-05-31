package log

import (
	"fmt"

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
	return &logPool{}
}

func (pool *logPool) addLog(log types.Log) {
	pool.allLogs = append(pool.allLogs, log)
}

func (pool *logPool) getLogsToSend(head models.Head, highestNumConfirmations uint64, finalityDepth uint64) []types.Log {
	latestBlockNum := uint64(head.Number)
	logsToReturn := pool.allLogs
	var logsToKeep []types.Log

	keptLogsDepth := finalityDepth
	if highestNumConfirmations > keptLogsDepth {
		keptLogsDepth = highestNumConfirmations
	}
	// deleting old logs that will never be sent for any listener anymore
	if len(pool.allLogs) > 0 {
		for _, log := range pool.allLogs {
			if int64(log.BlockNumber) >= int64(latestBlockNum)-int64(keptLogsDepth) {
				logsToKeep = append(logsToKeep, log)
			}
		}
		logger.Tracew(fmt.Sprintf("LogBroadcaster: Will delete %v older logs", len(pool.allLogs)-len(logsToKeep)),
			"latestBlockNum", latestBlockNum,
			"highestNumConfirmations", highestNumConfirmations,
			"remainingLogsCount", len(logsToKeep),
			"finalityDepth", finalityDepth,
		)
		pool.allLogs = logsToKeep
	}
	return logsToReturn
}
