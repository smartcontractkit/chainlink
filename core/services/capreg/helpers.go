package capreg

import "bytes"

func capabilityDiff(oldCCs, newCCs []CapabilityConfiguration) (
	removedCapabilities,
	newOrUpdatedCapabilities []CapabilityID,
) {
	// check for removed capabilities.
	for _, oldCC := range oldCCs {
		var found bool
		for _, newCC := range newCCs {
			if oldCC.CapabilityID == newCC.CapabilityID {
				found = true
				break
			}
		}
		if !found {
			removedCapabilities = append(removedCapabilities, oldCC.CapabilityID)
		}
	}

	// check for new or updated capabilities.
	for _, newCC := range newCCs {
		var found bool
		for _, oldCC := range oldCCs {
			if newCC.CapabilityID == oldCC.CapabilityID {
				found = true
				if !oldCC.Equal(newCC) {
					newOrUpdatedCapabilities = append(newOrUpdatedCapabilities, newCC.CapabilityID)
				}
				break
			}
		}
		if !found {
			newOrUpdatedCapabilities = append(newOrUpdatedCapabilities, newCC.CapabilityID)
		}
	}

	return
}

func nodesChanged(oldNodes, newNodes [][]byte) bool {
	if len(oldNodes) != len(newNodes) {
		return true
	}
	for i := range oldNodes {
		if !bytes.Equal(oldNodes[i], newNodes[i]) {
			return true
		}
	}
	return false
}

func filterRelevantDONs(localP2PID []byte, s State) (relevantDONs map[uint32]DON) {
	relevantDONs = make(map[uint32]DON)
	for _, don := range s.DONs {
		var isMember bool
		for _, member := range don.Nodes {
			if bytes.Equal(member, localP2PID) {
				isMember = true
				break
			}
		}
		if isMember {
			relevantDONs[don.ID] = don
		}
	}
	return relevantDONs
}
