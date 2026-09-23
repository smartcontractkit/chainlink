package stateview

import (
	"fmt"
	"sort"

	ccipshared "github.com/smartcontractkit/chainlink/deployment/ccip/shared"

	"github.com/smartcontractkit/ccip-owner-contracts/gethwrappers"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	evmstate "github.com/smartcontractkit/cld-changesets/legacy/pkg/family/evm"
)

// Aliases of the shared datastore helpers.
const (
	DefaultMCMSQualifier = ccipshared.DefaultMCMSQualifier
	SupersededLabel      = ccipshared.SupersededLabel
)

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

func dataStoreRefsForChain(e cldf.Environment, selector uint64) ([]datastore.AddressRef, error) {
	if e.DataStore == nil {
		return nil, fmt.Errorf("datastore not available for chain %d", selector)
	}
	refs, err := e.DataStore.Addresses().Fetch()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch address refs from datastore: %w", err)
	}
	chainRefs := make([]datastore.AddressRef, 0, len(refs))
	for _, ref := range refs {
		if ref.ChainSelector == selector && isActiveRef(ref) {
			chainRefs = append(chainRefs, ref)
		}
	}
	if err := ccipshared.CheckRefUniqueness(chainRefs); err != nil {
		return nil, fmt.Errorf("chain %d datastore is ambiguous: %w", selector, err)
	}
	return chainRefs, nil
}

// DataStoreTypeVersionsForChain returns the chain's active refs in slice-per-address form.
// It fails when two active refs claim the same (type, version, qualifier) identity.
func DataStoreTypeVersionsForChain(e cldf.Environment, selector uint64) (map[string][]cldf.TypeAndVersion, error) {
	refs, err := dataStoreRefsForChain(e, selector)
	if err != nil {
		return nil, err
	}
	return AddressRefsToTypeVersions(refs), nil
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

// loadMCMSWithTimelockForChain resolves the MCMS bundle for one chain from the datastore.
// Qualifier-isolated: only the requested qualifier is loaded.
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
	bundle, err := ccipshared.MCMSBundleRefs(refs, selector, qualifier)
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
