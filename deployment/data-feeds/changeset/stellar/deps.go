package stellar

import (
	"fmt"

	"github.com/Masterminds/semver/v3"

	cldfstellar "github.com/smartcontractkit/chainlink-deployments-framework/chain/stellar"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	stellardeploy "github.com/smartcontractkit/chainlink-stellar/deployment"

	"github.com/smartcontractkit/chainlink/deployment/data-feeds/changeset/stellar/operation"
)

// newStellarDeps builds the deploy/invoke dependencies for a CLDF Stellar chain.
// Package-level var so unit tests can substitute fakes (the Invoker seam).
var newStellarDeps = func(ch cldfstellar.Chain) (operation.StellarDeps, error) {
	d, err := stellardeploy.NewDeployerFromChain(ch)
	if err != nil {
		return operation.StellarDeps{}, err
	}
	return operation.StellarDeps{Deploy: d, Invoker: d}, nil
}

// stellarApplyDeps bundles a resolved contract's address with the chain deps
// needed to invoke an operation against it — shared by every stellar
// changeset that acts on a single already-deployed contract.
type stellarApplyDeps struct {
	deps       operation.StellarDeps
	contractID string
}

// verifyContractRef checks that chainSel names a known Stellar chain,
// versionStr parses, and an AddressRef exists for (chainSel, contractType,
// version, qualifier). This is the chain/version/ref-existence skeleton
// shared by every changeset's VerifyPreconditions; only the extra checks
// differ.
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

// resolveContractDeps looks up the AddressRef for (chainSel, contractType,
// version, qualifier) and bundles it with the chain's deploy/invoke deps,
// returning the AddressRef itself alongside for callers that need more than
// just its Address (e.g. LoadCacheClient/LoadProxyClient). version must
// already be parsed (see semver.MustParse in callers that have validated
// req.Version in VerifyPreconditions, typically via verifyContractRef).
func resolveContractDeps(env cldf.Environment, chainSel uint64, contractType datastore.ContractType, qualifier string, version *semver.Version) (stellarApplyDeps, datastore.AddressRef, error) {
	var zero stellarApplyDeps
	var zeroRef datastore.AddressRef
	ch, ok := env.BlockChains.StellarChains()[chainSel]
	if !ok {
		return zero, zeroRef, fmt.Errorf("stellar chain not found for chain selector %d", chainSel)
	}

	ref, err := env.DataStore.Addresses().Get(
		datastore.NewAddressRefKey(chainSel, contractType, version, qualifier),
	)
	if err != nil {
		return zero, zeroRef, fmt.Errorf("%s address ref not found for qualifier %q: %w", contractType, qualifier, err)
	}

	deps, err := newStellarDeps(ch)
	if err != nil {
		return zero, zeroRef, err
	}

	return stellarApplyDeps{deps: deps, contractID: ref.Address}, ref, nil
}
