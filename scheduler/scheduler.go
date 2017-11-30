package scheduler

import (
	"fmt"
	cronlib "github.com/robfig/cron"
	"github.com/smartcontractkit/chainlink-go/models"
)

type Scheduler struct {
	cron *cronlib.Cron
}

func Start() (*Scheduler, error) {
	sched := New()
	err := sched.Start()
	return sched, err
}

func New() *Scheduler {
	return &Scheduler{cronlib.New()}
}

func (self *Scheduler) Start() error {
	jobs, err := models.JobsWithCron()
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
}

func (self *Scheduler) AddJob(job models.Job) {
	cronStr := string(job.Schedule.Cron)
	self.cron.AddFunc(cronStr, func() { job.Run() })
}
