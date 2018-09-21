package services

import (
	"errors"
	"fmt"
	"sync"

	"github.com/smartcontractkit/chainlink/adapters"
	"github.com/smartcontractkit/chainlink/logger"
	"github.com/smartcontractkit/chainlink/store"
	"github.com/smartcontractkit/chainlink/store/models"
	"github.com/smartcontractkit/chainlink/utils"
	"go.uber.org/multierr"
)

// JobRunner safely handles coordinating job runs.
type JobRunner interface {
	Start() error
	Stop()
	resumeSleepingRuns() error
	channelForRun(string) chan<- store.RunRequest
	workerCount() int
}

type jobRunner struct {
	started     bool
	done        chan struct{}
	bootMutex   sync.Mutex
	store       *store.Store
	workerMutex sync.RWMutex
	workers     map[string]chan store.RunRequest
	wg          sync.WaitGroup
}

// NewJobRunner initializes a JobRunner.
func NewJobRunner(str *store.Store) JobRunner {
	return &jobRunner{
		store:   str,
		workers: make(map[string]chan store.RunRequest),
	}
}

// Start reinitializes runs and starts the execution of the store's runs.
func (rm *jobRunner) Start() error {
	rm.bootMutex.Lock()
	defer rm.bootMutex.Unlock()

	if rm.started {
		return errors.New("JobRunner already started")
	}
	rm.done = make(chan struct{})
	rm.started = true
	go rm.demultiplexRuns()
	return rm.resumeSleepingRuns()
}

// Stop closes all open worker channels.
func (rm *jobRunner) Stop() {
	rm.bootMutex.Lock()
	defer rm.bootMutex.Unlock()

	if !rm.started {
		return
	}
	rm.started = false
	close(rm.done)
	rm.wg.Wait()
}

func (rm *jobRunner) resumeSleepingRuns() error {
	pendingRuns, err := rm.store.JobRunsWithStatus(models.RunStatusPendingSleep)
	if err != nil {
		return err
	}
	for _, run := range pendingRuns {
		rm.store.RunChannel.Send(run.Result, nil)
	}
	return nil
}

func (rm *jobRunner) demultiplexRuns() {
	for {
		select {
		case <-rm.done:
			logger.Debug("JobRunner demultiplexing of job runs finished")
			return
		case rr, ok := <-rm.store.RunChannel.Receive():
			if !ok {
				logger.Warn("JobRunner RunChannel closed, demultiplexing of job runs finished")
				return
			}
			rm.channelForRun(rr.Input.JobRunID) <- rr
		}
	}
}

func (rm *jobRunner) channelForRun(runID string) chan<- store.RunRequest {
	rm.workerMutex.Lock()
	defer rm.workerMutex.Unlock()

	workerChannel, present := rm.workers[runID]
	if !present {
		workerChannel = make(chan store.RunRequest, 1000)
		rm.workers[runID] = workerChannel

		rm.wg.Add(1)
		go func() {
			defer rm.wg.Done()
			rm.workerLoop(runID, workerChannel)

			rm.workerMutex.Lock()
			delete(rm.workers, runID)
			rm.workerMutex.Unlock()
			logger.Debug("Worker finished for ", runID)
		}()
	}
	return workerChannel
}

func (rm *jobRunner) workerLoop(runID string, workerChannel chan store.RunRequest) {
	for {
		select {
		case rr := <-workerChannel:
			jr, err := rm.store.FindJobRun(runID)
			if err != nil {
				logger.Errorw(fmt.Sprint("Application Run Channel Executor: error finding run ", runID), jr.ForLogger("error", err)...)
			}
			if rr.BlockNumber != nil {
				logger.Debug("Woke up", jr.ID, "worker to process ", rr.BlockNumber.ToInt())
			}
			if jr, err = ExecuteRunAtBlock(jr, rm.store, rr.Input, rr.BlockNumber); err != nil {
				logger.Errorw(fmt.Sprint("Application Run Channel Executor: error executing run ", runID), jr.ForLogger("error", err)...)
			}

			if jr.Status.Finished() {
				return
			}
		case <-rm.done:
			logger.Debug("JobRunner worker loop for ", runID, " finished")
			return
		}
	}
}

func (rm *jobRunner) workerCount() int {
	rm.workerMutex.RLock()
	defer rm.workerMutex.RUnlock()

	return len(rm.workers)
}

// BeginRun creates a new run if the job is valid and starts the job.
func BeginRun(
	job models.JobSpec,
	initr models.Initiator,
	input models.RunResult,
	store *store.Store,
) (models.JobRun, error) {
	return BeginRunAtBlock(job, initr, input, store, nil)
}

// BeginRunAtBlock builds and executes a new run if the job is valid with the block number
// to determine if tasks should be resumed.
func BeginRunAtBlock(
	job models.JobSpec,
	initr models.Initiator,
	input models.RunResult,
	store *store.Store,
	bn *models.IndexableBlockNumber,
) (models.JobRun, error) {
	run, err := BuildRun(job, initr, store)
	if err != nil {
		return models.JobRun{}, err
	}
	if input.Amount != nil &&
		store.Config.MinimumContractPayment.Cmp(input.Amount) > 0 {
		err := fmt.Errorf(
			"Rejecting job %s with payment %s below minimum threshold (%s)",
			job.ID,
			input.Amount,
			store.Config.MinimumContractPayment.Text(10))
		run = run.ApplyResult(input.WithError(err))
		return run, multierr.Append(err, store.Save(&run))
	}
	return ExecuteRunAtBlock(run, store, input, bn)
}

// BuildRun checks to ensure the given job has not started or ended before
// creating a new run for the job.
func BuildRun(job models.JobSpec, i models.Initiator, store *store.Store) (models.JobRun, error) {
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
	return job.NewRun(i), nil
}

// ExecuteRun calls ExecuteRunAtBlock without an IndexableBlockNumber
func ExecuteRun(jr models.JobRun, store *store.Store, overrides models.RunResult) (models.JobRun, error) {
	return ExecuteRunAtBlock(jr, store, overrides, nil)
}

// ExecuteRunAtBlock starts the job and executes task runs within that job in the
// order defined in the run for as long as they do not return errors. Results
// are saved in the store (db).
func ExecuteRunAtBlock(
	jr models.JobRun,
	store *store.Store,
	overrides models.RunResult,
	bn *models.IndexableBlockNumber,
) (models.JobRun, error) {
	jr, err := prepareJobRun(jr, store, overrides, bn)
	if err != nil {
		return jr, wrapError(jr, err)
	}
	logger.Infow("Starting job", jr.ForLogger()...)
	unfinished := jr.UnfinishedTaskRuns()
	if len(unfinished) == 0 {
		return jr, wrapError(jr, errors.New("No unfinished tasks to run"))
	}
	offset := len(jr.TaskRuns) - len(unfinished)
	prevResult, err := unfinished[0].Result.Merge(jr.Overrides)
	if err != nil {
		return jr, wrapError(jr, err)
	}

	for i, taskRunTemplate := range unfinished {
		nextTaskRun, err := taskRunTemplate.MergeTaskParams(jr.Overrides.Data)
		if err != nil {
			return jr, wrapError(jr, err)
		}

		lastRun := markCompletedIfRunnable(startTask(jr, nextTaskRun, prevResult, bn, store))
		jr.TaskRuns[i+offset] = lastRun
		logTaskResult(lastRun, nextTaskRun, i)
		prevResult = lastRun.Result

		if err := store.Save(&jr); err != nil {
			return jr, wrapError(jr, err)
		}
		if !lastRun.Status.Runnable() {
			break
		}
	}

	jr = jr.ApplyResult(prevResult)
	logger.Infow("Finished current job run execution", jr.ForLogger()...)
	return jr, wrapError(jr, store.Save(&jr))
}

func prepareJobRun(jr models.JobRun, store *store.Store, overrides models.RunResult, bn *models.IndexableBlockNumber) (models.JobRun, error) {
	if jr.Status.CanStart() {
		jr.Status = models.RunStatusInProgress
	} else {
		return jr, fmt.Errorf("Unable to start with status %v", jr.Status)
	}
	var err error
	jr.Overrides, err = jr.Overrides.Merge(overrides)
	if err != nil {
		jr = jr.ApplyResult(jr.Result.WithError(err))
		return jr, multierr.Append(err, store.Save(&jr))
	}
	if err = store.Save(&jr); err != nil {
		return jr, err
	}
	if jr.Result.HasError() {
		return jr, jr.Result
	}
	return store.SaveCreationHeight(jr, bn)
}

func logTaskResult(lr models.TaskRun, tr models.TaskRun, i int) {
	logger.Debugw("Produced task run", "taskRun", lr)
	logger.Debugw(fmt.Sprintf("Task %v %v", tr.Task.Type, tr.Result.Status), tr.ForLogger("task", i, "result", lr.Result)...)
}

func markCompletedIfRunnable(tr models.TaskRun) models.TaskRun {
	if tr.Status.Runnable() {
		return tr.MarkCompleted()
	}
	return tr
}

func startTask(
	jr models.JobRun,
	tr models.TaskRun,
	input models.RunResult,
	bn *models.IndexableBlockNumber,
	store *store.Store,
) models.TaskRun {
	adapter, err := adapters.For(tr.Task, store)
	if err != nil {
		return tr.ApplyResult(tr.Result.WithError(err))
	}

	minConfs := utils.MaxUint64(
		store.Config.MinIncomingConfirmations,
		tr.Task.Confirmations,
		adapter.MinConfs())

	if !jr.Runnable(bn, minConfs) {
		tr = tr.MarkPendingConfirmations()
		tr.Result.Data = input.Data
		return tr
	}

	return tr.ApplyResult(adapter.Perform(input, store))
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
