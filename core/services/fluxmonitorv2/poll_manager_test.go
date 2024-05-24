package fluxmonitorv2_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/flux_aggregator_wrapper"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/fluxmonitorv2"
)

var (
	pollTickerDefaultDuration = 200 * time.Millisecond
	idleTickerDefaultDuration = 1 * time.Second // Setting this too low will cause the idle timer to fire before the assert
)

func newPollManager(t *testing.T) *fluxmonitorv2.PollManager {
	pm, err := fluxmonitorv2.NewPollManager(fluxmonitorv2.PollManagerConfig{
		IsHibernating:         false,
		PollTickerInterval:    pollTickerDefaultDuration,
		PollTickerDisabled:    false,
		IdleTimerPeriod:       idleTickerDefaultDuration,
		IdleTimerDisabled:     false,
		HibernationPollPeriod: 24 * time.Hour,
	}, logger.TestLogger(t))
	require.NoError(t, err)
	return pm
}

type tickChecker struct {
	pollTicked        bool
	idleTicked        bool
	roundTicked       bool
	hibernationTicked bool
	retryTicked       bool
	initialPoll       bool
}

// watchTicks watches the PollManager for ticks for the waitDuration
func watchTicks(t *testing.T, pm *fluxmonitorv2.PollManager, waitDuration time.Duration) tickChecker {
	ticks := tickChecker{
		pollTicked:        false,
		idleTicked:        false,
		roundTicked:       false,
		hibernationTicked: false,
		retryTicked:       false,
		initialPoll:       false,
	}

	waitCh := time.After(waitDuration)
	for {
		select {
		case <-pm.PollTickerTicks():
			ticks.pollTicked = true
		case <-pm.IdleTimerTicks():
			ticks.idleTicked = true
		case <-pm.RoundTimerTicks():
			ticks.roundTicked = true
		case <-pm.HibernationTimerTicks():
			ticks.hibernationTicked = true
		case <-pm.RetryTickerTicks():
			ticks.retryTicked = true
		case request := <-pm.Poll():
			switch request.Type {
			case fluxmonitorv2.PollRequestTypeInitial:
				ticks.initialPoll = true
			// Don't do anything with the other types for now
			default:
			}

		case <-waitCh:
			waitCh = nil
		}

		if waitCh == nil {
			break
		}
	}

	return ticks
}

func TestPollManager_PollTicker(t *testing.T) {
	t.Parallel()
	pm, err := fluxmonitorv2.NewPollManager(fluxmonitorv2.PollManagerConfig{
		PollTickerInterval:    pollTickerDefaultDuration,
		PollTickerDisabled:    false,
		IdleTimerPeriod:       idleTickerDefaultDuration,
		IdleTimerDisabled:     true,
		HibernationPollPeriod: 24 * time.Hour,
	}, logger.TestLogger(t))
	require.NoError(t, err)

	pm.Start(false, flux_aggregator_wrapper.OracleRoundState{})
	t.Cleanup(pm.Stop)

	ticks := watchTicks(t, pm, 2*time.Second)

	assert.True(t, ticks.pollTicked)
	assert.False(t, ticks.idleTicked)
	assert.False(t, ticks.roundTicked)
}

func TestPollManager_IdleTimer(t *testing.T) {
	t.Parallel()
	pm, err := fluxmonitorv2.NewPollManager(fluxmonitorv2.PollManagerConfig{
		PollTickerInterval:    100 * time.Millisecond,
		PollTickerDisabled:    true,
		IdleTimerPeriod:       idleTickerDefaultDuration,
		IdleTimerDisabled:     false,
		HibernationPollPeriod: 24 * time.Hour,
	}, logger.TestLogger(t))
	require.NoError(t, err)

	pm.Start(false, flux_aggregator_wrapper.OracleRoundState{
		StartedAt: uint64(time.Now().Unix()) - 10, // Even 10 seconds old the idle timer should tick
	})
	t.Cleanup(pm.Stop)

	ticks := watchTicks(t, pm, 2*time.Second)

	assert.False(t, ticks.pollTicked)
	assert.True(t, ticks.idleTicked)
	assert.False(t, ticks.roundTicked)
}

func TestPollManager_RoundTimer(t *testing.T) {
	t.Parallel()
	pm, err := fluxmonitorv2.NewPollManager(fluxmonitorv2.PollManagerConfig{
		PollTickerInterval:    pollTickerDefaultDuration,
		PollTickerDisabled:    true,
		IdleTimerPeriod:       idleTickerDefaultDuration,
		IdleTimerDisabled:     true,
		HibernationPollPeriod: 24 * time.Hour,
	}, logger.TestLogger(t))
	require.NoError(t, err)

	pm.Start(false, flux_aggregator_wrapper.OracleRoundState{
		StartedAt: uint64(time.Now().Unix()),
		Timeout:   1, // in seconds
	})
	t.Cleanup(pm.Stop)

	ticks := watchTicks(t, pm, 2*time.Second)

	assert.False(t, ticks.pollTicked)
	assert.False(t, ticks.idleTicked)
	assert.True(t, ticks.roundTicked)
}

func TestPollManager_RetryTimer(t *testing.T) {
	t.Parallel()
	pm, err := fluxmonitorv2.NewPollManager(fluxmonitorv2.PollManagerConfig{
		PollTickerInterval:      pollTickerDefaultDuration,
		PollTickerDisabled:      true,
		IdleTimerPeriod:         idleTickerDefaultDuration,
		IdleTimerDisabled:       true,
		HibernationPollPeriod:   24 * time.Hour,
		MinRetryBackoffDuration: 200 * time.Microsecond,
		MaxRetryBackoffDuration: 1 * time.Minute,
	}, logger.TestLogger(t))
	require.NoError(t, err)

	pm.Start(false, flux_aggregator_wrapper.OracleRoundState{
		StartedAt: uint64(time.Now().Unix()),
		Timeout:   10000, // in seconds. Don't timeout the round
	})
	t.Cleanup(pm.Stop)

	pm.StartRetryTicker()

	// Retry ticker fires
	ticks := watchTicks(t, pm, 2*time.Second)
	assert.False(t, ticks.pollTicked)
	assert.False(t, ticks.idleTicked)
	assert.False(t, ticks.roundTicked)
	assert.True(t, ticks.retryTicked)

	pm.StopRetryTicker()

	ticks = watchTicks(t, pm, 2*time.Second)
	assert.False(t, ticks.pollTicked)
	assert.False(t, ticks.idleTicked)
	assert.False(t, ticks.roundTicked)
	assert.False(t, ticks.retryTicked)
}

func TestPollManager_InitialPoll(t *testing.T) {
	t.Parallel()
	pm := newPollManager(t)
	pm.Start(false, flux_aggregator_wrapper.OracleRoundState{})

	ticks := watchTicks(t, pm, 1*time.Second)
	assert.True(t, ticks.initialPoll)
}

func TestPollManager_HibernationTimer(t *testing.T) {
	t.Parallel()
	pm, err := fluxmonitorv2.NewPollManager(fluxmonitorv2.PollManagerConfig{
		PollTickerInterval:    pollTickerDefaultDuration,
		PollTickerDisabled:    true,
		IdleTimerPeriod:       idleTickerDefaultDuration,
		IdleTimerDisabled:     true,
		HibernationPollPeriod: 1 * time.Second,
	}, logger.TestLogger(t))
	require.NoError(t, err)

	pm.Start(true, flux_aggregator_wrapper.OracleRoundState{
		StartedAt: uint64(time.Now().Unix()),
		Timeout:   1, // in seconds
	})
	t.Cleanup(pm.Stop)

	ticks := watchTicks(t, pm, 2*time.Second)

	assert.True(t, ticks.hibernationTicked)
}

func TestPollManager_HibernationOnStartThenAwaken(t *testing.T) {
	t.Parallel()
	pm, err := fluxmonitorv2.NewPollManager(fluxmonitorv2.PollManagerConfig{
		PollTickerInterval:    pollTickerDefaultDuration,
		PollTickerDisabled:    false,
		IdleTimerPeriod:       idleTickerDefaultDuration,
		IdleTimerDisabled:     false,
		HibernationPollPeriod: 24 * time.Hour,
	}, logger.TestLogger(t))
	require.NoError(t, err)
	t.Cleanup(pm.Stop)

	pm.Start(true, flux_aggregator_wrapper.OracleRoundState{
		StartedAt: uint64(time.Now().Unix()),
		Timeout:   1, // in seconds
	})

	ticks := watchTicks(t, pm, 2*time.Second)

	assert.False(t, ticks.pollTicked)
	assert.False(t, ticks.idleTicked)
	assert.False(t, ticks.roundTicked)

	pm.Awaken(flux_aggregator_wrapper.OracleRoundState{
		StartedAt: uint64(time.Now().Unix()),
		Timeout:   1,
	})

	ticks = watchTicks(t, pm, 2*time.Second)

	assert.True(t, ticks.pollTicked)
	assert.True(t, ticks.idleTicked)
	assert.True(t, ticks.roundTicked)
}

func TestPollManager_AwakeOnStartThenHibernate(t *testing.T) {
	t.Parallel()
	pm := newPollManager(t)

	pm.Start(false, flux_aggregator_wrapper.OracleRoundState{
		StartedAt: uint64(time.Now().Unix()),
		Timeout:   1,
	})
	t.Cleanup(pm.Stop)

	ticks := watchTicks(t, pm, 2*time.Second)

	assert.True(t, ticks.pollTicked)
	assert.True(t, ticks.idleTicked)
	assert.True(t, ticks.roundTicked)

	pm.Hibernate()

	ticks = watchTicks(t, pm, 2*time.Second)

	assert.False(t, ticks.pollTicked)
	assert.False(t, ticks.idleTicked)
	assert.False(t, ticks.roundTicked)
}

func TestPollManager_ShouldPerformInitialPoll(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name               string
		pollTickerDisabled bool
		idleTimerDisabled  bool
		isHibernating      bool
		want               bool
	}{
		{
			name:               "perform poll - all enabled",
			pollTickerDisabled: false,
			idleTimerDisabled:  false,
			isHibernating:      false,
			want:               true,
		},
		{
			name:               "don't perform poll - hibernating",
			pollTickerDisabled: false,
			idleTimerDisabled:  false,
			isHibernating:      true,
			want:               false,
		},
		{
			name:               "perform poll - only pollTickerDisabled",
			pollTickerDisabled: true,
			idleTimerDisabled:  false,
			isHibernating:      false,
			want:               true,
		},
		{
			name:               "perform poll - only idleTimerDisabled",
			pollTickerDisabled: false,
			idleTimerDisabled:  true,
			isHibernating:      false,
			want:               true,
		},
		{
			name:               "don't perform poll - idleTimerDisabled and pollTimerDisabled",
			pollTickerDisabled: true,
			idleTimerDisabled:  true,
			isHibernating:      false,
			want:               false,
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			pm, err := fluxmonitorv2.NewPollManager(fluxmonitorv2.PollManagerConfig{
				IsHibernating:         tc.isHibernating,
				HibernationPollPeriod: 24 * time.Hour,
				PollTickerInterval:    pollTickerDefaultDuration,
				PollTickerDisabled:    tc.pollTickerDisabled,
				IdleTimerPeriod:       idleTickerDefaultDuration,
				IdleTimerDisabled:     tc.idleTimerDisabled,
			}, logger.TestLogger(t))
			require.NoError(t, err)

			assert.Equal(t, tc.want, pm.ShouldPerformInitialPoll())
		})
	}
}

func TestPollManager_Stop(t *testing.T) {
	t.Parallel()
	pm := newPollManager(t)

	pm.Start(false, flux_aggregator_wrapper.OracleRoundState{
		StartedAt: uint64(time.Now().Unix()),
		Timeout:   1,
	})

	ticks := watchTicks(t, pm, 2*time.Second)

	assert.True(t, ticks.pollTicked)
	assert.True(t, ticks.idleTicked)
	assert.True(t, ticks.roundTicked)

	pm.Stop()

	ticks = watchTicks(t, pm, 2*time.Second)

	assert.False(t, ticks.pollTicked)
	assert.False(t, ticks.idleTicked)
	assert.False(t, ticks.roundTicked)
}

func TestPollManager_ResetIdleTimer(t *testing.T) {
	t.Parallel()
	pm := newPollManager(t)

	// Start again in awake mode
	pm.Start(false, flux_aggregator_wrapper.OracleRoundState{
		StartedAt: uint64(time.Now().Unix()),
		Timeout:   1,
	})
	t.Cleanup(pm.Stop)

	// Idle timer fires when not hibernating
	ticks := watchTicks(t, pm, 2*time.Second)
	assert.True(t, ticks.idleTicked)

	// Idle timer fires again after reset
	pm.ResetIdleTimer(uint64(time.Now().Unix()) + 1) // 1 second after now
	ticks = watchTicks(t, pm, 2*time.Second)
	assert.True(t, ticks.idleTicked)
}

func TestPollManager_ResetIdleTimerWhenHibernating(t *testing.T) {
	t.Parallel()
	pm := newPollManager(t)

	// Start in hibernation
	pm.Start(true, flux_aggregator_wrapper.OracleRoundState{
		StartedAt: uint64(time.Now().Unix()),
		Timeout:   1, // in seconds
	})
	t.Cleanup(pm.Stop)

	// Idle timer does not fire when hibernating
	ticks := watchTicks(t, pm, 2*time.Second)
	assert.False(t, ticks.idleTicked)

	// Idle timer does not reset because in hibernation, so it does not fire
	pm.ResetIdleTimer(uint64(time.Now().Unix()))
	ticks = watchTicks(t, pm, 2*time.Second)
	assert.False(t, ticks.idleTicked)
}

func TestPollManager_Reset(t *testing.T) {
	t.Parallel()
	pm := newPollManager(t)

	// Start again in awake mode
	pm.Start(false, flux_aggregator_wrapper.OracleRoundState{
		StartedAt: uint64(time.Now().Unix()),
		Timeout:   1,
	})
	t.Cleanup(pm.Stop)

	// Ticker/timers fires when not hibernating
	ticks := watchTicks(t, pm, 2*time.Second)
	assert.True(t, ticks.pollTicked)
	assert.True(t, ticks.idleTicked)
	assert.True(t, ticks.roundTicked)

	// Idle timer fires again after reset
	pm.Reset(flux_aggregator_wrapper.OracleRoundState{
		StartedAt: uint64(time.Now().Unix()),
		Timeout:   1,
	})
	ticks = watchTicks(t, pm, 2*time.Second)
	assert.True(t, ticks.pollTicked)
	assert.True(t, ticks.idleTicked)
	assert.True(t, ticks.roundTicked)
}

func TestPollManager_ResetWhenHibernating(t *testing.T) {
	t.Parallel()
	pm := newPollManager(t)

	// Start in hibernation
	pm.Start(true, flux_aggregator_wrapper.OracleRoundState{
		StartedAt: uint64(time.Now().Unix()),
		Timeout:   1, // in seconds
	})
	t.Cleanup(pm.Stop)

	// Ticker/timers do not fire when hibernating
	ticks := watchTicks(t, pm, 2*time.Second)
	assert.False(t, ticks.pollTicked)
	assert.False(t, ticks.idleTicked)
	assert.False(t, ticks.roundTicked)

	// Ticker/timers does not reset because in hibernation, so they do not fire
	pm.Reset(flux_aggregator_wrapper.OracleRoundState{
		StartedAt: uint64(time.Now().Unix()),
		Timeout:   1, // in seconds
	})
	ticks = watchTicks(t, pm, 2*time.Second)
	assert.False(t, ticks.pollTicked)
	assert.False(t, ticks.idleTicked)
	assert.False(t, ticks.roundTicked)
}
