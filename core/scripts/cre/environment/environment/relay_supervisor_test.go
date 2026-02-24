package environment

import (
	"testing"

	ns "github.com/smartcontractkit/chainlink-testing-framework/framework/components/simple_node_set"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	envconfig "github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/config"
)

func TestRelaySpecsFromConfig_AddsBootstrapPeeringPortForRemoteToLocalMixedDONs(t *testing.T) {
	cfg := &envconfig.Config{
		NodeSets: []*cre.NodeSet{
			{
				Input: &ns.Input{
					Name:               "workflow",
					Nodes:              2,
					HTTPPortRangeStart: 10100,
				},
				Placement: "local",
				NodeSpecs: []*cre.NodeSpecWithRole{
					{Roles: []string{cre.BootstrapNode}},
				},
			},
			{
				Input: &ns.Input{
					Name:  "capabilities",
					Nodes: 1,
				},
				Placement: "remote",
				NodeSpecs: []*cre.NodeSpecWithRole{
					{Roles: []string{cre.WorkerNode}},
				},
			},
		},
	}

	specs := relaySpecsFromConfig(cfg)
	got := map[int]bool{}
	for _, spec := range specs {
		got[spec.Port] = true
	}
	if !got[14100] || !got[14101] {
		t.Fatalf("expected relay specs to include per-node OCR relay ports 14100/14101, got %#v", specs)
	}
}

func TestRelaySpecsFromConfig_DoesNotAddBootstrapWhenNoRemoteNodeSets(t *testing.T) {
	cfg := &envconfig.Config{
		NodeSets: []*cre.NodeSet{
			{
				Placement: "local",
				NodeSpecs: []*cre.NodeSpecWithRole{
					{Roles: []string{cre.BootstrapNode}},
				},
			},
		},
	}

	specs := relaySpecsFromConfig(cfg)
	for _, spec := range specs {
		if spec.Port == 14100 || spec.Port == 5001 {
			t.Fatalf("did not expect OCR relay specs without remote nodesets, got %#v", specs)
		}
	}
}
