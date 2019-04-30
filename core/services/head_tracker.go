package services

import (
	"errors"
	"fmt"
	"sync"
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/store"
	strpkg "github.com/smartcontractkit/chainlink/core/store"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/store/presenters"
	"github.com/smartcontractkit/chainlink/core/utils"
	"github.com/tevino/abool"
)

// HeadTracker holds and stores the latest block number experienced by this particular node
// in a thread safe manner. Reconstitutes the last block number from the data
// store on reboot.
type HeadTracker struct {
	attachments           *attachmentCollection
	headers               chan models.BlockHeader
	headSubscription      models.EthSubscription
	store                 *strpkg.Store
	head                  *models.Head
	headMutex             sync.RWMutex
	connected             *abool.AtomicBool
	sleeper               utils.Sleeper
	done                  chan struct{}
	started               bool
	listenForNewHeadsWg   sync.WaitGroup
	subscriptionSucceeded chan struct{}
	bootMutex             sync.Mutex
}

// NewHeadTracker instantiates a new HeadTracker using the orm to persist new block numbers.
// Can be passed in an optional sleeper object that will dictate how often
// it tries to reconnect.
func NewHeadTracker(store *strpkg.Store, sleepers ...utils.Sleeper) *HeadTracker {
	var sleeper utils.Sleeper
	if len(sleepers) > 0 {
		sleeper = sleepers[0]
	} else {
		sleeper = utils.NewBackoffSleeper()
	}
	return &HeadTracker{
		store:       store,
		attachments: newAttachmentCollection(),
		sleeper:     sleeper,
		connected:   abool.New(),
	}
}

// Start retrieves the last persisted block number from the HeadTracker,
// subscribes to new heads, and if successful fires Connect on the
// HeadTrackable argument.
func (ht *HeadTracker) Start() error {
	ht.bootMutex.Lock()
	defer ht.bootMutex.Unlock()

	if ht.started {
		return nil
	}

	if err := ht.updateHeadFromDb(); err != nil {
		return err
	}
	number := ht.Head()
	if number != nil {
		logger.Debug("Tracking logs from last block ", presenters.FriendlyBigInt(number.ToInt()), " with hash ", number.Hash().Hex())
	}

	ht.done = make(chan struct{})
	ht.subscriptionSucceeded = make(chan struct{})

	ht.listenForNewHeadsWg.Add(1)
	go ht.listenForNewHeads()

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
	ht.listenForNewHeadsWg.Wait()
	close(ht.subscriptionSucceeded)
	ht.started = false
	return nil
}

// Save updates the latest block number, if indeed the latest, and persists
// this number in case of reboot. Thread safe.
func (ht *HeadTracker) Save(n *models.Head) error {
	if n == nil {
		return errors.New("Cannot save a nil block header")
	}

	ht.headMutex.Lock()
	if n.GreaterThan(ht.head) {
		copy := *n
		ht.head = &copy
		ht.headMutex.Unlock()
	} else {
		ht.headMutex.Unlock()
		msg := fmt.Sprintf("Cannot save new head confirmation %v because it's equal to or less than current head %v with hash %s", n, ht.head, n.Hash().Hex())
		return errBlockNotLater{msg}
	}
	return ht.store.SaveHead(n)
}

// Head returns the latest block header being tracked, or nil.
func (ht *HeadTracker) Head() *models.Head {
	ht.headMutex.RLock()
	defer ht.headMutex.RUnlock()
	return ht.head
}

// Attach registers an object that will have HeadTrackable events fired on occurence,
// such as Connect. If the HeadTracker is already connected, Connect will be
// called on the newly attached parameter.
func (ht *HeadTracker) Attach(t store.HeadTrackable) string {
	rval := ht.attachments.attach(t)
	if ht.connected.IsSet() {
		logger.WarnIf(t.Connect(ht.Head()))
	}
	return rval
}

// Detach deregisters an object from having HeadTrackable events fired.
func (ht *HeadTracker) Detach(id string) {
	t, present := ht.attachments.detach(id)
	if ht.connected.IsSet() && present {
		t.Disconnect()
	}
}

// Connected returns whether or not this HeadTracker is connected.
func (ht *HeadTracker) Connected() bool { return ht.connected.IsSet() }

func (ht *HeadTracker) connect(bn *models.Head) {
	ht.attachments.iter(func(t store.HeadTrackable) {
		logger.WarnIf(t.Connect(bn))
	})
}

func (ht *HeadTracker) disconnect() {
	ht.attachments.iter(func(t store.HeadTrackable) {
		t.Disconnect()
	})
}

func (ht *HeadTracker) onNewHead(head *models.Head) {
	ht.attachments.iter(func(t store.HeadTrackable) {
		t.OnNewHead(head)
	})
}

func (ht *HeadTracker) listenForNewHeads() {
	defer ht.listenForNewHeadsWg.Done()
	defer ht.unsubscribeFromHead()

	for {
		if !ht.subscribe() {
			return
		}
		if err := ht.receiveHeaders(); err != nil {
			logger.Errorw(fmt.Sprintf("Error in new head subscription, unsubscribed: %s", err.Error()), "err", err)
			continue
		} else {
			return
		}
	}
}

// subscribe periodically attempts to connect to the ethereum node via websocket.
// It returns true on success, and false if cut short by a done request and did not connect.
func (ht *HeadTracker) subscribe() bool {
	ht.sleeper.Reset()
	for {
		ht.unsubscribeFromHead()
		logger.Info("Connecting to node ", ht.store.Config.EthereumURL(), " in ", ht.sleeper.Duration())
		select {
		case <-ht.done:
			return false
		case <-time.After(ht.sleeper.After()):
			err := ht.subscribeToHead()
			if err != nil {
				logger.Warnw(fmt.Sprintf("Failed to connect to ethereum node %v", ht.store.Config.EthereumURL()), "err", err)
			} else {
				logger.Info("Connected to ethereum node ", ht.store.Config.EthereumURL())
				return true
			}
		}
	}
}

func (ht *HeadTracker) receiveHeaders() error {
	for {
		select {
		case <-ht.done:
			return nil
		case block, open := <-ht.headers:
			if !open {
				return errors.New("HeadTracker headers prematurely closed")
			}
			head := block.ToHead()
			logger.Debugw(fmt.Sprintf("Received header %v with hash %s", presenters.FriendlyBigInt(head.ToInt()), block.Hash().String()), "hash", head.Hash().Hex())
			if err := ht.Save(head); err != nil {
				switch err.(type) {
				case errBlockNotLater:
					logger.Warn(err)
				default:
					logger.Error(err)
				}
			} else {
				ht.onNewHead(head)
			}
		case err, open := <-ht.headSubscription.Err():
			if open && err != nil {
				return err
			}
		}
	}
}

func (ht *HeadTracker) subscribeToHead() error {
	ht.headers = make(chan models.BlockHeader)
	sub, err := ht.store.TxManager.SubscribeToNewHeads(ht.headers)
	if err != nil {
		return err
	}
	ht.headSubscription = sub
	ht.connected.Set()
	ht.connect(ht.Head())
	return nil
}

func (ht *HeadTracker) unsubscribeFromHead() error {
	if !ht.Connected() {
		return nil
	}

	timedUnsubscribe(ht.headSubscription)
	ht.disconnect()
	close(ht.headers)

	ht.connected.UnSet()
	return nil
}

func (ht *HeadTracker) updateHeadFromDb() error {
	number, err := ht.store.LastHead()
	if err != nil {
		return err
	}
	ht.headMutex.Lock()
	ht.head = number
	ht.headMutex.Unlock()
	return nil
}

// attachmentCollection is a thread safe ordered collection
// of HeadTrackables that are attached to HeadTracker.
type attachmentCollection struct {
	trackables map[string]store.HeadTrackable
	sortedIDs  []string // map order is non-deterministic, so keep order.
	mutex      *sync.RWMutex
}

func newAttachmentCollection() *attachmentCollection {
	return &attachmentCollection{
		trackables: map[string]strpkg.HeadTrackable{},
		sortedIDs:  []string{},
		mutex:      &sync.RWMutex{},
	}
}

func (a *attachmentCollection) attach(t store.HeadTrackable) string {
	id := uuid.NewV4().String()

	a.mutex.Lock()
	defer a.mutex.Unlock()

	a.sortedIDs = append(a.sortedIDs, id)
	a.trackables[id] = t
	return id
}

func (a *attachmentCollection) detach(id string) (store.HeadTrackable, bool) {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	t, present := a.trackables[id]
	if present {
		a.sortedIDs = removeTrackableID(id, a.sortedIDs, t)
		delete(a.trackables, id)
		return t, true
	}
	return nil, false
}

// iter iterates over the collection in an ordered thread safe manner, invoking
// the passed callback on each entry.
func (a *attachmentCollection) iter(cb func(store.HeadTrackable)) {
	a.mutex.RLock()
	defer a.mutex.RUnlock()

	for _, id := range a.sortedIDs {
		cb(a.trackables[id])
	}
}

func removeTrackableID(id string, old []string, t store.HeadTrackable) []string {
	idx := indexOf(id, old)
	if idx == -1 {
		logger.Panicf("invariant violated: id %s for %T exists in trackables but not in sortedIDs in attachmentCollection", id, t)
	}
	return append(old[:idx], old[idx+1:]...)
}

func indexOf(id string, arr []string) int {
	for i, v := range arr {
		if v == id {
			return i
		}
	}
	return -1
}

type errBlockNotLater struct {
	message string
}

func (e errBlockNotLater) Error() string {
	return e.message
}
