package services

import (
	"context"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/eth"
	strpkg "github.com/smartcontractkit/chainlink/core/store"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/store/presenters"
	"github.com/smartcontractkit/chainlink/core/utils"

	ethereum "github.com/ethereum/go-ethereum"
	gethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	promNumHeadsReceived = promauto.NewCounter(prometheus.CounterOpts{
		Name: "head_tracker_heads_received",
		Help: "The total number of heads seen",
	})
	promHeadsInQueue = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "head_tracker_heads_in_queue",
		Help: "The number of heads currently waiting to be executed. You can think of this as the 'load' on the head tracker. Should rarely or never be more than 0",
	})
	promCallbackDuration = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "head_tracker_callback_execution_duration",
		Help: "How long it took to execute all callbacks (ms)",
	})
	promCallbackDurationHist = promauto.NewHistogram(prometheus.HistogramOpts{
		Name:    "head_tracker_callback_execution_duration_hist",
		Help:    "How long it took to execute all callbacks (ms) histogram",
		Buckets: []float64{50, 100, 250, 500, 1000, 2000, 5000, 10000, 15000, 30000, 100000},
	})
)

const (
	// Log a warning if OnNewLongestChain callback execution takes longer than this amount of time
	callbackExecutionThreshold = 10 * time.Second

	// The size of the buffer for the headers channel. Note that callback
	// execution is synchronous and could take a non-trivial amount of time.
	headsBufferSize = 5

	// Maximum time we are allowed to spend backfilling heads
	backfillTimeBudget = 15 * time.Second
)

// HeadTracker holds and stores the latest block number experienced by this particular node
// in a thread safe manner. Reconstitutes the last block number from the data
// store on reboot.
type HeadTracker struct {
	callbacks             []strpkg.HeadTrackable
	headers               chan gethTypes.Header
	headSubscription      eth.Subscription
	highestSeenHead       *models.Head
	store                 *strpkg.Store
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

	if err := ht.setHighestSeenHeadFromDB(); err != nil {
		return err
	}
	if ht.highestSeenHead != nil {
		logger.Debug("Tracking logs from last block ", presenters.FriendlyBigInt(ht.highestSeenHead.ToInt()), " with hash ", ht.highestSeenHead.Hash.Hex())
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
func (ht *HeadTracker) Save(h models.Head) error {
	ht.headMutex.Lock()
	if h.GreaterThan(ht.highestSeenHead) {
		ht.highestSeenHead = &h
	}
	ht.headMutex.Unlock()

	err := ht.store.IdempotentInsertHead(h)
	if err != nil {
		return err
	}
	return ht.store.TrimOldHeads(ht.store.Config.EthHeadTrackerHistoryDepth())
}

// HighestSeenHead returns the block header with the highest number that has been seen, or nil
func (ht *HeadTracker) HighestSeenHead() *models.Head {
	ht.headMutex.RLock()
	defer ht.headMutex.RUnlock()

	return ht.highestSeenHead
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

func (ht *HeadTracker) listenForNewHeads() {
	defer ht.listenForNewHeadsWg.Done()
	defer func() {
		err := ht.unsubscribeFromHead()
		logger.ErrorIf(err, "failed when unsubscribe from head")
	}()

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
		err := ht.unsubscribeFromHead()
		if err != nil {
			logger.ErrorIf(err, "failed when unsubscribe from head")
			return false
		}

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

// This should be safe to run concurrently across multiple nodes connected to the same database
func (ht *HeadTracker) receiveHeaders() error {
	for {
		select {
		case <-ht.done:
			return nil
		case blockHeader, open := <-ht.headers:
			promNumHeadsReceived.Inc()
			promHeadsInQueue.Set(float64(len(ht.headers)))
			if !open {
				return errors.New("HeadTracker headers prematurely closed")
			}
			if err := ht.handleNewHead(blockHeader); err != nil {
				return err
			}
		case err, open := <-ht.headSubscription.Err():
			if open && err != nil {
				return err
			}
		}
	}
}

func (ht *HeadTracker) handleNewHead(bh gethTypes.Header) error {
	defer func(start time.Time, number int64) {
		elapsed := time.Since(start)
		ms := float64(elapsed.Milliseconds())
		promCallbackDuration.Set(ms)
		promCallbackDurationHist.Observe(ms)
		if elapsed > callbackExecutionThreshold {
			logger.Warnw(fmt.Sprintf("HeadTracker finished processing head %v in %s which exceeds callback execution threshold of %s", number, elapsed.String(), callbackExecutionThreshold.String()), "blockNumber", number, "time", elapsed)
		} else {
			logger.Debugw(fmt.Sprintf("HeadTracker finished processing head %v in %s", number, elapsed.String()), "blockNumber", number, "time", elapsed)
		}
	}(time.Now(), bh.Number.Int64())
	head := models.NewHeadFromBlockHeader(bh)
	prevHead := ht.HighestSeenHead()

	logger.Debugw(
		fmt.Sprintf("Received new head %v", presenters.FriendlyBigInt(head.ToInt())),
		"blockHeight", head.ToInt(),
		"blockHash", bh.Hash())

	if err := ht.Save(head); err != nil {
		return err
	}

	if prevHead == nil || head.Number > prevHead.Number {
		return ht.handleNewHighestHead(head)
	}
	if head.Number == prevHead.Number {
		if head.Hash != prevHead.Hash {
			logger.Debugw("HeadTracker: got duplicate head", "blockNum", head.Number, "gotHead", head.Hash.Hex(), "highestSeenHead", ht.highestSeenHead.Hash.Hex())
		} else {
			logger.Debugw("HeadTracker: head already in the database", "gotHead", head.Hash.Hex())
		}
	} else {
		logger.Debugw("HeadTracker: got out of order head", "blockNum", head.Number, "gotHead", head.Hash.Hex(), "highestSeenHead", ht.highestSeenHead.Number)
	}
	return nil
}

func (ht *HeadTracker) handleNewHighestHead(head models.Head) error {
	// NOTE: We must set a hard time limit on this, backfilling heads should
	// not block the head tracker
	ctx, cancel := context.WithTimeout(context.Background(), backfillTimeBudget)
	defer cancel()

	headWithChain, err := ht.GetChainWithBackfill(ctx, head, ht.store.Config.EthFinalityDepth())
	if err != nil {
		return err
	}

	ht.onNewLongestChain(headWithChain)
	return nil
}

// GetChainWithBackfill returns a chain of the given length, backfilling any
// heads that may be missing from the database
func (ht *HeadTracker) GetChainWithBackfill(ctx context.Context, head models.Head, depth uint) (models.Head, error) {
	head, err := ht.store.Chain(head.Hash, depth)
	if err != nil {
		return head, errors.Wrap(err, "GetChainWithBackfill failed fetching chain")
	}
	if uint(head.ChainLength()) >= depth {
		return head, nil
	}
	baseHeight := head.Number - int64(depth-1)
	if baseHeight < 0 {
		baseHeight = 0
	}

	if err := ht.backfill(ctx, head.EarliestInChain(), baseHeight); err != nil {
		return head, errors.Wrap(err, "GetChainWithBackfill failed backfilling")
	}
	return ht.store.Chain(head.Hash, depth)
}

// backfill fetches all missing heads up until the base height
func (ht *HeadTracker) backfill(ctx context.Context, head models.Head, baseHeight int64) error {
	if head.Number <= baseHeight {
		return nil
	}
	logger.Debugw("HeadTracker: backfill heads", "n", head.Number-baseHeight, "fromBlockHeight", baseHeight, "toBlockHeight", head.Number-1)

	fetched := 0

	for i := head.Number - 1; i >= baseHeight; i-- {
		// NOTE: Sequential requests here mean it's a potential performance bottleneck, be aware!
		existingHead, err := ht.store.HeadByHash(head.ParentHash)
		if err != nil {
			return errors.Wrap(err, "HeadByHash failed")
		}
		if existingHead != nil {
			head = *existingHead
			continue
		}
		head, err = ht.fetchAndSaveHead(ctx, i)
		fetched++
		if err != nil {
			if errors.Cause(err) == ethereum.NotFound {
				logger.Errorw("HeadTracker: backfill failed to fetch head (not found), chain will be truncated for this head", "headNum", i)
			} else if errors.Cause(err) == context.DeadlineExceeded {
				logger.Warnw("HeadTracker: backfill deadline exceeded, chain will be truncated for this head", "headNum", i)
			} else {
				logger.Errorw("HeadTracker: backfill encountered unknown error, chain will be truncated for this head", "headNum", i, "err", err)
			}
			break
		}
	}
	logger.Debugw("HeadTracker: backfill complete", "fetched", fetched)
	return nil
}

func (ht *HeadTracker) fetchAndSaveHead(ctx context.Context, n int64) (models.Head, error) {
	logger.Debugw("HeadTracker: fetching head", "blockHeight", n)
	var h *gethTypes.Header
	if err := ht.store.GethClientWrapper.GethClient(func(c eth.GethClient) error {
		var err error
		h, err = c.HeaderByNumber(ctx, big.NewInt(n))
		return err
	}); err != nil {
		if err.Error() == "not found" {
			return models.Head{}, ethereum.NotFound
		}
		return models.Head{}, errors.Wrap(err, "GethClient.HeaderByNumber failed")
	}

	if h == nil {
		return models.Head{}, errors.New("invariant violation: expected head to not be nil")
	}

	head := models.NewHeadFromBlockHeader(*h)
	if err := ht.store.IdempotentInsertHead(head); err != nil {
		return models.Head{}, err
	}
	return head, nil
}

func (ht *HeadTracker) onNewLongestChain(headWithChain models.Head) {
	ht.headMutex.Lock()
	defer ht.headMutex.Unlock()
	logger.Debugw("HeadTracker initiating callbacks", "headNum", headWithChain.Number, "chainLength", headWithChain.ChainLength())

	for _, trackable := range ht.callbacks {
		trackable.OnNewLongestChain(headWithChain)
	}
}

func (ht *HeadTracker) subscribeToHead() error {
	ht.headMutex.Lock()
	defer ht.headMutex.Unlock()

	ctx := context.Background()
	ht.headers = make(chan gethTypes.Header, headsBufferSize)
	sub, err := ht.store.TxManager.SubscribeToNewHeads(ctx, ht.headers)
	if err != nil {
		return errors.Wrap(err, "TxManager#SubscribeToNewHeads")
	}

	if err := verifyEthereumChainID(ht); err != nil {
		return errors.Wrap(err, "verifyEthereumChainID failed")
	}

	ht.headSubscription = sub
	ht.connected = true

	ht.connect(ht.highestSeenHead)
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

func (ht *HeadTracker) setHighestSeenHeadFromDB() error {
	head, err := ht.store.LastHead()
	if err != nil {
		return err
	}
	ht.highestSeenHead = head
	return nil
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
