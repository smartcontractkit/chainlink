package log

import (
	"context"
	"fmt"
	"math/big"
	"sync"
	"sync/atomic"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"

	evmclient "github.com/smartcontractkit/chainlink/core/chains/evm/client"
	httypes "github.com/smartcontractkit/chainlink/core/chains/evm/headtracker/types"
	evmtypes "github.com/smartcontractkit/chainlink/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/core/gethwrappers/generated"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/null"
	"github.com/smartcontractkit/chainlink/core/services"
	"github.com/smartcontractkit/chainlink/core/services/pg"
	"github.com/smartcontractkit/chainlink/core/utils"
)

//go:generate mockery --quiet --name Broadcaster --output ./mocks/ --case=underscore --structname Broadcaster --filename broadcaster.go

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
		services.ServiceCtx
		httypes.HeadTrackable

		// ReplayFromBlock enqueues a replay from the provided block number. If forceBroadcast is
		// set to true, the broadcaster will broadcast logs that were already marked consumed
		// previously by any subscribers.
		ReplayFromBlock(number int64, forceBroadcast bool)

		IsConnected() bool
		Register(listener Listener, opts ListenerOpts) (unsubscribe func())

		WasAlreadyConsumed(lb Broadcast, qopts ...pg.QOpt) (bool, error)
		MarkConsumed(lb Broadcast, qopts ...pg.QOpt) error

		// MarkManyConsumed marks all the provided log broadcasts as consumed.
		MarkManyConsumed(lbs []Broadcast, qopts ...pg.QOpt) error

		// NOTE: WasAlreadyConsumed, MarkConsumed and MarkManyConsumed MUST be used within a single goroutine in order for WasAlreadyConsumed to be accurate
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

	subscriberStatus int

	changeSubscriberStatus struct {
		newStatus subscriberStatus
		sub       *subscriber
	}

	replayRequest struct {
		fromBlock      int64
		forceBroadcast bool
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
		logPool       *logPool

		mailMon *utils.MailboxMonitor
		// Use the same channel for subs/unsubs so ordering is preserved
		// (unsubscribe must happen after subscribe)
		changeSubscriberStatus *utils.Mailbox[changeSubscriberStatus]
		newHeads               *utils.Mailbox[*evmtypes.Head]

		utils.StartStopOnce
		utils.DependentAwaiter

		chStop                chan struct{}
		wgDone                sync.WaitGroup
		trackedAddressesCount atomic.Uint32
		replayChannel         chan replayRequest
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
		// topic => topicValueFilters
		LogsWithTopics map[common.Hash][][]Topic

		ParseLog ParseLogFunc

		// Minimum number of block confirmations before the log is received
		MinIncomingConfirmations uint32

		// ReplayStartedCallback is called by the log broadcaster once a replay request is received.
		ReplayStartedCallback func()
	}

	ParseLogFunc func(log types.Log) (generated.AbigenLog, error)

	subscriber struct {
		listener Listener
		opts     ListenerOpts
	}

	Topic common.Hash
)

const (
	subscriberStatusSubscribe = iota
	subscriberStatusUnsubscribe
)

var _ Broadcaster = (*broadcaster)(nil)

// NewBroadcaster creates a new instance of the broadcaster
func NewBroadcaster(orm ORM, ethClient evmclient.Client, config Config, lggr logger.Logger, highestSavedHead *evmtypes.Head, mailMon *utils.MailboxMonitor) *broadcaster {
	chStop := make(chan struct{})
	lggr = lggr.Named("LogBroadcaster")
	id := ethClient.ChainID()
	return &broadcaster{
		orm:                    orm,
		config:                 config,
		logger:                 lggr,
		evmChainID:             *id,
		ethSubscriber:          newEthSubscriber(ethClient, config, lggr, chStop),
		registrations:          newRegistrations(lggr, *id),
		logPool:                newLogPool(lggr),
		mailMon:                mailMon,
		changeSubscriberStatus: utils.NewHighCapacityMailbox[changeSubscriberStatus](),
		newHeads:               utils.NewSingleMailbox[*evmtypes.Head](),
		DependentAwaiter:       utils.NewDependentAwaiter(),
		chStop:                 chStop,
		highestSavedHead:       highestSavedHead,
		replayChannel:          make(chan replayRequest, 1),
	}
}

func (b *broadcaster) Start(context.Context) error {
	return b.StartOnce("LogBroadcaster", func() error {
		b.wgDone.Add(2)
		go b.awaitInitialSubscribers()
		b.mailMon.Monitor(b.changeSubscriberStatus, "LogBroadcaster", "ChangeSubscriber", b.evmChainID.String())
		return nil
	})
}

// ReplayFromBlock implements the Broadcaster interface.
func (b *broadcaster) ReplayFromBlock(number int64, forceBroadcast bool) {
	b.logger.Infow("Replay requested", "block number", number, "force", forceBroadcast)
	select {
	case b.replayChannel <- replayRequest{
		fromBlock:      number,
		forceBroadcast: forceBroadcast,
	}:
	default:
	}
}

func (b *broadcaster) Close() error {
	return b.StopOnce("LogBroadcaster", func() error {
		close(b.chStop)
		b.wgDone.Wait()
		return b.changeSubscriberStatus.Close()
	})
}

func (b *broadcaster) Name() string {
	return b.logger.Name()
}

func (b *broadcaster) HealthReport() map[string]error {
	return map[string]error{b.Name(): b.Healthy()}
}

func (b *broadcaster) awaitInitialSubscribers() {
	defer b.wgDone.Done()
	b.logger.Debug("Starting to await initial subscribers until all dependents are ready...")
	for {
		select {
		case <-b.changeSubscriberStatus.Notify():
			b.onChangeSubscriberStatus()

		case <-b.DependentAwaiter.AwaitDependents():
			// ensure that any queued dependent subscriptions are registered first
			b.onChangeSubscriberStatus()
			go b.startResubscribeLoop()
			return

		case <-b.chStop:
			b.wgDone.Done() // because startResubscribeLoop won't be called
			return
		}
	}
}

func (b *broadcaster) Register(listener Listener, opts ListenerOpts) (unsubscribe func()) {
	// IfNotStopped RLocks the state mutex so LB cannot be closed until this
	// returns (no need to worry about listening for b.chStop)
	//
	// NOTE: We do not use IfStarted here because it is explicitly ok to
	// register listeners before starting, this allows us to register many
	// listeners then subscribe once on start, avoiding thrashing
	ok := b.IfNotStopped(func() {
		if len(opts.LogsWithTopics) == 0 {
			b.logger.Panic("Must supply at least 1 LogsWithTopics element to Register")
		}
		if opts.MinIncomingConfirmations <= 0 {
			b.logger.Warnw(fmt.Sprintf("LogBroadcaster requires that MinIncomingConfirmations must be at least 1 (got %v). Logs must have been confirmed in at least 1 block, it does not support reading logs from the mempool before they have been mined. MinIncomingConfirmations will be set to 1.", opts.MinIncomingConfirmations), "addr", opts.Contract.Hex(), "jobID", listener.JobID())
			opts.MinIncomingConfirmations = 1
		}

		sub := &subscriber{listener, opts}
		b.logger.Debugf("Registering subscriber %p with job ID %v", sub, sub.listener.JobID())
		wasOverCapacity := b.changeSubscriberStatus.Deliver(changeSubscriberStatus{subscriberStatusSubscribe, sub})
		if wasOverCapacity {
			b.logger.Panicf("LogBroadcaster subscribe: cannot subscribe %p with job ID %v; changeSubscriberStatus channel was full", sub, sub.listener.JobID())
		}

		// this is asynchronous but it shouldn't matter, since the channel is
		// ordered then it will work properly as long as you call unsubscribe
		// before subscribing a new listener with the same job/addr (e.g. on
		// replacement of the same job)
		unsubscribe = func() {
			b.logger.Debugf("Unregistering subscriber %p with job ID %v", sub, sub.listener.JobID())
			wasOverCapacity := b.changeSubscriberStatus.Deliver(changeSubscriberStatus{subscriberStatusUnsubscribe, sub})
			if wasOverCapacity {
				b.logger.Panicf("LogBroadcaster unsubscribe: cannot unsubscribe %p with job ID %v; changeSubscriberStatus channel was full", sub, sub.listener.JobID())
			}
		}
	})
	if !ok {
		b.logger.Panic("Register cannot be called on a stopped log broadcaster (this is an invariant violation because all dependent services should have unregistered themselves before logbroadcaster.Close was called)")
	}
	return
}

func (b *broadcaster) OnNewLongestChain(ctx context.Context, head *evmtypes.Head) {
	wasOverCapacity := b.newHeads.Deliver(head)
	if wasOverCapacity {
		b.logger.Debugw("Dropped the older head in the mailbox, while inserting latest (which is fine)", "latestBlockNumber", head.Number)
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
		// No need to worry about r.highestNumConfirmations here because it's
		// already at minimum this deep due to the latest seen head check above
		if !b.backfillBlockNumber.Valid || *backfillStart < b.backfillBlockNumber.Int64 {
			b.backfillBlockNumber.SetValid(*backfillStart)
		}
	}

	if b.backfillBlockNumber.Valid {
		b.logger.Debugw("Using an override as a start of the backfill",
			"blockNumber", b.backfillBlockNumber.Int64,
			"highestNumConfirmations", b.registrations.highestNumConfirmations,
			"blockBackfillDepth", b.config.BlockBackfillDepth(),
		)
	}

	var chRawLogs chan types.Log
	for {
		b.logger.Infow("Resubscribing and backfilling logs...")
		addresses, topics := b.registrations.addressesAndTopics()

		newSubscription, abort := b.ethSubscriber.createSubscription(addresses, topics)
		if abort {
			return
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
		// Replay requests take priority.
		select {
		case req := <-b.replayChannel:
			b.onReplayRequest(req)
			return true, nil
		default:
		}

		select {
		case rawLog := <-chRawLogs:
			b.logger.Debugw("Received a log",
				"blockNumber", rawLog.BlockNumber, "blockHash", rawLog.BlockHash, "address", rawLog.Address)
			b.onNewLog(rawLog)

		case <-b.newHeads.Notify():
			b.onNewHeads()

		case err := <-chErr:
			// The eth node connection was terminated so we need to backfill after resubscribing.
			lggr := b.logger
			// Do we have logs in the pool?
			// They are are invalid, since we may have missed 'removed' logs.
			if blockNum := b.invalidatePool(); blockNum > 0 {
				lggr = lggr.With("blockNumber", blockNum)
			}
			lggr.Debugw("Subscription terminated. Backfilling after resubscribing")
			return true, err

		case <-b.changeSubscriberStatus.Notify():
			needsResubscribe = b.onChangeSubscriberStatus() || needsResubscribe

		case req := <-b.replayChannel:
			b.onReplayRequest(req)
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

// onReplayRequest clears the pool and sets the block backfill number.
func (b *broadcaster) onReplayRequest(replayReq replayRequest) {
	// notify subscribers that we are about to replay.
	for subscriber := range b.registrations.registeredSubs {
		if subscriber.opts.ReplayStartedCallback != nil {
			subscriber.opts.ReplayStartedCallback()
		}
	}

	_ = b.invalidatePool()
	// NOTE: This ignores r.highestNumConfirmations, but it is
	// generally assumed that this will only be performed rarely and
	// manually by someone who knows what he is doing
	b.backfillBlockNumber.SetValid(replayReq.fromBlock)
	if replayReq.forceBroadcast {
		ctx, cancel := utils.ContextFromChan(b.chStop)
		defer cancel()

		// Use a longer timeout in the event that a very large amount of logs need to be marked
		// as consumed.
		err := b.orm.MarkBroadcastsUnconsumed(replayReq.fromBlock, pg.WithParentCtx(ctx), pg.WithLongQueryTimeout())
		if err != nil {
			b.logger.Errorw("Error marking broadcasts as unconsumed",
				"error", err, "fromBlock", replayReq.fromBlock)
		}
	}
	b.logger.Debugw(
		"Returning from the event loop to replay logs from specific block number",
		"fromBlock", replayReq.fromBlock,
		"forceBroadcast", replayReq.forceBroadcast,
	)
}

func (b *broadcaster) invalidatePool() int64 {
	if min := b.logPool.heap.FindMin(); min != nil {
		b.logPool = newLogPool(b.logger)
		// Note: even if we crash right now, PendingMinBlock is preserved in the database and we will backfill the same.
		blockNum := int64(min.(Uint64))
		b.backfillBlockNumber.SetValid(blockNum)
		return blockNum
	}
	return -1
}

func (b *broadcaster) onNewLog(log types.Log) {
	b.maybeWarnOnLargeBlockNumberDifference(int64(log.BlockNumber))

	if log.Removed {
		// Remove the whole block that contained this log.
		b.logger.Debugw("Found reverted log", "log", log)
		b.logPool.removeBlock(log.BlockHash, log.BlockNumber)
		return
	} else if !b.registrations.isAddressRegistered(log.Address) {
		b.logger.Debugw("Found unregistered address", "address", log.Address)
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
		head := b.newHeads.RetrieveLatestAndClear()
		if head == nil {
			break
		}
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

func (b *broadcaster) onChangeSubscriberStatus() (needsResubscribe bool) {
	for {
		change, exists := b.changeSubscriberStatus.Retrieve()
		if !exists {
			break
		}
		sub := change.sub

		if change.newStatus == subscriberStatusSubscribe {
			b.logger.Debugw("Subscribing listener", "requiredBlockConfirmations", sub.opts.MinIncomingConfirmations, "address", sub.opts.Contract, "jobID", sub.listener.JobID())
			needsResub := b.registrations.addSubscriber(sub)
			if needsResub {
				needsResubscribe = true
			}
		} else {
			b.logger.Debugw("Unsubscribing listener", "requiredBlockConfirmations", sub.opts.MinIncomingConfirmations, "address", sub.opts.Contract, "jobID", sub.listener.JobID())
			needsResub := b.registrations.removeSubscriber(sub)
			if needsResub {
				needsResubscribe = true
			}
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

// MarkManyConsumed marks the logs as having been successfully consumed by the subscriber
func (b *broadcaster) MarkManyConsumed(lbs []Broadcast, qopts ...pg.QOpt) (err error) {
	var (
		blockHashes  = make([]common.Hash, len(lbs))
		blockNumbers = make([]uint64, len(lbs))
		logIndexes   = make([]uint, len(lbs))
		jobIDs       = make([]int32, len(lbs))
	)
	for i := range lbs {
		blockHashes[i] = lbs[i].RawLog().BlockHash
		blockNumbers[i] = lbs[i].RawLog().BlockNumber
		logIndexes[i] = lbs[i].RawLog().Index
		jobIDs[i] = lbs[i].JobID()
	}
	return b.orm.MarkBroadcastsConsumed(blockHashes, blockNumbers, logIndexes, jobIDs, qopts...)
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

func topicsToHex(topics [][]Topic) [][]common.Hash {
	var topicsInHex [][]common.Hash
	for i := range topics {
		var hexes []common.Hash
		for j := range topics[i] {
			hexes = append(hexes, common.Hash(topics[i][j]))
		}
		topicsInHex = append(topicsInHex, hexes)
	}
	return topicsInHex
}

var _ BroadcasterInTest = &NullBroadcaster{}

type NullBroadcaster struct{ ErrMsg string }

func (n *NullBroadcaster) IsConnected() bool { return false }
func (n *NullBroadcaster) Register(listener Listener, opts ListenerOpts) (unsubscribe func()) {
	return func() {}
}

// ReplayFromBlock implements the Broadcaster interface.
func (n *NullBroadcaster) ReplayFromBlock(number int64, forceBroadcast bool) {}

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
func (n *NullBroadcaster) MarkManyConsumed(lbs []Broadcast, qopts ...pg.QOpt) error {
	return errors.New(n.ErrMsg)
}

func (n *NullBroadcaster) AddDependents(int) {}
func (n *NullBroadcaster) AwaitDependents() <-chan struct{} {
	ch := make(chan struct{})
	close(ch)
	return ch
}

// DependentReady does noop for NullBroadcaster.
func (n *NullBroadcaster) DependentReady() {}

func (n *NullBroadcaster) Name() string { return "" }

// Start does noop for NullBroadcaster.
func (n *NullBroadcaster) Start(context.Context) error                       { return nil }
func (n *NullBroadcaster) Close() error                                      { return nil }
func (n *NullBroadcaster) Healthy() error                                    { return nil }
func (n *NullBroadcaster) Ready() error                                      { return nil }
func (n *NullBroadcaster) HealthReport() map[string]error                    { return nil }
func (n *NullBroadcaster) OnNewLongestChain(context.Context, *evmtypes.Head) {}
func (n *NullBroadcaster) Pause()                                            {}
func (n *NullBroadcaster) Resume()                                           {}
func (n *NullBroadcaster) LogsFromBlock(common.Hash) int                     { return -1 }
