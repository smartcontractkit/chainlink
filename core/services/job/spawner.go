package job

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/store/models"
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
		orm                ORM
		jobTypeDelegates   map[Type]Delegate
		jobTypeDelegatesMu sync.RWMutex
		services           map[int32][]Service
		chStopJob          chan int32
		chStop             chan struct{}
		chDone             chan struct{}
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
	return &spawner{
		orm:              orm,
		jobTypeDelegates: make(map[Type]Delegate),
		services:         make(map[int32][]Service),
		chStopJob:        make(chan int32),
		chStop:           make(chan struct{}),
		chDone:           make(chan struct{}),
	}
}

func (js *spawner) Start() {
	go js.runLoop()
}

func (js *spawner) Stop() {
	close(js.chStop)
	<-js.chDone
}

func (js *spawner) RegisterDelegate(delegate Delegate) {
	js.jobTypeDelegatesMu.Lock()
	defer js.jobTypeDelegatesMu.Unlock()

	if _, exists := js.jobTypeDelegates[delegate.JobType()]; exists {
		panic("registered job type " + string(delegate.JobType()) + " more than once")
	}
	js.jobTypeDelegates[delegate.JobType()] = delegate
}

func (js *spawner) runLoop() {
	defer close(js.chDone)

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	js.startServicesForUnclaimedJobs()

	for {
		select {
		case <-js.chStop:
			js.stopAllServices()
			return

		case jobID := <-js.chStopJob:
			js.stopService(jobID)

		case <-ticker.C:
			js.startServicesForUnclaimedJobs()
		}
	}
}

func (js *spawner) startServicesForUnclaimedJobs() {
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
		// nodes in the cluster are already running, but it because Postgres
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
				logger.Errorw("Error creating services for job",
					"jobID", specDBRow.ID,
					"error", err,
				)
				continue
			}
			services = append(services, moreServices...)
		}

		logger.Infow("Starting services for job",
			"jobID", specDBRow.ID,
			"count", len(services),
		)

		for _, service := range services {
			err := service.Start()
			if err != nil {
				logger.Errorw("Error creating service for job",
					"jobID", specDBRow.ID,
					"error", err,
				)
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
			logger.Errorw("Error stopping job service",
				"jobID", jobID,
				"error", err,
			)
		} else {
			logger.Infow("Stopped job service",
				"jobID", jobID,
			)
		}
	}
	delete(js.services, jobID)
}

func (js *spawner) CreateJob(spec Spec) (int32, error) {
	delegate, exists := js.jobTypeDelegates[spec.JobType()]
	if !exists {
		panic(fmt.Sprintf("job type '%s' has not been registered with the job.Spawner", spec.JobType()))
	}

	specDBRow := delegate.ToDBRow(spec)
	err := js.orm.CreateJob(&specDBRow, spec.TaskDAG())
	return specDBRow.ID, err
}

func (js *spawner) DeleteJob(ctx context.Context, spec Spec) error {
	err := js.orm.DeleteJob(ctx, spec.JobID())
	if err != nil {
		return err
	}

	select {
	case <-js.chStop:
	case js.chStopJob <- spec.JobID():
	}

	return nil
}
