package services

import (
	"fmt"
	"time"

	"github.com/smartcontractkit/chainlink/core/store/orm"

	"github.com/smartcontractkit/chainlink/core/adapters"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/eth"
	"github.com/smartcontractkit/chainlink/core/services/keystore"
	"github.com/smartcontractkit/chainlink/core/services/synchronization"
	"github.com/smartcontractkit/chainlink/core/store"
	"github.com/smartcontractkit/chainlink/core/store/models"

	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	uuid "github.com/satori/go.uuid"
)

var (
	promAdapterCallsVec = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "adapter_perform_complete_total",
		Help: "The total number of adapters which have completed",
	},
		[]string{"job_spec_id", "task_type", "status"},
	)
)

//go:generate mockery --name RunExecutor --output ../internal/mocks/ --case=underscore

// RunExecutor handles the actual running of the job tasks
type RunExecutor interface {
	Execute(uuid.UUID) error
}

type runExecutor struct {
	store       *store.Store
	keyStore    *keystore.Master
	statsPusher synchronization.StatsPusher
}

// NewRunExecutor initializes a RunExecutor.
func NewRunExecutor(store *store.Store, keyStore *keystore.Master, statsPusher synchronization.StatsPusher) RunExecutor {
	return &runExecutor{
		store:       store,
		keyStore:    keyStore,
		statsPusher: statsPusher,
	}
}

// Execute performs the work associate with a job run
func (re *runExecutor) Execute(runID uuid.UUID) error {
	logger.Debugw("runExecutor woke up", "runID", runID.String())
	run, err := re.store.Unscoped().FindJobRun(runID)
	if err != nil {
		return errors.Wrapf(err, "error finding run %s", runID)
	}

	validated := false
	for taskIndex := range run.TaskRuns {
		taskRun := &run.TaskRuns[taskIndex]
		if !run.GetStatus().Runnable() {
			logger.Debugw("Run execution blocked", run.ForLogger("task", taskRun.ID.String())...)
			break
		}

		if taskRun.Status.Completed() {
			continue
		}

		if !meetsMinRequiredIncomingConfirmations(&run, taskRun, run.ObservedHeight) {
			logger.Debugw("Pausing run pending incoming confirmations",
				run.ForLogger("required_height", taskRun.MinRequiredIncomingConfirmations)...,
			)
			taskRun.Status = models.RunStatusPendingIncomingConfirmations
			run.SetStatus(models.RunStatusPendingIncomingConfirmations)

		} else if err := validateOnMainChainOnce(validated, &run, taskRun, re.store.EthClient); err != nil {
			logger.Warnw("Failure while trying to validate chain",
				run.ForLogger("error", err)...,
			)

			taskRun.SetError(err)
			run.SetError(err)

		} else {
			start := time.Now()

			// NOTE: adapters may define and return the new job run status in here
			result := re.executeTask(&run, *taskRun)

			taskRun.ApplyOutput(result)
			run.ApplyOutput(result)

			elapsed := time.Since(start).Seconds()

			logger.Debugw(fmt.Sprintf("Executed task %s", taskRun.TaskSpec.Type), run.ForLogger("task", taskRun.ID.String(), "elapsed", elapsed)...)
		}

		validated = true

		if err := re.store.ORM.SaveJobRun(&run); errors.Is(err, orm.ErrOptimisticUpdateConflict) {
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

func validateOnMainChainOnce(validated bool, run *models.JobRun, taskRun *models.TaskRun, ethClient eth.Client) error {
	if validated {
		return nil
	}
	return validateOnMainChain(run, taskRun, ethClient)
}

func (re *runExecutor) executeTask(run *models.JobRun, taskRun models.TaskRun) models.RunOutput {
	taskSpec := taskRun.TaskSpec

	params, err := models.Merge(run.RunRequest.RequestParams, taskSpec.Params)
	if err != nil {
		return models.NewRunOutputError(err)
	}
	taskSpec.Params = params

	adapter, err := adapters.For(taskSpec, re.store.Config, re.store.ORM)
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

	input := *models.NewRunInput(*run, taskRun.ID, data, taskRun.Status)
	result := adapter.Perform(input, re.store, re.keyStore)
	promAdapterCallsVec.WithLabelValues(run.JobSpecID.String(), string(adapter.TaskType()), string(result.Status())).Inc()

	return result
}

// NullRunExecutor implements Null pattern for RunExecutor interface
type NullRunExecutor struct{}

func (NullRunExecutor) Execute(uuid.UUID) error {
	return errors.New("NullRunExecutor#Execute should never be called")
}
