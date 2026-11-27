package standardcapabilities

import (
	"context"
	"maps"
	"slices"
	"sync"
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/types/core"
)

const defaultRegistryRemoveTimeout = 30 * time.Second

// registeredCapabilities is the capabilities registry as one plugin sees it. It records
// what the plugin adds, so that the host can remove it from the registry before it stops
// the plugin.
//
// A registry entry that outlives its plugin makes the calls into it wait for a reconnect
// that cannot succeed, and those calls block the registration of a replacement, for
// example when the job is proposed again on a running node.
type registeredCapabilities struct {
	core.CapabilitiesRegistry

	mu  sync.Mutex
	ids map[string]struct{}
}

func trackRegisteredCapabilities(registry core.CapabilitiesRegistry) *registeredCapabilities {
	return &registeredCapabilities{CapabilitiesRegistry: registry, ids: map[string]struct{}{}}
}

func (r *registeredCapabilities) Add(ctx context.Context, c capabilities.BaseCapability) error {
	info, err := c.Info(ctx)
	if err != nil {
		return err
	}
	if err := r.CapabilitiesRegistry.Add(ctx, c); err != nil {
		return err
	}
	r.mu.Lock()
	r.ids[info.ID] = struct{}{}
	r.mu.Unlock()
	return nil
}

func (r *registeredCapabilities) Remove(ctx context.Context, id string) error {
	if err := r.CapabilitiesRegistry.Remove(ctx, id); err != nil {
		return err
	}
	r.mu.Lock()
	delete(r.ids, id)
	r.mu.Unlock()
	return nil
}

// removeAll removes from the registry what the plugin added and did not remove itself.
// It must be called while the plugin is still running: the calls in flight can then
// finish, so the removal is not held up by them.
func (r *registeredCapabilities) removeAll(lggr logger.Logger, timeout time.Duration) {
	if r == nil {
		return
	}
	if timeout == 0 {
		timeout = defaultRegistryRemoveTimeout
	}
	r.mu.Lock()
	ids := slices.Sorted(maps.Keys(r.ids))
	r.mu.Unlock()

	for _, id := range ids {
		if !r.answers(id, timeout) {
			lggr.Errorw("Capability not removed from registry: the registry did not answer", "capabilityID", id, "timeout", timeout)
			return
		}
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		err := r.Remove(ctx, id)
		cancel()
		if err != nil {
			lggr.Warnw("Failed to remove capability from registry", "capabilityID", id, "err", err)
		}
	}
}

// answers reports whether a read of the registry returns within the timeout. A registry
// that does not observe a context while it waits for its locks cannot abandon a Remove,
// and one that completes after the plugin's replacement was added would remove the
// replacement. A read that completes late changes nothing, so the read is the call that
// is allowed to time out.
func (r *registeredCapabilities) answers(id string, timeout time.Duration) bool {
	answered := make(chan struct{})
	go func() {
		defer close(answered)
		// Only the return matters, not the result.
		_, _ = r.Get(context.Background(), id)
	}()
	select {
	case <-answered:
		return true
	case <-time.After(timeout):
		return false
	}
}
