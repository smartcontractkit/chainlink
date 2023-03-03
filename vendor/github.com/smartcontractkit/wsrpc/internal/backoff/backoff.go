// Package backoff implement the backoff strategy for wsrpc.
package backoff

import (
	"sync"
	"time"

	"github.com/cenkalti/backoff"
)

type Strategy interface {
	// NextBackOff returns the duration to wait before retrying the operation,
	// or backoff.
	NextBackOff() time.Duration

	// Reset to initial state.
	Reset()
}

// Config defines the configuration options for backoff.
type Config struct {
	// BaseDelay is the amount of time to backoff after the first failure.
	BaseDelay time.Duration
	// Multiplier is the factor with which to multiply backoffs after a
	// failed retry. Should ideally be greater than 1.
	Multiplier float64
	// Jitter is the factor with which backoffs are randomized.
	Jitter float64
	// MaxDelay is the upper bound of backoff delay.
	MaxDelay time.Duration
}

var defaultConfig = Config{
	BaseDelay:  1 * time.Second,
	Multiplier: 1.6,
	Jitter:     0.2,
	MaxDelay:   120 * time.Second,
}

// DefaultExponential is an exponential backoff implementation using the
// default values.
var DefaultExponential = NewExponential(defaultConfig)

// Exponential implements an exponential backoff algorithm.
type Exponential struct {
	mu sync.RWMutex
	backoff.BackOff
}

// NewExponential constructs a new Exponential Strategy from the configuration
// options.
func NewExponential(config Config) *Exponential {
	boff := &backoff.ExponentialBackOff{
		InitialInterval:     config.BaseDelay,
		Multiplier:          config.Multiplier,
		RandomizationFactor: config.Jitter,
		MaxInterval:         config.MaxDelay,
		MaxElapsedTime:      0, // Never stop
		Clock:               backoff.SystemClock,
	}
	// We have to reset to set the initial values
	boff.Reset()

	return &Exponential{
		BackOff: boff,
	}
}

// NextBackOff returns the amount of time to wait before the next retry.
func (es *Exponential) NextBackOff() time.Duration {
	es.mu.Lock()
	defer es.mu.Unlock()

	return es.BackOff.NextBackOff()
}

// Reset the interval back to the initial retry interval.
func (es *Exponential) Reset() {
	es.mu.Lock()
	defer es.mu.Unlock()

	es.BackOff.Reset()
}
