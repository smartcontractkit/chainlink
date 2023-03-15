package services

import "context"

// ServiceCtx represents a long-running service inside the Application.
//
// Typically, a ServiceCtx will leverage utils.StartStopOnce to implement these
// calls in a safe manner.
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
	Start(context.Context) error
	// Close stops the Service.
	// Invariants: Usually after this call the Service cannot be started
	// again, you need to build a new Service to do so.
	Close() error

	// Name returns the fully qualified name of the service
	Name() string

	Checkable
}
