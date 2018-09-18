package migration1537223654

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink/store/migrations/migration0"
	"github.com/smartcontractkit/chainlink/store/migrations/migration1537223654/old"
	"github.com/smartcontractkit/chainlink/store/orm"
)

type Migration struct{}

func (m Migration) Timestamp() string {
	return "1537223654"
}

func (m Migration) Migrate(orm *orm.ORM) error {
	var oldJobs []old.JobSpec
	if err := orm.All(&oldJobs); err != nil {
		return err
	}

	var oldInits []old.Initiator
	if err := orm.All(&oldInits); err != nil {
		return err
	}

	var oldRuns []old.JobRun
	if err := orm.All(&oldRuns); err != nil {
		return err
	}

	tx, err := orm.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for _, oj := range oldJobs {
		newJob := convert(oj)
		if err := tx.Save(&newJob); err != nil {
			return err
		}
	}

	for _, oi := range oldInits {
		newInit := convertInitiator(oi)
		if err := tx.Save(&newInit); err != nil {
			return err
		}
	}

	for _, oldRun := range oldRuns {
		newRun := convertJobRun(oldRun)
		if err := tx.Save(&newRun); err != nil {
			return err
		}
	}

	return tx.Commit()
}

func convert(oj old.JobSpec) JobSpec {
	return JobSpec{
		ID:         oj.ID,
		CreatedAt:  oj.CreatedAt,
		Initiators: convertInitiators(oj.Initiators),
		Tasks:      oj.Tasks,
		StartAt:    oj.StartAt,
		EndAt:      oj.EndAt,
	}
}

func convertInitiators(oldInits []old.Initiator) []Initiator {
	newInits := []Initiator{}
	for _, oi := range oldInits {
		newInits = append(newInits, convertInitiator(oi))
	}
	return newInits
}

func convertInitiator(ti old.Initiator) Initiator {
	return Initiator{
		ID:    ti.ID,
		JobID: ti.JobID,
		Type:  ti.Type,
		InitiatorParams: InitiatorParams{
			Schedule: ti.Schedule,
			Time:     ti.Time,
			Ran:      ti.Ran,
			Address:  ti.Address,
		},
	}
}

func convertJobRun(or old.JobRun) JobRun {
	ni := convertInitiator(or.Initiator)
	return JobRun{
		ID:             or.ID,
		JobID:          or.JobID,
		Result:         or.Result,
		Status:         or.Status,
		TaskRuns:       or.TaskRuns,
		CreatedAt:      or.CreatedAt,
		CompletedAt:    or.CompletedAt,
		Initiator:      ni,
		CreationHeight: or.CreationHeight,
		Overrides:      or.Overrides,
	}
}

type JobSpec struct {
	ID         migration0.Unchanged `json:"id" storm:"id,unique"`
	CreatedAt  migration0.Unchanged `json:"createdAt" storm:"index"`
	Initiators []Initiator          `json:"initiators"`
	Tasks      migration0.Unchanged `json:"tasks" storm:"inline"`
	StartAt    migration0.Unchanged `json:"startAt" storm:"index"`
	EndAt      migration0.Unchanged `json:"endAt" storm:"index"`
}

type Initiator struct {
	ID              int    `json:"id" storm:"id,increment"`
	JobID           string `json:"jobId" storm:"index"`
	Type            string `json:"type" storm:"index"`
	InitiatorParams `json:"params,omitempty"`
}

type InitiatorParams struct {
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
