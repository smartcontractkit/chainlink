package services

import (
	"sync"

	"github.com/smartcontractkit/chainlink/logger"
	"github.com/smartcontractkit/chainlink/store"
	"github.com/smartcontractkit/chainlink/store/models"
	"go.uber.org/multierr"
)

// JobSubscriber listens for push notifications of event logs from the ethereum
// node's websocket for specific jobs by subscribing to ethLogs.
type JobSubscriber interface {
	store.HeadTrackable
	AddJob(job models.JobSpec, bn *models.IndexableBlockNumber) error
	Jobs() []models.JobSpec
}

// jobSubscriber implementation
type jobSubscriber struct {
	store            *store.Store
	jobSubscriptions []JobSubscription
	jobsMutex        sync.RWMutex
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
	js.jobsMutex.RLock()
	defer js.jobsMutex.RUnlock()
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
	var merr error
	err := js.store.Jobs(func(j models.JobSpec) bool {
		merr = multierr.Append(merr, js.AddJob(j, bn))
		return true
	})
	return multierr.Append(merr, err)
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
	pendingRuns, err := js.store.JobRunsWithStatus(models.RunStatusPendingConfirmations)
	if err != nil {
		logger.Error("error fetching pending job runs:", err.Error())
	}

	ibn := head.ToIndexableBlockNumber().Number.ToHexUtilBig()
	logger.Debugw("Received new head",
		"current_height", ibn,
		"pending_run_count", len(pendingRuns),
	)
	for _, jr := range pendingRuns {
		_, err := ResumeConfirmingTask(&jr, js.store, ibn)
		if err != nil {
			logger.Error("JobSubscriber.OnNewHead: ", err.Error())
		}
	}
}
