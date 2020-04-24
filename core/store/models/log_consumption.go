package models

import (
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink/core/eth"
)

// LogConsumerTypeJob - LogConsumptions with this type were consumed by a job
const LogConsumerTypeJob = "job"

// LogConsumerTypes holds the list of valid consumer types
var LogConsumerTypes = [1]string{LogConsumerTypeJob}

// A LogConsumption is a unique record indicating that a particular consumer has
// already consumed a particular log. This record can be used to prevent consumers
// from re-processing duplicate logs
type LogConsumption struct {
	ID           *ID
	BlockHash    common.Hash
	LogIndex     uint
	ConsumerType string
	ConsumerID   *ID
	CreatedAt    time.Time
}

// A LogConsumer has a type and ID, and uniquely identifies a LogListener
type LogConsumer struct {
	Type string
	ID   *ID
}

// NewEmptyLogConsumption Creates a new LogConsumption
func NewEmptyLogConsumption() LogConsumption {
	return LogConsumption{
		ID:        NewID(),
		CreatedAt: time.Now(),
	}
}

// NewLogConsumption creates a new LogConsumption
func NewLogConsumption(log eth.RawLog, consumer LogConsumer) LogConsumption {
	lc := NewEmptyLogConsumption()
	lc.BlockHash = log.GetBlockHash()
	lc.LogIndex = log.GetIndex()
	lc.ConsumerType = consumer.Type
	lc.ConsumerID = consumer.ID
	return lc
}
