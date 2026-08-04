package stellar

import (
	"errors"
	"fmt"

	"github.com/Masterminds/semver/v3"

	cldfstellar "github.com/smartcontractkit/chainlink-deployments-framework/chain/stellar"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	stellardeploy "github.com/smartcontractkit/chainlink-stellar/deployment"

	"github.com/smartcontractkit/chainlink/deployment/data-feeds/changeset/stellar/operation"
)

// newStellarDeps builds the deploy/invoke dependencies for a CLDF Stellar
// chain. Package-level var so unit tests can substitute fakes.
var newStellarDeps = func(ch cldfstellar.Chain) (operation.StellarDeps, error) {
	d, err := stellardeploy.NewDeployerFromChain(ch)
	if err != nil {
		return operation.StellarDeps{}, err
	}
	return operation.StellarDeps{Deploy: d, Invoker: d}, nil
}

// chainDeps resolves the CLDF Stellar chain for chainSel and builds its
// deploy/invoke deps.
func chainDeps(env cldf.Environment, chainSel uint64) (cldfstellar.Chain, operation.StellarDeps, error) {
	ch, ok := env.BlockChains.StellarChains()[chainSel]
	if !ok {
		return cldfstellar.Chain{}, operation.StellarDeps{}, fmt.Errorf("stellar chain not found for chain selector %d", chainSel)
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
	deps       operation.StellarDeps
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
