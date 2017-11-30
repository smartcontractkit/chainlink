package scheduler_test

import (
	. "github.com/onsi/gomega"
	"github.com/smartcontractkit/chainlink-go/internal/cltest"
	"github.com/smartcontractkit/chainlink-go/models"
	"github.com/smartcontractkit/chainlink-go/scheduler"
	"testing"
)

func TestLoadingSavedSchedules(t *testing.T) {
	RegisterTestingT(t)
	cltest.SetUpDB()
	defer cltest.TearDownDB()

	j := models.NewJob()
	j.Schedule = models.Schedule{Cron: "* * * * *"}
	_ = models.Save(&j)

	sched := scheduler.New()
	_ = sched.Start()
	defer sched.Stop()

	jobRuns := []models.JobRun{}
	Eventually(func() []models.JobRun {
		_ = models.Where("JobID", j.ID, &jobRuns)
		return jobRuns
	}).Should(HaveLen(1))
}
