package v2

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/platform"
)

func TestEngine_EventLabelsPreservedAfterLabelRebuild(t *testing.T) {
	t.Parallel()

	cfg := &EngineConfig{
		WorkflowID:    "wf-1",
		WorkflowOwner: "owner-1",
	}
	engine := &Engine{cfg: cfg, orgID: "test-org"}
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
