// Package lazyinitservice provides an implementation of the job.ServiceCtx interface, LazyInitService.
//
// This implementation executes the service initialization lazily on the first Start method invocation.
// If the initialization fails, the service keeps trying to initialize the underlying service periodically until the first success.
// The initialization function can indicate that there is no point in retrying using the Unrecoverable error wrapper.
//
// # Testing
//
// If you want to simulate your service initialization failures, you can define TEST_<service_name>_INIT_FAILURES environment variable.
// The value of this environment variable should be an integer representing the number of initialization attempts before lazy service calls the real init function.
// For example, if you want service CCIP_Exec to fail three times before succeeding, set the TEST_CCIP_Exec_INIT_FAILURES env variable to "3".
package lazyinitservice

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"sync"

	"github.com/avast/retry-go/v4"

	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
)

var ErrNoService = errors.New("the service is permanently unavailable")
var ErrNotReady = errors.New("the service is not ready yet")
var ErrClosed = errors.New("the service is closed")

// An InitFunc represents an expensive blocking computation producing a service.
// Init functions must respect the context passed as the argument and quit promptly if the context is canceled.
type InitFunc = func(context.Context) (job.ServiceCtx, error)

// A LogErrorFunc is a callback for reporting background initialization and startup errors.
type LogErrorFunc = func(error)

type Option = func(*LazyInitService)

type LazyInitService struct {
	// name is the underlying service name.
	name string
	// lggr is the logger for the service.
	lggr logger.Logger
	// initFunc is the function creating the service.
	initFunc InitFunc
	// initComplete guards the initialization process allowing for a graceful shutdown.
	initComplete sync.WaitGroup
	// logErrorFunc is the function for logging errors occurring in background.
	logErrorFunc LogErrorFunc
	// cancelFunc is the function canceling the initialization process.
	cancelFunc context.CancelFunc
	// mu guards the fields below.
	mu sync.Mutex
	// initializedService contains the service the initFunc returns.
	initializedService job.ServiceCtx
	// err is the last error reported by the service.
	lastErr error
}

// WithLogErrorFunc instructs the service constructor to use the given function for error reporting.
func WithLogErrorFunc(f LogErrorFunc) Option {
	return func(s *LazyInitService) {
		s.logErrorFunc = f
	}
}

// New creates a new service with the given initialization function.
func New(lggr logger.Logger, name string, f InitFunc, opts ...Option) *LazyInitService {
	s := &LazyInitService{
		name:     name,
		lggr:     lggr,
		initFunc: f,
	}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

// Start initiates the underlying service initialization and starts it.
//
// Start ignores the given ctx cancellation.
// Use Close to stop the initialization process and the service.
func (s *LazyInitService) Start(_ctx context.Context) error {
	s.initComplete.Add(1)

	// We create a new context because the original context will be cancelled once `Start` returns.
	ctx, cancelFunc := context.WithCancel(context.Background())
	s.cancelFunc = cancelFunc
	go s.initAndRun(ctx)
	return nil
}

// initAndRun implements the lazy initialization logic.
func (s *LazyInitService) initAndRun(ctx context.Context) {
	defer s.initComplete.Done()
	s.setState(nil, ErrNotReady)

	initFailures := 0
	testEnvVar := fmt.Sprintf("TEST_%s_INIT_FAILURES", s.name)
	if v := os.Getenv(testEnvVar); v != "" {
		parsed, err := strconv.Atoi(v)
		if err != nil {
			s.lggr.Warnw("failed to parse environment variable", "service", s.name, "var", testEnvVar, "value", v, "err", err)
			s.reportError(fmt.Errorf("failed to parse env var %s value %s: %w", testEnvVar, v, err))
		}
		initFailures = parsed
	}
	n := 0
	service, err := retry.DoWithData[job.ServiceCtx](
		func() (job.ServiceCtx, error) {
			if n < initFailures {
				n++
				return nil, ErrNotReady
			}
			return s.initFunc(ctx)
		},
		retry.Context(ctx),
		retry.OnRetry(func(n uint, err error) {
			s.setState(nil, err)
			s.lggr.Warnw("service initialization failed", "service", s.name, "attempt", n, "err", err)
			s.reportError(fmt.Errorf("initialization attempt %d failed: %w", n, err))
		}),
	)
	if err != nil {
		s.setState(nil, err)
		s.lggr.Errorw("service initialization failed", "service", s.name, "err", err)
		s.reportError(err)
		return
	}
	if service == nil {
		s.setState(nil, ErrNoService)
		s.lggr.Errorw("service init function returned nil", "service", s.name)
		s.reportError(ErrNoService)
		return
	}
	s.setState(service, ErrNotReady)
	if err = service.Start(ctx); err != nil {
		s.setState(service, err)
		s.lggr.Errorw("service failed to start", "service", s.name, "err", err)
		s.reportError(fmt.Errorf("service failed to start: %w", err))
	}
	s.lggr.Infow("service started", "service", s.name)
	s.setState(service, nil)
}

// reportError records the given error using the service log error function.
func (s *LazyInitService) reportError(err error) {
	if s.logErrorFunc != nil {
		s.logErrorFunc(err)
	}
}

func (s *LazyInitService) setState(service job.ServiceCtx, err error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.initializedService = service
	s.lastErr = err
}

func (s *LazyInitService) getState() (service job.ServiceCtx, lastErr error) {
	s.mu.Lock()
	service, lastErr = s.initializedService, s.lastErr
	s.mu.Unlock()
	return
}

// Close implements graceful service shutdown logic.
func (s *LazyInitService) Close() error {
	s.lggr.Infow("closing service", "service", s.name)
	// First, cancel the context to break the initialization retry loop.
	if s.cancelFunc != nil {
		s.cancelFunc()
		s.cancelFunc = nil
	}
	// Second, wait for the initialization to complete.
	s.initComplete.Wait()

	service, _ := s.getState()

	// Now, we can close the internal service if it was initialized.
	if service != nil {
		err := service.Close()
		if err != nil {
			s.lggr.Warnw("service failed to close", "service", s.name)
			s.setState(service, err)
		} else {
			s.lggr.Infow("service closed", "service", s.name)
			s.setState(service, ErrClosed)
		}
		return err
	}
	return nil
}

func (s *LazyInitService) Ready() error {
	service, lastErr := s.getState()

	if service != nil {
		if r, ok := service.(services.HealthReporter); ok {
			return r.Ready()
		}
	}
	return lastErr
}

func (s *LazyInitService) HealthReport() map[string]error {
	service, lastErr := s.getState()

	if service != nil {
		if r, ok := service.(services.HealthReporter); ok {
			report := r.HealthReport()
			report[s.name] = lastErr
			return report
		}
	}
	return map[string]error{s.name: lastErr}
}

func (s *LazyInitService) Name() string {
	return s.name
}

// Unrecoverable wraps the given error into an error that signals to the retry mechanism to stop trying.
func Unrecoverable(err error) error {
	return retry.Unrecoverable(err)
}
