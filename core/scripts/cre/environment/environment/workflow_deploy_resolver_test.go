package environment

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	envconfig "github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/config"
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

func TestFinalizeWorkflowDonFamily(t *testing.T) {
	t.Parallel()

	t.Run("explicit family", func(t *testing.T) {
		t.Parallel()
		family, err := finalizeWorkflowDonFamily("feeds-zone-a", true)
		require.NoError(t, err)
		require.Equal(t, "feeds-zone-a", family)
	})

	t.Run("legacy default", func(t *testing.T) {
		t.Parallel()
		family, err := finalizeWorkflowDonFamily("", false)
		require.NoError(t, err)
		require.Equal(t, envconfig.DefaultDONFamily, family)
	})

	t.Run("pairing requires family", func(t *testing.T) {
		t.Parallel()
		_, err := finalizeWorkflowDonFamily("", true)
		require.Error(t, err)
	})
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

func TestResolveWorkflowDeployTargets_donFamilyIsolation(t *testing.T) {
	t.Parallel()

	topology := cre.NewDonFamilyGatewayPairingTestTopology()
	resolver := &LocalCREStateResolver{topology: topology}

	resolveForDON := func(t *testing.T, donName string) workflowDeployTargets {
		t.Helper()

		cmd := newTestDeployCmd(t)
		require.NoError(t, cmd.Flags().Set("workflow-don-name", donName))
		require.NoError(t, cmd.Flags().Set("don-family", donName))
		require.NoError(t, cmd.Flags().Set("gateway-url", "http://localhost:8080"))

		targets, err := resolveWorkflowDeployTargets(
			cmd,
			resolver,
			workflowDONSelector{ExplicitName: donName},
			0,
			donName,
			"http://localhost:8080",
		)
		require.NoError(t, err)
		return targets
	}

	targetsA := resolveForDON(t, "feeds-zone-a")
	targetsB := resolveForDON(t, "feeds-zone-b")

	require.Equal(t, uint32(1), targetsA.donID)
	require.Equal(t, "feeds-zone-a", targetsA.donFamily)
	require.Equal(t, uint32(2), targetsB.donID)
	require.Equal(t, "feeds-zone-b", targetsB.donFamily)
	require.NotEqual(t, targetsA.donFamily, targetsB.donFamily)
}
