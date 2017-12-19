package services_test

import (
	"testing"

	"github.com/smartcontractkit/chainlink-go/internal/cltest"
	"github.com/smartcontractkit/chainlink-go/store/models"
	"github.com/smartcontractkit/chainlink-go/services"
	"github.com/stretchr/testify/assert"
)

func TestRunningJob(t *testing.T) {
	t.Parallel()
	store := cltest.NewStore()
	defer store.Close()

	job := models.NewJob()
	job.Schedule = models.Schedule{Cron: "* * * * *"}
	job.Tasks = []models.Task{models.Task{Type: "NoOp"}}

	run := job.NewRun()
	assert.Equal(t, "", run.Status)

	services.StartJob(run, store)

	store.One("ID", run.ID, &run)
	assert.Equal(t, "completed", run.Status)
	assert.Equal(t, "completed", run.TaskRuns[0].Status)
}
