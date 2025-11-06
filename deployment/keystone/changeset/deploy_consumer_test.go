package changeset_test

import (
	"testing"

	chain_selectors "github.com/smartcontractkit/chain-selectors"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/engine/test/environment"
	"github.com/smartcontractkit/chainlink-deployments-framework/engine/test/runtime"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"

	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
)

func TestDeployFeedsConsumer(t *testing.T) {
	t.Parallel()

	registrySel := chain_selectors.TEST_90000001.Selector
	otherSel := chain_selectors.TEST_90000002.Selector
	rt, err := runtime.New(t.Context(), runtime.WithEnvOpts(
		environment.WithEVMSimulated(t, []uint64{registrySel, otherSel}),
	))
	require.NoError(t, err)

	err = rt.Exec(
		runtime.ChangesetTask(cldf.CreateLegacyChangeSet(changeset.DeployFeedsConsumerV2), &changeset.DeployRequestV2{
			ChainSel:  registrySel,
			Qualifier: "my-test-feeds-consumer",
		}),
	)
	require.NoError(t, err)

	// feeds consumer should be deployed on chain 0
	addrs, err := rt.State().AddressBook.AddressesForChain(registrySel)
	require.NoError(t, err)
	require.Len(t, addrs, 1)
	require.Len(t, rt.State().DataStore.Addresses().Filter(datastore.AddressRefByQualifier("my-test-feeds-consumer")), 1, "expected to find 'my-test-feeds-consumer' qualifier")

	// no feeds consumer registry on chain 1
	require.NotEqual(t, registrySel, otherSel)
	oaddrs, _ := rt.State().AddressBook.AddressesForChain(otherSel)
	require.Empty(t, oaddrs)
}
