package services

import (
	"context"
	"fmt"
	"math/big"
	"reflect"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/headtracker"
	httypes "github.com/smartcontractkit/chainlink/core/services/headtracker/types"
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

// HeadTracker holds and stores the latest block number experienced by this particular node
// in a thread safe manner. Reconstitutes the last block number from the data
// store on reboot.
type HeadTracker struct {
	log       *logger.Logger
	callbacks []httypes.HeadTrackable
	store     *strpkg.Store

	backfillMB   utils.Mailbox
	samplingMB   utils.Mailbox
	muLogger     sync.RWMutex
	headListener *headtracker.HeadListener
	headSaver    *headtracker.HeadSaver

	chStop chan struct{}
	wgDone *sync.WaitGroup
	utils.StartStopOnce
}

// NewHeadTracker instantiates a new HeadTracker using the orm to persist new block numbers.
// Can be passed in an optional sleeper object that will dictate how often
// it tries to reconnect.
func NewHeadTracker(l *logger.Logger, store *strpkg.Store, callbacks []httypes.HeadTrackable, sleepers ...utils.Sleeper) *HeadTracker {

	var wgDone sync.WaitGroup
	chStop := make(chan struct{})

	return &HeadTracker{
		store:        store,
		callbacks:    callbacks,
		log:          l,
		backfillMB:   *utils.NewMailbox(1),
		samplingMB:   *utils.NewMailbox(1),
		chStop:       chStop,
		wgDone:       &wgDone,
		headListener: headtracker.NewHeadListener(l, store.EthClient, store.Config, chStop, &wgDone, sleepers...),
		headSaver:    headtracker.NewHeadSaver(store),
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
	return ht.StartOnce("HeadTracker", func() error {
		ht.logger().Info("Starting HeadTracker")
		highestSeenHead, err := ht.headSaver.SetHighestSeenHeadFromDB()
		if err != nil {
			return err
		}
		if highestSeenHead != nil {
			ht.logger().Debugw(
				fmt.Sprintf("HeadTracker: Tracking logs from last block %v with hash %s", presenters.FriendlyBigInt(highestSeenHead.ToInt()), highestSeenHead.Hash.Hex()),
				"blockNumber", highestSeenHead.Number,
				"blockHash", highestSeenHead.Hash,
			)
		}

		ht.wgDone.Add(3)
		go ht.headListener.ListenForNewHeads(ht.handleNewHead, ht.handleConnected)
		go ht.backfiller()
		go ht.headSampler()

		return nil
	})
}

// Stop unsubscribes all connections and fires Disconnect.
func (ht *HeadTracker) Stop() error {
	return ht.StopOnce("HeadTracker", func() error {
		ht.logger().Info(fmt.Sprintf("HeadTracker disconnecting from %v", ht.store.Config.EthereumURL()))
		close(ht.chStop)
		ht.wgDone.Wait()
		return nil
	})
}

func (ht *HeadTracker) Save(ctx context.Context, h models.Head) error {
	return ht.headSaver.Save(ctx, h)
}

func (ht *HeadTracker) HighestSeenHead() *models.Head {
	return ht.headSaver.HighestSeenHead()
}

func (ht *HeadTracker) HighestSeenHeadFromDB() (*models.Head, error) {
	return ht.headSaver.HighestSeenHeadFromDB()
}

func (ht *HeadTracker) Chain(ctx context.Context, hash common.Hash, depth uint) (models.Head, error) {
	return ht.headSaver.Chain(ctx, hash, depth)
}

// Connected returns whether or not this HeadTracker is connected.
func (ht *HeadTracker) Connected() bool {
	return ht.headListener.Connected()
}

func (ht *HeadTracker) handleConnected() {
	ht.connect(ht.headSaver.HighestSeenHead())
}

func (ht *HeadTracker) connect(bn *models.Head) {
	for _, trackable := range ht.callbacks {
		if err := trackable.Connect(bn); err != nil {
			ht.logger().Warn(trackable.Connect(bn))
		}
	}
}

func (ht *HeadTracker) headSampler() {
	defer ht.wgDone.Done()

	debounceHead := time.NewTicker(ht.store.Config.EthHeadTrackerSamplingInterval())
	defer debounceHead.Stop()

	ctx, cancel := utils.ContextFromChan(ht.chStop)
	defer cancel()

	for {
		select {
		case <-ht.chStop:
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

func (ht *HeadTracker) backfiller() {
	defer ht.wgDone.Done()
	for {
		select {
		case <-ht.chStop:
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
					ctx, cancel := utils.ContextFromChan(ht.chStop)
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
	err = ht.headSaver.IdempotentInsertHead(ctx, *head)
	if ctx.Err() != nil {
		return models.Head{}, nil
	} else if err != nil {
		return models.Head{}, err
	}
	return *head, nil
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
			ht.logger().Debugw("HeadTracker: got duplicate head", "blockNum", head.Number, "gotHead", head.Hash.Hex(), "highestSeenHead", prevHead.Hash.Hex())
		} else {
			ht.logger().Debugw("HeadTracker: head already in the database", "gotHead", head.Hash.Hex())
		}
	} else {
		ht.logger().Debugw("HeadTracker: got out of order head", "blockNum", head.Number, "gotHead", head.Hash.Hex(), "highestSeenHead", prevHead.Number)
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
		go func(i int, t httypes.HeadTrackable) {
			start := time.Now()
			t.OnNewLongestChain(ctx, headWithChain)
			elapsed := time.Since(start)
			ht.logger().Debugw(fmt.Sprintf("HeadTracker: finished callback %v in %s", i, elapsed), "callbackType", reflect.TypeOf(t), "callbackIdx", i, "blockNumber", headWithChain.Number, "time", elapsed, "id", "head_tracker")
			wg.Done()
		}(idx, trackable)
	}
	wg.Wait()
}
