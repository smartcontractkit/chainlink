package services

import (
	"sync"

	"github.com/smartcontractkit/chainlink/logger"
	"github.com/smartcontractkit/chainlink/store"
	"github.com/smartcontractkit/chainlink/store/models"
	"go.uber.org/multierr"
)

// EthereumListener manages push notifications from the ethereum node's
// websocket to listen for new heads and log events.
type EthereumListener struct {
	Store            *store.Store
	jobSubscriptions []JobSubscription
	jobsMutex        sync.Mutex
}

// AddJob subscribes to ethereum log events for each "runlog" and "ethlog"
// initiator in the passed job spec.
func (el *EthereumListener) AddJob(job models.JobSpec, bn *models.IndexableBlockNumber) error {
	if !job.IsLogInitiated() {
		return nil
	}

	sub, err := StartJobSubscription(job, bn, el.Store)
	if err != nil {
		return err
	}
	el.addSubscription(sub)
	return nil
}

// Jobs returns the jobs being listened to.
func (el *EthereumListener) Jobs() []models.JobSpec {
	var jobs []models.JobSpec
	for _, js := range el.jobSubscriptions {
		jobs = append(jobs, js.Job)
	}
	return jobs
}

func (el *EthereumListener) addSubscription(sub JobSubscription) {
	el.jobsMutex.Lock()
	defer el.jobsMutex.Unlock()
	el.jobSubscriptions = append(el.jobSubscriptions, sub)
}

// Connect connects the jobs to the ethereum node by creating corresponding subscriptions.
func (el *EthereumListener) Connect(bn *models.IndexableBlockNumber) error {
	jobs, err := el.Store.Jobs()
	if err != nil {
		return err
	}
	for _, j := range jobs {
		err = multierr.Append(err, el.AddJob(j, bn))
	}
	return err
}

// Disconnect disconnects all subscriptions associated with jobs belonging to
// this listener.
func (el *EthereumListener) Disconnect() {
	el.jobsMutex.Lock()
	defer el.jobsMutex.Unlock()
	for _, sub := range el.jobSubscriptions {
		sub.Unsubscribe()
	}
	el.jobSubscriptions = []JobSubscription{}
}

// OnNewHead resumes all pending job runs based on the new head activity.
func (el *EthereumListener) OnNewHead(head *models.BlockHeader) {
	pendingRuns, err := el.Store.JobRunsWithStatus(models.RunStatusPendingConfirmations)
	if err != nil {
		logger.Error(err.Error())
	}
	for _, jr := range pendingRuns {
		if _, err := ExecuteRunAtBlock(jr, el.Store, jr.Result, head.ToIndexableBlockNumber()); err != nil {
			logger.Error(err.Error())
		}
	}
}
