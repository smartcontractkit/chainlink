package services

import (
	"fmt"
	"time"

	"github.com/smartcontractkit/chainlink/core/adapters"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/synchronization"
	"github.com/smartcontractkit/chainlink/core/store"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/store/orm"

	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	promAdapterCallsVec = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "adapter_perform_complete_total",
		Help: "The total number of adapters which have completed",
	},
		[]string{"job_spec_id", "task_type", "status"},
	)
)

//go:generate mockery -name RunExecutor -output ../internal/mocks/ -case=underscore

// RunExecutor handles the actual running of the job tasks
type RunExecutor interface {
	Execute(*models.ID) error
}

type runExecutor struct {
	store       *store.Store
	statsPusher synchronization.StatsPusher
}

// NewRunExecutor initializes a RunExecutor.
func NewRunExecutor(store *store.Store, statsPusher synchronization.StatsPusher) RunExecutor {
	return &runExecutor{
		store:       store,
		statsPusher: statsPusher,
	}
}

// Execute performs the work associate with a job run
func (re *runExecutor) Execute(runID *models.ID) error {
	run, err := re.store.Unscoped().FindJobRun(runID)
	if err != nil {
		return errors.Wrapf(err, "error finding run %s", runID)
	}

	for taskIndex := range run.TaskRuns {
		taskRun := &run.TaskRuns[taskIndex]
		if !run.GetStatus().Runnable() {
			logger.Debugw("Run execution blocked", run.ForLogger("task", taskRun.ID.String())...)
			break
		}

		if taskRun.Status.Completed() {
			continue
		}

		if meetsMinRequiredIncomingConfirmations(&run, taskRun, run.ObservedHeight) {
			start := time.Now()

			// NOTE: adapters may define and return the new job run status in here
			result := re.executeTask(&run, taskRun)

			taskRun.ApplyOutput(result)
			run.ApplyOutput(result)

			elapsed := time.Since(start).Seconds()

			logger.Debugw(fmt.Sprintf("Executed task %s", taskRun.TaskSpec.Type), run.ForLogger("task", taskRun.ID.String(), "elapsed", elapsed)...)

		} else {
			logger.Debugw("Pausing run pending incoming confirmations",
				run.ForLogger("required_height", taskRun.MinRequiredIncomingConfirmations)...,
			)
			taskRun.Status = models.RunStatusPendingIncomingConfirmations
			run.SetStatus(models.RunStatusPendingIncomingConfirmations)

		}

		if err := re.store.ORM.SaveJobRun(&run); errors.Cause(err) == orm.ErrOptimisticUpdateConflict {
			logger.Debugw("Optimistic update conflict while updating run", run.ForLogger()...)
			return nil
		} else if err != nil {
			return err
		}

		re.statsPusher.PushNow()
	}

	if run.GetStatus().Finished() {
		if run.GetStatus().Errored() {
			logger.Warnw("Task failed", run.ForLogger()...)
		} else {
			logger.Debugw("All tasks complete for run", run.ForLogger()...)
		}
	}
	return nil
}

func (re *runExecutor) executeTask(run *models.JobRun, taskRun *models.TaskRun) models.RunOutput {
	taskCopy := taskRun.TaskSpec // deliberately copied to keep mutations local

	params, err := models.Merge(run.RunRequest.RequestParams, taskCopy.Params)
	if err != nil {
		return models.NewRunOutputError(err)
	}
	taskCopy.Params = params

	adapter, err := adapters.For(taskCopy, re.store.Config, re.store.ORM)
	if err != nil {
		return models.NewRunOutputError(err)
	}

	previousTaskRun := run.PreviousTaskRun()

	previousTaskInput := models.JSON{}
	if previousTaskRun != nil {
		previousTaskInput = previousTaskRun.Result.Data
	}

	data, err := models.Merge(run.RunRequest.RequestParams, previousTaskInput, taskRun.Result.Data)
	if err != nil {
		return models.NewRunOutputError(err)
	}

	input := *models.NewRunInput(run.ID, data, taskRun.Status)
	result := adapter.Perform(input, re.store)
	promAdapterCallsVec.WithLabelValues(run.JobSpecID.String(), string(adapter.TaskType()), string(result.Status())).Inc()

	return result
}
