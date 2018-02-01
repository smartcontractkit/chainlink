package services_test

import (
	"testing"
	"time"

	. "github.com/onsi/gomega"
	"github.com/smartcontractkit/chainlink/internal/cltest"
	"github.com/smartcontractkit/chainlink/services"
	"github.com/smartcontractkit/chainlink/store/models"
	"github.com/smartcontractkit/chainlink/utils"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zapcore"
	null "gopkg.in/guregu/null.v3"
)

func TestLoadingSavedSchedules(t *testing.T) {
	t.Parallel()
	RegisterTestingT(t)
	store, cleanup := cltest.NewStore()
	defer cleanup()

	j := cltest.NewJob()
	j.Initiators = []models.Initiator{{Type: models.InitiatorCron, Schedule: "* * * * *"}}
	jobWoCron := models.NewJob()
	assert.Nil(t, store.SaveJob(j))
	assert.Nil(t, store.SaveJob(jobWoCron))

	sched := services.NewScheduler(store)
	err := sched.Start()
	assert.Nil(t, err)

	jobRuns := []models.JobRun{}
	Eventually(func() []models.JobRun {
		store.Where("JobID", j.ID, &jobRuns)
		return jobRuns
	}).Should(cltest.HaveLenAtLeast(1))

	sched.Stop()
	err = store.Where("JobID", jobWoCron.ID, &jobRuns)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(jobRuns), "No jobs should be created without the scheduler")
}

func TestRecurringAddJob(t *testing.T) {
	t.Parallel()
	RegisterTestingT(t)

	store, cleanup := cltest.NewStore()
	defer cleanup()

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

	for _, test := range tests {
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
			store.Where("JobID", j.ID, &jobRuns)
			assert.Equal(t, test.wantRuns, len(jobRuns))
		})
	}
}

func TestAddScheduledJobWhenStopped(t *testing.T) {
	t.Parallel()
	RegisterTestingT(t)
	store, cleanup := cltest.NewStore()
	defer cleanup()
	sched := services.NewScheduler(store)
	defer sched.Stop()

	j := cltest.NewJobWithSchedule("* * * * *")
	assert.Nil(t, store.SaveJob(j))
	sched.AddJob(j)

	jobRuns := []models.JobRun{}
	Consistently(func() []models.JobRun {
		store.Where("JobID", j.ID, &jobRuns)
		return jobRuns
	}).Should(HaveLen(0))

}

func TestOneTimeRunJobAt(t *testing.T) {
	RegisterTestingT(t)
	t.Parallel()

	store, cleanup := cltest.NewStore()
	defer cleanup()

	ot := services.OneTime{
		Clock: &cltest.NeverClock{},
		Store: store,
	}
	ot.Start()
	j := cltest.NewJob()
	assert.Nil(t, store.SaveJob(j))

	var finished bool
	go func() {
		ot.RunJobAt(models.Time{time.Now().Add(time.Hour)}, j)
		finished = true
	}()

	ot.Stop()

	Eventually(func() bool {
		return finished
	}).Should(Equal(true))
	jobRuns := []models.JobRun{}
	assert.Nil(t, store.Where("JobID", j.ID, &jobRuns))
	assert.Equal(t, 0, len(jobRuns))
}

func TestSchedulerAddingUnstartedJob(t *testing.T) {
	RegisterTestingT(t)
	logs := cltest.ObserveLogs()

	store, cleanupStore := cltest.NewStore()
	clock := cltest.UseSettableClock(store)

	startAt := utils.ParseISO8601("3000-01-01T00:00:00.000Z")
	j := cltest.NewJobWithSchedule("* * * * *")
	j.StartAt = utils.NullableTime(startAt)
	assert.Nil(t, store.Save(j))

	sched := services.NewScheduler(store)
	assert.Nil(t, sched.Start())
	defer sched.Stop()
	defer cleanupStore()

	Consistently(func() int {
		runs, err := store.JobRunsFor(j)
		assert.Nil(t, err)
		return len(runs)
	}, (2 * time.Second)).Should(Equal(0))

	clock.SetTime(startAt)

	Eventually(func() int {
		runs, err := store.JobRunsFor(j)
		assert.Nil(t, err)
		return len(runs)
	}).Should(Equal(2))

	for _, log := range logs.All() {
		assert.True(t, log.Level <= zapcore.WarnLevel)
	}
}
