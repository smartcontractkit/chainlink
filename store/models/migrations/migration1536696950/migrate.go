package migration1536696950

import (
	"time"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/smartcontractkit/chainlink/store/assets"
	"github.com/smartcontractkit/chainlink/store/models"
	"github.com/smartcontractkit/chainlink/store/models/migrations/migration1536521223"
	"github.com/smartcontractkit/chainlink/store/models/orm"
	null "gopkg.in/guregu/null.v3"
)

type Migration struct{}

func (m Migration) Timestamp() string {
	return "1536696950"
}

func (m Migration) Migrate(orm *orm.ORM) error {
	var jrs []migration1536521223.JobRun
	if err := orm.All(&jrs); err != nil {
		return err
	}
	for _, jr := range jrs {
		jr2 := m.Convert(jr)
		if err := orm.Save(&jr2); err != nil {
			return err
		}
	}
	return nil
}

func (m Migration) Convert(jr migration1536521223.JobRun) JobRun {
	return JobRun{
		ID:             jr.ID,
		JobID:          jr.JobID,
		Result:         convertRunResult(jr.Result),
		Status:         jr.Status,
		TaskRuns:       convertTaskRuns(jr.TaskRuns),
		CreatedAt:      jr.CreatedAt,
		Initiator:      jr.Initiator,
		CreationHeight: jr.CreationHeight,
		Overrides:      convertRunResult(jr.Overrides),
	}
}

func convertRunResult(rr migration1536521223.RunResult) RunResult {
	return RunResult{
		JobRunID:     rr.JobRunID,
		Data:         rr.Data,
		Status:       rr.Status,
		ErrorMessage: rr.ErrorMessage,
		Amount:       (*assets.Link)(rr.Amount),
	}
}

func convertTaskRuns(oldTRs []migration1536521223.TaskRun) []TaskRun {
	var trs []TaskRun
	for _, otr := range oldTRs {
		trs = append(trs, TaskRun{
			ID:     otr.ID,
			Result: convertRunResult(otr.Result),
			Status: otr.Status,
			Task:   otr.Task,
		})
	}
	return trs
}

type JobRun struct {
	ID             string                        `json:"id" storm:"id,unique"`
	JobID          string                        `json:"jobId" storm:"index"`
	Result         RunResult                     `json:"result" storm:"inline"`
	Status         migration1536521223.RunStatus `json:"status" storm:"index"`
	TaskRuns       []TaskRun                     `json:"taskRuns" storm:"inline"`
	CreatedAt      time.Time                     `json:"createdAt" storm:"index"`
	CompletedAt    null.Time                     `json:"completedAt"`
	Initiator      migration1536521223.Initiator `json:"initiator"`
	CreationHeight *hexutil.Big                  `json:"creationHeight"`
	Overrides      RunResult                     `json:"overrides"`
}

type TaskRun struct {
	ID     string                        `json:"id" storm:"id,unique"`
	Result RunResult                     `json:"result"`
	Status migration1536521223.RunStatus `json:"status"`
	Task   migration1536521223.TaskSpec  `json:"task"`
}

type RunResult struct {
	JobRunID     string                        `json:"jobRunId"`
	Data         models.JSON                   `json:"data"`
	Status       migration1536521223.RunStatus `json:"status"`
	ErrorMessage null.String                   `json:"error"`
	Amount       *assets.Link                  `json:"amount,omitempty"`
}
