package services

import (
	"context"
	"fmt"
	"math/big"
	"sync"

	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/service"
	httypes "github.com/smartcontractkit/chainlink/core/services/headtracker/types"
	"github.com/smartcontractkit/chainlink/core/store"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"go.uber.org/multierr"
)

var (
	numberJobSubscriptions = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "job_subscriber_subscriptions",
		Help: "The number of job subscriptions currently active",
	})
)

//go:generate mockery --name JobSubscriber  --output ../internal/mocks/ --case=underscore

// JobSubscriber listens for push notifications of event logs from the ethereum
// node's websocket for specific jobs by subscribing to ethLogs.
type JobSubscriber interface {
	httypes.HeadTrackable
	AddJob(job models.JobSpec, bn *models.Head) error
	RemoveJob(ID models.JobID) error
	Jobs() []models.JobSpec
	service.Service
}

// jobSubscriber implementation
type jobSubscriber struct {
	store            *store.Store
	jobSubscriptions map[string]JobSubscription
	jobsMutex        *sync.RWMutex
	runManager       RunManager
	jobResumer       utils.SleeperTask
	nextBlockWorker  *nextBlockWorker
}

type nextBlockWorker struct {
	runManager RunManager
	head       big.Int
	headMtx    sync.RWMutex
}

func (b *nextBlockWorker) getHead() big.Int {
	b.headMtx.RLock()
	defer b.headMtx.RUnlock()
	return b.head
}

func (b *nextBlockWorker) setHead(h big.Int) {
	b.headMtx.Lock()
	b.head = h
	b.headMtx.Unlock()
}

func (b *nextBlockWorker) Work() {
	head := b.getHead()
	err := b.runManager.ResumeAllPendingNextBlock(&head)
	if err != nil {
		logger.Errorw("Failed to resume confirming tasks on new head", "error", err)
	}
}

// NewJobSubscriber returns a new job subscriber.
func NewJobSubscriber(store *store.Store, runManager RunManager) JobSubscriber {
	b := &nextBlockWorker{runManager: runManager}
	js := &jobSubscriber{
		store:            store,
		runManager:       runManager,
		jobSubscriptions: map[string]JobSubscription{},
		jobsMutex:        &sync.RWMutex{},
		jobResumer:       utils.NewSleeperTask(b),
		nextBlockWorker:  b,
	}
	return js
}

func (js *jobSubscriber) Start() error {
	return nil
}

// Called on node shutdown, unsubscribe from everything
// and remove the subscriptions.
func (js *jobSubscriber) Close() error {
	js.jobsMutex.Lock()
	defer js.jobsMutex.Unlock()

	for _, sub := range js.jobSubscriptions {
		sub.Unsubscribe()
	}
	js.jobSubscriptions = map[string]JobSubscription{}
	return js.jobResumer.Stop()
}

func (js *jobSubscriber) Ready() error {
	return nil
}

func (js *jobSubscriber) Healthy() error {
	return nil
}

func (js *jobSubscriber) alreadySubscribed(jobID models.JobID) bool {
	js.jobsMutex.RLock()
	defer js.jobsMutex.RUnlock()
	if _, exists := js.jobSubscriptions[jobID.String()]; exists {
		return true
	}
	return false
}

// AddJob subscribes to ethereum log events for each "runlog" and "ethlog"
// initiator in the passed job spec.
func (js *jobSubscriber) AddJob(job models.JobSpec, bn *models.Head) error {
	if !job.IsLogInitiated() {
		return nil
	}
	if js.store.Config.EthereumDisabled() {
		logger.Errorw(fmt.Sprintf("ACTION REQUIRED: Attempted to add job with name '%s' but Ethereum was disabled. This job is NOT running.", job.Name), "job", job)
		return nil
	}

	if js.alreadySubscribed(job.ID) {
		return nil
	}
	// Create a new subscription for this job
	sub, err := StartJobSubscription(job, bn, js.store, js.runManager)
	if err != nil {
		js.store.UpsertErrorFor(job.ID, "Unable to start job subscription")
		return err
	}
	js.addSubscription(sub)
	return nil
}

// RemoveJob unsubscribes the job from a log subscription to trigger runs.
func (js *jobSubscriber) RemoveJob(ID models.JobID) error {
	js.jobsMutex.Lock()
	sub, ok := js.jobSubscriptions[ID.String()]
	delete(js.jobSubscriptions, ID.String())
	numberJobSubscriptions.Set(float64(len(js.jobSubscriptions)))
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
	numberJobSubscriptions.Set(float64(len(js.jobSubscriptions)))
}

// Connect connects the jobs to the ethereum node by creating corresponding subscriptions.
func (js *jobSubscriber) Connect(bn *models.Head) error {
	logger.Debugw("JobSubscriber connect", "head", bn)

	var merr error
	err := js.store.Jobs(
		func(j *models.JobSpec) bool {
			logger.Debugw("JobSubscriber adding job", "jobSpecID", j.ID)
			merr = multierr.Append(merr, js.AddJob(*j, bn))
			return true
		},
		models.InitiatorEthLog,
		models.InitiatorRandomnessLog,
		models.InitiatorRunLog,
	)
	return multierr.Append(merr, err)
}

// OnNewLongestChain resumes all pending job runs based on the new head activity.
func (js *jobSubscriber) OnNewLongestChain(ctx context.Context, head models.Head) {
	js.nextBlockWorker.setHead(*head.ToInt())
	js.jobResumer.WakeUp()
}

// NullJobSubscriber implements Null pattern for JobSubscriber interface
type NullJobSubscriber struct{}

func (NullJobSubscriber) Connect(head *models.Head) error                         { return nil }
func (NullJobSubscriber) OnNewLongestChain(ctx context.Context, head models.Head) {}
func (NullJobSubscriber) AddJob(job models.JobSpec, bn *models.Head) error {
	return errors.New("NullJobSubscriber#AddJob should never be called")
}
func (NullJobSubscriber) RemoveJob(ID models.JobID) error {
	return errors.New("NullJobSubscriber#RemoveJob should never be called")
}
func (NullJobSubscriber) Jobs() (j []models.JobSpec) { return }
func (NullJobSubscriber) Start() error               { return nil }
func (NullJobSubscriber) Close() error               { return nil }
func (NullJobSubscriber) Ready() error               { return nil }
func (NullJobSubscriber) Healthy() error             { return nil }
