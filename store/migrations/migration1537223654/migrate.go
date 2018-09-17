package migration1537223654

import (
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/smartcontractkit/chainlink/store/migrations/migration0"
	"github.com/smartcontractkit/chainlink/store/migrations/migration1536696950"
	"github.com/smartcontractkit/chainlink/store/migrations/migration1536764911"
	"github.com/smartcontractkit/chainlink/store/orm"
	null "gopkg.in/guregu/null.v3"
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
		newJob := convert(oj)
		err := tx.Save(&newJob)
		if err != nil {
			return err
		}
	}

	for _, oi := range oldInits {
		newInit := convertInitiator(oi)
		err := tx.Save(&newInit)
		if err != nil {
			return err
		}
	}

	for _, oldRun := range oldRuns {
		newRun := convertJobRun(oldRun)
		err := tx.Save(&newRun)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func convert(oj migration1536764911.JobSpec) JobSpec {
	return JobSpec{
		ID:         oj.ID,
		CreatedAt:  oj.CreatedAt,
		Initiators: convertInitiators(oj.Initiators),
		Tasks:      oj.Tasks,
		StartAt:    oj.StartAt,
		EndAt:      oj.EndAt,
	}
}

func convertInitiators(oldInits []migration0.Initiator) []Initiator {
	newInits := []Initiator{}
	for _, oi := range oldInits {
		newInits = append(newInits, convertInitiator(oi))
	}
	return newInits
}

func convertInitiator(oi migration0.Initiator) Initiator {
	return Initiator{
		ID:    oi.ID,
		JobID: oi.JobID,
		Type:  oi.Type,
		InitiatorParams: InitiatorParams{
			Schedule: oi.Schedule,
			Time:     oi.Time,
			Ran:      oi.Ran,
			Address:  oi.Address,
		},
	}
}

func convertJobRun(or migration1536696950.JobRun) JobRun {
	return JobRun{
		ID:             or.ID,
		JobID:          or.JobID,
		Result:         or.Result,
		Status:         or.Status,
		TaskRuns:       or.TaskRuns,
		CreatedAt:      or.CreatedAt,
		CompletedAt:    or.CompletedAt,
		Initiator:      convertInitiator(or.Initiator),
		CreationHeight: or.CreationHeight,
		Overrides:      or.Overrides,
	}
}

type JobSpec struct {
	ID         string                         `json:"id" storm:"id,unique"`
	CreatedAt  migration0.Time                `json:"createdAt" storm:"index"`
	Initiators []Initiator                    `json:"initiators"`
	Tasks      []migration1536764911.TaskSpec `json:"tasks" storm:"inline"`
	StartAt    null.Time                      `json:"startAt" storm:"index"`
	EndAt      null.Time                      `json:"endAt" storm:"index"`
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
	ID             string                        `json:"id" storm:"id,unique"`
	JobID          string                        `json:"jobId" storm:"index"`
	Result         migration1536696950.RunResult `json:"result" storm:"inline"`
	Status         migration0.RunStatus          `json:"status" storm:"index"`
	TaskRuns       []migration1536696950.TaskRun `json:"taskRuns" storm:"inline"`
	CreatedAt      time.Time                     `json:"createdAt" storm:"index"`
	CompletedAt    null.Time                     `json:"completedAt"`
	Initiator      Initiator                     `json:"initiator"`
	CreationHeight *hexutil.Big                  `json:"creationHeight"`
	Overrides      migration1536696950.RunResult `json:"overrides"`
}
