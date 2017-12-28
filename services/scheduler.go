package services

import (
	"errors"
	"fmt"

	cronlib "github.com/mrwonko/cron"
	"github.com/smartcontractkit/chainlink-go/logger"
	"github.com/smartcontractkit/chainlink-go/store"
	"github.com/smartcontractkit/chainlink-go/store/models"
)

type Scheduler struct {
	cron    *cronlib.Cron
	store   *store.Store
	started bool
}

func NewScheduler(store *store.Store) *Scheduler {
	return &Scheduler{store: store}
}

func (self *Scheduler) Start() error {
	if self.started {
		return errors.New("Scheduler already started")
	}
	self.cron = cronlib.New()
	jobs, err := self.store.JobsWithCron()
	if err != nil {
		return fmt.Errorf("Scheduler: %v", err)
	}

	self.started = true
	for _, j := range jobs {
		self.AddJob(j)
	}

	self.addResumer()
	self.cron.Start()
	return nil
}

func (self *Scheduler) Stop() {
	if self.started {
		self.cron.Stop()
		self.cron.Wait()
		self.started = false
	}
}

func (self *Scheduler) AddJob(job models.Job) {
	if !self.started {
		return
	}
	for _, initr := range job.Schedules() {
		cronStr := string(initr.Schedule)
		self.cron.AddFunc(cronStr, func() {
			_, err := StartJob(job.NewRun(), self.store)
			if err != nil {
				logger.Panic(err.Error())
			}
		})
	}
}

func (self *Scheduler) addResumer() {
	self.cron.AddFunc(self.store.Config.PollingSchedule, func() {
		pendingRuns, err := self.store.PendingJobRuns()
		if err != nil {
			logger.Panic(err.Error())
		}
		for _, jobRun := range pendingRuns {
			_, err := StartJob(jobRun, self.store)
			if err != nil {
				logger.Panic(err.Error())
			}
		}
	})
}
