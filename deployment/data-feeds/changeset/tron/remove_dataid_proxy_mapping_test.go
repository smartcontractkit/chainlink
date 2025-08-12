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

func TestRemoveDataIDProxyMapping(t *testing.T) {
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

	proxyAddress, err := address.Base58ToAddress("TYS5HCEnSU23FgSirvxqVqfwDoD5xHd9Bz")
	require.NoError(t, err)

	dataID := "0x01bb0467f50003040000000000000000"

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
			tron.UpdateDataIDProxyChangeset,
			types.UpdateDataIDProxyTronConfig{
				ChainSelector:  chainSelector,
				CacheAddress:   cacheAddress,
				ProxyAddresses: []address.Address{proxyAddress},
				DataIDs:        []string{dataID},
			},
		),
	)
	require.NoError(t, err)
	require.NotNil(t, resp)

	resp, err = commonChangesets.Apply(t, newEnv,
		commonChangesets.Configure(
			tron.RemoveFeedProxyMappingChangeset,
			types.RemoveFeedProxyTronConfig{
				ChainSelector:  chainSelector,
				CacheAddress:   cacheAddress,
				ProxyAddresses: []address.Address{proxyAddress},
			},
		),
	)
	require.NoError(t, err)
	require.NotNil(t, resp)
}
