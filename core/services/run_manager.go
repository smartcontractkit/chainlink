package services

import (
	"fmt"
	"math/big"
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/smartcontractkit/chainlink/core/adapters"
	"github.com/smartcontractkit/chainlink/core/assets"
	"github.com/smartcontractkit/chainlink/core/logger"
	clnull "github.com/smartcontractkit/chainlink/core/null"
	"github.com/smartcontractkit/chainlink/core/services/eth"
	"github.com/smartcontractkit/chainlink/core/services/synchronization"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/store/orm"
	"github.com/smartcontractkit/chainlink/core/utils"
	"gorm.io/gorm"

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
		jobSpecID models.JobID,
		initiator *models.Initiator,
		creationHeight *big.Int,
		runRequest *models.RunRequest) (*models.JobRun, error)
	CreateErrored(
		jobSpecID models.JobID,
		initiator models.Initiator,
		err error) (*models.JobRun, error)
	ResumePendingBridge(
		runID uuid.UUID,
		input models.BridgeRunResult) error
	Cancel(runID uuid.UUID) (*models.JobRun, error)

	ResumeAllInProgress() error
	ResumeAllPendingNextBlock(currentBlockHeight *big.Int) error
	ResumeAllPendingConnection() error
}

// runManager implements RunManager
type runManager struct {
	orm         *orm.ORM
	statsPusher synchronization.StatsPusher
	runQueue    RunQueue
	config      orm.ConfigReader
	clock       utils.AfterNower
	ethClient   eth.Client
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
	ethClient eth.Client,
	now time.Time) (*models.JobRun, []*adapters.PipelineAdapter) {

	run := models.MakeJobRun(job, now, initiator, currentHeight, runRequest)
	runAdapters := []*adapters.PipelineAdapter{}

	for i, task := range job.Tasks {
		adapter, err := adapters.For(task, config, orm, ethClient)
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
	clock utils.AfterNower) RunManager {
	return &runManager{
		orm:         orm,
		statsPusher: statsPusher,
		runQueue:    runQueue,
		config:      config,
		clock:       clock,
	}
}

// CreateErrored creates a run that is in the errored state. This is a
// special case where this job cannot run but we want to create the run record
// so the error is more visible to the node operator.
func (rm *runManager) CreateErrored(
	jobSpecID models.JobID,
	initiator models.Initiator,
	runErr error) (*models.JobRun, error) {
	job, err := rm.orm.Unscoped().FindJobSpec(jobSpecID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find job spec")
	}

	now := time.Now()
	run := models.JobRun{
		ID:          uuid.NewV4(),
		JobSpecID:   job.ID,
		CreatedAt:   now,
		UpdatedAt:   now,
		InitiatorID: initiator.ID,
	}

	run.SetError(runErr)
	defer rm.statsPusher.PushNow()
	return &run, rm.orm.CreateJobRun(&run)
}

// If we have seen the same runRequest already, then double the required incoming confs.
func (rm *runManager) MaybeDoubleMinIncomingConfs(run *models.JobRun) error {
	if run == nil {
		return logger.NewErrorw("RunManager: expected non-nil run")
	}
	if run.RunRequest.RequestID == nil {
		return logger.NewErrorw("RunManager: expected non-nil run request ID")
	}
	if len(run.TaskRuns) == 0 {
		return logger.NewErrorw("RunManager: expected non-empty task runs")
	}
	// We want the maximum number of random task minimum_confirmations
	// of all job runs with the same run_request.request_id
	var maxConfs uint32
	if err := rm.orm.DB.Raw(`
SELECT coalesce(max(task_runs.minimum_confirmations), 0) FROM job_runs 
	JOIN run_requests ON job_runs.run_request_id = run_requests.id 
	JOIN task_runs ON job_runs.id = task_runs.job_run_id
	JOIN task_specs ON task_runs.task_spec_id = task_specs.id
	WHERE run_requests.request_id = ? AND task_specs.type = ?
`, run.RunRequest.RequestID, adapters.TaskTypeRandom).Scan(&maxConfs).Error; err != nil {
		return logger.NewErrorw("RunManager: unable to check for duplicate requests", "err", err)
	}
	if maxConfs != 0 {
		for i := range run.TaskRuns {
			if run.TaskRuns[i].TaskSpec.Type == adapters.TaskTypeRandom && run.TaskRuns[i].MinRequiredIncomingConfirmations.Valid {
				// Sanity cap must be less than 256 otherwise the blockhash will not be available directly
				// We give it 56 blocks as a buffer to still be able to fulfill within that bound.
				newConfs := maxConfs * 2
				if newConfs > 200 {
					newConfs = 200
				}
				logger.Warnw("RunManager: duplicate VRF requestID seen, doubling incoming confirmations",
					"requestID", run.RunRequest.RequestID,
					"txHash", run.RunRequest.TxHash,
					"oldConfs", run.TaskRuns[i].MinRequiredIncomingConfirmations.Uint32,
					"newConfs", newConfs)
				run.TaskRuns[i].MinRequiredIncomingConfirmations.Uint32 = newConfs
			}
		}
	}
	return nil
}

// Create immediately persists a JobRun and sends it to the RunQueue for
// execution.
func (rm *runManager) Create(
	jobSpecID models.JobID,
	initiator *models.Initiator,
	creationHeight *big.Int,
	runRequest *models.RunRequest,
) (*models.JobRun, error) {
	job, err := rm.orm.Unscoped().FindJobSpec(jobSpecID)
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

	run, adapters := NewRun(&job, initiator, creationHeight, runRequest, rm.config, rm.orm, rm.ethClient, now)
	runCost := runCost(&job, rm.config, adapters)
	ValidateRun(run, runCost)

	if initiator.Type == models.InitiatorRandomnessLog {
		err = rm.MaybeDoubleMinIncomingConfs(run)
		if err != nil {
			return nil, err
		}
	}

	logger.Debugw(
		fmt.Sprintf("RunManager: creating new job run initiated by %s", run.Initiator.Type),
		run.ForLogger()...,
	)
	if err := rm.orm.CreateJobRun(run); err != nil {
		return nil, errors.Wrap(err, "CreateJobRun failed")
	}
	rm.statsPusher.PushNow()

	if run.GetStatus().Runnable() {
		logger.Debugw(
			fmt.Sprintf("RunManager: executing run initiated by %s", run.Initiator.Type),
			run.ForLogger()...,
		)
		rm.runQueue.Run(run.ID)
	}
	return run, nil
}

// ResumeAllPendingNextBlock wakes up all jobs that were sleeping because they
// were waiting for the next block
func (rm *runManager) ResumeAllPendingNextBlock(currentBlockHeight *big.Int) error {
	logger.Debugw("Resuming all runs pending next block", "currentBlockHeight", currentBlockHeight)

	observedHeight := utils.NewBig(currentBlockHeight)
	resumableRunStatuses := []models.RunStatus{
		models.RunStatusPendingConnection,
		models.RunStatusPendingOutgoingConfirmations,
		models.RunStatusPendingIncomingConfirmations,
	}
	resumableTaskStatuses := []models.RunStatus{
		models.RunStatusUnstarted,
		models.RunStatusPendingConnection,
		models.RunStatusPendingOutgoingConfirmations,
		models.RunStatusPendingIncomingConfirmations,
	}
	runIDs := []uuid.UUID{}

	err := rm.orm.Transaction(func(tx *gorm.DB) error {
		updateTaskRunsQuery := `
UPDATE task_runs
   SET status = ?, confirmations = (
     CASE
		 WHEN job_runs.creation_height IS NULL OR task_runs.minimum_confirmations IS NULL OR ?::bigint IS NULL THEN
       NULL
     ELSE
       GREATEST(0, LEAST(task_runs.minimum_confirmations, (? - job_runs.creation_height) + 1))
     END
   )
  FROM job_runs
 WHERE job_runs.status IN (?)
   AND task_runs.job_run_id = job_runs.id
   AND task_runs.status IN (?);`
		result := tx.Exec(updateTaskRunsQuery, models.RunStatusInProgress, observedHeight, observedHeight, resumableRunStatuses, resumableTaskStatuses)
		if result.Error != nil {
			return result.Error
		}

		updateJobRunsQuery := `
   UPDATE job_runs
      SET status = ?, observed_height = ?
    WHERE status IN (?)
RETURNING id;`
		result = tx.Raw(updateJobRunsQuery, models.RunStatusInProgress, observedHeight, resumableRunStatuses)
		if result.Error != nil {
			return result.Error
		}
		return result.Scan(&runIDs).Error
	})

	if err != nil {
		return err
	}

	for _, runID := range runIDs {
		rm.runQueue.Run(runID)
	}
	rm.statsPusher.PushNow()
	return nil
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
	runID uuid.UUID,
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
	queueRun := func(run *models.JobRun) { rm.runQueue.Run(run.ID) }
	return rm.orm.UnscopedJobRunsWithStatus(queueRun, models.RunStatusInProgress, models.RunStatusPendingSleep)
}

// Cancel suspends a running task.
func (rm *runManager) Cancel(runID uuid.UUID) (*models.JobRun, error) {
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
		rm.runQueue.Run(run.ID)
	}
	return nil
}

// NullRunManager implements Null pattern for RunManager interface
type NullRunManager struct{}

func (NullRunManager) Create(jobSpecID models.JobID, initiator *models.Initiator, creationHeight *big.Int, runRequest *models.RunRequest) (*models.JobRun, error) {
	return nil, errors.New("NullRunManager#Create should never be called")
}

func (NullRunManager) CreateErrored(jobSpecID models.JobID, initiator models.Initiator, err error) (*models.JobRun, error) {
	return nil, errors.New("NullJobSubscriber#CreateErrored should never be called")
}

func (NullRunManager) ResumePendingBridge(runID uuid.UUID, input models.BridgeRunResult) error {
	return errors.New("NullRunManager#ResumePendingBridge should never be called")
}

func (NullRunManager) Cancel(runID uuid.UUID) (*models.JobRun, error) {
	return nil, errors.New("NullJobSubscriber#Cancel should never be called")
}

func (NullRunManager) ResumeAllInProgress() error {
	return nil
}

func (NullRunManager) ResumeAllPendingNextBlock(currentBlockHeight *big.Int) error {
	return nil
}

func (NullRunManager) ResumeAllPendingConnection() error {
	return nil
}
