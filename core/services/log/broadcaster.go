package log

import (
	"context"
	"math/big"
	"reflect"
	"time"

	"github.com/smartcontractkit/chainlink/core/store/models"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/eth"
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
		Register(contract AbigenContract, listener Listener) (connected bool)
		Unregister(contract AbigenContract, listener Listener)
	}
}

// The Broadcast type wraps a models.Log but provides additional functionality
// for determining whether or not the log has been consumed and for marking
// the log as consumed
type Broadcast interface {
	DecodedLog() interface{}
	RawLog() types.Log
	SetDecodedLog(interface{})
	WasAlreadyConsumed() (bool, error)
	MarkConsumed() error
}

type broadcast struct {
	orm        ormInterface
	decodedLog interface{}
	rawLog     types.Log
	jobID      models.JobID
	jobIDV2    int32
	isV2       bool
}

func (lb *broadcast) DecodedLog() interface{} {
	return lb.decodedLog
}

func (lb *broadcast) RawLog() types.Log {
	return lb.rawLog
}

func (lb *broadcast) SetDecodedLog(newLog interface{}) {
	lb.decodedLog = newLog
}

func (lb *broadcast) WasAlreadyConsumed() (bool, error) {
	rawLog := lb.rawLog
	if lb.isV2 {
		return lb.orm.HasConsumedLogV2(rawLog.BlockHash, rawLog.Index, lb.jobIDV2)
	}
	return lb.orm.HasConsumedLog(rawLog.BlockHash, rawLog.Index, lb.jobID)
}

func (lb *broadcast) MarkConsumed() error {
	rawLog := lb.rawLog
	if lb.isV2 {
		return lb.orm.MarkLogConsumedV2(rawLog.BlockHash, rawLog.Index, lb.jobIDV2, rawLog.BlockNumber)
	}
	return lb.orm.MarkLogConsumed(rawLog.BlockHash, rawLog.Index, lb.jobID, rawLog.BlockNumber)
}

// A `registration` represents a Listener's subscription to the logs of a
// particular contract.
type registration struct {
	address  common.Address
	listener Listener
}

	AbigenContract interface {
		Address() common.Address
		ParseLog(log types.Log) (interface{}, error)
	}

	Config interface {
		BlockBackfillDepth() uint64
		TriggerFallbackDBPollInterval() time.Duration
	}
	go b.awaitInitialSubscribers()
	return nil
}

func (b *broadcaster) awaitInitialSubscribers() {
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

func (b *broadcaster) addresses() []common.Address {
	var addresses []common.Address
	for address := range b.listeners {
		addresses = append(addresses, address)
	}
	return addresses
}

func (b *broadcaster) Stop() error {
	if !b.OkayToStop() {
		return errors.New("Broadcaster is already stopped")
	}
	close(b.chStop)
	<-b.chDone
	return nil
}

func (b *broadcaster) Register(address common.Address, listener Listener) (connected bool) {
	select {
	case b.chAddListener <- registration{address, listener}:
	case <-b.chStop:
	}
	return b.connected.IsSet()
}

func (b *broadcaster) Unregister(address common.Address, listener Listener) {
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

func (b *broadcaster) backfillLogs() (chBackfilledLogs chan types.Log, abort bool) {
	if len(b.listeners) == 0 {
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
			logger.Errorw("Broadcaster backfill: could not fetch latest block header", "error", err)
			return true
		} else if latestBlock == nil {
			logger.Warn("got nil block header")
			return true
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

		logs, err := b.ethClient.FilterLogs(ctx, q)
		if err != nil {
			logger.Errorw("Broadcaster backfill: could not fetch logs", "error", err)
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

func (b *broadcaster) notifyConnect() {
	b.connected.Set()
	for _, listeners := range b.listeners {
		for listener := range listeners {
			listener.OnConnect()
		}
	}
}

func (b *broadcaster) notifyDisconnect() {
	b.connected.UnSet()
	for _, listeners := range b.listeners {
		for listener := range listeners {
			listener.OnDisconnect()
		}
	}
}

func (b *broadcaster) process(subscription managedSubscription, chRawLogs <-chan types.Log) (shouldResubscribe bool, _ error) {
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

func (b *broadcaster) onRawLog(rawLog types.Log) {
	for listener := range b.listeners[rawLog.Address] {
		// Ignore duplicate logs sent back due to reorgs
		if rawLog.Removed {
			continue
		}

		// Deep copy the log so that subscribers aren't sharing any state
		rawLogCopy := copyLog(rawLog)
		lb := &broadcast{
			rawLog:  rawLogCopy,
			orm:     b.orm,
			jobID:   listener.JobID(),
			jobIDV2: listener.JobIDV2(),
			isV2:    listener.IsV2Job(),
		}
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

func (b *broadcaster) onAddListener(r registration) (needsResubscribe bool) {
	_, knownAddress := b.listeners[r.address]
	if !knownAddress {
		b.listeners[r.address] = make(map[Listener]struct{})
	}
	if _, exists := b.listeners[r.address][r.listener]; exists {
		panic("registration already exists")
	}
	b.listeners[r.address][r.listener] = struct{}{}

	// Recreate the subscription with the new contract address
	return !knownAddress
}

func (b *broadcaster) onRemoveListener(r registration) (needsResubscribe bool) {
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
func (b *broadcaster) createSubscription() (sub managedSubscription, abort bool) {
	if len(b.listeners) == 0 {
		return newNoopSubscription(), false
	}

	ctx, cancel := utils.ContextFromChan(b.chStop)
	defer cancel()

	utils.RetryWithBackoff(ctx, func() (retry bool) {
		filterQuery := ethereum.FilterQuery{
			Addresses: b.addresses(),
		}
		chRawLogs := make(chan types.Log)

		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()
		innerSub, err := b.ethClient.SubscribeFilterLogs(ctx, filterQuery, chRawLogs)
		if err != nil {
			logger.Errorw("Broadcaster could not create subscription to Ethereum node", "error", err)
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

func (b *broadcaster) Register(contract AbigenContract, listener Listener) (connected bool) {
	b.subscriber.NotifyAddContract(contract.Address())
	b.relayer.NotifyAddListener(contract, listener)
	return b.subscriber.IsConnected()
}

func (b *broadcaster) Unregister(contract AbigenContract, listener Listener) {
	b.subscriber.NotifyRemoveContract(contract.Address())
	b.relayer.NotifyRemoveListener(contract, listener)
}

var _ Listener = (*decodingLogListener)(nil)

// NewDecodingLogListener creates a new decodingLogListener
func NewDecodingLogListener(codec eth.ContractCodec, nativeLogTypes map[common.Hash]interface{}, innerListener Listener) Listener {
	logTypes := make(map[common.Hash]reflect.Type)
	for eventID, logStruct := range nativeLogTypes {
		logTypes[eventID] = reflect.TypeOf(logStruct)
	}

	return &decodingLogListener{
		logTypes: logTypes,
		codec:    codec,
		Listener: innerListener,
	}
}

func (l *decodingLogListener) HandleLog(lb Broadcast, err error) {
	if err != nil {
		l.Listener.HandleLog(&broadcast{}, err)
		return
	}

	rawLog := lb.RawLog()

	if len(rawLog.Topics) == 0 {
		return
	}
	eventID := rawLog.Topics[0]
	logType, exists := l.logTypes[eventID]
	if !exists {
		// If a particular log type hasn't been registered with the decoder, we simply ignore it.
		return
	}

	var decodedLog interface{}
	if logType.Kind() == reflect.Ptr {
		decodedLog = reflect.New(logType.Elem()).Interface()
	} else {
		decodedLog = reflect.New(logType).Interface()
	}

	// Insert the raw log into the ".Log" field
	logStructV := reflect.ValueOf(decodedLog).Elem()
	logStructV.FieldByName("Log").Set(reflect.ValueOf(rawLog))

	// Decode the raw log into the struct
	event, err := l.codec.ABI().EventByID(eventID)
	if err != nil {
		l.Listener.HandleLog(nil, err)
		return
	}
	err = l.codec.UnpackLog(decodedLog, event.RawName, rawLog)
	if err != nil {
		l.Listener.HandleLog(nil, err)
		return
	}

	lb.SetDecodedLog(decodedLog)
	l.Listener.HandleLog(lb, nil)
}
