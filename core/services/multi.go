package services

import (
	"context"
	"io"
	"sync"

	"go.uber.org/multierr"
)

// StartClose is a subset of the ServiceCtx interface.
type StartClose interface {
	Start(context.Context) error
	Close() error
}

// MultiStart is a utility for starting multiple services together.
// The set of started services is tracked internally, so that they can be closed if any single service fails to start.
type MultiStart struct {
	started []StartClose
}

// Start attempts to Start all services. If any service fails to start, the previously started services will be
// Closed, and an error returned.
func (m *MultiStart) Start(ctx context.Context, srvcs ...StartClose) (err error) {
	for _, s := range srvcs {
		err = m.start(ctx, s)
		if err != nil {
			return err
		}
	}
	return
}

func (m *MultiStart) start(ctx context.Context, s StartClose) (err error) {
	err = s.Start(ctx)
	if err != nil {
		err = multierr.Append(err, m.Close())
	} else {
		m.started = append(m.started, s)
	}
	return
}

// Close closes all started services, in reverse order.
func (m *MultiStart) Close() (err error) {
	for i := len(m.started) - 1; i >= 0; i-- {
		s := m.started[i]
		err = multierr.Append(err, s.Close())
	}
	return
}

// CloseBecause calls Close and returns reason along with any additional errors.
func (m *MultiStart) CloseBecause(reason error) (err error) {
	return multierr.Append(reason, m.Close())
}

// MultiClose is a utility for closing multiple services concurrently.
type MultiClose []io.Closer

// Close closes alls service concurrently and collects any returned errors as a multierr.
func (m MultiClose) Close() (err error) {
	if len(m) == 0 {
		return nil
	}
	var wg sync.WaitGroup
	wg.Add(len(m))
	errs := make(chan error, len(m))
	for _, s := range m {
		go func(c io.Closer) {
			defer wg.Done()
			if e := c.Close(); e != nil {
				errs <- e
			}
		}(s)
	}
	wg.Wait()
	close(errs)
	for e := range errs {
		err = multierr.Append(err, e)
	}
	return
}
