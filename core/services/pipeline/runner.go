package pipeline

import (
	"context"
	"fmt"
	"runtime/debug"
	"sort"
	"sync"
	"time"

	"github.com/jpillora/backoff"
	"github.com/smartcontractkit/chainlink/core/store/models"

	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/postgres"
	"github.com/smartcontractkit/chainlink/core/utils"
	"gorm.io/gorm"
)

//go:generate mockery --name Runner --output ./mocks/ --case=underscore

// Runner checks the DB for incomplete TaskRuns and runs them.  For a
// TaskRun to be eligible to be run, its parent/input tasks must already
// all be complete.
type Runner interface {
	Start() error
	Close() error
	// We expect spec.JobID and spec.JobName to be set for logging/prometheus.
	ExecuteRun(ctx context.Context, spec Spec, meta JSONSerializable, l logger.Logger) (trrs TaskRunResults, err error)
	ExecuteAndInsertNewRun(ctx context.Context, spec Spec, meta JSONSerializable, l logger.Logger) (runID int64, finalResult FinalResult, err error)
	InsertFinishedRunWithResults(ctx context.Context, run Run, trrs TaskRunResults) (int64, error)

	// XXX: These methods are deprecated, you should be using the above
	CreateRun(ctx context.Context, jobID int32, meta map[string]interface{}) (runID int64, err error)
	AwaitRun(ctx context.Context, runID int64) error
	ResultsForRun(ctx context.Context, runID int64) ([]Result, error)
}

type runner struct {
	orm                             ORM
	config                          Config
	processIncompleteTaskRunsWorker utils.SleeperTask
	runReaperWorker                 utils.SleeperTask

	utils.StartStopOnce
	chStop  chan struct{}
	chDone  chan struct{}
	newRuns postgres.Subscription
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
	r.processIncompleteTaskRunsWorker = utils.NewSleeperTask(
		utils.SleeperTaskFuncWorker(r.processUnfinishedRuns),
	)
	r.runReaperWorker = utils.NewSleeperTask(
		utils.SleeperTaskFuncWorker(r.runReaper),
	)
	return r
}

func (r *runner) Start() error {
	if !r.OkayToStart() {
		return errors.New("Pipeline runner has already been started")
	}

	go r.runLoop()

	newRunsSubscription, err := r.orm.ListenForNewRuns()
	if err != nil {
		logger.Error("Pipeline runner could not subscribe to new run events, falling back to polling")
		return nil
	}
	r.newRuns = newRunsSubscription
	var newRunEvents = r.newRuns.Events()
	for i := 0; i < int(r.config.JobPipelineParallelism()); i++ {
		go func() {
			for {
				select {
				case <-newRunEvents:
					r.processUnfinishedRuns()
				case <-r.chStop:
					return
				}
			}
		}()
	}

	return nil
}

func (r *runner) Close() error {
	if !r.OkayToStop() {
		return errors.New("Pipeline runner has already been stopped")
	}

	close(r.chStop)
	<-r.chDone
	if r.newRuns != nil {
		r.newRuns.Close()
	}

	return nil
}

func (r *runner) destroy() {
	err := r.processIncompleteTaskRunsWorker.Stop()
	if err != nil {
		logger.Error(err)
	}
	err = r.runReaperWorker.Stop()
	if err != nil {
		logger.Error(err)
	}
}

func (r *runner) runLoop() {
	defer close(r.chDone)
	defer r.destroy()

	dbPollTicker := time.NewTicker(utils.WithJitter(r.config.TriggerFallbackDBPollInterval()))
	defer dbPollTicker.Stop()

	runReaperTicker := time.NewTicker(r.config.JobPipelineReaperInterval())
	defer runReaperTicker.Stop()

	for {
		select {
		case <-r.chStop:
			return
		case <-dbPollTicker.C:
			r.processIncompleteTaskRunsWorker.WakeUp()
		case <-runReaperTicker.C:
			r.runReaperWorker.WakeUp()
		}
	}
}

func (r *runner) CreateRun(ctx context.Context, jobID int32, meta map[string]interface{}) (int64, error) {
	runID, err := r.orm.CreateRun(ctx, jobID, meta)
	if err != nil {
		return 0, err
	}
	logger.Infow("Pipeline run created", "jobID", jobID, "runID", runID)
	return runID, nil
}

func (r *runner) AwaitRun(ctx context.Context, runID int64) error {
	ctx, cancel := utils.CombinedContext(r.chStop, ctx)
	defer cancel()
	return r.orm.AwaitRun(ctx, runID)
}

func (r *runner) ResultsForRun(ctx context.Context, runID int64) ([]Result, error) {
	ctx, cancel := utils.CombinedContext(r.chStop, ctx)
	defer cancel()
	return r.orm.ResultsForRun(ctx, runID)
}

// NOTE: This could potentially run on a different machine in the cluster than
// the one that originally added the job run.
func (r *runner) processUnfinishedRuns() {
	ctx, cancel := utils.CombinedContext(r.chStop, r.config.JobPipelineMaxRunDuration())
	defer cancel()

	_, err := r.orm.ProcessNextUnfinishedRun(ctx, r.executeRun)
	if err != nil {
		logger.Errorf("Error processing unfinished run: %v", err)
	}
}

type (
	memoryTaskRun struct {
		task          Task
		next          *memoryTaskRun
		nPredecessors int
		finished      bool
		inputs        []input
		predMu        sync.RWMutex
		finishMu      sync.Mutex
	}

	input struct {
		result Result
		index  int32
	}
)

// results returns the results sorted by index
// It is not thread-safe
func (m *memoryTaskRun) results() (a []Result) {
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

func (r *runner) ExecuteRun(ctx context.Context, spec Spec, meta JSONSerializable, l logger.Logger) (TaskRunResults, error) {
	var (
		trrs            TaskRunResults
		err             error
		retry           bool
		i               int
		numPanicRetries = 5
	)
	b := &backoff.Backoff{
		Min:    100 * time.Second,
		Max:    1 * time.Second,
		Factor: 2,
		Jitter: false,
	}
	for i = 0; i < numPanicRetries; i++ {
		trrs, retry, err = r.executeRun(ctx, r.orm.DB(), spec, meta, l)
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
	return trrs, err
}

// Generate a errored run from the spec.
func (r *runner) panickedRunResults(spec Spec) ([]TaskRunResult, error) {
	var panickedTrrs []TaskRunResult
	tasks, err := spec.TasksInDependencyOrder()
	if err != nil {
		return nil, err
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
	return panickedTrrs, nil
}

func (r *runner) executeRun(ctx context.Context, txdb *gorm.DB, spec Spec, meta JSONSerializable, l logger.Logger) (TaskRunResults, bool, error) {
	l.Debugw("Initiating tasks for pipeline run of spec", "job ID", spec.JobID, "job name", spec.JobName)
	var (
		err  error
		trrs TaskRunResults
	)
	startRun := time.Now()

	d := TaskDAG{}
	err = d.UnmarshalText([]byte(spec.DotDagSource))
	if err != nil {
		return trrs, false, err
	}

	// Find "firsts" and work forwards
	tasks, err := d.TasksInDependencyOrder()
	if err != nil {
		return nil, false, err
	}
	all := make(map[string]*memoryTaskRun)
	var graph []*memoryTaskRun
	txMu := new(sync.Mutex)
	for _, task := range tasks {
		if task.Type() == TaskTypeHTTP {
			task.(*HTTPTask).config = r.config
		}
		if task.Type() == TaskTypeBridge {
			task.(*BridgeTask).config = r.config
			task.(*BridgeTask).safeTx = SafeTx{txdb, txMu}
		}
		mtr := memoryTaskRun{
			nPredecessors: task.NPreds(),
			task:          task,
		}
		if mtr.nPredecessors == 0 {
			graph = append(graph, &mtr)
		}
		all[task.DotID()] = &mtr
	}

	// Populate next pointers
	for did, ts := range all {
		if ts.task.OutputTask() != nil {
			all[did].next = all[ts.task.OutputTask().DotID()]
		} else {
			all[did].next = nil
		}
	}

	// TODO: Test with multiple and single null successor IDs
	// https://www.pivotaltracker.com/story/show/176557536
	// 3. Execute tasks using "fan in" job processing
	var updateMu sync.Mutex
	var wg sync.WaitGroup
	var retry bool
	wg.Add(len(graph))
	for _, mtr := range graph {
		go func(m *memoryTaskRun) {
			defer func() {
				if err := recover(); err != nil {
					logger.Default.Errorw("goroutine panicked executing run", "panic", err, "stacktrace", string(debug.Stack()))
					// No mutex needed: if any goroutine panics, we retry the run.
					retry = true
				}
				wg.Done()
			}()
			for m != nil {
				m.predMu.RLock()
				nPredecessors := m.nPredecessors
				m.predMu.RUnlock()
				if nPredecessors > 0 {
					// This one is still waiting another chain, abandon this
					// goroutine and let the other handle it
					return
				}

				var finished bool

				// Avoid double execution, only one goroutine may finish the task
				m.finishMu.Lock()
				finished = m.finished
				if finished {
					m.finishMu.Unlock()
					return
				}
				m.finished = true
				m.finishMu.Unlock()

				startTaskRun := time.Now()

				result := r.executeTaskRun(ctx, spec, m.task, meta, m.results(), l)

				finishedAt := time.Now()

				trr := TaskRunResult{
					Task:       m.task,
					Result:     result,
					CreatedAt:  startTaskRun,
					FinishedAt: finishedAt,
					IsTerminal: m.next == nil,
				}

				updateMu.Lock()
				trrs = append(trrs, trr)
				updateMu.Unlock()

				elapsed := finishedAt.Sub(startTaskRun)

				promPipelineTaskExecutionTime.WithLabelValues(fmt.Sprintf("%d", spec.JobID), spec.JobName, string(m.task.Type())).Set(float64(elapsed))
				var status string
				if result.Error != nil {
					status = "error"
				} else {
					status = "completed"
				}
				promPipelineTasksTotalFinished.WithLabelValues(fmt.Sprintf("%d", spec.JobID), spec.JobName, string(m.task.Type()), status).Inc()
				if m.next == nil {
					return
				}

				m.next.predMu.Lock()
				m.next.inputs = append(m.next.inputs, input{result: result, index: m.task.OutputIndex()})
				m.next.nPredecessors--
				m.next.predMu.Unlock()

				m = m.next
			}
		}(mtr)
	}

	wg.Wait()

	runTime := time.Since(startRun)
	l.Debugw("Finished all tasks for pipeline run", "specID", spec.ID, "runTime", runTime)
	promPipelineRunTotalTimeToCompletion.WithLabelValues(fmt.Sprintf("%d", spec.JobID), spec.JobName).Set(float64(runTime))
	if retry || trrs.FinalResult().HasErrors() {
		promPipelineRunErrors.WithLabelValues(fmt.Sprintf("%d", spec.JobID), spec.JobName).Inc()
	}

	return trrs, retry, err
}

func (r *runner) executeTaskRun(ctx context.Context, spec Spec, task Task, meta JSONSerializable, inputs []Result, l logger.Logger) Result {
	loggerFields := []interface{}{
		"taskName", task.DotID(),
	}

	// Order of precedence for task timeout:
	// - Specific task timeout (task.TaskTimeout)
	// - Job level task timeout (spec.MaxTaskDuration)
	// - Passed in context
	taskTimeout, isSet := task.TaskTimeout()
	if isSet {
		var cancel context.CancelFunc
		ctx, cancel = utils.CombinedContext(r.chStop, taskTimeout)
		defer cancel()
	} else if spec.MaxTaskDuration != models.Interval(time.Duration(0)) {
		var cancel context.CancelFunc
		ctx, cancel = utils.CombinedContext(r.chStop, time.Duration(spec.MaxTaskDuration))
		defer cancel()
	}

	result := task.Run(ctx, meta, inputs)
	loggerFields = append(loggerFields, "result value", result.Value)
	loggerFields = append(loggerFields, "result error", result.Error)
	switch v := result.Value.(type) {
	case []byte:
		loggerFields = append(loggerFields, "resultString", fmt.Sprintf("%q", v))
		loggerFields = append(loggerFields, "resultHex", fmt.Sprintf("%x", v))
	}
	l.Debugw("Pipeline task completed", loggerFields...)

	return result
}

// ExecuteAndInsertNewRun bypasses the job pipeline entirely.
// It executes a run in memory then inserts the finished run/task run records, returning the final result
func (r *runner) ExecuteAndInsertNewRun(ctx context.Context, spec Spec, meta JSONSerializable, l logger.Logger) (runID int64, result FinalResult, err error) {
	var run Run
	run.PipelineSpecID = spec.ID
	run.CreatedAt = time.Now()
	trrs, err := r.ExecuteRun(ctx, spec, meta, l)
	if err != nil {
		return run.ID, result, errors.Wrapf(err, "error executing run for spec ID %v", spec.ID)
	}

	end := time.Now()
	run.FinishedAt = &end

	finalResult := trrs.FinalResult()
	run.Outputs = finalResult.OutputsDB()
	run.Errors = finalResult.ErrorsDB()

	if runID, err = r.orm.InsertFinishedRunWithResults(ctx, run, trrs); err != nil {
		return runID, result, errors.Wrapf(err, "error inserting finished results for spec ID %v", spec.ID)
	}

	return runID, finalResult, nil
}

func (r *runner) InsertFinishedRunWithResults(ctx context.Context, run Run, trrs TaskRunResults) (int64, error) {
	dbCtx, cancel := context.WithTimeout(ctx, r.config.DatabaseMaximumTxDuration())
	defer cancel()
	return r.orm.InsertFinishedRunWithResults(dbCtx, run, trrs)
}

func (r *runner) runReaper() {
	err := r.orm.DeleteRunsOlderThan(r.config.JobPipelineReaperThreshold())
	if err != nil {
		logger.Errorw("Pipeline run reaper failed", "error", err)
	}
}
