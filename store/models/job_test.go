package models_test

import (
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink/adapters"
	"github.com/smartcontractkit/chainlink/internal/cltest"
	"github.com/smartcontractkit/chainlink/store/models"
	"github.com/smartcontractkit/chainlink/utils"
	"github.com/stretchr/testify/assert"
	null "gopkg.in/guregu/null.v3"
)

func TestJobSave(t *testing.T) {
	t.Parallel()
	store, cleanup := cltest.NewStore()
	defer cleanup()

	j1 := cltest.NewJobWithSchedule("* * * * 7")
	assert.Nil(t, store.SaveJob(j1))

	store.Save(j1)
	j2, _ := store.FindJob(j1.ID)
	assert.Equal(t, j1.Initiators[0].Schedule, j2.Initiators[0].Schedule)
}

func TestJobNewRun(t *testing.T) {
	t.Parallel()

	job := cltest.NewJobWithSchedule("1 * * * *")
	job.Tasks = []models.Task{models.Task{Type: "NoOp"}}

	newRun := job.NewRun()
	assert.Equal(t, job.ID, newRun.JobID)
	assert.Equal(t, 1, len(newRun.TaskRuns))
	assert.Equal(t, "NoOp", job.Tasks[0].Type)
	assert.Nil(t, job.Tasks[0].Params)
	adapter, _ := adapters.For(job.Tasks[0])
	assert.NotNil(t, adapter)
}

func TestTimeDurationFromNow(t *testing.T) {
	future := models.Time{time.Now().Add(time.Second)}
	duration := future.DurationFromNow()
	assert.True(t, 0 < duration)
}

func TestJobEnded(t *testing.T) {
	t.Parallel()

	endAt := utils.ParseISO8601("3000-01-01T00:00:00.000Z")
	endAtNullable := null.Time{Time: endAt, Valid: true}

	tests := []struct {
		name    string
		endAt   null.Time
		current time.Time
		want    bool
	}{
		{"no end at", null.Time{Valid: false}, endAt, false},
		{"before end at", endAtNullable, endAt.Add(-time.Nanosecond), false},
		{"at end at", endAtNullable, endAt, false},
		{"after end at", endAtNullable, endAt.Add(time.Nanosecond), true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			job := cltest.NewJob()
			job.EndAt = test.endAt

			assert.Equal(t, test.want, job.Ended(test.current))
		})
	}
}
