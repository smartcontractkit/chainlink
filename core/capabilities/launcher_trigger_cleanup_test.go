package capabilities

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote"
)

// blockingCleanupPublisher is a pruned trigger publisher whose CloseAndUnregister blocks until
// release is closed or its ctx is cancelled.
type blockingCleanupPublisher struct {
	remote.TriggerPublisher // only the methods below are called
	release                 chan struct{}
	closeCalled             atomic.Bool // plain Close, without unregistering
	unregisterCalled        atomic.Bool
	unregistered            atomic.Bool
	unregisterCtxErr        atomic.Value
}

func newBlockingCleanupPublisher() *blockingCleanupPublisher {
	return &blockingCleanupPublisher{release: make(chan struct{})}
}

func (p *blockingCleanupPublisher) Close() error {
	p.closeCalled.Store(true)
	return nil
}

func (p *blockingCleanupPublisher) CloseAndUnregister(ctx context.Context) error {
	p.unregisterCalled.Store(true)
	select {
	case <-p.release:
		p.unregistered.Store(true)
		return nil
	case <-ctx.Done():
		p.unregisterCtxErr.Store(ctx.Err())
		return ctx.Err()
	}
}

type countingTrigger struct {
	capabilities.CapabilityInfo
	registrations atomic.Int32
}

func (c *countingTrigger) RegisterTrigger(context.Context, capabilities.TriggerRegistrationRequest) (<-chan capabilities.TriggerResponse, error) {
	c.registrations.Add(1)
	return make(chan capabilities.TriggerResponse), nil
}

func (c *countingTrigger) UnregisterTrigger(context.Context, capabilities.TriggerRegistrationRequest) error {
	return nil
}

func (c *countingTrigger) AckEvent(context.Context, string, string, string) error { return nil }

func newCleanupTestLauncher(t *testing.T) *launcher {
	return &launcher{lggr: logger.Sugared(logger.Test(t)), stopCh: make(services.StopChan)}
}

func registerAsync(t *testing.T, g gatedTriggerCapability) <-chan error {
	errCh := make(chan error, 1)
	go func() {
		_, err := g.RegisterTrigger(t.Context(), capabilities.TriggerRegistrationRequest{TriggerID: "trigger_0"})
		errCh <- err
	}()
	return errCh
}

func TestLauncher_TriggerCleanup_RegistrationWaitsForCleanupOfSameCapability(t *testing.T) {
	t.Parallel()
	w := newCleanupTestLauncher(t)
	const capID, otherCapID = "trigger@1.0.0", "other-trigger@1.0.0"

	pub := newBlockingCleanupPublisher()
	w.closeTriggerPublisherAsync(shimID{capID: capID, donID: 1, method: "Trigger"}, pub)

	underlying := &countingTrigger{}
	errCh := registerAsync(t, gatedTriggerCapability{TriggerCapability: underlying, capID: capID, launcher: w})

	// Other capabilities aren't gated.
	other := &countingTrigger{}
	_, err := gatedTriggerCapability{TriggerCapability: other, capID: otherCapID, launcher: w}.
		RegisterTrigger(t.Context(), capabilities.TriggerRegistrationRequest{TriggerID: "trigger_0"})
	require.NoError(t, err)
	assert.Equal(t, int32(1), other.registrations.Load())

	require.Never(t, func() bool { return underlying.registrations.Load() > 0 }, 100*time.Millisecond, 10*time.Millisecond)
	assert.True(t, pub.unregisterCalled.Load())

	close(pub.release)
	require.NoError(t, <-errCh)
	assert.True(t, pub.unregistered.Load())
	assert.Equal(t, int32(1), underlying.registrations.Load())

	w.triggerCleanupMu.Lock()
	assert.Empty(t, w.pendingTriggerCleanups)
	w.triggerCleanupMu.Unlock()
}

func TestLauncher_TriggerCleanup_RegistrationWaitsForAllPendingCleanups(t *testing.T) {
	t.Parallel()
	w := newCleanupTestLauncher(t)
	const capID = "trigger@1.0.0"

	pub1, pub2 := newBlockingCleanupPublisher(), newBlockingCleanupPublisher()
	w.closeTriggerPublisherAsync(shimID{capID: capID, donID: 1, method: "Trigger"}, pub1)
	w.closeTriggerPublisherAsync(shimID{capID: capID, donID: 2, method: "Trigger"}, pub2)

	underlying := &countingTrigger{}
	errCh := registerAsync(t, gatedTriggerCapability{TriggerCapability: underlying, capID: capID, launcher: w})

	close(pub1.release)
	require.Never(t, func() bool { return underlying.registrations.Load() > 0 }, 100*time.Millisecond, 10*time.Millisecond,
		"must keep waiting while another cleanup of the same capability is pending")

	close(pub2.release)
	require.NoError(t, <-errCh)
	assert.Equal(t, int32(1), underlying.registrations.Load())
}

func TestLauncher_TriggerCleanup_WaitingRegistrationReturnsOnCtxCancel(t *testing.T) {
	t.Parallel()
	w := newCleanupTestLauncher(t)
	const capID = "trigger@1.0.0"

	pub := newBlockingCleanupPublisher()
	w.closeTriggerPublisherAsync(shimID{capID: capID, donID: 1, method: "Trigger"}, pub)
	t.Cleanup(func() { close(pub.release); w.triggerCleanupWg.Wait() })

	underlying := &countingTrigger{}
	ctx, cancel := context.WithTimeout(t.Context(), 50*time.Millisecond)
	defer cancel()
	_, err := gatedTriggerCapability{TriggerCapability: underlying, capID: capID, launcher: w}.
		RegisterTrigger(ctx, capabilities.TriggerRegistrationRequest{TriggerID: "trigger_0"})
	require.ErrorIs(t, err, context.DeadlineExceeded)
	assert.Equal(t, int32(0), underlying.registrations.Load())
}

func TestLauncher_TriggerCleanup_ShutdownCancelsPendingAndClosesLaterPublishersOnly(t *testing.T) {
	t.Parallel()
	w := newCleanupTestLauncher(t)
	const capID = "trigger@1.0.0"

	pending := newBlockingCleanupPublisher()
	w.closeTriggerPublisherAsync(shimID{capID: capID, donID: 1, method: "Trigger"}, pending)

	close(w.stopCh)
	done := make(chan struct{})
	go func() {
		w.closeTriggerCleanups()
		close(done)
	}()
	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Fatal("closeTriggerCleanups didn't return after stopCh was closed")
	}
	require.ErrorIs(t, pending.unregisterCtxErr.Load().(error), context.Canceled)

	// After shutdown, pruned publishers are only closed, without unregistering.
	late := newBlockingCleanupPublisher()
	w.closeTriggerPublisherAsync(shimID{capID: capID, donID: 2, method: "Trigger"}, late)
	assert.True(t, late.closeCalled.Load())
	assert.False(t, late.unregisterCalled.Load())
}

func TestLauncher_TriggerCleanup_RegistrationWaitsForCleanupStartedAfterTheOneItWaitedOn(t *testing.T) {
	t.Parallel()
	w := newCleanupTestLauncher(t)
	const capID = "trigger@1.0.0"

	first := &triggerCleanup{pending: 1, done: make(chan struct{})}
	w.pendingTriggerCleanups = map[string]*triggerCleanup{capID: first}

	underlying := &countingTrigger{}
	errCh := registerAsync(t, gatedTriggerCapability{TriggerCapability: underlying, capID: capID, launcher: w})
	require.Never(t, func() bool { return underlying.registrations.Load() > 0 }, 50*time.Millisecond, 10*time.Millisecond)

	// The first cleanup finishes and a second one of the same capability starts before the waiter
	// wakes up. Both happen under the lock, so the waiter can only see the map after the swap.
	second := &triggerCleanup{pending: 1, done: make(chan struct{})}
	w.triggerCleanupMu.Lock()
	close(first.done)
	w.pendingTriggerCleanups[capID] = second
	w.triggerCleanupMu.Unlock()

	require.Never(t, func() bool { return underlying.registrations.Load() > 0 }, 100*time.Millisecond, 10*time.Millisecond,
		"must keep waiting for the newer cleanup")

	w.triggerCleanupMu.Lock()
	close(second.done)
	delete(w.pendingTriggerCleanups, capID)
	w.triggerCleanupMu.Unlock()

	require.NoError(t, <-errCh)
	assert.Equal(t, int32(1), underlying.registrations.Load())
}
