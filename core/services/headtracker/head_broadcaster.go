package headtracker

import (
	"context"
	"crypto/rand"
	"fmt"
	"reflect"
	"sync"
	"time"

	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/eth"
	httypes "github.com/smartcontractkit/chainlink/core/services/headtracker/types"
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
func NewHeadBroadcaster(logger logger.Logger) httypes.HeadBroadcaster {
	return &headBroadcaster{
		logger:        logger.Named("HeadBroadcaster"),
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
	logger    logger.Logger
	callbacks callbackSet
	mailbox   *utils.Mailbox
	mutex     *sync.Mutex
	chClose   chan struct{}
	wgDone    sync.WaitGroup
	utils.StartStopOnce
	latest *eth.Head
}

var _ httypes.HeadTrackable = (*headBroadcaster)(nil)

func (hb *headBroadcaster) Start() error {
	return hb.StartOnce("HeadBroadcaster", func() error {
		hb.wgDone.Add(1)
		go hb.run()
		return nil
	})
}

func (hb *headBroadcaster) Close() error {
	return hb.StopOnce("HeadBroadcaster", func() error {
		hb.mutex.Lock()
		// clear all callbacks
		hb.callbacks = make(callbackSet)
		hb.mutex.Unlock()

		close(hb.chClose)
		hb.wgDone.Wait()
		return nil
	})
}

func (hb *headBroadcaster) OnNewLongestChain(ctx context.Context, head *eth.Head) {
	hb.mailbox.Deliver(head)
}

// Subscribe - Subscribes to OnNewLongestChain and Connect until HeadBroadcaster is closed,
// or unsubscribe callback is called explicitly
func (hb *headBroadcaster) Subscribe(callback httypes.HeadTrackable) (currentLongestChain *eth.Head, unsubscribe func()) {
	if callback == nil {
		panic("callback must be non-nil func")
	}
	hb.mutex.Lock()
	defer hb.mutex.Unlock()
	currentLongestChain = hb.latest
	id, err := newID()
	if err != nil {
		hb.logger.Errorf("Unable to create ID for head relayble callback: %v", err)
		return
	}
	hb.callbacks[id] = callback
	unsubscribe = func() {
		hb.mutex.Lock()
		defer hb.mutex.Unlock()
		delete(hb.callbacks, id)
	}
	return
}

func (hb *headBroadcaster) run() {
	defer hb.wgDone.Done()
	for {
		select {
		case <-hb.chClose:
			return
		case <-hb.mailbox.Notify():
			hb.executeCallbacks()
		}
	}
}

// DEV: the head relayer makes no promises about head delivery! Subscribing
// Jobs should expect to the relayer to skip heads if there is a large number of listeners
// and all callbacks cannot be completed in the allotted time.
func (hb *headBroadcaster) executeCallbacks() {
	item, exists := hb.mailbox.Retrieve()
	if !exists {
		hb.logger.Info("No head to retrieve. It might have been skipped")
		return
	}
	head, ok := item.(*eth.Head)
	if !ok {
		hb.logger.Errorf("Expected `eth.Head`, got %T", head)
		return
	}
	hb.mutex.Lock()
	callbacks := hb.callbacks.clone()
	hb.latest = head
	hb.mutex.Unlock()

	hb.logger.Debugw("Initiating callbacks",
		"headNum", head.Number,
		"numCallbacks", len(callbacks),
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
			hb.logger.Debugw(fmt.Sprintf("Finished callback in %s", elapsed),
				"callbackType", reflect.TypeOf(trackable), "blockNumber", head.Number, "time", elapsed)
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

func (*NullBroadcaster) Start() error                                          { return nil }
func (*NullBroadcaster) Close() error                                          { return nil }
func (*NullBroadcaster) OnNewLongestChain(ctx context.Context, head *eth.Head) {}
func (*NullBroadcaster) Subscribe(callback httypes.HeadTrackable) (currentLongestChain *eth.Head, unsubscribe func()) {
	return nil, func() {}
}
func (n *NullBroadcaster) Healthy() error { return nil }
func (n *NullBroadcaster) Ready() error   { return nil }
