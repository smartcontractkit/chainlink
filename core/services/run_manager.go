package services

import (
	"fmt"
	"math/big"
	"time"

	"github.com/smartcontractkit/chainlink/core/adapters"
	"github.com/smartcontractkit/chainlink/core/assets"
	"github.com/smartcontractkit/chainlink/core/logger"
	clnull "github.com/smartcontractkit/chainlink/core/null"
	"github.com/smartcontractkit/chainlink/core/services/synchronization"
	"github.com/smartcontractkit/chainlink/core/store"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/store/orm"
	"github.com/smartcontractkit/chainlink/core/utils"

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

//go:generate mockery --name RunManager --output ../internal/mocks/ --case=underscore

// RunManager supplies methods for queueing, resuming and cancelling jobs in
// the RunQueue
type RunManager interface {
	Create(
		jobSpecID *models.ID,
		initiator *models.Initiator,
		creationHeight *big.Int,
		runRequest *models.RunRequest) (*models.JobRun, error)
	CreateErrored(
		jobSpecID *models.ID,
		initiator models.Initiator,
		err error) (*models.JobRun, error)
	ResumePendingBridge(
		runID *models.ID,
		input models.BridgeRunResult) error
	Cancel(runID *models.ID) (*models.JobRun, error)

	ResumeAllInProgress() error
	ResumeAllPendingNextBlock(currentBlockHeight *big.Int) error
	ResumeAllPendingConnection() error
}

// runManager implements RunManager
type runManager struct {
	orm         *orm.ORM
	statsPusher synchronization.StatsPusher
	runQueue    RunQueue
	txManager   store.TxManager
	config      orm.ConfigReader
	clock       utils.AfterNower
}

func runCost(job *models.JobSpec, config orm.ConfigReader, adapters []*adapters.PipelineAdapter) *assets.Link {
	minimumRunPayment := assets.NewLink(0)
	if job.MinPayment != nil {
		minimumRunPayment = job.MinPayment
		logger.Debugw("Using job's minimum payment", "required_payment", minimumRunPayment)
	} else if config.MinimumContractPayment() != nil {
		minimumRunPayment = config.MinimumContractPayment()
		logger.Debugw("Using configured minimum payment", "required_payment", minimumRunPayment)
	}

	for _, adapter := range adapters {
		minimumPayment := adapter.MinPayment()
		if minimumPayment != nil {
			minimumRunPayment = assets.NewLink(0).Add(minimumRunPayment, minimumPayment)
		}
	}

	return minimumRunPayment
}

// NewRun returns a complete run from a JobSpec
func NewRun(
	job *models.JobSpec,
	initiator *models.Initiator,
	currentHeight *big.Int,
	runRequest *models.RunRequest,
	config orm.ConfigReader,
	orm *orm.ORM,
	now time.Time) (*models.JobRun, []*adapters.PipelineAdapter) {

	run := models.MakeJobRun(job, now, initiator, currentHeight, runRequest)
	runAdapters := []*adapters.PipelineAdapter{}

	for i, task := range job.Tasks {
		adapter, err := adapters.For(task, config, orm)
		if err != nil {
			run.SetError(err)
			break
		}

		runAdapters = append(runAdapters, adapter)
		if currentHeight == nil {
			continue
		}

		run.TaskRuns[i].MinRequiredIncomingConfirmations = clnull.Uint32From(
			utils.MaxUint32(
				config.MinIncomingConfirmations(),
				task.MinRequiredIncomingConfirmations.Uint32,
				adapter.MinConfs()),
		)
	}

	return &run, runAdapters
}

// ValidateRun ensures that a run's initial preconditions have been met
func ValidateRun(run *models.JobRun, contractCost *assets.Link) {

	// payment is only present for runs triggered by runlogs
	if run.Payment != nil && contractCost.Cmp(run.Payment) > 0 {
		logger.Debugw("Rejecting run with insufficient payment",
			run.ForLogger("required_payment", contractCost.String())...)

		err := fmt.Errorf(
			"rejecting job %s with payment %s below minimum threshold (%s)",
			run.JobSpecID,
			run.Payment.Text(10),
			contractCost.Text(10))
		run.SetError(err)
		return
	}
}

// NewRunManager returns a new job manager
func NewRunManager(
	runQueue RunQueue,
	config orm.ConfigReader,
	orm *orm.ORM,
	statsPusher synchronization.StatsPusher,
	txManager store.TxManager,
	clock utils.AfterNower) RunManager {
	return &runManager{
		orm:         orm,
		statsPusher: statsPusher,
		runQueue:    runQueue,
		txManager:   txManager,
		config:      config,
		clock:       clock,
	}
}

// CreateErrored creates a run that is in the errored state. This is a
// special case where this job cannot run but we want to create the run record
// so the error is more visible to the node operator.
func (rm *runManager) CreateErrored(
	jobSpecID *models.ID,
	initiator models.Initiator,
	runErr error) (*models.JobRun, error) {
	job, err := rm.orm.Unscoped().FindJob(jobSpecID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find job spec")
	}

	now := time.Now()
	run := models.JobRun{
		ID:          models.NewID(),
		JobSpecID:   job.ID,
		CreatedAt:   now,
		UpdatedAt:   now,
		InitiatorID: initiator.ID,
	}

	run.SetError(runErr)
	defer rm.statsPusher.PushNow()
	return &run, rm.orm.CreateJobRun(&run)
}

// Create immediately persists a JobRun and sends it to the RunQueue for
// execution.
func (rm *runManager) Create(
	jobSpecID *models.ID,
	initiator *models.Initiator,
	creationHeight *big.Int,
	runRequest *models.RunRequest,
) (*models.JobRun, error) {
	logger.Debugw(fmt.Sprintf("New run triggered by %s", initiator.Type),
		"job", jobSpecID.String(),
		"creation_height", creationHeight.String(),
	)

	job, err := rm.orm.Unscoped().FindJob(jobSpecID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find job spec")
	}

	if job.Archived() {
		return nil, RecurringScheduleJobError{
			msg: fmt.Sprintf("Trying to run archived job %s", job.ID),
		}
	}

	now := rm.clock.Now()
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

	if len(job.Tasks) == 0 {
		return nil, fmt.Errorf("invariant for job %s: no tasks to run in NewRun", job.ID)
	}

	run, adapters := NewRun(&job, initiator, creationHeight, runRequest, rm.config, rm.orm, now)
	runCost := runCost(&job, rm.config, adapters)
	ValidateRun(run, runCost)

	if err := rm.orm.CreateJobRun(run); err != nil {
		return nil, errors.Wrap(err, "CreateJobRun failed")
	}
	rm.statsPusher.PushNow()

	if run.GetStatus().Runnable() {
		logger.Debugw(
			fmt.Sprintf("Executing run originally initiated by %s", run.Initiator.Type),
			run.ForLogger()...,
		)
		rm.runQueue.Run(run)
	}
	return run, nil
}

// ResumeAllPendingNextBlock wakes up all jobs that were sleeping because they
// were waiting for the next block
func (rm *runManager) ResumeAllPendingNextBlock(currentBlockHeight *big.Int) error {
	return rm.orm.UnscopedJobRunsWithStatus(func(run *models.JobRun) {
		currentTaskRun := run.NextTaskRun()
		if currentTaskRun == nil {
			err := rm.updateWithError(run, "Attempting to resume confirming run with no remaining tasks %s", run.ID)
			logger.ErrorIf(err, "failed when run manager updates with error")
			return
		}

		run.ObservedHeight = utils.NewBig(currentBlockHeight)
		logger.Debugw(fmt.Sprintf("New head #%s resuming run", currentBlockHeight), run.ForLogger()...)

		// Set jobRun status in progress if met minimum incoming confirmations
		// Task run status will be set later in runManager#Execute
		markInProgressIfSufficientIncomingConfirmations(run, currentTaskRun, run.ObservedHeight, rm.txManager)

		// Save job run and resume if status was set to InProgress
		err := rm.saveAndResumeIfInProgress(run)
		if err != nil {
			logger.Errorw("Error saving run", run.ForLogger("error", err)...)
		}
	},
		models.RunStatusPendingConnection,
		models.RunStatusPendingOutgoingConfirmations,
		models.RunStatusPendingIncomingConfirmations)
}

// ResumeAllPendingConnection wakes up all tasks that have gone to sleep because they
// needed an ethereum client connection.
func (rm *runManager) ResumeAllPendingConnection() error {
	return rm.orm.UnscopedJobRunsWithStatus(func(run *models.JobRun) {
		logger.Debugw("New connection resuming run", run.ForLogger()...)

		currentTaskRun := run.NextTaskRun()
		if currentTaskRun == nil {
			err := rm.updateWithError(run, "Attempting to resume connecting run with no remaining tasks %s", run.ID)
			logger.ErrorIf(err, "failed when run manager updates with error")
			return
		}

		run.SetStatus(models.RunStatusInProgress)

		err := rm.saveAndResumeIfInProgress(run)
		if err != nil {
			logger.Errorw("Error saving run", run.ForLogger("error", err)...)
		}
	},
		models.RunStatusPendingConnection, models.RunStatusPendingOutgoingConfirmations)
}

// ResumePendingBridgeTask wakes up a task that required a response from a bridge adapter.
func (rm *runManager) ResumePendingBridge(
	runID *models.ID,
	input models.BridgeRunResult,
) error {
	run, err := rm.orm.Unscoped().FindJobRun(runID)
	if err != nil {
		return err
	}

	logger.Debugw("External adapter resuming run", run.ForLogger("input_data", input.Data)...)

	if !run.GetStatus().PendingBridge() {
		return fmt.Errorf("attempting to resume non pending run %s", run.ID)
	}

	currentTaskRun := run.NextTaskRun()
	if currentTaskRun == nil {
		return rm.updateWithError(&run, "Attempting to resume pending run with no remaining tasks %s", run.ID)
	}

	data, err := models.Merge(run.RunRequest.RequestParams, input.Data)
	if err != nil {
		return rm.updateWithError(&run, "Error while merging onto RequestParams for run %s", run.ID)
	}
	run.RunRequest.RequestParams = data

	currentTaskRun.ApplyBridgeRunResult(input)
	run.ApplyBridgeRunResult(input)

	return rm.saveAndResumeIfInProgress(&run)
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
func (rm *runManager) ResumeAllInProgress() error {
	return rm.orm.UnscopedJobRunsWithStatus(rm.runQueue.Run, models.RunStatusInProgress, models.RunStatusPendingSleep)
}

// Cancel suspends a running task.
func (rm *runManager) Cancel(runID *models.ID) (*models.JobRun, error) {
	run, err := rm.orm.FindJobRun(runID)
	if err != nil {
		return nil, err
	}

	logger.Debugw("Cancelling run", run.ForLogger()...)
	if run.GetStatus().Finished() {
		return nil, fmt.Errorf("cannot cancel a run that has already finished")
	}

	run.Cancel()
	defer rm.statsPusher.PushNow()
	return &run, rm.orm.SaveJobRun(&run)
}

func (rm *runManager) updateWithError(run *models.JobRun, msg string, args ...interface{}) error {
	run.SetError(fmt.Errorf(msg, args...))
	logger.Error(fmt.Sprintf(msg, args...))

	if err := rm.orm.SaveJobRun(run); err != nil {
		logger.Errorw("Error saving run", run.ForLogger("error", err)...)
		return err
	}
	rm.statsPusher.PushNow()
	return nil
}

func (rm *runManager) saveAndResumeIfInProgress(run *models.JobRun) error {
	defer rm.statsPusher.PushNow()
	if err := rm.orm.SaveJobRun(run); err != nil {
		return err
	}
	rm.statsPusher.PushNow()
	if run.GetStatus() == models.RunStatusInProgress {
		rm.runQueue.Run(run)
	}
	return nil
}
