package fakes

import (
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/events"
)

// EventEmitterOption configures the emitter used by the manual trigger fakes.
type EventEmitterOption func(*eventEmitterConfig)

type eventEmitterConfig struct {
	eventEmitter *events.Emitter
}

// WithEventEmitter sets the emitter used for trigger execution events.
// The default is events.DefaultEmitter() (global beholder + global durable emitter).
func WithEventEmitter(e *events.Emitter) EventEmitterOption {
	return func(c *eventEmitterConfig) {
		c.eventEmitter = e
	}
}

func newEventEmitter(opts []EventEmitterOption) *events.Emitter {
	c := &eventEmitterConfig{eventEmitter: events.DefaultEmitter()}
	for _, opt := range opts {
		opt(c)
	}
	return c.eventEmitter
}
