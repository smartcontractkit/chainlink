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
		WorkflowDONIDs: []uint64{},
		DonsMetadata:   donsMetadata,
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

	return topology, nil
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
