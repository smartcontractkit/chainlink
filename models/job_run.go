package models

import (
	"time"

	"github.com/smartcontractkit/chainlink-go/models/adapters"
)

type JobRun struct {
	ID        string `storm:"id,index,unique"`
	JobID     string `storm:"index"`
	Status    string
	CreatedAt time.Time          `storm:"index"`
	Result    adapters.RunResult `storm:"inline"`
	TaskRuns  []TaskRun          `storm:"inline"`
}
