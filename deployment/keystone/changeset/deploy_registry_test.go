package changeset_test

import (
	"testing"

	chain_selectors "github.com/smartcontractkit/chain-selectors"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/engine/test/environment"
	"github.com/smartcontractkit/chainlink-deployments-framework/engine/test/runtime"

	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
)

func TestDeployCapabilityRegistry(t *testing.T) {
	t.Parallel()

	registrySel := chain_selectors.TEST_90000001.Selector
	otherSel := chain_selectors.TEST_90000002.Selector
	rt, err := runtime.New(t.Context(), runtime.WithEnvOpts(
		environment.WithEVMSimulated(t, []uint64{registrySel, otherSel}),
	))
	require.NoError(t, err)

	err = rt.Exec(
		runtime.ChangesetTask(cldf.CreateLegacyChangeSet(changeset.DeployCapabilityRegistryV2), &changeset.DeployRequestV2{
			ChainSel:  registrySel,
			Qualifier: "my-test-capabilities-registry",
		}),
	)
	require.NoError(t, err)

	// capabilities registry should be deployed on chain 0
	addrs, err := rt.State().AddressBook.AddressesForChain(registrySel)
	require.NoError(t, err)
	require.Len(t, addrs, 1)
	require.Len(t, rt.State().DataStore.Addresses().Filter(datastore.AddressRefByQualifier("my-test-capabilities-registry")), 1, "expected to find 'my-test-capabilities-registry' qualifier")

	// no capabilities registry on chain 1
	require.NotEqual(t, registrySel, otherSel)
	oaddrs, _ := rt.State().AddressBook.AddressesForChain(otherSel)
	require.Empty(t, oaddrs)
}
