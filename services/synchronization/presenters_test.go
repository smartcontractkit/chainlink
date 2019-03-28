package synchronization

import (
	"encoding/json"
	"testing"

	"github.com/smartcontractkit/chainlink/store/assets"
	"github.com/smartcontractkit/chainlink/store/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	null "gopkg.in/guregu/null.v3"
)

func TestSyncJobRunPresenter(t *testing.T) {
	jobRun := models.JobRun{
		ID:        "runID-411",
		JobSpecID: "jobSpecID-312",
		Status:    models.RunStatusInProgress,
		Result:    models.RunResult{Amount: assets.NewLink(2)},
		TaskRuns: []models.TaskRun{
			models.TaskRun{
				ID:     "task0RunID-938",
				Status: models.RunStatusPendingConfirmations,
			},
			models.TaskRun{
				ID:     "task1RunID-17",
				Status: models.RunStatusErrored,
				Result: models.RunResult{ErrorMessage: null.StringFrom("yikes fam")},
			},
		},
	}
	p := SyncJobRunPresenter{JobRun: &jobRun}

	bytes, err := p.MarshalJSON()
	require.NoError(t, err)

	var data map[string]interface{}
	err = json.Unmarshal(bytes, &data)
	require.NoError(t, err)

	assert.Equal(t, data["runID"], "runID-411")
	assert.Equal(t, data["jobID"], "jobSpecID-312")
	assert.Equal(t, data["status"], "in_progress")
	assert.Contains(t, data, "error")
	assert.Contains(t, data, "createdAt")
	assert.Equal(t, data["amount"], "2")
	assert.Equal(t, data["completedAt"], nil)
	assert.Contains(t, data, "tasks")

	tasks, ok := data["tasks"].([]interface{})
	require.True(t, ok)
	require.Len(t, tasks, 2)
	task0, ok := tasks[0].(map[string]interface{})
	require.True(t, ok)
	assert.Equal(t, task0["index"], float64(0))
	assert.Contains(t, task0, "type")
	assert.Equal(t, task0["status"], "pending_confirmations")
	assert.Equal(t, task0["error"], nil)
	task1, ok := tasks[1].(map[string]interface{})
	require.True(t, ok)
	assert.Equal(t, task1["index"], float64(1))
	assert.Contains(t, task1, "type")
	assert.Equal(t, task1["status"], "errored")
	assert.Equal(t, task1["error"], "yikes fam")
}
