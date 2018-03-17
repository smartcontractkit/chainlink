package services_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink/internal/cltest"
	"github.com/smartcontractkit/chainlink/services"
	"github.com/smartcontractkit/chainlink/store/models"
	"github.com/stretchr/testify/assert"
	null "gopkg.in/guregu/null.v3"
)

func TestJobRunner_ExecuteRun(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		input      string
		runResult  string
		wantStatus string
		wantData   string
	}{
		{"success", `{}`, `{"data":{"value":"100"}}`,
			models.StatusCompleted, `{"value":"100"}`},
		{"errored", `{}`, `{"error":"too much"}`,
			models.StatusErrored, `{}`},
		{"errored with a value", `{}`, `{"error":"too much", "data":{"value":"99"}}`,
			models.StatusErrored, `{"value":"99"}`},
		{"overriding bridge type params", `{"data":{"url":"hack"},"url":"hack"}`, `{"data":{"value":"100"}}`,
			models.StatusCompleted, `{"value":"100"}`},
		{"type parameter does not override", `{"data":{"type":"0"},"type":"0"}`, `{"data":{"value":"100"}}`,
			models.StatusCompleted, `{"value":"100"}`},
	}

	store, cleanup := cltest.NewStore()
	defer cleanup()

	for _, tt := range tests {
		test := tt
		t.Run(test.name, func(t *testing.T) {

			var run models.JobRun
			mockServer, cleanup := cltest.NewHTTPMockServer(t, 200, "POST", test.runResult,
				func(body string) {
					want := fmt.Sprintf(`{"id":"%v","data":%v}`, run.ID, test.input)
					assert.JSONEq(t, want, body)
				})
			defer cleanup()
			bt := cltest.NewBridgeType("auctionBidding", mockServer.URL)
			assert.Nil(t, store.Save(&bt))

			job := models.NewJob()
			job.Tasks = []models.TaskSpec{{Type: bt.Name}, {Type: "noop"}}
			assert.Nil(t, store.Save(&job))

			run = job.NewRun()
			input := models.RunResult{Data: cltest.JSONFromString(test.input)}
			run, err := services.ExecuteRun(run, store, input)
			assert.Nil(t, err)

			store.One("ID", run.ID, &run)
			assert.Equal(t, test.wantStatus, run.Status)
			assert.Equal(t, test.wantData, run.Result.Data.String())

			tr1 := run.TaskRuns[0]
			assert.Equal(t, test.wantStatus, tr1.Status)
			assert.Equal(t, test.wantData, tr1.Result.Data.String())

			if test.wantStatus == models.StatusCompleted {
				tr2 := run.TaskRuns[1]
				assert.Equal(t, test.wantData, tr2.Result.Data.String())
				assert.True(t, run.CompletedAt.Valid)
			}
		})
	}
}

func TestJobRunner_ExecuteRun_TransitionToPending(t *testing.T) {
	t.Parallel()
	store, cleanup := cltest.NewStore()
	defer cleanup()

	job := models.NewJob()
	job.Tasks = []models.TaskSpec{{Type: "NoOpPend"}}

	run, err := services.ExecuteRun(job.NewRun(), store, models.RunResult{})
	assert.Nil(t, err)

	store.One("ID", run.ID, &run)
	assert.Equal(t, models.StatusPending, run.Status)
}

func TestJobRunner_BeginRun(t *testing.T) {
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

	store, cleanup := cltest.NewStore()
	defer cleanup()

	for _, tt := range tests {
		test := tt
		t.Run(test.name, func(t *testing.T) {
			job := cltest.NewJob()
			job.StartAt = test.startAt
			job.EndAt = test.endAt
			assert.Nil(t, store.SaveJob(&job))

			_, err := services.BeginRun(job, store, models.RunResult{})

			if test.errored {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}
			jrs, err := store.JobRunsFor(job.ID)
			assert.Nil(t, err)
			assert.Equal(t, test.runCount, len(jrs))
		})
	}
}

func TestJobRunner_BuildRun(t *testing.T) {
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

	store, cleanup := cltest.NewStore()
	defer cleanup()
	clock := cltest.UseSettableClock(store)
	clock.SetTime(time.Now())

	for _, tt := range tests {
		test := tt
		t.Run(test.name, func(t *testing.T) {
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
