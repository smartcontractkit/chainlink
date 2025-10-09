package contracts_test

import (
	"fmt"
	"testing"

	chain_selectors "github.com/smartcontractkit/chain-selectors"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-deployments-framework/engine/test/environment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations/optest"
	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset/operations/contracts"
)

func Test_DeployRegistryOp(t *testing.T) {
	t.Parallel()

	registrySel := chain_selectors.TEST_90000001.Selector
	env, err := environment.New(t.Context(),
		environment.WithEVMSimulated(t, []uint64{registrySel}),
	)
	require.NoError(t, err)

	b := optest.NewBundle(t)
	deps := contracts.DeployCapabilityRegistryOpDeps{
		Env: env,
	}
	input := contracts.DeployCapabilityRegistryInput{
		ChainSelector: registrySel,
	}

	got, err := operations.ExecuteOperation(b, contracts.DeployCapabilityRegistryOp, deps, input)
	require.NoError(t, err)
	addrRefs, err := got.Output.Addresses.Fetch()
	require.NoError(t, err)
	require.Len(t, addrRefs, 1)

	fmt.Println(env.DataStore.Addresses())
}
