package synchronization

import (
	"encoding/json"

	"github.com/smartcontractkit/chainlink/store/models"
	"github.com/smartcontractkit/chainlink/utils"
)

// SyncJobRunPresenter presents a JobRun for synchronization purposes
type SyncJobRunPresenter struct {
	*models.JobRun
}

// MarshalJSON returns the JobRun as JSON
func (p SyncJobRunPresenter) MarshalJSON() ([]byte, error) {
	type SyncTaskRunPresenter struct {
		Index  int
		Type   string
		Status string
		Error  string
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
		JobID       string
		RunID       string
		Status      string
		Error       string
		CreatedAt   string
		Amount      string
		CompletedAt string
		Tasks       []interface{}
	}{
		RunID:       p.ID,
		JobID:       p.JobSpecID,
		Status:      string(p.Status),
		Error:       p.Result.ErrorMessage.ValueOrZero(),
		CreatedAt:   utils.ISO8601UTC(p.CreatedAt),
		Amount:      p.Result.Amount.String(),
		CompletedAt: utils.ISO8601UTC(p.CompletedAt.ValueOrZero()),
	})
}
