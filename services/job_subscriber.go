package services

import (
	"sync"
	"time"

	"github.com/smartcontractkit/chainlink/logger"
	"github.com/smartcontractkit/chainlink/store"
	"github.com/smartcontractkit/chainlink/store/models"
	"github.com/smartcontractkit/chainlink/utils"
	"go.uber.org/multierr"
)

// JobSubscriber listens for push notifications from the ethereum node's
// websocket for specific jobs.
type JobSubscriber interface {
	HeadTrackable
	AddJob(job models.JobSpec, bn *models.IndexableBlockNumber) error
	Jobs() []models.JobSpec
	Stop()
	WorkerChannelFor(jr models.JobRun) chan *models.IndexableBlockNumber
}

// jobSubscriber implementation
type jobSubscriber struct {
	Store            *store.Store
	jobSubscriptions []JobSubscription
	jobsMutex        sync.Mutex
	workerMutex      sync.Mutex
	workerWaiter     sync.WaitGroup
	workers          map[string]chan *models.IndexableBlockNumber
}

// NewJobSubscriber returns a new job subscriber.
func NewJobSubscriber(store *store.Store) JobSubscriber {
	return &jobSubscriber{
		Store:   store,
		workers: make(map[string]chan *models.IndexableBlockNumber),
	}
}

// AddJob subscribes to ethereum log events for each "runlog" and "ethlog"
// initiator in the passed job spec.
func (js *jobSubscriber) AddJob(job models.JobSpec, bn *models.IndexableBlockNumber) error {
	if !job.IsLogInitiated() {
		return nil
	}

	sub, err := StartJobSubscription(job, bn, js.Store)
	if err != nil {
		return err
	}
	js.addSubscription(sub)
	return nil
}

// Jobs returns the jobs being listened to.
func (js *jobSubscriber) Jobs() []models.JobSpec {
	var jobs []models.JobSpec
	for _, js := range js.jobSubscriptions {
		jobs = append(jobs, js.Job)
	}
	return jobs
}

func (js *jobSubscriber) addSubscription(sub JobSubscription) {
	js.jobsMutex.Lock()
	defer js.jobsMutex.Unlock()
	js.jobSubscriptions = append(js.jobSubscriptions, sub)
}

// Connect connects the jobs to the ethereum node by creating corresponding subscriptions.
func (js *jobSubscriber) Connect(bn *models.IndexableBlockNumber) error {
	jobs, err := js.Store.Jobs()
	if err != nil {
		return err
	}
	for _, j := range jobs {
		err = multierr.Append(err, js.AddJob(j, bn))
	}
	return err
}

// Disconnect disconnects all subscriptions associated with jobs belonging to
// this listener.
func (js *jobSubscriber) Disconnect() {
	js.jobsMutex.Lock()
	defer js.jobsMutex.Unlock()
	for _, sub := range js.jobSubscriptions {
		sub.Unsubscribe()
	}
	js.jobSubscriptions = []JobSubscription{}
}

// OnNewHead resumes all pending job runs based on the new head activity.
func (js *jobSubscriber) OnNewHead(head *models.BlockHeader) {
	pendingRuns, err := js.Store.JobRunsWithStatus(models.RunStatusPendingConfirmations, models.RunStatusInProgress)
	if err != nil {
		logger.Error(err.Error())
	}

	activeJobRunIDs := make(map[string]struct{})

	js.workerMutex.Lock()
	defer js.workerMutex.Unlock()
	for _, jr := range pendingRuns {
		activeJobRunIDs[jr.ID] = struct{}{}

		workerChannel := js.WorkerChannelFor(jr)
		blockNumber := head.ToIndexableBlockNumber()
		workerChannel <- blockNumber
	}

	//Stop any workers that didn't have corresponding pending confirmations
	for id, workerChannel := range js.workers {
		if _, ok := activeJobRunIDs[id]; !ok {
			close(workerChannel)
			delete(js.workers, id)
		}
	}
}

// workerChannelFor accepts a JobRun and returns a worker channel dedicated
// to that JobRun. The channel accepts new block heights for triggering runs,
// and ensures that the block height confirmations are run syncronously.
func (js *jobSubscriber) WorkerChannelFor(jr models.JobRun) chan *models.IndexableBlockNumber {
	workerChannel, present := js.workers[jr.ID]
	if !present {
		workerChannel = make(chan *models.IndexableBlockNumber, 100)
		js.workers[jr.ID] = workerChannel

		go func() {
			js.workerWaiter.Add(1)
			defer js.workerWaiter.Done()

			for blockNumber := range workerChannel {
				if blockNumber == nil {
					logger.Debug("Stopped worker for", jr.ID)
					break
				}

				logger.Debug("Woke up", jr.ID, "worker to process", blockNumber.ToInt())
				if _, err := ExecuteRunAtBlock(jr, js.Store, jr.Result, blockNumber); err != nil {
					logger.Error(err.Error())
				}
			}
		}()
	}
	return workerChannel
}

// Stop closes all workers that have been started to process Job Runs on new
// heads and waits for them to finish.
func (js *jobSubscriber) Stop() {
	js.workerMutex.Lock()
	for _, workerChannel := range js.workers {
		workerChannel <- nil
	}
	js.workerMutex.Unlock()
	utils.WaitTimeout(&js.workerWaiter, 10*time.Second)
}
