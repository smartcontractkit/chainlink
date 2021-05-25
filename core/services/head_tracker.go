package services

import (
	"context"
	"fmt"
	"math/big"
	"reflect"
	"sync"
	"time"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/postgres"
	strpkg "github.com/smartcontractkit/chainlink/core/store"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/store/presenters"
	"github.com/smartcontractkit/chainlink/core/utils"
)

var (
	promCurrentHead = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "head_tracker_current_head",
		Help: "The highest seen head number",
	})
	promNumHeadsReceived = promauto.NewCounter(prometheus.CounterOpts{
		Name: "head_tracker_heads_received",
		Help: "The total number of heads seen",
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
	promEthConnectionErrors = promauto.NewCounter(prometheus.CounterOpts{
		Name: "head_tracker_eth_connection_errors",
		Help: "The total number of eth node connection errors",
	})
)

// HeadTracker holds and stores the latest block number experienced by this particular node
// in a thread safe manner. Reconstitutes the last block number from the data
// store on reboot.
type HeadTracker struct {
	log                   *logger.Logger
	callbacks             []strpkg.HeadTrackable
	headers               chan *models.Head
	headSubscription      ethereum.Subscription
	highestSeenHead       *models.Head
	store                 *strpkg.Store
	headMutex             sync.RWMutex
	connected             bool
	sleeper               utils.Sleeper
	done                  chan struct{}
	started               bool
	wgDone                sync.WaitGroup
	backfillMB            utils.Mailbox
	samplingMB            utils.Mailbox
	subscriptionSucceeded chan struct{}
	muLogger              sync.RWMutex
}

// NewHeadTracker instantiates a new HeadTracker using the orm to persist new block numbers.
// Can be passed in an optional sleeper object that will dictate how often
// it tries to reconnect.
func NewHeadTracker(l *logger.Logger, store *strpkg.Store, callbacks []strpkg.HeadTrackable, sleepers ...utils.Sleeper) *HeadTracker {
	var sleeper utils.Sleeper
	if len(sleepers) > 0 {
		sleeper = sleepers[0]
	} else {
		sleeper = utils.NewBackoffSleeper()
	}
	return &HeadTracker{
		store:      store,
		callbacks:  callbacks,
		sleeper:    sleeper,
		log:        l,
		backfillMB: *utils.NewMailbox(1),
		samplingMB: *utils.NewMailbox(1),
		done:       make(chan struct{}),
	}
}

// SetLogger sets and reconfigures the log for the head tracker service
func (ht *HeadTracker) SetLogger(logger *logger.Logger) {
	ht.muLogger.Lock()
	defer ht.muLogger.Unlock()
	ht.log = logger
}

func (ht *HeadTracker) logger() *logger.Logger {
	ht.muLogger.RLock()
	defer ht.muLogger.RUnlock()
	return ht.log
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
		ht.logger().Debugw(
			fmt.Sprintf("Headtracker: Tracking logs from last block %v with hash %s", presenters.FriendlyBigInt(ht.highestSeenHead.ToInt()), ht.highestSeenHead.Hash.Hex()),
			"blockNumber", ht.highestSeenHead.Number,
			"blockHash", ht.highestSeenHead.Hash,
		)
	}

	ht.subscriptionSucceeded = make(chan struct{})

	ht.wgDone.Add(3)
	go ht.listenForNewHeads()
	go ht.backfiller()
	go ht.headSampler()

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

	ht.logger().Info(fmt.Sprintf("Head tracker disconnecting from %v", ht.store.Config.EthereumURL()))
	close(ht.done)
	close(ht.subscriptionSucceeded)
	ht.started = false
	ht.headMutex.Unlock()

	ht.wgDone.Wait()
	return nil
}

// Save updates the latest block number, if indeed the latest, and persists
// this number in case of reboot. Thread safe.
func (ht *HeadTracker) Save(ctx context.Context, h models.Head) error {
	ht.headMutex.Lock()
	if h.GreaterThan(ht.highestSeenHead) {
		ht.highestSeenHead = &h
	}
	ht.headMutex.Unlock()

	err := ht.store.IdempotentInsertHead(ctx, h)
	if ctx.Err() != nil {
		return nil
	} else if err != nil {
		return err
	}
	return ht.store.TrimOldHeads(ctx, ht.store.Config.EthHeadTrackerHistoryDepth())
}

// HighestSeenHead returns the block header with the highest number that has been seen, or nil
func (ht *HeadTracker) HighestSeenHead() *models.Head {
	ht.headMutex.RLock()
	defer ht.headMutex.RUnlock()

	if ht.highestSeenHead == nil {
		return nil
	}
	h := *ht.highestSeenHead
	return &h
}

// Connected returns whether or not this HeadTracker is connected.
func (ht *HeadTracker) Connected() bool {
	ht.headMutex.RLock()
	defer ht.headMutex.RUnlock()

	return ht.connected
}

func (ht *HeadTracker) connect(bn *models.Head) {
	for _, trackable := range ht.callbacks {
		if err := trackable.Connect(bn); err != nil {
			ht.logger().Warn(trackable.Connect(bn))
		}
	}
}

func (ht *HeadTracker) disconnect() {
	for _, trackable := range ht.callbacks {
		trackable.Disconnect()
	}
}

func (ht *HeadTracker) headSampler() {
	defer ht.wgDone.Done()

	debounceHead := time.NewTicker(ht.store.Config.EthHeadTrackerSamplingInterval())
	defer debounceHead.Stop()

	ctx, cancel := utils.ContextFromChan(ht.done)
	defer cancel()

	for {
		select {
		case <-ht.done:
			return
		case <-debounceHead.C:
			item, exists := ht.samplingMB.Retrieve()
			if !exists {
				continue
			}
			head, ok := item.(models.Head)
			if !ok {
				panic(fmt.Sprintf("expected `models.Head`, got %T", item))
			}

			ht.onNewLongestChain(ctx, head)
		}
	}
}

func (ht *HeadTracker) listenForNewHeads() {
	defer ht.wgDone.Done()
	defer func() {
		if err := ht.unsubscribeFromHead(); err != nil {
			ht.logger().Warn(errors.Wrap(err, "HeadTracker failed when unsubscribe from head"))
		}
	}()

	ctx, cancel := utils.ContextFromChan(ht.done)
	defer cancel()

	for {
		if !ht.subscribe() {
			break
		}
		err := ht.receiveHeaders(ctx)
		if ctx.Err() != nil {
			break
		} else if err != nil {
			ht.logger().Errorw("Error in new head subscription, unsubscribed", "err", err)
			ht.headers = nil
			continue
		} else {
			break
		}
	}
}

func (ht *HeadTracker) backfiller() {
	defer ht.wgDone.Done()
	for {
		select {
		case <-ht.done:
			return
		case <-ht.backfillMB.Notify():
			for {
				head, exists := ht.backfillMB.Retrieve()
				if !exists {
					break
				}
				h, is := head.(models.Head)
				if !is {
					panic(fmt.Sprintf("expected `models.Head`, got %T", head))
				}
				{
					ctx, cancel := utils.ContextFromChan(ht.done)
					err := ht.Backfill(ctx, h, ht.store.Config.EthFinalityDepth())
					defer cancel()
					if err != nil {
						ht.logger().Warnw("HeadTracker: unexpected error while backfilling heads", "err", err)
					} else if ctx.Err() != nil {
						break
					}
				}
			}
		}
	}
}

// Backfill given a head will fill in any missing heads up to the given depth
func (ht *HeadTracker) Backfill(ctx context.Context, headWithChain models.Head, depth uint) (err error) {
	if uint(headWithChain.ChainLength()) >= depth {
		return nil
	}
	baseHeight := headWithChain.Number - int64(depth-1)
	if baseHeight < 0 {
		baseHeight = 0
	}

	return ht.backfill(ctx, headWithChain.EarliestInChain(), baseHeight)
}

// backfill fetches all missing heads up until the base height
func (ht *HeadTracker) backfill(ctx context.Context, head models.Head, baseHeight int64) (err error) {
	if head.Number <= baseHeight {
		return nil
	}
	mark := time.Now()
	fetched := 0
	ht.logger().Debugw("HeadTracker: starting backfill",
		"blockNumber", head.Number,
		"id", "head_tracker",
		"n", head.Number-baseHeight,
		"fromBlockHeight", baseHeight,
		"toBlockHeight", head.Number-1)
	defer func() {
		if ctx.Err() != nil {
			return
		}
		ht.logger().Debugw("HeadTracker: finished backfill",
			"fetched", fetched,
			"blockNumber", head.Number,
			"time", time.Since(mark),
			"id", "head_tracker",
			"n", head.Number-baseHeight,
			"fromBlockHeight", baseHeight,
			"toBlockHeight", head.Number-1,
			"err", err)
	}()

	for i := head.Number - 1; i >= baseHeight; i-- {
		// NOTE: Sequential requests here mean it's a potential performance bottleneck, be aware!
		var existingHead *models.Head
		existingHead, err = ht.store.HeadByHash(ctx, head.ParentHash)
		if ctx.Err() != nil {
			break
		} else if err != nil {
			return errors.Wrap(err, "HeadByHash failed")
		}
		if existingHead != nil {
			head = *existingHead
			continue
		}
		head, err = ht.fetchAndSaveHead(ctx, i)
		fetched++
		if ctx.Err() != nil {
			break
		} else if err != nil {
			return errors.Wrap(err, "fetchAndSaveHead failed")
		}
	}
	return nil
}

func (ht *HeadTracker) fetchAndSaveHead(ctx context.Context, n int64) (models.Head, error) {
	ht.logger().Debugw("HeadTracker: fetching head", "blockHeight", n)
	head, err := ht.store.EthClient.HeaderByNumber(ctx, big.NewInt(n))
	if ctx.Err() != nil {
		return models.Head{}, nil
	} else if err != nil {
		return models.Head{}, err
	} else if head == nil {
		return models.Head{}, errors.New("got nil head")
	}
	err = ht.store.IdempotentInsertHead(ctx, *head)
	if ctx.Err() != nil {
		return models.Head{}, nil
	} else if err != nil {
		return models.Head{}, err
	}
	return *head, nil
}

// subscribe periodically attempts to connect to the ethereum node via websocket.
// It returns true on success, and false if cut short by a done request and did not connect.
func (ht *HeadTracker) subscribe() bool {
	ht.sleeper.Reset()
	for {
		if err := ht.unsubscribeFromHead(); err != nil {
			ht.logger().Error("failed when unsubscribe from head", err)
			return false
		}

		ht.logger().Info("HeadTracker: Connecting to ethereum node ", ht.store.Config.EthereumURL(), " in ", ht.sleeper.Duration())
		select {
		case <-ht.done:
			return false
		case <-time.After(ht.sleeper.After()):
			err := ht.subscribeToHead()
			if err != nil {
				promEthConnectionErrors.Inc()
				ht.logger().Warnw(fmt.Sprintf("HeadTracker: Failed to connect to ethereum node %v", ht.store.Config.EthereumURL()), "err", err)
			} else {
				ht.logger().Info("HeadTracker: Connected to ethereum node ", ht.store.Config.EthereumURL())
				return true
			}
		}
	}
}

// This should be safe to run concurrently across multiple nodes connected to the same database
// Note: returning nil from receiveHeaders will cause listenForNewHeads to exit completely
func (ht *HeadTracker) receiveHeaders(ctx context.Context) error {
	for {
		select {
		case <-ht.done:
			return nil
		case blockHeader, open := <-ht.headers:
			if !open {
				return errors.New("HeadTracker: headers prematurely closed")
			}
			timeBudget := ht.store.Config.HeadTimeBudget()
			{
				deadlineCtx, cancel := context.WithTimeout(ctx, timeBudget)
				defer cancel()

				if blockHeader == nil {
					ht.logger().Error("HeadTracker: got nil block header")
					continue
				}
				promNumHeadsReceived.Inc()

				err := ht.handleNewHead(ctx, *blockHeader)
				if ctx.Err() != nil {
					// the 'ctx' context is closed only on ht.done - on shutdown, so it's safe to return nil
					return nil
				} else if deadlineCtx.Err() != nil {
					ht.logger().Warnw("HeadTracker: handling of new head timed out", "error", ctx.Err(), "timeBudget", timeBudget.String())
					continue
				} else if err != nil {
					return err
				}
			}
		case err, open := <-ht.headSubscription.Err():
			if open && err != nil {
				return err
			}
		}
	}
}

func (ht *HeadTracker) handleNewHead(ctx context.Context, head models.Head) error {
	prevHead := ht.HighestSeenHead()

	ht.logger().Debugw(fmt.Sprintf("HeadTracker: Received new head %v", presenters.FriendlyBigInt(head.ToInt())),
		"blockHeight", head.ToInt(),
		"blockHash", head.Hash,
		"parentHeadHash", head.ParentHash,
	)

	err := ht.Save(ctx, head)
	if ctx.Err() != nil {
		return nil
	} else if err != nil {
		return err
	}

	if prevHead == nil || head.Number > prevHead.Number {
		promCurrentHead.Set(float64(head.Number))

		headWithChain, err := ht.store.Chain(ctx, head.Hash, ht.store.Config.EthFinalityDepth())
		if ctx.Err() != nil {
			return nil
		} else if err != nil {
			return errors.Wrap(err, "HeadTracker#handleNewHighestHead failed fetching chain")
		}

		ht.backfillMB.Deliver(headWithChain)
		ht.samplingMB.Deliver(headWithChain)
		return nil
	}
	if head.Number == prevHead.Number {
		if head.Hash != prevHead.Hash {
			ht.logger().Debugw("HeadTracker: got duplicate head", "blockNum", head.Number, "gotHead", head.Hash.Hex(), "highestSeenHead", ht.highestSeenHead.Hash.Hex())
		} else {
			ht.logger().Debugw("HeadTracker: head already in the database", "gotHead", head.Hash.Hex())
		}
	} else {
		ht.logger().Debugw("HeadTracker: got out of order head", "blockNum", head.Number, "gotHead", head.Hash.Hex(), "highestSeenHead", ht.highestSeenHead.Number)
	}
	return nil
}

// If total callback execution time exceeds this threshold we consider this to
// be a problem and will log a warning.
// Here we set it to the average time between blocks.
func (ht *HeadTracker) callbackExecutionThreshold() time.Duration {
	return ht.store.Config.HeadTimeBudget() / 2
}

func (ht *HeadTracker) onNewLongestChain(ctx context.Context, headWithChain models.Head) {
	defer func(start time.Time, number int64) {
		elapsed := time.Since(start)
		ms := float64(elapsed.Milliseconds())
		promCallbackDuration.Set(ms)
		promCallbackDurationHist.Observe(ms)
		if elapsed > ht.callbackExecutionThreshold() {
			ht.logger().Warnw(fmt.Sprintf("HeadTracker finished processing head %v in %s which exceeds callback execution threshold of %s", number, elapsed.String(), ht.store.Config.HeadTimeBudget().String()), "blockNumber", number, "time", elapsed, "id", "head_tracker")
		} else {
			ht.logger().Debugw(fmt.Sprintf("HeadTracker finished processing head %v in %s", number, elapsed.String()), "blockNumber", number, "time", elapsed, "id", "head_tracker")
		}
	}(time.Now(), headWithChain.Number)

	ht.headMutex.Lock()
	defer ht.headMutex.Unlock()

	ht.logger().Debugw("HeadTracker initiating callbacks",
		"headNum", headWithChain.Number,
		"chainLength", headWithChain.ChainLength(),
		"numCallbacks", len(ht.callbacks),
	)

	ht.concurrentlyExecuteCallbacks(ctx, headWithChain)
}

func (ht *HeadTracker) concurrentlyExecuteCallbacks(ctx context.Context, headWithChain models.Head) {
	wg := sync.WaitGroup{}
	wg.Add(len(ht.callbacks))
	for idx, trackable := range ht.callbacks {
		go func(i int, t strpkg.HeadTrackable) {
			start := time.Now()
			t.OnNewLongestChain(ctx, headWithChain)
			elapsed := time.Since(start)
			ht.logger().Debugw(fmt.Sprintf("HeadTracker: finished callback %v in %s", i, elapsed), "callbackType", reflect.TypeOf(t), "callbackIdx", i, "blockNumber", headWithChain.Number, "time", elapsed, "id", "head_tracker")
			wg.Done()
		}(idx, trackable)
	}
	wg.Wait()
}

func (ht *HeadTracker) subscribeToHead() error {
	ht.headMutex.Lock()
	defer ht.headMutex.Unlock()

	ht.headers = make(chan *models.Head)

	sub, err := ht.store.EthClient.SubscribeNewHead(context.Background(), ht.headers)
	if err != nil {
		return errors.Wrap(err, "EthClient#SubscribeNewHead")
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
	// ht.headers will be nil if subscription failed, channel closed, and
	// receiveHeaders returned from the loop. listenForNewHeads will set it to
	// nil in that case to avoid a double close panic.
	if ht.headers != nil {
		close(ht.headers)
	}
	return nil
}

func (ht *HeadTracker) setHighestSeenHeadFromDB() error {
	head, err := ht.HighestSeenHeadFromDB()
	if err != nil {
		return err
	}
	ht.highestSeenHead = head
	return nil
}

func (ht *HeadTracker) HighestSeenHeadFromDB() (*models.Head, error) {
	ctxQuery, _ := postgres.DefaultQueryCtx()
	ctx, cancel := utils.CombinedContext(ht.done, ctxQuery)
	defer cancel()
	return ht.store.LastHead(ctx)
}

// chainIDVerify checks whether or not the ChainID from the Chainlink config
// matches the ChainID reported by the ETH node connected to this Chainlink node.
func verifyEthereumChainID(ht *HeadTracker) error {
	ethereumChainID, err := ht.store.EthClient.ChainID(context.Background())
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
