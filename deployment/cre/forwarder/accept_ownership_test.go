package forwarder_test

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	chainsel "github.com/smartcontractkit/chain-selectors"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	cldfchain "github.com/smartcontractkit/chainlink-deployments-framework/chain"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/engine/test/environment"
	"github.com/smartcontractkit/chainlink-deployments-framework/engine/test/onchain"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"

	forwarderwrapper "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/forwarder_1_0_0"

	"github.com/smartcontractkit/chainlink/deployment/cre/contracts"
	"github.com/smartcontractkit/chainlink/deployment/cre/forwarder"
)

func TestAcceptOwnershipForwarder(t *testing.T) {
	t.Parallel()

	selector := chainsel.TEST_90000001.Selector

	// One additional funded account acts as the "new owner" that will accept ownership,
	// mirroring the real scenario where the previous owner transfers to the deployer EOA.
	env, err := environment.New(t.Context(),
		environment.WithEVMSimulatedWithConfig(t, []uint64{selector}, onchain.EVMSimLoaderConfig{
			NumAdditionalAccounts: 1,
		}),
		environment.WithLogger(logger.Test(t)),
	)
	require.NoError(t, err)

	// Deploy a KeystoneForwarder — deployer becomes the initial owner.
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
	require.NotEmpty(t, chain.Users, "expected at least one additional funded account")

	// Users[0] is the new owner — the deployer (current owner) transfers to it.
	newOwner := chain.Users[0]

	contract, err := forwarderwrapper.NewKeystoneForwarder(common.HexToAddress(refs[0].Address), chain.Client)
	require.NoError(t, err)

	tx, err := contract.TransferOwnership(chain.DeployerKey, newOwner.From)
	_, err = cldf.ConfirmIfNoError(chain, tx, err)
	require.NoError(t, err)

	// Rebuild the environment's BlockChains with newOwner as the deployer key
	// so the changeset signs acceptOwnership as the pending owner.
	chain.DeployerKey = newOwner
	evmChains := env.BlockChains.EVMChains()
	evmChains[selector] = chain
	blockChainMap := make(map[uint64]cldfchain.BlockChain, len(evmChains))
	for k, v := range evmChains {
		blockChainMap[k] = v
	}
	env.BlockChains = cldfchain.NewBlockChains(blockChainMap)

	// Apply the changeset — new owner accepts the pending ownership transfer.
	_, err = forwarder.AcceptOwnershipForwarder{}.Apply(*env, forwarder.AcceptOwnershipInput{
		ChainSelector: selector,
		Qualifier:     "test-accept-ownership",
	})
	require.NoError(t, err)

	owner, err := contract.Owner(nil)
	require.NoError(t, err)
	require.Equal(t, newOwner.From, owner)
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
