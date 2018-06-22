package adapters

import (
	"fmt"
	"time"

	"github.com/smartcontractkit/chainlink/store"
	"github.com/smartcontractkit/chainlink/store/models"
)

// Sleep adapter allows a job to do nothing for some amount of wall time.
type Sleep struct {
	EndAt models.Time `json:"until"`
}

const maxDuration = 30 * 24 * time.Hour

// Perform returns the input RunResult after waiting for the specified EndAt time.
func (adapter *Sleep) Perform(input models.RunResult, store *store.Store) models.RunResult {
	duration := adapter.EndAt.DurationFromNow()
	if duration <= 0 {
		input.Status = models.RunStatusCompleted
		return input
	}
	if duration > maxDuration {
		return input.WithError(fmt.Errorf("Sleep Adapter: %v is greater than max duration %v", duration, maxDuration))
	}

	input.Status = models.RunStatusPendingSleep
	go func() {
		<-store.Clock.After(duration)
		store.RunChannel <- input
	}()

	return input
}
