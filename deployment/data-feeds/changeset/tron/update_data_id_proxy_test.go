package tron_test

import (
	"testing"

	"github.com/fbsobreira/gotron-sdk/pkg/address"
	chain_selectors "github.com/smartcontractkit/chain-selectors"
	"github.com/stretchr/testify/require"

	cldf_tron "github.com/smartcontractkit/chainlink-deployments-framework/chain/tron"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/engine/test/environment"
	"github.com/smartcontractkit/chainlink-deployments-framework/engine/test/runtime"

	"github.com/smartcontractkit/chainlink/deployment/data-feeds/changeset/tron"
	"github.com/smartcontractkit/chainlink/deployment/data-feeds/changeset/types"
)

func TestUpdateDataIDProxy(t *testing.T) {
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

	proxyAddress, err := address.Base58ToAddress("TYS5HCEnSU23FgSirvxqVqfwDoD5xHd9Bz")
	require.NoError(t, err)

	dataID := "0x01bb0467f50003040000000000000000"

	err = rt.Exec(
		runtime.ChangesetTask(tron.SetFeedAdminChangeset, types.SetFeedAdminTronConfig{
			ChainSelector: selector,
			CacheAddress:  cacheAddress,
			AdminAddress:  chain.Address,
			IsAdmin:       true,
		}),
		runtime.ChangesetTask(tron.UpdateDataIDProxyChangeset, types.UpdateDataIDProxyTronConfig{
			ChainSelector:  selector,
			CacheAddress:   cacheAddress,
			ProxyAddresses: []address.Address{proxyAddress},
			DataIDs:        []string{dataID},
		}),
	)
	require.NoError(t, err)
}
