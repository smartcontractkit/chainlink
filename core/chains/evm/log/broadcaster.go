package log

import (
	"context"
	"math/big"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
	"go.uber.org/atomic"

	evmclient "github.com/smartcontractkit/chainlink/core/chains/evm/client"
	httypes "github.com/smartcontractkit/chainlink/core/chains/evm/headtracker/types"
	evmtypes "github.com/smartcontractkit/chainlink/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/null"
	"github.com/smartcontractkit/chainlink/core/services"
	"github.com/smartcontractkit/chainlink/core/services/pg"
	"github.com/smartcontractkit/chainlink/core/utils"
)

//go:generate mockery --name Broadcaster --output ./mocks/ --case=underscore --structname Broadcaster --filename broadcaster.go
//go:generate mockery --name Listener --output ./mocks/ --case=underscore --structname Listener --filename listener.go
//go:generate mockery --name Config --output ./mocks/ --case=underscore --structname Config --filename config.go

type (
	// The Broadcaster manages log subscription requests for the Chainlink node.  Instead
	// of creating a new subscription for each request, it multiplexes all subscriptions
	// to all of the relevant contracts over a single connection and forwards the logs to the
	// relevant subscribers.
	//
	// In case of node crash and/or restart, the logs will be backfilled for subscribers that are added before all
	// dependents of LogBroadcaster are done.
	//
	// The backfill starts from the earliest block of either:
	//  - Latest DB head minus BlockBackfillDepth and the maximum number of confirmations.
	//  - Earliest pending or unconsumed log broadcast from DB.
	//
	// If a subscriber is added after the LogBroadcaster does the initial backfill,
	// then it's possible/likely that the backfill fill only have depth: 1 (from latest head)
	//
	// Of course, these backfilled logs + any new logs will only be sent after the NumConfirmations for given subscriber.
	Broadcaster interface {
		utils.DependentAwaiter
		services.Service
		httypes.HeadTrackable
		ReplayFromBlock(number int64)

		IsConnected() bool
		Register(listener Listener, opts ListenerOpts) (unsubscribe func())

		WasAlreadyConsumed(lb Broadcast, qopts ...pg.QOpt) (bool, error)
		MarkConsumed(lb Broadcast, qopts ...pg.QOpt) error
		// NOTE: WasAlreadyConsumed and MarkConsumed MUST be used within a single goroutine in order for WasAlreadyConsumed to be accurate
	}

	BroadcasterInTest interface {
		Broadcaster
		BackfillBlockNumber() null.Int64
		TrackedAddressesCount() uint32
		// Pause pauses the eventLoop until Resume is called.
		Pause()
		// Resume resumes the eventLoop after calling Pause.
		Resume()
		LogsFromBlock(bh common.Hash) int
	}

	broadcaster struct {
		orm        ORM
		config     Config
		connected  atomic.Bool
		evmChainID big.Int

		// a block number to start backfill from
		backfillBlockNumber null.Int64

		ethSubscriber *ethSubscriber
		registrations *registrations
		logPool       iLogPool

		addSubscriber *utils.Mailbox
		rmSubscriber  *utils.Mailbox
		newHeads      *utils.Mailbox

		utils.StartStopOnce
		utils.DependentAwaiter

		chStop                chan struct{}
		wgDone                sync.WaitGroup
		trackedAddressesCount atomic.Uint32
		replayChannel         chan int64
		highestSavedHead      *evmtypes.Head
		lastSeenHeadNumber    atomic.Int64
		logger                logger.Logger

		// used for testing only
		testPause, testResume chan struct{}
	}

	Config interface {
		BlockBackfillDepth() uint64
		BlockBackfillSkip() bool
		EvmFinalityDepth() uint32
		EvmLogBackfillBatchSize() uint32
	}

	ListenerOpts struct {
		Contract common.Address

		// Event types to receive, with value filter for each field in the event
		// No filter or an empty filter for a given field position mean: all values allowed
		// the key should be a result of AbigenLog.Topic() call
		LogsWithTopics map[common.Hash][][]Topic

		ParseLog ParseLogFunc

		// Minimum number of block confirmations before the log is received
		MinIncomingConfirmations uint32
	}

	ParseLogFunc func(log types.Log) (generated.AbigenLog, error)

	registration struct {
		listener Listener
		opts     ListenerOpts
	}

	Topic common.Hash
)

var _ Broadcaster = (*broadcaster)(nil)

// NewBroadcaster creates a new instance of the broadcaster
func NewBroadcaster(orm ORM, ethClient evmclient.Client, config Config, lggr logger.Logger, highestSavedHead *evmtypes.Head) *broadcaster {
	chStop := make(chan struct{})
	lggr = lggr.Named("LogBroadcaster")
	return &broadcaster{
		orm:              orm,
		config:           config,
		logger:           lggr,
		evmChainID:       *ethClient.ChainID(),
		ethSubscriber:    newEthSubscriber(ethClient, config, lggr, chStop),
		registrations:    newRegistrations(lggr, *ethClient.ChainID()),
		logPool:          newLogPool(),
		addSubscriber:    utils.NewMailbox(0),
		rmSubscriber:     utils.NewMailbox(0),
		newHeads:         utils.NewMailbox(1),
		DependentAwaiter: utils.NewDependentAwaiter(),
		chStop:           chStop,
		highestSavedHead: highestSavedHead,
		replayChannel:    make(chan int64, 1),
	}
}

func (b *broadcaster) Start() error {
	return b.StartOnce("LogBroadcaster", func() error {
		b.wgDone.Add(2)
		go b.awaitInitialSubscribers()
		return nil
	})
}

func (b *broadcaster) ReplayFromBlock(number int64) {
	b.logger.Infof("Replay requested from block number: %v", number)
	select {
	case b.replayChannel <- number:
	default:
	}
}

func (b *broadcaster) Close() error {
	return b.StopOnce("LogBroadcaster", func() error {
		close(b.chStop)
		b.wgDone.Wait()
		return nil
	})
}

func (b *broadcaster) awaitInitialSubscribers() {
	defer b.wgDone.Done()
	b.logger.Debug("Starting to await initial subscribers until all dependents are ready...")
	for {
		select {
		case <-b.addSubscriber.Notify():
			b.onAddSubscribers()

		case <-b.rmSubscriber.Notify():
			b.onRmSubscribers()

		case <-b.DependentAwaiter.AwaitDependents():
			// ensure that any queued dependent subscriptions are registered first
			b.onAddSubscribers()
			go b.startResubscribeLoop()
			return

		case <-b.chStop:
			b.wgDone.Done() // because startResubscribeLoop won't be called
			return
		}
	}
}

func (b *broadcaster) Register(listener Listener, opts ListenerOpts) (unsubscribe func()) {
	if len(opts.LogsWithTopics) == 0 {
		b.logger.Panic("Must supply at least 1 LogsWithTopics element to Register")
	}

	reg := registration{listener, opts}
	wasOverCapacity := b.addSubscriber.Deliver(reg)
	if wasOverCapacity {
		b.logger.Error("Subscription mailbox is over capacity - dropped the oldest unprocessed subscription")
	}
	return func() {
		wasOverCapacity := b.rmSubscriber.Deliver(reg)
		if wasOverCapacity {
			b.logger.Error("Subscription removal mailbox is over capacity - dropped the oldest unprocessed removal")
		}
	}
}

func (b *broadcaster) OnNewLongestChain(ctx context.Context, head *evmtypes.Head) {
	wasOverCapacity := b.newHeads.Deliver(head)
	if wasOverCapacity {
		b.logger.Debugw("TRACE: Dropped the older head in the mailbox, while inserting latest (which is fine)", "latestBlockNumber", head.Number)
	}
}

func (b *broadcaster) IsConnected() bool {
	return b.connected.Load()
}

// The subscription is closed in two cases:
//   - intentionally, when the set of contracts we're listening to changes
//   - on a connection error
//
// This method recreates the subscription in both cases.  In the event of a connection
// error, it attempts to reconnect.  Any time there's a change in connection state, it
// notifies its subscribers.
func (b *broadcaster) startResubscribeLoop() {
	defer b.wgDone.Done()

	var subscription managedSubscription = newNoopSubscription()
	defer func() { subscription.Unsubscribe() }()

	if b.config.BlockBackfillSkip() && b.highestSavedHead != nil {
		b.logger.Warn("BlockBackfillSkip is set to true, preventing a deep backfill - some earlier chain events might be missed.")
	} else if b.highestSavedHead != nil {
		// The backfill needs to start at an earlier block than the one last saved in DB, to account for:
		// - keeping logs in the in-memory buffers in registration.go
		//   (which will be lost on node restart) for MAX(NumConfirmations of subscribers)
		// - HeadTracker saving the heads to DB asynchronously versus LogBroadcaster, where a head
		//   (or more heads on fast chains) may be saved but not yet processed by LB
		//   using BlockBackfillDepth makes sure the backfill will be dependent on the per-chain configuration
		from := b.highestSavedHead.Number -
			int64(b.registrations.highestNumConfirmations) -
			int64(b.config.BlockBackfillDepth())
		if from < 0 {
			from = 0
		}
		b.backfillBlockNumber = null.NewInt64(from, true)
	}

	// Remove leftover unconsumed logs, maybe update pending broadcasts, and backfill sooner if necessary.
	if backfillStart, abort := b.reinitialize(); abort {
		return
	} else if backfillStart != nil {
		if !b.backfillBlockNumber.Valid || *backfillStart < b.backfillBlockNumber.Int64 {
			b.backfillBlockNumber.SetValid(*backfillStart)
		}
	}

	var chRawLogs chan types.Log
	for {
		b.logger.Infow("Resubscribing and backfilling logs...")
		addresses, topics := b.registrations.addressesAndTopics()

		newSubscription, abort := b.ethSubscriber.createSubscription(addresses, topics)
		if abort {
			return
		}

		if b.backfillBlockNumber.Valid {
			b.logger.Debugw("Using an override as a start of the backfill",
				"blockNumber", b.backfillBlockNumber.Int64,
				"highestNumConfirmations", b.registrations.highestNumConfirmations,
				"blockBackfillDepth", b.config.BlockBackfillDepth(),
			)
		}

		chBackfilledLogs, abort := b.ethSubscriber.backfillLogs(b.backfillBlockNumber, addresses, topics)
		if abort {
			return
		}

		b.backfillBlockNumber.Valid = false

		// Each time this loop runs, chRawLogs is reconstituted as:
		// "remaining logs from last subscription <- backfilled logs <- logs from new subscription"
		// There will be duplicated logs in this channel.  It is the responsibility of subscribers
		// to account for this using the helpers on the Broadcast type.
		chRawLogs = b.appendLogChannel(chRawLogs, chBackfilledLogs)
		chRawLogs = b.appendLogChannel(chRawLogs, newSubscription.Logs())
		subscription.Unsubscribe()
		subscription = newSubscription

		b.connected.Store(true)

		b.trackedAddressesCount.Store(uint32(len(addresses)))

		shouldResubscribe, err := b.eventLoop(chRawLogs, subscription.Err())
		if err != nil {
			b.logger.Warnw("Error in the event loop - will reconnect", "err", err)
			b.connected.Store(false)
			continue
		} else if !shouldResubscribe {
			b.connected.Store(false)
			return
		}
	}
}

func (b *broadcaster) reinitialize() (backfillStart *int64, abort bool) {
	ctx, cancel := utils.ContextFromChan(b.chStop)
	defer cancel()

	utils.RetryWithBackoff(ctx, func() bool {
		var err error
		backfillStart, err = b.orm.Reinitialize(pg.WithParentCtx(ctx))
		if err != nil {
			b.logger.Errorw("Failed to reinitialize database", "err", err)
			return true
		}
		return false
	})

	select {
	case <-b.chStop:
		abort = true
	default:
	}
	return
}

func (b *broadcaster) eventLoop(chRawLogs <-chan types.Log, chErr <-chan error) (shouldResubscribe bool, _ error) {
	// We debounce requests to subscribe and unsubscribe to avoid making too many
	// RPC calls to the Ethereum node, particularly on startup.
	var needsResubscribe bool
	debounceResubscribe := time.NewTicker(1 * time.Second)
	defer debounceResubscribe.Stop()

	b.logger.Debug("Starting the event loop")
	for {
		select {
		case rawLog := <-chRawLogs:

			b.logger.Debugw("Received a log",
				"blockNumber", rawLog.BlockNumber, "blockHash", rawLog.BlockHash, "address", rawLog.Address)

			b.onNewLog(rawLog)

		case <-b.newHeads.Notify():
			b.onNewHeads()

		case err := <-chErr:
			// Note we'll get a message on this channel
			// if the eth node terminates the connection.
			return true, err

		case <-b.addSubscriber.Notify():
			needsResubscribe = b.onAddSubscribers() || needsResubscribe

		case <-b.rmSubscriber.Notify():
			needsResubscribe = b.onRmSubscribers() || needsResubscribe

		case blockNumber := <-b.replayChannel:
			b.backfillBlockNumber.SetValid(blockNumber)
			b.logger.Debugw("Returning from the event loop to replay logs from specific block number", "blockNumber", blockNumber)
			return true, nil

		case <-debounceResubscribe.C:
			if needsResubscribe {
				b.logger.Debug("Returning from the event loop to resubscribe")
				return true, nil
			}

		case <-b.chStop:
			return false, nil

		// testing only
		case <-b.testPause:
			select {
			case <-b.testResume:
			case <-b.chStop:
				return false, nil
			}
		}
	}
}

func (b *broadcaster) onNewLog(log types.Log) {
	b.maybeWarnOnLargeBlockNumberDifference(int64(log.BlockNumber))

	if log.Removed {
		// Remove the whole block that contained this log.
		b.logPool.removeBlock(log.BlockHash, log.BlockNumber)
		return
	} else if !b.registrations.isAddressRegistered(log.Address) {
		return
	}
	if b.logPool.addLog(log) {
		// First or new lowest block number
		ctx, cancel := utils.ContextFromChan(b.chStop)
		defer cancel()
		blockNumber := int64(log.BlockNumber)
		if err := b.orm.SetPendingMinBlock(&blockNumber, pg.WithParentCtx(ctx)); err != nil {
			b.logger.Errorw("Failed to set pending broadcasts number", "blockNumber", log.BlockNumber, "err", err)
		}
	}
}

func (b *broadcaster) onNewHeads() {
	var latestHead *evmtypes.Head
	for {
		// We only care about the most recent head
		item := b.newHeads.RetrieveLatestAndClear()
		if item == nil {
			break
		}
		head := evmtypes.AsHead(item)
		latestHead = head
	}

	// latestHead may sometimes be nil on high rate of heads,
	// when 'b.newHeads.Notify()' receives more times that the number of items in the mailbox
	// Some heads may be missed (which is fine for LogBroadcaster logic) but the latest one in a burst will be received
	if latestHead != nil {
		b.logger.Debugw("Received head", "blockNumber", latestHead.Number,
			"blockHash", latestHead.Hash, "parentHash", latestHead.ParentHash, "chainLen", latestHead.ChainLength())

		b.lastSeenHeadNumber.Store(latestHead.Number)

		keptLogsDepth := uint32(b.config.EvmFinalityDepth())
		if b.registrations.highestNumConfirmations > keptLogsDepth {
			keptLogsDepth = b.registrations.highestNumConfirmations
		}

		latestBlockNum := latestHead.Number
		keptDepth := latestBlockNum - int64(keptLogsDepth)
		if keptDepth < 0 {
			keptDepth = 0
		}

		ctx, cancel := utils.ContextFromChan(b.chStop)
		defer cancel()

		// if all subscribers requested 0 confirmations, we always get and delete all logs from the pool,
		// without comparing their block numbers to the current head's block number.
		if b.registrations.highestNumConfirmations == 0 {
			logs, lowest, highest := b.logPool.getAndDeleteAll()
			if len(logs) > 0 {
				broadcasts, err := b.orm.FindBroadcasts(lowest, highest)
				if err != nil {
					b.logger.Errorf("Failed to query for log broadcasts, %v", err)
					return
				}
				b.registrations.sendLogs(logs, *latestHead, broadcasts, b.orm)
				if err := b.orm.SetPendingMinBlock(nil, pg.WithParentCtx(ctx)); err != nil {
					b.logger.Errorw("Failed to set pending broadcasts number null", "err", err)
				}
			}
		} else {
			logs, minBlockNum := b.logPool.getLogsToSend(latestBlockNum)

			if len(logs) > 0 {
				broadcasts, err := b.orm.FindBroadcasts(minBlockNum, latestBlockNum)
				if err != nil {
					b.logger.Errorf("Failed to query for log broadcasts, %v", err)
					return
				}

				b.registrations.sendLogs(logs, *latestHead, broadcasts, b.orm)
			}
			newMin := b.logPool.deleteOlderLogs(keptDepth)
			if err := b.orm.SetPendingMinBlock(newMin); err != nil {
				b.logger.Errorw("Failed to set pending broadcasts number", "blockNumber", keptDepth, "err", err)
			}
		}
	}
}

func (b *broadcaster) onAddSubscribers() (needsResubscribe bool) {
	for {
		x, exists := b.addSubscriber.Retrieve()
		if !exists {
			break
		}
		reg, ok := x.(registration)
		if !ok {
			b.logger.Errorf("expected `registration`, got %T", x)
			continue
		}
		b.logger.Debugw("Subscribing listener", "requiredBlockConfirmations", reg.opts.MinIncomingConfirmations, "address", reg.opts.Contract)
		needsResub := b.registrations.addSubscriber(reg)
		if needsResub {
			needsResubscribe = true
		}
	}
	return
}

func (b *broadcaster) onRmSubscribers() (needsResubscribe bool) {
	for {
		x, exists := b.rmSubscriber.Retrieve()
		if !exists {
			break
		}
		reg, ok := x.(registration)
		if !ok {
			b.logger.Errorf("expected `registration`, got %T", x)
			continue
		}
		b.logger.Debugw("Unsubscribing listener", "requiredBlockConfirmations", reg.opts.MinIncomingConfirmations, "address", reg.opts.Contract)
		needsResub := b.registrations.removeSubscriber(reg)
		if needsResub {
			needsResubscribe = true
		}
	}
	return
}

func (b *broadcaster) appendLogChannel(ch1, ch2 <-chan types.Log) chan types.Log {
	if ch1 == nil && ch2 == nil {
		return nil
	}

	chCombined := make(chan types.Log)

	go func() {
		defer close(chCombined)
		if ch1 != nil {
			for rawLog := range ch1 {
				select {
				case chCombined <- rawLog:
				case <-b.chStop:
					return
				}
			}
		}
		if ch2 != nil {
			for rawLog := range ch2 {
				select {
				case chCombined <- rawLog:
				case <-b.chStop:
					return
				}
			}
		}
	}()

	return chCombined
}

func (b *broadcaster) maybeWarnOnLargeBlockNumberDifference(logBlockNumber int64) {
	lastSeenHeadNumber := b.lastSeenHeadNumber.Load()
	diff := logBlockNumber - lastSeenHeadNumber
	if diff < 0 {
		diff = -diff
	}

	if lastSeenHeadNumber > 0 && diff > 1000 {
		b.logger.Warnw("Detected a large block number difference between a log and recently seen head. "+
			"This may indicate a problem with data received from the chain or major network delays.",
			"lastSeenHeadNumber", lastSeenHeadNumber, "logBlockNumber", logBlockNumber, "diff", diff)
	}
}

// WasAlreadyConsumed reports whether the given consumer had already consumed the given log
func (b *broadcaster) WasAlreadyConsumed(lb Broadcast, qopts ...pg.QOpt) (bool, error) {
	return b.orm.WasBroadcastConsumed(lb.RawLog().BlockHash, lb.RawLog().Index, lb.JobID(), qopts...)
}

// MarkConsumed marks the log as having been successfully consumed by the subscriber
func (b *broadcaster) MarkConsumed(lb Broadcast, qopts ...pg.QOpt) error {
	return b.orm.MarkBroadcastConsumed(lb.RawLog().BlockHash, lb.RawLog().BlockNumber, lb.RawLog().Index, lb.JobID(), qopts...)
}

// test only
func (b *broadcaster) TrackedAddressesCount() uint32 {
	return b.trackedAddressesCount.Load()
}

// test only
func (b *broadcaster) BackfillBlockNumber() null.Int64 {
	return b.backfillBlockNumber
}

// test only
func (b *broadcaster) Pause() {
	select {
	case b.testPause <- struct{}{}:
	case <-b.chStop:
	}
}

// test only
func (b *broadcaster) Resume() {
	select {
	case b.testResume <- struct{}{}:
	case <-b.chStop:
	}
}

// test only
func (b *broadcaster) LogsFromBlock(bh common.Hash) int {
	return b.logPool.testOnly_getNumLogsForBlock(bh)
}

var _ BroadcasterInTest = &NullBroadcaster{}

type NullBroadcaster struct{ ErrMsg string }

func (n *NullBroadcaster) IsConnected() bool { return false }
func (n *NullBroadcaster) Register(listener Listener, opts ListenerOpts) (unsubscribe func()) {
	return func() {}
}

func (n *NullBroadcaster) ReplayFromBlock(number int64) {}

func (n *NullBroadcaster) BackfillBlockNumber() null.Int64 {
	return null.NewInt64(0, false)
}
func (n *NullBroadcaster) TrackedAddressesCount() uint32 {
	return 0
}
func (n *NullBroadcaster) WasAlreadyConsumed(lb Broadcast, qopts ...pg.QOpt) (bool, error) {
	return false, errors.New(n.ErrMsg)
}
func (n *NullBroadcaster) MarkConsumed(lb Broadcast, qopts ...pg.QOpt) error {
	return errors.New(n.ErrMsg)
}

func (n *NullBroadcaster) AddDependents(int) {}
func (n *NullBroadcaster) AwaitDependents() <-chan struct{} {
	ch := make(chan struct{})
	close(ch)
	return ch
}
func (n *NullBroadcaster) DependentReady()                                   {}
func (n *NullBroadcaster) Start() error                                      { return nil }
func (n *NullBroadcaster) Close() error                                      { return nil }
func (n *NullBroadcaster) Healthy() error                                    { return nil }
func (n *NullBroadcaster) Ready() error                                      { return nil }
func (n *NullBroadcaster) OnNewLongestChain(context.Context, *evmtypes.Head) {}
func (n *NullBroadcaster) Pause()                                            {}
func (n *NullBroadcaster) Resume()                                           {}
func (n *NullBroadcaster) LogsFromBlock(common.Hash) int                     { return -1 }
