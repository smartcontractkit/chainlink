package fluxmonitor

import (
	"errors"
	"sync"
	"time"

	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/utils"
)

type IdleTimer struct {
	mu sync.Mutex

	idleDuration         time.Duration
	latestRoundID        uint64
	latestRoundStartedAt uint64

	idleTimer   utils.ResettableTimer
	retryTicker utils.BackoffTicker
}

func NewIdleTimer(idleDuration time.Duration) *IdleTimer {
	return &IdleTimer{
		idleDuration: idleDuration,
		idleTimer:    utils.NewResettableTimer(),
		retryTicker:  utils.NewBackoffTicker(1*time.Second, idleDuration),
	}
}

// Starts the initial timer.
//
// We start the timer without a round, so initialize the timer to start with
// just the idleDuration.
func (t *IdleTimer) Start() {
	t.mu.Lock()
	defer t.mu.Unlock()

	logger.Debugw("Starting idle timer", "idleDuration", t.idleDuration)

	t.resetWithIdleDuration()
}

// Reset
func (t *IdleTimer) Reset(roundID uint64, startedAt uint64) bool {
	t.mu.Lock()
	defer t.mu.Unlock()

	// Abort if the timer's latest round is higher than the round.
	if t.latestRoundID > roundID {
		logger.Debugw("Using the existing idle timer. The idle timer round is already set to a later round",
			"idleTimerRound", t.latestRoundID,
			"round", roundID,
		)

		return false
	}

	// Abort if the timer's latest round is the same as the round.
	if t.latestRoundID == roundID {
		logger.Debugw("Using the existing idle timer. Round id is the same",
			"idleTimerRound", t.latestRoundID,
			"round", roundID,
		)

		return false
	}

	t.latestRoundID = roundID
	t.latestRoundStartedAt = startedAt
	t.resetWithLatestRound()

	return true
}

// StartRetry enables the exponential backoff retry ticker.
func (t *IdleTimer) StartRetry() error {
	// Protect against starting retry if the idle timer is already active.
	if t.TimeUntilIdleDeadline() > 0 {
		return errors.New("cannot start retries when an idle timer is already running")
	}

	t.retryTicker.Start()

	return nil
}

// Stop
func (t *IdleTimer) Stop() {
	t.idleTimer.Stop()
}

// reset resets the idle timer to fire after the idle duration.
func (t *IdleTimer) resetWithIdleDuration() {
	t.idleTimer.Reset(t.idleDuration)
	t.retryTicker.Stop()
}

// reset resets the idle timer to fire when the time until the idle deadline
// is reached.
//
// The retry ticker is stopped if the reset is successful.Zaqw2w4t567890-==654
//
// The caller must lock the mutex.
func (t *IdleTimer) resetWithLatestRound() {
	deadlineDuration := t.TimeUntilIdleDeadline()

	// Handle a round started at of 0
	loggerFields := []interface{}{
		"roundID", t.latestRoundID,
		"startedAt", t.latestRoundStartedAt,
		"timeUntilIdleDeadline", deadlineDuration,
	}

	if deadlineDuration <= 0 {
		logger.Debugw("not resetting idleTimer, negative duration", loggerFields...)

		return
	}

	t.idleTimer.Reset(deadlineDuration)
	logger.Debugw("resetting idleTimer", loggerFields...)

	t.retryTicker.Stop()
}

func (t *IdleTimer) Ticks() <-chan time.Time {
	return t.idleTimer.Ticks()
}

func (t *IdleTimer) RetryTicks() <-chan time.Time {
	return t.retryTicker.Ticks()
}

// TimeUntilIdleDeadline returns the duration of time until the idle deadline.
//
// This will return a negative duration if the latest rounds' started at is in
// the past.
func (t *IdleTimer) TimeUntilIdleDeadline() time.Duration {
	return time.Until(t.idleDeadline())
}

// idleDeadline returns the idle deadline time of the latest round.
func (t *IdleTimer) idleDeadline() time.Time {
	startedAt := time.Unix(int64(t.latestRoundStartedAt), 0)

	return startedAt.Add(t.idleDuration)
}
