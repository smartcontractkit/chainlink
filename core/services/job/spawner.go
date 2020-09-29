package job

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/smartcontractkit/chainlink/core/logger"
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
		CreateJob(spec Spec) (int32, error)
		DeleteJob(ctx context.Context, spec Spec) error
		RegisterDelegate(delegate Delegate)
	}

	spawner struct {
		orm                    ORM
		jobTypeDelegates       map[Type]Delegate
		jobTypeDelegatesMu     sync.RWMutex
		startUnclaimedServices utils.SleeperTask
		services               map[int32][]Service
		chStopJob              chan int32

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

var _ Spawner = (*spawner)(nil)

func NewSpawner(orm ORM) *spawner {
	s := &spawner{
		orm:              orm,
		jobTypeDelegates: make(map[Type]Delegate),
		services:         make(map[int32][]Service),
		chStopJob:        make(chan int32),
		chStop:           make(chan struct{}),
		chDone:           make(chan struct{}),
	}
	s.startUnclaimedServices = utils.NewSleeperTask(
		utils.SleeperTaskFuncWorker(s.startUnclaimedServicesWorker),
	)
	return s
}

func (js *spawner) Start() {
	js.AssertNeverStarted()
	go js.runLoop()
}

func (js *spawner) Stop() {
	js.AssertNeverStopped()
	close(js.chStop)
	<-js.chDone
}

func (js *spawner) RegisterDelegate(delegate Delegate) {
	js.jobTypeDelegatesMu.Lock()
	defer js.jobTypeDelegatesMu.Unlock()

	if _, exists := js.jobTypeDelegates[delegate.JobType()]; exists {
		panic("registered job type " + string(delegate.JobType()) + " more than once")
	}
	logger.Infof("Registered new job type '%v'", delegate.JobType())
	js.jobTypeDelegates[delegate.JobType()] = delegate
}

func (js *spawner) runLoop() {
	defer close(js.chDone)

	// Initialize the Postgres event listener for new jobs
	var chNewJobs <-chan string
	listener, err := js.orm.ListenForNewJobs()
	if err != nil {
		logger.Errorw("Job spawner failed to subscribe to 'new job' events, falling back to polling", "error", err)
	} else {
		chNewJobs = listener.Events()
	}

	// Initialize the DB poll ticker
	dbPollTicker := time.NewTicker(1 * time.Second)
	defer dbPollTicker.Stop()

	js.startUnclaimedServices.WakeUp()
	for {
		select {
		case <-chNewJobs:
			js.startUnclaimedServices.WakeUp()

		case <-dbPollTicker.C:
			js.startUnclaimedServices.WakeUp()

		case jobID := <-js.chStopJob:
			js.stopService(jobID)

		case <-js.chStop:
			if listener != nil {
				err := listener.Stop()
				if err != nil {
					logger.Errorw(`Error stopping pipeline runner's "new runs" listener`, "error", err)
				}
			}
			js.startUnclaimedServices.Stop()
			js.stopAllServices()
			return

		}
	}
}

func (js *spawner) startUnclaimedServicesWorker() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	specDBRows, err := js.orm.UnclaimedJobs(ctx)
	if err != nil {
		logger.Errorf("Couldn't fetch unclaimed jobs: %v", err)
		return
	}

	js.jobTypeDelegatesMu.RLock()
	defer js.jobTypeDelegatesMu.RUnlock()

	for _, specDBRow := range specDBRows {
		// `UnclaimedJobs` guarantees that we won't try to start jobs that other
		// nodes in the cluster are already running, but because Postgres
		// advisory locks are session-scoped, we have to manually guard against
		// trying to start jobs that we're already running locally.
		if _, exists := js.services[specDBRow.ID]; exists {
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
		err := service.Stop()
		if err != nil {
			logger.Errorw("Error stopping job service", "jobID", jobID, "error", err)
		} else {
			logger.Infow("Stopped job service", "jobID", jobID)
		}
	}
	delete(js.services, jobID)
}

func (js *spawner) CreateJob(spec Spec) (int32, error) {
	js.jobTypeDelegatesMu.Lock()
	defer js.jobTypeDelegatesMu.Unlock()

	delegate, exists := js.jobTypeDelegates[spec.JobType()]
	if !exists {
		panic(fmt.Sprintf("job type '%s' has not been registered with the job.Spawner", spec.JobType()))
	}

	specDBRow := delegate.ToDBRow(spec)
	err := js.orm.CreateJob(&specDBRow, spec.TaskDAG())
	if err != nil {
		logger.Errorw("Error creating job", "type", spec.JobType(), "error", err)
		return 0, err
	}

	logger.Infow("Created job", "type", spec.JobType(), "jobID", specDBRow.ID)
	return specDBRow.ID, err
}

func (js *spawner) DeleteJob(ctx context.Context, spec Spec) error {
	err := js.orm.DeleteJob(ctx, spec.JobID())
	if err != nil {
		logger.Errorw("Error deleting job", "type", spec.JobType(), "jobID", spec.JobID(), "error", err)
		return err
	}
	logger.Infow("Deleted job", "type", spec.JobType(), "jobID", spec.JobID())

	select {
	case <-js.chStop:
	case js.chStopJob <- spec.JobID():
	}

	return nil
}
