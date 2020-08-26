package services

import (
	"fmt"
	"math/big"
	"sync"

	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/store"
	"github.com/smartcontractkit/chainlink/core/store/models"

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
	store.HeadTrackable
	AddJob(job models.JobSpec, bn *models.Head) error
	RemoveJob(ID *models.ID) error
	Jobs() []models.JobSpec
	Stop() error
}

// jobSubscriber implementation
type jobSubscriber struct {
	store            *store.Store
	jobSubscriptions map[string]JobSubscription
	jobsMutex        *sync.RWMutex
	runManager       RunManager
	jobResumer       SleeperTask
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
		jobResumer:       NewSleeperTask(b),
		nextBlockWorker:  b,
	}
	return js
}

func (js *jobSubscriber) Stop() error {
	return js.jobResumer.Stop()
}

// AddJob subscribes to ethereum log events for each "runlog" and "ethlog"
// initiator in the passed job spec.
func (js *jobSubscriber) AddJob(job models.JobSpec, bn *models.Head) error {
	if !job.IsLogInitiated() {
		return nil
	}

	sub, err := StartJobSubscription(job, bn, js.store, js.runManager)
	if err != nil {
		js.store.UpsertErrorFor(job.ID, "Unable to start job subscription")
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
	var merr error
	err := js.store.Jobs(
		func(j *models.JobSpec) bool {
			merr = multierr.Append(merr, js.AddJob(*j, bn))
			return true
		},
		models.InitiatorEthLog,
		models.InitiatorRandomnessLog,
		models.InitiatorRunLog,
	)
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

// OnNewLongestChain resumes all pending job runs based on the new head activity.
func (js *jobSubscriber) OnNewLongestChain(head models.Head) {
	js.nextBlockWorker.setHead(*head.ToInt())
	js.jobResumer.WakeUp()
}
