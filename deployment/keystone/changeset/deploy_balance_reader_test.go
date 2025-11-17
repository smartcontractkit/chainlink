package changeset_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	chain_selectors "github.com/smartcontractkit/chain-selectors"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/engine/test/environment"
	"github.com/smartcontractkit/chainlink-deployments-framework/engine/test/runtime"

	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
)

func TestDeployBalanceReader(t *testing.T) {
	t.Parallel()

	registrySel := chain_selectors.TEST_90000001.Selector
	otherSel := chain_selectors.TEST_90000002.Selector
	rt, err := runtime.New(t.Context(), runtime.WithEnvOpts(
		environment.WithEVMSimulated(t, []uint64{registrySel, otherSel}),
	))
	require.NoError(t, err)

	t.Run("should deploy balancereader", func(t *testing.T) {
		qualifier := "my-balance-reader-qualifier"

		err = rt.Exec(
			runtime.ChangesetTask(cldf.CreateLegacyChangeSet(changeset.DeployBalanceReader), changeset.DeployBalanceReaderRequest{
				Qualifier: qualifier,
			}),
		)
		require.NoError(t, err)

		// registry, ocr3, balancereader should be deployed on registry chain
		addrs, err := rt.State().AddressBook.AddressesForChain(registrySel)
		require.NoError(t, err)
		assert.Len(t, addrs, 1)

		dsAddrs, err := rt.State().DataStore.Addresses().Fetch()
		require.NoError(t, err)
		assert.Len(t, dsAddrs, 2) // 2 balance readers, one per chain
		assert.Equal(t, qualifier, dsAddrs[0].Qualifier)
		assert.Equal(t, qualifier, dsAddrs[1].Qualifier)

		// only balancereader on chain 1
		require.NotEqual(t, registrySel, otherSel)
		oaddrs, err := rt.State().AddressBook.AddressesForChain(otherSel)
		require.NoError(t, err)
		assert.Len(t, oaddrs, 1)
	})
}
