package models_test

import (
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink/core/adapters"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	null "gopkg.in/guregu/null.v3"
)

func TestNewJobFromRequest(t *testing.T) {
	t.Parallel()
	store, cleanup := cltest.NewStore()
	defer cleanup()

	j1 := cltest.NewJobWithSchedule("* * * * 7")
	require.NoError(t, store.CreateJob(&j1))

	jsr := models.JobSpecRequest{
		Initiators: cltest.BuildInitiatorRequests(t, j1.Initiators),
		Tasks:      cltest.BuildTaskRequests(t, j1.Tasks),
		StartAt:    j1.StartAt,
		EndAt:      j1.EndAt,
	}

	j2 := models.NewJobFromRequest(jsr)
	require.NoError(t, store.CreateJob(&j2))

	fetched1, err := store.FindJob(j1.ID)
	assert.NoError(t, err)
	assert.Len(t, fetched1.Initiators, 1)
	assert.Len(t, fetched1.Tasks, 1)

	fetched2, err := store.FindJob(j2.ID)
	assert.NoError(t, err)
	assert.Len(t, fetched2.Initiators, 1)
	assert.Len(t, fetched2.Tasks, 1)
}

func TestJobSpec_Save(t *testing.T) {
	t.Parallel()
	store, cleanup := cltest.NewStore()
	defer cleanup()

	j1 := cltest.NewJobWithSchedule("* * * * 7")
	assert.NoError(t, store.CreateJob(&j1))
	initr := j1.Initiators[0]

	j2, err := store.FindJob(j1.ID)
	require.NoError(t, err)
	require.Len(t, j2.Initiators, 1)
	assert.Equal(t, initr.Schedule, j2.Initiators[0].Schedule)
}

func TestJobSpec_NewRun(t *testing.T) {
	t.Parallel()
	store, cleanup := cltest.NewStore()
	defer cleanup()

	job := cltest.NewJobWithSchedule("1 * * * *")
	initr := job.Initiators[0]
	job.Tasks = []models.TaskSpec{cltest.NewTask(t, "NoOp", `{"a":1}`)}

	run := job.NewRun(initr)

	assert.Equal(t, job.ID, run.JobSpecID)
	assert.Equal(t, 1, len(run.TaskRuns))

	taskRun := run.TaskRuns[0]
	assert.Equal(t, "noop", taskRun.TaskSpec.Type.String())
	adapter, _ := adapters.For(taskRun.TaskSpec, store)
	assert.NotNil(t, adapter)
	assert.JSONEq(t, `{"type":"NoOp","a":1}`, taskRun.TaskSpec.Params.String())

	assert.Equal(t, initr, run.Initiator)
}

func TestJobEnded(t *testing.T) {
	t.Parallel()

	endAt := cltest.ParseNullableTime("3000-01-01T00:00:00.000Z")

	tests := []struct {
		name    string
		endAt   null.Time
		current time.Time
		want    bool
	}{
		{"no end at", null.Time{Valid: false}, endAt.Time, false},
		{"before end at", endAt, endAt.Time.Add(-time.Nanosecond), false},
		{"at end at", endAt, endAt.Time, false},
		{"after end at", endAt, endAt.Time.Add(time.Nanosecond), true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			job := cltest.NewJob()
			job.EndAt = test.endAt

			assert.Equal(t, test.want, job.Ended(test.current))
		})
	}
}

func TestJobSpec_Started(t *testing.T) {
	t.Parallel()

	startAt := cltest.ParseNullableTime("3000-01-01T00:00:00.000Z")

	tests := []struct {
		name    string
		startAt null.Time
		current time.Time
		want    bool
	}{
		{"no start at", null.Time{Valid: false}, startAt.Time, true},
		{"before start at", startAt, startAt.Time.Add(-time.Nanosecond), false},
		{"at start at", startAt, startAt.Time, true},
		{"after start at", startAt, startAt.Time.Add(time.Nanosecond), true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			job := cltest.NewJob()
			job.StartAt = test.startAt

			assert.Equal(t, test.want, job.Started(test.current))
		})
	}
}

func TestNewTaskType(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   string
		want    string
		errored bool
	}{
		{"basic", "NoOp", "noop", false},
		{"special characters", "-_-", "-_-", false},
		{"invalid character", "NoOp!", "", true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := models.NewTaskType(test.input)

			if test.errored {
				assert.Error(t, err)
			} else {
				assert.Equal(t, models.TaskType(test.want), got)
				assert.NoError(t, err)
			}
		})
	}
}
