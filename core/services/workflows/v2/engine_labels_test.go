package v2

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink/v2/core/platform"
	p2ptypes "github.com/smartcontractkit/chainlink/v2/core/services/p2p/types"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/types"
)

func TestEngine_EventLabelsPreservedAfterLabelRebuild(t *testing.T) {
	t.Parallel()

	cfg := &EngineConfig{
		WorkflowID:    "wf-1",
		WorkflowOwner: "owner-1",
	}
	engine := &baseEngine{cfg: cfg, orgID: "test-org"}
	base := map[string]string{platform.KeyWorkflowID: "wf-1"}
	engine.storeLoggerLabels(base)

	labels := engine.eventLabels()
	require.Equal(t, "test-org", labels[platform.KeyOrganizationID])

	// Simulate localNodeSync rebuilding labels from buildLabels without org.
	engine.storeLoggerLabels(map[string]string{
		platform.KeyWorkflowID: "wf-1",
		platform.KeyDonID:      "11",
		platform.DonVersion:    "1",
	})

	labels = engine.eventLabels()
	require.Equal(t, "test-org", labels[platform.KeyOrganizationID])
	require.Equal(t, "11", labels[platform.KeyDonID])
}

func TestEngine_BuildLabelsIncludesShardAndWorkflowDon(t *testing.T) {
	t.Parallel()

	name, err := types.NewWorkflowName("test-workflow")
	require.NoError(t, err)
	engine := &baseEngine{cfg: &EngineConfig{ShardIndex: 2, WorkflowName: name}}
	peerID := p2ptypes.PeerID{1}
	node := &capabilities.Node{
		PeerID:      &peerID,
		WorkflowDON: capabilities.DON{ID: 7, F: 1, Members: []p2ptypes.PeerID{{1}, {2}, {3}, {4}}},
	}

	labels := map[string]string{}
	kvs := engine.buildLabels(node)
	for i := 0; i < len(kvs); i += 2 {
		labels[kvs[i].(string)] = kvs[i+1].(string)
	}

	require.Equal(t, "7", labels[platform.KeyDonID])
	require.Equal(t, labels[platform.KeyDonID], labels[platform.KeyWorkflowDonID])
	require.Equal(t, "2", labels[platform.KeyShardIndex])
}
