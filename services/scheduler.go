package services

import (
	"errors"
	"fmt"
	"time"

	cronlib "github.com/mrwonko/cron"
	"github.com/smartcontractkit/chainlink-go/logger"
	"github.com/smartcontractkit/chainlink-go/store"
	"github.com/smartcontractkit/chainlink-go/store/models"
)

type Scheduler struct {
	Recurring *Recurring
	OneTime   *OneTime
	store     *store.Store
	started   bool
}

func NewScheduler(store *store.Store) *Scheduler {
	return &Scheduler{
		Recurring: &Recurring{store: store},
		OneTime: &OneTime{
			Store: store,
			Clock: &Clock{},
		},
		store: store,
	}
}

func (self *Scheduler) Start() error {
	if self.started {
		return errors.New("Scheduler already started")
	}
	if err := self.OneTime.Start(); err != nil {
		return err
	}
	if err := self.Recurring.Start(); err != nil {
		return err
	}
	self.started = true

	jobs, err := self.store.Jobs()
	if err != nil {
		return fmt.Errorf("Scheduler: %v", err)
	}

	for _, j := range jobs {
		self.AddJob(j)
	}

	return nil
}

func (self *Scheduler) Stop() {
	if self.started {
		self.OneTime.Stop()
		self.Recurring.Stop()
		self.started = false
	}
}

func (self *Scheduler) AddJob(job models.Job) {
	if !self.started {
		return
	}
	self.Recurring.AddJob(job)
	self.OneTime.AddJob(job)
}

type Recurring struct {
	cron  *cronlib.Cron
	store *store.Store
}

func (self *Recurring) Start() error {
	self.cron = cronlib.New()
	self.addResumer()
	self.cron.Start()
	return nil
}

func (self *Recurring) Stop() {
	self.cron.Stop()
	self.cron.Wait()
}

func (self *Recurring) AddJob(job models.Job) {
	for _, initr := range job.InitiatorsFor("cron") {
		cronStr := string(initr.Schedule)
		self.cron.AddFunc(cronStr, func() {
			_, err := StartJob(job.NewRun(), self.store)
			if err != nil {
				logger.Panic(err.Error())
			}
		})
	}
}

func (self *Recurring) addResumer() {
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

type Afterer interface {
	After(d time.Duration) <-chan time.Time
}

type Clock struct{}

func (self *Clock) After(d time.Duration) <-chan time.Time {
	return time.After(d)
}

type OneTime struct {
	Store *store.Store
	Clock Afterer
	done  chan struct{}
}

func (self *OneTime) Start() error {
	self.done = make(chan struct{})
	return nil
}

func (self *OneTime) AddJob(job models.Job) {
	for _, initr := range job.InitiatorsFor("runAt") {
		go self.RunJobAt(initr.Time, job)
	}
}

func (self *OneTime) Stop() {
	close(self.done)
}

func (self *OneTime) RunJobAt(t models.Time, job models.Job) {
	select {
	case <-self.done:
	case <-self.Clock.After(t.DurationFromNow()):
		_, err := StartJob(job.NewRun(), self.Store)
		if err != nil {
			logger.Panic(err.Error())
		}
	}
}
