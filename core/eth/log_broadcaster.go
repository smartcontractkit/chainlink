package contracts

import (
	"math/big"
	"reflect"
	"sync"

	"github.com/ethereum/go-ethereum/common"
)

//go:generate mockery -name LogBroadcaster -output ../../internal/mocks/ -case=underscore
//go:generate mockery -name LogListener -output ../../internal/mocks/ -case=underscore

type LogBroadcaster interface {
	Start()
	Register(listener LogListener)
	Unregister(listener LogListener)
	Stop()
}

type LogListener interface {
	Receive(rawLog Log)
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
	blockIdx uint64
	txIdx    uint64
	logIdx   uint64
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
	go b.resubscribe()
}

func (b *logBroadcaster) Stop() {
	b.chStop <- struct{}{}
	<-b.chDone
}

func (b *logBroadcaster) Register(address common.Address, listener LogListener) {
	b.chAddListener <- registration{address, listener}
}

func (s *logBroadcaster) Unregister(address common.Address, listener LogListener) {
	s.chRemoveListener <- registration{address, listener}
}

func (s *logBroadcaster) resubscribe() {
	for {
		shouldResubscribe := s.subscribe()
		if !shouldResubscribe {
			return
		}
	}
}

func (s *logBroadcaster) subscribe() (shouldResubscribe bool, _ error) {
	subscription, chRawLogs, err := s.createSubscription()
	if err != nil {
		return false, err
	}
	defer subscription.Unsubscribe()

Loop:
	for {
		select {
		case rawLog := <-chRawLogs:
			// Skip logs that we've already seen
			if rawLog.BlockNumber < s.cursor.blockIdx ||
				(rawLog.BlockNumber == s.cursor.blockIdx && uint64(rawLog.TxIndex) < s.cursor.txIdx) ||
				(rawLog.BlockNumber == s.cursor.blockIdx && uint64(rawLog.TxIndex) == s.cursor.txIdx && uint64(rawLog.Index) < s.cursor.logIdx) {
				continue Loop
			}

			s.broadcast(rawLog)

			// @@TODO: persist logCursor in DB
			s.cursor.blockIdx = rawLog.BlockNumber
			s.cursor.txIdx = uint64(rawLog.TxIndex)
			s.cursor.logIdx = uint64(rawLog.Index)

		case r := <-s.chAddListener:
			_, knownAddress := s.listeners[r.address]
			if !knownAddress {
				s.listeners[r.address] = make(map[LogListener]struct{})
			}
			if _, exists := s.listeners[r.address][r.listener]; exists {
				panic("registration already exists")
			}
			s.listeners[r.address][r.listener] = struct{}{}

			if !knownAddress {
				// Recreate the subscription with the new contract address
				return true, nil
			}

		case r := <-s.chRemoveListener:
			delete(s.listeners[r.address], r.listener)
			if len(s.listeners[r.address]) == 0 {
				delete(s.listeners[r.address])
				// Recreate the subscription without this contract address
				return true, nil
			}

		case <-s.chClose:
			return false, nil
		}
	}
}

func (b *logBroadcaster) broadcast(rawLog Log) {
	var wg sync.WaitGroup
	wg.Add(len(b.listeners[rawLog.Address]))
	for listener := range b.listeners[rawLog.Address] {
		rawLogCopy := rawLog.Copy()
		go func() {
			defer wg.Done()
			listener.Receive(rawLogCopy)
		}()
	}
	wg.Wait()
}

func (s *logBroadcaster) createSubscription() (Subscription, chan Log, error) {
	var fromBlock *big.Int
	if s.cursor.blockIdx > 0 {
		fromBlock = big.NewInt(int64(s.cursor.blockIdx))
	}

	var addresses []common.Address
	for address := range s.listeners {
		addresses = append(addresses, address)
	}

	filterQuery := ethereum.FilterQuery{
		FromBlock: fromBlock,
		Addresses: addresses,
	}
	chRawLogs := make(chan Log)

	subscription, err := s.ethClient.SubscribeToLogs(chRawLogs, filterQuery)
	if err != nil {
		return nil, nil, err
	}
	return subscription, chRawLogs, nil
}

type DecodingLogListener struct {
	logTypes  map[common.Hash]reflect.Type
	contract  Contract
	handlerFn DecodedLogHandlerFunc
}

type DecodedLogHandlerFunc func(decodedLog interface{}, err error)

// Ensure that DecodingLogListener conforms to the LogListener interface
var _ LogListener = (DecodingLogListener)(nil)

func NewDecodingLogListener(contract Contract, nativeLogTypes map[common.Hash]interface{}, handlerFn DecodedLogHandlerFunc) *DecodingLogListener {
	logTypes := make(map[common.Hash]reflect.Type)
	for eventID, logStruct := range nativeLogTypes {
		logTypes[eventID] = reflect.TypeOf(logStruct)
	}

	return &DecodingLogListener{
		logTypes:  logTypes,
		contract:  contract,
		handlerFn: handlerFn,
	}
}

func (l *DecodingLogListener) Receive(rawLog Log) {
	eventID := rawLog.Topics[0]

	logType, exists := l.logTypes[eventID]
	if !exists {
		// @@TODO: ??
		return
	}

	var decodedLog interface{}
	if logType.Kind() == reflect.Ptr {
		decodedLog = reflect.New(logType.Elem()).Interface()
	} else {
		decodedLog = reflect.New(logType).Interface()
	}

	event, err := l.contract.ABI().EventByID(eventID)
	if err != nil {
		l.handlerFn(nil, err)
		return
	}

	err := l.contract.UnpackLog(decodedLog, eventName, rawLog)
	if err != nil {
		l.handlerFn(nil, err)
		return
	}

	l.handlerFn(decodedLog, nil)
}
