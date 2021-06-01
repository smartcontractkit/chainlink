package services

import (
	"context"
	"crypto/rand"
	"fmt"
	"reflect"
	"sync"
	"time"

	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"
)

const callbackTimeout = 2 * time.Second

type callbackID [256]byte

// HeadBroadcastable defines the interface for listeners
type HeadBroadcastable interface {
	OnNewLongestChain(ctx context.Context, head models.Head)
}

type callbackSet map[callbackID]HeadBroadcastable

func (set callbackSet) clone() callbackSet {
	cp := make(callbackSet)
	for id, callback := range set {
		cp[id] = callback
	}
	return cp
}

// NewHeadBroadcaster creates a new HeadBroadcaster
func NewHeadBroadcaster() *HeadBroadcaster {
	return &HeadBroadcaster{
		callbacks:     make(callbackSet),
		mailbox:       utils.NewMailbox(1),
		mutex:         &sync.RWMutex{},
		chClose:       make(chan struct{}),
		wgDone:        sync.WaitGroup{},
		StartStopOnce: utils.StartStopOnce{},
	}
}

// HeadBroadcaster relays heads from the head tracker to subscribed jobs, it is less robust against
// congestion than the head tracker, and missed heads should be expected by consuming jobs
type HeadBroadcaster struct {
	callbacks callbackSet
	mailbox   *utils.Mailbox
	mutex     *sync.RWMutex
	chClose   chan struct{}
	wgDone    sync.WaitGroup
	utils.StartStopOnce
}

var _ models.HeadTrackable = (*HeadBroadcaster)(nil)

func (hr *HeadBroadcaster) Start() error {
	return hr.StartOnce("HeadBroadcaster", func() error {
		hr.wgDone.Add(1)
		go hr.run()
		return nil
	})
}

func (hr *HeadBroadcaster) Close() error {
	return hr.StopOnce("HeadBroadcaster", func() error {
		close(hr.chClose)
		hr.wgDone.Wait()
		return nil
	})
}

func (hr *HeadBroadcaster) Connect(head *models.Head) error {
	return nil
}

func (hr *HeadBroadcaster) Disconnect() {}

func (hr *HeadBroadcaster) OnNewLongestChain(ctx context.Context, head models.Head) {
	hr.mailbox.Deliver(head)
}

func (hr *HeadBroadcaster) Subscribe(callback HeadBroadcastable) (unsubscribe func()) {
	hr.mutex.Lock()
	defer hr.mutex.Unlock()
	id, err := newID()
	if err != nil {
		logger.Errorf("HeadBroadcaster: Unable to create ID for head relayble callback: %v", err)
		return
	}
	hr.callbacks[id] = callback
	return func() {
		hr.mutex.Lock()
		defer hr.mutex.Unlock()
		delete(hr.callbacks, id)
	}
}

func (hr *HeadBroadcaster) run() {
	defer hr.wgDone.Done()
	for {
		select {
		case <-hr.chClose:
			return
		case <-hr.mailbox.Notify():
			hr.executeCallbacks()
		}
	}
}

// DEV: the head relayer makes no promises about head delivery! Subscribing
// Jobs should expect to the relayer to skip heads if there is a large number of listeners
// and all callbacks cannot be completed in the allotted time.
func (hr *HeadBroadcaster) executeCallbacks() {
	hr.mutex.RLock()
	callbacks := hr.callbacks.clone()
	hr.mutex.RUnlock()

	item, exists := hr.mailbox.Retrieve()
	if !exists {
		logger.Info("HeadBroadcaster: no head to retrieve. It might have been skipped")
		return
	}
	head, ok := item.(models.Head)
	if !ok {
		logger.Errorf("expected `models.Head`, got %T", head)
		return
	}

	wg := sync.WaitGroup{}
	wg.Add(len(hr.callbacks))

	for _, callback := range callbacks {
		go func(hr HeadBroadcastable) {
			defer wg.Done()
			start := time.Now()
			ctx, cancel := context.WithTimeout(context.Background(), callbackTimeout)
			defer cancel()
			hr.OnNewLongestChain(ctx, head)
			elapsed := time.Since(start)
			logger.Debugw(fmt.Sprintf("HeadBroadcaster: finished callback in %s", elapsed), "callbackType", reflect.TypeOf(hr), "blockNumber", head.Number, "time", elapsed, "id", "head_relayer")
		}(callback)
	}

	wg.Wait()
}

func newID() (id callbackID, _ error) {
	randBytes := make([]byte, 256)
	_, err := rand.Read(randBytes)
	if err != nil {
		return id, err
	}
	copy(id[:], randBytes)
	return id, nil
}
