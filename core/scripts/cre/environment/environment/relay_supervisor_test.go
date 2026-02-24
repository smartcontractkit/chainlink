package environment

import (
	"testing"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	envconfig "github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/config"
)

func TestRelaySpecsFromConfig_AddsBootstrapPeeringPortForRemoteToLocalMixedDONs(t *testing.T) {
	cfg := &envconfig.Config{
		NodeSets: []*cre.NodeSet{
			{
				Placement: "local",
				NodeSpecs: []*cre.NodeSpecWithRole{
					{Roles: []string{cre.BootstrapNode}},
				},
			},
			{
				Placement: "remote",
				NodeSpecs: []*cre.NodeSpecWithRole{
					{Roles: []string{cre.WorkerNode}},
				},
			},
		},
	}

	specs := relaySpecsFromConfig(cfg)
	foundBootstrap := false
	for _, spec := range specs {
		if spec.Name == "ocr-bootstrap" && spec.Port == 5001 {
			foundBootstrap = true
			break
		}
	}
	if !foundBootstrap {
		t.Fatalf("expected relay specs to include ocr-bootstrap:5001, got %#v", specs)
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
		if spec.Name == "ocr-bootstrap" && spec.Port == 5001 {
			t.Fatalf("did not expect ocr-bootstrap relay spec without remote nodesets, got %#v", specs)
		}
	}
}
