package stellar

import (
	"fmt"

	"github.com/Masterminds/semver/v3"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink/deployment/data-feeds/changeset/stellar/operation"
)

// SetProxyCacheRequest points an already-deployed DataFeedsProxy contract at a
// (possibly different) cache contract. Qualifier resolves the proxy being
// acted on (ChainSel, ProxyContract, Version, Qualifier); CacheQualifier
// resolves the new cache target (ChainSel, CacheContract, Version,
// CacheQualifier) — both lookups share this request's Version: per this
// suite's convention, cross-contract lookups use the acting changeset's
// datastore Version, so a cache recorded under a different version needs a
// matching-version record (an optional CacheVersion field is a recorded
// follow-up if operators need split versions).
type SetProxyCacheRequest struct {
	ChainSel       uint64
	Qualifier      string
	Version        string
	CacheQualifier string
}

var _ cldf.ChangeSetV2[*SetProxyCacheRequest] = SetProxyCache{}

// SetProxyCache points the proxy at a cache contract.
type SetProxyCache struct{}

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
	version := semver.MustParse(req.Version)

	proxyDeps, _, err := resolveContractDeps(env, req.ChainSel, ProxyContract, req.Qualifier, version)
	if err != nil {
		return out, err
	}
	cacheRef, err := env.DataStore.Addresses().Get(
		datastore.NewAddressRefKey(req.ChainSel, CacheContract, version, req.CacheQualifier),
	)
	if err != nil {
		return out, fmt.Errorf("cache address ref not found for qualifier %q: %w", req.CacheQualifier, err)
	}

	_, err = operations.ExecuteOperation(env.OperationsBundle, operation.SetProxyCache, proxyDeps.deps, operation.SetProxyCacheInput{
		ContractID: proxyDeps.contractID,
		Cache:      cacheRef.Address,
	})
	return out, err
}
