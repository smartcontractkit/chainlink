package services

import (
	"fmt"
	"time"

	"chainlink/core/adapters"
	"chainlink/core/logger"
	"chainlink/core/store"
	"chainlink/core/store/models"
	"chainlink/core/store/orm"

	"github.com/pkg/errors"
)

//go:generate mockery -name RunExecutor -output ../internal/mocks/ -case=underscore

// RunExecutor handles the actual running of the job tasks
type RunExecutor interface {
	Execute(*models.ID) error
}

type runExecutor struct {
	store *store.Store
}

// NewRunExecutor initializes a RunExecutor.
func NewRunExecutor(store *store.Store) RunExecutor {
	return &runExecutor{
		store: store,
	}
}

// Execute performs the work associate with a job run
func (je *runExecutor) Execute(runID *models.ID) error {
	run, err := je.store.Unscoped().FindJobRun(runID)
	if err != nil {
		return errors.Wrapf(err, "error finding run %s", runID)
	}

	for taskIndex := range run.TaskRuns {
		if !run.Status.Runnable() {
			break
		}

		taskRun := &run.TaskRuns[taskIndex]
		if taskRun.Status.Completed() {
			continue
		}

		if meetsMinimumConfirmations(&run, taskRun, run.ObservedHeight) {
			result := je.executeTask(&run, taskRun)

			taskRun.ApplyOutput(result)
			run.ApplyOutput(result)

			if !result.Status().Runnable() {
				logger.Debugw("Task execution blocked", run.ForLogger("task", taskRun.ID.String())...)
			}

		} else {
			logger.Debugw("Pausing run pending confirmations",
				run.ForLogger("required_height", taskRun.MinimumConfirmations)...,
			)
			taskRun.Status = models.RunStatusPendingConfirmations
			run.Status = models.RunStatusPendingConfirmations

		}

		if err := je.store.ORM.SaveJobRun(&run); errors.Cause(err) == orm.OptimisticUpdateConflictError {
			logger.Debugw("Optimistic update conflict while updating run", run.ForLogger()...)
			return nil
		} else if err != nil {
			return err
		}
	}

	if run.Status.Finished() {
		logger.Debugw("All tasks complete for run", run.ForLogger()...)
	}
	return nil
}

func (je *runExecutor) executeTask(run *models.JobRun, taskRun *models.TaskRun) models.RunOutput {
	taskCopy := taskRun.TaskSpec // deliberately copied to keep mutations local

	params, err := models.Merge(taskCopy.Params, run.Overrides)
	if err != nil {
		return models.NewRunOutputError(err)
	}
	taskCopy.Params = params

	adapter, err := adapters.For(taskCopy, je.store.Config, je.store.ORM)
	if err != nil {
		return models.NewRunOutputError(err)
	}

	previousTaskRun := run.PreviousTaskRun()

	previousTaskInput := models.JSON{}
	if previousTaskRun != nil {
		previousTaskInput = previousTaskRun.Result.Data
	}

	data, err := models.Merge(previousTaskInput, currentTaskRun.Result.Data, run.Overrides)
	if err != nil {
		return models.NewRunOutputError(err)
	}

	input := *models.NewRunInput(run.ID, data, taskRun.Status)

	start := time.Now()
	result := adapter.Perform(input, je.store)
	logger.Debugw(fmt.Sprintf("Executed task %s", taskCopy.Type), []interface{}{
		"task", taskRun.ID.String(),
		"result", result.Status(),
		"result_data", result.Data(),
		"elapsed", time.Since(start).Seconds(),
	}...)

	return result
}
