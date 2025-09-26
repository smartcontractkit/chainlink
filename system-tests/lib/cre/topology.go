package cre

import (
	"fmt"

	libc "github.com/smartcontractkit/chainlink/system-tests/lib/conversions"

	"github.com/smartcontractkit/chainlink/system-tests/lib/infra"
)

const (
	OCRPeeringPort          = 5001
	CapabilitiesPeeringPort = 6690
)

var (
	NodeTypeKey            = "type"
	HostLabelKey           = "host"
	IndexKey               = "node_index"
	ExtraRolesKey          = "extra_roles"
	NodeIDKey              = "node_id"
	NodeOCRFamiliesKey     = "node_ocr_families"
	NodeOCR2KeyBundleIDKey = "ocr2_key_bundle_id"
	NodeP2PIDKey           = "p2p_id"
	NodeDKGRecipientKey    = "dkg_recipient_key"
	DONIDKey               = "don_id"
	EnvironmentKey         = "environment"
	ProductKey             = "product"
	DONNameKey             = "don_name"
)

type Topology struct {
	WorkflowDONID           uint64                  `toml:"workflow_don_id" json:"workflow_don_id"`
	DonsMetadata            []*DonMetadata          `toml:"dons_metadata" json:"dons_metadata"`
	CapabilitiesPeeringData CapabilitiesPeeringData `toml:"capabilities_peering_data" json:"capabilities_peering_data"`
	OCRPeeringData          OCRPeeringData          `toml:"ocr_peering_data" json:"ocr_peering_data"`
	GatewayConnectorOutput  *GatewayConnectorOutput `toml:"gateway_connector_output" json:"gateway_connector_output"`
}

func NewTopology(nodeSetInput []*CapabilitiesAwareNodeSet, provider infra.Provider) (*Topology, error) {
	// TODO this setup is awkward, consider an withInfra opt to constructor
	dm := make([]*DonMetadata, len(nodeSetInput))
	for i := range nodeSetInput {
		// TODO take more care about the ID assignment, it should match what the capabilities registry will assign
		// currently we optimistically set the id to the that which the capabilities registry will assign it
		d := NewDonMetadata(nodeSetInput[i], libc.MustSafeUint64FromInt(i+1))
		d.labelNodes(provider)
		dm[i] = d
	}

	donsMetadata, err := NewDonsMetadata(dm, provider)
	if err != nil {
		return nil, fmt.Errorf("failed to create DONs metadata: %w", err)
	}

	wfDon, err := donsMetadata.GetWorkflowDON()
	if err != nil {
		return nil, fmt.Errorf("failed to get workflow DON: %w", err)
	}
	bt, err := wfDon.GetBootstrapNode()
	if err != nil {
		return nil, fmt.Errorf("failed to get workflow DON bootstrap node: %w", err)
	}
	capPeeringCfg, ocrPeeringCfg, err := peeringCfgs(bt)
	if err != nil {
		return nil, fmt.Errorf("failed to get peering data: %w", err)
	}

	topology := &Topology{
		WorkflowDONID: wfDon.ID,
		DonsMetadata:  dm,
		// TODO move to DON level
		CapabilitiesPeeringData: capPeeringCfg,
		OCRPeeringData:          ocrPeeringCfg,
	}

	if donsMetadata.GatewayRequired() {
		topology.GatewayConnectorOutput = NewGatewayConnectorOutput()
		for _, d := range donsMetadata.List() {
			if d.ContainsGatewayNode() {
				gc, err := d.GatewayConfig(provider)
				if err != nil {
					return nil, fmt.Errorf("failed to get gateway config for DON %s: %w", d.Name, err)
				}
				topology.GatewayConnectorOutput.Configurations = append(topology.GatewayConnectorOutput.Configurations, gc)
			}
		}
	}

	return topology, nil
}

// TODO i don't this think is actually used in so much as it seems to be overwritten in other places
// The keys should be set by the secret generator and are independent of the topology abstraction
func peeringCfgs(bt *NodeMetadata) (CapabilitiesPeeringData, OCRPeeringData, error) {
	return CapabilitiesPeeringData{
			GlobalBootstraperPeerID: bt.P2P,
			GlobalBootstraperHost:   bt.Host,
			Port:                    CapabilitiesPeeringPort,
		}, OCRPeeringData{
			OCRBootstraperPeerID: bt.P2P,
			OCRBootstraperHost:   bt.Host,
			Port:                 OCRPeeringPort,
		}, nil
}
