package old

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink/store/migrations/migration0"
)

type JobSpec struct {
	ID         migration0.Unchanged `json:"id" storm:"id,unique"`
	CreatedAt  migration0.Unchanged `json:"createdAt" storm:"index"`
	Initiators []Initiator          `json:"initiators"`
	Tasks      migration0.Unchanged `json:"tasks" storm:"inline"`
	StartAt    migration0.Unchanged `json:"startAt" storm:"index"`
	EndAt      migration0.Unchanged `json:"endAt" storm:"index"`
}

type Initiator struct {
	ID       int             `json:"id" storm:"id,increment"`
	JobID    string          `json:"jobId" storm:"index"`
	Type     string          `json:"type" storm:"index"`
	Schedule migration0.Cron `json:"schedule,omitempty"`
	Time     migration0.Time `json:"time,omitempty"`
	Ran      bool            `json:"ran,omitempty"`
	Address  common.Address  `json:"address,omitempty" storm:"index"`
}

type JobRun struct {
	ID             migration0.Unchanged `json:"id" storm:"id,unique"`
	JobID          migration0.Unchanged `json:"jobId" storm:"index"`
	Result         migration0.Unchanged `json:"result" storm:"inline"`
	Status         migration0.Unchanged `json:"status" storm:"index"`
	TaskRuns       migration0.Unchanged `json:"taskRuns" storm:"inline"`
	CreatedAt      migration0.Unchanged `json:"createdAt" storm:"index"`
	CompletedAt    migration0.Unchanged `json:"completedAt"`
	Initiator      Initiator            `json:"initiator"`
	CreationHeight migration0.Unchanged `json:"creationHeight"`
	Overrides      migration0.Unchanged `json:"overrides"`
}
