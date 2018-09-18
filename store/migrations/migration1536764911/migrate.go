package migration1536764911

import (
	"github.com/smartcontractkit/chainlink/store/migrations/migration0"
	"github.com/smartcontractkit/chainlink/store/migrations/migration1536764911/old"
	"github.com/smartcontractkit/chainlink/store/orm"
)

type Migration struct{}

func (m Migration) Timestamp() string {
	return "1536764911"
}

func (m Migration) Migrate(orm *orm.ORM) error {
	var jobSpecs []old.JobSpec
	if err := orm.All(&jobSpecs); err != nil {
		return err
	}

	var jobRuns []old.JobRun
	if err := orm.All(&jobRuns); err != nil {
		return err
	}

	tx, err := orm.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for _, oldJob := range jobSpecs {
		newJob := m.Convert(oldJob)
		if err := tx.Save(&newJob); err != nil {
			return err
		}
	}

	for _, oldRun := range jobRuns {
		newRun := convertJobRun(oldRun)
		if err := tx.Save(&newRun); err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (m Migration) Convert(js old.JobSpec) JobSpec {
	return JobSpec{
		ID:         js.ID,
		CreatedAt:  js.CreatedAt,
		Initiators: js.Initiators,
		Tasks:      convertTaskSpecs(js.Tasks),
		StartAt:    js.StartAt,
		EndAt:      js.EndAt,
	}
}

func convertTaskSpecs(oldSpecs []old.TaskSpec) []TaskSpec {
	ts := []TaskSpec{}
	for _, ot := range oldSpecs {
		ts = append(ts, convertTaskSpec(ot))
	}
	return ts
}

func convertTaskSpec(ot old.TaskSpec) TaskSpec {
	return (TaskSpec)(ot)
}

func convertJobRun(or old.JobRun) JobRun {
	return JobRun{
		ID:             or.ID,
		JobID:          or.JobID,
		Result:         or.Result,
		Status:         or.Status,
		TaskRuns:       convertTaskRuns(or.TaskRuns),
		CreatedAt:      or.CreatedAt,
		CompletedAt:    or.CompletedAt,
		Initiator:      or.Initiator,
		CreationHeight: or.CreationHeight,
		Overrides:      or.Overrides,
	}
}

func convertTaskRuns(ots []old.TaskRun) []TaskRun {
	nts := []TaskRun{}
	for _, ot := range ots {
		nts = append(nts, TaskRun{
			ID:     ot.ID,
			Result: ot.Result,
			Status: ot.Status,
			Task:   convertTaskSpec(ot.Task),
		})
	}
	return nts
}

type TaskSpec struct {
	Type          migration0.Unchanged `json:"type" storm:"index"`
	Confirmations migration0.Unchanged `json:"confirmations"`
	Params        migration0.JSON      `json:"params"`
}

type JobSpec struct {
	ID         migration0.Unchanged `json:"id" storm:"id,unique"`
	CreatedAt  migration0.Unchanged `json:"createdAt" storm:"index"`
	Initiators migration0.Unchanged `json:"initiators"`
	Tasks      []TaskSpec           `json:"tasks" storm:"inline"`
	StartAt    migration0.Unchanged `json:"startAt" storm:"index"`
	EndAt      migration0.Unchanged `json:"endAt" storm:"index"`
}

type JobRun struct {
	ID             migration0.Unchanged `json:"id" storm:"id,unique"`
	JobID          migration0.Unchanged `json:"jobId" storm:"index"`
	Result         migration0.Unchanged `json:"result" storm:"inline"`
	Status         migration0.Unchanged `json:"status" storm:"index"`
	TaskRuns       []TaskRun            `json:"taskRuns" storm:"inline"`
	CreatedAt      migration0.Unchanged `json:"createdAt" storm:"index"`
	CompletedAt    migration0.Unchanged `json:"completedAt"`
	Initiator      migration0.Unchanged `json:"initiator"`
	CreationHeight migration0.Unchanged `json:"creationHeight"`
	Overrides      migration0.Unchanged `json:"overrides"`
}

type TaskRun struct {
	ID     migration0.Unchanged `json:"id" storm:"id,unique"`
	Result migration0.Unchanged `json:"result"`
	Status migration0.Unchanged `json:"status"`
	Task   TaskSpec             `json:"task"`
}
