package changeset

import (
	"crypto/ecdsa"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	chain_selectors "github.com/smartcontractkit/chain-selectors"
	"github.com/stretchr/testify/require"

	cldfproposalutils "github.com/smartcontractkit/chainlink-deployments-framework/engine/cld/mcms/proposalutils"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/engine/test/environment"
	"github.com/smartcontractkit/chainlink-deployments-framework/engine/test/runtime"

	mcmschangesets "github.com/smartcontractkit/cld-changesets/legacy/mcms/changesets"

	cldftesthelpers "github.com/smartcontractkit/chainlink-deployments-framework/engine/cld/mcms/proposalutils/testhelpers"

	"github.com/smartcontractkit/chainlink/deployment/data-feeds/changeset/types"
)

func TestTransferOwnership(t *testing.T) {
	t.Parallel()

	selector := chain_selectors.TEST_90000001.Selector
	rt, err := runtime.New(t.Context(), runtime.WithEnvOpts(
		environment.WithEVMSimulated(t, []uint64{selector}),
		environment.WithLogger(logger.Test(t)),
	))
	require.NoError(t, err)

	chain := rt.Environment().BlockChains.EVMChains()[selector]

	MCMScfg := cldftesthelpers.SingleGroupTimelockConfig(t)
	MCMSQualifier := "MCMS_EVM_1"
	MCMScfg.Qualifier = &MCMSQualifier

	err = rt.Exec(
		runtime.ChangesetTask(cldf.CreateLegacyChangeSet(
			mcmschangesets.DeployMCMSWithTimelockV2), map[uint64]cldfproposalutils.MCMSWithTimelockConfig{
			selector: MCMScfg,
		}),
	)
	require.NoError(t, err)

	records := rt.Environment().DataStore.Addresses().Filter(datastore.AddressRefByType("RBACTimelock"))
	require.Len(t, records, 1)
	timeLockAddress := common.HexToAddress(records[0].Address)

	mcmsConfig := &types.MCMSConfig{
		MinDelay:          1,
		TimeLockQualifier: MCMSQualifier,
	}

	cache, err := DeployCache(chain, []string{})
	require.NoError(t, err)
	cacheAddress := cache.Contract.Address()

	// The cache is still owned by the deployer, so the precondition must reject the transfer.
	err = TransferOwnershipChangeset.VerifyPreconditions(rt.Environment(), types.TransferOwnershipConfig{
		ChainSelector:     selector,
		ContractAddresses: []common.Address{cacheAddress},
		NewOwnerAddress:   common.HexToAddress("0x1"),
		McmsConfig:        mcmsConfig,
	})
	require.ErrorContains(t, err, "not by the timelock")

	// Hand the cache to the timelock so there is something to transfer back.
	tx, err := cache.Contract.TransferOwnership(chain.DeployerKey, timeLockAddress)
	require.NoError(t, err)
	_, err = chain.Confirm(tx)
	require.NoError(t, err)

	err = rt.Exec(
		runtime.ChangesetTask(AcceptOwnershipChangeset, types.AcceptOwnershipConfig{
			ChainSelector:     selector,
			ContractAddresses: []common.Address{cacheAddress},
			McmsConfig:        mcmsConfig,
		}),
		runtime.SignAndExecuteProposalsTask([]*ecdsa.PrivateKey{cldftesthelpers.TestXXXMCMSSigner}),
	)
	require.NoError(t, err)

	owner, err := cache.Contract.Owner(nil)
	require.NoError(t, err)
	require.Equal(t, timeLockAddress, owner)

	t.Run("rejects invalid new owner", func(t *testing.T) {
		err := TransferOwnershipChangeset.VerifyPreconditions(rt.Environment(), types.TransferOwnershipConfig{
			ChainSelector:     selector,
			ContractAddresses: []common.Address{cacheAddress},
			NewOwnerAddress:   common.Address{},
			McmsConfig:        mcmsConfig,
		})
		require.ErrorContains(t, err, "zero address")

		err = TransferOwnershipChangeset.VerifyPreconditions(rt.Environment(), types.TransferOwnershipConfig{
			ChainSelector:     selector,
			ContractAddresses: []common.Address{cacheAddress},
			NewOwnerAddress:   timeLockAddress,
			McmsConfig:        mcmsConfig,
		})
		require.ErrorContains(t, err, "already owned by")

		err = TransferOwnershipChangeset.VerifyPreconditions(rt.Environment(), types.TransferOwnershipConfig{
			ChainSelector:     selector,
			ContractAddresses: []common.Address{},
			NewOwnerAddress:   chain.DeployerKey.From,
			McmsConfig:        mcmsConfig,
		})
		require.ErrorContains(t, err, "at least one contract address")
	})

	// Transfer ownership from the timelock back to the deployer via an MCMS proposal.
	err = rt.Exec(
		runtime.ChangesetTask(TransferOwnershipChangeset, types.TransferOwnershipConfig{
			ChainSelector:     selector,
			ContractAddresses: []common.Address{cacheAddress},
			NewOwnerAddress:   chain.DeployerKey.From,
			McmsConfig:        mcmsConfig,
		}),
		runtime.SignAndExecuteProposalsTask([]*ecdsa.PrivateKey{cldftesthelpers.TestXXXMCMSSigner}),
	)
	require.NoError(t, err)
	require.Len(t, rt.State().Proposals, 2)
	require.True(t, rt.State().Proposals[1].IsExecuted)

	// Ownership is two-step: the timelock still owns the cache until the deployer accepts.
	owner, err = cache.Contract.Owner(nil)
	require.NoError(t, err)
	require.Equal(t, timeLockAddress, owner)

	tx, err = cache.Contract.AcceptOwnership(chain.DeployerKey)
	require.NoError(t, err)
	_, err = chain.Confirm(tx)
	require.NoError(t, err)

	owner, err = cache.Contract.Owner(nil)
	require.NoError(t, err)
	require.Equal(t, chain.DeployerKey.From, owner)
}
