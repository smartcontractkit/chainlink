package stellar

import (
	"fmt"
	"os"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
)

// UpgradeRequest points the cache or proxy at a new WASM implementation.
type UpgradeRequest struct {
	ChainSel  uint64
	Qualifier string
	Version   string
	Contract  datastore.ContractType // CacheContract or ProxyContract
	WasmPath  string
}

var _ cldf.ChangeSetV2[*UpgradeRequest] = Upgrade{}

// Upgrade uploads a new WASM blob and points the contract at it.
type Upgrade struct{}

type uploadWASMInput struct {
	WasmPath string `json:"wasm_path"`
}

type uploadWASMOutput struct {
	WasmHash [32]byte `json:"wasm_hash"`
}

type upgradeContractInput struct {
	ContractID  string   `json:"contract_id"`
	IsProxy     bool     `json:"is_proxy"`
	NewWasmHash [32]byte `json:"new_wasm_hash"`
}

func (Upgrade) VerifyPreconditions(env cldf.Environment, req *UpgradeRequest) error {
	if err := validateContract(req.Contract); err != nil {
		return err
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
	d, err := resolveDeps(env, req.ChainSel, req.Contract, req.Qualifier, req.Version)
	if err != nil {
		return out, err
	}

	uploadReport, err := operations.ExecuteOperation(env.OperationsBundle, uploadWASMOp, d.deps, uploadWASMInput{
		WasmPath: req.WasmPath,
	})
	if err != nil {
		return out, err
	}

	_, err = operations.ExecuteOperation(env.OperationsBundle, upgradeContractOp, d.deps, upgradeContractInput{
		ContractID:  d.contractID,
		IsProxy:     req.Contract == ProxyContract,
		NewWasmHash: uploadReport.Output.WasmHash,
	})
	return out, err
}

var uploadWASMOp = operations.NewOperation(
	"df:upload-wasm", opVersion,
	"Uploads a WASM blob and returns its code hash (for upgrades)",
	func(b operations.Bundle, d StellarDeps, in uploadWASMInput) (uploadWASMOutput, error) {
		h, err := d.Deploy.UploadContractWASM(b.GetContext(), in.WasmPath)
		if err != nil {
			return uploadWASMOutput{}, err
		}
		return uploadWASMOutput{WasmHash: [32]byte(h)}, nil
	},
)

var upgradeContractOp = operations.NewOperation(
	"df:upgrade", opVersion,
	"Points the contract at a new WASM implementation",
	func(b operations.Bundle, d StellarDeps, in upgradeContractInput) (void, error) {
		return void{}, adminClient(d, in.ContractID, in.IsProxy).Upgrade(b.GetContext(), in.NewWasmHash)
	},
)
