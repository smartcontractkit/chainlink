package job

import (
	"context"
	"math"
	"reflect"
	"sync"

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

		// NOTE: Prefer to use CreateJob, this is only publicly exposed for use in tests
		// to start a job that was previously manually inserted into DB
		StartService(spec Job) error
	}

	spawner struct {
		orm              ORM
		config           Config
		jobTypeDelegates map[Type]Delegate
		activeJobs       map[int32]activeJob
		activeJobsMu     sync.RWMutex
		txm              postgres.TransactionManager

		utils.StartStopOnce
		chStop chan struct{}
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

var _ Spawner = (*spawner)(nil)

func NewSpawner(orm ORM, config Config, jobTypeDelegates map[Type]Delegate, txm postgres.TransactionManager) *spawner {
	s := &spawner{
		orm:              orm,
		config:           config,
		jobTypeDelegates: jobTypeDelegates,
		txm:              txm,
		activeJobs:       make(map[int32]activeJob),
		chStop:           make(chan struct{}),
	}
	return s
}

func (js *spawner) Start() error {
	return js.StartOnce("JobSpawner", func() error {
		js.startAllServices()
		return nil

	})
}

func (js *spawner) Close() error {
	return js.StopOnce("JobSpawner", func() error {
		close(js.chStop)
		js.stopAllServices()
		return nil

	})
}

func (js *spawner) startAllServices() {
	// TODO: rename to find AllJobs
	specs, _, err := js.orm.JobsV2(0, math.MaxUint32)
	if err != nil {
		logger.Errorf("Couldn't fetch unclaimed jobs: %v", err)
		return
	}

	for _, spec := range specs {
		if err = js.StartService(spec); err != nil {
			logger.Errorf("Couldn't start service %v: %v", spec.Name, err)
		}
	}
}

func (js *spawner) stopAllServices() {
	jobIDs := js.activeJobIDs()
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

func (js *spawner) StartService(spec Job) error {
	js.activeJobsMu.Lock()
	defer js.activeJobsMu.Unlock()

	delegate, exists := js.jobTypeDelegates[spec.Type]
	if !exists {
		logger.Errorw("Job type has not been registered with job.Spawner", "type", spec.Type, "jobID", spec.ID)
		return nil
	}
	// We always add the active job in the activeJob map, even in the case
	// that it fails to start. That way we have access to the delegate to call
	// OnJobDeleted before deleting. However, the activeJob will only have services
	// that it was able to start without an error.
	aj := activeJob{delegate: delegate, spec: spec}

	services, err := delegate.ServicesForSpec(spec)
	if err != nil {
		logger.Errorw("Error creating services for job", "jobID", spec.ID, "error", err)
		ctx, cancel := utils.ContextFromChan(js.chStop)
		defer cancel()
		js.orm.RecordError(ctx, spec.ID, err.Error())
		js.activeJobs[spec.ID] = aj
		return nil
	}

	logger.Debugw("JobSpawner: Starting services for job", "jobID", spec.ID, "count", len(services))

	for _, service := range services {
		err := service.Start()
		if err != nil {
			logger.Errorw("Error creating service for job", "jobID", spec.ID, "error", err)
			continue
		}
		aj.services = append(aj.services, service)
	}
	js.activeJobs[spec.ID] = aj
	return nil
}

// Should not get called before Start()
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

	if err = js.StartService(jb); err != nil {
		return jb, err
	}

	delegate.AfterJobCreated(jb)

	logger.Infow("Created job", "type", jb.Type, "jobID", jb.ID)
	return jb, err
}

// Should not get called before Start()
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

	combctx, cancel := utils.CombinedContext(js.chStop, ctx)
	defer cancel()

	// Inject the wrapping transaction into the combined context if it exists.
	// This allows rollbacks from other services.
	tx := postgres.GetTxFromContext(ctx)
	if tx != nil {
		combctx = postgres.InjectTxIntoContext(combctx, tx)
	}

	err := js.orm.DeleteJob(combctx, jobID)
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

func (js *spawner) activeJobIDs() []int32 {
	js.activeJobsMu.RLock()
	defer js.activeJobsMu.RUnlock()

	ids := make([]int32, 0, len(js.activeJobs))
	for jobID := range js.activeJobs {
		ids = append(ids, jobID)
	}
	return ids
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
