package stellar

import (
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	cache "github.com/smartcontractkit/chainlink-stellar/bindings/contracts/data_feeds_cache"
	proxy "github.com/smartcontractkit/chainlink-stellar/bindings/contracts/data_feeds_proxy"
)

// LoadCacheClient returns a generated cache client bound to the resolved
// CacheContract AddressRef for (chainSel, qualifier, version).
func LoadCacheClient(env cldf.Environment, chainSel uint64, qualifier, version string) (*cache.DataFeedsCacheClient, datastore.AddressRef, error) {
	d, ref, err := resolveContractDeps(env, chainSel, CacheContract, qualifier, version)
	if err != nil {
		return nil, datastore.AddressRef{}, err
	}
	return cache.NewDataFeedsCacheClient(d.deps.Invoker, d.contractID), ref, nil
}

// LoadProxyClient is LoadCacheClient's counterpart for ProxyContract.
func LoadProxyClient(env cldf.Environment, chainSel uint64, qualifier, version string) (*proxy.DataFeedsProxyClient, datastore.AddressRef, error) {
	d, ref, err := resolveContractDeps(env, chainSel, ProxyContract, qualifier, version)
	if err != nil {
		return nil, datastore.AddressRef{}, err
	}
	return proxy.NewDataFeedsProxyClient(d.deps.Invoker, d.contractID), ref, nil
}
