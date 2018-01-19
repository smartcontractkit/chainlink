package services

import (
	"errors"
	"fmt"
	"time"

	cronlib "github.com/mrwonko/cron"
	"github.com/smartcontractkit/chainlink/logger"
	"github.com/smartcontractkit/chainlink/store"
	"github.com/smartcontractkit/chainlink/store/models"
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
			Clock: store.Clock,
		},
		store: store,
	}
}

func (s *Scheduler) Start() error {
	if s.started {
		return errors.New("Scheduler already started")
	}
	if err := s.OneTime.Start(); err != nil {
		return err
	}
	if err := s.Recurring.Start(); err != nil {
		return err
	}
	s.started = true

	jobs, err := s.store.Jobs()
	if err != nil {
		return fmt.Errorf("Scheduler: %v", err)
	}

	for _, j := range jobs {
		s.AddJob(j)
	}

	return nil
}

func (s *Scheduler) Stop() {
	if s.started {
		s.Recurring.Stop()
		s.OneTime.Stop()
		s.started = false
	}
}

func (s *Scheduler) AddJob(job *models.Job) {
	if !s.started {
		return
	}
	s.Recurring.AddJob(job)
	s.OneTime.AddJob(job)
}

type Recurring struct {
	cron  *cronlib.Cron
	store *store.Store
}

func (r *Recurring) Start() error {
	r.cron = cronlib.New()
	r.addResumer()
	r.cron.Start()
	return nil
}

func (r *Recurring) Stop() {
	r.cron.Stop()
	r.cron.Wait()
}

func (r *Recurring) AddJob(job *models.Job) {
	for _, initr := range job.InitiatorsFor(models.InitiatorCron) {
		cronStr := string(initr.Schedule)
		r.cron.AddFunc(cronStr, func() {
			_, err := BeginRun(job, r.store)
			if err != nil && !expectedRecurringError(err) {
				logger.Panic(err.Error())
			}
		})
	}
}

func (r *Recurring) addResumer() {
	r.cron.AddFunc(r.store.Config.PollingSchedule, func() {
		pendingRuns, err := r.store.PendingJobRuns()
		if err != nil {
			logger.Panic(err.Error())
		}
		for _, jobRun := range pendingRuns {
			if err := ExecuteRun(jobRun, r.store); err != nil {
				logger.Panic(err.Error())
			}
		}
	})
}

type Afterer interface {
	After(d time.Duration) <-chan time.Time
}

type OneTime struct {
	Store *store.Store
	Clock Afterer
	done  chan struct{}
}

func (ot *OneTime) Start() error {
	ot.done = make(chan struct{})
	return nil
}

func (ot *OneTime) AddJob(job *models.Job) {
	for _, initr := range job.InitiatorsFor(models.InitiatorRunAt) {
		go ot.RunJobAt(initr.Time, job)
	}
}

func (ot *OneTime) Stop() {
	close(ot.done)
}

func (ot *OneTime) RunJobAt(t models.Time, job *models.Job) {
	select {
	case <-ot.done:
	case <-ot.Clock.After(t.DurationFromNow()):
		_, err := BeginRun(job, ot.Store)
		if err != nil {
			logger.Panic(err.Error())
		}
	}
}

func expectedRecurringError(err error) bool {
	switch err.(type) {
	case JobRunnerError:
		return true
	default:
		return false
	}
}
