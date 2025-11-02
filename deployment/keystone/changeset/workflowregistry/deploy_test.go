package workflowregistry

import (
	"testing"

	chain_selectors "github.com/smartcontractkit/chain-selectors"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/engine/test/environment"
	"github.com/smartcontractkit/chainlink-deployments-framework/engine/test/runtime"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Deploy(t *testing.T) {
	t.Parallel()

	registrySel := chain_selectors.TEST_90000001.Selector
	otherSel := chain_selectors.TEST_90000002.Selector
	rt, err := runtime.New(t.Context(), runtime.WithEnvOpts(
		environment.WithEVMSimulated(t, []uint64{registrySel, otherSel}),
	))
	require.NoError(t, err)

	err = rt.Exec(
		runtime.ChangesetTask(cldf.CreateLegacyChangeSet(Deploy), registrySel),
	)
	require.NoError(t, err)

	addrs, err := rt.State().AddressBook.AddressesForChain(registrySel)
	require.NoError(t, err)
	require.Len(t, addrs, 1)

	// assert nothing on chain 1
	require.NotEqual(t, registrySel, otherSel)
	oaddrs, _ := rt.State().AddressBook.AddressesForChain(otherSel)
	assert.Empty(t, oaddrs)
}
