package services

import (
	"errors"
	"fmt"
	"sync"

	"github.com/asdine/storm"
	uuid "github.com/satori/go.uuid"
	"github.com/smartcontractkit/chainlink/logger"
	"github.com/smartcontractkit/chainlink/store"
	"github.com/smartcontractkit/chainlink/store/models"
	"github.com/smartcontractkit/chainlink/store/presenters"
	"github.com/smartcontractkit/chainlink/utils"
)

// HeadTrackable represents any object that wishes to respond to ethereum events,
// after being attached to HeadTracker.
type HeadTrackable interface {
	Connect(*models.IndexableBlockNumber) error
	Disconnect()
	OnNewHead(*models.BlockHeader)
}

// HeadTracker holds and stores the latest block number experienced by this particular node
// in a thread safe manner. Reconstitutes the last block number from the data
// store on reboot.
type HeadTracker struct {
	trackers               map[string]HeadTrackable
	headers                chan models.BlockHeader
	headSubscription       models.EthSubscription
	store                  *store.Store
	number                 *models.IndexableBlockNumber
	headMutex              sync.RWMutex
	trackersMutex          sync.RWMutex
	connected              bool
	sleeper                utils.Sleeper
	listenToNewHeadsCloser chan struct{}
	listenToNewHeadsWaiter sync.WaitGroup
}

// NewHeadTracker instantiates a new HeadTracker using the orm to persist new block numbers.
// Can be passed in an optional sleeper object that will dictate how often
// it tries to reconnect.
func NewHeadTracker(store *store.Store, sleepers ...utils.Sleeper) *HeadTracker {
	var sleeper utils.Sleeper
	if len(sleepers) > 0 {
		sleeper = sleepers[0]
	} else {
		sleeper = utils.NewBackoffSleeper()
	}
	return &HeadTracker{
		store:                  store,
		trackers:               map[string]HeadTrackable{},
		sleeper:                sleeper,
		listenToNewHeadsCloser: make(chan struct{}),
	}
}

// Start retrieves the last persisted block number from the HeadTracker,
// subscribes to new heads, and if successful fires Connect on the
// HeadTrackable argument.
func (ht *HeadTracker) Start() error {
	numbers := []models.IndexableBlockNumber{}
	err := ht.store.Select().OrderBy("Digits", "Number").Limit(1).Reverse().Find(&numbers)
	if err != nil && err != storm.ErrNotFound {
		return err
	}
	if len(numbers) > 0 {
		ht.number = &numbers[0]
	}

	ht.headers = make(chan models.BlockHeader)
	sub, err := ht.subscribeToNewHeads(ht.headers)
	if err != nil {
		return err
	}
	ht.headSubscription = sub
	ht.connect(ht.number)
	go ht.updateBlockHeader()
	go ht.listenToNewHeads()
	return nil
}

// Stop unsubscribes all connections and fires Disconnect.
func (ht *HeadTracker) Stop() error {
	if ht.headSubscription == nil {
		return nil
	}
	ht.headSubscription.Unsubscribe()
	ht.headSubscription = nil

	ht.listenToNewHeadsCloser <- struct{}{}

	ht.listenToNewHeadsWaiter.Wait()

	close(ht.headers)
	ht.disconnect()
	return nil
}

// Save updates the latest block number, if indeed the latest, and persists
// this number in case of reboot. Thread safe.
func (ht *HeadTracker) Save(n *models.IndexableBlockNumber) error {
	if n == nil {
		return errors.New("Cannot save a nil block header")
	}

	ht.headMutex.Lock()
	if n.GreaterThan(ht.number) {
		copy := *n
		ht.number = &copy
	}
	ht.headMutex.Unlock()
	return ht.store.Save(n)
}

// LastRecord returns the latest block header being tracked, or nil.
func (ht *HeadTracker) LastRecord() *models.IndexableBlockNumber {
	ht.headMutex.RLock()
	defer ht.headMutex.RUnlock()
	return ht.number
}

// Attach registers an object that will have HeadTrackable events fired on occurence,
// such as Connect.
func (ht *HeadTracker) Attach(t HeadTrackable) string {
	ht.trackersMutex.Lock()
	defer ht.trackersMutex.Unlock()
	id := uuid.Must(uuid.NewV4()).String()
	ht.trackers[id] = t
	if ht.connected {
		t.Connect(ht.LastRecord())
	}
	return id
}

// Detach deregisters an object from having HeadTrackable events fired.
func (ht *HeadTracker) Detach(id string) {
	ht.trackersMutex.Lock()
	defer ht.trackersMutex.Unlock()
	t, present := ht.trackers[id]
	if ht.connected && present {
		t.Disconnect()
	}
	delete(ht.trackers, id)
}

// IsConnected returns whether or not this HeadTracker is connected.
func (ht *HeadTracker) IsConnected() bool { return ht.connected }

func (ht *HeadTracker) connect(bn *models.IndexableBlockNumber) {
	ht.trackersMutex.RLock()
	defer ht.trackersMutex.RUnlock()
	ht.connected = true
	for _, t := range ht.trackers {
		logger.WarnIf(t.Connect(bn))
	}
}

func (ht *HeadTracker) disconnect() {
	ht.trackersMutex.RLock()
	defer ht.trackersMutex.RUnlock()
	ht.connected = false
	for _, t := range ht.trackers {
		t.Disconnect()
	}
}

func (ht *HeadTracker) onNewHead(head *models.BlockHeader) {
	ht.trackersMutex.RLock()
	defer ht.trackersMutex.RUnlock()
	for _, t := range ht.trackers {
		t.OnNewHead(head)
	}
}

func (ht *HeadTracker) subscribeToNewHeads(headers chan models.BlockHeader) (models.EthSubscription, error) {
	sub, err := ht.store.TxManager.SubscribeToNewHeads(headers)
	if err != nil {
		return nil, err
	}
	go func() {
		err := <-sub.Err()
		if err != nil {
			logger.Errorw("Error in new head subscription, disconnected", "err", err)
			ht.reconnectLoop()
		}
	}()
	return sub, nil
}

func (ht *HeadTracker) updateBlockHeader() {
	header, err := ht.store.TxManager.GetBlockByNumber("latest")
	if err != nil {
		logger.Errorw("Unable to update latest block header", "err", err)
		return
	}

	bn := header.ToIndexableBlockNumber()
	if bn.GreaterThan(ht.LastRecord()) {
		logger.Debug("Fast forwarding to block header ", presenters.FriendlyBigInt(bn.ToInt()))
		ht.Save(bn)
	}
}

func (ht *HeadTracker) listenToNewHeads() {
	ht.listenToNewHeadsWaiter.Add(1)
	defer ht.listenToNewHeadsWaiter.Done()

	ht.headMutex.RLock()
	if ht.number != nil {
		logger.Debug("Tracking logs from last block ", presenters.FriendlyBigInt(ht.number.ToInt()), " with hash ", ht.number.Hash.String())
	}
	ht.headMutex.RUnlock()

	for {
		select {
		case header := <-ht.headers:
			number := header.ToIndexableBlockNumber()
			logger.Debugw(fmt.Sprintf("Received header %v with hash %s", presenters.FriendlyBigInt(number.ToInt()), header.Hash().String()), "hash", header.Hash())
			if err := ht.Save(number); err != nil {
				logger.Error(err.Error())
			} else {
				ht.onNewHead(&header)
			}
		case <-ht.listenToNewHeadsCloser:
			return
		}
	}
}

func (ht *HeadTracker) reconnectLoop() {
	ht.sleeper.Reset()
	for {
		ht.Stop()
		logger.Info("Reconnecting to node ", ht.store.Config.EthereumURL, " in ", ht.sleeper.Duration())
		ht.sleeper.Sleep()
		err := ht.Start()
		if err != nil {
			logger.Errorw(fmt.Sprintf("Error reconnecting to %v", ht.store.Config.EthereumURL), "err", err)
		} else {
			logger.Info("Reconnected to node ", ht.store.Config.EthereumURL)
			break
		}
	}
}
