package services

import (
	"errors"
	"fmt"
	"sync"

	"github.com/smartcontractkit/chainlink/adapters"
	"github.com/smartcontractkit/chainlink/logger"
	"github.com/smartcontractkit/chainlink/store"
	"github.com/smartcontractkit/chainlink/store/models"
)

// JobRunner safely handles coordinating job runs.
type JobRunner interface {
	Start() error
	Stop()
	resumeRunsSinceLastShutdown() error
	channelForRun(string) chan<- struct{}
	workerCount() int
}

type jobRunner struct {
	started              bool
	done                 chan struct{}
	bootMutex            sync.Mutex
	store                *store.Store
	workerMutex          sync.RWMutex
	workers              map[string]chan struct{}
	workersWg            sync.WaitGroup
	demultiplexStopperWg sync.WaitGroup
}

// NewJobRunner initializes a JobRunner.
func NewJobRunner(str *store.Store) JobRunner {
	return &jobRunner{
		store:   str,
		workers: make(map[string]chan struct{}),
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

	var starterWg sync.WaitGroup
	starterWg.Add(1)
	go rm.demultiplexRuns(&starterWg)
	starterWg.Wait()

	rm.demultiplexStopperWg.Add(1)
	return nil
}

// Stop closes all open worker channels.
func (rm *jobRunner) Stop() {
	rm.bootMutex.Lock()
	defer rm.bootMutex.Unlock()

	if !rm.started {
		return
	}
	close(rm.done)
	rm.started = false
	rm.demultiplexStopperWg.Wait()
}

// resumeRunsSinceLastShutdown queries the db for job runs that should be resumed
// since a previous node shutdown.
//
// As a result of its reliance on the database, it must run before anything
// persists a job RunStatus to the db to ensure that it only captures pending and in progress
// jobs as a result of the last shutdown, and not as a result of what's happening now.
//
// To recap: This must run before anything else writes job run status to the db,
// ie. tries to run a job.
// https://github.com/smartcontractkit/chainlink/pull/807
func (rm *jobRunner) resumeRunsSinceLastShutdown() error {
	sleepingRuns, err := rm.store.JobRunsWithStatus(models.RunStatusPendingSleep)
	if err != nil {
		return err
	}
	for _, run := range sleepingRuns {
		if err := QueueSleepingTask(&run, rm.store); err != nil {
			logger.Errorw("Error resuming sleeping job", "error", err)
		}
	}

	inProgressRuns, err := rm.store.JobRunsWithStatus(models.RunStatusInProgress)
	if err != nil {
		return err
	}
	for _, run := range inProgressRuns {
		rm.store.RunChannel.Send(run.ID)
	}
	return nil
}

func (rm *jobRunner) demultiplexRuns(starterWg *sync.WaitGroup) {
	starterWg.Done()
	defer rm.demultiplexStopperWg.Done()
	for {
		select {
		case <-rm.done:
			logger.Debug("JobRunner demultiplexing of job runs finished")
			rm.workersWg.Wait()
			return
		case rr, ok := <-rm.store.RunChannel.Receive():
			if !ok {
				logger.Panic("RunChannel closed before JobRunner, can no longer demultiplexing job runs")
				return
			}
			rm.channelForRun(rr.ID) <- struct{}{}
		}
	}
}

func (rm *jobRunner) channelForRun(runID string) chan<- struct{} {
	rm.workerMutex.Lock()
	defer rm.workerMutex.Unlock()

	workerChannel, present := rm.workers[runID]
	if !present {
		workerChannel = make(chan struct{}, 1)
		rm.workers[runID] = workerChannel
		rm.workersWg.Add(1)

		go func() {
			rm.workerLoop(runID, workerChannel)

			rm.workerMutex.Lock()
			delete(rm.workers, runID)
			rm.workersWg.Done()
			rm.workerMutex.Unlock()

			logger.Debug("Worker finished for ", runID)
		}()
	}
	return workerChannel
}

func (rm *jobRunner) workerLoop(runID string, workerChannel chan struct{}) {
	for {
		select {
		case <-workerChannel:
			run, err := rm.store.FindJobRun(runID)
			if err != nil {
				logger.Errorw(fmt.Sprint("Error finding run ", runID), run.ForLogger("error", err)...)
			}

			if err := executeRun(&run, rm.store); err != nil {
				logger.Errorw(fmt.Sprint("Error executing run ", runID), run.ForLogger("error", err)...)
				return
			}

			if run.Status.Finished() {
				logger.Debugw("All tasks complete for run", "run", run.ID)
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

func prepareTaskInput(run *models.JobRun, currentTaskRun *models.TaskRun) (models.RunResult, error) {
	input := currentTaskRun.Result
	previousTaskRun := run.PreviousTaskRun()

	if previousTaskRun != nil {
		input.Data = previousTaskRun.Result.Data.Merge(input.Data)
	}

	input.Data = run.Overrides.Data.Merge(input.Data)
	return input, nil
}

func executeTask(run *models.JobRun, currentTaskRun *models.TaskRun, store *store.Store) models.RunResult {
	taskCopy := currentTaskRun.TaskSpec // deliberately copied to keep mutations local
	taskCopy.Params = taskCopy.Params.Merge(run.Overrides.Data)

	adapter, err := adapters.For(taskCopy, store)
	if err != nil {
		currentTaskRun.Result.SetError(err)
		return currentTaskRun.Result
	}

	logger.Infow(fmt.Sprintf("Processing task %s", taskCopy.Type), []interface{}{"task", currentTaskRun.ID}...)

	input, err := prepareTaskInput(run, currentTaskRun)
	if err != nil {
		currentTaskRun.Result.SetError(err)
		return currentTaskRun.Result
	}

	result := adapter.Perform(input, store)

	logger.Infow(fmt.Sprintf("Finished processing task %s", taskCopy.Type), []interface{}{
		"task", currentTaskRun.ID,
		"result", result.Status,
		"result_data", result.Data,
	}...)

	return result
}

func executeRun(run *models.JobRun, store *store.Store) error {
	logger.Infow("Processing run", run.ForLogger()...)

	if !run.Status.Runnable() {
		return fmt.Errorf("Run triggered in non runnable state %s", run.Status)
	}

	currentTaskRun := run.NextTaskRun()
	if currentTaskRun == nil {
		return errors.New("Run triggered with no remaining tasks")
	}

	result := executeTask(run, currentTaskRun, store)

	currentTaskRun.ApplyResult(result)
	run.ApplyResult(result)

	if currentTaskRun.Status.PendingSleep() {
		logger.Debugw("Task is sleeping", []interface{}{"run", run.ID}...)
		if err := QueueSleepingTask(run, store); err != nil {
			return err
		}
	} else if !currentTaskRun.Status.Runnable() {
		logger.Debugw("Task execution blocked", []interface{}{"run", run.ID, "task", currentTaskRun.ID, "state", currentTaskRun.Result.Status}...)
	} else if futureTaskRun := run.NextTaskRun(); futureTaskRun != nil {
		if meetsMinimumConfirmations(run, futureTaskRun, run.ObservedHeight) {
			logger.Debugw("Adding next task to job run queue", []interface{}{"run", run.ID}...)
			run.Status = models.RunStatusInProgress
		} else {
			logger.Debugw("Blocking run pending incoming confirmations", []interface{}{"run", run.ID, "required_height", futureTaskRun.MinimumConfirmations}...)
			run.Status = models.RunStatusPendingConfirmations
		}
	}

	if err := updateAndTrigger(run, store); err != nil {
		return err
	}
	logger.Infow("Run finished processing", run.ForLogger()...)

	return nil
}
