package models_test

import (
	"testing"

	"github.com/smartcontractkit/chainlink-go/internal/cltest"
	"github.com/smartcontractkit/chainlink-go/models"
	"github.com/stretchr/testify/assert"
)

func TestJobSave(t *testing.T) {
	t.Parallel()
	store := cltest.Store()
	defer store.Close()

	j1 := models.NewJob()
	j1.Schedule = models.Schedule{Cron: "1 * * * *"}

	store.Save(&j1)

	var j2 models.Job
	store.One("ID", j1.ID, &j2)

	assert.Equal(t, j1.Schedule, j2.Schedule)
}

func TestJobNewRun(t *testing.T) {
	t.Parallel()
	store := cltest.Store()
	defer store.Close()

	job := models.NewJob()
	job.Schedule = models.Schedule{Cron: "1 * * * *"}
	job.Tasks = []models.Task{models.Task{Type: "NoOp"}}

	newRun := job.NewRun()
	assert.Equal(t, job.ID, newRun.JobID)
	assert.Equal(t, 1, len(newRun.TaskRuns))
	assert.Equal(t, "NoOp", job.Tasks[0].Type)
	assert.Nil(t, job.Tasks[0].Params)
	adapter, _ := job.Tasks[0].Adapter()
	assert.NotNil(t, adapter)
}
