package services_test

import (
	"testing"
	"time"

	"github.com/h2non/gock"
	"github.com/smartcontractkit/chainlink/internal/cltest"
	"github.com/smartcontractkit/chainlink/services"
	"github.com/smartcontractkit/chainlink/store/models"
	"github.com/stretchr/testify/assert"
	null "gopkg.in/guregu/null.v3"
)

func TestJobRunner_ExecuteRun(t *testing.T) {
	store, cleanup := cltest.NewStore()
	defer cleanup()

	tests := []struct {
		name       string
		input      string
		runResult  string
		wantStatus string
		wantOutput string
	}{
		{"success", `{}`, `{"output":{"value":"100"}}`, models.StatusCompleted,
			`{"value":"100"}`},
		{"errored", `{}`, `{"error":"too much"}`, models.StatusErrored, `{}`},
		{"errored with a value", `{}`, `{"error":"too much", "output":{"value":"99"}}`, models.StatusErrored,
			`{"value":"99"}`},
		{"overriding bridge type params", `{"url":"http://unsafe.com/hack"}`, `{"output":{"value":"100"}}`, models.StatusCompleted,
			`{"value":"100"}`},
	}

	bt := cltest.NewBridgeType("auctionBidding", "https://dbay.eth/api")
	assert.Nil(t, store.Save(bt))

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			gock.New("https://dbay.eth").
				Post("/api").
				JSON(test.input).
				Reply(200).
				JSON(test.runResult)

			job := models.NewJob()
			job.Tasks = []models.Task{{Type: bt.Name}, {Type: "noop"}}
			assert.Nil(t, store.Save(&job))

			run := job.NewRun()
			input := cltest.JSONFromString(test.input)
			run, err := services.ExecuteRun(run, store, input)
			assert.Nil(t, err)

			store.One("ID", run.ID, &run)
			assert.Equal(t, test.wantStatus, run.Status)
			assert.Equal(t, test.wantOutput, run.Result.Output.String())

			tr1 := run.TaskRuns[0]
			assert.Equal(t, test.wantStatus, tr1.Status)
			assert.Equal(t, test.wantOutput, tr1.Result.Output.String())

			if test.wantStatus == models.StatusCompleted {
				tr2 := run.TaskRuns[1]
				assert.Equal(t, test.wantOutput, tr2.Result.Output.String())
			}
		})
	}
}

func TestJobRunner_ExecuteRun_TransitionToPending(t *testing.T) {
	t.Parallel()
	store, cleanup := cltest.NewStore()
	defer cleanup()

	job := models.NewJob()
	job.Tasks = []models.Task{models.Task{Type: "NoOpPend"}}

	run, err := services.ExecuteRun(job.NewRun(), store, models.JSON{})
	assert.Nil(t, err)

	store.One("ID", run.ID, &run)
	assert.Equal(t, models.StatusPending, run.Status)
}

func TestJobRunner_BeginRun(t *testing.T) {
	t.Parallel()

	pastTime := cltest.ParseNullableTime("2000-01-01T00:00:00.000Z")
	futureTime := cltest.ParseNullableTime("3000-01-01T00:00:00.000Z")
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
			assert.Nil(t, store.SaveJob(&job))

			_, err := services.BeginRun(job, store, models.JSON{})

			if test.errored {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}
			jrs, err := store.JobRunsFor(job)
			assert.Nil(t, err)
			assert.Equal(t, test.runCount, len(jrs))
		})
	}
}

func TestJobRunner_BuildRun(t *testing.T) {
	t.Parallel()

	pastTime := cltest.ParseNullableTime("2000-01-01T00:00:00.000Z")
	futureTime := cltest.ParseNullableTime("3000-01-01T00:00:00.000Z")
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
			assert.Nil(t, store.SaveJob(&job))

			_, err := services.BuildRun(job, store)

			if test.errored {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}
