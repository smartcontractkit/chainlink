package services

import (
	"fmt"

	cronlib "github.com/mrwonko/cron"
	"github.com/smartcontractkit/chainlink-go/logger"
	"github.com/smartcontractkit/chainlink-go/models"
)

type Scheduler struct {
	cron *cronlib.Cron
	orm  models.ORM
}

func NewScheduler(orm models.ORM) *Scheduler {
	return &Scheduler{cronlib.New(), orm}
}

func (self *Scheduler) Start() error {
	jobs, err := self.orm.JobsWithCron()
	if err != nil {
		return fmt.Errorf("Scheduler: %v", err)
	}

	for _, j := range jobs {
		self.AddJob(j)
	}

	self.cron.Start()
	return nil
}

func (self *Scheduler) Stop() {
	self.cron.Stop()
	self.cron.Wait()
}

func (self *Scheduler) AddJob(job models.Job) {
	cronStr := string(job.Schedule.Cron)
	self.cron.AddFunc(cronStr, func() {
		err := StartJob(job.NewRun(), self.orm)
		if err != nil {
			logger.Panic(err.Error())
		}
	})
}
