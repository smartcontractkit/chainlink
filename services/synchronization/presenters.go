package synchronization

import (
	"encoding/json"

	"github.com/smartcontractkit/chainlink/store/models"
	"github.com/smartcontractkit/chainlink/utils"
	null "gopkg.in/guregu/null.v3"
)

// SyncJobRunPresenter presents a JobRun for synchronization purposes
type SyncJobRunPresenter struct {
	*models.JobRun
}

// MarshalJSON returns the JobRun as JSON
func (p SyncJobRunPresenter) MarshalJSON() ([]byte, error) {
	type SyncTaskRunPresenter struct {
		Index  int    `json:"index"`
		Type   string `json:"type"`
		Status string `json:"status"`
		Error  string `json:"error"`
	}
	tasks := []SyncTaskRunPresenter{}
	for index, task := range p.TaskRuns {
		tasks = append(tasks, SyncTaskRunPresenter{
			Index:  index,
			Type:   string(task.TaskSpec.Type),
			Status: string(task.Status),
			Error:  task.Result.ErrorMessage.ValueOrZero(),
		})
	}
	return json.Marshal(&struct {
		JobID       string                 `json:"jobID"`
		RunID       string                 `json:"runID"`
		Status      string                 `json:"status"`
		Error       null.String            `json:"error"`
		CreatedAt   string                 `json:"createdAt"`
		Amount      string                 `json:"amount"`
		CompletedAt null.Time              `json:"completedAt"`
		Tasks       []SyncTaskRunPresenter `json:"tasks"`
	}{
		RunID:       p.ID,
		JobID:       p.JobSpecID,
		Status:      string(p.Status),
		Error:       p.Result.ErrorMessage,
		CreatedAt:   utils.ISO8601UTC(p.CreatedAt),
		Amount:      p.Result.Amount.String(),
		CompletedAt: p.CompletedAt,
		Tasks:       tasks,
	})
}
