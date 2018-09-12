package migration1536764911

import (
	"github.com/smartcontractkit/chainlink/store/migrations/migration0"
	"github.com/smartcontractkit/chainlink/store/models"
	"github.com/smartcontractkit/chainlink/store/orm"
	null "gopkg.in/guregu/null.v3"
)

type Migration struct{}

func (m Migration) Timestamp() string {
	return "1536764911"
}

func (m Migration) Migrate(orm *orm.ORM) error {
	var jobSpecs []migration0.JobSpec
	err := orm.All(&jobSpecs)
	if err != nil {
		return err
	}

	tx, err := orm.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for _, oldJob := range jobSpecs {
		newJob := m.Convert(oldJob)
		err := tx.Save(&newJob)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (m Migration) Convert(js migration0.JobSpec) JobSpec {
	jsr := js.JobSpecRequest
	return JobSpec{
		ID:        js.ID,
		CreatedAt: js.CreatedAt,
		JobSpecRequest: JobSpecRequest{
			Initiators: jsr.Initiators,
			Tasks:      convertTaskSpecs(jsr.Tasks),
			StartAt:    jsr.StartAt,
			EndAt:      jsr.EndAt,
		},
	}
}

func convertTaskSpecs(oldSpecs []migration0.TaskSpec) []TaskSpec {
	var ts []TaskSpec
	for _, old := range oldSpecs {
		ts = append(ts, TaskSpec(old))
	}
	return ts
}

type TaskSpec struct {
	Type          migration0.TaskType `json:"type" storm:"index"`
	Confirmations uint64              `json:"confirmations"`
	Params        models.JSON         `json:"params"`
}

type JobSpec struct {
	ID        string          `json:"id" storm:"id,unique"`
	CreatedAt migration0.Time `json:"createdAt" storm:"index"`
	JobSpecRequest
}

// JobSpecRequest represents a schema for the incoming job spec request as used by the API.
type JobSpecRequest struct {
	Initiators []migration0.Initiator `json:"initiators"`
	Tasks      []TaskSpec             `json:"tasks" storm:"inline"`
	StartAt    null.Time              `json:"startAt" storm:"index"`
	EndAt      null.Time              `json:"endAt" storm:"index"`
}
