package adapters

import (
	"time"

	"chainlink/core/store"
	"chainlink/core/store/models"
	"chainlink/core/utils"
)

// Sleep adapter allows a job to do nothing for some amount of wall time.
type Sleep struct {
	Until models.AnyTime `json:"until"`
}

// Perform returns the input RunResult after waiting for the specified Until parameter.
func (adapter *Sleep) Perform(input models.RunResult, str *store.Store) models.RunResult {
	input.Status = models.RunStatusPendingSleep
	return input
}

// Duration returns the amount of sleeping this task should be paused for.
func (adapter *Sleep) Duration() time.Duration {
	return utils.DurationFromNow(adapter.Until.Time)
}
