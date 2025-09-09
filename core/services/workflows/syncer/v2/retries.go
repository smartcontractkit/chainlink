package v2

import (
	"math"
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
