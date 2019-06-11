package services

import (
	"errors"
	"fmt"
	"sync"

	"github.com/smartcontractkit/chainlink/core/adapters"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/store"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"go.uber.org/multierr"
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
		// Unscoped allows the processing of runs that are soft deleted asynchronously
		store:   str.Unscoped(),
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
	// Do all querying of run statuses since last shutdown before enqueuing
	// runs in progress and asleep, to prevent the following race condition:
	// 1. resume sleep, 2. awake from sleep, 3. in progress, 4. resume in progress (double enqueued).
	var merr error
	err := rm.store.UnscopedJobRunsWithStatus(func(run *models.JobRun) {

		if run.Result.Status == models.RunStatusPendingSleep {
			if err := QueueSleepingTask(run, rm.store.Unscoped()); err != nil {
				logger.Errorw("Error resuming sleeping job", "error", err)
			}
		} else {
			merr = multierr.Append(merr, rm.store.RunChannel.Send(run.ID))
		}

	}, models.RunStatusInProgress, models.RunStatusPendingSleep)

	if err != nil {
		return err
	}

	return merr
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

func prepareTaskInput(run *models.JobRun, input models.JSON) (models.JSON, error) {
	previousTaskRun := run.PreviousTaskRun()

	var err error
	if previousTaskRun != nil {
		if input, err = previousTaskRun.Result.Data.Merge(input); err != nil {
			return models.JSON{}, err
		}
	}

	if input, err = run.Overrides.Data.Merge(input); err != nil {
		return models.JSON{}, err
	}
	return input, nil
}

func executeTask(run *models.JobRun, currentTaskRun *models.TaskRun, store *store.Store) models.RunResult {
	taskCopy := currentTaskRun.TaskSpec // deliberately copied to keep mutations local

	var err error
	if taskCopy.Params, err = taskCopy.Params.Merge(run.Overrides.Data); err != nil {
		currentTaskRun.Result.SetError(err)
		return currentTaskRun.Result
	}

	adapter, err := adapters.For(taskCopy, store)
	if err != nil {
		currentTaskRun.Result.SetError(err)
		return currentTaskRun.Result
	}

	logger.Infow(fmt.Sprintf("Processing task %s", taskCopy.Type), []interface{}{"task", currentTaskRun.ID}...)

	data, err := prepareTaskInput(run, currentTaskRun.Result.Data)
	if err != nil {
		currentTaskRun.Result.SetError(err)
		return currentTaskRun.Result
	}

	currentTaskRun.Result.Data = data
	result := adapter.Perform(currentTaskRun.Result, store)

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
	} else if currentTaskRun.Status.Unstarted() {
		return fmt.Errorf("run %s task %s cannot return a status of empty string or Unstarted", run.ID, currentTaskRun.TaskSpec.Type)
	} else if futureTaskRun := run.NextTaskRun(); futureTaskRun != nil {
		validateMinimumConfirmations(run, futureTaskRun, run.ObservedHeight, store)
	}

	if run.Status.Finished() {
		run.SetFinishedAt()
	}

	if err := updateAndTrigger(run, store); err != nil {
		return err
	}
	logger.Infow("Run finished processing", run.ForLogger()...)

	return nil
}
