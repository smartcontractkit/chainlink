package services_test

import (
	"testing"

	"github.com/smartcontractkit/chainlink-go/internal/cltest"
	"github.com/smartcontractkit/chainlink-go/services"
	"github.com/smartcontractkit/chainlink-go/store/models"
	"github.com/stretchr/testify/assert"
)

func TestRunningJob(t *testing.T) {
	t.Parallel()
	store := cltest.NewStore()
	defer cltest.CleanUpStore(store)

	job := models.NewJob()
	job.Tasks = []models.Task{models.Task{Type: "NoOp"}}

	run := job.NewRun()
	services.StartJob(run, store)

	store.One("ID", run.ID, &run)
	assert.Equal(t, "completed", run.Status)
	assert.Equal(t, "completed", run.TaskRuns[0].Status)
}

func TestJobTransitionToPending(t *testing.T) {
	t.Parallel()
	store := cltest.NewStore()
	defer cltest.CleanUpStore(store)

	job := models.NewJob()
	job.Tasks = []models.Task{models.Task{Type: "NoOpPend"}}

	run := job.NewRun()
	services.StartJob(run, store)

	store.One("ID", run.ID, &run)
	assert.Equal(t, "pending", run.Status)
}
