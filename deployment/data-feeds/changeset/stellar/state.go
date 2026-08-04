package stellar

import (
	"fmt"

	"github.com/Masterminds/semver/v3"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	cache "github.com/smartcontractkit/chainlink-stellar/bindings/contracts/data_feeds_cache"
	proxy "github.com/smartcontractkit/chainlink-stellar/bindings/contracts/data_feeds_proxy"
)

// loadContractClientDeps parses version and resolves the contract's AddressRef
// plus chain deps — the shared lookup behind LoadCacheClient and LoadProxyClient.
func loadContractClientDeps(env cldf.Environment, chainSel uint64, contractType datastore.ContractType, qualifier, version string) (stellarApplyDeps, datastore.AddressRef, error) {
	v, err := semver.NewVersion(version)
	if err != nil {
		return stellarApplyDeps{}, datastore.AddressRef{}, fmt.Errorf("invalid version %q: %w", version, err)
	}
	return resolveContractDeps(env, chainSel, contractType, qualifier, v)
}

// LoadCacheClient returns a generated cache client bound to the resolved
// CacheContract AddressRef for (chainSel, qualifier, version).
func LoadCacheClient(env cldf.Environment, chainSel uint64, qualifier, version string) (*cache.DataFeedsCacheClient, datastore.AddressRef, error) {
	d, ref, err := loadContractClientDeps(env, chainSel, CacheContract, qualifier, version)
	if err != nil {
		return nil, datastore.AddressRef{}, err
	}
	return cache.NewDataFeedsCacheClient(d.deps.Invoker, d.contractID), ref, nil
}

// LoadProxyClient is LoadCacheClient's counterpart for ProxyContract.
func LoadProxyClient(env cldf.Environment, chainSel uint64, qualifier, version string) (*proxy.DataFeedsProxyClient, datastore.AddressRef, error) {
	d, ref, err := loadContractClientDeps(env, chainSel, ProxyContract, qualifier, version)
	if err != nil {
		return nil, datastore.AddressRef{}, err
	}
	return proxy.NewDataFeedsProxyClient(d.deps.Invoker, d.contractID), ref, nil
}
