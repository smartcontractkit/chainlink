package job

import (
	"context"
	"reflect"
	"strconv"
	"sync"
	"time"

	"github.com/pkg/errors"
	"gopkg.in/guregu/null.v4"

	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/service"
	"github.com/smartcontractkit/chainlink/core/services/postgres"
	"github.com/smartcontractkit/chainlink/core/utils"
)

//go:generate mockery --name Spawner --output ./mocks/ --case=underscore
//go:generate mockery --name Delegate --output ./mocks/ --case=underscore

type (
	// The job spawner manages the spinning up and spinning down of the long-running
	// services that perform the work described by job specs.  Each active job spec
	// has 1 or more of these services associated with it.
	Spawner interface {
		service.Service
		CreateJob(ctx context.Context, spec Job, name null.String) (Job, error)
		DeleteJob(ctx context.Context, jobID int32) error
		ActiveJobs() map[int32]Job
	}

	spawner struct {
		orm                          ORM
		config                       Config
		jobTypeDelegates             map[Type]Delegate
		startUnclaimedServicesWorker utils.SleeperTask
		activeJobs                   map[int32]activeJob
		activeJobsMu                 sync.RWMutex
		chStopJob                    chan int32
		txm                          postgres.TransactionManager

		utils.StartStopOnce
		chStop chan struct{}
		chDone chan struct{}
	}

	// TODO(spook): I can't wait for Go generics
	Delegate interface {
		JobType() Type
		// ServicesForSpec returns services to be started and stopped for this
		// job. In case a given job type relies upon well-defined startup/shutdown
		// ordering for services, they are started in the order they are given
		// and stopped in reverse order.
		ServicesForSpec(spec Job) ([]Service, error)
		AfterJobCreated(spec Job)
		BeforeJobDeleted(spec Job)
	}

	activeJob struct {
		delegate Delegate
		spec     Job
		services []Service
	}
)

const checkForDeletedJobsPollInterval = 5 * time.Minute

var _ Spawner = (*spawner)(nil)

func NewSpawner(orm ORM, config Config, jobTypeDelegates map[Type]Delegate, txm postgres.TransactionManager) *spawner {
	s := &spawner{
		orm:              orm,
		config:           config,
		jobTypeDelegates: jobTypeDelegates,
		txm:              txm,
		activeJobs:       make(map[int32]activeJob),
		chStopJob:        make(chan int32),
		chStop:           make(chan struct{}),
		chDone:           make(chan struct{}),
	}
	s.startUnclaimedServicesWorker = utils.NewSleeperTask(
		utils.SleeperTaskFuncWorker(s.startUnclaimedServices),
	)
	return s
}

func (js *spawner) Start() error {
	return js.StartOnce("JobSpawner", func() error {
		go js.runLoop()
		return nil

	})
}

func (js *spawner) Close() error {
	return js.StopOnce("JobSpawner", func() error {
		close(js.chStop)
		<-js.chDone
		return nil

	})
}

func (js *spawner) destroy() {
	js.stopAllServices()

	err := js.startUnclaimedServicesWorker.Stop()
	if err != nil {
		logger.Error(err)
	}
}

func (js *spawner) runLoop() {
	defer close(js.chDone)
	defer js.destroy()

	// Initialize the Postgres event listener for created and deleted jobs
	var newJobEvents <-chan postgres.Event
	newJobs, err := js.orm.ListenForNewJobs()
	if err != nil {
		logger.Warn("Job spawner could not subscribe to new job events, falling back to polling")
	} else {
		defer newJobs.Close()
		newJobEvents = newJobs.Events()
	}
	var pgDeletedJobEvents <-chan postgres.Event
	deletedJobs, err := js.orm.ListenForDeletedJobs()
	if err != nil {
		logger.Warn("Job spawner could not subscribe to deleted job events")
	} else {
		defer deletedJobs.Close()
		pgDeletedJobEvents = deletedJobs.Events()
	}

	// Initialize the DB poll ticker
	dbPollTicker := time.NewTicker(utils.WithJitter(js.config.TriggerFallbackDBPollInterval()))
	defer dbPollTicker.Stop()

	// Initialize the poll that checks for deleted jobs and removes them
	// This is only necessary as a fallback in case the event doesn't fire for some reason
	// It doesn't need to run very often
	deletedPollTicker := time.NewTicker(checkForDeletedJobsPollInterval)
	defer deletedPollTicker.Stop()

	ctx, cancel := utils.CombinedContext(js.chStop)
	defer cancel()

	js.startUnclaimedServicesWorker.WakeUp()
	for {
		select {
		case <-newJobEvents:
			js.startUnclaimedServicesWorker.WakeUp()

		case <-dbPollTicker.C:
			js.startUnclaimedServicesWorker.WakeUp()

		case jobID := <-js.chStopJob:
			js.stopService(jobID)

		case <-deletedPollTicker.C:
			js.checkForDeletedJobs(ctx)

		case deleteJobEvent := <-pgDeletedJobEvents:
			js.handlePGDeleteEvent(ctx, deleteJobEvent)

		case <-js.chStop:
			return
		}
	}
}

func (js *spawner) startUnclaimedServices() {
	ctx, cancel := utils.CombinedContext(js.chStop, 5*time.Second)
	defer cancel()

	specs, err := js.orm.ClaimUnclaimedJobs(ctx)
	if err != nil {
		logger.Errorf("Couldn't fetch unclaimed jobs: %v", err)
		return
	}

	js.activeJobsMu.Lock()
	defer js.activeJobsMu.Unlock()

	for _, spec := range specs {
		if _, exists := js.activeJobs[spec.ID]; exists {
			logger.Warnw("Job spawner ORM attempted to claim locally-claimed job, skipping", "jobID", spec.ID)
			continue
		}

		delegate, exists := js.jobTypeDelegates[spec.Type]
		if !exists {
			logger.Errorw("Job type has not been registered with job.Spawner", "type", spec.Type, "jobID", spec.ID)
			continue
		}
		services, err := delegate.ServicesForSpec(spec)
		if err != nil {
			logger.Errorw("Error creating services for job", "jobID", spec.ID, "error", err)
			js.orm.RecordError(ctx, spec.ID, err.Error())
			continue
		}

		logger.Debugw("JobSpawner: Starting services for job", "jobID", spec.ID, "count", len(services))

		aj := activeJob{delegate: delegate, spec: spec}
		for _, service := range services {
			err := service.Start()
			if err != nil {
				logger.Errorw("Error creating service for job", "jobID", spec.ID, "error", err)
				continue
			}
			aj.services = append(aj.services, service)
		}
		js.activeJobs[spec.ID] = aj
	}

	logger.Infow("JobSpawner: all jobs running", "count", len(specs))
}

func (js *spawner) stopAllServices() {
	var jobIDs []int32
	func() {
		js.activeJobsMu.RLock()
		defer js.activeJobsMu.RUnlock()

		for jobID := range js.activeJobs {
			jobIDs = append(jobIDs, jobID)
		}
	}()

	for _, jobID := range jobIDs {
		js.stopService(jobID)
	}
}

func (js *spawner) stopService(jobID int32) {
	js.activeJobsMu.Lock()
	defer js.activeJobsMu.Unlock()

	aj := js.activeJobs[jobID]

	for i := len(aj.services) - 1; i >= 0; i-- {
		service := aj.services[i]
		err := service.Close()
		if err != nil {
			logger.Errorw("Error stopping job service", "jobID", jobID, "error", err, "subservice", i, "serviceType", reflect.TypeOf(service))
		} else {
			logger.Infow("Stopped job service", "jobID", jobID, "subservice", i, "serviceType", reflect.TypeOf(service))
		}
	}
	delete(js.activeJobs, jobID)
}

func (js *spawner) checkForDeletedJobs(ctx context.Context) {
	jobIDs, err := js.orm.CheckForDeletedJobs(ctx)
	if err != nil {
		logger.Errorw("failed to CheckForDeletedJobs", "err", err)
		return
	}
	for _, jobID := range jobIDs {
		js.unloadDeletedJob(ctx, jobID)
	}
}

func (js *spawner) unloadDeletedJob(ctx context.Context, jobID int32) {
	logger.Infow("Unloading deleted job", "jobID", jobID)

	js.stopService(jobID)

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := js.orm.UnclaimJob(ctx, jobID); err != nil {
		logger.Errorw("Unexpected error unclaiming job", "jobID", jobID)
	}
}

func (js *spawner) handlePGDeleteEvent(ctx context.Context, ev postgres.Event) {
	jobIDString := ev.Payload
	jobID64, err := strconv.ParseInt(jobIDString, 10, 32)
	if err != nil {
		logger.Errorw("Unexpected error decoding deleted job event payload, expected 32-bit integer", "payload", jobIDString, "channel", ev.Channel)
	}
	jobID := int32(jobID64)
	js.unloadDeletedJob(ctx, jobID)
}

func (js *spawner) CreateJob(ctx context.Context, spec Job, name null.String) (Job, error) {
	var jb Job
	var err error
	delegate, exists := js.jobTypeDelegates[spec.Type]
	if !exists {
		logger.Errorf("job type '%s' has not been registered with the job.Spawner", spec.Type)
		return jb, errors.Errorf("job type '%s' has not been registered with the job.Spawner", spec.Type)
	}

	ctx, cancel := utils.CombinedContext(js.chStop, ctx)
	defer cancel()

	spec.Name = name

	ctx, cancel = context.WithTimeout(ctx, postgres.DefaultQueryTimeout)
	defer cancel()
	err = js.txm.TransactWithContext(ctx, func(context.Context) error {
		jb, err = js.orm.CreateJob(ctx, &spec, spec.Pipeline)
		if err != nil {
			logger.Errorw("Error creating job", "type", spec.Type, "error", err)

			return err
		}

		return nil
	})
	if err != nil {
		return jb, err
	}

	delegate.AfterJobCreated(jb)

	logger.Infow("Created job", "type", jb.Type, "jobID", jb.ID)
	return jb, err
}

func (js *spawner) DeleteJob(ctx context.Context, jobID int32) error {
	if jobID == 0 {
		return errors.New("will not delete job with 0 ID")
	}

	var aj activeJob
	var exists bool
	func() {
		js.activeJobsMu.RLock()
		defer js.activeJobsMu.RUnlock()
		aj, exists = js.activeJobs[jobID]
	}()
	if !exists {
		return errors.Errorf("job not found (id: %v)", jobID)
	}

	// Stop the service if we own the job.
	js.stopService(jobID)

	aj.delegate.BeforeJobDeleted(aj.spec)

	ctx, cancel := utils.CombinedContext(js.chStop, ctx)
	defer cancel()
	err := js.orm.DeleteJob(ctx, jobID)
	if err != nil {
		logger.Errorw("Error deleting job", "jobID", jobID, "error", err)
		return err
	}

	logger.Infow("Deleted job", "jobID", jobID)

	return nil
}

func (js *spawner) ActiveJobs() map[int32]Job {
	js.activeJobsMu.RLock()
	defer js.activeJobsMu.RUnlock()

	m := make(map[int32]Job, len(js.activeJobs))
	for jobID := range js.activeJobs {
		m[jobID] = js.activeJobs[jobID].spec
	}
	return m
}

var _ Delegate = &NullDelegate{}

type NullDelegate struct {
	Type Type
}

func (n *NullDelegate) JobType() Type {
	return n.Type
}

func (n *NullDelegate) ServicesForSpec(spec Job) (s []Service, err error) {
	return
}

func (*NullDelegate) AfterJobCreated(spec Job)  {}
func (*NullDelegate) BeforeJobDeleted(spec Job) {}
