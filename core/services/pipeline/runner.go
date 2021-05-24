package pipeline

import (
	"context"
	"fmt"
	"runtime/debug"
	"sort"
	"time"

	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
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

	// We expect spec.JobID and spec.JobName to be set for logging/prometheus.
	// ExecuteRun executes a new run in-memory according to a spec and returns the results.
	ExecuteRun(ctx context.Context, spec Spec, vars Vars, meta JSONSerializable, l logger.Logger) (run Run, trrs TaskRunResults, err error)
	// InsertFinishedRun saves the run results in the database.
	InsertFinishedRun(db *gorm.DB, run Run, trrs TaskRunResults, saveSuccessfulTaskRuns bool) (int64, error)

	// ExecuteAndInsertNewRun executes a new run in-memory according to a spec, persists and saves the results.
	// It is a combination of ExecuteRun and InsertFinishedRun.
	// Note that the spec MUST have a DOT graph for this to work.
	ExecuteAndInsertFinishedRun(ctx context.Context, spec Spec, vars Vars, meta JSONSerializable, l logger.Logger, saveSuccessfulTaskRuns bool) (runID int64, finalResult FinalResult, err error)

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
	chDone chan struct{}
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
		chDone:    make(chan struct{}),
	}
	r.runReaperWorker = utils.NewSleeperTask(
		utils.SleeperTaskFuncWorker(r.runReaper),
	)
	return r
}

func (r *runner) Start() error {
	return r.StartOnce("PipelineRunner", func() error {
		go r.runReaperLoop()
		return nil
	})
}

func (r *runner) Close() error {
	return r.StopOnce("PipelineRunner", func() error {
		close(r.chStop)
		<-r.chDone
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
	defer close(r.chDone)
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
	inputs []input
	vars   Vars
}

// Returns the results sorted by index. It is not thread-safe.
func (m *memoryTaskRun) inputsSorted() (a []Result) {
	inputs := make([]input, len(m.inputs))
	copy(inputs, m.inputs)
	sort.Slice(inputs, func(i, j int) bool {
		return inputs[i].index < inputs[j].index
	})
	a = make([]Result, len(inputs))
	for i, input := range inputs {
		a[i] = input.result
	}

	return
}

type input struct {
	result Result
	index  int32
}

// When a task panics, we catch the panic and wrap it in an error for reporting to the scheduler.
type ErrRunPanicked struct {
	v interface{}
}

func (err ErrRunPanicked) Error() string {
	return fmt.Sprintf("goroutine panicked when executing run: %v", err.v)
}

func (r *runner) ExecuteRun(
	ctx context.Context,
	spec Spec,
	vars Vars,
	meta JSONSerializable,
	l logger.Logger,
) (Run, TaskRunResults, error) {
	l.Debugw("Initiating tasks for pipeline run of spec", "job ID", spec.JobID, "job name", spec.JobName)

	var (
		startRun = time.Now()
		run      = Run{
			PipelineSpecID: spec.ID,
			CreatedAt:      startRun,
		}
	)

	pipeline, err := Parse(spec.DotDagSource)
	if err != nil {
		return run, nil, err
	}

	// initialize certain task params
	for _, task := range pipeline.Tasks {
		switch task.Type() {
		case TaskTypeHTTP:
			task.(*HTTPTask).config = r.config
		case TaskTypeBridge:
			task.(*BridgeTask).config = r.config
			task.(*BridgeTask).db = r.orm.DB()
		case TaskTypeETHCall:
			task.(*ETHCallTask).ethClient = r.ethClient
		case TaskTypeETHTx:
			task.(*ETHTxTask).txManager = r.txManager
		default:
		}
	}

	todo := context.TODO()
	scheduler := newScheduler(todo, pipeline, vars)
	go scheduler.Run()

	for taskRun := range scheduler.taskCh {
		// execute
		go func(taskRun *memoryTaskRun) {
			defer func() {
				if err := recover(); err != nil {
					logger.Default.Errorw("goroutine panicked executing run", "panic", err, "stacktrace", string(debug.Stack()))

					t := time.Now()
					scheduler.report(todo, TaskRunResult{
						Task:       taskRun.task,
						Result:     Result{Error: ErrRunPanicked{err}},
						FinishedAt: t,
						CreatedAt:  t, // TODO: more accurate start time
					})
				}
			}()
			result := r.executeTaskRun(ctx, spec, taskRun, meta, l)

			logTaskRunToPrometheus(result, spec)

			scheduler.report(todo, result)
		}(taskRun)
	}

	finishRun := time.Now()
	runTime := finishRun.Sub(startRun)
	l.Debugw("Finished all tasks for pipeline run", "specID", spec.ID, "runTime", runTime)
	PromPipelineRunTotalTimeToCompletion.WithLabelValues(fmt.Sprintf("%d", spec.JobID), spec.JobName).Set(float64(runTime))

	var taskRunResults TaskRunResults
	for _, result := range scheduler.results {
		taskRunResults = append(taskRunResults, result)
	}

	finalResult := taskRunResults.FinalResult()
	if finalResult.HasErrors() {
		PromPipelineRunErrors.WithLabelValues(fmt.Sprintf("%d", spec.JobID), spec.JobName).Inc()
	}
	run.Errors = finalResult.ErrorsDB()
	run.Outputs = finalResult.OutputsDB()
	run.FinishedAt = &finishRun

	return run, taskRunResults, err
}

func (r *runner) executeTaskRun(ctx context.Context, spec Spec, taskRun *memoryTaskRun, meta JSONSerializable, l logger.Logger) TaskRunResult {
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

	result := taskRun.task.Run(ctx, taskRun.vars, meta, taskRun.inputsSorted())
	loggerFields = append(loggerFields, "result value", result.Value)
	loggerFields = append(loggerFields, "result error", result.Error)
	switch v := result.Value.(type) {
	case []byte:
		loggerFields = append(loggerFields, "resultString", fmt.Sprintf("%q", v))
		loggerFields = append(loggerFields, "resultHex", fmt.Sprintf("%x", v))
	}
	l.Debugw("Pipeline task completed", loggerFields...)

	return TaskRunResult{
		Task:       taskRun.task,
		Result:     result,
		CreatedAt:  start,
		FinishedAt: time.Now(),
	}
}

func logTaskRunToPrometheus(trr TaskRunResult, spec Spec) {
	elapsed := trr.FinishedAt.Sub(trr.CreatedAt)

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
func (r *runner) ExecuteAndInsertFinishedRun(ctx context.Context, spec Spec, vars Vars, meta JSONSerializable, l logger.Logger, saveSuccessfulTaskRuns bool) (runID int64, finalResult FinalResult, err error) {
	run, trrs, err := r.ExecuteRun(ctx, spec, vars, meta, l)
	if err != nil {
		return run.ID, finalResult, errors.Wrapf(err, "error executing run for spec ID %v", spec.ID)
	}
	finalResult = trrs.FinalResult()
	runID, err = r.orm.InsertFinishedRun(r.orm.DB(), run, trrs, saveSuccessfulTaskRuns)
	if err != nil {
		return runID, finalResult, errors.Wrapf(err, "error inserting finished results for spec ID %v", spec.ID)
	}
	return runID, finalResult, nil
}

func (r *runner) InsertFinishedRun(db *gorm.DB, run Run, trrs TaskRunResults, saveSuccessfulTaskRuns bool) (int64, error) {
	return r.orm.InsertFinishedRun(db, run, trrs, saveSuccessfulTaskRuns)
}

func (r *runner) TestInsertFinishedRun(db *gorm.DB, jobID int32, jobName string, jobType string, specID int32) (int64, error) {
	t := time.Now()
	runID, err := r.InsertFinishedRun(db, Run{
		PipelineSpecID: specID,
		Errors:         RunErrors{null.String{}},
		Outputs:        JSONSerializable{Val: "queued eth transaction"},
		CreatedAt:      t,
		FinishedAt:     &t,
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
