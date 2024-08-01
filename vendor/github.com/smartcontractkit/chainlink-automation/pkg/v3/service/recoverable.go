package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"runtime/debug"
	"sync/atomic"
	"time"
)

var (
	ErrServiceAlreadyStarted   = fmt.Errorf("recoverable service already started")
	ErrServiceNotRunning       = fmt.Errorf("recoverable service not running")
	errServiceStopped          = fmt.Errorf("service stopped")
	errServiceContextCancelled = fmt.Errorf("service context cancelled")
)

const (
	PanicRestartWait = 10 * time.Second
)

// Recoverable is a service that a Recoverer can manage
type Recoverable interface {
	// Start is a function that is expected to block and only return on
	// completion either with an error or nil
	Start(context.Context) error
	// Close stops the execution of the blocking Start function and causes an
	// immediate return
	Close() error
}

// NewRecoverer creates a new configured recoverer
func NewRecoverer(svc Recoverable, logger *log.Logger) *recoverer {
	return &recoverer{
		service:  svc,
		log:      logger,
		stopped:  make(chan error, 1),
		coolDown: PanicRestartWait,
	}
}

// recoverer assists in managing an underlying service by recovering
// automatically from panics if they occur. This is intended to add a layer of
// resilience to an underlying service.
type recoverer struct {
	// dependencies
	service Recoverable
	log     *log.Logger

	// created by constructor
	stopped  chan error
	coolDown time.Duration

	// internal state
	running atomic.Bool
}

// Start starts the recoverable service and the recovery watcher and returns an
// error if the recoverer is already running
func (m *recoverer) Start(ctx context.Context) error {
	if m.running.Load() {
		return ErrServiceAlreadyStarted
	}

	go m.recoverableStart(ctx)

	m.serviceStart(ctx)

	return nil
}

// Stop stops the recoverable service and recovery watcher and returns an error
// if the recoverer is already stopped
func (m *recoverer) Close() error {
	if !m.running.Load() {
		return ErrServiceNotRunning
	}

	err := m.service.Close()

	select {
	case m.stopped <- errServiceContextCancelled:
	default:
	}

	return err
}

func (m *recoverer) serviceStart(ctx context.Context) {
	m.running.Store(true)

	for {
		select {
		case err := <-m.stopped:
			// restart the service
			if err != nil {
				if errors.Is(err, errServiceStopped) {
					<-time.After(m.coolDown)
					go m.recoverableStart(ctx)
				}

				if errors.Is(err, errServiceContextCancelled) {
					m.running.Store(false)
					return
				}
			}
		case <-ctx.Done():
			m.running.Store(false)
			return
		}
	}
}

func (m *recoverer) recoverableStart(ctx context.Context) {
	func(s Recoverable, l *log.Logger, chStop chan error, ctx context.Context) {
		defer func() {
			if err := recover(); err != nil {
				if l != nil {
					l.Println(err)
					l.Println(string(debug.Stack()))
				}

				chStop <- errServiceStopped
			}
		}()

		err := s.Start(ctx)

		if l != nil && err != nil {
			l.Println(err)
		}

		chStop <- err
	}(m.service, m.log, m.stopped, ctx)
}
