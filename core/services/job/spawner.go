package job

import (
	"sync"
	"time"

	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/core/logger"
)

//go:generate mockery --name Spawner --output ./mocks/ --case=underscore

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
		CreateJob(spec Spec) error
		DeleteJob(spec Spec) error
		RegisterJobType(registration Registration)
	}

	spawner struct {
		orm                    ORM
		jobTypeRegistrations   map[Type]Registration
		jobTypeRegistrationsMu sync.RWMutex
		services               map[int32][]Service
		chStopJob              chan int32
		chStop                 chan struct{}
		chDone                 chan struct{}
	}

	ServicesFactory func(spec Spec) ([]Service, error)

	Registration struct {
		JobType         Type
		Spec            Spec
		ServicesFactory ServicesFactory
	}
)

var _ Spawner = (*spawner)(nil)

func NewSpawner(orm ORM) *spawner {
	return &spawner{
		orm:                  orm,
		jobTypeRegistrations: make(map[Type]Registration),
		services:             make(map[int32][]Service),
		chStopJob:            make(chan int32),
		chStop:               make(chan struct{}),
		chDone:               make(chan struct{}),
	}
}

func (js *spawner) Start() {
	go js.runLoop()
}

func (js *spawner) Stop() {
	close(js.chStop)
	<-js.chDone
}

func (js *spawner) RegisterJobType(registration Registration) {
	js.jobTypeRegistrationsMu.Lock()
	defer js.jobTypeRegistrationsMu.Unlock()

	js.jobTypeRegistrations[registration.JobType] = registration
}

func (js *spawner) runLoop() {
	defer close(js.chDone)

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	js.startUnclaimedServices()

	for {
		select {
		case <-js.chStop:
			js.stopAllServices()
			return

		case jobID := <-js.chStopJob:
			js.stopService(jobID)

		case <-ticker.C:
			js.startUnclaimedServices()
		}
	}
}

func (js *spawner) startUnclaimedServices() {
	jobSpecs, err := js.orm.UnclaimedJobs(js.jobTypeRegistrations, js.services)
	if err != nil {
		logger.Errorf("error fetching unclaimed jobs: %v", err)
		return
	}

	for _, jobSpec := range jobSpecs {
		err := js.startService(jobSpec)
		if err != nil {
			logger.Errorw("error starting job service",
				"job type", jobSpec.JobType(),
				"job id", jobSpec.JobID(),
				"error", err,
			)
		}
	}
}

func (js *spawner) startService(jobSpec Spec) error {
	js.jobTypeRegistrationsMu.RLock()
	defer js.jobTypeRegistrationsMu.RUnlock()

	logger.Infow("Starting service for job",
		"jobID", jobSpec.JobID(),
		"jobType", jobSpec.JobType(),
	)

	reg, exists := js.jobTypeRegistrations[jobSpec.JobType()]
	if !exists {
		return errors.Errorf("Job Spawner got unknown job type '%v'", jobSpec.JobType())
	}

	services, err := reg.ServicesFactory(jobSpec)
	if err != nil {
		return err
	}

	for _, service := range services {
		err := service.Start()
		if err != nil {
			return err
		}
		js.services[*jobSpec.JobID()] = append(js.services[*jobSpec.JobID()], service)
	}
	return nil
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
			logger.Errorw("error stopping job",
				"job id", jobID,
				"error", err,
			)
		}
	}
	delete(js.services, jobID)
}

func (js *spawner) CreateJob(spec Spec) error {
	return js.orm.CreateJob(spec)
}

func (js *spawner) DeleteJob(spec Spec) error {
	if spec.JobID() == nil {
		return errors.New("Job Spawner could not delete job: got nil job ID")
	}

	err := js.orm.DeleteJob(spec)
	if err != nil {
		return err
	}

	select {
	case <-js.chStop:
	case js.chStopJob <- *spec.JobID():
	}

	return nil
}
