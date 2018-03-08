package services

import (
	"fmt"
	"time"

	"github.com/smartcontractkit/chainlink/adapters"
	"github.com/smartcontractkit/chainlink/logger"
	"github.com/smartcontractkit/chainlink/store"
	"github.com/smartcontractkit/chainlink/store/models"
	null "gopkg.in/guregu/null.v3"
)

// BeginRun creates a new run if the job is valid and starts the job.
func BeginRun(job models.JobSpec, store *store.Store, input models.RunResult) (models.JobRun, error) {
	run, err := BuildRun(job, store)
	if err != nil {
		return models.JobRun{}, err
	}
	return ExecuteRun(run, store, input)
}

// BuildRun checks to ensure the given job has not started or ended before
// creating a new run for the job.
func BuildRun(job models.JobSpec, store *store.Store) (models.JobRun, error) {
	now := store.Clock.Now()
	if !job.Started(now) {
		return models.JobRun{}, JobRunnerError{
			msg: fmt.Sprintf("Job runner: Job %v unstarted: %v before job's start time %v", job.ID, now, job.EndAt),
		}
	}
	if job.Ended(now) {
		return models.JobRun{}, JobRunnerError{
			msg: fmt.Sprintf("Job runner: Job %v ended: %v past job's end time %v", job.ID, now, job.EndAt),
		}
	}
	return job.NewRun(), nil
}

// ExecuteRun starts the job and executes task runs within that job in the
// order defined in the run for as long as they do not return errors. Results
// are saved in the store (db).
func ExecuteRun(run models.JobRun, store *store.Store, input models.RunResult) (models.JobRun, error) {
	run.Status = models.StatusInProgress
	if err := store.Save(&run); err != nil {
		return run, wrapError(run, err)
	}

	logger.Infow("Starting job", run.ForLogger()...)
	unfinished := run.UnfinishedTaskRuns()
	offset := len(run.TaskRuns) - len(unfinished)
	prevRun := unfinished[0]

	merged, err := prevRun.Result.Merge(input)
	if err != nil {
		return run, wrapError(run, err)
	}
	prevRun.Result = merged

	for i, taskRunTemplate := range unfinished {
		taskRun, err := taskRunTemplate.MergeTaskParams(input.Data)
		if err != nil {
			return run, wrapError(run, err)
		}
		prevRun = startTask(taskRun, prevRun.Result, store)
		logger.Debugw("Produced task run", "tr", prevRun)
		run.TaskRuns[i+offset] = prevRun
		if err := store.Save(&run); err != nil {
			return run, wrapError(run, err)
		}

		if prevRun.Result.Pending {
			logger.Infow(fmt.Sprintf("Task %v pending", taskRun.Task.Type), taskRun.ForLogger("task", i, "result", prevRun.Result)...)
			break
		}
		logger.Infow(fmt.Sprintf("Task %v finished", taskRun.Task.Type), taskRun.ForLogger("task", i, "result", prevRun.Result)...)
		if prevRun.Result.HasError() {
			break
		}
	}

	run.Result = prevRun.Result
	if run.Result.HasError() {
		run.Status = models.StatusErrored
	} else if run.Result.Pending {
		run.Status = models.StatusPending
	} else {
		run.Status = models.StatusCompleted
		run.CompletedAt = null.Time{Time: time.Now(), Valid: true}
	}

	logger.Infow("Finished current job run execution", run.ForLogger()...)
	return run, wrapError(run, store.Save(&run))
}

func startTask(
	run models.TaskRun,
	input models.RunResult,
	store *store.Store,
) models.TaskRun {
	run.Status = models.StatusInProgress

	adapter, err := adapters.For(run.Task, store)
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

func wrapError(run models.JobRun, err error) error {
	if err != nil {
		return fmt.Errorf("ExecuteRun: Job#%v: %v", run.JobID, err)
	}
	return nil
}

// JobRunnerError contains the field for the error message.
type JobRunnerError struct {
	msg string
}

// Error returns the error message for the run.
func (err JobRunnerError) Error() string {
	return err.msg
}
