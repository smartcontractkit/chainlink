package headtracker

import (
	"context"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/store/presenters"
	"github.com/smartcontractkit/chainlink/core/utils"
)

var (
	promCallbackDuration = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "fetcher_backend_callback_execution_duration",
		Help: "How long it took to execute all callbacks (ms)",
	})
	promCallbackDurationHist = promauto.NewHistogram(prometheus.HistogramOpts{
		Name:    "fetcher_backend_callback_execution_duration_hist",
		Help:    "How long it took to execute all callbacks (ms) histogram",
		Buckets: []float64{50, 100, 250, 500, 1000, 2000, 5000, 10000, 15000, 30000, 100000},
	})
	promCurrentHead = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "fetcher_backend_current_head",
		Help: "The highest seen head number",
	})
)

// HeadSubscription - Managing fetching of data for a single head
type HeadSubscription interface {
	Head() models.Head

	// Block - will deliver the block data
	Block() <-chan types.Block

	// Receipts - will deliver receipt data
	Receipts() <-chan []types.Receipt

	// Unsubscribe - stops inflow of data for this head
	Unsubscribe()
}

// Backend - Managing fetching of heads and block data
type Backend interface {

	// Subscribe - Start receiving newest heads and data for them
	Subscribe() <-chan HeadSubscription

	HeadByNumber(n *big.Int) HeadSubscription
}

type BackendConfig interface {
	HeadTimeBudget() time.Duration
	EthFinalityDepth() uint
}

type DefaultBackend struct {
	logger *logger.Logger

	config      BackendConfig
	listener    *HeadListener
	headTracker *services.HeadTracker

	startStopOnce utils.StartStopOnce
	chStop        chan struct{}
	wgDone        sync.WaitGroup
}

func newDefaultBackend(listener *HeadListener, headTracker *services.HeadTracker) *DefaultBackend {
	return &DefaultBackend{
		listener:    listener,
		headTracker: headTracker,
	}
}

func (hl *DefaultBackend) Start() error {
	return hl.StartOnce("DefaultBackend", func() error {
		hl.wgDone.Add(1)
		go hl.eventLoop()
		return nil
	})
}

func (hl *DefaultBackend) Stop() error {
	return hl.StopOnce("DefaultBackend", func() error {
		hl.logger.Info("DefaultBackend stopping")
		close(hl.chStop)
		hl.wgDone.Wait()
		return nil
	})
}

func (ht *DefaultBackend) HandleNewHead(ctx context.Context, head models.Head) error {

	defer func(start time.Time, number int64) {
		elapsed := time.Since(start)
		ms := float64(elapsed.Milliseconds())
		promCallbackDuration.Set(ms)
		promCallbackDurationHist.Observe(ms)

		callbackExecutionThreshold := ht.config.HeadTimeBudget() / 2
		if elapsed > callbackExecutionThreshold {
			ht.logger.Warnw(fmt.Sprintf("HeadTracker finished processing head %v in %s which exceeds callback execution threshold of %s", number, elapsed.String(), ht.store.Config.HeadTimeBudget().String()), "blockNumber", number, "time", elapsed, "id", "head_tracker")
		} else {
			ht.logger.Debugw(fmt.Sprintf("HeadTracker finished processing head %v in %s", number, elapsed.String()), "blockNumber", number, "time", elapsed, "id", "head_tracker")
		}
	}(time.Now(), int64(head.Number))
	prevHead := ht.headTracker.HighestSeenHead()

	ht.logger.Debugw(fmt.Sprintf("HeadTracker: Received new head %v", presenters.FriendlyBigInt(head.ToInt())),
		"blockHeight", head.ToInt(),
		"blockHash", head.Hash,
		"parentHeadHash", head.ParentHash,
	)

	err := ht.headTracker.Save(ctx, head)
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
			ht.logger.Debugw("HeadTracker: got duplicate head", "blockNum", head.Number, "gotHead", head.Hash.Hex(), "highestSeenHead", prevHead.Hash.Hex())
		} else {
			ht.logger.Debugw("HeadTracker: head already in the database", "gotHead", head.Hash.Hex())
		}
	} else {
		ht.logger.Debugw("HeadTracker: got out of order head", "blockNum", head.Number, "gotHead", head.Hash.Hex(), "highestSeenHead", prevHead.Number)
	}
	return nil

}

func (ht *DefaultBackend) handleNewHighestHead(ctx context.Context, head models.Head) error {
	promCurrentHead.Set(float64(head.Number))

	headWithChain, err := ht.headTracker.Chain(ctx, head.Hash, ht.config.EthFinalityDepth())
	if ctx.Err() != nil {
		return nil
	} else if err != nil {
		return errors.Wrap(err, "HeadTracker#handleNewHighestHead failed fetching chain")
	}

	//TODO:
	//ht.backfillMB.Deliver(headWithChain)

	_, err = ht.blockFetcher.fetchBlock(ctx, headWithChain, promBlockSizenHist, promBlockFetchDurationHist, promBlockBatchFetchDurationHist,
		promReceiptsFetchDurationHist, promReceiptsLimitedFetchDurationHist, promReceiptCount)

	if ctx.Err() != nil {
		return nil
	} else if err != nil {
		return errors.Wrap(err, "HeadTracker#handleNewHighestHead failed fetching whole block")
	}

	ht.onNewLongestChain(ctx, headWithChain)
	return nil
}

func (hl *DefaultBackend) eventLoop() {
	defer hl.wgDone.Done()

	//defer func() {
	//	if err := hl.unsubscribeFromHead(); err != nil {
	//		hl.logger.Warn(errors.Wrap(err, "HeadListener failed when unsubscribe from head"))
	//	}
	//}()
	//
	//ctx, cancel := utils.ContextFromChan(hl.chStop)
	//defer cancel()
	//
	//for {
	//	if !hl.subscribe() {
	//		break
	//	}
	//	err := hl.receiveHeaders(ctx)
	//	if ctx.Err() != nil {
	//		break
	//	} else if err != nil {
	//		hl.logger.Errorw(fmt.Sprintf("Error in new head subscription, unsubscribed: %s", err.Error()), "err", err)
	//		continue
	//	} else {
	//		break
	//	}
	//}
}
