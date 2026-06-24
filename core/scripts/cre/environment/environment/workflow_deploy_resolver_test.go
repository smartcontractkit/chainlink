// Tests for workflow deploy target resolution (donID, donFamily, gatewayURL).
package environment

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	envconfig "github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/config"
	"github.com/smartcontractkit/chainlink/system-tests/lib/infra"
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
		family, err := finalizeWorkflowDonFamily("feeds-zone-a")
		require.NoError(t, err)
		require.Equal(t, "feeds-zone-a", family)
	})

	t.Run("requires family", func(t *testing.T) {
		t.Parallel()
		_, err := finalizeWorkflowDonFamily("")
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

func TestResolveWorkflowDeployTargets_byDonFamilyOnly(t *testing.T) {
	t.Parallel()

	topology := deployTestTwoFamilyTopology(t)
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

func TestResolveWorkflowDeployTargets_donFamilyIsolation(t *testing.T) {
	t.Parallel()

	topology := deployTestTwoFamilyTopology(t)
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

func TestResolveWorkflowDeployTargets_defaultsGatewayURLFromFamily(t *testing.T) {
	t.Parallel()

	topology := deployTestTwoFamilyTopologyWithIncoming(t)
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

func TestResolveWorkflowDeployTargets_nilResolverRequiresDonFamily(t *testing.T) {
	t.Parallel()

	cmd := newTestDeployCmd(t)

	_, err := resolveWorkflowDeployTargets(
		cmd,
		nil,
		workflowDONSelector{},
		0,
		"",
		"",
	)
	require.Error(t, err)
}

func TestResolveWorkflowDeployTargets_infersDonFamilyFromWorkflowDonNameOnly(t *testing.T) {
	t.Parallel()

	topology := deployTestTwoFamilyTopologyWithIncoming(t)
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

func TestResolveWorkflowDeployTargets_shardedByDonFamilyAndShardIndex(t *testing.T) {
	t.Parallel()

	topology := deployTestShardedTopology(t)
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

func TestResolveWorkflowDeployTargets_preservesExplicitDonID(t *testing.T) {
	t.Parallel()

	topology := deployTestTwoFamilyTopologyWithIncoming(t)
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

func deployTestTwoFamilyTopology(t *testing.T) *cre.Topology {
	t.Helper()

	return mustDeployTestTopology(t,
		[]*cre.DonGatewayConfiguration{
			{GatewayConfiguration: &cre.GatewayConfiguration{AuthGatewayID: "gateway-node-0"}},
			{GatewayConfiguration: &cre.GatewayConfiguration{AuthGatewayID: "gateway-node-1"}},
		},
		&cre.DonMetadata{Name: "feeds-zone-a", ID: 1, DonFamily: "feeds-zone-a", Flags: []string{cre.WorkflowDON, cre.HTTPActionCapability}},
		&cre.DonMetadata{Name: "feeds-zone-b", ID: 2, DonFamily: "feeds-zone-b", Flags: []string{cre.WorkflowDON, cre.HTTPActionCapability}},
		&cre.DonMetadata{Name: "gateway-zone-a", DonFamily: "feeds-zone-a", NodesMetadata: []*cre.NodeMetadata{{Roles: []string{cre.GatewayNode}}}},
		&cre.DonMetadata{Name: "gateway-zone-b", DonFamily: "feeds-zone-b", NodesMetadata: []*cre.NodeMetadata{{Roles: []string{cre.GatewayNode}}}},
	)
}

func deployTestTwoFamilyTopologyWithIncoming(t *testing.T) *cre.Topology {
	t.Helper()

	return mustDeployTestTopology(t,
		[]*cre.DonGatewayConfiguration{
			deployTestGatewayConnector("gateway-node-0", "gateway-zone-a.local", 5002),
			deployTestGatewayConnector("gateway-node-1", "gateway-zone-b.local", 5004),
		},
		&cre.DonMetadata{Name: "feeds-zone-a", ID: 1, DonFamily: "feeds-zone-a", Flags: []string{cre.WorkflowDON, cre.HTTPActionCapability}},
		&cre.DonMetadata{Name: "feeds-zone-b", ID: 2, DonFamily: "feeds-zone-b", Flags: []string{cre.WorkflowDON, cre.HTTPActionCapability}},
		&cre.DonMetadata{Name: "gateway-zone-a", DonFamily: "feeds-zone-a", NodesMetadata: []*cre.NodeMetadata{{Roles: []string{cre.GatewayNode}}}},
		&cre.DonMetadata{Name: "gateway-zone-b", DonFamily: "feeds-zone-b", NodesMetadata: []*cre.NodeMetadata{{Roles: []string{cre.GatewayNode}}}},
	)
}

func deployTestShardedTopology(t *testing.T) *cre.Topology {
	t.Helper()

	return mustDeployTestTopology(t,
		[]*cre.DonGatewayConfiguration{
			deployTestGatewayConnector("gateway-node-0", "bootstrap-gateway.local", 5002),
		},
		&cre.DonMetadata{Name: "shard0", ID: 1, DonFamily: envconfig.DefaultDONFamily, ShardIndex: 0, Flags: []string{cre.WorkflowDON, cre.ShardDON, cre.HTTPActionCapability}},
		&cre.DonMetadata{Name: "shard1", ID: 2, DonFamily: envconfig.DefaultDONFamily, ShardIndex: 1, Flags: []string{cre.WorkflowDON, cre.ShardDON, cre.HTTPActionCapability}},
		&cre.DonMetadata{Name: "bootstrap-gateway", DonFamily: envconfig.DefaultDONFamily, NodesMetadata: []*cre.NodeMetadata{{Roles: []string{cre.GatewayNode}}}},
	)
}

func mustDeployTestTopology(t *testing.T, connectors []*cre.DonGatewayConfiguration, dons ...*cre.DonMetadata) *cre.Topology {
	t.Helper()

	allDONs := append([]*cre.DonMetadata{deployTestBootstrapDON()}, dons...)
	dm, err := cre.NewDonsMetadata(allDONs, infra.Provider{Type: infra.Docker})
	require.NoError(t, err)

	topology := &cre.Topology{
		DonsMetadata: dm,
		GatewayConnectors: &cre.GatewayConnectors{
			Configurations: connectors,
		},
	}
	require.NoError(t, topology.EnsureGatewayDonFamilyPairing())
	return topology
}

func deployTestBootstrapDON() *cre.DonMetadata {
	return &cre.DonMetadata{
		Name:          "bootstrap",
		DonFamily:     envconfig.DefaultDONFamily,
		NodesMetadata: []*cre.NodeMetadata{{Roles: []string{cre.BootstrapNode}}},
	}
}

func deployTestGatewayConnector(authID, host string, externalPort int) *cre.DonGatewayConfiguration {
	return &cre.DonGatewayConfiguration{
		GatewayConfiguration: &cre.GatewayConfiguration{
			AuthGatewayID: authID,
			Incoming: cre.Incoming{
				Protocol:     "http",
				Host:         host,
				Path:         "/",
				ExternalPort: externalPort,
			},
		},
	}
}
