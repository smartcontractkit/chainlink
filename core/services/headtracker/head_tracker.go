package headtracker

import (
	"context"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/chainlink/core/config"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/eth"
	httypes "github.com/smartcontractkit/chainlink/core/services/headtracker/types"
	"github.com/smartcontractkit/chainlink/core/utils"
)

var (
	promCurrentHead = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "head_tracker_current_head",
		Help: "The highest seen head number",
	}, []string{"evmChainID"})

	promOldHead = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "head_tracker_very_old_head",
		Help: "Counter is incremented every time we get a head that is much lower than the highest seen head ('much lower' is defined as a block that is ETH_FINALITY_DEPTH or greater below the highest seen head)",
	}, []string{"evmChainID"})
)

// HeadsBufferSize - The buffer is used when heads sampling is disabled, to ensure the callback is run for every head
const HeadsBufferSize = 10

// HeadTracker holds and stores the latest block number experienced by this particular node
// in a thread safe manner. Reconstitutes the last block number from the data
// store on reboot.
type HeadTracker struct {
	log             logger.Logger
	headBroadcaster httypes.HeadBroadcaster
	ethClient       eth.Client
	chainID         big.Int
	config          Config

	backfillMB   utils.Mailbox
	callbackMB   utils.Mailbox
	headListener HeadListener
	headSaver    HeadSaver
	ctx          context.Context
	cancel       context.CancelFunc
	chStop       chan struct{}
	wgDone       sync.WaitGroup
	utils.StartStopOnce
}

// NewHeadTracker instantiates a new HeadTracker using the orm to persist new block numbers.
// Can be passed in an optional sleeper object that will dictate how often
// it tries to reconnect.
func NewHeadTracker(
	l logger.Logger,
	ethClient eth.Client,
	config Config,
	orm ORM,
	headBroadcaster httypes.HeadBroadcaster,
) *HeadTracker {
	chStop := make(chan struct{})
	l = l.Named(logger.HeadTracker)
	ctx, cancel := context.WithCancel(context.Background())
	return &HeadTracker{
		headBroadcaster: headBroadcaster,
		ethClient:       ethClient,
		chainID:         *ethClient.ChainID(),
		config:          config,
		log:             l,
		backfillMB:      *utils.NewMailbox(1),
		callbackMB:      *utils.NewMailbox(HeadsBufferSize),
		ctx:             ctx,
		cancel:          cancel,
		chStop:          chStop,
		headListener:    NewHeadListener(l, ethClient, config, chStop),
		headSaver:       NewHeadSaver(l, orm, config),
	}
}

func (ht *HeadTracker) SetLogLevel(lvl zapcore.Level) {
	ht.log.SetLogLevel(lvl)
}

// Start retrieves the last persisted block number from the HeadTracker,
// subscribes to new heads, and if successful fires Connect on the
// HeadTrackable argument.
func (ht *HeadTracker) Start() error {
	return ht.StartOnce("HeadTracker", func() error {
		ht.log.Debugf("Starting HeadTracker with chain id: %v", ht.chainID.Int64())
		latestChain, err := ht.headSaver.LoadFromDB(context.Background())
		if err != nil {
			return err
		}
		if latestChain != nil {
			ht.log.Debugw(
				fmt.Sprintf("HeadTracker: Tracking logs from last block %v with hash %s", config.FriendlyBigInt(latestChain.ToInt()), latestChain.Hash.Hex()),
				"blockNumber", latestChain.Number,
				"blockHash", latestChain.Hash,
			)
		}

		// FIXME: Requests will block Start if they takes a long time. A future
		// improvement might allow Close() to cancel this context somehow.
		// https://app.shortcut.com/chainlinklabs/story/24187/ctrl-c-should-cancel-in-flight-requests-in-start-functions
		startCtx := context.Background()

		// NOTE: Always try to start the head tracker off with whatever the
		// latest head is, without waiting for the subscription to send us one.
		//
		// In some cases the subscription will send us the most recent head
		// anyway when we connect (but we should not rely on this because it is
		// not specced). If it happens this is fine, and the head will be
		// ignored as a duplicate.
		initialHead, err := ht.getInitialHead(startCtx)
		if err != nil {
			ht.log.Errorw("Error getting initial head", "err", err)
		} else if initialHead != nil {
			if err := ht.handleNewHead(startCtx, initialHead); err != nil {
				return errors.Wrap(err, "error handling initial head")
			}
		} else {
			ht.log.Debug("Got nil initial head")
		}

		ht.wgDone.Add(3)
		go ht.headListener.ListenForNewHeads(ht.handleNewHead, ht.wgDone.Done)
		go ht.backfiller()
		go ht.headCallbackLoop()

		return nil
	})
}

func (ht *HeadTracker) getInitialHead(ctx context.Context) (*eth.Head, error) {
	head, err := ht.ethClient.HeadByNumber(ctx, nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to fetch initial head")
	}
	loggerFields := []interface{}{"head", head}
	if head != nil {
		loggerFields = append(loggerFields, "blockNumber", head.Number, "blockHash", head.Hash)
	}
	ht.log.Debugw("Got initial current block", loggerFields...)
	return head, nil
}

// Stop unsubscribes all connections and fires Disconnect.
func (ht *HeadTracker) Stop() error {
	return ht.StopOnce("HeadTracker", func() error {
		ht.cancel()
		close(ht.chStop)
		ht.wgDone.Wait()
		return nil
	})
}

func (ht *HeadTracker) Save(ctx context.Context, h *eth.Head) error {
	return ht.headSaver.Save(ctx, h)
}

func (ht *HeadTracker) LatestChain() *eth.Head {
	return ht.headSaver.LatestChain()
}

func (ht *HeadTracker) HighestSeenHeadFromDB(ctx context.Context) (*eth.Head, error) {
	return ht.headSaver.LatestHeadFromDB(ctx)
}

// Connected returns whether or not this HeadTracker is connected.
func (ht *HeadTracker) Connected() bool {
	return ht.headListener.Connected()
}

func (ht *HeadTracker) headCallbackLoop() {
	defer ht.wgDone.Done()

	samplingInterval := ht.config.EvmHeadTrackerSamplingInterval()
	if samplingInterval > 0 {
		ht.log.Debugf("Head sampling is enabled - sampling interval is set to: %v", samplingInterval)
		debounceHead := time.NewTicker(samplingInterval)
		defer debounceHead.Stop()
		for {
			select {
			case <-ht.chStop:
				return
			case <-debounceHead.C:
				item := ht.callbackMB.RetrieveLatestAndClear()
				if item == nil {
					continue
				}
				ht.callbackOnLatestHead(item)
			}
		}
	} else {
		ht.log.Info("Head sampling is disabled - callback will be called on every head")
		for {
			select {
			case <-ht.chStop:
				return
			case <-ht.callbackMB.Notify():
				for {
					item, exists := ht.callbackMB.Retrieve()
					if !exists {
						break
					}
					ht.callbackOnLatestHead(item)
				}
			}
		}
	}
}

func (ht *HeadTracker) callbackOnLatestHead(item interface{}) {
	head := eth.AsHead(item)

	ht.headBroadcaster.OnNewLongestChain(ht.ctx, head)
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
				h := eth.AsHead(head)
				{
					err := ht.Backfill(ht.ctx, h, uint(ht.config.EvmFinalityDepth()))
					if err != nil {
						ht.log.Warnw("Unexpected error while backfilling heads", "err", err)
					} else if ht.ctx.Err() != nil {
						break
					}
				}
			}
		}
	}
}

// Backfill given a head will fill in any missing heads up to the given depth
func (ht *HeadTracker) Backfill(ctx context.Context, headWithChain *eth.Head, depth uint) (err error) {
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
func (ht *HeadTracker) backfill(ctx context.Context, head *eth.Head, baseHeight int64) (err error) {
	if head.Number <= baseHeight {
		return nil
	}
	mark := time.Now()
	fetched := 0
	l := ht.log.With("blockNumber", head.Number,
		"n", head.Number-baseHeight,
		"fromBlockHeight", baseHeight,
		"toBlockHeight", head.Number-1)
	l.Debug("Starting backfill")
	defer func() {
		if ctx.Err() != nil {
			l.Warnw("Backfill context error", "err", ctx.Err())
			return
		}
		l.Debugw("Finished backfill",
			"fetched", fetched,
			"time", time.Since(mark),
			"err", err)
	}()

	for i := head.Number - 1; i >= baseHeight; i-- {
		// NOTE: Sequential requests here mean it's a potential performance bottleneck, be aware!
		existingHead := ht.headSaver.Chain(head.ParentHash)
		if existingHead != nil {
			head = existingHead
			continue
		}
		head, err = ht.fetchAndSaveHead(ctx, i)
		fetched++
		if ctx.Err() != nil {
			ht.log.Debugw("context canceled, aborting backfill", "err", err, "ctx.Err", ctx.Err())
			break
		} else if err != nil {
			return errors.Wrap(err, "fetchAndSaveHead failed")
		}
	}
	return
}

func (ht *HeadTracker) fetchAndSaveHead(ctx context.Context, n int64) (*eth.Head, error) {
	ht.log.Debugw("Fetching head", "blockHeight", n)
	head, err := ht.ethClient.HeadByNumber(ctx, big.NewInt(n))
	if err != nil {
		return nil, err
	} else if head == nil {
		return nil, errors.New("got nil head")
	}
	err = ht.headSaver.Save(ctx, head)
	if err != nil {
		return nil, err
	}
	return head, nil
}

func (ht *HeadTracker) handleNewHead(ctx context.Context, head *eth.Head) error {
	prevHead := ht.LatestChain()

	ht.log.Debugw(fmt.Sprintf("Received new head %v", config.FriendlyBigInt(head.ToInt())),
		"blockHeight", head.ToInt(),
		"blockHash", head.Hash,
		"parentHeadHash", head.ParentHash,
	)

	err := ht.Save(ctx, head)
	if ctx.Err() != nil {
		return nil
	} else if err != nil {
		return errors.Wrapf(err, "failed to save head: %#v", head)
	}

	if prevHead == nil || head.Number > prevHead.Number {
		promCurrentHead.WithLabelValues(ht.chainID.String()).Set(float64(head.Number))

		headWithChain := ht.headSaver.Chain(head.Hash)
		if headWithChain == nil {
			return errors.Errorf("HeadTracker#handleNewHighestHead headWithChain was unexpectedly nil")
		}
		ht.backfillMB.Deliver(headWithChain)
		ht.callbackMB.Deliver(headWithChain)
		return nil
	}
	if head.Number == prevHead.Number {
		if head.Hash != prevHead.Hash {
			ht.log.Debugw("Got duplicate head", "blockNum", head.Number, "gotHead", head.Hash.Hex(), "highestSeenHead", prevHead.Hash.Hex())
		} else {
			ht.log.Debugw("Head already in the database", "gotHead", head.Hash.Hex())
		}
	} else {
		ht.log.Debugw("Got out of order head", "blockNum", head.Number, "gotHead", head.Hash.Hex(), "highestSeenHead", prevHead.Number)
		if head.Number < prevHead.Number-int64(ht.config.EvmFinalityDepth()) {
			promOldHead.WithLabelValues(ht.chainID.String()).Inc()
			ht.log.Errorf("Got very old block with number %d (highest seen was %d). This is a problem and either means a very deep re-org occurred, or the chain went backwards in block numbers. This node will not function correctly without manual intervention.", head.Number, prevHead.Number)
		}
	}
	return nil
}

func (ht *HeadTracker) Healthy() error {
	if !ht.headListener.ReceivingHeads() {
		return errors.New("Listener is not receiving heads")
	}
	if !ht.headListener.Connected() {
		return errors.New("Listener is not connected")
	}
	return nil
}

// Saver returns HeadSaver instance, exposed for testing.
// Consider removing this while refactoring HeadTracker.
func (ht *HeadTracker) Saver() HeadSaver {
	return ht.headSaver
}

var _ httypes.Tracker = &NullTracker{}

type NullTracker struct{}

func (n *NullTracker) HighestSeenHeadFromDB(context.Context) (*eth.Head, error) {
	return nil, nil
}
func (*NullTracker) Start() error   { return nil }
func (*NullTracker) Stop() error    { return nil }
func (*NullTracker) Ready() error   { return nil }
func (*NullTracker) Healthy() error { return nil }

func (*NullTracker) SetLogLevel(zapcore.Level) {}
