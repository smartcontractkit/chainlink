package adapters

import (
	"github.com/smartcontractkit/chainlink/store"
	"github.com/smartcontractkit/chainlink/store/models"
)

// Sleep adapter allows a job to do nothing for some amount of wall time.
type Sleep struct {
	EndAt models.Time `json:"until"`
}

// Perform returns the input RunResult after waiting for the specified EndAt time.
func (adapter *Sleep) Perform(input models.RunResult, store *store.Store) models.RunResult {
	duration := adapter.EndAt.DurationFromNow()
	if duration <= 0 {
		input.Status = models.RunStatusCompleted
		return input
	}

	input.Status = models.RunStatusPendingSleep
	go func() {
		<-store.Clock.After(duration)
		store.RunManager.Queue <- input
	}()

	return input
}
