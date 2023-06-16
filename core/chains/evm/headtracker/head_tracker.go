package headtracker

import (
	"context"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"go.uber.org/zap/zapcore"
	"golang.org/x/exp/maps"

	htrktypes "github.com/smartcontractkit/chainlink/v2/common/headtracker/types"
	txmgrtypes "github.com/smartcontractkit/chainlink/v2/common/txmgr/types"
	commontypes "github.com/smartcontractkit/chainlink/v2/common/types"
	evmclient "github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	httypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/headtracker/types"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
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

type headTracker[
	HTH htrktypes.Head[BLOCK_HASH, ID],
	S commontypes.Subscription,
	ID txmgrtypes.ID,
	BLOCK_HASH commontypes.Hashable,
] struct {
	log             logger.Logger
	headBroadcaster commontypes.HeadBroadcaster[HTH, BLOCK_HASH]
	headSaver       commontypes.HeadSaver[HTH, BLOCK_HASH]
	mailMon         *utils.MailboxMonitor
	ethClient       htrktypes.Client[HTH, S, ID, BLOCK_HASH]
	chainID         ID
	config          htrktypes.Config
	htConfig        htrktypes.HeadTrackerConfig

	backfillMB   *utils.Mailbox[HTH]
	broadcastMB  *utils.Mailbox[HTH]
	headListener commontypes.HeadListener[HTH, BLOCK_HASH]
	chStop       utils.StopChan
	wgDone       sync.WaitGroup
	utils.StartStopOnce
	getNilHead func() HTH
}

type evmHeadTracker = headTracker[*evmtypes.Head, ethereum.Subscription, *big.Int, common.Hash]

var _ commontypes.HeadTracker[*evmtypes.Head, common.Hash] = (*evmHeadTracker)(nil)

// NewHeadTracker instantiates a new HeadTracker using HeadSaver to persist new block numbers.
func NewHeadTracker[
	HTH htrktypes.Head[BLOCK_HASH, ID],
	S commontypes.Subscription,
	ID txmgrtypes.ID,
	BLOCK_HASH commontypes.Hashable,
](
	lggr logger.Logger,
	client htrktypes.Client[HTH, S, ID, BLOCK_HASH],
	config htrktypes.Config,
	htConfig htrktypes.HeadTrackerConfig,
	headBroadcaster commontypes.HeadBroadcaster[HTH, BLOCK_HASH],
	headSaver commontypes.HeadSaver[HTH, BLOCK_HASH],
	mailMon *utils.MailboxMonitor,
	getNilHead func() HTH,
) commontypes.HeadTracker[HTH, BLOCK_HASH] {
	chStop := make(chan struct{})
	lggr = lggr.Named("HeadTracker")
	return &headTracker[HTH, S, ID, BLOCK_HASH]{
		headBroadcaster: headBroadcaster,
		ethClient:       client,
		chainID:         client.ConfiguredChainID(),
		config:          config,
		htConfig:        htConfig,
		log:             lggr,
		backfillMB:      utils.NewSingleMailbox[HTH](),
		broadcastMB:     utils.NewMailbox[HTH](HeadsBufferSize),
		chStop:          chStop,
		headListener:    NewHeadListener[HTH, S, ID, BLOCK_HASH](lggr, client, config, chStop),
		headSaver:       headSaver,
		mailMon:         mailMon,
		getNilHead:      getNilHead,
	}
}

func NewEVMHeadTracker(
	lggr logger.Logger,
	ethClient evmclient.Client,
	config Config,
	htConfig HeadTrackerConfig,
	headBroadcaster httypes.HeadBroadcaster,
	headSaver httypes.HeadSaver,
	mailMon *utils.MailboxMonitor,
) httypes.HeadTracker {
	chStop := make(chan struct{})
	lggr = lggr.Named("HeadTracker")
	return &evmHeadTracker{
		headBroadcaster: headBroadcaster,
		ethClient:       ethClient,
		chainID:         ethClient.ConfiguredChainID(),
		config:          NewWrappedConfig(config),
		htConfig:        htConfig,
		log:             lggr,
		backfillMB:      utils.NewSingleMailbox[*evmtypes.Head](),
		broadcastMB:     utils.NewMailbox[*evmtypes.Head](HeadsBufferSize),
		chStop:          chStop,
		headListener:    NewEVMHeadListener(lggr, ethClient, config, chStop),
		headSaver:       headSaver,
		mailMon:         mailMon,
		getNilHead:      func() *evmtypes.Head { return nil },
	}
}

// Start starts HeadTracker service.
func (ht *headTracker[HTH, S, ID, BLOCK_HASH]) Start(ctx context.Context) error {
	return ht.StartOnce("HeadTracker", func() error {
		ht.log.Debugw("Starting HeadTracker", "chainID", ht.chainID)
		latestChain, err := ht.headSaver.Load(ctx)
		if err != nil {
			return err
		}
		if latestChain.IsValid() {
			ht.log.Debugw(
				fmt.Sprintf("HeadTracker: Tracking logs from last block %v with hash %s", config.FriendlyNumber(latestChain.BlockNumber()), latestChain.BlockHash()),
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
		initialHead, err := ht.getInitialHead(ctx)
		if err != nil {
			if errors.Is(err, ctx.Err()) {
				return nil
			}
			ht.log.Errorw("Error getting initial head", "err", err)
		} else if initialHead.IsValid() {
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
func (ht *headTracker[HTH, S, ID, BLOCK_HASH]) Close() error {
	return ht.StopOnce("HeadTracker", func() error {
		close(ht.chStop)
		ht.wgDone.Wait()
		return ht.broadcastMB.Close()
	})
}

func (ht *headTracker[HTH, S, ID, BLOCK_HASH]) Name() string {
	return ht.log.Name()
}

func (ht *headTracker[HTH, S, ID, BLOCK_HASH]) HealthReport() map[string]error {
	report := map[string]error{
		ht.Name(): ht.StartStopOnce.Healthy(),
	}
	maps.Copy(report, ht.headListener.HealthReport())
	return report
}

func (ht *headTracker[HTH, S, ID, BLOCK_HASH]) Backfill(ctx context.Context, headWithChain HTH, depth uint) (err error) {
	if uint(headWithChain.ChainLength()) >= depth {
		return nil
	}

	baseHeight := headWithChain.BlockNumber() - int64(depth-1)
	if baseHeight < 0 {
		baseHeight = 0
	}

	return ht.backfill(ctx, headWithChain.EarliestHeadInChain(), baseHeight)
}

func (ht *headTracker[HTH, S, ID, BLOCK_HASH]) LatestChain() HTH {
	return ht.headSaver.LatestChain()
}

func (ht *headTracker[HTH, S, ID, BLOCK_HASH]) getInitialHead(ctx context.Context) (HTH, error) {
	head, err := ht.ethClient.HeadByNumber(ctx, nil)
	if err != nil {
		return ht.getNilHead(), errors.Wrap(err, "failed to fetch initial head")
	}
	loggerFields := []interface{}{"head", head}
	if head.IsValid() {
		loggerFields = append(loggerFields, "blockNumber", head.BlockNumber(), "blockHash", head.BlockHash())
	}
	ht.log.Debugw("Got initial head", loggerFields...)
	return head, nil
}

func (ht *headTracker[HTH, S, ID, BLOCK_HASH]) handleNewHead(ctx context.Context, head HTH) error {
	prevHead := ht.headSaver.LatestChain()

	ht.log.Debugw(fmt.Sprintf("Received new head %v", config.FriendlyNumber(head.BlockNumber())),
		"blockHeight", head,
		"blockHash", head.BlockHash(),
		"parentHeadHash", head.GetParentHash(),
	)

	err := ht.headSaver.Save(ctx, head)
	if ctx.Err() != nil {
		return nil
	} else if err != nil {
		return errors.Wrapf(err, "failed to save head: %#v", head)
	}

	if !prevHead.IsValid() || head.BlockNumber() > prevHead.BlockNumber() {
		promCurrentHead.WithLabelValues(ht.chainID.String()).Set(float64(head.BlockNumber()))

		headWithChain := ht.headSaver.Chain(head.BlockHash())
		if !headWithChain.IsValid() {
			return errors.Errorf("HeadTracker#handleNewHighestHead headWithChain was unexpectedly nil")
		}
		ht.backfillMB.Deliver(headWithChain)
		ht.broadcastMB.Deliver(headWithChain)
	} else if head.BlockNumber() == prevHead.BlockNumber() {
		if head.BlockHash() != prevHead.BlockHash() {
			ht.log.Debugw("Got duplicate head", "blockNum", head.BlockNumber(), "head", head.BlockHash(), "prevHead", prevHead.BlockHash())
		} else {
			ht.log.Debugw("Head already in the database", "head", head.BlockHash())
		}
	} else {
		ht.log.Debugw("Got out of order head", "blockNum", head.BlockNumber(), "head", head.BlockHash(), "prevHead", prevHead.BlockNumber())
		if head.BlockNumber() < prevHead.BlockNumber()-int64(ht.config.FinalityDepth()) {
			promOldHead.WithLabelValues(ht.chainID.String()).Inc()
			ht.log.Criticalf("Got very old block with number %d (highest seen was %d). This is a problem and either means a very deep re-org occurred, one of the RPC nodes has gotten far out of sync, or the chain went backwards in block numbers. This node may not function correctly without manual intervention.", head.BlockNumber(), prevHead.BlockNumber())
			ht.SvcErrBuffer.Append(errors.New("got very old block"))
		}
	}
	return nil
}

func (ht *headTracker[HTH, S, ID, BLOCK_HASH]) broadcastLoop() {
	defer ht.wgDone.Done()

	samplingInterval := ht.htConfig.SamplingInterval()
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
				if !item.IsValid() {
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

func (ht *headTracker[HTH, S, ID, BLOCK_HASH]) backfillLoop() {
	defer ht.wgDone.Done()

	ctx, cancel := ht.chStop.NewCtx()
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
					err := ht.Backfill(ctx, head, uint(ht.config.FinalityDepth()))
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
func (ht *headTracker[HTH, S, ID, BLOCK_HASH]) backfill(ctx context.Context, head commontypes.Head[BLOCK_HASH], baseHeight int64) (err error) {
	if head.BlockNumber() <= baseHeight {
		return nil
	}
	mark := time.Now()
	fetched := 0
	l := ht.log.With("blockNumber", head.BlockNumber(),
		"n", head.BlockNumber()-baseHeight,
		"fromBlockHeight", baseHeight,
		"toBlockHeight", head.BlockNumber()-1)
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
		existingHead := ht.headSaver.Chain(head.GetParentHash())
		if existingHead.IsValid() {
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

func (ht *headTracker[HTH, S, ID, BLOCK_HASH]) fetchAndSaveHead(ctx context.Context, n int64) (HTH, error) {
	ht.log.Debugw("Fetching head", "blockHeight", n)
	head, err := ht.ethClient.HeadByNumber(ctx, big.NewInt(n))
	if err != nil {
		return ht.getNilHead(), err
	} else if !head.IsValid() {
		return ht.getNilHead(), errors.New("got nil head")
	}
	err = ht.headSaver.Save(ctx, head)
	if err != nil {
		return ht.getNilHead(), err
	}
	return head, nil
}

var NullTracker httypes.HeadTracker = &nullTracker{}

type nullTracker struct{}

func (*nullTracker) Start(context.Context) error    { return nil }
func (*nullTracker) Close() error                   { return nil }
func (*nullTracker) Ready() error                   { return nil }
func (*nullTracker) HealthReport() map[string]error { return map[string]error{} }
func (*nullTracker) Name() string                   { return "" }
func (*nullTracker) SetLogLevel(zapcore.Level)      {}
func (*nullTracker) Backfill(ctx context.Context, headWithChain *evmtypes.Head, depth uint) (err error) {
	return nil
}
func (*nullTracker) LatestChain() *evmtypes.Head { return nil }
