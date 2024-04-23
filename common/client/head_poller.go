package client

import (
	"sync"
	"time"

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
	channel         chan<- HEAD
	errCh           chan error

	stopCh chan struct{}
	wg     sync.WaitGroup
}

// NewHeadPoller creates a new HeadPoller instance
func NewHeadPoller[
	HEAD Head,
](pollingInterval time.Duration, pollingFunc func() (HEAD, error), channel chan<- HEAD) HeadPoller[HEAD] {
	return HeadPoller[HEAD]{
		pollingInterval: pollingInterval,
		pollingFunc:     pollingFunc,
		channel:         channel,
		errCh:           make(chan error),
		stopCh:          make(chan struct{}),
	}
}

var _ types.Subscription = &HeadPoller[Head]{}

func (p *HeadPoller[HEAD]) Start() error {
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
		close(p.errCh)
		return nil
	})
}

func (p *HeadPoller[HEAD]) Err() <-chan error {
	return p.errCh
}

func (p *HeadPoller[HEAD]) pollingLoop() {
	defer p.wg.Done()

	ticker := time.NewTicker(p.pollingInterval)
	defer ticker.Stop()

	for {
		select {
		case <-p.stopCh:
			return
		case <-ticker.C:
			result, err := p.pollingFunc()
			if err != nil {
				select {
				case p.errCh <- err:
					continue
				case <-p.stopCh:
					return
				}
			}

			select {
			case p.channel <- result:
			case <-p.stopCh:
				return
			}
		}
	}
}
