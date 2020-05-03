package services

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/smartcontractkit/chainlink/core/eth"
	"github.com/smartcontractkit/chainlink/core/logger"
	strpkg "github.com/smartcontractkit/chainlink/core/store"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/store/presenters"
	"github.com/smartcontractkit/chainlink/core/utils"

	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	numberHeadsReceived = promauto.NewCounter(prometheus.CounterOpts{
		Name: "head_tracker_heads_received",
		Help: "The total number of heads seen",
	})
)

// HeadTracker holds and stores the latest block number experienced by this particular node
// in a thread safe manner. Reconstitutes the last block number from the data
// store on reboot.
type HeadTracker struct {
	callbacks             []strpkg.HeadTrackable
	headers               chan eth.BlockHeader
	headSubscription      eth.Subscription
	store                 *strpkg.Store
	head                  *models.Head
	headMutex             sync.RWMutex
	connected             bool
	sleeper               utils.Sleeper
	done                  chan struct{}
	started               bool
	listenForNewHeadsWg   sync.WaitGroup
	subscriptionSucceeded chan struct{}
}

// NewHeadTracker instantiates a new HeadTracker using the orm to persist new block numbers.
// Can be passed in an optional sleeper object that will dictate how often
// it tries to reconnect.
func NewHeadTracker(store *strpkg.Store, callbacks []strpkg.HeadTrackable, sleepers ...utils.Sleeper) *HeadTracker {
	var sleeper utils.Sleeper
	if len(sleepers) > 0 {
		sleeper = sleepers[0]
	} else {
		sleeper = utils.NewBackoffSleeper()
	}
	return &HeadTracker{
		store:     store,
		callbacks: callbacks,
		sleeper:   sleeper,
	}
}

// Start retrieves the last persisted block number from the HeadTracker,
// subscribes to new heads, and if successful fires Connect on the
// HeadTrackable argument.
func (ht *HeadTracker) Start() error {
	ht.headMutex.Lock()
	defer ht.headMutex.Unlock()

	if ht.started {
		return nil
	}

	if err := ht.updateHeadFromDb(); err != nil {
		return err
	}
	number := ht.head
	if number != nil {
		logger.Debug("Tracking logs from last block ", presenters.FriendlyBigInt(number.ToInt()), " with hash ", number.Hash.Hex())
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
	ht.headMutex.Lock()

	if !ht.started {
		ht.headMutex.Unlock()
		return nil
	}

	if ht.connected {
		ht.connected = false
		ht.disconnect()
	}
	logger.Info(fmt.Sprintf("Head tracker disconnecting from %v", ht.store.Config.EthereumURL()))
	close(ht.done)
	close(ht.subscriptionSucceeded)
	ht.started = false
	ht.headMutex.Unlock()

	ht.listenForNewHeadsWg.Wait()
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
		msg := fmt.Sprintf("Cannot save new head confirmation %v because it's equal to or less than current head %v with hash %s", n, ht.head, n.Hash.Hex())
		return errBlockNotLater{msg}
	}
	return ht.store.CreateHead(n)
}

// Head returns the latest block header being tracked, or nil.
func (ht *HeadTracker) Head() *models.Head {
	ht.headMutex.RLock()
	defer ht.headMutex.RUnlock()

	return ht.head
}

// Connected returns whether or not this HeadTracker is connected.
func (ht *HeadTracker) Connected() bool {
	ht.headMutex.RLock()
	defer ht.headMutex.RUnlock()

	return ht.connected
}

func (ht *HeadTracker) connect(bn *models.Head) {
	for _, trackable := range ht.callbacks {
		logger.WarnIf(trackable.Connect(bn))
	}
}

func (ht *HeadTracker) disconnect() {
	for _, trackable := range ht.callbacks {
		trackable.Disconnect()
	}
}

func (ht *HeadTracker) onNewHead(head *models.Head) {
	numberHeadsReceived.Inc()

	ht.headMutex.Lock()
	defer ht.headMutex.Unlock()

	for _, trackable := range ht.callbacks {
		trackable.OnNewHead(head)
	}
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
		logger.Info("Connecting to ethereum node ", ht.store.Config.EthereumURL(), " in ", ht.sleeper.Duration())
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
			head := models.NewHead(block.Number.ToInt(), block.Hash())
			logger.Debugw(
				fmt.Sprintf("Received new head %v", presenters.FriendlyBigInt(head.ToInt())),
				"blockHeight", head.ToInt(),
				"blockHash", block.Hash(),
				"hash", head.Hash)
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
	ht.headMutex.Lock()
	defer ht.headMutex.Unlock()

	ctx := context.Background()
	ht.headers = make(chan eth.BlockHeader)
	sub, err := ht.store.TxManager.SubscribeToNewHeads(ctx, ht.headers)
	if err != nil {
		return errors.Wrap(err, "TxManager#SubscribeToNewHeads")
	}

	if err := verifyEthereumChainID(ht); err != nil {
		return errors.Wrap(err, "verifyEthereumChainID failed")
	}

	ht.headSubscription = sub
	ht.connected = true
	ht.connect(ht.head)
	return nil
}

func (ht *HeadTracker) unsubscribeFromHead() error {
	ht.headMutex.Lock()
	defer ht.headMutex.Unlock()

	if !ht.connected {
		return nil
	}

	timedUnsubscribe(ht.headSubscription)

	ht.connected = false
	ht.disconnect()
	close(ht.headers)
	return nil
}

func (ht *HeadTracker) updateHeadFromDb() error {
	number, err := ht.store.LastHead()
	if err != nil {
		return err
	}
	ht.head = number
	return nil
}

type errBlockNotLater struct {
	message string
}

func (e errBlockNotLater) Error() string {
	return e.message
}

// chainIDVerify checks whether or not the ChainID from the Chainlink config
// matches the ChainID reported by the ETH node connected to this Chainlink node.
func verifyEthereumChainID(ht *HeadTracker) error {
	ethereumChainID, err := ht.store.TxManager.GetChainID()
	if err != nil {
		return err
	}

	if ethereumChainID.Cmp(ht.store.Config.ChainID()) != 0 {
		return fmt.Errorf(
			"ethereum ChainID doesn't match chainlink config.ChainID: config ID=%d, eth RPC ID=%d",
			ht.store.Config.ChainID(),
			ethereumChainID,
		)
	}
	return nil
}
