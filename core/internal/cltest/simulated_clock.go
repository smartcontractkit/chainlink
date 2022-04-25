package cltest

import "time"

// SimulatedClock is used for tests.
// It wraps the time.Now() by setting a fixed time, controlled by tests.
// To change the time, tests should ONLY use the FastForwardBy() function.
// Tests SHOULD NOT directly set the currentTime field, as that can
// have unintended consequences, if the clock in moved back in time.
type SimulatedClock struct {
	currentTime time.Time
}

// Returns a new SimulatedClock set at the provided time.
func NewSimulatedClock(currentTime time.Time) SimulatedClock {
	return SimulatedClock{currentTime: currentTime}
}

// Now returns the current time on this clock.
func (c SimulatedClock) Now() time.Time {
	return c.currentTime
}

// Moves the clock ahead by provided duration
func (c SimulatedClock) FastForwardBy(duration time.Duration) time.Time {
	c.currentTime = c.currentTime.Add(duration)
	return c.currentTime
}
