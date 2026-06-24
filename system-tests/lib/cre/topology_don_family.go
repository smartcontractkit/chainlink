// Per-don_family gateway↔workflow pairing for local CRE topologies.
//
// Every topology with http-actions gateway wiring declares don_family on all nodesets (N >= 1
// families). env start validates the graph and scopes gateway connectors, gateway worker jobs,
// cap-reg/workflow-registry families, and deploy resolvers by family.
//
// Matching is by don_family string equality — not nodeset order or numeric index.
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

// gatewayDonFamilyPairingState is the validated pairing graph built at env start.
type gatewayDonFamilyPairingState struct {
	gatewayDONNamesByFamily  map[string][]string // don_family → gateway nodesets.name
	workflowDONNamesByFamily map[string][]string // don_family → workflow nodesets.name
	pairs                    []DonFamilyGatewayPair
}

// initDonFamilyGatewayPairing validates gateway↔workflow don_family wiring when the topology
// uses http-actions gateway. No-op when RequiresGateway() is false. Sets t.gatewayDonFamilyPairing
// used by GatewayConnectorsForDonFamily and GatewayServiceConfigsForGateway.
func (t *Topology) initDonFamilyGatewayPairing() error {
	if !t.DonsMetadata.RequiresGateway() {
		return nil
	}

	t.ensureGatewayConnectorIndex()

	state, err := t.buildDonFamilyPairingState()
	if err != nil {
		return err
	}
	t.gatewayDonFamilyPairing = state
	return nil
}

// EnsureGatewayDonFamilyPairing validates pairing and builds per-family lookup state.
func (t *Topology) EnsureGatewayDonFamilyPairing() error {
	return t.initDonFamilyGatewayPairing()
}

func (t *Topology) ensureGatewayConnectorIndex() {
	if t.GatewayConnectors == nil || len(t.GatewayConnectors.Configurations) == 0 {
		return
	}
	if t.gatewayConnectorsByDon == nil {
		t.gatewayConnectorsByDon = make(map[string]*DonGatewayConfiguration)
	}
	if len(t.gatewayConnectorsByDon) > 0 {
		return
	}

	gatewayCount := 0
	for _, d := range t.DonsMetadata.List() {
		if _, hasGateway := d.Gateway(); !hasGateway {
			continue
		}
		if gatewayCount >= len(t.GatewayConnectors.Configurations) {
			break
		}
		t.gatewayConnectorsByDon[d.Name] = t.GatewayConnectors.Configurations[gatewayCount]
		gatewayCount++
	}
}

// buildDonFamilyPairingState validates and materializes the pairing graph.
func (t *Topology) buildDonFamilyPairingState() (*gatewayDonFamilyPairingState, error) {
	wfDONs, err := t.DonsMetadata.WorkflowDONs()
	if err != nil {
		return nil, err
	}

	state := &gatewayDonFamilyPairingState{
		gatewayDONNamesByFamily:  make(map[string][]string),
		workflowDONNamesByFamily: make(map[string][]string),
	}

	for _, d := range t.DonsMetadata.List() {
		if _, hasGateway := d.Gateway(); !hasGateway {
			continue
		}
		if d.DonFamily == "" {
			return nil, fmt.Errorf("gateway DON %q has no don_family; set nodesets.don_family on every nodeset", d.Name)
		}
		state.gatewayDONNamesByFamily[d.DonFamily] = append(state.gatewayDONNamesByFamily[d.DonFamily], d.Name)
	}

	for _, wf := range wfDONs {
		if wf.DonFamily == "" {
			return nil, fmt.Errorf("workflow DON %q has no don_family; set nodesets.don_family on every nodeset", wf.Name)
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

// DonFamilyGatewayPairings returns a copy of resolved workflow→gateway pairs.
func (t *Topology) DonFamilyGatewayPairings() []DonFamilyGatewayPair {
	if t.gatewayDonFamilyPairing == nil {
		return nil
	}
	return slices.Clone(t.gatewayDonFamilyPairing.pairs)
}

// WorkflowDONFamilies returns distinct don_family values from workflow DONs.
func (t *Topology) WorkflowDONFamilies() []string {
	wfDONs, err := t.DonsMetadata.WorkflowDONs()
	if err != nil || len(wfDONs) == 0 {
		return nil
	}

	seen := make(map[string]struct{})
	families := make([]string, 0, len(wfDONs))
	for _, wf := range wfDONs {
		if wf.DonFamily == "" {
			continue
		}
		if _, ok := seen[wf.DonFamily]; ok {
			continue
		}
		seen[wf.DonFamily] = struct{}{}
		families = append(families, wf.DonFamily)
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
func (t *Topology) GatewayConnectorsForDonFamily(donFamily string) GatewayConnectors {
	if t.GatewayConnectors == nil || t.gatewayDonFamilyPairing == nil {
		return GatewayConnectors{}
	}

	configs := make([]*DonGatewayConfiguration, 0, len(t.gatewayDonFamilyPairing.gatewayDONNamesByFamily[donFamily]))
	for _, gwName := range t.gatewayDonFamilyPairing.gatewayDONNamesByFamily[donFamily] {
		if cfg, ok := t.gatewayConnectorsByDon[gwName]; ok {
			configs = append(configs, cfg)
		}
	}
	return GatewayConnectors{Configurations: configs}
}

// GatewayServiceConfigsForGateway scopes gateway worker service configs to DONs in the
// gateway nodeset's don_family (workflow and capabilities DONs, e.g. vault handler routing).
func (t *Topology) GatewayServiceConfigsForGateway(gatewayDONName string, services []GatewayServiceConfig) []GatewayServiceConfig {
	return t.gatewayServiceConfigsForDonFamily(t.DonFamilyForDON(gatewayDONName), services)
}

func (t *Topology) gatewayServiceConfigsForDonFamily(donFamily string, services []GatewayServiceConfig) []GatewayServiceConfig {
	if t.gatewayDonFamilyPairing == nil || donFamily == "" {
		return services
	}

	out := make([]GatewayServiceConfig, len(services))
	for i, svc := range services {
		out[i] = svc
		filtered := make([]string, 0, len(svc.DONs))
		for _, donName := range svc.DONs {
			if t.DonFamilyForDON(donName) == donFamily {
				filtered = append(filtered, donName)
			}
		}
		out[i].DONs = filtered
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
