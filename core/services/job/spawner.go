package job

import (
	"context"
	"strconv"
	"sync"
	"time"

	"github.com/pkg/errors"
	"gopkg.in/guregu/null.v4"

	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/postgres"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"
)

//go:generate mockery --name Spawner --output ./mocks/ --case=underscore
//go:generate mockery --name Delegate --output ./mocks/ --case=underscore

type (
	// The job spawner manages the spinning up and spinning down of the long-running
	// services that perform the work described by job specs.  Each active job spec
	// has 1 or more of these services associated with it.
	//
	// At present, Flux Monitor and Offchain Reporting jobs can only have a single
	// "initiator", meaning that they only require a single service.  But the older
	// "direct request" model allows for multiple initiators, which imply multiple
	// services.
	Spawner interface {
		Start()
		Stop()
		CreateJob(ctx context.Context, spec Spec, name null.String) (int32, error)
		DeleteJob(ctx context.Context, jobID int32) error
		RegisterDelegate(delegate Delegate)
	}

	spawner struct {
		orm                          ORM
		config                       Config
		jobTypeDelegates             map[Type]Delegate
		jobTypeDelegatesMu           sync.RWMutex
		startUnclaimedServicesWorker utils.SleeperTask
		services                     map[int32][]Service
		chStopJob                    chan int32

		utils.StartStopOnce
		chStop chan struct{}
		chDone chan struct{}
	}

	// TODO(spook): I can't wait for Go generics
	Delegate interface {
		JobType() Type
		ToDBRow(spec Spec) models.JobSpecV2
		FromDBRow(spec models.JobSpecV2) Spec
		ServicesForSpec(spec Spec) ([]Service, error)
	}
)

const checkForDeletedJobsPollInterval = 5 * time.Minute

var _ Spawner = (*spawner)(nil)

func NewSpawner(orm ORM, config Config) *spawner {
	s := &spawner{
		orm:              orm,
		config:           config,
		jobTypeDelegates: make(map[Type]Delegate),
		services:         make(map[int32][]Service),
		chStopJob:        make(chan int32),
		chStop:           make(chan struct{}),
		chDone:           make(chan struct{}),
	}
	s.startUnclaimedServicesWorker = utils.NewSleeperTask(
		utils.SleeperTaskFuncWorker(s.startUnclaimedServices),
	)
	return s
}

func (js *spawner) Start() {
	if !js.OkayToStart() {
		logger.Error("Job spawner has already been started")
		return
	}
	go js.runLoop()
}

func (js *spawner) Stop() {
	if !js.OkayToStop() {
		logger.Error("Job spawner has already been stopped")
		return
	}

	close(js.chStop)
	<-js.chDone
}

func (js *spawner) destroy() {
	js.stopAllServices()

	err := js.startUnclaimedServicesWorker.Stop()
	if err != nil {
		logger.Error(err)
	}
}

func (js *spawner) RegisterDelegate(delegate Delegate) {
	js.jobTypeDelegatesMu.Lock()
	defer js.jobTypeDelegatesMu.Unlock()

	if _, exists := js.jobTypeDelegates[delegate.JobType()]; exists {
		panic("registered job type " + string(delegate.JobType()) + " more than once")
	}
	logger.Infof("Registered job type '%v'", delegate.JobType())
	js.jobTypeDelegates[delegate.JobType()] = delegate
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

	specDBRows, err := js.orm.ClaimUnclaimedJobs(ctx)
	if err != nil {
		logger.Errorf("Couldn't fetch unclaimed jobs: %v", err)
		return
	}

	js.jobTypeDelegatesMu.RLock()
	defer js.jobTypeDelegatesMu.RUnlock()

	for _, specDBRow := range specDBRows {
		if _, exists := js.services[specDBRow.ID]; exists {
			logger.Warnw("Job spawner ORM attempted to claim locally-claimed job, skipping", "jobID", specDBRow.ID)
			continue
		}

		var services []Service
		for _, delegate := range js.jobTypeDelegates {
			spec := delegate.FromDBRow(specDBRow)
			if spec == nil {
				// This spec isn't owned by this delegate
				continue
			}

			moreServices, err := delegate.ServicesForSpec(spec)
			if err != nil {
				logger.Errorw("Error creating services for job", "jobID", specDBRow.ID, "error", err)
				continue
			}
			services = append(services, moreServices...)
		}

		logger.Infow("Starting services for job", "jobID", specDBRow.ID, "count", len(services))

		for _, service := range services {
			err := service.Start()
			if err != nil {
				logger.Errorw("Error creating service for job", "jobID", specDBRow.ID, "error", err)
				continue
			}
			js.services[specDBRow.ID] = append(js.services[specDBRow.ID], service)
		}
	}
}

func (js *spawner) stopAllServices() {
	for jobID := range js.services {
		js.stopService(jobID)
	}
}

func (js *spawner) stopService(jobID int32) {
	for _, service := range js.services[jobID] {
		err := service.Close()
		if err != nil {
			logger.Errorw("Error stopping job service", "jobID", jobID, "error", err)
		} else {
			logger.Infow("Stopped job service", "jobID", jobID)
		}
	}
	delete(js.services, jobID)
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

func (js *spawner) CreateJob(ctx context.Context, spec Spec, name null.String) (int32, error) {
	js.jobTypeDelegatesMu.Lock()
	defer js.jobTypeDelegatesMu.Unlock()

	delegate, exists := js.jobTypeDelegates[spec.JobType()]
	if !exists {
		logger.Errorf("job type '%s' has not been registered with the job.Spawner", spec.JobType())
		return 0, errors.Errorf("job type '%s' has not been registered with the job.Spawner", spec.JobType())
	}

	ctx, cancel := utils.CombinedContext(js.chStop, ctx)
	defer cancel()

	specDBRow := delegate.ToDBRow(spec)
	specDBRow.Name = name
	err := js.orm.CreateJob(ctx, &specDBRow, spec.TaskDAG())
	if err != nil {
		logger.Errorw("Error creating job", "type", spec.JobType(), "error", err)
		return 0, err
	}

	logger.Infow("Created job", "type", spec.JobType(), "jobID", specDBRow.ID)
	return specDBRow.ID, err
}

func (js *spawner) DeleteJob(ctx context.Context, jobID int32) error {
	if jobID == 0 {
		return errors.New("will not delete job with 0 ID")
	}

	ctx, cancel := utils.CombinedContext(js.chStop, ctx)
	defer cancel()

	err := js.orm.DeleteJob(ctx, jobID)
	if err != nil {
		logger.Errorw("Error deleting job", "jobID", jobID, "error", err)
		return err
	}
	logger.Infow("Deleted job", "jobID", jobID)

	select {
	case <-js.chStop:
	case js.chStopJob <- jobID:
	}

	return nil
}
