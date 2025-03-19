package infra

import (
	"fmt"

	cretypes "github.com/smartcontractkit/chainlink/system-tests/lib/cre/types"
	"github.com/smartcontractkit/chainlink/system-tests/lib/types"
)

func Host(nodeIndex int, nodeType cretypes.CapabilityFlag, donName string, infraDetails types.InfraInput) string {
	if infraDetails.InfraType == types.InfraType_CRIB {
		if nodeType == cretypes.BootstrapNode {
			return fmt.Sprintf("%s-bt-%d", donName, nodeIndex)
		}
		return fmt.Sprintf("%s-%d", donName, nodeIndex)
	}

	return fmt.Sprintf("%s-node%d", donName, nodeIndex)
}
