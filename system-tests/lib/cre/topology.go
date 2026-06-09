package cre

import (
	"fmt"
	"slices"
	"strings"

	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/deployment/cre/jobs/pkg"
	libc "github.com/smartcontractkit/chainlink/system-tests/lib/conversions"
	"github.com/smartcontractkit/chainlink/system-tests/lib/infra"
)

const (
	OCRPeeringPort          = 5001
	CapabilitiesPeeringPort = 6690
)

type Topology struct {
	WorkflowDONIDs        []uint64               `toml:"workflow_don_ids" json:"workflow_don_ids"`
	DonsMetadata          *DonsMetadata          `toml:"dons_metadata" json:"dons_metadata"`
	GatewayServiceConfigs []GatewayServiceConfig `toml:"gateway_service_configs" json:"gateway_service_configs"`
	GatewayConnectors     *GatewayConnectors     `toml:"gateway_connectors" json:"gateway_connectors"`

	gatewayConnectorsByDon map[string]*DonGatewayConfiguration
}

// DonFamilyGatewayPair links a workflow DON to a gateway DON in the same DON family.
type DonFamilyGatewayPair struct {
	DonFamily       string
	WorkflowDONName string
	GatewayDONName  string
}

func NewTopology(nodeSet []*NodeSet, provider infra.Provider, capabilityConfigs map[CapabilityFlag]CapabilityConfig) (*Topology, error) {
	dm := make([]*DonMetadata, len(nodeSet))
	for i := range nodeSet {
		// Use ContractDonID from NodeSet when set (resolved from Capabilities Registry contract).
		// Otherwise use optimistic i+1; the ID may be overwritten later when resolving from the contract.
		id := nodeSet[i].ContractDonID
		if id == 0 {
			id = libc.MustSafeUint64FromInt(i + 1)
		}
		d, err := NewDonMetadata(nodeSet[i], id, provider, capabilityConfigs)
		if err != nil {
			return nil, fmt.Errorf("failed to create DON metadata: %w", err)
		}
		dm[i] = d
	}

	donsMetadata, err := NewDonsMetadata(dm, provider)
	if err != nil {
		return nil, fmt.Errorf("failed to create DONs metadata: %w", err)
	}

	wfDONs, err := donsMetadata.WorkflowDONs()
	if err != nil {
		return nil, fmt.Errorf("failed to find any workflow DONs: %w", err)
	}

	topology := &Topology{
		WorkflowDONIDs:         []uint64{},
		DonsMetadata:           donsMetadata,
		gatewayConnectorsByDon: make(map[string]*DonGatewayConfiguration),
	}

	donNames := make([]string, 0, len(wfDONs))
	for _, wfDON := range wfDONs {
		donNames = append(donNames, wfDON.Name)
		topology.WorkflowDONIDs = append(topology.WorkflowDONIDs, wfDON.ID)
	}

	topology.GatewayServiceConfigs = append(topology.GatewayServiceConfigs, GatewayServiceConfig{
		ServiceName: pkg.ServiceNameWorkflows,
		Handlers:    []string{pkg.GatewayHandlerTypeWebAPICapabilities},
		DONs:        donNames,
	})

	if donsMetadata.RequiresGateway() {
		topology.GatewayConnectors = NewGatewayConnectorOutput()
		gatewayCount := 0
		for _, d := range donsMetadata.List() {
			if _, hasGateway := d.Gateway(); hasGateway {
				gc, err := d.GatewayConfig(provider, gatewayCount)
				if err != nil {
					return nil, fmt.Errorf("failed to get gateway config for DON %s: %w", d.Name, err)
				}
				topology.GatewayConnectors.Configurations = append(topology.GatewayConnectors.Configurations, gc)
				topology.gatewayConnectorsByDon[d.Name] = gc
				gatewayCount++
			}
		}
	}

	bootstrapNodesFound := 0
	for _, don := range topology.DonsMetadata.List() {
		if _, isBootstrap := don.Bootstrap(); isBootstrap {
			bootstrapNodesFound++
		}
	}

	if bootstrapNodesFound == 0 {
		return nil, errors.New("no bootstrap nodes found in topology. At least one bootstrap node is required")
	}

	if bootstrapNodesFound > 1 {
		return nil, errors.New("multiple bootstrap nodes found in topology. Only one bootstrap node is supported due to the limitations of the local environment")
	}

	if err := topology.validateDonFamilyGatewayPairing(); err != nil {
		return nil, err
	}

	return topology, nil
}

func (t *Topology) donByName(name string) *DonMetadata {
	for _, d := range t.DonsMetadata.List() {
		if d.Name == name {
			return d
		}
	}
	return nil
}

// DonFamilyGatewayPairingEnabled is true when any workflow or gateway DON has don_family set.
// When false, legacy all-to-all gateway wiring is used (single shared gateway topologies).
func (t *Topology) DonFamilyGatewayPairingEnabled() bool {
	if !t.DonsMetadata.RequiresGateway() {
		return false
	}
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

func (t *Topology) validateDonFamilyGatewayPairing() error {
	if !t.DonFamilyGatewayPairingEnabled() {
		return nil
	}

	wfDONs, err := t.DonsMetadata.WorkflowDONs()
	if err != nil {
		return err
	}

	gatewayFamilies := make(map[string][]string)
	for _, d := range t.DonsMetadata.List() {
		if _, hasGateway := d.Gateway(); !hasGateway {
			continue
		}
		if d.DonFamily == "" {
			return fmt.Errorf("gateway DON %q has no don_family; set nodesets.don_family on workflow and gateway nodesets", d.Name)
		}
		gatewayFamilies[d.DonFamily] = append(gatewayFamilies[d.DonFamily], d.Name)
	}

	for _, wf := range wfDONs {
		if wf.DonFamily == "" {
			return fmt.Errorf("workflow DON %q has no don_family; set nodesets.don_family on workflow and gateway nodesets", wf.Name)
		}
		if len(gatewayFamilies[wf.DonFamily]) == 0 {
			return fmt.Errorf("workflow DON %q is in don_family %q but no gateway DON is defined for that family", wf.Name, wf.DonFamily)
		}
	}

	return nil
}

// GatewayDonFamilyPairings returns workflow→gateway pairs grouped by don_family.
func (t *Topology) GatewayDonFamilyPairings() []DonFamilyGatewayPair {
	if !t.DonFamilyGatewayPairingEnabled() {
		return nil
	}

	wfDONs, err := t.DonsMetadata.WorkflowDONs()
	if err != nil {
		return nil
	}

	var pairs []DonFamilyGatewayPair
	for _, wf := range wfDONs {
		for _, gw := range t.DonsMetadata.List() {
			if _, hasGateway := gw.Gateway(); !hasGateway {
				continue
			}
			if gw.DonFamily == wf.DonFamily {
				pairs = append(pairs, DonFamilyGatewayPair{
					DonFamily:       wf.DonFamily,
					WorkflowDONName: wf.Name,
					GatewayDONName:  gw.Name,
				})
			}
		}
	}
	return pairs
}

// LogGatewayDonFamilyPairing prints resolved workflow→gateway pairs at env start.
func (t *Topology) LogGatewayDonFamilyPairing() {
	pairs := t.GatewayDonFamilyPairings()
	if len(pairs) == 0 {
		return
	}

	parts := make([]string, 0, len(pairs))
	for _, pair := range pairs {
		parts = append(parts, fmt.Sprintf("%s → %s (don_family=%s)", pair.WorkflowDONName, pair.GatewayDONName, pair.DonFamily))
	}
	fmt.Printf("Gateway don_family pairing enabled: %s\n", strings.Join(parts, ", "))
}

// GatewayConnectorsForDonFamily returns gateway connectors for workflow DONs in a don_family.
func (t *Topology) GatewayConnectorsForDonFamily(donFamily string) GatewayConnectors {
	if t.GatewayConnectors == nil {
		return GatewayConnectors{}
	}
	if !t.DonFamilyGatewayPairingEnabled() {
		return *t.GatewayConnectors
	}
	if donFamily == "" {
		return GatewayConnectors{}
	}

	configs := make([]*DonGatewayConfiguration, 0)
	for _, d := range t.DonsMetadata.List() {
		if _, hasGateway := d.Gateway(); !hasGateway || d.DonFamily != donFamily {
			continue
		}
		if cfg, ok := t.gatewayConnectorsByDon[d.Name]; ok {
			configs = append(configs, cfg)
		}
	}
	return GatewayConnectors{Configurations: configs}
}

// GatewayServiceConfigsForDonFamily scopes service DON lists to workflow DONs in don_family.
func (t *Topology) GatewayServiceConfigsForDonFamily(donFamily string, services []GatewayServiceConfig) []GatewayServiceConfig {
	if !t.DonFamilyGatewayPairingEnabled() {
		return services
	}
	if donFamily == "" {
		return services
	}

	wfDONs, err := t.DonsMetadata.WorkflowDONs()
	if err != nil {
		return services
	}

	workflowNames := make([]string, 0)
	for _, wf := range wfDONs {
		if wf.DonFamily == donFamily {
			workflowNames = append(workflowNames, wf.Name)
		}
	}
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

func (t *Topology) NodeSets() []*NodeSet {
	sets := make([]*NodeSet, len(t.DonsMetadata.List()))
	for i, d := range t.DonsMetadata.List() {
		ns := d.MustNodeSet()
		sets[i] = ns
	}
	return sets
}

func (t *Topology) DonsMetadataWithFlag(flag CapabilityFlag) []*DonMetadata {
	donsMetadata := make([]*DonMetadata, 0)
	for _, donMetadata := range t.DonsMetadata.List() {
		if !donMetadata.HasFlag(flag) {
			continue
		}
		donsMetadata = append(donsMetadata, donMetadata)
	}

	return donsMetadata
}

// BootstrapNode returns the metadata for the node that should be used as the bootstrap node for P2P peering
// Currently only one bootstrap is supported.
func (t *Topology) Bootstrap() (*NodeMetadata, bool) {
	return t.DonsMetadata.Bootstrap()
}

// AddGatewayHandlers adds the given handler names for the given DON.
// It updates service-centric GatewayServiceConfigs.
func (t *Topology) AddGatewayHandlers(donMetadata DonMetadata, handlers []string) error {
	for _, handlerName := range handlers {
		svcName := pkg.HandlerServiceName(handlerName)

		svcIdx := -1
		for i, svc := range t.GatewayServiceConfigs {
			if strings.EqualFold(svc.ServiceName, svcName) {
				svcIdx = i
				break
			}
		}

		if svcIdx == -1 {
			t.GatewayServiceConfigs = append(t.GatewayServiceConfigs, GatewayServiceConfig{
				ServiceName: svcName,
				Handlers:    []string{handlerName},
				DONs:        []string{donMetadata.Name},
			})
			continue
		}

		if !slices.ContainsFunc(t.GatewayServiceConfigs[svcIdx].Handlers, func(h string) bool {
			return strings.EqualFold(h, handlerName)
		}) {
			t.GatewayServiceConfigs[svcIdx].Handlers = append(t.GatewayServiceConfigs[svcIdx].Handlers, handlerName)
		}

		if !slices.Contains(t.GatewayServiceConfigs[svcIdx].DONs, donMetadata.Name) {
			t.GatewayServiceConfigs[svcIdx].DONs = append(t.GatewayServiceConfigs[svcIdx].DONs, donMetadata.Name)
		}
	}

	return nil
}

type PeeringNode interface {
	GetHost() string
	PeerID() string
}

func PeeringCfgs(bt PeeringNode) (CapabilitiesPeeringData, OCRPeeringData, error) {
	p := strings.TrimPrefix(bt.PeerID(), "p2p_")
	if p == "" {
		return CapabilitiesPeeringData{}, OCRPeeringData{}, errors.New("cannot create peering configs, node has no P2P key")
	}
	return CapabilitiesPeeringData{
			GlobalBootstraperPeerID: p,
			GlobalBootstraperHost:   bt.GetHost(),
			Port:                    CapabilitiesPeeringPort,
		}, OCRPeeringData{
			OCRBootstraperPeerID: p,
			OCRBootstraperHost:   bt.GetHost(),
			Port:                 OCRPeeringPort,
		}, nil
}
