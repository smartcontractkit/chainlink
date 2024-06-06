package pipeline

import (
	"context"
	"fmt"
	"net/http"
	"sort"
	"sync"
	"time"

	"github.com/google/uuid"
	pkgerrors "github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"gopkg.in/guregu/null.v4"

	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"
	commonutils "github.com/smartcontractkit/chainlink-common/pkg/utils"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/jsonserializable"

	"github.com/smartcontractkit/chainlink/v2/core/bridges"
	"github.com/smartcontractkit/chainlink/v2/core/chains/legacyevm"
	"github.com/smartcontractkit/chainlink/v2/core/config/env"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/recovery"
	"github.com/smartcontractkit/chainlink/v2/core/store/models"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

//go:generate mockery --quiet --name Runner --output ./mocks/ --case=underscore

type Runner interface {
	services.Service

	// Run is a blocking call that will execute the run until no further progress can be made.
	// If `incomplete` is true, the run is only partially complete and is suspended, awaiting to be resumed when more data comes in.
	// Note that `saveSuccessfulTaskRuns` value is ignored if the run contains async tasks.
	Run(ctx context.Context, run *Run, l logger.Logger, saveSuccessfulTaskRuns bool, fn func(tx sqlutil.DataSource) error) (incomplete bool, err error)
	ResumeRun(ctx context.Context, taskID uuid.UUID, value interface{}, err error) error

	// ExecuteRun executes a new run in-memory according to a spec and returns the results.
	// We expect spec.JobID and spec.JobName to be set for logging/prometheus.
	ExecuteRun(ctx context.Context, spec Spec, vars Vars, l logger.Logger) (run *Run, trrs TaskRunResults, err error)
	// InsertFinishedRun saves the run results in the database.
	// ds is an optional override, for example when executing a transaction.
	InsertFinishedRun(ctx context.Context, ds sqlutil.DataSource, run *Run, saveSuccessfulTaskRuns bool) error
	InsertFinishedRuns(ctx context.Context, ds sqlutil.DataSource, runs []*Run, saveSuccessfulTaskRuns bool) error

	// ExecuteAndInsertFinishedRun executes a new run in-memory according to a spec, persists and saves the results.
	// It is a combination of ExecuteRun and InsertFinishedRun.
	// Note that the spec MUST have a DOT graph for this to work.
	// This will persist the Spec in the DB if it doesn't have an ID.
	ExecuteAndInsertFinishedRun(ctx context.Context, spec Spec, vars Vars, l logger.Logger, saveSuccessfulTaskRuns bool) (runID int64, results TaskRunResults, err error)

	OnRunFinished(func(*Run))
	InitializePipeline(spec Spec) (*Pipeline, error)
}

type runner struct {
	services.StateMachine
	orm                    ORM
	btORM                  bridges.ORM
	config                 Config
	bridgeConfig           BridgeConfig
	legacyEVMChains        legacyevm.LegacyChainContainer
	ethKeyStore            ETHKeyStore
	vrfKeyStore            VRFKeyStore
	runReaperWorker        *commonutils.SleeperTask
	lggr                   logger.Logger
	httpClient             *http.Client
	unrestrictedHTTPClient *http.Client

	// test helper
	runFinished func(*Run)

	chStop services.StopChan
	wgDone sync.WaitGroup
}

var (
	// PromPipelineTaskExecutionTime reports how long each pipeline task took to execute
	// TODO: Make private again after
	// https://app.clubhouse.io/chainlinklabs/story/6065/hook-keeper-up-to-use-tasks-in-the-pipeline
	PromPipelineTaskExecutionTime = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "pipeline_task_execution_time",
		Help: "How long each pipeline task took to execute",
	},
		[]string{"job_id", "job_name", "task_id", "task_type"},
	)
	PromPipelineRunErrors = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "pipeline_run_errors",
		Help: "Number of errors for each pipeline spec",
	},
		[]string{"job_id", "job_name"},
	)
	PromPipelineRunTotalTimeToCompletion = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "pipeline_run_total_time_to_completion",
		Help: "How long each pipeline run took to finish (from the moment it was created)",
	},
		[]string{"job_id", "job_name"},
	)
	PromPipelineTasksTotalFinished = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "pipeline_tasks_total_finished",
		Help: "The total number of pipeline tasks which have finished",
	},
		[]string{"job_id", "job_name", "task_id", "task_type", "bridge_name", "status"},
	)
)

func NewRunner(orm ORM, btORM bridges.ORM, cfg Config, bridgeCfg BridgeConfig, legacyChains legacyevm.LegacyChainContainer, ethks ETHKeyStore, vrfks VRFKeyStore, lggr logger.Logger, httpClient, unrestrictedHTTPClient *http.Client) *runner {
	r := &runner{
		orm:                    orm,
		btORM:                  btORM,
		config:                 cfg,
		bridgeConfig:           bridgeCfg,
		legacyEVMChains:        legacyChains,
		ethKeyStore:            ethks,
		vrfKeyStore:            vrfks,
		chStop:                 make(chan struct{}),
		wgDone:                 sync.WaitGroup{},
		runFinished:            func(*Run) {},
		lggr:                   lggr.Named("PipelineRunner"),
		httpClient:             httpClient,
		unrestrictedHTTPClient: unrestrictedHTTPClient,
	}
	r.runReaperWorker = commonutils.NewSleeperTask(
		commonutils.SleeperFuncTask(r.runReaper, "PipelineRunnerReaper"),
	)
	return r
}

// Start starts Runner.
func (r *runner) Start(context.Context) error {
	return r.StartOnce("PipelineRunner", func() error {
		r.wgDone.Add(1)
		go r.scheduleUnfinishedRuns()
		if r.config.ReaperInterval() != time.Duration(0) {
			r.wgDone.Add(1)
			go r.runReaperLoop()
		}
		return nil
	})
}

func (r *runner) Close() error {
	return r.StopOnce("PipelineRunner", func() error {
		close(r.chStop)
		r.wgDone.Wait()
		return nil
	})
}

func (r *runner) Name() string {
	return r.lggr.Name()
}

func (r *runner) HealthReport() map[string]error {
	return map[string]error{r.Name(): r.Healthy()}
}

func (r *runner) destroy() {
	err := r.runReaperWorker.Stop()
	if err != nil {
		r.lggr.Error(err)
	}
}

func (r *runner) runReaperLoop() {
	defer r.wgDone.Done()
	defer r.destroy()
	if r.config.ReaperInterval() == 0 {
		return
	}

	runReaperTicker := time.NewTicker(utils.WithJitter(r.config.ReaperInterval()))
	defer runReaperTicker.Stop()
	for {
		select {
		case <-r.chStop:
			return
		case <-runReaperTicker.C:
			r.runReaperWorker.WakeUp()
			runReaperTicker.Reset(utils.WithJitter(r.config.ReaperInterval()))
		}
	}
}

type memoryTaskRun struct {
	task     Task
	inputs   []Result // sorted by input index
	vars     Vars
	attempts uint
}

// When a task panics, we catch the panic and wrap it in an error for reporting to the scheduler.
type ErrRunPanicked struct {
	v interface{}
}

func (err ErrRunPanicked) Error() string {
	return fmt.Sprintf("goroutine panicked when executing run: %v", err.v)
}

func NewRun(spec Spec, vars Vars) *Run {
	return &Run{
		State:          RunStatusRunning,
		JobID:          spec.JobID,
		PruningKey:     spec.JobID,
		PipelineSpec:   spec,
		PipelineSpecID: spec.ID,
		Inputs:         jsonserializable.JSONSerializable{Val: vars.vars, Valid: true},
		Outputs:        jsonserializable.JSONSerializable{Val: nil, Valid: false},
		CreatedAt:      time.Now(),
	}
}

func (r *runner) OnRunFinished(fn func(*Run)) {
	r.runFinished = fn
}

var (
	// github.com/smartcontractkit/libocr/offchainreporting2plus/internal/protocol.ReportingPluginTimeoutWarningGracePeriod
	overtime           = 100 * time.Millisecond
	overtimeThresholds = sqlutil.LogThresholds{
		Warn: func(timeout time.Duration) time.Duration {
			return timeout - (timeout / 5) // 80%
		},
		Error: func(timeout time.Duration) time.Duration {
			return timeout - (timeout / 10) // 90%
		},
	}
)

func init() {
	// undocumented escape hatch
	if v := env.PipelineOvertime.Get(); v != "" {
		d, err := time.ParseDuration(v)
		if err == nil {
			overtime = d
		}
	}
}

// overtimeContext returns a modified context for overtime work, since tasks are expected to keep running and return
// results, even after context cancellation.
func overtimeContext(ctx context.Context) (context.Context, context.CancelFunc) {
	ctx = overtimeThresholds.ContextWithValue(ctx)
	if d, ok := ctx.Deadline(); ok {
		// extend deadline
		return context.WithDeadline(context.WithoutCancel(ctx), d.Add(overtime))
	}
	// remove cancellation
	return context.WithoutCancel(ctx), func() {}
}

func (r *runner) ExecuteRun(
	ctx context.Context,
	spec Spec,
	vars Vars,
	l logger.Logger,
) (*Run, TaskRunResults, error) {
	// Pipeline runs may return results after the context is cancelled, so we modify the
	// deadline to give them time to return before the parent context deadline.
	var cancel func()
	ctx, cancel = commonutils.ContextWithDeadlineFn(ctx, func(orig time.Time) time.Time {
		if tenPct := time.Until(orig) / 10; overtime > tenPct {
			return orig.Add(-tenPct)
		}
		return orig.Add(-overtime)
	})
	defer cancel()

	var pipeline *Pipeline
	if spec.Pipeline != nil {
		// assume if set that it has been pre-initialized
		pipeline = spec.Pipeline
	} else {
		var err error
		pipeline, err = r.InitializePipeline(spec)
		if err != nil {
			return nil, nil, err
		}
	}

	run := NewRun(spec, vars)
	taskRunResults := r.run(ctx, pipeline, run, vars, l)

	if run.Pending {
		return run, nil, fmt.Errorf("unexpected async run for spec ID %v, tried executing via ExecuteRun", spec.ID)
	}

	return run, taskRunResults, nil
}

func (r *runner) InitializePipeline(spec Spec) (pipeline *Pipeline, err error) {
	pipeline, err = spec.GetOrParsePipeline()
	if err != nil {
		return
	}

	// initialize certain task params
	for _, task := range pipeline.Tasks {
		task.Base().uuid = uuid.New()

		switch task.Type() {
		case TaskTypeHTTP:
			task.(*HTTPTask).config = r.config
			task.(*HTTPTask).httpClient = r.httpClient
			task.(*HTTPTask).unrestrictedHTTPClient = r.unrestrictedHTTPClient
		case TaskTypeBridge:
			task.(*BridgeTask).config = r.config
			task.(*BridgeTask).bridgeConfig = r.bridgeConfig
			task.(*BridgeTask).orm = r.btORM
			task.(*BridgeTask).specId = spec.ID
			// URL is "safe" because it comes from the node's own database. We
			// must use the unrestrictedHTTPClient because some node operators
			// may run external adapters on their own hardware
			task.(*BridgeTask).httpClient = r.unrestrictedHTTPClient
		case TaskTypeETHCall:
			task.(*ETHCallTask).legacyChains = r.legacyEVMChains
			task.(*ETHCallTask).config = r.config
			task.(*ETHCallTask).specGasLimit = spec.GasLimit
			task.(*ETHCallTask).jobType = spec.JobType
		case TaskTypeVRF:
			task.(*VRFTask).keyStore = r.vrfKeyStore
		case TaskTypeVRFV2:
			task.(*VRFTaskV2).keyStore = r.vrfKeyStore
		case TaskTypeVRFV2Plus:
			task.(*VRFTaskV2Plus).keyStore = r.vrfKeyStore
		case TaskTypeEstimateGasLimit:
			task.(*EstimateGasLimitTask).legacyChains = r.legacyEVMChains
			task.(*EstimateGasLimitTask).specGasLimit = spec.GasLimit
			task.(*EstimateGasLimitTask).jobType = spec.JobType
		case TaskTypeETHTx:
			task.(*ETHTxTask).keyStore = r.ethKeyStore
			task.(*ETHTxTask).legacyChains = r.legacyEVMChains
			task.(*ETHTxTask).specGasLimit = spec.GasLimit
			task.(*ETHTxTask).jobType = spec.JobType
			task.(*ETHTxTask).forwardingAllowed = spec.ForwardingAllowed
		default:
		}
	}

	return pipeline, nil
}

func (r *runner) run(ctx context.Context, pipeline *Pipeline, run *Run, vars Vars, l logger.Logger) TaskRunResults {
	l = l.With("run.ID", run.ID, "executionID", uuid.New(), "specID", run.PipelineSpecID, "jobID", run.PipelineSpec.JobID, "jobName", run.PipelineSpec.JobName)
	l.Debug("Initiating tasks for pipeline run of spec")

	scheduler := newScheduler(pipeline, run, vars, l)
	go scheduler.Run()

	// This is "just in case" for cleaning up any stray reports.
	// Normally the scheduler loop doesn't stop until all in progress runs report back
	reportCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if pipelineTimeout := r.config.MaxRunDuration(); pipelineTimeout != 0 {
		ctx, cancel = context.WithTimeout(ctx, pipelineTimeout)
		defer cancel()
	}

	for taskRun := range scheduler.taskCh {
		taskRun := taskRun
		// execute
		go recovery.WrapRecoverHandle(l, func() {
			result := r.executeTaskRun(ctx, run.PipelineSpec, taskRun, l)

			logTaskRunToPrometheus(result, run.PipelineSpec)

			scheduler.report(reportCtx, result)
		}, func(err interface{}) {
			t := time.Now()
			scheduler.report(reportCtx, TaskRunResult{
				ID:         uuid.New(),
				Task:       taskRun.task,
				Result:     Result{Error: ErrRunPanicked{err}},
				FinishedAt: null.TimeFrom(t),
				CreatedAt:  t, // TODO: more accurate start time
			})
		})
	}

	// if the run is suspended, awaiting resumption
	run.Pending = scheduler.pending
	// scheduler.exiting = we had an error and the task was marked to failEarly
	run.FailSilently = scheduler.exiting
	run.State = RunStatusSuspended

	var runTime time.Duration
	if !scheduler.pending {
		run.FinishedAt = null.TimeFrom(time.Now())

		// NOTE: runTime can be very long now because it'll include suspend
		runTime = run.FinishedAt.Time.Sub(run.CreatedAt)
		PromPipelineRunTotalTimeToCompletion.WithLabelValues(fmt.Sprintf("%d", run.PipelineSpec.JobID), run.PipelineSpec.JobName).Set(float64(runTime))
	}

	// Update run results
	run.PipelineTaskRuns = nil
	for _, result := range scheduler.results {
		output := result.Result.OutputDB()
		run.PipelineTaskRuns = append(run.PipelineTaskRuns, TaskRun{
			ID:            result.ID,
			PipelineRunID: run.ID,
			Type:          result.Task.Type(),
			Index:         result.Task.OutputIndex(),
			Output:        output,
			Error:         result.Result.ErrorDB(),
			DotID:         result.Task.DotID(),
			CreatedAt:     result.CreatedAt,
			FinishedAt:    result.FinishedAt,
			task:          result.Task,
		})

		sort.Slice(run.PipelineTaskRuns, func(i, j int) bool {
			if run.PipelineTaskRuns[i].task.OutputIndex() == run.PipelineTaskRuns[j].task.OutputIndex() {
				return run.PipelineTaskRuns[i].FinishedAt.ValueOrZero().Before(run.PipelineTaskRuns[j].FinishedAt.ValueOrZero())
			}
			return run.PipelineTaskRuns[i].task.OutputIndex() < run.PipelineTaskRuns[j].task.OutputIndex()
		})
	}

	// Update run errors/outputs
	if run.FinishedAt.Valid {
		var errors []null.String
		var fatalErrors []null.String
		var outputs []interface{}
		for _, result := range run.PipelineTaskRuns {
			if result.Error.Valid {
				errors = append(errors, result.Error)
			}
			// skip non-terminal results
			if len(result.task.Outputs()) != 0 {
				continue
			}
			fatalErrors = append(fatalErrors, result.Error)
			outputs = append(outputs, result.Output.Val)
		}
		run.AllErrors = errors
		run.FatalErrors = fatalErrors
		run.Outputs = jsonserializable.JSONSerializable{Val: outputs, Valid: true}

		if run.HasFatalErrors() {
			run.State = RunStatusErrored
			PromPipelineRunErrors.WithLabelValues(fmt.Sprintf("%d", run.PipelineSpec.JobID), run.PipelineSpec.JobName).Inc()
		} else {
			run.State = RunStatusCompleted
		}
	}

	// TODO: drop this once we stop using TaskRunResults
	var taskRunResults TaskRunResults
	for _, result := range scheduler.results {
		taskRunResults = append(taskRunResults, result)
	}

	var idxs []int32
	for i := range taskRunResults {
		idxs = append(idxs, taskRunResults[i].Task.OutputIndex())
	}
	// Ensure that task run results are ordered by their output index
	sort.SliceStable(taskRunResults, func(i, j int) bool {
		return taskRunResults[i].Task.OutputIndex() < taskRunResults[j].Task.OutputIndex()
	})
	for i := range taskRunResults {
		idxs[i] = taskRunResults[i].Task.OutputIndex()
	}

	if r.config.VerboseLogging() {
		l = l.With(
			"run.PipelineTaskRuns", run.PipelineTaskRuns,
			"run.Outputs", run.Outputs,
			"run.CreatedAt", run.CreatedAt,
			"run.FinishedAt", run.FinishedAt,
			"run.Meta", run.Meta,
			"run.Inputs", run.Inputs,
		)
	}
	l = l.With("run.State", run.State, "fatal", run.HasFatalErrors(), "runTime", runTime)
	if run.HasFatalErrors() {
		// This will also log at error level in OCR if it fails Observe so the
		// level is appropriate
		l = l.With("run.FatalErrors", run.FatalErrors)
		l.Debugw("Completed pipeline run with fatal errors")
	} else if run.HasErrors() {
		l = l.With("run.AllErrors", run.AllErrors)
		l.Debugw("Completed pipeline run with errors")
	} else {
		l.Debugw("Completed pipeline run successfully")
	}

	return taskRunResults
}

func (r *runner) executeTaskRun(ctx context.Context, spec Spec, taskRun *memoryTaskRun, l logger.Logger) TaskRunResult {
	start := time.Now()
	l = l.With("taskName", taskRun.task.DotID(),
		"taskType", taskRun.task.Type(),
		"attempt", taskRun.attempts)

	// Task timeout will be whichever of the following timesout/cancels first:
	// - Pipeline-level timeout
	// - Specific task timeout (task.TaskTimeout)
	// - Job level task timeout (spec.MaxTaskDuration)
	// - Passed in context

	// CAUTION: Think twice before changing any of the context handling code
	// below. It has already been changed several times trying to "fix" a bug,
	// but actually introducing new ones. Please leave it as-is unless you have
	// an extremely good reason to change it.
	ctx, cancel := r.chStop.Ctx(ctx)
	defer cancel()
	if taskTimeout, isSet := taskRun.task.TaskTimeout(); isSet && taskTimeout > 0 {
		ctx, cancel = context.WithTimeout(ctx, taskTimeout)
		defer cancel()
	}
	if spec.MaxTaskDuration != models.Interval(time.Duration(0)) {
		ctx, cancel = context.WithTimeout(ctx, time.Duration(spec.MaxTaskDuration))
		defer cancel()
	}

	result, runInfo := taskRun.task.Run(ctx, l, taskRun.vars, taskRun.inputs)
	loggerFields := []interface{}{"runInfo", runInfo,
		"resultValue", result.Value,
		"resultError", result.Error,
		"resultType", fmt.Sprintf("%T", result.Value),
	}
	switch v := result.Value.(type) {
	case []byte:
		loggerFields = append(loggerFields, "resultString", fmt.Sprintf("%q", v))
		loggerFields = append(loggerFields, "resultHex", fmt.Sprintf("%x", v))
	}
	if r.config.VerboseLogging() {
		l.Tracew("Pipeline task completed", loggerFields...)
	}

	now := time.Now()

	var finishedAt null.Time
	if !runInfo.IsPending {
		finishedAt = null.TimeFrom(now)
	}
	return TaskRunResult{
		ID:         taskRun.task.Base().uuid,
		Task:       taskRun.task,
		Result:     result,
		CreatedAt:  start,
		FinishedAt: finishedAt,
		runInfo:    runInfo,
	}
}

func logTaskRunToPrometheus(trr TaskRunResult, spec Spec) {
	elapsed := trr.FinishedAt.Time.Sub(trr.CreatedAt)

	PromPipelineTaskExecutionTime.WithLabelValues(fmt.Sprintf("%d", spec.JobID), spec.JobName, trr.Task.DotID(), string(trr.Task.Type())).Set(float64(elapsed))
	var status string
	if trr.Result.Error != nil {
		status = "error"
	} else {
		status = "completed"
	}

	bridgeName := ""
	if bridgeTask, ok := trr.Task.(*BridgeTask); ok {
		bridgeName = bridgeTask.Name
	}

	PromPipelineTasksTotalFinished.WithLabelValues(fmt.Sprintf("%d", spec.JobID), spec.JobName, trr.Task.DotID(), string(trr.Task.Type()), bridgeName, status).Inc()
}

// ExecuteAndInsertFinishedRun executes a run in memory then inserts the finished run/task run records, returning the final result
func (r *runner) ExecuteAndInsertFinishedRun(ctx context.Context, spec Spec, vars Vars, l logger.Logger, saveSuccessfulTaskRuns bool) (runID int64, results TaskRunResults, err error) {
	run, trrs, err := r.ExecuteRun(ctx, spec, vars, l)
	if err != nil {
		return 0, trrs, pkgerrors.Wrapf(err, "error executing run for spec ID %v", spec.ID)
	}

	// don't insert if we exited early
	if run.FailSilently {
		return 0, trrs, nil
	}

	if spec.ID == 0 {
		err = r.orm.InsertFinishedRunWithSpec(ctx, run, saveSuccessfulTaskRuns)
	} else {
		err = r.orm.InsertFinishedRun(ctx, run, saveSuccessfulTaskRuns)
	}
	if err != nil {
		return 0, trrs, pkgerrors.Wrapf(err, "error inserting finished results for spec ID %v", run.PipelineSpecID)
	}
	return run.ID, trrs, nil
}

func (r *runner) Run(ctx context.Context, run *Run, l logger.Logger, saveSuccessfulTaskRuns bool, fn func(tx sqlutil.DataSource) error) (incomplete bool, err error) {
	pipeline, err := r.InitializePipeline(run.PipelineSpec)
	if err != nil {
		return false, err
	}

	// retain old UUID values
	for _, taskRun := range run.PipelineTaskRuns {
		task := pipeline.ByDotID(taskRun.DotID)
		if task == nil || task.Base() == nil {
			return false, pkgerrors.Errorf("failed to match a pipeline task for dot ID: %v", taskRun.DotID)
		}
		task.Base().uuid = taskRun.ID
	}

	preinsert := pipeline.RequiresPreInsert()

	err = r.orm.Transact(ctx, func(tx ORM) error {
		// OPTIMISATION: avoid an extra db write if there is no async tasks present or if this is a resumed run
		if preinsert && run.ID == 0 {
			now := time.Now()
			// initialize certain task params
			for _, task := range pipeline.Tasks {
				switch task.Type() {
				case TaskTypeETHTx:
					run.PipelineTaskRuns = append(run.PipelineTaskRuns, TaskRun{
						ID:            task.Base().uuid,
						PipelineRunID: run.ID,
						Type:          task.Type(),
						Index:         task.OutputIndex(),
						DotID:         task.DotID(),
						CreatedAt:     now,
					})
				default:
				}
			}
			if err = tx.CreateRun(ctx, run); err != nil {
				return err
			}
		}

		if fn != nil {
			return fn(tx.DataSource())
		}
		return nil
	})
	if err != nil {
		return false, err
	}

	for {
		r.run(ctx, pipeline, run, NewVarsFrom(run.Inputs.Val.(map[string]interface{})), l)

		if preinsert {
			// FailSilently = run failed and task was marked failEarly. skip StoreRun and instead delete all trace of it
			if run.FailSilently {
				if err = r.orm.DeleteRun(ctx, run.ID); err != nil {
					return false, pkgerrors.Wrap(err, "Run")
				}
				return false, nil
			}

			var restart bool
			restart, err = r.orm.StoreRun(ctx, run)
			if err != nil {
				return false, pkgerrors.Wrapf(err, "error storing run for spec ID %v state %v outputs %v errors %v finished_at %v",
					run.PipelineSpec.ID, run.State, run.Outputs, run.FatalErrors, run.FinishedAt)
			}

			if restart {
				// instant restart: new data is already available in the database
				continue
			}
		} else {
			if run.Pending {
				return false, pkgerrors.Wrapf(err, "a run without async returned as pending")
			}
			// don't insert if we exited early
			if run.FailSilently {
				return false, nil
			}

			if err = r.orm.InsertFinishedRun(ctx, run, saveSuccessfulTaskRuns); err != nil {
				return false, pkgerrors.Wrapf(err, "error inserting finished run for spec ID %v", run.PipelineSpec.ID)
			}
		}

		r.runFinished(run)

		return run.Pending, err
	}
}

func (r *runner) ResumeRun(ctx context.Context, taskID uuid.UUID, value interface{}, err error) error {
	run, start, err := r.orm.UpdateTaskRunResult(ctx, taskID, Result{
		Value: value,
		Error: err,
	})
	if err != nil {
		return fmt.Errorf("failed to update task run result: %w", err)
	}

	// TODO: Should probably replace this with a listener to update events
	// which allows to pass in a transactionalised database to this function
	if start {
		// start the runner again
		go func() {
			if _, err := r.Run(context.Background(), &run, r.lggr, false, nil); err != nil {
				r.lggr.Errorw("Resume run failure", "err", err)
			}
			r.lggr.Debug("Resume run success")
		}()
	}
	return nil
}

func (r *runner) InsertFinishedRun(ctx context.Context, ds sqlutil.DataSource, run *Run, saveSuccessfulTaskRuns bool) error {
	orm := r.orm
	if ds != nil {
		orm = orm.WithDataSource(ds)
	}
	return orm.InsertFinishedRun(ctx, run, saveSuccessfulTaskRuns)
}

func (r *runner) InsertFinishedRuns(ctx context.Context, ds sqlutil.DataSource, runs []*Run, saveSuccessfulTaskRuns bool) error {
	orm := r.orm
	if ds != nil {
		orm = orm.WithDataSource(ds)
	}
	return orm.InsertFinishedRuns(ctx, runs, saveSuccessfulTaskRuns)
}

func (r *runner) runReaper() {
	r.lggr.Debugw("Pipeline run reaper starting")
	ctx, cancel := r.chStop.CtxCancel(context.WithTimeout(context.Background(), r.config.ReaperInterval()))
	defer cancel()

	err := r.orm.DeleteRunsOlderThan(ctx, r.config.ReaperThreshold())
	if err != nil {
		r.lggr.Errorw("Pipeline run reaper failed", "err", err)
		r.SvcErrBuffer.Append(err)
	} else {
		r.lggr.Debugw("Pipeline run reaper completed successfully")
	}
}

// init task: Searches the database for runs stuck in the 'running' state while the node was previously killed.
// We pick up those runs and resume execution.
func (r *runner) scheduleUnfinishedRuns() {
	defer r.wgDone.Done()

	// limit using a createdAt < now() @ start of run to prevent executing new jobs
	now := time.Now()

	if r.config.ReaperInterval() > time.Duration(0) {
		// immediately run reaper so we don't consider runs that are too old
		r.runReaper()
	}

	ctx, cancel := r.chStop.NewCtx()
	defer cancel()

	var wgRunsDone sync.WaitGroup
	err := r.orm.GetUnfinishedRuns(ctx, now, func(run Run) error {
		wgRunsDone.Add(1)

		go func() {
			defer wgRunsDone.Done()

			_, err := r.Run(ctx, &run, r.lggr, false, nil)
			if ctx.Err() != nil {
				return
			} else if err != nil {
				r.lggr.Errorw("Pipeline run init job resumption failed", "err", err)
			}
		}()

		return nil
	})

	wgRunsDone.Wait()

	if ctx.Err() != nil {
		return
	} else if err != nil {
		r.lggr.Errorw("Pipeline run init job failed", "err", err)
	}
}
