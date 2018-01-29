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

// Scheduler contains fields for Recurring and OneTime for occurrences,
// a pointer to the store and a started field to indicate if the Scheduler
// has started or not.
type Scheduler struct {
	Recurring *Recurring
	OneTime   *OneTime
	store     *store.Store
	started   bool
}

// NewScheduler initializes the Scheduler instances with both Recurring
// and OneTime fields since jobs can contain tasks which utilize both.
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

// Start checks to ensure the Scheduler has not already started,
// calls the Start function for both Recurring and OneTime types,
// sets the started field to true, and adds jobs relevant to its
// initiator ("cron" and "runat").
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

// Stop is the governing function for both Recurring and OneTime
// Stop function. Sets the started field to false.
func (s *Scheduler) Stop() {
	if s.started {
		s.Recurring.Stop()
		s.OneTime.Stop()
		s.started = false
	}
}

// AddJob is the governing function for Recurring and OneTime,
// and will only execute if the Scheduler has not already started.
func (s *Scheduler) AddJob(job *models.Job) {
	if !s.started {
		return
	}
	s.Recurring.AddJob(job)
	s.OneTime.AddJob(job)
}

// Recurring is used for runs that need to execute on a schedule,
// and is configured with cron.
type Recurring struct {
	cron  *cronlib.Cron
	store *store.Store
}

// Start for Recurring types executes tasks with a "cron" initiator
// based on the configured schedule for the run.
func (r *Recurring) Start() error {
	r.cron = cronlib.New()
	r.addResumer()
	r.cron.Start()
	return nil
}

// Stop stops the cron scheduler and waits for running jobs to finish.
func (r *Recurring) Stop() {
	r.cron.Stop()
	r.cron.Wait()
}

// AddJob looks for "cron" initiators, adds them to cron's schedule
// for execution when specified.
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

// Afterer represents the time after a specified time.
type Afterer interface {
	After(d time.Duration) <-chan time.Time
}

// OneTime represents runs that are to be executed only once.
type OneTime struct {
	Store *store.Store
	Clock Afterer
	done  chan struct{}
}

// Start allocates a channel for the "done" field with an empty struct.
func (ot *OneTime) Start() error {
	ot.done = make(chan struct{})
	return nil
}

// AddJob runs the job at the time specified for the "runat" initiator.
func (ot *OneTime) AddJob(job *models.Job) {
	for _, initr := range job.InitiatorsFor(models.InitiatorRunAt) {
		go ot.RunJobAt(initr.Time, job)
	}
}

// Stop closes the "done" field's channel.
func (ot *OneTime) Stop() {
	close(ot.done)
}

// RunJobAt wait until the Stop() function has been called on the run
// or the specified time for the run is after the present time.
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
