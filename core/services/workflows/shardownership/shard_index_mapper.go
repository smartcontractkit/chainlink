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

// workflowFamilySuffix identifies a special DON family shared by every
// workflow DON in a zone (e.g. "zone-a_workflow").
const workflowFamilySuffix = "_workflow"

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
// node's shard-group family (the single DON family suffixed "_workflow" that
// its own workflow DON belongs to), then indexes every workflow DON sharing
// that family by the shard index encoded in its name.
func (s *ShardIndexMapper) OnNewRegistry(ctx context.Context, reg *registry.RegistryMetadata) error {
	localNode, err := reg.LocalNode(ctx)
	if err != nil {
		return fmt.Errorf("failed to resolve local node: %w", err)
	}
	if localNode.WorkflowDON.ID == 0 {
		return errors.New("local node does not belong to a workflow DON")
	}

	family, err := workflowFamilyOf(localNode.WorkflowDON)
	if err != nil {
		return err
	}

	byIndex, err := shardDONsByIndex(reg, family)
	if err != nil {
		return err
	}

	s.mu.Lock()
	s.byIndex = byIndex
	s.mu.Unlock()
	s.readyOnce.Do(func() { close(s.ready) })

	s.lggr.Debugw("refreshed shard DON index", "family", family, "workflowDONs", len(byIndex))
	return nil
}

// workflowFamilyOf returns the single DON family suffixed "_workflow" that
// don belongs to. It errors if don belongs to zero or more than one such
// family, since that leaves the shard group ambiguous.
func workflowFamilyOf(don commoncap.DON) (string, error) {
	var found string
	for _, family := range don.Families {
		if !strings.HasSuffix(family, workflowFamilySuffix) {
			continue
		}
		if found != "" {
			return "", fmt.Errorf("DON %q belongs to more than one %q family: %q and %q", don.Name, workflowFamilySuffix, found, family)
		}
		found = family
	}
	if found == "" {
		return "", fmt.Errorf("DON %q does not belong to any family suffixed %q", don.Name, workflowFamilySuffix)
	}
	return found, nil
}

// shardDONsByIndex finds every workflow DON in the registry that belongs to
// family (the local node's own shard group), and returns them indexed by the
// shard index encoded in their name (see shardIndexFromName).
func shardDONsByIndex(reg *registry.RegistryMetadata, family string) ([]commoncap.DON, error) {
	byIndex := make(map[uint32]commoncap.DON)
	maxIndex := uint32(0)
	for _, don := range reg.IDsToDONs {
		if !don.AcceptsWorkflows || !slices.Contains(don.Families, family) {
			continue
		}
		idx, err := shardIndexFromName(don.Name)
		if err != nil {
			return nil, fmt.Errorf("DON %q in family %q: %w", don.Name, family, err)
		}
		if existing, ok := byIndex[idx]; ok {
			return nil, fmt.Errorf("family %q has two DONs at shard index %d: %q and %q", family, idx, existing.Name, don.Name)
		}
		byIndex[idx] = don.DON
		if idx > maxIndex {
			maxIndex = idx
		}
	}
	if len(byIndex) == 0 {
		return nil, fmt.Errorf("no workflow DONs found in family %q", family)
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
