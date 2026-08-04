// Package operation contains the CLDF operations for the Stellar Data Feeds
// cache and proxy contracts. Each operation is one on-chain interaction; the
// changesets in the parent package compose and execute them.
package operation

import (
	"context"

	"github.com/Masterminds/semver/v3"

	cldfops "github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink-stellar/bindings"
	cache "github.com/smartcontractkit/chainlink-stellar/bindings/contracts/data_feeds_cache"
	proxy "github.com/smartcontractkit/chainlink-stellar/bindings/contracts/data_feeds_proxy"
	"github.com/smartcontractkit/chainlink-stellar/bindings/scval"
	"github.com/stellar/go-stellar-sdk/xdr"
)

var Version1_0_0 = semver.MustParse("1.0.0")

// SorobanContractDeployer is the deploy-time surface these operations need
// from chainlink-stellar's concrete deployment.Deployer, defined locally so
// the upstream package does not have to export a widened interface.
type SorobanContractDeployer interface {
	DeployContractWithArgs(ctx context.Context, wasmPath string, salt [32]byte, ctorArgs []xdr.ScVal) (string, error)
	UploadContractWASM(ctx context.Context, wasmPath string) (xdr.Hash, error)
}

// StellarDeps bundles deploy-time and runtime chain I/O for these operations.
// The same *deployment.Deployer satisfies both fields.
type StellarDeps struct {
	Deploy  SorobanContractDeployer
	Invoker bindings.Invoker
}

// Void is the output for operations that return no payload.
type Void struct{}

// DeployOutput carries a deployed contract address (C... strkey).
type DeployOutput struct {
	ContractID string `json:"contract_id"`
}

type DeployCacheInput struct {
	WasmPath string   `json:"wasm_path"`
	Salt     [32]byte `json:"salt"`
	Owner    string   `json:"owner"`
}

// DeployCache uploads the cache WASM and instantiates it via CreateContractV2
// with __constructor(owner). The cache constructor takes a single argument;
// the data-retention TTL is a hardcoded on-chain constant, not a ctor input.
var DeployCache = cldfops.NewOperation(
	"df-cache:deploy", Version1_0_0,
	"Deploys the DataFeedsCache Soroban contract",
	func(b cldfops.Bundle, d StellarDeps, in DeployCacheInput) (DeployOutput, error) {
		args := []xdr.ScVal{
			scval.AddressToScVal(in.Owner),
		}
		cid, err := d.Deploy.DeployContractWithArgs(b.GetContext(), in.WasmPath, in.Salt, args)
		if err != nil {
			return DeployOutput{}, err
		}
		return DeployOutput{ContractID: cid}, nil
	},
)

type DeployProxyInput struct {
	WasmPath string   `json:"wasm_path"`
	Salt     [32]byte `json:"salt"`
	Owner    string   `json:"owner"`
	Cache    string   `json:"cache"`
}

// DeployProxy instantiates the proxy via __constructor(owner, cache).
var DeployProxy = cldfops.NewOperation(
	"df-proxy:deploy", Version1_0_0,
	"Deploys the DataFeedsProxy Soroban contract",
	func(b cldfops.Bundle, d StellarDeps, in DeployProxyInput) (DeployOutput, error) {
		args := []xdr.ScVal{
			scval.AddressToScVal(in.Owner),
			scval.AddressToScVal(in.Cache),
		}
		cid, err := d.Deploy.DeployContractWithArgs(b.GetContext(), in.WasmPath, in.Salt, args)
		if err != nil {
			return DeployOutput{}, err
		}
		return DeployOutput{ContractID: cid}, nil
	},
)

type SetFeedConfigsInput struct {
	ContractID string                  `json:"contract_id"`
	Admin      string                  `json:"admin"`
	Entries    []cache.FeedConfigEntry `json:"entries"`
}

var SetFeedConfigs = cldfops.NewOperation(
	"df-cache:set-feed-configs", Version1_0_0,
	"Sets per-feed descriptions and workflow write-permissions on the cache",
	func(b cldfops.Bundle, d StellarDeps, in SetFeedConfigsInput) (Void, error) {
		c := cache.NewDataFeedsCacheClient(d.Invoker, in.ContractID)
		return Void{}, c.SetFeedConfigs(b.GetContext(), in.Admin, in.Entries)
	},
)

type RemoveFeedConfigsInput struct {
	ContractID string     `json:"contract_id"`
	Admin      string     `json:"admin"`
	DataIDs    [][16]byte `json:"data_ids"`
}

var RemoveFeedConfigs = cldfops.NewOperation(
	"df-cache:remove-feed-configs", Version1_0_0,
	"Removes feed configs from the cache",
	func(b cldfops.Bundle, d StellarDeps, in RemoveFeedConfigsInput) (Void, error) {
		c := cache.NewDataFeedsCacheClient(d.Invoker, in.ContractID)
		return Void{}, c.RemoveFeedConfigs(b.GetContext(), in.Admin, in.DataIDs)
	},
)

type FeedAdminInput struct {
	ContractID string `json:"contract_id"`
	Admin      string `json:"admin"`
}

var AddFeedAdmin = cldfops.NewOperation(
	"df-cache:add-feed-admin", Version1_0_0,
	"Grants feed-admin rights on the cache",
	func(b cldfops.Bundle, d StellarDeps, in FeedAdminInput) (Void, error) {
		c := cache.NewDataFeedsCacheClient(d.Invoker, in.ContractID)
		return Void{}, c.AddFeedAdmin(b.GetContext(), in.Admin)
	},
)

var RemoveFeedAdmin = cldfops.NewOperation(
	"df-cache:remove-feed-admin", Version1_0_0,
	"Revokes feed-admin rights on the cache",
	func(b cldfops.Bundle, d StellarDeps, in FeedAdminInput) (Void, error) {
		c := cache.NewDataFeedsCacheClient(d.Invoker, in.ContractID)
		return Void{}, c.RemoveFeedAdmin(b.GetContext(), in.Admin)
	},
)

// Ownership ops work for both contracts: the generated cache and proxy clients
// expose identical two-step-ownership methods, selected by ContractID + IsProxy.
type OwnershipInput struct {
	ContractID      string `json:"contract_id"`
	IsProxy         bool   `json:"is_proxy"`
	NewOwner        string `json:"new_owner"`
	LiveUntilLedger uint32 `json:"live_until_ledger"`
}

var TransferOwnership = cldfops.NewOperation(
	"df:transfer-ownership", Version1_0_0,
	"Begins two-step ownership transfer",
	func(b cldfops.Bundle, d StellarDeps, in OwnershipInput) (Void, error) {
		if in.IsProxy {
			c := proxy.NewDataFeedsProxyClient(d.Invoker, in.ContractID)
			return Void{}, c.TransferOwnership(b.GetContext(), in.NewOwner, in.LiveUntilLedger)
		}
		c := cache.NewDataFeedsCacheClient(d.Invoker, in.ContractID)
		return Void{}, c.TransferOwnership(b.GetContext(), in.NewOwner, in.LiveUntilLedger)
	},
)

var AcceptOwnership = cldfops.NewOperation(
	"df:accept-ownership", Version1_0_0,
	"Accepts a pending ownership transfer (caller must be the pending owner)",
	func(b cldfops.Bundle, d StellarDeps, in OwnershipInput) (Void, error) {
		if in.IsProxy {
			c := proxy.NewDataFeedsProxyClient(d.Invoker, in.ContractID)
			return Void{}, c.AcceptOwnership(b.GetContext())
		}
		c := cache.NewDataFeedsCacheClient(d.Invoker, in.ContractID)
		return Void{}, c.AcceptOwnership(b.GetContext())
	},
)

var RenounceOwnership = cldfops.NewOperation(
	"df:renounce-ownership", Version1_0_0,
	"Renounces ownership permanently",
	func(b cldfops.Bundle, d StellarDeps, in OwnershipInput) (Void, error) {
		if in.IsProxy {
			c := proxy.NewDataFeedsProxyClient(d.Invoker, in.ContractID)
			return Void{}, c.RenounceOwnership(b.GetContext())
		}
		c := cache.NewDataFeedsCacheClient(d.Invoker, in.ContractID)
		return Void{}, c.RenounceOwnership(b.GetContext())
	},
)

type UploadWASMInput struct {
	WasmPath string `json:"wasm_path"`
}
type UploadWASMOutput struct {
	WasmHash [32]byte `json:"wasm_hash"`
}

var UploadWASM = cldfops.NewOperation(
	"df:upload-wasm", Version1_0_0,
	"Uploads a WASM blob and returns its code hash (for upgrades)",
	func(b cldfops.Bundle, d StellarDeps, in UploadWASMInput) (UploadWASMOutput, error) {
		h, err := d.Deploy.UploadContractWASM(b.GetContext(), in.WasmPath)
		if err != nil {
			return UploadWASMOutput{}, err
		}
		return UploadWASMOutput{WasmHash: [32]byte(h)}, nil
	},
)

type UpgradeCacheInput struct {
	ContractID  string   `json:"contract_id"`
	NewWasmHash [32]byte `json:"new_wasm_hash"`
}

var UpgradeCache = cldfops.NewOperation(
	"df-cache:upgrade", Version1_0_0,
	"Points the cache at a new WASM implementation",
	func(b cldfops.Bundle, d StellarDeps, in UpgradeCacheInput) (Void, error) {
		c := cache.NewDataFeedsCacheClient(d.Invoker, in.ContractID)
		return Void{}, c.Upgrade(b.GetContext(), in.NewWasmHash)
	},
)

type RecoverTokensInput struct {
	ContractID string `json:"contract_id"`
	Token      string `json:"token"`
	To         string `json:"to"`
	Amount     int64  `json:"amount"`
}

var RecoverTokens = cldfops.NewOperation(
	"df-cache:recover-tokens", Version1_0_0,
	"Recovers tokens accidentally sent to the cache",
	func(b cldfops.Bundle, d StellarDeps, in RecoverTokensInput) (Void, error) {
		c := cache.NewDataFeedsCacheClient(d.Invoker, in.ContractID)
		return Void{}, c.RecoverTokens(b.GetContext(), in.Token, in.To, in.Amount)
	},
)

type SetProxyCacheInput struct {
	ContractID string `json:"contract_id"`
	Cache      string `json:"cache"`
}

var SetProxyCache = cldfops.NewOperation(
	"df-proxy:set-cache", Version1_0_0,
	"Points the proxy at a cache contract",
	func(b cldfops.Bundle, d StellarDeps, in SetProxyCacheInput) (Void, error) {
		c := proxy.NewDataFeedsProxyClient(d.Invoker, in.ContractID)
		return Void{}, c.SetCache(b.GetContext(), in.Cache)
	},
)
