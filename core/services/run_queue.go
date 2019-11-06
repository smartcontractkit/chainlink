package services

import (
	"fmt"
	"sync"
	"time"

	"chainlink/core/logger"
	"chainlink/core/store/models"
)

//go:generate mockery -name RunQueue -output ../internal/mocks/ -case=underscore

// RunQueue safely handles coordinating job runs.
type RunQueue interface {
	Start() error
	Stop()
	Run(*models.JobRun)

	WorkerCount() int
}

type runQueue struct {
	workersMutex sync.RWMutex
	workers      map[string]int
	workersWg    sync.WaitGroup

	runExecutor RunExecutor

	runsQueued   uint
	runsExecuted uint
	quit         chan struct{}
}

// NewRunQueue initializes a RunQueue.
func NewRunQueue(runExecutor RunExecutor) RunQueue {
	return &runQueue{
		quit:        make(chan struct{}),
		workers:     make(map[string]int),
		runExecutor: runExecutor,
	}
}

// Start prepares the job runner for accepting runs to execute.
func (rq *runQueue) Start() error {
	go rq.statisticsLogger()
	return nil
}

func (rq *runQueue) statisticsLogger() {
	ticker := time.NewTicker(5 * time.Minute)
	for {
		select {
		case <-ticker.C:
			rq.workersMutex.RLock()
			logger.Debugw("Run queue statistics", "runs_executed", rq.runsExecuted, "runs_queued", rq.runsQueued, "worker_count", len(rq.workers))
			rq.workersMutex.RUnlock()
		case <-rq.quit:
			ticker.Stop()
			return
		}
	}
}

// Stop closes all open worker channels.
func (rq *runQueue) Stop() {
	rq.quit <- struct{}{}
	rq.workersWg.Wait()
}

// Run tells the job runner to start executing a job
func (rq *runQueue) Run(run *models.JobRun) {
	runID := run.ID.String()

	rq.workersMutex.Lock()
	if queueCount, present := rq.workers[runID]; present {
		rq.runsQueued += 1
		rq.workers[runID] = queueCount + 1
		rq.workersMutex.Unlock()
		return
	}
	rq.runsExecuted += 1
	rq.workers[runID] = 1
	rq.workersMutex.Unlock()

	rq.workersWg.Add(1)
	go func() {
		for {
			rq.workersMutex.Lock()
			queueCount := rq.workers[runID]
			if queueCount <= 0 {
				delete(rq.workers, runID)
				rq.workersMutex.Unlock()
				break
			}
			rq.workers[runID] = queueCount - 1
			rq.workersMutex.Unlock()

			if err := rq.runExecutor.Execute(run.ID); err != nil {
				logger.Errorw(fmt.Sprint("Error executing run ", runID), "error", err)
			}
		}

		rq.workersWg.Done()
	}()
}

// WorkerCount returns the number of workers currently processing a job run
func (rq *runQueue) WorkerCount() int {
	rq.workersMutex.RLock()
	defer rq.workersMutex.RUnlock()

	return len(rq.workers)
}
