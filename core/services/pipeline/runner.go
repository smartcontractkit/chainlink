package pipeline

import (
	"context"
	"fmt"
	"net/http"
	"sort"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	uuid "github.com/satori/go.uuid"
	"gopkg.in/guregu/null.v4"

	"github.com/smartcontractkit/chainlink/core/chains/evm"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/recovery"
	"github.com/smartcontractkit/chainlink/core/services"
	"github.com/smartcontractkit/chainlink/core/services/pg"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"
)

//go:generate mockery --name Runner --output ./mocks/ --case=underscore

type Runner interface {
	services.ServiceCtx

	// Run is a blocking call that will execute the run until no further progress can be made.
	// If `incomplete` is true, the run is only partially complete and is suspended, awaiting to be resumed when more data comes in.
	// Note that `saveSuccessfulTaskRuns` value is ignored if the run contains async tasks.
	Run(ctx context.Context, run *Run, l logger.Logger, saveSuccessfulTaskRuns bool, fn func(tx pg.Queryer) error) (incomplete bool, err error)
	ResumeRun(taskID uuid.UUID, value interface{}, err error) error

	// We expect spec.JobID and spec.JobName to be set for logging/prometheus.
	// ExecuteRun executes a new run in-memory according to a spec and returns the results.
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
	chStop chan struct{}
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

func NewRunner(orm ORM, config Config, chainSet evm.ChainSet, ethks ETHKeyStore, vrfks VRFKeyStore, lggr logger.Logger, httpClient, unrestrictedHTTPClient *http.Client) *runner {
	r := &runner{
		orm:                    orm,
		config:                 config,
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
		r.wgDone.Add(2)
		go r.scheduleUnfinishedRuns()
		go r.runReaperLoop()
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

	taskRunResults, err := r.run(ctx, pipeline, &run, vars, l)
	if err != nil {
		return run, nil, err
	}

	if run.Pending {
		return run, nil, errors.Wrapf(err, "unexpected async run for spec ID %v, tried executing via ExecuteAndInsertFinishedRun", spec.ID)
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
		task.Base().uuid = uuid.NewV4()

		switch task.Type() {
		case TaskTypeHTTP:
			task.(*HTTPTask).config = r.config
			task.(*HTTPTask).httpClient = r.httpClient
			task.(*HTTPTask).unrestrictedHTTPClient = r.unrestrictedHTTPClient
		case TaskTypeBridge:
			task.(*BridgeTask).config = r.config
			task.(*BridgeTask).queryer = r.orm.GetQ()
			// URL is "safe" because it comes from the node's own database. We
			// must use the unrestrictedHTTPClient because some node operators
			// may run external adapters on their own hardware
			task.(*BridgeTask).httpClient = r.unrestrictedHTTPClient
		case TaskTypeETHCall:
			task.(*ETHCallTask).chainSet = r.chainSet
			task.(*ETHCallTask).config = r.config
		case TaskTypeVRF:
			task.(*VRFTask).keyStore = r.vrfKeyStore
		case TaskTypeVRFV2:
			task.(*VRFTaskV2).keyStore = r.vrfKeyStore
		case TaskTypeEstimateGasLimit:
			task.(*EstimateGasLimitTask).chainSet = r.chainSet
		case TaskTypeETHTx:
			task.(*ETHTxTask).keyStore = r.ethKeyStore
			task.(*ETHTxTask).chainSet = r.chainSet
		default:
		}
	}

	// retain old UUID values
	for _, taskRun := range run.PipelineTaskRuns {
		task := pipeline.ByDotID(taskRun.DotID)
		task.Base().uuid = taskRun.ID
	}

	return pipeline, nil
}

func (r *runner) run(
	ctx context.Context,
	pipeline *Pipeline,
	run *Run,
	vars Vars,
	l logger.Logger,
) (TaskRunResults, error) {
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
				ID:         uuid.NewV4(),
				Task:       taskRun.task,
				Result:     Result{Error: ErrRunPanicked{err}},
				FinishedAt: null.TimeFrom(t),
				CreatedAt:  t, // TODO: more accurate start time
			})
		})
	}

	// if the run is suspended, awaiting resumption
	run.Pending = scheduler.pending
	run.FailEarly = scheduler.exiting
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

	return taskRunResults, nil
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
	ctx, cancel := utils.WithCloseChan(ctx, r.chStop)
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
		return 0, finalResult, errors.Wrapf(err, "error executing run for spec ID %v", spec.ID)
	}

	finalResult = trrs.FinalResult(l)

	// don't insert if we exited early
	if run.FailEarly {
		return 0, finalResult, nil
	}

	if err = r.orm.InsertFinishedRun(&run, saveSuccessfulTaskRuns); err != nil {
		return 0, finalResult, errors.Wrapf(err, "error inserting finished results for spec ID %v", spec.ID)
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
		if _, err = r.run(ctx, pipeline, run, NewVarsFrom(run.Inputs.Val.(map[string]interface{})), l); err != nil {
			return false, errors.Wrapf(err, "failed to run for spec ID %v", run.PipelineSpec.ID)
		}

		if preinsert {
			// if run failed and it's failEarly, skip StoreRun and instead delete all trace of it
			if run.FailEarly {
				if err = r.orm.DeleteRun(run.ID); err != nil {
					return false, errors.Wrap(err, "Run")
				}
				return false, nil
			}

			var restart bool
			restart, err = r.orm.StoreRun(run)
			if err != nil {
				return false, errors.Wrapf(err, "error storing run for spec ID %v state %v outputs %v errors %v finished_at %v",
					run.PipelineSpec.ID, run.State, run.Outputs, run.FatalErrors, run.FinishedAt)
			}

			if restart {
				// instant restart: new data is already available in the database
				continue
			}
		} else {
			if run.Pending {
				return false, errors.Wrapf(err, "a run without async returned as pending")
			}
			// don't insert if we exited early
			if run.FailEarly {
				return false, nil
			}

			if err = r.orm.InsertFinishedRun(run, saveSuccessfulTaskRuns, pg.WithParentCtx(ctx)); err != nil {
				return false, errors.Wrapf(err, "error storing run for spec ID %v", run.PipelineSpec.ID)
			}
		}

		r.runFinished(run)

		return run.Pending, err
	}
}

func (r *runner) ResumeRun(taskID uuid.UUID, value interface{}, err error) error {
	result := Result{
		Value: value,
		Error: err,
	}
	run, start, err := r.orm.UpdateTaskRunResult(taskID, result)
	if err != nil {
		return err
	}

	// TODO: Should probably replace this with a listener to update events
	// which allows to pass in a transactionalised database to this function
	if start {
		// start the runner again
		go func() {
			if _, err := r.Run(context.Background(), &run, r.lggr, false, nil); err != nil {
				r.lggr.Errorw("Resume", "err", err)
			}
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
	ctx, cancel := utils.ContextFromChan(r.chStop)
	defer cancel()

	err := r.orm.DeleteRunsOlderThan(ctx, r.config.JobPipelineReaperThreshold())
	if err != nil {
		r.lggr.Errorw("Pipeline run reaper failed", "error", err)
	}
}

// init task: Searches the database for runs stuck in the 'running' state while the node was previously killed.
// We pick up those runs and resume execution.
func (r *runner) scheduleUnfinishedRuns() {
	defer r.wgDone.Done()

	// limit using a createdAt < now() @ start of run to prevent executing new jobs
	now := time.Now()

	// immediately run reaper so we don't consider runs that are too old
	r.runReaper()

	ctx, cancel := utils.ContextFromChan(r.chStop)
	defer cancel()

	err := r.orm.GetUnfinishedRuns(ctx, now, func(run Run) error {
		go func() {
			_, err := r.Run(ctx, &run, r.lggr, false, nil)
			if ctx.Err() != nil {
				return
			} else if err != nil {
				r.lggr.Errorw("Pipeline run init job resumption failed", "error", err)
			}
		}()
		return nil
	})
	if ctx.Err() != nil {
		return
	} else if err != nil {
		r.lggr.Errorw("Pipeline run init job failed", "error", err)
	}
}
