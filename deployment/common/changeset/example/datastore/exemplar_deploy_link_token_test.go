package example

import (
	"testing"

	chainselectors "github.com/smartcontractkit/chain-selectors"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-deployments-framework/engine/test/environment"
)

func Test_ExemplarDeployLinkToken(t *testing.T) {
	t.Parallel()

	selector := chainselectors.TEST_90000001.Selector
	env, err := environment.New(t.Context(), environment.WithEVMSimulated(t, []uint64{selector}))
	require.NoError(t, err)

	result, err := ExemplarDeployLinkToken{}.Apply(*env, selector)
	require.NoError(t, err)

	// Check that one address ref was created
	addresRefs, err := result.DataStore.Addresses().Fetch()
	require.NoError(t, err)
	require.Len(t, addresRefs, 1)

	// Check that one contract metadata ref was created
	contractMetadata, err := result.DataStore.ContractMetadata().Fetch()
	require.NoError(t, err)
	require.Len(t, contractMetadata, 1)

	// Check that env metadata was set correctly
	_, err = result.DataStore.EnvMetadata().Get()
	require.NoError(t, err)
}
