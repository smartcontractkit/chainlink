package services

import (
	"sync"

	"github.com/smartcontractkit/chainlink/logger"
	"github.com/smartcontractkit/chainlink/store"
	"github.com/smartcontractkit/chainlink/store/models"
	"go.uber.org/multierr"
)

// JobSubscriber listens for push notifications from the ethereum node's
// websocket for specific jobs.
type JobSubscriber interface {
	HeadTrackable
	AddJob(job models.JobSpec, bn *models.IndexableBlockNumber) error
	Jobs() []models.JobSpec
}

// jobSubscriber implementation
type jobSubscriber struct {
	store            *store.Store
	jobSubscriptions []JobSubscription
	jobsMutex        sync.Mutex
}

// NewJobSubscriber returns a new job subscriber.
func NewJobSubscriber(store *store.Store) JobSubscriber {
	return &jobSubscriber{store: store}
}

// AddJob subscribes to ethereum log events for each "runlog" and "ethlog"
// initiator in the passed job spec.
func (js *jobSubscriber) AddJob(job models.JobSpec, bn *models.IndexableBlockNumber) error {
	if !job.IsLogInitiated() {
		return nil
	}

	sub, err := StartJobSubscription(job, bn, js.store)
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
	jobs, err := js.store.Jobs()
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
	pendingRuns, err := js.store.JobRunsWithStatus(models.RunStatusPendingConfirmations, models.RunStatusInProgress)
	if err != nil {
		logger.Error("error fetching pending job runs:", err.Error())
	}

	ibn := head.ToIndexableBlockNumber()
	for _, jr := range pendingRuns {
		if err := js.store.RunChannel.Send(jr.ID, ibn); err != nil {
			logger.Error("JobSubscriber.OnNewHead: ", err.Error())
		}
	}
}
