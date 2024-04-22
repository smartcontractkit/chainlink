package client

import (
	"sync"
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"

	"github.com/smartcontractkit/chainlink/v2/common/types"
)

// HeadPoller is a component that polls a function at a given interval
// and delivers the result to a channel. It is used to poll for new heads
// and implements the Subscription interface.
type HeadPoller[
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

// TODO: write an error to the Err()<- chan if the parent context is canceled/closed?
// TODO: Do we want to add ctx to the NewHeadPoller constructor?
// TODO: Should we start the polling loop right away or wait for Subsribe?

// NewHeadPoller creates a new HeadPoller instance
func NewHeadPoller[
	HEAD Head,
](pollingInterval time.Duration, pollingFunc func() (HEAD, error), channel chan<- HEAD, logger logger.Logger) HeadPoller[HEAD] {
	return HeadPoller[HEAD]{
		pollingInterval: pollingInterval,
		pollingFunc:     pollingFunc,
		logger:          logger,
		channel:         channel,
		stopCh:          make(chan struct{}),
	}
}

var _ types.Subscription = &HeadPoller[Head]{}

// Subscribe starts the polling process
func (p *HeadPoller[HEAD]) Subscribe() error {
	return p.StartOnce("HeadPoller", func() error {
		p.wg.Add(1)
		go p.pollingLoop()
		return nil
	})
}

// Unsubscribe cancels the sending of events to the data channel
func (p *HeadPoller[HEAD]) Unsubscribe() {
	_ = p.StopOnce("HeadPoller", func() error {
		close(p.stopCh)
		p.wg.Wait()
		return nil
	})
	close(p.errCh)
}

func (p *HeadPoller[HEAD]) Err() <-chan error {
	return p.errCh
}

func (p *HeadPoller[HEAD]) pollingLoop() {
	ticker := time.NewTicker(p.pollingInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			result, err := p.pollingFunc()
			if err != nil {
				p.logger.Error("error occurred when calling polling function:", err)
				p.errCh <- err
				continue
			}
			p.channel <- result
		case <-p.stopCh:
			p.wg.Done()
			return
		}
	}
}
