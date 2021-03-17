package fluxmonitorv2

import (
	"time"

	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/flux_aggregator_wrapper"
)

type PollManagerConfig struct {
	IsHibernating      bool
	PollTickerInterval time.Duration
	PollTickerDisabled bool
	IdleTimerInterval  time.Duration
	IdleTimerDisabled  bool
}

// PollManager manages the tickers/timers which cause the Flux Monitor to start
// a poll
type PollManager struct {
	cfg PollManagerConfig

	isHibernating bool
	pollTicker    *PollTicker
	idleTimer     *IdleTimer
}

// NewPollManager initializes a new PollManager
func NewPollManager(cfg PollManagerConfig) *PollManager {
	return &PollManager{
		cfg:           cfg,
		isHibernating: cfg.IsHibernating,

		pollTicker: NewPollTicker(cfg.PollTickerInterval, cfg.PollTickerDisabled),
		idleTimer:  NewIdleTimer(cfg.IdleTimerInterval, cfg.IdleTimerDisabled),
	}
}

// PollTickerTicks ticks on a given interval
func (pm *PollManager) PollTickerTicks() <-chan time.Time {
	return pm.pollTicker.Ticks()
}

// IdleTimerTicks ticks after a given period
func (pm *PollManager) IdleTimerTicks() <-chan time.Time {
	return pm.idleTimer.Ticks()
}

func (pm *PollManager) Stop() {
	pm.pollTicker.Stop()
	pm.idleTimer.Stop()
}

// Reset resets all tickers/timers
func (pm *PollManager) Reset(roundState flux_aggregator_wrapper.OracleRoundState) {
	pm.ResetPollTicker()
	pm.ResetIdleTimer(roundStateTimesOutAt(roundState))
}

// ResetPollTicker resets the poll ticker if enabled and not hibernating
func (pm *PollManager) ResetPollTicker() {
	if pm.pollTicker.IsEnabled() && !pm.isHibernating {
		pm.pollTicker.Resume()
	} else {
		pm.pollTicker.Pause()
	}
}

func (pm *PollManager) IsPollTickerDisabled() bool {
	return pm.pollTicker.IsDisabled()
}

func (pm *PollManager) IsIdleTimerDisabled() bool {
	return pm.idleTimer.IsDisabled()
}

func (pm *PollManager) ResetIdleTimer(roundStartedAtUTC uint64) {
	// Stop the timer if hibernating or disabled
	if pm.isHibernating || pm.idleTimer.IsDisabled() {
		pm.idleTimer.Stop()

		return
	}

	// There is no active round, so keep using the idleTimer we already have
	if roundStartedAtUTC == 0 {
		return
	}

	startedAt := time.Unix(int64(roundStartedAtUTC), 0)
	idleDeadline := startedAt.Add(pm.idleTimer.Period())
	timeUntilIdleDeadline := time.Until(idleDeadline)

	// loggerFields := fm.loggerFields(
	// 	"startedAt", roundStartedAtUTC,
	// 	"timeUntilIdleDeadline", timeUntilIdleDeadline,
	// )

	if timeUntilIdleDeadline <= 0 {
		// fm.logger.Debugw("not resetting idleTimer, negative duration", loggerFields...)

		return
	}

	pm.idleTimer.Reset(timeUntilIdleDeadline)
	// 	fm.logger.Debugw("resetting idleTimer", loggerFields...)
}

// func (fm *FluxMonitor) resetIdleTimer(roundStartedAtUTC uint64) {
// 	if fm.isHibernating || fm.idleTimer.IsDisabled() {
// 		fm.idleTimer.Stop()
// 		return
// 	} else if roundStartedAtUTC == 0 {
// 		// There is no active round, so keep using the idleTimer we already have
// 		return
// 	}

// 	startedAt := time.Unix(int64(roundStartedAtUTC), 0)
// 	idleDeadline := startedAt.Add(fm.idleTimer.Period())
// 	timeUntilIdleDeadline := time.Until(idleDeadline)
// 	loggerFields := fm.loggerFields(
// 		"startedAt", roundStartedAtUTC,
// 		"timeUntilIdleDeadline", timeUntilIdleDeadline,
// 	)

// 	if timeUntilIdleDeadline <= 0 {
// 		fm.logger.Debugw("not resetting idleTimer, negative duration", loggerFields...)
// 		return
// 	}
// 	fm.idleTimer.Reset(timeUntilIdleDeadline)
// 	fm.logger.Debugw("resetting idleTimer", loggerFields...)
// }

// Hibernate sets hibernation to true and resets all ticker/timers
func (pm *PollManager) Hibernate(roundState flux_aggregator_wrapper.OracleRoundState) {
	pm.isHibernating = true
	pm.Reset(roundState)
}

// Awaken sets hibernation to false and resets all ticker/timers
func (pm *PollManager) Awaken(roundState flux_aggregator_wrapper.OracleRoundState) {
	pm.isHibernating = false
	pm.Reset(roundState)
}
