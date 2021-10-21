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
func NewHeadBroadcaster(logger *logger.Logger) httypes.HeadBroadcaster {
	return &headBroadcaster{
		logger:        logger,
		callbacks:     make(callbackSet),
		mailbox:       utils.NewMailbox(1),
		mutex:         &sync.Mutex{},
		chClose:       make(chan struct{}),
		wgDone:        sync.WaitGroup{},
		StartStopOnce: utils.StartStopOnce{},
	}
}

// headBroadcaster relays heads from the head tracker to subscribed jobs, it is less robust against
// congestion than the head tracker, and missed heads should be expected by consuming jobs
type headBroadcaster struct {
	logger    *logger.Logger
	callbacks callbackSet
	mailbox   *utils.Mailbox
	mutex     *sync.Mutex
	chClose   chan struct{}
	wgDone    sync.WaitGroup
	utils.StartStopOnce
	latest *models.Head
}

var _ httypes.HeadTrackable = (*headBroadcaster)(nil)

func (hr *headBroadcaster) Start() error {
	return hr.StartOnce("HeadBroadcaster", func() error {
		hr.wgDone.Add(1)
		go hr.run()
		return nil
	})
}

func (hr *headBroadcaster) Close() error {
	return hr.StopOnce("HeadBroadcaster", func() error {
		hr.mutex.Lock()
		// clear all callbacks
		hr.callbacks = make(callbackSet)
		hr.mutex.Unlock()

		close(hr.chClose)
		hr.wgDone.Wait()
		return nil
	})
}

func (hr *headBroadcaster) OnNewLongestChain(ctx context.Context, head models.Head) {
	hr.mailbox.Deliver(head)
}

// Subscribe - Subscribes to OnNewLongestChain and Connect until HeadBroadcaster is closed,
// or unsubscribe callback is called explicitly
func (hr *headBroadcaster) Subscribe(callback httypes.HeadTrackable) (currentLongestChain *models.Head, unsubscribe func()) {
	hr.mutex.Lock()
	defer hr.mutex.Unlock()
	currentLongestChain = hr.latest
	id, err := newID()
	if err != nil {
		hr.logger.Errorf("HeadBroadcaster: Unable to create ID for head relayble callback: %v", err)
		return
	}
	hr.callbacks[id] = callback
	unsubscribe = func() {
		hr.mutex.Lock()
		defer hr.mutex.Unlock()
		delete(hr.callbacks, id)
	}
	return
}

func (hr *headBroadcaster) run() {
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
func (hr *headBroadcaster) executeCallbacks() {
	item, exists := hr.mailbox.Retrieve()
	if !exists {
		hr.logger.Info("HeadBroadcaster: no head to retrieve. It might have been skipped")
		return
	}
	head, ok := item.(models.Head)
	if !ok {
		hr.logger.Errorf("expected `models.Head`, got %T", head)
		return
	}
	hr.mutex.Lock()
	callbacks := hr.callbacks.clone()
	hr.latest = &head
	hr.mutex.Unlock()

	hr.logger.Debugw("HeadBroadcaster initiating callbacks",
		"headNum", head.Number,
		"numCallbacks", len(hr.callbacks),
	)

	wg := sync.WaitGroup{}
	wg.Add(len(callbacks))

	for _, callback := range callbacks {
		go func(trackable httypes.HeadTrackable) {
			defer wg.Done()
			start := time.Now()
			ctx, cancel := context.WithTimeout(context.Background(), callbackTimeout)
			defer cancel()
			trackable.OnNewLongestChain(ctx, head)
			elapsed := time.Since(start)
			hr.logger.Debugw(fmt.Sprintf("HeadBroadcaster: finished callback in %s", elapsed), "callbackType", reflect.TypeOf(hr), "blockNumber", head.Number, "time", elapsed, "id", "head_relayer")
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

type NullBroadcaster struct{}

func (*NullBroadcaster) Start() error                                            { return nil }
func (*NullBroadcaster) Close() error                                            { return nil }
func (*NullBroadcaster) OnNewLongestChain(ctx context.Context, head models.Head) {}
func (*NullBroadcaster) Subscribe(callback httypes.HeadTrackable) (currentLongestChain *models.Head, unsubscribe func()) {
	return nil, func() {}
}
func (n *NullBroadcaster) Healthy() error { return nil }
func (n *NullBroadcaster) Ready() error   { return nil }
