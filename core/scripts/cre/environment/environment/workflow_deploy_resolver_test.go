// Tests for workflow deploy target resolution (donID, donFamily, gatewayURL).
//
// Covers resolveWorkflowDeployTargets, finalizeWorkflowDonFamily, and deploy-time
// flag cross-checks. Exercises multi-zone and sharded topologies via in-memory
// fixtures from system-tests/lib/cre/test_topology_helpers.go.
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
	cmd.Flags().Uint("shard-index", 0, "")
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

	t.Run("requires family when local state", func(t *testing.T) {
		t.Parallel()
		_, err := finalizeWorkflowDonFamily("", true)
		require.Error(t, err)
	})
}

// validateWorkflowDeployFlags only cross-checks when both --workflow-don-name and --don-family are set.
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

func TestValidateWorkflowDeployFlags_nameMismatch(t *testing.T) {
	t.Parallel()

	cmd := newTestDeployCmd(t)
	require.NoError(t, cmd.Flags().Set("workflow-don-name", "feeds-zone-b"))
	require.NoError(t, cmd.Flags().Set("don-family", "feeds-zone-a"))

	donMeta := &cre.DonMetadata{Name: "feeds-zone-a", DonFamily: "feeds-zone-a"}
	err := validateWorkflowDeployFlags(cmd, donMeta, workflowDONSelector{ExplicitName: "feeds-zone-b"}, "feeds-zone-a")
	require.Error(t, err)
	require.Contains(t, err.Error(), `--workflow-don-name "feeds-zone-b"`)
}

// DF-style deploy: --don-family alone resolves the workflow DON and fills donID from state.
func TestResolveWorkflowDeployTargets_byDonFamilyOnly(t *testing.T) {
	t.Parallel()

	topology := cre.NewDonFamilyGatewayPairingTestTopology()
	resolver := &LocalCREStateResolver{topology: topology}

	cmd := newTestDeployCmd(t)
	require.NoError(t, cmd.Flags().Set("don-family", "feeds-zone-a"))
	require.NoError(t, cmd.Flags().Set("gateway-url", "http://localhost:8080"))

	targets, err := resolveWorkflowDeployTargets(
		cmd,
		resolver,
		workflowDONSelector{DonFamily: "feeds-zone-a"},
		0,
		"feeds-zone-a",
		"http://localhost:8080",
	)
	require.NoError(t, err)
	require.Equal(t, uint32(1), targets.donID)
	require.Equal(t, "feeds-zone-a", targets.donFamily)
}

// Each zone gets its own cap-reg donID and registry family — no cross-zone bleed.
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

// When --gateway-url is omitted, deploy defaults to the gateway paired with the target don_family.
func TestResolveWorkflowDeployTargets_defaultsGatewayURLFromFamily(t *testing.T) {
	t.Parallel()

	topology := cre.NewDonFamilyGatewayPairingTestTopologyWithIncoming()
	resolver := &LocalCREStateResolver{topology: topology}

	resolveURL := func(t *testing.T, family string, wantURL string) {
		t.Helper()

		cmd := newTestDeployCmd(t)
		require.NoError(t, cmd.Flags().Set("don-family", family))

		targets, err := resolveWorkflowDeployTargets(
			cmd,
			resolver,
			workflowDONSelector{DonFamily: family},
			0,
			family,
			"",
		)
		require.NoError(t, err)
		require.Equal(t, wantURL, targets.gatewayURL)
	}

	resolveURL(t, "feeds-zone-a", "http://gateway-zone-a.local:5002/")
	resolveURL(t, "feeds-zone-b", "http://gateway-zone-b.local:5004/")
}

// Deploy without local CRE state on disk falls back to DefaultDONFamily.
func TestResolveWorkflowDeployTargets_nilResolverUsesDefaultDONFamily(t *testing.T) {
	t.Parallel()

	cmd := newTestDeployCmd(t)

	targets, err := resolveWorkflowDeployTargets(
		cmd,
		nil,
		workflowDONSelector{},
		0,
		"",
		"",
	)
	require.NoError(t, err)
	require.Equal(t, envconfig.DefaultDONFamily, targets.donFamily)
}

// --workflow-don-name alone is enough to infer donFamily and donID from saved topology.
func TestResolveWorkflowDeployTargets_infersDonFamilyFromWorkflowDonNameOnly(t *testing.T) {
	t.Parallel()

	topology := cre.NewDonFamilyGatewayPairingTestTopologyWithIncoming()
	resolver := &LocalCREStateResolver{topology: topology}

	cmd := newTestDeployCmd(t)
	require.NoError(t, cmd.Flags().Set("workflow-don-name", "feeds-zone-b"))

	targets, err := resolveWorkflowDeployTargets(
		cmd,
		resolver,
		workflowDONSelector{ExplicitName: "feeds-zone-b"},
		0,
		"",
		"",
	)
	require.NoError(t, err)
	require.Equal(t, uint32(2), targets.donID)
	require.Equal(t, "feeds-zone-b", targets.donFamily)
}

// Sharded topologies: don_family + shard index resolve shard1's cap-reg donID for UpsertWorkflow.
func TestResolveWorkflowDeployTargets_shardedByDonFamilyAndShardIndex(t *testing.T) {
	t.Parallel()

	topology := cre.NewShardedWorkflowTestTopology()
	resolver := &LocalCREStateResolver{topology: topology}

	cmd := newTestDeployCmd(t)
	require.NoError(t, cmd.Flags().Set("don-family", envconfig.DefaultDONFamily))

	shard1 := uint(1)
	targets, err := resolveWorkflowDeployTargets(
		cmd,
		resolver,
		workflowDONSelector{DonFamily: envconfig.DefaultDONFamily, ShardIndex: &shard1},
		0,
		envconfig.DefaultDONFamily,
		"",
	)
	require.NoError(t, err)
	require.Equal(t, uint32(2), targets.donID)
	require.Equal(t, envconfig.DefaultDONFamily, targets.donFamily)
}

// Explicit --don-id is not overwritten by topology metadata when the flag was set.
func TestResolveWorkflowDeployTargets_preservesExplicitDonID(t *testing.T) {
	t.Parallel()

	topology := cre.NewDonFamilyGatewayPairingTestTopologyWithIncoming()
	resolver := &LocalCREStateResolver{topology: topology}

	cmd := newTestDeployCmd(t)
	require.NoError(t, cmd.Flags().Set("workflow-don-name", "feeds-zone-a"))
	require.NoError(t, cmd.Flags().Set("don-id", "99"))

	targets, err := resolveWorkflowDeployTargets(
		cmd,
		resolver,
		workflowDONSelector{ExplicitName: "feeds-zone-a"},
		99,
		"",
		"",
	)
	require.NoError(t, err)
	require.Equal(t, uint32(99), targets.donID)
}
