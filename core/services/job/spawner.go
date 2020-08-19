package job

import (
	"github.com/pkg/errors"
	"sync"

	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/store/models"
)

//go:generate mockery --name Spawner --output ../../internal/mocks/ --case=underscore
//go:generate mockery --name JobSpec --output ../../internal/mocks/ --case=underscore
//go:generate mockery --name JobService --output ../../internal/mocks/ --case=underscore
//go:generate mockery --name JobSpawnerORM --output ../../internal/mocks/ --case=underscore

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
		Start() error
		Stop()
		AddJob(jobSpec JobSpec) error
		RemoveJob(jobID *models.ID)
		RegisterJobType(jobType string, factory JobSpecToJobServiceFunc)
	}

	spawner struct {
		orm                   JobSpawnerORM
		jobServiceFactories   map[string]JobSpecToJobServiceFunc
		jobServiceFactoriesMu sync.RWMutex
		chAdd                 chan addEntry
		chRemove              chan models.ID
		chStop                chan struct{}
		chDone                chan struct{}
	}

	JobSpecToJobServiceFunc func(jobSpec JobSpec) ([]JobService, error)

	addEntry struct {
		jobSpec  JobSpec
		services []JobService
	}

	JobSpec interface {
		JobID() *models.ID
		JobType() string
	}

	JobService interface {
		Start() error
		Stop() error
	}

	JobSpawnerORM interface {
		JobsAsInterfaces(fn func(jobSpec JobSpec) bool) error
		UpsertErrorFor(jobID *models.ID, err string)
	}
)

func NewSpawner(orm JobSpawnerORM) *spawner {
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
			jobID := entry.jobSpec.JobID()
			if jobID == nil {
				logger.Errorf("%v job spec has nil job ID", entry.jobSpec.JobType())
				continue
			} else if _, ok := jobMap[*jobID]; ok {
				logger.Errorf("%v job '%s' has already been added", entry.jobSpec.JobType(), jobID.String())
				continue
			}
			for _, service := range entry.services {
				service.Start()
			}
			jobMap[*jobID] = entry.services

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

	services, err := factory(jobSpec)
	if err != nil {
		return err
	} else if len(services) == 0 {
		return nil
	}

	js.chAdd <- addEntry{jobSpec, services}
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
