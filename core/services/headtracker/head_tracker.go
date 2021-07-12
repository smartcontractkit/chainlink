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
	"github.com/smartcontractkit/chainlink/core/logger"
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
)

// HeadTracker holds and stores the latest block number experienced by this particular node
// in a thread safe manner. Reconstitutes the last block number from the data
// store on reboot.
type HeadTracker struct {
	log             *logger.Logger
	headBroadcaster httypes.HeadBroadcaster
	store           *strpkg.Store

	backfillMB   utils.Mailbox
	samplingMB   utils.Mailbox
	muLogger     sync.RWMutex
	headListener *HeadListener
	headSaver    *HeadSaver
	chStop       chan struct{}
	wgDone       *sync.WaitGroup
	utils.StartStopOnce
}

// NewHeadTracker instantiates a new HeadTracker using the orm to persist new block numbers.
// Can be passed in an optional sleeper object that will dictate how often
// it tries to reconnect.
func NewHeadTracker(
	l *logger.Logger,
	store *strpkg.Store,
	headBroadcaster httypes.HeadBroadcaster,
	sleepers ...utils.Sleeper,
) *HeadTracker {

	var wgDone sync.WaitGroup
	chStop := make(chan struct{})

	return &HeadTracker{
		store:           store,
		headBroadcaster: headBroadcaster,
		log:             l,
		backfillMB:      *utils.NewMailbox(1),
		samplingMB:      *utils.NewMailbox(1),
		chStop:          chStop,
		wgDone:          &wgDone,
		headListener:    NewHeadListener(l, store.EthClient, store.Config, chStop, &wgDone, sleepers...),
		headSaver:       NewHeadSaver(store),
	}
}

// SetLogger sets and reconfigures the log for the head tracker service
func (ht *HeadTracker) SetLogger(logger *logger.Logger) {
	ht.muLogger.Lock()
	defer ht.muLogger.Unlock()
	ht.log = logger
	ht.headListener.SetLogger(logger)
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
		ht.logger().Debug("Starting HeadTracker")
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

// Connected returns whether or not this HeadTracker is connected.
func (ht *HeadTracker) Connected() bool {
	return ht.headListener.Connected()
}

func (ht *HeadTracker) handleConnected() {
	ht.connect(ht.headSaver.HighestSeenHead())
}

func (ht *HeadTracker) connect(bn *models.Head) {
	if err := ht.headBroadcaster.Connect(bn); err != nil {
		ht.logger().Warn(err)
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

			ht.headBroadcaster.OnNewLongestChain(ctx, head)
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
	head, err := ht.store.EthClient.HeadByNumber(ctx, big.NewInt(n))
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
		if head.Number < prevHead.Number-int64(ht.store.Config.EthFinalityDepth()) {
			ht.logger().Errorf("HeadTracker: got very old block with number %d (highest seen was %d). This is a problem and either means a very deep re-org occurred, or the chain went backwards in block numbers. This node will not function correctly without manual intervention.", head.Number, prevHead.Number)
		}
	}
	return nil
}
