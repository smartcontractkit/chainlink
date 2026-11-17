package capabilities

import (
	"context"
	"fmt"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote"
)

// triggerCleanup tracks the pruned trigger publishers of one capability whose registrations are still
// being unregistered from the underlying. done is closed when the last of them finishes.
type triggerCleanup struct {
	pending int
	done    chan struct{}
}

// closeTriggerPublisherAsync closes a pruned trigger publisher and unregisters its triggers from the
// underlying in the background, so a large unregister doesn't block the registry sync under reconcileMu.
//
// Until it finishes, new RegisterTrigger calls for the same capability wait (see gatedTriggerCapability).
// Otherwise a late unregister could remove a registration that another publisher of this capability
// (e.g. the new DON's) has just made for the same trigger ID.
func (w *launcher) closeTriggerPublisherAsync(key shimID, pub remote.TriggerPublisher) {
	w.triggerCleanupMu.Lock()
	if w.triggerCleanupsClosed {
		w.triggerCleanupMu.Unlock()
		// The launcher is shutting down, and so is the underlying; nothing to release.
		if err := pub.Close(); err != nil {
			w.lggr.Errorw("failed to close stale trigger publisher", "capabilityID", key.capID, "donID", key.donID, "method", key.method, "error", err)
		}
		return
	}
	if w.pendingTriggerCleanups == nil {
		w.pendingTriggerCleanups = make(map[string]*triggerCleanup)
	}
	c, ok := w.pendingTriggerCleanups[key.capID]
	if !ok {
		c = &triggerCleanup{done: make(chan struct{})}
		w.pendingTriggerCleanups[key.capID] = c
	}
	c.pending++
	w.triggerCleanupWg.Add(1)
	w.triggerCleanupMu.Unlock()

	go func() {
		defer w.triggerCleanupWg.Done()
		// Not the reconcile ctx: that one may be cancelled as soon as onNewRegistry returns.
		ctx, cancel := w.stopCh.NewCtx()
		defer cancel()
		if err := pub.CloseAndUnregister(ctx); err != nil {
			w.lggr.Errorw("failed to close stale trigger publisher", "capabilityID", key.capID, "donID", key.donID, "method", key.method, "error", err)
		}

		w.triggerCleanupMu.Lock()
		defer w.triggerCleanupMu.Unlock()
		c.pending--
		if c.pending == 0 {
			close(c.done)
			delete(w.pendingTriggerCleanups, key.capID)
		}
	}()
}

// waitForTriggerCleanups blocks until no pruned publisher of capID is still unregistering, or ctx is done.
// It checks the map again after each wake-up: a newer cleanup of the same capability may have started
// after the one it waited on.
func (w *launcher) waitForTriggerCleanups(ctx context.Context, capID string) error {
	for {
		w.triggerCleanupMu.Lock()
		c := w.pendingTriggerCleanups[capID]
		w.triggerCleanupMu.Unlock()
		if c == nil {
			return nil
		}
		select {
		case <-c.done:
		case <-ctx.Done():
			return fmt.Errorf("waiting for stale trigger registrations of %s to be released: %w", capID, ctx.Err())
		}
	}
}

// closeTriggerCleanups stops new async cleanups and waits for the running ones. Pending cleanups are
// cancelled through stopCh, which must already be closed.
func (w *launcher) closeTriggerCleanups() {
	w.triggerCleanupMu.Lock()
	w.triggerCleanupsClosed = true
	w.triggerCleanupMu.Unlock()
	w.triggerCleanupWg.Wait()
}

// gatedTriggerCapability is the underlying handed to trigger publishers. RegisterTrigger first waits
// for pending cleanups of the same capability. UnregisterTrigger and AckEvent are not gated, so a
// pruned publisher's own cleanup never waits on itself.
type gatedTriggerCapability struct {
	capabilities.TriggerCapability
	capID    string
	launcher *launcher
}

func (g gatedTriggerCapability) RegisterTrigger(ctx context.Context, request capabilities.TriggerRegistrationRequest) (<-chan capabilities.TriggerResponse, error) {
	if err := g.launcher.waitForTriggerCleanups(ctx, g.capID); err != nil {
		return nil, err
	}
	return g.TriggerCapability.RegisterTrigger(ctx, request)
}
