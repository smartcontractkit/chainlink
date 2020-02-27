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
	promNumberRunsQueued = promauto.NewCounter(prometheus.CounterOpts{
		Name: "run_queue_runs_queued",
		Help: "The total number of runs that have been queued",
	})
	promNumberRunQueueWorkers = promauto.NewGauge(prometheus.GaugeOpts{
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
	chStop   chan struct{}
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
				worker = &singleJobSpecWorker{make(chan struct{}), make(chan struct{}), make(chan struct{})}
				rq.workers[*run.ID] = worker
				atomic.AddInt32(&rq.numWorkers, 1)
				go rq.runSingleJobSpecWorker(run.ID, worker)
				worker.chAddRun <- struct{}{}
			} else {
				select {
				case worker.chAddRun <- struct{}{}:
				case <-worker.chDone:
					// If the worker is just spinning down as new work is coming in, allow it
					// to die and spin up a new one.
					worker = &singleJobSpecWorker{make(chan struct{}), make(chan struct{}), make(chan struct{})}
					rq.workers[*run.ID] = worker
					atomic.AddInt32(&rq.numWorkers, 1)
					go rq.runSingleJobSpecWorker(run.ID, worker)
					worker.chAddRun <- struct{}{}
				}
			}

		case jobID := <-rq.chWorkerDone:
			delete(rq.workers, jobID)
			atomic.AddInt32(&rq.numWorkers, -1)

		case <-rq.chStop:
			for _, run := range rq.workers {
				run.chStop <- struct{}{}
				<-run.chDone
				atomic.AddInt32(&rq.numWorkers, -1)
			}
			return
		}
	}
}

func (rq *runQueue) runSingleJobSpecWorker(runID *models.ID, worker *singleJobSpecWorker) {
	defer close(worker.chDone)

	promNumberRunQueueWorkers.Inc()

	var (
		startOnce       sync.Once
		chWorkerStarted = make(chan struct{})
		chResume        = make(chan int)
		chWorkComplete  = make(chan struct{})
	)

	// The worker goroutine accepts job run requests from the coordinator loop below.  It
	// can accept further requests after work has begun.  Once its queue of work is empty,
	// it sends a "work finished" message to that loop, which may or may not cause it to
	// spin down, depending on whether further requests are already enqueued.
	go func() {
		for numRunsRequested := range chResume {
			startOnce.Do(func() { close(chWorkerStarted) })

			for i := 0; i < numRunsRequested; i++ {
				select {
				case <-worker.chStop:
					return
				default:
				}

				promNumberRunsQueued.Inc()

				if err := rq.runExecutor.Execute(runID); err != nil {
					logger.Errorw(fmt.Sprint("Error executing run ", *runID), "error", err)
				}
			}
			select {
			case chWorkComplete <- struct{}{}:
			case <-worker.chStop:
				return
			}
		}
	}()

	// The coordinator loop accepts "do job run" requests from orchestrateWorkers, "work finished"
	// signals from the worker goroutine, and "stop" messages from the Chainlink application,
	// coordinating them such that we can avoid races when messages arrive simultaneously.
	numRunsRequested := 0
	for {
		select {
		case <-worker.chAddRun:
			numRunsRequested++

			select {
			case chResume <- numRunsRequested:
				numRunsRequested = 0
			case <-chWorkerStarted:
				// If we couldn't send the worker a work request, make sure it's not because it
				// simply hasn't had a chance to start yet.  Avoids a race when .Stop() is called
				// too quickly.
			}

		case <-chWorkComplete:
			if numRunsRequested > 0 {
				// If we've queued up more runs while the worker was working, keep it going.
				chResume <- numRunsRequested
				numRunsRequested = 0
			} else {
				// If we hit 0 runs, shut down the worker.
				select {
				case rq.chWorkerDone <- *runID:
				case <-worker.chStop:
				}
				close(chResume)
				return
			}

		case <-worker.chStop:
			close(chResume)
			return
		}
	}
}

// WorkerCount returns the number of workers currently processing a job run
func (rq *runQueue) WorkerCount() int {
	return int(atomic.LoadInt32(&rq.numWorkers))
}
