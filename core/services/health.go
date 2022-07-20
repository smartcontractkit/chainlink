package services

import (
	"sync"
	"time"

	"github.com/pkg/errors"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/smartcontractkit/chainlink/core/static"
	"github.com/smartcontractkit/chainlink/core/utils"
)

//go:generate mockery --name Checker --output ./mocks/ --case=underscore
type (
	// Checker provides a service which can be probed for system health.
	Checker interface {
		// Register a service for health checks.
		Register(name string, service Checkable) error
		// Unregister a service.
		Unregister(name string) error
		// IsReady returns the current readiness of the system.
		// A system is considered ready if all checks are passing (no errors)
		IsReady() (ready bool, errors map[string]error)
		// IsHealthy returns the current health of the system.
		// A system is considered healthy if all checks are passing (no errors)
		IsHealthy() (healthy bool, errors map[string]error)

		Start() error
		Close() error
	}

	checker struct {
		srvMutex   sync.RWMutex
		services   map[string]Checkable
		stateMutex sync.RWMutex
		state      map[string]State

		chStop chan struct{}
		chDone chan struct{}

		utils.StartStopOnce
	}

	State struct {
		ready   error
		healthy error
	}

	Status string
)

var _ Checker = (*checker)(nil)

const (
	StatusPassing Status = "passing"
	StatusFailing Status = "failing"

	interval = 15 * time.Second
)

var (
	healthStatus = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "health",
			Help: "Health status by service",
		},
		[]string{"service_id"},
	)
	uptimeSeconds = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "uptime_seconds",
			Help: "Uptime of the application measured in seconds",
		},
	)
	nodeVersion = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "version",
			Help: "Node version information",
		},
		[]string{"version", "commit"},
	)
)

func NewChecker() Checker {
	c := &checker{
		services: make(map[string]Checkable, 10),
		state:    make(map[string]State, 10),
		chStop:   make(chan struct{}),
		chDone:   make(chan struct{}),
	}

	return c
}

func (c *checker) Start() error {
	return c.StartOnce("HealthCheck", func() error {
		nodeVersion.WithLabelValues(static.Version, static.Sha).Inc()

		// update immediately
		c.update()

		go c.run()

		return nil
	})
}

func (c *checker) Close() error {
	return c.StopOnce("HealthCheck", func() error {
		close(c.chStop)
		<-c.chDone
		return nil
	})
}

func (c *checker) run() {
	defer close(c.chDone)

	ticker := time.NewTicker(interval)

	for {
		select {
		case <-ticker.C:
			c.update()
		case <-c.chStop:
			return
		}
	}

}

func (c *checker) update() {
	state := make(map[string]State, len(c.services))

	c.srvMutex.RLock()
	// copy services into a new map to avoid lock contention while doing checks
	services := make(map[string]Checkable, len(c.services))
	for name, s := range c.services {
		services[name] = s
	}
	c.srvMutex.RUnlock()

	// now, do all the checks
	for name, s := range services {
		ready := s.Ready()
		healthy := s.Healthy()

		state[name] = State{ready, healthy}
	}

	// we use a separate lock to avoid holding the lock over state while talking
	// to services
	c.stateMutex.Lock()
	defer c.stateMutex.Unlock()

	for name, state := range state {
		c.state[name] = state

		value := 0
		if state.healthy == nil {
			value = 1
		}

		// report metrics to prometheus
		healthStatus.WithLabelValues(name).Set(float64(value))
	}
	uptimeSeconds.Add(interval.Seconds())
}

func (c *checker) Register(name string, service Checkable) error {
	if service == nil || name == "" {
		return errors.Errorf("misconfigured check %#v for %v", name, service)
	}

	c.srvMutex.Lock()
	defer c.srvMutex.Unlock()
	c.services[name] = service
	return nil
}

func (c *checker) Unregister(name string) error {
	if name == "" {
		return errors.Errorf("name cannot be empty")
	}

	c.srvMutex.Lock()
	defer c.srvMutex.Unlock()
	delete(c.services, name)
	healthStatus.DeleteLabelValues(name)
	return nil
}

func (c *checker) IsReady() (ready bool, errors map[string]error) {
	c.stateMutex.RLock()
	defer c.stateMutex.RUnlock()

	ready = true
	errors = make(map[string]error, len(c.services))

	for name, state := range c.state {
		errors[name] = state.ready

		if state.ready != nil {
			ready = false
		}
	}

	return
}

func (c *checker) IsHealthy() (healthy bool, errors map[string]error) {
	c.stateMutex.RLock()
	defer c.stateMutex.RUnlock()

	healthy = true
	errors = make(map[string]error, len(c.services))

	for name, state := range c.state {
		errors[name] = state.healthy

		if state.healthy != nil {
			healthy = false
		}
	}

	return
}
