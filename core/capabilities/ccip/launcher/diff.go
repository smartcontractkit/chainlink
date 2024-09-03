package launcher

import (
	"fmt"

	ragep2ptypes "github.com/smartcontractkit/libocr/ragep2p/types"

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
	capabilityID string,
	oldState,
	newState registrysyncer.LocalRegistry,
) (diffResult, error) {
	ccipCapability, err := checkCapabilityPresence(capabilityID, newState)
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
			// If its in the current state and the config count for the DON has changed, mark as updated.
			// Since the registry returns the full state we need to compare the config count.
			if don.ConfigVersion > currDONState.ConfigVersion {
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
	ccipCapability registrysyncer.Capability,
	state registrysyncer.LocalRegistry,
) (map[registrysyncer.DonID]registrysyncer.DON, error) {
	ccipDONs := make(map[registrysyncer.DonID]registrysyncer.DON)
	for _, don := range state.IDsToDONs {
		_, ok := don.CapabilityConfigurations[ccipCapability.ID]
		if ok {
			ccipDONs[registrysyncer.DonID(don.ID)] = don
		}
	}

	return ccipDONs, nil
}

// checkCapabilityPresence checks if the capability with the given capabilityID
// is present in the given capability registry state.
func checkCapabilityPresence(
	capabilityID string,
	state registrysyncer.LocalRegistry,
) (registrysyncer.Capability, error) {
	// Sanity check to make sure the capability registry has the capability we are looking for.
	ccipCapability, ok := state.IDsToCapabilities[capabilityID]
	if !ok {
		return registrysyncer.Capability{},
			fmt.Errorf("failed to find capability with capabilityID %s in capability registry state", capabilityID)
	}

	return ccipCapability, nil
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
