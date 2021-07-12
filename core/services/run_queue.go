package services

import (
	"fmt"
	"sync"

	uuid "github.com/satori/go.uuid"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/service"

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

//go:generate mockery --name RunQueue --output ../internal/mocks/ --case=underscore

// RunQueue safely handles coordinating job runs.
type RunQueue interface {
	service.Service

	Run(uuid.UUID)

	WorkerCount() int
}

type runQueue struct {
	workersMutex  sync.RWMutex
	workers       map[string]int
	workersWg     sync.WaitGroup
	stopRequested bool

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

// Close closes all open worker channels.
func (rq *runQueue) Close() error {
	rq.workersMutex.Lock()
	rq.stopRequested = true
	rq.workersMutex.Unlock()
	rq.workersWg.Wait()
	return nil
}

func (rq *runQueue) Ready() error {
	return nil
}

func (rq *runQueue) Healthy() error {
	return nil
}

func (rq *runQueue) incrementQueue(runID string) bool {
	defer rq.workersMutex.Unlock()
	rq.workersMutex.Lock()
	numberRunsQueued.Inc()

	wasEmpty := rq.workers[runID] == 0
	rq.workers[runID]++
	numberRunQueueWorkers.Set(float64(len(rq.workers)))
	return wasEmpty
}

func (rq *runQueue) decrementQueue(runID string) bool {
	defer rq.workersMutex.Unlock()
	rq.workersMutex.Lock()

	rq.workers[runID]--
	isEmpty := rq.workers[runID] <= 0
	if isEmpty {
		delete(rq.workers, runID)
	}

	numberRunQueueWorkers.Set(float64(len(rq.workers)))
	return isEmpty
}

// Run tells the job runner to start executing a job
func (rq *runQueue) Run(runID uuid.UUID) {
	rq.workersMutex.Lock()
	if rq.stopRequested {
		rq.workersMutex.Unlock()
		return
	}
	rq.workersMutex.Unlock()

	id := runID.String()
	if !rq.incrementQueue(id) {
		return
	}

	rq.workersWg.Add(1)
	go func() {
		defer rq.workersWg.Done()

		for {
			if err := rq.runExecutor.Execute(runID); err != nil {
				logger.Errorw(fmt.Sprint("Error executing run ", id), "error", err)
			}

			if rq.decrementQueue(id) {
				return
			}
		}
	}()
}

// WorkerCount returns the number of workers currently processing a job run
func (rq *runQueue) WorkerCount() int {
	rq.workersMutex.RLock()
	defer rq.workersMutex.RUnlock()

	return len(rq.workers)
}

// NullRunQueue implements Null pattern for RunQueue interface
type NullRunQueue struct{}

func (NullRunQueue) Start() error   { return nil }
func (NullRunQueue) Close() error   { return nil }
func (NullRunQueue) Ready() error   { return nil }
func (NullRunQueue) Healthy() error { return nil }
func (NullRunQueue) Run(uuid.UUID) {
	panic("NullRunQueue#Run should never be called")
}
func (NullRunQueue) WorkerCount() int {
	panic("NullRunQueue#WorkerCount should never be called")
}
