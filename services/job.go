package services

import (
	"fmt"

	"github.com/smartcontractkit/chainlink-go/adapters"
	"github.com/smartcontractkit/chainlink-go/logger"
	"github.com/smartcontractkit/chainlink-go/models"
)

func StartJob(run models.JobRun, orm *models.ORM) error {
	run.Status = "in progress"
	err := orm.Save(&run)
	if err != nil {
		return runJobError(run, err)
	}

	logger.Infow("Starting job", run.ForLogger()...)
	var prevRun models.TaskRun
	for i, taskRun := range run.TaskRuns {
		prevRun = startTask(taskRun, prevRun.Result)
		run.TaskRuns[i] = prevRun
		err = orm.Save(&run)
		if err != nil {
			return runJobError(run, err)
		}

		logger.Infow("Task finished", run.ForLogger("task", i, "result", prevRun.Result)...)
		if prevRun.Result.Error != nil {
			break
		}
	}

	run.Result = prevRun.Result
	if run.Result.Error != nil {
		run.Status = "errored"
	} else {
		run.Status = "completed"
	}

	logger.Infow("Finished job", run.ForLogger()...)
	return runJobError(run, orm.Save(&run))
}

func startTask(run models.TaskRun, input models.RunResult) models.TaskRun {
	run.Status = "in progress"
	adapter, err := adapters.For(run.Task)

	if err != nil {
		run.Status = "errored"
		run.Result.Error = err
		return run
	}
	run.Result = adapter.Perform(input)

	if run.Result.Error != nil {
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
