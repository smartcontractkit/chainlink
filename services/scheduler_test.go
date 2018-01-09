package services_test

import (
	"testing"
	"time"

	. "github.com/onsi/gomega"
	"github.com/smartcontractkit/chainlink-go/internal/cltest"
	"github.com/smartcontractkit/chainlink-go/services"
	"github.com/smartcontractkit/chainlink-go/store/models"
	"github.com/stretchr/testify/assert"
)

func TestLoadingSavedSchedules(t *testing.T) {
	t.Parallel()
	RegisterTestingT(t)
	store := cltest.NewStore()
	defer cltest.CleanUpStore(store)

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
	}).Should(HaveLen(1))

	sched.Stop()
	err = store.Where("JobID", jobWoCron.ID, &jobRuns)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(jobRuns), "No jobs should be created without the scheduler")
}

func TestAddJob(t *testing.T) {
	t.Parallel()
	RegisterTestingT(t)
	store := cltest.NewStore()

	sched := services.NewScheduler(store)
	sched.Start()

	defer sched.Stop()
	defer cltest.CleanUpStore(store)

	j := cltest.NewJobWithSchedule("* * * * *")
	err := store.SaveJob(j)
	assert.Nil(t, err)
	sched.AddJob(j)

	jobRuns := []models.JobRun{}
	Eventually(func() []models.JobRun {
		err = store.Where("JobID", j.ID, &jobRuns)
		assert.Nil(t, err)
		return jobRuns
	}).Should(HaveLen(1))
}

func TestAddJobWhenStopped(t *testing.T) {
	t.Parallel()
	RegisterTestingT(t)
	store := cltest.NewStore()
	sched := services.NewScheduler(store)

	defer sched.Stop()
	defer cltest.CleanUpStore(store)

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
	}).Should(HaveLen(1))
}

func TestOneTimeRunJobAt(t *testing.T) {
	RegisterTestingT(t)
	t.Parallel()

	store := cltest.NewStore()
	defer cltest.CleanUpStore(store)

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
