package services

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/google/uuid"
	"golang.org/x/exp/maps"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

// ServiceCtx represents a long-running service inside the Application.
//
// Typically, a ServiceCtx will leverage utils.StartStopOnce to implement these
// calls in a safe manner, or bootstrap via New.
//
// # Template
//
// Mockable Foo service with a run loop
//
//	//go:generate mockery --quiet --name Foo --output ../internal/mocks/ --case=underscore
//	type (
//		// Expose a public interface so we can mock the service.
//		Foo interface {
//			service.ServiceCtx
//
//			// ...
//		}
//
//		foo struct {
//			// ...
//
//			stop chan struct{}
//			done chan struct{}
//
//			utils.StartStopOnce
//		}
//	)
//
//	var _ Foo = (*foo)(nil)
//
//	func NewFoo() Foo {
//		f := &foo{
//			// ...
//		}
//
//		return f
//	}
//
//	func (f *foo) Start(ctx context.Context) error {
//		return f.StartOnce("Foo", func() error {
//			go f.run()
//
//			return nil
//		})
//	}
//
//	func (f *foo) Close() error {
//		return f.StopOnce("Foo", func() error {
//			// trigger goroutine cleanup
//			close(f.stop)
//			// wait for cleanup to complete
//			<-f.done
//			return nil
//		})
//	}
//
//	func (f *foo) run() {
//		// signal cleanup completion
//		defer close(f.done)
//
//		for {
//			select {
//			// ...
//			case <-f.stop:
//				// stop the routine
//				return
//			}
//		}
//
//	}
type ServiceCtx interface {
	// Start the service. Must quit immediately if the context is cancelled.
	// The given context applies to Start function only and must not be retained.
	//
	// See MultiStart
	Start(context.Context) error
	// Close stops the Service.
	// Invariants: Usually after this call the Service cannot be started
	// again, you need to build a new Service to do so.
	//
	// See MultiCloser
	Close() error

	Checkable
}

// Group tracks a group of goroutines and provides shutdown signals.
type Group struct {
	wg sync.WaitGroup
	Health
	utils.StopChan
	logger.Logger
}

// Go calls fn in a tracked goroutine that will block closing the service.
// fn should yield to StopChan promptly.
func (g *Group) Go(fn func()) {
	g.wg.Add(1)
	go func() {
		defer g.wg.Done()
		fn()
	}()
}

type Health struct {
	m  map[string]error
	mu sync.RWMutex
}

// SetUnwell records a condition key and an error, which causes an unhealthy report, until SetWell(condition) is called.
func (h *Health) SetUnwell(condition string, err error) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.m[condition] = fmt.Errorf("%s: %w", condition, err)
}

// SetWell removes a condition and error recorded by SetUnwell.
func (h *Health) SetWell(condition string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	delete(h.m, condition)
}

// Unwell causes an unhealthy report, until the returned func() is called.
// Use this for simple cases where the func() can be kept in scope, and prefer to defer it inline if possible:
//
//	defer Unwell(fmt.Errorf("foo bar: %w", err))()
//
// See SetUnwell for an alternative API.
func (h *Health) Unwell(err error) func() {
	cond := uuid.NewString()
	h.SetUnwell(cond, err)
	return func() { h.SetWell(cond) }
}

func (h *Health) healthy() (err error) {
	h.mu.RLock()
	errs := maps.Values(h.m)
	h.mu.RUnlock()
	return errors.Join(errs...)
}

// Spec specifies a service for New().
type Spec struct {
	Name        string
	Start       func(context.Context) error
	SubServices []ServiceCtx
}

type service struct {
	utils.StartStopOnce
	g    Group
	spec Spec
}

// New returns a new ServiceCtx defined by Spec and a Group for managing goroutines and logging.
// You *should* embed the ServiceCtx (to inherit methods), but *not* the Group:
//
//	type example struct {
//		ServiceCtx
//		g *Group
//	}
func New(spec Spec, lggr logger.Logger) (ServiceCtx, *Group) {
	s := &service{
		g: Group{
			StopChan: make(utils.StopChan),
			Logger:   lggr.Named(spec.Name),
			Health:   Health{m: make(map[string]error)},
		},
		spec: spec,
	}
	return s, &s.g
}

// Ready implements [Checkable.Ready] and overrides and extends [utils.StartStopOnce.Ready()] to include [Spec.SubServices]
// readiness as well.
func (s *service) Ready() (err error) {
	err = s.StartStopOnce.Ready()
	for _, sub := range s.spec.SubServices {
		err = errors.Join(err, sub.Ready())
	}
	return
}

// Healthy overrides [utils.StartStopOnce.Healthy] and extends it to include Group errors as well.
// Do not override this method in your service. Instead, report errors via the Group.
func (s *service) Healthy() (err error) {
	err = s.StartStopOnce.Healthy()
	if err == nil {
		err = s.g.healthy()
	}
	return
}

func (s *service) HealthReport() map[string]error {
	m := map[string]error{s.Name(): s.Healthy()}
	for _, sub := range s.spec.SubServices {
		CopyHealth(m, sub.HealthReport())
	}
	return m
}

func (s *service) Name() string { return s.g.Logger.Name() }

func (s *service) Start(ctx context.Context) error {
	return s.StartOnce(s.spec.Name, func() error {
		var ms MultiStart
		s.g.Logger.Debug("Starting sub-services")
		for _, sub := range s.spec.SubServices {
			if err := ms.Start(ctx, sub); err != nil {
				s.g.Logger.Errorw("Failed to start sub-service", "error", err)
				return fmt.Errorf("failed to start sub-service: %w", err)
			}
		}
		return s.spec.Start(ctx)
	})
}

func (s *service) Close() error {
	return s.StopOnce(s.spec.Name, func() (err error) {
		s.g.Logger.Debug("Stopping sub-services")
		close(s.g.StopChan)
		defer s.g.wg.Wait()

		return MultiCloser(s.spec.SubServices).Close()
	})
}
