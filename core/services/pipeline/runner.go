package pipeline

import (
	"context"
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
		processTasks utils.SleeperTask
		orm          ORM
		config       Config
		chStop       chan struct{}
		chDone       chan struct{}
	}
)

func NewRunner(orm ORM, config Config) *runner {
	r := &runner{
		orm:    orm,
		config: config,
		chStop: make(chan struct{}),
		chDone: make(chan struct{}),
	}
	r.processTasks = utils.NewSleeperTask(
		utils.SleeperTaskFuncWorker(r.processIncompleteTaskRuns),
	)
	return r
}

func (r *runner) Start() {
	go func() {
		defer close(r.chDone)

		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-r.chStop:
				return
			case <-ticker.C:
				r.processIncompleteTaskRuns()
			}
		}
	}()
}

func (r *runner) Stop() {
	close(r.chStop)
	<-r.chDone
}

func (r *runner) CreateRun(jobID int32) (int64, error) {
	runID, err := r.orm.CreateRun(jobID)
	if err != nil {
		return 0, err
	}
	r.processTasks.WakeUp()
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
func (r *runner) processIncompleteTaskRuns() {
	var done bool
	var err error
	for !done {
		// var runID int64
		done, err = r.orm.ProcessNextUnclaimedTaskRun(func(jobID int32, taskRun TaskRun, predecessors []TaskRun) (result Result) {
			// runID = taskRun.PipelineRunID

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
				logger.Errorw("Pipeline task run errored", append(loggerFields, "error", err)...)
				return Result{Error: err}
			}

			result = task.Run(inputs)
			if result.Error != nil {
				logger.Errorw("Pipeline task run errored", append(loggerFields, "error", result.Error)...)
			}
			return result
		})
		if err != nil {
			logger.Errorf("Error processing incomplete task runs: %v", err)
			return
		}

		// err = r.orm.NotifyCompletion(runID)
		// if err != nil {
		// 	logger.Errorf("Error calling pg_notify for run %v: %v", runID, err)
		// }
	}
}
