package services_test

import (
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink/internal/cltest"
	"github.com/smartcontractkit/chainlink/services"
	"github.com/smartcontractkit/chainlink/store/models"
	"github.com/smartcontractkit/chainlink/utils"
	"github.com/stretchr/testify/assert"
	null "gopkg.in/guregu/null.v3"
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

func TestBeginRun(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		endAt    null.Time
		runCount int
		errored  bool
	}{
		{"job ended", null.Time{Time: utils.ParseISO8601("1999-12-31T23:59:59.000Z"), Valid: true}, 0, true},
		{"job not ended", null.Time{Valid: false}, 1, false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			store, cleanup := cltest.NewStore()
			defer cleanup()

			job := cltest.NewJob()
			job.EndAt = test.endAt
			assert.Nil(t, store.SaveJob(job))

			run, err := services.BeginRun(job, store)

			if test.errored {
				assert.Nil(t, run)
				assert.NotNil(t, err)
			} else {
				assert.NotNil(t, run)
				assert.Nil(t, err)
			}
			jrs, err := store.JobRunsFor(job)
			assert.Nil(t, err)
			assert.Equal(t, test.runCount, len(jrs))
		})
	}
}

func TestNewRun(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		endAt   null.Time
		errored bool
	}{
		{"job ended", null.Time{Time: utils.ParseISO8601("1999-12-31T23:59:59.000Z"), Valid: true}, true},
		{"job not ended", null.Time{Valid: false}, false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			store, cleanup := cltest.NewStore()
			clock := cltest.UseSettableClock(store)
			clock.SetTime(time.Now())
			defer cleanup()

			job := cltest.NewJob()
			job.EndAt = test.endAt
			assert.Nil(t, store.SaveJob(job))

			run, err := services.NewRun(job, store)

			if test.errored {
				assert.Nil(t, run)
				assert.NotNil(t, err)
			} else {
				assert.NotNil(t, run)
				assert.Nil(t, err)
			}
		})
	}
}
