package services

import (
	"errors"
	"sync"

	"github.com/mrwonko/cron"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/store"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"
)

// Scheduler contains fields for Recurring and OneTime for occurrences,
// a pointer to the store and a started field to indicate if the Scheduler
// has started or not.
type Scheduler struct {
	Recurring    *Recurring
	OneTime      *OneTime
	store        *store.Store
	runManager   RunManager
	startedMutex sync.RWMutex
	started      bool
}

// NewScheduler initializes the Scheduler instances with both Recurring
// and OneTime fields since jobs can contain tasks which utilize both.
func NewScheduler(store *store.Store, runManager RunManager) *Scheduler {
	return &Scheduler{
		Recurring: NewRecurring(runManager),
		OneTime: &OneTime{
			Store:      store,
			Clock:      store.Clock,
			RunManager: runManager,
		},
		store:      store,
		runManager: runManager,
	}
}

// Start checks to ensure the Scheduler has not already started,
// calls the Start function for both Recurring and OneTime types,
// sets the started field to true, and adds jobs relevant to its
// initiator ("cron" and "runat").
func (s *Scheduler) Start() error {
	s.startedMutex.Lock()
	defer s.startedMutex.Unlock()
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

	return s.store.Jobs(func(j *models.JobSpec) bool {
		s.addJob(j)
		return true
	})
}

// Stop is the governing function for both Recurring and OneTime
// Stop function. Sets the started field to false.
func (s *Scheduler) Stop() {
	s.startedMutex.Lock()
	defer s.startedMutex.Unlock()
	if s.started {
		s.Recurring.Stop()
		s.OneTime.Stop()
		s.started = false
	}
}

func (s *Scheduler) addJob(job *models.JobSpec) {
	s.Recurring.AddJob(*job)
	s.OneTime.AddJob(*job)
}

// AddJob is the governing function for Recurring and OneTime,
// and will only execute if the Scheduler has not already started.
func (s *Scheduler) AddJob(job models.JobSpec) {
	s.startedMutex.RLock()
	defer s.startedMutex.RUnlock()
	if !s.started {
		return
	}
	s.addJob(&job)
}

// Recurring is used for runs that need to execute on a schedule,
// and is configured with cron.
// Instances of Recurring must be initialized using NewRecurring().
type Recurring struct {
	Cron       Cron
	Clock      utils.Nower
	runManager RunManager
}

// NewRecurring create a new instance of Recurring, ready to use.
func NewRecurring(runManager RunManager) *Recurring {
	return &Recurring{
		runManager: runManager,
	}
}

// Start for Recurring types executes tasks with a "cron" initiator
// based on the configured schedule for the run.
func (r *Recurring) Start() error {
	r.Cron = newChainlinkCron()
	r.Cron.Start()
	return nil
}

// Stop stops the cron scheduler and waits for running jobs to finish.
func (r *Recurring) Stop() {
	r.Cron.Stop()
}

// AddJob looks for "cron" initiators, adds them to cron's schedule
// for execution when specified.
func (r *Recurring) AddJob(job models.JobSpec) {
	for _, initr := range job.InitiatorsFor(models.InitiatorCron) {
		r.Cron.AddFunc(string(initr.Schedule), func() {
			_, err := r.runManager.Create(job.ID, &initr, &models.JSON{}, nil, &models.RunRequest{})
			if err != nil && !expectedRecurringScheduleJobError(err) {
				logger.Errorw(err.Error())
			}
		})
	}
}

// OneTime represents runs that are to be executed only once.
type OneTime struct {
	Store      *store.Store
	Clock      utils.Afterer
	RunManager RunManager
	done       chan struct{}
}

// Start allocates a channel for the "done" field with an empty struct.
func (ot *OneTime) Start() error {
	ot.done = make(chan struct{})
	return nil
}

// AddJob runs the job at the time specified for the "runat" initiator.
func (ot *OneTime) AddJob(job models.JobSpec) {
	for _, initiator := range job.InitiatorsFor(models.InitiatorRunAt) {
		if !initiator.Time.Valid {
			logger.Errorf("RunJobAt: JobSpec %s must have initiator with valid run at time: %v", job.ID, initiator)
			continue
		}

		go ot.RunJobAt(initiator, job)
	}
}

// Stop closes the "done" field's channel.
func (ot *OneTime) Stop() {
	close(ot.done)
}

// RunJobAt wait until the Stop() function has been called on the run
// or the specified time for the run is after the present time.
func (ot *OneTime) RunJobAt(initiator models.Initiator, job models.JobSpec) {
	select {
	case <-ot.done:
	case <-ot.Clock.After(utils.DurationFromNow(initiator.Time.Time)):
		_, err := ot.RunManager.Create(job.ID, &initiator, &models.JSON{}, nil, &models.RunRequest{})
		if err != nil {
			logger.Error(err.Error())
			return
		}
		if err := ot.Store.MarkRan(&initiator, true); err != nil {
			logger.Error(err.Error())
		}
	}
}

func expectedRecurringScheduleJobError(err error) bool {
	switch err.(type) {
	case RecurringScheduleJobError:
		return true
	default:
		return false
	}
}

// Cron is an interface for scheduling recurring functions to run.
// Cron's schedule format is similar to the standard cron format
// but with an extra field at the beginning for seconds.
type Cron interface {
	Start()
	Stop()
	AddFunc(string, func()) error
}

type chainlinkCron struct {
	*cron.Cron
}

func newChainlinkCron() *chainlinkCron {
	return &chainlinkCron{cron.New()}
}

func (cc *chainlinkCron) Stop() {
	cc.Cron.Stop()
	cc.Cron.Wait()
}
