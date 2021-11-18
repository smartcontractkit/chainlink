package job

import (
	"context"
	"math"
	"reflect"
	"sync"

	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/service"
	"github.com/smartcontractkit/chainlink/core/services/pg"
	"github.com/smartcontractkit/chainlink/core/utils"
	"github.com/smartcontractkit/sqlx"
)

//go:generate mockery --name Spawner --output ./mocks/ --case=underscore
//go:generate mockery --name Delegate --output ./mocks/ --case=underscore

type (
	// The job spawner manages the spinning up and spinning down of the long-running
	// services that perform the work described by job specs.  Each active job spec
	// has 1 or more of these services associated with it.
	Spawner interface {
		service.Service
		CreateJob(jb *Job, qopts ...pg.QOpt) error
		DeleteJob(jobID int32, qopts ...pg.QOpt) error
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
		db               *sqlx.DB
		lggr             logger.Logger

		utils.StartStopOnce
		chStop              chan struct{}
		lbDependentAwaiters []utils.DependentAwaiter
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

func NewSpawner(orm ORM, config Config, jobTypeDelegates map[Type]Delegate, db *sqlx.DB, lggr logger.Logger, lbDependentAwaiters []utils.DependentAwaiter) *spawner {
	s := &spawner{
		orm:                 orm,
		config:              config,
		jobTypeDelegates:    jobTypeDelegates,
		db:                  db,
		lggr:                lggr.Named("JobSpawner"),
		activeJobs:          make(map[int32]activeJob),
		chStop:              make(chan struct{}),
		lbDependentAwaiters: lbDependentAwaiters,
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
	specs, _, err := js.orm.FindJobs(0, math.MaxUint32)
	if err != nil {
		js.lggr.Errorf("Couldn't fetch unclaimed jobs: %v", err)
		return
	}

	for _, spec := range specs {
		if err = js.StartService(spec); err != nil {
			js.lggr.Errorf("Couldn't start service %v: %v", spec.Name, err)
		}
	}
	// Log Broadcaster fully starts after all initial Register calls are done from other starting services
	// to make sure the initial backfill covers those subscribers.
	for _, lbd := range js.lbDependentAwaiters {
		lbd.DependentReady()
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
			js.lggr.Errorw("Error stopping job service", "jobID", jobID, "error", err, "subservice", i, "serviceType", reflect.TypeOf(service))
		} else {
			js.lggr.Infow("Stopped job service", "jobID", jobID, "subservice", i, "serviceType", reflect.TypeOf(service))
		}
	}

	delete(js.activeJobs, jobID)
}

func (js *spawner) StartService(spec Job) error {
	js.activeJobsMu.Lock()
	defer js.activeJobsMu.Unlock()

	delegate, exists := js.jobTypeDelegates[spec.Type]
	if !exists {
		js.lggr.Errorw("Job type has not been registered with job.Spawner", "type", spec.Type, "jobID", spec.ID)
		return nil
	}
	// We always add the active job in the activeJob map, even in the case
	// that it fails to start. That way we have access to the delegate to call
	// OnJobDeleted before deleting. However, the activeJob will only have services
	// that it was able to start without an error.
	aj := activeJob{delegate: delegate, spec: spec}

	services, err := delegate.ServicesForSpec(spec)
	if err != nil {
		js.lggr.Errorw("Error creating services for job", "jobID", spec.ID, "error", err)
		ctx, cancel := utils.ContextFromChan(js.chStop)
		defer cancel()
		js.orm.TryRecordError(spec.ID, err.Error(), pg.WithParentCtx(ctx))
		js.activeJobs[spec.ID] = aj
		return nil
	}

	js.lggr.Debugw("JobSpawner: Starting services for job", "jobID", spec.ID, "count", len(services))

	for _, service := range services {
		err := service.Start()
		if err != nil {
			js.lggr.Errorw("Error creating service for job", "jobID", spec.ID, "error", err)
			continue
		}
		aj.services = append(aj.services, service)
	}
	js.activeJobs[spec.ID] = aj
	return nil
}

// Should not get called before Start()
func (js *spawner) CreateJob(jb *Job, qopts ...pg.QOpt) error {
	delegate, exists := js.jobTypeDelegates[jb.Type]
	if !exists {
		js.lggr.Errorf("job type '%s' has not been registered with the job.Spawner", jb.Type)
		return errors.Errorf("job type '%s' has not been registered with the job.Spawner", jb.Type)
	}

	q := pg.NewQ(js.db, qopts...)
	if q.ParentCtx != nil {
		ctx, cancel := utils.CombinedContext(js.chStop, q.ParentCtx)
		defer cancel()
		q.ParentCtx = ctx
	} else {
		ctx, cancel := utils.ContextFromChan(js.chStop)
		defer cancel()
		q.ParentCtx = ctx
	}

	err := js.orm.CreateJob(jb, pg.WithQueryer(q))
	if err != nil {
		js.lggr.Errorw("Error creating job", "type", jb.Type, "error", err)
		return err
	}

	if err = js.StartService(*jb); err != nil {
		return err
	}

	delegate.AfterJobCreated(*jb)

	js.lggr.Infow("Created job", "type", jb.Type, "jobID", jb.ID)
	return err
}

// Should not get called before Start()
func (js *spawner) DeleteJob(jobID int32, qopts ...pg.QOpt) error {
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

	var cancel context.CancelFunc
	defer func() {
		if cancel != nil {
			cancel()
		}
	}()
	setCtx := func(parentCtx context.Context) (ctx context.Context) {
		if parentCtx == nil {
			ctx, cancel = utils.ContextFromChan(js.chStop)
		} else {
			ctx, cancel = utils.CombinedContext(js.chStop, parentCtx)
		}
		return ctx
	}
	err := js.orm.DeleteJob(jobID, append(qopts, pg.MergeCtx(setCtx))...)
	if err != nil {
		js.lggr.Errorw("Error deleting job", "jobID", jobID, "error", err)
		return err
	}

	js.lggr.Infow("Deleted job", "jobID", jobID)

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
