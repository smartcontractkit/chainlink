package log

import (
	"context"
	"math/big"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum"
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

		registrations  map[common.Address]map[common.Hash]map[Listener]struct{} // contractAddress => logTopic => Listener
		decoders       map[common.Address]AbigenContract
		addSubscriber  *utils.Mailbox
		rmSubscriber   *utils.Mailbox
		newHeads       *utils.Mailbox
		logs           map[common.Address][]types.Log
		logsMu         sync.Mutex
		canonicalChain map[common.Hash]struct{}

		utils.StartStopOnce
		utils.DependentAwaiter
		chStop chan struct{}
		chDone chan struct{}
	}

	// The Listener responds to log events through HandleLog, and contains setup/tear-down
	// callbacks in the On* functions.
	Listener interface {
		OnConnect()
		OnDisconnect()
		HandleLog(b Broadcast)
		JobID() models.JobID
		JobIDV2() int32
		IsV2Job() bool
	}

	Config interface {
		BlockBackfillDepth() uint64
		TriggerFallbackDBPollInterval() time.Duration
	}

	ListenerOpts struct {
		Contract AbigenContract
		Logs     []generated.AbigenLog
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
		registrations:    make(map[common.Address]map[common.Hash]map[Listener]struct{}),
		decoders:         make(map[common.Address]AbigenContract),
		addSubscriber:    utils.NewMailbox(0),
		rmSubscriber:     utils.NewMailbox(0),
		newHeads:         utils.NewMailbox(1),
		logs:             make(map[common.Address][]types.Log),
		DependentAwaiter: utils.NewDependentAwaiter(),
		chStop:           make(chan struct{}),
		chDone:           make(chan struct{}),
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
		return nil
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
	if len(opts.Logs) < 1 {
		logger.Fatal("Must supply at least 1 Log to Register")
	}
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
		newSubscription, abort := b.createSubscription()
		if abort {
			return
		}

		chBackfilledLogs, abort := b.backfillLogs()
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

		shouldResubscribe, err := b.process(chRawLogs, subscription.Err())
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

func (b *broadcaster) backfillLogs() (chBackfilledLogs chan types.Log, abort bool) {
	if len(b.registrations) == 0 {
		ch := make(chan types.Log)
		close(ch)
		return ch, false
	}

	ctx, cancel := utils.ContextFromChan(b.chStop)
	defer cancel()

	utils.RetryWithBackoff(ctx, func() (retry bool) {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		latestBlock, err := b.ethClient.HeaderByNumber(ctx, nil)
		if err != nil {
			logger.Errorw("Log subscriber backfill: could not fetch latest block header", "error", err)
			return true
		} else if latestBlock == nil {
			logger.Warn("got nil block header")
			return true
		}
		currentHeight := uint64(latestBlock.Number)

		// Backfill from `backfillDepth` blocks ago.  It'b up to the subscribers to
		// filter out logs they've already dealt with.
		fromBlock := currentHeight - b.config.BlockBackfillDepth()
		if fromBlock > currentHeight {
			fromBlock = 0 // Overflow protection
		}

		addresses, topics := b.registeredAddressesAndTopics()

		q := ethereum.FilterQuery{
			FromBlock: big.NewInt(int64(fromBlock)),
			Addresses: addresses,
			Topics:    [][]common.Hash{topics},
		}

		logs, err := b.ethClient.FilterLogs(ctx, q)
		if err != nil {
			logger.Errorw("Log subscriber backfill: could not fetch logs", "error", err)
			return true
		}

		chBackfilledLogs = make(chan types.Log)
		go b.deliverBackfilledLogs(logs, chBackfilledLogs)

		return false
	})
	select {
	case <-b.chStop:
		abort = true
	default:
		abort = false
	}
	return
}

func (b *broadcaster) deliverBackfilledLogs(logs []types.Log, chBackfilledLogs chan<- types.Log) {
	defer close(chBackfilledLogs)
	for _, log := range logs {
		select {
		case chBackfilledLogs <- log:
		case <-b.chStop:
			return
		}
	}
}

func (b *broadcaster) process(chRawLogs <-chan types.Log, chErr <-chan error) (shouldResubscribe bool, _ error) {
	// We debounce requests to subscribe and unsubscribe to avoid making too many
	// RPC calls to the Ethereum node, particularly on startup.
	var needsResubscribe bool
	debounceResubscribe := time.NewTicker(1 * time.Second)
	defer debounceResubscribe.Stop()

	for {
		select {
		case rawLog := <-chRawLogs:
			b.processLog(rawLog)

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

func (b *broadcaster) processLog(log types.Log) {
	if log.Removed {
		return
	} else if _, exists := b.registrations[log.Address]; !exists {
		return
	}

	b.logsMu.Lock()
	defer b.logsMu.Unlock()
	b.logs[log.Address] = append(b.logs[log.Address], log)
}

func (b *broadcaster) onNewHeads() {
	for {
		x := b.newHeads.Retrieve()
		if x == nil {
			break
		}
		head, ok := x.(models.Head)
		if !ok {
			logger.Errorf("expected `models.Head`, got %T", x)
			continue
		}
		b.latestBlock = uint64(head.Number)

		b.canonicalChain = make(map[common.Hash]struct{})
		for h := &head; h.Parent != nil; h = h.Parent {
			b.canonicalChain[h.Hash] = struct{}{}
		}
	}

	// Now that we're caught up, broadcast all pending logs
	b.broadcastPendingLogs()
}

func (b *broadcaster) broadcastPendingLogs() {
	b.logsMu.Lock()
	defer b.logsMu.Unlock()

	for _, logs := range b.logs {
		for _, log := range logs {
			// Skip logs that have been reorged away
			if _, exists := b.canonicalChain[log.BlockHash]; !exists {
				continue
			}

			var wg sync.WaitGroup
			for listener := range b.registrations[log.Address][log.Topics[0]] {
				listener := listener

				logCopy := copyLog(log)
				decodedLog, err := b.decoders[log.Address].ParseLog(logCopy)
				if err != nil {
					logger.Errorw("could not parse contract log", "error", err)
					continue
				}

				wg.Add(1)
				go func() {
					defer wg.Done()
					listener.HandleLog(&broadcast{
						orm:        b.orm,
						rawLog:     logCopy,
						decodedLog: decodedLog,
						jobID:      listener.JobID(),
						jobIDV2:    listener.JobIDV2(),
						isV2:       listener.IsV2Job(),
					})
				}()
			}
			wg.Wait()
		}
	}
	b.logs = make(map[common.Address][]types.Log)
}

func copyLog(l types.Log) types.Log {
	var cpy types.Log
	cpy.Address = l.Address
	if l.Topics != nil {
		cpy.Topics = make([]common.Hash, len(l.Topics))
		copy(cpy.Topics, l.Topics)
	}
	if l.Data != nil {
		cpy.Data = make([]byte, len(l.Data))
		copy(cpy.Data, l.Data)
	}
	cpy.BlockNumber = l.BlockNumber
	cpy.TxHash = l.TxHash
	cpy.TxIndex = l.TxIndex
	cpy.BlockHash = l.BlockHash
	cpy.Index = l.Index
	cpy.Removed = l.Removed
	return cpy
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

		addr := reg.opts.Contract.Address()
		b.decoders[addr] = reg.opts.Contract

		if _, exists := b.registrations[addr]; !exists {
			b.registrations[addr] = make(map[common.Hash]map[Listener]struct{})
		}

		topics := make([]common.Hash, len(reg.opts.Logs))
		for i, log := range reg.opts.Logs {
			topic := log.Topic()
			topics[i] = topic

			if _, exists := b.registrations[addr][topic]; !exists {
				b.registrations[addr][topic] = make(map[Listener]struct{})
				needsResubscribe = true
			}
			b.registrations[addr][topic][reg.listener] = struct{}{}
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

		addr := reg.opts.Contract.Address()

		if _, exists := b.registrations[addr]; !exists {
			continue
		}
		for _, logType := range reg.opts.Logs {
			topic := logType.Topic()

			if _, exists := b.registrations[addr][topic]; !exists {
				continue
			}

			delete(b.registrations[addr][topic], reg.listener)

			if len(b.registrations[addr][topic]) == 0 {
				needsResubscribe = true
				delete(b.registrations[addr], topic)
			}
			if len(b.registrations[addr]) == 0 {
				delete(b.registrations, addr)
			}
		}
	}
	return
}

// createSubscription creates a new log subscription starting at the current block.  If previous logs
// are needed, they must be obtained through backfilling, as subscriptions can only be started from
// the current head.
func (b *broadcaster) createSubscription() (sub managedSubscription, abort bool) {
	if len(b.registrations) == 0 {
		return newNoopSubscription(), false
	}

	ctx, cancel := utils.ContextFromChan(b.chStop)
	defer cancel()

	utils.RetryWithBackoff(ctx, func() (retry bool) {
		addresses, topics := b.registeredAddressesAndTopics()

		filterQuery := ethereum.FilterQuery{
			Addresses: addresses,
			Topics:    [][]common.Hash{topics},
		}
		chRawLogs := make(chan types.Log)

		ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
		defer cancel()

		innerSub, err := b.ethClient.SubscribeFilterLogs(ctx, filterQuery, chRawLogs)
		if err != nil {
			logger.Errorw("Log subscriber could not create subscription to Ethereum node", "error", err)
			return true
		}

		sub = managedSubscriptionImpl{
			subscription: innerSub,
			chRawLogs:    chRawLogs,
		}
		return false
	})
	select {
	case <-b.chStop:
		abort = true
	default:
		abort = false
	}
	return
}

func (b *broadcaster) registeredAddressesAndTopics() ([]common.Address, []common.Hash) {
	var addresses []common.Address
	var topics []common.Hash
	for addr := range b.registrations {
		addresses = append(addresses, addr)
		for topic := range b.registrations[addr] {
			topics = append(topics, topic)
		}
	}
	return addresses, topics
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

// A managedSubscription acts as wrapper for the Subscription. Specifically, the
// managedSubscription closes the log channel as soon as the unsubscribe request is made
type managedSubscription interface {
	Err() <-chan error
	Logs() chan types.Log
	Unsubscribe()
}

type managedSubscriptionImpl struct {
	subscription ethereum.Subscription
	chRawLogs    chan types.Log
}

func (sub managedSubscriptionImpl) Err() <-chan error {
	return sub.subscription.Err()
}

func (sub managedSubscriptionImpl) Logs() chan types.Log {
	return sub.chRawLogs
}

func (sub managedSubscriptionImpl) Unsubscribe() {
	sub.subscription.Unsubscribe()
	close(sub.chRawLogs)
}

type noopSubscription struct {
	chRawLogs chan types.Log
}

func newNoopSubscription() noopSubscription {
	return noopSubscription{make(chan types.Log)}
}

func (b noopSubscription) Err() <-chan error    { return nil }
func (b noopSubscription) Logs() chan types.Log { return b.chRawLogs }
func (b noopSubscription) Unsubscribe()         { close(b.chRawLogs) }

// ListenerJobID returns the appropriate job ID for a listener
func ListenerJobID(listener Listener) interface{} {
	if listener.IsV2Job() {
		return listener.JobIDV2()
	}
	return listener.JobID()
}
