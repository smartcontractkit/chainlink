package tron_test

import (
	"testing"

	chain_selectors "github.com/smartcontractkit/chain-selectors"
	"github.com/stretchr/testify/require"

	cldf_tron "github.com/smartcontractkit/chainlink-deployments-framework/chain/tron"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	"github.com/smartcontractkit/chainlink-deployments-framework/engine/test/environment"
	"github.com/smartcontractkit/chainlink-deployments-framework/engine/test/runtime"

	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset/tron"
)

func TestDeployForwarder(t *testing.T) {
	t.Parallel()

	registrySel := chain_selectors.TRON_DEVNET.Selector
	rt, err := runtime.New(t.Context(), runtime.WithEnvOpts(
		environment.WithTronContainer(t, []uint64{registrySel}),
	))
	require.NoError(t, err)

	t.Run("should deploy forwarder", func(t *testing.T) {
		deployOptions := cldf_tron.DefaultDeployOptions()
		deployOptions.FeeLimit = 1_000_000_000

		err = rt.Exec(
			runtime.ChangesetTask(tron.DeployForwarder{},
				&tron.DeployForwarderRequest{
					ChainSelectors: []uint64{registrySel},
					Qualifier:      "my-test-forwarder",
					DeployOptions:  deployOptions,
				},
			),
		)
		require.NoError(t, err)

		addrs := rt.State().DataStore.Addresses().Filter(datastore.AddressRefByQualifier("my-test-forwarder"))
		require.Len(t, addrs, 1)
	})
}
