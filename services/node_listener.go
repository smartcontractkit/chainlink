package services

import (
	"fmt"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/rpc"
	"github.com/smartcontractkit/chainlink/logger"
	"github.com/smartcontractkit/chainlink/store"
	"github.com/smartcontractkit/chainlink/store/models"
	"github.com/smartcontractkit/chainlink/utils"
	"go.uber.org/multierr"
)

// NodeListener manages push notifications from the ethereum node's
// websocket to listen for new heads and log events.
type NodeListener struct {
	Store            *store.Store
	jobSubscriptions []JobSubscription
	headers          chan models.BlockHeader
	headSubscription *rpc.ClientSubscription
	jobsMutex        sync.Mutex
	started          bool
}

// Start obtains the jobs from the store and subscribes to logs and newHeads
// in order to start and resume jobs waiting on events or confirmations.
func (nl *NodeListener) Start() error {
	nl.started = true
	nl.headers = make(chan models.BlockHeader)
	if err := nl.subscribeToNewHeads(); err != nil {
		return err
	}

	if err := nl.subscribeJobs(); err != nil {
		return err
	}

	go nl.listenToNewHeads()
	return nil
}

// Stop gracefully closes its access to the store's EthNotifications.
func (nl *NodeListener) Stop() error {
	nl.started = false
	if nl.headSubscription != nil && nl.headSubscription.Err() != nil {
		nl.headSubscription.Unsubscribe()
	}
	if nl.headers != nil {
		close(nl.headers)
	}
	nl.unsubscribeJobs()
	return nil
}

func (nl *NodeListener) subscribeJobs() error {
	jobs, err := nl.Store.Jobs()
	if err != nil {
		return err
	}
	for _, j := range jobs {
		err = multierr.Append(err, nl.AddJob(j))
	}
	return err
}

// AddJob looks for "runlog" and "ethlog" Initiators for a given job
// and watches the Ethereum blockchain for the addresses in the job.
func (nl *NodeListener) AddJob(job models.Job) error {
	if !nl.started || !job.IsLogInitiated() {
		return nil
	}

	sub, err := StartJobSubscription(job, nl.Store)
	if err != nil {
		return err
	}
	nl.addSubscription(sub)
	return nil
}

func (nl *NodeListener) subscribeToNewHeads() error {
	sub, err := nl.Store.TxManager.SubscribeToNewHeads(nl.headers)
	if err != nil {
		return err
	}
	nl.headSubscription = sub
	go func() {
		err := <-sub.Err()
		logger.Warnw("Error in new head subscription, disconnected", "err", err)
		nl.Stop()
		nl.reconnectLoop()
	}()
	return nil
}

func (nl *NodeListener) reconnectLoop() {
	b := utils.NewBackoff()
	for {
		t := b.Duration()
		logger.Info("Reconnecting to node in ", t)
		time.Sleep(t)
		err := nl.Start()
		if err != nil {
			logger.Warnw("Error reconnecting", "err", err)
		} else {
			logger.Info("Reconnected to node")
			break
		}
	}
}

func (nl *NodeListener) listenToNewHeads() {
	for header := range nl.headers {
		number := header.IndexableBlockNumber()
		logger.Debugw(fmt.Sprintf("Received header %v", number.FriendlyString()), "hash", header.Hash())
		if err := nl.Store.HeadTracker.Save(number); err != nil {
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

func (nl *NodeListener) addSubscription(sub JobSubscription) {
	nl.jobsMutex.Lock()
	defer nl.jobsMutex.Unlock()
	nl.jobSubscriptions = append(nl.jobSubscriptions, sub)
}

func (nl *NodeListener) unsubscribeJobs() {
	nl.jobsMutex.Lock()
	defer nl.jobsMutex.Unlock()
	for _, sub := range nl.jobSubscriptions {
		sub.Unsubscribe()
	}
}
