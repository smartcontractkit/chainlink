package scheduler_test

import (
	. "github.com/onsi/gomega"
	"github.com/smartcontractkit/chainlink-go/internal/cltest"
	"github.com/smartcontractkit/chainlink-go/models"
	"github.com/smartcontractkit/chainlink-go/scheduler"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLoadingSavedSchedules(t *testing.T) {
	RegisterTestingT(t)
	cltest.SetUpDB()
	defer cltest.TearDownDB()

	j := models.NewJob()
	j.Schedule = models.Schedule{Cron: "* * * * *"}
	jobWoCron := models.NewJob()
	_ = models.Save(&j)
	_ = models.Save(&jobWoCron)

	sched := scheduler.New()
	err := sched.Start()
	assert.Nil(t, err)
	defer sched.Stop()

	jobRuns := []models.JobRun{}
	Eventually(func() []models.JobRun {
		_ = models.Where("JobID", j.ID, &jobRuns)
		return jobRuns
	}).Should(HaveLen(1))

	err = models.Where("JobID", jobWoCron.ID, &jobRuns)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(jobRuns), "No jobs should be created without the scheduler")
}
