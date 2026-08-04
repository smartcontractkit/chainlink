package stellar

import (
	"errors"
	"fmt"
	"os"

	"github.com/Masterminds/semver/v3"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	stellardeploy "github.com/smartcontractkit/chainlink-stellar/deployment"

	"github.com/smartcontractkit/chainlink/deployment/data-feeds/changeset/stellar/operation"
)

// DeployProxyRequest configures a DataFeedsProxy deployment. The cache passed
// to __constructor(owner, cache) is resolved from the datastore by (ChainSel,
// CacheContract, Version, CacheQualifier) — it must already be recorded.
// Cross-contract lookups share the request's Version (see README).
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
		return errors.New("cache qualifier must be set")
	}
	return nil
}

func (DeployProxy) Apply(env cldf.Environment, req *DeployProxyRequest) (cldf.ChangesetOutput, error) {
	var out cldf.ChangesetOutput
	ch, deps, err := chainDeps(env, req.ChainSel)
	if err != nil {
		return out, err
	}
	cacheRef, err := env.DataStore.Addresses().Get(
		datastore.NewAddressRefKey(req.ChainSel, CacheContract, semver.MustParse(req.Version), req.CacheQualifier),
	)
	if err != nil {
		return out, fmt.Errorf("cache address ref not found for qualifier %q: %w", req.CacheQualifier, err)
	}
	owner, err := ownerOrSigner(ch, req.Owner)
	if err != nil {
		return out, err
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
	return recordAddress(report.Output.ContractID, req.ChainSel, ProxyContract, req.Qualifier, req.Version, req.LabelSet)
}
