package scheduler

import (
	"fmt"
	cronlib "github.com/robfig/cron"
	"github.com/smartcontractkit/chainlink-go/models"
)

type Scheduler struct {
	cron *cronlib.Cron
}

func New() *Scheduler {
	return &Scheduler{cronlib.New()}
}

func (self *Scheduler) Start() error {
	jobs := []models.Job{}
	err := models.All(&jobs)
	if err != nil {
		return fmt.Errorf("Scheduler: ", err)
	}

	for _, j := range jobs {
		var job = j
		cronStr := string(job.Schedule.Cron)
		self.cron.AddFunc(cronStr, func() { job.Run() })
	}

	self.cron.Start()
	return nil
}

func (self *Scheduler) Stop() {
	self.cron.Stop()
}
