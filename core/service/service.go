package service

import "github.com/smartcontractkit/chainlink/core/services/health"

type (
	// Service represents a long running service inside the
	// Application.
	//
	// Typically a Service will leverage utils.StartStopOnce to implement these
	// calls in a safe manner.
	//
	// Template
	//
	// Mockable Foo service with a run loop
	//  //go:generate mockery --name Foo --output ../../internal/mocks/ --case=underscore
	//  type (
	//  	// Expose a public interface so we can mock the service.
	//  	Foo interface {
	//  		service.Service
	//
	//  		// ...
	//  	}
	//
	//  	foo struct {
	//  		// ...
	//
	//  		stop chan struct{}
	//  		done chan struct{}
	//
	//  		utils.StartStopOnce
	//  	}
	//  )
	//
	//  var _ Foo = (*foo)(nil)
	//
	//  func NewFoo() Foo {
	//  	f := &foo{
	//  		// ...
	//  	}
	//
	//  	return f
	//  }
	//
	//  func (f *foo) Start() error {
	//  	return f.StartOnce("Foo", func() error {
	//  		go f.run()
	//
	//  		return nil
	//  	})
	//  }
	//
	//  func (f *foo) Close() error {
	//  	return f.StopOnce("Foo", func() error {
	//  		// trigger goroutine cleanup
	//  		close(f.stop)
	//  		// wait for cleanup to complete
	//  		<-f.done
	//  		return nil
	//  	})
	//  }
	//
	//  func (f *foo) run() {
	//  	// signal cleanup completion
	//  	defer close(f.done)
	//
	//  	for {
	//  		select {
	//  		// ...
	//  		case <-f.stop:
	//  			// stop the routine
	//  			return
	//  		}
	//  	}
	//
	//  }
	Service interface {
		// Start the service.
		Start() error
		// Stop the Service.
		// Invariants: Usually after this call the Service cannot be started
		// again, you need to build a new Service to do so.
		Close() error

		health.Checkable
	}
)
