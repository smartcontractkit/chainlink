package headmanager

import (
	"context"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"golang.org/x/exp/maps"

	hmtypes "github.com/smartcontractkit/chainlink/v2/common/headmanager/types"
	"github.com/smartcontractkit/chainlink/v2/common/types"
	"github.com/smartcontractkit/chainlink/v2/core/config"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

var (
	promCurrentHead = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "head_tracker_current_head",
		Help: "The highest seen head number",
	}, []string{"evmChainID"})

	promOldHead = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "head_tracker_very_old_head",
		Help: "Counter is incremented every time we get a head that is much lower than the highest seen head ('much lower' is defined as a block that is EVM.FinalityDepth or greater below the highest seen head)",
	}, []string{"evmChainID"})
)

// HeadsBufferSize - The buffer is used when heads sampling is disabled, to ensure the callback is run for every head
const HeadsBufferSize = 10

type Tracker[
	HMH hmtypes.Head[BLOCK_HASH, ID],
	S types.Subscription,
	ID types.ID,
	BLOCK_HASH types.Hashable,
] struct {
	log             logger.Logger
	headBroadcaster types.Broadcaster[HMH, BLOCK_HASH]
	headSaver       types.Saver[HMH, BLOCK_HASH]
	mailMon         *utils.MailboxMonitor
	client          hmtypes.Client[HMH, S, ID, BLOCK_HASH]
	chainID         ID
	config          hmtypes.Config
	htConfig        hmtypes.HeadTrackerConfig

	backfillMB   *utils.Mailbox[HMH]
	broadcastMB  *utils.Mailbox[HMH]
	headListener types.Listener[HMH, BLOCK_HASH]
	chStop       utils.StopChan
	wgDone       sync.WaitGroup
	utils.StartStopOnce
	getNilHead func() HMH
}

// NewTracker instantiates a new Tracker using HeadSaver to persist new block numbers.
func NewTracker[
	HMH hmtypes.Head[BLOCK_HASH, ID],
	S types.Subscription,
	ID types.ID,
	BLOCK_HASH types.Hashable,
](
	lggr logger.Logger,
	client hmtypes.Client[HMH, S, ID, BLOCK_HASH],
	config hmtypes.Config,
	htConfig hmtypes.HeadTrackerConfig,
	headBroadcaster types.Broadcaster[HMH, BLOCK_HASH],
	headSaver types.Saver[HMH, BLOCK_HASH],
	mailMon *utils.MailboxMonitor,
	getNilHead func() HMH,
) *Tracker[HMH, S, ID, BLOCK_HASH] {
	chStop := make(chan struct{})
	lggr = lggr.Named("Tracker")
	return &Tracker[HMH, S, ID, BLOCK_HASH]{
		headBroadcaster: headBroadcaster,
		client:          client,
		chainID:         client.ConfiguredChainID(),
		config:          config,
		htConfig:        htConfig,
		log:             lggr,
		backfillMB:      utils.NewSingleMailbox[HMH](),
		broadcastMB:     utils.NewMailbox[HMH](HeadsBufferSize),
		chStop:          chStop,
		headListener:    NewListener[HMH, S, ID, BLOCK_HASH](lggr, client, config, chStop),
		headSaver:       headSaver,
		mailMon:         mailMon,
		getNilHead:      getNilHead,
	}
}

// Start starts Tracker service.
func (t *Tracker[HMH, S, ID, BLOCK_HASH]) Start(ctx context.Context) error {
	return t.StartOnce("Tracker", func() error {
		t.log.Debugw("Starting Tracker", "chainID", t.chainID)
		latestChain, err := t.headSaver.Load(ctx)
		if err != nil {
			return err
		}
		if latestChain.IsValid() {
			t.log.Debugw(
				fmt.Sprintf("Tracker: Tracking logs from last block %v with hash %s", config.FriendlyNumber(latestChain.BlockNumber()), latestChain.BlockHash()),
				"blockNumber", latestChain.BlockNumber(),
				"blockHash", latestChain.BlockHash(),
			)
		}

		// NOTE: Always try to start the head tracker off with whatever the
		// latest head is, without waiting for the subscription to send us one.
		//
		// In some cases the subscription will send us the most recent head
		// anyway when we connect (but we should not rely on this because it is
		// not specced). If it happens this is fine, and the head will be
		// ignored as a duplicate.
		initialHead, err := t.getInitialHead(ctx)
		if err != nil {
			if errors.Is(err, ctx.Err()) {
				return nil
			}
			t.log.Errorw("Error getting initial head", "err", err)
		} else if initialHead.IsValid() {
			if err := t.handleNewHead(ctx, initialHead); err != nil {
				return errors.Wrap(err, "error handling initial head")
			}
		} else {
			t.log.Debug("Got nil initial head")
		}

		t.wgDone.Add(3)
		go t.headListener.ListenForNewHeads(t.handleNewHead, t.wgDone.Done)
		go t.backfillLoop()
		go t.broadcastLoop()

		t.mailMon.Monitor(t.broadcastMB, "Tracker", "Broadcast", t.chainID.String())

		return nil
	})
}

// Close stops Tracker service.
func (t *Tracker[HMH, S, ID, BLOCK_HASH]) Close() error {
	return t.StopOnce("Tracker", func() error {
		close(t.chStop)
		t.wgDone.Wait()
		return t.broadcastMB.Close()
	})
}

func (t *Tracker[HMH, S, ID, BLOCK_HASH]) Name() string {
	return t.log.Name()
}

func (t *Tracker[HMH, S, ID, BLOCK_HASH]) HealthReport() map[string]error {
	report := map[string]error{
		t.Name(): t.StartStopOnce.Healthy(),
	}
	maps.Copy(report, t.headListener.HealthReport())
	return report
}

func (t *Tracker[HMH, S, ID, BLOCK_HASH]) Backfill(ctx context.Context, headWithChain HMH, depth uint) (err error) {
	if uint(headWithChain.ChainLength()) >= depth {
		return nil
	}

	baseHeight := headWithChain.BlockNumber() - int64(depth-1)
	if baseHeight < 0 {
		baseHeight = 0
	}

	return t.backfill(ctx, headWithChain.EarliestHeadInChain(), baseHeight)
}

func (t *Tracker[HMH, S, ID, BLOCK_HASH]) LatestChain() HMH {
	return t.headSaver.LatestChain()
}

func (t *Tracker[HMH, S, ID, BLOCK_HASH]) getInitialHead(ctx context.Context) (HMH, error) {
	head, err := t.client.HeadByNumber(ctx, nil)
	if err != nil {
		return t.getNilHead(), errors.Wrap(err, "failed to fetch initial head")
	}
	loggerFields := []interface{}{"head", head}
	if head.IsValid() {
		loggerFields = append(loggerFields, "blockNumber", head.BlockNumber(), "blockHash", head.BlockHash())
	}
	t.log.Debugw("Got initial head", loggerFields...)
	return head, nil
}

func (t *Tracker[HMH, S, ID, BLOCK_HASH]) handleNewHead(ctx context.Context, head HMH) error {
	prevHead := t.headSaver.LatestChain()

	t.log.Debugw(fmt.Sprintf("Received new head %v", config.FriendlyNumber(head.BlockNumber())),
		"blockHeight", head.BlockNumber(),
		"blockHash", head.BlockHash(),
		"parentHeadHash", head.GetParentHash(),
	)

	err := t.headSaver.Save(ctx, head)
	if ctx.Err() != nil {
		return nil
	} else if err != nil {
		return errors.Wrapf(err, "failed to save head: %#v", head)
	}

	if !prevHead.IsValid() || head.BlockNumber() > prevHead.BlockNumber() {
		promCurrentHead.WithLabelValues(t.chainID.String()).Set(float64(head.BlockNumber()))

		headWithChain := t.headSaver.Chain(head.BlockHash())
		if !headWithChain.IsValid() {
			return errors.Errorf("Tracker#handleNewHighestHead headWithChain was unexpectedly nil")
		}
		t.backfillMB.Deliver(headWithChain)
		t.broadcastMB.Deliver(headWithChain)
	} else if head.BlockNumber() == prevHead.BlockNumber() {
		if head.BlockHash() != prevHead.BlockHash() {
			t.log.Debugw("Got duplicate head", "blockNum", head.BlockNumber(), "head", head.BlockHash(), "prevHead", prevHead.BlockHash())
		} else {
			t.log.Debugw("Head already in the database", "head", head.BlockHash())
		}
	} else {
		t.log.Debugw("Got out of order head", "blockNum", head.BlockNumber(), "head", head.BlockHash(), "prevHead", prevHead.BlockNumber())
		prevUnFinalizedHead := prevHead.BlockNumber() - int64(t.config.FinalityDepth())
		if head.BlockNumber() < prevUnFinalizedHead {
			promOldHead.WithLabelValues(t.chainID.String()).Inc()
			t.log.Criticalf("Got very old block with number %d (highest seen was %d). This is a problem and either means a very deep re-org occurred, one of the RPC nodes has gotten far out of sync, or the chain went backwards in block numbers. This node may not function correctly without manual intervention.", head.BlockNumber(), prevHead.BlockNumber())
			t.SvcErrBuffer.Append(errors.New("got very old block"))
		}
	}
	return nil
}

func (t *Tracker[HMH, S, ID, BLOCK_HASH]) broadcastLoop() {
	defer t.wgDone.Done()

	samplingInterval := t.htConfig.SamplingInterval()
	if samplingInterval > 0 {
		t.log.Debugf("Head sampling is enabled - sampling interval is set to: %v", samplingInterval)
		debounceHead := time.NewTicker(samplingInterval)
		defer debounceHead.Stop()
		for {
			select {
			case <-t.chStop:
				return
			case <-debounceHead.C:
				item := t.broadcastMB.RetrieveLatestAndClear()
				if !item.IsValid() {
					continue
				}
				t.headBroadcaster.BroadcastNewLongestChain(item)
			}
		}
	} else {
		t.log.Info("Head sampling is disabled - callback will be called on every head")
		for {
			select {
			case <-t.chStop:
				return
			case <-t.broadcastMB.Notify():
				for {
					item, exists := t.broadcastMB.Retrieve()
					if !exists {
						break
					}
					t.headBroadcaster.BroadcastNewLongestChain(item)
				}
			}
		}
	}
}

func (t *Tracker[HMH, S, ID, BLOCK_HASH]) backfillLoop() {
	defer t.wgDone.Done()

	ctx, cancel := t.chStop.NewCtx()
	defer cancel()

	for {
		select {
		case <-t.chStop:
			return
		case <-t.backfillMB.Notify():
			for {
				head, exists := t.backfillMB.Retrieve()
				if !exists {
					break
				}
				{
					err := t.Backfill(ctx, head, uint(t.config.FinalityDepth()))
					if err != nil {
						t.log.Warnw("Unexpected error while backfilling heads", "err", err)
					} else if ctx.Err() != nil {
						break
					}
				}
			}
		}
	}
}

// backfill fetches all missing heads up until the base height
func (t *Tracker[HMH, S, ID, BLOCK_HASH]) backfill(ctx context.Context, head types.Head[BLOCK_HASH], baseHeight int64) (err error) {
	headBlockNumber := head.BlockNumber()
	if headBlockNumber <= baseHeight {
		return nil
	}
	mark := time.Now()
	fetched := 0
	l := t.log.With("blockNumber", headBlockNumber,
		"n", headBlockNumber-baseHeight,
		"fromBlockHeight", baseHeight,
		"toBlockHeight", headBlockNumber-1)
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

	for i := head.BlockNumber() - 1; i >= baseHeight; i-- {
		// NOTE: Sequential requests here mean it's a potential performance bottleneck, be aware!
		existingHead := t.headSaver.Chain(head.GetParentHash())
		if existingHead.IsValid() {
			head = existingHead
			continue
		}
		head, err = t.fetchAndSaveHead(ctx, i)
		fetched++
		if ctx.Err() != nil {
			t.log.Debugw("context canceled, aborting backfill", "err", err, "ctx.Err", ctx.Err())
			break
		} else if err != nil {
			return errors.Wrap(err, "fetchAndSaveHead failed")
		}
	}
	return
}

func (t *Tracker[HMH, S, ID, BLOCK_HASH]) fetchAndSaveHead(ctx context.Context, n int64) (HMH, error) {
	t.log.Debugw("Fetching head", "blockHeight", n)
	head, err := t.client.HeadByNumber(ctx, big.NewInt(n))
	if err != nil {
		return t.getNilHead(), err
	} else if !head.IsValid() {
		return t.getNilHead(), errors.New("got nil head")
	}
	err = t.headSaver.Save(ctx, head)
	if err != nil {
		return t.getNilHead(), err
	}
	return head, nil
}
