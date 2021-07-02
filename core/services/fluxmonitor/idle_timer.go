package fluxmonitor

import (
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

	idleTimer utils.ResettableTimer
}

func NewIdleTimer(idleDuration time.Duration) *IdleTimer {
	return &IdleTimer{
		idleDuration: idleDuration,
		idleTimer:    utils.NewResettableTimer(),
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

// Stop
func (t *IdleTimer) Stop() {
	t.idleTimer.Stop()
}

// reset resets the idle timer to fire after the idle duration.
func (t *IdleTimer) resetWithIdleDuration() {
	t.idleTimer.Reset(t.idleDuration)
}

// reset resets the idle timer to fire when the time until the idle deadline
// is reached.
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
}

func (t *IdleTimer) Ticks() <-chan time.Time {
	return t.idleTimer.Ticks()
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

// func (p *PollingDeviationChecker) resetIdleTimer(roundStartedAtUTC uint64) {
// 	if p.isHibernating || p.initr.IdleTimer.Disabled {
// 		p.idleTimer.Stop()
// 		return
// 	} else if roundStartedAtUTC == 0 {
// 		// There is no active round, so keep using the idleTimer we already have
// 		return
// 	}

// 	startedAt := time.Unix(int64(roundStartedAtUTC), 0)
// 	idleDeadline := startedAt.Add(p.initr.IdleTimer.Duration.Duration())
// 	timeUntilIdleDeadline := time.Until(idleDeadline)
// 	loggerFields := p.loggerFields(
// 		"startedAt", roundStartedAtUTC,
// 		"timeUntilIdleDeadline", timeUntilIdleDeadline,
// 	)

// 	if timeUntilIdleDeadline <= 0 {
// 		logger.Debugw("not resetting idleTimer, negative duration", loggerFields...)
// 		return
// 	}
// 	p.idleTimer.Reset(timeUntilIdleDeadline)
// 	logger.Debugw("resetting idleTimer", loggerFields...)
// }

// // Reset stops a ResettableTimer
// // and resets it with a new duration
// func (t *ResettableTimer) Reset(duration time.Duration) {
// 	t.mu.Lock()
// 	defer t.mu.Unlock()
// 	if t.timer != nil {
// 		t.timer.Stop()
// 	}
// 	t.timer = time.NewTimer(duration)
// }
