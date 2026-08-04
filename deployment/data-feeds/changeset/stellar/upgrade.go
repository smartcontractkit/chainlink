package stellar

import (
	"fmt"
	"os"

	"github.com/Masterminds/semver/v3"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink/deployment/data-feeds/changeset/stellar/operation"
)

// UpgradeRequest points an already-deployed cache or proxy contract at a new
// WASM implementation.
type UpgradeRequest struct {
	ChainSel  uint64
	Qualifier string
	Version   string
	Contract  datastore.ContractType // CacheContract or ProxyContract
	WasmPath  string
}

var _ cldf.ChangeSetV2[*UpgradeRequest] = Upgrade{}

// Upgrade uploads a new WASM blob and points the contract at it. Apply chains
// two operations: UploadWASM produces the code hash that UpgradeContract then
// applies on-chain.
type Upgrade struct{}

func (Upgrade) VerifyPreconditions(env cldf.Environment, req *UpgradeRequest) error {
	if req.Contract != CacheContract && req.Contract != ProxyContract {
		return fmt.Errorf("unsupported contract type %q: must be %q or %q", req.Contract, CacheContract, ProxyContract)
	}
	if err := verifyContractRef(env, req.ChainSel, req.Contract, req.Qualifier, req.Version); err != nil {
		return err
	}
	if _, err := os.Stat(req.WasmPath); err != nil {
		return fmt.Errorf("wasm path: %w", err)
	}
	return nil
}

func (Upgrade) Apply(env cldf.Environment, req *UpgradeRequest) (cldf.ChangesetOutput, error) {
	var out cldf.ChangesetOutput
	version := semver.MustParse(req.Version)
	d, _, err := resolveContractDeps(env, req.ChainSel, req.Contract, req.Qualifier, version)
	if err != nil {
		return out, err
	}

	uploadReport, err := operations.ExecuteOperation(env.OperationsBundle, operation.UploadWASM, d.deps, operation.UploadWASMInput{
		WasmPath: req.WasmPath,
	})
	if err != nil {
		return out, err
	}

	_, err = operations.ExecuteOperation(env.OperationsBundle, operation.UpgradeContract, d.deps, operation.UpgradeContractInput{
		ContractID:  d.contractID,
		IsProxy:     req.Contract == ProxyContract,
		NewWasmHash: uploadReport.Output.WasmHash,
	})
	return out, err
}
