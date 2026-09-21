package stateview

import (
	"errors"
	"fmt"
	"sort"

	"github.com/smartcontractkit/chainlink/deployment"
	ccipshared "github.com/smartcontractkit/chainlink/deployment/ccip/shared"

	"github.com/smartcontractkit/ccip-owner-contracts/gethwrappers"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	evmstate "github.com/smartcontractkit/cld-changesets/legacy/pkg/family/evm"

	deployutil "github.com/smartcontractkit/chainlink-ccip/deployment/utils"
)

// DefaultMCMSQualifier qualifies the default CCIP MCMS bundle; further bundles (RMN,
// ultra-fast-curse, operator) resolve by their own qualifier.
const DefaultMCMSQualifier = deployutil.CLLQualifier

// SupersededLabel marks rows an operator declared historical; they never load as active contracts.
const SupersededLabel = "superseded"

// isActiveRef reports whether a ref may load into chain state.
func isActiveRef(ref datastore.AddressRef) bool {
	return ref.Version != nil && !ref.Labels.Contains(SupersededLabel)
}

// AddressRefsToTypeVersions reshapes refs into the slice-per-address form LoadChainState
// consumes, preserving labels and skipping versionless refs.
func AddressRefsToTypeVersions(refs []datastore.AddressRef) map[string][]cldf.TypeAndVersion {
	addresses := make(map[string][]cldf.TypeAndVersion)
	for _, ref := range refs {
		if ref.Version == nil {
			continue
		}
		tv := cldf.TypeAndVersion{
			Type:    cldf.ContractType(ref.Type),
			Version: *ref.Version,
		}
		if !ref.Labels.IsEmpty() {
			tv.Labels = cldf.NewLabelSet(ref.Labels.List()...)
		}
		addresses[ref.Address] = append(addresses[ref.Address], tv)
	}
	return addresses
}

// AddressBookTypeVersionsForChain returns the address book refs of one chain in the
// slice-per-address form LoadChainState consumes. A chain missing from the address book
// yields an empty map, mirroring AddressesForChain's ErrChainNotFound tolerance.
func AddressBookTypeVersionsForChain(e cldf.Environment, selector uint64) (map[string][]cldf.TypeAndVersion, error) {
	flat, err := e.ExistingAddresses.AddressesForChain(selector) //nolint:staticcheck // AddressBook remains a supported source until Phase 3.
	if err != nil {
		if errors.Is(err, cldf.ErrChainNotFound) {
			return map[string][]cldf.TypeAndVersion{}, nil
		}
		return nil, fmt.Errorf("failed to get addresses for chain %d: %w", selector, err)
	}
	return TypeVersionsToSlices(flat), nil
}

// DataStoreTypeVersionsForChain returns the chain's active refs in slice-per-address form.
// Fails when two active refs share (type, version, qualifier) but differ by address — the
// qualifier is the instance identity, so that is corruption. Same (type, version) under
// different qualifiers is legitimate (lanes, pools, MCMS bundles, FeeQuoter instances).
func DataStoreTypeVersionsForChain(e cldf.Environment, selector uint64) (map[string][]cldf.TypeAndVersion, error) {
	if e.DataStore == nil {
		return nil, fmt.Errorf("datastore not available for chain %d", selector)
	}
	refs, err := e.DataStore.Addresses().Fetch()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch address refs from datastore: %w", err)
	}
	var chainRefs []datastore.AddressRef
	for _, ref := range refs {
		if ref.ChainSelector == selector && isActiveRef(ref) {
			chainRefs = append(chainRefs, ref)
		}
	}
	if err := checkRefUniqueness(chainRefs); err != nil {
		return nil, fmt.Errorf("chain %d datastore is ambiguous: %w", selector, err)
	}
	return AddressRefsToTypeVersions(chainRefs), nil
}

// checkRefUniqueness rejects two active refs claiming one (type, version, qualifier) identity
// with different addresses.
func checkRefUniqueness(refs []datastore.AddressRef) error {
	seen := make(map[string]string)
	for _, ref := range refs {
		key := fmt.Sprintf("%s %s %s", ref.Type, ref.Version, ref.Qualifier)
		if prev, dup := seen[key]; dup && prev != ref.Address {
			return fmt.Errorf("%s %s qualified %q resolves to both %s and %s", ref.Type, ref.Version, ref.Qualifier, prev, ref.Address)
		}
		seen[key] = ref.Address
	}
	return nil
}

// TypeVersionsToSlices reshapes a flat address→TypeAndVersion map into the slice-per-address
// form LoadChainState consumes.
func TypeVersionsToSlices(addresses map[string]cldf.TypeAndVersion) map[string][]cldf.TypeAndVersion {
	sliced := make(map[string][]cldf.TypeAndVersion, len(addresses))
	for addr, tv := range addresses {
		sliced[addr] = append(sliced[addr], tv)
	}
	return sliced
}

// flattenTypeVersions collapses the slice-per-address form back to one TypeAndVersion per
// address, for helpers that still consume the address-book shape.
func flattenTypeVersions(addresses map[string][]cldf.TypeAndVersion) map[string]cldf.TypeAndVersion {
	flat := make(map[string]cldf.TypeAndVersion, len(addresses))
	for addr, tvs := range addresses {
		if len(tvs) > 0 {
			flat[addr] = tvs[0]
		}
	}
	return flat
}

// singularTypeVersions lists the (type, version) pairs whose loader case assigns a singular
// state field: a second ref at a different address is unrecoverable ambiguity. Types absent
// from this set are either keyed by symbol/label/remote selector, resolved by qualifier, or
// not modeled at all (skipped in datastore mode) — none can corrupt state by multiplicity.
// USDCTokenPoolProxy is deliberately absent: it is per-token multi-instance in the datastore;
// the legacy state keeps one deterministically (sorted order) and the parity harness validates
// the pick matches what the address book resolves.
var singularTypeVersions = map[string]bool{
	dispatchKey(cldf.NewTypeAndVersion(ccipshared.OnRamp, deployment.Version1_6_0)):                         true,
	dispatchKey(cldf.NewTypeAndVersion(ccipshared.OffRamp, deployment.Version1_6_0)):                        true,
	dispatchKey(cldf.NewTypeAndVersion(ccipshared.ARMProxy, deployment.Version1_0_0)):                       true,
	dispatchKey(cldf.NewTypeAndVersion(ccipshared.RMNRemote, deployment.Version1_6_0)):                      true,
	dispatchKey(cldf.NewTypeAndVersion(ccipshared.RMNHome, deployment.Version1_6_0)):                        true,
	dispatchKey(cldf.NewTypeAndVersion(ccipshared.WETH9, deployment.Version1_0_0)):                          true,
	dispatchKey(cldf.NewTypeAndVersion(ccipshared.NonceManager, deployment.Version1_6_0)):                   true,
	dispatchKey(cldf.NewTypeAndVersion(ccipshared.TokenAdminRegistry, deployment.Version1_5_0)):             true,
	dispatchKey(cldf.NewTypeAndVersion(ccipshared.TokenPoolFactory, deployment.Version1_5_1)):               true,
	dispatchKey(cldf.NewTypeAndVersion(ccipshared.Router, deployment.Version1_2_0)):                         true,
	dispatchKey(cldf.NewTypeAndVersion(ccipshared.TestRouter, deployment.Version1_2_0)):                     true,
	dispatchKey(cldf.NewTypeAndVersion(ccipshared.CapabilitiesRegistry, deployment.Version1_0_0)):           true,
	dispatchKey(cldf.NewTypeAndVersion(ccipshared.CCIPHome, deployment.Version1_6_0)):                       true,
	dispatchKey(cldf.NewTypeAndVersion(ccipshared.CCIPReceiver, deployment.Version1_0_0)):                   true,
	dispatchKey(cldf.NewTypeAndVersion(ccipshared.LogMessageDataReceiver, deployment.Version1_0_0)):         true,
	dispatchKey(cldf.NewTypeAndVersion(ccipshared.Multicall3, deployment.Version1_0_0)):                     true,
	dispatchKey(cldf.NewTypeAndVersion(ccipshared.USDCToken, deployment.Version1_0_0)):                      true,
	dispatchKey(cldf.NewTypeAndVersion(ccipshared.USDCTokenPool, deployment.Version1_5_1)):                  true,
	dispatchKey(cldf.NewTypeAndVersion(ccipshared.USDCTokenPool, deployment.Version1_6_2)):                  true,
	dispatchKey(cldf.NewTypeAndVersion(ccipshared.HybridLockReleaseUSDCTokenPool, deployment.Version1_5_1)): true,
	dispatchKey(cldf.NewTypeAndVersion(ccipshared.HybridLockReleaseUSDCTokenPool, deployment.Version1_6_2)): true,
	dispatchKey(cldf.NewTypeAndVersion(ccipshared.CCTPMessageTransmitterProxy, deployment.Version1_6_2)):    true,
	dispatchKey(cldf.NewTypeAndVersion(ccipshared.USDCMockTransmitter, deployment.Version1_0_0)):            true,
	dispatchKey(cldf.NewTypeAndVersion(ccipshared.USDCTokenMessenger, deployment.Version1_0_0)):             true,
	dispatchKey(cldf.NewTypeAndVersion(ccipshared.FactoryBurnMintERC20Token, deployment.Version1_6_2)):      true,
	dispatchKey(cldf.NewTypeAndVersion(ccipshared.FactoryBurnMintERC20Token, deployment.Version1_5_1)):      true,
	dispatchKey(cldf.NewTypeAndVersion(ccipshared.FeeAggregator, deployment.Version1_0_0)):                  true,
	dispatchKey(cldf.NewTypeAndVersion(ccipshared.DonIDClaimer, deployment.Version1_6_1)):                   true,
	dispatchKey(cldf.NewTypeAndVersion(ccipshared.EVMSignerRegistry, deployment.Version1_0_0)):              true,
	dispatchKey(cldf.NewTypeAndVersion(ccipshared.MockRMN, deployment.Version1_0_0)):                        true,
	dispatchKey(cldf.NewTypeAndVersion(ccipshared.PriceRegistry, deployment.Version1_2_0)):                  true,
	dispatchKey(cldf.NewTypeAndVersion(ccipshared.RMN, deployment.Version1_5_0)):                            true,
}

// loadMCMSWithTimelockForChain resolves the MCMS bundle for one chain from the datastore.
// Qualifier-isolated: requested qualifier exclusively, falling back to the empty qualifier
// (legacy/test deployments) only when it holds no rows — never a mix. No refs yield an empty
// state; a partial or ambiguous bundle errors.
func loadMCMSWithTimelockForChain(e cldf.Environment, selector uint64, qualifier string) (*evmstate.MCMSWithTimelockState, error) {
	chain, ok := e.BlockChains.EVMChains()[selector]
	if !ok {
		return nil, fmt.Errorf("chain %d not found", selector)
	}
	if e.DataStore == nil {
		return nil, fmt.Errorf("datastore not available for chain %d", selector)
	}
	refs, err := e.DataStore.Addresses().Fetch()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch address refs from datastore: %w", err)
	}
	bundle, err := mcmsBundleRefs(refs, selector, qualifier)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve mcms bundle for chain %d with qualifier %s: %w", selector, qualifier, err)
	}
	if len(bundle) == 0 {
		return &evmstate.MCMSWithTimelockState{}, nil
	}
	state, err := evmstate.MaybeLoadMCMSWithTimelockChainStateFromRefs(chain, bundle)
	if err != nil {
		return nil, fmt.Errorf("failed to load mcms state from datastore for chain %d with qualifier %s: %w", selector, qualifier, err)
	}
	return state, nil
}

// mcmsBundleRefs selects the chain's bundle rows. Falls back to the empty qualifier (legacy/
// test deployments) only for the default bundle; custom qualifiers fail closed. Skips
// versionless and superseded refs (the ref loader dereferences Version).
func mcmsBundleRefs(refs []datastore.AddressRef, selector uint64, qualifier string) ([]datastore.AddressRef, error) {
	active := make([]datastore.AddressRef, 0, len(refs))
	for _, ref := range refs {
		if ref.ChainSelector == selector && isActiveRef(ref) {
			active = append(active, ref)
		}
	}
	qualifiers := []string{qualifier}
	if qualifier == DefaultMCMSQualifier {
		qualifiers = append(qualifiers, "")
	}
	for _, q := range qualifiers {
		var bundle []datastore.AddressRef
		for _, ref := range active {
			if ref.Qualifier == q {
				bundle = append(bundle, ref)
			}
		}
		if len(bundle) > 0 {
			sortAddressRefs(bundle)
			if err := checkRefUniqueness(bundle); err != nil {
				return nil, err
			}
			return bundle, nil
		}
	}
	if qualifier != DefaultMCMSQualifier {
		return nil, fmt.Errorf("no mcms refs for chain %d with qualifier %q", selector, qualifier)
	}
	return nil, nil
}

// applyMCMSWithTimelockState writes a loaded MCMS state into the chain state, keeping the ABI
// registry in sync.
func (c CCIPOnChainState) applyMCMSWithTimelockState(selector uint64, mcmsState *evmstate.MCMSWithTimelockState) {
	if mcmsState == nil {
		return
	}
	chainState, ok := c.EVMChainState(selector)
	if !ok {
		return
	}
	chainState.MCMSWithTimelockState = *mcmsState
	if mcmsState.ProposerMcm != nil {
		chainState.ABIByAddress[mcmsState.ProposerMcm.Address().Hex()] = gethwrappers.ManyChainMultiSigABI
	}
	if mcmsState.CancellerMcm != nil {
		chainState.ABIByAddress[mcmsState.CancellerMcm.Address().Hex()] = gethwrappers.ManyChainMultiSigABI
	}
	if mcmsState.BypasserMcm != nil {
		chainState.ABIByAddress[mcmsState.BypasserMcm.Address().Hex()] = gethwrappers.ManyChainMultiSigABI
	}
	if mcmsState.Timelock != nil {
		chainState.ABIByAddress[mcmsState.Timelock.Address().Hex()] = gethwrappers.RBACTimelockABI
	}
	if mcmsState.CallProxy != nil {
		chainState.ABIByAddress[mcmsState.CallProxy.Address().Hex()] = gethwrappers.CallProxyABI
	}
	c.WriteEVMChainState(selector, chainState)
}

// sortAddressRefs orders refs deterministically (chain, type, version, qualifier, address) so
// callers that pick one ref from many never depend on map iteration order.
func sortAddressRefs(refs []datastore.AddressRef) {
	sort.Slice(refs, func(i, j int) bool {
		a, b := refs[i], refs[j]
		if a.ChainSelector != b.ChainSelector {
			return a.ChainSelector < b.ChainSelector
		}
		if a.Type != b.Type {
			return a.Type < b.Type
		}
		if a.Version != nil && b.Version != nil && !a.Version.Equal(b.Version) {
			return a.Version.LessThan(b.Version)
		}
		if a.Qualifier != b.Qualifier {
			return a.Qualifier < b.Qualifier
		}
		return a.Address < b.Address
	})
}
