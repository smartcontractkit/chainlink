package client

import (
	"context"
	"sync"
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"

	"github.com/smartcontractkit/chainlink/v2/common/types"
)

// Poller is a component that polls a function at a given interval
// and delivers the result to a channel. It is used by multinode to poll
// for new heads and implements the Subscription interface.
type Poller[T any] struct {
	services.StateMachine
	pollingInterval time.Duration
	pollingFunc     func(ctx context.Context) (T, error)
	pollingTimeout  time.Duration
	logger          logger.Logger
	channel         chan<- T
	errCh           chan error

	stopCh services.StopChan
	wg     sync.WaitGroup
}

// NewPoller creates a new Poller instance and returns a channel to receive the polled data
func NewPoller[
	T any,
](pollingInterval time.Duration, pollingFunc func(ctx context.Context) (T, error), pollingTimeout time.Duration, logger logger.Logger) (Poller[T], <-chan T) {
	channel := make(chan T)
	return Poller[T]{
		pollingInterval: pollingInterval,
		pollingFunc:     pollingFunc,
		pollingTimeout:  pollingTimeout,
		channel:         channel,
		logger:          logger,
		errCh:           make(chan error),
		stopCh:          make(chan struct{}),
	}, channel
}

var _ types.Subscription = &Poller[any]{}

func (p *Poller[T]) Start() error {
	return p.StartOnce("Poller", func() error {
		p.wg.Add(1)
		go p.pollingLoop()
		return nil
	})
}

// Unsubscribe cancels the sending of events to the data channel
func (p *Poller[T]) Unsubscribe() {
	_ = p.StopOnce("Poller", func() error {
		close(p.stopCh)
		p.wg.Wait()
		close(p.errCh)
		close(p.channel)
		return nil
	})
}

func (p *Poller[T]) Err() <-chan error {
	return p.errCh
}

func (p *Poller[T]) pollingLoop() {
	defer p.wg.Done()

	ticker := time.NewTicker(p.pollingInterval)
	defer ticker.Stop()

	for {
		select {
		case <-p.stopCh:
			return
		case <-ticker.C:
			// Set polling timeout
			pollingCtx, cancelPolling := p.stopCh.CtxCancel(context.WithTimeout(context.Background(), p.pollingTimeout))
			// Execute polling function
			result, err := p.pollingFunc(pollingCtx)
			cancelPolling()
			if err != nil {
				p.logger.Warnf("polling error: %v", err)
				continue
			}
			// Send result to channel or block if channel is full
			select {
			case p.channel <- result:
			case <-p.stopCh:
				return
			}
		}
	}
}
