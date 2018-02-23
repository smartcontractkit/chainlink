package services_test

import (
	"testing"
	"time"

	"github.com/onsi/gomega"
	"github.com/smartcontractkit/chainlink/internal/cltest"
	"github.com/smartcontractkit/chainlink/services"
	"github.com/smartcontractkit/chainlink/store/models"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zapcore"
	null "gopkg.in/guregu/null.v3"
)

func TestScheduler_Start_LoadingRecurringJobs(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore()
	defer cleanup()

	jobWCron := cltest.NewJob()
	jobWCron.Initiators = []models.Initiator{{Type: models.InitiatorCron, Schedule: "* * * * * *"}}
	assert.Nil(t, store.SaveJob(&jobWCron))
	jobWoCron := cltest.NewJob()
	assert.Nil(t, store.SaveJob(&jobWoCron))

	sched := services.NewScheduler(store)
	assert.Nil(t, sched.Start())
	defer sched.Stop()

	cltest.WaitForRuns(t, jobWCron, store, 1)
	cltest.WaitForRuns(t, jobWoCron, store, 0)
}

func TestRecurring_AddJob(t *testing.T) {
	nullTime := cltest.NullTime(nil)
	pastTime := cltest.NullTime("2000-01-01T00:00:00.000Z")
	futureTime := cltest.NullTime("3000-01-01T00:00:00.000Z")
	tests := []struct {
		name        string
		startAt     null.Time
		endAt       null.Time
		wantEntries int
		wantRuns    int
	}{
		{"before start at", futureTime, nullTime, 1, 0},
		{"before end at", nullTime, futureTime, 1, 1},
		{"after start at", pastTime, nullTime, 1, 1},
		{"after end at", nullTime, pastTime, 0, 0},
		{"no range", nullTime, nullTime, 1, 1},
		{"start at after end at", futureTime, pastTime, 0, 0},
	}

	store, cleanup := cltest.NewStore()
	defer cleanup()
	for _, tt := range tests {
		test := tt
		t.Run(test.name, func(t *testing.T) {
			r := services.NewRecurring(store)
			cron := cltest.NewMockCron()
			r.Cron = cron
			defer r.Stop()

			j := cltest.NewJobWithSchedule("* * * * *")
			j.StartAt = test.startAt
			j.EndAt = test.endAt

			r.AddJob(j)

			assert.Equal(t, test.wantEntries, len(cron.Entries))

			cron.RunEntries()
			jobRuns := []models.JobRun{}
			assert.Nil(t, store.Where("JobID", j.ID, &jobRuns))
			assert.Equal(t, test.wantRuns, len(jobRuns))
		})
	}
}

func TestScheduler_AddJob_WhenStopped(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore()
	defer cleanup()
	sched := services.NewScheduler(store)
	defer sched.Stop()

	j := cltest.NewJobWithSchedule("* * * * *")
	assert.Nil(t, store.SaveJob(&j))
	sched.AddJob(j)

	cltest.WaitForRuns(t, j, store, 0)
}

func TestOneTime_RunJobAt(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore()
	defer cleanup()

	ot := services.OneTime{
		Clock: &cltest.NeverClock{},
		Store: store,
	}
	ot.Start()
	j := cltest.NewJob()
	assert.Nil(t, store.SaveJob(&j))

	var finished bool
	go func() {
		ot.RunJobAt(models.Time{time.Now().Add(time.Hour)}, j)
		finished = true
	}()

	ot.Stop()

	gomega.NewGomegaWithT(t).Eventually(func() bool {
		return finished
	}).Should(gomega.Equal(true))
	jobRuns := []models.JobRun{}
	assert.Nil(t, store.Where("JobID", j.ID, &jobRuns))
	assert.Equal(t, 0, len(jobRuns))
}

func TestScheduler_Start_AddingUnstartedJob(t *testing.T) {
	logs := cltest.ObserveLogs()

	store, cleanupStore := cltest.NewStore()
	clock := cltest.UseSettableClock(store)

	startAt := cltest.ParseISO8601("3000-01-01T00:00:00.000Z")
	j := cltest.NewJobWithSchedule("* * * * *")
	j.StartAt = cltest.NullableTime(startAt)
	assert.Nil(t, store.Save(&j))

	sched := services.NewScheduler(store)
	assert.Nil(t, sched.Start())
	defer sched.Stop()
	defer cleanupStore()

	gomega.NewGomegaWithT(t).Consistently(func() int {
		runs, err := store.JobRunsFor(j.ID)
		assert.Nil(t, err)
		return len(runs)
	}, (2 * time.Second)).Should(gomega.Equal(0))

	clock.SetTime(startAt)

	cltest.WaitForRuns(t, j, store, 2)

	for _, log := range logs.All() {
		assert.True(t, log.Level <= zapcore.WarnLevel)
	}
}
