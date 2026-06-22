package environment

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
)

// workflowDONSelector identifies which workflow DON to target during local deploy.
// ContainerPattern is used for Docker artifact copy (substring match) and may also
// identify a DON when it maps unambiguously. ExplicitName (--workflow-don-name) is
// the preferred selector for multi-DON topologies.
type workflowDONSelector struct {
	ExplicitName     string
	ContainerPattern string
}

// ResolveWorkflowDONMetadata selects the workflow DON for registry registration and
// vault DB polling. Multi-DON topologies require --workflow-don-name or a container
// pattern that matches exactly one DON.
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

// containerPatternMatchesDON reports whether a Docker container name pattern identifies
// a specific workflow DON. This is stricter than Docker substring matching: the pattern
// must equal the DON name, equal "<don>-node", or start with "<don>-node" followed by
// an optional node suffix (e.g. feeds-zone-a-node-0).
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
