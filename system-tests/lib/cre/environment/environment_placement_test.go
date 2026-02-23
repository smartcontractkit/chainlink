package environment

import (
	"testing"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
)

func TestSummarizeNodeSetPlacement_AllowsMixedPlacements(t *testing.T) {
	nodeSets := []*cre.NodeSet{
		{Placement: "local"},
		{Placement: "remote"},
	}

	summary, err := summarizeNodeSetPlacement(nodeSets)
	if err != nil {
		t.Fatalf("summarizeNodeSetPlacement returned error: %v", err)
	}
	if !summary.HasLocalTargets || !summary.HasRemoteTargets {
		t.Fatalf("expected both local and remote placements, got %+v", summary)
	}
}
