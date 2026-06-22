package environment

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
)

// workflowDONSelector picks which workflow DON a deploy targets.
//
// Deploy uses two related but distinct inputs from workflow.go:
//   - ContainerPattern (--container-name-pattern): substring match for copying WASM into Docker containers.
//   - ExplicitName (--workflow-don-name): topology DON name for registry registration and vault DB polling.
//
// This type is passed to ResolveWorkflowDONMetadata. When both are set, ExplicitName wins for DON
// selection; ContainerPattern may still be used separately for artifact copy in deployWorkflow.
type workflowDONSelector struct {
	ExplicitName     string // nodesets.name, e.g. "feeds-zone-a"
	ContainerPattern string // e.g. "feeds-zone-a-node"; may infer the DON when it matches exactly one
}

// ResolveWorkflowDONMetadata returns the workflow DON metadata from the topology saved by env start.
//
// Used by workflow deploy (registry targets, vault DB port) and related resolvers. Selection priority:
//  1. --workflow-don-name when set — exact match on nodesets.name.
//  2. --container-name-pattern when set — map pattern to a DON via containerPatternMatchesDON; error if
//     zero matches (except single-DON legacy fallback) or more than one match.
//  3. When neither flag is set — OK only if the topology has exactly one workflow DON.
//
// Multi-DON topologies (feeds-zone-a + feeds-zone-b) must pass (1) or an unambiguous (2). Generic
// patterns like "workflow-node" are allowed for single-DON legacy CI only.
func (r *LocalCREStateResolver) ResolveWorkflowDONMetadata(sel workflowDONSelector) (*cre.DonMetadata, error) {
	wfDONs, err := r.topology.DonsMetadata.WorkflowDONs()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get workflow DONs from local CRE state")
	}
	if len(wfDONs) == 0 {
		return nil, errors.New("no workflow DON found in local CRE state")
	}

	if name := strings.TrimSpace(sel.ExplicitName); name != "" {
		for _, wf := range wfDONs {
			if wf.Name == name {
				return wf, nil
			}
		}
		return nil, fmt.Errorf("workflow DON %q not found in local CRE state", name)
	}

	if pattern := strings.TrimSpace(sel.ContainerPattern); pattern != "" {
		matches := workflowDONsMatchingContainerPattern(wfDONs, pattern)
		switch len(matches) {
		case 1:
			return matches[0], nil
		case 0:
			// Single-DON legacy: generic Docker patterns (e.g. "workflow-node") often do not embed the
			// nodeset name; allow the lone workflow DON when there is no ambiguity.
			if len(wfDONs) == 1 {
				return wfDONs[0], nil
			}
			return nil, fmt.Errorf(
				"container pattern %q does not identify a workflow DON among %d workflow DONs; set --workflow-don-name (e.g. feeds-zone-a)",
				pattern,
				len(wfDONs),
			)
		default:
			return nil, fmt.Errorf(
				"container pattern %q matches multiple workflow DONs (%s); set --workflow-don-name",
				pattern,
				strings.Join(workflowDONNames(matches), ", "),
			)
		}
	}

	if len(wfDONs) == 1 {
		return wfDONs[0], nil
	}

	return nil, fmt.Errorf(
		"local CRE state has %d workflow DONs; set --workflow-don-name (e.g. feeds-zone-a)",
		len(wfDONs),
	)
}

// containerPatternMatchesDON reports whether a container-name-pattern identifies a specific workflow DON.
//
// Stricter than the Docker copy substring match in deployWorkflow: the pattern must equal the DON
// name, equal "<don>-node", or start with "<don>-node" plus an optional suffix (e.g. feeds-zone-a-node-0).
// HasPrefix on "<don>-node" avoids false positives where a shorter DON name is a prefix of another
// (e.g. pattern "feeds-zone-a-node" must not match DON "feeds").
func containerPatternMatchesDON(containerPattern string, don *cre.DonMetadata) bool {
	if containerPattern == don.Name {
		return true
	}
	if strings.TrimSuffix(containerPattern, "-node") == don.Name {
		return true
	}
	return strings.HasPrefix(containerPattern, don.Name+"-node")
}

func workflowDONsMatchingContainerPattern(wfDONs []*cre.DonMetadata, containerPattern string) []*cre.DonMetadata {
	matches := make([]*cre.DonMetadata, 0, 1)
	for _, wf := range wfDONs {
		if containerPatternMatchesDON(containerPattern, wf) {
			matches = append(matches, wf)
		}
	}
	return matches
}

func workflowDONNames(wfDONs []*cre.DonMetadata) []string {
	names := make([]string, len(wfDONs))
	for i, wf := range wfDONs {
		names[i] = wf.Name
	}
	return names
}
