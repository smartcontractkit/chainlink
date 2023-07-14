package headmanager

import (
	"context"
	"fmt"
	"reflect"
	"sync"
	"time"

	"github.com/smartcontractkit/chainlink/v2/common/types"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

const TrackableCallbackTimeout = 2 * time.Second

type callbackSet[H types.Head[BLOCK_HASH], BLOCK_HASH types.Hashable] map[int]types.HeadTrackable[H, BLOCK_HASH]

func (set callbackSet[H, BLOCK_HASH]) values() []types.HeadTrackable[H, BLOCK_HASH] {
	var values []types.HeadTrackable[H, BLOCK_HASH]
	for _, callback := range set {
		values = append(values, callback)
	}
	return values
}

type Broadcaster[H types.Head[BLOCK_HASH], BLOCK_HASH types.Hashable] struct {
	logger    logger.Logger
	callbacks callbackSet[H, BLOCK_HASH]
	mailbox   *utils.Mailbox[H]
	mutex     *sync.Mutex
	chClose   utils.StopChan
	wgDone    sync.WaitGroup
	utils.StartStopOnce
	latest         H
	lastCallbackID int
}

// NewBroadcaster creates a new Broadcaster
func NewBroadcaster[
	H types.Head[BLOCK_HASH],
	BLOCK_HASH types.Hashable,
](
	lggr logger.Logger,
) *Broadcaster[H, BLOCK_HASH] {
	return &Broadcaster[H, BLOCK_HASH]{
		logger:        lggr.Named("Broadcaster"),
		callbacks:     make(callbackSet[H, BLOCK_HASH]),
		mailbox:       utils.NewSingleMailbox[H](),
		mutex:         &sync.Mutex{},
		chClose:       make(chan struct{}),
		wgDone:        sync.WaitGroup{},
		StartStopOnce: utils.StartStopOnce{},
	}
}

func (b *Broadcaster[H, BLOCK_HASH]) Start(context.Context) error {
	return b.StartOnce("Broadcaster", func() error {
		b.wgDone.Add(1)
		go b.run()
		return nil
	})
}

func (b *Broadcaster[H, BLOCK_HASH]) Close() error {
	return b.StopOnce("Broadcaster", func() error {
		b.mutex.Lock()
		// clear all callbacks
		b.callbacks = make(callbackSet[H, BLOCK_HASH])
		b.mutex.Unlock()

		close(b.chClose)
		b.wgDone.Wait()
		return nil
	})
}

func (b *Broadcaster[H, BLOCK_HASH]) Name() string {
	return b.logger.Name()
}

func (b *Broadcaster[H, BLOCK_HASH]) HealthReport() map[string]error {
	return map[string]error{b.Name(): b.StartStopOnce.Healthy()}
}

func (b *Broadcaster[H, BLOCK_HASH]) BroadcastNewLongestChain(head H) {
	b.mailbox.Deliver(head)
}

// Subscribe subscribes to OnNewLongestChain and Connect until Broadcaster is closed,
// or unsubscribe callback is called explicitly
func (b *Broadcaster[H, BLOCK_HASH]) Subscribe(callback types.HeadTrackable[H, BLOCK_HASH]) (currentLongestChain H, unsubscribe func()) {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	currentLongestChain = b.latest

	b.lastCallbackID++
	callbackID := b.lastCallbackID
	b.callbacks[callbackID] = callback
	unsubscribe = func() {
		b.mutex.Lock()
		defer b.mutex.Unlock()
		delete(b.callbacks, callbackID)
	}

	return
}

func (b *Broadcaster[H, BLOCK_HASH]) run() {
	defer b.wgDone.Done()

	for {
		select {
		case <-b.chClose:
			return
		case <-b.mailbox.Notify():
			b.executeCallbacks()
		}
	}
}

// DEV: the head relayer makes no promises about head delivery! Subscribing
// Jobs should expect to the relayer to skip heads if there is a large number of listeners
// and all callbacks cannot be completed in the allotted time.
func (b *Broadcaster[H, BLOCK_HASH]) executeCallbacks() {
	head, exists := b.mailbox.Retrieve()
	if !exists {
		b.logger.Info("No head to retrieve. It might have been skipped")
		return
	}

	b.mutex.Lock()
	callbacks := b.callbacks.values()
	b.latest = head
	b.mutex.Unlock()

	b.logger.Debugw("Initiating callbacks",
		"headNum", head.BlockNumber(),
		"numCallbacks", len(callbacks),
	)

	wg := sync.WaitGroup{}
	wg.Add(len(callbacks))

	ctx, cancel := b.chClose.NewCtx()
	defer cancel()

	for _, callback := range callbacks {
		go func(trackable types.HeadTrackable[H, BLOCK_HASH]) {
			defer wg.Done()
			start := time.Now()
			cctx, cancel := context.WithTimeout(ctx, TrackableCallbackTimeout)
			defer cancel()
			trackable.OnNewLongestChain(cctx, head)
			elapsed := time.Since(start)
			b.logger.Debugw(fmt.Sprintf("Finished callback in %s", elapsed),
				"callbackType", reflect.TypeOf(trackable), "blockNumber", head.BlockNumber(), "time", elapsed)
		}(callback)
	}

	wg.Wait()
}
