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

	// kovanChainID is the Chain ID for Kovan test network
	kovanChainID = big.NewInt(42)
)

// headRingBuffer is a small goroutine that sits between the eth client and the
// head tracker and drops the oldest head if necessary in order to keep to a fixed
// queue size (defined by the buffer size of out channel)
type headRingBuffer struct {
	in    <-chan *models.Head
	out   chan models.Head
	start sync.Once
}

func newHeadRingBuffer(in <-chan *models.Head, size int) (r *headRingBuffer, out chan models.Head) {
	out = make(chan models.Head, size)
	return &headRingBuffer{
		in:    in,
		out:   out,
		start: sync.Once{},
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
			logger.Error("HeadTracker: got nil block header")
			continue
		}
		promNumHeadsReceived.Inc()
		hInQueue := len(r.out)
		promHeadsInQueue.Set(float64(hInQueue))
		if hInQueue > 0 {
			logger.Infof("HeadTracker: Head %v is lagging behind, there are %v more heads in the queue. Your node is operating close to its maximum capacity and may start to miss jobs.", h.Number, hInQueue)
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
				logger.Errorf("HeadTracker: dropping head %v with hash 0x%x because queue is full. WARNING: Your node is overloaded and may start missing jobs.", dropped.Number, h.Hash)
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
		case blockHeader, open := <-ht.outHeaders:
			if !open {
				return errors.New("HeadTracker: outHeaders prematurely closed")
			}
			ctx, cancel := context.WithTimeout(context.Background(), ht.totalNewHeadTimeBudget())
			if err := ht.handleNewHead(ctx, blockHeader); err != nil {
				cancel()
				return err
			}
			cancel()
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
			logger.Warnw(fmt.Sprintf("HeadTracker finished processing head %v in %s which exceeds callback execution threshold of %s", number, elapsed.String(), ht.callbackExecutionThreshold().String()), "blockNumber", number, "time", elapsed, "id", "head_tracker")
		} else {
			logger.Debugw(fmt.Sprintf("HeadTracker finished processing head %v in %s", number, elapsed.String()), "blockNumber", number, "time", elapsed, "id", "head_tracker")
		}
	}(time.Now(), int64(head.Number))
	prevHead := ht.HighestSeenHead()

	logger.Debugw(fmt.Sprintf("Received new head %v", presenters.FriendlyBigInt(head.ToInt())),
		"blockHeight", head.ToInt(),
		"blockHash", head.Hash,
	)

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
	promCurrentHead.Set(float64(head.Number))
	// NOTE: We must set a hard time limit on this, backfilling heads should
	// not block the head tracker
	ctx, cancel := context.WithTimeout(context.Background(), ht.backfillTimeBudget())
	defer cancel()

	headWithChain, err := ht.GetChainWithBackfill(ctx, head, ht.store.Config.EthFinalityDepth())
	if err != nil {
		return err
	}

	ht.onNewLongestChain(ctx, headWithChain)
	return nil
}

func (ht *HeadTracker) isKovan() bool {
	return ht.store.Config.ChainID().Cmp(kovanChainID) == 0
}

// totalNewHeadTimeBudget is the timeout on the shared context for all
// requests triggered by a new head
//
// These values are chosen to be roughly 2 * block time (to give some leeway
// for temporary overload). They are by no means set in stone and may require
// adjustment based on real world feedback.
func (ht *HeadTracker) totalNewHeadTimeBudget() time.Duration {
	if ht.isKovan() {
		return 8 * time.Second
	}
	return 26 * time.Second
}

// Maximum time we are allowed to spend backfilling heads. This should be
// somewhat shorter than the average time between heads to ensure we
// don't starve the runqueue.
func (ht *HeadTracker) backfillTimeBudget() time.Duration {
	if ht.isKovan() {
		return 3 * time.Second
	}
	return 10 * time.Second
}

// If total callback execution time exceeds this threshold we consider this to
// be a problem and will log a warning.
// Here we set it to the average time between blocks.
func (ht *HeadTracker) callbackExecutionThreshold() time.Duration {
	if ht.isKovan() {
		return 4 * time.Second
	}
	return 13 * time.Second
}

// GetChainWithBackfill returns a chain of the given length, backfilling any
// heads that may be missing from the database
func (ht *HeadTracker) GetChainWithBackfill(ctx context.Context, head models.Head, depth uint) (models.Head, error) {
	ctx, cancel := context.WithTimeout(ctx, ht.backfillTimeBudget())
	defer cancel()

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
	mark := time.Now()
	fetched := 0
	defer func() {
		logger.Debugw("HeadTracker: finished backfill",
			"fetched", fetched,
			"blockNumber", head.Number,
			"time", time.Since(mark),
			"id", "head_tracker",
			"n", head.Number-baseHeight,
			"fromBlockHeight", baseHeight,
			"toBlockHeight", head.Number-1)
	}()

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
				logger.Infow("HeadTracker: backfill deadline exceeded, chain will be truncated for this head", "headNum", i)
			} else {
				logger.Errorw("HeadTracker: backfill encountered unknown error, chain will be truncated for this head", "headNum", i, "err", err)
			}
			break
		}
	}
	return nil
}

func (ht *HeadTracker) fetchAndSaveHead(ctx context.Context, n int64) (models.Head, error) {
	logger.Debugw("HeadTracker: fetching head", "blockHeight", n)
	head, err := ht.store.EthClient.HeaderByNumber(ctx, big.NewInt(n))
	if err != nil {
		return models.Head{}, err
	} else if head == nil {
		return models.Head{}, errors.New("got nil head")
	}
	if err := ht.store.IdempotentInsertHead(*head); err != nil {
		return models.Head{}, err
	}
	return *head, nil
}

func (ht *HeadTracker) onNewLongestChain(ctx context.Context, headWithChain models.Head) {
	ht.headMutex.Lock()
	defer ht.headMutex.Unlock()

	logger.Debugw("HeadTracker initiating callbacks",
		"headNum", headWithChain.Number,
		"chainLength", headWithChain.ChainLength(),
		"numCallbacks", len(ht.callbacks),
	)

	if ht.store.Config.EnableBulletproofTxManager() {
		ht.concurrentlyExecuteCallbacks(ctx, headWithChain)
	} else {
		// NOTE: Legacy tx manager probably has implicit ordering requirements, so it's not safe to parallelise
		ht.seriallyExecuteCallbacks(ctx, headWithChain)
	}
}

func (ht *HeadTracker) concurrentlyExecuteCallbacks(ctx context.Context, headWithChain models.Head) {
	wg := sync.WaitGroup{}
	wg.Add(len(ht.callbacks))
	for idx, trackable := range ht.callbacks {
		go func(i int, t strpkg.HeadTrackable) {
			start := time.Now()
			t.OnNewLongestChain(ctx, headWithChain)
			elapsed := time.Since(start)
			logger.Debugw(fmt.Sprintf("HeadTracker: finished callback %v in %s", i, elapsed), "callbackType", reflect.TypeOf(t), "callbackIdx", i, "blockNumber", headWithChain.Number, "time", elapsed, "id", "head_tracker")
			wg.Done()
		}(idx, trackable)
	}
	wg.Wait()
}

func (ht *HeadTracker) seriallyExecuteCallbacks(ctx context.Context, headWithChain models.Head) {
	for i, t := range ht.callbacks {
		start := time.Now()
		t.OnNewLongestChain(ctx, headWithChain)
		elapsed := time.Since(start)
		logger.Debugw(fmt.Sprintf("HeadTracker: finished callback %v in %s", i, elapsed), "callbackType", reflect.TypeOf(t), "callbackIdx", i, "blockNumber", headWithChain.Number, "time", elapsed, "id", "head_tracker")
	}
}

func (ht *HeadTracker) subscribeToHead() error {
	ht.headMutex.Lock()
	defer ht.headMutex.Unlock()

	ht.inHeaders = make(chan *models.Head)
	var rb *headRingBuffer
	rb, ht.outHeaders = newHeadRingBuffer(ht.inHeaders, int(ht.store.Config.EthHeadTrackerMaxBufferSize()))
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
	ethereumChainID, err := ht.store.EthClient.ChainID(context.TODO())
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
