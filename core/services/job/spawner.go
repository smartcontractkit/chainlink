package job

import (
	"fmt"
	"github.com/pkg/errors"
	"sync"

	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/eth"
	"github.com/smartcontractkit/chainlink/core/store/models"
)

// The job spawner manages the spinning up and spinning down of the long-running
// services that perform the work described by job specs.  Each active job spec
// has 1 or more of these services associated with it.
//
// At present, Flux Monitor and Offchain Reporting jobs can only have a single
// "initiator", meaning that they only require a single service.  But the older
// "direct request" model allows for multiple initiators, which imply multiple
// services.
type (
	Spawner interface {
		Start() error
		Stop()
		AddJob(jobSpec JobSpec) error
		RemoveJob(jobID *models.ID)
		RegisterJobType(jobType string, factory JobSpecToJobServiceFunc)
	}

	spawner struct {
		orm                   ormInterface
		jobServiceFactories   map[string]JobSpecToJobServiceFunc
		jobServiceFactoriesMu sync.RWMutex
		chAdd                 chan addEntry
		chRemove              chan models.ID
		chStop                chan struct{}
		chDone                chan struct{}
	}

	JobSpecToJobServiceFunc func(jobSpec JobSpec) (JobService, error)

	addEntry struct {
		jobSpec  JobSpec
		services []JobService
	}

	JobSpec interface {
		JobID() *models.ID
		JobType() string
	}

	JobService interface {
		Start()
		Stop()
	}

	ormInterface interface {
		JobsAsInterfaces(fn func(jobSpec JobSpec) bool) error
		UpsertErrorFor(jobID *models.ID, err string)
	}
)

func NewSpawner(orm ormInterface) *spawner {
	return &spawner{
		orm:                 orm,
		jobServiceFactories: make(map[string]JobSpecToJobServiceFunc),
		chAdd:               make(chan addEntry),
		chRemove:            make(chan models.ID),
		chStop:              make(chan struct{}),
		chDone:              make(chan struct{}),
	}
}

func (js *spawner) Start() error {
	go js.runLoop()

	// Add all of the jobs that we already have in the DB
	var wg sync.WaitGroup
	err := js.orm.JobsAsInterfaces(func(jobSpec JobSpec) bool {
		if jobSpec == nil {
			err := errors.New("received nil job")
			logger.Error(err)
			return true
		}

		wg.Add(1)
		go func() {
			defer wg.Done()

			err := js.AddJob(jobSpec)
			if err != nil {
				logger.Errorf("error adding %v job: %v", jobSpec.JobType(), err)
			}
		}()
		return true
	})

	wg.Wait()

	return err
}

func (js *spawner) Stop() {
	close(js.chStop)
	<-js.chDone
}

func (js *spawner) runLoop() {
	defer close(js.chDone)

	jobMap := map[models.ID][]JobService{}

	for {
		select {
		case entry := <-js.chAdd:
			if _, ok := jobMap[entry.jobSpec.JobID()]; ok {
				logger.Errorf("%v job '%s' has already been added", entry.jobSpec.JobType(), entry.jobSpec.JobID().String())
				continue
			}
			for _, service := range entry.services {
				service.Start()
			}
			jobMap[entry.jobSpec.JobID()] = entry.services

		case jobID := <-js.chRemove:
			services, ok := jobMap[jobID]
			if !ok {
				logger.Debugf("job '%s' is missing", jobID.String())
				continue
			}
			for _, service := range services {
				service.Stop()
			}
			delete(jobMap, jobID)

		case <-js.chStop:
			for _, services := range jobMap {
				for _, service := range services {
					service.Stop()
				}
			}
			return
		}
	}
}

func (js *spawner) AddJob(jobSpec JobSpec) error {
	if jobSpec.JobID() == nil {
		err := errors.New("Job Spawner received job with nil ID")
		logger.Error(err)
		js.orm.UpsertErrorFor(jobSpec.JobID(), "Unable to add job - job has nil ID")
		return err
	}

	js.jobServiceFactoriesMu.RLock()
	defer js.jobServiceFactoriesMu.RUnlock()

	factory, exists := js.jobServiceFactories[jobSpec.JobType()]
	if !exists {
		return errors.Errorf("Job Spawner got unknown job type '%v'", jobSpec.JobType())
	}

	services := factory(jobSpec)

	if len(services) == 0 {
		return nil
	}

	js.chAdd <- addEntry{*jobSpec.ID, services}
	return nil
}

func (js *spawner) RemoveJob(id *models.ID) {
	if id == nil {
		logger.Warn("nil job ID passed to Spawner#RemoveJob")
		return
	}
	js.chRemove <- *id
}

func (js *spawner) RegisterJobType(jobType string, factory JobSpecToJobServiceFunc) {
	js.jobServiceFactoriesMu.Lock()
	defer js.jobServiceFactoriesMu.Unlock()

	js.jobServiceFactories[jobType] = factory
}
