package services

import (
	"fmt"
	"math/big"

	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/adapters"
	"github.com/smartcontractkit/chainlink/core/logger"
	clnull "github.com/smartcontractkit/chainlink/core/null"
	"github.com/smartcontractkit/chainlink/core/store"
	"github.com/smartcontractkit/chainlink/core/store/assets"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"
)

// ExecuteJob saves and immediately begins executing a run for a specified job
// if it is ready.
func ExecuteJob(
	job models.JobSpec,
	initiator models.Initiator,
	input models.RunResult,
	creationHeight *big.Int,
	store *store.Store) (*models.JobRun, error) {
	return ExecuteJobWithRunRequest(
		job,
		initiator,
		input,
		creationHeight,
		store,
		models.NewRunRequest(),
	)
}

// ExecuteJobWithRunRequest saves and immediately begins executing a run
// for a specified job if it is ready, assigning the passed initiator run.
func ExecuteJobWithRunRequest(
	job models.JobSpec,
	initiator models.Initiator,
	input models.RunResult,
	creationHeight *big.Int,
	store *store.Store,
	runRequest models.RunRequest) (*models.JobRun, error) {

	logger.Debugw(fmt.Sprintf("New run triggered by %s", initiator.Type),
		"job", job.ID,
		"input_status", input.Status,
		"creation_height", creationHeight,
	)

	run, err := NewRun(job, initiator, input, creationHeight, store, runRequest.Payment)
	if err != nil {
		return nil, errors.Wrap(err, "NewRun failed")
	}

	run.RunRequest = runRequest
	return run, createAndTrigger(run, store)
}

// MeetsMinimumPayment is a helper that returns true if jobrun received
// sufficient payment (more than jobspec's MinimumPayment) to be considered successful
func MeetsMinimumPayment(
	expectedMinJobPayment *assets.Link,
	actualRunPayment *assets.Link) bool {
	// input.Payment is always present for runs triggered by ethlogs
	if actualRunPayment == nil || expectedMinJobPayment == nil || expectedMinJobPayment.IsZero() {
		return true
	}
	return expectedMinJobPayment.Cmp(actualRunPayment) < 1
}

// NewRun returns a run from an input job, in an initial state ready for
// processing by the job runner system
func NewRun(
	job models.JobSpec,
	initiator models.Initiator,
	input models.RunResult,
	currentHeight *big.Int,
	store *store.Store,
	payment *assets.Link) (*models.JobRun, error) {

	now := store.Clock.Now()
	if !job.Started(now) {
		return nil, RecurringScheduleJobError{
			msg: fmt.Sprintf("Job runner: Job %s unstarted: %v before job's start time %v", job.ID.String(), now, job.EndAt),
		}
	}

	if job.Ended(now) {
		return nil, RecurringScheduleJobError{
			msg: fmt.Sprintf("Job runner: Job %s ended: %v past job's end time %v", job.ID.String(), now, job.EndAt),
		}
	}

	run := job.NewRun(initiator)

	run.Overrides = input.Data
	run.ApplyResult(input)
	run.CreationHeight = models.NewBig(currentHeight)
	run.ObservedHeight = models.NewBig(currentHeight)

	if !MeetsMinimumPayment(job.MinPayment, payment) {
		logger.Infow("Rejecting run with insufficient payment", []interface{}{
			"run", run.ID,
			"job", run.JobSpecID,
			"input_payment", payment,
			"required_payment", job.MinPayment,
		}...)

		err := fmt.Errorf(
			"Rejecting job %s with payment %s below job-specific-minimum threshold (%s)",
			job.ID,
			payment,
			job.MinPayment.Text(10))
		run.SetError(err)
	}

	cost := assets.NewLink(0)
	for i, taskRun := range run.TaskRuns {
		adapter, err := adapters.For(taskRun.TaskSpec, store)

		if err != nil {
			run.SetError(err)
			return &run, nil
		}

		mp := adapter.MinContractPayment()
		if mp != nil {
			cost.Add(cost, mp)
		}

		if currentHeight != nil {
			run.TaskRuns[i].MinimumConfirmations = clnull.Uint32From(
				utils.MaxUint32(
					store.Config.MinIncomingConfirmations(),
					taskRun.TaskSpec.Confirmations.Uint32,
					adapter.MinConfs()),
			)
		}
	}

	// payment is always present for runs triggered by ethlogs
	if payment != nil {
		if cost.Cmp(payment) > 0 {
			logger.Debugw("Rejecting run with insufficient payment", []interface{}{
				"run", run.ID,
				"job", run.JobSpecID,
				"input_payment", payment,
				"required_payment", cost,
			}...)

			err := fmt.Errorf(
				"Rejecting job %s with payment %s below minimum threshold (%s)",
				job.ID,
				payment,
				store.Config.MinimumContractPayment().Text(10))
			run.SetError(err)
		}
	}

	if len(run.TaskRuns) == 0 {
		run.SetError(fmt.Errorf("invariant for job %s: no tasks to run in NewRun", job.ID.String()))
	}

	if !run.Status.Runnable() {
		return &run, nil
	}

	initialTask := run.TaskRuns[0]
	validateMinimumConfirmations(&run, &initialTask, run.CreationHeight, store)
	return &run, nil
}

// ResumeConfirmingTask resumes a confirming run if the minimum confirmations have been met
func ResumeConfirmingTask(
	run *models.JobRun,
	store *store.Store,
	currentBlockHeight *big.Int,
) error {

	logger.Debugw("New head resuming run", run.ForLogger()...)

	if !run.Status.PendingConfirmations() && !run.Status.PendingConnection() {
		return fmt.Errorf("Attempt to resume non confirming task")
	}

	currentTaskRun := run.NextTaskRun()
	if currentTaskRun == nil {
		return fmt.Errorf("Attempting to resume confirming run with no remaining tasks %s", run.ID.String())
	}

	if currentBlockHeight == nil {
		return fmt.Errorf("Attempting to resume confirming run with no currentBlockHeight %s", run.ID.String())
	}

	run.ObservedHeight = models.NewBig(currentBlockHeight)

	validateMinimumConfirmations(run, currentTaskRun, run.ObservedHeight, store)
	return updateAndTrigger(run, store)
}

// ResumeConnectingTask resumes a run that was left in pending_connection.
func ResumeConnectingTask(
	run *models.JobRun,
	store *store.Store,
) error {

	logger.Debugw("New connection resuming run", run.ForLogger()...)

	if !run.Status.PendingConnection() {
		return fmt.Errorf("Attempt to resume non connecting task")
	}

	currentTaskRun := run.NextTaskRun()
	if currentTaskRun == nil {
		return fmt.Errorf("Attempting to resume connecting run with no remaining tasks %s", run.ID.String())
	}

	run.Status = models.RunStatusInProgress
	return updateAndTrigger(run, store)
}

// ResumePendingTask takes the body provided from an external adapter,
// saves it for the next task to process, then tells the job runner to execute
// it
func ResumePendingTask(
	run *models.JobRun,
	store *store.Store,
	input models.RunResult,
) error {
	logger.Debugw("External adapter resuming job", []interface{}{
		"run", run.ID.String(),
		"job", run.JobSpecID.String(),
		"status", run.Status,
		"input_data", input.Data,
		"input_result", input.Status,
	}...)

	if !run.Status.PendingBridge() {
		return fmt.Errorf("Attempting to resume non pending run %s", run.ID.String())
	}

	currentTaskRun := run.NextTaskRun()
	if currentTaskRun == nil {
		return fmt.Errorf("Attempting to resume pending run with no remaining tasks %s", run.ID.String())
	}

	run.Overrides.Merge(input.Data)

	currentTaskRun.ApplyResult(input)
	if currentTaskRun.Status.Finished() && run.TasksRemain() {
		run.Status = models.RunStatusInProgress
	} else if currentTaskRun.Status.Finished() {
		run.ApplyResult(input)
		run.SetFinishedAt()
	} else {
		run.ApplyResult(input)
	}

	return updateAndTrigger(run, store)
}

func prepareAdapter(
	taskRun *models.TaskRun,
	data models.JSON,
	store *store.Store,
) (*adapters.PipelineAdapter, error) {
	taskCopy := taskRun.TaskSpec // deliberately copied to keep mutations local

	merged, err := taskCopy.Params.Merge(data)
	if err != nil {
		return nil, err
	}
	taskCopy.Params = merged

	return adapters.For(taskCopy, store)
}

// QueueSleepingTask creates a go routine which will wake up the job runner
// once the sleep's time has elapsed
func QueueSleepingTask(
	run *models.JobRun,
	store *store.Store,
) error {
	if !run.Status.PendingSleep() {
		return fmt.Errorf("Attempting to resume non sleeping run %s", run.ID.String())
	}

	currentTaskRun := run.NextTaskRun()
	if currentTaskRun == nil {
		return fmt.Errorf("Attempting to resume sleeping run with no remaining tasks %s", run.ID.String())
	}

	if !currentTaskRun.Status.PendingSleep() {
		return fmt.Errorf("Attempting to resume sleeping run with non sleeping task %s", run.ID.String())
	}

	adapter, err := prepareAdapter(currentTaskRun, run.Overrides, store)
	if err != nil {
		currentTaskRun.SetError(err)
		run.SetError(err)
		return store.SaveJobRun(run)
	}

	if sleepAdapter, ok := adapter.BaseAdapter.(*adapters.Sleep); ok {
		return performTaskSleep(run, currentTaskRun, sleepAdapter, store)
	}

	return fmt.Errorf("Attempting to resume non sleeping task for run %s (%s)", run.ID.String(), currentTaskRun.TaskSpec.Type)
}

func performTaskSleep(
	run *models.JobRun,
	task *models.TaskRun,
	adapter *adapters.Sleep,
	store *store.Store) error {

	duration := adapter.Duration()
	if duration <= 0 {
		logger.Debugw("Sleep duration has already elapsed, completing task", run.ForLogger()...)
		task.Status = models.RunStatusCompleted
		run.Status = models.RunStatusInProgress
		return updateAndTrigger(run, store)
	}

	// XXX: This is to eliminate data race that occurs because slices share their
	// underlying array even in copies
	runCopy := *run
	runCopy.TaskRuns = make([]models.TaskRun, len(run.TaskRuns))
	copy(runCopy.TaskRuns, run.TaskRuns)

	go func(run models.JobRun) {
		logger.Debugw("Task sleeping...", run.ForLogger()...)

		<-store.Clock.After(duration)

		task := run.NextTaskRun()
		task.Status = models.RunStatusCompleted
		run.Status = models.RunStatusInProgress

		logger.Debugw("Waking job up after sleep", run.ForLogger()...)

		if err := updateAndTrigger(&run, store); err != nil {
			logger.Errorw("Error resuming sleeping job:", "error", err)
		}
	}(runCopy)

	return nil
}

func validateMinimumConfirmations(
	run *models.JobRun,
	taskRun *models.TaskRun,
	currentHeight *models.Big,
	store *store.Store) {

	updateTaskRunConfirmations(currentHeight, run, taskRun)
	if !meetsMinimumConfirmations(run, taskRun, run.ObservedHeight) {
		logger.Debugw("Run cannot continue because it lacks sufficient confirmations", []interface{}{"run", run.ID.String(), "required_height", taskRun.MinimumConfirmations}...)
		run.Status = models.RunStatusPendingConfirmations
	} else if err := validateOnMainChain(run, taskRun, store); err != nil {
		run.SetError(err)
	} else {
		logger.Debugw("Adding next task to job run queue", []interface{}{"run", run.ID.String(), "nextTask", taskRun.TaskSpec.Type}...)
		run.Status = models.RunStatusInProgress
	}
}

func updateTaskRunConfirmations(currentHeight *models.Big, jr *models.JobRun, taskRun *models.TaskRun) {
	if !taskRun.MinimumConfirmations.Valid || jr.CreationHeight == nil || currentHeight == nil {
		return
	}

	confs := blockConfirmations(currentHeight, jr.CreationHeight)
	diff := utils.MinBigs(confs, big.NewInt(int64(taskRun.MinimumConfirmations.Uint32)))

	// diff's ceiling is guaranteed to be MaxUint32 since MinimumConfirmations
	// ceiling is MaxUint32.
	taskRun.Confirmations = clnull.Uint32From(uint32(diff.Int64()))
}

func validateOnMainChain(jr *models.JobRun, taskRun *models.TaskRun, store *store.Store) error {
	txhash := jr.RunRequest.TxHash
	if txhash == nil || !taskRun.MinimumConfirmations.Valid || taskRun.MinimumConfirmations.Uint32 == 0 {
		return nil
	}

	receipt, err := store.TxManager.GetTxReceipt(*txhash)
	if err != nil {
		return err
	}
	if invalidRequest(jr.RunRequest, receipt) {
		return fmt.Errorf(
			"TxHash %s initiating run %s not on main chain; presumably has been uncled",
			txhash.Hex(),
			jr.ID.String(),
		)
	}
	return nil
}

func invalidRequest(request models.RunRequest, receipt *models.TxReceipt) bool {
	return receipt.Unconfirmed() ||
		(request.BlockHash != nil && *request.BlockHash != *receipt.BlockHash)
}

func meetsMinimumConfirmations(
	run *models.JobRun,
	taskRun *models.TaskRun,
	currentHeight *models.Big) bool {
	if !taskRun.MinimumConfirmations.Valid || run.CreationHeight == nil || currentHeight == nil {
		return true
	}

	diff := blockConfirmations(currentHeight, run.CreationHeight)
	return diff.Cmp(big.NewInt(int64(taskRun.MinimumConfirmations.Uint32))) >= 0
}

func blockConfirmations(currentHeight, creationHeight *models.Big) *big.Int {
	bigDiff := new(big.Int).Sub(currentHeight.ToInt(), creationHeight.ToInt())
	confs := bigDiff.Add(bigDiff, big.NewInt(1)) // creation of runlog alone warrants 1 confirmation
	if confs.Cmp(big.NewInt(0)) < 0 {            // negative, so floor at 0
		confs.SetUint64(0)
	}
	return confs
}

func updateAndTrigger(run *models.JobRun, store *store.Store) error {
	if err := store.SaveJobRun(run); err != nil {
		return err
	}
	return triggerIfReady(run, store)
}

func createAndTrigger(run *models.JobRun, store *store.Store) error {
	if err := store.CreateJobRun(run); err != nil {
		return errors.Wrap(err, "CreateJobRun failed")
	}
	return triggerIfReady(run, store)
}

func triggerIfReady(run *models.JobRun, store *store.Store) error {
	if run.Status == models.RunStatusInProgress {
		logger.Debugw(fmt.Sprintf("Executing run originally initiated by %s", run.Initiator.Type), run.ForLogger()...)
		return store.RunChannel.Send(run.ID)
	}

	logger.Debugw(fmt.Sprintf("Pausing run originally initiated by %s", run.Initiator.Type), run.ForLogger()...)
	return nil
}

// RecurringScheduleJobError contains the field for the error message.
type RecurringScheduleJobError struct {
	msg string
}

// Error returns the error message for the run.
func (err RecurringScheduleJobError) Error() string {
	return err.msg
}
