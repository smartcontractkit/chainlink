package common

import (
	"context"
	"errors"
	"sync"

	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers"
)

type Callback struct {
	c    chan handlers.UserCallbackPayload
	mu   sync.Mutex
	sent bool
}

func (c *Callback) SendResponse(ctx context.Context, payload handlers.UserCallbackPayload) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.sent {
		return errors.New("response already sent: each callback can only be used once")
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	case c.c <- payload:
		c.sent = true
		return nil
	}
}

func (c *Callback) Wait(ctx context.Context) (handlers.UserCallbackPayload, error) {
	select {
	case <-ctx.Done():
		return handlers.UserCallbackPayload{}, ctx.Err()
	case r := <-c.c:
		return r, nil
	}
}

func NewCallback() *Callback {
	ch := make(chan handlers.UserCallbackPayload, 1)
	return &Callback{c: ch}
}
