package cre

import (
	"fmt"
	"slices"
	"strings"
)

// DonFamilyGatewayPair links a workflow DON to a gateway DON in the same DON family.
type DonFamilyGatewayPair struct {
	DonFamily       string
	WorkflowDONName string
	GatewayDONName  string
}

// donFamilyPairingState holds resolved workflow↔gateway pairings keyed by don_family.
// Built at env start when pairing is enabled; legacy topologies leave Topology.donFamilyPairing nil.
type donFamilyPairingState struct {
	gatewayDONNamesByFamily  map[string][]string
	workflowDONNamesByFamily map[string][]string
	pairs                    []DonFamilyGatewayPair
}

func (t *Topology) initDonFamilyGatewayPairing() error {
	// Opt-in: pairing activates only when at least one workflow/gateway nodeset sets don_family.
	// Legacy single-gateway topologies skip pairing resolution and keep all-to-all wiring.
	t.donFamilyPairingEnabled = t.computeDonFamilyGatewayPairingEnabled()
	if !t.donFamilyPairingEnabled {
		return nil
	}

	state, err := t.buildDonFamilyPairingState()
	if err != nil {
		return err
	}
	t.donFamilyPairing = state
	return nil
}

// DonFamilyGatewayPairingEnabled is true when any workflow or gateway DON has don_family set.
// When false, legacy all-to-all gateway wiring is used (single shared gateway topologies).
// Once pairing is enabled, every workflow and gateway DON must declare don_family or env start fails.
func (t *Topology) DonFamilyGatewayPairingEnabled() bool {
	return t.donFamilyPairingEnabled
}

func (t *Topology) computeDonFamilyGatewayPairingEnabled() bool {
	if !t.DonsMetadata.RequiresGateway() {
		return false
	}
	// Any don_family on a workflow or gateway DON enables pairing for the whole topology.
	for _, d := range t.DonsMetadata.List() {
		if d.DonFamily == "" {
			continue
		}
		if d.IsWorkflowDON() {
			return true
		}
		if _, hasGateway := d.Gateway(); hasGateway {
			return true
		}
	}
	return false
}

// buildDonFamilyPairingState validates that every workflow/gateway DON declares don_family,
// matches workflow DONs to gateway DONs with the same don_family, and stores the result for
// scoped connector/service lookups. Failures surface at env start rather than as silent
// cross-family gateway wiring at runtime.
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

// DonFamilyGatewayPairings returns workflow→gateway pairs grouped by don_family.
func (t *Topology) DonFamilyGatewayPairings() []DonFamilyGatewayPair {
	if t.donFamilyPairing == nil {
		return nil
	}
	return slices.Clone(t.donFamilyPairing.pairs)
}

// WorkflowDONFamilies returns distinct non-empty don_family values from workflow DONs.
// Legacy topologies without don_family leave this empty; callers fall back to DefaultDONFamily.
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

// DonFamilyGatewayPairingSummary returns a human-readable summary of resolved workflow→gateway pairs.
// Empty when pairing is disabled or no pairs were resolved.
func (t *Topology) DonFamilyGatewayPairingSummary() string {
	pairs := t.DonFamilyGatewayPairings()
	if len(pairs) == 0 {
		return ""
	}

	parts := make([]string, 0, len(pairs))
	for _, pair := range pairs {
		parts = append(parts, fmt.Sprintf("%s → %s (don_family=%s)", pair.WorkflowDONName, pair.GatewayDONName, pair.DonFamily))
	}
	return "Gateway don_family pairing enabled: " + strings.Join(parts, ", ")
}

// GatewayConnectorsForDonFamily returns gateway connectors scoped to donFamily when pairing is
// enabled; otherwise returns all connectors (legacy all-to-all wiring).
func (t *Topology) GatewayConnectorsForDonFamily(donFamily string) GatewayConnectors {
	if t.GatewayConnectors == nil {
		return GatewayConnectors{}
	}
	if !t.donFamilyPairingEnabled {
		return *t.GatewayConnectors
	}
	if donFamily == "" || t.donFamilyPairing == nil {
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

// GatewayServiceConfigsForDonFamily scopes each service entry's DON list to workflow DONs in
// donFamily when pairing is enabled.
func (t *Topology) GatewayServiceConfigsForDonFamily(donFamily string, services []GatewayServiceConfig) []GatewayServiceConfig {
	if !t.donFamilyPairingEnabled || donFamily == "" || t.donFamilyPairing == nil {
		return services
	}

	workflowNames := t.donFamilyPairing.workflowDONNamesByFamily[donFamily]
	if len(workflowNames) == 0 {
		return services
	}

	out := make([]GatewayServiceConfig, len(services))
	for i, svc := range services {
		out[i] = svc
		out[i].DONs = slices.Clone(workflowNames)
	}
	return out
}

// DonFamilyForDON returns don_family for a DON name in the topology.
func (t *Topology) DonFamilyForDON(donName string) string {
	if d := t.donByName(donName); d != nil {
		return d.DonFamily
	}
	return ""
}
