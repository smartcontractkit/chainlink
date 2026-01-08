package changeset

import (
	"testing"

	chainsel "github.com/smartcontractkit/chain-selectors"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	"github.com/smartcontractkit/chainlink-deployments-framework/engine/test/environment"
	"github.com/smartcontractkit/chainlink-deployments-framework/engine/test/runtime"

	"github.com/smartcontractkit/chainlink/deployment/cre/contracts"
)

func TestDeployShardConfig(t *testing.T) {
	t.Parallel()

	selector := chainsel.TEST_90000001.Selector
	rt, err := runtime.New(t.Context(), runtime.WithEnvOpts(
		environment.WithEVMSimulated(t, []uint64{selector}),
		environment.WithLogger(logger.Test(t)),
	))
	require.NoError(t, err)

	task := runtime.ChangesetTask(DeployShardConfig{}, DeployShardConfigInput{
		ChainSelector:     selector,
		InitialShardCount: 10,
		Qualifier:         "test-shard-config",
	})

	err = rt.Exec(task)
	require.NoError(t, err)

	// Verify deployment
	addrs := rt.State().DataStore.Addresses().Filter(
		datastore.AddressRefByType(datastore.ContractType(contracts.ShardConfig)),
	)
	require.Len(t, addrs, 1)
	require.Equal(t, "test-shard-config", addrs[0].Qualifier)
}

func TestDeployShardConfig_ValidationErrors(t *testing.T) {
	t.Parallel()

	selector := chainsel.TEST_90000001.Selector
	rt, err := runtime.New(t.Context(), runtime.WithEnvOpts(
		environment.WithEVMSimulated(t, []uint64{selector}),
		environment.WithLogger(logger.Test(t)),
	))
	require.NoError(t, err)

	t.Run("missing chain selector", func(t *testing.T) {
		err := DeployShardConfig{}.VerifyPreconditions(
			rt.Environment(),
			DeployShardConfigInput{InitialShardCount: 10},
		)
		require.Error(t, err)
		require.Contains(t, err.Error(), "chain selector")
	})

	t.Run("zero shard count", func(t *testing.T) {
		err := DeployShardConfig{}.VerifyPreconditions(
			rt.Environment(),
			DeployShardConfigInput{ChainSelector: selector},
		)
		require.Error(t, err)
		require.Contains(t, err.Error(), "shard count")
	})
}
