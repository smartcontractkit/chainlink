package services_test

import (
	"testing"
	"time"

	"github.com/onsi/gomega"
	"github.com/smartcontractkit/chainlink/core/services"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/tools/cltest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tevino/abool"
	"go.uber.org/zap/zapcore"
	null "gopkg.in/guregu/null.v3"
)

func TestScheduler_Start_LoadingRecurringJobs(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore()
	defer cleanup()

	jobWCron := cltest.NewJobWithSchedule("* * * * * *")
	require.NoError(t, store.CreateJob(&jobWCron))
	jobWoCron := cltest.NewJob()
	require.NoError(t, store.CreateJob(&jobWoCron))

	sched := services.NewScheduler(store)
	require.NoError(t, sched.Start())
	defer sched.Stop()

	cltest.WaitForRunsAtLeast(t, jobWCron, store, 1)
	cltest.WaitForRuns(t, jobWoCron, store, 0)
}

func TestScheduler_AddJob_WhenStopped(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore()
	defer cleanup()
	sched := services.NewScheduler(store)
	defer sched.Stop()

	j := cltest.NewJobWithSchedule("* * * * *")
	assert.Nil(t, store.CreateJob(&j))
	sched.AddJob(j)

	cltest.WaitForRuns(t, j, store, 0)
}

func TestScheduler_Start_AddingUnstartedJob(t *testing.T) {
	logs := cltest.ObserveLogs()

	store, cleanupStore := cltest.NewStore()
	defer cleanupStore()
	clock := cltest.UseSettableClock(store)

	startAt := cltest.ParseISO8601("3000-01-01T00:00:00.000Z")
	j := cltest.NewJobWithSchedule("* * * * *")
	j.StartAt = cltest.NullableTime(startAt)
	assert.Nil(t, store.CreateJob(&j))

	sched := services.NewScheduler(store)
	defer sched.Stop()
	assert.Nil(t, sched.Start())

	gomega.NewGomegaWithT(t).Consistently(func() int {
		runs, err := store.JobRunsFor(j.ID)
		assert.NoError(t, err)
		return len(runs)
	}, (2 * time.Second)).Should(gomega.Equal(0))

	clock.SetTime(startAt)

	cltest.WaitForRunsAtLeast(t, j, store, 2)

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

			job := cltest.NewJobWithSchedule("* * * * *")
			job.StartAt = test.startAt
			job.EndAt = test.endAt

			require.NoError(t, store.CreateJob(&job))
			r.AddJob(job)

			assert.Equal(t, test.wantEntries, len(cron.Entries))

			cron.RunEntries()
			jobRuns, err := store.JobRunsFor(job.ID)
			assert.NoError(t, err)
			assert.Equal(t, test.wantRuns, len(jobRuns))
		})
	}
}

func TestRecurring_AddJob_Archived(t *testing.T) {
	store, cleanup := cltest.NewStore()
	defer cleanup()
	r := services.NewRecurring(store)
	cron := cltest.NewMockCron()
	r.Cron = cron
	defer r.Stop()

	job := cltest.NewJobWithSchedule("* * * * *")
	require.NoError(t, store.CreateJob(&job))
	r.AddJob(job)

	cron.RunEntries()
	count, err := store.Unscoped().JobRunsCountFor(job.ID)
	require.NoError(t, err)
	assert.Equal(t, 1, count)

	require.NoError(t, store.ArchiveJob(job.ID))
	cron.RunEntries()

	count, err = store.Unscoped().JobRunsCountFor(job.ID)
	require.NoError(t, err)
	assert.Equal(t, 1, count)
}

func TestOneTime_AddJob(t *testing.T) {
	nullTime := cltest.NullTime(nil)
	pastTime := cltest.NullTime("2000-01-01T00:00:00.000Z")
	futureTime := cltest.NullTime("3000-01-01T00:00:00.000Z")
	pastRunTime := time.Now().Add(time.Hour * -1)
	tests := []struct {
		name          string
		startAt       null.Time
		endAt         null.Time
		runAt         time.Time
		wantCompleted bool
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
	jobRunner, cleanup := cltest.NewJobRunner(store)
	defer cleanup()
	jobRunner.Start()

	for _, tt := range tests {
		test := tt
		t.Run(test.name, func(t *testing.T) {
			ot := services.OneTime{
				Clock: store.Clock,
				Store: store,
			}
			require.NoError(t, ot.Start())
			defer ot.Stop()

			j := cltest.NewJobWithRunAtInitiator(test.runAt)
			require.Nil(t, store.CreateJob(&j))

			j.StartAt = test.startAt
			j.EndAt = test.endAt

			ot.AddJob(j)

			tester := func() bool {
				completed := false
				jobRuns, err := store.JobRunsFor(j.ID)
				require.NoError(t, err)
				if (len(jobRuns) > 0) && (jobRuns[0].Status == models.RunStatusCompleted) {
					completed = true
				}
				return completed
			}

			if test.wantCompleted {
				gomega.NewGomegaWithT(t).Eventually(tester).Should(gomega.Equal(true))
			} else {
				gomega.NewGomegaWithT(t).Consistently(tester).Should(gomega.Equal(false))
			}

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
	j := cltest.NewJobWithRunAtInitiator(time.Now().Add(time.Hour))
	assert.Nil(t, store.CreateJob(&j))
	initr := j.Initiators[0]

	finished := abool.New()
	go func() {
		ot.RunJobAt(initr, j)
		finished.Set()
	}()

	ot.Stop()

	gomega.NewGomegaWithT(t).Eventually(func() bool {
		return finished.IsSet()
	}).Should(gomega.Equal(true))
	jobRuns, err := store.JobRunsFor(j.ID)
	assert.NoError(t, err)
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
	j := cltest.NewJobWithRunAtInitiator(time.Now().Add(time.Hour * -1))
	assert.Nil(t, store.CreateJob(&j))
	initr := j.Initiators[0]
	initr.ID = j.Initiators[0].ID
	initr.JobSpecID = j.ID

	finished := abool.New()
	go func() {
		ot.RunJobAt(initr, j)
		finished.Set()
	}()

	gomega.NewGomegaWithT(t).Eventually(func() bool {
		return finished.IsSet()
	}).Should(gomega.Equal(true))
	jobRuns, err := store.JobRunsFor(j.ID)
	assert.NoError(t, err)
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

	j := cltest.NewJobWithRunAtInitiator(time.Now())
	assert.NoError(t, store.CreateJob(&j))
	ot.RunJobAt(j.Initiators[0], j)

	j2, err := ot.Store.FindJob(j.ID)
	require.NoError(t, err)
	require.Len(t, j2.Initiators, 1)
	ot.RunJobAt(j2.Initiators[0], j2)

	jobRuns, err := store.JobRunsFor(j.ID)
	assert.NoError(t, err)
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

	j := cltest.NewJobWithRunAtInitiator(time.Now())
	j.EndAt = cltest.NullTime("2000-01-01T00:10:00.000Z")
	assert.NoError(t, store.CreateJob(&j))

	ot.RunJobAt(j.Initiators[0], j)

	j2, err := store.FindJob(j.ID)
	require.NoError(t, err)
	require.Len(t, j2.Initiators, 1)
	assert.Equal(t, false, j2.Initiators[0].Ran)
}

func TestOneTime_RunJobAt_ArchivedRun(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore()
	defer cleanup()

	ot := services.OneTime{
		Clock: cltest.InstantClock{},
		Store: store,
	}

	j := cltest.NewJobWithRunAtInitiator(time.Now())
	j.EndAt = cltest.NullTime("2000-01-01T00:10:00.000Z")
	require.NoError(t, store.CreateJob(&j))
	require.NoError(t, store.ArchiveJob(j.ID))

	ot.RunJobAt(j.Initiators[0], j)

	unscoped := store.Unscoped()
	j2, err := unscoped.FindJob(j.ID)
	require.NoError(t, err)
	require.Len(t, j2.Initiators, 1)
	assert.Equal(t, false, j2.Initiators[0].Ran)
	count, err := unscoped.JobRunsCountFor(j.ID)
	require.NoError(t, err)
	assert.Equal(t, 0, count)
}
