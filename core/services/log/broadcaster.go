package log

import (
	"context"
	"sync"
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
		Register(listener Listener, opts ListenerOpts) (connected bool, unsubscribe func())
	}

	broadcaster struct {
		orm         ORM
		ethClient   eth.Client
		config      Config
		connected   *abool.AtomicBool
		latestBlock uint64

		ethSubscriber *ethSubscriber
		registrations registrations
		logPool       logPool

		addSubscriber *utils.Mailbox
		rmSubscriber  *utils.Mailbox
		newHeads      *utils.Mailbox

		logsByHeight map[common.Address]map[int64][]types.Log
		logsMu       sync.Mutex
		seenChains   []*models.Head

		utils.StartStopOnce
		utils.DependentAwaiter
		chStop chan struct{}
		chDone chan struct{}
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
	return &broadcaster{
		orm:              orm,
		ethClient:        ethClient,
		config:           config,
		connected:        abool.New(),
		ethSubscriber:    newEthSubscriber(ethClient, config),
		registrations:    newRegistrations(),
		logPool:          newLogPool(),
		addSubscriber:    utils.NewMailbox(0),
		rmSubscriber:     utils.NewMailbox(0),
		newHeads:         utils.NewMailbox(1),
		logsByHeight:     make(map[common.Address]map[int64][]types.Log),
		DependentAwaiter: utils.NewDependentAwaiter(),
		chStop:           make(chan struct{}),
		chDone:           make(chan struct{}),
		seenChains:       make([]*models.Head, 0),
	}
}

func (b *broadcaster) Start() error {
	return b.StartOnce("Log broadcaster", func() error {
		go b.awaitInitialSubscribers()
		return nil
	})
}

func (b *broadcaster) Stop() error {
	return b.StopOnce("Log broadcaster", func() error {
		close(b.chStop)
		<-b.chDone
		return b.ethSubscriber.Stop()
	})
}

func (b *broadcaster) awaitInitialSubscribers() {
	for {
		select {
		case <-b.addSubscriber.Notify():
			b.onAddSubscribers()

		case <-b.rmSubscriber.Notify():
			b.onRmSubscribers()

		case <-b.DependentAwaiter.AwaitDependents():
			go b.startResubscribeLoop()
			return

		case <-b.chStop:
			close(b.chDone)
			return
		}
	}
}

func (b *broadcaster) Register(listener Listener, opts ListenerOpts) (connected bool, unsubscribe func()) {
	b.addSubscriber.Deliver(registration{listener, opts})
	return b.IsConnected(), func() {
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
	defer close(b.chDone)

	var subscription managedSubscription = newNoopSubscription()
	defer func() { subscription.Unsubscribe() }()

	var chRawLogs chan types.Log
	for {
		addresses, topics := b.registrations.addressesAndTopics()

		newSubscription, abort := b.ethSubscriber.createSubscription(addresses, topics)
		if abort {
			return
		}

		chBackfilledLogs, abort := b.ethSubscriber.backfillLogs(addresses, topics)
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
				return true, nil
			}

		case <-b.chStop:
			return false, nil
		}
	}
}

func (b *broadcaster) onNewLog(log types.Log) {
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

		b.updateSeenChains(&head)
	}
	latestBlockNumber := b.seenChains[len(b.seenChains)-1].Number

	confirmationDepths := b.registrations.getDistinctConfirmationDepths()

	logs := b.logPool.getLogsToSend(b.seenChains, confirmationDepths)

	for _, log := range logs {
		logger.Infof("log: %v", log.BlockNumber)
		b.registrations.sendLog(log, b.orm, latestBlockNumber)
	}

	logger.Warnf("///////////////////////// heads")
}

func (b *broadcaster) updateSeenChains(head *models.Head) {
	if len(b.seenChains) == 0 {
		b.seenChains = append(b.seenChains, head)
		return
	}

	lastSeen := b.seenChains[len(b.seenChains)-1]
	if head.Parent != nil && head.Parent.Hash == lastSeen.Hash {
		// just a continuation of the previous chain so replace it
		b.seenChains[len(b.seenChains)-1] = head
		return
	}

	b.seenChains = append(b.seenChains, head)

	if b.seenChains[0].Number < head.Number-10 {
		b.seenChains = b.seenChains[1:]
	}
}

func (b *broadcaster) currentChain() *models.Head {
	if len(b.seenChains) == 0 {
		return nil
	}
	return b.seenChains[len(b.seenChains)-1]
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
