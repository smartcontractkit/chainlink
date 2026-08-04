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

// DeployCacheRequest configures a DataFeedsCache deployment. The real cache
// contract's __constructor(owner) takes a single argument — there is no
// RetentionTTLLedgers constructor input (that's a hardcoded on-chain constant).
type DeployCacheRequest struct {
	ChainSel  uint64
	WasmPath  string
	Owner     string // defaults to the chain's deployer address when empty
	Qualifier string
	Version   string
	LabelSet  datastore.LabelSet
}

var _ cldf.ChangeSetV2[*DeployCacheRequest] = DeployCache{}

// DeployCache uploads and instantiates the DataFeedsCache contract and records
// its address in the datastore under CacheContract.
type DeployCache struct{}

func (DeployCache) VerifyPreconditions(env cldf.Environment, req *DeployCacheRequest) error {
	if _, ok := env.BlockChains.StellarChains()[req.ChainSel]; !ok {
		return fmt.Errorf("stellar chain not found for chain selector %d", req.ChainSel)
	}
	if _, err := semver.NewVersion(req.Version); err != nil {
		return fmt.Errorf("invalid version %q: %w", req.Version, err)
	}
	if _, err := os.Stat(req.WasmPath); err != nil {
		return fmt.Errorf("wasm path: %w", err)
	}
	if req.Owner != "" {
		if err := validateAddress(req.Owner); err != nil {
			return err
		}
	}
	return nil
}

func (DeployCache) Apply(env cldf.Environment, req *DeployCacheRequest) (cldf.ChangesetOutput, error) {
	var out cldf.ChangesetOutput
	ch, ok := env.BlockChains.StellarChains()[req.ChainSel]
	if !ok {
		return out, fmt.Errorf("stellar chain not found for chain selector %d", req.ChainSel)
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
	salt := stellardeploy.GenerateDeterministicSalt(owner, "data_feeds_cache-"+req.Qualifier)

	report, err := operations.ExecuteOperation(env.OperationsBundle, operation.DeployCache, deps, operation.DeployCacheInput{
		WasmPath: req.WasmPath,
		Salt:     salt,
		Owner:    owner,
	})
	if err != nil {
		return out, err
	}

	out.DataStore = datastore.NewMemoryDataStore()
	return out, out.DataStore.Addresses().Add(datastore.AddressRef{
		Address:       report.Output.ContractID,
		ChainSelector: req.ChainSel,
		Type:          CacheContract,
		Version:       semver.MustParse(req.Version),
		Qualifier:     req.Qualifier,
		Labels:        req.LabelSet,
	})
}
