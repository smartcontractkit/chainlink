package service

import (
	"reflect"
	"sync"

	"github.com/smartcontractkit/chainlink/core/logger"
	"go.uber.org/multierr"
)

// Usage
//
// mgr := NewManager()
// mgr.Register(promReporter)
// mgr.Register(fluxmonitorSvc, RunConcurrently())
//
// mgr.Run()
//
// mgr.Shutdown(merr)

// ServiceRunner defines how a service is to be run.
type ServiceRunner struct {
	Svc Service

	runConcurrently bool
	chStop          chan struct{}
}

type Manager struct {
	// wgDone waits until all concurrent services have been closed
	wgDone sync.WaitGroup

	// Public for now until we can move health checking into here.
	Runners []ServiceRunner
}

func NewManager() *Manager {
	return &Manager{
		Runners: []ServiceRunner{},
	}
}

type serviceRunnerOption func(runner *ServiceRunner)

func (s *Manager) Register(svc Service, opts ...serviceRunnerOption) {
	runner := ServiceRunner{
		Svc: svc,
	}

	for _, opt := range opts {
		opt(&runner)
	}

	s.Runners = append(s.Runners, runner)
}

// RunConcurrently runs the service in a goroutine
func RunConcurrently() serviceRunnerOption {
	return func(runner *ServiceRunner) {
		runner.runConcurrently = true
		runner.chStop = make(chan struct{})
	}
}

// Run starts all services that have been registered.
func (s *Manager) Run() error {
	var err error

	for _, runner := range s.Runners {
		logger.Debugw("Starting service", "serviceType", reflect.TypeOf(runner.Svc))

		if runner.runConcurrently {
			s.wgDone.Add(1)
			go func() {
				defer s.wgDone.Done()

				err = runner.Svc.Start()

				<-runner.chStop
			}()
		} else {
			err = runner.Svc.Start()
		}
	}

	// TODO - Do we need to clean up the goroutines that have already been
	// started successfully? We do not currently do this.
	if err != nil {
		return err
	}

	return nil
}

// Shutdown stops the running services in the reverse order from which they were
// started.
//
// It currently takes err as an argument becuase we want to append to the
// multierrors. When all services get moved into here, we can remove this.
func (s *Manager) Shutdown(merr error) error {
	for i := len(s.Runners) - 1; i >= 0; i-- {
		runner := s.Runners[i]
		logger.Debugw("Closing service...", "serviceType", reflect.TypeOf(runner.Svc))

		err := runner.Svc.Close()
		merr = multierr.Append(merr, err)

		if runner.runConcurrently {
			close(runner.chStop)
		}
	}

	// Wait for all concurrent services running in go routines to finish
	// stopping.
	s.wgDone.Wait()

	return merr
}
