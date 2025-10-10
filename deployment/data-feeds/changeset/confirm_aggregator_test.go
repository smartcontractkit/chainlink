package changeset_test

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
	"github.com/smartcontractkit/chainlink/deployment/data-feeds/changeset"
	"github.com/smartcontractkit/chainlink/deployment/data-feeds/changeset/types"
)

func TestConfirmAggregator(t *testing.T) {
	t.Parallel()

	selector := chain_selectors.TEST_90000001.Selector
	rt, err := runtime.New(t.Context(), runtime.WithEnvOpts(
		environment.WithEVMSimulated(t, []uint64{selector}),
		environment.WithLogger(logger.Test(t)),
	))
	require.NoError(t, err)

	// without MCMS

	err = rt.Exec(
		runtime.ChangesetTask(changeset.DeployCacheChangeset, types.DeployConfig{
			ChainsToDeploy: []uint64{selector},
			Labels:         []string{"data-feeds"},
		}),
		runtime.ChangesetTask(changeset.DeployAggregatorProxyChangeset, types.DeployAggregatorProxyConfig{
			ChainsToDeploy:   []uint64{selector},
			AccessController: []common.Address{common.HexToAddress("0x")},
		}),
	)
	require.NoError(t, err)

	records := rt.Environment().DataStore.Addresses().Filter(datastore.AddressRefByType("AggregatorProxy"))
	require.Len(t, records, 1)
	proxyAddress := records[0].Address

	err = rt.Exec(
		runtime.ChangesetTask(changeset.ProposeAggregatorChangeset, types.ProposeConfirmAggregatorConfig{
			ChainSelector:        selector,
			ProxyAddress:         common.HexToAddress(proxyAddress),
			NewAggregatorAddress: common.HexToAddress("0x123"),
		}),
		runtime.ChangesetTask(changeset.ConfirmAggregatorChangeset, types.ProposeConfirmAggregatorConfig{
			ChainSelector:        selector,
			ProxyAddress:         common.HexToAddress(proxyAddress),
			NewAggregatorAddress: common.HexToAddress("0x123"),
		}),
		runtime.ChangesetTask(cldf.CreateLegacyChangeSet(commonChangesets.DeployMCMSWithTimelockV2), map[uint64]commonTypes.MCMSWithTimelockConfigV2{
			selector: proposalutils.SingleGroupTimelockConfigV2(t),
		}),
	)
	require.NoError(t, err)

	// with MCMS

	err = rt.Exec(
		runtime.ChangesetTask(changeset.ProposeAggregatorChangeset, types.ProposeConfirmAggregatorConfig{
			ChainSelector:        selector,
			ProxyAddress:         common.HexToAddress(proxyAddress),
			NewAggregatorAddress: common.HexToAddress("0x124"),
		}),
		runtime.ChangesetTask(cldf.CreateLegacyChangeSet(commonChangesets.TransferToMCMSWithTimelockV2), commonChangesets.TransferToMCMSWithTimelockConfig{
			ContractsByChain: map[uint64][]common.Address{
				selector: {common.HexToAddress(proxyAddress)},
			},
			MCMSConfig: proposalutils.TimelockConfig{MinDelay: 0},
		}),
		runtime.SignAndExecuteProposalsTask([]*ecdsa.PrivateKey{proposalutils.TestXXXMCMSSigner}),
	)
	require.NoError(t, err)
	require.Len(t, rt.State().Proposals, 1)
	require.True(t, rt.State().Proposals[0].IsExecuted)

	err = rt.Exec(
		runtime.ChangesetTask(changeset.ConfirmAggregatorChangeset, types.ProposeConfirmAggregatorConfig{
			ChainSelector:        selector,
			ProxyAddress:         common.HexToAddress(proxyAddress),
			NewAggregatorAddress: common.HexToAddress("0x124"),
			McmsConfig: &types.MCMSConfig{
				MinDelay: 0,
			},
		}),
		runtime.SignAndExecuteProposalsTask([]*ecdsa.PrivateKey{proposalutils.TestXXXMCMSSigner}),
	)
	require.NoError(t, err)
	require.Len(t, rt.State().Proposals, 2)
	require.True(t, rt.State().Proposals[1].IsExecuted)

	// We expect 8 outputs for the 8 changesets we ran.
	require.Len(t, rt.State().Outputs, 8)
}
