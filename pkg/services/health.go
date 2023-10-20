package services

import (
	"errors"
	"fmt"
	"maps"
	"runtime/debug"
	"sync"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// CopyHealth copies health statuses from src to dest. Useful when implementing HealthReporter.HealthReport.
// If duplicate names are encountered, the errors are joined, unless testing in which case a panic is thrown.
func CopyHealth(dest, src map[string]error) {
	for name, err := range src {
		errOrig, ok := dest[name]
		if ok {
			if testing.Testing() {
				panic("service names must be unique: duplicate name: " + name)
			}
			if errOrig != nil {
				dest[name] = errors.Join(errOrig, err)
				continue
			}
		}
		dest[name] = err
	}
}

// HealthChecker is a services.Service which monitors other services and can be probed for system health.
type HealthChecker struct {
	StateMachine
	chStop chan struct{}
	chDone chan struct{}

	servicesMu sync.RWMutex
	services   map[string]HealthReporter

	stateMu sync.RWMutex
	healthy map[string]error
	ready   map[string]error
}

const interval = 15 * time.Second

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
	version = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "version",
			Help: "Application version information",
		},
		[]string{"version", "commit"},
	)
)

func NewChecker() *HealthChecker {
	return &HealthChecker{
		services: make(map[string]HealthReporter, 10),
		healthy:  make(map[string]error, 10),
		ready:    make(map[string]error, 10),
		chStop:   make(chan struct{}),
		chDone:   make(chan struct{}),
	}
}

func (c *HealthChecker) Start() error {
	return c.StartOnce("HealthCheck", func() error {
		if bi, ok := debug.ReadBuildInfo(); ok {
			hash := bi.Main.Sum
			if len(hash) > 7 {
				hash = hash[:7]
			}
			version.WithLabelValues(bi.Main.Version, hash).Inc()
		}

		// update immediately
		c.update()

		go c.run()

		return nil
	})
}

func (c *HealthChecker) Close() error {
	return c.StopOnce("HealthCheck", func() error {
		close(c.chStop)
		<-c.chDone
		return nil
	})
}

func (c *HealthChecker) run() {
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

func (c *HealthChecker) update() {
	// copy services into a new map to avoid lock contention while doing checks
	c.servicesMu.RLock()
	l := len(c.services)
	services := make(map[string]HealthReporter, l)
	maps.Copy(services, c.services)
	c.servicesMu.RUnlock()

	ready := make(map[string]error, l)
	healthy := make(map[string]error, l)

	// now, do all the checks
	for name, s := range services {
		ready[name] = s.Ready()
		for n, err := range s.HealthReport() {
			healthy[n] = err
			value := 0
			if err == nil {
				value = 1
			}

			// report metrics to prometheus
			healthStatus.WithLabelValues(name).Set(float64(value))
		}
	}
	uptimeSeconds.Add(interval.Seconds())

	// save state
	c.stateMu.Lock()
	defer c.stateMu.Unlock()
	maps.Copy(c.ready, ready)
	maps.Copy(c.healthy, healthy)
}

// Register a service for health checks.
func (c *HealthChecker) Register(service HealthReporter) error {
	name := service.Name()
	if name == "" {
		return fmt.Errorf("misconfigured check %#v for %v", name, service)
	}

	c.servicesMu.Lock()
	defer c.servicesMu.Unlock()
	if testing.Testing() {
		if orig, ok := c.services[name]; ok {
			panic(fmt.Errorf("duplicate name %q: service names must be unique: types %T & %T", name, service, orig))
		}
	}
	c.services[name] = service
	return nil
}

// Unregister a service.
func (c *HealthChecker) Unregister(name string) error {
	if name == "" {
		return fmt.Errorf("name cannot be empty")
	}

	c.servicesMu.Lock()
	defer c.servicesMu.Unlock()
	delete(c.services, name)
	healthStatus.DeleteLabelValues(name)
	return nil
}

// IsReady returns the current readiness of the system.
// A system is considered ready if all checks are passing (no errors)
func (c *HealthChecker) IsReady() (ready bool, errors map[string]error) {
	c.stateMu.RLock()
	defer c.stateMu.RUnlock()

	ready = true
	errors = make(map[string]error, len(c.ready))

	for name, state := range c.ready {
		errors[name] = state

		if state != nil {
			ready = false
		}
	}

	return
}

// IsHealthy returns the current health of the system.
// A system is considered healthy if all checks are passing (no errors)
func (c *HealthChecker) IsHealthy() (healthy bool, errors map[string]error) {
	c.stateMu.RLock()
	defer c.stateMu.RUnlock()

	healthy = true
	errors = make(map[string]error, len(c.healthy))

	for name, state := range c.healthy {
		errors[name] = state

		if state != nil {
			healthy = false
		}
	}

	return
}
