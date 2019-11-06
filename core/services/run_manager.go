package services

import (
	"fmt"
	"math/big"

	"chainlink/core/adapters"
	"chainlink/core/logger"
	clnull "chainlink/core/null"
	"chainlink/core/store"
	"chainlink/core/store/assets"
	"chainlink/core/store/models"
	"chainlink/core/store/orm"
	"chainlink/core/utils"

	"github.com/pkg/errors"
)

// RecurringScheduleJobError contains the field for the error message.
type RecurringScheduleJobError struct {
	msg string
}

// Error returns the error message for the run.
func (err RecurringScheduleJobError) Error() string {
	return err.msg
}

//go:generate mockery -name RunManager -output ../internal/mocks/ -case=underscore

// RunManager supplies methods for queueing, resuming and cancelling jobs in
// the RunQueue
type RunManager interface {
	Create(
		jobSpecID *models.ID,
		initiator *models.Initiator,
		data *models.JSON,
		creationHeight *big.Int,
		runRequest *models.RunRequest) (*models.JobRun, error)
	CreateErrored(
		jobSpecID *models.ID,
		initiator models.Initiator,
		err error) (*models.JobRun, error)
	ResumePending(
		runID *models.ID,
		input models.BridgeRunResult) error
	Cancel(runID *models.ID) error

	ResumeAllInProgress() error
	ResumeAllConfirming(currentBlockHeight *big.Int) error
	ResumeAllConnecting() error
}

// runManager implements RunManager
type runManager struct {
	orm       *orm.ORM
	runQueue  RunQueue
	txManager store.TxManager
	config    orm.ConfigReader
	clock     utils.AfterNower
}

func newRun(
	job *models.JobSpec,
	initiator *models.Initiator,
	data *models.JSON,
	currentHeight *big.Int,
	payment *assets.Link,
	config orm.ConfigReader,
	orm *orm.ORM,
	clock utils.AfterNower,
	txManager store.TxManager) (*models.JobRun, error) {

	now := clock.Now()
	if !job.Started(now) {
		return nil, RecurringScheduleJobError{
			msg: fmt.Sprintf("Job runner: Job %v unstarted: %v before job's start time %v", job.ID, now, job.EndAt),
		}
	}

	if job.Ended(now) {
		return nil, RecurringScheduleJobError{
			msg: fmt.Sprintf("Job runner: Job %v ended: %v past job's end time %v", job.ID, now, job.EndAt),
		}
	}

	run := job.NewRun(*initiator)
	run.Overrides = *data
	run.CreationHeight = models.NewBig(currentHeight)
	run.ObservedHeight = models.NewBig(currentHeight)

	if !MeetsMinimumPayment(job.MinPayment, payment) {
		logger.Infow("Rejecting run with insufficient payment",
			run.ForLogger(
				"input_payment", payment,
				"required_payment", job.MinPayment)...)

		err := fmt.Errorf(
			"Rejecting job %s with payment %s below job-specific-minimum threshold (%s)",
			job.ID,
			payment,
			job.MinPayment.Text(10))
		run.SetError(err)
	}

	cost := &assets.Link{}
	cost.Set(job.MinPayment)
	for i, taskRun := range run.TaskRuns {
		adapter, err := adapters.For(taskRun.TaskSpec, config, orm)

		if err != nil {
			run.SetError(err)
			return &run, nil
		}

		if job.MinPayment.IsZero() {
			mp := adapter.MinContractPayment()
			if mp != nil {
				cost.Add(cost, mp)
			}
		}

		if currentHeight != nil {
			run.TaskRuns[i].MinimumConfirmations = clnull.Uint32From(
				utils.MaxUint32(
					config.MinIncomingConfirmations(),
					taskRun.TaskSpec.Confirmations.Uint32,
					adapter.MinConfs()),
			)
		}
	}

	// payment is only present for runs triggered by runlogs
	if payment != nil {
		if cost.Cmp(payment) > 0 {
			logger.Debugw("Rejecting run with insufficient payment",
				run.ForLogger(
					"input_payment", payment,
					"required_payment", cost)...)

			err := fmt.Errorf(
				"Rejecting job %s with payment %s below minimum threshold (%s)",
				job.ID,
				payment,
				config.MinimumContractPayment().Text(10))
			run.SetError(err)
		}
	}

	if len(run.TaskRuns) == 0 {
		run.SetError(fmt.Errorf("invariant for job %s: no tasks to run in NewRun", job.ID))
	}

	if !run.Status.Runnable() {
		return &run, nil
	}

	initialTask := run.TaskRuns[0]
	validateMinimumConfirmations(&run, &initialTask, run.CreationHeight, txManager)
	return &run, nil
}

// NewRunManager returns a new job manager
func NewRunManager(
	runQueue RunQueue,
	config orm.ConfigReader,
	orm *orm.ORM,
	txManager store.TxManager,
	clock utils.AfterNower) RunManager {
	return &runManager{
		orm:       orm,
		runQueue:  runQueue,
		txManager: txManager,
		config:    config,
		clock:     clock,
	}
}

// CreateErrored creates a run that is in the errored state. This is a
// special case where this job cannot run but we want to create the run record
// so the error is more visible to the node operator.
func (jm *runManager) CreateErrored(
	jobSpecID *models.ID,
	initiator models.Initiator,
	runErr error) (*models.JobRun, error) {
	job, err := jm.orm.Unscoped().FindJob(jobSpecID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find job spec")
	}

	run := job.NewRun(initiator)
	run.SetError(runErr)
	return &run, jm.orm.CreateJobRun(&run)
}

// Create immediately persists a JobRun and sends it to the RunQueue for
// execution.
func (jm *runManager) Create(
	jobSpecID *models.ID,
	initiator *models.Initiator,
	data *models.JSON,
	creationHeight *big.Int,
	runRequest *models.RunRequest,
) (*models.JobRun, error) {
	logger.Debugw(fmt.Sprintf("New run triggered by %s", initiator.Type),
		"job", jobSpecID.String(),
		"creation_height", creationHeight.String(),
	)

	job, err := jm.orm.Unscoped().FindJob(jobSpecID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find job spec")
	}

	if job.Archived() {
		return nil, RecurringScheduleJobError{
			msg: fmt.Sprintf("Trying to run archived job %s", job.ID),
		}
	}

	run, err := newRun(&job, initiator, data, creationHeight, runRequest.Payment, jm.config, jm.orm, jm.clock, jm.txManager)
	if err != nil {
		return nil, errors.Wrap(err, "newRun failed")
	}

	run.RunRequest = *runRequest

	if err := jm.orm.CreateJobRun(run); err != nil {
		return nil, errors.Wrap(err, "CreateJobRun failed")
	}

	if run.Status == models.RunStatusInProgress {
		logger.Debugw(
			fmt.Sprintf("Executing run originally initiated by %s", run.Initiator.Type),
			run.ForLogger()...,
		)
		jm.runQueue.Run(run)
	}

	return run, nil
}

// ResumeAllConfirming wakes up all jobs that were sleeping because they were
// waiting for block confirmations.
func (jm *runManager) ResumeAllConfirming(currentBlockHeight *big.Int) error {
	return jm.orm.UnscopedJobRunsWithStatus(func(run *models.JobRun) {
		currentTaskRun := run.NextTaskRun()
		if currentTaskRun == nil {
			jm.updateWithError(run, "Attempting to resume confirming run with no remaining tasks %s", run.ID)
			return
		}

		run.ObservedHeight = models.NewBig(currentBlockHeight)
		logger.Debugw(fmt.Sprintf("New head #%s resuming run", currentBlockHeight), run.ForLogger()...)

		validateMinimumConfirmations(run, currentTaskRun, run.ObservedHeight, jm.txManager)

		err := jm.updateAndTrigger(run)
		if err != nil {
			logger.Error("Error saving job run", "error", err)
		}
	}, models.RunStatusPendingConnection, models.RunStatusPendingConfirmations)
}

// ResumeAllConnecting wakes up all tasks that have gone to sleep because they
// needed an ethereum client connection.
func (jm *runManager) ResumeAllConnecting() error {
	return jm.orm.UnscopedJobRunsWithStatus(func(run *models.JobRun) {
		logger.Debugw("New connection resuming run", run.ForLogger()...)

		currentTaskRun := run.NextTaskRun()
		if currentTaskRun == nil {
			jm.updateWithError(run, "Attempting to resume connecting run with no remaining tasks %s", run.ID)
			return
		}

		run.Status = models.RunStatusInProgress
		err := jm.updateAndTrigger(run)
		if err != nil {
			logger.Error("Error saving job run", "error", err)
		}
	}, models.RunStatusPendingConnection, models.RunStatusPendingConfirmations)
}

// ResumePendingTask wakes up a task that required a response from a bridge adapter.
func (jm *runManager) ResumePending(
	runID *models.ID,
	input models.BridgeRunResult,
) error {
	run, err := jm.orm.Unscoped().FindJobRun(runID)
	if err != nil {
		return err
	}

	logger.Debugw("External adapter resuming job", run.ForLogger("input_data", input.Data)...)

	if !run.Status.PendingBridge() {
		return fmt.Errorf("Attempting to resume non pending run %s", run.ID)
	}

	currentTaskRun := run.NextTaskRun()
	if currentTaskRun == nil {
		return jm.updateWithError(&run, "Attempting to resume pending run with no remaining tasks %s", run.ID)
	}

	run.Overrides.Merge(input.Data)

	currentTaskRun.ApplyBridgeRunResult(input)
	run.ApplyBridgeRunResult(input)

	return jm.updateAndTrigger(&run)
}

// ResumeAllInProgress queries the db for job runs that should be resumed
// since a previous node shutdown.
//
// As a result of its reliance on the database, it must run before anything
// persists a job RunStatus to the db to ensure that it only captures pending and in progress
// jobs as a result of the last shutdown, and not as a result of what's happening now.
//
// To recap: This must run before anything else writes job run status to the db,
// ie. tries to run a job.
// https://chainlink/pull/807
func (jm *runManager) ResumeAllInProgress() error {
	return jm.orm.UnscopedJobRunsWithStatus(jm.runQueue.Run, models.RunStatusInProgress, models.RunStatusPendingSleep)
}

// Cancel suspends a running task.
func (jm *runManager) Cancel(runID *models.ID) error {
	run, err := jm.orm.FindJobRun(runID)
	if err != nil {
		return err
	}

	logger.Debugw("Cancelling job", run.ForLogger()...)
	if !run.Status.Runnable() {
		return fmt.Errorf("Cannot cancel a non runnable job")
	}

	currentTaskRun := run.NextTaskRun()
	if currentTaskRun != nil {
		currentTaskRun.Status = models.RunStatusCancelled
	}

	run.Status = models.RunStatusCancelled
	return jm.orm.SaveJobRun(&run)
}

func (jm *runManager) updateWithError(run *models.JobRun, msg string, args ...interface{}) error {
	run.SetError(fmt.Errorf(msg, args...))
	logger.Error(fmt.Sprintf(msg, args...))

	if err := jm.orm.SaveJobRun(run); err != nil {
		logger.Error("Error saving job run", "error", err)
		return err
	}
	return nil
}

func (jm *runManager) updateAndTrigger(run *models.JobRun) error {
	if err := jm.orm.SaveJobRun(run); err != nil {
		return err
	}
	if run.Status == models.RunStatusInProgress {
		jm.runQueue.Run(run)
	}
	return nil
}
