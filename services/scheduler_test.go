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
)

func TestLoadingSavedSchedules(t *testing.T) {
	t.Parallel()
	RegisterTestingT(t)
	store, cleanup := cltest.NewStore()
	defer cleanup()

	j := cltest.NewJob()
	j.Initiators = []models.Initiator{{Type: "cron", Schedule: "* * * * *"}}
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

func TestAddScheduledJob(t *testing.T) {
	t.Parallel()
	RegisterTestingT(t)
	store, cleanup := cltest.NewStore()
	defer cleanup()

	sched := services.NewScheduler(store)
	sched.Start()
	defer sched.Stop()

	j := cltest.NewJobWithSchedule("* * * * *")
	err := store.SaveJob(j)
	assert.Nil(t, err)
	sched.AddJob(j)

	jobRuns := []models.JobRun{}
	Eventually(func() []models.JobRun {
		err = store.Where("JobID", j.ID, &jobRuns)
		assert.Nil(t, err)
		return jobRuns
	}).Should(cltest.HaveLenAtLeast(1))
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

	assert.Nil(t, sched.Start())
	Eventually(func() []models.JobRun {
		store.Where("JobID", j.ID, &jobRuns)
		return jobRuns
	}).Should(cltest.HaveLenAtLeast(1))
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
	t.Parallel()
	RegisterTestingT(t)

	store, cleanup := cltest.NewStore()
	defer cleanup()
	clock := cltest.UseSettableClock(store)

	startAt := utils.ParseISO8601("3000-01-01T00:00:00.000Z")
	j := cltest.NewJobWithSchedule("* * * * *")
	j.StartAt = utils.NullableTime(startAt)
	assert.Nil(t, store.Save(j))

	sched := services.NewScheduler(store)
	assert.Nil(t, sched.Start())

	sched.AddJob(j)
	Consistently(func() int {
		runs, err := store.JobRunsFor(j)
		assert.Nil(t, err)
		return len(runs)
	}).Should(Equal(0))

	clock.SetTime(startAt)

	Eventually(func() int {
		runs, err := store.JobRunsFor(j)
		assert.Nil(t, err)
		return len(runs)
	}).Should(Equal(2))
}
