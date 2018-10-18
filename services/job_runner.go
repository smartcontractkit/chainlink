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
	resumeRuns() error
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
	return rm.resumeRuns()
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

func (rm *jobRunner) resumeRuns() error {
	sleepingRuns, err := rm.store.JobRunsWithStatus(models.RunStatusPendingSleep)
	if err != nil {
		return err
	}
	for _, run := range sleepingRuns {
		if _, err := QueueSleepingTask(&run, rm.store); err != nil {
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

			if run, err := executeRun(&run, rm.store); err != nil {
				logger.Errorw(fmt.Sprint("Error executing run ", runID), run.ForLogger("error", err)...)
				return
			}

			if run.Status.Finished() {
				logger.Debugw("All tasks complete for run", []interface{}{"run", run.ID}...)
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

func executeTask(run *models.JobRun, currentTaskRun *models.TaskRun, store *store.Store) models.RunResult {
	var err error
	if currentTaskRun.Task.Params, err = currentTaskRun.Task.Params.Merge(run.Overrides.Data); err != nil {
		return currentTaskRun.Result.WithError(err)
	}

	input := models.RunResult{Data: models.JSON{}}
	previousTaskRun := run.PreviousTaskRun()
	if previousTaskRun != nil {
		input = previousTaskRun.Result
	}
	if input.Data, err = run.Overrides.Data.Merge(input.Data); err != nil {
		return currentTaskRun.Result.WithError(err)
	}
	input.JobRunID = run.ID

	adapter, err := adapters.For(currentTaskRun.Task, store)
	if err != nil {
		return currentTaskRun.Result.WithError(err)
	}

	logger.Infow(fmt.Sprintf("Processing task %s", currentTaskRun.Task.Type), []interface{}{"task", currentTaskRun.ID}...)

	result := adapter.Perform(input, store)

	logger.Infow(fmt.Sprintf("Finished processing task %s", currentTaskRun.Task.Type), []interface{}{
		"task", currentTaskRun.ID,
		"result", result.Status,
	}...)

	return result
}

func executeRun(run *models.JobRun, store *store.Store) (*models.JobRun, error) {
	logger.Infow("Processing run", run.ForLogger()...)

	if !run.Status.Runnable() {
		return run, fmt.Errorf("Run triggered in non runnable state %s", run.Status)
	}

	if !run.TasksRemain() {
		return run, errors.New("Run triggered with no remaining tasks")
	}

	currentTaskRunIndex, _ := run.NextTaskRunIndex()
	currentTaskRun := run.TaskRuns[currentTaskRunIndex]

	result := executeTask(run, &currentTaskRun, store)

	currentTaskRun = currentTaskRun.ApplyResult(result)
	run.TaskRuns[currentTaskRunIndex] = currentTaskRun
	if currentTaskRun.Status.PendingSleep() {
		run.TaskRuns[currentTaskRunIndex] = currentTaskRun
		*run = run.ApplyResult(result)

		logger.Debugw("Task is sleeping", []interface{}{"run", run.ID}...)
		if run, err := QueueSleepingTask(run, store); err != nil {
			return run, err
		}
	} else if !currentTaskRun.Status.Runnable() {
		logger.Debugw("Task execution blocked", []interface{}{"run", run.ID, "task", currentTaskRun.ID, "state", currentTaskRun.Result.Status}...)
		*run = run.ApplyResult(result)
	} else if !run.TasksRemain() {
		logger.Debugw("All tasks completed, marking run complete", []interface{}{"run", run.ID, "task", currentTaskRun.ID}...)
		*run = run.ApplyResult(currentTaskRun.Result)
		*run = run.MarkCompleted()
	}

	if !run.Status.Finished() && run.Status.Runnable() {
		run = queueNextTask(run, store)
	}

	if err := saveAndTrigger(run, store); err != nil {
		return run, err
	}
	logger.Infow("Run finished processing", run.ForLogger()...)

	return run, nil
}

func queueNextTask(run *models.JobRun, store *store.Store) *models.JobRun {
	currentTaskRunIndex, _ := run.NextTaskRunIndex()
	futureTaskRun := run.TaskRuns[currentTaskRunIndex]

	if meetsMinimumConfirmations(run, &futureTaskRun, run.ObservedHeight) {
		logger.Debugw("Adding next task to job run queue", []interface{}{"run", run.ID}...)
		run.Status = models.RunStatusInProgress
	} else {
		logger.Debugw("Blocking run pending incoming confirmations", []interface{}{"run", run.ID, "required_height", futureTaskRun.MinimumConfirmations}...)
		run.Status = models.RunStatusPendingConfirmations
	}

	run.TaskRuns[currentTaskRunIndex] = futureTaskRun
	return run
}
