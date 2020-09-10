package pipeline

import (
	"github.com/smartcontractkit/chainlink/core/null"
	"sync"
	"time"

	"github.com/guregu/null.v4"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/store/models"
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
		Start() error
		Stop()
		CreateJob(jobSpec Spec) error
		StopJob(jobID models.ID)
		RegisterJobType(jobType JobType, factory JobSpecToJobServiceFunc)
	}

	spawner struct {
		orm                   ORM
		jobServiceFactories   map[JobType]JobSpecToJobServiceFunc
		jobServiceFactoriesMu sync.RWMutex
		jobServices           map[models.ID][]JobService
		chStopJob             chan models.ID
		chStop                chan struct{}
		chDone                chan struct{}
	}

	JobSpecToJobServiceFunc func(jobSpec JobSpec) ([]JobService, error)
)

func NewSpawner(orm ORM) *spawner {
	return &spawner{
		orm:                 orm,
		jobServiceFactories: make(map[JobType]JobSpecToJobServiceFunc),
		jobServices:         make(map[models.ID][]JobService),
		chStopJob:           make(chan models.ID),
		chStop:              make(chan struct{}),
		chDone:              make(chan struct{}),
	}
}

func (js *spawner) Start() {
	go js.runLoop()
}

func (js *spawner) Stop() {
	close(js.chStop)
	<-js.chDone
}

func (js *spawner) runLoop() {
	defer close(js.chDone)

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	js.startUnclaimedJobServices()

	for {
		select {
		case <-js.chStop:
			js.stopAllJobServices()
			return
		case jobID := <-js.chStopJob:
			js.stopJobService(jobID)
		case <-ticker.C:
			js.startUnclaimedJobServices()
		}
	}
}

func (js *spawner) startUnclaimedJobServices() {
	// NOTE: .UnclaimedJobs() should automatically lock/claim the jobs
	jobSpecs, err := js.orm.UnclaimedJobs()
	if err != nil {
		logger.Errorf("error fetching unclaimed jobs: %v", err)
		return
	}

	for _, jobSpec := range jobSpecs {
		err := js.startJobService(jobSpec)
		if err != nil {
			logger.Errorw("error starting job service",
				"job type", jobSpec.JobType(),
				"job id", jobSpec.JobID(),
				"error", err,
			)
			// TODO: un-claim the job
		}
	}
}

func (js *spawner) startJobService(jobSpec JobSpec) error {
	js.jobServiceFactoriesMu.RLock()
	defer js.jobServiceFactoriesMu.RUnlock()

	factory, exists := js.jobServiceFactories[jobSpec.JobType()]
	if !exists {
		return errors.Errorf("Job Spawner got unknown job type '%v'", jobSpec.JobType())
	}

	services, err := factory(jobSpec)
	if err != nil {
		return err
	}

	for _, service := range services {
		err := service.Start()
		if err != nil {
			return err
		}
		js.jobServices[*jobSpec.JobID()] = append(js.jobServices[*jobSpec.JobID()], service)
	}
	return nil
}

func (js *spawner) stopAllJobServices() {
	for jobID := range js.jobServices {
		js.stopJobService(jobID)
	}
}

func (js *spawner) stopJobService(jobID models.ID) {
	for _, service := range js.jobServices[jobID] {
		err := service.Stop()
		if err != nil {
			logger.Errorw("error stopping job",
				"job id", jobID,
				"error", err,
			)
		}
	}
	delete(js.jobServices, jobID)
}

func (js *spawner) CreateJob(spec JobSpec) error {
	return js.db.Transaction(func(tx *gorm.DB) error {
		// Save the spec to the DB
		err := tx.Create(spec)
		if err != nil {
			return err
		}

		// Convert the task DAG into TaskSpec DB rows
		taskSpecs := []PipelineTaskSpec{}
		taskSpecIDs := make(map[Task]int64)
		err = spec.TaskDAG().ReverseWalkTasks(func(task Task) error {
			var successorID null.Int64
			if len(task.OutputTasks()) > 1 {
				return errors.New("task has > 1 output task")

			} else if len(task.OutputTasks()) == 1 {
				successor := task.OutputTasks()[0]
				successorID = null.Int64From(taskSpecIDs[successor])
			}

			taskSpec := PipelineTaskSpec{
				TaskJson:    JSONSerializable{task},
				SuccessorID: successorID,
			}

			err := tx.Create(&taskSpec).Error
			if err != nil {
				return err
			}

			taskSpecIDs[task] = taskSpec
			taskSpecs = append(taskSpecs, taskSpec)
			return nil
		})

		pipelineSpec := PipelineSpec{
			JobSpecID:    spec.JobID(),
			SourceDotDag: spec.TaskDAG().DOTSource,
			TaskSpecs:    taskSpecs,
		}
		return tx.Create(&pipelineSpec).Error
	})
}

func (js *spawner) StopJob(id models.ID) {
	select {
	case <-js.chStop:
	case js.chStopJob <- id:
	}
}

func (js *spawner) RegisterJobType(jobType JobType, factory JobSpecToJobServiceFunc) {
	js.jobServiceFactoriesMu.Lock()
	defer js.jobServiceFactoriesMu.Unlock()

	js.jobServiceFactories[jobType] = factory
}
