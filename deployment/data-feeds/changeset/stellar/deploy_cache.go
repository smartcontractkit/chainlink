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

// DeployCacheRequest configures a DataFeedsCache deployment.
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
	ch, deps, err := chainDeps(env, req.ChainSel)
	if err != nil {
		return out, err
	}
	owner, err := ownerOrSigner(ch, req.Owner)
	if err != nil {
		return out, err
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
	return recordAddress(report.Output.ContractID, req.ChainSel, CacheContract, req.Qualifier, req.Version, req.LabelSet)
}
