package models

import (
	"time"

	"github.com/ethereum/go-ethereum/common"
)

// A LogConsumption is a unique record indicating that a particular job has
// already consumed a particular log. This record can be used to prevent consumers
// from re-processing duplicate logs
type LogConsumption struct {
	ID          uint
	BlockHash   common.Hash
	LogIndex    uint
	JobID       *ID
	JobIDV2     int32 `gorm:"column:job_id_v2"`
	BlockNumber uint64
	CreatedAt   time.Time
}

// NewLogConsumption creates a new LogConsumption
func NewLogConsumption(blockHash common.Hash, logIndex uint, jobID *ID, jobIDV2 int32, blockNumber uint64) LogConsumption {
	return LogConsumption{
		BlockHash:   blockHash,
		LogIndex:    logIndex,
		JobID:       jobID,
		JobIDV2:     jobIDV2,
		BlockNumber: blockNumber,
	}
}
