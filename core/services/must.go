package services

import (
	"context"

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
