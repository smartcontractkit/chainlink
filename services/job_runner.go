package services

import (
	"fmt"

	"github.com/smartcontractkit/chainlink-go/adapters"
	"github.com/smartcontractkit/chainlink-go/logger"
	"github.com/smartcontractkit/chainlink-go/store/models"
	"github.com/smartcontractkit/chainlink-go/store"
)

func StartJob(run models.JobRun, store *store.Store) error {
	run.Status = "in progress"
	if err := store.Save(&run); err != nil {
		return runJobError(run, err)
	}

	logger.Infow("Starting job", run.ForLogger()...)
	var prevRun models.TaskRun
	for i, taskRun := range run.TaskRuns {
		prevRun = startTask(taskRun, prevRun.Result, store)
		run.TaskRuns[i] = prevRun
		if err := store.Save(&run); err != nil {
			return runJobError(run, err)
		}

		logger.Infow("Task finished", run.ForLogger("task", i, "result", prevRun.Result)...)
		if prevRun.Result.HasError() {
			break
		}
	}

	run.Result = prevRun.Result
	if run.Result.HasError() {
		run.Status = "errored"
	} else {
		run.Status = "completed"
	}

	logger.Infow("Finished job", run.ForLogger()...)
	return runJobError(run, store.Save(&run))
}

func startTask(run models.TaskRun, input models.RunResult, store *store.Store) models.TaskRun {
	run.Status = "in progress"
	adapter, err := adapters.For(run.Task, store)

	if err != nil {
		run.Status = "errored"
		run.Result.SetError(err)
		return run
	}
	run.Result = adapter.Perform(input)

	if run.Result.HasError() {
		run.Status = "errored"
	} else {
		run.Status = "completed"
	}

	return run
}

func runJobError(run models.JobRun, err error) error {
	if err != nil {
		return fmt.Errorf("StartJob#%v: %v", run.JobID, err)
	}
	return nil
}
