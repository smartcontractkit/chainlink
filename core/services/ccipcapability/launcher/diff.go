package launcher

import (
	"fmt"

	gethCommon "github.com/ethereum/go-ethereum/common"
	ragep2ptypes "github.com/smartcontractkit/libocr/ragep2p/types"

	kcr "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/capabilities_registry"
	"github.com/smartcontractkit/chainlink/v2/core/services/ccipcapability/common"
	"github.com/smartcontractkit/chainlink/v2/core/services/registrysyncer"
)

// diffResult contains the added, removed and updated CCIP DONs.
// It is determined by using the `diff` function below.
type diffResult struct {
	added   map[registrysyncer.DonID]registrysyncer.DON
	removed map[registrysyncer.DonID]registrysyncer.DON
	updated map[registrysyncer.DonID]registrysyncer.DON
}

// diff compares the old and new state and returns the added, removed and updated CCIP DONs.
func diff(
	capabilityVersion,
	capabilityLabelledName string,
	oldState,
	newState registrysyncer.LocalRegistry,
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
	newCCIPDONs map[registrysyncer.DonID]registrysyncer.DON,
) (
	dr diffResult,
	err error,
) {
	added := make(map[registrysyncer.DonID]registrysyncer.DON)
	removed := make(map[registrysyncer.DonID]registrysyncer.DON)
	updated := make(map[registrysyncer.DonID]registrysyncer.DON)

	for id, don := range newCCIPDONs {
		if currDONState, ok := currCCIPDONs[id]; !ok {
			// Not in current state, so mark as added.
			added[id] = don
		} else {
			fmt.Println(currDONState)
			// If its in the current state and the config count for the DON has changed, mark as updated.
			// Since the registry returns the full state we need to compare the config count.
			//if don.ConfigCount > currDONState.ConfigCount {
			//	updated[id] = don
			//}
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
	state registrysyncer.LocalRegistry,
) (map[registrysyncer.DonID]registrysyncer.DON, error) {
	ccipDONs := make(map[registrysyncer.DonID]registrysyncer.DON)
	for _, don := range state.IDsToDONs {
		for cid, _ := range don.CapabilityConfigurations {
			capability, ok := state.IDsToCapabilities[cid]
			if !ok {
				return nil, fmt.Errorf("could not find capability matching id %s", cid)
			}
			hid, err := common.HashedCapabilityID(ccipCapability.LabelledName, ccipCapability.Version)
			hidStr := gethCommon.Bytes2Hex(hid[:])
			if err != nil {
				return nil, fmt.Errorf("failed to hash capability id: %w", err)
			}
			if capability.ID == hidStr {
				ccipDONs[registrysyncer.DonID(don.ID)] = don
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
	state registrysyncer.LocalRegistry,
) (kcr.CapabilitiesRegistryCapabilityInfo, error) {
	// Sanity check to make sure the capability registry has the capability we are looking for.
	//hid, err := common.HashedCapabilityID(capabilityLabelledName, capabilityVersion)
	//if err != nil {
	//	return kcr.CapabilitiesRegistryCapabilityInfo{}, fmt.Errorf("failed to hash capability id: %w", err)
	//}
	//hidStr := gethCommon.Bytes2Hex(hid[:])
	//ccipCapability, ok := state.IDsToCapabilities[hid]
	//if !ok {
	//	return kcr.CapabilitiesRegistryCapabilityInfo{},
	//		fmt.Errorf("failed to find capability with name %s and version %s in capability registry state",
	//			capabilityLabelledName, capabilityVersion)
	//}

	//return ccipCapability, nil
	return kcr.CapabilitiesRegistryCapabilityInfo{}, nil
}

// isMemberOfDON returns true if and only if the given p2pID is a member of the given DON.
func isMemberOfDON(don registrysyncer.DON, p2pID ragep2ptypes.PeerID) bool {
	for _, node := range don.Members {
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
