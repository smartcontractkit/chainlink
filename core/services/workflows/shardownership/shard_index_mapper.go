package shardownership

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strings"
	"sync"

	commoncap "github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/registry"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
)

// shardNameMarker is the marker a workflow DON's name is suffixed with to
// encode its shard index (e.g. "workflow-1-zone-a-shard-1" -> shard 1). A
// name with no such suffix (e.g. "workflow-1-zone-a") is shard index 0.
const shardNameMarker = "shard-"

// ShardIndexMapper tracks the workflow DONs currently known to the Capabilities
// Registry and detects their shard indices (0, 1, 2, ...).
// It's a global singleton object used by ManualShardResolver and, in consequence,
// every per-workflow instance of ShardFailoverManager.
// It implements registrysyncer/v2.Listener and refreshes on every registry
// update.
type ShardIndexMapper struct {
	lggr logger.Logger

	mu      sync.RWMutex
	byIndex []commoncap.DON

	readyOnce sync.Once
	ready     chan struct{}
}

func NewShardIndexMapper(lggr logger.Logger) *ShardIndexMapper {
	return &ShardIndexMapper{lggr: logger.Named(lggr, "ShardIndexMapper"), ready: make(chan struct{})}
}

// OnNewRegistry implements registrysyncer/v2.Listener. It determines this
// node's shard-group name prefix (its own DON name with any "shard-X" suffix
// stripped), then indexes every other workflow DON that shares at least one
// family with its own workflow DON and shares that name prefix, by the shard
// index encoded in its name.
func (s *ShardIndexMapper) OnNewRegistry(ctx context.Context, reg *registry.RegistryMetadata) error {
	localNode, err := reg.LocalNode(ctx)
	if err != nil {
		return fmt.Errorf("failed to resolve local node: %w", err)
	}
	if localNode.WorkflowDON.ID == 0 {
		return errors.New("local node does not belong to a workflow DON")
	}

	namePrefix := shardGroupNamePrefix(localNode.WorkflowDON.Name)

	byIndex, err := shardDONsByIndex(reg, localNode.WorkflowDON, namePrefix)
	if err != nil {
		return err
	}

	s.mu.Lock()
	s.byIndex = byIndex
	s.mu.Unlock()
	s.readyOnce.Do(func() { close(s.ready) })

	s.lggr.Debugw("refreshed shard DON index", "namePrefix", namePrefix, "workflowDONs", len(byIndex))
	return nil
}

// sharesFamily reports whether a and b have at least one family in common.
func sharesFamily(a, b []string) bool {
	for _, family := range a {
		if slices.Contains(b, family) {
			return true
		}
	}
	return false
}

// shardDONsByIndex finds every workflow DON in the registry that shares at
// least one family with localDON (the local node's own shard group) and
// shares namePrefix (the local node's own DON name with any shard suffix
// stripped), and returns them indexed by the shard index encoded in their
// name (see shardIndexFromName).
func shardDONsByIndex(reg *registry.RegistryMetadata, localDON commoncap.DON, namePrefix string) ([]commoncap.DON, error) {
	byIndex := make(map[uint32]commoncap.DON)
	maxIndex := uint32(0)
	for _, don := range reg.IDsToDONs {
		if !don.AcceptsWorkflows || !sharesFamily(don.Families, localDON.Families) || shardGroupNamePrefix(don.Name) != namePrefix {
			continue
		}
		idx, err := shardIndexFromName(don.Name)
		if err != nil {
			return nil, fmt.Errorf("DON %q: %w", don.Name, err)
		}
		if existing, ok := byIndex[idx]; ok {
			return nil, fmt.Errorf("name prefix %q has two DONs at shard index %d: %q and %q", namePrefix, idx, existing.Name, don.Name)
		}
		byIndex[idx] = don.DON
		if idx > maxIndex {
			maxIndex = idx
		}
	}
	if len(byIndex) == 0 {
		return nil, fmt.Errorf("no workflow DONs found sharing a family with local DON %q and name prefix %q", localDON.Name, namePrefix)
	}

	result := make([]commoncap.DON, maxIndex+1)
	for idx, don := range byIndex {
		result[idx] = don
	}
	return result, nil
}

// shardIndexFromName parses the single-digit shard index encoded in a
// workflow DON's name (e.g. "workflow-1-zone-a-shard-1" -> 1). A name with no
// "shard-" suffix is shard index 0.
func shardIndexFromName(name string) (uint32, error) {
	idx := strings.LastIndex(name, shardNameMarker)
	if idx == -1 {
		return 0, nil
	}
	suffix := name[idx+len(shardNameMarker):]
	if len(suffix) != 1 || suffix[0] < '0' || suffix[0] > '9' {
		return 0, fmt.Errorf("expected DON name %q to end with %q followed by a single digit, got suffix %q", name, shardNameMarker, suffix)
	}
	return uint32(suffix[0] - '0'), nil
}

// shardGroupNamePrefix strips the "_shard-X" or "-shard-X" suffix (and its
// separating underscore or hyphen) from a workflow DON's name, leaving the
// base name shared by every DON in its shard group (e.g.
// "workflow-1-zone-a-shard-1" -> "workflow-1-zone-a"). A name with no such
// suffix is returned unchanged.
func shardGroupNamePrefix(name string) string {
	idx := strings.LastIndex(name, shardNameMarker)
	if idx == -1 {
		return name
	}
	prefix := name[:idx]
	if strings.HasSuffix(prefix, "_") || strings.HasSuffix(prefix, "-") {
		prefix = prefix[:len(prefix)-1]
	}
	return prefix
}

// WaitReady blocks until the first registry snapshot has been processed (so
// DonByShardIndex reflects at least one real registry view rather than the
// empty startup state), or until ctx is done, whichever comes first.
func (s *ShardIndexMapper) WaitReady(ctx context.Context) error {
	select {
	case <-s.ready:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// DonByShardIndex returns the DON assigned to shardIndex (0-based, ascending
// DON ID order among the current workflow DONs), or nil if shardIndex is out
// of range for the last registry snapshot seen, or unassigned (a gap in the
// shard indices seen).
func (s *ShardIndexMapper) DonByShardIndex(_ context.Context, shardIndex uint32) *commoncap.DON {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if int(shardIndex) >= len(s.byIndex) {
		return nil
	}
	don := s.byIndex[shardIndex]
	if don.ID == 0 {
		return nil
	}
	return &don
}
