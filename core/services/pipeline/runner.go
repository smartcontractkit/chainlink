package pipeline

import (
	"context"
	"database/sql"
	"fmt"
	"runtime/debug"
	"sort"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	uuid "github.com/satori/go.uuid"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/utils"
	"gopkg.in/guregu/null.v4"
	"gorm.io/gorm"

	"github.com/smartcontractkit/chainlink/core/service"
	"github.com/smartcontractkit/chainlink/core/services/eth"
	"github.com/smartcontractkit/chainlink/core/store/models"
)

//go:generate mockery --name Runner --output ./mocks/ --case=underscore

type Runner interface {
	service.Service

	// Run is a blocking call that will execute the run until no further progress can be made.
	// If `incomplete` is true, the run is only partially complete and is suspended, awaiting to be resumed when more data comes in.
	// Note that `saveSuccessfulTaskRuns` value is ignored if the run contains async tasks.
	Run(ctx context.Context, run *Run, l logger.Logger, saveSuccessfulTaskRuns bool) (incomplete bool, err error)

	// We expect spec.JobID and spec.JobName to be set for logging/prometheus.
	// ExecuteRun executes a new run in-memory according to a spec and returns the results.
	ExecuteRun(ctx context.Context, spec Spec, vars Vars, l logger.Logger) (run Run, trrs TaskRunResults, err error)
	// InsertFinishedRun saves the run results in the database.
	InsertFinishedRun(db *gorm.DB, run Run, trrs TaskRunResults, saveSuccessfulTaskRuns bool) (int64, error)

	// ExecuteAndInsertNewRun executes a new run in-memory according to a spec, persists and saves the results.
	// It is a combination of ExecuteRun and InsertFinishedRun.
	// Note that the spec MUST have a DOT graph for this to work.
	ExecuteAndInsertFinishedRun(ctx context.Context, spec Spec, vars Vars, l logger.Logger, saveSuccessfulTaskRuns bool) (runID int64, finalResult FinalResult, err error)

	// Test method for inserting completed non-pipeline job runs
	TestInsertFinishedRun(db *gorm.DB, jobID int32, jobName string, jobType string, specID int32) (int64, error)
}

type runner struct {
	orm             ORM
	config          Config
	ethClient       eth.Client
	txManager       TxManager
	runReaperWorker utils.SleeperTask

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
		[]string{"job_id", "job_name", "task_type"},
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
		[]string{"job_id", "job_name", "task_type", "status"},
	)
)

func NewRunner(orm ORM, config Config, ethClient eth.Client, txManager TxManager) *runner {
	r := &runner{
		orm:       orm,
		config:    config,
		ethClient: ethClient,
		txManager: txManager,
		chStop:    make(chan struct{}),
		wgDone:    sync.WaitGroup{},
	}
	r.runReaperWorker = utils.NewSleeperTask(
		utils.SleeperTaskFuncWorker(r.runReaper),
	)
	return r
}

func (r *runner) Start() error {
	return r.StartOnce("PipelineRunner", func() error {
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
		logger.Error(err)
	}
}

func (r *runner) runReaperLoop() {
	r.wgDone.Add(1)
	defer r.wgDone.Done()
	defer r.destroy()

	runReaperTicker := time.NewTicker(r.config.JobPipelineReaperInterval())
	defer runReaperTicker.Stop()
	for {
		select {
		case <-r.chStop:
			return
		case <-runReaperTicker.C:
			r.runReaperWorker.WakeUp()
		}
	}
}

type memoryTaskRun struct {
	task   Task
	inputs []Result // sorted by input index
	vars   Vars
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
		Inputs:         JSONSerializable{Val: vars.vars, Null: false},
		Outputs:        JSONSerializable{Val: nil, Null: true},
		CreatedAt:      time.Now(),
	}
}

func (r *runner) ExecuteRun(
	ctx context.Context,
	spec Spec,
	vars Vars,
	l logger.Logger,
) (Run, TaskRunResults, error) {
	l.Debugw("Initiating tasks for pipeline run of spec", "job ID", spec.JobID, "job name", spec.JobName)

	run := NewRun(spec, vars)

	taskRunResults, err := r.run(ctx, &run, vars, l)
	if err != nil {
		return run, nil, err
	}

	if run.Async && run.Pending {
		return run, nil, errors.Wrapf(err, "unexpected async run for spec ID %v, tried executing via ExecuteAndInsertFinishedRun", spec.ID)
	}

	finalResult := taskRunResults.FinalResult()
	if finalResult.HasErrors() {
		PromPipelineRunErrors.WithLabelValues(fmt.Sprintf("%d", spec.JobID), spec.JobName).Inc()
	}

	return run, taskRunResults, nil
}

func (r *runner) run(
	ctx context.Context,
	run *Run,
	vars Vars,
	l logger.Logger,
) (TaskRunResults, error) {
	l.Debugw("Initiating tasks for pipeline run of spec", "job ID", run.PipelineSpec.JobID, "job name", run.PipelineSpec.JobName)

	pipeline, err := Parse(run.PipelineSpec.DotDagSource)
	if err != nil {
		return nil, err
	}

	// initialize certain task params
	for _, task := range pipeline.Tasks {
		switch task.Type() {
		case TaskTypeHTTP:
			task.(*HTTPTask).config = r.config
		case TaskTypeBridge:
			task.(*BridgeTask).config = r.config
			task.(*BridgeTask).db = r.orm.DB()
			task.(*BridgeTask).id = uuid.NewV4()
		case TaskTypeETHCall:
			task.(*ETHCallTask).ethClient = r.ethClient
		case TaskTypeETHTx:
			task.(*ETHTxTask).txManager = r.txManager
		default:
		}
	}

	// avoid an extra db write if there is no async tasks present or if this is a resumed run
	if pipeline.HasAsync() {
		run.Async = true
		if run.ID == 0 {
			if err = r.orm.CreateRun(r.orm.DB(), run); err != nil {
				return nil, err
			}
		}
	}

	todo := context.TODO()
	scheduler := newScheduler(todo, pipeline, run, vars)
	go scheduler.Run()

	for taskRun := range scheduler.taskCh {
		// execute
		go func(taskRun *memoryTaskRun) {
			defer func() {
				if err := recover(); err != nil {
					logger.Default.Errorw("goroutine panicked executing run", "panic", err, "stacktrace", string(debug.Stack()))

					t := time.Now()
					scheduler.report(todo, TaskRunResult{
						ID:         uuid.NewV4(),
						Task:       taskRun.task,
						Result:     Result{Error: ErrRunPanicked{err}},
						FinishedAt: null.TimeFrom(t),
						CreatedAt:  t, // TODO: more accurate start time
					})
				}
			}()
			result := r.executeTaskRun(ctx, run.PipelineSpec, taskRun, l)

			logTaskRunToPrometheus(result, run.PipelineSpec)

			scheduler.report(todo, result)
		}(taskRun)
	}

	// if the run is suspended, awaiting resumption
	run.Pending = scheduler.pending
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
			Output:        &output,
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
		var outputs []interface{}
		for _, result := range run.PipelineTaskRuns {
			// skip non-terminal results
			if len(result.task.Outputs()) != 0 {
				continue
			}
			errors = append(errors, result.Error)
			outputs = append(outputs, result.Output.Val)
		}
		run.Errors = errors
		run.Outputs = JSONSerializable{Val: outputs, Null: false}

		if run.HasErrors() {
			run.State = RunStatusErrored
		} else {
			run.State = RunStatusCompleted
		}
	}

	// TODO: drop this once we stop using TaskRunResults
	var taskRunResults TaskRunResults
	for _, result := range scheduler.results {
		taskRunResults = append(taskRunResults, result)
	}

	return taskRunResults, err
}

func (r *runner) executeTaskRun(ctx context.Context, spec Spec, taskRun *memoryTaskRun, l logger.Logger) TaskRunResult {
	start := time.Now()
	loggerFields := []interface{}{
		"taskName", taskRun.task.DotID(),
	}

	// Order of precedence for task timeout:
	// - Specific task timeout (task.TaskTimeout)
	// - Job level task timeout (spec.MaxTaskDuration)
	// - Passed in context
	taskTimeout, isSet := taskRun.task.TaskTimeout()
	if isSet {
		var cancel context.CancelFunc
		ctx, cancel = utils.CombinedContext(r.chStop, taskTimeout)
		defer cancel()
	} else if spec.MaxTaskDuration != models.Interval(time.Duration(0)) {
		var cancel context.CancelFunc
		ctx, cancel = utils.CombinedContext(r.chStop, time.Duration(spec.MaxTaskDuration))
		defer cancel()
	}

	result := taskRun.task.Run(ctx, taskRun.vars, taskRun.inputs)
	loggerFields = append(loggerFields, "result value", result.Value)
	loggerFields = append(loggerFields, "result error", result.Error)
	switch v := result.Value.(type) {
	case []byte:
		loggerFields = append(loggerFields, "resultString", fmt.Sprintf("%q", v))
		loggerFields = append(loggerFields, "resultHex", fmt.Sprintf("%x", v))
	}
	l.Debugw("Pipeline task completed", loggerFields...)

	now := time.Now()

	var id uuid.UUID
	if taskRun.task.Type() == TaskTypeBridge {
		id = taskRun.task.(*BridgeTask).id
	} else {
		id = uuid.NewV4()
	}

	return TaskRunResult{
		ID:         id,
		Task:       taskRun.task,
		Result:     result,
		CreatedAt:  start,
		FinishedAt: null.TimeFrom(now),
	}
}

func logTaskRunToPrometheus(trr TaskRunResult, spec Spec) {
	elapsed := trr.FinishedAt.Time.Sub(trr.CreatedAt)

	PromPipelineTaskExecutionTime.WithLabelValues(fmt.Sprintf("%d", spec.JobID), spec.JobName, string(trr.Task.Type())).Set(float64(elapsed))
	var status string
	if trr.Result.Error != nil {
		status = "error"
	} else {
		status = "completed"
	}
	PromPipelineTasksTotalFinished.WithLabelValues(fmt.Sprintf("%d", spec.JobID), spec.JobName, string(trr.Task.Type()), status).Inc()
}

// ExecuteAndInsertNewRun executes a run in memory then inserts the finished run/task run records, returning the final result
func (r *runner) ExecuteAndInsertFinishedRun(ctx context.Context, spec Spec, vars Vars, l logger.Logger, saveSuccessfulTaskRuns bool) (runID int64, finalResult FinalResult, err error) {
	run, trrs, err := r.ExecuteRun(ctx, spec, vars, l)
	if err != nil {
		return 0, finalResult, errors.Wrapf(err, "error executing run for spec ID %v", spec.ID)
	}

	finalResult = trrs.FinalResult()
	if runID, err = r.orm.InsertFinishedRun(r.orm.DB(), run, trrs, saveSuccessfulTaskRuns); err != nil {
		return runID, finalResult, errors.Wrapf(err, "error inserting finished results for spec ID %v", spec.ID)
	}
	return runID, finalResult, nil

}

func (r *runner) Run(ctx context.Context, run *Run, l logger.Logger, saveSuccessfulTaskRuns bool) (incomplete bool, err error) {
	for {
		trrs, err := r.run(ctx, run, NewVarsFrom(run.Inputs.Val.(map[string]interface{})), l)
		if err != nil {
			return false, errors.Wrapf(err, "failed to run for spec ID %v", run.PipelineSpec.ID)
		}

		if run.Async {
			var db *sql.DB
			db, err = r.orm.DB().DB()
			if err != nil {
				return false, errors.Wrap(err, "unable to retrieve sql.DB")
			}

			var restart bool

			restart, err = r.orm.StoreRun(db, run)
			if err != nil {
				return false, errors.Wrapf(err, "error storing run for spec ID %v", run.PipelineSpec.ID)
			}

			if restart {
				// instant restart: new data is already available in the database
				continue
			}
		} else {
			if run.Pending {
				return false, errors.Wrapf(err, "a run without async returned as pending")
			}
			if _, err = r.orm.InsertFinishedRun(r.orm.DB(), *run, trrs, saveSuccessfulTaskRuns); err != nil {
				return false, errors.Wrapf(err, "error storing run for spec ID %v", run.PipelineSpec.ID)
			}
		}

		return run.Pending, err
	}
}

func (r *runner) InsertFinishedRun(db *gorm.DB, run Run, trrs TaskRunResults, saveSuccessfulTaskRuns bool) (int64, error) {
	return r.orm.InsertFinishedRun(db, run, trrs, saveSuccessfulTaskRuns)
}

func (r *runner) TestInsertFinishedRun(db *gorm.DB, jobID int32, jobName string, jobType string, specID int32) (int64, error) {
	t := time.Now()
	runID, err := r.InsertFinishedRun(db, Run{
		State:          RunStatusCompleted,
		PipelineSpecID: specID,
		Errors:         RunErrors{null.String{}},
		Outputs:        JSONSerializable{Val: "queued eth transaction"},
		CreatedAt:      t,
		FinishedAt:     null.TimeFrom(t),
	}, nil, false)
	elapsed := time.Since(t)

	// For testing metrics.
	id := fmt.Sprintf("%d", jobID)
	PromPipelineTaskExecutionTime.WithLabelValues(id, jobName, jobType).Set(float64(elapsed))
	var status string
	if err != nil {
		status = "error"
		PromPipelineRunErrors.WithLabelValues(id, jobName).Inc()
	} else {
		status = "completed"
	}
	PromPipelineRunTotalTimeToCompletion.WithLabelValues(id, jobName).Set(float64(elapsed))
	PromPipelineTasksTotalFinished.WithLabelValues(id, jobName, jobType, status).Inc()
	return runID, err
}

func (r *runner) runReaper() {
	err := r.orm.DeleteRunsOlderThan(r.config.JobPipelineReaperThreshold())
	if err != nil {
		logger.Errorw("Pipeline run reaper failed", "error", err)
	}
}

// init task: Searches the database for runs stuck in the 'running' state while the node was previously killed.
// We pick up those runs and resume execution.
func (r *runner) scheduleUnfinishedRuns() {
	r.wgDone.Add(1)
	defer r.wgDone.Done()

	// limit using a createdAt < now() @ start of run to prevent executing new jobs
	now := time.Now()

	// immediately run reaper so we don't consider runs that are too old
	r.runReaper()

	err := r.orm.GetUnfinishedRuns(now, func(run Run) error {
		go func() {
			if _, err := r.Run(context.TODO(), &run, *logger.Default, false); err != nil {
				logger.Errorw("Pipeline run init job resumption failed", "error", err)
			}
		}()
		return nil
	})
	if err != nil {
		logger.Errorw("Pipeline run init job failed", "error", err)
	}
}
