package stellar

import (
	"context"
	"errors"
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.com/stellar/go-stellar-sdk/xdr"

	cldfstellar "github.com/smartcontractkit/chainlink-deployments-framework/chain/stellar"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-stellar/bindings"
	cache "github.com/smartcontractkit/chainlink-stellar/bindings/contracts/data_feeds_cache"
	proxy "github.com/smartcontractkit/chainlink-stellar/bindings/contracts/data_feeds_proxy"
	stellardeploy "github.com/smartcontractkit/chainlink-stellar/deployment"
)

// StellarDeps bundles deploy-time and runtime chain I/O for the operations.
// The same *deployment.Deployer satisfies both fields.
type StellarDeps struct {
	Deploy  SorobanContractDeployer
	Invoker bindings.Invoker
}

// SorobanContractDeployer is the deploy-time surface the operations need from
// chainlink-stellar's concrete deployment.Deployer, defined locally so the
// upstream package does not have to export a widened interface.
type SorobanContractDeployer interface {
	DeployContractWithArgs(ctx context.Context, wasmPath string, salt [32]byte, ctorArgs []xdr.ScVal) (string, error)
	UploadContractWASM(ctx context.Context, wasmPath string) (xdr.Hash, error)
}

// void is the output for operations that return no payload.
type void struct{}

var opVersion = semver.MustParse("1.0.0")

// contractAdmin is the ownership/upgrade/recovery surface shared by the
// generated cache and proxy clients.
type contractAdmin interface {
	TransferOwnership(ctx context.Context, newOwner string, liveUntilLedger uint32) error
	AcceptOwnership(ctx context.Context) error
	RenounceOwnership(ctx context.Context) error
	Upgrade(ctx context.Context, newWasmHash [32]byte) error
	RecoverTokens(ctx context.Context, token, to string, amount int64) error
}

// adminClient returns the generated client selected by isProxy.
func adminClient(d StellarDeps, contractID string, isProxy bool) contractAdmin {
	if isProxy {
		return proxy.NewDataFeedsProxyClient(d.Invoker, contractID)
	}
	return cache.NewDataFeedsCacheClient(d.Invoker, contractID)
}

// newStellarDeps builds the deploy/invoke dependencies for a CLDF Stellar
// chain. Package-level var so unit tests can substitute fakes.
var newStellarDeps = func(ch cldfstellar.Chain) (StellarDeps, error) {
	d, err := stellardeploy.NewDeployerFromChain(ch)
	if err != nil {
		return StellarDeps{}, err
	}
	return StellarDeps{Deploy: d, Invoker: d}, nil
}

// chainDeps resolves the CLDF Stellar chain for chainSel and builds its
// deploy/invoke deps.
func chainDeps(env cldf.Environment, chainSel uint64) (cldfstellar.Chain, StellarDeps, error) {
	ch, ok := env.BlockChains.StellarChains()[chainSel]
	if !ok {
		return cldfstellar.Chain{}, StellarDeps{}, fmt.Errorf("stellar chain not found for chain selector %d", chainSel)
	}
	deps, err := newStellarDeps(ch)
	return ch, deps, err
}

// ownerOrSigner defaults an empty owner to the chain's deployer signer.
func ownerOrSigner(ch cldfstellar.Chain, owner string) (string, error) {
	if owner != "" {
		return owner, nil
	}
	if ch.Signer == nil {
		return "", errors.New("owner not set and chain has no signer")
	}
	return ch.Signer.Address(), nil
}

// recordAddress returns a ChangesetOutput whose datastore holds the deployed
// contract's AddressRef.
func recordAddress(address string, chainSel uint64, contractType datastore.ContractType, qualifier, version string, labels datastore.LabelSet) (cldf.ChangesetOutput, error) {
	var out cldf.ChangesetOutput
	out.DataStore = datastore.NewMemoryDataStore()
	return out, out.DataStore.Addresses().Add(datastore.AddressRef{
		Address:       address,
		ChainSelector: chainSel,
		Type:          contractType,
		Version:       semver.MustParse(version),
		Qualifier:     qualifier,
		Labels:        labels,
	})
}

// stellarApplyDeps bundles a resolved contract's address with the chain deps
// needed to invoke operations against it.
type stellarApplyDeps struct {
	deps       StellarDeps
	contractID string
}

// verifyContractRef checks that chainSel names a known Stellar chain,
// versionStr parses, and an AddressRef exists for (chainSel, contractType,
// version, qualifier).
func verifyContractRef(env cldf.Environment, chainSel uint64, contractType datastore.ContractType, qualifier, versionStr string) error {
	if _, ok := env.BlockChains.StellarChains()[chainSel]; !ok {
		return fmt.Errorf("stellar chain not found for chain selector %d", chainSel)
	}
	version, err := semver.NewVersion(versionStr)
	if err != nil {
		return fmt.Errorf("invalid version %q: %w", versionStr, err)
	}
	if _, err := env.DataStore.Addresses().Get(
		datastore.NewAddressRefKey(chainSel, contractType, version, qualifier),
	); err != nil {
		return fmt.Errorf("%s address ref not found for qualifier %q: %w", contractType, qualifier, err)
	}
	return nil
}

// resolveDeps is resolveContractDeps for callers that don't need the AddressRef.
func resolveDeps(env cldf.Environment, chainSel uint64, contractType datastore.ContractType, qualifier, version string) (stellarApplyDeps, error) {
	d, _, err := resolveContractDeps(env, chainSel, contractType, qualifier, version)
	return d, err
}

// resolveContractDeps looks up the AddressRef for (chainSel, contractType,
// version, qualifier) and bundles it with the chain's deploy/invoke deps.
func resolveContractDeps(env cldf.Environment, chainSel uint64, contractType datastore.ContractType, qualifier, version string) (stellarApplyDeps, datastore.AddressRef, error) {
	v, err := semver.NewVersion(version)
	if err != nil {
		return stellarApplyDeps{}, datastore.AddressRef{}, fmt.Errorf("invalid version %q: %w", version, err)
	}
	_, deps, err := chainDeps(env, chainSel)
	if err != nil {
		return stellarApplyDeps{}, datastore.AddressRef{}, err
	}
	ref, err := env.DataStore.Addresses().Get(
		datastore.NewAddressRefKey(chainSel, contractType, v, qualifier),
	)
	if err != nil {
		return stellarApplyDeps{}, datastore.AddressRef{}, fmt.Errorf("%s address ref not found for qualifier %q: %w", contractType, qualifier, err)
	}
	return stellarApplyDeps{deps: deps, contractID: ref.Address}, ref, nil
}
