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
		return fmt.Errorf("Error finding run %s", runID.String())
	}

	if !run.Status.Runnable() {
		return fmt.Errorf("Run triggered in non runnable state %s", run.Status)
	}

	for run.Status.Runnable() {
		currentTaskRun := run.NextTaskRun()
		if currentTaskRun == nil {
			return errors.New("Run triggered with no remaining tasks")
		}

		result := je.executeTask(&run, currentTaskRun)

		currentTaskRun.ApplyOutput(result)
		run.ApplyOutput(result)

		if !result.Status().Runnable() {
			logger.Debugw("Task execution blocked", []interface{}{"run", run.ID, "task", currentTaskRun.ID.String(), "state", currentTaskRun.Status}...)
		} else if currentTaskRun.Status.Unstarted() {
			return fmt.Errorf("run %s task %s cannot return a status of empty string or Unstarted", run.ID, currentTaskRun.TaskSpec.Type)
		} else if futureTaskRun := run.NextTaskRun(); futureTaskRun != nil {
			validateMinimumConfirmations(&run, futureTaskRun, run.ObservedHeight, je.store.TxManager)
		}

		if err := je.store.ORM.SaveJobRun(&run); errors.Cause(err) == orm.OptimisticUpdateConflictError {
			return nil
		} else if err != nil {
			return err
		}

		if run.Status.Finished() {
			logger.Debugw("All tasks complete for run", run.ForLogger()...)
			break
		}
	}

	return nil
}

func (je *runExecutor) executeTask(run *models.JobRun, currentTaskRun *models.TaskRun) models.RunOutput {
	taskCopy := currentTaskRun.TaskSpec // deliberately copied to keep mutations local

	var err error
	if taskCopy.Params, err = taskCopy.Params.Merge(run.Overrides); err != nil {
		return models.NewRunOutputError(err)
	}

	adapter, err := adapters.For(taskCopy, je.store.Config, je.store.ORM)
	if err != nil {
		return models.NewRunOutputError(err)
	}

	previousTaskRun := run.PreviousTaskRun()

	data := models.JSON{}
	if previousTaskRun != nil {
		if data, err = previousTaskRun.Result.Data.Merge(currentTaskRun.Result.Data); err != nil {
			return models.NewRunOutputError(err)
		}
	}

	if data, err = run.Overrides.Merge(data); err != nil {
		return models.NewRunOutputError(err)
	}

	input := *models.NewRunInput(run.ID, data, currentTaskRun.Status)

	start := time.Now()
	result := adapter.Perform(input, je.store)
	logger.Debugw(fmt.Sprintf("Executed task %s", taskCopy.Type), []interface{}{
		"task", currentTaskRun.ID.String(),
		"result", result.Status(),
		"result_data", result.Data(),
		"elapsed", time.Since(start).Seconds(),
	}...)

	return result
}
