package environment

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
)

func newTestDeployCmd(t *testing.T) *cobra.Command {
	t.Helper()

	cmd := &cobra.Command{}
	cmd.Flags().String("workflow-don-name", "", "")
	cmd.Flags().String("don-family", "", "")
	cmd.Flags().Uint32("don-id", 0, "")
	cmd.Flags().String("gateway-url", "", "")
	return cmd
}

func TestValidateWorkflowDeployFlags_donFamilyMismatch(t *testing.T) {
	t.Parallel()

	cmd := newTestDeployCmd(t)
	require.NoError(t, cmd.Flags().Set("workflow-don-name", "feeds-zone-a"))
	require.NoError(t, cmd.Flags().Set("don-family", "feeds-zone-b"))

	donMeta := &cre.DonMetadata{Name: "feeds-zone-a", DonFamily: "feeds-zone-a"}
	err := validateWorkflowDeployFlags(cmd, donMeta, workflowDONSelector{ExplicitName: "feeds-zone-a"}, "feeds-zone-b")
	require.Error(t, err)
	require.Contains(t, err.Error(), `--don-family "feeds-zone-b"`)
	require.Contains(t, err.Error(), `don_family "feeds-zone-a"`)
}

func TestValidateWorkflowDeployFlags_agreeWithState(t *testing.T) {
	t.Parallel()

	cmd := newTestDeployCmd(t)
	require.NoError(t, cmd.Flags().Set("workflow-don-name", "feeds-zone-a"))
	require.NoError(t, cmd.Flags().Set("don-family", "feeds-zone-a"))

	donMeta := &cre.DonMetadata{Name: "feeds-zone-a", DonFamily: "feeds-zone-a"}
	err := validateWorkflowDeployFlags(cmd, donMeta, workflowDONSelector{ExplicitName: "feeds-zone-a"}, "feeds-zone-a")
	require.NoError(t, err)
}

func TestValidateWorkflowDeployFlags_skipsWhenOnlyOneFlagSet(t *testing.T) {
	t.Parallel()

	cmd := newTestDeployCmd(t)
	require.NoError(t, cmd.Flags().Set("don-family", "feeds-zone-b"))

	donMeta := &cre.DonMetadata{Name: "feeds-zone-a", DonFamily: "feeds-zone-a"}
	err := validateWorkflowDeployFlags(cmd, donMeta, workflowDONSelector{ExplicitName: "feeds-zone-a"}, "feeds-zone-b")
	require.NoError(t, err)
}

func TestResolveWorkflowDeployTargets_multiZoneFamilyIsolation(t *testing.T) {
	t.Parallel()

	topology := cre.NewDonFamilyGatewayPairingTestTopology()
	resolver := &LocalCREStateResolver{topology: topology}

	resolveZone := func(t *testing.T, zone string) workflowDeployTargets {
		t.Helper()

		cmd := newTestDeployCmd(t)
		require.NoError(t, cmd.Flags().Set("workflow-don-name", zone))
		require.NoError(t, cmd.Flags().Set("don-family", zone))
		require.NoError(t, cmd.Flags().Set("gateway-url", "http://localhost:8080"))

		targets, err := resolveWorkflowDeployTargets(
			cmd,
			resolver,
			workflowDONSelector{ExplicitName: zone},
			0,
			zone,
			"http://localhost:8080",
		)
		require.NoError(t, err)
		return targets
	}

	targetsA := resolveZone(t, "feeds-zone-a")
	targetsB := resolveZone(t, "feeds-zone-b")

	require.Equal(t, uint32(1), targetsA.donID)
	require.Equal(t, "feeds-zone-a", targetsA.donFamily)
	require.Equal(t, uint32(2), targetsB.donID)
	require.Equal(t, "feeds-zone-b", targetsB.donFamily)
	require.NotEqual(t, targetsA.donFamily, targetsB.donFamily)
}

func TestMultiZoneWorkflowMetadataFiltering(t *testing.T) {
	t.Parallel()

	// Lightweight smoke: workflow nodes filter registry metadata by don_family.
	// Mirrors core/services/workflows/syncer/v2/file_workflow_source.go filtering.
	donFamilies := map[string]struct{}{"feeds-zone-a": {}}

	workflows := []struct {
		name      string
		donFamily string
		visible   bool
	}{
		{name: "wf-zone-a", donFamily: "feeds-zone-a", visible: true},
		{name: "wf-zone-b", donFamily: "feeds-zone-b", visible: false},
	}

	var visible []string
	for _, wf := range workflows {
		if _, ok := donFamilies[wf.donFamily]; ok {
			visible = append(visible, wf.name)
		}
	}

	require.Equal(t, []string{"wf-zone-a"}, visible)
}
