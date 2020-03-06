package services

import (
	"bytes"
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
		Help: "The current number of run queue workers",
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
	workers      map[models.ID]*singleJobSpecWorker // Access to this field is serialized through the select in orchestrateWorkers().
	numWorkers   int32                              // Note: every access to this field must be atomic, as it is shared between goroutines
	chRuns       chan *models.JobRun                // Run requests arrive from the Chainlink application via this channel
	chStop       chan struct{}                      // A shutdown request arrives from the Chainlink application via this channel
	chDone       chan struct{}                      // When the runQueue has finished shutting down, it sends a message on this channel to unblock Stop()
	chWorkerDone chan models.ID                     // When an individual job spec worker runs out of work, it shuts down and sends a message on this channel
}

// NewRunQueue initializes a RunQueue.
func NewRunQueue(runExecutor RunExecutor) RunQueue {
	return &runQueue{
		runExecutor:  runExecutor,
		workers:      make(map[models.ID]*singleJobSpecWorker),
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
			rq.handleIncomingRun(run)
		case jobSpecID := <-rq.chWorkerDone:
			rq.forgetDeadWorker(jobSpecID)
		case <-rq.chStop:
			rq.terminateAllWorkers()
			return
		}
	}
}

func (rq *runQueue) handleIncomingRun(run *models.JobRun) {
	worker, exists := rq.workers[*run.JobSpecID]
	if !exists {
		// If there's no worker for this job spec, start one.
		rq.startNewWorkerWithRun(run)
		return
	}

	select {
	case worker.chAddRun <- run:
		// If the worker is accepting work, just send the run.
	case <-worker.chDone:
		// If the worker is just spinning down as new work is coming in, allow it
		// to die and spin up a new one.
		rq.startNewWorkerWithRun(run)
	}
}

func (rq *runQueue) startNewWorkerWithRun(run *models.JobRun) {
	worker := newSingleJobSpecWorker(*run.JobSpecID, rq)
	rq.workers[*run.JobSpecID] = worker
	atomic.AddInt32(&rq.numWorkers, 1)
	worker.start()
	worker.chAddRun <- run
}

func (rq *runQueue) forgetDeadWorker(jobSpecID models.ID) {
	delete(rq.workers, jobSpecID)
	if n := atomic.AddInt32(&rq.numWorkers, -1); n < 0 {
		panic("numWorkers should never be < 0")
	}
}

func (rq *runQueue) terminateAllWorkers() {
	for _, worker := range rq.workers {
		worker.stop()
		rq.forgetDeadWorker(worker.jobSpecID)
	}
}

// WorkerCount returns the number of workers currently processing a job run
func (rq *runQueue) WorkerCount() int {
	return int(atomic.LoadInt32(&rq.numWorkers))
}

type singleJobSpecWorker struct {
	jobSpecID           models.ID
	runQueue            *runQueue
	chStartedProcessing chan struct{}
	chAddRun            chan *models.JobRun
	chBatch             chan []models.ID
	chBatchComplete     chan struct{}
	chStop              chan struct{}
	chDone              chan struct{}
	runsRequested       []models.ID // Access to this slice is serialized through the select in enqueueWork()
}

func newSingleJobSpecWorker(jobSpecID models.ID, runQueue *runQueue) *singleJobSpecWorker {
	return &singleJobSpecWorker{
		jobSpecID:           jobSpecID,
		runQueue:            runQueue,
		chStartedProcessing: make(chan struct{}),
		chAddRun:            make(chan *models.JobRun),
		chBatch:             make(chan []models.ID),
		chBatchComplete:     make(chan struct{}),
		chStop:              make(chan struct{}),
		chDone:              make(chan struct{}),
	}
}

func (worker *singleJobSpecWorker) start() {
	promNumberRunQueueWorkers.Inc()
	go worker.enqueueWork()
	go worker.processWork()
}

func (worker *singleJobSpecWorker) stop() {
	worker.chStop <- struct{}{}
	<-worker.chDone
	promNumberRunQueueWorkers.Dec()
}

func (worker *singleJobSpecWorker) enqueueWork() {
	defer close(worker.chDone)

	// The coordinator loop accepts "do job run" requests from orchestrateWorkers, "work finished"
	// signals from the worker goroutine, and "stop" messages from the Chainlink application,
	// coordinating them such that we can avoid races when messages arrive simultaneously.
	for {
		select {
		case run := <-worker.chAddRun:
			worker.enqueueRunAndResume(run)

		case <-worker.chBatchComplete:
			shutdown := worker.resumeOrShutdown()
			if shutdown {
				worker.stopProcessingWork()
				return
			}

		case <-worker.chStop:
			worker.stopProcessingWork()
			return
		}
	}
}

func (worker *singleJobSpecWorker) stopProcessingWork() {
	close(worker.chBatch)
}

func (worker *singleJobSpecWorker) enqueueRunAndResume(run *models.JobRun) {
	if !bytes.Equal(run.JobSpecID.Bytes(), worker.jobSpecID.Bytes()) {
		panic("worker received a run request from another job spec")
	}

	promNumberRunsQueued.Inc()
	worker.runsRequested = append(worker.runsRequested, *run.ID)

	select {
	case worker.chBatch <- worker.runsRequested:
		worker.runsRequested = nil
	case <-worker.chStartedProcessing:
		// If we couldn't send the worker a work request, make sure it's not because it
		// simply hasn't had a chance to start yet.  Avoids a race when .Stop() is called
		// too quickly.
	}
}

func (worker *singleJobSpecWorker) resumeOrShutdown() (shutdown bool) {
	if len(worker.runsRequested) > 0 {
		// If we've queued up more runs while the worker was working, keep it going.
		worker.chBatch <- worker.runsRequested
		worker.runsRequested = nil
		return false
	}

	// If we hit 0 runs, shut down the worker.
	select {
	case worker.runQueue.chWorkerDone <- worker.jobSpecID:
	case <-worker.chStop:
	}
	return true
}

// The processWork goroutine accepts job run requests from the enqueueWork goroutine.
// Once its queue of work is empty, it sends a "work finished" message to the enqueueWork
// goroutine, which may or may not allow it to spin down, depending on whether additional
// requests have been enqueued in the meantime.
func (worker *singleJobSpecWorker) processWork() {
	var startOnce sync.Once

	for runsRequested := range worker.chBatch {
		// We have to wait until processWork has actually received a job before indicating
		// to the enqueueWork goroutine that work has started.  If we do this before the worker's
		// for loop starts, there's a race that can allow enqueueWork to receive a chStop
		// message and shut down the worker even though runs are queued.
		startOnce.Do(func() { close(worker.chStartedProcessing) })

		for _, runID := range runsRequested {
			select {
			case <-worker.chStop:
				return
			default:
				worker.executeRun(runID)
			}
		}
		select {
		case worker.chBatchComplete <- struct{}{}:
		case <-worker.chStop:
			return
		}
	}
}

func (worker *singleJobSpecWorker) executeRun(runID models.ID) {
	if err := worker.runQueue.runExecutor.Execute(&runID); err != nil {
		logger.Errorw(fmt.Sprint("Error executing run ", runID.String()), "error", err)
	}
}
