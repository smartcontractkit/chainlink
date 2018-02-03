package services_test

import (
	"testing"
	"time"

	"github.com/h2non/gock"
	"github.com/smartcontractkit/chainlink/internal/cltest"
	"github.com/smartcontractkit/chainlink/services"
	"github.com/smartcontractkit/chainlink/store/models"
	"github.com/smartcontractkit/chainlink/utils"
	"github.com/stretchr/testify/assert"
	null "gopkg.in/guregu/null.v3"
)

func TestRunningJob(t *testing.T) {
	store, cleanup := cltest.NewStore()
	defer cleanup()

	tests := []struct {
		name         string
		runResult    string
		wantedStatus string
	}{
		{"success", `{"output":{"value":"100"}}`, models.StatusCompleted},
		{"errored", `{"error":"too much"}`, models.StatusErrored},
		{"errored with a value", `{"error":"too much", "output":{"value":"99"}}`, models.StatusErrored},
	}

	bt := cltest.NewBridgeType("auctionBidding", "https://dbay.eth/api")
	assert.Nil(t, store.Save(bt))

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			gock.New("https://dbay.eth").Post("/api").Reply(200).JSON(test.runResult)

			job := models.NewJob()
			job.Tasks = []models.Task{{Type: bt.Name}}
			assert.Nil(t, store.Save(job))

			run := job.NewRun()
			assert.Nil(t, services.ExecuteRun(run, store, models.Output{}))

			store.One("ID", run.ID, &run)
			assert.Equal(t, test.wantedStatus, run.TaskRuns[0].Status)
			assert.Equal(t, test.wantedStatus, run.Status)
		})
	}
}

func TestJobTransitionToPending(t *testing.T) {
	t.Parallel()
	store, cleanup := cltest.NewStore()
	defer cleanup()

	job := models.NewJob()
	job.Tasks = []models.Task{models.Task{Type: "NoOpPend"}}

	run := job.NewRun()
	services.ExecuteRun(run, store, models.Output{})

	store.One("ID", run.ID, &run)
	assert.Equal(t, models.StatusPending, run.Status)
}

func TestJobRunnerBeginRun(t *testing.T) {
	t.Parallel()

	pastTime := utils.ParseNullableTime("2000-01-01T00:00:00.000Z")
	futureTime := utils.ParseNullableTime("3000-01-01T00:00:00.000Z")
	nullTime := null.Time{Valid: false}

	tests := []struct {
		name     string
		startAt  null.Time
		endAt    null.Time
		errored  bool
		runCount int
	}{
		{"job not started", futureTime, nullTime, true, 0},
		{"job started", pastTime, futureTime, false, 1},
		{"job with no time range", nullTime, nullTime, false, 1},
		{"job ended", nullTime, pastTime, true, 0},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			store, cleanup := cltest.NewStore()
			defer cleanup()

			job := cltest.NewJob()
			job.StartAt = test.startAt
			job.EndAt = test.endAt
			assert.Nil(t, store.SaveJob(job))

			run, err := services.BeginRun(job, store, models.Output{})

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

func TestJobRunnerBuildRun(t *testing.T) {
	t.Parallel()

	pastTime := utils.ParseNullableTime("2000-01-01T00:00:00.000Z")
	futureTime := utils.ParseNullableTime("3000-01-01T00:00:00.000Z")
	nullTime := null.Time{Valid: false}

	tests := []struct {
		name    string
		startAt null.Time
		endAt   null.Time
		errored bool
	}{
		{"job not started", futureTime, nullTime, true},
		{"job started", pastTime, futureTime, false},
		{"job with no time range", nullTime, nullTime, false},
		{"job ended", nullTime, pastTime, true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			store, cleanup := cltest.NewStore()
			clock := cltest.UseSettableClock(store)
			clock.SetTime(time.Now())
			defer cleanup()

			job := cltest.NewJob()
			job.StartAt = test.startAt
			job.EndAt = test.endAt
			assert.Nil(t, store.SaveJob(job))

			run, err := services.BuildRun(job, store)

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
