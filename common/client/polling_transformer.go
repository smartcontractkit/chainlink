package client

import (
	"context"
	"sync"
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
)

// PollingTransformer is a component that polls a function at a given interval
// and delivers the result to subscribers
type PollingTransformer[HEAD Head] struct {
	interval time.Duration
	pollFunc func() (HEAD, error)

	logger logger.Logger

	subscribers []chan<- HEAD

	isPolling bool
	stopCh    services.StopChan
	wg        sync.WaitGroup
}

func NewPollingTransformer[HEAD Head](pollInterval time.Duration, pollFunc func() (HEAD, error), logger logger.Logger) *PollingTransformer[HEAD] {
	return &PollingTransformer[HEAD]{
		interval:  pollInterval,
		pollFunc:  pollFunc,
		logger:    logger,
		isPolling: false,
	}
}

// Subscribe adds a Subscriber to the polling transformer
func (pt *PollingTransformer[HEAD]) Subscribe(sub chan<- HEAD) {
	pt.subscribers = append(pt.subscribers, sub)
}

// Unsubscribe removes a Subscriber from the polling transformer
func (pt *PollingTransformer[HEAD]) Unsubscribe(sub chan<- HEAD) {
	for i, s := range pt.subscribers {
		if s == sub {
			close(s)
			pt.subscribers = append(pt.subscribers[:i], pt.subscribers[i+1:]...)
			return
		}
	}
}

// StartPolling starts the polling loop and delivers the polled value to subscribers
func (pt *PollingTransformer[HEAD]) StartPolling() {
	pt.stopCh = make(chan struct{})
	pt.wg.Add(1)
	go pt.pollingLoop(pt.stopCh.NewCtx())
	pt.isPolling = true
}

// pollingLoop polls the pollFunc at the given interval and delivers the result to subscribers
func (pt *PollingTransformer[HEAD]) pollingLoop(ctx context.Context, cancel context.CancelFunc) {
	defer pt.wg.Done()
	defer cancel()

	pollT := time.NewTicker(pt.interval)
	defer pollT.Stop()

	for {
		select {
		case <-ctx.Done():
			for _, subscriber := range pt.subscribers {
				close(subscriber)
			}
			return
		case <-pollT.C:
			head, err := pt.pollFunc()
			if err != nil {
				// TODO: handle error
			}
			pt.logger.Debugw("PollingTransformer: polled value", "head", head)
			for _, subscriber := range pt.subscribers {
				select {
				case subscriber <- head:
					// Successfully sent head
				default:
					// Subscriber's channel is closed
					pt.Unsubscribe(subscriber)
				}
			}
		}
	}
}

// StopPolling stops the polling loop
func (pt *PollingTransformer[HEAD]) StopPolling() {
	close(pt.stopCh)
	pt.wg.Wait()
}
