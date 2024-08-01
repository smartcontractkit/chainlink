package services

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/google/uuid"
	"golang.org/x/exp/maps"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/timeutil"
)

// Service represents a long-running service inside the Application.
//
// The simplest way to implement a Service is via NewService.
//
// For other cases, consider embedding a services.StateMachine to implement these
// calls in a safe manner.
type Service interface {
	// Start the service.
	//  - Must return promptly if the context is cancelled.
	//  - Must not retain the context after returning (only applies to start-up)
	//  - Must not depend on external resources (no blocking network calls)
	Start(context.Context) error
	// Close stops the Service.
	// Invariants: Usually after this call the Service cannot be started
	// again, you need to build a new Service to do so.
	//
	// See MultiCloser
	Close() error

	HealthReporter
}

// Engine manages service internals like health, goroutine tracking, and shutdown signals.
// See Config.NewServiceEngine
type Engine struct {
	StopChan
	logger.SugaredLogger

	wg sync.WaitGroup

	emitHealthErr func(error)
	conds         map[string]error
	condsMu       sync.RWMutex
}

// Go runs fn in a tracked goroutine that will block closing the service.
func (e *Engine) Go(fn func(context.Context)) {
	e.wg.Add(1)
	go func() {
		defer e.wg.Done()
		ctx, cancel := e.StopChan.NewCtx()
		defer cancel()
		fn(ctx)
	}()
}

// GoTick is like Go but calls fn for each tick.
//
//	v.e.GoTick(services.NewTicker(time.Minute), v.method)
func (e *Engine) GoTick(ticker *timeutil.Ticker, fn func(context.Context)) {
	e.Go(func(ctx context.Context) {
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				fn(ctx)
			}
		}
	})
}

// EmitHealthErr records an error to be reported via the next call to Healthy().
func (e *Engine) EmitHealthErr(err error) { e.emitHealthErr(err) }

// SetHealthCond records a condition key and an error, which causes an unhealthy report, until ClearHealthCond(condition) is called.
// condition keys are for internal use only, and do not show up in the health report.
func (e *Engine) SetHealthCond(condition string, err error) {
	e.condsMu.Lock()
	defer e.condsMu.Unlock()
	e.conds[condition] = fmt.Errorf("%s: %e", condition, err)
}

// ClearHealthCond removes a condition and error recorded by SetHealthCond.
func (e *Engine) ClearHealthCond(condition string) {
	e.condsMu.Lock()
	defer e.condsMu.Unlock()
	delete(e.conds, condition)
}

// NewHealthCond causes an unhealthy report, until the returned clear func() is called.
// Use this for simple cases where the func() can be kept in scope, and prefer to defer it inline if possible:
//
//	defer NewHealthCond(fmt.Errorf("foo bar: %i", err))()
//
// See SetHealthCond for an alternative API.
func (e *Engine) NewHealthCond(err error) (clear func()) {
	cond := uuid.NewString()
	e.SetHealthCond(cond, err)
	return func() { e.ClearHealthCond(cond) }
}

func (e *Engine) clearCond() error {
	e.condsMu.RLock()
	errs := maps.Values(e.conds)
	e.condsMu.RUnlock()
	return errors.Join(errs...)
}

// Config is a configuration for constructing a Service, typically with an Engine, to be embedded and extended as part
// of a Service implementation.
type Config struct {
	// Name is required. It will be logged shorthand on Start and Close, and appended to the logger name.
	// It must be unique among services sharing the same logger, in order to ensure uniqueness of the fully qualified name.
	Name string
	// NewSubServices is an optional constructor for dependent Services to Start and Close along with this one.
	NewSubServices func(logger.Logger) []Service
	// Start is an optional hook called after starting SubServices.
	Start func(context.Context) error
	// Close is an optional hook called before closing SubServices.
	Close func() error
}

// NewServiceEngine returns a new Service defined by Config, and an Engine for managing health, goroutines, and logging.
//   - You *should* embed the Service, in order to inherit the methods.
//   - You *should not* embed the Engine. Use an unexported field instead.
//
// For example:
//
//	type myType struct {
//		services.Service
//		env *service.Engine
//	}
//	t := myType{}
//	t.Service, t.eng = service.Config{
//		Name: "MyType",
//		Start: t.start,
//		Close: t.close,
//	}.NewServiceEngine(lggr)
func (c Config) NewServiceEngine(lggr logger.Logger) (Service, *Engine) {
	s := c.new(logger.Sugared(lggr))
	return s, &s.eng
}

// NewService returns a new Service defined by Config.
//   - You *should* embed the Service, in order to inherit the methods.
func (c Config) NewService(lggr logger.Logger) Service {
	return c.new(logger.Sugared(lggr))
}

func (c Config) new(lggr logger.SugaredLogger) *service {
	lggr = lggr.Named(c.Name)
	s := &service{
		cfg: c,
		eng: Engine{
			StopChan:      make(StopChan),
			SugaredLogger: lggr,
			conds:         make(map[string]error),
		},
	}
	s.eng.emitHealthErr = s.StateMachine.SvcErrBuffer.Append
	if c.NewSubServices != nil {
		s.subs = c.NewSubServices(lggr)
	}
	return s
}

type service struct {
	StateMachine
	cfg  Config
	eng  Engine
	subs []Service
}

// Ready implements [HealthReporter.Ready] and overrides and extends [utils.StartStopOnce.Ready()] to include [Config.SubServices]
// readiness as well.
func (s *service) Ready() (err error) {
	err = s.StateMachine.Ready()
	for _, sub := range s.subs {
		err = errors.Join(err, sub.Ready())
	}
	return
}

// Healthy overrides [StateMachine.Healthy] and extends it to include Engine errors as well.
// Do not override this method in your service. Instead, report errors via the Engine.
func (s *service) Healthy() (err error) {
	err = s.StateMachine.Healthy()
	if err == nil {
		err = s.eng.clearCond()
	}
	return
}

func (s *service) HealthReport() map[string]error {
	m := map[string]error{s.Name(): s.Healthy()}
	for _, sub := range s.subs {
		CopyHealth(m, sub.HealthReport())
	}
	return m
}

func (s *service) Name() string { return s.eng.SugaredLogger.Name() }

func (s *service) Start(ctx context.Context) error {
	return s.StartOnce(s.cfg.Name, func() error {
		s.eng.Info("Starting")
		if len(s.subs) > 0 {
			var ms MultiStart
			s.eng.Infof("Starting %d sub-services", len(s.subs))
			for _, sub := range s.subs {
				if err := ms.Start(ctx, sub); err != nil {
					s.eng.Errorw("Failed to start sub-service", "error", err)
					return fmt.Errorf("failed to start sub-service of %s: %w", s.cfg.Name, err)
				}
			}
		}
		if s.cfg.Start != nil {
			if err := s.cfg.Start(ctx); err != nil {
				return fmt.Errorf("failed to start service %s: %w", s.cfg.Name, err)
			}
		}
		s.eng.Info("Started")
		return nil
	})
}

func (s *service) Close() error {
	return s.StopOnce(s.cfg.Name, func() (err error) {
		s.eng.Info("Closing")
		defer s.eng.Info("Closed")

		close(s.eng.StopChan)
		s.eng.wg.Wait()

		if s.cfg.Close != nil {
			err = s.cfg.Close()
		}

		if len(s.subs) == 0 {
			return
		}
		s.eng.Infof("Closing %d sub-services", len(s.subs))
		err = errors.Join(err, MultiCloser(s.subs).Close())
		return
	})
}
