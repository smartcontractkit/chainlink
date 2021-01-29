package log

import (
	"context"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"
)

type relayer struct {
	orm       ORM
	listeners map[common.Address]map[Listener]struct{}

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
	address  common.Address
	listener Listener
}

type connectionEvent int

const (
	disconnected connectionEvent = iota
	connected
)

func newRelayer(orm ORM, dependentAwaiter utils.DependentAwaiter) *relayer {
	return &relayer{
		orm:              orm,
		listeners:        make(map[common.Address]map[Listener]struct{}),
		addListener:      utils.NewMailbox(500),
		rmListener:       utils.NewMailbox(500),
		newHeads:         utils.NewMailbox(1),
		newLogs:          utils.NewMailbox(1000),
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
			r.onRemoveListeners()

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

func (r *relayer) NotifyAddListener(contractAddr common.Address, listener Listener) {
	r.addListener.Deliver(registration{contractAddr, listener})
}

func (r *relayer) NotifyRemoveListener(contractAddr common.Address, listener Listener) {
	r.rmListener.Deliver(registration{contractAddr, listener})
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
	for {
		select {
		case <-r.addListener.Notify():
			r.onAddListeners()

		case <-r.rmListener.Notify():
			r.onRemoveListeners()

		case <-r.newLogs.Notify():
			r.onNewLogs()

		case <-r.newHeads.Notify():
			r.onNewHeads()

		case <-r.connectionEvents.Notify():
			r.onConnectionEvents()
		}
	}
}

func (r *relayer) onAddListeners() {
	for {
		x := r.addListener.Retrieve()
		if x == nil {
			break
		}
		reg := x.(registration)
		_, knownAddress := r.listeners[reg.address]
		if !knownAddress {
			r.listeners[reg.address] = make(map[Listener]struct{})
		}
		if _, exists := r.listeners[reg.address][reg.listener]; exists {
			panic("registration already exists")
		}
		r.listeners[reg.address][reg.listener] = struct{}{}
	}
}

func (r *relayer) onRemoveListeners() {
	for {
		x := r.rmListener.Retrieve()
		if x == nil {
			break
		}
		reg := x.(registration)
		reg.listener.OnDisconnect()
		delete(r.listeners[reg.address], reg.listener)
		if len(r.listeners[reg.address]) == 0 {
			delete(r.listeners, reg.address)
		}
	}
}

func (r *relayer) onNewLogs() {
	for {
		x := r.newLogs.Retrieve()
		if x == nil {
			break
		}
		log := x.(types.Log)
		for listener := range r.listeners[log.Address] {
			err := r.orm.UpsertUnconsumedLogBroadcastForListener(log, listener)
			if err != nil {
				logger.Errorw("could not upsert log consumption record",
					"contract", log.Address,
					"block", log.BlockHash,
					"tx", log.TxHash,
					"logIndex", log.Index,
					"removed", log.Removed,
					"jobID", listener.JobID(),
					"jobIDV2", listener.JobIDV2(),
					"error", err,
				)
			}
		}
	}
}

func (r *relayer) onNewHeads() {
	for {
		x := r.newHeads.Retrieve()
		if x == nil {
			break
		}
		head := x.(models.Head)

		logs, err := r.orm.UnconsumedLogsPriorToBlock(uint64(head.Number) - r.fewestRequestedConfirmations)
		if err != nil {
			logger.Errorw("could not fetch logs to broadcast", "error", err)
			return
		}
		for _, log := range logs {
			r.broadcast(log)
		}
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
			lb := &broadcast{
				orm:     r.orm,
				rawLog:  logCopy,
				jobID:   listener.JobID(),
				jobIDV2: listener.JobIDV2(),
				isV2:    listener.IsV2Job(),
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
		evt := x.(connectionEvent)
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
