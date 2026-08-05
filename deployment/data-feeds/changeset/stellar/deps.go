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

// StellarDeps bundles deploy and invoke chain I/O for the operations. The
// same chainlink-stellar Deployer satisfies both fields.
type StellarDeps struct {
	Deploy  SorobanContractDeployer
	Invoker bindings.Invoker
}

// SorobanContractDeployer is the deploy surface needed from chainlink-stellar's
// Deployer, defined locally so upstream need not export it.
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
	Upgrade(ctx context.Context, newWasmHash [32]byte) error
	RecoverTokens(ctx context.Context, token, to string, amount int64) error
}

func adminClient(d StellarDeps, contractID string, isProxy bool) contractAdmin {
	if isProxy {
		return proxy.NewDataFeedsProxyClient(d.Invoker, contractID)
	}
	return cache.NewDataFeedsCacheClient(d.Invoker, contractID)
}

// newStellarDeps is a package-level var so unit tests can substitute fakes.
var newStellarDeps = func(ch cldfstellar.Chain) (StellarDeps, error) {
	d, err := stellardeploy.NewDeployerFromChain(ch)
	if err != nil {
		return StellarDeps{}, err
	}
	return StellarDeps{Deploy: d, Invoker: d}, nil
}

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

func recordAddress(address string, chainSel uint64, contractType datastore.ContractType, qualifier, version string, labels datastore.LabelSet) (cldf.ChangesetOutput, error) {
	var out cldf.ChangesetOutput
	v, err := semver.NewVersion(version)
	if err != nil {
		return out, fmt.Errorf("invalid version %q: %w", version, err)
	}
	out.DataStore = datastore.NewMemoryDataStore()
	return out, out.DataStore.Addresses().Add(datastore.AddressRef{
		Address:       address,
		ChainSelector: chainSel,
		Type:          contractType,
		Version:       v,
		Qualifier:     qualifier,
		Labels:        labels,
	})
}

// stellarApplyDeps pairs a resolved contract address with its chain deps.
type stellarApplyDeps struct {
	deps       StellarDeps
	contractID string
}

// verifyContractRef checks the chain exists, the version parses, and an
// AddressRef exists for the contract.
func verifyContractRef(env cldf.Environment, chainSel uint64, contractType datastore.ContractType, qualifier, version string) error {
	if _, ok := env.BlockChains.StellarChains()[chainSel]; !ok {
		return fmt.Errorf("stellar chain not found for chain selector %d", chainSel)
	}
	_, err := getAddressRef(env, chainSel, contractType, qualifier, version)
	return err
}

// getAddressRef parses version and fetches the contract's AddressRef.
func getAddressRef(env cldf.Environment, chainSel uint64, contractType datastore.ContractType, qualifier, version string) (datastore.AddressRef, error) {
	v, err := semver.NewVersion(version)
	if err != nil {
		return datastore.AddressRef{}, fmt.Errorf("invalid version %q: %w", version, err)
	}
	ref, err := env.DataStore.Addresses().Get(
		datastore.NewAddressRefKey(chainSel, contractType, v, qualifier),
	)
	if err != nil {
		return datastore.AddressRef{}, fmt.Errorf("%s address ref not found for qualifier %q: %w", contractType, qualifier, err)
	}
	return ref, nil
}

// resolveContractDeps resolves the contract's AddressRef and bundles it with
// the chain deps.
func resolveContractDeps(env cldf.Environment, chainSel uint64, contractType datastore.ContractType, qualifier, version string) (stellarApplyDeps, datastore.AddressRef, error) {
	_, deps, err := chainDeps(env, chainSel)
	if err != nil {
		return stellarApplyDeps{}, datastore.AddressRef{}, err
	}
	ref, err := getAddressRef(env, chainSel, contractType, qualifier, version)
	if err != nil {
		return stellarApplyDeps{}, datastore.AddressRef{}, err
	}
	return stellarApplyDeps{deps: deps, contractID: ref.Address}, ref, nil
}
