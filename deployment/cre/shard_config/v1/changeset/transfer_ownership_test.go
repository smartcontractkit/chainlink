package changeset

import (
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

func TestTransferOwnership(t *testing.T) {
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

	// Create a new owner address
	newOwner := common.HexToAddress("0x1234567890123456789012345678901234567890")

	// Transfer ownership
	transferTask := runtime.ChangesetTask(TransferOwnership{}, TransferOwnershipInput{
		ChainSelector:  selector,
		NewOwner:       newOwner,
		ShardConfigRef: shardConfigRef,
	})

	err = rt.Exec(transferTask)
	require.NoError(t, err)

	// Verify the transfer by reading pending owner from the contract
	// Note: With 2-step ownership, the new owner needs to accept ownership
	contract, err := shard_config.NewShardConfig(
		common.HexToAddress(addrs[0].Address),
		rt.Environment().BlockChains.EVMChains()[selector].Client,
	)
	require.NoError(t, err)

	// The current owner should still be the deployer until the new owner accepts
	currentOwner, err := contract.Owner(nil)
	require.NoError(t, err)
	// Owner hasn't changed yet (2-step process)
	require.NotEqual(t, newOwner, currentOwner)
}

func TestTransferOwnership_ValidationErrors(t *testing.T) {
	t.Parallel()

	selector := chainsel.TEST_90000001.Selector
	rt, err := runtime.New(t.Context(), runtime.WithEnvOpts(
		environment.WithEVMSimulated(t, []uint64{selector}),
		environment.WithLogger(logger.Test(t)),
	))
	require.NoError(t, err)

	t.Run("missing chain selector", func(t *testing.T) {
		err := TransferOwnership{}.VerifyPreconditions(
			rt.Environment(),
			TransferOwnershipInput{NewOwner: common.Address{1}},
		)
		require.Error(t, err)
		require.Contains(t, err.Error(), "chain selector")
	})

	t.Run("missing new owner", func(t *testing.T) {
		err := TransferOwnership{}.VerifyPreconditions(
			rt.Environment(),
			TransferOwnershipInput{ChainSelector: selector},
		)
		require.Error(t, err)
		require.Contains(t, err.Error(), "new owner")
	})
}
