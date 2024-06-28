package client

import (
	"github.com/ethereum/go-ethereum"
)

var _ ethereum.Subscription = &subForwarder[any]{}

// subForwarder wraps a subscription in order to intercept and augment each result before forwarding.
type subForwarder[T any] struct {
	destCh chan<- T

	srcCh  chan T
	srcSub ethereum.Subscription

	interceptResult func(T) T
	interceptError  func(error) error

	done  chan struct{}
	err   chan error
	unSub chan struct{}
}

func newSubForwarder[T any](destCh chan<- T, interceptResult func(T) T, interceptError func(error) error) *subForwarder[T] {
	return &subForwarder[T]{
		interceptResult: interceptResult,
		interceptError:  interceptError,
		destCh:          destCh,
		srcCh:           make(chan T),
		done:            make(chan struct{}),
		err:             make(chan error),
		unSub:           make(chan struct{}, 1),
	}
}

// start spawns the forwarding loop for sub.
func (c *subForwarder[T]) start(sub ethereum.Subscription, err error) error {
	if err != nil {
		close(c.srcCh)
		return err
	}
	c.srcSub = sub
	go c.forwardLoop()
	return nil
}

// forwardLoop receives from src, adds the chainID, and then sends to dest.
// It also handles Unsubscribing, which may interrupt either forwarding operation.
func (c *subForwarder[T]) forwardLoop() {
	// the error channel must be closed when unsubscribing
	defer close(c.err)
	defer close(c.done)

	for {
		select {
		case err := <-c.srcSub.Err():
			if c.interceptError != nil {
				err = c.interceptError(err)
			}
			select {
			case c.err <- err:
			case <-c.unSub:
				c.srcSub.Unsubscribe()
			}
			return

		case h := <-c.srcCh:
			if c.interceptResult != nil {
				h = c.interceptResult(h)
			}
			select {
			case c.destCh <- h:
			case <-c.unSub:
				c.srcSub.Unsubscribe()
				return
			}

		case <-c.unSub:
			c.srcSub.Unsubscribe()
			return
		}
	}
}

func (c *subForwarder[T]) Unsubscribe() {
	// tell forwardLoop to unsubscribe
	select {
	case c.unSub <- struct{}{}:
	default:
		// already triggered
	}
	// wait for forwardLoop to complete
	<-c.done
}

func (c *subForwarder[T]) Err() <-chan error {
	return c.err
}
