package environment

import (
	"strings"
	"testing"

	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/clnode"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/simple_node_set"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
)

func TestBuildRemoteNodeSetInputRequiresImageOrBuildFields(t *testing.T) {
	nodeSet := &cre.NodeSet{
		Input: &simple_node_set.Input{
			Name: "remote-don",
		},
		NodeSpecs: []*cre.NodeSpecWithRole{
			{
				Input: &clnode.Input{
					Node: &clnode.NodeInput{
						Image: "",
					},
				},
			},
		},
	}

	_, err := buildRemoteNodeSetInput(nodeSet)
	if err == nil {
		t.Fatal("expected missing image/build validation error")
	}
	if !strings.Contains(err.Error(), "must set node.image or docker build fields") {
		t.Fatalf("expected image validation error, got: %v", err)
	}
}

func TestBuildRemoteNodeSetInputRejectsImageAndBuildFieldsTogether(t *testing.T) {
	nodeSet := &cre.NodeSet{
		Input: &simple_node_set.Input{
			Name: "remote-don",
		},
		NodeSpecs: []*cre.NodeSpecWithRole{
			{
				Input: &clnode.Input{
					Node: &clnode.NodeInput{
						Image:          "repo/chainlink:tag",
						DockerContext:  "../../../..",
						DockerFilePath: "core/chainlink.Dockerfile",
					},
				},
			},
		},
	}

	_, err := buildRemoteNodeSetInput(nodeSet)
	if err == nil {
		t.Fatal("expected image+build conflict validation error")
	}
	if !strings.Contains(err.Error(), "either node.image or docker build fields") {
		t.Fatalf("expected image/build conflict error, got: %v", err)
	}
}
