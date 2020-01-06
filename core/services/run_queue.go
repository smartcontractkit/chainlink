package services

import (
	"fmt"
	"sync"

	"chainlink/core/logger"
	"chainlink/core/store/models"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	numberRunsQueued = promauto.NewCounter(prometheus.CounterOpts{
		Name: "run_queue_runs_queued",
		Help: "The total number of runs that have been queued",
	})
	numberRunQueueWorkers = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "run_queue_queue_size",
		Help: "The size of the run queue",
	})
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
}

// NewRunQueue initializes a RunQueue.
func NewRunQueue(runExecutor RunExecutor) RunQueue {
	return &runQueue{
		workers:     make(map[string]int),
		runExecutor: runExecutor,
	}
}

// Start prepares the job runner for accepting runs to execute.
func (rq *runQueue) Start() error {
	return nil
}

// Stop closes all open worker channels.
func (rq *runQueue) Stop() {
	rq.workersWg.Wait()
}

// Run tells the job runner to start executing a job
func (rq *runQueue) Run(run *models.JobRun) {
	runID := run.ID.String()

	defer numberRunsQueued.Inc()

	rq.workersMutex.Lock()
	if queueCount, present := rq.workers[runID]; present {
		rq.workers[runID] = queueCount + 1
		rq.workersMutex.Unlock()
		return
	}
	rq.workers[runID] = 1
	numberRunQueueWorkers.Set(float64(len(rq.workers)))
	rq.workersMutex.Unlock()

	rq.workersWg.Add(1)
	go func() {
		for {
			rq.workersMutex.Lock()
			queueCount := rq.workers[runID]
			if queueCount <= 0 {
				delete(rq.workers, runID)
				numberRunQueueWorkers.Set(float64(len(rq.workers)))
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
