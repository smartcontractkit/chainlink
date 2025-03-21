package job

import (
	"context"
	"fmt"
	"math"
	"reflect"
	"sync"

	pkgerrors "github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"
	"github.com/smartcontractkit/chainlink-common/pkg/utils"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

type (
	// Spawner manages the spinning up and down of the long-running
	// services that perform the work described by job specs.  Each active job spec
	// has 1 or more of these services associated with it.
	Spawner interface {
		services.Service

		// CreateJob creates a new job and starts services.
		// All services must start without errors for the job to be active.
		CreateJob(ctx context.Context, ds sqlutil.DataSource, jb *Job) (err error)
		// DeleteJob deletes a job and stops any active services.
		DeleteJob(ctx context.Context, ds sqlutil.DataSource, jobID int32) error
		// ActiveJobs returns a map of jobs with active services (started without error).
		ActiveJobs() map[int32]Job

		// StartService starts services for the given job spec.
		// NOTE: Prefer to use CreateJob, this is only publicly exposed for use in tests
		// to start a job that was previously manually inserted into DB
		StartService(ctx context.Context, spec Job) error
	}

	Checker interface {
		Register(service services.HealthReporter) error
		Unregister(name string) error
	}

	spawner struct {
		services.StateMachine
		orm              ORM
		config           Config
		checker          Checker
		jobTypeDelegates map[Type]Delegate
		activeJobs       map[int32]activeJob
		activeJobsMu     sync.RWMutex
		lggr             logger.Logger

		chStop              services.StopChan
		lbDependentAwaiters []utils.DependentAwaiter
	}

	// TODO(spook): I can't wait for Go generics
	Delegate interface {
		JobType() Type
		// BeforeJobCreated is only called once on first time job create.
		BeforeJobCreated(Job)
		// ServicesForSpec returns services to be started and stopped for this
		// job. In case a given job type relies upon well-defined startup/shutdown
		// ordering for services, they are started in the order they are given
		// and stopped in reverse order.
		ServicesForSpec(context.Context, Job) ([]ServiceCtx, error)
		AfterJobCreated(Job)
		BeforeJobDeleted(Job)
		// OnDeleteJob will be called from within DELETE db transaction.  Any db
		// commands issued within OnDeleteJob() should be performed first, before any
		// non-db side effects.  This is required in order to guarantee mutual atomicity between
		// all tasks intended to happen during job deletion.  For the same reason, the job will
		// not show up in the db within OnDeleteJob(), even though it is still actively running.
		OnDeleteJob(ctx context.Context, jb Job) error
	}

	activeJob struct {
		delegate Delegate
		spec     Job
		services []ServiceCtx
	}
)

var _ Spawner = (*spawner)(nil)

func NewSpawner(orm ORM, config Config, checker Checker, jobTypeDelegates map[Type]Delegate, lggr logger.Logger, lbDependentAwaiters []utils.DependentAwaiter) *spawner {
	namedLogger := lggr.Named("JobSpawner")
	s := &spawner{
		orm:                 orm,
		config:              config,
		checker:             checker,
		jobTypeDelegates:    jobTypeDelegates,
		lggr:                namedLogger,
		activeJobs:          make(map[int32]activeJob),
		chStop:              make(services.StopChan),
		lbDependentAwaiters: lbDependentAwaiters,
	}
	return s
}

// Start starts Spawner.
func (js *spawner) Start(ctx context.Context) error {
	return js.StartOnce("JobSpawner", func() error {
		js.startAllServices(ctx)
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

func (js *spawner) Name() string {
	return js.lggr.Name()
}

func (js *spawner) HealthReport() map[string]error {
	return map[string]error{js.Name(): js.Healthy()}
}

func (js *spawner) startAllServices(ctx context.Context) {
	// TODO: rename to find AllJobs
	jbs, _, err := js.orm.FindJobs(ctx, 0, math.MaxUint32)
	if err != nil {
		werr := fmt.Errorf("couldn't fetch unclaimed jobs: %w", err)
		js.lggr.Critical(werr.Error())
		js.SvcErrBuffer.Append(werr)
		return
	}

	jobIDs := make([]int32, len(jbs))
	for i, jb := range jbs {
		jobIDs[i] = jb.ID
	}
	js.lggr.Debugw("Starting jobs...", "jobIDs", jobIDs)
	wg := sync.WaitGroup{}
	wg.Add(len(jbs))
	for _, jb := range jbs {
		go func(jb Job) {
			defer wg.Done()
			if err = js.StartService(ctx, jb); err != nil {
				js.lggr.Errorw("Couldn't start job", "jobID", jb.ID, "jobName", jb.Name.ValueOrZero(), "err", err)
			}
		}(jb)
	}
	wg.Wait()
	js.lggr.Debugw("Started jobs", "jobIDs", jobIDs)
	// Log Broadcaster fully starts after all initial Register calls are done from other starting services
	// to make sure the initial backfill covers those subscribers.
	for _, lbd := range js.lbDependentAwaiters {
		lbd.DependentReady()
	}
}

func (js *spawner) stopAllServices() {
	jobIDs := js.activeJobIDs()
	wg := sync.WaitGroup{}
	js.lggr.Debugw("Stopping jobs...", "jobIDs", jobIDs)
	wg.Add(len(jobIDs))
	for _, jobID := range jobIDs {
		go func(jobID int32) {
			defer wg.Done()
			js.stopService(jobID)
		}(jobID)
	}
	wg.Wait()
	js.lggr.Debugw("Stopped jobs", "jobIDs", jobIDs)
}

// stopService removes the job from memory and stop the services.
// It will always delete the job from memory even if closing the services fail.
func (js *spawner) stopService(jobID int32) {
	lggr := js.lggr.With("jobID", jobID)
	js.activeJobsMu.Lock()
	defer js.activeJobsMu.Unlock()

	aj := js.activeJobs[jobID]

	for i := len(aj.services) - 1; i >= 0; i-- {
		service := aj.services[i]
		sLggr := lggr.With("subservice", i, "serviceType", reflect.TypeOf(service))
		if c, ok := service.(services.HealthReporter); ok {
			if err := js.checker.Unregister(c.Name()); err != nil {
				sLggr.Warnw("Failed to unregister service from health checker", "err", err)
			}
		}
		if err := service.Close(); err != nil {
			sLggr.Criticalw("Error stopping job service", "err", err)
			js.SvcErrBuffer.Append(pkgerrors.Wrap(err, "error stopping job service"))
		}
	}

	delete(js.activeJobs, jobID)
}

func (js *spawner) StartService(ctx context.Context, jb Job) error {
	lggr := js.lggr.With("jobID", jb.ID)
	js.activeJobsMu.Lock()
	defer js.activeJobsMu.Unlock()

	delegate, exists := js.jobTypeDelegates[jb.Type]
	if !exists {
		lggr.Errorw("Job type has not been registered with job.Spawner", "type", jb.Type)
		return pkgerrors.Errorf("unregistered type %q for job: %d", jb.Type, jb.ID)
	}
	// We always add the active job in the activeJob map, even in the case
	// that it fails to start. That way we have access to the delegate to call
	// OnJobDeleted before deleting. However, the activeJob will only have services
	// that it was able to start without an error.
	aj := activeJob{delegate: delegate, spec: jb}

	jb.PipelineSpec.JobName = jb.Name.ValueOrZero()
	jb.PipelineSpec.JobID = jb.ID
	jb.PipelineSpec.JobType = string(jb.Type)
	jb.PipelineSpec.ForwardingAllowed = jb.ForwardingAllowed
	if jb.GasLimit.Valid {
		jb.PipelineSpec.GasLimit = &jb.GasLimit.Uint32
	}

	srvs, err := delegate.ServicesForSpec(ctx, jb)
	if err != nil {
		lggr.Errorw("Error creating services for job", "err", err)
		cctx, cancel := js.chStop.NewCtx()
		defer cancel()
		js.orm.TryRecordError(cctx, jb.ID, err.Error())
		js.activeJobs[jb.ID] = aj
		return pkgerrors.Wrapf(err, "failed to create services for job: %d", jb.ID)
	}

	var ms services.MultiStart
	for _, srv := range srvs {
		err = ms.Start(ctx, srv)
		if err != nil {
			lggr.Criticalw("Error starting service for job", "err", err)
			return err
		}
		if c, ok := srv.(services.HealthReporter); ok {
			err = js.checker.Register(c)
			if err != nil {
				lggr.Errorw("Error registering service with health checker", "err", err)
				return err
			}
		}
		aj.services = append(aj.services, srv)
	}
	js.activeJobs[jb.ID] = aj
	return nil
}

// Should not get called before Start()
func (js *spawner) CreateJob(ctx context.Context, ds sqlutil.DataSource, jb *Job) (err error) {
	orm := js.orm
	if ds != nil {
		orm = orm.WithDataSource(ds)
	}
	delegate, exists := js.jobTypeDelegates[jb.Type]
	if !exists {
		js.lggr.Errorf("job type '%s' has not been registered with the job.Spawner", jb.Type)
		err = pkgerrors.Errorf("job type '%s' has not been registered with the job.Spawner", jb.Type)
		return
	}

	err = orm.CreateJob(ctx, jb)
	if err != nil {
		js.lggr.Errorw("Error creating job", "type", jb.Type, "err", err)
		return
	}
	js.lggr.Infow("Created job", "type", jb.Type, "jobID", jb.ID)

	delegate.BeforeJobCreated(*jb)
	err = js.StartService(ctx, *jb)
	if err != nil {
		js.lggr.Errorw("Error starting job services", "type", jb.Type, "jobID", jb.ID, "err", err)
	} else {
		js.lggr.Infow("Started job services", "type", jb.Type, "jobID", jb.ID)
	}

	delegate.AfterJobCreated(*jb)

	return err
}

// Should not get called before Start()
func (js *spawner) DeleteJob(ctx context.Context, ds sqlutil.DataSource, jobID int32) error {
	if ds == nil {
		ds = js.orm.DataSource()
	}
	if jobID == 0 {
		return pkgerrors.New("will not delete job with 0 ID")
	}

	lggr := js.lggr.With("jobID", jobID)
	lggr.Debugw("Deleting job")

	var aj activeJob
	var exists bool
	func() {
		js.activeJobsMu.RLock()
		defer js.activeJobsMu.RUnlock()
		aj, exists = js.activeJobs[jobID]
	}()

	if !exists { // inactive, so look up the spec and delegate
		jb, err := js.orm.WithDataSource(ds).FindJob(ctx, jobID)
		if err != nil {
			return pkgerrors.Wrapf(err, "job %d not found", jobID)
		}
		aj.spec = jb
		if !func() (ok bool) {
			js.activeJobsMu.RLock()
			defer js.activeJobsMu.RUnlock()
			aj.delegate, ok = js.jobTypeDelegates[jb.Type]
			return ok
		}() {
			js.lggr.Errorw("Job type has not been registered with job.Spawner", "type", jb.Type, "jobID", jb.ID)
			return pkgerrors.Errorf("unregistered type %q for job: %d", jb.Type, jb.ID)
		}
	}

	aj.delegate.BeforeJobDeleted(aj.spec)

	err := sqlutil.Transact(ctx, js.orm.WithDataSource, ds, nil, func(tx ORM) error {
		err := tx.DeleteJob(ctx, jobID, aj.spec.Type)
		if err != nil {
			js.lggr.Errorw("Error deleting job", "jobID", jobID, "err", err)
			return err
		}
		// This comes after calling orm.DeleteJob(), so that any non-db side effects inside it only get executed if
		// we know the DELETE will succeed.  The DELETE will be finalized only if all db transactions in OnDeleteJob()
		// succeed.  If either of those fails, the job will not be stopped and everything will be rolled back.
		err = aj.delegate.OnDeleteJob(ctx, aj.spec)
		if err != nil {
			return err
		}

		return nil
	})

	if exists {
		// Stop the service and remove the job from memory, which will always happen even if closing the services fail.
		js.stopService(jobID)
	}
	lggr.Infow("Stopped and deleted job")

	return err
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

// ServicesForSpec does no-op.
func (n *NullDelegate) ServicesForSpec(ctx context.Context, spec Job) (s []ServiceCtx, err error) {
	return
}

func (n *NullDelegate) BeforeJobCreated(spec Job) {}
func (n *NullDelegate) AfterJobCreated(spec Job)  {}
func (n *NullDelegate) BeforeJobDeleted(spec Job) {}
func (n *NullDelegate) OnDeleteJob(context.Context, Job) error {
	return nil
}
