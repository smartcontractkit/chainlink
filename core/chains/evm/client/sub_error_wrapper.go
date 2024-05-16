package client

import (
	"fmt"

	commontypes "github.com/smartcontractkit/chainlink/v2/common/types"
)

// subErrorWrapper - adds specified prefix to a subscription error
type subErrorWrapper struct {
	sub         commontypes.Subscription
	errorPrefix string

	done    chan struct{}
	unSub   chan struct{}
	errorCh chan error
}

func newSubscriptionErrorWrapper(sub commontypes.Subscription, errorPrefix string) *subErrorWrapper {
	s := &subErrorWrapper{
		sub:         sub,
		errorPrefix: errorPrefix,
		done:        make(chan struct{}),
		unSub:       make(chan struct{}),
		errorCh:     make(chan error),
	}

	go func() {
		for {
			select {
			// sub.Err channel is closed by sub.Unsubscribe
			case err, ok := <-sub.Err():
				if !ok {
					// might only happen if someone terminated wrapped subscription
					// in any case - do our best to release resources
					// we can't call Unsubscribe on root sub as this might cause panic
					close(s.errorCh)
					close(s.done)
					return
				}

				select {
				case s.errorCh <- fmt.Errorf("%s: %w", s.errorPrefix, err):
				case <-s.unSub:
					s.close()
					return
				}
			case <-s.unSub:
				s.close()
				return
			}
		}
	}()

	return s
}

func (s *subErrorWrapper) close() {
	s.sub.Unsubscribe()
	close(s.errorCh)
	close(s.done)
}

func (s *subErrorWrapper) Unsubscribe() {
	select {
	// already unsubscribed
	case <-s.done:
	// signal unsubscribe
	case s.unSub <- struct{}{}:
		// wait for unsubscribe to complete
		<-s.done
	}
}

func (s *subErrorWrapper) Err() <-chan error {
	return s.errorCh
}
