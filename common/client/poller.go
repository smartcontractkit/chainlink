package client

import (
	"sync"
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"

	"github.com/smartcontractkit/chainlink/v2/common/types"
)

// Poller is a component that polls a function at a given interval
// and delivers the result to a channel. It is used to poll for new heads
// and implements the Subscription interface.
type Poller[
	HEAD Head,
] struct {
	services.StateMachine
	pollingInterval time.Duration
	pollingFunc     func() (HEAD, error)
	logger          logger.Logger

	channel chan<- HEAD
	errCh   chan error

	stopCh chan struct{}
	wg     sync.WaitGroup
}

// NewPoller creates a new Poller instance
func NewPoller[
	HEAD Head,
](pollingInterval time.Duration, pollingFunc func() (HEAD, error), channel chan<- HEAD, logger logger.Logger) Poller[HEAD] {
	return Poller[HEAD]{
		pollingInterval: pollingInterval,
		pollingFunc:     pollingFunc,
		logger:          logger,
		channel:         channel,
		stopCh:          make(chan struct{}),
	}
}

var _ types.Subscription = &Poller[Head]{}

// Subscribe starts the polling process
func (p *Poller[HEAD]) Subscribe() error {
	return p.StartOnce("Poller", func() error {
		p.wg.Add(1)
		go p.pollingLoop()
		return nil
	})
}

// Unsubscribe cancels the sending of events to the data channel
func (p *Poller[HEAD]) Unsubscribe() {
	_ = p.StopOnce("Poller", func() error {
		close(p.stopCh)
		p.wg.Wait()
		return nil
	})
	close(p.errCh)
}

func (p *Poller[HEAD]) Err() <-chan error {
	return p.errCh
}

func (p *Poller[HEAD]) pollingLoop() {
	ticker := time.NewTicker(p.pollingInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			result, err := p.pollingFunc()
			if err != nil {
				p.logger.Error("error occurred when calling polling function:", err)
				continue
			}
			p.channel <- result
		case <-p.stopCh:
			p.wg.Done()
			return
		}
	}
}
