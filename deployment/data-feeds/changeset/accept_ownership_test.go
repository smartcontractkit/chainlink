package changeset

import (
	"crypto/ecdsa"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	chain_selectors "github.com/smartcontractkit/chain-selectors"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/engine/test/environment"
	"github.com/smartcontractkit/chainlink-deployments-framework/engine/test/runtime"

	commonChangesets "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	commonTypes "github.com/smartcontractkit/chainlink/deployment/common/types"
	"github.com/smartcontractkit/chainlink/deployment/data-feeds/changeset/types"
)

func TestAcceptOwnership(t *testing.T) {
	t.Parallel()

	selector := chain_selectors.TEST_90000001.Selector
	rt, err := runtime.New(t.Context(), runtime.WithEnvOpts(
		environment.WithEVMSimulated(t, []uint64{selector}),
		environment.WithLogger(logger.Test(t)),
	))
	require.NoError(t, err)

	chain := rt.Environment().BlockChains.EVMChains()[selector]

	err = rt.Exec(
		runtime.ChangesetTask(cldf.CreateLegacyChangeSet(
			commonChangesets.DeployMCMSWithTimelockV2), map[uint64]commonTypes.MCMSWithTimelockConfigV2{
			selector: proposalutils.SingleGroupTimelockConfigV2(t),
		}),
	)
	require.NoError(t, err)

	records := rt.Environment().DataStore.Addresses().Filter(datastore.AddressRefByType("RBACTimelock"))
	require.Len(t, records, 1)
	timeLockAddress := records[0].Address

	cache, _ := DeployCache(chain, []string{})
	tx, _ := cache.Contract.TransferOwnership(chain.DeployerKey, common.HexToAddress(timeLockAddress))
	_, err = chain.Confirm(tx)
	require.NoError(t, err)

	err = rt.Exec(
		runtime.ChangesetTask(AcceptOwnershipChangeset, types.AcceptOwnershipConfig{
			ChainSelector:     selector,
			ContractAddresses: []common.Address{cache.Contract.Address()},
			McmsConfig: &types.MCMSConfig{
				MinDelay: 1,
			},
		}),
		runtime.SignAndExecuteProposalsTask([]*ecdsa.PrivateKey{proposalutils.TestXXXMCMSSigner}),
	)
	require.NoError(t, err)
	require.Len(t, rt.State().Proposals, 1)
	require.True(t, rt.State().Proposals[0].IsExecuted)
}
