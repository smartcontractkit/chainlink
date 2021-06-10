package pipeline

import (
	"context"
	"fmt"
	"runtime/debug"
	"sort"
	"sync"
	"time"

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
	promPipelineTaskExecutionTime = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "pipeline_task_execution_time",
		Help: "How long each pipeline task took to execute",
	},
		[]string{"job_id", "job_name", "task_type"},
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
	task                  Task
	next                  *memoryTaskRun
	predecessorsRemaining int
	inputs                []input
	predMu                sync.RWMutex
	claimed               bool
	claimedMu             sync.Mutex
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

func (m *memoryTaskRun) numPredecessorsRemaining() int {
	m.predMu.RLock()
	defer m.predMu.RUnlock()
	return m.predecessorsRemaining
}

func (m *memoryTaskRun) tryToClaim() (ok bool) {
	m.claimedMu.Lock()
	defer m.claimedMu.Unlock()
	if m.claimed {
		return false
	}
	m.claimed = true
	return true
}

func (m *memoryTaskRun) predecessorFinished(result Result, index int32) {
	m.predMu.Lock()
	defer m.predMu.Unlock()
	m.inputs = append(m.inputs, input{result: result, index: index})
	m.predecessorsRemaining--
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
		} else {
			break
		}
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
	tasks, err := spec.TasksInDependencyOrder()
	if err != nil {
		return run, nil, err
	}
	f := time.Now()
	for _, task := range tasks {
		panickedTrrs = append(panickedTrrs, TaskRunResult{
			Task:       task,
			CreatedAt:  f,
			Result:     Result{Value: nil, Error: ErrRunPanicked},
			FinishedAt: time.Now(),
			IsTerminal: task.OutputTask() == nil,
		})
	}
	run.Outputs = TaskRunResults(panickedTrrs).FinalResult().OutputsDB()
	run.Errors = TaskRunResults(panickedTrrs).FinalResult().ErrorsDB()
	return run, panickedTrrs, nil
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

	var taskDAG TaskDAG
	err := taskDAG.UnmarshalText([]byte(spec.DotDagSource))
	if err != nil {
		return run, nil, false, err
	}

	headTaskRuns, err := r.memoryTaskRunDAGFromTaskDAG(taskDAG, pipelineInput, txdb)
	if err != nil {
		return run, nil, false, err
	}

	// TODO: Test with multiple and single null successor IDs
	// https://www.pivotaltracker.com/story/show/176557536

	// We execute tasks using "fan in" job processing because tasks have only one
	// output but potentially many inputs. Start one goroutine per "head task",
	// having it proceed down its chain of successors until:
	//  - it encounters a task with unfulfilled predecessors
	//  - it encounters a task that has already been claimed/finished
	//  - it processes the final task

	var (
		vars             = NewVarsFrom(map[string]interface{}{"input": pipelineInput})
		taskRunResults   TaskRunResults
		taskRunResultsMu sync.Mutex
		wg               sync.WaitGroup
		retry            bool
	)
	wg.Add(len(headTaskRuns))
	for _, mtr := range headTaskRuns {
		go func(taskRun *memoryTaskRun) {
			defer func() {
				if err := recover(); err != nil {
					logger.Default.Errorw("goroutine panicked executing run", "panic", err, "stacktrace", string(debug.Stack()))
					// No mutex needed: if any goroutine panics, we retry the run.
					retry = true
				}
				wg.Done()
			}()
			for taskRun != nil {
				// If this one is still waiting another chain, abandon this
				// goroutine and let the other handle it
				if taskRun.numPredecessorsRemaining() > 0 {
					return
				}
				// Avoid double execution, only one goroutine may finish the task
				if ok := taskRun.tryToClaim(); !ok {
					return
				}

				taskRunResult := r.executeTaskRun(ctx, spec, vars, taskRun, meta, l)

				if taskRunResult.Result.Error != nil {
					vars.Set(taskRunResult.Task.DotID(), taskRunResult.Result.Error)
				} else {
					vars.Set(taskRunResult.Task.DotID(), taskRunResult.Result.Value)
				}

				taskRunResultsMu.Lock()
				taskRunResults = append(taskRunResults, taskRunResult)
				taskRunResultsMu.Unlock()

				logTaskRunToPrometheus(taskRunResult, spec)

				if taskRun.next == nil {
					return
				}
				taskRun.next.predecessorFinished(taskRunResult.Result, taskRun.task.OutputIndex())
				taskRun = taskRun.next
			}
		}(mtr)
	}

	wg.Wait()
	finishRun := time.Now()
	runTime := finishRun.Sub(startRun)
	l.Debugw("Finished all tasks for pipeline run", "specID", spec.ID, "runTime", runTime)
	promPipelineRunTotalTimeToCompletion.WithLabelValues(fmt.Sprintf("%d", spec.JobID), spec.JobName).Set(float64(runTime))
	if retry || taskRunResults.FinalResult().HasErrors() {
		promPipelineRunErrors.WithLabelValues(fmt.Sprintf("%d", spec.JobID), spec.JobName).Inc()
	}
	if !retry {
		run.Errors = taskRunResults.FinalResult().ErrorsDB()
		run.Outputs = taskRunResults.FinalResult().OutputsDB()
		run.FinishedAt = &finishRun
	}
	return run, taskRunResults, retry, err
}

func (r *runner) memoryTaskRunDAGFromTaskDAG(taskDAG TaskDAG, pipelineInput interface{}, txdb *gorm.DB) (headTaskRuns []*memoryTaskRun, _ error) {
	tasks, err := taskDAG.TasksInDependencyOrder()
	if err != nil {
		return nil, err
	}

	all := make(map[string]*memoryTaskRun)
	txMu := &sync.Mutex{}

	for _, task := range tasks {
		if task.Type() == TaskTypeHTTP {
			task.(*HTTPTask).config = r.config
		} else if task.Type() == TaskTypeBridge {
			task.(*BridgeTask).config = r.config
			task.(*BridgeTask).safeTx = SafeTx{txdb, txMu}
		}

		mtr := memoryTaskRun{
			predecessorsRemaining: task.NumPredecessors(),
			task:                  task,
		}
		if mtr.numPredecessorsRemaining() == 0 {
			headTaskRuns = append(headTaskRuns, &mtr)
			mtr.inputs = append(mtr.inputs, input{index: 0, result: Result{Value: pipelineInput}})
		}
		all[task.DotID()] = &mtr
	}

	// Populate next pointers
	for dotID, mtr := range all {
		if mtr.task.OutputTask() != nil {
			all[dotID].next = all[mtr.task.OutputTask().DotID()]
		} else {
			all[dotID].next = nil
		}
	}

	return headTaskRuns, nil
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
		IsTerminal: taskRun.next == nil,
	}
}

func logTaskRunToPrometheus(trr TaskRunResult, spec Spec) {
	elapsed := trr.FinishedAt.Sub(trr.CreatedAt)

	promPipelineTaskExecutionTime.WithLabelValues(fmt.Sprintf("%d", spec.JobID), spec.JobName, string(trr.Task.Type())).Set(float64(elapsed))
	var status string
	if trr.Result.Error != nil {
		status = "error"
	} else {
		status = "completed"
	}
	promPipelineTasksTotalFinished.WithLabelValues(fmt.Sprintf("%d", spec.JobID), spec.JobName, string(trr.Task.Type()), status).Inc()
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

func (r *runner) runReaper() {
	err := r.orm.DeleteRunsOlderThan(r.config.JobPipelineReaperThreshold())
	if err != nil {
		logger.Errorw("Pipeline run reaper failed", "error", err)
	}
}
