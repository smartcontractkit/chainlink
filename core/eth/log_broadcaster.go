package eth

import (
	"context"
	"fmt"
	"math/big"
	"reflect"
	"sync"
	"time"

	"chainlink/core/logger"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
)

//go:generate mockery -name LogBroadcaster -output ../internal/mocks/ -case=underscore
//go:generate mockery -name LogListener -output ../internal/mocks/ -case=underscore

type LogBroadcaster interface {
	Start()
	Register(address common.Address, listener LogListener)
	Unregister(address common.Address, listener LogListener)
	Stop()
}

type LogListener interface {
	OnConnect()
	OnDisconnect()
	HandleLog(log interface{}, err error)
}

type logBroadcaster struct {
	ethClient Client
	cursor    logCursor

	listeners        map[common.Address]map[LogListener]struct{}
	chAddListener    chan registration
	chRemoveListener chan registration

	chStop chan struct{}
	chDone chan struct{}
}

type logCursor struct {
	initialized bool
	blockIdx    uint64
	logIdx      uint64
}

type registration struct {
	address  common.Address
	listener LogListener
}

func NewLogBroadcaster(ethClient Client) LogBroadcaster {
	return &logBroadcaster{
		ethClient:        ethClient,
		listeners:        make(map[common.Address]map[LogListener]struct{}),
		chAddListener:    make(chan registration),
		chRemoveListener: make(chan registration),
		chStop:           make(chan struct{}),
		chDone:           make(chan struct{}),
	}
}

func (b *logBroadcaster) Start() {
	go b.startResubscribeLoop()
}

func (b *logBroadcaster) Stop() {
	close(b.chStop)
	<-b.chDone
}

func (b *logBroadcaster) Register(address common.Address, listener LogListener) {
	select {
	case b.chAddListener <- registration{address, listener}:
	case <-b.chStop:
	}
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
			default:
				// Don't hammer the Ethereum node with subscription requests in case of an error.
				// A configurable timeout might be useful here.
				time.Sleep(10 * time.Second)
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
	for _, listeners := range b.listeners {
		for listener := range listeners {
			listener.OnConnect()
		}
	}
}

func (b *logBroadcaster) notifyDisconnect() {
	for _, listeners := range b.listeners {
		for listener := range listeners {
			listener.OnDisconnect()
		}
	}
}

func (b *logBroadcaster) process(subscription Subscription, chRawLogs <-chan Log) (shouldResubscribe bool, _ error) {
	defer subscription.Unsubscribe()
ProcessLoop:
	for {
		select {
		case err := <-subscription.Err():
			return true, err

		case rawLog := <-chRawLogs:
			// Skip logs that we've already seen
			if b.cursor.initialized &&
				(rawLog.BlockNumber < b.cursor.blockIdx ||
					(rawLog.BlockNumber == b.cursor.blockIdx && uint64(rawLog.Index) <= b.cursor.logIdx)) {
				continue ProcessLoop
			}

			b.broadcast(rawLog)

			// @@TODO: persist logCursor in DB
			b.cursor.initialized = true
			b.cursor.blockIdx = rawLog.BlockNumber
			b.cursor.logIdx = uint64(rawLog.Index)

		case r := <-b.chAddListener:
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
				return true, nil
			}

		case r := <-b.chRemoveListener:
			delete(b.listeners[r.address], r.listener)
			if len(b.listeners[r.address]) == 0 {
				delete(b.listeners, r.address)
				// Recreate the subscription without this contract address
				return true, nil
			}

		case <-b.chStop:
			return false, nil
		}
	}
}

func (b *logBroadcaster) broadcast(rawLog Log) {
	var wg sync.WaitGroup
	wg.Add(len(b.listeners[rawLog.Address]))
	for listener := range b.listeners[rawLog.Address] {
		rawLogCopy := rawLog.Copy()
		listener := listener
		go func() {
			defer wg.Done()
			listener.HandleLog(rawLogCopy, nil)
		}()
	}
	wg.Wait()
}

func (b *logBroadcaster) createSubscription() (Subscription, chan Log, error) {
	if len(b.listeners) == 0 {
		return noopSubscription{}, nil, nil
	}

	var fromBlock *big.Int
	if b.cursor.blockIdx > 0 {
		fromBlock = big.NewInt(int64(b.cursor.blockIdx))
	}

	var addresses []common.Address
	for address := range b.listeners {
		addresses = append(addresses, address)
	}

	filterQuery := ethereum.FilterQuery{
		FromBlock: fromBlock,
		Addresses: addresses,
	}
	chRawLogs := make(chan Log)

	subscription, err := b.ethClient.SubscribeToLogs(context.Background(), chRawLogs, filterQuery)
	if err != nil {
		return nil, nil, err
	}
	return subscription, chRawLogs, nil
}

type noopSubscription struct{}

func (s noopSubscription) Err() <-chan error { return nil }
func (s noopSubscription) Unsubscribe()      {}

type decodingLogListener struct {
	logTypes map[common.Hash]reflect.Type
	codec    ContractCodec
	LogListener
}

// Ensure that DecodingLogListener conforms to the LogListener interface
var _ LogListener = (*decodingLogListener)(nil)

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

func (l *decodingLogListener) HandleLog(log interface{}, err error) {
	if err != nil {
		l.LogListener.HandleLog(nil, err)
		return
	}

	rawLog, is := log.(Log)
	if !is {
		panic("DecodingLogListener expects to receive an eth.Log")
	}

	eventID := rawLog.Topics[0]

	logType, exists := l.logTypes[eventID]
	if !exists {
		panic(fmt.Sprintf("DecodingLogListener got unknown log with topic %v", eventID.Hex()))
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
