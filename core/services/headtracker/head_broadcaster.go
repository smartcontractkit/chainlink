package headtracker

import (
	"context"
	"crypto/rand"
	"fmt"
	"reflect"
	"sync"
	"time"

	"github.com/smartcontractkit/chainlink/core/logger"
	httypes "github.com/smartcontractkit/chainlink/core/services/headtracker/types"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"
)

const callbackTimeout = 2 * time.Second

type callbackID [256]byte

type callbackSet map[callbackID]httypes.HeadTrackable

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
		toUnsubscribe: make([]func(), 0),
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
	callbacks     callbackSet
	toUnsubscribe []func()
	mailbox       *utils.Mailbox
	mutex         *sync.RWMutex
	chClose       chan struct{}
	wgDone        sync.WaitGroup
	utils.StartStopOnce
}

var _ httypes.HeadTrackable = (*HeadBroadcaster)(nil)

func (hr *HeadBroadcaster) Start() error {
	return hr.StartOnce("HeadBroadcaster", func() error {
		hr.wgDone.Add(1)
		go hr.run()
		return nil
	})
}

func (hr *HeadBroadcaster) Close() error {
	return hr.StopOnce("HeadBroadcaster", func() error {

		for _, unsubscribe := range hr.toUnsubscribe {
			unsubscribe()
		}

		close(hr.chClose)
		hr.wgDone.Wait()
		return nil
	})
}

func (hr *HeadBroadcaster) Connect(head *models.Head) error {
	hr.mutex.RLock()
	callbacks := hr.callbacks.clone()
	hr.mutex.RUnlock()

	for i, callback := range callbacks {
		err := callback.Connect(head)
		if err != nil {
			logger.Errorf("HeadBroadcaster: Failed Connect callback at index %v: %v", i, err)
		}
	}

	return nil
}

func (hr *HeadBroadcaster) OnNewLongestChain(ctx context.Context, head models.Head) {
	hr.mailbox.Deliver(head)
}

func (hr *HeadBroadcaster) SubscribeUntilClose(callback httypes.HeadTrackable) {
	hr.toUnsubscribe = append(hr.toUnsubscribe, hr.Subscribe(callback))
}
func (hr *HeadBroadcaster) SubscribeForConnectUntilClose(onConnect func() error) {
	callback := &httypes.HeadTrackableCallback{OnConnect: onConnect}
	hr.toUnsubscribe = append(hr.toUnsubscribe, hr.Subscribe(callback))
}

func (hr *HeadBroadcaster) Subscribe(callback httypes.HeadTrackable) (unsubscribe func()) {
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

	logger.Debugw("HeadBroadcaster initiating callbacks",
		"headNum", head.Number,
		"chainLength", head.ChainLength(),
		"numCallbacks", len(hr.callbacks),
	)

	wg := sync.WaitGroup{}
	wg.Add(len(hr.callbacks))

	for _, callback := range callbacks {
		go func(hr httypes.HeadTrackable) {
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
