package onchain

import (
	stderrors "errors"
	"fmt"

	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/core/scripts/cre/reconciler/internal/domain"
	"github.com/smartcontractkit/chainlink/deployment"
)

// ValidateNodeLabels checks every chart-declared node's JD labels (p2p_id,
// environment, type) against what's expected, per phase_plan.md D5: the chart
// remains the source of truth for DON membership/role/namespace, and this
// validates JD's own view is consistent with it. Runs unconditionally, using a
// lightweight JD-only client (buildOffchainClient) rather than the full CLDF
// environment, so it can run before/independently of on-chain work — including
// on a run where on-chain work is already complete and Apply would otherwise
// be skipped entirely.
//
// Requires discovery to have already populated state.NodeRuntime[name].CSAKey
// (from this run or a prior one) for every chart node.
func ValidateNodeLabels(desired *domain.DesiredState, cv *domain.ChartValues, state *domain.StateFile) error {
	offchainClient, err := buildOffchainClient(desired.JD)
	if err != nil {
		return errors.Wrap(err, "failed to build JD client for label preflight")
	}
	if offchainClient == nil {
		// JD isn't configured (or the access token isn't set) — matches
		// Reconciler.requireJDAccessToken's existing tolerance of a JD-less run.
		return nil
	}

	nodeByCSA := make(map[string]domain.ChartNodeInfo, len(cv.Nodes))
	lookupIDs := make([]string, 0, len(cv.Nodes))
	var errs []error
	for _, node := range cv.Nodes {
		info, ok := state.NodeRuntime[node.Name]
		if !ok || info.CSAKey == "" {
			errs = append(errs, fmt.Errorf("node %s: no discovered CSA key — cannot validate JD labels", node.Name))
			continue
		}
		lookupIDs = append(lookupIDs, info.CSAKey)
		nodeByCSA[info.CSAKey] = node
	}
	if len(errs) > 0 {
		return stderrors.Join(errs...)
	}

	nodes, err := deployment.NodeInfo(lookupIDs, offchainClient)
	if err != nil {
		return errors.Wrap(err, "failed to load node info from JD for label preflight")
	}

	seen := make(map[string]bool, len(nodes))
	for _, n := range nodes {
		chartNode, ok := nodeByCSA[n.CSAKey]
		if !ok {
			continue
		}
		seen[chartNode.Name] = true

		labels := make(map[string]string, len(n.Labels))
		for _, l := range n.Labels {
			if l.Value != nil {
				labels[l.Key] = *l.Value
			}
		}

		if labels["p2p_id"] == "" {
			errs = append(errs, fmt.Errorf("node %s: JD label \"p2p_id\" is missing", chartNode.Name))
		}
		if got := labels["environment"]; got != desired.JD.Environment {
			errs = append(errs, fmt.Errorf("node %s: JD label \"environment\" is %q, want %q", chartNode.Name, got, desired.JD.Environment))
		}
		wantType := expectedJDTypeLabel(chartNode.NodeType)
		if got := labels["type"]; got != wantType {
			errs = append(errs, fmt.Errorf("node %s: JD label \"type\" is %q, want %q", chartNode.Name, got, wantType))
		}
	}
	for _, node := range cv.Nodes {
		if !seen[node.Name] {
			errs = append(errs, fmt.Errorf("node %s: not found in JD (looked up by discovered CSA key)", node.Name))
		}
	}

	return stderrors.Join(errs...)
}

// expectedJDTypeLabel maps a chart-declared node role to the JD "type" label
// value it should have. Worker/standard nodes are labeled "plugin" in JD, not
// "standard" — see system-tests/lib/cre/don.go's Role constants.
func expectedJDTypeLabel(role domain.NodeRole) string {
	switch role {
	case domain.RoleBootstrap:
		return "bootstrap"
	case domain.RoleGateway:
		return "gateway"
	default:
		return "plugin"
	}
}
