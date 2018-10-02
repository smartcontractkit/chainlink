package services

import (
	"errors"
	"fmt"
	"sync"
	"time"

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
	trackers                    map[string]HeadTrackable
	headers                     chan models.BlockHeader
	headSubscription            models.EthSubscription
	store                       *store.Store
	head                        *models.IndexableBlockNumber
	headMutex                   sync.RWMutex
	trackersMutex               sync.RWMutex
	connected                   bool
	sleeper                     utils.Sleeper
	done                        chan struct{}
	disconnected                chan struct{}
	disconnectedWg              sync.WaitGroup
	started                     bool
	connectedWg                 sync.WaitGroup
	connectionRequestListenerWg sync.WaitGroup
	connectionRequest           chan struct{}
	connectionSucceeded         chan struct{}
	bootMutex                   sync.Mutex
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
		store:    store,
		trackers: map[string]HeadTrackable{},
		sleeper:  sleeper,
	}
}

// Start retrieves the last persisted block number from the HeadTracker,
// subscribes to new heads, and if successful fires Connect on the
// HeadTrackable argument.
func (ht *HeadTracker) Start() error {
	ht.bootMutex.Lock()
	defer ht.bootMutex.Unlock()

	numbers := []models.IndexableBlockNumber{}
	err := ht.store.Select().OrderBy("Digits", "Number").Limit(1).Reverse().Find(&numbers)
	if err != nil && err != storm.ErrNotFound {
		return err
	}
	if len(numbers) > 0 {
		ht.headMutex.Lock()
		ht.head = &numbers[0]
		ht.headMutex.Unlock()
	}

	ht.done = make(chan struct{})
	ht.connectionRequest = make(chan struct{})
	ht.connectionSucceeded = make(chan struct{})

	ht.connectionRequestListenerWg.Add(1)
	go ht.connectionRequestListener()
	ht.connectionRequest <- struct{}{}
	<-ht.connectionSucceeded
	ht.started = true
	return nil
}

// Stop unsubscribes all connections and fires Disconnect.
func (ht *HeadTracker) Stop() error {
	ht.bootMutex.Lock()
	defer ht.bootMutex.Unlock()

	if !ht.started {
		return nil
	}

	close(ht.done)
	ht.connectionRequestListenerWg.Wait()
	close(ht.connectionRequest)
	close(ht.connectionSucceeded)
	ht.started = false
	return nil
}

// Save updates the latest block number, if indeed the latest, and persists
// this number in case of reboot. Thread safe.
func (ht *HeadTracker) Save(n *models.IndexableBlockNumber) error {
	if n == nil {
		return errors.New("Cannot save a nil block header")
	}

	ht.headMutex.Lock()
	if n.GreaterThan(ht.head) {
		copy := *n
		ht.head = &copy
	}
	ht.headMutex.Unlock()
	return ht.store.Save(n)
}

// Head returns the latest block header being tracked, or nil.
func (ht *HeadTracker) Head() *models.IndexableBlockNumber {
	ht.headMutex.RLock()
	defer ht.headMutex.RUnlock()
	return ht.head
}

// Attach registers an object that will have HeadTrackable events fired on occurence,
// such as Connect.
func (ht *HeadTracker) Attach(t HeadTrackable) string {
	ht.trackersMutex.Lock()
	defer ht.trackersMutex.Unlock()
	id := uuid.Must(uuid.NewV4()).String()
	ht.trackers[id] = t
	if ht.connected {
		t.Connect(ht.Head())
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
	for _, t := range ht.trackers {
		logger.WarnIf(t.Connect(bn))
	}
}

func (ht *HeadTracker) disconnect() {
	ht.trackersMutex.RLock()
	defer ht.trackersMutex.RUnlock()
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

func (ht *HeadTracker) connectToHead() error {
	ht.disconnected = make(chan struct{})
	ht.headers = make(chan models.BlockHeader)
	sub, err := ht.store.TxManager.SubscribeToNewHeads(ht.headers)
	if err != nil {
		return err
	}
	ht.headSubscription = sub
	ht.connectedWg.Add(2)
	go ht.listenToSubscriptionErrors()
	go ht.listenToNewHeads()
	ht.connectedWg.Wait()
	ht.disconnectedWg.Add(2)
	ht.connected = true
	ht.connect(ht.Head())
	return nil
}

func (ht *HeadTracker) disconnectFromHead() error {
	if !ht.connected {
		return nil
	}

	ht.headSubscription.Unsubscribe()
	close(ht.disconnected)
	ht.disconnectedWg.Wait()

	close(ht.headers)
	ht.connected = false
	ht.disconnect()
	return nil
}

func (ht *HeadTracker) listenToSubscriptionErrors() {
	ht.connectedWg.Done()
	select {
	case err := <-ht.headSubscription.Err():
		ht.disconnectedWg.Done()
		if err != nil {
			logger.Errorw("Error in new head subscription, disconnected", "err", err)
			ht.connectionRequest <- struct{}{}
		}
	case <-ht.disconnected:
		ht.disconnectedWg.Done()
	}
}

func (ht *HeadTracker) updateBlockHeader() {
	header, err := ht.store.TxManager.GetBlockByNumber("latest")
	if err != nil {
		logger.Errorw("Unable to update latest block header", "err", err)
		return
	}

	bn := header.ToIndexableBlockNumber()
	if bn.GreaterThan(ht.Head()) {
		logger.Debug("Fast forwarding to block header ", presenters.FriendlyBigInt(bn.ToInt()))
		ht.Save(bn)
	}
}

func (ht *HeadTracker) listenToNewHeads() {
	defer ht.disconnectedWg.Done()

	ht.updateBlockHeader()
	number := ht.Head()
	if number != nil {
		logger.Debug("Tracking logs from last block ", presenters.FriendlyBigInt(number.ToInt()), " with hash ", number.Hash.String())
	}

	ht.connectedWg.Done()
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
		case <-ht.disconnected:
			return
		}
	}
}

func (ht *HeadTracker) connectionRequestListener() {
	defer ht.connectionRequestListenerWg.Done()
	for {
		select {
		case <-ht.done:
			ht.disconnectFromHead()
			return
		case <-ht.connectionRequest:
			if !ht.reconnectionPoll() {
				return
			}
		}
	}
}

// reconnectionPoll periodically attempts to connect to an ethereum node. It returns true
// on success, and false if cut short by a done request and did not connect.
func (ht *HeadTracker) reconnectionPoll() bool {
	ht.sleeper.Reset()
	for {
		ht.disconnectFromHead()
		logger.Info("Connecting to node ", ht.store.Config.EthereumURL, " in ", ht.sleeper.Duration())
		select {
		case <-ht.done:
			return false
		case <-time.After(ht.sleeper.After()):
			err := ht.connectToHead()
			if err != nil {
				logger.Errorw(fmt.Sprintf("Error connecting to %v", ht.store.Config.EthereumURL), "err", err)
			} else {
				logger.Info("Connected to node ", ht.store.Config.EthereumURL)
				select {
				case ht.connectionSucceeded <- struct{}{}:
				default:
				}
				return true
			}
		}
	}
}
