package onchain

import (
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/core/scripts/cre/reconciler/internal/domain"
	keystone_changeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
	crecontracts "github.com/smartcontractkit/chainlink/system-tests/lib/cre/contracts"
)

func addressRef(contractType string, chainSelector uint64) domain.AddressRef {
	return domain.AddressRef{
		ChainSelector: domain.ChainSelector(chainSelector),
		Address:       "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb9226",
		Type:          contractType,
		Version:       crecontracts.V2Version.String(),
	}
}

func TestContractsFullyDeployed_NeitherPresent(t *testing.T) {
	t.Parallel()

	state := &domain.StateFile{}
	require.False(t, contractsFullyDeployed(state))
}

func TestContractsFullyDeployed_OnlyCapRegPresent(t *testing.T) {
	t.Parallel()

	state := &domain.StateFile{}
	state.SetAddress(addressRef(keystone_changeset.CapabilitiesRegistry.String(), 1337))
	require.False(t, contractsFullyDeployed(state))
}

func TestContractsFullyDeployed_OnlyWorkflowRegPresent(t *testing.T) {
	t.Parallel()

	state := &domain.StateFile{}
	state.SetAddress(addressRef(keystone_changeset.WorkflowRegistry.String(), 1337))
	require.False(t, contractsFullyDeployed(state))
}

func TestContractsFullyDeployed_BothPresent(t *testing.T) {
	t.Parallel()

	state := &domain.StateFile{}
	state.SetAddress(addressRef(keystone_changeset.CapabilitiesRegistry.String(), 1337))
	state.SetAddress(addressRef(keystone_changeset.WorkflowRegistry.String(), 1337))
	require.True(t, contractsFullyDeployed(state))
}

func TestDeployContracts_SkipsWhenBothPresent(t *testing.T) {
	t.Parallel()

	const chainSelector = 1337
	state := &domain.StateFile{}
	state.SetAddress(addressRef(keystone_changeset.CapabilitiesRegistry.String(), chainSelector))
	state.SetAddress(addressRef(keystone_changeset.WorkflowRegistry.String(), chainSelector))

	d := NewDeployer(nil, "", zerolog.Nop(), nil)
	env := &cldf.Environment{}

	err := d.deployContracts(env, chainSelector, state)
	require.NoError(t, err)
	require.NotNil(t, env.DataStore)
}
