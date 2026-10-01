package cresettings

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/capabilities/globalconfig"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/globalconfig/globalconfigtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
)

func TestCapRegistryProjector_PicksUpCommitAfterTrigger(t *testing.T) {
	t.Parallel()

	committed := &committedStore{}
	gc := globalconfig.New()
	p := newCapRegistryProjector(logger.TestLogger(t), committed.load, gc)
	p.idleInterval = time.Hour // only the trigger-driven fast polling can apply the change
	p.fastInterval = 10 * time.Millisecond
	require.NoError(t, p.Start(t.Context()))
	t.Cleanup(func() { require.NoError(t, p.Close()) })

	// Hint arrives while the transaction is still open (nothing committed yet).
	p.Trigger()
	require.Never(t, func() bool { _, v := gc.Load(); return v != 0 }, 100*time.Millisecond, 10*time.Millisecond)

	// The transaction commits shortly afterwards.
	committed.set(&job.CRESettingsSpec{ConfigType: ConfigTypeCapRegistry, OffchainConfig: `{"version":2}`, Hash: "h2"})
	require.Eventually(t, func() bool { _, v := gc.Load(); return v == 2 }, 5*time.Second, 10*time.Millisecond)

	// A committed delete is picked up the same way and notifies subscribers.
	updates, unsubscribe := gc.Subscribe()
	defer unsubscribe()
	committed.set(nil)
	p.Trigger()
	select {
	case <-updates:
	case <-time.After(5 * time.Second):
		t.Fatal("no notification after committed delete")
	}
	reg, v := gc.LoadParsed()
	assert.Nil(t, reg)
	assert.Equal(t, uint64(0), v)
}

func TestCapRegistryProjector_IdlePollingWithoutTrigger(t *testing.T) {
	t.Parallel()

	committed := &committedStore{}
	gc := globalconfig.New()
	p := newCapRegistryProjector(logger.TestLogger(t), committed.load, gc)
	p.idleInterval = 20 * time.Millisecond
	require.NoError(t, p.Start(t.Context()))
	t.Cleanup(func() { require.NoError(t, p.Close()) })

	committed.set(&job.CRESettingsSpec{ConfigType: ConfigTypeCapRegistry, OffchainConfig: `{"version":4}`, Hash: "h4"})
	require.Eventually(t, func() bool { _, v := gc.Load(); return v == 4 }, 5*time.Second, 10*time.Millisecond)
}

func TestCapRegistryProjector_Metrics(t *testing.T) {
	t.Parallel()

	committed := &committedStore{}
	p := newCapRegistryProjector(logger.TestLogger(t), committed.load, globalconfig.New())
	m, reader := globalconfigtest.NewMetrics(t)
	p.metrics = m
	prod := globalconfigtest.Series("cre", "prod")

	committed.set(&job.CRESettingsSpec{ConfigType: ConfigTypeCapRegistry, OffchainConfig: `{"domain":"cre","env":"prod","version":7}`, Hash: "h7"})
	require.NoError(t, p.Refresh(t.Context()))
	assert.Equal(t, int64(7), globalconfigtest.Collect(t, reader)[globalconfig.MetricAppliedVersion][prod])

	// Committed state older than what is applied cannot be applied (the DB prevents this; the
	// runtime check is a backstop) and counts as an apply error.
	committed.set(&job.CRESettingsSpec{ConfigType: ConfigTypeCapRegistry, OffchainConfig: `{"domain":"cre","env":"prod","version":6}`, Hash: "h6"})
	require.Error(t, p.Refresh(t.Context()))

	// Failing to read committed state is an apply error for the last applied series.
	committed.mu.Lock()
	committed.err = assert.AnError
	committed.mu.Unlock()
	require.ErrorIs(t, p.Refresh(t.Context()), assert.AnError)
	committed.mu.Lock()
	committed.err = nil
	committed.mu.Unlock()

	// Withdrawn: the same series drops to 0.
	committed.set(nil)
	require.NoError(t, p.Refresh(t.Context()))

	got := globalconfigtest.Collect(t, reader)
	assert.Equal(t, int64(0), got[globalconfig.MetricAppliedVersion][prod])
	assert.Equal(t, int64(2), got[globalconfig.MetricApplyErrors][prod])
}
