package adapters

import (
	"github.com/smartcontractkit/chainlink/logger"
	"github.com/smartcontractkit/chainlink/store"
	"github.com/smartcontractkit/chainlink/store/models"
)

// Sleep adapter allows a job to do nothing for some amount of wall time.
type Sleep struct {
	Until models.Time `json:"until"`
}

// Perform returns the input RunResult after waiting for the specified Until parameter.
func (adapter *Sleep) Perform(input models.RunResult, str *store.Store) models.RunResult {
	duration := adapter.Until.DurationFromNow()
	if duration <= 0 {
		input.Status = models.RunStatusCompleted
		return input
	}

	input.Status = models.RunStatusPendingSleep
	go func() {
		<-str.Clock.After(duration)
		if err := str.RunChannel.Send(input.JobRunID, nil); err != nil {
			logger.Error("Sleep Adapter Perform:", err.Error())
		}
	}()

	return input
}
