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

	jobWCron, _ := cltest.NewJobWithSchedule("* * * * * *")
	assert.Nil(t, store.SaveJob(&jobWCron))
	jobWoCron := cltest.NewJob()
	assert.Nil(t, store.SaveJob(&jobWoCron))

	sched := services.NewScheduler(store)
	assert.Nil(t, sched.Start())
	defer sched.Stop()

	cltest.WaitForRuns(t, jobWCron, store, 1)
	cltest.WaitForRuns(t, jobWoCron, store, 0)
}

func TestScheduler_AddJob_WhenStopped(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore()
	defer cleanup()
	sched := services.NewScheduler(store)
	defer sched.Stop()

	j, _ := cltest.NewJobWithSchedule("* * * * *")
	assert.Nil(t, store.SaveJob(&j))
	sched.AddJob(j)

	cltest.WaitForRuns(t, j, store, 0)
}

func TestScheduler_Start_AddingUnstartedJob(t *testing.T) {
	logs := cltest.ObserveLogs()

	store, cleanupStore := cltest.NewStore()
	clock := cltest.UseSettableClock(store)

	startAt := cltest.ParseISO8601("3000-01-01T00:00:00.000Z")
	j, _ := cltest.NewJobWithSchedule("* * * * *")
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

			j, _ := cltest.NewJobWithSchedule("* * * * *")
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

func TestOneTime_AddJob(t *testing.T) {
	nullTime := cltest.NullTime(nil)
	pastTime := cltest.NullTime("2000-01-01T00:00:00.000Z")
	futureTime := cltest.NullTime("3000-01-01T00:00:00.000Z")
	pastRunTime := time.Now().Add(time.Hour * -1)
	tests := []struct {
		name     string
		startAt  null.Time
		endAt    null.Time
		runAt    time.Time
		wantRuns bool
	}{
		{"run at before start at", futureTime, nullTime, pastRunTime, false},
		{"run at before end at", nullTime, futureTime, pastRunTime, true},
		{"run at after start at", pastTime, nullTime, pastRunTime, true},
		{"run at after end at", nullTime, pastTime, pastRunTime, false},
		{"no range", nullTime, nullTime, pastRunTime, true},
		{"start at after end at", futureTime, pastTime, pastRunTime, false},
	}

	store, cleanup := cltest.NewStore()
	defer cleanup()
	for _, tt := range tests {
		test := tt
		t.Run(test.name, func(t *testing.T) {
			ot := services.OneTime{
				Clock: store.Clock,
				Store: store,
			}

			j, _ := cltest.NewJobWithRunAtInitiator(test.runAt)
			assert.Nil(t, store.SaveJob(&j))

			j.StartAt = test.startAt
			j.EndAt = test.endAt

			ot.AddJob(j)

			gomega.NewGomegaWithT(t).Eventually(func() bool {
				jobRuns := []models.JobRun{}
				completed := false
				assert.Nil(t, store.Where("JobID", j.ID, &jobRuns))
				if (len(jobRuns) > 0) && (jobRuns[0].Status == models.RunStatusCompleted) {
					completed = true
				}
				return completed
			}).Should(gomega.Equal(test.wantRuns))
		})
	}
}

func TestOneTime_RunJobAt_StopJobBeforeExecution(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore()
	defer cleanup()

	ot := services.OneTime{
		Clock: &cltest.NeverClock{},
		Store: store,
	}
	ot.Start()
	j, initr := cltest.NewJobWithRunAtInitiator(time.Now().Add(time.Hour))
	assert.Nil(t, store.SaveJob(&j))

	var finished bool
	go func() {
		ot.RunJobAt(initr, j)
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

func TestOneTime_RunJobAt_ExecuteLateJob(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore()
	defer cleanup()

	ot := services.OneTime{
		Clock: store.Clock,
		Store: store,
	}
	j, initr := cltest.NewJobWithRunAtInitiator(time.Now().Add(time.Hour * -1))
	assert.Nil(t, store.SaveJob(&j))
	initr.ID = j.Initiators[0].ID
	initr.JobID = j.ID

	var finished bool
	go func() {
		ot.RunJobAt(initr, j)
		finished = true
	}()

	gomega.NewGomegaWithT(t).Eventually(func() bool {
		return finished
	}).Should(gomega.Equal(true))
	jobRuns := []models.JobRun{}
	assert.Nil(t, store.Where("JobID", j.ID, &jobRuns))
	assert.Equal(t, 1, len(jobRuns))
}

func TestOneTime_RunJobAt_RunTwice(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore()
	defer cleanup()

	ot := services.OneTime{
		Clock: store.Clock,
		Store: store,
	}

	j, _ := cltest.NewJobWithRunAtInitiator(time.Now())
	assert.Nil(t, store.SaveJob(&j))

	var initrs []models.Initiator
	store.Where("JobID", j.ID, &initrs)

	ot.RunJobAt(initrs[0], j)
	j2, err := ot.Store.FindJob(j.ID)
	assert.Nil(t, err)

	var initrs2 []models.Initiator
	store.Where("JobID", j.ID, &initrs2)

	ot.RunJobAt(initrs2[0], j2)

	jobRuns := []models.JobRun{}
	assert.Nil(t, store.Where("JobID", j.ID, &jobRuns))
	assert.Equal(t, 1, len(jobRuns))
}

func TestOneTime_RunJobAt_UnstartedRun(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore()
	defer cleanup()

	ot := services.OneTime{
		Clock: store.Clock,
		Store: store,
	}

	j, _ := cltest.NewJobWithRunAtInitiator(time.Now())
	j.EndAt = cltest.NullTime("2000-01-01T00:10:00.000Z")
	assert.Nil(t, store.SaveJob(&j))

	ot.RunJobAt(j.Initiators[0], j)

	var initrs2 []models.Initiator
	store.Where("JobID", j.ID, &initrs2)

	assert.Equal(t, false, initrs2[0].Ran)
}
