package models_test

import (
	"testing"

	"github.com/smartcontractkit/chainlink-go/adapters"
	"github.com/smartcontractkit/chainlink-go/internal/cltest"
	"github.com/smartcontractkit/chainlink-go/store/models"
	"github.com/stretchr/testify/assert"
)

func TestJobSave(t *testing.T) {
	t.Parallel()
	store := cltest.NewStore()
	defer store.Close()

	j1 := cltest.NewJobWithSchedule("* * * * *")

	store.Save(&j1)

	var j2 models.Job
	store.One("ID", j1.ID, &j2)

	assert.Equal(t, j1.Initiators[0].Schedule, j2.Initiators[0].Schedule)
}

func TestJobNewRun(t *testing.T) {
	t.Parallel()
	store := cltest.NewStore()
	defer store.Close()

	job := cltest.NewJobWithSchedule("1 * * * *")
	job.Tasks = []models.Task{models.Task{Type: "NoOp"}}

	newRun := job.NewRun()
	assert.Equal(t, job.ID, newRun.JobID)
	assert.Equal(t, 1, len(newRun.TaskRuns))
	assert.Equal(t, "NoOp", job.Tasks[0].Type)
	assert.Nil(t, job.Tasks[0].Params)
	adapter, _ := adapters.For(job.Tasks[0], store)
	assert.NotNil(t, adapter)
}
