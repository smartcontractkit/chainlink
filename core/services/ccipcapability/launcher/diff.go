package launcher

import (
	"fmt"

	"github.com/ethereum/go-ethereum/crypto"

	ragep2ptypes "github.com/smartcontractkit/libocr/ragep2p/types"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils"
	kcr "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/capabilities_registry"
	"github.com/smartcontractkit/chainlink/v2/core/services/registrysyncer"
)

// diffResult contains the added, removed and updated CCIP DONs.
// It is determined by using the `diff` function below.
type diffResult struct {
	added   map[registrysyncer.DonID]kcr.CapabilitiesRegistryDONInfo
	removed map[registrysyncer.DonID]kcr.CapabilitiesRegistryDONInfo
	updated map[registrysyncer.DonID]kcr.CapabilitiesRegistryDONInfo
}

// diff compares the old and new state and returns the added, removed and updated CCIP DONs.
func diff(
	capabilityVersion,
	capabilityLabelledName string,
	oldState,
	newState registrysyncer.State,
) (diffResult, error) {
	ccipCapability, err := checkCapabilityPresence(capabilityVersion, capabilityLabelledName, newState)
	if err != nil {
		return diffResult{}, fmt.Errorf("failed to check capability presence: %w", err)
	}

	newCCIPDONs, err := filterCCIPDONs(ccipCapability, newState)
	if err != nil {
		return diffResult{}, fmt.Errorf("failed to filter CCIP DONs from new state: %w", err)
	}

	currCCIPDONs, err := filterCCIPDONs(ccipCapability, oldState)
	if err != nil {
		return diffResult{}, fmt.Errorf("failed to filter CCIP DONs from old state: %w", err)
	}

	// compare curr with new and launch or update OCR instances as needed
	diffRes, err := compareDONs(currCCIPDONs, newCCIPDONs)
	if err != nil {
		return diffResult{}, fmt.Errorf("failed to compare CCIP DONs: %w", err)
	}

	return diffRes, nil
}

// compareDONs compares the current and new CCIP DONs and returns the added, removed and updated DONs.
func compareDONs(
	currCCIPDONs,
	newCCIPDONs map[registrysyncer.DonID]kcr.CapabilitiesRegistryDONInfo,
) (
	dr diffResult,
	err error,
) {
	added := make(map[registrysyncer.DonID]kcr.CapabilitiesRegistryDONInfo)
	removed := make(map[registrysyncer.DonID]kcr.CapabilitiesRegistryDONInfo)
	updated := make(map[registrysyncer.DonID]kcr.CapabilitiesRegistryDONInfo)

	for id, don := range newCCIPDONs {
		if currDONState, ok := currCCIPDONs[id]; !ok {
			// Not in current state, so mark as added.
			added[id] = don
		} else {
			// If its in the current state and the config count for the DON has changed, mark as updated.
			// Since the registry returns the full state we need to compare the config count.
			if don.ConfigCount > currDONState.ConfigCount {
				updated[id] = don
			}
		}
	}

	for id, don := range currCCIPDONs {
		if _, ok := newCCIPDONs[id]; !ok {
			// In current state but not in latest registry state, so should remove.
			removed[id] = don
		}
	}

	return diffResult{
		added:   added,
		removed: removed,
		updated: updated,
	}, nil
}

// filterCCIPDONs filters the CCIP DONs from the given state.
func filterCCIPDONs(
	ccipCapability kcr.CapabilitiesRegistryCapabilityInfo,
	state registrysyncer.State,
) (map[registrysyncer.DonID]kcr.CapabilitiesRegistryDONInfo, error) {
	ccipDONs := make(map[registrysyncer.DonID]kcr.CapabilitiesRegistryDONInfo)
	for _, don := range state.IDsToDONs {
		for _, donCapabilities := range don.CapabilityConfigurations {
			hid, err := hashedCapabilityId(ccipCapability.LabelledName, ccipCapability.Version)
			if err != nil {
				return nil, fmt.Errorf("failed to hash capability id: %w", err)
			}
			if donCapabilities.CapabilityId == hid {
				ccipDONs[registrysyncer.DonID(don.Id)] = don
			}
		}
	}

	return ccipDONs, nil
}

// checkCapabilityPresence checks if the capability with the given version and
// labelled name is present in the given capability registry state.
func checkCapabilityPresence(
	capabilityVersion,
	capabilityLabelledName string,
	state registrysyncer.State,
) (kcr.CapabilitiesRegistryCapabilityInfo, error) {
	// Sanity check to make sure the capability registry has the capability we are looking for.
	hid, err := hashedCapabilityId(capabilityLabelledName, capabilityVersion)
	if err != nil {
		return kcr.CapabilitiesRegistryCapabilityInfo{}, fmt.Errorf("failed to hash capability id: %w", err)
	}
	ccipCapability, ok := state.IDsToCapabilities[hid]
	if !ok {
		return kcr.CapabilitiesRegistryCapabilityInfo{},
			fmt.Errorf("failed to find capability with name %s and version %s in capability registry state",
				capabilityLabelledName, capabilityVersion)
	}

	return ccipCapability, nil
}

// hashedCapabilityId returns the hashed capability id in a manner equivalent to the capability registry.
func hashedCapabilityId(capabilityLabelledName, capabilityVersion string) (r [32]byte, err error) {
	tabi := `[{"type": "string"}, {"type": "string"}]`
	abiEncoded, err := utils.ABIEncode(tabi, capabilityLabelledName, capabilityVersion)
	if err != nil {
		return r, fmt.Errorf("failed to ABI encode capability version and labelled name: %w", err)
	}

	h := crypto.Keccak256(abiEncoded)
	copy(r[:], h)
	return r, nil
}

// mustHashedCapabilityId is a helper function that panics if the hashedCapabilityId function returns an error.
// should only use in tests.
func mustHashedCapabilityId(capabilityLabelledName, capabilityVersion string) [32]byte {
	r, err := hashedCapabilityId(capabilityLabelledName, capabilityVersion)
	if err != nil {
		panic(err)
	}
	return r
}

// isMemberOfDON returns true if and only if the given p2pID is a member of the given DON.
func isMemberOfDON(don kcr.CapabilitiesRegistryDONInfo, p2pID ragep2ptypes.PeerID) bool {
	for _, node := range don.NodeP2PIds {
		if node == p2pID {
			return true
		}
	}
	return false
}

// isMemberOfBootstrapSubcommittee returns true if and only if the given p2pID is a member of the given bootstrap subcommittee.
func isMemberOfBootstrapSubcommittee(
	bootstrapP2PIDs [][32]byte,
	p2pID ragep2ptypes.PeerID,
) bool {
	for _, bootstrapID := range bootstrapP2PIDs {
		if bootstrapID == p2pID {
			return true
		}
	}
	return false
}
