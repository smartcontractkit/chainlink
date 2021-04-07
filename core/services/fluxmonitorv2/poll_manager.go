package fluxmonitorv2

import (
	"time"

	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/flux_aggregator_wrapper"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/utils"
)

type PollManagerConfig struct {
	IsHibernating           bool
	PollTickerInterval      time.Duration
	PollTickerDisabled      bool
	IdleTimerPeriod         time.Duration
	IdleTimerDisabled       bool
	HibernationPollPeriod   time.Duration
	MinRetryBackoffDuration time.Duration
	MaxRetryBackoffDuration time.Duration
}

// PollManager manages the tickers/timers which cause the Flux Monitor to start
// a poll. It contains 4 types of tickers and timers which determine when to
// initiate a poll
//
// HibernationTimer - The PollManager can be set to hibernate, which disables all
// other ticker/timers, and enables the hibernation timer. Upon expiry of the
// hibernation timer, a poll is requested. When the PollManager is awakened, the
// other tickers and timers are enabled with the current round state, and the
// hibernation timer is disabled.
//
// PollTicker - The poll ticker requests a poll at a given interval defined in
// PollManagerConfig. Disabling this through config will permanently disable
// the ticker, even through a resets.
//
// IdleTimer - The idle timer requests a poll after no poll has taken place
// since the last round was start and the IdleTimerPeriod has elapsed. This can
// also be known as a heartbeat.
//
// RoundTimer - The round timer requests a poll when the round state provided by
// the contract has timed out.
//
// RetryTicker - The retry ticker requests a poll with a backoff duration. This
// is started when the idle timer fails, and will poll with a maximum backoff
// of either 1 hour or the idle timer period if it is lower
type PollManager struct {
	cfg PollManagerConfig

	hibernationTimer utils.ResettableTimer
	pollTicker       utils.PausableTicker
	idleTimer        utils.ResettableTimer
	roundTimer       utils.ResettableTimer
	retryTicker      utils.BackoffTicker

	logger *logger.Logger
}

// NewPollManager initializes a new PollManager
func NewPollManager(cfg PollManagerConfig, logger *logger.Logger) *PollManager {
	minBackoffDuration := cfg.MinRetryBackoffDuration
	if cfg.IdleTimerPeriod < minBackoffDuration {
		minBackoffDuration = cfg.IdleTimerPeriod
	}
	maxBackoffDuration := cfg.MaxRetryBackoffDuration
	if cfg.IdleTimerPeriod < maxBackoffDuration {
		maxBackoffDuration = cfg.IdleTimerPeriod
	}

	return &PollManager{
		cfg:    cfg,
		logger: logger,

		hibernationTimer: utils.NewResettableTimer(),
		pollTicker:       utils.NewPausableTicker(cfg.PollTickerInterval),
		idleTimer:        utils.NewResettableTimer(),
		roundTimer:       utils.NewResettableTimer(),
		retryTicker:      utils.NewBackoffTicker(minBackoffDuration, maxBackoffDuration),
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

// HibernationTimerTicks ticks after a given period
func (pm *PollManager) HibernationTimerTicks() <-chan time.Time {
	return pm.hibernationTimer.Ticks()
}

// RoundTimerTicks ticks after a given period
func (pm *PollManager) RoundTimerTicks() <-chan time.Time {
	return pm.roundTimer.Ticks()
}

func (pm *PollManager) RetryTickerTicks() <-chan time.Time {
	return pm.retryTicker.Ticks()
}

// Start initializes all the timers and determines whether to go into immediate
// hibernation.
func (pm *PollManager) Start(hibernate bool, roundState flux_aggregator_wrapper.OracleRoundState) {
	if hibernate {
		pm.Hibernate()
	} else {
		pm.Awaken(roundState)
	}
}

// ShouldPerformInitialPoll determines whether to perform an initial poll
func (pm *PollManager) ShouldPerformInitialPoll() bool {
	return (!pm.cfg.PollTickerDisabled || !pm.cfg.IdleTimerDisabled) && !pm.cfg.IsHibernating
}

// Reset resets the timers except for the hibernation timer. Will not reset if
// hibernating.
func (pm *PollManager) Reset(roundState flux_aggregator_wrapper.OracleRoundState) {
	if !pm.cfg.IsHibernating {
		pm.startPollTicker()
		pm.startIdleTimer(roundState.StartedAt)
		pm.startRoundTimer(roundStateTimesOutAt(roundState))
	}
}

// Reset resets the idle timer unless hibernating
func (pm *PollManager) ResetIdleTimer(roundStartedAtUTC uint64) {
	if !pm.cfg.IsHibernating {
		pm.startIdleTimer(roundStartedAtUTC)
	}
}

// StartRetryTicker starts the retry ticker
func (pm *PollManager) StartRetryTicker() {
	pm.retryTicker.Start()
}

// StopRetryTicker stops the retry ticker
func (pm *PollManager) StopRetryTicker() {
	pm.retryTicker.Stop()
}

// Stop stops all timers/tickers
func (pm *PollManager) Stop() {
	pm.hibernationTimer.Stop()
	pm.pollTicker.Destroy()
	pm.idleTimer.Stop()
	pm.roundTimer.Stop()
}

// Hibernate sets hibernation to true, starts the hibernation timer and stops
// all other ticker/timers
func (pm *PollManager) Hibernate() {
	pm.logger.Info("entering hibernation mode")

	// Start the hibernation timer
	pm.cfg.IsHibernating = true
	pm.hibernationTimer.Reset(pm.cfg.HibernationPollPeriod)

	// Stop the other tickers
	pm.pollTicker.Pause()
	pm.idleTimer.Stop()
	pm.roundTimer.Stop()
}

// Awaken sets hibernation to false, stops the hibernation timer and starts all
// other tickers
func (pm *PollManager) Awaken(roundState flux_aggregator_wrapper.OracleRoundState) {
	pm.logger.Info("exiting hibernation mode, reactivating contract")

	// Stop the hibernation timer
	pm.cfg.IsHibernating = false
	pm.hibernationTimer.Stop()

	// Start the other tickers
	pm.startPollTicker()
	pm.startIdleTimer(roundState.StartedAt)
	pm.startRoundTimer(roundStateTimesOutAt(roundState))
}

// startPollTicker starts the poll ticker if it is enabled
func (pm *PollManager) startPollTicker() {
	if pm.cfg.PollTickerDisabled {
		pm.pollTicker.Pause()

		return
	}

	pm.pollTicker.Resume()
}

// startIdleTimer starts the idle timer if it is enabled
func (pm *PollManager) startIdleTimer(roundStartedAtUTC uint64) {
	// Stop the retry timer when the idle timer is started
	pm.retryTicker.Stop()

	if pm.cfg.IdleTimerDisabled {
		pm.idleTimer.Stop()

		return
	}

	// Keep using the idleTimer we already have
	if roundStartedAtUTC == 0 {
		pm.logger.Debugw("keeping existing timer, no active round")

		return
	}

	startedAt := time.Unix(int64(roundStartedAtUTC), 0)
	deadline := startedAt.Add(pm.cfg.IdleTimerPeriod)
	deadlineDuration := time.Until(deadline)

	log := pm.logger.With(
		"pollFrequency", pm.cfg.PollTickerInterval,
		"idleDuration", pm.cfg.IdleTimerPeriod,
		"startedAt", roundStartedAtUTC,
		"timeUntilIdleDeadline", deadlineDuration,
	)

	if deadlineDuration <= 0 {
		log.Debugw("not resetting idleTimer, negative duration")

		return
	}

	pm.idleTimer.Reset(deadlineDuration)
	log.Debugw("resetting idleTimer")
}

// startRoundTimer starts the round timer
func (pm *PollManager) startRoundTimer(roundTimesOutAt uint64) {
	log := pm.logger.With(
		"pollFrequency", pm.cfg.PollTickerInterval,
		"idleDuration", pm.cfg.IdleTimerPeriod,
		"timesOutAt", roundTimesOutAt,
	)

	if roundTimesOutAt == 0 {
		log.Debugw("disabling roundTimer, no active round")
		pm.roundTimer.Stop()

		return
	}

	timesOutAt := time.Unix(int64(roundTimesOutAt), 0)
	timeoutDuration := time.Until(timesOutAt)

	if timeoutDuration <= 0 {
		log.Debugw("roundTimer has run down; disabling")
		pm.roundTimer.Stop()

		return
	}

	pm.roundTimer.Reset(timeoutDuration)
	log.Debugw("updating roundState.TimesOutAt", "value", roundTimesOutAt)
}

func roundStateTimesOutAt(rs flux_aggregator_wrapper.OracleRoundState) uint64 {
	return rs.StartedAt + rs.Timeout
}
