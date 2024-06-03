package headtracker

import (
	"context"
	"fmt"
	"reflect"
	"sync"
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/mailbox"

	"github.com/smartcontractkit/chainlink/v2/common/types"
)

const TrackableCallbackTimeout = 2 * time.Second

type callbackSet[H types.Head[BLOCK_HASH], BLOCK_HASH types.Hashable] map[int]HeadTrackable[H, BLOCK_HASH]

func (set callbackSet[H, BLOCK_HASH]) values() []HeadTrackable[H, BLOCK_HASH] {
	var values []HeadTrackable[H, BLOCK_HASH]
	for _, callback := range set {
		values = append(values, callback)
	}
	return values
}

// HeadTrackable is implemented by the core txm to be able to receive head events from any chain.
// Chain implementations should notify head events to the core txm via this interface.
//
//go:generate mockery --quiet --name HeadTrackable --output ./mocks/ --case=underscore
type HeadTrackable[H types.Head[BLOCK_HASH], BLOCK_HASH types.Hashable] interface {
	// OnNewLongestChain sends a new head when it becomes available. Subscribers can recursively trace the parent
	// of the head to the finalized block back.
	OnNewLongestChain(ctx context.Context, head H)
}

// HeadBroadcaster relays new Heads to all subscribers.
//
//go:generate mockery --quiet --name HeadBroadcaster --output ./mocks/ --case=underscore
type HeadBroadcaster[H types.Head[BLOCK_HASH], BLOCK_HASH types.Hashable] interface {
	services.Service
	BroadcastNewLongestChain(H)
	Subscribe(callback HeadTrackable[H, BLOCK_HASH]) (currentLongestChain H, unsubscribe func())
}

type headBroadcaster[H types.Head[BLOCK_HASH], BLOCK_HASH types.Hashable] struct {
	services.StateMachine
	logger         logger.Logger
	callbacks      callbackSet[H, BLOCK_HASH]
	mailbox        *mailbox.Mailbox[H]
	mutex          sync.Mutex
	chClose        services.StopChan
	wgDone         sync.WaitGroup
	latest         H
	lastCallbackID int
}

// NewHeadBroadcaster creates a new HeadBroadcaster
func NewHeadBroadcaster[
	H types.Head[BLOCK_HASH],
	BLOCK_HASH types.Hashable,
](
	lggr logger.Logger,
) HeadBroadcaster[H, BLOCK_HASH] {
	return &headBroadcaster[H, BLOCK_HASH]{
		logger:    logger.Named(lggr, "HeadBroadcaster"),
		callbacks: make(callbackSet[H, BLOCK_HASH]),
		mailbox:   mailbox.NewSingle[H](),
		chClose:   make(chan struct{}),
	}
}

func (hb *headBroadcaster[H, BLOCK_HASH]) Start(context.Context) error {
	return hb.StartOnce("HeadBroadcaster", func() error {
		hb.wgDone.Add(1)
		go hb.run()
		return nil
	})
}

func (hb *headBroadcaster[H, BLOCK_HASH]) Close() error {
	return hb.StopOnce("HeadBroadcaster", func() error {
		hb.mutex.Lock()
		// clear all callbacks
		hb.callbacks = make(callbackSet[H, BLOCK_HASH])
		hb.mutex.Unlock()

		close(hb.chClose)
		hb.wgDone.Wait()
		return nil
	})
}

func (hb *headBroadcaster[H, BLOCK_HASH]) Name() string {
	return hb.logger.Name()
}

func (hb *headBroadcaster[H, BLOCK_HASH]) HealthReport() map[string]error {
	return map[string]error{hb.Name(): hb.Healthy()}
}

func (hb *headBroadcaster[H, BLOCK_HASH]) BroadcastNewLongestChain(head H) {
	hb.mailbox.Deliver(head)
}

// Subscribe subscribes to OnNewLongestChain and Connect until HeadBroadcaster is closed,
// or unsubscribe callback is called explicitly
func (hb *headBroadcaster[H, BLOCK_HASH]) Subscribe(callback HeadTrackable[H, BLOCK_HASH]) (currentLongestChain H, unsubscribe func()) {
	hb.mutex.Lock()
	defer hb.mutex.Unlock()

	currentLongestChain = hb.latest

	hb.lastCallbackID++
	callbackID := hb.lastCallbackID
	hb.callbacks[callbackID] = callback
	unsubscribe = func() {
		hb.mutex.Lock()
		defer hb.mutex.Unlock()
		delete(hb.callbacks, callbackID)
	}

	return
}

func (hb *headBroadcaster[H, BLOCK_HASH]) run() {
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
func (hb *headBroadcaster[H, BLOCK_HASH]) executeCallbacks() {
	head, exists := hb.mailbox.Retrieve()
	if !exists {
		hb.logger.Info("No head to retrieve. It might have been skipped")
		return
	}

	hb.mutex.Lock()
	callbacks := hb.callbacks.values()
	hb.latest = head
	hb.mutex.Unlock()

	hb.logger.Debugw("Initiating callbacks",
		"headNum", head.BlockNumber(),
		"numCallbacks", len(callbacks),
	)

	wg := sync.WaitGroup{}
	wg.Add(len(callbacks))

	ctx, cancel := hb.chClose.NewCtx()
	defer cancel()

	for _, callback := range callbacks {
		go func(trackable HeadTrackable[H, BLOCK_HASH]) {
			defer wg.Done()
			start := time.Now()
			cctx, cancel := context.WithTimeout(ctx, TrackableCallbackTimeout)
			defer cancel()
			trackable.OnNewLongestChain(cctx, head)
			elapsed := time.Since(start)
			hb.logger.Debugw(fmt.Sprintf("Finished callback in %s", elapsed),
				"callbackType", reflect.TypeOf(trackable), "blockNumber", head.BlockNumber(), "time", elapsed)
		}(callback)
	}

	wg.Wait()
}
