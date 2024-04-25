package client

import (
	"context"
	"sync"
	"time"

	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"

	"github.com/smartcontractkit/chainlink/v2/common/types"
)

// Poller is a component that polls a function at a given interval
// and delivers the result to a channel. It is used to poll for new heads
// and implements the Subscription interface.
type Poller[T any] struct {
	services.StateMachine
	pollingInterval time.Duration
	pollingFunc     func(ctx context.Context, args ...interface{}) (T, error)
	pollingArgs     []interface{}
	pollingTimeout  *time.Duration
	logger          *logger.Logger
	channel         chan<- T
	errCh           chan error

	stopCh chan struct{}
	wg     sync.WaitGroup
}

// NewPoller creates a new Poller instance
func NewPoller[
	T any,
](pollingInterval time.Duration, pollingFunc func(ctx context.Context, args ...interface{}) (T, error), pollingTimeout *time.Duration, channel chan<- T, logger *logger.Logger, args ...interface{}) Poller[T] {
	return Poller[T]{
		pollingInterval: pollingInterval,
		pollingFunc:     pollingFunc,
		pollingArgs:     args,
		pollingTimeout:  pollingTimeout,
		channel:         channel,
		logger:          logger,
		errCh:           make(chan error),
		stopCh:          make(chan struct{}),
	}
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
			pollingCtx := context.Background()
			cancelPolling := context.CancelFunc(func() {})
			if p.pollingTimeout != nil {
				pollingCtx, cancelPolling = context.WithTimeout(pollingCtx, *p.pollingTimeout)
			}

			// Execute polling function in goroutine
			var result T
			var err error
			pollingDone := make(chan struct{})
			go func() {
				defer func() {
					if r := recover(); r != nil {
						err = errors.Errorf("panic: %v", r)
					}
					close(pollingDone)
				}()
				result, err = p.pollingFunc(pollingCtx, p.pollingArgs...)
			}()

			// Wait for polling to complete or timeout
			select {
			case <-pollingCtx.Done():
				cancelPolling()
				p.logError(errors.New("polling timeout exceeded"))
			case <-pollingDone:
				cancelPolling()
				if err != nil {
					p.logError(err)
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
}

func (p *Poller[T]) logError(err error) {
	if p.logger != nil {
		(*p.logger).Errorf("polling error: %v", err)
	}
}
