package fluxmonitorv2

import (
	"fmt"
	"sync/atomic"
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/generated/flux_aggregator_wrapper"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

type PollManagerConfig struct {
	IsHibernating           bool
	PollTickerInterval      time.Duration
	PollTickerDisabled      bool
	IdleTimerPeriod         time.Duration
	IdleTimerDisabled       bool
	DrumbeatSchedule        string
	DrumbeatEnabled         bool
	DrumbeatRandomDelay     time.Duration
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

	isHibernating    atomic.Bool
	hibernationTimer utils.ResettableTimer
	pollTicker       utils.PausableTicker
	idleTimer        utils.ResettableTimer
	roundTimer       utils.ResettableTimer
	retryTicker      utils.BackoffTicker
	drumbeat         utils.CronTicker
	chPoll           chan PollRequest

	logger logger.Logger
}

// NewPollManager initializes a new PollManager
func NewPollManager(cfg PollManagerConfig, lggr logger.Logger) (*PollManager, error) {
	minBackoffDuration := cfg.MinRetryBackoffDuration
	if cfg.IdleTimerPeriod < minBackoffDuration {
		minBackoffDuration = cfg.IdleTimerPeriod
	}
	maxBackoffDuration := cfg.MaxRetryBackoffDuration
	if cfg.IdleTimerPeriod < maxBackoffDuration {
		maxBackoffDuration = cfg.IdleTimerPeriod
	}
	// Always initialize the idle timer so that no matter what it has a ticker
	// and won't get starved by an old startedAt timestamp from the oracle state on boot.
	var idleTimer = utils.NewResettableTimer()
	if !cfg.IdleTimerDisabled {
		idleTimer.Reset(cfg.IdleTimerPeriod)
	}

	p := &PollManager{
		cfg:    cfg,
		logger: logger.Named(lggr, "PollManager"),

		hibernationTimer: utils.NewResettableTimer(),
		pollTicker:       utils.NewPausableTicker(cfg.PollTickerInterval),
		idleTimer:        idleTimer,
		roundTimer:       utils.NewResettableTimer(),
		retryTicker:      utils.NewBackoffTicker(minBackoffDuration, maxBackoffDuration),
		chPoll:           make(chan PollRequest),
	}
	var err error
	if cfg.DrumbeatEnabled {
		p.drumbeat, err = utils.NewCronTicker(cfg.DrumbeatSchedule)
		if err != nil {
			return nil, err
		}
	}
	p.isHibernating.Store(cfg.IsHibernating)
	return p, nil
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

// RetryTickerTicks ticks with a backoff when the retry ticker is activated
func (pm *PollManager) RetryTickerTicks() <-chan time.Time {
	return pm.retryTicker.Ticks()
}

// DrumbeatTicks ticks on a cron schedule when the drumbeat ticker is activated
func (pm *PollManager) DrumbeatTicks() <-chan time.Time {
	return pm.drumbeat.Ticks()
}

// Poll returns a channel which the manager will use to send polling requests
//
// Note: In the future, we should change the tickers above to send their request
// through this channel to simplify the listener.
func (pm *PollManager) Poll() <-chan PollRequest {
	return pm.chPoll
}

// Start initializes all the timers and determines whether to go into immediate
// hibernation.
func (pm *PollManager) Start(hibernate bool, roundState flux_aggregator_wrapper.OracleRoundState) {
	pm.isHibernating.Store(hibernate)

	if pm.ShouldPerformInitialPoll() {
		// We want this to be non blocking but if there is no received for the
		// polling channel, this go routine would hang around forever. Since we
		// should always have a receiver for the polling channel, set a timeout
		// of 5 seconds to kill the goroutine.
		go func() {
			select {
			case pm.chPoll <- PollRequest{PollRequestTypeInitial, time.Now()}:
			case <-time.After(5 * time.Second):
				pm.logger.Warn("Start up poll was not consumed")
			}
		}()
	}

	pm.maybeWarnAboutIdleAndPollIntervals()

	if hibernate {
		pm.Hibernate()
	} else {
		pm.Awaken(roundState)
	}
}

// ShouldPerformInitialPoll determines whether to perform an initial poll
func (pm *PollManager) ShouldPerformInitialPoll() bool {
	return (!pm.cfg.PollTickerDisabled || !pm.cfg.IdleTimerDisabled) && !pm.isHibernating.Load()
}

// Reset resets the timers except for the hibernation timer. Will not reset if
// hibernating.
func (pm *PollManager) Reset(roundState flux_aggregator_wrapper.OracleRoundState) {
	if pm.isHibernating.Load() {
		pm.hibernationTimer.Reset(pm.cfg.HibernationPollPeriod)
	} else {
		pm.startPollTicker()
		pm.startIdleTimer(roundState.StartedAt)
		pm.startRoundTimer(roundStateTimesOutAt(roundState))
		pm.startDrumbeat()
	}
}

// ResetIdleTimer resets the idle timer unless hibernating
func (pm *PollManager) ResetIdleTimer(roundStartedAtUTC uint64) {
	if !pm.isHibernating.Load() {
		pm.startIdleTimer(roundStartedAtUTC)
	}
}

// StartRetryTicker starts the retry ticker
func (pm *PollManager) StartRetryTicker() bool {
	return pm.retryTicker.Start()
}

// StopRetryTicker stops the retry ticker
func (pm *PollManager) StopRetryTicker() {
	if pm.retryTicker.Stop() {
		pm.logger.Debug("stopped retry ticker")
	}
}

// Stop stops all timers/tickers
func (pm *PollManager) Stop() {
	pm.hibernationTimer.Stop()
	pm.pollTicker.Destroy()
	pm.idleTimer.Stop()
	pm.roundTimer.Stop()
	pm.drumbeat.Stop()
}

// Hibernate sets hibernation to true, starts the hibernation timer and stops
// all other ticker/timers
func (pm *PollManager) Hibernate() {
	pm.logger.Infof("entering hibernation mode (period: %v)", pm.cfg.HibernationPollPeriod)

	// Start the hibernation timer
	pm.isHibernating.Store(true)
	pm.hibernationTimer.Reset(pm.cfg.HibernationPollPeriod)

	// Stop the other tickers
	pm.pollTicker.Pause()
	pm.idleTimer.Stop()
	pm.roundTimer.Stop()
	pm.drumbeat.Stop()
	pm.StopRetryTicker()
}

// Awaken sets hibernation to false, stops the hibernation timer and starts all
// other tickers
func (pm *PollManager) Awaken(roundState flux_aggregator_wrapper.OracleRoundState) {
	pm.logger.Info("exiting hibernation mode, reactivating contract")

	// Stop the hibernation timer
	pm.isHibernating.Store(false)
	pm.hibernationTimer.Stop()

	// Start the other tickers
	pm.startPollTicker()
	pm.startIdleTimer(roundState.StartedAt)
	pm.startRoundTimer(roundStateTimesOutAt(roundState))
	pm.startDrumbeat()
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
	if pm.cfg.IdleTimerDisabled {
		pm.idleTimer.Stop()

		return
	}

	// Keep using the idleTimer we already have
	if roundStartedAtUTC == 0 {
		pm.logger.Debugw("not resetting idleTimer, no active round")

		return
	}

	startedAt := time.Unix(int64(roundStartedAtUTC), 0)
	deadline := startedAt.Add(pm.cfg.IdleTimerPeriod)
	deadlineDuration := time.Until(deadline)

	log := logger.With(pm.logger,
		"pollFrequency", pm.cfg.PollTickerInterval,
		"idleDuration", pm.cfg.IdleTimerPeriod,
		"startedAt", roundStartedAtUTC,
		"timeUntilIdleDeadline", deadlineDuration,
	)

	if deadlineDuration <= 0 {
		log.Debugw("not resetting idleTimer, round was started further in the past than idle timer period")
		return
	}

	// Stop the retry timer when the idle timer is started
	if pm.retryTicker.Stop() {
		pm.logger.Debugw("stopped the retryTicker")
	}

	pm.idleTimer.Reset(deadlineDuration)
	log.Debugw("resetting idleTimer")
}

// startRoundTimer starts the round timer
func (pm *PollManager) startRoundTimer(roundTimesOutAt uint64) {
	log := logger.With(pm.logger,
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
		log.Debugw(fmt.Sprintf("disabling roundTimer, as the round is already past its timeout by %v", -timeoutDuration))
		pm.roundTimer.Stop()

		return
	}

	pm.roundTimer.Reset(timeoutDuration)
	log.Debugw("updating roundState.TimesOutAt", "value", roundTimesOutAt)
}

// startDrumbeat starts the drumbeat ticker if it is enabled
func (pm *PollManager) startDrumbeat() {
	if !pm.cfg.DrumbeatEnabled {
		if pm.drumbeat.Stop() {
			pm.logger.Debug("disabled drumbeat ticker")
		}
		return
	}

	if pm.drumbeat.Start() {
		pm.logger.Debugw("started drumbeat ticker", "schedule", pm.cfg.DrumbeatSchedule)
	}
}

func roundStateTimesOutAt(rs flux_aggregator_wrapper.OracleRoundState) uint64 {
	return rs.StartedAt + rs.Timeout
}

// ShouldPerformInitialPoll determines whether to perform an initial poll
func (pm *PollManager) maybeWarnAboutIdleAndPollIntervals() {
	if !pm.cfg.IdleTimerDisabled && !pm.cfg.PollTickerDisabled && pm.cfg.IdleTimerPeriod < pm.cfg.PollTickerInterval {
		pm.logger.Warnw("The value of IdleTimerPeriod is lower than PollTickerInterval. The idle timer should usually be less frequent that poll",
			"IdleTimerPeriod", pm.cfg.IdleTimerPeriod, "PollTickerInterval", pm.cfg.PollTickerInterval)
	}
}
