package pipeline

import (
	"context"
	"fmt"
	"runtime/debug"
	"sort"
	"sync"
	"time"

	"gopkg.in/guregu/null.v4"

	"github.com/jpillora/backoff"
	"github.com/smartcontractkit/chainlink/core/service"
	"github.com/smartcontractkit/chainlink/core/store/models"

	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/utils"
	"gorm.io/gorm"
)

//go:generate mockery --name Runner --output ./mocks/ --case=underscore

type Runner interface {
	service.Service

	// We expect spec.JobID and spec.JobName to be set for logging/prometheus.
	// ExecuteRun executes a new run in-memory according to a spec and returns the results.
	ExecuteRun(ctx context.Context, spec Spec, pipelineInput interface{}, meta JSONSerializable, l logger.Logger) (run Run, trrs TaskRunResults, err error)
	// InsertFinishedRun saves the run results in the database.
	InsertFinishedRun(db *gorm.DB, run Run, trrs TaskRunResults, saveSuccessfulTaskRuns bool) (int64, error)

	// ExecuteAndInsertNewRun executes a new run in-memory according to a spec, persists and saves the results.
	// It is a combination of ExecuteRun and InsertFinishedRun.
	// Note that the spec MUST have a DOT graph for this to work.
	ExecuteAndInsertFinishedRun(ctx context.Context, spec Spec, pipelineInput interface{}, meta JSONSerializable, l logger.Logger, saveSuccessfulTaskRuns bool) (runID int64, finalResult FinalResult, err error)

	// Test method for inserting completed non-pipeline job runs
	TestInsertFinishedRun(db *gorm.DB, jobID int32, jobName string, jobType string, specID int32) (int64, error)
}

type runner struct {
	orm             ORM
	config          Config
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
	ErrRunPanicked = errors.New("pipeline run panicked")
)

func NewRunner(orm ORM, config Config) *runner {
	r := &runner{
		orm:    orm,
		config: config,
		chStop: make(chan struct{}),
		chDone: make(chan struct{}),
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

func (r *runner) ExecuteRun(ctx context.Context, spec Spec, pipelineInput interface{}, meta JSONSerializable, l logger.Logger) (Run, TaskRunResults, error) {
	var (
		trrs            TaskRunResults
		err             error
		retry           bool
		i               int
		numPanicRetries = 5
		run             Run
	)
	b := &backoff.Backoff{
		Min:    100 * time.Second,
		Max:    1 * time.Second,
		Factor: 2,
		Jitter: false,
	}
	for i = 0; i < numPanicRetries; i++ {
		run, trrs, retry, err = r.executeRun(ctx, r.orm.DB(), spec, pipelineInput, meta, l)
		if retry {
			time.Sleep(b.Duration())
			continue
		}
		break
	}
	if i == numPanicRetries {
		return r.panickedRunResults(spec)
	}
	return run, trrs, err
}

// Generate a errored run from the spec.
func (r *runner) panickedRunResults(spec Spec) (Run, []TaskRunResult, error) {
	var panickedTrrs []TaskRunResult
	var run Run
	run.PipelineSpecID = spec.ID
	run.CreatedAt = time.Now()
	run.FinishedAt = &run.CreatedAt
	p, err := spec.Pipeline()
	if err != nil {
		return run, nil, err
	}
	f := time.Now()
	for _, task := range p.Tasks {
		panickedTrrs = append(panickedTrrs, TaskRunResult{
			Task:       task,
			CreatedAt:  f,
			Result:     Result{Value: nil, Error: ErrRunPanicked},
			FinishedAt: time.Now(),
		})
	}
	run.Outputs = TaskRunResults(panickedTrrs).FinalResult().OutputsDB()
	run.Errors = TaskRunResults(panickedTrrs).FinalResult().ErrorsDB()
	return run, panickedTrrs, nil
}

type scheduler struct {
	pipeline     *Pipeline
	dependencies map[int64]uint
	input        interface{}
	waiting      uint
	// roots are the tasks at the start of the pipeline
	roots   []Task
	results map[int64]TaskRunResult

	taskCh   chan *memoryTaskRun
	resultCh chan TaskRunResult
}

func newScheduler(p *Pipeline, i interface{}) *scheduler {
	dependencies := make(map[int64]uint, len(p.Tasks))
	var roots []Task

	for id, task := range p.Tasks {
		i := len(task.Inputs())
		dependencies[id] = uint(i)

		// no inputs: this is a root
		if i == 0 {
			roots = append(roots, task)
		}
	}
	s := &scheduler{
		pipeline:     p,
		dependencies: dependencies,
		input:        i,
		results:      make(map[int64]TaskRunResult, len(p.Tasks)),

		// taskCh should never block
		taskCh:   make(chan *memoryTaskRun, len(dependencies)),
		resultCh: make(chan TaskRunResult),
	}

	for _, task := range roots {
		run := &memoryTaskRun{task: task}
		// fill in the inputs
		run.inputs = append(run.inputs, input{index: 0, result: Result{Value: s.input}})

		s.taskCh <- run
		s.waiting++
	}

	go s.run()

	return s
}

func (s *scheduler) run() {
	for result := range s.resultCh {
		s.waiting--

		// mark job as complete
		s.results[result.Task.ID()] = result

		for _, output := range result.Task.Outputs() {
			id := output.ID()
			s.dependencies[id]--

			// if all dependencies are done, schedule task run
			if s.dependencies[id] == 0 {
				task := s.pipeline.Tasks[id]
				run := &memoryTaskRun{task: task}

				// fill in the inputs
				for _, i := range task.Inputs() {
					run.inputs = append(run.inputs, input{index: int32(i.OutputIndex()), result: s.results[i.ID()].Result})
				}

				s.taskCh <- run
				s.waiting++
			}
		}

		// if we are done, stop execution
		if s.waiting == 0 {
			close(s.taskCh)
			break
		}
	}
}

// When a task panics, we catch the panic and wrap it in an error for reporting to the scheduler.
type panicError struct {
	v interface{}
}

func (err panicError) Error() string {
	return fmt.Sprintf("goroutine panicked when executing run: %v", err.v)
}

func (r *runner) executeRun(
	ctx context.Context,
	txdb *gorm.DB,
	spec Spec,
	pipelineInput interface{},
	meta JSONSerializable,
	l logger.Logger,
) (Run, TaskRunResults, bool, error) {
	l.Debugw("Initiating tasks for pipeline run of spec", "job ID", spec.JobID, "job name", spec.JobName)

	var (
		startRun = time.Now()
		run      = Run{
			PipelineSpecID: spec.ID,
			CreatedAt:      startRun,
		}
	)

	pipeline, err := Parse([]byte(spec.DotDagSource))
	if err != nil {
		return run, nil, false, err
	}

	// initialize certain task params
	txMu := &sync.Mutex{}
	for _, task := range pipeline.Tasks {
		if task.Type() == TaskTypeHTTP {
			task.(*HTTPTask).config = r.config
		} else if task.Type() == TaskTypeBridge {
			task.(*BridgeTask).config = r.config
			task.(*BridgeTask).safeTx = SafeTx{txdb, txMu}
		}
	}

	scheduler := newScheduler(pipeline, pipelineInput)

	// TODO: Test with multiple and single null successor IDs
	// https://www.pivotaltracker.com/story/show/176557536

	var (
		vars  = NewVarsFrom(map[string]interface{}{"input": pipelineInput})
		retry bool
	)

	for taskRun := range scheduler.taskCh {
		// execute
		go func(taskRun *memoryTaskRun) {
			defer func() {
				if err := recover(); err != nil {
					logger.Default.Errorw("goroutine panicked executing run", "panic", err, "stacktrace", string(debug.Stack()))

					// No mutex needed: if any goroutine panics, we retry the run.
					retry = true

					scheduler.resultCh <- TaskRunResult{
						Task:       taskRun.task,
						Result:     Result{Error: panicError{err}},
						FinishedAt: time.Now(),
						// TODO: CreatedAt
					}
				}
			}()
			result := r.executeTaskRun(ctx, spec, vars, taskRun, meta, l)

			// TODO: remove vars locks by constructing Vars instances per task?
			if result.Result.Error != nil {
				vars.Set(result.Task.DotID(), result.Result.Error)
			} else {
				vars.Set(result.Task.DotID(), result.Result.Value)
			}

			logTaskRunToPrometheus(result, spec)

			// report the result back to the scheduler
			scheduler.resultCh <- result
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

	var errors bool
	if retry {
		errors = true
	} else {
		finalResult := taskRunResults.FinalResult()
		if finalResult.HasErrors() {
			errors = true
		}
		run.Errors = finalResult.ErrorsDB()
		run.Outputs = finalResult.OutputsDB()
		run.FinishedAt = &finishRun
	}

	if errors {
		PromPipelineRunErrors.WithLabelValues(fmt.Sprintf("%d", spec.JobID), spec.JobName).Inc()
	}

	return run, taskRunResults, retry, err
}

func (r *runner) executeTaskRun(ctx context.Context, spec Spec, vars Vars, taskRun *memoryTaskRun, meta JSONSerializable, l logger.Logger) TaskRunResult {
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

	result := taskRun.task.Run(ctx, vars, meta, taskRun.inputsSorted())
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
func (r *runner) ExecuteAndInsertFinishedRun(ctx context.Context, spec Spec, pipelineInput interface{}, meta JSONSerializable, l logger.Logger, saveSuccessfulTaskRuns bool) (runID int64, finalResult FinalResult, err error) {
	run, trrs, err := r.ExecuteRun(ctx, spec, pipelineInput, meta, l)
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
