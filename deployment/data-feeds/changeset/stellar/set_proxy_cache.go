package stellar

import (
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
	CacheVersion   string // cache's datastore version; defaults to Version
}

// cacheVersion defaults an empty CacheVersion to Version: cache and proxy
// usually share a release.
func (req *SetProxyCacheRequest) cacheVersion() string {
	if req.CacheVersion != "" {
		return req.CacheVersion
	}
	return req.Version
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
	return verifyContractRef(env, req.ChainSel, CacheContract, req.CacheQualifier, req.cacheVersion())
}

func (SetProxyCache) Apply(env cldf.Environment, req *SetProxyCacheRequest) (cldf.ChangesetOutput, error) {
	var out cldf.ChangesetOutput
	proxyDeps, _, err := resolveContractDeps(env, req.ChainSel, ProxyContract, req.Qualifier, req.Version)
	if err != nil {
		return out, err
	}
	cacheRef, err := getAddressRef(env, req.ChainSel, CacheContract, req.CacheQualifier, req.cacheVersion())
	if err != nil {
		return out, err
	}

	_, err = operations.ExecuteOperation(env.OperationsBundle, setProxyCacheOp, proxyDeps.deps, setProxyCacheInput{
		ContractID: proxyDeps.contractID,
		Cache:      cacheRef.Address,
	})
	if err != nil {
		return out, err
	}
	return metadataOutput(env, req.ChainSel, proxyDeps.contractID, func(m *ContractMetadata) {
		m.Cache = cacheRef.Address
	})
}

var setProxyCacheOp = operations.NewOperation(
	"df-proxy:set-cache", opVersion,
	"Points the proxy at a cache contract",
	func(b operations.Bundle, d StellarDeps, in setProxyCacheInput) (void, error) {
		c := proxy.NewDataFeedsProxyClient(d.Invoker, in.ContractID)
		return void{}, c.SetCache(b.GetContext(), in.Cache)
	},
)
