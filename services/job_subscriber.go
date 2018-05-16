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
type JobSubscriber struct {
	Store            *store.Store
	jobSubscriptions []JobSubscription
	jobsMutex        sync.Mutex
}

// AddJob subscribes to ethereum log events for each "runlog" and "ethlog"
// initiator in the passed job spec.
func (js *JobSubscriber) AddJob(job models.JobSpec, bn *models.IndexableBlockNumber) error {
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
func (js *JobSubscriber) Jobs() []models.JobSpec {
	var jobs []models.JobSpec
	for _, js := range js.jobSubscriptions {
		jobs = append(jobs, js.Job)
	}
	return jobs
}

func (js *JobSubscriber) addSubscription(sub JobSubscription) {
	js.jobsMutex.Lock()
	defer js.jobsMutex.Unlock()
	js.jobSubscriptions = append(js.jobSubscriptions, sub)
}

// Connect connects the jobs to the ethereum node by creating corresponding subscriptions.
func (js *JobSubscriber) Connect(bn *models.IndexableBlockNumber) error {
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
func (js *JobSubscriber) Disconnect() {
	js.jobsMutex.Lock()
	defer js.jobsMutex.Unlock()
	for _, sub := range js.jobSubscriptions {
		sub.Unsubscribe()
	}
	js.jobSubscriptions = []JobSubscription{}
}

// OnNewHead resumes all pending job runs based on the new head activity.
func (js *JobSubscriber) OnNewHead(head *models.BlockHeader) {
	pendingRuns, err := js.Store.JobRunsWithStatus(models.RunStatusPendingConfirmations)
	if err != nil {
		logger.Error(err.Error())
	}
	for _, jr := range pendingRuns {
		if _, err := ExecuteRunAtBlock(jr, js.Store, jr.Result, head.ToIndexableBlockNumber()); err != nil {
			logger.Error(err.Error())
		}
	}
}
