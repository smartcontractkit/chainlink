// Per-don_family gateway↔workflow pairing for local CRE topologies.
//
// Every topology with http-actions gateway wiring declares don_family on workflow and gateway
// nodesets (N >= 1 families). env start validates the graph and scopes gateway connectors,
// gateway worker jobs, cap-reg/workflow-registry families, and deploy resolvers by family.
//
// Matching is by don_family string equality — not nodeset order, zone name, or numeric index.
package cre

import (
	"fmt"
	"slices"
	"strings"
)

// DonFamilyGatewayPair is one resolved workflow DON linked to one gateway DON sharing don_family.
type DonFamilyGatewayPair struct {
	DonFamily       string
	WorkflowDONName string
	GatewayDONName  string
}

// donFamilyPairingState is the validated pairing graph built at env start.
type donFamilyPairingState struct {
	gatewayDONNamesByFamily  map[string][]string // don_family → gateway nodesets.name
	workflowDONNamesByFamily map[string][]string // don_family → workflow nodesets.name
	pairs                    []DonFamilyGatewayPair
}

// initDonFamilyGatewayPairing validates gateway↔workflow don_family wiring when the topology
// uses http-actions gateway. No-op when RequiresGateway() is false. Sets t.donFamilyPairing
// used by GatewayConnectorsForDonFamily and GatewayServiceConfigsForDonFamily.
func (t *Topology) initDonFamilyGatewayPairing() error {
	if !t.DonsMetadata.RequiresGateway() {
		return nil
	}

	state, err := t.buildDonFamilyPairingState()
	if err != nil {
		return err
	}
	t.donFamilyPairing = state
	return nil
}

// buildDonFamilyPairingState validates and materializes the pairing graph.
//
// Steps:
//  1. Index gateway DONs by don_family (each must declare don_family).
//  2. For each workflow DON, require don_family and at least one gateway in that family.
//  3. Record workflow↔gateway pairs for logging and downstream lookups.
//
// Errors here fail env start instead of leaving cross-family gateway wiring at runtime.
func (t *Topology) buildDonFamilyPairingState() (*donFamilyPairingState, error) {
	wfDONs, err := t.DonsMetadata.WorkflowDONs()
	if err != nil {
		return nil, err
	}

	state := &donFamilyPairingState{
		gatewayDONNamesByFamily:  make(map[string][]string),
		workflowDONNamesByFamily: make(map[string][]string),
	}

	for _, d := range t.DonsMetadata.List() {
		if _, hasGateway := d.Gateway(); !hasGateway {
			continue
		}
		if d.DonFamily == "" {
			return nil, fmt.Errorf("gateway DON %q has no don_family; set nodesets.don_family on workflow and gateway nodesets", d.Name)
		}
		// A don_family may list multiple gateway nodesets; all are indexed here.
		state.gatewayDONNamesByFamily[d.DonFamily] = append(state.gatewayDONNamesByFamily[d.DonFamily], d.Name)
	}

	for _, wf := range wfDONs {
		if wf.DonFamily == "" {
			return nil, fmt.Errorf("workflow DON %q has no don_family; set nodesets.don_family on workflow and gateway nodesets", wf.Name)
		}
		if len(state.gatewayDONNamesByFamily[wf.DonFamily]) == 0 {
			return nil, fmt.Errorf("workflow DON %q is in don_family %q but no gateway DON is defined for that family", wf.Name, wf.DonFamily)
		}
		state.workflowDONNamesByFamily[wf.DonFamily] = append(state.workflowDONNamesByFamily[wf.DonFamily], wf.Name)
		// One workflow DON pairs with every gateway nodeset in its don_family.
		for _, gwName := range state.gatewayDONNamesByFamily[wf.DonFamily] {
			state.pairs = append(state.pairs, DonFamilyGatewayPair{
				DonFamily:       wf.DonFamily,
				WorkflowDONName: wf.Name,
				GatewayDONName:  gwName,
			})
		}
	}

	return state, nil
}

// DonFamilyGatewayPairings returns a copy of resolved workflow→gateway pairs.
func (t *Topology) DonFamilyGatewayPairings() []DonFamilyGatewayPair {
	if t.donFamilyPairing == nil {
		return nil
	}
	return slices.Clone(t.donFamilyPairing.pairs)
}

// WorkflowDONFamilies returns distinct don_family values from workflow DONs.
func (t *Topology) WorkflowDONFamilies() []string {
	if t.donFamilyPairing == nil {
		return nil
	}

	families := make([]string, 0, len(t.donFamilyPairing.workflowDONNamesByFamily))
	for family := range t.donFamilyPairing.workflowDONNamesByFamily {
		families = append(families, family)
	}
	return families
}

// DonFamilyGatewayPairingSummary returns a one-line log of resolved pairs, or "" when no gateway topology.
func (t *Topology) DonFamilyGatewayPairingSummary() string {
	pairs := t.DonFamilyGatewayPairings()
	if len(pairs) == 0 {
		return ""
	}

	parts := make([]string, 0, len(pairs))
	for _, pair := range pairs {
		parts = append(parts, fmt.Sprintf("%s → %s (don_family=%s)", pair.WorkflowDONName, pair.GatewayDONName, pair.DonFamily))
	}
	return "Gateway don_family pairs: " + strings.Join(parts, ", ")
}

// GatewayConnectorsForDonFamily returns gateway connector configs for workflow node TOML injection.
//
// Returns empty when donFamily is unknown or the topology has no gateway pairing state.
func (t *Topology) GatewayConnectorsForDonFamily(donFamily string) GatewayConnectors {
	if t.GatewayConnectors == nil || t.donFamilyPairing == nil {
		return GatewayConnectors{}
	}
	if donFamily == "" {
		return GatewayConnectors{}
	}

	configs := make([]*DonGatewayConfiguration, 0, len(t.donFamilyPairing.gatewayDONNamesByFamily[donFamily]))
	for _, gwName := range t.donFamilyPairing.gatewayDONNamesByFamily[donFamily] {
		if cfg, ok := t.gatewayConnectorsByDon[gwName]; ok {
			configs = append(configs, cfg)
		}
	}
	return GatewayConnectors{Configurations: configs}
}

// GatewayServiceConfigsForDonFamily rewrites each gateway service entry's DON list to workflow
// DON names in donFamily. Gateway workers use this so each gateway only handles its family's workflows.
func (t *Topology) GatewayServiceConfigsForDonFamily(donFamily string, services []GatewayServiceConfig) []GatewayServiceConfig {
	if t.donFamilyPairing == nil || donFamily == "" {
		return services
	}

	workflowNames := t.donFamilyPairing.workflowDONNamesByFamily[donFamily]
	if len(workflowNames) == 0 {
		return services
	}

	out := make([]GatewayServiceConfig, len(services))
	for i, svc := range services {
		out[i] = svc
		// Replace global DON list with only workflow nodesets in this family.
		out[i].DONs = slices.Clone(workflowNames)
	}
	return out
}

// DonFamilyForDON returns nodesets.don_family for a topology DON name, or "" if unknown.
func (t *Topology) DonFamilyForDON(donName string) string {
	if d := t.donByName(donName); d != nil {
		return d.DonFamily
	}
	return ""
}
