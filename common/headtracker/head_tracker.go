package headtracker

import (
	"context"
	"errors"
	"fmt"
	"math/big"
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
type HeadTracker[H types.Head[BLOCK_HASH], BLOCK_HASH types.Hashable] interface {
	services.Service
	// Backfill given a head will fill in any missing heads up to latestFinalized
	Backfill(ctx context.Context, headWithChain H) (err error)
	LatestChain() H
	// LatestAndFinalizedBlock - returns latest and latest finalized blocks.
	// NOTE: Returns latest finalized block as is, ignoring the FinalityTagBypass feature flag.
	LatestAndFinalizedBlock(ctx context.Context) (latest, finalized H, err error)
}

type headTracker[
	HTH htrktypes.Head[BLOCK_HASH, ID],
	S types.Subscription,
	ID types.ID,
	BLOCK_HASH types.Hashable,
] struct {
	services.Service
	eng *services.Engine

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
	ht := &headTracker[HTH, S, ID, BLOCK_HASH]{
		headBroadcaster: headBroadcaster,
		client:          client,
		chainID:         client.ConfiguredChainID(),
		config:          config,
		htConfig:        htConfig,
		backfillMB:      mailbox.NewSingle[HTH](),
		broadcastMB:     mailbox.New[HTH](HeadsBufferSize),
		headSaver:       headSaver,
		mailMon:         mailMon,
		getNilHead:      getNilHead,
	}
	ht.Service, ht.eng = services.Config{
		Name: "HeadTracker",
		NewSubServices: func(lggr logger.Logger) []services.Service {
			ht.headListener = NewHeadListener[HTH, S, ID, BLOCK_HASH](lggr, client, config,
				// NOTE: Always try to start the head tracker off with whatever the
				// latest head is, without waiting for the subscription to send us one.
				//
				// In some cases the subscription will send us the most recent head
				// anyway when we connect (but we should not rely on this because it is
				// not specced). If it happens this is fine, and the head will be
				// ignored as a duplicate.
				func(ctx context.Context) {
					err := ht.handleInitialHead(ctx)
					if err != nil {
						ht.log.Errorw("Error handling initial head", "err", err.Error())
					}
				}, ht.handleNewHead)
			return []services.Service{ht.headListener}
		},
		Start: ht.start,
		Close: ht.close,
	}.NewServiceEngine(lggr)
	ht.log = logger.Sugared(ht.eng)
	return ht
}

// Start starts HeadTracker service.
func (ht *headTracker[HTH, S, ID, BLOCK_HASH]) start(context.Context) error {
	ht.eng.Go(ht.backfillLoop)
	ht.eng.Go(ht.broadcastLoop)

	ht.mailMon.Monitor(ht.broadcastMB, "HeadTracker", "Broadcast", ht.chainID.String())

	return nil
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

	latestFinalized, err := ht.calculateLatestFinalized(ctx, initialHead, ht.htConfig.FinalityTagBypass())
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

func (ht *headTracker[HTH, S, ID, BLOCK_HASH]) close() error {
	return ht.broadcastMB.Close()
}

func (ht *headTracker[HTH, S, ID, BLOCK_HASH]) Backfill(ctx context.Context, headWithChain HTH) (err error) {
	latestFinalized, err := ht.calculateLatestFinalized(ctx, headWithChain, ht.htConfig.FinalityTagBypass())
	if err != nil {
		return fmt.Errorf("failed to calculate finalized block: %w", err)
	}

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

	if headWithChain.BlockNumber()-latestFinalized.BlockNumber() > int64(ht.htConfig.MaxAllowedFinalityDepth()) {
		return fmt.Errorf("gap between latest finalized block (%d) and current head (%d) is too large (> %d)",
			latestFinalized.BlockNumber(), headWithChain.BlockNumber(), ht.htConfig.MaxAllowedFinalityDepth())
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

	if err := ht.headSaver.Save(ctx, head); ctx.Err() != nil {
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
		prevLatestFinalized := prevHead.LatestFinalizedHead()

		if prevLatestFinalized != nil && head.BlockNumber() <= prevLatestFinalized.BlockNumber() {
			promOldHead.WithLabelValues(ht.chainID.String()).Inc()
			err := fmt.Errorf("got very old block with number %d (highest seen was %d)", head.BlockNumber(), prevHead.BlockNumber())
			ht.log.Critical("Got very old block. Either a very deep re-org occurred, one of the RPC nodes has gotten far out of sync, or the chain went backwards in block numbers. This node may not function correctly without manual intervention.", "err", err)
			ht.eng.EmitHealthErr(err)
		}
	}
	return nil
}

func (ht *headTracker[HTH, S, ID, BLOCK_HASH]) broadcastLoop(ctx context.Context) {
	samplingInterval := ht.htConfig.SamplingInterval()
	if samplingInterval > 0 {
		ht.log.Debugf("Head sampling is enabled - sampling interval is set to: %v", samplingInterval)
		debounceHead := time.NewTicker(samplingInterval)
		defer debounceHead.Stop()
		for {
			select {
			case <-ctx.Done():
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
			case <-ctx.Done():
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

func (ht *headTracker[HTH, S, ID, BLOCK_HASH]) backfillLoop(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-ht.backfillMB.Notify():
			for {
				head, exists := ht.backfillMB.Retrieve()
				if !exists {
					break
				}
				{
					err := ht.Backfill(ctx, head)
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

// LatestAndFinalizedBlock - returns latest and latest finalized blocks.
// NOTE: Returns latest finalized block as is, ignoring the FinalityTagBypass feature flag.
// TODO: BCI-3321 use cached values instead of making RPC requests
func (ht *headTracker[HTH, S, ID, BLOCK_HASH]) LatestAndFinalizedBlock(ctx context.Context) (latest, finalized HTH, err error) {
	latest, err = ht.client.HeadByNumber(ctx, nil)
	if err != nil {
		err = fmt.Errorf("failed to get latest block: %w", err)
		return
	}

	if !latest.IsValid() {
		err = fmt.Errorf("expected latest block to be valid")
		return
	}

	finalized, err = ht.calculateLatestFinalized(ctx, latest, false)
	if err != nil {
		err = fmt.Errorf("failed to calculate latest finalized block: %w", err)
		return
	}
	if !finalized.IsValid() {
		err = fmt.Errorf("expected finalized block to be valid")
		return
	}

	return
}

func (ht *headTracker[HTH, S, ID, BLOCK_HASH]) getHeadAtHeight(ctx context.Context, chainHeadHash BLOCK_HASH, blockHeight int64) (HTH, error) {
	chainHead := ht.headSaver.Chain(chainHeadHash)
	if chainHead.IsValid() {
		// check if provided chain contains a block of specified height
		headAtHeight, err := chainHead.HeadAtHeight(blockHeight)
		if err == nil {
			// we are forced to reload the block due to type mismatched caused by generics
			hthAtHeight := ht.headSaver.Chain(headAtHeight.BlockHash())
			// ensure that the block was not removed from the chain by another goroutine
			if hthAtHeight.IsValid() {
				return hthAtHeight, nil
			}
		}
	}

	return ht.client.HeadByNumber(ctx, big.NewInt(blockHeight))
}

// calculateLatestFinalized - returns latest finalized block. It's expected that currentHeadNumber - is the head of
// canonical chain. There is no guaranties that returned block belongs to the canonical chain. Additional verification
// must be performed before usage.
func (ht *headTracker[HTH, S, ID, BLOCK_HASH]) calculateLatestFinalized(ctx context.Context, currentHead HTH, finalityTagBypass bool) (HTH, error) {
	if ht.config.FinalityTagEnabled() && !finalityTagBypass {
		latestFinalized, err := ht.client.LatestFinalizedBlock(ctx)
		if err != nil {
			return latestFinalized, fmt.Errorf("failed to get latest finalized block: %w", err)
		}

		if !latestFinalized.IsValid() {
			return latestFinalized, fmt.Errorf("failed to get valid latest finalized block")
		}

		if ht.config.FinalizedBlockOffset() == 0 {
			return latestFinalized, nil
		}

		finalizedBlockNumber := max(latestFinalized.BlockNumber()-int64(ht.config.FinalizedBlockOffset()), 0)
		return ht.getHeadAtHeight(ctx, latestFinalized.BlockHash(), finalizedBlockNumber)
	}
	// no need to make an additional RPC call on chains with instant finality
	if ht.config.FinalityDepth() == 0 && ht.config.FinalizedBlockOffset() == 0 {
		return currentHead, nil
	}
	finalizedBlockNumber := currentHead.BlockNumber() - int64(ht.config.FinalityDepth()) - int64(ht.config.FinalizedBlockOffset())
	if finalizedBlockNumber <= 0 {
		finalizedBlockNumber = 0
	}
	return ht.getHeadAtHeight(ctx, currentHead.BlockHash(), finalizedBlockNumber)
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
