package stellar

import (
	"fmt"

	"github.com/Masterminds/semver/v3"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	proxy "github.com/smartcontractkit/chainlink-stellar/bindings/contracts/data_feeds_proxy"
)

// SetProxyCacheRequest points the proxy at a new cache. Qualifier resolves
// the proxy; CacheQualifier resolves the cache.
type SetProxyCacheRequest struct {
	ChainSel       uint64
	Qualifier      string
	Version        string
	CacheQualifier string
}

var _ cldf.ChangeSetV2[*SetProxyCacheRequest] = SetProxyCache{}

// SetProxyCache points the proxy at a cache contract.
type SetProxyCache struct{}

type setProxyCacheInput struct {
	ContractID string `json:"contract_id"`
	Cache      string `json:"cache"`
}

func (SetProxyCache) VerifyPreconditions(env cldf.Environment, req *SetProxyCacheRequest) error {
	if err := verifyContractRef(env, req.ChainSel, ProxyContract, req.Qualifier, req.Version); err != nil {
		return err
	}
	if err := verifyContractRef(env, req.ChainSel, CacheContract, req.CacheQualifier, req.Version); err != nil {
		return err
	}
	return nil
}

func (SetProxyCache) Apply(env cldf.Environment, req *SetProxyCacheRequest) (cldf.ChangesetOutput, error) {
	var out cldf.ChangesetOutput
	proxyDeps, err := resolveDeps(env, req.ChainSel, ProxyContract, req.Qualifier, req.Version)
	if err != nil {
		return out, err
	}
	cacheRef, err := env.DataStore.Addresses().Get(
		datastore.NewAddressRefKey(req.ChainSel, CacheContract, semver.MustParse(req.Version), req.CacheQualifier),
	)
	if err != nil {
		return out, fmt.Errorf("cache address ref not found for qualifier %q: %w", req.CacheQualifier, err)
	}

	_, err = operations.ExecuteOperation(env.OperationsBundle, setProxyCacheOp, proxyDeps.deps, setProxyCacheInput{
		ContractID: proxyDeps.contractID,
		Cache:      cacheRef.Address,
	})
	return out, err
}

var setProxyCacheOp = operations.NewOperation(
	"df-proxy:set-cache", opVersion,
	"Points the proxy at a cache contract",
	func(b operations.Bundle, d StellarDeps, in setProxyCacheInput) (void, error) {
		c := proxy.NewDataFeedsProxyClient(d.Invoker, in.ContractID)
		return void{}, c.SetCache(b.GetContext(), in.Cache)
	},
)
