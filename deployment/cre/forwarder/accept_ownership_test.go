package forwarder_test

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	chainsel "github.com/smartcontractkit/chain-selectors"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/engine/test/environment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"

	forwarderwrapper "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/forwarder_1_0_0"

	"github.com/smartcontractkit/chainlink/deployment/cre/contracts"
	"github.com/smartcontractkit/chainlink/deployment/cre/forwarder"
)

func TestAcceptOwnershipForwarder(t *testing.T) {
	t.Parallel()

	selector := chainsel.TEST_90000001.Selector

	env, err := environment.New(t.Context(),
		environment.WithEVMSimulated(t, []uint64{selector}),
		environment.WithLogger(logger.Test(t)),
	)
	require.NoError(t, err)

	// Deploy a KeystoneForwarder contract
	deployOut, err := operations.ExecuteSequence(env.OperationsBundle, forwarder.DeploySequence,
		forwarder.DeploySequenceDeps{Env: env},
		forwarder.DeploySequenceInput{
			Targets:   []uint64{selector},
			Qualifier: "test-accept-ownership",
		},
	)
	require.NoError(t, err)

	env.DataStore = deployOut.Output.Datastore

	refs := env.DataStore.Addresses().Filter(
		datastore.AddressRefByChainSelector(selector),
		datastore.AddressRefByType(datastore.ContractType(contracts.KeystoneForwarder)),
	)
	require.Len(t, refs, 1)

	chain := env.BlockChains.EVMChains()[selector]

	// Simulate a pending ownership transfer: the current owner (deployer) calls
	// transferOwnership(deployer), making itself the pending owner. This mirrors the
	// real scenario where a previous owner has called transferOwnership(<deployer_eoa>)
	// before this changeset is run.
	contract, err := forwarderwrapper.NewKeystoneForwarder(common.HexToAddress(refs[0].Address), chain.Client)
	require.NoError(t, err)

	tx, err := contract.TransferOwnership(chain.DeployerKey, chain.DeployerKey.From)
	_, err = cldf.ConfirmIfNoError(chain, tx, err)
	require.NoError(t, err)

	// Apply the changeset — deployer accepts the pending ownership
	_, err = forwarder.AcceptOwnershipForwarder{}.Apply(*env, forwarder.AcceptOwnershipInput{
		ChainSelector: selector,
		Qualifier:     "test-accept-ownership",
	})
	require.NoError(t, err)

	owner, err := contract.Owner(nil)
	require.NoError(t, err)
	require.Equal(t, chain.DeployerKey.From, owner)
}

func TestAcceptOwnershipForwarder_VerifyPreconditions(t *testing.T) {
	t.Parallel()

	selector := chainsel.TEST_90000001.Selector

	env, err := environment.New(t.Context(),
		environment.WithEVMSimulated(t, []uint64{selector}),
		environment.WithLogger(logger.Test(t)),
	)
	require.NoError(t, err)

	t.Run("unknown chain selector", func(t *testing.T) {
		err := forwarder.AcceptOwnershipForwarder{}.VerifyPreconditions(*env, forwarder.AcceptOwnershipInput{
			ChainSelector: 0,
		})
		require.Error(t, err)
		require.Contains(t, err.Error(), "not found in environment")
	})

	t.Run("no forwarder in datastore", func(t *testing.T) {
		err := forwarder.AcceptOwnershipForwarder{}.VerifyPreconditions(*env, forwarder.AcceptOwnershipInput{
			ChainSelector: selector,
		})
		require.Error(t, err)
		require.Contains(t, err.Error(), "no KeystoneForwarder found")
	})
}
