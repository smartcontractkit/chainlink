package services

import (
	"context"
	"fmt"
	"math/big"
	"reflect"
	"sync"
	"time"

	"github.com/smartcontractkit/chainlink/core/logger"
	strpkg "github.com/smartcontractkit/chainlink/core/store"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/store/presenters"
	"github.com/smartcontractkit/chainlink/core/utils"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
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
	promNumHeadsDropped = promauto.NewCounter(prometheus.CounterOpts{
		Name: "head_tracker_num_heads_dropped",
		Help: "The total number of heads dropped",
	})
	promEthConnectionErrors = promauto.NewCounter(prometheus.CounterOpts{
		Name: "head_tracker_eth_connection_errors",
		Help: "The total number of eth node connection errors",
	})
)

// headRingBuffer is a small goroutine that sits between the eth client and the
// head tracker and drops the oldest head if necessary in order to keep to a fixed
// queue size (defined by the buffer size of out channel)
type headRingBuffer struct {
	in     <-chan *models.Head
	out    chan models.Head
	start  sync.Once
	logger *logger.Logger
}

func newHeadRingBuffer(in <-chan *models.Head, size int, logger *logger.Logger) (r *headRingBuffer, out chan models.Head) {
	out = make(chan models.Head, size)
	return &headRingBuffer{
		in:     in,
		out:    out,
		start:  sync.Once{},
		logger: logger,
	}, out
}

// Start the headRingBuffer goroutine
// It will be stopped implicitly by closing the in channel
func (r *headRingBuffer) Start() {
	r.start.Do(func() {
		go r.run()
	})
}

func (r *headRingBuffer) run() {
	for h := range r.in {
		if h == nil {
			r.logger.Error("HeadTracker: got nil block header")
			continue
		}
		promNumHeadsReceived.Inc()
		hInQueue := len(r.out)
		promHeadsInQueue.Set(float64(hInQueue))
		if hInQueue > 0 {
			r.logger.Infof("HeadTracker: Head %v is lagging behind, there are %v more heads in the queue. Your node is operating close to its maximum capacity and may start to miss jobs.", h.Number, hInQueue)
		}
		select {
		case r.out <- *h:
		default:
			// Need to select/default here because it's conceivable (although
			// improbable) that between the previous select and now, all heads were drained
			// from r.out by another goroutine
			//
			// NOTE: In this unlikely event, we may drop an extra head unnecessarily.
			// The probability of this seems vanishingly small, and only hits
			// if the queue was already full anyway, so we can live with this
			select {
			case dropped := <-r.out:
				promNumHeadsDropped.Inc()
				r.logger.Errorf("HeadTracker: dropping head %v with hash 0x%x because queue is full. WARNING: Your node is overloaded and may start missing jobs.", dropped.Number, h.Hash)
				r.out <- *h
			default:
				r.out <- *h
			}
		}
	}
	close(r.out)
}

// HeadTracker holds and stores the latest block number experienced by this particular node
// in a thread safe manner. Reconstitutes the last block number from the data
// store on reboot.
type HeadTracker struct {
	callbacks             []strpkg.HeadTrackable
	inHeaders             chan *models.Head
	outHeaders            chan models.Head
	headSubscription      ethereum.Subscription
	highestSeenHead       *models.Head
	store                 *strpkg.Store
	headMutex             sync.RWMutex
	connected             bool
	sleeper               utils.Sleeper
	done                  chan struct{}
	started               bool
	listenForNewHeadsWg   sync.WaitGroup
	backfillMB            utils.Mailbox
	subscriptionSucceeded chan struct{}
	logger                *logger.Logger
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
	l := logger.CreateLogger(logger.Default.With(
		"id", "head_tracker",
	))
	return &HeadTracker{
		store:      store,
		callbacks:  callbacks,
		sleeper:    sleeper,
		logger:     l,
		backfillMB: *utils.NewMailbox(1),
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
		ht.logger.Debugw(
			fmt.Sprintf("Headtracker: Tracking logs from last block %v with hash %s", presenters.FriendlyBigInt(ht.highestSeenHead.ToInt()), ht.highestSeenHead.Hash.Hex()),
			"blockNumber", ht.highestSeenHead.Number,
			"blockHash", ht.highestSeenHead.Hash,
		)
	}

	ht.done = make(chan struct{})
	ht.subscriptionSucceeded = make(chan struct{})

	ht.listenForNewHeadsWg.Add(2)
	go ht.listenForNewHeads()
	go ht.backfiller()

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

	ht.logger.Info(fmt.Sprintf("Head tracker disconnecting from %v", ht.store.Config.EthereumURL()))
	close(ht.done)
	close(ht.subscriptionSucceeded)
	ht.started = false
	ht.headMutex.Unlock()

	ht.listenForNewHeadsWg.Wait()
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
		ht.logger.WarnIf(trackable.Connect(bn))
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
		logger.WarnIf(errors.Wrap(err, "HeadTracker failed when unsubscribe from head"))
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
			ht.logger.Errorw(fmt.Sprintf("Error in new head subscription, unsubscribed: %s", err.Error()), "err", err)
			continue
		} else {
			break
		}
	}
}

func (ht *HeadTracker) backfiller() {
	defer ht.listenForNewHeadsWg.Done()
	for {
		select {
		case <-ht.done:
			return
		case <-ht.backfillMB.Notify():
			for {
				head := ht.backfillMB.Retrieve()
				if head == nil {
					break
				}
				h, is := head.(models.Head)
				if !is {
					logger.Errorf("HeadTracker: invariant violation, expected %T but got %T", models.Head{}, head)
					continue
				}
				{
					ctx, cancel := utils.ContextFromChan(ht.done)
					err := ht.Backfill(ctx, h, ht.store.Config.EthFinalityDepth())
					defer cancel()
					if err != nil {
						logger.Warnw("HeadTracker: unexpected error while backfilling heads", "err", err)
					} else if ctx.Err() != nil {
						break
					}
				}
			}
		}
	}
}

// Backfill given a head will fill in any missing heads up to the given depth
func (ht *HeadTracker) Backfill(ctx context.Context, head models.Head, depth uint) (err error) {
	if uint(head.ChainLength()) >= depth {
		return nil
	}
	baseHeight := head.Number - int64(depth-1)
	if baseHeight < 0 {
		baseHeight = 0
	}

	return ht.backfill(ctx, head.EarliestInChain(), baseHeight)
}

// backfill fetches all missing heads up until the base height
func (ht *HeadTracker) backfill(ctx context.Context, head models.Head, baseHeight int64) (err error) {
	if head.Number <= baseHeight {
		return nil
	}
	mark := time.Now()
	fetched := 0
	logger.Debugw("HeadTracker: starting backfill",
		"blockNumber", head.Number,
		"id", "head_tracker",
		"n", head.Number-baseHeight,
		"fromBlockHeight", baseHeight,
		"toBlockHeight", head.Number-1)
	defer func() {
		if ctx.Err() != nil {
			return
		}
		logger.Debugw("HeadTracker: finished backfill",
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
	logger.Debugw("HeadTracker: fetching head", "blockHeight", n)
	head, err := ht.store.EthClient.HeaderByNumber(ctx, big.NewInt(n))
	if ctx.Err() != nil {
		return models.Head{}, nil
	} else if err != nil {
		return models.Head{}, err
	} else if head == nil {
		return models.Head{}, errors.New("got nil head")
	}
	if err := ht.store.IdempotentInsertHead(ctx, *head); err != nil {
		return models.Head{}, err
	}
	return *head, nil
}

// subscribe periodically attempts to connect to the ethereum node via websocket.
// It returns true on success, and false if cut short by a done request and did not connect.
func (ht *HeadTracker) subscribe() bool {
	ht.sleeper.Reset()
	for {
		err := ht.unsubscribeFromHead()
		if err != nil {
			ht.logger.ErrorIf(err, "failed when unsubscribe from head")
			return false
		}

		ht.logger.Info("HeadTracker: Connecting to ethereum node ", ht.store.Config.EthereumURL(), " in ", ht.sleeper.Duration())
		select {
		case <-ht.done:
			return false
		case <-time.After(ht.sleeper.After()):
			err := ht.subscribeToHead()
			if err != nil {
				promEthConnectionErrors.Inc()
				ht.logger.Warnw(fmt.Sprintf("HeadTracker: Failed to connect to ethereum node %v", ht.store.Config.EthereumURL()), "err", err)
			} else {
				ht.logger.Info("HeadTracker: Connected to ethereum node ", ht.store.Config.EthereumURL())
				return true
			}
		}
	}
}

// This should be safe to run concurrently across multiple nodes connected to the same database
func (ht *HeadTracker) receiveHeaders(ctx context.Context) error {
	for {
		select {
		case <-ht.done:
			return nil
		case blockHeader, open := <-ht.outHeaders:
			if !open {
				return errors.New("HeadTracker: outHeaders prematurely closed")
			}
			timeBudget := ht.store.Config.HeadTimeBudget()
			{
				deadlineCtx, cancel := context.WithTimeout(ctx, timeBudget)
				defer cancel()

				err := ht.handleNewHead(ctx, blockHeader)
				if ctx.Err() != nil {
					return nil
				} else if deadlineCtx.Err() != nil {
					logger.Warnw("HeadTracker: handling of new head timed out", "error", ctx.Err(), "timeBudget", timeBudget.String())
					return err
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
	defer func(start time.Time, number int64) {
		elapsed := time.Since(start)
		ms := float64(elapsed.Milliseconds())
		promCallbackDuration.Set(ms)
		promCallbackDurationHist.Observe(ms)
		if elapsed > ht.callbackExecutionThreshold() {
			ht.logger.Warnw(fmt.Sprintf("HeadTracker finished processing head %v in %s which exceeds callback execution threshold of %s", number, elapsed.String(), ht.store.Config.HeadTimeBudget().String()), "blockNumber", number, "time", elapsed, "id", "head_tracker")
		} else {
			ht.logger.Debugw(fmt.Sprintf("HeadTracker finished processing head %v in %s", number, elapsed.String()), "blockNumber", number, "time", elapsed, "id", "head_tracker")
		}
	}(time.Now(), int64(head.Number))
	prevHead := ht.HighestSeenHead()

	ht.logger.Debugw(fmt.Sprintf("HeadTracker: Received new head %v", presenters.FriendlyBigInt(head.ToInt())),
		"blockHeight", head.ToInt(),
		"blockHash", head.Hash,
	)

	err := ht.Save(ctx, head)
	if ctx.Err() != nil {
		return nil
	} else if err != nil {
		return err
	}

	if prevHead == nil || head.Number > prevHead.Number {
		return ht.handleNewHighestHead(ctx, head)
	}
	if head.Number == prevHead.Number {
		if head.Hash != prevHead.Hash {
			ht.logger.Debugw("HeadTracker: got duplicate head", "blockNum", head.Number, "gotHead", head.Hash.Hex(), "highestSeenHead", ht.highestSeenHead.Hash.Hex())
		} else {
			ht.logger.Debugw("HeadTracker: head already in the database", "gotHead", head.Hash.Hex())
		}
	} else {
		ht.logger.Debugw("HeadTracker: got out of order head", "blockNum", head.Number, "gotHead", head.Hash.Hex(), "highestSeenHead", ht.highestSeenHead.Number)
	}
	return nil
}

func (ht *HeadTracker) handleNewHighestHead(ctx context.Context, head models.Head) error {
	promCurrentHead.Set(float64(head.Number))

	headWithChain, err := ht.store.Chain(ctx, head.Hash, ht.store.Config.EthFinalityDepth())
	if ctx.Err() != nil {
		return nil
	} else if err != nil {
		return errors.Wrap(err, "HeadTracker#handleNewHighestHead failed fetching chain")
	}
	ht.backfillMB.Deliver(headWithChain)

	ht.onNewLongestChain(ctx, headWithChain)
	return nil
}

// If total callback execution time exceeds this threshold we consider this to
// be a problem and will log a warning.
// Here we set it to the average time between blocks.
func (ht *HeadTracker) callbackExecutionThreshold() time.Duration {
	return ht.store.Config.HeadTimeBudget() / 2
}

func (ht *HeadTracker) onNewLongestChain(ctx context.Context, headWithChain models.Head) {
	ht.headMutex.Lock()
	defer ht.headMutex.Unlock()

	ht.logger.Debugw("HeadTracker initiating callbacks",
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
			ht.logger.Debugw(fmt.Sprintf("HeadTracker: finished callback %v in %s", i, elapsed), "callbackType", reflect.TypeOf(t), "callbackIdx", i, "blockNumber", headWithChain.Number, "time", elapsed, "id", "head_tracker")
			wg.Done()
		}(idx, trackable)
	}
	wg.Wait()
}

func (ht *HeadTracker) subscribeToHead() error {
	ht.headMutex.Lock()
	defer ht.headMutex.Unlock()

	ht.inHeaders = make(chan *models.Head)
	var rb *headRingBuffer
	rb, ht.outHeaders = newHeadRingBuffer(ht.inHeaders, int(ht.store.Config.EthHeadTrackerMaxBufferSize()), ht.logger)
	// It will autostop when we close inHeaders channel
	rb.Start()

	sub, err := ht.store.EthClient.SubscribeNewHead(context.Background(), ht.inHeaders)
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
	close(ht.inHeaders)
	// Drain channel and wait for ringbuffer to close it
	for range ht.outHeaders {
	}
	return nil
}

func (ht *HeadTracker) setHighestSeenHeadFromDB() error {
	head, err := ht.store.LastHead(context.Background())
	if err != nil {
		return err
	}
	ht.highestSeenHead = head
	return nil
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
