package stellar

import (
	"fmt"
	"os"

	"github.com/Masterminds/semver/v3"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink/deployment/data-feeds/changeset/stellar/operation"
)

// UpgradeCacheRequest points an already-deployed DataFeedsCache contract at a
// new WASM implementation.
type UpgradeCacheRequest struct {
	ChainSel  uint64
	Qualifier string
	Version   string
	WasmPath  string
}

var _ cldf.ChangeSetV2[*UpgradeCacheRequest] = UpgradeCache{}

// UpgradeCache uploads a new WASM blob and points the cache at it. Apply
// chains two operations: UploadWASM produces the code hash that UpgradeCache
// then applies on-chain.
type UpgradeCache struct{}

func (UpgradeCache) VerifyPreconditions(env cldf.Environment, req *UpgradeCacheRequest) error {
	if err := verifyContractRef(env, req.ChainSel, CacheContract, req.Qualifier, req.Version); err != nil {
		return err
	}
	if _, err := os.Stat(req.WasmPath); err != nil {
		return fmt.Errorf("wasm path: %w", err)
	}
	return nil
}

func (UpgradeCache) Apply(env cldf.Environment, req *UpgradeCacheRequest) (cldf.ChangesetOutput, error) {
	var out cldf.ChangesetOutput
	version := semver.MustParse(req.Version)
	d, _, err := resolveContractDeps(env, req.ChainSel, CacheContract, req.Qualifier, version)
	if err != nil {
		return out, err
	}

	uploadReport, err := operations.ExecuteOperation(env.OperationsBundle, operation.UploadWASM, d.deps, operation.UploadWASMInput{
		WasmPath: req.WasmPath,
	})
	if err != nil {
		return out, err
	}

	_, err = operations.ExecuteOperation(env.OperationsBundle, operation.UpgradeCache, d.deps, operation.UpgradeCacheInput{
		ContractID:  d.contractID,
		NewWasmHash: uploadReport.Output.WasmHash,
	})
	return out, err
}
