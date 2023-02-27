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

	evmclient "github.com/smartcontractkit/chainlink/core/chains/evm/client"
	httypes "github.com/smartcontractkit/chainlink/core/chains/evm/headtracker/types"
	evmtypes "github.com/smartcontractkit/chainlink/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/core/config"
	"github.com/smartcontractkit/chainlink/core/logger"
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

type headTracker struct {
	log             logger.Logger
	headBroadcaster httypes.HeadBroadcaster
	headSaver       httypes.HeadSaver
	mailMon         *utils.MailboxMonitor
	ethClient       evmclient.Client
	chainID         big.Int
	config          Config

	backfillMB   *utils.Mailbox[*evmtypes.Head]
	broadcastMB  *utils.Mailbox[*evmtypes.Head]
	headListener httypes.HeadListener
	chStop       chan struct{}
	wgDone       sync.WaitGroup
	utils.StartStopOnce
}

// NewHeadTracker instantiates a new HeadTracker using HeadSaver to persist new block numbers.
func NewHeadTracker(
	lggr logger.Logger,
	ethClient evmclient.Client,
	config Config,
	headBroadcaster httypes.HeadBroadcaster,
	headSaver httypes.HeadSaver,
	mailMon *utils.MailboxMonitor,
) httypes.HeadTracker {
	chStop := make(chan struct{})
	lggr = lggr.Named("HeadTracker")
	id := ethClient.ChainID()
	return &headTracker{
		headBroadcaster: headBroadcaster,
		ethClient:       ethClient,
		chainID:         *id,
		config:          config,
		log:             lggr,
		backfillMB:      utils.NewSingleMailbox[*evmtypes.Head](),
		broadcastMB:     utils.NewMailbox[*evmtypes.Head](HeadsBufferSize),
		chStop:          chStop,
		headListener:    NewHeadListener(lggr, ethClient, config, chStop),
		headSaver:       headSaver,
		mailMon:         mailMon,
	}
}

// Start starts HeadTracker service.
func (ht *headTracker) Start(ctx context.Context) error {
	return ht.StartOnce("HeadTracker", func() error {
		ht.log.Debugf("Starting HeadTracker with chain id: %v", ht.chainID.Int64())
		latestChain, err := ht.headSaver.LoadFromDB(ctx)
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

		// NOTE: Always try to start the head tracker off with whatever the
		// latest head is, without waiting for the subscription to send us one.
		//
		// In some cases the subscription will send us the most recent head
		// anyway when we connect (but we should not rely on this because it is
		// not specced). If it happens this is fine, and the head will be
		// ignored as a duplicate.
		initialHead, err := ht.getInitialHead(ctx)
		if err != nil {
			if errors.Is(err, ctx.Err()) {
				return nil
			}
			ht.log.Errorw("Error getting initial head", "err", err)
		} else if initialHead != nil {
			if err := ht.handleNewHead(ctx, initialHead); err != nil {
				return errors.Wrap(err, "error handling initial head")
			}
		} else {
			ht.log.Debug("Got nil initial head")
		}

		ht.wgDone.Add(3)
		go ht.headListener.ListenForNewHeads(ht.handleNewHead, ht.wgDone.Done)
		go ht.backfillLoop()
		go ht.broadcastLoop()

		ht.mailMon.Monitor(ht.broadcastMB, "HeadTracker", "Broadcast", ht.chainID.String())

		return nil
	})
}

// Close stops HeadTracker service.
func (ht *headTracker) Close() error {
	return ht.StopOnce("HeadTracker", func() error {
		close(ht.chStop)
		ht.wgDone.Wait()
		return ht.broadcastMB.Close()
	})
}

func (ht *headTracker) Healthy() error {
	if !ht.headListener.ReceivingHeads() {
		return errors.New("Listener is not receiving heads")
	}
	if !ht.headListener.Connected() {
		return errors.New("Listener is not connected")
	}
	return nil
}

func (ht *headTracker) Name() string {
	return ht.log.Name()
}

func (ht *headTracker) HealthReport() map[string]error {
	return map[string]error{ht.Name(): ht.Healthy()}
}

func (ht *headTracker) Backfill(ctx context.Context, headWithChain *evmtypes.Head, depth uint) (err error) {
	if uint(headWithChain.ChainLength()) >= depth {
		return nil
	}

	baseHeight := headWithChain.Number - int64(depth-1)
	if baseHeight < 0 {
		baseHeight = 0
	}

	return ht.backfill(ctx, headWithChain.EarliestInChain(), baseHeight)
}

func (ht *headTracker) LatestChain() *evmtypes.Head {
	return ht.headSaver.LatestChain()
}

func (ht *headTracker) getInitialHead(ctx context.Context) (*evmtypes.Head, error) {
	head, err := ht.ethClient.HeadByNumber(ctx, nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to fetch initial head")
	}
	loggerFields := []interface{}{"head", head}
	if head != nil {
		loggerFields = append(loggerFields, "blockNumber", head.Number, "blockHash", head.Hash)
	}
	ht.log.Debugw("Got initial head", loggerFields...)
	return head, nil
}

func (ht *headTracker) handleNewHead(ctx context.Context, head *evmtypes.Head) error {
	prevHead := ht.headSaver.LatestChain()

	ht.log.Debugw(fmt.Sprintf("Received new head %v", config.FriendlyBigInt(head.ToInt())),
		"blockHeight", head.ToInt(),
		"blockHash", head.Hash,
		"parentHeadHash", head.ParentHash,
	)

	err := ht.headSaver.Save(ctx, head)
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
		ht.broadcastMB.Deliver(headWithChain)
	} else if head.Number == prevHead.Number {
		if head.Hash != prevHead.Hash {
			ht.log.Debugw("Got duplicate head", "blockNum", head.Number, "head", head.Hash.Hex(), "prevHead", prevHead.Hash.Hex())
		} else {
			ht.log.Debugw("Head already in the database", "head", head.Hash.Hex())
		}
	} else {
		ht.log.Debugw("Got out of order head", "blockNum", head.Number, "head", head.Hash.Hex(), "prevHead", prevHead.Number)
		if head.Number < prevHead.Number-int64(ht.config.EvmFinalityDepth()) {
			promOldHead.WithLabelValues(ht.chainID.String()).Inc()
			ht.log.Criticalf("Got very old block with number %d (highest seen was %d). This is a problem and either means a very deep re-org occurred, one of the RPC nodes has gotten far out of sync, or the chain went backwards in block numbers. This node may not function correctly without manual intervention.", head.Number, prevHead.Number)
		}
	}
	return nil
}

func (ht *headTracker) broadcastLoop() {
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
				item := ht.broadcastMB.RetrieveLatestAndClear()
				if item == nil {
					continue
				}
				ht.headBroadcaster.BroadcastNewLongestChain(item)
			}
		}
	} else {
		ht.log.Info("Head sampling is disabled - callback will be called on every head")
		for {
			select {
			case <-ht.chStop:
				return
			case <-ht.broadcastMB.Notify():
				for {
					item, exists := ht.broadcastMB.Retrieve()
					if !exists {
						break
					}
					ht.headBroadcaster.BroadcastNewLongestChain(item)
				}
			}
		}
	}
}

func (ht *headTracker) backfillLoop() {
	defer ht.wgDone.Done()

	ctx, cancel := utils.ContextFromChan(ht.chStop)
	defer cancel()

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
				{
					err := ht.Backfill(ctx, head, uint(ht.config.EvmFinalityDepth()))
					if err != nil {
						ht.log.Warnw("Unexpected error while backfilling heads", "err", err)
					} else if ctx.Err() != nil {
						break
					}
				}
			}
		}
	}
}

// backfill fetches all missing heads up until the base height
func (ht *headTracker) backfill(ctx context.Context, head *evmtypes.Head, baseHeight int64) (err error) {
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

func (ht *headTracker) fetchAndSaveHead(ctx context.Context, n int64) (*evmtypes.Head, error) {
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

var NullTracker httypes.HeadTracker = &nullTracker{}

type nullTracker struct{}

func (*nullTracker) Start(context.Context) error    { return nil }
func (*nullTracker) Close() error                   { return nil }
func (*nullTracker) Ready() error                   { return nil }
func (*nullTracker) Healthy() error                 { return nil }
func (*nullTracker) HealthReport() map[string]error { return map[string]error{} }
func (*nullTracker) Name() string                   { return "" }
func (*nullTracker) SetLogLevel(zapcore.Level)      {}
func (*nullTracker) Backfill(ctx context.Context, headWithChain *evmtypes.Head, depth uint) (err error) {
	return nil
}
func (*nullTracker) LatestChain() *evmtypes.Head { return nil }
