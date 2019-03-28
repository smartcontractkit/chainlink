package presenters

import (
	"testing"

	"github.com/smartcontractkit/chainlink/store/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_NewTx(t *testing.T) {
	t.Parallel()

	tx := models.Tx{
		GasLimit: uint64(5000),
		Nonce:    uint64(100),
		SentAt:   uint64(300),
	}
	ptx := NewTx(&tx)

	assert.Equal(t, "5000", ptx.GasLimit)
	assert.Equal(t, "100", ptx.Nonce)
	assert.Equal(t, "300", ptx.SentAt)
}

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
