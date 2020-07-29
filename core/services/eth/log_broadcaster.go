package eth

import (
	"context"
	"github.com/pkg/errors"
	"math/big"
	"reflect"
	"time"

	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

//go:generate mockery --name LogBroadcaster --output ../../internal/mocks/ --case=underscore
//go:generate mockery --name LogListener --output ../../internal/mocks/ --case=underscore
//go:generate mockery --name LogBroadcast --output ../../internal/mocks/ --case=underscore

// The LogBroadcaster manages log subscription requests for the Chainlink node.  Instead
// of creating a new websocket subscription for each request, it multiplexes all subscriptions
// to all of the relevant contracts over a single connection and forwards the logs to the
// relevant subscribers.
type LogBroadcaster interface {
	utils.DependentAwaiter
	Start() error
	Register(address common.Address, listener LogListener) (connected bool)
	Unregister(address common.Address, listener LogListener)
	Stop()
}

// The LogListener responds to log events through HandleLog, and contains setup/tear-down
// callbacks in the On* functions.
type LogListener interface {
	OnConnect()
	OnDisconnect()
	HandleLog(lb LogBroadcast, err error)
	JobID() *models.ID
}

type ormInterface interface {
	HasConsumedLog(blockHash common.Hash, logIndex uint, jobID *models.ID) (bool, error)
	MarkLogConsumed(blockHash common.Hash, logIndex uint, jobID *models.ID) error
}

type logBroadcaster struct {
	ethClient     Client
	orm           ormInterface
	backfillDepth uint64
	connected     bool
	started       bool

	listeners        map[common.Address]map[LogListener]struct{}
	chAddListener    chan registration
	chRemoveListener chan registration

	utils.DependentAwaiter
	chStop chan struct{}
	chDone chan struct{}
}

// NewLogBroadcaster creates a new instance of the logBroadcaster
func NewLogBroadcaster(ethClient Client, orm ormInterface, backfillDepth uint64) LogBroadcaster {
	return &logBroadcaster{
		ethClient:        ethClient,
		orm:              orm,
		backfillDepth:    backfillDepth,
		listeners:        make(map[common.Address]map[LogListener]struct{}),
		chAddListener:    make(chan registration),
		chRemoveListener: make(chan registration),
		chStop:           make(chan struct{}),
		chDone:           make(chan struct{}),
		DependentAwaiter: utils.NewDependentAwaiter(),
	}
}

// The LogBroadcast type wraps an models.Log but provides additional functionality
// for determining whether or not the log has been consumed and for marking
// the log as consumed
type LogBroadcast interface {
	Log() Log
	UpdateLog(Log)
	WasAlreadyConsumed() (bool, error)
	MarkConsumed() error
}

type Log interface {
	RawLog() types.Log
}

type GethRawLog struct {
	types.Log
}

func (rl GethRawLog) RawLog() types.Log { return rl.Log }

type logBroadcast struct {
	orm        ormInterface
	log        Log
	consumerID *models.ID
}

func (lb *logBroadcast) Log() Log {
	return lb.log
}

func (lb *logBroadcast) UpdateLog(newLog Log) {
	lb.log = newLog
}

func (lb *logBroadcast) WasAlreadyConsumed() (bool, error) {
	rawLog := lb.log.RawLog()
	return lb.orm.HasConsumedLog(rawLog.BlockHash, rawLog.Index, lb.consumerID)
}

func (lb *logBroadcast) MarkConsumed() error {
	rawLog := lb.log.RawLog()
	return lb.orm.MarkLogConsumed(rawLog.BlockHash, rawLog.Index, lb.consumerID)
}

// A `registration` represents a LogListener's subscription to the logs of a
// particular contract.
type registration struct {
	address  common.Address
	listener LogListener
}

func (b *logBroadcaster) Start() error {
	go b.awaitInitialSubscribers()
	b.started = true
	return nil
}

func (b *logBroadcaster) awaitInitialSubscribers() {
	for {
		select {
		case r := <-b.chAddListener:
			b.onAddListener(r)

		case <-b.DependentAwaiter.AwaitDependents():
			go b.startResubscribeLoop()
			return

		case <-b.chStop:
			close(b.chDone)
			return
		}
	}
}

func (b *logBroadcaster) addresses() []common.Address {
	var addresses []common.Address
	for address := range b.listeners {
		addresses = append(addresses, address)
	}
	return addresses
}

func (b *logBroadcaster) Stop() {
	close(b.chStop)
	if b.started {
		<-b.chDone
		b.started = false

	}
}

func (b *logBroadcaster) Register(address common.Address, listener LogListener) (connected bool) {
	select {
	case b.chAddListener <- registration{address, listener}:
	case <-b.chStop:
	}
	return b.connected
}

func (b *logBroadcaster) Unregister(address common.Address, listener LogListener) {
	select {
	case b.chRemoveListener <- registration{address, listener}:
	case <-b.chStop:
	}
}

// The subscription is closed in two cases:
//   - intentionally, when the set of contracts we're listening to changes
//   - on a connection error
//
// This method recreates the subscription in both cases.  In the event of a connection
// error, it attempts to reconnect.  Any time there's a change in connection state, it
// notifies its subscribers.
func (b *logBroadcaster) startResubscribeLoop() {
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
		//     remaining logs from last subscription <= backfilled logs <= logs from new subscription
		// There will be duplicated logs in this channel.  It is the responsibility of subscribers
		// to account for this using the helpers on the LogBroadcast type.
		chRawLogs = appendLogChannel(chRawLogs, chBackfilledLogs)
		chRawLogs = appendLogChannel(chRawLogs, newSubscription.Logs())
		subscription.Unsubscribe()
		subscription = newSubscription

		b.notifyConnect()
		shouldResubscribe, err := b.process(subscription, chRawLogs)
		if err != nil {
			logger.Error(err)
			b.notifyDisconnect()
			continue
		} else if !shouldResubscribe {
			b.notifyDisconnect()
			return
		}
	}
}

func (b *logBroadcaster) backfillLogs() (chBackfilledLogs chan types.Log, abort bool) {
	if len(b.listeners) == 0 {
		ch := make(chan types.Log)
		close(ch)
		return ch, false
	}

	abort = utils.RetryWithBackoff(b.chStop, "backfilling logs", func() error {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		latestBlock, err := b.ethClient.HeaderByNumber(ctx, nil)
		if err != nil {
			return err
		} else if latestBlock == nil {
			logger.Warn("got nil block header")
			return errors.New("got nil block header")
		}
		currentHeight := uint64(latestBlock.Number.Int64())

		// Backfill from `backfillDepth` blocks ago.  It's up to the subscribers to
		// filter out logs they've already dealt with.
		fromBlock := currentHeight - b.backfillDepth
		if fromBlock > currentHeight {
			fromBlock = 0 // Overflow protection
		}

		q := ethereum.FilterQuery{
			FromBlock: big.NewInt(int64(fromBlock)),
			Addresses: b.addresses(),
		}

		logs, err := b.ethClient.FilterLogs(ctx, q)
		if err != nil {
			return err
		}

		chBackfilledLogs = make(chan types.Log)
		go b.deliverBackfilledLogs(logs, chBackfilledLogs)
		return nil

	})
	return
}

func (b *logBroadcaster) deliverBackfilledLogs(logs []types.Log, chBackfilledLogs chan<- types.Log) {
	defer close(chBackfilledLogs)
	for _, log := range logs {
		select {
		case chBackfilledLogs <- log:
		case <-b.chStop:
			return
		}
	}
}

func (b *logBroadcaster) notifyConnect() {
	b.connected = true
	for _, listeners := range b.listeners {
		for listener := range listeners {
			listener.OnConnect()
		}
	}
}

func (b *logBroadcaster) notifyDisconnect() {
	b.connected = false
	for _, listeners := range b.listeners {
		for listener := range listeners {
			listener.OnDisconnect()
		}
	}
}

func (b *logBroadcaster) process(subscription managedSubscription, chRawLogs <-chan types.Log) (shouldResubscribe bool, _ error) {
	// We debounce requests to subscribe and unsubscribe to avoid making too many
	// RPC calls to the Ethereum node, particularly on startup.
	var needsResubscribe bool
	debounceResubscribe := time.NewTicker(1 * time.Second)
	defer debounceResubscribe.Stop()

	for {
		select {
		case rawLog := <-chRawLogs:
			b.onRawLog(rawLog)

		case r := <-b.chAddListener:
			needsResubscribe = b.onAddListener(r) || needsResubscribe

		case r := <-b.chRemoveListener:
			needsResubscribe = b.onRemoveListener(r) || needsResubscribe

		case <-debounceResubscribe.C:
			if needsResubscribe {
				return true, nil
			}

		case err := <-subscription.Err():
			return true, err

		case <-b.chStop:
			return false, nil
		}
	}
}

func (b *logBroadcaster) onRawLog(rawLog types.Log) {
	for listener := range b.listeners[rawLog.Address] {
		// Ignore duplicate logs sent back due to reorgs
		if rawLog.Removed {
			continue
		}

		// Deep copy the log so that subscribers aren't sharing any state
		rawLogCopy := copyLog(rawLog)
		lb := &logBroadcast{log: GethRawLog{rawLogCopy}, orm: b.orm, consumerID: listener.JobID()}
		listener.HandleLog(lb, nil)
	}
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

func (b *logBroadcaster) onAddListener(r registration) (needsResubscribe bool) {
	_, knownAddress := b.listeners[r.address]
	if !knownAddress {
		b.listeners[r.address] = make(map[LogListener]struct{})
	}
	if _, exists := b.listeners[r.address][r.listener]; exists {
		panic("registration already exists")
	}
	b.listeners[r.address][r.listener] = struct{}{}

	// Recreate the subscription with the new contract address
	return !knownAddress
}

func (b *logBroadcaster) onRemoveListener(r registration) (needsResubscribe bool) {
	r.listener.OnDisconnect()
	delete(b.listeners[r.address], r.listener)
	if len(b.listeners[r.address]) == 0 {
		delete(b.listeners, r.address)
		// Recreate the subscription without this contract address
		return true
	}
	return false
}

// createSubscription creates a new log subscription starting at the current block.  If previous logs
// are needed, they must be obtained through backfilling, as subscriptions can only be started from
// the current head.
func (b *logBroadcaster) createSubscription() (sub managedSubscription, abort bool) {
	if len(b.listeners) == 0 {
		return newNoopSubscription(), false
	}

	abort = utils.RetryWithBackoff(b.chStop, "creating subscription to Ethereum node", func() error {
		filterQuery := ethereum.FilterQuery{
			Addresses: b.addresses(),
		}
		chRawLogs := make(chan types.Log)

		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()
		innerSub, err := b.ethClient.SubscribeFilterLogs(ctx, filterQuery, chRawLogs)
		if err != nil {
			return err
		}

		sub = managedSubscriptionImpl{
			subscription: innerSub,
			chRawLogs:    chRawLogs,
		}
		return nil
	})
	return
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

func (s noopSubscription) Err() <-chan error    { return nil }
func (s noopSubscription) Logs() chan types.Log { return s.chRawLogs }
func (s noopSubscription) Unsubscribe()         { close(s.chRawLogs) }

// DecodingLogListener receives raw logs from the LogBroadcaster and decodes them into
// Go structs using the provided ContractCodec (a simple wrapper around a go-ethereum
// ABI type).
type decodingLogListener struct {
	logTypes map[common.Hash]reflect.Type
	codec    ContractCodec
	LogListener
}

var _ LogListener = (*decodingLogListener)(nil)

// NewDecodingLogListener creates a new decodingLogListener
func NewDecodingLogListener(codec ContractCodec, nativeLogTypes map[common.Hash]Log, innerListener LogListener) LogListener {
	logTypes := make(map[common.Hash]reflect.Type)
	for eventID, logStruct := range nativeLogTypes {
		logTypes[eventID] = reflect.TypeOf(logStruct)
	}

	return &decodingLogListener{
		logTypes:    logTypes,
		codec:       codec,
		LogListener: innerListener,
	}
}

func (l *decodingLogListener) HandleLog(lb LogBroadcast, err error) {
	if err != nil {
		l.LogListener.HandleLog(&logBroadcast{}, err)
		return
	}

	rawLog, is := lb.Log().(GethRawLog)
	if !is {
		panic("DecodingLogListener expects to receive a logBroadcast with a GethRawLog")
	}

	if len(rawLog.Topics) == 0 {
		return
	}
	eventID := rawLog.Topics[0]
	logType, exists := l.logTypes[eventID]
	if !exists {
		// If a particular log type hasn't been registered with the decoder, we simply ignore it.
		return
	}

	var decodedLog Log
	var ok bool
	if logType.Kind() == reflect.Ptr {
		decodedLog, ok = reflect.New(logType.Elem()).Interface().(Log)
	} else {
		decodedLog, ok = reflect.New(logType).Interface().(Log)
	}
	if !ok {
		panic("DecodingLogListener expects a Rawlog logType")
	}

	// Insert the raw log into the ".Log" field
	logStructV := reflect.ValueOf(decodedLog).Elem()
	logStructV.FieldByName("GethRawLog").Set(reflect.ValueOf(rawLog))

	// Decode the raw log into the struct
	event, err := l.codec.ABI().EventByID(eventID)
	if err != nil {
		l.LogListener.HandleLog(nil, err)
		return
	}
	err = l.codec.UnpackLog(decodedLog, event.RawName, rawLog.Log)
	if err != nil {
		l.LogListener.HandleLog(nil, err)
		return
	}

	lb.UpdateLog(decodedLog)
	l.LogListener.HandleLog(lb, nil)
}

func appendLogChannel(ch1, ch2 <-chan types.Log) chan types.Log {
	if ch1 == nil && ch2 == nil {
		return nil
	}

	chCombined := make(chan types.Log)

	go func() {
		defer close(chCombined)
		if ch1 != nil {
			for rawLog := range ch1 {
				chCombined <- rawLog
			}
		}
		if ch2 != nil {
			for rawLog := range ch2 {
				chCombined <- rawLog
			}
		}
	}()

	return chCombined
}
