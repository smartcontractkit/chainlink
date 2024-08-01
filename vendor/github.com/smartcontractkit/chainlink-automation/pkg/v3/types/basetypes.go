package types

import (
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/types/automation"
)

type UpkeepType uint8

const (
	// Exploratory AUTO 4335: add type for unknown
	ConditionTrigger UpkeepType = iota
	LogTrigger
)

// RetryRecord is a record of a payload that can be retried after a certain interval.
type RetryRecord struct {
	// payload is the desired unit of work to be retried
	Payload automation.UpkeepPayload
	// Interval is the time interval after which the same payload can be retried.
	Interval time.Duration
}
