package stellar

import (
	"errors"
	"fmt"
	"os"

	"github.com/Masterminds/semver/v3"
	"github.com/stellar/go-stellar-sdk/xdr"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink-stellar/bindings/scval"
	stellardeploy "github.com/smartcontractkit/chainlink-stellar/deployment"
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
	report, err := operations.ExecuteOperation(env.OperationsBundle, deployCacheOp, deps, deployCacheInput{
		WasmPath: req.WasmPath,
		Salt:     salt,
		Owner:    owner,
	})
	if err != nil {
		return out, err
	}
	return recordAddress(report.Output.ContractID, req.ChainSel, CacheContract, req.Qualifier, req.Version, req.LabelSet)
}

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
	report, err := operations.ExecuteOperation(env.OperationsBundle, deployProxyOp, deps, deployProxyInput{
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

type deployOutput struct {
	ContractID string `json:"contract_id"`
}

type deployCacheInput struct {
	WasmPath string   `json:"wasm_path"`
	Salt     [32]byte `json:"salt"`
	Owner    string   `json:"owner"`
}

// deployCacheOp uploads the cache WASM and instantiates it via CreateContractV2
// with __constructor(owner); the data-retention TTL is an on-chain constant,
// not a constructor input.
var deployCacheOp = operations.NewOperation(
	"df-cache:deploy", opVersion,
	"Deploys the DataFeedsCache Soroban contract",
	func(b operations.Bundle, d StellarDeps, in deployCacheInput) (deployOutput, error) {
		args := []xdr.ScVal{
			scval.AddressToScVal(in.Owner),
		}
		cid, err := d.Deploy.DeployContractWithArgs(b.GetContext(), in.WasmPath, in.Salt, args)
		if err != nil {
			return deployOutput{}, err
		}
		return deployOutput{ContractID: cid}, nil
	},
)

type deployProxyInput struct {
	WasmPath string   `json:"wasm_path"`
	Salt     [32]byte `json:"salt"`
	Owner    string   `json:"owner"`
	Cache    string   `json:"cache"`
}

// deployProxyOp instantiates the proxy via __constructor(owner, cache).
var deployProxyOp = operations.NewOperation(
	"df-proxy:deploy", opVersion,
	"Deploys the DataFeedsProxy Soroban contract",
	func(b operations.Bundle, d StellarDeps, in deployProxyInput) (deployOutput, error) {
		args := []xdr.ScVal{
			scval.AddressToScVal(in.Owner),
			scval.AddressToScVal(in.Cache),
		}
		cid, err := d.Deploy.DeployContractWithArgs(b.GetContext(), in.WasmPath, in.Salt, args)
		if err != nil {
			return deployOutput{}, err
		}
		return deployOutput{ContractID: cid}, nil
	},
)
