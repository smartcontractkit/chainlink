package headtracker

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/mailbox"

	htrktypes "github.com/smartcontractkit/chainlink/v2/common/headtracker/types"
	"github.com/smartcontractkit/chainlink/v2/common/types"
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

// HeadTracker holds and stores the block experienced by a particular node in a thread safe manner.
//
//go:generate mockery --quiet --name HeadTracker --output ./mocks/ --case=underscore
type HeadTracker[H types.Head[BLOCK_HASH], BLOCK_HASH types.Hashable] interface {
	services.Service
	// Backfill given a head will fill in any missing heads up to latestFinalized
	Backfill(ctx context.Context, headWithChain, latestFinalized H) (err error)
	LatestChain() H
}

type headTracker[
	HTH htrktypes.Head[BLOCK_HASH, ID],
	S types.Subscription,
	ID types.ID,
	BLOCK_HASH types.Hashable,
] struct {
	services.StateMachine
	log             logger.SugaredLogger
	headBroadcaster HeadBroadcaster[HTH, BLOCK_HASH]
	headSaver       HeadSaver[HTH, BLOCK_HASH]
	mailMon         *mailbox.Monitor
	client          htrktypes.Client[HTH, S, ID, BLOCK_HASH]
	chainID         types.ID
	config          htrktypes.Config
	htConfig        htrktypes.HeadTrackerConfig

	backfillMB   *mailbox.Mailbox[HTH]
	broadcastMB  *mailbox.Mailbox[HTH]
	headListener HeadListener[HTH, BLOCK_HASH]
	chStop       services.StopChan
	wgDone       sync.WaitGroup
	getNilHead   func() HTH
}

// NewHeadTracker instantiates a new HeadTracker using HeadSaver to persist new block numbers.
func NewHeadTracker[
	HTH htrktypes.Head[BLOCK_HASH, ID],
	S types.Subscription,
	ID types.ID,
	BLOCK_HASH types.Hashable,
](
	lggr logger.Logger,
	client htrktypes.Client[HTH, S, ID, BLOCK_HASH],
	config htrktypes.Config,
	htConfig htrktypes.HeadTrackerConfig,
	headBroadcaster HeadBroadcaster[HTH, BLOCK_HASH],
	headSaver HeadSaver[HTH, BLOCK_HASH],
	mailMon *mailbox.Monitor,
	getNilHead func() HTH,
) HeadTracker[HTH, BLOCK_HASH] {
	chStop := make(chan struct{})
	lggr = logger.Named(lggr, "HeadTracker")
	return &headTracker[HTH, S, ID, BLOCK_HASH]{
		headBroadcaster: headBroadcaster,
		client:          client,
		chainID:         client.ConfiguredChainID(),
		config:          config,
		htConfig:        htConfig,
		log:             logger.Sugared(lggr),
		backfillMB:      mailbox.NewSingle[HTH](),
		broadcastMB:     mailbox.New[HTH](HeadsBufferSize),
		chStop:          chStop,
		headListener:    NewHeadListener[HTH, S, ID, BLOCK_HASH](lggr, client, config, chStop),
		headSaver:       headSaver,
		mailMon:         mailMon,
		getNilHead:      getNilHead,
	}
}

// Start starts HeadTracker service.
func (ht *headTracker[HTH, S, ID, BLOCK_HASH]) Start(ctx context.Context) error {
	return ht.StartOnce("HeadTracker", func() error {
		ht.log.Debugw("Starting HeadTracker", "chainID", ht.chainID)
		// NOTE: Always try to start the head tracker off with whatever the
		// latest head is, without waiting for the subscription to send us one.
		//
		// In some cases the subscription will send us the most recent head
		// anyway when we connect (but we should not rely on this because it is
		// not specced). If it happens this is fine, and the head will be
		// ignored as a duplicate.
		err := ht.handleInitialHead(ctx)
		if err != nil {
			if ctx.Err() != nil {
				return ctx.Err()
			}
			ht.log.Errorw("Error handling initial head", "err", err)
		}

		ht.wgDone.Add(3)
		go ht.headListener.ListenForNewHeads(ht.handleNewHead, ht.wgDone.Done)
		go ht.backfillLoop()
		go ht.broadcastLoop()

		ht.mailMon.Monitor(ht.broadcastMB, "HeadTracker", "Broadcast", ht.chainID.String())

		return nil
	})
}

func (ht *headTracker[HTH, S, ID, BLOCK_HASH]) handleInitialHead(ctx context.Context) error {
	initialHead, err := ht.client.HeadByNumber(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to fetch initial head: %w", err)
	}

	if !initialHead.IsValid() {
		ht.log.Warnw("Got nil initial head", "head", initialHead)
		return nil
	}
	ht.log.Debugw("Got initial head", "head", initialHead, "blockNumber", initialHead.BlockNumber(), "blockHash", initialHead.BlockHash())

	latestFinalized, err := ht.calculateLatestFinalized(ctx, initialHead)
	if err != nil {
		return fmt.Errorf("failed to calculate latest finalized head: %w", err)
	}

	if !latestFinalized.IsValid() {
		return fmt.Errorf("latest finalized block is not valid")
	}

	latestChain, err := ht.headSaver.Load(ctx, latestFinalized.BlockNumber())
	if err != nil {
		return fmt.Errorf("failed to initialized headSaver: %w", err)
	}

	if latestChain.IsValid() {
		earliest := latestChain.EarliestHeadInChain()
		ht.log.Debugw(
			"Loaded chain from DB",
			"latest_blockNumber", latestChain.BlockNumber(),
			"latest_blockHash", latestChain.BlockHash(),
			"earliest_blockNumber", earliest.BlockNumber(),
			"earliest_blockHash", earliest.BlockHash(),
		)
	}
	if err := ht.handleNewHead(ctx, initialHead); err != nil {
		return fmt.Errorf("error handling initial head: %w", err)
	}

	return nil
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
	report := map[string]error{ht.Name(): ht.Healthy()}
	services.CopyHealth(report, ht.headListener.HealthReport())
	return report
}

func (ht *headTracker[HTH, S, ID, BLOCK_HASH]) Backfill(ctx context.Context, headWithChain, latestFinalized HTH) (err error) {
	if !latestFinalized.IsValid() {
		return errors.New("can not perform backfill without a valid latestFinalized head")
	}

	if headWithChain.BlockNumber() < latestFinalized.BlockNumber() {
		const errMsg = "invariant violation: expected head of canonical chain to be ahead of the latestFinalized"
		ht.log.With("head_block_num", headWithChain.BlockNumber(),
			"latest_finalized_block_number", latestFinalized.BlockNumber()).
			Criticalf(errMsg)
		return errors.New(errMsg)
	}

	return ht.backfill(ctx, headWithChain, latestFinalized)
}

func (ht *headTracker[HTH, S, ID, BLOCK_HASH]) LatestChain() HTH {
	return ht.headSaver.LatestChain()
}

func (ht *headTracker[HTH, S, ID, BLOCK_HASH]) handleNewHead(ctx context.Context, head HTH) error {
	prevHead := ht.headSaver.LatestChain()

	ht.log.Debugw(fmt.Sprintf("Received new head %v", head.BlockNumber()),
		"blockHash", head.BlockHash(),
		"parentHeadHash", head.GetParentHash(),
		"blockTs", head.GetTimestamp(),
		"blockTsUnix", head.GetTimestamp().Unix(),
		"blockDifficulty", head.BlockDifficulty(),
	)

	err := ht.headSaver.Save(ctx, head)
	if ctx.Err() != nil {
		return nil
	} else if err != nil {
		return fmt.Errorf("failed to save head: %#v: %w", head, err)
	}

	if !prevHead.IsValid() || head.BlockNumber() > prevHead.BlockNumber() {
		promCurrentHead.WithLabelValues(ht.chainID.String()).Set(float64(head.BlockNumber()))

		headWithChain := ht.headSaver.Chain(head.BlockHash())
		if !headWithChain.IsValid() {
			return fmt.Errorf("HeadTracker#handleNewHighestHead headWithChain was unexpectedly nil")
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
		prevUnFinalizedHead := prevHead.BlockNumber() - int64(ht.config.FinalityDepth())
		if head.BlockNumber() < prevUnFinalizedHead {
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
					latestFinalized, err := ht.calculateLatestFinalized(ctx, head)
					if err != nil {
						ht.log.Warnw("Failed to calculate finalized block", "err", err)
						continue
					}

					err = ht.Backfill(ctx, head, latestFinalized)
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

// calculateLatestFinalized - returns latest finalized block. It's expected that currentHeadNumber - is the head of
// canonical chain. There is no guaranties that returned block belongs to the canonical chain. Additional verification
// must be performed before usage.
func (ht *headTracker[HTH, S, ID, BLOCK_HASH]) calculateLatestFinalized(ctx context.Context, currentHead HTH) (h HTH, err error) {
	if ht.config.FinalityTagEnabled() {
		return ht.client.LatestFinalizedBlock(ctx)
	}
	// no need to make an additional RPC call on chains with instant finality
	if ht.config.FinalityDepth() == 0 {
		return currentHead, nil
	}
	finalizedBlockNumber := currentHead.BlockNumber() - int64(ht.config.FinalityDepth())
	if finalizedBlockNumber <= 0 {
		finalizedBlockNumber = 0
	}
	return ht.client.HeadByNumber(ctx, big.NewInt(finalizedBlockNumber))
}

// backfill fetches all missing heads up until the latestFinalizedHead
func (ht *headTracker[HTH, S, ID, BLOCK_HASH]) backfill(ctx context.Context, head, latestFinalizedHead HTH) (err error) {
	headBlockNumber := head.BlockNumber()
	mark := time.Now()
	fetched := 0
	baseHeight := latestFinalizedHead.BlockNumber()
	l := ht.log.With("blockNumber", headBlockNumber,
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
		existingHead := ht.headSaver.Chain(head.GetParentHash())
		if existingHead.IsValid() {
			head = existingHead
			continue
		}
		head, err = ht.fetchAndSaveHead(ctx, i, head.GetParentHash())
		fetched++
		if ctx.Err() != nil {
			ht.log.Debugw("context canceled, aborting backfill", "err", err, "ctx.Err", ctx.Err())
			return fmt.Errorf("fetchAndSaveHead failed: %w", ctx.Err())
		} else if err != nil {
			return fmt.Errorf("fetchAndSaveHead failed: %w", err)
		}
	}

	if head.BlockHash() != latestFinalizedHead.BlockHash() {
		const errMsg = "expected finalized block to be present in canonical chain"
		ht.log.With("finalized_block_number", latestFinalizedHead.BlockNumber(), "finalized_hash", latestFinalizedHead.BlockHash(),
			"canonical_chain_block_number", head.BlockNumber(), "canonical_chain_hash", head.BlockHash()).Criticalf(errMsg)
		return fmt.Errorf(errMsg)
	}

	l = l.With("latest_finalized_block_hash", latestFinalizedHead.BlockHash(),
		"latest_finalized_block_number", latestFinalizedHead.BlockNumber())

	err = ht.headSaver.MarkFinalized(ctx, latestFinalizedHead)
	if err != nil {
		l.Debugw("failed to mark block as finalized", "err", err)
		return nil
	}

	l.Debugw("marked block as finalized")

	return
}

func (ht *headTracker[HTH, S, ID, BLOCK_HASH]) fetchAndSaveHead(ctx context.Context, n int64, hash BLOCK_HASH) (HTH, error) {
	ht.log.Debugw("Fetching head", "blockHeight", n, "blockHash", hash)
	head, err := ht.client.HeadByHash(ctx, hash)
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
