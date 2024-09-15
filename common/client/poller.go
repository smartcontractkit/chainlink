package client

import (
	"context"
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"

	"github.com/smartcontractkit/chainlink/v2/common/types"
)

// Poller is a component that polls a function at a given interval
// and delivers the result to a channel. It is used by multinode to poll
// for new heads and implements the Subscription interface.
type Poller[T any] struct {
	services.Service
	eng *services.Engine

	pollingInterval time.Duration
	pollingFunc     func(ctx context.Context) (T, error)
	pollingTimeout  time.Duration
	channel         chan<- T
	errCh           chan error
}

// NewPoller creates a new Poller instance and returns a channel to receive the polled data
func NewPoller[
	T any,
](pollingInterval time.Duration, pollingFunc func(ctx context.Context) (T, error), pollingTimeout time.Duration, lggr logger.Logger) (Poller[T], <-chan T) {
	channel := make(chan T)
	p := Poller[T]{
		pollingInterval: pollingInterval,
		pollingFunc:     pollingFunc,
		pollingTimeout:  pollingTimeout,
		channel:         channel,
		errCh:           make(chan error),
	}
	p.Service, p.eng = services.Config{
		Name:  "Poller",
		Start: p.start,
		Close: p.close,
	}.NewServiceEngine(lggr)
	return p, channel
}

var _ types.Subscription = &Poller[any]{}

func (p *Poller[T]) start(ctx context.Context) error {
	p.eng.Go(p.pollingLoop)
	return nil
}

// Unsubscribe cancels the sending of events to the data channel
func (p *Poller[T]) Unsubscribe() {
	_ = p.Close()
}

func (p *Poller[T]) close() error {
	close(p.errCh)
	close(p.channel)
	return nil
}

func (p *Poller[T]) Err() <-chan error {
	return p.errCh
}

func (p *Poller[T]) pollingLoop(ctx context.Context) {
	ticker := time.NewTicker(p.pollingInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			// Set polling timeout
			pollingCtx, cancelPolling := context.WithTimeout(ctx, p.pollingTimeout)
			// Execute polling function
			result, err := p.pollingFunc(pollingCtx)
			cancelPolling()
			if err != nil {
				p.eng.Warnf("polling error: %v", err)
				continue
			}
			// Send result to channel or block if channel is full
			select {
			case p.channel <- result:
			case <-ctx.Done():
				return
			}
		}
	}
}
