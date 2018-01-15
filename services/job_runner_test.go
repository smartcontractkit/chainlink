package services_test

import (
	"testing"

	"github.com/smartcontractkit/chainlink/internal/cltest"
	"github.com/smartcontractkit/chainlink/services"
	"github.com/smartcontractkit/chainlink/store/models"
	"github.com/stretchr/testify/assert"
)

func TestRunningJob(t *testing.T) {
	t.Parallel()
	store, cleanup := cltest.NewStore()
	defer cleanup()

	job := models.NewJob()
	job.Tasks = []models.Task{models.Task{Type: "NoOp"}}

	run := job.NewRun()
	services.ResumeRun(run, store)

	store.One("ID", run.ID, &run)
	assert.Equal(t, models.StatusCompleted, run.Status)
	assert.Equal(t, models.StatusCompleted, run.TaskRuns[0].Status)
}

func TestJobTransitionToPending(t *testing.T) {
	t.Parallel()
	store, cleanup := cltest.NewStore()
	defer cleanup()

	job := models.NewJob()
	job.Tasks = []models.Task{models.Task{Type: "NoOpPend"}}

	run := job.NewRun()
	services.ResumeRun(run, store)

	store.One("ID", run.ID, &run)
	assert.Equal(t, models.StatusPending, run.Status)
}
