package log

import (
	"context"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"
)

type relayer struct {
	orm       ORM
	config    Config
	listeners map[common.Address]map[Listener]struct{}
	decoders  map[common.Address]AbigenContract

	latestBlock                  uint64
	fewestRequestedConfirmations uint64

	addListener      *utils.Mailbox
	rmListener       *utils.Mailbox
	newHeads         *utils.Mailbox
	newLogs          *utils.Mailbox
	connectionEvents *utils.Mailbox

	utils.StartStopOnce
	utils.DependentAwaiter
	chStop chan struct{}
	chDone chan struct{}
}

// A `registration` represents a Listener's subscription to the logs of a
// particular contract.
type registration struct {
	contract AbigenContract
	listener Listener
}

type connectionEvent int

const (
	disconnected connectionEvent = iota
	connected
)

func newRelayer(orm ORM, config Config, dependentAwaiter utils.DependentAwaiter) *relayer {
	return &relayer{
		orm:              orm,
		config:           config,
		listeners:        make(map[common.Address]map[Listener]struct{}),
		decoders:         make(map[common.Address]AbigenContract),
		addListener:      utils.NewMailbox(0),
		rmListener:       utils.NewMailbox(0),
		newHeads:         utils.NewMailbox(1),
		newLogs:          utils.NewMailbox(10000),
		connectionEvents: utils.NewMailbox(1),
		DependentAwaiter: dependentAwaiter,
		chStop:           make(chan struct{}),
		chDone:           make(chan struct{}),
	}
}

func (r *relayer) Start() error {
	return r.StartOnce("Log relayer", func() (err error) {
		go r.awaitInitialSubscribers()
		return nil
	})
}

func (r *relayer) awaitInitialSubscribers() {
	for {
		select {
		case <-r.addListener.Notify():
			r.onAddListeners()

		case <-r.rmListener.Notify():
			r.onRmListeners()

		case <-r.DependentAwaiter.AwaitDependents():
			go r.runLoop()
			return

		case <-r.chStop:
			close(r.chDone)
			return
		}
	}
}

func (r *relayer) Stop() error {
	return r.StopOnce("Log relayer", func() (err error) {
		close(r.chStop)
		<-r.chDone
		return nil
	})
}

func (r *relayer) NotifyAddListener(contract AbigenContract, listener Listener) {
	r.addListener.Deliver(registration{contract, listener})
}

func (r *relayer) NotifyRemoveListener(contract AbigenContract, listener Listener) {
	r.rmListener.Deliver(registration{contract, listener})
}

func (r *relayer) NotifyNewLog(log types.Log) {
	r.newLogs.Deliver(log)
}

func (r *relayer) OnNewLongestChain(ctx context.Context, head models.Head) {
	r.newHeads.Deliver(head)
}

func (r *relayer) NotifyConnected() {
	r.connectionEvents.Deliver(connected)
}

func (r *relayer) NotifyDisconnected() {
	r.connectionEvents.Deliver(disconnected)
}

func (r *relayer) runLoop() {
	defer close(r.chDone)

	// DB polling is an absolute worst-case fallback. It should
	// not be relied upon to guarantee 100% log delivery.
	dbPoll := time.NewTicker(r.config.TriggerFallbackDBPollInterval())
	defer dbPoll.Stop()

	for {
		select {
		case <-r.addListener.Notify():
			r.onAddListeners()

		case <-r.rmListener.Notify():
			r.onRmListeners()

		case <-r.newLogs.Notify():
			r.onNewLogs()

		case <-r.newHeads.Notify():
			r.onNewHeads()

		case <-dbPoll.C:
			r.onDBPoll()

		case <-r.connectionEvents.Notify():
			r.onConnectionEvents()

		case <-r.chStop:
			return
		}
	}
}

func (r *relayer) onAddListeners() {
	for {
		x := r.addListener.Retrieve()
		if x == nil {
			break
		}
		reg, ok := x.(registration)
		if !ok {
			logger.Errorf("expected `registration`, got %T", x)
			continue
		}
		_, knownAddress := r.listeners[reg.contract.Address()]
		if !knownAddress {
			r.listeners[reg.contract.Address()] = make(map[Listener]struct{})
		}
		if _, exists := r.listeners[reg.contract.Address()][reg.listener]; exists {
			panic("registration already exists")
		}
		r.listeners[reg.contract.Address()][reg.listener] = struct{}{}
		r.decoders[reg.contract.Address()] = reg.contract

		err := r.orm.UpsertBroadcastsForListenerSinceBlock(r.latestBlock, reg.contract.Address(), ListenerJobID(reg.listener))
		if err != nil {
			logger.Errorw("error upserting log broadcast",
				"error", err,
				"contract", reg.contract.Address(),
				"latestBlock", r.latestBlock,
				"jobID", reg.listener.JobID(),
				"jobIDV2", reg.listener.JobIDV2(),
			)
		}
	}
	r.broadcastAllUnconsumed()
}

func (r *relayer) onRmListeners() {
	for {
		x := r.rmListener.Retrieve()
		if x == nil {
			break
		}
		reg, ok := x.(registration)
		if !ok {
			logger.Errorf("expected `registration`, got %T", x)
			continue
		}
		reg.listener.OnDisconnect()
		delete(r.listeners[reg.contract.Address()], reg.listener)
		if len(r.listeners[reg.contract.Address()]) == 0 {
			delete(r.listeners, reg.contract.Address())
		}

		err := r.orm.DeleteUnconsumedBroadcastsForListener(ListenerJobID(reg.listener))
		if err != nil {
			logger.Errorw("could not delete unconsumed log broadcasts for unregistering listener", "error", err)
		}
	}
}

func (r *relayer) onNewLogs() {
	for {
		x := r.newLogs.Retrieve()
		if x == nil {
			break
		}
		log, ok := x.(types.Log)
		if !ok {
			logger.Errorf("expected `types.Log`, got %T", x)
			continue
		}
		for listener := range r.listeners[log.Address] {
			err := r.orm.UpsertBroadcastForListener(log, ListenerJobID(listener))
			if err != nil {
				logger.Errorw("could not upsert log consumption record",
					"contract", log.Address,
					"block", log.BlockHash,
					"blockNumber", log.BlockNumber,
					"tx", log.TxHash,
					"logIndex", log.Index,
					"removed", log.Removed,
					"jobID", listener.JobID(),
					"jobIDV2", listener.JobIDV2(),
					"error", err,
				)
			}
		}
		r.broadcastAllUnconsumed()
	}
}

func (r *relayer) onNewHeads() {
	for {
		x := r.newHeads.Retrieve()
		if x == nil {
			break
		}
		head, ok := x.(models.Head)
		if !ok {
			logger.Errorf("expected `models.Head`, got %T", x)
			continue
		}
		r.latestBlock = uint64(head.Number)
	}
	r.broadcastAllUnconsumed()
}

func (r *relayer) onDBPoll() {
	r.broadcastAllUnconsumed()
}

func (r *relayer) broadcastAllUnconsumed() {
	logs, err := r.orm.UnconsumedLogsPriorToBlock(r.latestBlock + 1 - r.fewestRequestedConfirmations)
	if err != nil {
		logger.Errorw("could not fetch logs to broadcast", "error", err)
		return
	}
	for _, log := range logs {
		r.broadcast(log)
	}
}

func (r *relayer) broadcast(log types.Log) {
	var wg sync.WaitGroup
	wg.Add(len(r.listeners[log.Address]))

	for listener := range r.listeners[log.Address] {
		listener := listener
		go func() {
			defer wg.Done()

			// Deep copy the log so that subscribers aren't sharing any state
			logCopy := copyLog(log)

			// Decode the log
			decoder, exists := r.decoders[log.Address]
			if !exists || decoder == nil {
				logger.Errorw("log decoder for contract is not registered", "contract", log.Address)
				return
			}
			decodedLog, err := decoder.ParseLog(logCopy)
			if err != nil {
				logger.Errorw("could not decode log", "contract", log.Address, "topic", log.Topics[0])
				return
			}

			lb := &broadcast{
				orm:        r.orm,
				rawLog:     logCopy,
				decodedLog: decodedLog,
				jobID:      listener.JobID(),
				jobIDV2:    listener.JobIDV2(),
				isV2:       listener.IsV2Job(),
			}
			listener.HandleLog(lb, nil)
		}()
	}
	wg.Wait()
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

func (r *relayer) onConnectionEvents() {
	for {
		x := r.connectionEvents.Retrieve()
		if x == nil {
			break
		}
		evt, ok := x.(connectionEvent)
		if !ok {
			logger.Errorf("expected `connectedEvent`, got %T", x)
			continue
		}
		if evt == connected {
			for _, listeners := range r.listeners {
				for listener := range listeners {
					listener.OnConnect()
				}
			}

		} else if evt == disconnected {
			for _, listeners := range r.listeners {
				for listener := range listeners {
					listener.OnDisconnect()
				}
			}

		} else {
			logger.Errorw("got unknown connection event", "event", evt)
		}
	}
}
