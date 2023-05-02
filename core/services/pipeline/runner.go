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

	"github.com/smartcontractkit/chainlink/v2/core/bridges"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/recovery"
	"github.com/smartcontractkit/chainlink/v2/core/services"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
	"github.com/smartcontractkit/chainlink/v2/core/store/models"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

//go:generate mockery --quiet --name Runner --output ./mocks/ --case=underscore

type Runner interface {
	services.ServiceCtx

	// Run is a blocking call that will execute the run until no further progress can be made.
	// If `incomplete` is true, the run is only partially complete and is suspended, awaiting to be resumed when more data comes in.
	// Note that `saveSuccessfulTaskRuns` value is ignored if the run contains async tasks.
	Run(ctx context.Context, run *Run, l logger.Logger, saveSuccessfulTaskRuns bool, fn func(tx pg.Queryer) error) (incomplete bool, err error)
	ResumeRun(taskID uuid.UUID, value interface{}, err error) error

	// ExecuteRun executes a new run in-memory according to a spec and returns the results.
	// We expect spec.JobID and spec.JobName to be set for logging/prometheus.
	ExecuteRun(ctx context.Context, spec Spec, vars Vars, l logger.Logger) (run Run, trrs TaskRunResults, err error)
	// InsertFinishedRun saves the run results in the database.
	InsertFinishedRun(run *Run, saveSuccessfulTaskRuns bool, qopts ...pg.QOpt) error
	InsertFinishedRuns(runs []*Run, saveSuccessfulTaskRuns bool, qopts ...pg.QOpt) error

	// ExecuteAndInsertFinishedRun executes a new run in-memory according to a spec, persists and saves the results.
	// It is a combination of ExecuteRun and InsertFinishedRun.
	// Note that the spec MUST have a DOT graph for this to work.
	ExecuteAndInsertFinishedRun(ctx context.Context, spec Spec, vars Vars, l logger.Logger, saveSuccessfulTaskRuns bool) (runID int64, finalResult FinalResult, err error)

	OnRunFinished(func(*Run))
}

type runner struct {
	orm                    ORM
	btORM                  bridges.ORM
	config                 Config
	chainSet               evm.ChainSet
	ethKeyStore            ETHKeyStore
	vrfKeyStore            VRFKeyStore
	runReaperWorker        utils.SleeperTask
	lggr                   logger.Logger
	httpClient             *http.Client
	unrestrictedHTTPClient *http.Client

	// test helper
	runFinished func(*Run)

	utils.StartStopOnce
	chStop utils.StopChan
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
		[]string{"job_id", "job_name", "task_id", "task_type", "status"},
	)
)

func NewRunner(orm ORM, btORM bridges.ORM, cfg Config, chainSet evm.ChainSet, ethks ETHKeyStore, vrfks VRFKeyStore, lggr logger.Logger, httpClient, unrestrictedHTTPClient *http.Client) *runner {
	r := &runner{
		orm:                    orm,
		btORM:                  btORM,
		config:                 cfg,
		chainSet:               chainSet,
		ethKeyStore:            ethks,
		vrfKeyStore:            vrfks,
		chStop:                 make(chan struct{}),
		wgDone:                 sync.WaitGroup{},
		runFinished:            func(*Run) {},
		lggr:                   lggr.Named("PipelineRunner"),
		httpClient:             httpClient,
		unrestrictedHTTPClient: unrestrictedHTTPClient,
	}
	r.runReaperWorker = utils.NewSleeperTask(
		utils.SleeperFuncTask(r.runReaper, "PipelineRunnerReaper"),
	)
	return r
}

// Start starts Runner.
func (r *runner) Start(context.Context) error {
	return r.StartOnce("PipelineRunner", func() error {
		r.wgDone.Add(1)
		go r.scheduleUnfinishedRuns()
		if r.config.JobPipelineReaperInterval() != time.Duration(0) {
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
	return map[string]error{r.Name(): r.StartStopOnce.Healthy()}
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
	if r.config.JobPipelineReaperInterval() == 0 {
		return
	}

	runReaperTicker := time.NewTicker(utils.WithJitter(r.config.JobPipelineReaperInterval()))
	defer runReaperTicker.Stop()
	for {
		select {
		case <-r.chStop:
			return
		case <-runReaperTicker.C:
			r.runReaperWorker.WakeUp()
			runReaperTicker.Reset(utils.WithJitter(r.config.JobPipelineReaperInterval()))
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

func NewRun(spec Spec, vars Vars) Run {
	return Run{
		State:          RunStatusRunning,
		PipelineSpec:   spec,
		PipelineSpecID: spec.ID,
		Inputs:         JSONSerializable{Val: vars.vars, Valid: true},
		Outputs:        JSONSerializable{Val: nil, Valid: false},
		CreatedAt:      time.Now(),
	}
}

func (r *runner) OnRunFinished(fn func(*Run)) {
	r.runFinished = fn
}

// Be careful with the ctx passed in here: it applies to requests in individual
// tasks but should _not_ apply to the scheduler or run itself
func (r *runner) ExecuteRun(
	ctx context.Context,
	spec Spec,
	vars Vars,
	l logger.Logger,
) (Run, TaskRunResults, error) {
	run := NewRun(spec, vars)

	pipeline, err := r.initializePipeline(&run)

	if err != nil {
		return run, nil, err
	}

	taskRunResults := r.run(ctx, pipeline, &run, vars, l)

	if run.Pending {
		return run, nil, pkgerrors.Wrapf(err, "unexpected async run for spec ID %v, tried executing via ExecuteAndInsertFinishedRun", spec.ID)
	}

	return run, taskRunResults, nil
}

func (r *runner) initializePipeline(run *Run) (*Pipeline, error) {
	pipeline, err := Parse(run.PipelineSpec.DotDagSource)
	if err != nil {
		return nil, err
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
			task.(*BridgeTask).orm = r.btORM
			task.(*BridgeTask).specId = run.PipelineSpec.ID
			// URL is "safe" because it comes from the node's own database. We
			// must use the unrestrictedHTTPClient because some node operators
			// may run external adapters on their own hardware
			task.(*BridgeTask).httpClient = r.unrestrictedHTTPClient
		case TaskTypeETHCall:
			task.(*ETHCallTask).chainSet = r.chainSet
			task.(*ETHCallTask).config = r.config
			task.(*ETHCallTask).specGasLimit = run.PipelineSpec.GasLimit
			task.(*ETHCallTask).jobType = run.PipelineSpec.JobType
		case TaskTypeVRF:
			task.(*VRFTask).keyStore = r.vrfKeyStore
		case TaskTypeVRFV2:
			task.(*VRFTaskV2).keyStore = r.vrfKeyStore
		case TaskTypeEstimateGasLimit:
			task.(*EstimateGasLimitTask).chainSet = r.chainSet
			task.(*EstimateGasLimitTask).specGasLimit = run.PipelineSpec.GasLimit
			task.(*EstimateGasLimitTask).jobType = run.PipelineSpec.JobType
		case TaskTypeETHTx:
			task.(*ETHTxTask).keyStore = r.ethKeyStore
			task.(*ETHTxTask).chainSet = r.chainSet
			task.(*ETHTxTask).specGasLimit = run.PipelineSpec.GasLimit
			task.(*ETHTxTask).jobType = run.PipelineSpec.JobType
			task.(*ETHTxTask).forwardingAllowed = run.PipelineSpec.ForwardingAllowed
		default:
		}
	}

	// retain old UUID values
	for _, taskRun := range run.PipelineTaskRuns {
		task := pipeline.ByDotID(taskRun.DotID)
		if task != nil && task.Base() != nil {
			task.Base().uuid = taskRun.ID
		} else {
			return nil, pkgerrors.Errorf("failed to match a pipeline task for dot ID: %v", taskRun.DotID)
		}
	}

	return pipeline, nil
}

func (r *runner) run(ctx context.Context, pipeline *Pipeline, run *Run, vars Vars, l logger.Logger) TaskRunResults {
	l = l.With("jobID", run.PipelineSpec.JobID, "jobName", run.PipelineSpec.JobName)
	l.Debug("Initiating tasks for pipeline run of spec")

	scheduler := newScheduler(pipeline, run, vars, l)
	go scheduler.Run()

	// This is "just in case" for cleaning up any stray reports.
	// Normally the scheduler loop doesn't stop until all in progress runs report back
	reportCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if pipelineTimeout := r.config.JobPipelineMaxRunDuration(); pipelineTimeout != 0 {
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

	if !scheduler.pending {
		run.FinishedAt = null.TimeFrom(time.Now())

		// NOTE: runTime can be very long now because it'll include suspend
		runTime := run.FinishedAt.Time.Sub(run.CreatedAt)
		l.Debugw("Finished all tasks for pipeline run", "specID", run.PipelineSpecID, "runTime", runTime)
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
		run.Outputs = JSONSerializable{Val: outputs, Valid: true}

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
	l.Debugw("Pipeline task completed", loggerFields...)

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
	PromPipelineTasksTotalFinished.WithLabelValues(fmt.Sprintf("%d", spec.JobID), spec.JobName, trr.Task.DotID(), string(trr.Task.Type()), status).Inc()
}

// ExecuteAndInsertFinishedRun executes a run in memory then inserts the finished run/task run records, returning the final result
func (r *runner) ExecuteAndInsertFinishedRun(ctx context.Context, spec Spec, vars Vars, l logger.Logger, saveSuccessfulTaskRuns bool) (runID int64, finalResult FinalResult, err error) {
	run, trrs, err := r.ExecuteRun(ctx, spec, vars, l)
	if err != nil {
		return 0, finalResult, pkgerrors.Wrapf(err, "error executing run for spec ID %v", spec.ID)
	}

	finalResult = trrs.FinalResult(l)

	// don't insert if we exited early
	if run.FailSilently {
		return 0, finalResult, nil
	}

	if err = r.orm.InsertFinishedRun(&run, saveSuccessfulTaskRuns); err != nil {
		return 0, finalResult, pkgerrors.Wrapf(err, "error inserting finished results for spec ID %v", spec.ID)
	}
	return run.ID, finalResult, nil

}

func (r *runner) Run(ctx context.Context, run *Run, l logger.Logger, saveSuccessfulTaskRuns bool, fn func(tx pg.Queryer) error) (incomplete bool, err error) {
	pipeline, err := r.initializePipeline(run)
	if err != nil {
		return false, err
	}

	preinsert := pipeline.RequiresPreInsert()

	q := r.orm.GetQ().WithOpts(pg.WithParentCtx(ctx))
	err = q.Transaction(func(tx pg.Queryer) error {
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
			if err = r.orm.CreateRun(run, pg.WithQueryer(tx)); err != nil {
				return err
			}
		}

		if fn != nil {
			return fn(tx)
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
				if err = r.orm.DeleteRun(run.ID); err != nil {
					return false, pkgerrors.Wrap(err, "Run")
				}
				return false, nil
			}

			var restart bool
			restart, err = r.orm.StoreRun(run)
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

			if err = r.orm.InsertFinishedRun(run, saveSuccessfulTaskRuns, pg.WithParentCtx(ctx)); err != nil {
				return false, pkgerrors.Wrapf(err, "error storing run for spec ID %v", run.PipelineSpec.ID)
			}
		}

		r.runFinished(run)

		return run.Pending, err
	}
}

func (r *runner) ResumeRun(taskID uuid.UUID, value interface{}, err error) error {
	run, start, err := r.orm.UpdateTaskRunResult(taskID, Result{
		Value: value,
		Error: err,
	})
	if err != nil {
		return err
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

func (r *runner) InsertFinishedRun(run *Run, saveSuccessfulTaskRuns bool, qopts ...pg.QOpt) error {
	return r.orm.InsertFinishedRun(run, saveSuccessfulTaskRuns, qopts...)
}

func (r *runner) InsertFinishedRuns(runs []*Run, saveSuccessfulTaskRuns bool, qopts ...pg.QOpt) error {
	return r.orm.InsertFinishedRuns(runs, saveSuccessfulTaskRuns, qopts...)
}

func (r *runner) runReaper() {
	r.lggr.Debugw("Pipeline run reaper starting")
	ctx, cancel := r.chStop.CtxCancel(context.WithTimeout(context.Background(), r.config.JobPipelineReaperInterval()))
	defer cancel()

	err := r.orm.DeleteRunsOlderThan(ctx, r.config.JobPipelineReaperThreshold())
	if err != nil {
		r.lggr.Errorw("Pipeline run reaper failed", "error", err)
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

	if r.config.JobPipelineReaperInterval() > time.Duration(0) {
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
				r.lggr.Errorw("Pipeline run init job resumption failed", "error", err)
			}
		}()

		return nil
	})

	wgRunsDone.Wait()

	if ctx.Err() != nil {
		return
	} else if err != nil {
		r.lggr.Errorw("Pipeline run init job failed", "error", err)
	}
}
