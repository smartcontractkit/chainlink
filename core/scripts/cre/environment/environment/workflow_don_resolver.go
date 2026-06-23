package environment

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
)

// workflowDONSelector picks which workflow DON a deploy targets.
//
// Deploy uses --workflow-don-name (nodesets.name) or --don-family with optional --shard-index.
// --container-name-pattern is copy-only; see workflow.go.
type workflowDONSelector struct {
	ExplicitName string // --workflow-don-name, e.g. "feeds-zone-a" or "shard0"
	DonFamily    string // --don-family; resolves when exactly one workflow DON matches (or one per shard index)
	ShardIndex   *uint  // --shard-index when set; required for multiple shard DONs in the same family
}

// ResolveWorkflowDONMetadata returns the workflow DON metadata from the topology saved by env start.
//
// Used by workflow deploy (registry don ID, vault DB polling) and deploy target resolution.
// Selection priority:
//  1. --workflow-don-name when set — exact match on nodesets.name.
//  2. --don-family when set — unique workflow DON with that don_family; use --shard-index when
//     several shard DONs share the family.
//  3. When the topology has exactly one workflow DON — use it.
//
// Multi-DON topologies must pass (1) or an unambiguous (2).
func (r *LocalCREStateResolver) ResolveWorkflowDONMetadata(sel workflowDONSelector) (*cre.DonMetadata, error) {
	wfDONs, err := r.topology.DonsMetadata.WorkflowDONs()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get workflow DONs from local CRE state")
	}
	if len(wfDONs) == 0 {
		return nil, errors.New("no workflow DON found in local CRE state")
	}

	// Explicit nodesets.name always wins — escape hatch for shards and ambiguous families.
	if name := strings.TrimSpace(sel.ExplicitName); name != "" {
		for _, wf := range wfDONs {
			if wf.Name == name {
				return wf, nil
			}
		}
		return nil, fmt.Errorf("workflow DON %q not found in local CRE state", name)
	}

	if family := strings.TrimSpace(sel.DonFamily); family != "" {
		return resolveWorkflowDONByFamily(wfDONs, family, sel.ShardIndex)
	}

	// Single-DON CI topologies need no selector flags.
	if len(wfDONs) == 1 {
		return wfDONs[0], nil
	}

	return nil, fmt.Errorf(
		"local CRE state has %d workflow DONs; set --don-family (and --shard-index for sharded topologies) or --workflow-don-name",
		len(wfDONs),
	)
}

// resolveWorkflowDONByFamily finds the workflow DON for a don_family CLI value.
//
// When shardIndex is non-nil, matches nodesets.shard_index (shard0/shard1 under the same family).
// When shardIndex is nil, succeeds only if exactly one workflow DON has that family, unless all
// matches are shard DONs — then --shard-index is required.
func resolveWorkflowDONByFamily(wfDONs []*cre.DonMetadata, family string, shardIndex *uint) (*cre.DonMetadata, error) {
	matches := make([]*cre.DonMetadata, 0, 1)
	for _, wf := range wfDONs {
		if wf.DonFamily == family {
			matches = append(matches, wf)
		}
	}
	if len(matches) == 0 {
		return nil, fmt.Errorf("no workflow DON with don_family %q in local CRE state", family)
	}

	if shardIndex != nil {
		// Narrow to the shard DON when several nodesets share don_family (shard0, shard1, …).
		filtered := make([]*cre.DonMetadata, 0, 1)
		for _, wf := range matches {
			if wf.ShardIndex == *shardIndex {
				filtered = append(filtered, wf)
			}
		}
		switch len(filtered) {
		case 1:
			return filtered[0], nil
		case 0:
			return nil, fmt.Errorf(
				"no workflow DON with don_family %q and shard_index %d; set --workflow-don-name",
				family,
				*shardIndex,
			)
		default:
			return nil, fmt.Errorf(
				"don_family %q and shard_index %d match multiple workflow DONs (%s); set --workflow-don-name",
				family,
				*shardIndex,
				strings.Join(workflowDONNames(filtered), ", "),
			)
		}
	}

	switch len(matches) {
	case 1:
		return matches[0], nil
	default:
		// e.g. shard0 + shard1 both test-don-family — disambiguate by shard index or name.
		if workflowDONsAreAllSharded(matches) {
			return nil, fmt.Errorf(
				"don_family %q matches multiple shard workflow DONs (%s); set --shard-index or --workflow-don-name",
				family,
				strings.Join(workflowDONNames(matches), ", "),
			)
		}
		// Multiple non-shard workflow DONs sharing a family — rare; require explicit name.
		return nil, fmt.Errorf(
			"don_family %q matches multiple workflow DONs (%s); set --workflow-don-name",
			family,
			strings.Join(workflowDONNames(matches), ", "),
		)
	}
}

// workflowDONsAreAllSharded reports whether every DON in the list is a shard workflow DON.
func workflowDONsAreAllSharded(wfDONs []*cre.DonMetadata) bool {
	for _, wf := range wfDONs {
		if !wf.IsShardDON() {
			return false
		}
	}
	return true
}

// workflowContainerPatternForDON returns the default docker cp substring for a workflow DON's nodes.
//
// Matches how env start names containers: fmt.Sprintf("%s-node%d", nodesets.name, nodeIndex).
// Substring match covers all nodes in the nodeset (e.g. feeds-zone-a-node matches feeds-zone-a-node0).
func workflowContainerPatternForDON(don *cre.DonMetadata) string {
	return don.Name + "-node"
}

// workflowDONNames extracts nodesets.name values for error messages.
func workflowDONNames(wfDONs []*cre.DonMetadata) []string {
	names := make([]string, len(wfDONs))
	for i, wf := range wfDONs {
		names[i] = wf.Name
	}
	return names
}
