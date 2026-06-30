package changeset_test

import (
	"testing"

	chainsel "github.com/smartcontractkit/chain-selectors"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	cldfchain "github.com/smartcontractkit/chainlink-deployments-framework/chain"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/engine/test/environment"
	"github.com/smartcontractkit/chainlink-deployments-framework/engine/test/onchain"

	forwarderwrapper "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/forwarder_1_0_0"

	"github.com/smartcontractkit/chainlink/deployment/cre/common/changeset"
)

func TestAcceptOwnershipEOA(t *testing.T) {
	t.Parallel()

	selector := chainsel.TEST_90000001.Selector

	env, err := environment.New(t.Context(),
		environment.WithEVMSimulatedWithConfig(t, []uint64{selector}, onchain.EVMSimLoaderConfig{
			NumAdditionalAccounts: 1,
		}),
		environment.WithLogger(logger.Test(t)),
	)
	require.NoError(t, err)

	chain := env.BlockChains.EVMChains()[selector]

	// Deploy a KeystoneForwarder — deployer becomes the initial owner.
	addr, deployTx, contract, err := forwarderwrapper.DeployKeystoneForwarder(chain.DeployerKey, chain.Client)
	_, err = cldf.ConfirmIfNoError(chain, deployTx, err)
	require.NoError(t, err)

	require.NotEmpty(t, chain.Users, "expected at least one additional funded account")
	newOwner := chain.Users[0]

	// Deployer transfers ownership to the new owner EOA.
	tx, err := contract.TransferOwnership(chain.DeployerKey, newOwner.From)
	_, err = cldf.ConfirmIfNoError(chain, tx, err)
	require.NoError(t, err)

	// Rebuild the environment so acceptOwnership is signed by the pending owner.
	chain.DeployerKey = newOwner
	evmChains := env.BlockChains.EVMChains()
	evmChains[selector] = chain
	blockChainMap := make(map[uint64]cldfchain.BlockChain, len(evmChains))
	for k, v := range evmChains {
		blockChainMap[k] = v
	}
	env.BlockChains = cldfchain.NewBlockChains(blockChainMap)

	_, err = changeset.AcceptOwnershipEOA{}.Apply(*env, changeset.AcceptOwnershipInput{
		ChainSelector:   selector,
		ContractAddress: addr.Hex(),
	})
	require.NoError(t, err)

	owner, err := contract.Owner(nil)
	require.NoError(t, err)
	require.Equal(t, newOwner.From, owner)
}

func TestAcceptOwnershipEOA_VerifyPreconditions(t *testing.T) {
	t.Parallel()

	selector := chainsel.TEST_90000001.Selector

	env, err := environment.New(t.Context(),
		environment.WithEVMSimulated(t, []uint64{selector}),
		environment.WithLogger(logger.Test(t)),
	)
	require.NoError(t, err)

	t.Run("unknown chain selector", func(t *testing.T) {
		err := changeset.AcceptOwnershipEOA{}.VerifyPreconditions(*env, changeset.AcceptOwnershipInput{
			ChainSelector: 0,
		})
		require.Error(t, err)
		require.Contains(t, err.Error(), "not found in environment")
	})

	t.Run("empty contract address", func(t *testing.T) {
		err := changeset.AcceptOwnershipEOA{}.VerifyPreconditions(*env, changeset.AcceptOwnershipInput{
			ChainSelector:   selector,
			ContractAddress: "",
		})
		require.Error(t, err)
		require.Contains(t, err.Error(), "contractAddress is required")
	})
}
