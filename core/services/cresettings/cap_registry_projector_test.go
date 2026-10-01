package cresettings

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/capabilities/globalconfig"
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
