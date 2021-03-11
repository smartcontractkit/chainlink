package fluxmonitorv2

import (
	"time"

	"github.com/smartcontractkit/chainlink/core/utils"
)

// PollTicker defines a PausableTicker which can be disabled
type PollTicker struct {
	ticker   utils.PausableTicker
	interval time.Duration
	disabled bool
}

// NewPollTicker constructs a new PollTicker
func NewPollTicker(interval time.Duration, disabled bool) *PollTicker {
	return &PollTicker{
		ticker:   utils.NewPausableTicker(interval),
		interval: interval,
		disabled: disabled,
	}
}

// Interval gets the ticker interval
func (t *PollTicker) Interval() time.Duration {
	return t.interval
}

// IsEnabled determines if the picker is enabled
func (t *PollTicker) IsEnabled() bool {
	return !t.disabled
}

// IsDisabled determines if the picker is disabled
func (t *PollTicker) IsDisabled() bool {
	return t.disabled
}

// Ticks ticks on a given interval
func (t *PollTicker) Ticks() <-chan time.Time {
	return t.ticker.Ticks()
}

// Resume resumes the ticker if it is enabled
func (t *PollTicker) Resume() {
	if t.IsEnabled() {
		t.ticker.Resume()
	}
}

// Pause pauses the ticker
func (t *PollTicker) Pause() {
	t.ticker.Pause()
}

// Stop stops the ticker permanently
func (t *PollTicker) Stop() {
	t.ticker.Destroy()
}
