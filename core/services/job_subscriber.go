package services

import (
	"fmt"
	"sync"

	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/store"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"go.uber.org/multierr"
)

// JobSubscriber listens for push notifications of event logs from the ethereum
// node's websocket for specific jobs by subscribing to ethLogs.
type JobSubscriber interface {
	store.HeadTrackable
	AddJob(job models.JobSpec, bn *models.Head) error
	RemoveJob(ID *models.ID) error
	Jobs() []models.JobSpec
}

// jobSubscriber implementation
type jobSubscriber struct {
	store            *store.Store
	jobSubscriptions map[string]JobSubscription
	jobsMutex        *sync.RWMutex
}

// NewJobSubscriber returns a new job subscriber.
func NewJobSubscriber(store *store.Store) JobSubscriber {
	return &jobSubscriber{
		store:            store,
		jobSubscriptions: map[string]JobSubscription{},
		jobsMutex:        &sync.RWMutex{},
	}
}

// AddJob subscribes to ethereum log events for each "runlog" and "ethlog"
// initiator in the passed job spec.
func (js *jobSubscriber) AddJob(job models.JobSpec, bn *models.Head) error {
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

// RemoveJob unsubscribes the job from a log subscription to trigger runs.
func (js *jobSubscriber) RemoveJob(ID *models.ID) error {
	js.jobsMutex.Lock()
	sub, ok := js.jobSubscriptions[ID.String()]
	delete(js.jobSubscriptions, ID.String())
	js.jobsMutex.Unlock()
	if !ok {
		return fmt.Errorf("JobSubscriber#RemoveJob: job %s not found", ID)
	}
	sub.Unsubscribe()
	return nil
}

// Jobs returns the jobs being listened to.
func (js *jobSubscriber) Jobs() []models.JobSpec {
	js.jobsMutex.RLock()
	defer js.jobsMutex.RUnlock()
	var jobs []models.JobSpec
	for _, sub := range js.jobSubscriptions {
		jobs = append(jobs, sub.Job)
	}
	return jobs
}

func (js *jobSubscriber) addSubscription(sub JobSubscription) {
	js.jobsMutex.Lock()
	defer js.jobsMutex.Unlock()
	js.jobSubscriptions[sub.Job.ID.String()] = sub
}

// Connect connects the jobs to the ethereum node by creating corresponding subscriptions.
func (js *jobSubscriber) Connect(bn *models.Head) error {
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
	js.jobSubscriptions = map[string]JobSubscription{}
}

// OnNewHead resumes all pending job runs based on the new head activity.
func (js *jobSubscriber) OnNewHead(head *models.Head) {
	height := head.ToInt()

	err := js.store.UnscopedJobRunsWithStatus(func(run *models.JobRun) {
		err := ResumeConfirmingTask(run, js.store.Unscoped(), height)
		if err != nil {
			logger.Errorf("JobSubscriber.OnNewHead: %v", err)
		}

	}, models.RunStatusPendingConnection, models.RunStatusPendingConfirmations)

	if err != nil {
		logger.Errorf("error fetching pending job runs: %v", err)
	}
}
