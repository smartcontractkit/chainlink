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
	AddJob(job models.JobSpec, bn *models.IndexableBlockNumber) error
	Connect(bn *models.IndexableBlockNumber) error
	Disconnect()
	Jobs() []models.JobSpec
	OnNewHead(head *models.BlockHeader)
	Stop()
}

// jobSubscriber implementation
type jobSubscriber struct {
	Store            *store.Store
	jobSubscriptions []JobSubscription
	jobsMutex        sync.Mutex
	workerMutex      sync.Mutex
	workerWaiter     sync.WaitGroup
}

// NewJobSubscriber returns a new job subscriber.
func NewJobSubscriber(store *store.Store) JobSubscriber {
	return &jobSubscriber{
		Store: store,
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

		workerChannel := js.Store.RunManager.WorkerChannelFor(jr.ID)
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
		blockNumber := head.ToIndexableBlockNumber()
		workerChannel <- blockNumber
	}

	//Stop any workers that didn't have corresponding pending confirmations
	for id, workerChannel := range js.Store.RunManager.Workers {
		if _, ok := activeJobRunIDs[id]; !ok {
			close(workerChannel)
			delete(js.Store.RunManager.Workers, id)
		}
	}
}

// Stop closes all workers that have been started to process Job Runs on new
// heads and waits for them to finish.
func (js *jobSubscriber) Stop() {
	js.workerMutex.Lock()
	for _, workerChannel := range js.Store.RunManager.Workers {
		workerChannel <- nil
	}
	js.workerMutex.Unlock()
	utils.WaitTimeout(&js.workerWaiter, 10*time.Second)
}
