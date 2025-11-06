package tron_test

import (
	"testing"

	"github.com/fbsobreira/gotron-sdk/pkg/address"

	"github.com/stretchr/testify/require"

	chain_selectors "github.com/smartcontractkit/chain-selectors"

	cldf_tron "github.com/smartcontractkit/chainlink-deployments-framework/chain/tron"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/engine/test/environment"
	"github.com/smartcontractkit/chainlink-deployments-framework/engine/test/runtime"
	"github.com/smartcontractkit/chainlink/deployment/data-feeds/changeset/tron"

	"github.com/smartcontractkit/chainlink/deployment/data-feeds/changeset/types"
)

func TestSetFeedConfig(t *testing.T) {
	t.Parallel()

	selector := chain_selectors.TRON_DEVNET.Selector
	rt, err := runtime.New(t.Context(), runtime.WithEnvOpts(
		environment.WithTronContainer(t, []uint64{selector}),
	))
	require.NoError(t, err)

	chain := rt.Environment().BlockChains.TronChains()[selector]

	deployOptions := cldf_tron.DefaultDeployOptions()
	deployOptions.FeeLimit = 1_000_000_000

	err = rt.Exec(
		runtime.ChangesetTask(tron.DeployCacheChangeset, types.DeployTronConfig{
			ChainsToDeploy: []uint64{selector},
			Labels:         []string{"data-feeds"},
			Qualifier:      "tron",
			DeployOptions:  deployOptions,
		}),
	)
	require.NoError(t, err)

	cacheAddressStr, err := cldf.SearchAddressBook(rt.State().AddressBook, selector, "DataFeedsCache")
	require.NoError(t, err)

	cacheAddress, err := address.Base58ToAddress(cacheAddressStr)
	require.NoError(t, err)

	dataID := "0x01cb0467f50003040000000000000000"

	allowedSender1, err := address.Base58ToAddress("TYS5HCEnSU23FgSirvxqVqfwDoD5xHd9Bz")
	require.NoError(t, err)
	allowedWorkflowOwner1, err := address.Base58ToAddress("TJatHg7jd3BJ21czkeA1WM76nfaLQ1RUFr")
	require.NoError(t, err)

	allowedSender2, err := address.Base58ToAddress("TSvJFKyg8ZrFyt46mEQTUfwQmY5rTAoCHY")
	require.NoError(t, err)
	allowedWorkflowOwner2, err := address.Base58ToAddress("TV3xgF64Q5bWD4rZjXB2MbKKuXqZuE71Nc")
	require.NoError(t, err)

	workflowMetadata := []types.DataFeedsCacheTronWorkflowMetadata{
		{
			AllowedSender:        allowedSender1,
			AllowedWorkflowOwner: allowedWorkflowOwner1,
			AllowedWorkflowName:  [10]byte{'T', 'e', 's', 't', 'W', 'o', 'r', 'd', '1'},
		},
		{
			AllowedSender:        allowedSender2,
			AllowedWorkflowOwner: allowedWorkflowOwner2,
			AllowedWorkflowName:  [10]byte{'T', 'e', 's', 't', 'W', 'o', 'r', 'd', '2'},
		},
	}

	triggerOpts := cldf_tron.DefaultTriggerOptions()
	triggerOpts.FeeLimit = 1_000_000_000

	err = rt.Exec(
		runtime.ChangesetTask(tron.SetFeedAdminChangeset, types.SetFeedAdminTronConfig{
			ChainSelector: selector,
			CacheAddress:  cacheAddress,
			AdminAddress:  chain.Address,
			IsAdmin:       true,
		}),
		runtime.ChangesetTask(tron.SetFeedConfigChangeset, types.SetFeedDecimalTronConfig{
			ChainSelector:    selector,
			CacheAddress:     cacheAddress,
			DataIDs:          []string{dataID},
			Descriptions:     []string{"test description"},
			WorkflowMetadata: workflowMetadata,
			TriggerOptions:   triggerOpts,
		}),
	)
	require.NoError(t, err)
}
