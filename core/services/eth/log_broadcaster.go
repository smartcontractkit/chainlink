package eth

import (
	"context"
	"math/big"
	"reflect"
	"time"

	"github.com/smartcontractkit/chainlink/core/eth"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/store/orm"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/jinzhu/gorm"
)

//go:generate mockery -name LogBroadcaster -output ../../internal/mocks/ -case=underscore
//go:generate mockery -name LogListener -output ../../internal/mocks/ -case=underscore

// The LogBroadcaster manages log subscription requests for the Chainlink node.  Instead
// of creating a new websocket subscription for each request, it multiplexes all subscriptions
// to all of the relevant contracts over a single connection and forwards the logs to the
// relevant subscribers.
type LogBroadcaster interface {
	Start()
	Register(address common.Address, listener LogListener) (connected bool)
	Unregister(address common.Address, listener LogListener)
	Stop()
}

type LogListener interface {
	OnConnect()
	OnDisconnect()
	HandleLog(log interface{}, err error)
}

type logBroadcaster struct {
	ethClient eth.Client
	orm       *orm.ORM
	cursor    models.LogCursor
	connected bool

	listeners        map[common.Address]map[LogListener]struct{}
	chAddListener    chan registration
	chRemoveListener chan registration

	chStop chan struct{}
	chDone chan struct{}
}

type registration struct {
	address  common.Address
	listener LogListener
}

func NewLogBroadcaster(ethClient eth.Client, orm *orm.ORM) LogBroadcaster {
	return &logBroadcaster{
		ethClient:        ethClient,
		orm:              orm,
		listeners:        make(map[common.Address]map[LogListener]struct{}),
		chAddListener:    make(chan registration),
		chRemoveListener: make(chan registration),
		chStop:           make(chan struct{}),
		chDone:           make(chan struct{}),
	}
}

const logBroadcasterCursorName = "logBroadcaster"

func (b *logBroadcaster) Start() {
	// Grab the current on-chain block height
	currentHeight, abort := b.getOnChainBlockHeight()
	if abort {
		return
	}

	// Grab the cursor from the DB
	cursor, err := b.orm.FindLogCursor(logBroadcasterCursorName)
	if err != nil && !gorm.IsRecordNotFoundError(err) {
		logger.Errorf("error fetching log cursor: %v", err)
	}
	b.cursor = cursor

	// If the latest block is newer than the one in the cursor (or if we have
	// no cursor), start from that block height.
	if currentHeight > cursor.BlockIndex {
		b.updateLogCursor(currentHeight, 0)
	}

	go b.startResubscribeLoop()
}

func (b *logBroadcaster) getOnChainBlockHeight() (_ uint64, abort bool) {
	var currentHeight uint64
	for {
		var err error
		currentHeight, err = b.ethClient.GetBlockHeight()
		if err == nil {
			break
		}

		logger.Errorf("error fetching current block height: %v", err)
		select {
		case <-b.chStop:
			return 0, true
		case <-time.After(10 * time.Second):
		}
		continue
	}
	return currentHeight, false
}

func (b *logBroadcaster) Stop() {
	close(b.chStop)
	<-b.chDone
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
ResubscribeLoop:
	for {
		subscription, chRawLogs, err := b.createSubscription()
		if err != nil {
			logger.Errorf("error creating subscription to Ethereum node: %v", err)

			select {
			case <-b.chStop:
				return
			case <-time.After(10 * time.Second):
				// Don't hammer the Ethereum node with subscription requests in case of an error.
				// A configurable timeout might be useful here.
				continue ResubscribeLoop
			}
		}

		b.notifyConnect()

		shouldResubscribe, err := b.process(subscription, chRawLogs)
		if err != nil {
			logger.Error(err)
			b.notifyDisconnect()
			continue ResubscribeLoop
		} else if !shouldResubscribe {
			b.notifyDisconnect()
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

func (b *logBroadcaster) updateLogCursor(blockIdx, logIdx uint64) {
	b.cursor.Initialized = true
	b.cursor.Name = logBroadcasterCursorName
	b.cursor.BlockIndex = blockIdx
	b.cursor.LogIndex = logIdx

	err := b.orm.SaveLogCursor(&b.cursor)
	if err != nil {
		logger.Error("can't save log cursor to DB:", err)
	}
}

func (b *logBroadcaster) process(subscription eth.Subscription, chRawLogs <-chan eth.Log) (shouldResubscribe bool, _ error) {
	defer subscription.Unsubscribe()

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

func (b *logBroadcaster) onRawLog(rawLog eth.Log) {
	// Skip logs that we've already seen
	if b.cursor.Initialized &&
		(rawLog.BlockNumber < b.cursor.BlockIndex ||
			(rawLog.BlockNumber == b.cursor.BlockIndex && uint64(rawLog.Index) <= b.cursor.LogIndex)) {
		return
	}

	for listener := range b.listeners[rawLog.Address] {
		// Make a copy of the log for each listener to avoid data races
		listener.HandleLog(rawLog.Copy(), nil)
	}

	b.updateLogCursor(rawLog.BlockNumber, uint64(rawLog.Index))
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

	if !knownAddress {
		// Recreate the subscription with the new contract address
		return true
	}
	return false
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

func (b *logBroadcaster) createSubscription() (eth.Subscription, chan eth.Log, error) {
	if len(b.listeners) == 0 {
		return noopSubscription{}, nil, nil
	}

	var addresses []common.Address
	for address := range b.listeners {
		addresses = append(addresses, address)
	}

	filterQuery := ethereum.FilterQuery{
		FromBlock: big.NewInt(int64(b.cursor.BlockIndex)),
		Addresses: addresses,
	}
	chRawLogs := make(chan eth.Log)

	subscription, err := b.ethClient.SubscribeToLogs(context.Background(), chRawLogs, filterQuery)
	if err != nil {
		return nil, nil, err
	}
	return subscription, chRawLogs, nil
}

type noopSubscription struct{}

func (s noopSubscription) Err() <-chan error { return nil }
func (s noopSubscription) Unsubscribe()      {}

// DecodingLogListener receives raw logs from the LogBroadcaster and decodes them into
// Go structs using the provided ContractCodec (a simple wrapper around a go-ethereum
// ABI type).
type decodingLogListener struct {
	logTypes map[common.Hash]reflect.Type
	codec    eth.ContractCodec
	LogListener
}

// Ensure that DecodingLogListener conforms to the LogListener interface
var _ LogListener = (*decodingLogListener)(nil)

func NewDecodingLogListener(codec eth.ContractCodec, nativeLogTypes map[common.Hash]interface{}, innerListener LogListener) LogListener {
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

func (l *decodingLogListener) HandleLog(log interface{}, err error) {
	if err != nil {
		l.LogListener.HandleLog(nil, err)
		return
	}

	rawLog, is := log.(eth.Log)
	if !is {
		panic("DecodingLogListener expects to receive an eth.Log")
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
		l.LogListener.HandleLog(nil, err)
		return
	}
	err = l.codec.UnpackLog(decodedLog, event.RawName, rawLog)
	if err != nil {
		l.LogListener.HandleLog(nil, err)
		return
	}

	l.LogListener.HandleLog(decodedLog, nil)
}
