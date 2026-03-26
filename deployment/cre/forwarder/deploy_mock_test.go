package forwarder_test

import (
	"testing"

	chainsel "github.com/smartcontractkit/chain-selectors"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	"github.com/smartcontractkit/chainlink-deployments-framework/engine/test/environment"
	"github.com/smartcontractkit/chainlink-deployments-framework/engine/test/runtime"

	"github.com/smartcontractkit/chainlink/deployment/cre/contracts"
	"github.com/smartcontractkit/chainlink/deployment/cre/forwarder"
)

func TestDeployMockForwarder(t *testing.T) {
	t.Parallel()

	registrySel := chainsel.TEST_90000001.Selector
	rt, err := runtime.New(t.Context(), runtime.WithEnvOpts(
		environment.WithEVMSimulated(t, []uint64{registrySel}),
	))
	require.NoError(t, err)

	err = rt.Exec(
		runtime.ChangesetTask(forwarder.DeployMockForwarders{},
			forwarder.DeployMockForwardersInput{
				Targets:   []uint64{registrySel},
				Qualifier: "my-test-mock-forwarder",
			},
		),
	)
	require.NoError(t, err)

	addrs := rt.State().DataStore.Addresses().Filter(
		datastore.AddressRefByChainSelector(registrySel),
	)
	require.Len(t, addrs, 1)

	mockAddrs := rt.State().DataStore.Addresses().Filter(
		datastore.AddressRefByType(datastore.ContractType(contracts.MockKeystoneForwarder)),
	)
	require.Len(t, mockAddrs, 1)
	require.Equal(t, "my-test-mock-forwarder", mockAddrs[0].Qualifier)

	labels := mockAddrs[0].Labels.List()
	require.Len(t, labels, 2)
	require.Contains(t, labels[0], forwarder.DeploymentBlockLabel)
	require.Contains(t, labels[1], forwarder.DeploymentHashLabel)
}
