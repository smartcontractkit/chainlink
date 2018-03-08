package services

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/asdine/storm"
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
	HeadTracker      *HeadTracker
	jobSubscriptions []JobSubscription
	jobsMutex        sync.Mutex
	started          bool
}

// Start obtains the jobs from the store and subscribes to logs and newHeads
// in order to start and resume jobs waiting on events or confirmations.
func (nl *NodeListener) Start() error {
	ht, err := NewHeadTracker(nl.Store, nl)
	if err != nil {
		return err
	}
	nl.HeadTracker = ht
	nl.started = true
	return nl.HeadTracker.Start()
}

// Stop gracefully closes its access to the store's EthNotifications and resets
// resources.
func (nl *NodeListener) Stop() error {
	if nl.started {
		nl.started = false
		nl.HeadTracker.Stop()
	}
	return nil
}

// AddJob looks for "runlog" and "ethlog" Initiators for a given job
// and watches the Ethereum blockchain for the addresses in the job.
func (nl *NodeListener) AddJob(job models.Job) error {
	if !nl.started || !job.IsLogInitiated() {
		return nil
	}

	sub, err := StartJobSubscription(job, nl.HeadTracker.Get(), nl.Store)
	if err != nil {
		return err
	}
	nl.addSubscription(sub)
	return nil
}

func (nl *NodeListener) Jobs() []models.Job {
	var jobs []models.Job
	for _, js := range nl.jobSubscriptions {
		jobs = append(jobs, js.Job)
	}
	return jobs
}

func (nl *NodeListener) addSubscription(sub JobSubscription) {
	nl.jobsMutex.Lock()
	defer nl.jobsMutex.Unlock()
	nl.jobSubscriptions = append(nl.jobSubscriptions, sub)
}

func (nl *NodeListener) Connected() error {
	jobs, err := nl.Store.Jobs()
	if err != nil {
		return err
	}
	for _, j := range jobs {
		err = multierr.Append(err, nl.AddJob(j))
	}
	return err
}

func (nl *NodeListener) Disconnected() {
	nl.jobsMutex.Lock()
	defer nl.jobsMutex.Unlock()
	for _, sub := range nl.jobSubscriptions {
		sub.Unsubscribe()
	}
	nl.jobSubscriptions = []JobSubscription{}
}

func (nl *NodeListener) OnNewHead(_ *models.BlockHeader) {
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

type HeadTrackable interface {
	Connected() error
	Disconnected()
	OnNewHead(*models.BlockHeader)
}

type NoOpHeadTrackable struct{}

func (NoOpHeadTrackable) Connected() error              { return nil }
func (NoOpHeadTrackable) Disconnected()                 {}
func (NoOpHeadTrackable) OnNewHead(*models.BlockHeader) {}

// Holds and stores the latest block number experienced by this particular node
// in a thread safe manner. Reconstitutes the last block number from the data
// store on reboot.
type HeadTracker struct {
	Tracker          HeadTrackable
	headers          chan models.BlockHeader
	headSubscription *rpc.ClientSubscription
	store            *store.Store
	number           *models.IndexableBlockNumber
	mutex            sync.RWMutex
}

// Instantiates a new HeadTracker using the orm to persist new block numbers
func NewHeadTracker(store *store.Store, tracker ...HeadTrackable) (*HeadTracker, error) {
	ht := &HeadTracker{store: store}
	numbers := []models.IndexableBlockNumber{}
	err := store.Select().OrderBy("Digits", "Number").Limit(1).Reverse().Find(&numbers)
	if err != nil && err != storm.ErrNotFound {
		return nil, err
	}
	if len(numbers) > 0 {
		ht.number = &numbers[0]
		logger.Info("Tracking logs from block ", ht.number.FriendlyString(), " with hash ", ht.number.Hash.String())
	}
	if len(tracker) > 0 {
		ht.Tracker = tracker[0]
	} else {
		ht.Tracker = NoOpHeadTrackable{}
	}

	return ht, nil
}

func (ht *HeadTracker) Start() error {
	ht.headers = make(chan models.BlockHeader)
	sub, err := ht.subscribeToNewHeads()
	if err != nil {
		return err
	}
	ht.headSubscription = sub
	if err = ht.Tracker.Connected(); err != nil {
		return err
	}
	go ht.listenToNewHeads()
	return nil
}

func (ht *HeadTracker) Stop() error {
	if ht.headSubscription != nil && ht.headSubscription.Err() != nil {
		ht.headSubscription.Unsubscribe()
		ht.headSubscription = nil
	}
	if ht.headers != nil {
		close(ht.headers)
		ht.headers = nil
	}
	ht.Tracker.Disconnected()
	return nil
}

// Updates the latest block number, if indeed the latest, and persists
// this number in case of reboot. Thread safe.
func (ht *HeadTracker) Save(n *models.IndexableBlockNumber) error {
	if n == nil {
		return errors.New("Cannot save a nil block header")
	}

	ht.mutex.Lock()
	if ht.number == nil || ht.number.ToInt().Cmp(n.ToInt()) < 0 {
		copy := *n
		ht.number = &copy
	}
	ht.mutex.Unlock()
	return ht.store.Save(n)
}

// Returns the latest block header being tracked, or nil.
func (ht *HeadTracker) Get() *models.IndexableBlockNumber {
	ht.mutex.RLock()
	defer ht.mutex.RUnlock()
	return ht.number
}

func (ht *HeadTracker) subscribeToNewHeads() (*rpc.ClientSubscription, error) {
	sub, err := ht.store.TxManager.SubscribeToNewHeads(ht.headers)
	if err != nil {
		return nil, err
	}
	go func() {
		err := <-sub.Err()
		if err != nil {
			logger.Warnw("Error in new head subscription, disconnected", "err", err)
			ht.Stop()
			ht.reconnectLoop()
		}
	}()
	return sub, nil
}

func (ht *HeadTracker) listenToNewHeads() {
	for header := range ht.headers {
		number := header.IndexableBlockNumber()
		logger.Debugw(fmt.Sprintf("Received header %v", number.FriendlyString()), "hash", header.Hash())
		if err := ht.Save(number); err != nil {
			logger.Error(err.Error())
		} else {
			ht.Tracker.OnNewHead(&header)
		}
	}
}

func (ht *HeadTracker) reconnectLoop() {
	b := utils.NewBackoff()
	for {
		t := b.Duration()
		logger.Info("Reconnecting to node ", ht.store.Config.EthereumURL, " in ", t)
		time.Sleep(t)
		err := ht.Start()
		if err != nil {
			logger.Warnw(fmt.Sprintf("Error reconnecting to %v", ht.store.Config.EthereumURL), "err", err)
			ht.Stop()
		} else {
			logger.Info("Reconnected to node ", ht.store.Config.EthereumURL)
			break
		}
	}
}
