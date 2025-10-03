package view

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	chainselectors "github.com/smartcontractkit/chain-selectors"
	"github.com/stretchr/testify/require"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/engine/test/environment"
	"github.com/smartcontractkit/chainlink-deployments-framework/engine/test/runtime"

	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	commontypes "github.com/smartcontractkit/chainlink/deployment/common/types"
	"github.com/smartcontractkit/chainlink/deployment/vault/changeset"
	"github.com/smartcontractkit/chainlink/deployment/vault/changeset/types"
)

func TestVault_NoChains(t *testing.T) {
	t.Parallel()

	env, err := environment.New(t.Context())
	require.NoError(t, err)

	viewMarshaler, err := Vault(*env, nil)
	require.NoError(t, err)
	require.NotNil(t, viewMarshaler)

	view := viewMarshaler.(*VaultView)
	require.Empty(t, view.TimelockBalances)
	require.Empty(t, view.WhitelistedAddresses)
	require.Empty(t, view.MCMSWithTimelock)
}

func TestGenerateVaultView_WithoutTimelock(t *testing.T) {
	t.Parallel()

	selectors := []uint64{chainselectors.TEST_90000001.Selector}
	env, err := environment.New(t.Context(),
		environment.WithEVMSimulated(t, selectors),
	)
	require.NoError(t, err)

	view, err := GenerateVaultView(*env, selectors)
	require.NoError(t, err)
	require.NotNil(t, view)

	require.Empty(t, view.TimelockBalances)

	require.Len(t, view.WhitelistedAddresses, len(selectors))
	for _, sel := range selectors {
		require.Empty(t, view.WhitelistedAddresses[sel])
	}

	require.Empty(t, view.MCMSWithTimelock)
}

func TestGenerateVaultView_WithMCMSAndWhitelist(t *testing.T) {
	t.Parallel()

	selectors := []uint64{chainselectors.TEST_90000001.Selector, chainselectors.TEST_90000002.Selector}
	rt, err := runtime.New(t.Context(), runtime.WithEnvOpts(
		environment.WithEVMSimulated(t, selectors),
	))
	require.NoError(t, err)

	setupMCMS(t, rt, selectors)

	whitelistByChain := map[uint64][]types.WhitelistAddress{}
	for i, sel := range selectors {
		addr := common.HexToAddress("0x1111111111111111111111111111111111111111")
		if i == 1 {
			addr = common.HexToAddress("0x2222222222222222222222222222222222222222")
		}
		whitelistByChain[sel] = []types.WhitelistAddress{{
			Address:     addr.Hex(),
			Description: "recipient",
			Labels:      []string{"test"},
		}}
	}

	err = rt.Exec(
		runtime.ChangesetTask(
			changeset.SetWhitelistChangeset,
			types.SetWhitelistConfig{WhitelistByChain: whitelistByChain},
		),
	)
	require.NoError(t, err)

	view, err := GenerateVaultView(rt.Environment(), selectors)
	require.NoError(t, err)
	require.NotNil(t, view)

	require.Len(t, view.TimelockBalances, len(selectors))
	require.Len(t, view.WhitelistedAddresses, len(selectors))
	for _, sel := range selectors {
		entries := view.WhitelistedAddresses[sel]
		require.Len(t, entries, 1)
		require.Equal(t, whitelistByChain[sel][0].Address, entries[0].Address)
	}
	require.Len(t, view.MCMSWithTimelock, len(selectors))
}

func setupMCMS(t *testing.T, rt *runtime.Runtime, chainSelectors []uint64) {
	t.Helper()

	timelockCfgs := make(map[uint64]commontypes.MCMSWithTimelockConfigV2)
	for _, sel := range chainSelectors {
		timelockCfgs[sel] = proposalutils.SingleGroupTimelockConfigV2(t)
	}

	err := rt.Exec(
		runtime.ChangesetTask(
			cldf.CreateLegacyChangeSet(commonchangeset.DeployMCMSWithTimelockV2),
			timelockCfgs,
		),
	)
	require.NoError(t, err)
}
