package fluxmonitorv2

import (
	"time"

	"github.com/smartcontractkit/chainlink/core/utils"
)

// IdleTimer defines a ResettableTimer which can be disabled
type IdleTimer struct {
	timer    utils.ResettableTimer
	period   time.Duration
	disabled bool
}

// NewIdleTimer constructs a new IdleTimer
func NewIdleTimer(period time.Duration, disabled bool) *IdleTimer {
	return &IdleTimer{
		timer:    utils.NewResettableTimer(),
		period:   period,
		disabled: disabled,
	}
}

// Period gets the timer period
func (t *IdleTimer) Period() time.Duration {
	return t.period
}

// IsEnabled determines if the timer is enabled
func (t *IdleTimer) IsEnabled() bool {
	return !t.disabled
}

// IsDisabled determines if the timer is disabled
func (t *IdleTimer) IsDisabled() bool {
	return t.disabled
}

// Ticks ticks on a given interval
func (t *IdleTimer) Ticks() <-chan time.Time {
	return t.timer.Ticks()
}

// Reset resets the timer
func (t *IdleTimer) Reset(period time.Duration) {
	t.timer.Reset(period)
}

// Stop stops the timer permanently
func (t *IdleTimer) Stop() {
	t.timer.Stop()
}
