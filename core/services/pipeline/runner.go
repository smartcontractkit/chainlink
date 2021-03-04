package pipeline

import (
	"context"
	"fmt"
	"sort"
	"sync"
	"time"

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
	CreateRun(ctx context.Context, jobID int32, meta map[string]interface{}) (runID int64, err error)
	ExecuteRun(ctx context.Context, run Run, l logger.Logger) (trrs TaskRunResults, err error)
	ExecuteAndInsertNewRun(ctx context.Context, spec Spec, l logger.Logger) (runID int64, finalResult FinalResult, err error)
	AwaitRun(ctx context.Context, runID int64) error
	ResultsForRun(ctx context.Context, runID int64) ([]Result, error)
	InsertFinishedRunWithResults(ctx context.Context, run Run, trrs TaskRunResults) (int64, error)
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
		[]string{"pipeline_spec_id", "task_type"},
	)
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
					_, err = r.processRun()
					if err != nil {
						logger.Errorf("Error processing incomplete task runs: %v", err)
					}
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

// NewRun creates an in-memory Run along with its TaskRuns for the provided job
// It does not interact with the database
func NewRun(spec Spec, startedAt time.Time) (run Run, err error) {
	if len(spec.PipelineTaskSpecs) == 0 {
		return run, errors.New("spec.PipelineTaskSpecs was empty")
	}
	taskRuns := make([]TaskRun, len(spec.PipelineTaskSpecs))
	for i, taskSpec := range spec.PipelineTaskSpecs {
		taskRuns[i] = TaskRun{
			Type:               taskSpec.Type,
			PipelineTaskSpecID: taskSpec.ID,
			PipelineTaskSpec:   taskSpec,
			CreatedAt:          startedAt,
		}
	}
	run = Run{
		PipelineSpecID:   spec.ID,
		PipelineSpec:     spec,
		PipelineTaskRuns: taskRuns,
		CreatedAt:        startedAt,
	}

	return run, nil
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
	_, err := r.processRun()
	if err != nil {
		logger.Errorf("Error processing unfinished run: %v", err)
	}
}

func (r *runner) processRun() (anyRemaining bool, err error) {
	ctx, cancel := utils.CombinedContext(r.chStop, r.config.JobPipelineMaxRunDuration())
	defer cancel()

	return r.orm.ProcessNextUnfinishedRun(ctx, r.executeRun)
}

type (
	memoryTaskRun struct {
		taskRun       TaskRun
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

func (r *runner) ExecuteRun(ctx context.Context, run Run, l logger.Logger) (trrs TaskRunResults, err error) {
	return r.executeRun(ctx, r.orm.DB(), run, l)
}

func (r *runner) executeRun(ctx context.Context, txdb *gorm.DB, run Run, l logger.Logger) (trrs TaskRunResults, err error) {
	l.Debugw("Initiating tasks for pipeline run", "runID", run.ID)
	startRun := time.Now()

	if run.PipelineSpec.ID == 0 {
		return nil, errors.Errorf("run.PipelineSpec.ID may not be 0: %#v", run)
	}

	// Find "firsts" and work forwards
	// 1. Make map of all memory task runs keyed by task spec id
	all := make(map[int32]*memoryTaskRun)
	for _, tr := range run.PipelineTaskRuns {
		if tr.PipelineTaskSpec.ID == 0 {
			return nil, errors.Errorf("taskRun.PipelineTaskSpec.ID may not be 0: %#v", tr)
		}
		mtr := memoryTaskRun{
			taskRun: tr,
		}
		all[tr.PipelineTaskSpec.ID] = &mtr
	}

	var graph []*memoryTaskRun

	// TODO: Test with multiple and single null successor IDs
	// https://www.pivotaltracker.com/story/show/176557536

	// 2. Fill in predecessor count and next, append firsts to work graph
	for id, mtr := range all {
		for _, pred := range all {
			if !pred.taskRun.PipelineTaskSpec.SuccessorID.IsZero() && pred.taskRun.PipelineTaskSpec.SuccessorID.ValueOrZero() == int64(id) {
				mtr.nPredecessors++
			}
		}

		if mtr.taskRun.PipelineTaskSpec.SuccessorID.IsZero() {
			mtr.next = nil
		} else {
			mtr.next = all[int32(mtr.taskRun.PipelineTaskSpec.SuccessorID.ValueOrZero())]
		}

		if mtr.nPredecessors == 0 {
			// No predecessors so this is the first one
			graph = append(graph, mtr)
		}
	}

	// 3. Execute tasks using "fan in" job processing

	var updateMu sync.Mutex
	var wg sync.WaitGroup
	wg.Add(len(graph))

	// HACK: This mutex is necessary to work around a bug in the pq driver that
	// causes concurrent database calls inside the same transaction to fail
	// with a mysterious `pq: unexpected Parse response 'C'` error
	// FIXME: Get rid of this by replacing pq with pgx
	var txdbMutex sync.Mutex

	for _, mtr := range graph {
		go func(m *memoryTaskRun) {
			defer wg.Done()
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

				result := r.executeTaskRun(ctx, txdb, run.PipelineSpec, m.taskRun, m.results(), &txdbMutex, l)

				finishedAt := time.Now()

				trr := TaskRunResult{
					ID:         m.taskRun.ID,
					TaskSpecID: m.taskRun.PipelineTaskSpec.ID,
					Result:     result,
					FinishedAt: finishedAt,
					IsTerminal: m.next == nil,
				}

				updateMu.Lock()
				trrs = append(trrs, trr)
				updateMu.Unlock()

				elapsed := finishedAt.Sub(startTaskRun)

				promPipelineTaskExecutionTime.WithLabelValues(fmt.Sprintf("%d", run.PipelineSpec.ID), string(m.taskRun.PipelineTaskSpec.Type)).Set(float64(elapsed))
				var status string
				if result.Error != nil {
					status = "error"
				} else {
					status = "completed"
				}
				promPipelineTasksTotalFinished.WithLabelValues(fmt.Sprintf("%d", run.PipelineSpec.ID), string(m.taskRun.PipelineTaskSpec.Type), status).Inc()

				if m.next == nil {
					return
				}

				m.next.predMu.Lock()
				m.next.inputs = append(m.next.inputs, input{result: result, index: m.taskRun.PipelineTaskSpec.Index})
				m.next.nPredecessors--
				m.next.predMu.Unlock()

				m = m.next
			}
		}(mtr)
	}

	wg.Wait()

	runTime := time.Since(startRun)
	l.Debugw("Finished all tasks for pipeline run", "runID", run.ID, "runTime", runTime)

	return trrs, err
}

func (r *runner) executeTaskRun(ctx context.Context, txdb *gorm.DB, spec Spec, taskRun TaskRun, inputs []Result, txdbMutex *sync.Mutex, l logger.Logger) Result {
	loggerFields := []interface{}{
		"taskName", taskRun.PipelineTaskSpec.DotID,
		"taskID", taskRun.PipelineTaskSpecID,
		"runID", taskRun.PipelineRunID,
		"taskRunID", taskRun.ID,
		"specID", taskRun.PipelineTaskSpec.PipelineSpecID,
	}

	task, err := UnmarshalTaskFromMap(
		taskRun.PipelineTaskSpec.Type,
		taskRun.PipelineTaskSpec.JSON.Val,
		taskRun.PipelineTaskSpec.DotID,
		r.config,
		txdb,
		txdbMutex,
	)
	if err != nil {
		l.Errorw("Pipeline task run could not be unmarshaled", append(loggerFields, "error", err)...)
		return Result{Error: err}
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

	result := task.Run(ctx, taskRun, inputs)
	if _, is := result.Error.(FinalErrors); !is && result.Error != nil {
		f := append(loggerFields, "error", result.Error)
		l.Warnw("Pipeline task run errored", f...)
	} else {
		f := append(loggerFields, "result", result.Value)
		switch v := result.Value.(type) {
		case []byte:
			f = append(f, "resultString", fmt.Sprintf("%q", v))
			f = append(f, "resultHex", fmt.Sprintf("%x", v))
		}
		l.Debugw("Pipeline task completed", f...)
	}

	return result
}

// ExecuteAndInsertNewRun bypasses the job pipeline entirely.
// It executes a run in memory then inserts the finished run/task run records, returning the final result
func (r *runner) ExecuteAndInsertNewRun(ctx context.Context, spec Spec, l logger.Logger) (runID int64, result FinalResult, err error) {
	start := time.Now()

	run, err := NewRun(spec, start)
	if err != nil {
		return run.ID, result, errors.Wrapf(err, "error creating new run for spec ID %v", spec.ID)
	}

	trrs, err := r.ExecuteRun(ctx, run, l)
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
