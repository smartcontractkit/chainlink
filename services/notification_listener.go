package services

import (
	"fmt"
	"sync"

	"github.com/ethereum/go-ethereum/rpc"
	"github.com/smartcontractkit/chainlink/logger"
	"github.com/smartcontractkit/chainlink/store"
	"github.com/smartcontractkit/chainlink/store/models"
	"go.uber.org/multierr"
)

// NotificationListener manages push notifications from the ethereum node's
// websocket to listen for new heads and log events.
type NotificationListener struct {
	Store             *store.Store
	subscriptions     []Subscription
	headNotifications chan models.BlockHeader
	headSubscription  *rpc.ClientSubscription
	subMutx           sync.Mutex
	started           bool
}

// Start obtains the jobs from the store and subscribes to logs and newHeads
// in order to start and resume jobs waiting on events or confirmations.
func (nl *NotificationListener) Start() error {
	nl.started = true
	nl.headNotifications = make(chan models.BlockHeader)
	if err := nl.subscribeToNewHeads(); err != nil {
		return err
	}

	jobs, err := nl.Store.Jobs()
	if err != nil {
		return err
	}
	if err := nl.subscribeJobs(jobs); err != nil {
		return err
	}

	go nl.listenToNewHeads()
	return nil
}

// Stop gracefully closes its access to the store's EthNotifications.
func (nl *NotificationListener) Stop() error {
	nl.started = false
	if nl.headSubscription != nil && nl.headSubscription.Err() != nil {
		nl.headSubscription.Unsubscribe()
	}
	if nl.headNotifications != nil {
		close(nl.headNotifications)
	}
	nl.unsubscribeJobs()
	return nil
}

func (nl *NotificationListener) subscribeJobs(jobs []models.Job) error {
	var err error
	for _, j := range jobs {
		err = multierr.Append(err, nl.AddJob(j))
	}
	return err
}

// AddJob looks for "runlog" and "ethlog" Initiators for a given job
// and watches the Ethereum blockchain for the addresses in the job.
func (nl *NotificationListener) AddJob(job models.Job) error {
	if !nl.started || !job.IsLogInitiated() {
		return nil
	}

	sub, err := StartSubscription(job, nl.Store)
	if err != nil {
		return err
	}
	nl.addSubscription(sub)
	return nil
}

func (nl *NotificationListener) subscribeToNewHeads() error {
	sub, err := nl.Store.TxManager.SubscribeToNewHeads(nl.headNotifications)
	if err != nil {
		return err
	}
	nl.headSubscription = sub
	go func() {
		err := <-sub.Err()
		logger.Errorw("Error in new head subscription", "err", err)
	}()
	return nil
}

func (nl *NotificationListener) listenToNewHeads() {
	for head := range nl.headNotifications {
		logger.Debugw(fmt.Sprintf("Received new blockchain head %v", head.Number.String()), "newHead", head.Number)
		if err := nl.Store.HeadTracker.Save(&head); err != nil {
			logger.Error(err.Error())
		}
		pendingRuns, err := nl.Store.PendingJobRuns()
		if err != nil {
			logger.Error(err.Error())
		}
		for _, jr := range pendingRuns {
			if _, err := ExecuteRun(jr, nl.Store, models.RunResult{}); err != nil {
				logger.Error(err.Error())
			}
		}
	}
}

func (nl *NotificationListener) addSubscription(sub Subscription) {
	nl.subMutx.Lock()
	defer nl.subMutx.Unlock()
	nl.subscriptions = append(nl.subscriptions, sub)
}

func (nl *NotificationListener) unsubscribeJobs() {
	nl.subMutx.Lock()
	defer nl.subMutx.Unlock()
	for _, sub := range nl.subscriptions {
		sub.Unsubscribe()
	}
}
