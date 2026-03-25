package environment

import (
	"testing"

	"github.com/Masterminds/semver/v3"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	keystone_changeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	envconfig "github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/config"
	"github.com/smartcontractkit/chainlink/system-tests/lib/infra"
)

func TestResolverAddressRef(t *testing.T) {
	t.Parallel()

	cfg := &envconfig.Config{}
	require.NoError(t, cfg.SetAddresses([]datastore.AddressRef{{
		Address:       "0x123",
		Type:          datastore.ContractType(keystone_changeset.WorkflowRegistry),
		ChainSelector: 1337,
		Version:       semver.MustParse("2.0.0"),
	}}))

	resolver := &LocalCREStateResolver{cfg: cfg}

	addrRef, err := resolver.AddressRef(keystone_changeset.WorkflowRegistry)
	require.NoError(t, err)
	require.Equal(t, "0x123", addrRef.Address)
	require.Equal(t, uint64(1337), addrRef.ChainSelector)
	require.Equal(t, semver.MustParse("2.0.0").String(), addrRef.Version.String())
}

func TestResolverWorkflowDONMetadata(t *testing.T) {
	t.Parallel()

	donsMetadata, err := cre.NewDonsMetadata([]*cre.DonMetadata{{
		Name:  "workflow-don",
		ID:    7,
		Flags: []string{cre.WorkflowDON},
		NodesMetadata: []*cre.NodeMetadata{{
			Roles: []string{cre.BootstrapNode},
		}},
	}}, infra.Provider{Type: infra.Docker})
	require.NoError(t, err)

	resolver := &LocalCREStateResolver{
		topology: &cre.Topology{DonsMetadata: donsMetadata},
	}

	workflowDON, err := resolver.WorkflowDONMetadata()
	require.NoError(t, err)
	require.Equal(t, "workflow-don", workflowDON.Name)

	workflowDONID, err := resolver.WorkflowDONID()
	require.NoError(t, err)
	require.Equal(t, uint32(7), workflowDONID)

	workflowDONName, err := resolver.WorkflowDONName()
	require.NoError(t, err)
	require.Equal(t, "workflow-don", workflowDONName)
}

func TestResolverGatewayURLFallsBackToProviderHost(t *testing.T) {
	t.Parallel()

	resolver := &LocalCREStateResolver{
		cfg: &envconfig.Config{
			Infra: &infra.Provider{Type: infra.Docker},
		},
		topology: &cre.Topology{
			GatewayConnectors: &cre.GatewayConnectors{
				Configurations: []*cre.DonGatewayConfiguration{{
					GatewayConfiguration: &cre.GatewayConfiguration{
						Incoming: cre.Incoming{
							Protocol:     "http",
							ExternalPort: 5002,
							Path:         "/",
						},
					},
				}},
			},
		},
	}

	gatewayURL, err := resolver.GatewayURL()
	require.NoError(t, err)
	require.Equal(t, "http://localhost:5002/", gatewayURL)
}

func TestResolveContractAddressAndVersion(t *testing.T) {
	t.Parallel()

	makeCmd := func(address string) *cobra.Command {
		cmd := newCobraCommand()
		cmd.Flags().String("workflow-registry-address", "", "")
		require.NoError(t, cmd.Flags().Set("workflow-registry-address", address))
		return cmd
	}

	t.Run("uses state when flag not changed", func(t *testing.T) {
		cfg := &envconfig.Config{}
		require.NoError(t, cfg.SetAddresses([]datastore.AddressRef{{
			Address: "0x456",
			Type:    datastore.ContractType(keystone_changeset.WorkflowRegistry),
			Version: semver.MustParse("1.1.0"),
		}}))

		cmd := newCobraCommand()
		cmd.Flags().String("workflow-registry-address", "", "")

		address, version, err := resolveContractAddressAndVersion(cmd, &LocalCREStateResolver{cfg: cfg}, keystone_changeset.WorkflowRegistry, "", "2.0.0", "workflow-registry-address")
		require.NoError(t, err)
		require.Equal(t, "0x456", address)
		require.Equal(t, "1.1.0", version.String())
	})

	t.Run("uses explicit override when flag changed", func(t *testing.T) {
		cmd := makeCmd("0xabc")

		address, version, err := resolveContractAddressAndVersion(cmd, nil, keystone_changeset.WorkflowRegistry, "0xabc", "2.0.0", "workflow-registry-address")
		require.NoError(t, err)
		require.Equal(t, "0xabc", address)
		require.Equal(t, "2.0.0", version.String())
	})
}

func TestToDockerHostRPC(t *testing.T) {
	t.Parallel()

	require.Equal(t, "http://host.docker.internal:8545", toDockerHostRPC("http://localhost:8545"))
	require.Equal(t, "http://host.docker.internal:8545", toDockerHostRPC("http://127.0.0.1:8545"))
}

func newCobraCommand() *cobra.Command {
	return &cobra.Command{Use: "test"}
}
