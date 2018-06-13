package adapters

import (
	"fmt"
	"time"

	"github.com/smartcontractkit/chainlink/store"
	"github.com/smartcontractkit/chainlink/store/models"
)

// Sleep adapter allows a job to do nothing for some amount of wall time.
type Sleep struct {
	Seconds int `json:"seconds"`
}

const maxDuration = 72 * time.Hour

// Perform returns the empty RunResult
func (adaptor *Sleep) Perform(input models.RunResult, _ *store.Store) models.RunResult {
	duration := time.Duration(adaptor.Seconds) * time.Second
	if duration > maxDuration {
		return input.WithError(fmt.Errorf("Duration %s exceeds maximum of %s", duration, maxDuration))
	}

	time.Sleep(duration)
	input.Status = models.RunStatusCompleted
	return input
}
