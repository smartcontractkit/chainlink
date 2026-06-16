package v2

import (
	"math"
	"sync"
	"time"

	"github.com/jonboulle/clockwork"
)

type reconciliationEvent struct {
	Event
	id          string
	signature   string
	nextRetryAt time.Time
	retryCount  int
}

func (r *reconciliationEvent) updateNextRetryFor(clock clockwork.Clock, retryInterval time.Duration, maxRetryInterval time.Duration) {
	r.retryCount++
	nextRetry := math.Pow(2, float64(r.retryCount)) * float64(retryInterval)
	nextRetry = math.Min(float64(maxRetryInterval), nextRetry)
	r.nextRetryAt = clock.Now().Add(time.Duration(nextRetry))
}

func activationRetriesExhausted(retryCount, maxActivationRetries int) bool {
	return maxActivationRetries > 0 && retryCount >= maxActivationRetries
}

type droppedActivations struct {
	mu       sync.RWMutex
	bySource map[string]map[string]string // source -> workflowID -> signature
}

func newDroppedActivations() *droppedActivations {
	return &droppedActivations{
		bySource: make(map[string]map[string]string),
	}
}

func (d *droppedActivations) isDropped(source, workflowID, signature string) bool {
	d.mu.RLock()
	defer d.mu.RUnlock()
	sig, ok := d.bySource[source][workflowID]
	return ok && sig == signature
}

func (d *droppedActivations) drop(source, workflowID, signature string) {
	d.mu.Lock()
	defer d.mu.Unlock()
	if d.bySource[source] == nil {
		d.bySource[source] = make(map[string]string)
	}
	d.bySource[source][workflowID] = signature
}

func (d *droppedActivations) clear(source, workflowID string) {
	d.mu.Lock()
	defer d.mu.Unlock()
	if m := d.bySource[source]; m != nil {
		delete(m, workflowID)
	}
}

type reconcileReport struct {
	// events is a map of event type to the number of events of that type
	NumEventsByType map[string]int
	// id -> nextRetry time
	Backoffs map[string]time.Time
}

func newReconcileReport() *reconcileReport {
	return &reconcileReport{
		NumEventsByType: map[string]int{},
		Backoffs:        map[string]time.Time{},
	}
}
