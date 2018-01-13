package models_test

import (
	"fmt"
	"testing"

	"github.com/smartcontractkit/chainlink/internal/cltest"
	"github.com/smartcontractkit/chainlink/services"
	"github.com/smartcontractkit/chainlink/store/models"
	"github.com/stretchr/testify/assert"
)

func TestRetrievingJobRunsWithErrorsFromDB(t *testing.T) {
	t.Parallel()
	store, cleanup := cltest.NewStore()
	defer cleanup()

	job := models.NewJob()
	jr := job.NewRun()
	jr.Result = models.RunResultWithError(fmt.Errorf("bad idea"))
	err := store.Save(&jr)
	assert.Nil(t, err)

	run := models.JobRun{}
	err = store.One("ID", jr.ID, &run)
	assert.Nil(t, err)
	assert.True(t, run.Result.HasError())
	assert.Equal(t, "bad idea", run.Result.Error())
}

func TestTaskRunsToRun(t *testing.T) {
	t.Parallel()
	store, cleanup := cltest.NewStore()
	defer cleanup()

	j := models.NewJob()
	j.Tasks = []models.Task{
		{Type: "NoOp"},
		{Type: "NoOpPend"},
		{Type: "NoOp"},
	}
	assert.Nil(t, store.SaveJob(j))
	jr := j.NewRun()
	assert.Equal(t, jr.TaskRuns, jr.UnfinishedTaskRuns())

	jr, err := services.StartJob(jr, store)
	assert.Nil(t, err)
	assert.Equal(t, jr.TaskRuns[1:], jr.UnfinishedTaskRuns())
}
