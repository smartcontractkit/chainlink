package services

import (
	"fmt"

	"github.com/smartcontractkit/chainlink/adapters"
	"github.com/smartcontractkit/chainlink/logger"
	"github.com/smartcontractkit/chainlink/store"
	"github.com/smartcontractkit/chainlink/store/models"
)

func BeginRun(job *models.Job, store *store.Store) (*models.JobRun, error) {
	run, err := BuildRun(job, store)
	if err != nil {
		return nil, fmt.Errorf("BeginRun: %v", err.Error())
	}
	return run, ExecuteRun(run, store)
}

func BuildRun(job *models.Job, store *store.Store) (*models.JobRun, error) {
	now := store.Clock.Now()
	if !job.Started(now) {
		return nil, fmt.Errorf("BeginRun: %v before job's start time(%v)", now, job.StartAt)
	}
	if job.Ended(now) {
		return nil, fmt.Errorf("BeginRun: %v past job's end time(%v)", now, job.EndAt)
	}
	return job.NewRun(), nil
}

func ExecuteRun(run *models.JobRun, store *store.Store) error {
	run.Status = models.StatusInProgress
	if err := store.Save(run); err != nil {
		return runJobError(run, err)
	}

	logger.Infow("Starting job", run.ForLogger()...)
	unfinished := run.UnfinishedTaskRuns()
	offset := len(run.TaskRuns) - len(unfinished)
	prevRun := run.NextTaskRun()
	for i, taskRun := range unfinished {
		prevRun = startTask(taskRun, prevRun.Result, store)
		run.TaskRuns[i+offset] = prevRun
		if err := store.Save(run); err != nil {
			return runJobError(run, err)
		}

		if prevRun.Result.Pending {
			logger.Infow("Task pending", run.ForLogger("task", i, "result", prevRun.Result)...)
			break
		} else {
			logger.Infow("Task finished", run.ForLogger("task", i, "result", prevRun.Result)...)
			if prevRun.Result.HasError() {
				break
			}
		}
	}

	run.Result = prevRun.Result
	if run.Result.HasError() {
		run.Status = models.StatusErrored
	} else if run.Result.Pending {
		run.Status = models.StatusPending
	} else {
		run.Status = models.StatusCompleted
	}

	logger.Infow("Finished job", run.ForLogger()...)
	return runJobError(run, store.Save(run))
}

func startTask(
	run models.TaskRun,
	input models.RunResult,
	store *store.Store,
) models.TaskRun {
	run.Status = models.StatusInProgress
	adapter, err := adapters.For(run.Task)

	if err != nil {
		run.Status = models.StatusErrored
		run.Result.SetError(err)
		return run
	}
	run.Result = adapter.Perform(input, store)

	if run.Result.HasError() {
		run.Status = models.StatusErrored
	} else if run.Result.Pending {
		run.Status = models.StatusPending
	} else {
		run.Status = models.StatusCompleted
	}

	return run
}

func runJobError(run *models.JobRun, err error) error {
	if err != nil {
		return fmt.Errorf("executeRun#%v: %v", run.JobID, err)
	}
	return nil
}
