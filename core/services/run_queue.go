package services

import (
	"fmt"
	"sync"
	"sync/atomic"

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
	runExecutor  RunExecutor
	workers      map[models.ID]*singleJobSpecWorker
	numWorkers   int32
	chRuns       chan *models.JobRun
	chStop       chan struct{}
	chDone       chan struct{}
	chWorkerDone chan models.ID
}

type singleJobSpecWorker struct {
	chAddRun chan struct{}
	chDone   chan struct{}
}

// NewRunQueue initializes a RunQueue.
func NewRunQueue(runExecutor RunExecutor) RunQueue {
	return &runQueue{
		workers:      make(map[models.ID]*singleJobSpecWorker),
		runExecutor:  runExecutor,
		chRuns:       make(chan *models.JobRun),
		chStop:       make(chan struct{}),
		chDone:       make(chan struct{}),
		chWorkerDone: make(chan models.ID),
	}
}

// Start prepares the job runner for accepting runs to execute.
func (rq *runQueue) Start() error {
	go rq.orchestrateWorkers()
	return nil
}

// Stop closes all open worker channels.
func (rq *runQueue) Stop() {
	close(rq.chStop)
	<-rq.chDone
}

// Run tells the job runner to start executing a job
func (rq *runQueue) Run(run *models.JobRun) {
	select {
	case rq.chRuns <- run:
	case <-rq.chStop:
	}
}

func (rq *runQueue) orchestrateWorkers() {
	defer close(rq.chDone)
	for {
		select {
		case run := <-rq.chRuns:
			worker, exists := rq.workers[*run.ID]
			if !exists {
				worker = &singleJobSpecWorker{make(chan struct{}), make(chan struct{})}
				rq.workers[*run.ID] = worker
				atomic.AddInt32(&rq.numWorkers, 1)
				go rq.singleJobSpecWorkerLoop(run, worker)
			}
			worker.chAddRun <- struct{}{}

		case jobID := <-rq.chWorkerDone:
			delete(rq.workers, jobID)
			atomic.AddInt32(&rq.numWorkers, -1)

		case <-rq.chStop:
			for _, run := range rq.workers {
				<-run.chDone
			}
			return
		}
	}
}

func (rq *runQueue) singleJobSpecWorkerLoop(run *models.JobRun, worker *singleJobSpecWorker) {
	defer close(worker.chDone)

	var (
		startOnce       sync.Once
		chWorkerStarted = make(chan struct{})
		chResume        = make(chan int)
		chRunComplete   = make(chan struct{})
	)

	// The worker goroutine accepts job run requests from the CoordinateWorkAndDeath loop
	// below.  It can accept further requests after work has begun.  Once its queue of work
	// is empty, it sends a "die" request to that loop, which may or may not be permitted
	// depending on whether further requests are already inbound.
	go func() {
		for n := range chResume {
			startOnce.Do(func() { close(chWorkerStarted) })

			for i := 0; i < n; i++ {
				select {
				case <-rq.chStop:
					return
				default:
					if err := rq.runExecutor.Execute(run.ID); err != nil {
						logger.Errorw(fmt.Sprint("Error executing run ", *run.ID), "error", err)
					}
				}
			}
			select {
			case chRunComplete <- struct{}{}:
			case <-rq.chStop:
				return
			}
		}
	}()

	// The CoordinateWorkAndDeath loop accepts job run requests from orchestrateWorkers and
	// "die" requests from the worker, coordinating them such that we can avoid a race in the
	// edge case where both types of request arrive simultaneously.
	n := 0
CoordinateWorkAndDeath:
	for {
		select {
		case <-worker.chAddRun:
			select {
			case chResume <- n + 1:
				n = 0
			case <-chWorkerStarted:
				n++
			}
		case <-chRunComplete:
			if n > 0 {
				chResume <- n
				n = 0
			} else {
				select {
				case rq.chWorkerDone <- *run.ID:
					close(chResume)
					return
				default:
					continue CoordinateWorkAndDeath
				}
			}
		case <-rq.chStop:
			return
		}
	}
}

// WorkerCount returns the number of workers currently processing a job run
func (rq *runQueue) WorkerCount() int {
	return int(atomic.LoadInt32(&rq.numWorkers))
}
