package eth

import (
	"context"
	"math/big"
	"reflect"
	"time"

	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/tevino/abool"
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

type orm interface {
	HasConsumedLog(rawLog models.RawLog, JobID *models.ID) (bool, error)
	CreateLogConsumption(lc *models.LogConsumption) error
}

type logBroadcaster struct {
	ethClient     Client
	orm           orm
	backfillDepth uint64
	connected     *abool.AtomicBool
	started       *abool.AtomicBool

	listeners        map[common.Address]map[LogListener]struct{}
	chAddListener    chan registration
	chRemoveListener chan registration

	utils.DependentAwaiter
	chStop chan struct{}
	chDone chan struct{}
}

// NewLogBroadcaster creates a new instance of the logBroadcaster
func NewLogBroadcaster(ethClient Client, orm orm, backfillDepth uint64) LogBroadcaster {
	return &logBroadcaster{
		ethClient:        ethClient,
		orm:              orm,
		backfillDepth:    backfillDepth,
		connected:        abool.New(),
		started:          abool.New(),
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
	Log() interface{}
	UpdateLog(models.RawLog)
	WasAlreadyConsumed() (bool, error)
	MarkConsumed() error
}

type logBroadcast struct {
	orm   orm
	log   models.RawLog
	jobID *models.ID
}

func (lb *logBroadcast) Log() interface{} {
	return lb.log
}

func (lb *logBroadcast) UpdateLog(newLog models.RawLog) {
	lb.log = newLog
}

func (lb *logBroadcast) WasAlreadyConsumed() (bool, error) {
	return lb.orm.HasConsumedLog(lb.log, lb.jobID)
}

func (lb *logBroadcast) MarkConsumed() error {
	lc := models.NewLogConsumption(lb.log, lb.jobID)
	return lb.orm.CreateLogConsumption(&lc)
}

type registration struct {
	address  common.Address
	listener LogListener
}

// A ManagedSubscription acts as wrapper for the Subscription. Specifically, the
// ManagedSubscription closes the log channel as soon as the unsubscribe request is made
type ManagedSubscription interface {
	Err() <-chan error
	Logs() chan models.Log
	Unsubscribe()
}

type managedSubscription struct {
	subscription Subscription
	chRawLogs    chan models.Log
}

func (sub managedSubscription) Err() <-chan error {
	return sub.subscription.Err()
}

func (sub managedSubscription) Logs() chan models.Log {
	return sub.chRawLogs
}

func (sub managedSubscription) Unsubscribe() {
	sub.subscription.Unsubscribe()
	close(sub.chRawLogs)
}

func (b *logBroadcaster) Start() error {
	go b.awaitInitialSubscribers()
	b.started.Set()
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
	if b.started.IsSet() {
		<-b.chDone
		b.started.UnSet()

	}
}

func (b *logBroadcaster) Register(address common.Address, listener LogListener) (connected bool) {
	select {
	case b.chAddListener <- registration{address, listener}:
	case <-b.chStop:
	}
	return b.connected.IsSet()
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

	var subscription ManagedSubscription = newNoopSubscription()
	defer func() { subscription.Unsubscribe() }()

	var chRawLogs chan models.Log
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

func (b *logBroadcaster) backfillLogs() (chBackfilledLogs chan models.Log, abort bool) {
	if len(b.listeners) == 0 {
		ch := make(chan models.Log)
		close(ch)
		return ch, false
	}

	abort = utils.RetryWithBackoff(b.chStop, "backfilling logs", func() error {
		latestBlock, err := b.ethClient.GetLatestBlock()
		if err != nil {
			return err
		}
		currentHeight := uint64(latestBlock.Number)

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

		logs, err := b.ethClient.GetLogs(q)
		if err != nil {
			return err
		}

		chBackfilledLogs = make(chan models.Log)
		go b.deliverBackfilledLogs(logs, chBackfilledLogs)
		return nil

	})
	return
}

func (b *logBroadcaster) deliverBackfilledLogs(logs []models.Log, chBackfilledLogs chan<- models.Log) {
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
	b.connected.Set()
	for _, listeners := range b.listeners {
		for listener := range listeners {
			listener.OnConnect()
		}
	}
}

func (b *logBroadcaster) notifyDisconnect() {
	b.connected.UnSet()
	for _, listeners := range b.listeners {
		for listener := range listeners {
			listener.OnDisconnect()
		}
	}
}

func (b *logBroadcaster) process(subscription Subscription, chRawLogs <-chan models.Log) (shouldResubscribe bool, _ error) {
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

func (b *logBroadcaster) onRawLog(rawLog models.Log) {
	for listener := range b.listeners[rawLog.Address] {
		// Ignore duplicate logs sent back due to reorgs
		if rawLog.Removed {
			continue
		}

		rawLogCopy := rawLog.Copy()
		lb := logBroadcast{b.orm, &rawLogCopy, listener.JobID()}
		listener.HandleLog(&lb, nil)
	}
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
func (b *logBroadcaster) createSubscription() (sub ManagedSubscription, abort bool) {
	if len(b.listeners) == 0 {
		return newNoopSubscription(), false
	}

	abort = utils.RetryWithBackoff(b.chStop, "creating subscription to Ethereum node", func() error {
		filterQuery := ethereum.FilterQuery{
			Addresses: b.addresses(),
		}
		chRawLogs := make(chan models.Log)

		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()
		innerSub, err := b.ethClient.SubscribeToLogs(ctx, chRawLogs, filterQuery)
		if err != nil {
			return err
		}

		sub = managedSubscription{
			subscription: innerSub,
			chRawLogs:    chRawLogs,
		}
		return nil
	})
	return
}

type noopSubscription struct {
	chRawLogs chan models.Log
}

func newNoopSubscription() noopSubscription {
	return noopSubscription{make(chan models.Log)}
}

func (s noopSubscription) Err() <-chan error     { return nil }
func (s noopSubscription) Logs() chan models.Log { return s.chRawLogs }
func (s noopSubscription) Unsubscribe()          { close(s.chRawLogs) }

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
func NewDecodingLogListener(codec ContractCodec, nativeLogTypes map[common.Hash]interface{}, innerListener LogListener) LogListener {
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

	rawLog, is := lb.Log().(*models.Log)
	if !is {
		panic("DecodingLogListener expects to receive a logBroadcast with a *models.Log")
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

	var decodedLog models.RawLog
	var ok bool
	if logType.Kind() == reflect.Ptr {
		decodedLog, ok = reflect.New(logType.Elem()).Interface().(models.RawLog)
	} else {
		decodedLog, ok = reflect.New(logType).Interface().(models.RawLog)
	}
	if !ok {
		panic("DecodingLogListener expects a Rawlog logType")
	}

	// Insert the raw log into the ".Log" field
	logStructV := reflect.ValueOf(decodedLog).Elem()
	logStructV.FieldByName("Log").Set(reflect.ValueOf(*rawLog))

	// Decode the raw log into the struct
	event, err := l.codec.ABI().EventByID(eventID)
	if err != nil {
		l.LogListener.HandleLog(nil, err)
		return
	}
	err = l.codec.UnpackLog(decodedLog, event.RawName, *rawLog)
	if err != nil {
		l.LogListener.HandleLog(nil, err)
		return
	}

	lb.UpdateLog(decodedLog)
	l.LogListener.HandleLog(lb, nil)
}

func appendLogChannel(ch1, ch2 <-chan models.Log) chan models.Log {
	if ch1 == nil && ch2 == nil {
		return nil
	}

	chCombined := make(chan models.Log)

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
