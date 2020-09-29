package pipeline

import (
	"context"
	"sync"
	"time"

	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/utils"
)

type (
	// Runner polls the DB occasionally for incomplete TaskRuns and runs them.
	// For a TaskRun to be eligible to be run, its parent/input tasks must
	// already all be complete.
	Runner interface {
		Start()
		Stop()
		CreateRun(jobID int32) (int64, error)
		AwaitRun(ctx context.Context, runID int64) error
		ResultsForRun(runID int64) ([]Result, error)
	}

	runner struct {
		processIncompleteTaskRuns utils.SleeperTask
		orm                       ORM
		config                    Config

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
	r.processIncompleteTaskRuns = utils.NewSleeperTask(
		utils.SleeperTaskFuncWorker(r.processIncompleteTaskRunsWorker),
	)
	return r
}

func (r *runner) Start() {
	r.AssertNeverStarted()
	go r.runLoop()
}

func (r *runner) Stop() {
	r.AssertNeverStopped()
	close(r.chStop)
	<-r.chDone
}

func (r *runner) runLoop() {
	defer close(r.chDone)

	newRunListener, err := r.orm.ListenForNewRuns()
	if err != nil {
		logger.Errorw(`Pipeline runner failed to subscribe to "new run" events, falling back to polling`, "error", err)
	}

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-r.chStop:
			if newRunListener != nil {
				err := newRunListener.Stop()
				if err != nil {
					logger.Errorw(`Error stopping pipeline runner's "new runs" listener`, "error", err)
				}
			}
			r.processIncompleteTaskRuns.Stop()
			return

		case <-newRunListener.Events():
			r.processIncompleteTaskRuns.WakeUp()
		case <-ticker.C:
			r.processIncompleteTaskRuns.WakeUp()
		}
	}
}

func (r *runner) CreateRun(jobID int32) (int64, error) {
	logger.Infow("Creating new pipeline run", "jobID", jobID)

	runID, err := r.orm.CreateRun(jobID)
	if err != nil {
		logger.Errorw("Error creating new pipeline run", "jobID", jobID, "error", err)
		return 0, err
	}
	return runID, nil
}

func (r *runner) AwaitRun(ctx context.Context, runID int64) error {
	return r.orm.AwaitRun(ctx, runID)
}

func (r *runner) ResultsForRun(runID int64) ([]Result, error) {
	return r.orm.ResultsForRun(runID)
}

// NOTE: This could potentially run on a different machine in the cluster than
// the one that originally added the task runs.
func (r *runner) processIncompleteTaskRunsWorker() {
	threads := int(r.config.PipelineRunnerParallelism())

	var wg sync.WaitGroup
	wg.Add(threads)

	for i := 0; i < threads; i++ {
		go func() {
			defer wg.Done()

			var done bool
			var err error
			for !done {
				done, err = r.orm.ProcessNextUnclaimedTaskRun(func(jobID int32, taskRun TaskRun, predecessors []TaskRun) Result {
					loggerFields := []interface{}{
						"jobID", jobID,
						"taskName", taskRun.DotID,
						"taskID", taskRun.PipelineTaskSpecID,
						"runID", taskRun.PipelineRunID,
						"taskRunID", taskRun.ID,
					}

					logger.Infow("Running pipeline task", loggerFields...)

					inputs := make([]Result, len(predecessors))
					for i, predecessor := range predecessors {
						inputs[i] = predecessor.Result()
					}

					task, err := UnmarshalTaskFromMap(taskRun.PipelineTaskSpec.Type, taskRun.PipelineTaskSpec.JSON.Val, taskRun.DotID, r.orm, r.config)
					if err != nil {
						logger.Errorw("Pipeline task run could not be unmarshaled", append(loggerFields, "error", err)...)
						return Result{Error: err}
					}

					result := task.Run(inputs)
					if result.Error != nil {
						logger.Errorw("Pipeline task run errored", append(loggerFields, "error", result.Error)...)
					}
					return result
				})
				if err != nil {
					logger.Errorf("Error processing incomplete task runs: %v", err)
					return
				}
			}
		}()
	}
	wg.Wait()
}
