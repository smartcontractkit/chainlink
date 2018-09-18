package migration1536764911

import (
	"github.com/smartcontractkit/chainlink/store/migrations/migration0"
	"github.com/smartcontractkit/chainlink/store/orm"
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
	return JobSpec{
		ID:         js.ID,
		CreatedAt:  js.CreatedAt,
		Initiators: js.Initiators,
		Tasks:      convertTaskSpecs(js.Tasks),
		StartAt:    js.StartAt,
		EndAt:      js.EndAt,
	}
}

func convertTaskSpecs(oldSpecs []migration0.TaskSpec) []TaskSpec {
	ts := []TaskSpec{}
	for _, old := range oldSpecs {
		new := TaskSpec{
			Type:          old.Type,
			Confirmations: old.Confirmations,
			Params:        old.Params,
		}
		ts = append(ts, new)
	}
	return ts
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
