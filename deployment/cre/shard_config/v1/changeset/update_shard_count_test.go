package changeset

import (
	"math/big"
	"testing"

	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/common"
	chainsel "github.com/smartcontractkit/chain-selectors"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	"github.com/smartcontractkit/chainlink-deployments-framework/engine/test/environment"
	"github.com/smartcontractkit/chainlink-deployments-framework/engine/test/runtime"

	"github.com/smartcontractkit/chainlink/deployment/cre/contracts"

	shard_config "github.com/smartcontractkit/chainlink-evm/contracts/cre/gobindings/shardconfig/generated/v1_0_0/shard_config"
)

func TestUpdateShardCount(t *testing.T) {
	t.Parallel()

	selector := chainsel.TEST_90000001.Selector
	rt, err := runtime.New(t.Context(), runtime.WithEnvOpts(
		environment.WithEVMSimulated(t, []uint64{selector}),
		environment.WithLogger(logger.Test(t)),
	))
	require.NoError(t, err)

	// First deploy the ShardConfig contract
	deployTask := runtime.ChangesetTask(DeployShardConfig{}, DeployShardConfigInput{
		ChainSelector:     selector,
		InitialShardCount: 10,
		Qualifier:         "test-shard-config",
	})

	err = rt.Exec(deployTask)
	require.NoError(t, err)

	// Get the deployed contract address
	addrs := rt.State().DataStore.Addresses().Filter(
		datastore.AddressRefByType(datastore.ContractType(contracts.ShardConfig)),
	)
	require.Len(t, addrs, 1)

	// Create the address reference key
	shardConfigRef := datastore.NewAddressRefKey(
		selector,
		datastore.ContractType(contracts.ShardConfig),
		semver.MustParse("1.0.0"),
		"test-shard-config",
	)

	// Update the shard count
	updateTask := runtime.ChangesetTask(UpdateShardCount{}, UpdateShardCountInput{
		ChainSelector:  selector,
		NewShardCount:  20,
		ShardConfigRef: shardConfigRef,
	})

	err = rt.Exec(updateTask)
	require.NoError(t, err)

	// Verify the update by reading from the contract
	contract, err := shard_config.NewShardConfig(
		common.HexToAddress(addrs[0].Address),
		rt.Environment().BlockChains.EVMChains()[selector].Client,
	)
	require.NoError(t, err)

	count, err := contract.GetDesiredShardCount(nil)
	require.NoError(t, err)
	require.Equal(t, big.NewInt(20), count)
}

func TestUpdateShardCount_ValidationErrors(t *testing.T) {
	t.Parallel()

	selector := chainsel.TEST_90000001.Selector
	rt, err := runtime.New(t.Context(), runtime.WithEnvOpts(
		environment.WithEVMSimulated(t, []uint64{selector}),
		environment.WithLogger(logger.Test(t)),
	))
	require.NoError(t, err)

	t.Run("missing chain selector", func(t *testing.T) {
		err := UpdateShardCount{}.VerifyPreconditions(
			rt.Environment(),
			UpdateShardCountInput{NewShardCount: 10},
		)
		require.Error(t, err)
		require.Contains(t, err.Error(), "chain selector")
	})

	t.Run("zero shard count", func(t *testing.T) {
		err := UpdateShardCount{}.VerifyPreconditions(
			rt.Environment(),
			UpdateShardCountInput{ChainSelector: selector},
		)
		require.Error(t, err)
		require.Contains(t, err.Error(), "shard count")
	})
}
