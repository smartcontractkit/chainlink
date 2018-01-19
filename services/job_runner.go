package services

import (
	"fmt"
	"time"

	"github.com/smartcontractkit/chainlink/adapters"
	"github.com/smartcontractkit/chainlink/logger"
	"github.com/smartcontractkit/chainlink/store"
	"github.com/smartcontractkit/chainlink/store/models"
)

func BeginRun(job *models.Job, store *store.Store) (*models.JobRun, error) {
	run, err := BuildRun(job, store)
	if err != nil {
		return nil, err
	}
	return run, ExecuteRun(run, store)
}

func BuildRun(job *models.Job, store *store.Store) (*models.JobRun, error) {
	now := store.Clock.Now()
	if !job.Started(now) {
		return nil, NewJobUnstartedError(job, now)
	}
	if job.Ended(now) {
		return nil, NewJobEndedError(job, now)
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

type JobEndedError struct {
	msg string
}

func (err JobEndedError) Error() string {
	return err.msg
}

func NewJobEndedError(j *models.Job, t time.Time) JobEndedError {
	return JobEndedError{
		fmt.Sprintf("Job runner: Job %v ended: %v past job's end time %v", j.ID, t, j.EndAt),
	}
}

type JobUnstartedError struct {
	msg string
}

func (err JobUnstartedError) Error() string {
	return err.msg
}

func NewJobUnstartedError(j *models.Job, t time.Time) JobUnstartedError {
	return JobUnstartedError{
		fmt.Sprintf("Job runner: Job %v unstarted: %v before job's start time %v", j.ID, t, j.EndAt),
	}
}
