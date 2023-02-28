package common

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/smartcontractkit/chainlink/core/utils"
)

// HealthError is type to track error for a given system
type HealthError struct {
	System    string
	Err       error
	createdAt time.Time
}

func NewHealthError(systemName string, err error) *HealthError {
	return &HealthError{
		System: systemName,
		Err:    err,
	}
}

// TODO - name, maybe ErrorCache is better...
// HealthReporter is TTL cache for HealthErrors
type HealthReporter struct {
	mu    sync.RWMutex
	cache map[string][]*HealthError

	wg               sync.WaitGroup
	quitCh           chan struct{}
	lookBackDuration time.Duration

	utils.StartStopOnce
}

// NewHealtherReport creates a HealthReporter with a given lookup back duration
func NewHealthReporter(lookBack time.Duration) *HealthReporter {
	return &HealthReporter{
		cache:            make(map[string][]*HealthError),
		quitCh:           make(chan struct{}),
		lookBackDuration: lookBack,
	}
}

func (h *HealthReporter) Start(ctx context.Context) {
	h.StartOnce("health-reporter", func() error {

		h.wg.Add(1)
		go h.enforceTTL(ctx)
		return nil
	})
}

// Report coasceleces multiple errors per system over the lookback duration
func (h *HealthReporter) Report() map[string]error {
	result := make(map[string]error)
	h.mu.RLock()
	defer h.mu.RUnlock()
	for name, herrs := range h.cache {
		errs := make([]error, len(herrs))

		for i, herr := range herrs {
			errs[i] = herr.Err
		}

		result[name] = errors.Join(errs...)
	}

	return result
}

// Add puts an error into the cache
func (h *HealthReporter) Add(herr *HealthError) {
	h.mu.Lock()
	defer h.mu.Unlock()
	herr.createdAt = time.Now()

	val, exists := h.cache[herr.System]
	if !exists {
		val = make([]*HealthError, 0)
	}
	val = append(val, herr)
	h.cache[herr.System] = val
}

// prune drops errors created outside the lookback duration
func (h *HealthReporter) prune() {
	h.mu.Lock()
	defer h.mu.Unlock()

	lookupBack := time.Now().Add(-1 * h.lookBackDuration)
	// by construction the elements in the cache are time ordered
	for name, herrs := range h.cache {
		dropCnt := 0
		for _, herr := range herrs {
			if herr.createdAt.Before(lookupBack) {
				dropCnt++
			} else {
				break
			}
		}
		h.cache[name] = h.cache[name][dropCnt:]
	}
}

// enforceTTL prunes the cache after every lookup back period
func (h *HealthReporter) enforceTTL(ctx context.Context) {
	defer h.wg.Done()

	ttlTicker := time.NewTicker(h.lookBackDuration)
	defer ttlTicker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-h.quitCh:
			return
		case <-ttlTicker.C:
			h.prune()
		}
	}

}

// Close releases all resources and blocks until they are released
func (h *HealthReporter) Close() error {
	return h.StopOnce("health-reporter", func() error {
		close(h.quitCh)
		h.wg.Wait()
		return nil
	})

}
