package common

import (
	"context"
	"errors"
	"sync/atomic"

	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers"
)

type Callback struct {
	ch chan handlers.UserCallbackPayload

	sent       atomic.Bool
	waitCalled atomic.Bool
}

func (c *Callback) SendResponse(payload handlers.UserCallbackPayload) error {
	if !c.sent.CompareAndSwap(false, true) {
		return errors.New("response already sent: each callback can only be used once")
	}
	// The channel is initialized with a buffer size of 1,
	// so this send will not block.
	c.ch <- payload
	return nil
}

func (c *Callback) Wait(ctx context.Context) (handlers.UserCallbackPayload, error) {
	if !c.waitCalled.CompareAndSwap(false, true) {
		return handlers.UserCallbackPayload{}, errors.New("Wait can only be called once per Callback instance")
	}
	select {
	case <-ctx.Done():
		return handlers.UserCallbackPayload{}, ctx.Err()
	case r := <-c.ch:
		return r, nil
	}
}

func NewCallback() *Callback {
	ch := make(chan handlers.UserCallbackPayload, 1)
	return &Callback{ch: ch}
}
