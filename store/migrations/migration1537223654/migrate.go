package migration1537223654

import (
	"encoding/json"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink/store/migrations/migration0"
	"github.com/smartcontractkit/chainlink/store/migrations/migration1536696950"
	"github.com/smartcontractkit/chainlink/store/migrations/migration1536764911"
	"github.com/smartcontractkit/chainlink/store/orm"
)

type Migration struct{}

func (m Migration) Timestamp() string {
	return "1537223654"
}

func (m Migration) Migrate(orm *orm.ORM) error {
	var oldJobs []migration1536764911.JobSpec
	err := orm.All(&oldJobs)
	if err != nil {
		return err
	}

	var oldInits []migration0.Initiator
	err = orm.All(&oldInits)
	if err != nil {
		return err
	}

	var oldRuns []migration1536696950.JobRun
	err = orm.All(&oldRuns)
	if err != nil {
		return err
	}

	tx, err := orm.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for _, oj := range oldJobs {
		newJob, err := convert(oj)
		if err != nil {
			return err
		}
		err = tx.Save(&newJob)
		if err != nil {
			return err
		}
	}

	for _, oi := range oldInits {
		newInit, err := convertInitiator(oi)
		if err != nil {
			return err
		}
		err = tx.Save(&newInit)
		if err != nil {
			return err
		}
	}

	for _, oldRun := range oldRuns {
		newRun, err := convertJobRun(oldRun)
		if err != nil {
			return err
		}
		err = tx.Save(&newRun)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func convert(oj migration1536764911.JobSpec) (JobSpec, error) {
	newInits, err := convertInitiators(oj.Initiators)
	return JobSpec{
		ID:         oj.ID,
		CreatedAt:  oj.CreatedAt,
		Initiators: newInits,
		Tasks:      oj.Tasks,
		StartAt:    oj.StartAt,
		EndAt:      oj.EndAt,
	}, err
}

func convertInitiators(unchanged migration0.Unchanged) ([]Initiator, error) {
	oldInits, ok := unchanged.([]interface{})
	if !ok {
		return []Initiator{}, fmt.Errorf("convertInitiators: Unable to convert %v of type %T to []migration0.interface{}", unchanged, unchanged)
	}

	newInits := []Initiator{}
	for _, oi := range oldInits {
		ni, err := convertInitiator(migration0.Unchanged(oi))
		if err != nil {
			return newInits, err
		}
		newInits = append(newInits, ni)
	}
	return newInits, nil
}

func UnchangedToInitiator(uc migration0.Unchanged) (migration0.Initiator, error) {
	b, err := json.Marshal(uc)
	if err != nil {
		return migration0.Initiator{}, err
	}
	var ti migration0.Initiator
	if err = json.Unmarshal(b, &ti); err != nil {
		return migration0.Initiator{}, err
	}
	return ti, nil
}

func convertInitiator(uc migration0.Unchanged) (Initiator, error) {
	ti, err := UnchangedToInitiator(uc)
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
	}, err
}

func convertJobRun(or migration1536696950.JobRun) (JobRun, error) {
	ni, err := convertInitiator(or.Initiator)
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
	}, err
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
