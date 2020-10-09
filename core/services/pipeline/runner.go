package pipeline

import (
	"context"
	"sync"
	"time"

	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/utils"
)

type (
	// Runner checks the DB for incomplete TaskRuns and runs them.  For a
	// TaskRun to be eligible to be run, its parent/input tasks must already
	// all be complete.
	Runner interface {
		Start()
		Stop()
		CreateRun(ctx context.Context, jobID int32, meta map[string]interface{}) (int64, error)
		AwaitRun(ctx context.Context, runID int64) error
		ResultsForRun(ctx context.Context, runID int64) ([]interface{}, error)
	}

	runner struct {
		orm                             ORM
		config                          Config
		processIncompleteTaskRunsWorker utils.SleeperTask
		newRunsListener                 *utils.PostgresEventListener

		utils.StartStopOnce
		chStop chan struct{}
		chDone chan struct{}
	}
)

func NewRunner(orm ORM, config Config) *runner {
	r := &runner{
		orm:    orm,
		config: config,
		chStop: make(chan struct{}),
		chDone: make(chan struct{}),
	}
	r.processIncompleteTaskRunsWorker = utils.NewSleeperTask(
		utils.SleeperTaskFuncWorker(r.processIncompleteTaskRuns),
	)
	return r
}

func (r *runner) Start() {
	if !r.OkayToStart() {
		logger.Error("Pipeline runner has already been started")
		return
	}
	go r.runLoop()
}

func (r *runner) Stop() {
	if !r.OkayToStop() {
		logger.Error("Pipeline runner has already been stopped")
		return
	}

	r.processIncompleteTaskRunsWorker.Stop()

	if r.newRunsListener != nil {
		err := r.newRunsListener.Stop()
		if err != nil {
			logger.Errorw(`Error stopping pipeline runner's "new runs" listener`, "error", err)
		}
	}
	close(r.chStop)
	<-r.chDone
}

func (r *runner) runLoop() {
	defer close(r.chDone)

	var err error
	r.newRunsListener, err = r.orm.ListenForNewRuns()
	if err != nil {
		logger.Errorw(`Pipeline runner failed to subscribe to "new run" events, falling back to polling`, "error", err)
	}

	ticker := time.NewTicker(r.config.JobPipelineDBPollInterval())
	defer ticker.Stop()

	for {
		select {
		case <-r.chStop:
			return
		case <-r.newRunsListener.Events():
			r.processIncompleteTaskRunsWorker.WakeUp()
		case <-ticker.C:
			r.processIncompleteTaskRunsWorker.WakeUp()
		}
	}
}

func (r *runner) CreateRun(ctx context.Context, jobID int32, meta map[string]interface{}) (int64, error) {
	logger.Infow("Creating new pipeline run", "jobID", jobID)

	runID, err := r.orm.CreateRun(ctx, jobID, meta)
	if err != nil {
		logger.Errorw("Error creating new pipeline run", "jobID", jobID, "error", err)
		return 0, err
	}
	return runID, nil
}

func (r *runner) AwaitRun(ctx context.Context, runID int64) error {
	ctx, cancel := utils.CombinedContext(r.chStop, ctx)
	defer cancel()
	return r.orm.AwaitRun(ctx, runID)
}

func (r *runner) ResultsForRun(ctx context.Context, runID int64) ([]interface{}, error) {
	ctx, cancel := utils.CombinedContext(r.chStop, ctx)
	defer cancel()
	return r.orm.ResultsForRun(ctx, runID)
}

// NOTE: This could potentially run on a different machine in the cluster than
// the one that originally added the task runs.
func (r *runner) processIncompleteTaskRuns() {
	threads := int(r.config.JobPipelineParallelism())

	var wg sync.WaitGroup
	wg.Add(threads)

	for i := 0; i < threads; i++ {
		go func() {
			defer wg.Done()
			for {
				select {
				case <-r.chStop:
					return
				default:
				}

				anyRemaining, err := r.processRun()
				if err != nil {
					logger.Errorf("Error processing incomplete task runs: %v", err)
					return
				} else if !anyRemaining {
					return
				}
			}
		}()
	}
	wg.Wait()
}

func (r *runner) processRun() (anyRemaining bool, err error) {
	ctx, cancel := utils.CombinedContext(r.chStop, r.config.JobPipelineMaxTaskDuration())
	defer cancel()

	return r.orm.ProcessNextUnclaimedTaskRun(ctx, func(jobID int32, taskRun TaskRun, predecessors []TaskRun) Result {
		loggerFields := []interface{}{
			"jobID", jobID,
			"taskName", taskRun.PipelineTaskSpec.DotID,
			"taskID", taskRun.PipelineTaskSpecID,
			"runID", taskRun.PipelineRunID,
			"taskRunID", taskRun.ID,
		}

		logger.Infow("Running pipeline task", loggerFields...)

		inputs := make([]Result, len(predecessors))
		for i, predecessor := range predecessors {
			inputs[i] = predecessor.Result()
		}

		task, err := UnmarshalTaskFromMap(
			taskRun.PipelineTaskSpec.Type,
			taskRun.PipelineTaskSpec.JSON.Val,
			taskRun.PipelineTaskSpec.DotID,
			r.orm,
			r.config,
		)
		if err != nil {
			logger.Errorw("Pipeline task run could not be unmarshaled", append(loggerFields, "error", err)...)
			return Result{Error: err}
		}

		result := task.Run(taskRun, inputs)
		if result.Error != nil {
			logger.Errorw("Pipeline task run errored", append(loggerFields, "error", result.Error)...)
		} else {
			logger.Infow("Pipeline task completed", append(loggerFields, "result", result.Value)...)
		}

		return result
	})
}
