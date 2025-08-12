package tron_test

import (
	"testing"

	"github.com/fbsobreira/gotron-sdk/pkg/address"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	chain_selectors "github.com/smartcontractkit/chain-selectors"

	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"
	cldf_tron "github.com/smartcontractkit/chainlink-deployments-framework/chain/tron"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink/deployment/data-feeds/changeset/tron"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	commonChangesets "github.com/smartcontractkit/chainlink/deployment/common/changeset"

	"github.com/smartcontractkit/chainlink/deployment/data-feeds/changeset/types"
	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
)

func TestSetFeedConfig(t *testing.T) {
	t.Parallel()
	lggr := logger.Test(t)
	cfg := memory.MemoryEnvironmentConfig{
		TronChains: 1,
	}
	env := memory.NewMemoryEnvironment(t, lggr, zapcore.DebugLevel, cfg)

	chainSelector := env.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chain_selectors.FamilyTron))[0]

	deployOptions := cldf_tron.DefaultDeployOptions()
	deployOptions.FeeLimit = 1_000_000_000

	newEnv, err := commonChangesets.Apply(t, env, commonChangesets.Configure(
		tron.DeployCacheChangeset,
		types.DeployTronConfig{
			ChainsToDeploy: []uint64{chainSelector},
			Labels:         []string{"data-feeds"},
			Qualifier:      "tron",
			DeployOptions:  deployOptions,
		},
	))
	require.NoError(t, err)

	cacheAddressStr, err := cldf.SearchAddressBook(newEnv.ExistingAddresses, chainSelector, "DataFeedsCache")
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

	resp, err := commonChangesets.Apply(t, newEnv,
		commonChangesets.Configure(
			tron.SetFeedAdminChangeset,
			types.SetFeedAdminTronConfig{
				ChainSelector: chainSelector,
				CacheAddress:  cacheAddress,
				AdminAddress:  env.BlockChains.TronChains()[chainSelector].Address,
				IsAdmin:       true,
			},
		),
		commonChangesets.Configure(
			tron.SetFeedConfigChangeset,
			types.SetFeedDecimalTronConfig{
				ChainSelector:    chainSelector,
				CacheAddress:     cacheAddress,
				DataIDs:          []string{dataID},
				Descriptions:     []string{"test description"},
				WorkflowMetadata: workflowMetadata,
				TriggerOptions:   triggerOpts,
			},
		),
	)
	require.NoError(t, err)
	require.NotNil(t, resp)
}
