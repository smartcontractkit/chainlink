package log

import (
	"context"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/tevino/abool"

	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/eth"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"
)

//go:generate mockery --name Broadcaster --output ./mocks/ --case=underscore --structname Broadcaster --filename broadcaster.go
//go:generate mockery --name Listener --output ./mocks/ --case=underscore --structname Listener --filename listener.go
//go:generate mockery --name AbigenContract --output ./mocks/ --case=underscore --structname AbigenContract --filename abigen_contract.go

type (
	// The Broadcaster manages log subscription requests for the Chainlink node.  Instead
	// of creating a new subscription for each request, it multiplexes all subscriptions
	// to all of the relevant contracts over a single connection and forwards the logs to the
	// relevant subscribers.
	Broadcaster interface {
		utils.DependentAwaiter
		Start() error
		Stop() error
		IsConnected() bool
		Register(listener Listener, opts ListenerOpts) (unsubscribe func())
		SetLatestHeadFromStorage(head *models.Head)
		LatestHead() *models.Head
	}

	broadcaster struct {
		orm        ORM
		config     Config
		connected  *abool.AtomicBool
		latestHead *models.Head

		ethSubscriber *ethSubscriber
		registrations *registrations
		logPool       *logPool

		addSubscriber *utils.Mailbox
		rmSubscriber  *utils.Mailbox
		newHeads      *utils.Mailbox

		utils.StartStopOnce
		utils.DependentAwaiter

		headFromStorageAwaiter utils.DependentAwaiter
		chStop                 chan struct{}
		chDone                 chan struct{}
	}

	Config interface {
		BlockBackfillDepth() uint64
		TriggerFallbackDBPollInterval() time.Duration
	}

	ListenerOpts struct {
		Contract         AbigenContract
		Logs             []generated.AbigenLog
		NumConfirmations uint64
	}

	AbigenContract interface {
		Address() common.Address
		ParseLog(log types.Log) (generated.AbigenLog, error)
	}

	registration struct {
		listener Listener
		opts     ListenerOpts
	}
)

var _ Broadcaster = (*broadcaster)(nil)

// NewBroadcaster creates a new instance of the broadcaster
func NewBroadcaster(orm ORM, ethClient eth.Client, config Config) *broadcaster {
	headAwaiter := utils.NewDependentAwaiter()
	headAwaiter.AddDependents(1)
	chStop := make(chan struct{})
	return &broadcaster{
		orm:                    orm,
		config:                 config,
		connected:              abool.New(),
		ethSubscriber:          newEthSubscriber(ethClient, config, chStop),
		registrations:          newRegistrations(),
		logPool:                newLogPool(),
		addSubscriber:          utils.NewMailbox(0),
		rmSubscriber:           utils.NewMailbox(0),
		newHeads:               utils.NewMailbox(1),
		DependentAwaiter:       utils.NewDependentAwaiter(),
		headFromStorageAwaiter: headAwaiter,
		chStop:                 chStop,
		chDone:                 make(chan struct{}),
	}
}

func (b *broadcaster) Start() error {
	return b.StartOnce("Log broadcaster", func() error {
		go b.awaitInitialSubscribers()
		return nil
	})
}

func (b *broadcaster) SetLatestHeadFromStorage(head *models.Head) {
	b.latestHead = head
	b.headFromStorageAwaiter.DependentReady()
}

func (b *broadcaster) LatestHead() *models.Head {
	return b.latestHead
}

func (b *broadcaster) Stop() error {
	return b.StopOnce("Log broadcaster", func() error {
		close(b.chStop)
		<-b.chDone
		return nil
	})
}

func (b *broadcaster) awaitInitialSubscribers() {
	defer close(b.chDone)
	for {
		select {
		case <-b.addSubscriber.Notify():
			b.onAddSubscribers()

		case <-b.rmSubscriber.Notify():
			b.onRmSubscribers()

		case <-b.DependentAwaiter.AwaitDependents():
			<-b.headFromStorageAwaiter.AwaitDependents()
			go b.startResubscribeLoop()
			return

		case <-b.chStop:
			return
		}
	}
}

func (b *broadcaster) Register(listener Listener, opts ListenerOpts) (unsubscribe func()) {
	if len(opts.Logs) < 1 {
		logger.Fatal("Must supply at least 1 Log to Register")
	}
	b.addSubscriber.Deliver(registration{listener, opts})
	return func() {
		b.rmSubscriber.Deliver(registration{listener, opts})
	}
}

func (b *broadcaster) Connect(head *models.Head) error { return nil }
func (b *broadcaster) Disconnect()                     {}

func (b *broadcaster) OnNewLongestChain(ctx context.Context, head models.Head) {
	b.newHeads.Deliver(head)
}

func (b *broadcaster) IsConnected() bool {
	return b.connected.IsSet()
}

// The subscription is closed in two cases:
//   - intentionally, when the set of contracts we're listening to changes
//   - on a connection error
//
// This method recreates the subscription in both cases.  In the event of a connection
// error, it attempts to reconnect.  Any time there'b a change in connection state, it
// notifies its subscribers.
func (b *broadcaster) startResubscribeLoop() {

	var subscription managedSubscription = newNoopSubscription()
	defer func() { subscription.Unsubscribe() }()

	var chRawLogs chan types.Log
	for {
		logger.Debugf("LogBroadcaster: resubscribing and backfilling logs...")
		addresses, topics := b.registrations.addressesAndTopics()

		newSubscription, abort := b.ethSubscriber.createSubscription(addresses, topics)
		if abort {
			return
		}

		chBackfilledLogs, abort := b.ethSubscriber.backfillLogs(b.latestHead, addresses, topics)
		if abort {
			return
		}

		// Each time this loop runs, chRawLogs is reconstituted as:
		//     remaining logs from last subscription <- backfilled logs <- logs from new subscription
		// There will be duplicated logs in this channel.  It is the responsibility of subscribers
		// to account for this using the helpers on the Broadcast type.
		chRawLogs = b.appendLogChannel(chRawLogs, chBackfilledLogs)
		chRawLogs = b.appendLogChannel(chRawLogs, newSubscription.Logs())
		subscription.Unsubscribe()
		subscription = newSubscription

		b.connected.Set()

		shouldResubscribe, err := b.eventLoop(chRawLogs, subscription.Err())
		if err != nil {
			logger.Warn(err)
			b.connected.UnSet()
			continue
		} else if !shouldResubscribe {
			b.connected.UnSet()
			return
		}
	}
}

func (b *broadcaster) eventLoop(chRawLogs <-chan types.Log, chErr <-chan error) (shouldResubscribe bool, _ error) {
	// We debounce requests to subscribe and unsubscribe to avoid making too many
	// RPC calls to the Ethereum node, particularly on startup.
	var needsResubscribe bool
	debounceResubscribe := time.NewTicker(1 * time.Second)
	defer debounceResubscribe.Stop()

	logger.Debugf("LogBroadcaster: starting the event loop")
	for {
		select {
		case rawLog := <-chRawLogs:
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

		case <-debounceResubscribe.C:
			if needsResubscribe {
				logger.Debugf("LogBroadcaster: returning from the event loop to resubscribe")
				return true, nil
			}

		case <-b.chStop:
			return false, nil
		}
	}
}

func (b *broadcaster) onNewLog(log types.Log) {
	logger.Tracef("========== onNewLog %v, %v, %v", log.BlockNumber, log.BlockHash, log.Topics)
	if log.Removed {
		return
	} else if !b.registrations.isAddressRegistered(log.Address) {
		return
	}
	b.logPool.addLog(log)
}

func (b *broadcaster) onNewHeads() {
	for {
		// We only care about the most recent head
		x := b.newHeads.RetrieveLatestAndClear()
		if x == nil {
			// This should never happen
			break
		}
		head, ok := x.(models.Head)
		if !ok {
			logger.Errorf("expected `models.Head`, got %T", x)
			continue
		}
		logger.Tracef("///////////////////////// onNewHeads %v, %v", head.Number, head.Hash)
		b.latestHead = &head
	}

	logs := b.logPool.getLogsToSend(b.latestHead, b.registrations.highestNumConfirmations)
	b.registrations.sendLogs(logs, b.orm, b.latestHead)
}

func (b *broadcaster) onAddSubscribers() (needsResubscribe bool) {
	for {
		x := b.addSubscriber.Retrieve()
		if x == nil {
			break
		}
		reg, ok := x.(registration)
		if !ok {
			logger.Errorf("expected `registration`, got %T", x)
			continue
		}
		logger.Debugf("LogBroadcaster: Subscribing listener with %v required block confirmations", reg.opts.NumConfirmations)
		needsResub := b.registrations.addSubscriber(reg)
		if needsResub {
			needsResubscribe = true
		}
	}
	return
}

func (b *broadcaster) onRmSubscribers() (needsResubscribe bool) {
	for {
		x := b.rmSubscriber.Retrieve()
		if x == nil {
			break
		}
		reg, ok := x.(registration)
		if !ok {
			logger.Errorf("expected `registration`, got %T", x)
			continue
		}
		logger.Debugf("LogBroadcaster: Unsubscribing listener with %v required block confirmations", reg.opts.NumConfirmations)
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
