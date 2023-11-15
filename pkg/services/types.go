package services

import "context"

// Service represents a long-running service inside the Application.
//
// Typically, a Service will leverage utils.StateMachine to implement these
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
//			service.Service
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
type Service interface {
	// Start the service.
	//  - Must return promptly if the context is cancelled.
	//  - Must not retain the context after returning (only applies to start-up)
	//  - Must not depend on external resources (no blocking network calls)
	Start(context.Context) error
	// Close stops the service.
	// Invariants: Usually after this call the Service cannot be started
	// again, you need to build a new Service to do so.
	Close() error

	HealthReporter
}

// HealthReporter should be implemented by any type requiring health checks.
type HealthReporter interface {
	// Ready should return nil if ready, or an error message otherwise. From the k8s docs:
	// > ready means it’s initialized and healthy means that it can accept traffic in kubernetes
	// See: https://kubernetes.io/docs/tasks/configure-pod-container/configure-liveness-readiness-startup-probes/
	Ready() error
	// HealthReport returns a full health report of the callee including its dependencies.
	// key is the dep name, value is nil if healthy, or error message otherwise.
	// Use CopyHealth to collect reports from sub-services.
	HealthReport() map[string]error
	// Name returns the fully qualified name of the component. Usually the logger name.
	Name() string
}
