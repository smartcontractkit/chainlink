package stellar

import (
	"fmt"
	"os"

	"github.com/Masterminds/semver/v3"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	stellardeploy "github.com/smartcontractkit/chainlink-stellar/deployment"

	"github.com/smartcontractkit/chainlink/deployment/data-feeds/changeset/stellar/operation"
)

// DeployProxyRequest configures a DataFeedsProxy deployment. The proxy's
// __constructor(owner, cache) resolves the cache address from the datastore
// by (ChainSel, CacheContract, Version, CacheQualifier) — the cache must
// already be deployed and recorded. The cache lookup reuses this request's
// own Version (there is no separate CacheVersion field): per this suite's
// convention, cross-contract lookups share the acting changeset's datastore
// Version, so a cache recorded under a different version needs a
// matching-version record (an optional CacheVersion field is a recorded
// follow-up if operators need split versions).
type DeployProxyRequest struct {
	ChainSel       uint64
	WasmPath       string
	Owner          string // defaults to the chain's deployer address when empty
	CacheQualifier string
	Qualifier      string
	Version        string
	LabelSet       datastore.LabelSet
}

var _ cldf.ChangeSetV2[*DeployProxyRequest] = DeployProxy{}

// DeployProxy uploads and instantiates the DataFeedsProxy contract and records
// its address in the datastore under ProxyContract.
type DeployProxy struct{}

func (DeployProxy) VerifyPreconditions(env cldf.Environment, req *DeployProxyRequest) error {
	if err := verifyContractRef(env, req.ChainSel, CacheContract, req.CacheQualifier, req.Version); err != nil {
		return err
	}
	if _, err := os.Stat(req.WasmPath); err != nil {
		return fmt.Errorf("wasm path: %w", err)
	}
	if req.Owner != "" {
		if err := validateAddress(req.Owner); err != nil {
			return err
		}
	}
	if req.CacheQualifier == "" {
		return fmt.Errorf("cache qualifier must be set")
	}
	return nil
}

func (DeployProxy) Apply(env cldf.Environment, req *DeployProxyRequest) (cldf.ChangesetOutput, error) {
	var out cldf.ChangesetOutput
	ch, ok := env.BlockChains.StellarChains()[req.ChainSel]
	if !ok {
		return out, fmt.Errorf("stellar chain not found for chain selector %d", req.ChainSel)
	}

	version := semver.MustParse(req.Version)
	cacheRef, err := env.DataStore.Addresses().Get(
		datastore.NewAddressRefKey(req.ChainSel, CacheContract, version, req.CacheQualifier),
	)
	if err != nil {
		return out, fmt.Errorf("cache address ref not found for qualifier %q: %w", req.CacheQualifier, err)
	}

	deps, err := newStellarDeps(ch)
	if err != nil {
		return out, err
	}

	owner := req.Owner
	if owner == "" {
		if ch.Signer == nil {
			return out, fmt.Errorf("owner not set and chain has no signer")
		}
		owner = ch.Signer.Address()
	}
	salt := stellardeploy.GenerateDeterministicSalt(owner, "data_feeds_proxy-"+req.Qualifier)

	report, err := operations.ExecuteOperation(env.OperationsBundle, operation.DeployProxy, deps, operation.DeployProxyInput{
		WasmPath: req.WasmPath,
		Salt:     salt,
		Owner:    owner,
		Cache:    cacheRef.Address,
	})
	if err != nil {
		return out, err
	}

	out.DataStore = datastore.NewMemoryDataStore()
	return out, out.DataStore.Addresses().Add(datastore.AddressRef{
		Address:       report.Output.ContractID,
		ChainSelector: req.ChainSel,
		Type:          ProxyContract,
		Version:       version,
		Qualifier:     req.Qualifier,
		Labels:        req.LabelSet,
	})
}
