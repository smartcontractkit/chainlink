package services_test

import (
	"testing"

	. "github.com/onsi/gomega"
	"github.com/smartcontractkit/chainlink-go/internal/cltest"
	"github.com/smartcontractkit/chainlink-go/models"
	"github.com/smartcontractkit/chainlink-go/services"
	"github.com/stretchr/testify/assert"
)

func TestLoadingSavedSchedules(t *testing.T) {
	t.Parallel()
	RegisterTestingT(t)
	store := cltest.NewStore()
	defer store.Close()

	j := models.NewJob()
	j.Schedule = models.Schedule{Cron: "* * * * *"}
	jobWoCron := models.NewJob()
	_ = store.Save(&j)
	_ = store.Save(&jobWoCron)

	sched := services.NewScheduler(store)
	err := sched.Start()
	assert.Nil(t, err)

	jobRuns := []models.JobRun{}
	Eventually(func() []models.JobRun {
		_ = store.Where("JobID", j.ID, &jobRuns)
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
	defer store.Close()

	sched := services.NewScheduler(store)
	_ = sched.Start()
	defer sched.Stop()

	j := models.NewJob()
	j.Schedule = models.Schedule{Cron: "* * * * *"}
	_ = store.Save(&j)
	sched.AddJob(j)

	jobRuns := []models.JobRun{}
	Eventually(func() []models.JobRun {
		_ = store.Where("JobID", j.ID, &jobRuns)
		return jobRuns
	}).Should(HaveLen(1))
}

func TestAddJobWhenStopped(t *testing.T) {
	t.Parallel()
	RegisterTestingT(t)
	store := cltest.NewStore()
	defer store.Close()

	sched := services.NewScheduler(store)

	j := models.NewJob()
	j.Schedule = models.Schedule{Cron: "* * * * *"}
	_ = store.Save(&j)
	sched.AddJob(j)

	jobRuns := []models.JobRun{}
	Consistently(func() []models.JobRun {
		_ = store.Where("JobID", j.ID, &jobRuns)
		return jobRuns
	}).Should(HaveLen(0))

	sched.Start()
	Eventually(func() []models.JobRun {
		_ = store.Where("JobID", j.ID, &jobRuns)
		return jobRuns
	}).Should(HaveLen(1))
}
