package models

import (
	"time"

	"github.com/ethereum/go-ethereum/common"
)

// TODO - RYAN

// LogConsumerTypeJob ...
const LogConsumerTypeJob = "job"

// LogConsumerTypes ...
var LogConsumerTypes = [1]string{LogConsumerTypeJob}

// LogConsumption ...
type LogConsumption struct {
	ID           *ID
	BlockHash    common.Hash
	LogIndex     uint
	ConsumerType string // ["job", ...?]
	ConsumerID   *ID
	CreatedAt    time.Time
}

// NewLogConsumption ...
func NewLogConsumption() LogConsumption {
	return LogConsumption{
		ID:        NewID(),
		CreatedAt: time.Now(),
	}
}
