package adapters

import (
	"time"

	"chainlink/core/logger"
	"chainlink/core/store"
	"chainlink/core/store/models"
	"chainlink/core/utils"
)

// Sleep adapter allows a job to do nothing for some amount of wall time.
type Sleep struct {
	Until models.AnyTime `json:"until"`
}

// Perform returns the input RunResult after waiting for the specified Until parameter.
func (adapter *Sleep) Perform(input models.RunInput, str *store.Store) models.RunOutput {
	duration := adapter.Duration()
	if duration > 0 {
		logger.Debugw("Task sleeping...", "duration", duration)
		<-str.Clock.After(duration)
	}

	return models.NewRunOutputComplete(models.JSON{})
}

// Duration returns the amount of sleeping this task should be paused for.
func (adapter *Sleep) Duration() time.Duration {
	return utils.DurationFromNow(adapter.Until.Time)
}
