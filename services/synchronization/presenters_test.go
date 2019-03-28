package synchronization

import (
	"testing"

	"github.com/smartcontractkit/chainlink/store/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSyncJobRunPresenter(t *testing.T) {
	jobRun := models.JobRun{
		ID:        "AOC",
		JobSpecID: "BSFTW",
		Status:    models.RunStatusInProgress,
	}
	p := SyncJobRunPresenter{JobRun: &jobRun}

	bytes, err := p.MarshalJSON()
	require.NoError(t, err)
	json := string(bytes)

	assert.Contains(t, json, `"RunID":"AOC"`)
	assert.Contains(t, json, `"JobID":"BSFTW"`)
	assert.Contains(t, json, `"Status":"in_progress"`)
	assert.Contains(t, json, "Error")
	assert.Contains(t, json, "CreatedAt")
	assert.Contains(t, json, "Amount")
	assert.Contains(t, json, "CompletedAt")
	//assert.Contains(t, json, "Sender")
	assert.Contains(t, json, "Tasks")

	//Task fields:
	//- index
	//- type
	//- status
	//- error
}
